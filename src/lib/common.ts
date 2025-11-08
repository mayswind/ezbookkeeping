import { keys, keysIfValueEquals, values } from '@/core/base.ts';
import type { NameValue, TypeAndName, TypeAndDisplayName} from '@/core/base.ts';

// eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
export function isFunction(val: unknown): val is Function {
    return typeof(val) === 'function';
}

export function isDefined<T>(val: T | null | undefined): val is T {
    return val !== null && typeof(val) !== 'undefined';
}

export function isObject(val: unknown): val is object {
    return val !== null && typeof(val) === 'object' && !isArray(val);
}

export function isArray<T>(val: unknown): val is T[] {
    if (isFunction(Array.isArray)) {
        return Array.isArray(val);
    }

    return Object.prototype.toString.call(val) === '[object Array]';
}

export function isString(val: unknown): val is string {
    return typeof(val) === 'string';
}

export function isNumber(val: unknown): val is number {
    return typeof(val) === 'number';
}

export function isInteger(val: unknown): val is number {
    return Number.isInteger(val);
}

export function isBoolean(val: unknown): val is boolean {
    return typeof(val) === 'boolean';
}

export function isYearMonth(val: unknown): val is string {
    if (!isString(val)) {
        return false;
    }

    const items = val.split('-');

    if (items.length !== 2) {
        return false;
    }

    return !!parseInt(items[0] as string) && !!parseInt(items[1] as string);
}

export function isEquals<T>(obj1: T, obj2: T): boolean {
    if (obj1 === obj2) {
        return true;
    }

    if (isArray(obj1) && isArray(obj2)) {
        const arr1 = obj1;
        const arr2 = obj2;

        if (arr1.length !== arr2.length) {
            return false;
        }

        for (let i = 0; i < arr1.length; i++) {
            if (!isEquals(arr1[i], arr2[i])) {
                return false;
            }
        }

        return true;
    } else if (isObject(obj1) && isObject(obj2)) {
        const keys1 = Object.keys(obj1);
        const keys2 = Object.keys(obj2);

        if (keys1.length !== keys2.length) {
            return false;
        }

        const keyExistsMap2 = new Map<string, boolean>();

        for (let i = 0; i < keys2.length; i++) {
            const key = keys2[i] as string;

            keyExistsMap2.set(key, true);
        }

        for (let i = 0; i < keys1.length; i++) {
            const key = keys1[i] as string;

            if (!keyExistsMap2.get(key)) {
                return false;
            }

            if (!isEquals((obj1 as Record<string, unknown>)[key], (obj2 as Record<string, unknown>)[key])) {
                return false;
            }
        }

        return true;
    } else {
        return obj1 === obj2;
    }
}

export function isYearMonthEquals(val1: string, val2: string): boolean {
    const items1 = val1.split('-');
    const items2 = val2.split('-');

    if (items1.length !== 2 || items2.length !== 2) {
        return false;
    }

    return (!!parseInt(items1[0] as string) && !!parseInt(items1[1] as string)) && (parseInt(items1[0] as string) === parseInt(items2[0] as string)) && (parseInt(items1[1] as string) === parseInt(items2[1] as string));
}

export function isArray1SubsetOfArray2<T>(array1: T[], array2: T[]): boolean {
    if (array1.length > array2.length) {
        return false;
    }

    const array2ValuesMap: Map<T, boolean> = new Map<T, boolean>();

    for (let i = 0; i < array2.length; i++) {
        array2ValuesMap.set(array2[i] as T, true);
    }

    for (let i = 0; i < array1.length; i++) {
        if (!array2ValuesMap.get(array1[i] as T)) {
            return false;
        }
    }

    return true;
}

export function isObjectEmpty(obj: object): boolean {
    if (!obj) {
        return true;
    }

    for (const _ of keys(obj)) {
        return false;
    }

    return true;
}

export function ofObject<T>(object: T): T {
    return object;
}

export function getNumberValue(value: unknown, defaultValue: number): number {
    if (isString(value)) {
        return parseInt(value, 10);
    } else if (isNumber(value)) {
        return value;
    } else {
        return defaultValue;
    }
}

export function sortNumbersArray(array: number[]): number[] {
    return array.sort(function (num1, num2) {
        return num1 - num2;
    });
}

export function getObjectOwnFieldCount(object: object): number {
    let count = 0;

    if (!object || !isObject(object)) {
        return count;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    for (const _ of keys(object)) {
        count++;
    }

    return count;
}

export function replaceAll(value: string, originalValue: string, targetValue: string): string {
    // Escape special characters in originalValue to safely use it in a regex pattern.
    // This ensures that characters like . (dot), * (asterisk), +, ?, etc. are treated literally,
    // rather than as special regex symbols.
    const escapedOriginalValue = originalValue.replace(/([.*+?^=!:${}()|\-/\\])/g, '\\$1');

    return value.replace(new RegExp(escapedOriginalValue, 'g'), targetValue);
}

export function removeAll(value: string, originalValue: string): string {
    return replaceAll(value, originalValue, '');
}

export function limitText(value: string, maxLength: number): string {
    let length = 0;

    for (let i = 0; i < value.length; i++) {
        const c = value.charCodeAt(i);

        if ((c >= 0x0001 && c <= 0x007e) || (0xff60 <= c && c <= 0xff9f)) {
            length++;
        } else {
            length += 2;
        }
    }

    if (length <= maxLength || maxLength <= 3) {
        return value;
    }

    return value.substring(0, maxLength - 3) + '...';
}

export function base64encode(arrayBuffer: ArrayBuffer): string {
    return btoa(String.fromCharCode.apply(null, Array.from(new Uint8Array(arrayBuffer))));
}

export function base64decode(str: string): string {
    if (!str) {
        return '';
    }

    return atob(str);
}

export function arrayBufferToString(arrayBuffer: ArrayBuffer): string {
    return String.fromCharCode.apply(null, Array.from(new Uint8Array(arrayBuffer)));
}

export function stringToArrayBuffer(str: string): ArrayBuffer {
    return Uint8Array.from(str, c => c.charCodeAt(0)).buffer;
}

export function getFirstVisibleItem<T>(items: Record<string, T>[] | Record<string, Record<string, T>>, hiddenField: string): Record<string, T> | null {
    if (isArray(items) && items.length > 0) {
        const arr = items as Record<string, T>[];

        for (const item of arr) {
            if (hiddenField && item[hiddenField]) {
                continue;
            }

            return item;
        }
    } else if (isObject(items)) {
        const obj = items as Record<string, Record<string, T>>;

        for (const item of values(obj)) {
            if (hiddenField && item[hiddenField]) {
                continue;
            }

            return item;
        }
    }

    return null;
}

export function getItemByKeyValue<T>(src: Record<string, T>[] | Record<string, Record<string, T>>, value: T, keyField: string): Record<string, T> | null {
    if (isArray(src)) {
        const arr = src as Record<string, T>[];

        for (const item of arr) {
            if (item[keyField] === value) {
                return item;
            }
        }
    } else if (isObject(src)) {
        const obj = src as Record<string, Record<string, T>>;

        for (const item of values(obj)) {
            if (item[keyField] === value) {
                return item;
            }
        }
    }

    return null;
}

export function findNameByValue(items: NameValue[], value: string): string | null {
    for (const item of items) {
        if (item.value === value) {
            return item.name;
        }
    }

    return null;
}

export function findNameByType(items: TypeAndName[], type: number): string | null {
    for (const item of items) {
        if (item.type === type) {
            return item.name;
        }
    }

    return null;
}

export function findDisplayNameByType(items: TypeAndDisplayName[], type: number): string | null {
    for (const item of items) {
        if (item.type === type) {
            return item.displayName;
        }
    }

    return null;
}

export function getNameByKeyValue<V, N>(src: Record<string, N | V>[] | Record<string, Record<string, N | V>>, value: V, keyField: string | null, nameField: string, defaultName?: N): N | undefined {
    if (isArray(src)) {
        const arr = src as Record<string, N | V>[];

        if (keyField) {
            for (const option of arr) {
                if (option[keyField] === value) {
                    return option[nameField] as N;
                }
            }
        } else if (isNumber(value)) {
            const index = value;

            if (arr[index]) {
                const option = arr[index];

                return option[nameField] as N;
            }
        }
    } else if (isObject(src)) {
        const obj = src as Record<string, Record<string, N | V>>;

        if (keyField) {
            for (const option of values(obj)) {
                if (option[keyField] === value) {
                    return option[nameField] as N;
                }
            }
        } else if (isString(value)) {
            const key = value;

            if (obj[key]) {
                const option = obj[key];

                return option[nameField] as N;
            }
        }
    }

    return defaultName;
}

export function arrayContainsFieldValue<T>(array: Record<string, T>[], fieldName: string, value: T): boolean {
    if (!value || !array || !array.length) {
        return false;
    }

    for (const item of array) {
        if (item[fieldName] === value) {
            return true;
        }
    }

    return false;
}

export function objectFieldToArrayItem(object: object): string[] {
    const ret: string[] = [];

    for (const field of keys(object)) {
        ret.push(field);
    }

    return ret;
}

export function objectFieldWithValueToArrayItem<T>(object: Record<string, T>, value: T): string[] {
    const ret: string[] = [];

    for (const field of keysIfValueEquals(object, value)) {
        ret.push(field);
    }

    return ret;
}

export function arrayItemToObjectField<T>(array: string[], value: T): Record<string, T> {
    const ret: Record<string, T> = {};

    for (const item of array) {
        ret[item] = value;
    }

    return ret;
}

export function splitItemsToMap(str: string | undefined | null, separator: string): Record<string, boolean> {
    const ret: Record<string, boolean> = {};

    if (!str) {
        return ret;
    }

    const items = str.split(separator);

    for (const item of items) {
        if (item) {
            ret[item] = true;
        }
    }

    return ret;
}

export function countSplitItems(str: string | undefined | null, separator: string): number {
    if (!str) {
        return 0;
    }

    const items = str.split(separator);
    let count = 0;

    for (const item of items) {
        if (item) {
            count++;
        }
    }

    return count;
}

export function categorizedArrayToPlainArray<T>(object: Record<string, T[]>): T[] {
    const ret: T[] = [];

    for (const array of values(object)) {
        for (const item of array) {
            ret.push(item);
        }
    }

    return ret;
}

export function selectAll(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string } }): void {
    for (const itemId of keys(filterItemIds)) {
        const item = allItemsMap[itemId];

        if (item) {
            filterItemIds[item.id] = false;
        }
    }
}

export function selectNone(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string } }): void {
    for (const itemId of keys(filterItemIds)) {
        const item = allItemsMap[itemId];

        if (item) {
            filterItemIds[item.id] = true;
        }
    }
}

export function selectInvert(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string } }): void {
    for (const itemId of keys(filterItemIds)) {
        const item = allItemsMap[itemId];

        if (item) {
            filterItemIds[item.id] = !filterItemIds[item.id];
        }
    }
}

export function selectAllVisible(filterItemIds: Record<string, boolean>, allItemsMap: { [key: string]: { id: string, hidden?: boolean } }): void {
    for (const itemId of keys(filterItemIds)) {
        const item = allItemsMap[itemId];

        if (item && !item.hidden) {
            filterItemIds[item.id] = false;
        }
    }
}

export function isPrimaryItemHasSecondaryValue(primaryItem: Record<string, Record<string, unknown>[]>, primarySubItemsField: string, secondaryValueField: string | undefined, secondaryHiddenField: string | undefined, secondaryValue: unknown): boolean {
    const secondaryItems = primaryItem[primarySubItemsField];

    if (!secondaryItems || secondaryItems.length < 1) {
        return false;
    }

    for (const secondaryItem of secondaryItems) {
        if (secondaryHiddenField && secondaryItem[secondaryHiddenField]) {
            continue;
        }

        if (secondaryValueField && secondaryItem[secondaryValueField] === secondaryValue) {
            return true;
        } else if (!secondaryValueField && secondaryItem === secondaryValue) {
            return true;
        }
    }

    return false;
}

export function getPrimaryValueBySecondaryValue<T>(items: Record<string, Record<string, T>[]>[] | Record<string, Record<string, Record<string, T>[]>>, primarySubItemsField: string | undefined, primaryValueField: string | undefined, primaryHiddenField: string | undefined, secondaryValueField: string | undefined, secondaryHiddenField: string | undefined, secondaryValue: T): Record<string, T>[] | Record<string, Record<string, T>[]> | null | undefined {
    if (primarySubItemsField) {
        if (isArray(items)) {
            const arr = items as Record<string, Record<string, T>[]>[];

            for (const primaryItem of arr) {
                if (primaryHiddenField && primaryItem[primaryHiddenField]) {
                    continue;
                }

                if (isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryHiddenField, secondaryValue)) {
                    if (primaryValueField) {
                        return primaryItem[primaryValueField];
                    } else {
                        return primaryItem;
                    }
                }
            }
        } else {
            const obj = items as Record<string, Record<string, Record<string, T>[]>>;

            for (const primaryItem of values(obj)) {
                if (primaryHiddenField && primaryItem[primaryHiddenField]) {
                    continue;
                }

                if (isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryHiddenField, secondaryValue)) {
                    if (primaryValueField) {
                        return primaryItem[primaryValueField];
                    } else {
                        return primaryItem;
                    }
                }
            }
        }
    }

    return null;
}

export function arrangeArrayWithNewStartIndex<T>(array: T[], startIndex: number): T[] {
    if (startIndex <= 0 || startIndex >= array.length) {
        return array;
    }

    const newArray: T[] = [];

    for (let i = startIndex; i < array.length; i++) {
        newArray.push(array[i] as T);
    }

    for (let i = 0; i < startIndex; i++) {
        newArray.push(array[i] as T);
    }

    return newArray;
}
