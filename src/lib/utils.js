import CryptoJS from "crypto-js";
import uaParser from 'ua-parser-js';
import accountConstants from '../consts/account.js';

function isFunction(val) {
    return typeof(val) === 'function';
}

function isObject(val) {
    return val != null && typeof(val) === 'object';
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
    base64encode,
    arrayBufferToString,
    stringToArrayBuffer,
    generateRandomString,
    parseUserAgent,
    getCategorizedAccounts,
    getAccountByAccountId,
    getAllFilteredAccountsBalance,
};
