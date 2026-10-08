import { describe, expect, test } from 'vitest';

import { formatQty, parseQty, parseSignedQty, QTY_SCALE } from '@/ext/qty.ts';

describe('parseQty', () => {
    test('parses whole numbers and decimals into scaled integers', () => {
        expect(parseQty('2')).toBe(2000);
        expect(parseQty('2.5')).toBe(2500);
        expect(parseQty('0.125')).toBe(125);
        expect(parseQty(' 7 ')).toBe(7000);
        expect(parseQty('0')).toBe(0);
    });

    test('accepts a decimal comma', () => {
        expect(parseQty('2,5')).toBe(2500);
    });

    test('does not suffer from floating point errors', () => {
        expect(parseQty('0.1')).toBe(100);
        expect(parseQty('0.7')).toBe(700);
        expect(parseQty('1.005')).toBe(1005);
        expect(parseQty('19.999')).toBe(19999);
    });

    test('rejects anything that is not a plain number with at most 3 decimals', () => {
        for (const bad of ['', ' ', 'abc', '1.2345', '-1', '1e3', '1.', '.5', '1..2', '1 000', '0x10']) {
            expect(parseQty(bad), bad).toBeNull();
        }
    });

    test('rejects numbers too large to stay exact', () => {
        expect(parseQty('99999999999999999999')).toBeNull();
    });
});

describe('formatQty', () => {
    test('drops needless decimals', () => {
        expect(formatQty(2500)).toBe('2.5');
        expect(formatQty(3000)).toBe('3');
        expect(formatQty(125)).toBe('0.125');
        expect(formatQty(0)).toBe('0');
        expect(formatQty(1005)).toBe('1.005');
        expect(formatQty(10)).toBe('0.01');
    });

    test('handles negative quantities', () => {
        expect(formatQty(-2500)).toBe('-2.5');
        expect(formatQty(-1)).toBe('-0.001');
    });

    test('round trips with parseQty', () => {
        for (const value of [0, 1, 10, 100, 999, 1000, 1234, 99999, 123456789]) {
            expect(parseQty(formatQty(value))).toBe(value);
        }
    });
});

describe('parseSignedQty', () => {
    test('accepts a leading minus', () => {
        expect(parseSignedQty('-2.5')).toBe(-2500);
        expect(parseSignedQty('3')).toBe(3000);
    });

    test('rejects zero adjustments and bad input', () => {
        expect(parseSignedQty('-0')).toBeNull();
        expect(parseSignedQty('-')).toBeNull();
        expect(parseSignedQty('--1')).toBeNull();
    });
});

test('scale matches the backend', () => {
    expect(QTY_SCALE).toBe(1000);
});
