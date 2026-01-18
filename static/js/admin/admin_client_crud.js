/**
 * @file admin_client_crud.js
 * @description Lógica de alta y edición para el modelo Client.
 */

document.addEventListener('DOMContentLoaded', async () => {
    const params = new URLSearchParams(window.location.search);
    const clientId = params.get('id');
    const pageTitle = document.getElementById('page-title');

    console.log(" [%] Inicializando CRUD Clientes...");

    if (clientId) {
        console.log(` [MODE] Editando Cliente ID: ${clientId}`);
        if(pageTitle) pageTitle.textContent = "Actualizar Cliente";
        document.getElementById('btn-delete')?.classList.remove('hidden');
        loadClientData(clientId);
    } else {
        console.log(" [MODE] Creación de nuevo cliente.");
        if(pageTitle) pageTitle.textContent = "Registrar Cliente";
    }

    document.getElementById('crud-form').addEventListener('submit', (e) => {
        e.preventDefault();
        saveClient(clientId);
    });

    document.getElementById('btn-delete')?.addEventListener('click', () => {
        if(confirm("¿Deseas eliminar este cliente de forma permanente?")) deleteClient(clientId);
    });
});

async function loadClientData(id) {
    try {
        const res = await fetch(`/api/v1/clients/${id}`);
        if (!res.ok) throw new Error("No se pudo obtener el cliente");
        
        const data = await res.json();
        console.log(" [DATA LOADED]:", data);
        
        document.getElementById('user_id').value = data.user_id;
        document.getElementById('full_name').value = data.full_name;
        document.getElementById('email').value = data.email;
        document.getElementById('phone').value = data.phone;
        
        // Formatear JSON de preferencias para el textarea
        document.getElementById('preferences').value = typeof data.preferences === 'string' 
            ? data.preferences 
            : JSON.stringify(data.preferences || {}, null, 2);

    } catch (e) {
        console.error(" [ERROR]:", e);
        alert("Error al cargar los datos del cliente.");
    }
}

async function saveClient(id) {
    console.log(" [SYSTEM] Iniciando proceso de guardado...");
    
    // Preparar objeto de datos basado en el modelo GORM
    const clientData = {
        user_id: parseInt(document.getElementById('user_id').value),
        full_name: document.getElementById('full_name').value,
        email: document.getElementById('email').value,
        phone: document.getElementById('phone').value,
        preferences: document.getElementById('preferences').value || "{}"
    };

    const url = id ? `/api/v1/clients/${id}` : '/api/v1/clients';
    const method = id ? 'PUT' : 'POST';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(clientData)
        });

        const result = await res.json();

        if (res.ok) {
            console.log(" [SUCCESS]:", result);
            alert("✅ Cliente guardado correctamente.");
            window.location.href = "/dashboard/admin/clients/";
        } else {
            console.error(" [SERVER ERROR]:", result);
            alert("❌ Error: " + (result.error || "No se pudo completar la operación"));
        }
    } catch (e) {
        console.error(" [FETCH ERROR]:", e);
        alert("❌ Error crítico de conexión.");
    }
}

async function deleteClient(id) {
    try {
        const res = await fetch(`/api/v1/clients/${id}`, { method: 'DELETE' });
        if (res.ok) {
            console.log(" [DELETE SUCCESS]");
            window.location.href = "/dashboard/admin/clients/";
        } else {
            alert("Error al eliminar.");
        }
    } catch (e) {
        console.error(" [DELETE ERROR]:", e);
    }
}