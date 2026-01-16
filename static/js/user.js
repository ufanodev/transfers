/**
 * @file user.js
 * @description Gestión de la tabla de usuarios ultra-compacta.
 * Maneja carga de datos, paginación local y redirección al CRUD.
 */

let allUsers = [];

/**
 * Inicialización al cargar el DOM
 */
document.addEventListener('DOMContentLoaded', () => {
    console.log("%c[SYSTEM] Módulo User iniciado.", "color: #10b981; font-weight: bold;");
    loadUsers();
});

/**
 * Carga la lista de usuarios desde la API de Go
 */
async function loadUsers() {
    const paginationInfo = document.getElementById('pagination-info');
    try {
        console.log("[API] Solicitando /api/v1/users...");
        const response = await fetch('/api/v1/users');
        
        if (response.status === 401) {
            window.location.href = "/";
            return;
        }

        // Verificar si el servidor devolvió HTML por error (redirección fallida)
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("text/html")) {
            console.error("[FATAL] El servidor devolvió HTML en lugar de JSON.");
            if (paginationInfo) paginationInfo.textContent = "ERROR: BUCLE DE REDIRECCIÓN";
            return;
        }

        allUsers = await response.json();
        console.log("Usuarios cargados correctamente:", allUsers);
        renderTable();
    } catch (error) {
        console.error("Error al cargar usuarios:", error);
        if (paginationInfo) paginationInfo.textContent = "ERROR DE CONEXIÓN";
    }
}

/**
 * Renderiza la tabla con normalización de campos (PascalCase a camelCase)
 */
function renderTable() {
    const pageSizeSelect = document.getElementById('pageSize');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 10;
    const tbody = document.getElementById('user-table-body');
    const paginationInfo = document.getElementById('pagination-info');
    
    // Filtrar por tamaño de página (0 es ALL)
    const displayUsers = (pageSize === 0) ? allUsers : allUsers.slice(0, pageSize);
    
    if (paginationInfo) {
        paginationInfo.textContent = `REGISTROS: ${displayUsers.length} / ${allUsers.length}`;
    }

    if (!allUsers || allUsers.length === 0) {
        tbody.innerHTML = `<tr><td colspan="4" class="p-10 text-center text-gray-400 uppercase text-[9px] font-black italic">No hay datos encontrados</td></tr>`;
        return;
    }

    tbody.innerHTML = displayUsers.map(u => {
        // --- NORMALIZACIÓN CRÍTICA ---
        // Esto permite que el JS lea los campos sin importar si Go los mandó como "ID" o "id"
        const userId = u.ID || u.id;
        const username = u.Username || u.username || "Sin Nombre";
        const email = u.Email || u.email || "Sin Email";
        const role = u.Role || u.role || "client";
        const isActive = (u.IsActive !== undefined) ? u.IsActive : u.is_active;
        
        return `
        <tr class="hover:bg-slate-50 transition-colors border-b border-gray-100">
            <td class="px-3 py-1.5 overflow-hidden">
                <div class="flex flex-col leading-none">
                    <span class="font-black text-black uppercase text-[10px] truncate" title="${username}">
                        ${username}
                    </span>
                    <span class="text-[8px] text-gray-400 font-bold lowercase truncate" title="${email}">
                        ${email}
                    </span>
                </div>
            </td>
            <td class="px-3 py-1.5 text-center">
                <span class="px-1.5 py-0.5 rounded-[4px] text-[8px] font-black uppercase shadow-sm ${getRoleStyle(role)}">
                    ${role}
                </span>
            </td>
            <td class="px-3 py-1.5 text-center">
                <span class="text-[8px] font-black italic ${isActive ? 'text-emerald-600' : 'text-red-500'}">
                    ${isActive ? '● ACTIVO' : '○ INACTIVO'}
                </span>
            </td>
            <td class="px-3 py-1.5 text-right">
                <div class="flex justify-end gap-1.5">
                    <button onclick="window.location.href='/dashboard/users/manage?id=${userId}'" 
                        class="p-1 bg-black text-[#10b981] rounded hover:bg-[#10b981] hover:text-black transition-all shadow-sm active:scale-90"
                        title="Editar">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                            <path d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                    </button>
                    <button onclick="quickDelete(${userId})" 
                        class="p-1 bg-red-50 text-red-600 rounded hover:bg-red-600 hover:text-white transition-all shadow-sm active:scale-90"
                        title="Borrar">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                            <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Estilos para badges de roles
 */
function getRoleStyle(role) {
    // Normalizar a minúsculas para la comparación
    const r = role.toLowerCase();
    const s = {
        'admin': 'bg-black text-[#10b981] border border-[#10b981]/30',
        'company': 'bg-slate-700 text-white',
        'driver': 'bg-[#10b981] text-black',
        'client': 'bg-slate-100 text-slate-500'
    };
    return s[r] || 'bg-gray-100 text-gray-400';
}

/**
 * Eliminación rápida
 */
async function quickDelete(id) {
    if (!id || !confirm("¿Eliminar este registro permanentemente?")) return;

    try {
        const response = await fetch(`/api/v1/users/${id}`, { method: 'DELETE' });
        if (response.ok) {
            console.log(`Usuario ${id} eliminado.`);
            loadUsers(); 
        } else {
            const err = await response.json();
            alert("Error: " + (err.error || "No autorizado"));
        }
    } catch (e) {
        console.error("Error en el borrado:", e);
    }
}

function exportData(type) {
    alert(`Generando Reporte ${type.toUpperCase()}...\nRegistros: ${allUsers.length}`);
}