import { describe, expect, test } from 'vitest';

import { toCsv } from '@/ext/csv.ts';

describe('toCsv', () => {
    test('plain rows', () => {
        expect(toCsv([['Item', 'Qty'], ['Rice', 2.5]])).toBe('Item,Qty\r\nRice,2.5');
    });

    test('quotes cells with commas, quotes and line breaks', () => {
        expect(toCsv([['a,b', 'say "hi"', 'two\nlines']])).toBe('"a,b","say ""hi""","two\nlines"');
    });

    test('keeps accented and non-latin text', () => {
        expect(toCsv([['Zoë', '日本語']])).toBe('Zoë,日本語');
    });

    test('defuses spreadsheet formulas typed into names or notes', () => {
        expect(toCsv([['=HYPERLINK("http://evil")', '+1', '-2', '@SUM(A1)']])).toBe('"\'=HYPERLINK(""http://evil"")",\'+1,\'-2,\'@SUM(A1)');
    });

    test('numbers, including negative ones, are written as they are', () => {
        expect(toCsv([[-5, 0, 12.5]])).toBe('-5,0,12.5');
    });

    test('empty table', () => {
        expect(toCsv([])).toBe('');
    });
});
