/**
 * @file admin_booking.js
 * @description Listado dinámico de reservas TakeUs
 */

let allBookings = [];

document.addEventListener('DOMContentLoaded', async () => {
    await loadBookings();
});

async function loadBookings() {
    const tbody = document.getElementById('booking-table-body');
    try {
        const res = await fetch('/api/v1/bookings');
        if (!res.ok) throw new Error();
        allBookings = await res.json();
        renderTable(allBookings);
    } catch (e) {
        tbody.innerHTML = `<tr><td colspan="6" class="p-8 text-center text-red-400 font-bold uppercase text-[10px]">Error al sincronizar</td></tr>`;
    }
}

function renderTable(data) {
    const tbody = document.getElementById('booking-table-body');
    const colors = {
        'pending': 'bg-amber-50 text-amber-600 border-amber-100',
        'confirmed': 'bg-blue-50 text-blue-600 border-blue-100',
        'cancelled': 'bg-red-50 text-red-600 border-red-100'
    };

    tbody.innerHTML = data.map(b => `
        <tr class="hover:bg-slate-50 transition-colors border-b border-slate-100 group">
            <td class="px-5 py-3 font-black text-slate-700 uppercase">
                ${b.scheduled_date.split('T')[0]} <br>
                <span class="text-admin-accent text-[9px]">${b.scheduled_time}</span>
            </td>
            <td class="px-5 py-3">
                <div class="font-black text-slate-800 uppercase italic text-[11px]">${b.client_name}</div>
                <div class="text-[9px] text-slate-400 font-bold">${b.client_phone}</div>
            </td>
            <td class="px-5 py-3">
                <div class="text-[10px] text-slate-400 font-bold uppercase truncate max-w-[180px]">O: ${b.origin_address}</div>
                <div class="text-[10px] text-slate-800 font-black uppercase truncate max-w-[180px]">D: ${b.dest_address}</div>
            </td>
            <td class="px-5 py-3 text-center">
                <div class="flex flex-col items-center justify-center bg-slate-50 rounded-lg py-1 px-2 border border-slate-100">
                    <span class="font-black text-slate-700 text-[10px]">${b.pax} PAX</span>
                    <span class="text-slate-500 text-[9px] font-bold">${b.maletas} MAL.</span>
                    ${b.animal ? '<i class="fa-solid fa-dog text-[9px] text-amber-500 mt-0.5"></i>' : ''}
                </div>
            </td>
            <td class="px-5 py-3">
                <span class="px-3 py-1 rounded-full text-[8px] font-black uppercase border ${colors[b.status] || 'bg-slate-100'}">
                    ${b.status}
                </span>
            </td>
            <td class="px-5 py-3 text-right">
                <div class="flex justify-end gap-2">
                    <button onclick="window.location.href='/admin/bookings/manage?id=${b.id || b.ID}'" class="w-8 h-8 flex items-center justify-center bg-white text-slate-400 rounded-lg hover:bg-admin-accent hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-pen-to-square text-[11px]"></i>
                    </button>
                    <button onclick="deleteBooking(${b.id || b.ID})" class="w-8 h-8 flex items-center justify-center bg-white text-red-300 rounded-lg hover:bg-red-600 hover:text-white transition-all shadow-sm border border-slate-100">
                        <i class="fa-solid fa-trash-can text-[11px]"></i>
                    </button>
                </div>
            </td>
        </tr>
    `).join('');
}

function filterBookings() {
    const q = document.getElementById('bookingSearch').value.toLowerCase();
    const filtered = allBookings.filter(b => 
        b.client_name.toLowerCase().includes(q) || 
        b.dest_address.toLowerCase().includes(q)
    );
    renderTable(filtered);
}

async function deleteBooking(id) {
    if(!confirm("¿Desea cancelar esta reserva?")) return;
    try {
        const res = await fetch(`/api/v1/bookings/${id}`, { method: 'DELETE' });
        if(res.ok) loadBookings();
    } catch(e) { console.error(e); }
}