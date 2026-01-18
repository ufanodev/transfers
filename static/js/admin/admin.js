/**
 * @file admin.js
 * @description Lógica central del Panel de Administración de TakeUs.
 * Gestiona telemetría, estado de red y utilidades de depuración.
 */

const ADMIN_CONFIG = {
    refreshInterval: 30000, // Actualizar stats cada 30 segundos
    apiBase: '/api/v1'
};

document.addEventListener('DOMContentLoaded', () => {
    const adminStyle = "color: #3b82f6; font-weight: bold; background: #1e293b; padding: 4px 10px; border-radius: 5px; border: 1px solid #3b82f6;";
    console.log("%c TakeUs ADMIN %c Panel de Control Operativo Listo", adminStyle, "color: #64748b; font-style: italic;");

    // Inicialización
    initDashboard();

    // Configurar refresco automático
    setInterval(() => {
        refreshDatabaseStats();
        checkSystemStatus();
    }, ADMIN_CONFIG.refreshInterval);
});

/**
 * Orquestador de carga inicial
 */
async function initDashboard() {
    await checkSystemStatus();
    await refreshDatabaseStats();
}

/**
 * Consulta la API para obtener el conteo de registros en cada tabla.
 */
async function refreshDatabaseStats() {
    const statusText = document.getElementById('admin-status');
    console.log("[ADMIN] 🔃 Sincronizando estadísticas...");

    const resources = [
        { id: 'count-bookings', endpoint: 'bookings' },
        { id: 'count-drivers', endpoint: 'drivers' },
        { id: 'count-payments', endpoint: 'payments' },
        { id: 'count-vehicles', endpoint: 'vehicles' },
        { id: 'count-companies', endpoint: 'companies' }
    ];

    // Ejecutar todas las peticiones en paralelo para mayor velocidad
    const promises = resources.map(async (item) => {
        try {
            const response = await fetch(`${ADMIN_CONFIG.apiBase}/${item.endpoint}`);
            
            if (response.status === 401) {
                handleUnauthorized();
                return;
            }

            if (response.ok) {
                const data = await response.json();
                const countElement = document.getElementById(item.id);
                if (countElement) {
                    const count = Array.isArray(data) ? data.length : (data.count || 0);
                    
                    // Efecto de transición si el valor cambia
                    if (countElement.textContent !== count.toString()) {
                        countElement.textContent = count;
                        countElement.classList.add('animate-pulse', 'text-blue-400');
                        setTimeout(() => countElement.classList.remove('animate-pulse', 'text-blue-400'), 2000);
                    }
                }
            }
        } catch (error) {
            console.error(`[ADMIN] Error en ${item.endpoint}:`, error);
            const el = document.getElementById(item.id);
            if (el) el.textContent = "!";
        }
    });

    await Promise.all(promises);
}

/**
 * Verifica latencia y disponibilidad del servidor Go
 */
async function checkSystemStatus() {
    const statusLabel = document.getElementById('admin-status');
    if (!statusLabel) return;

    try {
        const start = performance.now();
        const res = await fetch(`${ADMIN_CONFIG.apiBase}/users/me`);
        const end = performance.now();
        
        if (res.ok) {
            const latency = Math.round(end - start);
            statusLabel.innerHTML = `
                <span class="flex items-center gap-2">
                    <span class="relative flex h-2 w-2">
                        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                        <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
                    </span>
                    ONLINE <span class="text-[10px] text-slate-400 font-mono">(${latency}ms)</span>
                </span>`;
        } else if (res.status === 401) {
            handleUnauthorized();
        }
    } catch (err) {
        statusLabel.innerHTML = `<span class="text-red-500 font-bold">● OFFLINE</span>`;
    }
}

/**
 * Manejo centralizado de expulsión por sesión expirada
 */
function handleUnauthorized() {
    console.warn("🚫 Sesión administrativa expirada o no válida.");
    // Pequeño delay para que el admin pueda ver el log antes de redirigir
    setTimeout(() => {
        window.location.href = "/logout";
    }, 500);
}

/**
 * Utilidad para depuración rápida desde consola: adminTestAPI('rides')
 */
window.adminTestAPI = async function(endpoint) {
    console.group(`🔍 Inspeccionando: ${endpoint.toUpperCase()}`);
    try {
        const res = await fetch(`${ADMIN_CONFIG.apiBase}/${endpoint}`);
        const data = await res.json();
        console.table(data);
        return data;
    } catch (err) {
        console.error("Error en consulta manual:", err);
    }
    console.groupEnd();
};