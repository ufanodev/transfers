/**
 * @file admin_user.js
 * @description Gestión de la tabla de usuarios con auditoría de red y manejo de redirecciones.
 */

let allUsers = [];
let sortDirection = true; // true = ASC, false = DESC

/**
 * Inicialización al cargar el DOM
 */
document.addEventListener('DOMContentLoaded', () => {
    console.log("%c[SYSTEM] Módulo Admin_User iniciado.", "color: #10b981; font-weight: bold;");
    loadUsers();
});

/**
 * Carga la lista de usuarios desde la API de Go
 */
async function loadUsers() {
    const paginationInfo = document.getElementById('pagination-info');
    console.log("%c[API] Solicitando lista a: /api/v1/users", "color: #3b82f6;");

    try {
        // Usamos redirect: 'follow' pero validaremos el tipo de contenido
        const response = await fetch('/api/v1/users');
        
        console.log(`[NETWORK] Status: ${response.status} | URL: ${response.url}`);

        // 1. Verificar si la sesión expiró (401)
        if (response.status === 401) {
            console.warn("%c[AUTH] Sesión expirada o no válida (401). Redirigiendo...", "color: #f59e0b;");
            window.location.href = "/";
            return;
        }

        // 2. Detectar si el servidor nos redirigió a un HTML (Bucle de redirección)
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("text/html")) {
            console.error("%c[FATAL] El servidor devolvió HTML en lugar de JSON. Probablemente un bucle de redirección en el Middleware.", "color: #ef4444; font-weight: bold;");
            if (paginationInfo) paginationInfo.textContent = "ERROR: BUCLE DE REDIRECCIÓN";
            return;
        }

        if (!response.ok) throw new Error(`Error del servidor: ${response.status}`);

        allUsers = await response.json();
        console.log("%c[DATA] Registros cargados correctamente:", "color: #10b981;", allUsers.length);
        renderTable();

    } catch (error) {
        console.error("%c[FETCH ERROR] Fallo crítico en la comunicación:", "color: #ef4444;", error);
        if (paginationInfo) {
            paginationInfo.textContent = "ERROR DE CONEXIÓN / REDIRECCIÓN";
            paginationInfo.style.color = "#ef4444";
        }
    }
}

/**
 * Ordena los datos localmente
 */
function sortBy(field) {
    sortDirection = !sortDirection;
    console.log(`[SORT] Ordenando por: ${field} | ASC: ${sortDirection}`);
    
    allUsers.sort((a, b) => {
        let valA = (a[field] || "").toString().toLowerCase();
        let valB = (b[field] || "").toString().toLowerCase();
        if (valA < valB) return sortDirection ? -1 : 1;
        if (valA > valB) return sortDirection ? 1 : -1;
        return 0;
    });
    renderTable();
}

/**
 * Renderiza la tabla de usuarios
 */
function renderTable() {
    const pageSizeSelect = document.getElementById('pageSize');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 10;
    const tbody = document.getElementById('user-table-body');
    const info = document.getElementById('pagination-info');
    
    const displayData = (pageSize === 0) ? allUsers : allUsers.slice(0, pageSize);
    
    if (info) info.textContent = `REGISTROS: ${displayData.length} / ${allUsers.length}`;

    if (allUsers.length === 0) {
        tbody.innerHTML = `<tr><td colspan="4" class="p-10 text-center text-gray-400 uppercase text-[10px] font-black italic">No hay datos disponibles</td></tr>`;
        return;
    }

    tbody.innerHTML = displayData.map(u => {
        const userId = u.id || u.ID;
        return `
        <tr class="hover:bg-emerald-50/40 transition-colors border-b border-gray-100">
            <td class="px-3 py-2">
                <div class="flex flex-col leading-tight">
                    <span class="font-black text-black uppercase text-[10px] truncate">${u.username}</span>
                    <span class="text-[8px] text-gray-400 font-bold lowercase truncate italic">${u.email}</span>
                </div>
            </td>
            <td class="px-3 py-2 text-center">
                <span class="px-2 py-0.5 rounded text-[8px] font-black uppercase shadow-sm ${getRoleStyle(u.role)}">
                    ${u.role}
                </span>
            </td>
            <td class="px-3 py-2 text-center">
                <span class="text-[8px] font-black italic ${u.is_active ? 'text-emerald-600' : 'text-red-500'}">
                    ${u.is_active ? '● ACTIVO' : '○ BLOQUEADO'}
                </span>
            </td>
            <td class="px-3 py-2 text-right">
                <button onclick="window.location.href='/admin/users/manage?id=${userId}'" 
                    class="p-1.5 bg-black text-[#10b981] rounded hover:bg-[#10b981] hover:text-black transition-all active:scale-90">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                        <path d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                    </svg>
                </button>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Estilos de Tailwind para los roles
 */
function getRoleStyle(role) {
    const styles = {
        'admin':   'bg-black text-[#10b981] border border-[#10b981]/30',
        'driver':  'bg-emerald-100 text-emerald-800',
        'company': 'bg-slate-800 text-white',
        'client':  'bg-slate-100 text-slate-500'
    };
    return styles[role] || 'bg-gray-100 text-gray-400';
}

function exportData(type) {
    console.log(`[EXPORT] Preparando ${type}...`);
    alert(`Exportando ${allUsers.length} registros a ${type.toUpperCase()}`);
}