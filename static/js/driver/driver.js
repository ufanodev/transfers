/**
 * @file driver_dashboard.js
 * @description Lógica central del Panel de Conductor.
 * Gestiona el estado del turno, perfil del conductor y telemetría de API.
 */

document.addEventListener('DOMContentLoaded', () => {
    // Estilo visual para la terminal de depuración
    const driverStyle = "color: #10b981; font-weight: bold; background: #1a1a1a; padding: 3px 8px; border-radius: 5px; border: 1px solid #10b981;";
    console.log("%c TakeUs DRIVER %c Sistema operativo y conectado.", driverStyle, "color: #64748b; font-style: italic;");
    
    // 1. Cargar datos del conductor (Nombre, Matrícula, Empresa)
    initDriverProfile();

    // 2. Gestión de parámetros en URL (Rutas directas)
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has('ride_id')) {
        console.log("📍 Ruta específica detectada para seguimiento:", urlParams.get('ride_id'));
    }
});

/**
 * Carga la información del conductor desde la API de Go
 */
async function initDriverProfile() {
    const nameDisplay = document.getElementById('driver-name');
    
    try {
        // Llamada al endpoint de "mi perfil"
        const res = await fetch('/api/v1/users/me'); 
        
        if (res.ok) {
            const user = await res.json();
            if (nameDisplay) nameDisplay.textContent = user.first_name || 'Conductor';
            console.log(`%c[AUTH] ✅ Identidad confirmada: ${user.email}`, "color: #10b981");
        } else if (res.status === 401) {
            handleAuthError();
        }
    } catch (err) {
        console.error("[DRIVER] 💥 Fallo al cargar perfil:", err);
    }
}

/**
 * Realiza una consulta a la API y muestra el resultado en la terminal visual.
 * @param {string} endpoint - Recurso a consultar (rides, bookings, events)
 */
async function testAPI(endpoint) {
    const now = new Date().toLocaleTimeString();
    console.log(`%c[DRIVER-LOG] %cSolicitando: ${endpoint.toUpperCase()} [${now}]`, "color: #10b981; font-weight: bold", "color: #94a3b8");

    const resultDiv = document.getElementById('api-result');
    const pre = resultDiv?.querySelector('pre');
    
    try {
        console.time(`Latencia API: ${endpoint}`);
        const response = await fetch(`/api/v1/${endpoint}`);
        console.timeEnd(`Latencia API: ${endpoint}`);

        if (response.status === 401) {
            handleAuthError();
            return;
        }

        const data = await response.json();

        // Actualización de la terminal en el HTML (si existe)
        if (resultDiv && pre) {
            resultDiv.classList.remove('hidden');
            pre.textContent = `// TERMINAL TAKEUS DRIVER\n// SYNC: ${new Date().toLocaleString()}\n// DATA: /api/v1/${endpoint}\n\n` + 
                              JSON.stringify(data, null, 4);
            
            resultDiv.scrollIntoView({ behavior: 'smooth', block: 'start' });
        } else {
            // Si no hay div en el HTML, lo sacamos por consola de forma elegante
            console.table(data);
        }

    } catch (error) {
        console.error("%c[ERROR CRÍTICO] %cNo se pudo conectar con la API:", "color: #ef4444; font-weight: bold", "color: red", error);
        alert("Fallo de red. Compruebe su conexión.");
    }
}

/**
 * Simulación de cambio de estado operativo (Disponible/Fuera de Turno)
 */
async function toggleAvailability() {
    console.log("📡 Sincronizando disponibilidad con el servidor...");
    // Futura implementación: PATCH /api/v1/drivers/status
    const statusText = document.querySelector('.text-green-500');
    if (statusText) {
        statusText.classList.toggle('text-green-500');
        statusText.classList.toggle('text-gray-500');
        statusText.innerHTML = statusText.classList.contains('text-green-500') ? 
            '<div class="w-2 h-2 rounded-full bg-green-500"></div> En Línea' : 
            '<div class="w-2 h-2 rounded-full bg-gray-500"></div> Fuera de Servicio';
    }
}

/**
 * Redirección centralizada por fallo de autenticación
 */
function handleAuthError() {
    console.warn("🚫 Sesión no válida. Expulsando del panel operativo.");
    window.location.href = "/";
}