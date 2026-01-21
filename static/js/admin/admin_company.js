/**
 * @file admin_company.js
 * @description Gestión de listado y eliminación con triple validación para Empresas.
 * Ubicación física: /static/js/admin/admin_company.js
 */

let allCompanies = [];
let currentSort = { field: 'name', asc: true };

document.addEventListener('DOMContentLoaded', async () => {
    console.log("%c[LOAD] admin_company.js inicializado", "color: #3b82f6; font-weight: bold;");
    await loadCompanies();

    const searchInput = document.getElementById('companySearch');
    if (searchInput) {
        searchInput.addEventListener('input', () => renderTable());
    }
});

/**
 * Carga los datos desde la API
 */
async function loadCompanies() {
    const tbody = document.getElementById('company-table-body');
    try {
        const response = await fetch('/api/v1/companies');
        if (response.status === 401) { window.location.href = "/"; return; }
        if (!response.ok) throw new Error("Fallo al obtener empresas");

        allCompanies = await response.json();
        renderTable();

    } catch (error) {
        console.error("[ERROR]", error);
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="4" class="p-10 text-center text-red-500 font-black uppercase text-[10px]">Error de conexión</td></tr>`;
        }
    }
}

/**
 * Renderiza la tabla con filtros y acciones
 */
function renderTable() {
    const tbody = document.getElementById('company-table-body');
    const searchInput = document.getElementById('companySearch');
    const pageSizeSelect = document.getElementById('pageSize');
    const paginationInfo = document.getElementById('pagination-info');

    if (!tbody) return;

    let filtered = [...allCompanies];

    // 1. Filtro
    const term = searchInput ? searchInput.value.toLowerCase().trim() : "";
    if (term) {
        filtered = filtered.filter(c => 
            (c.name && c.name.toLowerCase().includes(term)) ||
            (c.tax_id && c.tax_id.toLowerCase().includes(term)) ||
            (c.address && c.address.toLowerCase().includes(term))
        );
    }

    // 2. Orden
    filtered.sort((a, b) => {
        let valA = a[currentSort.field] ? a[currentSort.field].toString().toLowerCase() : '';
        let valB = b[currentSort.field] ? b[currentSort.field].toString().toLowerCase() : '';
        return currentSort.asc ? valA.localeCompare(valB) : valB.localeCompare(valA);
    });

    // 3. Paginación
    const size = pageSizeSelect ? parseInt(pageSizeSelect.value) : 10;
    const displayData = (size === 0) ? filtered : filtered.slice(0, size);

    if (paginationInfo) {
        paginationInfo.textContent = `EMPRESAS REGISTRADAS: ${filtered.length} (Mostrando ${displayData.length})`;
    }

    if (displayData.length === 0) {
        tbody.innerHTML = `<tr><td colspan="4" class="p-10 text-center text-slate-400 font-black uppercase italic text-[10px]">Sin registros</td></tr>`;
        return;
    }

    // 4. Renderizado con botones de Editar y Borrar
    tbody.innerHTML = displayData.map(c => `
        <tr class="hover:bg-slate-50 transition-colors border-b border-slate-100 group">
            <td class="px-6 py-4">
                <div class="flex flex-col">
                    <span class="font-black text-slate-700 uppercase tracking-tighter">${c.name}</span>
                    <span class="text-[9px] text-slate-400 font-bold italic">ID: ${c.id || c.ID} | CIF: ${c.tax_id || 'N/A'}</span>
                </div>
            </td>
            <td class="px-6 py-4">
                <div class="flex flex-col leading-tight">
                    <span class="text-slate-600 font-bold text-[10px] uppercase truncate max-w-[250px]">${c.address || 'N/A'}</span>
                    <span class="text-[9px] text-admin-accent font-medium">${c.website || ''}</span>
                </div>
            </td>
            <td class="px-6 py-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[8px] font-black uppercase bg-emerald-50 text-emerald-600 border border-emerald-100 shadow-sm">Activo</span>
            </td>
            <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/admin/companies/manage?id=${c.id || c.ID}'" 
                        class="w-8 h-8 flex items-center justify-center bg-white text-slate-400 rounded-lg hover:bg-admin-accent hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-pen-to-square text-[11px]"></i>
                    </button>
                    <button onclick="openDeleteModal(${c.id || c.ID}, '${c.name}')" 
                        class="w-8 h-8 flex items-center justify-center bg-white text-red-300 rounded-lg hover:bg-red-600 hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-trash-can text-[11px]"></i>
                    </button>
                </div>
            </td>
        </tr>
    `).join('');
}

// --- LÓGICA DEL MODAL DE BORRADO ---

window.openDeleteModal = (id, name) => {
    document.getElementById('delete-id').value = id;
    document.getElementById('confirm-item-name').textContent = name;
    document.getElementById('confirm-item-id').textContent = id;
    document.getElementById('delete-modal').classList.remove('hidden');
};

window.closeDeleteModal = () => {
    document.getElementById('delete-modal').classList.add('hidden');
    document.getElementById('confirm-input-name').value = '';
    document.getElementById('confirm-input-id').value = '';
    document.getElementById('confirm-key').value = '';
};

window.confirmFinalDelete = async () => {
    const id = document.getElementById('delete-id').value;
    const nameInput = document.getElementById('confirm-input-name').value.trim().toUpperCase();
    const idInput = document.getElementById('confirm-input-id').value.trim();
    const keyInput = document.getElementById('confirm-key').value.trim().toUpperCase();
    
    const realName = document.getElementById('confirm-item-name').textContent.trim().toUpperCase();

    // Triple Validación
    if (nameInput === realName && idInput === id && keyInput === 'ELIMINAR') {
        try {
            const response = await fetch(`/api/v1/companies/${id}`, { method: 'DELETE' });
            if (response.ok) {
                closeDeleteModal();
                await loadCompanies(); // Recarga la tabla tras borrar
            } else {
                alert("Error: El servidor no permitió borrar el registro.");
            }
        } catch (e) {
            console.error(e);
            alert("Error de red.");
        }
    } else {
        alert("Los datos de confirmación no coinciden. Verifique Nombre, ID y la clave 'ELIMINAR'.");
    }
};

// --- UTILIDADES ---
window.filterCompanies = () => renderTable();
window.sortBy = (field) => {
    if (currentSort.field === field) { currentSort.asc = !currentSort.asc; }
    else { currentSort.field = field; currentSort.asc = true; }
    renderTable();
};
window.exportData = (type) => alert(`Iniciando exportación a ${type.toUpperCase()}...`);