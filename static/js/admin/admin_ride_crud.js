/**
 * @file admin_ride_crud.js
 * @description Gestión de despacho con herencia automática de reservas.
 */

document.addEventListener('DOMContentLoaded', async () => {
    // 1. Cargar datos maestros
    await Promise.all([
        loadBookings(),
        loadDrivers(),
        loadVehicles()
    ]);

    // 2. Modo Edición
    const urlParams = new URLSearchParams(window.location.search);
    const rideId = urlParams.get('id');
    if (rideId) {
        await setupEditMode(rideId);
    }

    document.getElementById('crud-form').addEventListener('submit', handleFormSubmit);
});

async function loadBookings() {
    const res = await fetch('/api/v1/bookings');
    const data = await res.json();
    const select = document.getElementById('booking_id');
    select.innerHTML = '<option value="">-- Seleccionar Booking --</option>';
    data.forEach(b => {
        // Solo mostramos reservas que no estén canceladas o ya completadas
        if (b.status !== 'cancelled' && b.status !== 'completed') {
            const opt = document.createElement('option');
            opt.value = b.id || b.ID;
            opt.textContent = `ID:${b.id || b.ID} | ${b.client_name} | ${b.scheduled_time}`;
            select.appendChild(opt);
        }
    });
}

/**
 * FUNCIÓN CLAVE: Rellena el formulario al seleccionar una reserva
 */
async function autoFillFromBooking(bookingId) {
    if (!bookingId) return;
    try {
        const res = await fetch(`/api/v1/bookings/${bookingId}`);
        const b = await res.json();

        // Rellenar campos heredados
        document.getElementById('client_name').value = b.client_name || '';
        document.getElementById('client_phone').value = b.client_phone || '';
        document.getElementById('pax').value = b.pax || 4;
        document.getElementById('maletas').value = b.maletas || 3;
        document.getElementById('animal').checked = b.animal || false;
        document.getElementById('origin_address').value = b.origin_address || '';
        document.getElementById('destination_address').value = b.dest_address || '';
        
        console.log("Datos de reserva #"+bookingId+" heredados con éxito.");
    } catch (e) { console.error("Error al heredar datos:", e); }
}

async function loadDrivers() {
    const res = await fetch('/api/v1/drivers');
    const data = await res.json();
    const select = document.getElementById('driver_id');
    select.innerHTML = '<option value="">-- Seleccionar Conductor --</option>';
    data.forEach(d => {
        const opt = document.createElement('option');
        opt.value = d.id || d.ID;
        opt.textContent = d.full_name.toUpperCase();
        select.appendChild(opt);
    });
}

async function loadVehicles() {
    const res = await fetch('/api/v1/vehicles');
    const data = await res.json();
    const select = document.getElementById('vehicle_id');
    select.innerHTML = '<option value="">-- Vehículo --</option>';
    data.forEach(v => {
        const opt = document.createElement('option');
        opt.value = v.id || v.ID;
        opt.textContent = `${v.plate_number} (${v.make} ${v.model})`;
        select.appendChild(opt);
    });
}

async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/rides/${id}`);
        const r = await res.json();

        document.getElementById('rideId').value = r.id || r.ID;
        document.getElementById('booking_id').value = r.booking_id;
        document.getElementById('driver_id').value = r.driver_id;
        document.getElementById('vehicle_id').value = r.vehicle_id;
        document.getElementById('voucher_number').value = r.voucher_number || '';
        
        document.getElementById('client_name').value = r.client_name || '';
        document.getElementById('client_phone').value = r.client_phone || '';
        document.getElementById('pax').value = r.pax || 4;
        document.getElementById('maletas').value = r.maletas || 3;
        document.getElementById('animal').checked = r.animal || false;

        document.getElementById('origin_address').value = r.origin_address || '';
        document.getElementById('destination_address').value = r.destination_address || '';
        
        document.getElementById('km_start').value = r.km_start || 0;
        document.getElementById('total_amount').value = r.total_amount || 0;
        document.getElementById('is_finished').checked = r.is_finished;

        document.getElementById('page-title').textContent = "Modificar Servicio en Curso";
        document.getElementById('btn-save').textContent = "SINCRONIZAR CAMBIOS OPERATIVOS";
    } catch (e) { console.error(e); }
}

async function handleFormSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('rideId').value;

    const payload = {
        booking_id: parseInt(document.getElementById('booking_id').value),
        driver_id: parseInt(document.getElementById('driver_id').value),
        vehicle_id: parseInt(document.getElementById('vehicle_id').value),
        voucher_number: document.getElementById('voucher_number').value.trim(),
        
        client_name: document.getElementById('client_name').value,
        client_phone: document.getElementById('client_phone').value,
        pax: parseInt(document.getElementById('pax').value),
        maletas: parseInt(document.getElementById('maletas').value),
        animal: document.getElementById('animal').checked,

        origin_address: document.getElementById('origin_address').value,
        destination_address: document.getElementById('destination_address').value,
        
        km_start: parseFloat(document.getElementById('km_start').value) || 0,
        km_end: parseFloat(document.getElementById('km_end').value) || 0,
        total_amount: parseFloat(document.getElementById('total_amount').value) || 0,
        is_finished: document.getElementById('is_finished').checked
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `/api/v1/rides/${id}` : '/api/v1/rides';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert("Carrera despachada correctamente");
            window.location.href = "/admin/rides";
        } else {
            alert("Fallo en la sincronización.");
        }
    } catch (e) { alert("Error de servidor"); }
}