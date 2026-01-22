/**
 * @file admin_booking_event_crud.js
 * @description Envío de notas manuales y rectificaciones a la auditoría.
 */

document.addEventListener('DOMContentLoaded', () => {
    const eventForm = document.getElementById('event-form');

    eventForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        // Obtener datos del formulario
        const bookingId = document.getElementById('booking_id').value;
        const eventType = document.getElementById('event_type').value;
        const description = document.getElementById('description').value.trim();

        // Validar que el ID de reserva sea válido
        if (!bookingId || bookingId <= 0) {
            alert("Por favor, ingrese un ID de reserva válido.");
            return;
        }

        const payload = {
            booking_id: parseInt(bookingId),
            event_type: eventType,
            description: description,
            created_by: "ADMIN_MANUAL" // Esto podría venir del token JWT en el futuro
        };

        try {
            const res = await fetch('/api/v1/events', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                // Éxito: Volver al log
                window.location.href = '/admin/events';
            } else {
                const errData = await res.json();
                alert("Error al registrar evento: " + (errData.error || "Fallo desconocido"));
            }
        } catch (error) {
            console.error("Error en la petición:", error);
            alert("No se pudo conectar con el servidor de auditoría.");
        }
    });
});