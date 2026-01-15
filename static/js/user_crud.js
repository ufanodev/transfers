/**
 * @file user_crud.js
 * @description Lógica para gestionar un único usuario (POST/PUT/DELETE)
 */

document.addEventListener('DOMContentLoaded', async () => {
    // 1. Obtener ID de la URL si existe (?id=XXX)
    const params = new URLSearchParams(window.location.search);
    const userId = params.get('id');

    if (userId) {
        // MODO EDICIÓN
        prepareEditMode(userId);
    } else {
        // MODO CREACIÓN
        prepareCreateMode();
    }

    // 2. Manejar el envío del formulario
    document.getElementById('crud-form').addEventListener('submit', (e) => {
        e.preventDefault();
        saveUser(userId);
    });

    // 3. Manejar eliminación
    document.getElementById('btn-delete').addEventListener('click', () => {
        if(confirm("¿Estás seguro de eliminar este usuario? Esta acción es irreversible.")) {
            deleteUser(userId);
        }
    });
});

async function prepareEditMode(id) {
    document.getElementById('page-title').textContent = "Editar Usuario";
    document.getElementById('form-subtitle').textContent = "Actualizar Perfil";
    document.getElementById('pass-hint').classList.remove('hidden');
    document.getElementById('btn-delete').classList.remove('hidden');
    document.getElementById('userId').value = id;

    // Cargar datos actuales
    try {
        const res = await fetch(`/api/v1/users/${id}`);
        if (res.ok) {
            const user = await res.json();
            document.getElementById('username').value = user.username;
            document.getElementById('email').value = user.email;
            document.getElementById('role').value = user.role;
            document.getElementById('is_active').checked = user.is_active;
        }
    } catch (e) { console.error("Error al cargar usuario", e); }
}

function prepareCreateMode() {
    document.getElementById('page-title').textContent = "Nuevo Usuario";
    document.getElementById('form-subtitle').textContent = "Crear Perfil";
}

async function saveUser(id) {
    const userData = {
        username: document.getElementById('username').value,
        email: document.getElementById('email').value,
        role: document.getElementById('role').value,
        is_active: document.getElementById('is_active').checked
    };

    const pass = document.getElementById('password').value;
    if (pass) userData.password = pass;

    const url = id ? `/api/v1/users/${id}` : '/api/v1/users';
    const method = id ? 'PUT' : 'POST';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        });

        if (res.ok) {
            alert("Operación exitosa");
            window.location.href = "/dashboard/users";
        } else {
            const err = await res.json();
            alert("Error: " + err.error);
        }
    } catch (e) { console.error(e); }
}

async function deleteUser(id) {
    const res = await fetch(`/api/v1/users/${id}`, { method: 'DELETE' });
    if (res.ok) {
        window.location.href = "/dashboard/users";
    }
}