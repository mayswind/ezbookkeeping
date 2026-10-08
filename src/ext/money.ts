// Sale arithmetic. It mirrors pkg/ext/services/finance.go exactly (integers only, rounding half up), so the total
// shown on screen is the total the server books.

/** Quantity (scaled by 1000) times unit price per whole unit, in minor currency units, rounded half up. */
export function lineTotal(qty: number, unitPrice: number): number {
    if (!Number.isSafeInteger(qty) || !Number.isSafeInteger(unitPrice) || qty < 0 || unitPrice < 0) {
        return 0;
    }

    // BigInt keeps the product exact even when qty * price exceeds 2^53
    const product = BigInt(qty) * BigInt(unitPrice) + 500n;
    const total = product / 1000n;

    return total > BigInt(Number.MAX_SAFE_INTEGER) ? Number.MAX_SAFE_INTEGER : Number(total);
}

export type PayMode = 'full' | 'partial' | 'credit';

export interface SaleTotals {
    readonly subtotal: number;
    readonly discount: number;
    readonly total: number;
    readonly paid: number;
    readonly credit: number;
}

/**
 * Works out what is paid now and what goes on credit.
 * The discount is clamped to the subtotal, and a partial payment to the total, so the numbers can never be
 * ones the server would refuse.
 */
export function computeTotals(lineTotals: number[], discount: number, mode: PayMode, partialPaid: number): SaleTotals {
    const subtotal = lineTotals.reduce((sum, value) => sum + value, 0);
    const finalDiscount = Math.min(Math.max(discount, 0), subtotal);
    const total = subtotal - finalDiscount;

    let paid = total;

    if (mode === 'credit') {
        paid = 0;
    } else if (mode === 'partial') {
        paid = Math.min(Math.max(partialPaid, 0), total);
    }

    return { subtotal, discount: finalDiscount, total, paid, credit: total - paid };
}
