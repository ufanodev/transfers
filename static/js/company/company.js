/**
 * @file company_dashboard.js
 * @description Lógica de control para el Panel de Empresa (Business).
 * Gestiona la monitorización de flota, métricas financieras y estados operativos.
 */

document.addEventListener('DOMContentLoaded', () => {
    // Log de inicialización con estilo corporativo
    const companyStyle = "color: #6366f1; font-weight: bold; background: #0f172a; padding: 3px 8px; border-radius: 5px; border: 1px solid #6366f1;";
    console.log("%c TakeUs BUSINESS %c Terminal de administración activa.", companyStyle, "color: #64748b; font-style: italic;");

    // 1. Cargar datos de la empresa y KPIs
    initCompanyProfile();
    updateDashboardMetrics();

    // 2. Escuchar eventos de navegación interna
    setupBusinessNavigation();
});

/**
 * Obtiene los datos de la empresa y el administrador logueado
 */
async function initCompanyProfile() {
    const companyDisplay = document.getElementById('company-name');
    
    try {
        const res = await fetch('/api/v1/users/me'); 
        
        if (res.ok) {
            const data = await res.json();
            // Si el usuario pertenece a una empresa, mostramos el nombre comercial
            if (companyDisplay) {
                companyDisplay.textContent = data.company_name || "TakeUs Partner";
            }
            console.log(`%c[CORP] ✅ Identidad verificada para: ${data.email}`, "color: #6366f1");
        } else if (res.status === 401) {
            handleAuthError();
        }
    } catch (err) {
        console.error("[COMPANY] 💥 Error en carga de perfil corporativo:", err);
    }
}

/**
 * Función para actualizar los contadores del Dashboard (Drivers, Reservas, etc.)
 */
async function updateDashboardMetrics() {
    console.log("[METRICS] 📊 Actualizando indicadores de rendimiento...");
    
    try {
        // Ejemplo de llamadas paralelas para optimizar velocidad
        const [resDrivers, resBookings] = await Promise.all([
            fetch('/api/v1/drivers'),
            fetch('/api/v1/bookings')
        ]);

        if (resDrivers.ok && resBookings.ok) {
            const drivers = await resDrivers.json();
            const bookings = await resBookings.json();

            // Aquí inyectarías los datos en las cards del HTML
            // document.getElementById('total-drivers').textContent = drivers.length;
            console.log(`[DATA] 🚛 Conductores: ${drivers.length} | 📅 Reservas: ${bookings.length}`);
        }
    } catch (err) {
        console.warn("[METRICS] No se pudieron sincronizar los KPIs en tiempo real.");
    }
}

/**
 * Ejecuta pruebas de API específicas de administración
 * @param {string} endpoint - Recurso (drivers, vehicles, payments)
 */
async function testAPI(endpoint) {
    const resultDiv = document.getElementById('api-result');
    const pre = resultDiv?.querySelector('pre');
    
    console.log(`%c[ADMIN-API] %cConsultando recurso: ${endpoint.toUpperCase()}`, "color: #6366f1; font-weight: bold", "color: #94a3b8");

    try {
        const response = await fetch(`/api/v1/${endpoint}`);
        
        if (response.status === 401) {
            handleAuthError();
            return;
        }

        const data = await response.json();

        if (resultDiv && pre) {
            resultDiv.classList.remove('hidden');
            pre.textContent = `// BUSINESS INTEL LOG\n// DATE: ${new Date().toLocaleString()}\n// RESOURCE: /api/v1/${endpoint}\n\n` + 
                              JSON.stringify(data, null, 4);
            resultDiv.scrollIntoView({ behavior: 'smooth' });
        } else {
            console.table(data);
        }

    } catch (error) {
        console.error("[ERROR] Fallo en la consulta administrativa:", error);
    }
}

/**
 * Configuración de interacción para el menú lateral
 */
function setupBusinessNavigation() {
    const navLinks = document.querySelectorAll('aside nav a');
    navLinks.forEach(link => {
        link.addEventListener('click', () => {
            // Podrías añadir lógica para guardar el último tab visitado
            const section = link.innerText.trim();
            console.log(`[NAV] Accediendo a sección: ${section}`);
        });
    });
}

/**
 * Gestión de cierre de sesión por expiración
 */
function handleAuthError() {
    alert("La sesión corporativa ha expirado por seguridad.");
    window.location.href = "/";
}