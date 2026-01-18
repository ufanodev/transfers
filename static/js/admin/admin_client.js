/**
 * @file client.js
 * @description Gestión de Clientes con ordenación dinámica y logs.
 */

let allClients = [];
let sortDirection = true; // true = ASC, false = DESC

document.addEventListener('DOMContentLoaded', () => {
    console.log(" [SYSTEM] Módulo Clients iniciado.");
    loadClients();
});

async function loadClients() {
    console.log(" [API] Cargando clientes...");
    try {
        const response = await fetch('/api/v1/clients');
        if (response.status === 401) { window.location.href = "/"; return; }
        
        allClients = await response.json();
        console.log(` [OK] ${allClients.length} clientes recibidos.`);
        renderTable();
    } catch (error) {
        console.error(" [ERROR] Fallo en carga:", error);
    }
}

function sortData(field) {
    sortDirection = !sortDirection;
    console.log(` [SORT] Ordenando por ${field} | Dirección: ${sortDirection ? 'ASC' : 'DESC'}`);
    
    allClients.sort((a, b) => {
        let valA = a[field];
        let valB = b[field];
        
        if (typeof valA === 'string') {
            return sortDirection 
                ? valA.localeCompare(valB) 
                : valB.localeCompare(valA);
        }
        return sortDirection ? valA - valB : valB - valA;
    });
    renderTable();
}

function renderTable() {
    const size = parseInt(document.getElementById('pageSize').value);
    const tbody = document.getElementById('client-table-body');
    const info = document.getElementById('pagination-info');
    
    const displayData = (size === 0) ? allClients : allClients.slice(0, size);
    if(info) info.textContent = `CLIENTES: ${displayData.length} / ${allClients.length}`;

    tbody.innerHTML = displayData.map(c => `
        <tr class="hover:bg-emerald-50/50 transition-colors border-b border-gray-50">
            <td class="px-3 py-1.5 font-bold text-gray-400">#${c.id}</td>
            <td class="px-3 py-1.5">
                <div class="flex flex-col">
                    <span class="font-black text-black uppercase text-[10px]">${c.full_name}</span>
                    <span class="text-[8px] text-gray-400 font-bold italic">User ID: ${c.user_id}</span>
                </div>
            </td>
            <td class="px-3 py-1.5">
                <div class="flex flex-col leading-tight">
                    <span class="text-black font-bold">${c.email}</span>
                    <span class="text-[8px] text-[#10b981] font-black uppercase">${c.phone}</span>
                </div>
            </td>
            <td class="px-3 py-1.5 text-center">
                <span class="bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded text-[8px] font-black uppercase">
                    ${c.preferences ? 'Personalizado' : 'Estándar'}
                </span>
            </td>
            <td class="px-3 py-1.5 text-right">
                <div class="flex justify-end gap-1.5">
                    <button onclick="window.location.href='/dashboard/clients/manage?id=${c.id}'" 
                        class="p-1 bg-black text-[#10b981] rounded hover:bg-[#10b981] hover:text-black transition-all shadow-sm">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24"><path d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>
                    </button>
                </div>
            </td>
        </tr>
    `).join('');
}