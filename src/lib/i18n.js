import moment from 'moment-timezone';

import { defaultLanguage, allLanguages } from '@/locales/index.js';
import datetime from '@/consts/datetime.js';
import timezone from '@/consts/timezone.js';
import currency from '@/consts/currency.js';
import account from '@/consts/account.js';
import category from '@/consts/category.js';
import statistics from '@/consts/statistics.js';

import {
    isString,
    isNumber,
    getTextBefore,
    getTextAfter,
    copyObjectTo,
    copyArrayTo
} from './common.js';

import {
    isPM,
    parseDateFromUnixTime,
    formatUnixTime,
    formatTime,
    getCurrentDateTime,
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
} from './datetime.js';

import {
    numericCurrencyToString
} from './currency.js';

import {
    getCategorizedAccounts,
    getAllFilteredAccountsBalance
} from '@/lib/account.js';

import logger from './logger.js';
import services from './services.js';

const apiNotFoundErrorCode = 100001;
const specifiedApiNotFoundErrors = {
    '/api/register.json': {
        message: 'User registration is disabled'
    }
};
const validatorErrorCode = 200000;
const parameterizedErrors = [
    {
        localeKey: 'parameter invalid',
        regex: /^parameter "(\w+)" is invalid$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter required',
        regex: /^parameter "(\w+)" is required$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter too large',
        regex: /^parameter "(\w+)" must be less than (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'number',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too long',
        regex: /^parameter "(\w+)" must be less than (\d+) characters$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too small',
        regex: /^parameter "(\w+)" must be more than (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'number',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too short',
        regex: /^parameter "(\w+)" must be more than (\d+) characters$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter length not equal',
        regex: /^parameter "(\w+)" length is not equal to (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter cannot be blank',
        regex: /^parameter "(\w+)" cannot be blank$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid username format',
        regex: /^parameter "(\w+)" is invalid username format$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid email format',
        regex: /^parameter "(\w+)" is invalid email format$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid currency',
        regex: /^parameter "(\w+)" is invalid currency$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid color',
        regex: /^parameter "(\w+)" is invalid color$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    }
];

function getAllLanguageInfos() {
    return allLanguages;
}

function getAllLanguageInfoArray(translateFn, includeSystemDefault) {
    const ret = [];

    for (const code in allLanguages) {
        if (!Object.prototype.hasOwnProperty.call(allLanguages, code)) {
            continue;
        }

        const languageInfo = allLanguages[code];

        ret.push({
            code: code,
            displayName: languageInfo.displayName
        });
    }

    ret.sort(function (lang1, lang2) {
        return lang1.code.localeCompare(lang2.code);
    });

    if (includeSystemDefault) {
        ret.splice(0, 0, {
            code: '',
            displayName: translateFn('System Default')
        });
    }

    return ret;
}

function getLanguageInfo(locale) {
    return allLanguages[locale];
}

function getDefaultLanguage() {
    if (!window || !window.navigator) {
        return defaultLanguage;
    }

    let browserLocale = window.navigator.browserLanguage || window.navigator.language;

    if (!browserLocale) {
        return defaultLanguage;
    }

    if (!allLanguages[browserLocale]) {
        const locale = getLocaleFromLanguageAlias(browserLocale);

        if (locale) {
            browserLocale = locale;
        }
    }

    if (!allLanguages[browserLocale] && browserLocale.split('-').length > 1) { // maybe language-script-region
        const localeParts = browserLocale.split('-');
        browserLocale = localeParts[0] + '-' + localeParts[1];

        if (!allLanguages[browserLocale]) {
            const locale = getLocaleFromLanguageAlias(browserLocale);

            if (locale) {
                browserLocale = locale;
            }
        }

        if (!allLanguages[browserLocale]) {
            browserLocale = localeParts[0];
            const locale = getLocaleFromLanguageAlias(browserLocale);

            if (locale) {
                browserLocale = locale;
            }
        }
    }

    if (!allLanguages[browserLocale]) {
        return defaultLanguage;
    }

    return browserLocale;
}

function getLocaleFromLanguageAlias(alias) {
    for (let locale in allLanguages) {
        if (!Object.prototype.hasOwnProperty.call(allLanguages, locale)) {
            continue;
        }

        if (locale.toLowerCase() === alias.toLowerCase()) {
            return locale;
        }

        const lang = allLanguages[locale];
        const aliases = lang.aliases;

        if (!aliases || aliases.length < 1) {
            continue;
        }

        for (let i = 0; i < aliases.length; i++) {
            if (aliases[i].toLowerCase() === alias.toLowerCase()) {
                return locale;
            }
        }
    }

    return null;
}

function getCurrentLanguageCode(i18nGlobal) {
    return i18nGlobal.locale;
}

function getCurrentLanguageInfo(i18nGlobal) {
    const locale = getLanguageInfo(i18nGlobal.locale);

    if (locale) {
        return locale;
    }

    return getLanguageInfo(getDefaultLanguage());
}

function getCurrentLanguageDisplayName(i18nGlobal) {
    const currentLanguageInfo = getCurrentLanguageInfo(i18nGlobal);
    return currentLanguageInfo.displayName;
}

function getDefaultCurrency(translateFn) {
    return translateFn('default.currency');
}

function getDefaultFirstDayOfWeek(translateFn) {
    return translateFn('default.firstDayOfWeek');
}

function getCurrencyName(currencyCode, translateFn) {
    return translateFn(`currency.${currencyCode}`);
}

function getAllMeridiemIndicatorNames(translateFn) {
    return [
        translateFn('datetime.AM.content'),
        translateFn('datetime.PM.content')
    ];
}

function getAllLongMonthNames(translateFn) {
    return [
        translateFn('datetime.January.long'),
        translateFn('datetime.February.long'),
        translateFn('datetime.March.long'),
        translateFn('datetime.April.long'),
        translateFn('datetime.May.long'),
        translateFn('datetime.June.long'),
        translateFn('datetime.July.long'),
        translateFn('datetime.August.long'),
        translateFn('datetime.September.long'),
        translateFn('datetime.October.long'),
        translateFn('datetime.November.long'),
        translateFn('datetime.December.long')
    ];
}

function getAllShortMonthNames(translateFn) {
    return [
        translateFn('datetime.January.short'),
        translateFn('datetime.February.short'),
        translateFn('datetime.March.short'),
        translateFn('datetime.April.short'),
        translateFn('datetime.May.short'),
        translateFn('datetime.June.short'),
        translateFn('datetime.July.short'),
        translateFn('datetime.August.short'),
        translateFn('datetime.September.short'),
        translateFn('datetime.October.short'),
        translateFn('datetime.November.short'),
        translateFn('datetime.December.short')
    ];
}

function getAllLongWeekdayNames(translateFn) {
    return [
        translateFn('datetime.Sunday.long'),
        translateFn('datetime.Monday.long'),
        translateFn('datetime.Tuesday.long'),
        translateFn('datetime.Wednesday.long'),
        translateFn('datetime.Thursday.long'),
        translateFn('datetime.Friday.long'),
        translateFn('datetime.Saturday.long')
    ];
}

function getAllShortWeekdayNames(translateFn) {
    return [
        translateFn('datetime.Sunday.short'),
        translateFn('datetime.Monday.short'),
        translateFn('datetime.Tuesday.short'),
        translateFn('datetime.Wednesday.short'),
        translateFn('datetime.Thursday.short'),
        translateFn('datetime.Friday.short'),
        translateFn('datetime.Saturday.short')
    ];
}

function getAllMinWeekdayNames(translateFn) {
    return [
        translateFn('datetime.Sunday.min'),
        translateFn('datetime.Monday.min'),
        translateFn('datetime.Tuesday.min'),
        translateFn('datetime.Wednesday.min'),
        translateFn('datetime.Thursday.min'),
        translateFn('datetime.Friday.min'),
        translateFn('datetime.Saturday.min')
    ];
}

function getAllLongDateFormats(translateFn) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormats(translateFn, datetime.allLongDateFormat, datetime.allLongDateFormatArray, 'format.longDate', defaultLongDateFormatTypeName, datetime.defaultLongDateFormat);
}

function getAllShortDateFormats(translateFn) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormats(translateFn, datetime.allShortDateFormat, datetime.allShortDateFormatArray, 'format.shortDate', defaultShortDateFormatTypeName, datetime.defaultShortDateFormat);
}

function getAllLongTimeFormats(translateFn) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    return getDateTimeFormats(translateFn, datetime.allLongTimeFormat, datetime.allLongTimeFormatArray, 'format.longTime', defaultLongTimeFormatTypeName, datetime.defaultLongTimeFormat);
}

function getAllShortTimeFormats(translateFn) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    return getDateTimeFormats(translateFn, datetime.allShortTimeFormat, datetime.allShortTimeFormatArray, 'format.shortTime', defaultShortTimeFormatTypeName, datetime.defaultShortTimeFormat);
}

function getMonthShortName(month, translateFn) {
    return translateFn(`datetime.${month}.short`);
}

function getMonthLongName(month, translateFn) {
    return translateFn(`datetime.${month}.long`);
}

function getWeekdayShortName(weekDay, translateFn) {
    return translateFn(`datetime.${weekDay}.short`);
}

function getWeekdayLongName(weekDay, translateFn) {
    return translateFn(`datetime.${weekDay}.long`);
}

function getI18nLongDateFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, datetime.allLongDateFormat, datetime.allLongDateFormatArray, 'format.longDate', defaultLongDateFormatTypeName, datetime.defaultLongDateFormat, formatTypeValue);
}

function getI18nShortDateFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, datetime.allShortDateFormat, datetime.allShortDateFormatArray, 'format.shortDate', defaultShortDateFormatTypeName, datetime.defaultShortDateFormat, formatTypeValue);
}

function getI18nLongYearFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, datetime.allLongDateFormat, datetime.allLongDateFormatArray, 'format.longYear', defaultLongDateFormatTypeName, datetime.defaultLongDateFormat, formatTypeValue);
}

function getI18nShortYearFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, datetime.allShortDateFormat, datetime.allShortDateFormatArray, 'format.shortYear', defaultShortDateFormatTypeName, datetime.defaultShortDateFormat, formatTypeValue);
}

function getI18nLongYearMonthFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, datetime.allLongDateFormat, datetime.allLongDateFormatArray, 'format.longYearMonth', defaultLongDateFormatTypeName, datetime.defaultLongDateFormat, formatTypeValue);
}

function getI18nShortYearMonthFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, datetime.allShortDateFormat, datetime.allShortDateFormatArray, 'format.shortYearMonth', defaultShortDateFormatTypeName, datetime.defaultShortDateFormat, formatTypeValue);
}

function getI18nLongMonthDayFormat(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormat(translateFn, datetime.allLongDateFormat, datetime.allLongDateFormatArray, 'format.longMonthDay', defaultLongDateFormatTypeName, datetime.defaultLongDateFormat, formatTypeValue);
}

function getI18nShortMonthDayFormat(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormat(translateFn, datetime.allShortDateFormat, datetime.allShortDateFormatArray, 'format.shortMonthDay', defaultShortDateFormatTypeName, datetime.defaultShortDateFormat, formatTypeValue);
}

function isLongDateMonthAfterYear(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    const type = getDateTimeFormatType(datetime.allLongDateFormat, datetime.allLongDateFormatArray, defaultLongDateFormatTypeName, datetime.defaultLongDateFormat, formatTypeValue);
    return type.isMonthAfterYear;
}

function isShortDateMonthAfterYear(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    const type = getDateTimeFormatType(datetime.allShortDateFormat, datetime.allShortDateFormatArray, defaultShortDateFormatTypeName, datetime.defaultShortDateFormat, formatTypeValue);
    return type.isMonthAfterYear;
}

function getI18nLongTimeFormat(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    return getDateTimeFormat(translateFn, datetime.allLongTimeFormat, datetime.allLongTimeFormatArray, 'format.longTime', defaultLongTimeFormatTypeName, datetime.defaultLongTimeFormat, formatTypeValue);
}

function getI18nShortTimeFormat(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    return getDateTimeFormat(translateFn, datetime.allShortTimeFormat, datetime.allShortTimeFormatArray, 'format.shortTime', defaultShortTimeFormatTypeName, datetime.defaultShortTimeFormat, formatTypeValue);
}

function isLongTime24HourFormat(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    const type = getDateTimeFormatType(datetime.allLongTimeFormat, datetime.allLongTimeFormatArray, defaultLongTimeFormatTypeName, datetime.defaultLongTimeFormat, formatTypeValue);
    return type.is24HourFormat;
}

function isLongTimeMeridiemIndicatorFirst(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    const type = getDateTimeFormatType(datetime.allLongTimeFormat, datetime.allLongTimeFormatArray, defaultLongTimeFormatTypeName, datetime.defaultLongTimeFormat, formatTypeValue);
    return type.isMeridiemIndicatorFirst;
}

function isShortTime24HourFormat(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    const type = getDateTimeFormatType(datetime.allShortTimeFormat, datetime.allShortTimeFormatArray, defaultShortTimeFormatTypeName, datetime.defaultShortTimeFormat, formatTypeValue);
    return type.is24HourFormat;
}

function isShortTimeMeridiemIndicatorFirst(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    const type = getDateTimeFormatType(datetime.allShortTimeFormat, datetime.allShortTimeFormatArray, defaultShortTimeFormatTypeName, datetime.defaultShortTimeFormat, formatTypeValue);
    return type.isMeridiemIndicatorFirst;
}

function getDateTimeFormats(translateFn, allFormatMap, allFormatArray, localeFormatPathPrefix, localeDefaultFormatTypeName, systemDefaultFormatType) {
    const defaultFormat = getDateTimeFormat(translateFn, allFormatMap, allFormatArray,
        localeFormatPathPrefix, localeDefaultFormatTypeName, systemDefaultFormatType, datetime.defaultDateTimeFormatValue);
    const ret = [];

    ret.push({
        type: datetime.defaultDateTimeFormatValue,
        format: defaultFormat,
        displayName: `${translateFn('Language Default')} (${formatTime(getCurrentDateTime(), defaultFormat)})`
    });

    for (let i = 0; i < allFormatArray.length; i++) {
        const formatType = allFormatArray[i];
        const format = translateFn(`${localeFormatPathPrefix}.${formatType.key}`);

        ret.push({
            type: formatType.type,
            format: format,
            displayName: formatTime(getCurrentDateTime(), format)
        });
    }

    return ret;
}

function getDateTimeFormat(translateFn, allFormatMap, allFormatArray, localeFormatPathPrefix, localeDefaultFormatTypeName, systemDefaultFormatType, formatTypeValue) {
    const type = getDateTimeFormatType(allFormatMap, allFormatArray,
        localeDefaultFormatTypeName, systemDefaultFormatType, formatTypeValue);
    return translateFn(`${localeFormatPathPrefix}.${type.key}`);
}

function getAllTimezones(includeSystemDefault, translateFn) {
    const defaultTimezoneOffset = getBrowserTimezoneOffset();
    const defaultTimezoneOffsetMinutes = getBrowserTimezoneOffsetMinutes();
    const allTimezones = timezone.all;
    const allTimezoneInfos = [];

    for (let i = 0; i < allTimezones.length; i++) {
        const utcOffset = (allTimezones[i].timezoneName !== timezone.utcTimezoneName ? getTimezoneOffset(allTimezones[i].timezoneName) : '');
        const displayName = translateFn(`timezone.${allTimezones[i].displayName}`);

        allTimezoneInfos.push({
            name: allTimezones[i].timezoneName,
            utcOffset: utcOffset,
            utcOffsetMinutes: getTimezoneOffsetMinutes(allTimezones[i].timezoneName),
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

function getAllCurrencies(translateFn) {
    const allCurrencyCodes = currency.all;
    const allCurrencies = [];

    for (let currencyCode in allCurrencyCodes) {
        if (!Object.prototype.hasOwnProperty.call(allCurrencyCodes, currencyCode)) {
            return;
        }

        allCurrencies.push({
            code: currencyCode,
            displayName: getCurrencyName(currencyCode, translateFn)
        });
    }

    allCurrencies.sort(function(c1, c2) {
        return c1.displayName.localeCompare(c2.displayName);
    })

    return allCurrencies;
}

function getAllWeekDays(translateFn) {
    const allWeekDays = [];

    for (let i = 0; i < datetime.allWeekDaysArray.length; i++) {
        const weekDay = datetime.allWeekDaysArray[i];

        allWeekDays.push({
            type: weekDay.type,
            displayName: translateFn(`datetime.${weekDay.name}.long`)
        });
    }

    return allWeekDays;
}

function getAllDateRanges(includeCustom, translateFn) {
    const allDateRanges = [];

    for (let dateRangeField in datetime.allDateRanges) {
        if (!Object.prototype.hasOwnProperty.call(datetime.allDateRanges, dateRangeField)) {
            continue;
        }

        const dateRangeType = datetime.allDateRanges[dateRangeField];

        if (includeCustom || dateRangeType.type !== datetime.allDateRanges.Custom.type) {
            allDateRanges.push({
                type: dateRangeType.type,
                displayName: translateFn(dateRangeType.name)
            });
        }
    }

    return allDateRanges;
}

function getAllRecentMonthDateRanges(userStore, includeAll, includeCustom, translateFn) {
    const allRecentMonthDateRanges = [];
    const recentDateRanges = getRecentMonthDateRanges(12);

    if (includeAll) {
        allRecentMonthDateRanges.push({
            dateType: datetime.allDateRanges.All.type,
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
            dateType: datetime.allDateRanges.Custom.type,
            minTime: 0,
            maxTime: 0,
            displayName: translateFn('Custom Date')
        });
    }

    return allRecentMonthDateRanges;
}

function getDateRangeDisplayName(userStore, dateType, startTime, endTime, translateFn) {
    if (dateType === datetime.allDateRanges.All.type) {
        return translateFn(datetime.allDateRanges.All.name);
    }

    for (let dateRangeField in datetime.allDateRanges) {
        if (!Object.prototype.hasOwnProperty.call(datetime.allDateRanges, dateRangeField)) {
            continue;
        }

        const dateRange = datetime.allDateRanges[dateRangeField];

        if (dateRange && dateRange.type !== datetime.allDateRanges.Custom.type && dateRange.type === dateType && dateRange.name) {
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

function getAllAccountCategories(translateFn) {
    const allAccountCategories = [];

    for (let i = 0; i < account.allCategories.length; i++) {
        const accountCategory = account.allCategories[i];

        allAccountCategories.push({
            id: accountCategory.id,
            displayName: translateFn(accountCategory.name),
            defaultAccountIconId: accountCategory.defaultAccountIconId
        });
    }

    return allAccountCategories;
}

function getAllAccountTypes(translateFn) {
    const allAccountTypes = [];

    for (let i = 0; i < account.allAccountTypesArray.length; i++) {
        const accountType = account.allAccountTypesArray[i];

        allAccountTypes.push({
            id: accountType.id,
            displayName: translateFn(accountType.name)
        });
    }

    return allAccountTypes;
}

function getAllStatisticsChartDataTypes(translateFn) {
    const allChartDataTypes = [];

    for (const dataTypeField in statistics.allChartDataTypes) {
        if (!Object.prototype.hasOwnProperty.call(statistics.allChartDataTypes, dataTypeField)) {
            return;
        }

        const chartDataType = statistics.allChartDataTypes[dataTypeField];

        allChartDataTypes.push({
            type: chartDataType.type,
            displayName: translateFn(chartDataType.name)
        });
    }

    return allChartDataTypes;
}

function getAllStatisticsSortingTypes(translateFn) {
    const allSortingTypes = [];

    for (const sortingTypeField in statistics.allSortingTypes) {
        if (!Object.prototype.hasOwnProperty.call(statistics.allSortingTypes, sortingTypeField)) {
            return;
        }

        const sortingType = statistics.allSortingTypes[sortingTypeField];

        allSortingTypes.push({
            type: sortingType.type,
            displayName: translateFn(sortingType.name),
            displayFullName: translateFn(sortingType.fullName)
        });
    }

    return allSortingTypes;
}

function getAllTransactionEditScopeTypes(translateFn) {
    return [{
        type: 0,
        displayName: translateFn('None')
    }, {
        type: 1,
        displayName: translateFn('All')
    }, {
        type: 2,
        displayName: translateFn('Today or later')
    }, {
        type: 3,
        displayName: translateFn('Recent 24 hours or later')
    }, {
        type: 4,
        displayName: translateFn('This week or later')
    }, {
        type: 5,
        displayName: translateFn('This month or later')
    }, {
        type: 6,
        displayName: translateFn('This year or later')
    }];
}

function getAllTransactionDefaultCategories(categoryType, locale, translateFn) {
    const allCategories = {};
    const categoryTypes = [];

    if (categoryType === 0) {
        for (let i = category.allCategoryTypes.Income; i <= category.allCategoryTypes.Transfer; i++) {
            categoryTypes.push(i);
        }
    } else {
        categoryTypes.push(categoryType);
    }

    for (let i = 0; i < categoryTypes.length; i++) {
        const categories = [];
        const categoryType = categoryTypes[i];
        let defaultCategories = [];

        if (categoryType === category.allCategoryTypes.Income) {
            defaultCategories = copyArrayTo(category.defaultIncomeCategories, []);
        } else if (categoryType === category.allCategoryTypes.Expense) {
            defaultCategories = copyArrayTo(category.defaultExpenseCategories, []);
        } else if (categoryType === category.allCategoryTypes.Transfer) {
            defaultCategories = copyArrayTo(category.defaultTransferCategories, []);
        }

        for (let j = 0; j < defaultCategories.length; j++) {
            const category = defaultCategories[j];

            const submitCategory = {
                name: translateFn('category.' + category.name, locale),
                type: categoryType,
                icon: category.categoryIconId,
                color: category.color,
                subCategories: []
            }

            for (let k = 0; k < category.subCategories.length; k++) {
                const subCategory = category.subCategories[k];
                submitCategory.subCategories.push({
                    name: translateFn('category.' + subCategory.name, locale),
                    type: categoryType,
                    icon: subCategory.categoryIconId,
                    color: subCategory.color
                });
            }

            categories.push(submitCategory);
        }

        allCategories[categoryType] = categories;
    }

    return allCategories;
}

function getAllDisplayExchangeRates(exchangeRatesData, translateFn) {
    if (!exchangeRatesData || !exchangeRatesData.exchangeRates) {
        return [];
    }

    const availableExchangeRates = [];

    for (let i = 0; i < exchangeRatesData.exchangeRates.length; i++) {
        const exchangeRate = exchangeRatesData.exchangeRates[i];

        availableExchangeRates.push({
            currencyCode: exchangeRate.currency,
            currencyDisplayName: getCurrencyName(exchangeRate.currency, translateFn),
            rate: exchangeRate.rate
        });
    }

    availableExchangeRates.sort(function(c1, c2) {
        return c1.currencyDisplayName.localeCompare(c2.currencyDisplayName);
    })

    return availableExchangeRates;
}

function getEnableDisableOptions(translateFn) {
    return [{
        value: true,
        displayName: translateFn('Enable')
    },{
        value: false,
        displayName: translateFn('Disable')
    }];
}

function getDisplayCurrency(value, currencyCode, options, translateFn) {
    if (!isNumber(value) && !isString(value)) {
        return value;
    }

    if (isNumber(value)) {
        value = value.toString();
    }

    if (!options) {
        options = {};
    }

    if (!options.notConvertValue) {
        const hasIncompleteFlag = isString(value) && value.charAt(value.length - 1) === '+';

        if (hasIncompleteFlag) {
            value = value.substring(0, value.length - 1);
        }

        value = numericCurrencyToString(value, options.enableThousandsSeparator);

        if (hasIncompleteFlag) {
            value = value + '+';
        }
    }

    const currencyDisplayMode = options.currencyDisplayMode;

    if (currencyCode && currencyDisplayMode === currency.allCurrencyDisplayModes.Symbol) {
        const currencyInfo = currency.all[currencyCode];
        let currencySymbol = currency.defaultCurrencySymbol;

        if (currencyInfo && currencyInfo.symbol) {
            currencySymbol = currencyInfo.symbol;
        } else if (currencyInfo && currencyInfo.code) {
            currencySymbol = currencyInfo.code;
        }

        return translateFn('format.currency.symbol', {
            amount: value,
            symbol: currencySymbol
        });
    } else if (currencyCode && currencyDisplayMode === currency.allCurrencyDisplayModes.Code) {
        return `${value} ${currencyCode}`;
    } else if (currencyCode && currencyDisplayMode === currency.allCurrencyDisplayModes.Name) {
        const currencyName = getCurrencyName(currencyCode, translateFn);
        return `${value} ${currencyName}`;
    } else {
        return value;
    }
}

function getDisplayCurrencyPrependAndAppendText(currencyCode, currencyDisplayMode, translateFn) {
    const options = {
        currencyDisplayMode: currencyDisplayMode,
        notConvertValue: true
    };

    const placeholder = '***';
    const finalText = getDisplayCurrency(placeholder, currencyCode, options, translateFn);

    if (!finalText) {
        return null;
    }

    let prependText = getTextBefore(finalText, placeholder);

    if (prependText) {
        prependText = prependText.trim();
    }

    let appendText = getTextAfter(finalText, placeholder);

    if (appendText) {
        appendText = appendText.trim();
    }

    return {
        prependText: prependText,
        appendText: appendText
    };
}

function getCategorizedAccountsWithDisplayBalance(exchangeRatesStore, allVisibleAccounts, showAccountBalance, defaultCurrency, options, translateFn) {
    const categorizedAccounts = copyObjectTo(getCategorizedAccounts(allVisibleAccounts), {});

    for (let category in categorizedAccounts) {
        if (!Object.prototype.hasOwnProperty.call(categorizedAccounts, category)) {
            continue;
        }

        const accountCategory = categorizedAccounts[category];

        if (accountCategory.accounts) {
            for (let i = 0; i < accountCategory.accounts.length; i++) {
                const account = accountCategory.accounts[i];

                if (showAccountBalance && account.isAsset) {
                    account.displayBalance = getDisplayCurrency(account.balance, account.currency, options, translateFn);
                } else if (showAccountBalance && account.isLiability) {
                    account.displayBalance = getDisplayCurrency(-account.balance, account.currency, options, translateFn);
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

            accountCategory.displayBalance = getDisplayCurrency(totalBalance, defaultCurrency, options, translateFn);
        } else {
            accountCategory.displayBalance = '***';
        }
    }

    return categorizedAccounts;
}

function joinMultiText(textArray, translateFn) {
    if (!textArray || !textArray.length) {
        return '';
    }

    const separator = translateFn('format.misc.multiTextJoinSeparator');

    return textArray.join(separator);
}

function getLocalizedError(error) {
    if (error.errorCode === apiNotFoundErrorCode && specifiedApiNotFoundErrors[error.path]) {
        return {
            message: `${specifiedApiNotFoundErrors[error.path].message}`
        };
    }

    if (error.errorCode !== validatorErrorCode) {
        return {
            message: `error.${error.errorMessage}`
        };
    }

    for (let i = 0; i < parameterizedErrors.length; i++) {
        const errorInfo = parameterizedErrors[i];
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

function setLanguage(i18nGlobal, locale, force) {
    if (!locale) {
        locale = getDefaultLanguage();
        logger.info(`No specified language, use browser default language ${locale}`);
    }

    if (!getLanguageInfo(locale)) {
        locale = getDefaultLanguage();
        logger.warn(`Not found language ${locale}, use browser default language ${locale}`);
    }

    if (!force && i18nGlobal.locale === locale) {
        logger.info(`Current locale is already ${locale}`);
        return null;
    }

    logger.info(`Apply current language to ${locale}`);

    i18nGlobal.locale = locale;
    moment.updateLocale(locale, {
        months : getAllLongMonthNames(i18nGlobal.t),
        monthsShort : getAllShortMonthNames(i18nGlobal.t),
        weekdays : getAllLongWeekdayNames(i18nGlobal.t),
        weekdaysShort : getAllShortWeekdayNames(i18nGlobal.t),
        weekdaysMin : getAllMinWeekdayNames(i18nGlobal.t),
        meridiem: function (hours) {
            if (isPM(hours)) {
                return i18nGlobal.t('datetime.PM.content');
            } else {
                return i18nGlobal.t('datetime.AM.content');
            }
        }
    });
    services.setLocale(locale);
    document.querySelector('html').setAttribute('lang', locale);

    const defaultCurrency = getDefaultCurrency(i18nGlobal.t);
    const defaultFirstDayOfWeekName = getDefaultFirstDayOfWeek(i18nGlobal.t);
    let defaultFirstDayOfWeek = datetime.defaultFirstDayOfWeek;

    if (datetime.allWeekDays[defaultFirstDayOfWeekName]) {
        defaultFirstDayOfWeek = datetime.allWeekDays[defaultFirstDayOfWeekName].type;
    }

    return {
        defaultCurrency: defaultCurrency,
        defaultFirstDayOfWeek: defaultFirstDayOfWeek
    };
}

function setTimeZone(timezone) {
    if (timezone) {
        moment.tz.setDefault(timezone);
    } else {
        moment.tz.setDefault();
    }
}

function initLocale(i18nGlobal, lastUserLanguage, timezone) {
    let localeDefaultSettings = null;

    if (lastUserLanguage && getLanguageInfo(lastUserLanguage)) {
        logger.info(`Last user language is ${lastUserLanguage}`);
        localeDefaultSettings = setLanguage(i18nGlobal, lastUserLanguage, true);
    } else {
        localeDefaultSettings = setLanguage(i18nGlobal, null, true);
    }

    if (timezone) {
        logger.info(`Current timezone is ${timezone}`);
        setTimeZone(timezone);
    } else {
        logger.info(`No timezone is set, use browser default ${getTimezoneOffset()} (maybe ${moment.tz.guess(true)})`);
    }

    return localeDefaultSettings;
}

export function getI18nOptions() {
    return {
        legacy: true,
        locale: defaultLanguage,
        fallbackLocale: defaultLanguage,
        formatFallbackMessages: true,
        messages: (function () {
            const messages = {};

            for (let locale in allLanguages) {
                if (!Object.prototype.hasOwnProperty.call(allLanguages, locale)) {
                    continue;
                }

                const lang = allLanguages[locale];
                messages[locale] = lang.content;
            }

            return messages;
        })()
    };
}

export function translateIf(text, isTranslate, translateFn) {
    if (isTranslate) {
        return translateFn(text);
    }

    return text;
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
        getAllLanguageInfos: getAllLanguageInfos,
        getAllLanguageInfoArray: (includeSystemDefault) => getAllLanguageInfoArray(i18nGlobal.t, includeSystemDefault),
        getLanguageInfo: getLanguageInfo,
        getDefaultLanguage: getDefaultLanguage,
        getCurrentLanguageCode: () => getCurrentLanguageCode(i18nGlobal),
        getCurrentLanguageInfo: () => getCurrentLanguageInfo(i18nGlobal),
        getCurrentLanguageDisplayName: () => getCurrentLanguageDisplayName(i18nGlobal),
        getDefaultCurrency: () => getDefaultCurrency(i18nGlobal.t),
        getDefaultFirstDayOfWeek: () => getDefaultFirstDayOfWeek(i18nGlobal.t),
        getCurrencyName: (currencyCode) => getCurrencyName(currencyCode, i18nGlobal.t),
        getAllMeridiemIndicatorNames: () => getAllMeridiemIndicatorNames(i18nGlobal.t),
        getAllLongMonthNames: () => getAllLongMonthNames(i18nGlobal.t),
        getAllShortMonthNames: () => getAllShortMonthNames(i18nGlobal.t),
        getAllLongWeekdayNames: () => getAllLongWeekdayNames(i18nGlobal.t),
        getAllShortWeekdayNames: () => getAllShortWeekdayNames(i18nGlobal.t),
        getAllMinWeekdayNames: () => getAllMinWeekdayNames(i18nGlobal.t),
        getAllLongDateFormats: () => getAllLongDateFormats(i18nGlobal.t),
        getAllShortDateFormats: () => getAllShortDateFormats(i18nGlobal.t),
        getAllLongTimeFormats: () => getAllLongTimeFormats(i18nGlobal.t),
        getAllShortTimeFormats: () => getAllShortTimeFormats(i18nGlobal.t),
        getMonthShortName: (month) => getMonthShortName(month, i18nGlobal.t),
        getMonthLongName: (month) => getMonthLongName(month, i18nGlobal.t),
        getWeekdayShortName: (weekDay) => getWeekdayShortName(weekDay, i18nGlobal.t),
        getWeekdayLongName: (weekDay) => getWeekdayLongName(weekDay, i18nGlobal.t),
        formatUnixTimeToLongDateTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongDateFormat(i18nGlobal.t, userStore.currentUserLongDateFormat) + ' ' + getI18nLongTimeFormat(i18nGlobal.t, userStore.currentUserLongTimeFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortDateTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortDateFormat(i18nGlobal.t, userStore.currentUserShortDateFormat) + ' ' + getI18nShortTimeFormat(i18nGlobal.t, userStore.currentUserShortTimeFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongDate: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongDateFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortDate: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortDateFormat(i18nGlobal.t, userStore.currentUserShortDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongYear: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongYearFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortYear: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortYearFormat(i18nGlobal.t, userStore.currentUserShortDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongYearMonth: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongYearMonthFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortYearMonth: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortYearMonthFormat(i18nGlobal.t, userStore.currentUserShortDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongMonthDay: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongMonthDayFormat(i18nGlobal.t, userStore.currentUserLongDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortMonthDay: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortMonthDayFormat(i18nGlobal.t, userStore.currentUserShortDateFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToLongTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nLongTimeFormat(i18nGlobal.t, userStore.currentUserLongTimeFormat), utcOffset, currentUtcOffset),
        formatUnixTimeToShortTime: (userStore, unixTime, utcOffset, currentUtcOffset) => formatUnixTime(unixTime, getI18nShortTimeFormat(i18nGlobal.t, userStore.currentUserShortTimeFormat), utcOffset, currentUtcOffset),
        formatTimeToLongYearMonth: (userStore, dateTime) => formatTime(dateTime, getI18nLongYearMonthFormat(i18nGlobal.t, userStore.currentUserLongDateFormat)),
        formatTimeToShortYearMonth: (userStore, dateTime) => formatTime(dateTime, getI18nShortYearMonthFormat(i18nGlobal.t, userStore.currentUserShortDateFormat)),
        isLongDateMonthAfterYear: (userStore) => isLongDateMonthAfterYear(i18nGlobal.t, userStore.currentUserLongDateFormat),
        isShortDateMonthAfterYear: (userStore) => isShortDateMonthAfterYear(i18nGlobal.t, userStore.currentUserShortDateFormat),
        isLongTime24HourFormat: (userStore) => isLongTime24HourFormat(i18nGlobal.t, userStore.currentUserLongTimeFormat),
        isLongTimeMeridiemIndicatorFirst: (userStore) => isLongTimeMeridiemIndicatorFirst(i18nGlobal.t, userStore.currentUserLongTimeFormat),
        isShortTime24HourFormat: (userStore) => isShortTime24HourFormat(i18nGlobal.t, userStore.currentUserShortTimeFormat),
        isShortTimeMeridiemIndicatorFirst: (userStore) => isShortTimeMeridiemIndicatorFirst(i18nGlobal.t, userStore.currentUserShortTimeFormat),
        getAllTimezones: (includeSystemDefault) => getAllTimezones(includeSystemDefault, i18nGlobal.t),
        getTimezoneDifferenceDisplayText: (utcOffset) => getTimezoneDifferenceDisplayText(utcOffset, i18nGlobal.t),
        getAllCurrencies: () => getAllCurrencies(i18nGlobal.t),
        getAllWeekDays: () => getAllWeekDays(i18nGlobal.t),
        getAllDateRanges: (includeCustom) => getAllDateRanges(includeCustom, i18nGlobal.t),
        getAllRecentMonthDateRanges: (userStore, includeAll, includeCustom) => getAllRecentMonthDateRanges(userStore, includeAll, includeCustom, i18nGlobal.t),
        getDateRangeDisplayName: (userStore, dateType, startTime, endTime) => getDateRangeDisplayName(userStore, dateType, startTime, endTime, i18nGlobal.t),
        getAllAccountCategories: () => getAllAccountCategories(i18nGlobal.t),
        getAllAccountTypes: () => getAllAccountTypes(i18nGlobal.t),
        getAllStatisticsChartDataTypes: () => getAllStatisticsChartDataTypes(i18nGlobal.t),
        getAllStatisticsSortingTypes: () => getAllStatisticsSortingTypes(i18nGlobal.t),
        getAllTransactionEditScopeTypes: () => getAllTransactionEditScopeTypes(i18nGlobal.t),
        getAllTransactionDefaultCategories: (categoryType, locale) => getAllTransactionDefaultCategories(categoryType, locale, i18nGlobal.t),
        getAllDisplayExchangeRates: (exchangeRatesData) => getAllDisplayExchangeRates(exchangeRatesData, i18nGlobal.t),
        getEnableDisableOptions: () => getEnableDisableOptions(i18nGlobal.t),
        getDisplayCurrency: (value, currencyCode, options) => getDisplayCurrency(value, currencyCode, options, i18nGlobal.t),
        getDisplayCurrencyPrependAndAppendText: (currencyCode, currencyDisplayMode) => getDisplayCurrencyPrependAndAppendText(currencyCode, currencyDisplayMode, i18nGlobal.t),
        getCategorizedAccountsWithDisplayBalance: (exchangeRatesStore, allVisibleAccounts, showAccountBalance, defaultCurrency, options) => getCategorizedAccountsWithDisplayBalance(exchangeRatesStore, allVisibleAccounts, showAccountBalance, defaultCurrency, options, i18nGlobal.t),
        joinMultiText: (textArray) => joinMultiText(textArray, i18nGlobal.t),
        setLanguage: (locale, force) => setLanguage(i18nGlobal, locale, force),
        setTimeZone: (timezone) => setTimeZone(timezone),
        initLocale: (lastUserLanguage, timezone) => initLocale(i18nGlobal, lastUserLanguage, timezone)
    };
}
