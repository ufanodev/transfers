/**
 * @file admin.js
 * @description Gestión de la interfaz del Panel de Administración.
 * Proporciona funcionalidad para la inspección rápida de datos de los modelos
 * a través de la API v1 con logs de auditoría en consola.
 */

/**
 * Realiza una petición GET a la API para un endpoint específico y muestra los resultados
 * en la terminal visual del dashboard.
 * @param {string} endpoint - El recurso a consultar (ej. 'clients', 'bookings', etc.)
 */
async function testAPI(endpoint) {
    const timestamp = new Date().toLocaleTimeString();
    console.log(`%c[API CALL] %cIniciando petición a /api/v1/${endpoint} a las ${timestamp}`, "color: #10b981; font-weight: bold", "color: gray");

    const resultDiv = document.getElementById('api-result');
    const pre = resultDiv.querySelector('pre');
    const title = document.getElementById('api-title');
    
    try {
        // Marcamos inicio de carga visualmente si fuera necesario
        console.time(`Tiempo de respuesta para ${endpoint}`);

        // Realizamos la llamada asíncrona a la API de Go
        const res = await fetch(`/api/v1/${endpoint}`);
        
        console.timeEnd(`Tiempo de respuesta para ${endpoint}`);
        console.log(`%c[API RES] %cStatus: ${res.status} ${res.statusText}`, "color: #10b981; font-weight: bold", "color: white");

        // Verificamos si la sesión sigue siendo válida (Middleware JWT de Go)
        if (res.status === 401) {
            console.warn("%c[AUTH] %cSesión expirada o no autorizada. Redirigiendo al login...", "color: #ef4444; font-weight: bold", "color: gray");
            window.location.href = "/";
            return;
        }

        if (!res.ok) {
            throw new Error(`Error en la respuesta del servidor: ${res.status}`);
        }
        
        // Parseamos la respuesta JSON
        const data = await res.json();
        console.log(`%c[DATA] %cRegistros obtenidos: ${Array.isArray(data) ? data.length : '1'}`, "color: #10b981; font-weight: bold", "color: gray");
        
        // Hacemos visible el contenedor de resultados
        if (resultDiv) {
            resultDiv.classList.remove('hidden');
            
            // Actualizamos el título de la terminal
            if (title) title.textContent = "INSPECCIÓN: " + endpoint.toUpperCase();
            
            // Formateamos el JSON para que sea legible
            pre.textContent = `// Resultado de la consulta: ${endpoint.toUpperCase()}\n// Fecha: ${new Date().toLocaleString()}\n\n` + JSON.stringify(data, null, 4);
            
            // Realizamos un scroll suave hacia el área de resultados
            resultDiv.scrollIntoView({ behavior: 'smooth' });
        }

    } catch (err) {
        // Capturamos y mostramos errores de red o del servidor
        console.error(`%c[FATAL ERROR] %cOcurrió un problema con ${endpoint}:`, "color: #ef4444; font-weight: bold", "color: red", err);
        alert("Ocurrió un error al intentar recuperar los datos de " + endpoint);
    }
}

/**
 * Log de auditoría inicial para confirmar carga del módulo
 */
(() => {
    const statusStyle = "color: #10b981; font-weight: bold; background: #000; padding: 2px 5px; border-radius: 3px;";
    console.log("%c TakeUs %c Módulo Admin cargado y listo.", statusStyle, "color: gray; font-style: italic;");
    console.log(`[SYS] Entorno: ${window.location.hostname} | Path: ${window.location.pathname}`);
})();