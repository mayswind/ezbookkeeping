// Quantities are stored as integers scaled by 1000, so 2.5 units is 2500 and nothing ever touches floating point.

export const QTY_SCALE = 1000;
const QTY_DECIMALS = 3;

/**
 * Parses text typed by a person ("2", "2.5", "2,5", " 0.125 ") into a scaled quantity.
 * Returns null for anything that is not a plain non-negative number with at most 3 decimals.
 */
export function parseQty(text: string): number | null {
    const trimmed = text.trim().replace(',', '.');

    if (!/^\d+(\.\d{1,3})?$/.test(trimmed)) {
        return null;
    }

    const [whole = '0', fraction = ''] = trimmed.split('.');
    const scaled = Number(whole) * QTY_SCALE + Number(fraction.padEnd(QTY_DECIMALS, '0'));

    return Number.isSafeInteger(scaled) ? scaled : null;
}

/** Formats a scaled quantity for display, dropping needless decimals: 2500 -> "2.5", 3000 -> "3". */
export function formatQty(scaled: number): string {
    const negative = scaled < 0;
    const abs = Math.abs(scaled);
    const whole = Math.floor(abs / QTY_SCALE);
    const fraction = String(abs % QTY_SCALE).padStart(QTY_DECIMALS, '0').replace(/0+$/, '');
    const text = fraction ? `${whole}.${fraction}` : String(whole);

    return negative ? `-${text}` : text;
}

/** Like parseQty but accepts a leading minus, for stock adjustments. */
export function parseSignedQty(text: string): number | null {
    const trimmed = text.trim();

    if (trimmed.startsWith('-')) {
        const value = parseQty(trimmed.substring(1));
        return value === null || value === 0 ? null : -value;
    }

    return parseQty(trimmed);
}
