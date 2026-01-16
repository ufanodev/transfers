/**
 * @file admin_client_crud.js
 * @description Gestión de Clientes con creación automática de cuenta de usuario.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("🚀 Módulo CRUD Clientes (Auto-User) Iniciado");
    
    const params = new URLSearchParams(window.location.search);
    const clientId = params.get('id');

    const pageTitle = document.getElementById('page-title');
    const formSubtitle = document.getElementById('form-subtitle');

    if (clientId) {
        if(pageTitle) pageTitle.textContent = "EDITAR CLIENTE";
        if(formSubtitle) formSubtitle.textContent = `Ficha técnica del registro #${clientId}`;
        await loadClientData(clientId);
    } else {
        if(pageTitle) pageTitle.textContent = "NUEVO CLIENTE";
        if(formSubtitle) formSubtitle.textContent = "Se creará automáticamente una cuenta de acceso con el email";
    }

    const crudForm = document.getElementById('crud-form');
    if (crudForm) {
        crudForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            await saveClient(clientId);
        });
    }
});

/**
 * Carga los datos del cliente para edición
 */
async function loadClientData(id) {
    try {
        const res = await fetch(`/api/v1/clients/${id}`);
        if (!res.ok) throw new Error("Cliente no encontrado");
        
        const data = await res.json();
        
        // Mapeo de campos a los inputs del HTML
        if(document.getElementById('full_name')) document.getElementById('full_name').value = data.full_name || data.FullName || "";
        if(document.getElementById('email')) document.getElementById('email').value = data.email || data.Email || "";
        if(document.getElementById('phone')) document.getElementById('phone').value = data.phone || data.Phone || "";
        
        const prefs = data.preferences || data.Preferences;
        const elPrefs = document.getElementById('preferences');
        if(elPrefs) {
            // Intentamos mostrar el texto limpio si es un JSON de notas
            try {
                const parsed = typeof prefs === 'string' ? JSON.parse(prefs) : prefs;
                elPrefs.value = parsed.notes ? parsed.notes : JSON.stringify(parsed, null, 2);
            } catch (e) {
                elPrefs.value = prefs || "";
            }
        }
    } catch (e) {
        console.error("❌ Error al cargar datos:", e);
    }
}

/**
 * Guarda el cliente (POST crea usuario + cliente / PUT solo actualiza cliente)
 */
async function saveClient(id) {
    const btnSave = document.getElementById('btn-save');
    
    // Captura de valores
    const valFullName = document.getElementById('full_name')?.value;
    const valEmail = document.getElementById('email')?.value;
    const valPhone = document.getElementById('phone')?.value;
    const valPrefsRaw = document.getElementById('preferences')?.value.trim();

    // Validación básica en frontend
    if (!valFullName || !valEmail || !valPhone) {
        alert("⚠️ Nombre, Email y Teléfono son obligatorios para crear la cuenta.");
        return;
    }

    if(btnSave) btnSave.disabled = true;

    // Preparamos el payload. 
    // NOTA: No enviamos user_id, el servidor lo generará en el POST.
    const clientData = {
        full_name: valFullName,
        email: valEmail,
        phone: valPhone,
        preferences: valPrefsRaw || "{}"
    };

    const url = id ? `/api/v1/clients/${id}` : '/api/v1/clients';
    const method = id ? 'PUT' : 'POST';

    console.log(`[API] Iniciando petición ${method}...`);

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(clientData)
        });

        const result = await res.json();

        if (res.ok) {
            if (!id) {
                // Si es un registro nuevo, informamos de las credenciales
                alert(`✅ ¡Éxito! Cliente y Usuario creados.\n\nUsuario: ${valEmail}\nPassword: ${valPhone}`);
            } else {
                alert("✅ Datos actualizados correctamente.");
            }
            window.location.href = "/dashboard/clients";
        } else {
            // Manejo de errores controlados (Email duplicado, etc)
            alert("❌ Error: " + (result.error || "No se pudo procesar el registro"));
        }
    } catch (e) {
        console.error("[FETCH ERROR]", e);
        alert("❌ Error de red: El servidor no responde.");
    } finally {
        if(btnSave) btnSave.disabled = false;
    }
}