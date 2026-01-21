/**
 * @file admin_company_crud.js
 * @description Lógica unificada para Crear (POST) y Editar (PUT) empresas.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[INIT] Formulario CRUD Empresa", "color: #3b82f6; font-weight: bold;");
    
    // 1. Cargar usuarios primero (necesario tanto para crear como para editar)
    await loadUsersForSelection();

    // 2. Comprobar si estamos en modo EDICIÓN (?id=X)
    const urlParams = new URLSearchParams(window.location.search);
    const companyId = urlParams.get('id');

    if (companyId) {
        console.log(`[MODE] Editando Empresa ID: ${companyId}`);
        await setupEditMode(companyId);
    } else {
        console.log("[MODE] Creando Nueva Empresa");
    }

    // 3. Vincular evento submit
    const form = document.getElementById('crud-form');
    if (form) {
        form.addEventListener('submit', handleFormSubmit);
    }
});

/**
 * Obtiene usuarios con rol 'company' para el select
 */
async function loadUsersForSelection() {
    try {
        const res = await fetch('/api/v1/users');
        const users = await res.json();
        const select = document.getElementById('user_id');
        
        const companyUsers = users.filter(u => u.role === 'company');
        
        select.innerHTML = '<option value="">-- Seleccionar Usuario Responsable --</option>';
        companyUsers.forEach(u => {
            const opt = document.createElement('option');
            opt.value = u.id || u.ID;
            opt.textContent = `${u.email} (ID: ${u.id || u.ID})`;
            select.appendChild(opt);
        });
    } catch (e) {
        console.error("Error cargando usuarios:", e);
    }
}

/**
 * Carga los datos de la empresa para rellenar el formulario
 */
async function setupEditMode(id) {
    try {
        const res = await fetch(`/api/v1/companies/${id}`);
        if (!res.ok) throw new Error("Empresa no encontrada");
        
        const c = await res.json();

        // Rellenar campos del HTML
        document.getElementById('companyId').value = c.id || c.ID;
        document.getElementById('user_id').value = c.user_id;
        document.getElementById('name').value = c.name;
        document.getElementById('tax_id').value = c.tax_id;
        document.getElementById('address').value = c.address;
        document.getElementById('postal_code').value = c.postal_code;
        document.getElementById('website').value = c.website;

        // Cambiar estética a modo Edición
        document.getElementById('page-title').textContent = "Modificar Empresa";
        const btn = document.getElementById('btn-save');
        btn.textContent = "ACTUALIZAR DATOS";
        btn.classList.replace('bg-admin-dark', 'bg-admin-accent');

    } catch (e) {
        console.error(e);
        alert("Error al cargar datos de la empresa.");
    }
}

/**
 * Envía los datos (POST si no hay ID, PUT si hay ID)
 */
async function handleFormSubmit(e) {
    e.preventDefault();

    const id = document.getElementById('companyId').value;
    const userIdValue = document.getElementById('user_id').value;

    if (!userIdValue) {
        alert("Seleccione un usuario");
        return;
    }

    const payload = {
        user_id: parseInt(userIdValue),
        name: document.getElementById('name').value.trim(),
        tax_id: document.getElementById('tax_id').value.trim(),
        address: document.getElementById('address').value.trim(),
        postal_code: document.getElementById('postal_code').value.trim(),
        website: document.getElementById('website').value.trim()
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `/api/v1/companies/${id}` : '/api/v1/companies';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert("Operación exitosa");
            window.location.href = "/admin/companies";
        } else {
            const err = await res.json();
            alert("Error: " + (err.error || "Fallo en el servidor"));
        }
    } catch (e) {
        alert("Error de conexión con la API");
    }
}