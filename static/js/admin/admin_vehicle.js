/**
 * @file admin_vehicle.js
 * @description Lógica de listado, filtrado y eliminación de vehículos.
 */

let allVehicles = [];
let allCompanies = [];

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[LOAD] admin_vehicle.js cargado correctamente", "color: #3b82f6; font-weight: bold;");
    
    // Carga paralela de empresas y vehículos
    await Promise.all([loadCompanies(), loadVehicles()]);
});

async function loadCompanies() {
    try {
        const res = await fetch('/api/v1/companies');
        allCompanies = await res.json();
    } catch (e) { console.error("Error cargando empresas:", e); }
}

async function loadVehicles() {
    const tbody = document.getElementById('vehicle-table-body');
    try {
        const response = await fetch('/api/v1/vehicles');
        if (response.status === 401) { window.location.href = "/"; return; }
        if (!response.ok) throw new Error("Fallo al obtener vehículos");

        allVehicles = await response.json();
        renderTable();

    } catch (error) {
        console.error("[ERROR]", error);
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="5" class="p-10 text-center text-red-500 font-black uppercase text-[10px]">Error de conexión con el servidor</td></tr>`;
        }
    }
}

function renderTable() {
    const tbody = document.getElementById('vehicle-table-body');
    const searchInput = document.getElementById('vehicleSearch');
    const paginationInfo = document.getElementById('pagination-info');

    if (!tbody) return;

    let filtered = [...allVehicles];
    const term = searchInput.value.toLowerCase().trim();

    if (term) {
        filtered = filtered.filter(v => 
            v.plate_number.toLowerCase().includes(term) || 
            (v.make && v.make.toLowerCase().includes(term)) ||
            (v.model_name && v.model_name.toLowerCase().includes(term))
        );
    }

    paginationInfo.textContent = `FLOTA ACTIVA: ${filtered.length} VEHÍCULOS`;

    if (filtered.length === 0) {
        tbody.innerHTML = `<tr><td colspan="5" class="p-10 text-center text-slate-400 font-black uppercase italic text-[10px]">No hay vehículos registrados</td></tr>`;
        return;
    }

    tbody.innerHTML = filtered.map(v => {
        const company = allCompanies.find(c => (c.id || c.ID) === v.company_id);
        const companyName = company ? company.name : `Empresa ID: ${v.company_id}`;

        return `
        <tr class="hover:bg-slate-50 transition-colors border-b border-slate-100 group">
            <td class="px-6 py-4">
                <div class="flex flex-col">
                    <span class="font-black text-slate-700 uppercase tracking-tighter">${v.plate_number}</span>
                    <span class="text-[9px] text-slate-400 font-bold italic">${v.make} ${v.model_name}</span>
                </div>
            </td>
            <td class="px-6 py-4">
                <div class="flex flex-col leading-tight">
                    <span class="text-slate-600 font-bold text-[10px] uppercase">${v.vehicle_type}</span>
                    <span class="text-[9px] text-admin-accent font-medium">${v.category || 'N/A'} - Cap: ${v.capacity_pax} pax</span>
                </div>
            </td>
            <td class="px-6 py-4">
                <span class="text-slate-700 font-black uppercase text-[10px] border-b-2 border-slate-100">${companyName}</span>
            </td>
            <td class="px-6 py-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[8px] font-black uppercase ${v.is_active ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-red-50 text-red-600 border-red-100'} border shadow-sm">
                    ${v.is_active ? 'Operativo' : 'Inactivo'}
                </span>
            </td>
            <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/admin/vehicles/manage?id=${v.id || v.ID}'" 
                        class="w-8 h-8 flex items-center justify-center bg-white text-slate-400 rounded-lg hover:bg-admin-accent hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-pen-to-square text-[11px]"></i>
                    </button>
                    <button onclick="openDeleteModal(${v.id || v.ID}, '${v.plate_number}')" 
                        class="w-8 h-8 flex items-center justify-center bg-white text-red-300 rounded-lg hover:bg-red-600 hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-trash-can text-[11px]"></i>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

// --- MODAL LÓGICA DE BORRADO ---
window.openDeleteModal = (id, plate) => {
    document.getElementById('delete-id').value = id;
    document.getElementById('confirm-item-name').textContent = plate;
    document.getElementById('confirm-item-id').textContent = id;
    document.getElementById('delete-modal').classList.remove('hidden');
}

window.closeDeleteModal = () => document.getElementById('delete-modal').classList.add('hidden');

window.confirmFinalDelete = async () => {
    const id = document.getElementById('delete-id').value;
    const nameInput = document.getElementById('confirm-input-name').value.trim().toUpperCase();
    const idInput = document.getElementById('confirm-input-id').value.trim();
    const keyInput = document.getElementById('confirm-key').value.trim().toUpperCase();
    const realPlate = document.getElementById('confirm-item-name').textContent.trim().toUpperCase();

    if (nameInput === realPlate && idInput === id && keyInput === 'ELIMINAR') {
        const res = await fetch(`/api/v1/vehicles/${id}`, { method: 'DELETE' });
        if (res.ok) { 
            closeDeleteModal(); 
            loadVehicles(); 
        } else {
            alert("Error al eliminar del servidor");
        }
    } else { 
        alert("Los datos de validación no coinciden."); 
    }
}

window.filterVehicles = () => renderTable();
window.exportData = (type) => alert(`Iniciando exportación a ${type.toUpperCase()}...`);