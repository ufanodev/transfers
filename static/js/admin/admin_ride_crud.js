/**
 * @file admin_ride_crud.js
 * @description Lógica para la creación y edición de Rides (Servicios de Transporte).
 * Gestiona la vinculación entre reservas, conductores y vehículos.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[RIDE-CRUD] Inicializando gestor de asignaciones", "color: #3b82f6; font-weight: bold;");

    // 1. Cargar dependencias de los selectores en paralelo
    await Promise.all([
        loadBookingsForSelection(),
        loadDriversForSelection(),
        loadVehiclesForSelection()
    ]);

    // 2. Determinar si es Edición o Creación
    const urlParams = new URLSearchParams(window.location.search);
    const rideId = urlParams.get('id');

    if (rideId) {
        console.log(`[MODE] Editando Ride ID: ${rideId}`);
        await setupEditMode(rideId);
    } else {
        console.log("[MODE] Nueva Asignación Logística");
    }

    // 3. Vincular el evento de envío del formulario
    const form = document.getElementById('crud-form');
    if (form) {
        form.addEventListener('submit', handleFormSubmit);
    }
});

/**
 * Carga reservas (Bookings) disponibles para ser asignadas
 */
async function loadBookingsForSelection() {
    try {
        const res = await fetch('/api/v1/bookings');
        const data = await res.json();
        const select = document.getElementById('booking_id');
        
        select.innerHTML = '<option value="">-- Seleccionar Reserva Pendiente --</option>';
        data.forEach(b => {
            const opt = document.createElement('option');
            opt.value = b.id || b.ID;
            // Mostramos ID y una parte de la ubicación para identificarla
            opt.textContent = `RES-${b.id || b.ID} | Destino: ${b.dropoff_location.substring(0, 25)}...`;
            select.appendChild(opt);
        });
    } catch (e) { console.error("Error al cargar reservas:", e); }
}

/**
 * Carga la lista de conductores disponibles
 */
async function loadDriversForSelection() {
    try {
        const res = await fetch('/api/v1/drivers');
        const data = await res.json();
        const select = document.getElementById('driver_id');
        
        select.innerHTML = '<option value="">-- Seleccionar Conductor --</option>';
        data.forEach(d => {
            const opt = document.createElement('option');
            opt.value = d.id || d.ID;
            opt.textContent = `${d.full_name.toUpperCase()} (${d.is_available ? 'DISPONIBLE' : 'OCUPADO'})`;
            select.appendChild(opt);
        });
    } catch (e) { console.error("Error al cargar conductores:", e); }
}

/**
 * Carga la flota de vehículos
 */
async function loadVehiclesForSelection() {
    try {
        const res = await fetch('/api/v1/vehicles');
        const data = await res.json();
        const select = document.getElementById('vehicle_id');
        
        select.innerHTML = '<option value="">-- Seleccionar Vehículo --</option>';
        data.forEach(v => {
            const opt = document.createElement('option');
            opt.value = v.id || v.ID;
            opt.textContent = `${v.plate_number} - ${v.make} ${v.model_name} (${v.vehicle_type})`;
            select.appendChild(opt);
        });
    } catch (e) { console.error("Error al cargar vehículos:", e); }
}

/**
 * Carga los datos de un Ride existente para modificarlo
 */
async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/rides/${id}`);
        if (!res.ok) throw new Error("Ride no encontrado");
        const r = await res.json();

        // Rellenar campos ocultos y visibles
        document.getElementById('rideId').value = r.id || r.ID;
        document.getElementById('booking_id').value = r.booking_id;
        document.getElementById('driver_id').value = r.driver_id;
        document.getElementById('vehicle_id').value = r.vehicle_id;
        document.getElementById('status').value = r.status;
        document.getElementById('external_id').value = r.external_id || '';

        // Actualizar UI
        document.getElementById('page-title').textContent = "Gestionar Asignación Logística";
        document.getElementById('btn-text').textContent = "ACTUALIZAR ESTADO DEL VIAJE";
        const btn = document.getElementById('btn-save');
        btn.classList.replace('bg-admin-dark', 'bg-admin-accent');

    } catch (e) {
        console.error("Error cargando el Ride:", e);
        alert("No se pudo cargar la información del servicio seleccionado.");
    }
}

/**
 * Procesa el envío del formulario (POST / PUT)
 */
async function handleFormSubmit(e) {
    e.preventDefault();

    const id = document.getElementById('rideId').value;
    
    // Construcción del objeto para el backend en Go
    const payload = {
        booking_id: parseInt(document.getElementById('booking_id').value),
        driver_id: parseInt(document.getElementById('driver_id').value),
        vehicle_id: parseInt(document.getElementById('vehicle_id').value),
        status: document.getElementById('status').value,
        external_id: document.getElementById('external_id').value.trim()
    };

    // Validaciones básicas
    if (!payload.booking_id || !payload.driver_id || !payload.vehicle_id) {
        alert("Por favor, complete la asignación de Conductor y Vehículo.");
        return;
    }

    const method = id ? 'PUT' : 'POST';
    const url = id ? `/api/v1/rides/${id}` : '/api/v1/rides';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert("Operación logística registrada con éxito");
            window.location.href = "/admin/rides";
        } else {
            const err = await res.json();
            alert("Error del servidor: " + (err.error || "Fallo desconocido"));
        }
    } catch (e) {
        console.error(e);
        alert("Error de conexión al intentar guardar la asignación.");
    }
}