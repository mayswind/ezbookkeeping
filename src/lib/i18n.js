import { useI18n as useVueI18n } from 'vue-i18n';
import moment from 'moment-timezone';

import { defaultLanguage, allLanguages } from '@/locales/index.ts';

import { Month, WeekDay, MeridiemIndicator, LongDateFormat, ShortDateFormat, LongTimeFormat, ShortTimeFormat, DateRange, LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { DecimalSeparator, DigitGroupingSymbol, DigitGroupingType } from '@/core/numeral.ts';
import { CurrencyDisplayType, CurrencySortingType } from '@/core/currency.ts';
import { PresetAmountColor } from '@/core/color.ts';
import { AccountType, AccountCategory } from '@/core/account.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionEditScopeType, TransactionTagFilterType } from '@/core/transaction.ts';
import { ScheduledTemplateFrequencyType } from '@/core/template.ts';
import { CategoricalChartType, TrendChartType, ChartDataType, ChartSortingType, ChartDateAggregationType } from '@/core/statistics.ts';

import { UTC_TIMEZONE, ALL_TIMEZONES } from '@/consts/timezone.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { SUPPORTED_IMPORT_FILE_TYPES } from '@/consts/file.ts';
import { DEFAULT_EXPENSE_CATEGORIES, DEFAULT_INCOME_CATEGORIES, DEFAULT_TRANSFER_CATEGORIES } from '@/consts/category.ts';
import { KnownErrorCode, SPECIFIED_API_NOT_FOUND_ERRORS, PARAMETERIZED_ERRORS } from '@/consts/api.ts';

import {
    isString,
    isNumber,
    isBoolean,
    copyObjectTo,
    copyArrayTo
} from './common.ts';

import {
    isPM,
    parseDateFromUnixTime,
    formatUnixTime,
    formatCurrentTime,
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
} from './datetime.ts';

import {
    appendDigitGroupingSymbol,
    parseAmount,
    formatAmount,
    formatExchangeRateAmount,
    getAdaptiveDisplayAmountRate
} from './numeral.ts';

import {
    getCurrencyFraction,
    appendCurrencySymbol,
    getAmountPrependAndAppendCurrencySymbol
} from './currency.ts';

import {
    getCategorizedAccountsMap,
    getAllFilteredAccountsBalance
} from './account.js';

import logger from './logger.ts';
import services from './services.js';

function getLanguageDisplayName(translateFn, languageName) {
    return translateFn(`language.${languageName}`);
}

function getAllLanguageInfoArray(translateFn, includeSystemDefault) {
    const ret = [];

    for (const languageTag in allLanguages) {
        if (!Object.prototype.hasOwnProperty.call(allLanguages, languageTag)) {
            continue;
        }

        const languageInfo = allLanguages[languageTag];
        let displayName = languageInfo.displayName;
        let languageNameInCurrentLanguage = getLanguageDisplayName(translateFn, languageInfo.name);

        if (languageNameInCurrentLanguage && languageNameInCurrentLanguage !== displayName) {
            displayName = `${languageNameInCurrentLanguage} (${displayName})`;
        }

        ret.push({
            languageTag: languageTag,
            displayName: displayName
        });
    }

    ret.sort(function (lang1, lang2) {
        return lang1.languageTag.localeCompare(lang2.languageTag);
    });

    if (includeSystemDefault) {
        ret.splice(0, 0, {
            languageTag: '',
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

function getCurrentLanguageTag(i18nGlobal) {
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

function getDefaultCurrency(translateFn) {
    return translateFn('default.currency');
}

function getDefaultFirstDayOfWeek(translateFn) {
    return translateFn('default.firstDayOfWeek');
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

function getAllMeridiemIndicators(translateFn) {
    const allMeridiemIndicators = MeridiemIndicator.values();
    const meridiemIndicatorNames = [];
    const localizedMeridiemIndicatorNames = [];

    for (let i = 0; i < allMeridiemIndicators.length; i++) {
        const indicator = allMeridiemIndicators[i];

        meridiemIndicatorNames.push(indicator.name);
        localizedMeridiemIndicatorNames.push(translateFn(`datetime.${indicator.name}.content`));
    }

    return {
        values: meridiemIndicatorNames,
        displayValues: localizedMeridiemIndicatorNames
    };
}

function getAllLongMonthNames(translateFn) {
    const ret = [];
    const allMonths = Month.values();

    for (let i = 0; i < allMonths.length; i++) {
        const month = allMonths[i];
        ret.push(translateFn(`datetime.${month.name}.long`));
    }

    return ret;
}

function getAllShortMonthNames(translateFn) {
    const ret = [];
    const allMonths = Month.values();

    for (let i = 0; i < allMonths.length; i++) {
        const month = allMonths[i];
        ret.push(translateFn(`datetime.${month.name}.short`));
    }

    return ret;
}

function getAllLongWeekdayNames(translateFn) {
    const ret = [];
    const allWeekDays = WeekDay.values();

    for (let i = 0; i < allWeekDays.length; i++) {
        const weekDay = allWeekDays[i];
        ret.push(translateFn(`datetime.${weekDay.name}.long`));
    }

    return ret;
}

function getAllShortWeekdayNames(translateFn) {
    const ret = [];
    const allWeekDays = WeekDay.values();

    for (let i = 0; i < allWeekDays.length; i++) {
        const weekDay = allWeekDays[i];
        ret.push(translateFn(`datetime.${weekDay.name}.short`));
    }

    return ret;
}

function getAllMinWeekdayNames(translateFn) {
    const ret = [];
    const allWeekDays = WeekDay.values();

    for (let i = 0; i < allWeekDays.length; i++) {
        const weekDay = allWeekDays[i];
        ret.push(translateFn(`datetime.${weekDay.name}.min`));
    }

    return ret;
}

function getAllLongDateFormats(translateFn) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    return getDateTimeFormats(translateFn, LongDateFormat.all(), LongDateFormat.values(), 'format.longDate', defaultLongDateFormatTypeName, LongDateFormat.Default);
}

function getAllShortDateFormats(translateFn) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    return getDateTimeFormats(translateFn, ShortDateFormat.all(), ShortDateFormat.values(), 'format.shortDate', defaultShortDateFormatTypeName, ShortDateFormat.Default);
}

function getAllLongTimeFormats(translateFn) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    return getDateTimeFormats(translateFn, LongTimeFormat.values(), LongTimeFormat.values(), 'format.longTime', defaultLongTimeFormatTypeName, LongTimeFormat.Default);
}

function getAllShortTimeFormats(translateFn) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    return getDateTimeFormats(translateFn, ShortTimeFormat.values(), ShortTimeFormat.values(), 'format.shortTime', defaultShortTimeFormatTypeName, ShortTimeFormat.Default);
}

function getMonthShortName(monthName, translateFn) {
    return translateFn(`datetime.${monthName}.short`);
}

function getMonthLongName(monthName, translateFn) {
    return translateFn(`datetime.${monthName}.long`);
}

function getMonthdayOrdinal(monthDay, translateFn) {
    return translateFn(`datetime.monthDayOrdinal.${monthDay}`);
}

function getMonthdayShortName(monthDay, translateFn) {
    return translateFn('format.misc.monthDay', {
        ordinal: getMonthdayOrdinal(monthDay, translateFn)
    });
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

function isLongDateMonthAfterYear(translateFn, formatTypeValue) {
    const defaultLongDateFormatTypeName = translateFn('default.longDateFormat');
    const type = getDateTimeFormatType(LongDateFormat.all(), LongDateFormat.values(), defaultLongDateFormatTypeName, LongDateFormat.Default, formatTypeValue);
    return type.isMonthAfterYear;
}

function isShortDateMonthAfterYear(translateFn, formatTypeValue) {
    const defaultShortDateFormatTypeName = translateFn('default.shortDateFormat');
    const type = getDateTimeFormatType(ShortDateFormat.all(), ShortDateFormat.values(), defaultShortDateFormatTypeName, ShortDateFormat.Default, formatTypeValue);
    return type.isMonthAfterYear;
}

function getI18nLongTimeFormat(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    return getDateTimeFormat(translateFn, LongTimeFormat.values(), LongTimeFormat.values(), 'format.longTime', defaultLongTimeFormatTypeName, LongTimeFormat.Default, formatTypeValue);
}

function getI18nShortTimeFormat(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    return getDateTimeFormat(translateFn, ShortTimeFormat.values(), ShortTimeFormat.values(), 'format.shortTime', defaultShortTimeFormatTypeName, ShortTimeFormat.Default, formatTypeValue);
}

function formatYearQuarter(translateFn, year, quarter) {
    if (1 <= quarter && quarter <= 4) {
        return translateFn('format.yearQuarter.q' + quarter, {
            year: year,
            quarter: quarter
        });
    } else {
        return '';
    }
}

function isLongTime24HourFormat(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    const type = getDateTimeFormatType(LongTimeFormat.all(), LongTimeFormat.values(), defaultLongTimeFormatTypeName, LongTimeFormat.Default, formatTypeValue);
    return type.is24HourFormat;
}

function isLongTimeMeridiemIndicatorFirst(translateFn, formatTypeValue) {
    const defaultLongTimeFormatTypeName = translateFn('default.longTimeFormat');
    const type = getDateTimeFormatType(LongTimeFormat.all(), LongTimeFormat.values(), defaultLongTimeFormatTypeName, LongTimeFormat.Default, formatTypeValue);
    return type.isMeridiemIndicatorFirst;
}

function isShortTime24HourFormat(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    const type = getDateTimeFormatType(ShortTimeFormat.all(), ShortTimeFormat.values(), defaultShortTimeFormatTypeName, ShortTimeFormat.Default, formatTypeValue);
    return type.is24HourFormat;
}

function isShortTimeMeridiemIndicatorFirst(translateFn, formatTypeValue) {
    const defaultShortTimeFormatTypeName = translateFn('default.shortTimeFormat');
    const type = getDateTimeFormatType(ShortTimeFormat.all(), ShortTimeFormat.values(), defaultShortTimeFormatTypeName, ShortTimeFormat.Default, formatTypeValue);
    return type.isMeridiemIndicatorFirst;
}

function getDateTimeFormats(translateFn, allFormatMap, allFormatArray, localeFormatPathPrefix, localeDefaultFormatTypeName, systemDefaultFormatType) {
    const defaultFormat = getDateTimeFormat(translateFn, allFormatMap, allFormatArray,
        localeFormatPathPrefix, localeDefaultFormatTypeName, systemDefaultFormatType, LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE);
    const ret = [];

    ret.push({
        type: LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE,
        format: defaultFormat,
        displayName: `${translateFn('Language Default')} (${formatCurrentTime(defaultFormat)})`
    });

    for (let i = 0; i < allFormatArray.length; i++) {
        const formatType = allFormatArray[i];
        const format = translateFn(`${localeFormatPathPrefix}.${formatType.key}`);

        ret.push({
            type: formatType.type,
            format: format,
            displayName: formatCurrentTime(format)
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

function getAllCurrencies(translateFn) {
    const allCurrencies = [];

    for (let currencyCode in ALL_CURRENCIES) {
        if (!Object.prototype.hasOwnProperty.call(ALL_CURRENCIES, currencyCode)) {
            continue;
        }

        allCurrencies.push({
            currencyCode: currencyCode,
            displayName: getCurrencyName(currencyCode, translateFn)
        });
    }

    allCurrencies.sort(function(c1, c2) {
        return c1.displayName.localeCompare(c2.displayName);
    })

    return allCurrencies;
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

function getAllTimezoneTypesUsedForStatistics(currentTimezone, translateFn) {
    const currentTimezoneOffset = getTimezoneOffset(currentTimezone);

    return [
        {
            displayName: translateFn(TimezoneTypeForStatistics.ApplicationTimezone.name) + ` (UTC${currentTimezoneOffset})`,
            type: TimezoneTypeForStatistics.ApplicationTimezone.type
        },
        {
            displayName: translateFn(TimezoneTypeForStatistics.TransactionTimezone.name),
            type: TimezoneTypeForStatistics.TransactionTimezone.type
        }
    ];
}

function getAllDecimalSeparators(translateFn) {
    const defaultDecimalSeparatorTypeName = translateFn('default.decimalSeparator');
    return getNumeralSeparatorFormats(translateFn, DecimalSeparator.values(), DecimalSeparator.parse(defaultDecimalSeparatorTypeName), DecimalSeparator.Default, DecimalSeparator.LanguageDefaultType);
}

function getAllDigitGroupingSymbols(translateFn) {
    const defaultDigitGroupingSymbolTypeName = translateFn('default.digitGroupingSymbol');
    return getNumeralSeparatorFormats(translateFn, DigitGroupingSymbol.values(), DigitGroupingSymbol.parse(defaultDigitGroupingSymbolTypeName), DigitGroupingSymbol.Default, DigitGroupingSymbol.LanguageDefaultType);
}

function getNumeralSeparatorFormats(translateFn, allSeparatorArray, localeDefaultType, systemDefaultType, languageDefaultValue) {
    let defaultSeparatorType = localeDefaultType;

    if (!defaultSeparatorType) {
        defaultSeparatorType = systemDefaultType;
    }

    const ret = [];

    ret.push({
        type: languageDefaultValue,
        symbol: defaultSeparatorType.symbol,
        displayName: `${translateFn('Language Default')} (${defaultSeparatorType.symbol})`
    });

    for (let i = 0; i < allSeparatorArray.length; i++) {
        const type = allSeparatorArray[i];

        ret.push({
            type: type.type,
            symbol: type.symbol,
            displayName: `${translateFn('numeral.' + type.name)} (${type.symbol})`
        });
    }

    return ret;
}

function getAllDigitGroupingTypes(translateFn) {
    const defaultDigitGroupingTypeName = translateFn('default.digitGrouping');
    let defaultDigitGroupingType = DigitGroupingType.parse(defaultDigitGroupingTypeName);

    if (!defaultDigitGroupingType) {
        defaultDigitGroupingType = DigitGroupingType.Default;
    }

    const ret = [];

    ret.push({
        type: DigitGroupingType.LanguageDefaultType,
        enabled: defaultDigitGroupingType.enabled,
        displayName: `${translateFn('Language Default')} (${translateFn('numeral.' + defaultDigitGroupingType.name)})`
    });

    const allDigitGroupingTypes = DigitGroupingType.values();

    for (let i = 0; i < allDigitGroupingTypes.length; i++) {
        const type = allDigitGroupingTypes[i];

        ret.push({
            type: type.type,
            enabled: type.enabled,
            displayName: translateFn('numeral.' + type.name)
        });
    }

    return ret;
}

function getAllCurrencyDisplayTypes(userStore, settingsStore, translateFn) {
    const defaultCurrencyDisplayTypeName = translateFn('default.currencyDisplayType');
    let defaultCurrencyDisplayType = CurrencyDisplayType.parse(defaultCurrencyDisplayTypeName);

    if (!defaultCurrencyDisplayType) {
        defaultCurrencyDisplayType = CurrencyDisplayType.Default;
    }

    const defaultCurrency = userStore.currentUserDefaultCurrency;

    const ret = [];
    const defaultSampleValue = getFormattedAmountWithCurrency(12345, defaultCurrency, translateFn, userStore, settingsStore, false, defaultCurrencyDisplayType);

    ret.push({
        type: CurrencyDisplayType.LanguageDefaultType,
        displayName: `${translateFn('Language Default')} (${defaultSampleValue})`
    });

    const allCurrencyDisplayTypes = CurrencyDisplayType.values();

    for (let i = 0; i < allCurrencyDisplayTypes.length; i++) {
        const type = allCurrencyDisplayTypes[i];
        const sampleValue = getFormattedAmountWithCurrency(12345, defaultCurrency, translateFn, userStore, settingsStore, false, type);
        const displayName = `${translateFn(type.name)} (${sampleValue})`

        ret.push({
            type: type.type,
            displayName: displayName
        });
    }

    return ret;
}

function getAllCurrencySortingTypes(translateFn) {
    return getLocalizedDisplayNameAndType(CurrencySortingType.values(), translateFn);
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

function getNumberWithDigitGroupingSymbol(value, translateFn, userStore) {
    const numberFormatOptions = getNumberFormatOptions(translateFn, userStore);
    return appendDigitGroupingSymbol(value, numberFormatOptions);
}

function getParsedAmountNumber(value, translateFn, userStore) {
    const numberFormatOptions = getNumberFormatOptions(translateFn, userStore);
    return parseAmount(value, numberFormatOptions);
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

function getFormattedExchangeRateAmount(value, translateFn, userStore) {
    const numberFormatOptions = getNumberFormatOptions(translateFn, userStore);
    return formatExchangeRateAmount(value, numberFormatOptions);
}

function getAdaptiveAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, translateFn, userStore) {
    const numberFormatOptions = getNumberFormatOptions(translateFn, userStore);
    return getAdaptiveDisplayAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, numberFormatOptions);
}

function getAmountPrependAndAppendText(currencyCode, userStore, settingsStore, isPlural, translateFn) {
    const currencyDisplayType = getCurrentCurrencyDisplayType(translateFn, userStore);
    const currencyUnit = getCurrencyUnitName(currencyCode, isPlural, translateFn);
    const currencyName = getCurrencyName(currencyCode, translateFn);
    return getAmountPrependAndAppendCurrencySymbol(currencyDisplayType, currencyCode, currencyUnit, currencyName, isPlural);
}

function getAllExpenseIncomeAmountColors(translateFn, expenseOrIncome) {
    const ret = [];
    let defaultAmountName = '';

    if (expenseOrIncome === 1) { // expense
        defaultAmountName = PresetAmountColor.DefaultExpenseColor.name;
    } else if (expenseOrIncome === 2) { // income
        defaultAmountName = PresetAmountColor.DefaultIncomeColor.name;
    }

    if (defaultAmountName) {
        defaultAmountName = translateFn('color.amount.' + defaultAmountName);
    }

    ret.push({
        type: PresetAmountColor.SystemDefaultType,
        displayName: translateFn('System Default') + (defaultAmountName ? ` (${defaultAmountName})` : '')
    });

    const allPresetAmountColors = PresetAmountColor.values();

    for (let i = 0; i < allPresetAmountColors.length; i++) {
        const amountColor = allPresetAmountColors[i];

        ret.push({
            type: amountColor.type,
            displayName: translateFn('color.amount.' + amountColor.name)
        });
    }

    return ret;
}

function getAllAccountCategories(translateFn) {
    const ret = [];
    const allCategories = AccountCategory.values();

    for (let i = 0; i < allCategories.length; i++) {
        const accountCategory = allCategories[i];

        ret.push({
            type: accountCategory.type,
            displayName: translateFn(accountCategory.name),
            defaultAccountIconId: accountCategory.defaultAccountIconId
        });
    }

    return ret;
}

function getAllAccountTypes(translateFn) {
    return getLocalizedDisplayNameAndType(AccountType.values(), translateFn);
}

function getAllCategoricalChartTypes(translateFn) {
    return getLocalizedDisplayNameAndType(CategoricalChartType.values(), translateFn);
}

function getAllTrendChartTypes(translateFn) {
    return getLocalizedDisplayNameAndType(TrendChartType.values(), translateFn);
}

function getAllStatisticsChartDataTypes(translateFn, analysisType) {
    return getLocalizedDisplayNameAndType(ChartDataType.values(analysisType), translateFn);
}

function getAllStatisticsSortingTypes(translateFn) {
    return getLocalizedDisplayNameAndType(ChartSortingType.values(), translateFn);
}

function getAllStatisticsDateAggregationTypes(translateFn) {
    return getLocalizedDisplayNameAndType(ChartDateAggregationType.values(), translateFn);
}

function getAllTransactionEditScopeTypes(translateFn) {
    return getLocalizedDisplayNameAndType(TransactionEditScopeType.values(), translateFn);
}

function getAllTransactionTagFilterTypes(translateFn) {
    return getLocalizedDisplayNameAndType(TransactionTagFilterType.values(), translateFn);
}

function getAllTransactionScheduledFrequencyTypes(translateFn) {
    return getLocalizedDisplayNameAndType(ScheduledTemplateFrequencyType.values(), translateFn);
}

function getAllTransactionDefaultCategories(categoryType, locale, translateFn) {
    const allCategories = {};
    const categoryTypes = [];

    if (categoryType === 0) {
        for (let i = CategoryType.Income; i <= CategoryType.Transfer; i++) {
            categoryTypes.push(i);
        }
    } else {
        categoryTypes.push(categoryType);
    }

    for (let i = 0; i < categoryTypes.length; i++) {
        const categories = [];
        const categoryType = categoryTypes[i];
        let defaultCategories = [];

        if (categoryType === CategoryType.Income) {
            defaultCategories = copyArrayTo(DEFAULT_INCOME_CATEGORIES, []);
        } else if (categoryType === CategoryType.Expense) {
            defaultCategories = copyArrayTo(DEFAULT_EXPENSE_CATEGORIES, []);
        } else if (categoryType === CategoryType.Transfer) {
            defaultCategories = copyArrayTo(DEFAULT_TRANSFER_CATEGORIES, []);
        }

        for (let j = 0; j < defaultCategories.length; j++) {
            const category = defaultCategories[j];

            const submitCategory = {
                name: translateFn('category.' + category.name, {}, { locale: locale }),
                type: categoryType,
                icon: category.categoryIconId,
                color: category.color,
                subCategories: []
            }

            for (let k = 0; k < category.subCategories.length; k++) {
                const subCategory = category.subCategories[k];
                submitCategory.subCategories.push({
                    name: translateFn('category.' + subCategory.name, {}, { locale: locale }),
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

function getAllDisplayExchangeRates(settingsStore, exchangeRatesData, translateFn) {
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

    if (settingsStore.appSettings.currencySortByInExchangeRatesPage === CurrencySortingType.Name.type) {
        availableExchangeRates.sort(function(c1, c2) {
            return c1.currencyDisplayName.localeCompare(c2.currencyDisplayName);
        });
    } else if (settingsStore.appSettings.currencySortByInExchangeRatesPage === CurrencySortingType.CurrencyCode.type) {
        availableExchangeRates.sort(function(c1, c2) {
            return c1.currencyCode.localeCompare(c2.currencyCode);
        });
    } else if (settingsStore.appSettings.currencySortByInExchangeRatesPage === CurrencySortingType.ExchangeRate.type) {
        availableExchangeRates.sort(function(c1, c2) {
            const rate1 = parseFloat(c1.rate);
            const rate2 = parseFloat(c2.rate);

            if (rate1 > rate2) {
                return 1;
            } else if (rate1 < rate2) {
                return -1;
            } else {
                return 0;
            }
        });
    }

    return availableExchangeRates;
}

function getAllSupportedImportFileTypes(i18nGlobal, translateFn) {
    const allSupportedImportFileTypes = [];

    for (let i = 0; i < SUPPORTED_IMPORT_FILE_TYPES.length; i++) {
        const fileType = SUPPORTED_IMPORT_FILE_TYPES[i];
        let document = {
            language: '',
            displayLanguageName: '',
            anchor: ''
        };

        if (fileType.document) {
            if (fileType.document.supportMultiLanguages === true) {
                document.language = getCurrentLanguageTag(i18nGlobal);
                document.anchor = translateFn(`document.anchor.export_and_import.${fileType.document.anchor}`);
            } else if (isString(fileType.document.supportMultiLanguages) && allLanguages[fileType.document.supportMultiLanguages]) {
                document.language = fileType.document.supportMultiLanguages;

                if (document.language !== getCurrentLanguageTag(i18nGlobal)) {
                    document.displayLanguageName = getLanguageDisplayName(translateFn, allLanguages[fileType.document.supportMultiLanguages].name);
                }

                document.anchor = fileType.document.anchor;
            }

            if (document.language) {
                document.language = document.language.replace(/-/g, '_');
            }

            if (document.anchor) {
                document.anchor = document.anchor.toLowerCase().replace(/ /g, '-');
            }

            if (document.language === defaultLanguage) {
                document.language = '';
            }
        } else {
            document = null;
        }

        const subTypes = [];

        if (fileType.subTypes) {
            for (let i = 0; i < fileType.subTypes.length; i++) {
                const subType = fileType.subTypes[i];

                subTypes.push({
                    type: subType.type,
                    displayName: translateFn(subType.name),
                    extensions: subType.extensions
                });
            }
        }

        allSupportedImportFileTypes.push({
            type: fileType.type,
            displayName: translateFn(fileType.name),
            extensions: fileType.extensions,
            subTypes: subTypes.length ? subTypes : undefined,
            document: document
        });
    }

    return allSupportedImportFileTypes;
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

function getServerTipContent(tipConfig, i18nGlobal) {
    if (!tipConfig) {
        return '';
    }

    const currentLanguage = getCurrentLanguageTag(i18nGlobal);

    if (tipConfig[currentLanguage]) {
        return tipConfig[currentLanguage];
    }

    return tipConfig.default || '';
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
                return i18nGlobal.t(`datetime.${MeridiemIndicator.PM.name}.content`);
            } else {
                return i18nGlobal.t(`datetime.${MeridiemIndicator.AM.name}.content`);
            }
        }
    });

    services.setLocale(locale);
    document.querySelector('html').setAttribute('lang', locale);

    const defaultCurrency = getDefaultCurrency(i18nGlobal.t);
    const defaultFirstDayOfWeekName = getDefaultFirstDayOfWeek(i18nGlobal.t);
    let defaultFirstDayOfWeek = WeekDay.DefaultFirstDay.type;

    if (WeekDay.parse(defaultFirstDayOfWeekName)) {
        defaultFirstDayOfWeek = WeekDay.parse(defaultFirstDayOfWeekName).type;
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
        getAllLanguageInfoArray: (includeSystemDefault) => getAllLanguageInfoArray(i18nGlobal.t, includeSystemDefault),
        getLanguageInfo: getLanguageInfo,
        getDefaultLanguage: getDefaultLanguage,
        getCurrentLanguageTag: () => getCurrentLanguageTag(i18nGlobal),
        getCurrentLanguageInfo: () => getCurrentLanguageInfo(i18nGlobal),
        getCurrentLanguageDisplayName: () => getCurrentLanguageDisplayName(i18nGlobal),
        getDefaultCurrency: () => getDefaultCurrency(i18nGlobal.t),
        getDefaultFirstDayOfWeek: () => getDefaultFirstDayOfWeek(i18nGlobal.t),
        getCurrencyName: (currencyCode) => getCurrencyName(currencyCode, i18nGlobal.t),
        getAllMeridiemIndicators: () => getAllMeridiemIndicators(i18nGlobal.t),
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
        getMonthdayOrdinal: (monthDay) => getMonthdayOrdinal(monthDay, i18nGlobal.t),
        getMonthdayShortName: (monthDay) => getMonthdayShortName(monthDay, i18nGlobal.t),
        getWeekdayShortName: (weekDay) => getWeekdayShortName(weekDay, i18nGlobal.t),
        getWeekdayLongName: (weekDay) => getWeekdayLongName(weekDay, i18nGlobal.t),
        getMultiMonthdayShortNames: (monthdays) => getMultiMonthdayShortNames(monthdays, i18nGlobal.t),
        getMultiWeekdayLongNames: (weekdayTypes, firstDayOfWeek) => getMultiWeekdayLongNames(weekdayTypes, firstDayOfWeek, i18nGlobal.t),
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
        formatYearQuarter: (year, quarter) => formatYearQuarter(i18nGlobal.t, year, quarter),
        isLongDateMonthAfterYear: (userStore) => isLongDateMonthAfterYear(i18nGlobal.t, userStore.currentUserLongDateFormat),
        isShortDateMonthAfterYear: (userStore) => isShortDateMonthAfterYear(i18nGlobal.t, userStore.currentUserShortDateFormat),
        isLongTime24HourFormat: (userStore) => isLongTime24HourFormat(i18nGlobal.t, userStore.currentUserLongTimeFormat),
        isLongTimeMeridiemIndicatorFirst: (userStore) => isLongTimeMeridiemIndicatorFirst(i18nGlobal.t, userStore.currentUserLongTimeFormat),
        isShortTime24HourFormat: (userStore) => isShortTime24HourFormat(i18nGlobal.t, userStore.currentUserShortTimeFormat),
        isShortTimeMeridiemIndicatorFirst: (userStore) => isShortTimeMeridiemIndicatorFirst(i18nGlobal.t, userStore.currentUserShortTimeFormat),
        getAllTimezones: (includeSystemDefault) => getAllTimezones(includeSystemDefault, i18nGlobal.t),
        getTimezoneDifferenceDisplayText: (utcOffset) => getTimezoneDifferenceDisplayText(utcOffset, i18nGlobal.t),
        getAllCurrencies: () => getAllCurrencies(i18nGlobal.t),
        getAllWeekDays: (firstDayOfWeek) => getAllWeekDays(firstDayOfWeek, i18nGlobal.t),
        getAllDateRanges: (scene, includeCustom, includeBillingCycle) => getAllDateRanges(scene, includeCustom, includeBillingCycle, i18nGlobal.t),
        getAllRecentMonthDateRanges: (userStore, includeAll, includeCustom) => getAllRecentMonthDateRanges(userStore, includeAll, includeCustom, i18nGlobal.t),
        getDateRangeDisplayName: (userStore, dateType, startTime, endTime) => getDateRangeDisplayName(userStore, dateType, startTime, endTime, i18nGlobal.t),
        getAllTimezoneTypesUsedForStatistics: (currentTimezone) => getAllTimezoneTypesUsedForStatistics(currentTimezone, i18nGlobal.t),
        getAllDecimalSeparators: () => getAllDecimalSeparators(i18nGlobal.t),
        getAllDigitGroupingSymbols: () => getAllDigitGroupingSymbols(i18nGlobal.t),
        getAllDigitGroupingTypes: () => getAllDigitGroupingTypes(i18nGlobal.t),
        getAllCurrencyDisplayTypes: (settingsStore, userStore) => getAllCurrencyDisplayTypes(userStore, settingsStore, i18nGlobal.t),
        getAllCurrencySortingTypes: () => getAllCurrencySortingTypes(i18nGlobal.t),
        getCurrentDecimalSeparator: (userStore) => getCurrentDecimalSeparator(i18nGlobal.t, userStore.currentUserDecimalSeparator),
        getCurrentDigitGroupingSymbol: (userStore) => getCurrentDigitGroupingSymbol(i18nGlobal.t, userStore.currentUserDigitGroupingSymbol),
        getCurrentDigitGroupingType: (userStore) => getCurrentDigitGroupingType(i18nGlobal.t, userStore.currentUserDigitGrouping),
        appendDigitGroupingSymbol: (userStore, value) => getNumberWithDigitGroupingSymbol(value, i18nGlobal.t, userStore),
        parseAmount: (userStore, value) => getParsedAmountNumber(value, i18nGlobal.t, userStore),
        formatAmount: (userStore, value, currencyCode) => getFormattedAmount(value, i18nGlobal.t, userStore, currencyCode),
        formatAmountWithCurrency: (settingsStore, userStore, value, currencyCode) => getFormattedAmountWithCurrency(value, currencyCode, i18nGlobal.t, userStore, settingsStore),
        formatExchangeRateAmount: (userStore, value) => getFormattedExchangeRateAmount(value, i18nGlobal.t, userStore),
        getAdaptiveAmountRate: (userStore, amount1, amount2, fromExchangeRate, toExchangeRate) => getAdaptiveAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, i18nGlobal.t, userStore),
        getAmountPrependAndAppendText: (settingsStore, userStore, currencyCode, isPlural) => getAmountPrependAndAppendText(currencyCode, userStore, settingsStore, isPlural, i18nGlobal.t),
        getAllExpenseAmountColors: () => getAllExpenseIncomeAmountColors(i18nGlobal.t, 1),
        getAllIncomeAmountColors: () => getAllExpenseIncomeAmountColors(i18nGlobal.t, 2),
        getAllAccountCategories: () => getAllAccountCategories(i18nGlobal.t),
        getAllAccountTypes: () => getAllAccountTypes(i18nGlobal.t),
        getAllCategoricalChartTypes: () => getAllCategoricalChartTypes(i18nGlobal.t),
        getAllTrendChartTypes: () => getAllTrendChartTypes(i18nGlobal.t),
        getAllStatisticsChartDataTypes: (analysisType) => getAllStatisticsChartDataTypes(i18nGlobal.t, analysisType),
        getAllStatisticsSortingTypes: () => getAllStatisticsSortingTypes(i18nGlobal.t),
        getAllStatisticsDateAggregationTypes: () => getAllStatisticsDateAggregationTypes(i18nGlobal.t),
        getAllTransactionEditScopeTypes: () => getAllTransactionEditScopeTypes(i18nGlobal.t),
        getAllTransactionTagFilterTypes: () => getAllTransactionTagFilterTypes(i18nGlobal.t),
        getAllTransactionScheduledFrequencyTypes: () => getAllTransactionScheduledFrequencyTypes(i18nGlobal.t),
        getAllTransactionDefaultCategories: (categoryType, locale) => getAllTransactionDefaultCategories(categoryType, locale, i18nGlobal.t),
        getAllDisplayExchangeRates: (settingsStore, exchangeRatesData) => getAllDisplayExchangeRates(settingsStore, exchangeRatesData, i18nGlobal.t),
        getAllSupportedImportFileTypes: () => getAllSupportedImportFileTypes(i18nGlobal, i18nGlobal.t),
        getEnableDisableOptions: () => getEnableDisableOptions(i18nGlobal.t),
        getCategorizedAccountsWithDisplayBalance: (allVisibleAccounts, showAccountBalance, defaultCurrency, settingsStore, userStore, exchangeRatesStore) => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts, showAccountBalance, defaultCurrency, userStore, settingsStore, exchangeRatesStore, i18nGlobal.t),
        getServerTipContent: (tipConfig) => getServerTipContent(tipConfig, i18nGlobal),
        joinMultiText: (textArray) => joinMultiText(textArray, i18nGlobal.t),
        tt: (...args) => i18nGlobal.t(...args),
        ti: (text, isTranslate) => translateIf(text, isTranslate, i18nGlobal.t),
        te: (message) => translateError(message, i18nGlobal.t),
        setLanguage: (locale, force) => setLanguage(i18nGlobal, locale, force),
        setTimeZone: (timezone) => setTimeZone(timezone),
        initLocale: (lastUserLanguage, timezone) => initLocale(i18nGlobal, lastUserLanguage, timezone)
    };
}

export function useI18n() {
    const i18nGlobal = useVueI18n();
    return i18nFunctions(i18nGlobal);
}
