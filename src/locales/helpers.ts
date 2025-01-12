import { useI18n as useVueI18n } from 'vue-i18n';
import moment from 'moment-timezone';

import type {TypeAndName, TypeAndDisplayName, LocalizedSwitchOption } from '@/core/base.ts';

import { type LanguageInfo, allLanguages, DEFAULT_LANGUAGE } from '@/locales/index.ts';

import {
    type DateFormat,
    type TimeFormat,
    type LocalizedMeridiemIndicator,
    type LocalizedDateRange,
    type LocalizedRecentMonthDateRange,
    Month,
    WeekDay,
    MeridiemIndicator,
    LongDateFormat,
    ShortDateFormat,
    LongTimeFormat,
    ShortTimeFormat,
    DateRange,
    DateRangeScene
} from '@/core/datetime.ts';

import {
    type LocalizedTimezoneInfo,
    TimezoneTypeForStatistics
} from '@/core/timezone.ts';

import {
    type NumberFormatOptions,
    type NumeralSymbolType,
    type LocalizedNumeralSymbolType,
    type LocalizedDigitGroupingType,
    DecimalSeparator,
    DigitGroupingSymbol,
    DigitGroupingType
} from '@/core/numeral.ts';

import {
    type CurrencyPrependAndAppendText,
    CurrencyDisplayType,
    CurrencySortingType
} from '@/core/currency.ts';

import {
    type LocalizedAccountCategory,
    AccountType,
    AccountCategory
} from '@/core/account.ts';

import {
    TransactionEditScopeType,
    TransactionTagFilterType
} from '@/core/transaction.ts';

import {
    ScheduledTemplateFrequencyType
} from '@/core/template.ts';

import {
    StatisticsAnalysisType,
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import type { LocaleDefaultSettings } from '@/core/setting.ts';
import type { ErrorResponse } from '@/core/api.ts';

import { UTC_TIMEZONE, ALL_TIMEZONES } from '@/consts/timezone.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { KnownErrorCode, SPECIFIED_API_NOT_FOUND_ERRORS, PARAMETERIZED_ERRORS } from '@/consts/api.ts';

import type { LatestExchangeRateResponse, LocalizedLatestExchangeRate } from '@/models/exchange_rate.ts';

import {
    isObject,
    isString,
    isNumber,
    isBoolean
} from '@/lib/common.ts';

import {
    isPM,
    formatUnixTime,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffset,
    getBrowserTimezoneOffsetMinutes,
    getTimeDifferenceHoursAndMinutes,
    getDateTimeFormatType,
    getRecentMonthDateRanges
} from '@/lib/datetime.ts';

import {
    appendDigitGroupingSymbol,
    parseAmount,
    formatAmount,
    formatExchangeRateAmount,
    getAdaptiveDisplayAmountRate
} from '@/lib/numeral.ts';

import {
    getCurrencyFraction,
    appendCurrencySymbol,
    getAmountPrependAndAppendCurrencySymbol
} from '@/lib/currency.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

export interface LocalizedErrorParameter {
    key: string;
    localized: boolean;
    value: string;
}

export interface LocalizedError {
    message: string;
    parameters?: LocalizedErrorParameter[];
}

export function getI18nOptions(): object {
    return {
        legacy: true,
        locale: DEFAULT_LANGUAGE,
        fallbackLocale: DEFAULT_LANGUAGE,
        formatFallbackMessages: true,
        messages: (function () {
            const messages: Record<string, object> = {};

            for (const languageKey in allLanguages) {
                if (!Object.prototype.hasOwnProperty.call(allLanguages, languageKey)) {
                    continue;
                }

                const languageInfo = allLanguages[languageKey];
                messages[languageKey] = languageInfo.content;
            }

            return messages;
        })()
    };
}

export function useI18n() {
    const { t, locale } = useVueI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();

    // private functions
    function getLanguageInfo(languageKey: string): LanguageInfo | undefined {
        return allLanguages[languageKey];
    }

    function getDefaultLanguage(): string {
        if (!window || !window.navigator) {
            return DEFAULT_LANGUAGE;
        }

        let browserLanguage = window.navigator.browserLanguage || window.navigator.language;

        if (!browserLanguage) {
            return DEFAULT_LANGUAGE;
        }

        if (!allLanguages[browserLanguage]) {
            const languageKey = getLanguageKeyFromLanguageAlias(browserLanguage);

            if (languageKey) {
                browserLanguage = languageKey;
            }
        }

        if (!allLanguages[browserLanguage] && browserLanguage.split('-').length > 1) { // maybe language-script-region
            const languageTagParts = browserLanguage.split('-');
            browserLanguage = languageTagParts[0] + '-' + languageTagParts[1];

            if (!allLanguages[browserLanguage]) {
                const languageKey = getLanguageKeyFromLanguageAlias(browserLanguage);

                if (languageKey) {
                    browserLanguage = languageKey;
                }
            }

            if (!allLanguages[browserLanguage]) {
                browserLanguage = languageTagParts[0];
                const languageKey = getLanguageKeyFromLanguageAlias(browserLanguage);

                if (languageKey) {
                    browserLanguage = languageKey;
                }
            }
        }

        if (!allLanguages[browserLanguage]) {
            return DEFAULT_LANGUAGE;
        }

        return browserLanguage;
    }

    function getLanguageKeyFromLanguageAlias(alias: string): string | null {
        for (const languageKey in allLanguages) {
            if (!Object.prototype.hasOwnProperty.call(allLanguages, languageKey)) {
                continue;
            }

            if (languageKey.toLowerCase() === alias.toLowerCase()) {
                return languageKey;
            }

            const langInfo = allLanguages[languageKey];
            const aliases = langInfo.aliases;

            if (!aliases || aliases.length < 1) {
                continue;
            }

            for (let i = 0; i < aliases.length; i++) {
                if (aliases[i].toLowerCase() === alias.toLowerCase()) {
                    return languageKey;
                }
            }
        }

        return null;
    }

    function getLocalizedError(error: ErrorResponse): LocalizedError {
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

    function getLocalizedErrorParameters(parameters?: LocalizedErrorParameter[]): Record<string, string> {
        const localizedParameters: Record<string, string> = {};

        if (parameters) {
            for (let i = 0; i < parameters.length; i++) {
                const parameter = parameters[i];

                if (parameter.localized) {
                    localizedParameters[parameter.key] = t(`parameter.${parameter.value}`);
                } else {
                    localizedParameters[parameter.key] = parameter.value;
                }
            }
        }

        return localizedParameters;
    }

    function getAllCurrencyDisplayTypes(): TypeAndDisplayName[] {
        const defaultCurrencyDisplayTypeName = t('default.currencyDisplayType');
        let defaultCurrencyDisplayType = CurrencyDisplayType.parse(defaultCurrencyDisplayTypeName);

        if (!defaultCurrencyDisplayType) {
            defaultCurrencyDisplayType = CurrencyDisplayType.Default;
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;

        const ret = [];
        const defaultSampleValue = getFormattedAmountWithCurrency(12345, defaultCurrency, false, defaultCurrencyDisplayType);

        ret.push({
            type: CurrencyDisplayType.LanguageDefaultType,
            displayName: `${t('Language Default')} (${defaultSampleValue})`
        });

        const allCurrencyDisplayTypes = CurrencyDisplayType.values();

        for (let i = 0; i < allCurrencyDisplayTypes.length; i++) {
            const type = allCurrencyDisplayTypes[i];
            const sampleValue = getFormattedAmountWithCurrency(12345, defaultCurrency, false, type);
            const displayName = `${t(type.name)} (${sampleValue})`

            ret.push({
                type: type.type,
                displayName: displayName
            });
        }

        return ret;
    }

    function getLocalizedDisplayNameAndType(typeAndNames: TypeAndName[]): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];

        for (let i = 0; i < typeAndNames.length; i++) {
            const nameAndType = typeAndNames[i];

            ret.push({
                type: nameAndType.type,
                displayName: t(nameAndType.name)
            });
        }

        return ret;
    }

    function getLocalizedNumeralSeparatorFormats<T extends NumeralSymbolType>(allSeparatorArray: T[], localeDefaultType: T | undefined, systemDefaultType: T, languageDefaultValue: number): LocalizedNumeralSymbolType[] {
        let defaultSeparatorType: T | undefined = localeDefaultType;

        if (!defaultSeparatorType) {
            defaultSeparatorType = systemDefaultType;
        }

        const ret: LocalizedNumeralSymbolType[] = [];

        ret.push({
            type: languageDefaultValue,
            symbol: defaultSeparatorType.symbol,
            displayName: `${t('Language Default')} (${defaultSeparatorType.symbol})`
        });

        for (let i = 0; i < allSeparatorArray.length; i++) {
            const type = allSeparatorArray[i];

            ret.push({
                type: type.type,
                symbol: type.symbol,
                displayName: `${t('numeral.' + type.name)} (${type.symbol})`
            });
        }

        return ret;
    }

    function getAllMonthNames(type: string): string[] {
        const ret = [];
        const allMonths = Month.values();

        for (let i = 0; i < allMonths.length; i++) {
            const month = allMonths[i];
            ret.push(t(`datetime.${month.name}.${type}`));
        }

        return ret;
    }

    function getAllWeekdayNames(type: string): string[] {
        const ret = [];
        const allWeekDays = WeekDay.values();

        for (let i = 0; i < allWeekDays.length; i++) {
            const weekDay = allWeekDays[i];
            ret.push(t(`datetime.${weekDay.name}.${type}`));
        }

        return ret;
    }

    function getLocalizedDateTimeType<T extends DateFormat | TimeFormat>(allFormatMap: Record<string, T>, allFormatArray: T[], formatTypeValue: number, languageDefaultTypeNameKey: string, systemDefaultFormatType: T): T {
        return getDateTimeFormatType(allFormatMap, allFormatArray, formatTypeValue, t(`default.${languageDefaultTypeNameKey}`), systemDefaultFormatType);
    }

    function getLocalizedDateTimeFormat<T extends DateFormat | TimeFormat>(type: string, allFormatMap: Record<string, T>, allFormatArray: T[], formatTypeValue: number, languageDefaultTypeNameKey: string, systemDefaultFormatType: T): string {
        const formatType = getLocalizedDateTimeType(allFormatMap, allFormatArray, formatTypeValue, languageDefaultTypeNameKey, systemDefaultFormatType);
        return t(`format.${type}.${formatType.key}`);
    }

    function getLocalizedLongDateFormat(): string {
        return getLocalizedDateTimeFormat<LongDateFormat>('longDate', LongDateFormat.all(), LongDateFormat.values(), userStore.currentUserLongDateFormat, 'longDateFormat', LongDateFormat.Default);
    }

    function getLocalizedShortDateFormat(): string {
        return getLocalizedDateTimeFormat<ShortDateFormat>('shortDate', ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default);
    }

    function getLocalizedLongYearFormat(): string {
        return getLocalizedDateTimeFormat<LongDateFormat>('longYear', LongDateFormat.all(), LongDateFormat.values(), userStore.currentUserLongDateFormat, 'longDateFormat', LongDateFormat.Default);
    }

    function getLocalizedShortYearFormat(): string {
        return getLocalizedDateTimeFormat<ShortDateFormat>('shortYear', ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default);
    }

    function getLocalizedLongYearMonthFormat(): string {
        return getLocalizedDateTimeFormat<LongDateFormat>('longYearMonth', LongDateFormat.all(), LongDateFormat.values(), userStore.currentUserLongDateFormat, 'longDateFormat', LongDateFormat.Default);
    }

    function getLocalizedShortYearMonthFormat(): string {
        return getLocalizedDateTimeFormat<ShortDateFormat>('shortYearMonth', ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default);
    }

    function getLocalizedLongMonthDayFormat(): string {
        return getLocalizedDateTimeFormat<LongDateFormat>('longMonthDay', LongDateFormat.all(), LongDateFormat.values(), userStore.currentUserLongDateFormat, 'longDateFormat', LongDateFormat.Default);
    }

    function getLocalizedShortMonthDayFormat(): string {
        return getLocalizedDateTimeFormat<ShortDateFormat>('shortMonthDay', ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default);
    }

    function getLocalizedLongTimeFormat(): string {
        return getLocalizedDateTimeFormat<LongTimeFormat>('longTime', LongTimeFormat.all(), LongTimeFormat.values(), userStore.currentUserLongTimeFormat, 'longTimeFormat', LongTimeFormat.Default);
    }

    function getLocalizedShortTimeFormat(): string {
        return getLocalizedDateTimeFormat<ShortTimeFormat>('shortTime', ShortTimeFormat.all(), ShortTimeFormat.values(), userStore.currentUserShortTimeFormat, 'shortTimeFormat', ShortTimeFormat.Default);
    }

    function getNumberFormatOptions(currencyCode?: string): NumberFormatOptions {
        return {
            decimalSeparator: getCurrentDecimalSeparator(),
            decimalNumberCount: getCurrencyFraction(currencyCode),
            digitGroupingSymbol: getCurrentDigitGroupingSymbol(),
            digitGrouping: getCurrentDigitGroupingType(),
        };
    }

    function getCurrencyUnitName(currencyCode: string, isPlural: boolean): string {
        const currencyInfo = ALL_CURRENCIES[currencyCode];

        if (currencyInfo && currencyInfo.unit) {
            if (isPlural) {
                return t(`currency.unit.${currencyInfo.unit}.plural`);
            } else {
                return t(`currency.unit.${currencyInfo.unit}.normal`);
            }
        }

        return '';
    }

    function getCurrentCurrencyDisplayType(): CurrencyDisplayType {
        let currencyDisplayType = CurrencyDisplayType.valueOf(userStore.currentUserCurrencyDisplayType);

        if (!currencyDisplayType) {
            const defaultCurrencyDisplayTypeName = t('default.currencyDisplayType');
            currencyDisplayType = CurrencyDisplayType.parse(defaultCurrencyDisplayTypeName);
        }

        if (!currencyDisplayType) {
            currencyDisplayType = CurrencyDisplayType.Default;
        }

        return currencyDisplayType;
    }

    // public functions
    function translateIf(text: string, isTranslate: boolean): string {
        if (isTranslate) {
            return t(text);
        }

        return text;
    }

    function translateError(message: string | { error: ErrorResponse }): string {
        let finalMessage = '';
        let parameters = {};

        if (isObject(message) && isObject((message as { error: ErrorResponse }).error)) {
            const localizedError = getLocalizedError((message as { error: ErrorResponse }).error);
            finalMessage = localizedError.message;
            parameters = getLocalizedErrorParameters(localizedError.parameters);
        } else if (isString(message)) {
            finalMessage = message as string;
        } else {
            return '';
        }

        return t(finalMessage, parameters);
    }

    function joinMultiText(textArray: string[]): string {
        if (!textArray || !textArray.length) {
            return '';
        }

        const separator = t('format.misc.multiTextJoinSeparator');

        return textArray.join(separator);
    }

    function getCurrentLanguageTag(): string {
        return locale.value;
    }

    function getCurrentLanguageInfo(): LanguageInfo {
        const languageInfo = getLanguageInfo(locale.value);

        if (languageInfo) {
            return languageInfo;
        }

        return getLanguageInfo(getDefaultLanguage()) as LanguageInfo;
    }

    function getCurrentLanguageDisplayName() {
        const currentLanguageInfo = getCurrentLanguageInfo();
        return currentLanguageInfo.displayName;
    }

    function getDefaultCurrency(): string {
        return t('default.currency');
    }

    function getDefaultFirstDayOfWeek(): string {
        return t('default.firstDayOfWeek');
    }

    function getAllEnableDisableOptions(): LocalizedSwitchOption[] {
        return [{
            value: true,
            displayName: t('Enable')
        },{
            value: false,
            displayName: t('Disable')
        }];
    }

    function getAllMeridiemIndicators(): LocalizedMeridiemIndicator {
        const allMeridiemIndicators = MeridiemIndicator.values();
        const meridiemIndicatorNames = [];
        const localizedMeridiemIndicatorNames = [];

        for (let i = 0; i < allMeridiemIndicators.length; i++) {
            const indicator = allMeridiemIndicators[i];

            meridiemIndicatorNames.push(indicator.name);
            localizedMeridiemIndicatorNames.push(t(`datetime.${indicator.name}.content`));
        }

        return {
            values: meridiemIndicatorNames,
            displayValues: localizedMeridiemIndicatorNames
        };
    }

    function getAllLongMonthNames(): string[] {
        return getAllMonthNames('long');
    }

    function getAllShortMonthNames(): string[] {
        return getAllMonthNames('short');
    }

    function getAllLongWeekdayNames(): string[] {
        return getAllWeekdayNames('long');
    }

    function getAllShortWeekdayNames(): string[] {
        return getAllWeekdayNames('short');
    }

    function getAllMinWeekdayNames(): string[] {
        return getAllWeekdayNames('min');
    }

    function getAllWeekDays(firstDayOfWeek: number): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];
        const allWeekDays = WeekDay.values();

        if (!isNumber(firstDayOfWeek)) {
            firstDayOfWeek = WeekDay.DefaultFirstDay.type;
        }

        for (let i = firstDayOfWeek; i < allWeekDays.length; i++) {
            const weekDay = allWeekDays[i];

            ret.push({
                type: weekDay.type,
                displayName: t(`datetime.${weekDay.name}.long`)
            });
        }

        for (let i = 0; i < firstDayOfWeek; i++) {
            const weekDay = allWeekDays[i];

            ret.push({
                type: weekDay.type,
                displayName: t(`datetime.${weekDay.name}.long`)
            });
        }

        return ret;
    }

    function getAllDateRanges(scene: DateRangeScene, includeCustom?: boolean, includeBillingCycle?: boolean): LocalizedDateRange[] {
        const ret: LocalizedDateRange[] = [];
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
                        displayName: t(dateRange.name),
                        isBillingCycle: dateRange.isBillingCycle
                    });
                }

                continue;
            }

            if (includeCustom || dateRange.type !== DateRange.Custom.type) {
                ret.push({
                    type: dateRange.type,
                    displayName: t(dateRange.name)
                });
            }
        }

        return ret;
    }

    function getAllRecentMonthDateRanges(includeAll: boolean, includeCustom: boolean): LocalizedRecentMonthDateRange[] {
        const allRecentMonthDateRanges: LocalizedRecentMonthDateRange[] = [];
        const recentDateRanges = getRecentMonthDateRanges(12);

        if (includeAll) {
            allRecentMonthDateRanges.push({
                dateType: DateRange.All.type,
                minTime: 0,
                maxTime: 0,
                displayName: t('All')
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
                displayName: formatUnixTime(recentDateRange.minTime, getLocalizedLongYearMonthFormat())
            });
        }

        if (includeCustom) {
            allRecentMonthDateRanges.push({
                dateType: DateRange.Custom.type,
                minTime: 0,
                maxTime: 0,
                displayName: t('Custom Date')
            });
        }

        return allRecentMonthDateRanges;
    }

    function getAllTimezones(includeSystemDefault?: boolean): LocalizedTimezoneInfo[] {
        const defaultTimezoneOffset = getBrowserTimezoneOffset();
        const defaultTimezoneOffsetMinutes = getBrowserTimezoneOffsetMinutes();
        const allTimezoneInfos: LocalizedTimezoneInfo[] = [];

        for (let i = 0; i < ALL_TIMEZONES.length; i++) {
            const utcOffset = (ALL_TIMEZONES[i].timezoneName !== UTC_TIMEZONE.timezoneName ? getTimezoneOffset(ALL_TIMEZONES[i].timezoneName) : '');
            const displayName = t(`timezone.${ALL_TIMEZONES[i].displayName}`);

            allTimezoneInfos.push({
                name: ALL_TIMEZONES[i].timezoneName,
                utcOffset: utcOffset,
                utcOffsetMinutes: getTimezoneOffsetMinutes(ALL_TIMEZONES[i].timezoneName),
                displayName: displayName,
                displayNameWithUtcOffset: `(UTC${utcOffset}) ${displayName}`
            });
        }

        if (includeSystemDefault) {
            const defaultDisplayName = t('System Default');

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

    function getAllTimezoneTypesUsedForStatistics(currentTimezone?: string): TypeAndDisplayName[] {
        const currentTimezoneOffset = getTimezoneOffset(currentTimezone);

        return [
            {
                type: TimezoneTypeForStatistics.ApplicationTimezone.type,
                displayName: t(TimezoneTypeForStatistics.ApplicationTimezone.name) + ` (UTC${currentTimezoneOffset})`
            },
            {
                type: TimezoneTypeForStatistics.TransactionTimezone.type,
                displayName: t(TimezoneTypeForStatistics.TransactionTimezone.name)
            }
        ];
    }

    function getAllDigitGroupingTypes(): LocalizedDigitGroupingType[] {
        const defaultDigitGroupingTypeName = t('default.digitGrouping');
        let defaultDigitGroupingType = DigitGroupingType.parse(defaultDigitGroupingTypeName);

        if (!defaultDigitGroupingType) {
            defaultDigitGroupingType = DigitGroupingType.Default;
        }

        const ret: LocalizedDigitGroupingType[] = [];

        ret.push({
            type: DigitGroupingType.LanguageDefaultType,
            enabled: defaultDigitGroupingType.enabled,
            displayName: `${t('Language Default')} (${t('numeral.' + defaultDigitGroupingType.name)})`
        });

        const allDigitGroupingTypes = DigitGroupingType.values();

        for (let i = 0; i < allDigitGroupingTypes.length; i++) {
            const type = allDigitGroupingTypes[i];

            ret.push({
                type: type.type,
                enabled: type.enabled,
                displayName: t('numeral.' + type.name)
            });
        }

        return ret;
    }

    function getAllAccountCategories(): LocalizedAccountCategory[] {
        const ret: LocalizedAccountCategory[] = [];
        const allCategories = AccountCategory.values();

        for (let i = 0; i < allCategories.length; i++) {
            const accountCategory = allCategories[i];

            ret.push({
                type: accountCategory.type,
                displayName: t(accountCategory.name),
                defaultAccountIconId: accountCategory.defaultAccountIconId
            });
        }

        return ret;
    }

    function getAllDisplayExchangeRates(exchangeRatesData?: LatestExchangeRateResponse): LocalizedLatestExchangeRate[] {
        const availableExchangeRates: LocalizedLatestExchangeRate[] = [];

        if (!exchangeRatesData || !exchangeRatesData.exchangeRates) {
            return availableExchangeRates;
        }

        for (let i = 0; i < exchangeRatesData.exchangeRates.length; i++) {
            const exchangeRate = exchangeRatesData.exchangeRates[i];

            availableExchangeRates.push({
                currencyCode: exchangeRate.currency,
                currencyDisplayName: getCurrencyName(exchangeRate.currency),
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

    function getMonthShortName(monthName: string): string {
        return t(`datetime.${monthName}.short`);
    }

    function getMonthLongName(monthName: string): string {
        return t(`datetime.${monthName}.long`);
    }

    function getMonthdayOrdinal(monthDay: number): string {
        return t(`datetime.monthDayOrdinal.${monthDay}`);
    }

    function getMonthdayShortName(monthDay: number): string {
        return t('format.misc.monthDay', {
            ordinal: getMonthdayOrdinal(monthDay)
        });
    }

    function getWeekdayShortName(weekDayName: string): string {
        return t(`datetime.${weekDayName}.short`);
    }

    function getWeekdayLongName(weekDayName: string): string {
        return t(`datetime.${weekDayName}.long`);
    }

    function getMultiMonthdayShortNames(monthDays: number[]): string {
        if (!monthDays) {
            return '';
        }

        if (monthDays.length === 1) {
            return t('format.misc.monthDay', {
                ordinal: getMonthdayOrdinal(monthDays[0])
            });
        } else {
            return t('format.misc.monthDays', {
                multiMonthDays: joinMultiText(monthDays.map(monthDay =>
                    t('format.misc.eachMonthDayInMonthDays', {
                        ordinal: getMonthdayOrdinal(monthDay)
                    })))
            });
        }
    }

    function getMultiWeekdayLongNames(weekdayTypes: number[], firstDayOfWeek?: number): string {
        const weekdayTypesMap: Record<number, boolean> = {};

        if (!isNumber(firstDayOfWeek)) {
            firstDayOfWeek = WeekDay.DefaultFirstDay.type;
        }

        for (let i = 0; i < weekdayTypes.length; i++) {
            weekdayTypesMap[weekdayTypes[i]] = true;
        }

        const allWeekDays = getAllWeekDays(firstDayOfWeek as number);
        const finalWeekdayNames = [];

        for (let i = 0; i < allWeekDays.length; i++) {
            const weekDay = allWeekDays[i];

            if (weekdayTypesMap[weekDay.type]) {
                finalWeekdayNames.push(weekDay.displayName);
            }
        }

        return joinMultiText(finalWeekdayNames);
    }

    function getCurrentDecimalSeparator(): string {
        let decimalSeparatorType = DecimalSeparator.valueOf(userStore.currentUserDecimalSeparator);

        if (!decimalSeparatorType) {
            const defaultDecimalSeparatorTypeName = t('default.decimalSeparator');
            decimalSeparatorType = DecimalSeparator.parse(defaultDecimalSeparatorTypeName);

            if (!decimalSeparatorType) {
                decimalSeparatorType = DecimalSeparator.Default;
            }
        }

        return decimalSeparatorType.symbol;
    }

    function getCurrentDigitGroupingSymbol(): string {
        let digitGroupingSymbolType = DigitGroupingSymbol.valueOf(userStore.currentUserDigitGroupingSymbol);

        if (!digitGroupingSymbolType) {
            const defaultDigitGroupingSymbolTypeName = t('default.digitGroupingSymbol');
            digitGroupingSymbolType = DigitGroupingSymbol.parse(defaultDigitGroupingSymbolTypeName);

            if (!digitGroupingSymbolType) {
                digitGroupingSymbolType = DigitGroupingSymbol.Default;
            }
        }

        return digitGroupingSymbolType.symbol;
    }

    function getCurrentDigitGroupingType(): number {
        let digitGroupingType = DigitGroupingType.valueOf(userStore.currentUserDigitGrouping);

        if (!digitGroupingType) {
            const defaultDigitGroupingTypeName = t('default.digitGrouping');
            digitGroupingType = DigitGroupingType.parse(defaultDigitGroupingTypeName);

            if (!digitGroupingType) {
                digitGroupingType = DigitGroupingType.Default;
            }
        }

        return digitGroupingType.type;
    }

    function getCurrencyName(currencyCode: string): string {
        return t(`currency.name.${currencyCode}`);
    }

    function isLongDateMonthAfterYear() {
        return getLocalizedDateTimeType(LongDateFormat.all(), LongDateFormat.values(), userStore.currentUserLongDateFormat, 'longDateFormat', LongDateFormat.Default).isMonthAfterYear;
    }

    function isShortDateMonthAfterYear() {
        return getLocalizedDateTimeType(ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default).isMonthAfterYear;
    }

    function isLongTime24HourFormat() {
        return getLocalizedDateTimeType(LongTimeFormat.all(), LongTimeFormat.values(), userStore.currentUserLongTimeFormat, 'longTimeFormat', LongTimeFormat.Default).is24HourFormat;
    }

    function isLongTimeMeridiemIndicatorFirst() {
        return getLocalizedDateTimeType(LongTimeFormat.all(), LongTimeFormat.values(), userStore.currentUserLongTimeFormat, 'longTimeFormat', LongTimeFormat.Default).isMeridiemIndicatorFirst;
    }

    function isShortTime24HourFormat() {
        return getLocalizedDateTimeType(ShortTimeFormat.all(), ShortTimeFormat.values(), userStore.currentUserShortTimeFormat, 'shortTimeFormat', ShortTimeFormat.Default).is24HourFormat;
    }

    function isShortTimeMeridiemIndicatorFirst() {
        return getLocalizedDateTimeType(ShortTimeFormat.all(), ShortTimeFormat.values(), userStore.currentUserShortTimeFormat, 'shortTimeFormat', ShortTimeFormat.Default).isMeridiemIndicatorFirst;
    }

    function formatYearQuarter(year: number, quarter: number): string {
        if (1 <= quarter && quarter <= 4) {
            return t('format.yearQuarter.q' + quarter, {
                year: year,
                quarter: quarter
            });
        } else {
            return '';
        }
    }

    function getTimezoneDifferenceDisplayText(utcOffset: number): string {
        const defaultTimezoneOffset = getTimezoneOffsetMinutes();
        const offsetTime = getTimeDifferenceHoursAndMinutes(utcOffset - defaultTimezoneOffset);

        if (utcOffset > defaultTimezoneOffset) {
            if (offsetTime.offsetMinutes) {
                return t('format.misc.hoursMinutesAheadOfDefaultTimezone', {
                    hours: offsetTime.offsetHours,
                    minutes: offsetTime.offsetMinutes
                });
            } else {
                return t('format.misc.hoursAheadOfDefaultTimezone', {
                    hours: offsetTime.offsetHours
                });
            }
        } else if (utcOffset < defaultTimezoneOffset) {
            if (offsetTime.offsetMinutes) {
                return t('format.misc.hoursMinutesBehindDefaultTimezone', {
                    hours: offsetTime.offsetHours,
                    minutes: offsetTime.offsetMinutes
                });
            } else {
                return t('format.misc.hoursBehindDefaultTimezone', {
                    hours: offsetTime.offsetHours
                });
            }
        } else {
            return t('Same time as default timezone');
        }
    }

    function getNumberWithDigitGroupingSymbol(value: number | string): string {
        const numberFormatOptions = getNumberFormatOptions();
        return appendDigitGroupingSymbol(value, numberFormatOptions);
    }

    function getParsedAmountNumber(value: string): number {
        const numberFormatOptions = getNumberFormatOptions();
        return parseAmount(value, numberFormatOptions);
    }

    function getFormattedAmount(value: number | string, currencyCode?: string): string {
        const numberFormatOptions = getNumberFormatOptions(currencyCode);
        return formatAmount(value, numberFormatOptions);
    }

    function getFormattedAmountWithCurrency(value: number | string, currencyCode?: string, notConvertValue?: boolean, currencyDisplayType?: CurrencyDisplayType): string | null {
        if (!isNumber(value) && !isString(value)) {
            return null;
        }

        if (isNumber(value)) {
            value = value.toString();
        }

        let textualValue = value as string;
        const isPlural: boolean = textualValue !== '100' && textualValue !== '-100';

        if (!notConvertValue) {
            const numberFormatOptions = getNumberFormatOptions();
            const hasIncompleteFlag = isString(textualValue) && textualValue.charAt(textualValue.length - 1) === '+';

            if (hasIncompleteFlag) {
                textualValue = textualValue.substring(0, textualValue.length - 1);
            }

            textualValue = formatAmount(textualValue, numberFormatOptions);

            if (hasIncompleteFlag) {
                textualValue = textualValue + '+';
            }
        }

        if (!isBoolean(currencyCode) && !currencyCode) {
            currencyCode = userStore.currentUserDefaultCurrency;
        } else if (isBoolean(currencyCode) && !currencyCode) {
            currencyCode = '';
        }

        if (!currencyCode) {
            return textualValue;
        }

        if (!currencyDisplayType) {
            currencyDisplayType = getCurrentCurrencyDisplayType();
        }

        const currencyUnit = getCurrencyUnitName(currencyCode, isPlural);
        const currencyName = getCurrencyName(currencyCode);
        return appendCurrencySymbol(textualValue, currencyDisplayType, currencyCode, currencyUnit, currencyName, isPlural);
    }

    function getFormattedExchangeRateAmount(value: number | string): string {
        const numberFormatOptions = getNumberFormatOptions();
        return formatExchangeRateAmount(value, numberFormatOptions);
    }

    function getAdaptiveAmountRate(amount1: number, amount2: number, fromExchangeRate: { rate: string }, toExchangeRate: { rate: string }): string | null {
        const numberFormatOptions = getNumberFormatOptions();
        return getAdaptiveDisplayAmountRate(amount1, amount2, fromExchangeRate, toExchangeRate, numberFormatOptions);
    }

    function getAmountPrependAndAppendText(currencyCode: string, isPlural: boolean): CurrencyPrependAndAppendText | null {
        const currencyDisplayType = getCurrentCurrencyDisplayType();
        const currencyUnit = getCurrencyUnitName(currencyCode, isPlural);
        const currencyName = getCurrencyName(currencyCode);
        return getAmountPrependAndAppendCurrencySymbol(currencyDisplayType, currencyCode, currencyUnit, currencyName, isPlural);
    }

    function setLanguage(languageKey: string | null, force?: boolean): LocaleDefaultSettings | null {
        if (!languageKey) {
            languageKey = getDefaultLanguage();
            logger.info(`No specified language, use browser default language ${languageKey}`);
        }

        if (!getLanguageInfo(languageKey)) {
            languageKey = getDefaultLanguage();
            logger.warn(`Not found language ${languageKey}, use browser default language ${languageKey}`);
        }

        if (!force && locale.value === languageKey) {
            logger.info(`Current locale is already ${languageKey}`);
            return null;
        }

        logger.info(`Apply current language to ${languageKey}`);

        locale.value = languageKey;
        moment.updateLocale(languageKey, {
            months : getAllLongMonthNames(),
            monthsShort : getAllShortMonthNames(),
            weekdays : getAllLongWeekdayNames(),
            weekdaysShort : getAllShortWeekdayNames(),
            weekdaysMin : getAllMinWeekdayNames(),
            meridiem: function (hours) {
                if (isPM(hours)) {
                    return t(`datetime.${MeridiemIndicator.PM.name}.content`);
                } else {
                    return t(`datetime.${MeridiemIndicator.AM.name}.content`);
                }
            }
        });

        services.setLocale(languageKey);
        document.querySelector('html')?.setAttribute('lang', languageKey);

        const defaultCurrency = getDefaultCurrency();
        const defaultFirstDayOfWeekName = getDefaultFirstDayOfWeek();
        let defaultFirstDayOfWeek = WeekDay.DefaultFirstDay.type;

        if (WeekDay.parse(defaultFirstDayOfWeekName)) {
            defaultFirstDayOfWeek = (WeekDay.parse(defaultFirstDayOfWeekName) as WeekDay).type;
        }

        return {
            currency: defaultCurrency,
            firstDayOfWeek: defaultFirstDayOfWeek
        };
    }

    function setTimeZone(timezone: string): void {
        if (timezone) {
            moment.tz.setDefault(timezone);
        } else {
            moment.tz.setDefault();
        }
    }

    function initLocale(lastUserLanguage: string, timezone: string): LocaleDefaultSettings | null {
        let localeDefaultSettings: LocaleDefaultSettings | null = null;

        if (lastUserLanguage && getLanguageInfo(lastUserLanguage)) {
            logger.info(`Last user language is ${lastUserLanguage}`);
            localeDefaultSettings = setLanguage(lastUserLanguage, true);
        } else {
            localeDefaultSettings = setLanguage(null, true);
        }

        if (timezone) {
            logger.info(`Current timezone is ${timezone}`);
            setTimeZone(timezone);
        } else {
            logger.info(`No timezone is set, use browser default ${getTimezoneOffset()} (maybe ${moment.tz.guess(true)})`);
        }

        return localeDefaultSettings;
    }

    return {
        // common functions
        tt: t,
        ti: translateIf,
        te: translateError,
        joinMultiText,
        // get current language info
        getCurrentLanguageTag,
        getCurrentLanguageInfo,
        getCurrentLanguageDisplayName,
        // get localization default type
        getDefaultCurrency,
        getDefaultFirstDayOfWeek,
        // get all localized info of specified type
        getAllEnableDisableOptions,
        getAllMeridiemIndicators,
        getAllLongMonthNames,
        getAllShortMonthNames,
        getAllLongWeekdayNames,
        getAllShortWeekdayNames,
        getAllMinWeekdayNames,
        getAllWeekDays,
        getAllDateRanges,
        getAllRecentMonthDateRanges,
        getAllTimezones,
        getAllTimezoneTypesUsedForStatistics,
        getAllDecimalSeparators: () => getLocalizedNumeralSeparatorFormats(DecimalSeparator.values(), DecimalSeparator.parse(t('default.decimalSeparator')), DecimalSeparator.Default, DecimalSeparator.LanguageDefaultType),
        getAllDigitGroupingSymbols: () => getLocalizedNumeralSeparatorFormats(DigitGroupingSymbol.values(), DigitGroupingSymbol.parse(t('default.digitGroupingSymbol')), DigitGroupingSymbol.Default, DigitGroupingSymbol.LanguageDefaultType),
        getAllDigitGroupingTypes,
        getAllCurrencyDisplayTypes,
        getAllCurrencySortingTypes: () => getLocalizedDisplayNameAndType(CurrencySortingType.values()),
        getAllAccountCategories,
        getAllAccountTypes: () => getLocalizedDisplayNameAndType(AccountType.values()),
        getAllCategoricalChartTypes: () => getLocalizedDisplayNameAndType(CategoricalChartType.values()),
        getAllTrendChartTypes: () => getLocalizedDisplayNameAndType(TrendChartType.values()),
        getAllStatisticsChartDataTypes: (analysisType: StatisticsAnalysisType) => getLocalizedDisplayNameAndType(ChartDataType.values(analysisType)),
        getAllStatisticsSortingTypes: () => getLocalizedDisplayNameAndType(ChartSortingType.values()),
        getAllStatisticsDateAggregationTypes: () => getLocalizedDisplayNameAndType(ChartDateAggregationType.values()),
        getAllTransactionEditScopeTypes: () => getLocalizedDisplayNameAndType(TransactionEditScopeType.values()),
        getAllTransactionTagFilterTypes: () => getLocalizedDisplayNameAndType(TransactionTagFilterType.values()),
        getAllTransactionScheduledFrequencyTypes: () => getLocalizedDisplayNameAndType(ScheduledTemplateFrequencyType.values()),
        getAllDisplayExchangeRates,
        // get localized info
        getMonthShortName,
        getMonthLongName,
        getMonthdayOrdinal,
        getMonthdayShortName,
        getWeekdayShortName,
        getWeekdayLongName,
        getMultiMonthdayShortNames,
        getMultiWeekdayLongNames,
        getCurrentDecimalSeparator,
        getCurrentDigitGroupingSymbol,
        getCurrentDigitGroupingType,
        getCurrencyName,
        isLongDateMonthAfterYear,
        isShortDateMonthAfterYear,
        isLongTime24HourFormat,
        isLongTimeMeridiemIndicatorFirst,
        isShortTime24HourFormat,
        isShortTimeMeridiemIndicatorFirst,
        // format functions
        formatUnixTimeToLongDateTime: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongDateFormat() + ' ' + getLocalizedLongTimeFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortDateTime: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortDateFormat() + ' ' + getLocalizedShortTimeFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongDate: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongDateFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortDate: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortDateFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongYear: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongYearFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortYear: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortYearFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongYearMonth: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongYearMonthFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortYearMonth: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortYearMonthFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongMonthDay: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongMonthDayFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortMonthDay: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortMonthDayFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongTime: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongTimeFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortTime: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortTimeFormat(), utcOffset, currentUtcOffset),
        formatYearQuarter,
        getTimezoneDifferenceDisplayText,
        appendDigitGroupingSymbol: getNumberWithDigitGroupingSymbol,
        parseAmount: getParsedAmountNumber,
        formatAmount: getFormattedAmount,
        formatAmountWithCurrency: getFormattedAmountWithCurrency,
        formatExchangeRateAmount: getFormattedExchangeRateAmount,
        getAdaptiveAmountRate,
        getAmountPrependAndAppendText,
        // localization setting functions
        setLanguage,
        setTimeZone,
        initLocale
    };
}
