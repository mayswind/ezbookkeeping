export function isFunction(val) {
    return typeof(val) === 'function';
}

export function isObject(val) {
    return val != null && typeof(val) === 'object' && !isArray(val);
}

export function isArray(val) {
    if (isFunction(Array.isArray)) {
        return Array.isArray(val);
    }

    return Object.prototype.toString.call(val) === '[object Array]';
}

export function isString(val) {
    return typeof(val) === 'string';
}

export function isNumber(val) {
    return typeof(val) === 'number';
}

export function isBoolean(val) {
    return typeof(val) === 'boolean';
}

export function isYearMonth(val) {
    if (typeof(val) !== 'string') {
        return false;
    }

    const items = val.split('-');

    if (items.length !== 2) {
        return false;
    }

    return isNumber(items[0]) && isNumber(items[1]);
}

export function isEquals(obj1, obj2) {
    if (obj1 === obj2) {
        return true;
    }

    if (isArray(obj1) && isArray(obj2)) {
        if (obj1.length !== obj2.length) {
            return false;
        }

        for (let i = 0; i < obj1.length; i++) {
            if (!isEquals(obj1[i], obj2[i])) {
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

        const keyExistsMap2 = {};

        for (let i = 0; i < keys2.length; i++) {
            const key = keys2[i];

            keyExistsMap2[key] = true;
        }

        for (let i = 0; i < keys1.length; i++) {
            const key = keys1[i];

            if (!keyExistsMap2[key]) {
                return false;
            }

            if (!isEquals(obj1[key], obj2[key])) {
                return false;
            }
        }

        return true;
    } else {
        return obj1 === obj2;
    }
}

export function getObjectOwnFieldCount(object) {
    let count = 0;

    if (!object || !isObject(object)) {
        return count;
    }

    for (let field in object) {
        if (!Object.prototype.hasOwnProperty.call(object, field)) {
            continue;
        }

        count++;
    }

    return count;
}

export function appendThousandsSeparator(value, enable) {
    if (!enable || value.length <= 3) {
        return value;
    }

    const negative = value.charAt(0) === '-';

    if (negative) {
        value = value.substring(1);
    }

    const dotPos = value.indexOf('.');
    const integer = dotPos < 0 ? value : value.substring(0, dotPos);
    const decimals = dotPos < 0 ? '' : value.substring(dotPos + 1, value.length);

    const finalChars = [];

    for (let i = 0; i < integer.length; i++) {
        if (i % 3 === 0 && i > 0) {
            finalChars.push(',');
        }

        finalChars.push(integer.charAt(integer.length - 1 - i));
    }

    finalChars.reverse();

    let newInteger = finalChars.join('');

    if (negative) {
        newInteger = `-${newInteger}`;
    }

    if (dotPos < 0) {
        return newInteger;
    } else {
        return `${newInteger}.${decimals}`;
    }
}

export function formatPercent(value, precision, lowPrecisionValue) {
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.floor(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        return lowPrecisionValue + '%';
    }

    const result = normalizedValue / ratio;
    return result + '%';
}

export function limitText(value, maxLength) {
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

export function getTextBefore(fullText, text) {
    if (!text) {
        return fullText;
    }

    const index = fullText.indexOf(text);

    if (index >= 0) {
        return fullText.substring(0, index);
    }

    return '';
}

export function getTextAfter(fullText, text) {
    if (!text) {
        return fullText;
    }

    let index = fullText.indexOf(text);

    if (index >= 0) {
        index += text.length;
        return fullText.substring(index);
    }

    return '';
}

export function base64encode(arrayBuffer) {
    if (!arrayBuffer || arrayBuffer.length === 0) {
        return null;
    }

    return btoa(String.fromCharCode.apply(null, new Uint8Array(arrayBuffer)));
}

export function base64decode(str) {
    if (!str) {
        return '';
    }

    return atob(str);
}

export function arrayBufferToString(arrayBuffer) {
    return String.fromCharCode.apply(null, new Uint8Array(arrayBuffer));
}

export function stringToArrayBuffer(str){
    return Uint8Array.from(str, c => c.charCodeAt(0)).buffer;
}

export function getFirstItem(items) {
    if (isArray(items)) {
        if (items.length > 0) {
            return items[0];
        }
    } else {
        for (let field in items) {
            if (!Object.prototype.hasOwnProperty.call(items, field)) {
                continue;
            }

            return items[field];
        }
    }

    return null;
}

export function getItemByKeyValue(src, value, keyField) {
    if (isArray(src)) {
        for (let i = 0; i < src.length; i++) {
            const item = src[i];

            if (item[keyField] === value) {
                return item;
            }
        }
    } else {
        for (let field in src) {
            if (!Object.prototype.hasOwnProperty.call(src, field)) {
                continue;
            }

            const item = src[field];

            if (item[keyField] === value) {
                return item;
            }
        }
    }

    return null;
}

export function getNameByKeyValue(src, value, keyField, nameField, defaultName) {
    if (isArray(src)) {
        if (keyField) {
            for (let i = 0; i < src.length; i++) {
                const option = src[i];

                if (option[keyField] === value) {
                    return option[nameField];
                }
            }
        } else {
            if (src[value]) {
                const option = src[value];

                return option[nameField];
            }
        }
    } else if (isObject(src)) {
        if (keyField) {
            for (let key in src) {
                if (!Object.prototype.hasOwnProperty.call(src, key)) {
                    continue;
                }

                const option = src[key];

                if (option[keyField] === value) {
                    return option[nameField];
                }
            }
        } else {
            if (src[value]) {
                const option = src[value];

                return option[nameField];
            }
        }
    }

    return defaultName;
}

export function copyObjectTo(fromObject, toObject) {
    if (!isObject(fromObject)) {
        return toObject;
    }

    if (!isObject(toObject)) {
        toObject = {};
    }

    for (let key in fromObject) {
        if (!Object.prototype.hasOwnProperty.call(fromObject, key)) {
            continue;
        }

        const fromValue = fromObject[key];
        const toValue = toObject[key];

        if (isArray(fromValue)) {
            toObject[key] = copyArrayTo(fromValue, toValue);
        } else if (isObject(fromValue)) {
            toObject[key] = copyObjectTo(fromValue, toValue);
        } else {
            if (fromValue !== toValue) {
                toObject[key] = fromValue;
            }
        }
    }

    return toObject;
}

export function copyArrayTo(fromArray, toArray) {
    if (!isArray(fromArray)) {
        return toArray;
    }

    if (!isArray(toArray)) {
        toArray = [];
    }

    for (let i = 0; i < fromArray.length; i++) {
        const fromValue = fromArray[i];

        if (toArray.length > i) {
            const toValue = toArray[i];

            if (isArray(fromValue)) {
                toArray[i] = copyArrayTo(fromValue, toValue);
            } else if (isObject(fromValue)) {
                toArray[i] = copyObjectTo(fromValue, toValue);
            } else {
                if (fromValue !== toValue) {
                    toArray[i] = fromValue;
                }
            }
        } else {
            if (isArray(fromValue)) {
                toArray.push(copyArrayTo(fromValue, []));
            } else if (isObject(fromValue)) {
                toArray.push(copyObjectTo(fromValue, {}));
            } else {
                toArray.push(fromValue);
            }
        }
    }

    return toArray;
}

export function arrayContainsFieldValue(array, fieldName, value) {
    if (!value || !array || !array.length) {
        return false;
    }

    for (let i = 0; i < array.length; i++) {
        if (array[i][fieldName] === value) {
            return true;
        }
    }

    return false;
}

export function categorizedArrayToPlainArray(object) {
    const ret = [];

    for (let field in object) {
        if (!Object.prototype.hasOwnProperty.call(object, field)) {
            continue;
        }

        const array = object[field];

        for (let i = 0; i < array.length; i++) {
            ret.push(array[i]);
        }
    }

    return ret;
}

export function isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryValue) {
    for (let i = 0; i < primaryItem[primarySubItemsField].length; i++) {
        const secondaryItem = primaryItem[primarySubItemsField][i];

        if (secondaryValueField && secondaryItem[secondaryValueField] === secondaryValue) {
            return true;
        } else if (!secondaryValueField && secondaryItem === secondaryValue) {
            return true;
        }
    }

    return false;
}

export function getPrimaryValueBySecondaryValue(items, primarySubItemsField, primaryValueField, secondaryValueField, secondaryValue) {
    if (primarySubItemsField) {
        if (isArray(items)) {
            for (let i = 0; i < items.length; i++) {
                const primaryItem = items[i];

                if (isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryValue)) {
                    if (primaryValueField) {
                        return primaryItem[primaryValueField];
                    } else {
                        return primaryItem;
                    }
                }
            }
        } else {
            for (let field in items) {
                if (!Object.prototype.hasOwnProperty.call(items, field)) {
                    continue;
                }

                const primaryItem = items[field];

                if (isPrimaryItemHasSecondaryValue(primaryItem, primarySubItemsField, secondaryValueField, secondaryValue)) {
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

export function arrangeArrayWithNewStartIndex(array, startIndex) {
    if (startIndex <= 0 || startIndex >= array.length) {
        return array;
    }

    const newArray = [];

    for (let i = startIndex; i < array.length; i++) {
        newArray.push(array[i]);
    }

    for (let i = 0; i < startIndex; i++) {
        newArray.push(array[i]);
    }

    return newArray;
}
