/**
 * @file admin_booking_crud.js
 * @description Gestión de creación y edición de reservas TakeUs con auto-creación de cliente
 */

document.addEventListener('DOMContentLoaded', async () => {
    await loadClientsForSelect();

    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get('id');

    if (id) {
        await setupEditMode(id);
    }

    document.getElementById('crud-form').addEventListener('submit', handleSave);
});

async function loadClientsForSelect() {
    try {
        const res = await fetch('/api/v1/clients');
        const clients = await res.json();
        const select = document.getElementById('client_id');
        select.innerHTML = '<option value="">-- CREAR NUEVO CLIENTE --</option>';
        clients.forEach(c => {
            const opt = document.createElement('option');
            opt.value = c.id || c.ID;
            opt.textContent = c.full_name.toUpperCase();
            select.appendChild(opt);
        });
    } catch (e) { console.error("Error cargando clientes:", e); }
}

async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/bookings/${id}`);
        if (!res.ok) throw new Error();
        const b = await res.json();

        document.getElementById('bookingId').value = b.id || b.ID;
        document.getElementById('client_id').value = b.client_id;
        document.getElementById('client_name').value = b.client_name;
        document.getElementById('client_phone').value = b.client_phone;
        document.getElementById('pax').value = b.pax;
        document.getElementById('maletas').value = b.maletas;
        document.getElementById('animal').checked = b.animal;
        document.getElementById('scheduled_date').value = b.scheduled_date.split('T')[0];
        document.getElementById('scheduled_time').value = b.scheduled_time;
        document.getElementById('origin_address').value = b.origin_address;
        document.getElementById('dest_address').value = b.dest_address;
        document.getElementById('status').value = b.status;
        document.getElementById('notes').value = b.notes || "";

        document.getElementById('page-title').textContent = "Modificar Expediente Reserva";
        document.getElementById('btn-save').textContent = "ACTUALIZAR DATOS RESERVA";
    } catch (e) { console.error("Error al cargar modo edición", e); }
}

async function handleSave(e) {
    e.preventDefault();
    
    const bookingId = document.getElementById('bookingId').value;
    const date = document.getElementById('scheduled_date').value;
    const time = document.getElementById('scheduled_time').value;
    let clientId = document.getElementById('client_id').value;

    const clientName = document.getElementById('client_name').value.trim();
    const clientPhone = document.getElementById('client_phone').value.trim();

    // --- LÓGICA DE AUTO-CREACIÓN DE CLIENTE ---
    if (!clientId && !bookingId) { // Solo creamos cliente nuevo si es una reserva nueva
        console.log("Detectado cliente nuevo, registrando en base de datos...");
        try {
            const clientRes = await fetch('/api/v1/clients', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    full_name: clientName,
                    phone: clientPhone,
                    email: "modificar@email.es",
                    preferences: "Creado desde despacho rápido"
                })
            });

            if (clientRes.ok) {
                const newClient = await clientRes.json();
                clientId = newClient.id || newClient.ID;
                console.log("Cliente creado con ID:", clientId);
            } else {
                alert("No se pudo crear el cliente nuevo automáticamente.");
                return;
            }
        } catch (err) {
            console.error("Error creando cliente:", err);
            return;
        }
    }

    const payload = {
        client_id: parseInt(clientId),
        client_name: clientName,
        client_phone: clientPhone,
        pax: parseInt(document.getElementById('pax').value),
        maletas: parseInt(document.getElementById('maletas').value),
        animal: document.getElementById('animal').checked,
        scheduled_date: new Date(date).toISOString(),
        scheduled_time: time,
        scheduled_at: new Date(`${date}T${time}:00Z`).toISOString(),
        origin_address: document.getElementById('origin_address').value.trim(),
        origin_lat: 0.0,
        origin_lng: 0.0,
        dest_address: document.getElementById('dest_address').value.trim(),
        dest_lat: 0.0,
        dest_lng: 0.0,
        status: document.getElementById('status').value,
        notes: document.getElementById('notes').value.trim()
    };

    const method = bookingId ? 'PUT' : 'POST';
    const url = bookingId ? `/api/v1/bookings/${bookingId}` : '/api/v1/bookings';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert("Reserva y Cliente sincronizados con éxito");
            window.location.href = "/admin/bookings";
        } else {
            const err = await res.json();
            alert("Error del sistema: " + (err.error || "Fallo al guardar"));
        }
    } catch (e) { 
        alert("Error crítico de comunicación"); 
    }
}