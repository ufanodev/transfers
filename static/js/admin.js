/**
 * @file admin.js
 * @description Gestión de la interfaz del Panel de Administración.
 * Proporciona funcionalidad para la inspección rápida de datos de los modelos
 * a través de la API v1.
 */

/**
 * Realiza una petición GET a la API para un endpoint específico y muestra los resultados
 * en la terminal visual del dashboard.
 * * @param {string} endpoint - El recurso a consultar (ej. 'clients', 'bookings', etc.)
 */
async function testAPI(endpoint) {
    const resultDiv = document.getElementById('api-result');
    const pre = resultDiv.querySelector('pre');
    const title = document.getElementById('api-title');
    
    try {
        // Realizamos la llamada asíncrona a la API de Go
        const res = await fetch(`/api/v1/${endpoint}`);
        
        // Verificamos si la sesión sigue siendo válida (Middleware JWT de Go)
        if (res.status === 401) {
            // Si el servidor indica que no hay autorización, redirigimos al login
            window.location.href = "/";
            return;
        }
        
        // Parseamos la respuesta JSON
        const data = await res.json();
        
        // Hacemos visible el contenedor de resultados
        resultDiv.classList.remove('hidden');
        
        // Actualizamos el título de la terminal con el recurso consultado
        title.textContent = "VISTA PREVIA DE DATOS: " + endpoint.toUpperCase();
        
        // Formateamos el JSON con 4 espacios de sangría para que sea legible
        pre.textContent = JSON.stringify(data, null, 4);
        
        // Realizamos un scroll suave hacia el área de resultados
        resultDiv.scrollIntoView({ behavior: 'smooth' });

    } catch (err) {
        // Capturamos y mostramos errores de red o del servidor
        console.error("Error Admin API:", err);
        alert("Ocurrió un error al intentar recuperar los datos de " + endpoint);
    }
}

// Log para confirmar la carga del script en la consola del navegador
console.log("Admin module loaded successfully.");