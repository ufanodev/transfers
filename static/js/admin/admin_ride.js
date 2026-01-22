/**
 * @file admin_ride.js
 * @description Gestión del monitor de viajes activos y despacho.
 */

let allRides = [];

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[RIDE-MONITOR] Sincronizando operaciones...", "color: #10b981; font-weight: bold;");
    await loadRides();
});

async function loadRides() {
    const tbody = document.getElementById('ride-table-body');
    const rideCount = document.getElementById('ride-count');
    
    try {
        const res = await fetch('/api/v1/rides');
        if (!res.ok) throw new Error("Error en la red");
        allRides = await res.json();

        rideCount.textContent = `${allRides.length} SERVICIOS EN SISTEMA`;
        renderTable(allRides);
    } catch (err) {
        tbody.innerHTML = `<tr><td colspan="4" class="p-10 text-center text-red-500 font-black uppercase">Fallo al conectar con el servidor logístico</td></tr>`;
    }
}

function renderTable(data) {
    const tbody = document.getElementById('ride-table-body');
    
    const statusClasses = {
        'scheduled': 'bg-blue-50 text-blue-600 border-blue-100',
        'in_progress': 'bg-emerald-50 text-emerald-600 border-emerald-100',
        'completed': 'bg-slate-100 text-slate-500 border-slate-200',
        'cancelled': 'bg-red-50 text-red-600 border-red-100'
    };

    tbody.innerHTML = data.map(r => {
        // Determinamos el estado basado en el campo is_finished si no existe status string
        let status = r.status || (r.is_finished ? 'completed' : 'in_progress');
        
        return `
        <tr class="hover:bg-slate-50 transition-all border-b border-slate-100 group">
            <td class="px-6 py-4">
                <div class="flex flex-col">
                    <span class="font-black text-slate-800 uppercase text-[12px] tracking-tighter">
                        ${r.voucher_number || '#S/REF'}
                    </span>
                    <span class="text-[10px] text-slate-400 font-bold uppercase truncate max-w-[200px]">
                        ${r.origin_address || 'Dirección no definida'}
                    </span>
                </div>
            </td>
            <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center text-slate-400">
                        <i class="fa-solid fa-user-tie text-xs"></i>
                    </div>
                    <div class="flex flex-col">
                        <span class="font-black text-slate-700 text-[11px] uppercase italic">
                            ${r.client_name || 'Pasajero General'}
                        </span>
                        <span class="text-[9px] text-admin-accent font-bold">
                            Cond: ${r.driver_id ? 'Asignado' : 'PENDIENTE'}
                        </span>
                    </div>
                </div>
            </td>
            <td class="px-6 py-4 text-center">
                <span class="px-3 py-1 rounded-full text-[9px] font-black uppercase border ${statusClasses[status] || 'bg-slate-50'}">
                    ${status.replace('_', ' ')}
                </span>
            </td>
            <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/admin/rides/manage?id=${r.id || r.ID}'" 
                            class="w-8 h-8 flex items-center justify-center bg-white text-slate-400 rounded-xl hover:bg-admin-accent hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-sliders text-[10px]"></i>
                    </button>
                    <button onclick="openDeleteModal(${r.id || r.ID}, '${r.voucher_number || r.id}')" 
                            class="w-8 h-8 flex items-center justify-center bg-white text-red-300 rounded-xl hover:bg-red-600 hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-trash-can text-[10px]"></i>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

function filterRides() {
    const term = document.getElementById('rideSearch').value.toLowerCase();
    const status = document.getElementById('statusFilter').value;
    
    const filtered = allRides.filter(r => {
        const matchesTerm = (r.voucher_number && r.voucher_number.toLowerCase().includes(term)) || 
                          (r.client_name && r.client_name.toLowerCase().includes(term)) ||
                          (r.origin_address && r.origin_address.toLowerCase().includes(term));
        
        const currentStatus = r.status || (r.is_finished ? 'completed' : 'in_progress');
        const matchesStatus = status === "" || currentStatus === status;
        
        return matchesTerm && matchesStatus;
    });
    
    renderTable(filtered);
}

// Lógica del Modal de Borrado
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
        const res = await fetch(`/api/v1/rides/${id}`, { method: 'DELETE' });
        if (res.ok) {
            closeDeleteModal();
            await loadRides();
        }
    } catch (e) { alert("Error al eliminar"); }
}

function exportData(format) {
    alert(`Exportando monitor a ${format.toUpperCase()}...`);
}