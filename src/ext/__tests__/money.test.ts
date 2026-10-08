import { describe, expect, test } from 'vitest';

import { computeTotals, lineTotal } from '@/ext/money.ts';

describe('lineTotal (must match the Go implementation)', () => {
    test('whole quantities', () => {
        expect(lineTotal(2000, 1500)).toBe(3000);
        expect(lineTotal(1000, 999)).toBe(999);
    });

    test('rounds half up like the server', () => {
        expect(lineTotal(500, 999)).toBe(500); // 0.5 x 9.99 = 4.995 -> 500 minor units
        expect(lineTotal(1, 999)).toBe(1); // 0.001 x 999 = 0.999 -> 1
        expect(lineTotal(1, 400)).toBe(0); // 0.4 -> 0
        expect(lineTotal(1, 500)).toBe(1); // 0.5 -> 1
    });

    test('is exact for huge products', () => {
        expect(lineTotal(999_999_999_999, 9_999_999_999_999)).toBe(Number.MAX_SAFE_INTEGER);
        expect(lineTotal(3_000_000_000, 3)).toBe(9_000_000);
    });

    test('invalid input is zero', () => {
        expect(lineTotal(-1, 100)).toBe(0);
        expect(lineTotal(100, -1)).toBe(0);
        expect(lineTotal(1.5, 100)).toBe(0);
        expect(lineTotal(Number.NaN, 100)).toBe(0);
    });
});

describe('computeTotals', () => {
    test('pay in full', () => {
        expect(computeTotals([3000, 2000], 0, 'full', 0)).toEqual({ subtotal: 5000, discount: 0, total: 5000, paid: 5000, credit: 0 });
    });

    test('everything on credit', () => {
        expect(computeTotals([5000], 0, 'credit', 999)).toEqual({ subtotal: 5000, discount: 0, total: 5000, paid: 0, credit: 5000 });
    });

    test('part now, rest on credit', () => {
        expect(computeTotals([5000], 0, 'partial', 2000)).toEqual({ subtotal: 5000, discount: 0, total: 5000, paid: 2000, credit: 3000 });
    });

    test('a discount reduces the total before payment is split', () => {
        expect(computeTotals([5000], 400, 'partial', 2000)).toEqual({ subtotal: 5000, discount: 400, total: 4600, paid: 2000, credit: 2600 });
        expect(computeTotals([5000], 400, 'full', 0).paid).toBe(4600);
    });

    test('values the server would refuse are clamped', () => {
        expect(computeTotals([5000], 9999, 'full', 0).total).toBe(0);
        expect(computeTotals([5000], -50, 'full', 0).discount).toBe(0);
        expect(computeTotals([5000], 0, 'partial', 9999).paid).toBe(5000);
        expect(computeTotals([5000], 0, 'partial', -5).paid).toBe(0);
    });

    test('an empty cart is zero', () => {
        expect(computeTotals([], 0, 'full', 0)).toEqual({ subtotal: 0, discount: 0, total: 0, paid: 0, credit: 0 });
    });
});
