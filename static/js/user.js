/**
 * @file user.js
 * @description Gestión de la lista de usuarios y tabla dinámica.
 */

let allUsers = [];

async function loadUsers() {
    const paginationInfo = document.getElementById('pagination-info');
    try {
        const res = await fetch('/api/v1/users');
        if (res.status === 401) {
            window.location.href = "/";
            return;
        }
        allUsers = await res.json();
        renderTable();
    } catch (e) {
        console.error("Error cargando usuarios:", e);
        paginationInfo.textContent = "ERROR AL CARGAR";
    }
}

function renderTable() {
    const pageSize = parseInt(document.getElementById('pageSize').value);
    const tbody = document.getElementById('user-table-body');
    const paginationInfo = document.getElementById('pagination-info');
    
    // Filtro de paginación
    const displayUsers = (pageSize === 0) ? allUsers : allUsers.slice(0, pageSize);
    
    paginationInfo.textContent = `Mostrando ${displayUsers.length} de ${allUsers.length} registros`;

    if (allUsers.length === 0) {
        tbody.innerHTML = `<tr><td colspan="5" class="p-10 text-center text-gray-400 font-bold uppercase tracking-widest text-xs">No se encontraron registros</td></tr>`;
        return;
    }

    tbody.innerHTML = displayUsers.map(u => `
        <tr class="hover:bg-brand-light transition-all group">
            <td class="p-6 font-bold text-slate-800">${u.username}</td>
            <td class="p-6 text-gray-500 text-xs">${u.email}</td>
            <td class="p-6 text-center">
                <span class="px-3 py-1 rounded text-[9px] font-black uppercase ${getRoleClass(u.role)}">
                    ${u.role}
                </span>
            </td>
            <td class="p-6 text-center">
                <span class="text-[10px] font-black ${u.is_active ? 'text-brand-primary' : 'text-red-500'}">
                    ${u.is_active ? '● ACTIVO' : '○ INACTIVO'}
                </span>
            </td>
            <td class="p-6 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/dashboard/users/manage?id=${u.ID}'" 
                        class="p-2 bg-brand-accent text-brand-dark rounded-lg hover:bg-brand-primary hover:text-white transition-all shadow-sm" title="Editar">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                    </button>
                    <button onclick="quickDelete(${u.ID})" 
                        class="p-2 bg-red-50 text-red-500 rounded-lg hover:bg-red-500 hover:text-white transition-all shadow-sm" title="Borrar">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                    </button>
                </div>
            </td>
        </tr>
    `).join('');
}

function getRoleClass(role) {
    const s = { 
        admin: 'bg-black text-white', 
        company: 'bg-brand-dark text-white', 
        driver: 'bg-brand-accent text-brand-dark', 
        client: 'bg-gray-100 text-gray-500' 
    };
    return s[role] || 'bg-gray-100';
}

async function quickDelete(id) {
    if(!confirm("¿Estás seguro de que deseas eliminar este usuario permanentemente?")) return;
    try {
        const res = await fetch(`/api/v1/users/${id}`, { method: 'DELETE' });
        if (res.ok) loadUsers();
        else alert("Error al eliminar");
    } catch (e) { console.error(e); }
}

function exportData(type) {
    alert(`Generando reporte ${type.toUpperCase()}...`);
}

document.addEventListener('DOMContentLoaded', loadUsers);