// Download a table as a CSV file that opens correctly in Excel and Google Sheets.

/** Cells that start with these characters would be run as formulas by a spreadsheet, so they get a quote in front. */
const FORMULA_START = /^[=+\-@\t\r]/;

function cell(value: string | number): string {
    let text = String(value);

    if (typeof value === 'string' && FORMULA_START.test(text)) {
        text = `'${text}`;
    }

    return /[",\n\r]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text;
}

/** Turns rows into CSV text (lines end with CRLF as the format asks). */
export function toCsv(rows: (string | number)[][]): string {
    return rows.map(row => row.map(cell).join(',')).join('\r\n');
}

/** Saves rows as a file in the browser. The byte order mark makes Excel read accented letters correctly. */
export function downloadCsv(fileName: string, rows: (string | number)[][]): void {
    const blob = new Blob(['﻿' + toCsv(rows)], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');

    link.href = url;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
}
