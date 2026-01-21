/**
 * @file admin_ride_crud.js
 * @description Gestión de carreras vinculando Bookings, Drivers y Vehicles con GPS y métricas.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[RIDE-CRUD] Sincronizando modelo...", "color: #3b82f6; font-weight: bold;");

    await Promise.all([
        loadBookings(),
        loadDrivers(),
        loadVehicles()
    ]);

    const urlParams = new URLSearchParams(window.location.search);
    const rideId = urlParams.get('id');

    if (rideId) {
        await setupEditMode(rideId);
    }

    document.getElementById('crud-form').addEventListener('submit', handleFormSubmit);
});

// Carga de selectores (Igual que antes)
async function loadBookings() {
    const res = await fetch('/api/v1/bookings');
    const data = await res.json();
    const select = document.getElementById('booking_id');
    select.innerHTML = '<option value="">-- Seleccionar Booking --</option>';
    data.forEach(b => {
        const opt = document.createElement('option');
        opt.value = b.id || b.ID;
        opt.textContent = `ID:${b.id || b.ID} - ${b.pickup_location.substring(0,20)}...`;
        select.appendChild(opt);
    });
}

async function loadDrivers() {
    const res = await fetch('/api/v1/drivers');
    const data = await res.json();
    const select = document.getElementById('driver_id');
    select.innerHTML = '<option value="">-- Conductor --</option>';
    data.forEach(d => {
        const opt = document.createElement('option');
        opt.value = d.id || d.ID;
        opt.textContent = d.full_name;
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
        opt.textContent = `${v.plate_number} (${v.make})`;
        select.appendChild(opt);
    });
}

async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/rides/${id}`);
        const r = await res.json();

        // 1. IDs
        document.getElementById('rideId').value = r.id || r.ID;
        document.getElementById('booking_id').value = r.booking_id;
        document.getElementById('driver_id').value = r.driver_id;
        document.getElementById('vehicle_id').value = r.vehicle_id;
        document.getElementById('voucher_number').value = r.voucher_number || '';

        // 2. Tiempos (Tratamiento de fechas para input datetime-local)
        const formatDate = (dateStr) => dateStr ? dateStr.slice(0, 16) : '';
        document.getElementById('start_time_real').value = formatDate(r.start_time_real);
        document.getElementById('end_time_real').value = formatDate(r.end_time_real);
        document.getElementById('wait_time_minutes').value = r.wait_time_minutes || 0;

        // 3. Kilometraje y Dinero
        document.getElementById('km_start').value = r.km_start || 0;
        document.getElementById('km_end').value = r.km_end || 0;
        document.getElementById('km_total').value = r.km_total || 0;
        document.getElementById('is_holiday').checked = r.is_holiday;
        document.getElementById('is_night_shift').checked = r.is_night_shift;
        document.getElementById('extra_charges').value = r.extra_charges || 0;
        document.getElementById('total_amount').value = r.total_amount || 0;
        document.getElementById('is_finished').checked = r.is_finished;

        // 4. GPS
        document.getElementById('origin_address').value = r.origin_address || '';
        document.getElementById('origin_lat').value = r.origin_lat || '';
        document.getElementById('origin_lng').value = r.origin_lng || '';
        document.getElementById('destination_address').value = r.destination_address || '';
        document.getElementById('destination_lat').value = r.destination_lat || '';
        document.getElementById('destination_lng').value = r.destination_lng || '';

        document.getElementById('page-title').textContent = "Modificar Carrera Realizada";

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

        // Fechas
        start_time_real: document.getElementById('start_time_real').value ? new Date(document.getElementById('start_time_real').value).toISOString() : null,
        end_time_real: document.getElementById('end_time_real').value ? new Date(document.getElementById('end_time_real').value).toISOString() : null,
        wait_time_minutes: parseInt(document.getElementById('wait_time_minutes').value) || 0,

        // Kms
        km_start: parseFloat(document.getElementById('km_start').value) || 0,
        km_end: parseFloat(document.getElementById('km_end').value) || 0,
        km_total: parseFloat(document.getElementById('km_total').value) || 0,

        // Flags y Dinero
        is_holiday: document.getElementById('is_holiday').checked,
        is_night_shift: document.getElementById('is_night_shift').checked,
        extra_charges: parseFloat(document.getElementById('extra_charges').value) || 0,
        total_amount: parseFloat(document.getElementById('total_amount').value) || 0,
        is_finished: document.getElementById('is_finished').checked,

        // GPS
        origin_address: document.getElementById('origin_address').value.trim(),
        origin_lat: parseFloat(document.getElementById('origin_lat').value) || 0,
        origin_lng: parseFloat(document.getElementById('origin_lng').value) || 0,
        destination_address: document.getElementById('destination_address').value.trim(),
        destination_lat: parseFloat(document.getElementById('destination_lat').value) || 0,
        destination_lng: parseFloat(document.getElementById('destination_lng').value) || 0
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
            alert("Carrera guardada con éxito");
            window.location.href = "/admin/rides";
        } else {
            const err = await res.json();
            alert("Error: " + (err.error || "Fallo al guardar"));
        }
    } catch (e) { alert("Error de servidor"); }
}