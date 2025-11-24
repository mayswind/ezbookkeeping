// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type PartialRecord<K extends keyof any, T> = {
    [P in K]?: T;
}

export function* itemAndIndex<T>(arr: T[]): Iterable<[T, number]> {
    for (let i = 0; i < arr.length; i++) {
        yield [arr[i], i] as [T, number];
    }
}

export function* reversed<T>(arr: T[]): Iterable<T> {
    for (let i = arr.length - 1; i >= 0; i--) {
        yield arr[i] as T;
    }
}

export function* reversedItemAndIndex<T>(arr: T[]): Iterable<[T, number]> {
    for (let i = arr.length - 1; i >= 0; i--) {
        yield [arr[i], i] as [T, number];
    }
}

export function* entries<K extends string | number | symbol, V>(obj: Record<K, V>): Iterable<[string, V]> {
    for (const key in obj) {
        if (!Object.prototype.hasOwnProperty.call(obj, key)) {
            continue;
        }

        yield [key, obj[key]] as [string, V];
    }
}

export function* keys<K extends string | number | symbol, V>(obj: Record<K, V>): Iterable<string> {
    for (const key in obj) {
        if (!Object.prototype.hasOwnProperty.call(obj, key)) {
            continue;
        }

        yield key;
    }
}

export function* keysIfValueEquals<K extends string | number | symbol, V>(obj: Record<K, V>, value: V): Iterable<string> {
    for (const key in obj) {
        if (!Object.prototype.hasOwnProperty.call(obj, key)) {
            continue;
        }

        if (obj[key] !== value) {
            continue;
        }

        yield key;
    }
}

export function* values<K extends string | number | symbol, V>(obj: Record<K, V>): Iterable<V> {
    for (const key in obj) {
        if (!Object.prototype.hasOwnProperty.call(obj, key)) {
            continue;
        }

        yield obj[key] as V;
    }
}

export interface NameValue {
    readonly name: string;
    readonly value: string;
}

export interface NameNumeralValue {
    readonly name: string;
    readonly value: number;
}

export interface KeyAndName {
    readonly key: string;
    readonly name: string;
}

export interface TypeAndName {
    readonly type: number;
    readonly name: string;
}

export interface TypeAndDisplayName {
    readonly type: number;
    readonly displayName: string;
}

export interface LocalizedSwitchOption {
    readonly value: boolean;
    readonly displayName: string;
}

export type BeforeResolveFunction = (callback: () => void) => void;
