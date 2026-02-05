/**
 * @file admin_ride.js
 * @description Torre de Control TakeUs - Gestión del flujo operativo de viajes.
 */

let allRides = [];

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[TAKEUS-OPS] Monitor de Operaciones Iniciado", "color: #3b82f6; font-weight: bold;");
    await loadRides();
});

/**
 * Carga todos los viajes desde el backend
 */
async function loadRides() {
    const tbody = document.getElementById('ride-table-body');
    const rideCount = document.getElementById('ride-count');
    
    try {
        const res = await fetch('/api/v1/rides');
        if (!res.ok) throw new Error("Error en la comunicación con el servidor");
        allRides = await res.json();

        rideCount.textContent = `${allRides.length} SERVICIOS EN MONITOR`;
        renderTable(allRides);
    } catch (err) {
        console.error(err);
        tbody.innerHTML = `<tr><td colspan="7" class="p-10 text-center text-red-500 font-black uppercase">Error de conexión con la central logística</td></tr>`;
    }
}

/**
 * Renderiza la tabla con los 10 campos clave y orden de prioridad
 */
function renderTable(data) {
    const tbody = document.getElementById('ride-table-body');
    
    // 5. STATUS - Configuración visual según el flujo real
    const statusConfig = {
        'scheduled':   { class: 'bg-slate-100 text-slate-600 border-slate-200', label: 'Programado' },
        'assigned':    { class: 'bg-blue-50 text-blue-600 border-blue-200', label: 'Asignado' },
        'on_route':    { class: 'bg-indigo-50 text-indigo-600 border-indigo-200', label: 'En Camino' },
        'arrived':     { class: 'bg-amber-50 text-amber-600 border-amber-200', label: 'En Puerta' },
        'in_progress': { class: 'bg-emerald-50 text-emerald-600 border-emerald-200', label: 'En Curso' },
        'completed':   { class: 'bg-gray-800 text-white border-black', label: 'Completado' },
        'cancelled':   { class: 'bg-red-50 text-red-600 border-red-200', label: 'Cancelado' }
    };

    // Orden de prioridad: Lo que está ocurriendo AHORA arriba
    const priority = { 'in_progress': 1, 'arrived': 2, 'on_route': 3, 'assigned': 4, 'scheduled': 5, 'completed': 6, 'cancelled': 7 };
    
    data.sort((a, b) => {
        const pA = priority[a.status] || 99;
        const pB = priority[b.status] || 99;
        return pA - pB;
    });

    tbody.innerHTML = data.map(r => {
        const conf = statusConfig[r.status] || statusConfig['scheduled'];
        
        return `
        <tr class="hover:bg-blue-50/30 transition-all border-b border-slate-100 group">
            <td class="px-4 py-3">
                <div class="flex flex-col">
                    <span class="font-black text-slate-800 text-[12px]">#R-${r.id}</span>
                    <span class="text-[9px] text-slate-400 font-bold uppercase">BK-${r.booking_id}</span>
                    <span class="text-[9px] text-admin-accent font-mono truncate w-20" title="${r.voucher_number}">${r.voucher_number || 'S/REF'}</span>
                </div>
            </td>
            
            <td class="px-4 py-3">
                <div class="flex flex-col">
                    <span class="font-bold text-slate-700 uppercase text-[11px]">${r.client_name || 'Pasajero'}</span>
                    <div class="flex items-center gap-2 mt-1">
                        <span class="text-[10px] ${r.driver_id ? 'text-admin-accent' : 'text-red-500'} font-black italic">
                            <i class="fa-solid fa-user-tie mr-1"></i> 
                            ${r.driver_id ? 'ID: ' + r.driver_id : 'PENDIENTE'}
                        </span>
                        <span class="text-[9px] text-slate-400 font-bold">
                            <i class="fa-solid fa-car"></i> ${r.vehicle_id ? 'ID: ' + r.vehicle_id : '--'}
                        </span>
                    </div>
                </div>
            </td>

            <td class="px-4 py-3">
                <div class="flex flex-col max-w-[200px]">
                    <span class="text-[10px] font-bold truncate text-slate-600" title="${r.pickup_address}">
                        <i class="fa-solid fa-circle-dot text-emerald-500 mr-1"></i> ${r.pickup_address}
                    </span>
                    <span class="text-[10px] truncate text-slate-400 italic" title="${r.dropoff_address}">
                        <i class="fa-solid fa-location-dot text-red-400 mr-1"></i> ${r.dropoff_address}
                    </span>
                </div>
            </td>

            <td class="px-4 py-3 text-center">
                <div class="flex flex-col items-center">
                    <span class="text-[10px] font-black text-slate-700">${r.start_time_real ? formatTime(r.start_time_real) : '--:--'}</span>
                    <span class="text-[8px] uppercase text-slate-400">${r.is_finished ? 'Finalizado' : 'En Cola'}</span>
                </div>
            </td>

            <td class="px-4 py-3 text-center">
                <span class="px-2 py-1 rounded-md text-[9px] font-black uppercase border ${conf.class}">
                    ${conf.label}
                </span>
            </td>

            <td class="px-4 py-3 text-right">
                <span class="font-black text-slate-800 text-[12px]">${r.total_amount?.toFixed(2) || '0.00'}€</span>
            </td>

            <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-1">
                    <button onclick="window.location.href='/admin/rides/manage?id=${r.id}'" 
                            class="w-8 h-8 flex items-center justify-center bg-white border border-slate-200 text-slate-600 rounded-lg hover:bg-admin-accent hover:text-white transition-all shadow-sm" title="Ver y Gestionar">
                        <i class="fa-solid fa-eye text-[10px]"></i>
                    </button>
                    <button onclick="openDeleteModal(${r.id}, '${r.voucher_number || r.id}')" 
                            class="w-8 h-8 flex items-center justify-center bg-white border border-slate-200 text-red-300 rounded-lg hover:bg-red-600 hover:text-white transition-all shadow-sm">
                        <i class="fa-solid fa-trash-can text-[10px]"></i>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Utilidad para formatear fechas ISO a HH:mm
 */
function formatTime(dateStr) {
    if (!dateStr) return '--:--';
    const date = new Date(dateStr);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

/**
 * Filtro de búsqueda en tiempo real
 */
function filterRides() {
    const term = document.getElementById('rideSearch').value.toLowerCase();
    const status = document.getElementById('statusFilter').value;
    
    const filtered = allRides.filter(r => {
        const matchesTerm = (r.voucher_number?.toLowerCase().includes(term)) || 
                          (r.client_name?.toLowerCase().includes(term)) ||
                          (r.pickup_address?.toLowerCase().includes(term));
        
        const matchesStatus = status === "" || r.status === status;
        return matchesTerm && matchesStatus;
    });
    
    renderTable(filtered);
}

/**
 * Modales de borrado
 */
function openDeleteModal(id, name) {
    document.getElementById('delete-id').value = id;
    document.getElementById('confirm-item-name').textContent = name;
    document.getElementById('confirm-key').value = '';
    document.getElementById('delete-modal').classList.remove('hidden');
}

function closeDeleteModal() {
    document.getElementById('delete-modal').classList.add('hidden');
}

async function confirmFinalDelete() {
    const id = document.getElementById('delete-id').value;
    const key = document.getElementById('confirm-key').value;

    if (key !== 'ELIMINAR') {
        alert("Escriba la palabra correctamente");
        return;
    }

    try {
        const res = await fetch(`/api/v1/rides/${id}`, { 
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
        });
        if (res.ok) {
            closeDeleteModal();
            await loadRides();
        }
    } catch (e) { alert("Error al eliminar"); }
}

function exportData(format) {
    alert(`Generando reporte de servicios en ${format.toUpperCase()}...`);
}