/**
 * @file admin_user.js
 * @description Gestión de la tabla de usuarios con búsqueda dinámica y auditoría.
 */

let allUsers = [];      // Datos originales de la API
let filteredUsers = []; // Datos tras aplicar filtros
let sortDirection = true; 

document.addEventListener('DOMContentLoaded', () => {
    const adminStyle = "color: #3b82f6; font-weight: bold; background: #1e293b; padding: 3px 8px; border-radius: 5px;";
    console.log("%c[SYSTEM] Módulo Usuarios Admin iniciado.", adminStyle);
    loadUsers();
});

/**
 * Carga la lista de usuarios desde la API
 */
async function loadUsers() {
    const info = document.getElementById('pagination-info');
    
    try {
        const response = await fetch('/api/v1/users');
        
        if (response.status === 401) {
            window.location.href = "/";
            return;
        }

        if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);

        allUsers = await response.json();
        filteredUsers = [...allUsers]; // Inicialmente son iguales
        
        console.log(`[DATA] ${allUsers.length} usuarios cargados.`);
        renderTable();

    } catch (error) {
        console.error("[FETCH ERROR]", error);
        if (info) {
            info.textContent = "ERROR DE CONEXIÓN CON API";
            info.classList.replace('text-admin-accent', 'text-red-500');
        }
    }
}

/**
 * Filtra los usuarios en tiempo real basándose en el input de búsqueda
 */
function filterUsers() {
    const searchTerm = document.getElementById('userSearch').value.toLowerCase();
    
    filteredUsers = allUsers.filter(u => {
        return (
            u.username.toLowerCase().includes(searchTerm) ||
            u.email.toLowerCase().includes(searchTerm) ||
            u.role.toLowerCase().includes(searchTerm)
        );
    });

    renderTable();
}

/**
 * Ordena los datos filtrados
 */
function sortBy(field) {
    sortDirection = !sortDirection;
    
    filteredUsers.sort((a, b) => {
        let valA = (a[field] || "").toString().toLowerCase();
        let valB = (b[field] || "").toString().toLowerCase();
        if (valA < valB) return sortDirection ? -1 : 1;
        if (valA > valB) return sortDirection ? 1 : -1;
        return 0;
    });
    
    renderTable();
}

/**
 * Pinta las filas en el tbody
 */
function renderTable() {
    const tbody = document.getElementById('user-table-body');
    const info = document.getElementById('pagination-info');
    const pageSize = parseInt(document.getElementById('pageSize').value) || 0;
    
    // Aplicar paginación visual
    const displayData = (pageSize === 0) ? filteredUsers : filteredUsers.slice(0, pageSize);
    
    if (info) {
        info.textContent = `MOSTRANDO: ${displayData.length} DE ${filteredUsers.length} RESULTADOS`;
    }

    if (filteredUsers.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="4" class="p-12 text-center">
                    <i class="fa-solid fa-user-slash text-4xl text-slate-200 mb-4 block"></i>
                    <p class="text-slate-400 font-bold uppercase text-xs tracking-widest">No se encontraron usuarios</p>
                </td>
            </tr>`;
        return;
    }

    tbody.innerHTML = displayData.map(u => {
        const userId = u.id || u.ID;
        return `
        <tr class="hover:bg-slate-50 transition-all border-b border-slate-100 group">
            <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-400 text-[10px] font-bold border border-slate-200 group-hover:bg-admin-accent group-hover:text-white transition-colors">
                        ${u.username.substring(0,2).toUpperCase()}
                    </div>
                    <div class="flex flex-col">
                        <span class="font-bold text-slate-700 text-sm">${u.username}</span>
                        <span class="text-[11px] text-slate-400 font-medium">${u.email}</span>
                    </div>
                </div>
            </td>
            <td class="px-6 py-4 text-center">
                <span class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-tighter ${getRoleStyle(u.role)}">
                    ${u.role}
                </span>
            </td>
            <td class="px-6 py-4 text-center">
                <div class="flex items-center justify-center gap-2">
                    <span class="w-2 h-2 rounded-full ${u.is_active ? 'bg-emerald-500 animate-pulse' : 'bg-red-500'}"></span>
                    <span class="text-[11px] font-bold ${u.is_active ? 'text-emerald-600' : 'text-red-500'} uppercase">
                        ${u.is_active ? 'Activo' : 'Bloqueado'}
                    </span>
                </div>
            </td>
            <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/admin/users/manage?id=${userId}'" 
                        class="w-8 h-8 flex items-center justify-center bg-slate-100 text-slate-600 rounded-lg hover:bg-admin-accent hover:text-white transition-all shadow-sm"
                        title="Editar Usuario">
                        <i class="fa-solid fa-pen-to-square text-xs"></i>
                    </button>
                    <button onclick="deleteUser(${userId})" 
                        class="w-8 h-8 flex items-center justify-center bg-red-50 text-red-500 rounded-lg hover:bg-red-500 hover:text-white transition-all shadow-sm"
                        title="Eliminar">
                        <i class="fa-solid fa-trash-can text-xs"></i>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Estilos específicos por Rol (Actualizados a Slate/Blue)
 */
function getRoleStyle(role) {
    const styles = {
        'admin':   'bg-slate-900 text-white border border-slate-700',
        'driver':  'bg-emerald-100 text-emerald-700',
        'company': 'bg-blue-100 text-blue-700',
        'client':  'bg-slate-100 text-slate-600'
    };
    return styles[role] || 'bg-gray-100 text-gray-500';
}

/**
 * Placeholder para exportación
 */
function exportData(type) {
    alert(`Preparando exportación a ${type.toUpperCase()} de ${filteredUsers.length} registros...`);
}

/**
 * Placeholder para borrado
 */
async function deleteUser(id) {
    if(confirm('¿Está seguro de eliminar este usuario? Esta acción no se puede deshacer.')) {
        console.log("Eliminando usuario:", id);
        // Aquí iría el fetch DELETE /api/v1/users/:id
    }
}