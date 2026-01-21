let allDrivers = [];
let allCompanies = [];

document.addEventListener('DOMContentLoaded', async () => {
    // Cargamos empresas primero para poder cruzar los nombres en la tabla
    await Promise.all([loadCompanies(), loadDrivers()]);
});

async function loadCompanies() {
    const res = await fetch('/api/v1/companies');
    allCompanies = await res.json();
}

async function loadDrivers() {
    const tbody = document.getElementById('driver-table-body');
    try {
        const response = await fetch('/api/v1/drivers');
        if (response.status === 401) { window.location.href = "/"; return; }
        allDrivers = await response.json();
        renderTable();
    } catch (error) {
        tbody.innerHTML = `<tr><td colspan="4" class="p-10 text-center text-red-500 font-black uppercase text-[10px]">Error de conexión</td></tr>`;
    }
}

function renderTable() {
    const tbody = document.getElementById('driver-table-body');
    const term = document.getElementById('driverSearch').value.toLowerCase();
    
    let filtered = allDrivers.filter(d => 
        d.full_name.toLowerCase().includes(term) || 
        d.license_number.toLowerCase().includes(term)
    );

    document.getElementById('pagination-info').textContent = `FLOTA ACTIVA: ${filtered.length} CONDUCTORES`;

    tbody.innerHTML = filtered.map(d => {
        // Buscamos el nombre de la empresa
        const company = allCompanies.find(c => (c.id || c.ID) === d.company_id);
        const companyName = company ? company.name : `Empresa ID: ${d.company_id}`;
        const expiryDate = new Date(d.license_expiry_date).toLocaleDateString();

        return `
        <tr class="hover:bg-slate-50 transition-colors border-b border-slate-100 group">
            <td class="px-6 py-4">
                <div class="flex flex-col">
                    <span class="font-black text-slate-700 uppercase tracking-tighter">${d.full_name}</span>
                    <span class="text-[9px] text-slate-400 font-bold italic">Lic: ${d.license_number} | Expira: ${expiryDate}</span>
                </div>
            </td>
            <td class="px-6 py-4">
                <div class="flex flex-col leading-tight">
                    <span class="text-slate-600 font-bold text-[10px] uppercase">${companyName}</span>
                    <span class="text-[9px] text-admin-accent font-medium">${d.email} | ${d.phone}</span>
                </div>
            </td>
            <td class="px-6 py-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[8px] font-black uppercase ${d.is_available ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-red-50 text-red-600 border-red-100'} border shadow-sm">
                    ${d.is_available ? 'Disponible' : 'En Viaje / Off'}
                </span>
            </td>
            <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/admin/drivers/manage?id=${d.id || d.ID}'" class="w-8 h-8 flex items-center justify-center bg-white text-slate-400 rounded-lg hover:bg-admin-accent hover:text-white transition-all shadow-sm border border-slate-100"><i class="fa-solid fa-pen-to-square text-[11px]"></i></button>
                    <button onclick="openDeleteModal(${d.id || d.ID}, '${d.full_name}')" class="w-8 h-8 flex items-center justify-center bg-white text-red-300 rounded-lg hover:bg-red-600 hover:text-white transition-all shadow-sm border border-slate-100"><i class="fa-solid fa-trash-can text-[11px]"></i></button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

// --- LÓGICA MODAL ---
window.openDeleteModal = (id, name) => {
    document.getElementById('delete-id').value = id;
    document.getElementById('confirm-item-name').textContent = name;
    document.getElementById('confirm-item-id').textContent = id;
    document.getElementById('delete-modal').classList.remove('hidden');
}

window.closeDeleteModal = () => document.getElementById('delete-modal').classList.add('hidden');

window.confirmFinalDelete = async () => {
    const id = document.getElementById('delete-id').value;
    const nameInput = document.getElementById('confirm-input-name').value.trim().toUpperCase();
    const idInput = document.getElementById('confirm-input-id').value.trim();
    const keyInput = document.getElementById('confirm-key').value.trim().toUpperCase();
    const realName = document.getElementById('confirm-item-name').textContent.trim().toUpperCase();

    if (nameInput === realName && idInput === id && keyInput === 'ELIMINAR') {
        const res = await fetch(`/api/v1/drivers/${id}`, { method: 'DELETE' });
        if (res.ok) { closeDeleteModal(); loadDrivers(); }
    } else { alert("Confirmación incorrecta."); }
}

window.filterDrivers = () => renderTable();
window.exportData = (type) => alert("Exportando...");