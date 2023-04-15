import CryptoJS from 'crypto-js';
import moment from 'moment';
import Clipboard from 'clipboard';
import uaParser from 'ua-parser-js';

import dateTimeConstants from '../consts/datetime.js';
import accountConstants from '../consts/account.js';
import categoryConstants from '../consts/category.js';
import transactionConstants from '../consts/transaction.js';
import settings from './settings.js';

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

function isEquals(obj1, obj2) {
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

function getUtcOffsetMinutesByUtcOffset(utcOffset) {
    if (!utcOffset) {
        return 0;
    }

    const parts = utcOffset.split(':');

    if (parts.length !== 2) {
        return 0;
    }

    return parseInt(parts[0]) * 60 + parseInt(parts[1]);
}

function getUtcOffsetByUtcOffsetMinutes(utcOffsetMinutes) {
    let offsetHours = parseInt(Math.abs(utcOffsetMinutes) / 60);
    let offsetMinutes = Math.abs(utcOffsetMinutes) - offsetHours * 60;

    if (offsetHours < 10) {
        offsetHours = '0' + offsetHours;
    }

    if (offsetMinutes < 10) {
        offsetMinutes = '0' + offsetMinutes;
    }

    if (utcOffsetMinutes >= 0) {
        return `+${offsetHours}:${offsetMinutes}`;
    } else if (utcOffsetMinutes < 0) {
        return `-${offsetHours}:${offsetMinutes}`;
    }
}

function getTimezoneOffset(timezone) {
    if (timezone) {
        return moment().tz(timezone).format('Z');
    } else {
        return moment().format('Z');
    }
}

function getTimezoneOffsetMinutes(timezone) {
    const utcOffset = getTimezoneOffset(timezone);
    return getUtcOffsetMinutesByUtcOffset(utcOffset);
}

function getBrowserTimezoneOffsetMinutes() {
    return -new Date().getTimezoneOffset();
}

function getLocalDatetimeFromUnixTime(unixTime) {
    return new Date(unixTime * 1000);
}

function getUnixTimeFromLocalDatetime(datetime) {
    return datetime.getTime() / 1000;
}

function getActualUnixTimeForStore(unixTime, utcOffset, currentUtcOffset) {
    return unixTime - (utcOffset - currentUtcOffset) * 60;
}

function getDummyUnixTimeForLocalUsage(unixTime, utcOffset, currentUtcOffset) {
    return unixTime + (utcOffset - currentUtcOffset) * 60;
}

function getCurrentUnixTime() {
    return moment().unix();
}

function parseDateFromUnixTime(unixTime, utcOffset, currentUtcOffset) {
    if (isNumber(utcOffset)) {
        if (!isNumber(currentUtcOffset)) {
            currentUtcOffset = getTimezoneOffsetMinutes();
        }

        unixTime = getDummyUnixTimeForLocalUsage(unixTime, utcOffset, currentUtcOffset);
    }

    return moment.unix(unixTime);
}

function formatUnixTime(unixTime, format, utcOffset, currentUtcOffset) {
    return parseDateFromUnixTime(unixTime, utcOffset, currentUtcOffset).format(format);
}

function formatTime(dateTime, format) {
    return moment(dateTime).format(format);
}

function getUnixTime(date) {
    return moment(date).unix();
}

function getYear(date) {
    return moment(date).year();
}

function getMonth(date) {
    return moment(date).month() + 1;
}

function getYearAndMonth(date) {
    const year = getYear(date);
    let month = getMonth(date);

    if (month < 10) {
        month = '0' + month;
    }

    return `${year}-${month}`;
}

function getDay(date) {
    return moment(date).date();
}

function getDayOfWeekName(date) {
    const dayOfWeek = moment(date).days();
    return dateTimeConstants.allWeekDaysArray[dayOfWeek].name;
}

function getHour(date) {
    return moment(date).hour();
}

function getMinute(date) {
    return moment(date).minute();
}

function getSecond(date) {
    return moment(date).second();
}

function getUnixTimeBeforeUnixTime(unixTime, amount, unit) {
    return moment.unix(unixTime).subtract(amount, unit).unix();
}

function getUnixTimeAfterUnixTime(unixTime, amount, unit) {
    return moment.unix(unixTime).add(amount, unit).unix();
}

function getMinuteFirstUnixTime(date) {
    const datetime = moment(date);
    return datetime.set({ second: 0, millisecond: 0 }).unix();
}

function getMinuteLastUnixTime(date) {
    return moment.unix(getMinuteFirstUnixTime(date)).add(1, 'minutes').subtract(1, 'seconds').unix();
}

function getTodayFirstUnixTime() {
    return moment().set({ hour: 0, minute: 0, second: 0, millisecond: 0 }).unix();
}

function getTodayLastUnixTime() {
    return moment.unix(getTodayFirstUnixTime()).add(1, 'days').subtract(1, 'seconds').unix();
}

function getThisWeekFirstUnixTime(firstDayOfWeek) {
    const today = moment.unix(getTodayFirstUnixTime());

    if (!isNumber(firstDayOfWeek)) {
        firstDayOfWeek = 0;
    }

    let dayOfWeek = today.day() - firstDayOfWeek;

    if (dayOfWeek < 0) {
        dayOfWeek += 7;
    }

    return today.subtract(dayOfWeek, 'days').unix();
}

function getThisWeekLastUnixTime(firstDayOfWeek) {
    return moment.unix(getThisWeekFirstUnixTime(firstDayOfWeek)).add(7, 'days').subtract(1, 'seconds').unix();
}

function getThisMonthFirstUnixTime() {
    const today = moment.unix(getTodayFirstUnixTime());
    return today.subtract(today.date() - 1, 'days').unix();
}

function getThisMonthLastUnixTime() {
    return moment.unix(getThisMonthFirstUnixTime()).add(1, 'months').subtract(1, 'seconds').unix();
}

function getThisYearFirstUnixTime() {
    const today = moment.unix(getTodayFirstUnixTime());
    return today.subtract(today.dayOfYear() - 1, 'days').unix();
}

function getThisYearLastUnixTime() {
    return moment.unix(getThisYearFirstUnixTime()).add(1, 'years').subtract(1, 'seconds').unix();
}

function getShiftedDateRange(minTime, maxTime, scale) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ second: 0, millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ second: 59, millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    if (firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix()) {
        const months = getYear(maxDateTime) * 12 + getMonth(maxDateTime) - getYear(minDateTime) * 12 - getMonth(minDateTime) + 1;
        const newMinDateTime = minDateTime.add(months * scale, 'months');
        const newMaxDateTime = newMinDateTime.clone().add(months, 'months').subtract(1, 'seconds');

        return {
            minTime: newMinDateTime.unix(),
            maxTime: newMaxDateTime.unix()
        };
    }

    const range = (maxTime - minTime + 1) * scale;

    return {
        minTime: minTime + range,
        maxTime: maxTime + range
    };
}

function getDateRangeByDateType(dateType, firstDayOfWeek) {
    let maxTime = 0;
    let minTime = 0;

    if (dateType === dateTimeConstants.allDateRanges.All.type) { // All
        maxTime = 0;
        minTime = 0;
    } else if (dateType === dateTimeConstants.allDateRanges.Today.type) { // Today
        maxTime = getTodayLastUnixTime();
        minTime = getTodayFirstUnixTime();
    } else if (dateType === dateTimeConstants.allDateRanges.Yesterday.type) { // Yesterday
        maxTime = getUnixTimeBeforeUnixTime(getTodayLastUnixTime(), 1, 'days');
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 1, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.LastSevenDays.type) { // Last 7 days
        maxTime = getTodayLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 6, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.LastThirtyDays.type) { // Last 30 days
        maxTime = getTodayLastUnixTime();
        minTime = getUnixTimeBeforeUnixTime(getTodayFirstUnixTime(), 29, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.ThisWeek.type) { // This week
        maxTime = getThisWeekLastUnixTime(firstDayOfWeek);
        minTime = getThisWeekFirstUnixTime(firstDayOfWeek);
    } else if (dateType === dateTimeConstants.allDateRanges.LastWeek.type) { // Last week
        maxTime = getUnixTimeBeforeUnixTime(getThisWeekLastUnixTime(firstDayOfWeek), 7, 'days');
        minTime = getUnixTimeBeforeUnixTime(getThisWeekFirstUnixTime(firstDayOfWeek), 7, 'days');
    } else if (dateType === dateTimeConstants.allDateRanges.ThisMonth.type) { // This month
        maxTime = getThisMonthLastUnixTime();
        minTime = getThisMonthFirstUnixTime();
    } else if (dateType === dateTimeConstants.allDateRanges.LastMonth.type) { // Last month
        maxTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'seconds');
        minTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'months');
    } else if (dateType === dateTimeConstants.allDateRanges.ThisYear.type) { // This year
        maxTime = getThisYearLastUnixTime();
        minTime = getThisYearFirstUnixTime();
    } else if (dateType === dateTimeConstants.allDateRanges.LastYear.type) { // Last year
        maxTime = getUnixTimeBeforeUnixTime(getThisYearLastUnixTime(), 1, 'years');
        minTime = getUnixTimeBeforeUnixTime(getThisYearFirstUnixTime(), 1, 'years');
    } else {
        return null;
    }

    return {
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    };
}

function isDateRangeMatchFullYears(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ second: 0, millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ second: 59, millisecond: 999 });

    const firstDayOfYear = minDateTime.clone().startOf('year');
    const lastDayOfYear = maxDateTime.clone().endOf('year');

    return firstDayOfYear.unix() === minDateTime.unix() && lastDayOfYear.unix() === maxDateTime.unix();
}

function isDateRangeMatchFullMonths(minTime, maxTime) {
    const minDateTime = parseDateFromUnixTime(minTime).set({ second: 0, millisecond: 0 });
    const maxDateTime = parseDateFromUnixTime(maxTime).set({ second: 59, millisecond: 999 });

    const firstDayOfMonth = minDateTime.clone().startOf('month');
    const lastDayOfMonth = maxDateTime.clone().endOf('month');

    return firstDayOfMonth.unix() === minDateTime.unix() && lastDayOfMonth.unix() === maxDateTime.unix();
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

function numericCurrencyToString(num) {
    let str = num.toString();
    const negative = str.charAt(0) === '-';

    if (negative) {
        str = str.substring(1);
    }

    if (str.length === 0) {
        str = '0.00';
    } else if (str.length === 1) {
        str = '0.0' + str;
    } else if (str.length === 2) {
        str = '0.' + str;
    } else {
        let integer = str.substring(0, str.length - 2);
        let decimals = str.substring(str.length - 2);

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
        str = str.substring(1);
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

    const integer = str.substring(0, dotPos);
    const decimals = str.substring(dotPos + 1, str.length);

    if (decimals.length < 1) {
        return sign * parseInt(integer) * 100;
    } else if (decimals.length === 1) {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals) * 10;
    } else if (decimals.length === 2) {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals);
    } else {
        return sign * parseInt(integer) * 100 + sign * parseInt(decimals.substring(0, 2));
    }
}

function getExchangedAmount(amount, fromRate, toRate) {
    const exchangeRate = parseFloat(toRate) / parseFloat(fromRate);

    if (!isNumber(exchangeRate)) {
        return null;
    }

    return amount * exchangeRate;
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

function getNameByKeyValue(src, value, keyField, nameField, defaultName) {
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

function generateRandomString() {
    const baseString = 'ebk_' + Math.round(new Date().getTime() / 1000) + '_' + Math.random();
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

function parseDeviceInfo(ua) {
    const uaInfo = parseUserAgent(ua);
    let result = '';

    if (uaInfo.device.model) {
        result = uaInfo.device.model;
    } else if (uaInfo.os.name) {
        result = uaInfo.os.name;

        if (uaInfo.os.version) {
            result += ' ' + uaInfo.os.version;
        }
    }

    if (uaInfo.browser.name) {
        let browserInfo = uaInfo.browser.name;

        if (uaInfo.browser.version) {
            browserInfo += ' ' + uaInfo.browser.version;
        }

        if (result) {
            result += ' (' + browserInfo + ')';
        } else {
            result = browserInfo;
        }
    }

    if (!result) {
        return 'Unknown Device';
    }

    return result;
}

function transactionTypeToCategroyType(transactionType) {
    if (transactionType === transactionConstants.allTransactionTypes.Income) {
        return categoryConstants.allCategoryTypes.Income;
    } else if (transactionType === transactionConstants.allTransactionTypes.Expense) {
        return categoryConstants.allCategoryTypes.Expense;
    } else if (transactionType === transactionConstants.allTransactionTypes.Transfer) {
        return categoryConstants.allCategoryTypes.Transfer;
    } else {
        return null;
    }
}

function categroyTypeToTransactionType(categoryType) {
    if (categoryType === categoryConstants.allCategoryTypes.Income) {
        return transactionConstants.allTransactionTypes.Income;
    } else if (categoryType === categoryConstants.allCategoryTypes.Expense) {
        return transactionConstants.allTransactionTypes.Expense;
    } else if (categoryType === categoryConstants.allCategoryTypes.Transfer) {
        return transactionConstants.allTransactionTypes.Transfer;
    } else {
        return null;
    }
}

function getCategoryInfo(categoryId) {
    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        if (accountConstants.allCategories[i].id === categoryId) {
            return accountConstants.allCategories[i];
        }
    }

    return null;
}

function getCategorizedAccounts(allAccounts) {
    const ret = {};

    for (let i = 0; i < allAccounts.length; i++) {
        const account = allAccounts[i];

        if (!ret[account.category]) {
            const categoryInfo = getCategoryInfo(account.category);

            if (categoryInfo) {
                ret[account.category] = {
                    category: account.category,
                    name: categoryInfo.name,
                    icon: categoryInfo.defaultAccountIconId,
                    accounts: []
                };
            }
        }

        if (ret[account.category]) {
            const accountList = ret[account.category].accounts;
            accountList.push(account);
        }
    }

    return ret;
}

function getAllFilteredAccountsBalance(categorizedAccounts, accountFilter) {
    const allAccountCategories = accountConstants.allCategories;
    const ret = [];

    for (let categoryIdx = 0; categoryIdx < allAccountCategories.length; categoryIdx++) {
        const accountCategory = allAccountCategories[categoryIdx];

        if (!categorizedAccounts[accountCategory.id] || !categorizedAccounts[accountCategory.id].accounts) {
            continue;
        }

        for (let accountIdx = 0; accountIdx < categorizedAccounts[accountCategory.id].accounts.length; accountIdx++) {
            const account = categorizedAccounts[accountCategory.id].accounts[accountIdx];

            if (account.hidden || !accountFilter(account)) {
                continue;
            }

            if (account.type === accountConstants.allAccountTypes.SingleAccount) {
                ret.push({
                    balance: account.balance,
                    isAsset: account.isAsset,
                    isLiability: account.isLiability,
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
                        isAsset: subAccount.isAsset,
                        isLiability: subAccount.isLiability,
                        currency: subAccount.currency
                    });
                }
            }
        }
    }

    return ret;
}

function makeButtonCopyToClipboard({ text, el, successCallback, errorCallback }) {
    const clipboard = new Clipboard(el, {
        text: function () {
            return text;
        }
    });

    clipboard.on('success', (e) => {
        if (successCallback) {
            successCallback(e);
        }
    });

    clipboard.on('error', (e) => {
        if (errorCallback) {
            errorCallback(e);
        }
    });

    return clipboard;
}

function changeClipboardObjectText(clipboard, text) {
    if (!clipboard) {
        return;
    }

    clipboard.text = function () {
        return text;
    };
}

export default {
    isFunction,
    isObject,
    isArray,
    isString,
    isNumber,
    isBoolean,
    isEquals,
    getUtcOffsetMinutesByUtcOffset,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getUnixTimeFromLocalDatetime,
    getActualUnixTimeForStore,
    getDummyUnixTimeForLocalUsage,
    getCurrentUnixTime,
    parseDateFromUnixTime,
    formatUnixTime,
    formatTime,
    getUnixTime,
    getYear,
    getMonth,
    getYearAndMonth,
    getDay,
    getDayOfWeekName,
    getHour,
    getMinute,
    getSecond,
    getUnixTimeBeforeUnixTime,
    getUnixTimeAfterUnixTime,
    getMinuteFirstUnixTime,
    getMinuteLastUnixTime,
    getTodayFirstUnixTime,
    getTodayLastUnixTime,
    getThisWeekFirstUnixTime,
    getThisWeekLastUnixTime,
    getThisMonthFirstUnixTime,
    getThisMonthLastUnixTime,
    getThisYearFirstUnixTime,
    getThisYearLastUnixTime,
    getShiftedDateRange,
    getDateRangeByDateType,
    isDateRangeMatchFullYears,
    isDateRangeMatchFullMonths,
    copyObjectTo,
    copyArrayTo,
    appendThousandsSeparator,
    numericCurrencyToString,
    stringCurrencyToNumeric,
    getExchangedAmount,
    base64encode,
    arrayBufferToString,
    stringToArrayBuffer,
    getNameByKeyValue,
    generateRandomString,
    parseUserAgent,
    parseDeviceInfo,
    transactionTypeToCategroyType,
    categroyTypeToTransactionType,
    getCategoryInfo,
    getCategorizedAccounts,
    getAllFilteredAccountsBalance,
    makeButtonCopyToClipboard,
    changeClipboardObjectText,
};
