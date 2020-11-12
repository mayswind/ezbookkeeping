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

export default {
    isFunction,
    isObject,
    isArray,
    isString,
    isNumber,
    isBoolean,
    getCategorizedAccounts,
    getAccountByAccountId,
};
