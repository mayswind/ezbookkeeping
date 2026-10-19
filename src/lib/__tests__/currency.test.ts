import { describe, expect, it } from 'vitest';

import { parseQuickSwitchCurrencies } from '@/lib/currency.ts';

describe('parseQuickSwitchCurrencies', () => {
    it('keeps valid unique currencies in order', () => {
        expect(parseQuickSwitchCurrencies('USD,EUR,USD,XXXX,,JPY')).toEqual(['USD', 'EUR', 'JPY']);
    });

    it('returns empty list for empty value', () => {
        expect(parseQuickSwitchCurrencies('')).toEqual([]);
    });

    it('limits the number of currencies', () => {
        expect(parseQuickSwitchCurrencies('USD,EUR,JPY,GBP,CNY,CHF,CAD')).toEqual(['USD', 'EUR', 'JPY', 'GBP', 'CNY', 'CHF']);
    });
});
