import { reversed } from '@/core/base.ts';

export function mean<T>(values: T[], valueFn: (item: T) => number): number {
    if (values.length < 1) {
        return 0;
    }

    let sum: number = 0;

    for (const item of values) {
        sum += valueFn(item);
    }

    return sum / values.length;
}

export function median<T>(sortedValues: T[], valueFn: (item: T) => number): number {
    if (sortedValues.length < 1) {
        return 0;
    }

    const mid: number = Math.floor(sortedValues.length / 2);

    if (sortedValues.length % 2 === 0) {
        return (valueFn(sortedValues[mid - 1] as T) + valueFn(sortedValues[mid] as T)) / 2;
    } else {
        return valueFn(sortedValues[mid] as T);
    }
}

export function percentile<T>(sortedValues: T[], percentile: number, valueFn: (item: T) => number): number {
    if (sortedValues.length < 1 || percentile < 0 || percentile > 1) {
        return 0;
    }

    const index: number = (sortedValues.length - 1) * percentile + 1;
    const indexFloor: number = Math.floor(index);
    const indexCeil: number = Math.ceil(index);

    if (indexFloor === indexCeil) {
        return valueFn(sortedValues[indexFloor - 1] as T);
    } else {
        const value1: number = valueFn(sortedValues[indexFloor - 1] as T);
        const value2: number = valueFn(sortedValues[indexCeil - 1] as T);
        return value1 + (index - indexFloor) * (value2 - value1);
    }
}

export function sumMaxN<T>(sortedValues: T[], n: number, valueFn: (item: T) => number): number {
    if (sortedValues.length < 1 || n <= 0) {
        return 0;
    }

    let sum: number = 0;
    const count: number = Math.min(n, sortedValues.length);
    const startIndex: number = sortedValues.length - count;

    for (let i = sortedValues.length - 1; i >= startIndex; i--) {
        sum += valueFn(sortedValues[i] as T);
    }

    return sum;
}

export function cumulativePercentage<T>(sortedValues: T[], percentageThreshold: number, totalValue: number, valueFn: (item: T) => number): number {
    if (sortedValues.length < 1 || percentageThreshold < 0 || percentageThreshold > 1) {
        return 0;
    }

    const thresholdValue: number = percentageThreshold * totalValue;
    let cumulativeValue: number = 0;
    let cumulativeCount: number = 0;

    for (const item of reversed(sortedValues)) {
        cumulativeValue += valueFn(item);
        cumulativeCount++;

        if (cumulativeValue >= thresholdValue) {
            return 100.0 * cumulativeCount / sortedValues.length;
        }
    }

    return 0;
}

export function meanAbsoluteDeviation<T>(values: T[], meanValue: number, valueFn: (item: T) => number): number {
    if (values.length < 1) {
        return 0;
    }

    let sumOfAbsoluteDifferences: number = 0;

    for (const item of values) {
        const difference: number = Math.abs(valueFn(item) - meanValue);
        sumOfAbsoluteDifferences += difference;
    }

    return sumOfAbsoluteDifferences / values.length;
}

export function medianAbsoluteDeviation<T>(sortedValues: T[], medianValue: number, valueFn: (item: T) => number): number {
    if (sortedValues.length < 1) {
        return 0;
    }

    const absoluteDeviations: number[] = sortedValues.map(item => Math.abs(valueFn(item) - medianValue));
    absoluteDeviations.sort((a, b) => a - b);

    return median(absoluteDeviations, x => x);
}

export function varianceAndStandardDeviation<T>(values: T[], meanValue: number, valueFn: (item: T) => number): { variance: number; standardDeviation: number } {
    if (values.length < 1) {
        return { variance: 0, standardDeviation: 0 };
    }

    let sumOfSquaredDifferences: number = 0;

    for (const item of values) {
        const difference: number = valueFn(item) - meanValue;
        sumOfSquaredDifferences += difference * difference;
    }

    const variance: number = sumOfSquaredDifferences / values.length;
    const standardDeviation: number = Math.sqrt(variance);

    return { variance, standardDeviation };
}

export function coefficientOfVariation(standardDeviation: number, meanValue: number): number | undefined {
    if (meanValue === 0) {
        return undefined;
    }

    return standardDeviation / meanValue;
}

export function skewness<T>(values: T[], meanValue: number, standardDeviation: number, valueFn: (item: T) => number): number {
    if (values.length < 1 || standardDeviation === 0) {
        return 0;
    }

    let sumOfCubedDifferences: number = 0;

    for (const item of values) {
        const difference: number = valueFn(item) - meanValue;
        sumOfCubedDifferences += Math.pow(difference, 3);
    }

    return sumOfCubedDifferences / (values.length * Math.pow(standardDeviation, 3));
}

export function kurtosis<T>(values: T[], meanValue: number, variance: number, valueFn: (item: T) => number): number {
    if (values.length < 1 || variance === 0) {
        return 0;
    }

    let sumOfQuarticDifferences: number = 0;

    for (const item of values) {
        const difference: number = valueFn(item) - meanValue;
        sumOfQuarticDifferences += Math.pow(difference, 4);
    }

    return sumOfQuarticDifferences / (values.length * Math.pow(variance, 2));
}
