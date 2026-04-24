import { describe, expect, it } from 'vitest';

import { mean, median, percentile, sumMaxN } from '@/lib/math.ts';

describe('mean', () => {
	it('should return zero for empty array', () => {
		expect(mean([], item => item)).toBeCloseTo(0);
	});

	it('should return the average for positive values', () => {
		expect(mean([1, 2, 3, 4], item => item)).toBeCloseTo(2.5);
	});

	it('should return the average for negative and positive values', () => {
		expect(mean([-10, 0, 20], item => item)).toBeCloseTo(10 / 3);
	});
});

describe('median', () => {
	it('should return zero for empty sorted array', () => {
		expect(median([], item => item)).toBeCloseTo(0);
	});

	it('should return the middle value for odd-length sorted array', () => {
		expect(median([1, 3, 5], item => item)).toBeCloseTo(3);
	});

	it('should return the average of the two middle values for even-length sorted array', () => {
		expect(median([1, 3, 5, 7], item => item)).toBeCloseTo(4);
	});
});

describe('percentile', () => {
	it('should return zero for empty sorted array', () => {
		expect(percentile([], 0.5, item => item)).toBeCloseTo(0);
	});

	it('should return zero when percentile is smaller than zero', () => {
		expect(percentile([1, 2, 3], -0.1, item => item)).toBeCloseTo(0);
	});

	it('should return zero when percentile is larger than one', () => {
		expect(percentile([1, 2, 3], 1.1, item => item)).toBeCloseTo(0);
	});

	it('should return the minimum value for zero percentile', () => {
		expect(percentile([5, 10, 15, 20], 0, item => item)).toBeCloseTo(5);
	});

	it('should return the maximum value for one percentile', () => {
		expect(percentile([5, 10, 15, 20], 1, item => item)).toBeCloseTo(20);
	});

	it('should return the exact indexed value when percentile maps to an integer index', () => {
		expect(percentile([10, 20, 30, 40, 50], 0.25, item => item)).toBeCloseTo(20);
	});

	it('should interpolate between neighboring values when percentile maps to a fractional index', () => {
		expect(percentile([10, 20, 30, 40, 50, 60, 70, 80], 0.25, item => item)).toBeCloseTo(27.5);
	});
});

describe('sumMaxN', () => {
	it('should return zero for empty sorted array', () => {
		expect(sumMaxN([], 3, item => item)).toBe(0);
	});

	it('should return zero when n is zero', () => {
		expect(sumMaxN([1, 2, 3], 0, item => item)).toBe(0);
	});

	it('should return the sum of the largest n values', () => {
		expect(sumMaxN([1, 2, 3, 4, 5], 2, item => item)).toBe(9);
	});

	it('should return the sum of all values when n is larger than array length', () => {
		expect(sumMaxN([1, 2, 3, 4], 10, item => item)).toBe(10);
	});
});
