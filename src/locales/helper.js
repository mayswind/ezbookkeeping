import { WeekDay, LongDateFormat, ShortDateFormat, LongTimeFormat, ShortTimeFormat, DateRange } from '@/core/datetime.ts';
import { DecimalSeparator, DigitGroupingSymbol, DigitGroupingType } from '@/core/numeral.ts';
import { CurrencyDisplayType } from '@/core/currency.ts'
import { AccountCategory } from '@/core/account.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';

import { UTC_TIMEZONE, ALL_TIMEZONES } from '@/consts/timezone.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { KnownErrorCode, SPECIFIED_API_NOT_FOUND_ERRORS, PARAMETERIZED_ERRORS } from '@/consts/api.ts';

import {
    isString,
    isNumber,
    isBoolean,
    copyObjectTo
} from '@/lib/common.ts';

import {
    parseDateFromUnixTime,
    formatUnixTime,
    getYear,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffset,
    getBrowserTimezoneOffsetMinutes,
    getTimeDifferenceHoursAndMinutes,
    getDateTimeFormatType,
    getRecentMonthDateRanges,
    isDateRangeMatchFullYears,
    isDateRangeMatchFullMonths
} from '@/lib/datetime.ts';

import {
    formatAmount,
    getAdaptiveDisplayAmountRate
} from '@/lib/numeral.ts';

import {
    getCurrencyFraction,
    appendCurrencySymbol
} from '@/lib/currency.ts';

import {
    getCategorizedAccountsMap,
    getAllFilteredAccountsBalance
} from '@/lib/account.ts';

function getLocalizedDisplayNameAndType(typeAndNames, translateFn) {
    const ret = [];

    for (let i = 0; i < typeAndNames.length; i++) {
        const nameAndType = typeAndNames[i];

        ret.push({
            type: nameAndType.type,
            displayName: translateFn(nameAndType.name)
        });
    }

    return ret;
}

function getCurrencyName(currencyCode, translateFn) {
    return translateFn(`currency.name.${currencyCode}`);
}

function getCurrencyUnitName(currencyCode, isPlural, translateFn) {
    const currencyInfo = ALL_CURRENCIES[currencyCode];

    if (currencyInfo && currencyInfo.unit) {
        if (isPlural) {
            return translateFn(`currency.unit.${currencyInfo.unit}.plural`);
        } else {
            return translateFn(`currency.unit.${currencyInfo.unit}.normal`);
        }
    }

    return '';
}

function getMonthdayOrdinal(monthDay, translateFn) {
    return translateFn(`datetime.monthDayOrdinal.${monthDay}`);
}

function getWeekdayShortName(weekDayName, translateFn) {
    return translateFn(`datetime.${weekDayName}.short`);
}

function getWeekdayLongName(weekDayName, translateFn) {
    return translateFn(`datetime.${weekDayName}.long`);
}

function getMultiMonthdayShortNames(monthDays, translateFn) {
    if (!monthDays) {
        return '';
    }

    if (monthDays.length === 1) {
        return translateFn('format.misc.monthDay', {
            ordinal: getMonthdayOrdinal(monthDays[0], translateFn)
        });
    } else {
        return translateFn('format.misc.monthDays', {
            multiMonthDays: joinMultiText(monthDays.map(monthDay =>
                translateFn('format.misc.eachMonthDayInMonthDays', {
                    ordinal: getMonthdayOrdinal(monthDay, translateFn)
                })), translateFn)
        });
    }
}

function getMultiWeekdayLongNames(weekdayTypes, firstDayOfWeek, translateFn) {
    const weekdayTypesMap = {};

    if (!isNumber(firstDayOfWeek)) {
        firstDayOfWeek = WeekDay.DefaultFirstDay.type;
    }

    for (let i = 0; i < weekdayTypes.length; i++) {
        weekdayTypesMap[weekdayTypes[i]] = true;
    }

    const allWeekDays = getAllWeekDays(firstDayOfWeek, translateFn);
    const finalWeekdayNames = [];

    for (let i = 0; i < allWeekDays.length; i++) {
        const weekDay = allWeekDays[i];

        if (weekdayTypesMap[weekDay.type]) {
            finalWeekdayNames.push(weekDay.displayName);
        }
    }

    return joinMultiText(finalWeekdayNames, translateFn);
}

function getI18nLongDateFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, LongDateFormat.all(), LongDateFormat.values(), 'format.longDate', defaultLongDateFormatTypeName, LongDateFormat.Default, formatTypeValue);
}

function getI18nShortDateFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, ShortDateFormat.all(), ShortDateFormat.values(), 'format.shortDate', defaultShortDateFormatTypeName, ShortDateFormat.Default, formatTypeValue);
}

function getI18nLongYearFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, LongDateFormat.all(), LongDateFormat.values(), 'format.longYear', defaultLongDateFormatTypeName, LongDateFormat.Default, formatTypeValue);
}

function getI18nShortYearFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, ShortDateFormat.all(), ShortDateFormat.values(), 'format.shortYear', defaultShortDateFormatTypeName, ShortDateFormat.Default, formatTypeValue);
}

function getI18nLongYearMonthFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, LongDateFormat.all(), LongDateFormat.values(), 'format.longYearMonth', defaultLongDateFormatTypeName, LongDateFormat.Default, formatTypeValue);
}

function getI18nShortYearMonthFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, ShortDateFormat.all(), ShortDateFormat.values(), 'format.shortYearMonth', defaultShortDateFormatTypeName, ShortDateFormat.Default, formatTypeValue);
}

function getI18nLongMonthDayFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, LongDateFormat.all(), LongDateFormat.values(), 'format.longMonthDay', defaultLongDateFormatTypeName, LongDateFormat.Default, formatTypeValue);
}

function getI18nShortMonthDayFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, ShortDateFormat.all(), ShortDateFormat.values(), 'format.shortMonthDay', defaultShortDateFormatTypeName, ShortDateFormat.Default, formatTypeValue);
}

function getI18nLongTimeFormat(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    return getDateTimeFormat(translateFn, LongTimeFormat.all(), LongTimeFormat.values(), 'format.longTime', defaultLongTimeFormatTypeName, LongTimeFormat.Default, formatTypeValue);
}

function getI18nShortTimeFormat(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    return getDateTimeFormat(translateFn, ShortTimeFormat.all(), ShortTimeFormat.values(), 'format.shortTime', defaultShortTimeFormatTypeName, ShortTimeFormat.Default, formatTypeValue);
}

function getDateTimeFormat(translateFn, allFormatMap, allFormatArray, localeFormatPathPrefix, localeDefaultFormatTypeName, systemDefaultFormatType, formatTypeValue) {
    const type = getDateTimeFormatType(allFormatMap, allFormatArray, formatTypeValue, localeDefaultFormatTypeName, systemDefaultFormatType);
    return translateFn(`${localeFormatPathPrefix}.${type.key}`);
}

function getAllTimezones(includeSystemDefault, translateFn) {
    const defaultTimezoneOffset = getBrowserTimezoneOffset();
    const defaultTimezoneOffsetMinutes = getBrowserTimezoneOffsetMinutes();
    const allTimezoneInfos = [];

    for (let i = 0; i < ALL_TIMEZONES.length; i++) {
        const utcOffset = (ALL_TIMEZONES[i].timezoneName !== UTC_TIMEZONE.timezoneName ? getTimezoneOffset(ALL_TIMEZONES[i].timezoneName) : '');
        const displayName = translateFn(`timezone.${ALL_TIMEZONES[i].displayName}`);

        allTimezoneInfos.push({
            name: ALL_TIMEZONES[i].timezoneName,
            utcOffset: utcOffset,
            utcOffsetMinutes: getTimezoneOffsetMinutes(ALL_TIMEZONES[i].timezoneName),
            displayName: displayName,
            displayNameWithUtcOffset: `(UTC${utcOffset}) ${displayName}`
        });
    }

    if (includeSystemDefault) {
        const defaultDisplayName = translateFn('System Default');

        allTimezoneInfos.push({
            name: '',
            utcOffset: defaultTimezoneOffset,
            utcOffsetMinutes: defaultTimezoneOffsetMinutes,
            displayName: defaultDisplayName,
            displayNameWithUtcOffset: `(UTC${defaultTimezoneOffset}) ${defaultDisplayName}`
        });
    }

    allTimezoneInfos.sort(function(c1, c2) {
        const utcOffset1 = parseInt(c1.utcOffset.replace(':', ''));
        const utcOffset2 = parseInt(c2.utcOffset.replace(':', ''));

        if (utcOffset1 !== utcOffset2) {
            return utcOffset1 - utcOffset2;
        }

        return c1.displayName.localeCompare(c2.displayName);
    })

    return allTimezoneInfos;
}

function getTimezoneDifferenceDisplayText(utcOffset, translateFn) {
    const defaultTimezoneOffset = getTimezoneOffsetMinutes();
    const offsetTime = getTimeDifferenceHoursAndMinutes(utcOffset - defaultTimezoneOffset);

    if (utcOffset > defaultTimezoneOffset) {
        if (offsetTime.offsetMinutes) {
            return translateFn('format.misc.hoursMinutesAheadOfDefaultTimezone', {
                hours: offsetTime.offsetHours,
                minutes: offsetTime.offsetMinutes
            });
        } else {
            return translateFn('format.misc.hoursAheadOfDefaultTimezone', {
                hours: offsetTime.offsetHours
            });
        }
    } else if (utcOffset < defaultTimezoneOffset) {
        if (offsetTime.offsetMinutes) {
            return translateFn('format.misc.hoursMinutesBehindDefaultTimezone', {
                hours: offsetTime.offsetHours,
                minutes: offsetTime.offsetMinutes
            });
        } else {
            return translateFn('format.misc.hoursBehindDefaultTimezone', {
                hours: offsetTime.offsetHours
            });
        }
    } else {
        return translateFn('Same time as default timezone');
    }
}

function getAllWeekDays(firstDayOfWeek, translateFn) {
    const ret = [];
    const allWeekDays = WeekDay.values();

    if (!isNumber(firstDayOfWeek)) {
        firstDayOfWeek = WeekDay.DefaultFirstDay.type;
    }

    for (let i = firstDayOfWeek; i < allWeekDays.length; i++) {
        const weekDay = allWeekDays[i];

        ret.push({
            type: weekDay.type,
            displayName: translateFn(`datetime.${weekDay.name}.long`)
        });
    }

    for (let i = 0; i < firstDayOfWeek; i++) {
        const weekDay = allWeekDays[i];

        ret.push({
            type: weekDay.type,
            displayName: translateFn(`datetime.${weekDay.name}.long`)
        });
    }

    return ret;
}

function getAllDateRanges(scene, includeCustom, includeBillingCycle, translateFn) {
    const ret = [];
    const allDateRanges = DateRange.values();

    for (let i = 0; i < allDateRanges.length; i++) {
        const dateRange = allDateRanges[i];

        if (!dateRange.isAvailableForScene(scene)) {
            continue;
        }

        if (dateRange.isBillingCycle) {
            if (includeBillingCycle) {
                ret.push({
                    type: dateRange.type,
                    displayName: translateFn(dateRange.name),
                    isBillingCycle: dateRange.isBillingCycle
                });
            }

            continue;
        }

        if (includeCustom || dateRange.type !== DateRange.Custom.type) {
            ret.push({
                type: dateRange.type,
                displayName: translateFn(dateRange.name)
            });
        }
    }

    return ret;
}

function getAllRecentMonthDateRanges(userStore, includeAll, includeCustom, translateFn) {
    const allRecentMonthDateRanges = [];
    const recentDateRanges = getRecentMonthDateRanges(12);

    if (includeAll) {
        allRecentMonthDateRanges.push({
            dateType: DateRange.All.type,
            minTime: 0,
            maxTime: 0,
            displayName: translateFn('All')
        });
    }

    for (let i = 0; i < recentDateRanges.length; i++) {
        const recentDateRange = recentDateRanges[i];

        allRecentMonthDateRanges.push({
            dateType: recentDateRange.dateType,
            minTime: recentDateRange.minTime,
            maxTime: recentDateRange.maxTime,
            year: recentDateRange.year,
            month: recentDateRange.month,
            isPreset: true,
            displayName: formatUnixTime(recentDateRange.minTime, getI18nLongYearMonthFormat(translateFn, userStore.currentUserLongDateFormat))
        });
    }

    if (includeCustom) {
        allRecentMonthDateRanges.push({
            dateType: DateRange.Custom.type,
            minTime: 0,
            maxTime: 0,
            displayName: translateFn('Custom Date')
        });
    }

    return allRecentMonthDateRanges;
}

function getDateRangeDisplayName(userStore, dateType, startTime, endTime, translateFn) {
    if (dateType === DateRange.All.type) {
        return translateFn(DateRange.All.name);
    }

    const allDateRanges = DateRange.values();

    for (let i = 0; i < allDateRanges.length; i++) {
        const dateRange = allDateRanges[i];

        if (dateRange && dateRange.type !== DateRange.Custom.type && dateRange.type === dateType && dateRange.name) {
            return translateFn(dateRange.name);
        }
    }

    if (isDateRangeMatchFullYears(startTime, endTime)) {
        const displayStartTime = formatUnixTime(startTime, getI18nShortYearFormat(translateFn, userStore.currentUserShortDateFormat));
        const displayEndTime = formatUnixTime(endTime, getI18nShortYearFormat(translateFn, userStore.currentUserShortDateFormat));

        return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
    }

    if (isDateRangeMatchFullMonths(startTime, endTime)) {
        const displayStartTime = formatUnixTime(startTime, getI18nShortYearMonthFormat(translateFn, userStore.currentUserShortDateFormat));
        const displayEndTime = formatUnixTime(endTime, getI18nShortYearMonthFormat(translateFn, userStore.currentUserShortDateFormat));

        return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
    }

    const startTimeYear = getYear(parseDateFromUnixTime(startTime));
    const endTimeYear = getYear(parseDateFromUnixTime(endTime));

    const displayStartTime = formatUnixTime(startTime, getI18nShortDateFormat(translateFn, userStore.currentUserShortDateFormat));
    const displayEndTime = formatUnixTime(endTime, getI18nShortDateFormat(translateFn, userStore.currentUserShortDateFormat));

    if (displayStartTime === displayEndTime) {
        return displayStartTime;
    } else if (startTimeYear === endTimeYear) {
        const displayShortEndTime = formatUnixTime(endTime, getI18nShortMonthDayFormat(translateFn, userStore.currentUserShortDateFormat));
        return `${displayStartTime} ~ ${displayShortEndTime}`;
    }

    return `${displayStartTime} ~ ${displayEndTime}`;
}

function getCurrentDecimalSeparator(translateFn, decimalSeparator) {
    let decimalSeparatorType = DecimalSeparator.valueOf(decimalSeparator);

    if (!decimalSeparatorType) {
        const defaultDecimalSeparatorTypeName = translateFn('default.decimalSeparator');
        decimalSeparatorType = DecimalSeparator.parse(defaultDecimalSeparatorTypeName);

        if (!decimalSeparatorType) {
            decimalSeparatorType = DecimalSeparator.Default;
        }
    }

    return decimalSeparatorType.symbol;
}

function getCurrentDigitGroupingSymbol(translateFn, digitGroupingSymbol) {
    let digitGroupingSymbolType = DigitGroupingSymbol.valueOf(digitGroupingSymbol);

    if (!digitGroupingSymbolType) {
        const defaultDigitGroupingSymbolTypeName = translateFn('default.digitGroupingSymbol');
        digitGroupingSymbolType = DigitGroupingSymbol.parse(defaultDigitGroupingSymbolTypeName);

        if (!digitGroupingSymbolType) {
            digitGroupingSymbolType = DigitGroupingSymbol.Default;
        }
    }

    return digitGroupingSymbolType.symbol;
}

function getCurrentDigitGroupingType(translateFn, digitGrouping) {
    let digitGroupingType = DigitGroupingType.valueOf(digitGrouping);

    if (!digitGroupingType) {
        const defaultDigitGroupingTypeName = translateFn('default.digitGrouping');
        digitGroupingType = DigitGroupingType.parse(defaultDigitGroupingTypeName);

        if (!digitGroupingType) {
            digitGroupingType = DigitGroupingType.Default;
        }
    }

    return digitGroupingType.type;
}

function getNumberFormatOptions(translateFn, userStore, currencyCode) {
    return {
        decimalSeparator: getCurrentDecimalSeparator(translateFn, userStore.currentUserDecimalSeparator),
        decimalNumberCount: getCurrencyFraction(currencyCode),
        digitGroupingSymbol: getCurrentDigitGroupingSymbol(translateFn, userStore.currentUserDigitGroupingSymbol),
        digitGrouping: getCurrentDigitGroupingType(translateFn, userStore.currentUserDigitGrouping),
    };
}

function getFormattedAmount(value, translateFn, userStore, currencyCode) {
    const numberFormatOptions = getNumberFormatOptions(translateFn, userStore, currencyCode);
    return formatAmount(value, numberFormatOptions);
}

function getCurrentCurrencyDisplayType(translateFn, userStore) {
    let currencyDisplayType = CurrencyDisplayType.valueOf(userStore.currentUserCurrencyDisplayType);

    if (!currencyDisplayType) {
        const defaultCurrencyDisplayTypeName = translateFn('default.currencyDisplayType');
        currencyDisplayType = CurrencyDisplayType.parse(defaultCurrencyDisplayTypeName);
    }

    if (!currencyDisplayType) {
        currencyDisplayType = CurrencyDisplayType.Default;
    }

    return currencyDisplayType;
}

function getFormattedAmountWithCurrency(value, currencyCode, translateFn, userStore, settingsStore, notConvertValue, currencyDisplayType) {
    if (!isNumber(value) && !isString(value)) {
        return value;
    }

    if (isNumber(value)) {
        value = value.toString();
    }

    const isPlural = value !== '100' && value !== '-100';

    if (!notConvertValue) {
        const numberFormatOptions = getNumberFormatOptions(translateFn, userStore, currencyCode);
        const hasIncompleteFlag = isString(value) && value.charAt(value.length - 1) === '+';

        if (hasIncompleteFlag) {
            value = value.substring(0, value.length - 1);
        }

        value = formatAmount(value, numberFormatOptions);

        if (hasIncompleteFlag) {
            value = value + '+';
        }
    }

    if (!isBoolean(currencyCode) && !currencyCode) {
        currencyCode = userStore.currentUserDefaultCurrency;
    } else if (isBoolean(currencyCode) && !currencyCode) {
        currencyCode = '';
    }

    if (!currencyCode) {
        return value;
    }

    if (!currencyDisplayType) {
        currencyDisplayType = getCurrentCurrencyDisplayType(translateFn, userStore);
    }

    const currencyUnit = getCurrencyUnitName(currencyCode, isPlural, translateFn);
    const currencyName = getCurrencyName(currencyCode, translateFn);
    return appendCurrencySymbol(value, currencyDisplayType, currencyCode, currencyUnit, currencyName, isPlural);
}

function getAdaptiveAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, translateFn, userStore) {
    const numberFormatOptions = getNumberFormatOptions(translateFn, userStore);
    return getAdaptiveDisplayAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, numberFormatOptions);
}

function getAllTransactionTagFilterTypes(translateFn) {
    return getLocalizedDisplayNameAndType(TransactionTagFilterType.values(), translateFn);
}

function getCategorizedAccountsWithDisplayBalance(allVisibleAccounts, showAccountBalance, defaultCurrency, userStore, settingsStore, exchangeRatesStore, translateFn) {
    const ret = [];
    const allCategories = AccountCategory.values();
    const categorizedAccounts = copyObjectTo(getCategorizedAccountsMap(allVisibleAccounts), {});

    for (let i = 0; i < allCategories.length; i++) {
        const category = allCategories[i];

        if (!categorizedAccounts[category.type]) {
            continue;
        }

        const accountCategory = categorizedAccounts[category.type];

        if (accountCategory.accounts) {
            for (let i = 0; i < accountCategory.accounts.length; i++) {
                const account = accountCategory.accounts[i];

                if (showAccountBalance && account.isAsset) {
                    account.displayBalance = getFormattedAmountWithCurrency(account.balance, account.currency, translateFn, userStore, settingsStore);
                } else if (showAccountBalance && account.isLiability) {
                    account.displayBalance = getFormattedAmountWithCurrency(-account.balance, account.currency, translateFn, userStore, settingsStore);
                } else {
                    account.displayBalance = '***';
                }
            }
        }

        if (showAccountBalance) {
            const accountsBalance = getAllFilteredAccountsBalance(categorizedAccounts, account => account.category === accountCategory.category);
            let totalBalance = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === defaultCurrency) {
                    if (accountsBalance[i].isAsset) {
                        totalBalance += accountsBalance[i].balance;
                    } else if (accountsBalance[i].isLiability) {
                        totalBalance -= accountsBalance[i].balance;
                    }
                } else {
                    const balance = exchangeRatesStore.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, defaultCurrency);

                    if (!isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    if (accountsBalance[i].isAsset) {
                        totalBalance += Math.floor(balance);
                    } else if (accountsBalance[i].isLiability) {
                        totalBalance -= Math.floor(balance);
                    }
                }
            }

            if (hasUnCalculatedAmount) {
                totalBalance = totalBalance + '+';
            }

            accountCategory.displayBalance = getFormattedAmountWithCurrency(totalBalance, defaultCurrency, translateFn, userStore, settingsStore);
        } else {
            accountCategory.displayBalance = '***';
        }

        ret.push(accountCategory);
    }

    return ret;
}

function joinMultiText(textArray, translateFn) {
    if (!textArray || !textArray.length) {
        return '';
    }

    const separator = translateFn('format.misc.multiTextJoinSeparator');

    return textArray.join(separator);
}

function getLocalizedError(error) {
    if (error.errorCode === KnownErrorCode.ApiNotFound && SPECIFIED_API_NOT_FOUND_ERRORS[error.path]) {
        return {
            message: `${SPECIFIED_API_NOT_FOUND_ERRORS[error.path].message}`
        };
    }

    if (error.errorCode !== KnownErrorCode.ValidatorError) {
        return {
            message: `error.${error.errorMessage}`
        };
    }

    for (let i = 0; i < PARAMETERIZED_ERRORS.length; i++) {
        const errorInfo = PARAMETERIZED_ERRORS[i];
        const matches = error.errorMessage.match(errorInfo.regex);

        if (matches && matches.length === errorInfo.parameters.length + 1) {
            return {
                message: `parameterizedError.${errorInfo.localeKey}`,
                parameters: errorInfo.parameters.map((param, index) => {
                    return {
                        key: param.field,
                        localized: param.localized,
                        value: matches[index + 1]
                    }
                })
            };
        }
    }

    return {
        message: `error.${error.errorMessage}`
    };
}

function getLocalizedErrorParameters(parameters, i18nFunc) {
    let localizedParameters = {};

    if (parameters) {
        for (let i = 0; i < parameters.length; i++) {
            const parameter = parameters[i];

            if (parameter.localized) {
                localizedParameters[parameter.key] = i18nFunc(`parameter.${parameter.value}`);
            } else {
                localizedParameters[parameter.key] = parameter.value;
            }
        }
    }

    return localizedParameters;
}

export function translateError(message, translateFn) {
    let parameters = {};

    if (message && message.error) {
        const localizedError = getLocalizedError(message.error);
        message = localizedError.message;
        parameters = getLocalizedErrorParameters(localizedError.parameters, translateFn);
    }

    return translateFn(message, parameters);
}

export function i18nFunctions(i18nGlobal) {
    return {
        getWeekdayShortName: (weekDay) => getWeekdayShortName(weekDay, i18nGlobal.t),
        getWeekdayLongName: (weekDay) => getWeekdayLongName(weekDay, i18nGlobal.t),
        getMultiMonthdayShortNames: (monthdays) => getMultiMonthdayShortNames(monthdays, i18nGlobal.t),
        getMultiWeekdayLongNames: (weekdayTypes, firstDayOfWeek) => getMultiWeekdayLongNames(weekdayTypes, firstDayOfWeek, i18nGlobal.t),
        formatUnixTimeToLongDateTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongDateFormat(i18nGlobal.t, userStore.currentUserLongDateFormat) + ' ' + getI18nLongTimeFormat(i18nGlobal.t, userStore.currentUserLongTimeFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongDate: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongDateFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongYear: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongYearFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongYearMonth: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongYearMonthFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongMonthDay: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongMonthDayFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongTimeFormat(i18nGlobal.t, userStore.currentUserLongTimeFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortTimeFormat(i18nGlobal.t, userStore.currentUserShortTimeFormat), utcOffset, currentUtcOffset),
        getAllTimezones: (includeSystemDefault) => getAllTimezones(includeSystemDefault, i18nGlobal.t),
        getTimezoneDifferenceDisplayText: (utcOffset) => getTimezoneDifferenceDisplayText(utcOffset, i18nGlobal.t),
        getAllDateRanges: (scene, includeCustom, includeBillingCycle) => getAllDateRanges(scene, includeCustom, includeBillingCycle, i18nGlobal.t),
        getAllRecentMonthDateRanges: (userStore, includeAll, includeCustom) => getAllRecentMonthDateRanges(userStore, includeAll, includeCustom, i18nGlobal.t),
        getDateRangeDisplayName: (userStore, dateType, startTime, endTime) => getDateRangeDisplayName(userStore, dateType, startTime, endTime, i18nGlobal.t),
        formatAmount: (userStore, value, currencyCode) => getFormattedAmount(value, i18nGlobal.t, userStore, currencyCode),
        formatAmountWithCurrency: (settingsStore, userStore, value, currencyCode) => getFormattedAmountWithCurrency(value, currencyCode, i18nGlobal.t, userStore, settingsStore),
        getAdaptiveAmountRate: (userStore, amount1, amount2, fromExchangeRate, toExchangeRate) => getAdaptiveAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, i18nGlobal.t, userStore),
        getAllTransactionTagFilterTypes: () => getAllTransactionTagFilterTypes(i18nGlobal.t),
        getCategorizedAccountsWithDisplayBalance: (allVisibleAccounts, showAccountBalance, defaultCurrency, settingsStore, userStore, exchangeRatesStore) => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts, showAccountBalance, defaultCurrency, userStore, settingsStore, exchangeRatesStore, i18nGlobal.t)
    };
}
