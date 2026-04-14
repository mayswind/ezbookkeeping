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
