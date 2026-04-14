import { describe, expect, test } from '@jest/globals';

import { mean, median, percentile, sumMaxN } from '@/lib/math.ts';

type TestNumberItem = {
	value: number;
};

function createNumberItems(values: number[]): TestNumberItem[] {
	return values.map(value => ({ value }));
}

function getTestTitle(functionName: string, title: string): string {
	return `${functionName}: ${title}`;
}

type MeanTestCase = {
	title: string;
	values: number[];
	expected: number;
};

const TEST_CASES_MEAN: MeanTestCase[] = [
	{
		title: 'returns zero for empty array',
		values: [],
		expected: 0,
	},
	{
		title: 'returns the average for positive values',
		values: [1, 2, 3, 4],
		expected: 2.5,
	},
	{
		title: 'returns the average for negative and positive values',
		values: [-10, 0, 20],
		expected: 10 / 3,
	},
];

describe('mean', () => {
	TEST_CASES_MEAN.forEach((testCase) => {
		test(getTestTitle('mean', testCase.title), () => {
			const result = mean(createNumberItems(testCase.values), item => item.value);
			expect(result).toBeCloseTo(testCase.expected);
		});
	});
});

type MedianTestCase = {
	title: string;
	values: number[];
	expected: number;
};

const TEST_CASES_MEDIAN: MedianTestCase[] = [
	{
		title: 'returns zero for empty sorted array',
		values: [],
		expected: 0,
	},
	{
		title: 'returns the middle value for odd-length sorted array',
		values: [1, 3, 5],
		expected: 3,
	},
	{
		title: 'returns the average of the two middle values for even-length sorted array',
		values: [1, 3, 5, 7],
		expected: 4,
	},
];

describe('median', () => {
	TEST_CASES_MEDIAN.forEach((testCase) => {
		test(getTestTitle('median', testCase.title), () => {
			const result = median(createNumberItems(testCase.values), item => item.value);
			expect(result).toBeCloseTo(testCase.expected);
		});
	});
});

type PercentileTestCase = {
	title: string;
	values: number[];
	percentileValue: number;
	expected: number;
};

const TEST_CASES_PERCENTILE: PercentileTestCase[] = [
	{
		title: 'returns zero for empty sorted array',
		values: [],
		percentileValue: 0.5,
		expected: 0,
	},
	{
		title: 'returns zero when percentile is smaller than zero',
		values: [1, 2, 3],
		percentileValue: -0.1,
		expected: 0,
	},
	{
		title: 'returns zero when percentile is larger than one',
		values: [1, 2, 3],
		percentileValue: 1.1,
		expected: 0,
	},
	{
		title: 'returns the minimum value for zero percentile',
		values: [5, 10, 15, 20],
		percentileValue: 0,
		expected: 5,
	},
	{
		title: 'returns the maximum value for one percentile',
		values: [5, 10, 15, 20],
		percentileValue: 1,
		expected: 20,
	},
	{
		title: 'returns the exact indexed value when percentile maps to an integer index',
		values: [10, 20, 30, 40, 50],
		percentileValue: 0.25,
		expected: 20,
	},
	{
		title: 'interpolates between neighboring values when percentile maps to a fractional index',
		values: [10, 20, 30, 40, 50, 60, 70, 80],
		percentileValue: 0.25,
		expected: 27.5,
	},
];

describe('percentile', () => {
	TEST_CASES_PERCENTILE.forEach((testCase) => {
		test(getTestTitle('percentile', testCase.title), () => {
			const result = percentile(
				createNumberItems(testCase.values),
				testCase.percentileValue,
				item => item.value
			);

			expect(result).toBeCloseTo(testCase.expected);
		});
	});
});

type SumMaxNTestCase = {
	title: string;
	values: number[];
	n: number;
	expected: number;
};

const TEST_CASES_SUM_MAX_N: SumMaxNTestCase[] = [
	{
		title: 'returns zero for empty sorted array',
		values: [],
		n: 3,
		expected: 0,
	},
	{
		title: 'returns zero when n is zero',
		values: [1, 2, 3],
		n: 0,
		expected: 0,
	},
	{
		title: 'returns the sum of the largest n values',
		values: [1, 2, 3, 4, 5],
		n: 2,
		expected: 9,
	},
	{
		title: 'returns the sum of all values when n is larger than array length',
		values: [1, 2, 3, 4],
		n: 10,
		expected: 10,
	},
];

describe('sumMaxN', () => {
	TEST_CASES_SUM_MAX_N.forEach((testCase) => {
		test(getTestTitle('sumMaxN', testCase.title), () => {
			const result = sumMaxN(createNumberItems(testCase.values), testCase.n, item => item.value);
			expect(result).toBe(testCase.expected);
		});
	});
});
