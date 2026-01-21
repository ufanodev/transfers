/**
 * @file admin_vehicle_crud.js
 * @description Lógica de gestión para el alta y modificación de vehículos de la flota.
 * @dependencies Requiere endpoints de /api/v1/companies y /api/v1/vehicles.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[VEHICLE-CRUD] Inicializando sistema de gestión de flota", "color: #3b82f6; font-weight: bold;");

    // 1. Cargar empresas para el selector
    await loadCompaniesForSelection();

    // 2. Determinar modo (Crear o Editar)
    const urlParams = new URLSearchParams(window.location.search);
    const vehicleId = urlParams.get('id');

    if (vehicleId) {
        await setupEditMode(vehicleId);
    }

    // 3. Listener del formulario
    document.getElementById('crud-form').addEventListener('submit', handleFormSubmit);
});

/**
 * Carga todas las empresas disponibles para asignar el vehículo
 */
async function loadCompaniesForSelection() {
    try {
        const res = await fetch('/api/v1/companies');
        const companies = await res.json();
        const select = document.getElementById('company_id');
        
        select.innerHTML = '<option value="">-- Seleccionar Empresa Propietaria --</option>';
        companies.forEach(c => {
            const opt = document.createElement('option');
            opt.value = c.id || c.ID;
            opt.textContent = c.name.toUpperCase();
            select.appendChild(opt);
        });
    } catch (e) {
        console.error("Error cargando empresas:", e);
    }
}

/**
 * Prepara el formulario con datos existentes del vehículo
 */
async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/vehicles/${id}`);
        if (!res.ok) throw new Error("Vehículo no encontrado");
        const v = await res.json();

        // Rellenar campos
        document.getElementById('vehicleId').value = v.id || v.ID;
        document.getElementById('company_id').value = v.company_id;
        document.getElementById('make').value = v.make;
        document.getElementById('model_name').value = v.model_name;
        document.getElementById('plate_number').value = v.plate_number;
        document.getElementById('license_number').value = v.license_number;
        
        if (v.license_expiry_date) {
            document.getElementById('license_expiry_date').value = v.license_expiry_date.split('T')[0];
        }

        document.getElementById('vehicle_type').value = v.vehicle_type;
        document.getElementById('category').value = v.category;
        document.getElementById('capacity_pax').value = v.capacity_pax;
        document.getElementById('base_fare').value = v.base_fare;
        document.getElementById('is_active').checked = v.is_active;

        // Cambiar interfaz
        document.getElementById('page-title').textContent = "Modificar Vehículo";
        const btn = document.getElementById('btn-save');
        btn.textContent = "ACTUALIZAR FICHA TÉCNICA";
        btn.classList.replace('bg-admin-dark', 'bg-admin-accent');

    } catch (e) {
        console.error("Error en modo edición:", e);
        alert("No se pudieron cargar los datos del vehículo.");
    }
}

/**
 * Gestiona el envío de datos (POST o PUT)
 */
async function handleFormSubmit(e) {
    e.preventDefault();

    const id = document.getElementById('vehicleId').value;
    
    const payload = {
        company_id: parseInt(document.getElementById('company_id').value),
        make: document.getElementById('make').value.trim(),
        model_name: document.getElementById('model_name').value.trim(),
        plate_number: document.getElementById('plate_number').value.trim(),
        license_number: document.getElementById('license_number').value.trim(),
        license_expiry_date: new Date(document.getElementById('license_expiry_date').value).toISOString(),
        vehicle_type: document.getElementById('vehicle_type').value,
        category: document.getElementById('category').value.trim(),
        capacity_pax: parseInt(document.getElementById('capacity_pax').value),
        base_fare: parseFloat(document.getElementById('base_fare').value) || 0,
        is_active: document.getElementById('is_active').checked
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `/api/v1/vehicles/${id}` : '/api/v1/vehicles';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert("Vehículo guardado correctamente");
            window.location.href = "/admin/vehicles";
        } else {
            const err = await res.json();
            alert("Error: " + (err.error || "No se pudo procesar la solicitud"));
        }
    } catch (e) {
        alert("Error de conexión con el servidor");
    }
}