/**
 * @file admin_driver_crud.js
 * @description Lógica de gestión para el alta y modificación de conductores.
 * @dependencies Requiere endpoints de /api/v1/users, /api/v1/companies y /api/v1/drivers.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[DRIVER-CRUD] Inicializando sistema de gestión de flota", "color: #3b82f6; font-weight: bold;");

    // 1. Carga de dependencias necesarias (Selects)
    await Promise.all([
        loadUsersForDriverSelection(),
        loadCompaniesForSelection()
    ]);

    // 2. Determinar modo (Crear o Editar)
    const urlParams = new URLSearchParams(window.location.search);
    const driverId = urlParams.get('id');

    if (driverId) {
        await setupEditMode(driverId);
    }

    // 3. Listener del formulario
    document.getElementById('crud-form').addEventListener('submit', handleFormSubmit);
});

/**
 * Carga usuarios con rol 'driver' que no están vinculados (opcionalmente filtrados en backend)
 */
async function loadUsersForDriverSelection() {
    try {
        const res = await fetch('/api/v1/users');
        const users = await res.json();
        const select = document.getElementById('user_id');
        
        // Filtramos solo conductores
        const drivers = users.filter(u => u.role === 'driver');
        
        select.innerHTML = '<option value="">-- Vincular cuenta de acceso --</option>';
        drivers.forEach(u => {
            const opt = document.createElement('option');
            opt.value = u.id || u.ID;
            opt.textContent = `${u.email} (ID: ${u.id || u.ID})`;
            select.appendChild(opt);
        });
    } catch (e) { console.error("Error cargando usuarios driver:", e); }
}

/**
 * Carga todas las empresas disponibles para asignar al conductor
 */
async function loadCompaniesForSelection() {
    try {
        const res = await fetch('/api/v1/companies');
        const companies = await res.json();
        const select = document.getElementById('company_id');
        
        select.innerHTML = '<option value="">-- Seleccionar Empresa --</option>';
        companies.forEach(c => {
            const opt = document.createElement('option');
            opt.value = c.id || c.ID;
            opt.textContent = c.name.toUpperCase();
            select.appendChild(opt);
        });
    } catch (e) { console.error("Error cargando empresas:", e); }
}

/**
 * Prepara el formulario con datos existentes del conductor
 */
async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/drivers/${id}`);
        const d = await res.json();

        document.getElementById('driverId').value = d.id || d.ID;
        document.getElementById('user_id').value = d.user_id;
        document.getElementById('company_id').value = d.company_id;
        document.getElementById('full_name').value = d.full_name;
        document.getElementById('email').value = d.email;
        document.getElementById('phone').value = d.phone;
        document.getElementById('license_number').value = d.license_number;
        
        // Formatear fecha para input date (YYYY-MM-DD)
        if (d.license_expiry_date) {
            document.getElementById('license_expiry_date').value = d.license_expiry_date.split('T')[0];
        }

        document.getElementById('is_available').checked = d.is_available;

        // UI Update
        document.getElementById('page-title').textContent = "Modificar Conductor";
        const btn = document.getElementById('btn-save');
        btn.textContent = "GUARDAR CAMBIOS EN EXPEDIENTE";
        btn.classList.replace('bg-admin-dark', 'bg-admin-accent');
    } catch (e) { console.error("Error en modo edición:", e); }
}

/**
 * Gestiona el envío de datos a la API (POST o PUT)
 */
async function handleFormSubmit(e) {
    e.preventDefault();

    const id = document.getElementById('driverId').value;
    
    // Preparar el objeto con los tipos de datos correctos para Go
    const payload = {
        user_id: parseInt(document.getElementById('user_id').value),
        company_id: parseInt(document.getElementById('company_id').value),
        full_name: document.getElementById('full_name').value.trim(),
        email: document.getElementById('email').value.trim(),
        phone: document.getElementById('phone').value.trim(),
        license_number: document.getElementById('license_number').value.trim(),
        license_expiry_date: new Date(document.getElementById('license_expiry_date').value).toISOString(),
        is_available: document.getElementById('is_available').checked
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `/api/v1/drivers/${id}` : '/api/v1/drivers';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert("Conductor guardado correctamente");
            window.location.href = "/admin/drivers";
        } else {
            const err = await res.json();
            alert("Error: " + (err.error || "No se pudo procesar la solicitud"));
        }
    } catch (e) {
        alert("Error crítico de comunicación con el servidor");
    }
}