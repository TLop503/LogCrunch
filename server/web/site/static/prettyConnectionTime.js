// Parse long verbose date into simple format.
document.addEventListener("DOMContentLoaded", () => {
    const rows = document.querySelectorAll("table tr:not(:first-child)");

    rows.forEach(row => {
        const td = row.querySelectorAll("td")[3];
        if (!td) return;

        const raw = td.textContent.trim();

        // Extract timestamp with fractional seconds (if present)
        const match = raw.match(/^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})(\.\d+)?/);
        if (!match) return;

        const baseTime = match[1]; // e.g. "2025-07-13 23:19:16"
        const fractional = match[2] || ""; // e.g. ".891612889" or ""

        // Truncate fractional seconds to 3 digits
        const truncatedFrac = fractional ? fractional.slice(0, 4) : ""; // including dot + 3 digits

        // Construct a string to parse (replace space with T for ISO format)
        const isoStr = (baseTime + truncatedFrac).replace(' ', 'T');

        const dateObj = new Date(isoStr);

        // failsafe
        if (isNaN(dateObj.getTime())) {
            td.textContent = raw; // fallback original text
            return;
        }

        // Helper to pad numbers
        const pad = (n) => n.toString().padStart(2, '0');

        // Extract milliseconds and pad to 3 digits
        const ms = dateObj.getMilliseconds().toString().padStart(3, '0');

        td.textContent = `${dateObj.getFullYear()}-${pad(dateObj.getMonth() + 1)}-${pad(dateObj.getDate())} `
            + `${pad(dateObj.getHours())}:${pad(dateObj.getMinutes())}:${pad(dateObj.getSeconds())}.${ms}`;
    });
});
