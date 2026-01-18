/**
 * @file client.js
 * @description Lógica central del Panel de Cliente.
 * Gestiona la carga de datos del perfil, navegación y pruebas de API.
 * @version 1.4.0
 */

document.addEventListener('DOMContentLoaded', () => {
    console.log("%c TakeUs %c Panel de Cliente Iniciado", "color: #10b981; font-weight: bold", "color: gray");

    // 1. Inicializar la interfaz con datos del usuario
    loadClientProfile();

    // 2. Escuchar clics en el menú lateral para efectos visuales o carga dinámica
    setupNavigation();
});

/**
 * Carga los datos básicos del cliente desde la API para personalizar el Dashboard
 */
async function loadClientProfile() {
    const nameDisplay = document.getElementById('user-name');
    
    try {
        // Asumiendo que tenemos un endpoint para 'mi perfil'
        const res = await fetch('/api/v1/users/me'); 
        
        if (res.ok) {
            const user = await res.json();
            if (nameDisplay) nameDisplay.textContent = user.first_name || 'Cliente';
            console.log("[CLIENT] ✅ Perfil cargado correctamente.");
        } else {
            // Si falla el fetch de usuario, es probable que la sesión haya expirado
            if (res.status === 401) handleSessionExpired();
            if (nameDisplay) nameDisplay.textContent = "Usuario";
        }
    } catch (err) {
        console.error("[CLIENT] 💥 Error al conectar con la API de perfil:", err);
    }
}

/**
 * Función genérica para testear endpoints desde la consola o botones
 * @param {string} endpoint - El recurso a consultar (ej: 'bookings', 'payments')
 */
async function testAPI(endpoint) {
    const resultDiv = document.getElementById('api-result');
    // Si no existe el div de resultados en el HTML, lo buscamos o creamos uno temporal
    if (!resultDiv) {
        console.warn(`[API-TEST] Intentando consultar ${endpoint}, pero no hay contenedor de resultados.`);
        return;
    }

    const pre = resultDiv.querySelector('pre');
    console.log(`[API-TEST] Consultando: /api/v1/${endpoint}...`);

    try {
        const res = await fetch(`/api/v1/${endpoint}`);
        
        if (res.status === 401) {
            handleSessionExpired();
            return;
        }

        const data = await res.json();
        
        // Mostrar visualmente el JSON
        resultDiv.classList.remove('hidden');
        if (pre) {
            pre.textContent = `// RESULTADO DE: ${endpoint.toUpperCase()}\n` + JSON.stringify(data, null, 4);
        }
        
        resultDiv.scrollIntoView({ behavior: 'smooth' });

    } catch (err) {
        console.error(`[API-TEST] ❌ Error en endpoint ${endpoint}:`, err);
        alert("Fallo al conectar con el servidor.");
    }
}

/**
 * Maneja la redirección cuando la sesión ya no es válida
 */
function handleSessionExpired() {
    console.warn("[AUTH] Sesión expirada o inválida. Redirigiendo al login...");
    // Podríamos añadir un parámetro para volver aquí después del login: /?redirect=/client
    window.location.href = "/";
}

/**
 * Configura comportamientos para los links del sidebar
 */
function setupNavigation() {
    const links = document.querySelectorAll('aside nav a');
    
    links.forEach(link => {
        link.addEventListener('click', (e) => {
            // Si el link es solo para pruebas (sin URL real), podríamos prevenir el default
            // Por ahora solo marcamos el "activo" visualmente
            if (link.getAttribute('href') === '#') {
                e.preventDefault();
                console.log("[NAV] Opción en desarrollo: " + link.innerText);
            }
        });
    });
}