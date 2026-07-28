import { describe, expect, it } from 'vitest';

import { evaluateExpressionToAmount } from '@/lib/evaluator.ts';

describe('evaluateExpressionToAmount', () => {
    it.each([
        ['', undefined],
        ['1+2', '300'],
        ['(1+2)*3', '900'],
        ['-1+2', '100'],
        ['1.5+2.5', '400'],
        ['1+2*3-(4/2)', '500'],
        ['2*-3-3/-2', '-450'],
        ['-1.2-3.4*(-5.6/7.8*(9.0-1.2))', '1784'],
        ['(((2+3)))*(((((-5+7)))))', '1000'],
        ['3.5+0.1', '360'],
        ['3.55+0.11', '366'],
        ['3.555+0.111', '366'],
        ['0.1234567+1', '112'],
        ['9999999999999.99-9999999999999.98', '1'],
        ['9999999999999.99/9999999999999.99', '100'],
        ['9999999999999.99*9223-9999999999999.99*9222', '999999999999999'],
        ['-9999999999999.99*9223+9999999999999.99*9222', '-999999999999999'],
    ])('evaluates valid expression %j', (expression, expected) => {
        expect(evaluateExpressionToAmount(expression as string)?.toString()).toBe(expected);
    });

    it.each([
        '1++2',
        '1^2',
        '+-*/',
        'a+b',
        '1/0',
        '1+(2*3',
        '1+2*3)',
        '1+((((2*3)))',
        '1+2(3)',
        '1)*(2',
        '0.abcd+1',
    ])('rejects invalid expression %j', expression => {
        expect(evaluateExpressionToAmount(expression)).toBeUndefined();
    });

    it.each([
        '9999999999999.99+0.01',
        '-9999999999999.99-0.01',
    ])('rejects numeric overflow in expression %j', expression => {
        expect(() => evaluateExpressionToAmount(expression)).toThrow('Numeric Overflow');
    });
});
