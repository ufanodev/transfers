document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const userId = urlParams.get('id');

    if (userId) {
        console.log("[CRUD] 📝 Modo: MODIFICAR USUARIO (ID:", userId, ")");
        await setupEditMode(userId);
    } else {
        console.log("[CRUD] ✨ Modo: NUEVO USUARIO");
        document.getElementById('page-title').textContent = "Nuevo Usuario";
        document.getElementById('btn-save').textContent = "Crear Usuario";
    }

    // Listener del Formulario
    document.getElementById('crud-form').addEventListener('submit', handleFormSubmit);
});

/**
 * Prepara el formulario para editar
 */
async function setupEditMode(id) {
    try {
        const response = await fetch(`/api/v1/users/${id}`);
        if (!response.ok) throw new Error("No se pudo obtener el usuario");
        
        const u = await response.json();

        // Llenar campos
        document.getElementById('userId').value = u.id || u.ID;
        document.getElementById('username').value = u.username;
        document.getElementById('email').value = u.email;
        document.getElementById('role').value = u.role;
        document.getElementById('is_active').checked = u.is_active;

        // Cambiar textos visuales
        document.getElementById('page-title').textContent = "Modificar Usuario";
        document.getElementById('btn-save').textContent = "Guardar Cambios";
        document.getElementById('pass-hint').classList.remove('hidden');
        document.getElementById('btn-delete-trigger').classList.remove('hidden');
        
        // Configurar modal de borrado
        document.getElementById('confirm-user-name').textContent = u.username;
        document.getElementById('btn-delete-trigger').onclick = () => openDeleteModal(u);

    } catch (err) {
        alert("Error cargando datos: " + err.message);
    }
}

/**
 * Envío de datos (POST para crear / PUT para modificar)
 */
async function handleFormSubmit(e) {
    e.preventDefault();
    const userId = document.getElementById('userId').value;
    
    const payload = {
        username: document.getElementById('username').value,
        email: document.getElementById('email').value,
        role: document.getElementById('role').value,
        is_active: document.getElementById('is_active').checked
    };

    const password = document.getElementById('password').value;
    if (password) payload.password = password;

    const method = userId ? 'PUT' : 'POST';
    const url = userId ? `/api/v1/users/${userId}` : '/api/v1/users';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.ok) {
            alert(userId ? "✅ Usuario actualizado" : "✅ Usuario creado");
            window.location.href = "/admin/users";
        } else {
            const err = await res.json();
            alert("Error: " + err.error);
        }
    } catch (err) {
        alert("Fallo de red");
    }
}

// LÓGICA DE BORRADO
function openDeleteModal() {
    document.getElementById('delete-modal').classList.remove('hidden');
}

function closeDeleteModal() {
    document.getElementById('delete-modal').classList.add('hidden');
}

async function confirmFinalDelete() {
    const userId = document.getElementById('userId').value;
    const inputUser = document.getElementById('confirm-username').value;
    const inputEmail = document.getElementById('confirm-email').value;
    const inputKey = document.getElementById('confirm-key').value;

    const realUser = document.getElementById('username').value;
    const realEmail = document.getElementById('email').value;

    if (inputUser === realUser && inputEmail === realEmail && inputKey.toUpperCase() === 'ELIMINAR') {
        try {
            const res = await fetch(`/api/v1/users/${userId}`, { method: 'DELETE' });
            if (res.ok) {
                alert("🗑️ Usuario eliminado correctamente");
                window.location.href = "/admin/users";
            }
        } catch (err) { alert("Error al borrar"); }
    } else {
        alert("❌ Los datos de confirmación no coinciden.");
    }
}