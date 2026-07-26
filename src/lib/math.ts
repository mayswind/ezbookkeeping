import { reversed } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from './numeral.ts';

export function min(value1: BigDecimal | undefined | null, value2: BigDecimal | undefined | null): BigDecimal {
    if (value1 && value2) {
        if (value1.lessThanOrEqual(value2)) {
            return value1;
        } else {
            return value2;
        }
    } else if (!value1 && value2) {
        return value2;
    } else if (!value2 && value1) {
        return value1;
    } else {
        return BIG_DECIMAL_ZERO;
    }
}

export function max(value1: BigDecimal | undefined | null, value2: BigDecimal | undefined | null): BigDecimal {
    if (value1 && value2) {
        if (value1.greaterThanOrEqual(value2)) {
            return value1;
        } else {
            return value2;
        }
    } else if (!value1 && value2) {
        return value2;
    } else if (!value2 && value1) {
        return value1;
    } else {
        return BIG_DECIMAL_ZERO;
    }
}

export function mean<T>(values: T[], valueFn: (item: T) => BigDecimal): BigDecimal {
    if (values.length < 1) {
        return BIG_DECIMAL_ZERO;
    }

    let sum: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of values) {
        sum = sum.add(valueFn(item));
    }

    return sum.divide(values.length);
}

export function median<T>(sortedValues: T[], valueFn: (item: T) => BigDecimal): BigDecimal {
    if (sortedValues.length < 1) {
        return BIG_DECIMAL_ZERO;
    }

    const mid: number = Math.floor(sortedValues.length / 2);

    if (sortedValues.length % 2 === 0) {
        return valueFn(sortedValues[mid - 1] as T).add(valueFn(sortedValues[mid] as T)).divide(2); // (value1 + value2) / 2
    } else {
        return valueFn(sortedValues[mid] as T);
    }
}

export function percentile<T>(sortedValues: T[], percentile: number, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (sortedValues.length < 1 || percentile < 0 || percentile > 1) {
        return BIG_DECIMAL_ZERO;
    }

    const index: number = (sortedValues.length - 1) * percentile + 1;
    const indexFloor: number = Math.floor(index);
    const indexCeil: number = Math.ceil(index);

    if (indexFloor === indexCeil) {
        return valueFn(sortedValues[indexFloor - 1] as T);
    } else {
        const value1: BigDecimal = valueFn(sortedValues[indexFloor - 1] as T);
        const value2: BigDecimal = valueFn(sortedValues[indexCeil - 1] as T);
        return value1.add(value2.subtract(value1).multiply(index - indexFloor)); // value1 + (value2 - value1) * (index - indexFloor)
    }
}

export function sumMaxN<T>(sortedValues: T[], n: number, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (sortedValues.length < 1 || n <= 0) {
        return BIG_DECIMAL_ZERO;
    }

    let sum: BigDecimal = BIG_DECIMAL_ZERO;
    const count: number = Math.min(n, sortedValues.length);
    const startIndex: number = sortedValues.length - count;

    for (let i = sortedValues.length - 1; i >= startIndex; i--) {
        sum = sum.add(valueFn(sortedValues[i] as T));
    }

    return sum;
}

export function cumulativePercentage<T>(sortedValues: T[], percentageThreshold: number, totalValue: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (sortedValues.length < 1 || percentageThreshold < 0 || percentageThreshold > 1) {
        return BIG_DECIMAL_ZERO;
    }

    const thresholdValue: BigDecimal = totalValue.multiply(percentageThreshold);
    let cumulativeValue: BigDecimal = BIG_DECIMAL_ZERO;
    let cumulativeCount: number = 0;

    for (const item of reversed(sortedValues)) {
        cumulativeValue = cumulativeValue.add(valueFn(item));
        cumulativeCount++;

        if (cumulativeValue.greaterThanOrEqual(thresholdValue)) {
            return parseBigDecimal(cumulativeCount).divide(sortedValues.length).multiply(100); // (cumulativeCount / totalCount) * 100
        }
    }

    return BIG_DECIMAL_ZERO;
}

export function meanAbsoluteDeviation<T>(values: T[], meanValue: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (values.length < 1) {
        return BIG_DECIMAL_ZERO;
    }

    let sumOfAbsoluteDifferences: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of values) {
        const difference: BigDecimal = valueFn(item).subtract(meanValue).abs(); // abs(value - mean)
        sumOfAbsoluteDifferences = sumOfAbsoluteDifferences.add(difference);
    }

    return sumOfAbsoluteDifferences.divide(values.length);
}

export function medianAbsoluteDeviation<T>(sortedValues: T[], medianValue: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (sortedValues.length < 1) {
        return BIG_DECIMAL_ZERO;
    }

    const absoluteDeviations: BigDecimal[] = sortedValues.map(item => valueFn(item).subtract(medianValue).abs()); // abs(value - median)
    absoluteDeviations.sort((a, b) => a.compareTo(b));

    return median(absoluteDeviations, x => x);
}

export function varianceAndStandardDeviation<T>(values: T[], meanValue: BigDecimal, valueFn: (item: T) => BigDecimal): { variance: BigDecimal; standardDeviation: BigDecimal } {
    if (values.length < 1) {
        return { variance: BIG_DECIMAL_ZERO, standardDeviation: BIG_DECIMAL_ZERO };
    }

    let sumOfSquaredDifferences: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of values) {
        const difference: BigDecimal = valueFn(item).subtract(meanValue); // (value - mean)
        sumOfSquaredDifferences = sumOfSquaredDifferences.add(difference.pow(2)); // (value - mean)^2
    }

    const variance: BigDecimal = sumOfSquaredDifferences.divide(values.length); // sumOfSquaredDifferences / n
    const standardDeviation: BigDecimal = variance.sqrt();

    return { variance, standardDeviation };
}

export function coefficientOfVariation(standardDeviation: BigDecimal, meanValue: BigDecimal): BigDecimal | undefined {
    if (meanValue.isZero()) {
        return undefined;
    }

    return standardDeviation.divide(meanValue); // standardDeviation / meanValue
}

export function skewness<T>(values: T[], meanValue: BigDecimal, standardDeviation: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (values.length < 1 || standardDeviation.isZero()) {
        return BIG_DECIMAL_ZERO;
    }

    let sumOfCubedDifferences: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of values) {
        const difference: BigDecimal = valueFn(item).subtract(meanValue); // (value - mean)
        sumOfCubedDifferences = sumOfCubedDifferences.add(difference.pow(3)); // (value - mean)^3
    }

    return sumOfCubedDifferences.divide(standardDeviation.pow(3).multiply(values.length)); // sumOfCubedDifferences / (n * standardDeviation^3)
}

export function kurtosis<T>(values: T[], meanValue: BigDecimal, variance: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (values.length < 1 || variance.isZero()) {
        return BIG_DECIMAL_ZERO;
    }

    let sumOfQuarticDifferences: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of values) {
        const difference: BigDecimal = valueFn(item).subtract(meanValue); // (value - mean)
        sumOfQuarticDifferences = sumOfQuarticDifferences.add(difference.pow(4)); // (value - mean)^4
    }

    return sumOfQuarticDifferences.divide(variance.pow(2).multiply(values.length)); // sumOfQuarticDifferences / (n * variance^2)
}

export function giniCoefficient<T>(sortedValues: T[], totalValue: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (sortedValues.length < 1 || totalValue.isZero()) {
        return BIG_DECIMAL_ZERO;
    }

    const n: number = sortedValues.length;
    let weightedSum: BigDecimal = BIG_DECIMAL_ZERO;

    for (let i = 0; i < n; i++) {
        weightedSum = weightedSum.add(valueFn(sortedValues[i] as T).multiply(i + 1)); // (i + 1) * value
    }

    return weightedSum.multiply(2).divide(n).divide(totalValue).subtract(parseBigDecimal(n + 1).divide(n)); // (2 * weightedSum / (n * totalValue)) - ((n + 1) / n)
}

export function herfindahlHirschmanIndex<T>(values: T[], totalValue: BigDecimal, valueFn: (item: T) => BigDecimal): BigDecimal {
    if (values.length < 1 || totalValue.isZero()) {
        return BIG_DECIMAL_ZERO;
    }

    let hhi: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of values) {
        const share: BigDecimal = valueFn(item).divide(totalValue); // item / totalValue
        hhi = hhi.add(share.pow(2));
    }

    return hhi;
}
