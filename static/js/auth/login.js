/**
 * @file login.js
 * @description Gestión de autenticación con Logs de Auditoría.
 * Captura credenciales, valida presencia de datos y gestiona redirección por roles.
 * @version 1.3.0
 */

document.addEventListener('DOMContentLoaded', () => {
    
    // Referencias a elementos del DOM
    const loginForm = document.getElementById('login-form');
    const errorBox = document.getElementById('error-message');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const submitBtn = loginForm?.querySelector('button[type="submit"]');

    if (loginForm) {
        console.log("[JS-AUTH] 🚀 Sistema de login inicializado.");

        loginForm.addEventListener('submit', async (e) => {
            // 1. Evitar el refresco de página por defecto
            e.preventDefault();
            
            // 2. Limpieza de estado visual y mensajes de error
            errorBox.classList.add('hidden');
            const errorSpan = errorBox.querySelector('span');

            // 3. Captura y Auditoría de valores (Trim para eliminar espacios accidentales)
            const emailValue = emailInput.value.trim();
            const passwordValue = passwordInput.value;

            console.log("[JS-AUTH] 📝 Capturando datos...");
            console.log(`   - Email: [${emailValue}]`);
            console.log(`   - Password: [${passwordValue ? "********" : "VACÍO"}]`);

            // Validación previa en cliente para ahorrar peticiones al servidor
            if (!emailValue || !passwordValue) {
                console.warn("[JS-AUTH] ⚠️ Intento de envío incompleto.");
                errorBox.classList.remove('hidden');
                if (errorSpan) errorSpan.textContent = "Por favor, rellene todos los campos.";
                return;
            }

            /**
             * 4. Preparación del Payload
             * IMPORTANTE: Usamos 'username' para que el struct LoginRequest en Go 
             * lo reconozca correctamente (binding:"required").
             */
            const payload = {
                username: emailValue, 
                password: passwordValue
            };

            console.log("[JS-AUTH] 📤 Enviando petición POST a /login:", payload);
            
            // Feedback visual: deshabilitar botón mientras se procesa la red
            if (submitBtn) {
                submitBtn.disabled = true;
                submitBtn.textContent = 'Verificando...';
            }

            try {
                // 5. Llamada asíncrona al Backend de Go
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json' 
                    },
                    body: JSON.stringify(payload)
                });

                console.log(`[JS-AUTH] 📡 Respuesta del servidor: Status ${response.status}`);

                const data = await response.json();

                if (response.ok) {
                    /**
                     * ÉXITO: El servidor Go ha validado y establecido la Cookie HttpOnly.
                     * Redirigimos al semáforo central (/dashboard) para que Go decida la vista final.
                     */
                    console.log("[JS-AUTH] ✅ Credenciales correctas. Redirigiendo a /dashboard...");
                    window.location.href = "/dashboard";
                } else {
                    // ERROR: Credenciales incorrectas, usuario no encontrado o fallo de servidor
                    console.error("[JS-AUTH] 🚫 Error:", data.error);
                    errorBox.classList.remove('hidden');
                    if (errorSpan) errorSpan.textContent = data.error || "Acceso denegado.";
                    
                    // Restaurar botón para reintento
                    if (submitBtn) {
                        submitBtn.disabled = false;
                        submitBtn.textContent = 'Iniciar sesión';
                    }
                }

            } catch (err) {
                // FALLO TÉCNICO: Servidor apagado, error de red o timeout
                console.error("[JS-AUTH] 💥 Fallo crítico en la comunicación:", err);
                alert("Error de conexión con el servidor Go. Verifique que el servidor esté corriendo.");
                
                if (submitBtn) {
                    submitBtn.disabled = false;
                    submitBtn.textContent = 'Iniciar sesión';
                }
            }
        });
    } else {
        console.error("[JS-AUTH] ❌ No se encontró el elemento 'login-form' en el DOM.");
    }
});