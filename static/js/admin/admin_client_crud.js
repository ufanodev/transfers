/**
 * @file admin_client_crud.js
 * @description Gestión de Clientes con logs de auditoría para detectar IDs nulos.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[DEBUG] 🚀 Iniciando script CRUD Clientes...", "color: #3b82f6; font-weight: bold;");
    
    const urlParams = new URLSearchParams(window.location.search);
    const clientId = urlParams.get('id');

    const form = document.getElementById('crud-form');
    if (!form) {
        console.error("❌ ERROR: No se encontró el formulario con ID 'crud-form'");
        return;
    }

    if (clientId) {
        console.log(`[MODE] 📝 Modo Edición para ID: ${clientId}`);
        await setupEditMode(clientId);
    } else {
        console.log("[MODE] ✨ Modo Creación activo.");
        document.getElementById('page-title').textContent = "Nuevo Cliente";
        document.getElementById('btn-save').textContent = "Crear Ficha Cliente";
    }

    form.addEventListener('submit', handleFormSubmit);
});

async function setupEditMode(id) {
    try {
        const response = await fetch(`/api/v1/clients/${id}`);
        if (!response.ok) throw new Error("Cliente no encontrado.");
        const c = await response.json();

        // Función segura para rellenar campos
        const fill = (id, val) => {
            const el = document.getElementById(id);
            if (el) el.value = val || "";
            else console.warn(`[WARN] No se pudo rellenar '${id}' porque no existe en el HTML.`);
        };

        fill('clientId', c.id || c.ID);
        fill('full_name', c.full_name);
        fill('email', c.email);
        fill('phone', c.phone);
        
        if (document.getElementById('preferences')) {
            document.getElementById('preferences').value = typeof c.preferences === 'object' 
                ? JSON.stringify(c.preferences, null, 2) 
                : c.preferences;
        }

        document.getElementById('page-title').textContent = "Modificar Cliente";
        document.getElementById('btn-save').textContent = "Guardar Cambios";
        
        if (document.getElementById('btn-delete-trigger')) {
            document.getElementById('btn-delete-trigger').classList.remove('hidden');
        }

    } catch (err) {
        console.error("❌ Error en carga:", err);
    }
}

async function handleFormSubmit(e) {
    e.preventDefault();
    console.log("%c[SYSTEM] 🛡️ Validando elementos antes de enviar...", "color: #10b981; font-weight: bold;");

    // LISTA DE IDs QUE EL JS BUSCA EN EL HTML
    const requiredIds = ['clientId', 'full_name', 'email', 'phone', 'preferences'];
    const data = {};

    // AUDITORÍA DE ELEMENTOS
    for (const id of requiredIds) {
        const el = document.getElementById(id);
        if (!el) {
            console.error(`%c[CRÍTICO] ❌ El elemento HTML con id="${id}" NO EXISTE.`, "background: red; color: white; padding: 4px;");
            alert(`Error técnico: Falta el campo con ID '${id}' en el HTML.`);
            return; // Aquí es donde se evitaba el error de "null reading value"
        }
        data[id] = el.value.trim();
    }

    const payload = {
        full_name: data.full_name,
        email: data.email,
        phone: data.phone,
        preferences: data.preferences || "{}"
    };

    const method = data.clientId ? 'PUT' : 'POST';
    const url = data.clientId ? `/api/v1/clients/${data.clientId}` : '/api/v1/clients';

    console.log(`[API-SEND] 📤 Enviando datos vía ${method}...`);

    try {
        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            console.log("[SUCCESS] ✅ Cliente guardado correctamente.");
            window.location.href = "/admin/clients";
        } else {
            const errData = await response.json();
            alert("Error: " + (errData.error || "Fallo en el servidor"));
        }
    } catch (err) {
        console.error("❌ Error de red:", err);
    }
}

// Funciones de Modal
function openDeleteModal() { document.getElementById('delete-modal').classList.remove('hidden'); }
function closeDeleteModal() { document.getElementById('delete-modal').classList.add('hidden'); }