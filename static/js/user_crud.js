document.addEventListener('DOMContentLoaded', async () => {
    const params = new URLSearchParams(window.location.search);
    const userId = params.get('id');

    // Evitar errores de null si los elementos no existen
    const pageTitle = document.getElementById('page-title');
    const formSubtitle = document.getElementById('form-subtitle');

    if (userId) {
        console.log("Iniciando modo edición para ID:", userId);
        if(pageTitle) pageTitle.textContent = "EDITAR USUARIO";
        if(formSubtitle) formSubtitle.textContent = "Actualizar datos";
        
        document.getElementById('pass-hint')?.classList.remove('hidden');
        document.getElementById('btn-delete')?.classList.remove('hidden');
        
        loadUserData(userId);
    } else {
        console.log("Iniciando modo creación");
        if(pageTitle) pageTitle.textContent = "NUEVO USUARIO";
        if(formSubtitle) formSubtitle.textContent = "Crear perfil";
    }

    document.getElementById('crud-form').addEventListener('submit', (e) => {
        e.preventDefault();
        saveUser(userId);
    });

    document.getElementById('btn-delete')?.addEventListener('click', () => {
        if(confirm("¿Eliminar usuario?")) deleteUser(userId);
    });
});

async function loadUserData(id) {
    try {
        const res = await fetch(`/api/v1/users/${id}`);
        const data = await res.json();
        console.log("Datos cargados:", data);
        
        document.getElementById('username').value = data.username;
        document.getElementById('email').value = data.email;
        document.getElementById('role').value = data.role;
        document.getElementById('is_active').checked = data.is_active;
    } catch (e) {
        console.error("Error cargando usuario:", e);
    }
}

async function saveUser(id) {
    console.log("Enviando formulario...");
    
    const userData = {
        username: document.getElementById('username').value,
        email: document.getElementById('email').value,
        role: document.getElementById('role').value,
        is_active: document.getElementById('is_active').checked
    };

    const password = document.getElementById('password').value;
    if (password) userData.password = password;

    const url = id ? `/api/v1/users/${id}` : '/api/v1/users';
    const method = id ? 'PUT' : 'POST';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        });

        const result = await res.json();

        if (res.ok) {
            console.log("Respuesta Exitosa:", result);
            alert("✅ ¡Guardado con éxito!");
            window.location.href = "/dashboard/users";
        } else {
            console.error("Respuesta Error:", result);
            alert("❌ Error: " + (result.error || "No se pudo guardar"));
        }
    } catch (e) {
        console.error("Error de red:", e);
        alert("❌ Error crítico de conexión");
    }
}