/**
 * @file login.js
 * @description Gestión de la autenticación en el lado del cliente.
 * Este script captura las credenciales del formulario de acceso, las envía al 
 * backend de Go mediante una petición asíncrona y gestiona la redirección 
 * inicial basándose en la respuesta del servidor.
 * * @author Sistema Transfers
 * @version 1.0.0
 */

// Espera a que el DOM esté completamente cargado para asegurar que los elementos existan
document.addEventListener('DOMContentLoaded', () => {
    
    // Referencias a elementos del DOM
    const loginForm = document.getElementById('login-form');
    const errorBox = document.getElementById('error-message');

    // Verificación de seguridad para evitar errores si el formulario no está en la página
    if (loginForm) {
        
        // Intercepta el evento de envío del formulario
        loginForm.addEventListener('submit', async (e) => {
            
            // Evita el refresco automático de la página (comportamiento por defecto del submit)
            e.preventDefault();
            
            // Limpia estados de error previos ocultando la caja de alerta
            errorBox.classList.add('hidden');

            // Captura de valores de los inputs
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            try {
                // Envío de datos al endpoint de Go mediante FETCH API
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json' 
                    },
                    // Conversión del objeto JS a cadena JSON para el cuerpo de la petición
                    body: JSON.stringify({ email, password })
                });

                // Si el servidor responde con un status 2xx (Éxito)
                if (response.ok) {
                    /**
                     * NOTA: En este punto, el servidor ya ha inyectado la Cookie HttpOnly 
                     * en el navegador. Redirigimos al dashboard centralizado donde Go 
                     * decidirá qué vista servir según el rol del usuario.
                     */
                    window.location.href = "/dashboard/";
                } else {
                    // Si el servidor responde con error (401, 400, etc.), muestra la alerta visual
                    errorBox.classList.remove('hidden');
                }
                
            } catch (error) {
                // Gestión de errores críticos de comunicación o red
                console.error("Error de comunicación con el servidor:", error);
            }
        });
    }
});