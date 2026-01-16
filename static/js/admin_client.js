/**
 * @file admin_client.js
 * @description Gestión de la tabla de clientes (Listado y Eliminación).
 */

let allClients = [];

document.addEventListener('DOMContentLoaded', () => {
    console.log("%c[SYSTEM] Módulo Clientes iniciado.", "color: #10b981; font-weight: bold;");
    loadClients();
});

/**
 * Obtiene los clientes desde la API de backend
 */
async function loadClients() {
    const paginationInfo = document.getElementById('pagination-info');
    try {
        const response = await fetch('/api/v1/clients');
        
        // Manejo de sesión expirada o no autorizada
        if (response.status === 401) {
            window.location.href = "/";
            return;
        }

        // Seguridad: Verificar que el servidor no responda con la página de Login (HTML)
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("text/html")) {
            console.error("[ERROR] La API devolvió HTML inesperado. Revisa la sesión.");
            if (paginationInfo) paginationInfo.textContent = "SESIÓN EXPIRADA";
            return;
        }

        allClients = await response.json();
        console.log("[DATA] Clientes recibidos:", allClients);
        renderTable();

    } catch (error) {
        console.error("[FETCH ERROR]:", error);
        if (paginationInfo) paginationInfo.textContent = "ERROR DE CONEXIÓN";
    }
}

/**
 * Renderiza los datos en el tbody de la tabla
 */
function renderTable() {
    const pageSizeSelect = document.getElementById('pageSize');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 10;
    const tbody = document.getElementById('client-table-body');
    const info = document.getElementById('pagination-info');
    
    // Filtro de paginación local (0 significa "ALL")
    const displayData = (pageSize === 0) ? allClients : allClients.slice(0, pageSize);
    
    if (info) info.textContent = `CLIENTES: ${displayData.length} / ${allClients.length}`;

    // Si no hay datos, mostrar mensaje amigable
    if (!allClients || allClients.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="4" class="p-10 text-center text-gray-400 uppercase text-[9px] font-black italic">
                    No hay clientes registrados en el sistema
                </td>
            </tr>`;
        return;
    }

    tbody.innerHTML = displayData.map(c => {
        // NORMALIZACIÓN DE CAMPOS (Soporta GORM PascalCase y JSON camelCase)
        const id = c.id || c.ID;
        const fullName = c.full_name || c.FullName || "Sin Nombre";
        const email = c.email || c.Email || "N/A";
        const phone = c.phone || c.Phone || "N/A";
        const userId = c.user_id || c.UserID || "—";

        return `
        <tr class="hover:bg-emerald-50/40 transition-colors border-b border-gray-100">
            <td class="px-3 py-2 overflow-hidden">
                <div class="flex flex-col leading-none">
                    <span class="font-black text-black uppercase text-[10px] truncate">${fullName}</span>
                    <span class="text-[8px] text-gray-400 font-bold lowercase truncate">${email}</span>
                </div>
            </td>
            <td class="px-3 py-2 text-center text-[9px] font-black text-slate-500 uppercase">
                ${phone}
            </td>
            <td class="px-3 py-2 text-center">
                <span class="text-[8px] font-black italic text-emerald-600 uppercase tracking-tighter">
                    ID Usuario: ${userId}
                </span>
            </td>
            <td class="px-3 py-2 text-right">
                <div class="flex justify-end gap-1.5">
                    <button onclick="window.location.href='/dashboard/clients/manage?id=${id}'" 
                        class="p-1.5 bg-black text-[#10b981] rounded-md hover:bg-[#10b981] hover:text-black transition-all shadow-sm active:scale-90">
                        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                            <path d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                    </button>
                    <button onclick="deleteClient(${id})" 
                        class="p-1.5 bg-red-50 text-red-600 rounded-md hover:bg-red-600 hover:text-white transition-all shadow-sm active:scale-90">
                        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24">
                            <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Elimina un cliente mediante la API
 */
async function deleteClient(id) {
    if (!id || !confirm("¿Eliminar este cliente de forma permanente?")) return;

    try {
        const response = await fetch(`/api/v1/clients/${id}`, { 
            method: 'DELETE' 
        });

        if (response.ok) {
            console.log(`%c[DELETE] Cliente ${id} borrado.`, "color: #ef4444");
            loadClients(); 
        } else {
            const errorData = await response.json();
            alert("No se pudo eliminar el registro: " + (errorData.error || "Error desconocido"));
        }
    } catch (e) {
        console.error("Error en operación DELETE:", e);
        alert("Error crítico al intentar conectar con el servidor.");
    }
}

/**
 * Función global para los botones de exportación del Nav
 */
function exportData(type) {
    if (allClients.length === 0) {
        alert("No hay datos para exportar.");
        return;
    }
    alert(`Exportando ${allClients.length} registros a ${type.toUpperCase()}...`);
}