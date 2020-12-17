import CryptoJS from "crypto-js";
import moment from 'moment';
import uaParser from 'ua-parser-js';

import accountConstants from '../consts/account.js';
import settings from "./settings.js";

function isFunction(val) {
    return typeof(val) === 'function';
}

function isObject(val) {
    return val != null && typeof(val) === 'object' && !isArray(val);
}

function isArray(val) {
    if (isFunction(Array.isArray)) {
        return Array.isArray(val);
    }

    return Object.prototype.toString.call(val) === '[object Array]';
}

function isString(val) {
    return typeof(val) === 'string';
}

function isNumber(val) {
    return typeof(val) === 'number';
}

function isBoolean(val) {
    return typeof(val) === 'boolean';
}

function formatDate(date, format) {
    return moment(date).format(format);
}

function getUnixTime(date) {
    return moment(date).unix();
}

function copyObjectTo(fromObject, toObject) {
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
            toObject[key] = this.copyArrayTo(fromValue, toValue);
        } else if (isObject(fromValue)) {
            toObject[key] = this.copyObjectTo(fromValue, toValue);
        } else {
            if (fromValue !== toValue) {
                toObject[key] = fromValue;
            }
        }
    }

    return toObject;
}

function copyArrayTo(fromArray, toArray) {
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
                toArray[i] = this.copyArrayTo(fromValue, toValue);
            } else if (isObject(fromValue)) {
                toArray[i] = this.copyObjectTo(fromValue, toValue);
            } else {
                if (fromValue !== toValue) {
                    toArray[i] = fromValue;
                }
            }
        } else {
            if (isArray(fromValue)) {
                toArray.push(this.copyArrayTo(fromValue, []));
            } else if (isObject(fromValue)) {
                toArray.push(this.copyObjectTo(fromValue, {}));
            } else {
                toArray.push(fromValue);
            }
        }
    }

    return toArray;
}

function appendThousandsSeparator(value) {
    if (!settings.isEnableThousandsSeparator() || value.length <= 3) {
        return value;
    }

    const negative = value.charAt(0) === '-';

    if (negative) {
        value = value.substr(1);
    }

    const dotPos = value.indexOf('.');
    const integer = dotPos < 0 ? value : value.substr(0, dotPos);
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

function numericCurrencyToString(num) {
    let str = num.toString();
    const negative = str.charAt(0) === '-';

    if (negative) {
        str = str.substr(1);
    }

    if (str.length === 0) {
        str = '0.00';
    } else if (str.length === 1) {
        str = '0.0' + str;
    } else if (str.length === 2) {
        str = '0.' + str;
    } else {
        let integer = str.substr(0, str.length - 2);
        let decimals = str.substr(str.length - 2, 2);

        integer = appendThousandsSeparator(integer);

        str = `${integer}.${decimals}`;
    }

    if (negative) {
        str = `-${str}`;
    }

    return str;
}

function stringCurrencyToNumeric(str) {
    if (!str || str.length < 1) {
        return 0;
    }

    const negative = str.charAt(0) === '-';

    if (negative) {
        str = str.substr(1);
    }

    if (!str || str.length < 1) {
        return 0;
    }

    const sign = negative ? -1 : 1;

    if (str.indexOf(',')) {
        str = str.replaceAll(/,/g, '');
    }

    let dotPos = str.indexOf('.');

    if (dotPos < 0) {
        return sign * parseInt(str) * 100;
    } else if (dotPos === 0) {
        str = '0' + str;
        dotPos++;
    }

    const integer = str.substr(0, dotPos);
    const decimals = str.substring(dotPos + 1, str.length);

    if (decimals.length < 1) {
        return sign * parseInt(integer) * 100;
    } else if (decimals.length === 1) {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals) * 10;
    } else if (decimals.length === 2) {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals);
    } else {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals.substr(0, 2));
    }
}

function base64encode(arrayBuffer) {
    if (!arrayBuffer || arrayBuffer.length === 0) {
        return null;
    }

    return btoa(String.fromCharCode.apply(null, new Uint8Array(arrayBuffer)));
}

function arrayBufferToString(arrayBuffer) {
    return String.fromCharCode.apply(null, new Uint8Array(arrayBuffer));
}

function stringToArrayBuffer(str){
    return Uint8Array.from(str, c => c.charCodeAt(0)).buffer;
}

function generateRandomString() {
    const baseString = 'lab_' + Math.round(new Date().getTime() / 1000) + '_' + Math.random();
    return CryptoJS.SHA256(baseString).toString();
}

function parseUserAgent(ua) {
    const uaParseRet = uaParser(ua);

    return {
        device: {
            vendor: uaParseRet.device.vendor,
            model: uaParseRet.device.model,
            type: uaParseRet.device.type
        },
        os: {
            name: uaParseRet.os.name,
            version: uaParseRet.os.version
        },
        browser: {
            name: uaParseRet.browser.name,
            version: uaParseRet.browser.version
        }
    };
}

function getCategorizedAccounts(allAccounts) {
    const ret = {};

    for (let i = 0; i < allAccounts.length; i++) {
        const account = allAccounts[i];

        if (!ret[account.category]) {
            ret[account.category] = [];
        }

        const accountList = ret[account.category];
        accountList.push(account);
    }

    return ret;
}

function getAccountByAccountId(categorizedAccounts, accountId) {
    for (let category in categorizedAccounts) {
        if (!Object.prototype.hasOwnProperty.call(categorizedAccounts, category)) {
            continue;
        }

        const accountList = categorizedAccounts[category];

        for (let i = 0; i < accountList.length; i++) {
            if (accountList[i].id === accountId) {
                return accountList[i];
            }
        }
    }

    return null;
}

function getAllFilteredAccountsBalance(categorizedAccounts, accountFilter) {
    const allAccountCategories = accountConstants.allCategories;
    const ret = [];

    for (let categoryIdx = 0; categoryIdx < allAccountCategories.length; categoryIdx++) {
        const accountCategory = allAccountCategories[categoryIdx];

        if (!categorizedAccounts[accountCategory.id]) {
            continue;
        }

        for (let accountIdx = 0; accountIdx < categorizedAccounts[accountCategory.id].length; accountIdx++) {
            const account = categorizedAccounts[accountCategory.id][accountIdx];

            if (account.hidden || !accountFilter(account)) {
                continue;
            }

            if (account.type === accountConstants.allAccountTypes.SingleAccount) {
                ret.push({
                    balance: account.balance,
                    currency: account.currency
                });
            } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
                for (let subAccountIdx = 0; subAccountIdx < account.subAccounts.length; subAccountIdx++) {
                    const subAccount = account.subAccounts[subAccountIdx];

                    if (subAccount.hidden || !accountFilter(subAccount)) {
                        continue;
                    }

                    ret.push({
                        balance: subAccount.balance,
                        currency: subAccount.currency
                    });
                }
            }
        }
    }

    return ret;
}

export default {
    isFunction,
    isObject,
    isArray,
    isString,
    isNumber,
    isBoolean,
    formatDate,
    getUnixTime,
    copyObjectTo,
    copyArrayTo,
    appendThousandsSeparator,
    numericCurrencyToString,
    stringCurrencyToNumeric,
    base64encode,
    arrayBufferToString,
    stringToArrayBuffer,
    generateRandomString,
    parseUserAgent,
    getCategorizedAccounts,
    getAccountByAccountId,
    getAllFilteredAccountsBalance,
};
