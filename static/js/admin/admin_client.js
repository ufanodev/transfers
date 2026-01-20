/**
 * @file admin_client.js
 * @description Gestión de la tabla de clientes: Carga, Búsqueda, Paginación y Acciones.
 */

let allClients = [];      // Datos brutos de la API
let filteredClients = []; // Datos tras filtros de búsqueda
let sortDirection = true; // true = ASC, false = DESC

document.addEventListener('DOMContentLoaded', () => {
    const adminStyle = "color: #3b82f6; font-weight: bold; background: #1e293b; padding: 3px 8px; border-radius: 5px;";
    console.log("%c[SYSTEM] Módulo Clientes Admin iniciado.", adminStyle);
    
    // Carga inicial
    loadClients();
});

/**
 * Carga la lista completa de clientes desde la API de Go
 */
async function loadClients() {
    const info = document.getElementById('pagination-info');
    console.log("[API] 📡 Solicitando: /api/v1/clients");

    try {
        const response = await fetch('/api/v1/clients');

        // 1. Manejo de Sesión Expirada
        if (response.status === 401) {
            console.warn("[AUTH] Sesión no válida. Redirigiendo...");
            window.location.href = "/";
            return;
        }

        if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);

        // 2. Almacenamiento de datos
        allClients = await response.json();
        filteredClients = [...allClients];
        
        console.log(`[DATA] ✅ ${allClients.length} clientes cargados.`);
        renderTable();

    } catch (error) {
        console.error("[FETCH ERROR]", error);
        if (info) {
            info.textContent = "ERROR DE CONEXIÓN CON API";
            info.style.color = "#ef4444";
        }
    }
}

/**
 * Filtra los clientes por Nombre, Email o Teléfono
 */
function filterClients() {
    const searchTerm = document.getElementById('clientSearch').value.toLowerCase();
    
    filteredClients = allClients.filter(c => {
        return (
            (c.full_name || "").toLowerCase().includes(searchTerm) ||
            (c.email || "").toLowerCase().includes(searchTerm) ||
            (c.phone || "").toLowerCase().includes(searchTerm)
        );
    });

    renderTable();
}

/**
 * Ordena los datos localmente
 * @param {string} field - Campo de la base de datos (full_name, email, etc)
 */
function sortBy(field) {
    sortDirection = !sortDirection;
    console.log(`[SORT] Ordenando por: ${field} | ASC: ${sortDirection}`);

    filteredClients.sort((a, b) => {
        let valA = (a[field] || "").toString().toLowerCase();
        let valB = (b[field] || "").toString().toLowerCase();
        
        if (valA < valB) return sortDirection ? -1 : 1;
        if (valA > valB) return sortDirection ? 1 : -1;
        return 0;
    });

    renderTable();
}

/**
 * Procesa y dibuja las filas en el tbody de la tabla
 */
function renderTable() {
    const tbody = document.getElementById('client-table-body');
    const info = document.getElementById('pagination-info');
    const pageSizeSelect = document.getElementById('pageSize');
    const pageSize = pageSizeSelect ? parseInt(pageSizeSelect.value) : 10;
    
    // Lógica de Paginación Visual
    const displayData = (pageSize === 0) ? filteredClients : filteredClients.slice(0, pageSize);
    
    if (info) {
        info.textContent = `Viendo ${displayData.length} de ${filteredClients.length} Clientes`;
    }

    // Caso: No hay resultados
    if (filteredClients.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="4" class="p-12 text-center">
                    <i class="fa-solid fa-address-book text-3xl text-slate-200 mb-2 block"></i>
                    <p class="text-slate-400 font-bold uppercase text-[9px] tracking-widest italic">No se encontraron registros</p>
                </td>
            </tr>`;
        return;
    }

    // Mapeo de Filas
    tbody.innerHTML = displayData.map(c => {
        const clientId = c.id || c.ID;
        return `
        <tr class="hover:bg-slate-50 transition-all border-b border-slate-100 group">
            <td class="px-5 py-2.5">
                <div class="flex items-center gap-3">
                    <div class="w-7 h-7 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 text-[9px] font-black border border-slate-200 group-hover:bg-admin-accent group-hover:text-white transition-all shadow-sm">
                        ${getInitials(c.full_name)}
                    </div>
                    <div class="flex flex-col leading-tight">
                        <span class="font-bold text-slate-700 uppercase tracking-tighter">${c.full_name || 'SIN NOMBRE'}</span>
                        <span class="text-[9px] text-slate-400 lowercase font-medium italic">${c.email || 'sin email'}</span>
                    </div>
                </div>
            </td>
            <td class="px-5 py-2.5 text-slate-500">
                <div class="flex items-center gap-1.5 font-bold">
                    <i class="fa-solid fa-phone text-[8px] opacity-30"></i>
                    <span>${c.phone || 'N/A'}</span>
                </div>
            </td>
            <td class="px-5 py-2.5 text-center">
                <span class="px-2 py-0.5 rounded-full bg-slate-100 text-slate-500 text-[8px] font-black uppercase border border-slate-200">
                    ${c.preferred_language || 'ES'}
                </span>
            </td>
            <td class="px-5 py-2.5 text-right">
                <div class="flex justify-end gap-1.5">
                    <button onclick="window.location.href='/admin/clients/manage?id=${clientId}'" 
                        class="w-7 h-7 flex items-center justify-center bg-white text-slate-400 border border-slate-200 rounded-lg hover:bg-admin-accent hover:text-white hover:border-admin-accent transition-all shadow-sm active:scale-90"
                        title="Editar Cliente">
                        <i class="fa-solid fa-pen-to-square text-[10px]"></i>
                    </button>
                    <button onclick="deleteClient(${clientId})" 
                        class="w-7 h-7 flex items-center justify-center bg-white text-red-300 border border-red-100 rounded-lg hover:bg-red-500 hover:text-white hover:border-red-500 transition-all shadow-sm active:scale-90"
                        title="Eliminar">
                        <i class="fa-solid fa-trash-can text-[10px]"></i>
                    </button>
                </div>
            </td>
        </tr>`;
    }).join('');
}

/**
 * Genera iniciales del nombre para el avatar
 */
function getInitials(name) {
    if (!name) return "??";
    return name.split(" ").map(n => n[0]).slice(0, 2).join("").toUpperCase();
}

/**
 * Placeholder para exportar datos
 */
function exportData(type) {
    console.log(`[EXPORT] Formato: ${type}`);
    alert(`Iniciando descarga de ${filteredClients.length} clientes en formato ${type.toUpperCase()}...`);
}

/**
 * Función para borrar (Redirige al CRUD para usar el modal de confirmación)
 */
function deleteClient(id) {
    // Redirigimos al manage con el ID para que el admin pueda usar el borrado seguro allí
    window.location.href = `/admin/clients/manage?id=${id}&action=delete`;
}