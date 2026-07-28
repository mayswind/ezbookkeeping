import { describe, expect, it } from 'vitest';

import {
    BIG_DECIMAL_ZERO,
    BIG_DECIMAL_NEGATIVE_INFINITY,
    BIG_DECIMAL_POSITIVE_INFINITY,
    parseBigDecimal,
    isBigDecimal
} from '@/lib/numeral.ts';

describe('parseBigDecimal', () => {
    it.each([
        ['0', '0'],
        ['-0', '0'],
        ['00123.4500', '123.45'],
        ['1e+21', '1000000000000000000000'],
        ['-99999999999999999999.99999999999999999999', '-99999999999999999999.99999999999999999999'],
    ])('should parse string %j without losing precision', (value, expected) => {
        expect(parseBigDecimal(value).toString()).toBe(expected);
    });

    it('should parse numbers', () => {
        expect(parseBigDecimal(123.45).toString()).toBe('123.45');
    });

    it('should return NaN for unsupported runtime values', () => {
        expect(parseBigDecimal(null as unknown as string).isNaN()).toBe(true);
    });

    it('should reject invalid decimal strings', () => {
        expect(() => parseBigDecimal('not-a-number')).toThrow();
    });
});

describe('isBigDecimal', () => {
    it('should identify values created by parseBigDecimal', () => {
        expect(isBigDecimal(parseBigDecimal('1'))).toBe(true);
    });

    it.each([0, '1', null, undefined, {}, Number.NaN])('should reject non-BigDecimal value %j', value => {
        expect(isBigDecimal(value)).toBe(false);
    });
});

describe('BigDecimal predicates', () => {
    it.each([
        ['0', true, true, false, false, true, true],
        ['-0', true, true, false, false, true, true],
        ['1', false, true, true, false, true, false],
        ['-1', false, true, false, true, false, true],
        ['Infinity', false, false, true, false, true, false],
        ['-Infinity', false, false, false, true, false, true],
        ['NaN', false, false, false, false, false, false],
    ])(
        'should report the state of %s',
        (value, zero, finite, positive, negative, positiveOrZero, negativeOrZero) => {
            const decimal = parseBigDecimal(value);

            expect(decimal.isZero()).toBe(zero);
            expect(decimal.isFinite()).toBe(finite);
            expect(decimal.isPositive()).toBe(positive);
            expect(decimal.isNegative()).toBe(negative);
            expect(decimal.isPositiveOrZero()).toBe(positiveOrZero);
            expect(decimal.isNegativeOrZero()).toBe(negativeOrZero);
            expect(decimal.isPositiveInfinity()).toBe(value === 'Infinity');
            expect(decimal.isNegativeInfinity()).toBe(value === '-Infinity');
            expect(decimal.isNaN()).toBe(value === 'NaN');
        }
    );

    it('should expose zero and infinity constants', () => {
        expect(BIG_DECIMAL_ZERO.isZero()).toBe(true);
        expect(BIG_DECIMAL_POSITIVE_INFINITY.isPositiveInfinity()).toBe(true);
        expect(BIG_DECIMAL_NEGATIVE_INFINITY.isNegativeInfinity()).toBe(true);
    });
});

describe('BigDecimal comparisons', () => {
    it('should compare with BigDecimal operands', () => {
        const value = parseBigDecimal('1.00000000000000000001');
        const smaller = parseBigDecimal('1.00000000000000000000');

        expect(value.equals(parseBigDecimal('1.00000000000000000001'))).toBe(true);
        expect(value.notEquals(smaller)).toBe(true);
        expect(value.compareTo(smaller)).toBe(1);
        expect(value.greaterThan(smaller)).toBe(true);
        expect(value.greaterThanOrEqual(value)).toBe(true);
        expect(smaller.lessThan(value)).toBe(true);
        expect(smaller.lessThanOrEqual(value)).toBe(true);
    });

    it('should compare with number operands', () => {
        const value = parseBigDecimal('2');

        expect(value.equals(2)).toBe(true);
        expect(value.notEquals(3)).toBe(true);
        expect(value.compareTo(2)).toBe(0);
        expect(value.greaterThan(1)).toBe(true);
        expect(value.greaterThanOrEqual(2)).toBe(true);
        expect(value.lessThan(3)).toBe(true);
        expect(value.lessThanOrEqual(2)).toBe(true);
    });

    it('should handle undefined equality operands', () => {
        const value = parseBigDecimal('1');

        expect(value.equals(undefined)).toBe(false);
        expect(value.notEquals(undefined)).toBe(true);
    });

    it('should include both boundaries in between checks', () => {
        const value = parseBigDecimal('2');

        expect(value.between(parseBigDecimal('2'), 3)).toBe(true);
        expect(value.between(1, parseBigDecimal('2'))).toBe(true);
        expect(value.between(3, 1)).toBe(false);
    });

    it('should make comparisons involving NaN unordered', () => {
        const value = parseBigDecimal('NaN');

        expect(value.equals(value)).toBe(false);
        expect(value.notEquals(value)).toBe(true);
        expect(value.compareTo(0)).toBeNaN();
        expect(value.greaterThan(0)).toBe(false);
        expect(value.greaterThanOrEqual(0)).toBe(false);
        expect(value.lessThan(0)).toBe(false);
        expect(value.lessThanOrEqual(0)).toBe(false);
        expect(value.between(-1, 1)).toBe(false);
    });
});

describe('BigDecimal arithmetic', () => {
    it('should perform exact arithmetic with BigDecimal operands', () => {
        const value = parseBigDecimal('0.1');
        const operand = parseBigDecimal('0.2');

        expect(value.add(operand).toString()).toBe('0.3');
        expect(value.subtract(operand).toString()).toBe('-0.1');
        expect(value.multiply(operand).toString()).toBe('0.02');
        expect(value.divide(operand).toString()).toBe('0.5');
    });

    it('should perform arithmetic with number operands', () => {
        const value = parseBigDecimal('10');

        expect(value.add(2).toString()).toBe('12');
        expect(value.subtract(2).toString()).toBe('8');
        expect(value.multiply(2).toString()).toBe('20');
        expect(value.divide(2).toString()).toBe('5');
    });

    it('should preserve the original immutable value', () => {
        const value = parseBigDecimal('1');

        value.add(1);
        value.subtract(1);
        value.multiply(2);
        value.divide(2);

        expect(value.toString()).toBe('1');
    });

    it('should handle division by zero according to decimal.js semantics', () => {
        expect(parseBigDecimal('1').divide(0).isPositiveInfinity()).toBe(true);
        expect(parseBigDecimal('-1').divide(BIG_DECIMAL_ZERO).isNegativeInfinity()).toBe(true);
        expect(BIG_DECIMAL_ZERO.divide(0).isNaN()).toBe(true);
    });

    it('should negate values and return their signs', () => {
        expect(parseBigDecimal('12.3').negate().toString()).toBe('-12.3');
        expect(parseBigDecimal('-12.3').negate().toString()).toBe('12.3');
        expect(parseBigDecimal('9').sign().toString()).toBe('1');
        expect(parseBigDecimal('-9').sign().toString()).toBe('-1');
        expect(BIG_DECIMAL_ZERO.sign().toString()).toBe('0');
        expect(parseBigDecimal('-0').sign().toString()).toBe('0');
        expect(parseBigDecimal('NaN').sign().toString()).toBe('0');
    });
});

describe('BigDecimal mathematical functions', () => {
    it('should calculate powers with BigDecimal and number exponents', () => {
        expect(parseBigDecimal('2').pow(parseBigDecimal('10')).toString()).toBe('1024');
        expect(parseBigDecimal('2').pow(-2).toString()).toBe('0.25');
        expect(BIG_DECIMAL_ZERO.pow(0).toString()).toBe('1');
        expect(parseBigDecimal('-2').pow(parseBigDecimal('0.5')).isNaN()).toBe(true);
    });

    it('should calculate square roots and reject negative radicands', () => {
        expect(parseBigDecimal('2.25').sqrt().toString()).toBe('1.5');
        expect(BIG_DECIMAL_ZERO.sqrt().toString()).toBe('0');
        expect(parseBigDecimal('-1').sqrt().isNaN()).toBe(true);
    });

    it('should calculate natural logarithms at their boundaries', () => {
        expect(parseBigDecimal('1').log().toString()).toBe('0');
        expect(BIG_DECIMAL_ZERO.log().isNegativeInfinity()).toBe(true);
        expect(parseBigDecimal('-1').log().isNaN()).toBe(true);
    });

    it('should calculate exponentials', () => {
        expect(BIG_DECIMAL_ZERO.exp().toString()).toBe('1');
        expect(BIG_DECIMAL_NEGATIVE_INFINITY.exp().toString()).toBe('0');
        expect(BIG_DECIMAL_POSITIVE_INFINITY.exp().isPositiveInfinity()).toBe(true);
    });

    it('should calculate absolute and truncated values', () => {
        expect(parseBigDecimal('-12.34').abs().toString()).toBe('12.34');
        expect(parseBigDecimal('12.99').truncate().toString()).toBe('12');
        expect(parseBigDecimal('-12.99').truncate().toString()).toBe('-12');
        expect(BIG_DECIMAL_POSITIVE_INFINITY.truncate().isPositiveInfinity()).toBe(true);
    });
});

describe('BigDecimal.toSafeIntegerNumber', () => {
    it.each([
        ['9007199254740991', 9007199254740991],
        ['-9007199254740991', -9007199254740991],
        ['1', 1],
        ['-1', -1],
        ['0', 0],
        ['-0', -0],
    ])('should convert safe integer %s to a number', (value, expected) => {
        const result = parseBigDecimal(value).toSafeIntegerNumber();

        if (Object.is(expected, -0)) {
            expect(Object.is(result, -0)).toBe(true);
        } else {
            expect(result).toBe(expected);
        }
    });

    it('should convert a BigDecimal created from a number', () => {
        expect(parseBigDecimal(Number.MAX_SAFE_INTEGER).toSafeIntegerNumber()).toBe(Number.MAX_SAFE_INTEGER);
        expect(parseBigDecimal(Number.MIN_SAFE_INTEGER).toSafeIntegerNumber()).toBe(Number.MIN_SAFE_INTEGER);
    });

    it.each([
        '9007199254740992',
        '-9007199254740992',
        '1.5',
        'NaN',
        'Infinity',
        '-Infinity',
    ])('should reject unsafe integer conversion for %s', value => {
        expect(() => parseBigDecimal(value).toSafeIntegerNumber()).toThrow(
            'cannot be converted to a safe integer number'
        );
    });

    it('should validate the rounded double result near the safe integer boundary', () => {
        expect(parseBigDecimal('9007199254740991.1').toSafeIntegerNumber()).toBe(Number.MAX_SAFE_INTEGER);
        expect(parseBigDecimal('-9007199254740991.1').toSafeIntegerNumber()).toBe(Number.MIN_SAFE_INTEGER);
    });

    it('should include the converted value in the error message', () => {
        expect(() => parseBigDecimal('1.5').toSafeIntegerNumber()).toThrow(
            'big decimal value "1.5" cannot be converted to a safe integer number'
        );
        expect(() => parseBigDecimal('NaN').toSafeIntegerNumber()).toThrow(
            'big decimal value "NaN" cannot be converted to a safe integer number'
        );
    });
});

describe('BigDecimal.toDoubleNumber', () => {
    it.each([
        ['0', 0],
        ['-1.25', -1.25],
        ['1.25', 1.25],
        ['5e-324', Number.MIN_VALUE],
        ['1.7976931348623157e+308', Number.MAX_VALUE],
        ['Infinity', Number.POSITIVE_INFINITY],
        ['-Infinity', Number.NEGATIVE_INFINITY],
    ])('should convert %s to the expected double', (value, expected) => {
        expect(parseBigDecimal(value).toDoubleNumber()).toBe(expected);
    });

    it('should preserve negative zero', () => {
        expect(Object.is(parseBigDecimal('-0').toDoubleNumber(), -0)).toBe(true);
    });

    it('should return NaN for a decimal NaN', () => {
        expect(parseBigDecimal('NaN').toDoubleNumber()).toBeNaN();
    });

    it('should expose double precision loss', () => {
        expect(parseBigDecimal('9007199254740993').toDoubleNumber()).toBe(9007199254740992);
        expect(parseBigDecimal('1.00000000000000000001').toDoubleNumber()).toBe(1);
    });

    it('should overflow and underflow according to JavaScript number boundaries', () => {
        expect(parseBigDecimal('1e+309').toDoubleNumber()).toBe(Number.POSITIVE_INFINITY);
        expect(parseBigDecimal('-1e+309').toDoubleNumber()).toBe(Number.NEGATIVE_INFINITY);
        expect(parseBigDecimal('1e-325').toDoubleNumber()).toBe(0);
        expect(Object.is(parseBigDecimal('-1e-325').toDoubleNumber(), -0)).toBe(true);
    });
});

describe('BigDecimal.toString', () => {
    it('should serialize without exponential notation', () => {
        expect(parseBigDecimal('1e+25').toString()).toBe('10000000000000000000000000');
        expect(parseBigDecimal('1e-25').toString()).toBe('0.0000000000000000000000001');
    });

    it.each([
        ['NaN', 'NaN'],
        ['Infinity', 'Infinity'],
        ['-Infinity', '-Infinity'],
        ['-0', '0'],
    ])('should serialize special value %s as %s', (value, expected) => {
        expect(parseBigDecimal(value).toString()).toBe(expected);
    });
});