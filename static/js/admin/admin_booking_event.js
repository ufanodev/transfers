/**
 * @file admin_booking_event.js
 * @description Listado de la caja negra (auditoría) del sistema.
 */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[AUDITORÍA] Sincronizando logs...", "color: #3b82f6; font-weight: bold;");
    await loadEvents();
});

async function loadEvents() {
    const tbody = document.getElementById('event-table-body');
    
    try {
        const res = await fetch('/api/v1/events');
        if (!res.ok) throw new Error("Error en red");
        const events = await res.json();

        if (events.length === 0) {
            tbody.innerHTML = `<tr><td colspan="5" class="p-10 text-center text-slate-400 italic">No hay eventos registrados aún.</td></tr>`;
            return;
        }

        renderTable(events);
    } catch (err) {
        tbody.innerHTML = `<tr><td colspan="5" class="p-10 text-center text-red-500 font-black">FALLO EN LA CONEXIÓN CON EL SERVIDOR DE AUDITORÍA</td></tr>`;
    }
}

function renderTable(data) {
    const tbody = document.getElementById('event-table-body');
    
    // Configuración de estilos por tipo de evento
    const badgeStyles = {
        'incident': 'text-red-600 bg-red-50 border-red-100',
        'rectification': 'text-amber-600 bg-amber-50 border-amber-100',
        'manual_note': 'text-blue-600 bg-blue-50 border-blue-100',
        'status_change': 'text-emerald-600 bg-emerald-50 border-emerald-100',
        'created': 'text-slate-600 bg-slate-50 border-slate-200'
    };

    tbody.innerHTML = data.map(ev => {
        const dateObj = new Date(ev.created_at);
        const dateStr = dateObj.toLocaleDateString();
        const timeStr = dateObj.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        const style = badgeStyles[ev.event_type] || badgeStyles['created'];

        return `
        <tr class="hover:bg-slate-50 transition-all border-b border-slate-100 group">
            <td class="px-5 py-3">
                <div class="flex flex-col">
                    <span class="font-black text-slate-700 uppercase">${dateStr}</span>
                    <span class="text-[9px] text-admin-accent font-bold">${timeStr}</span>
                </div>
            </td>
            <td class="px-5 py-3 font-black text-slate-400">
                <span class="bg-slate-100 px-2 py-1 rounded text-slate-600 border border-slate-200">#${ev.booking_id}</span>
            </td>
            <td class="px-5 py-3">
                <span class="px-2 py-0.5 rounded-full text-[8px] font-black uppercase border ${style}">
                    ${ev.event_type.replace('_', ' ')}
                </span>
            </td>
            <td class="px-5 py-3">
                <p class="text-slate-600 font-medium leading-tight max-w-md">${ev.description}</p>
            </td>
            <td class="px-5 py-3 text-right">
                <span class="text-[9px] font-black text-slate-400 uppercase italic tracking-widest">
                    ${ev.created_by || 'SISTEMA'}
                </span>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Función Placeholder para exportación (PDF/Excel)
 */
function exportData(type) {
    alert(`Generando reporte de auditoría en formato ${type.toUpperCase()}...`);
    // Aquí iría la lógica de jsPDF o SheetJS
}