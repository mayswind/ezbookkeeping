import { useI18n as useVueI18n } from 'vue-i18n';
import moment from 'moment-timezone';

import type { PartialRecord, TypeAndName, TypeAndDisplayName, LocalizedSwitchOption } from '@/core/base.ts';

import { type LanguageInfo, type LanguageOption, ALL_LANGUAGES, DEFAULT_LANGUAGE } from '@/locales/index.ts';

import {
    type DateFormat,
    type TimeFormat,
    type LocalizedMeridiemIndicator,
    type LocalizedDateTimeFormat,
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
    DateRangeScene,
    LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE
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
    type LocalizedCurrencyInfo,
    type CurrencyPrependAndAppendText,
    CurrencyDisplayType,
    CurrencySortingType
} from '@/core/currency.ts';

import {
    FiscalYearStart
} from '@/core/fiscalyear.ts';

import {
    PresetAmountColor
} from '@/core/color.ts';

import {
    type LocalizedAccountCategory,
    AccountType,
    AccountCategory
} from '@/core/account.ts';

import {
    type PresetCategory,
    type LocalizedPresetCategory,
    type LocalizedPresetSubCategory,
    CategoryType,
    ALL_CATEGORY_TYPES
} from '@/core/category.ts';

import {
    TransactionEditScopeType,
    TransactionTagFilterType,
    ImportTransactionColumnType
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

import {
    type LocalizedImportFileType,
    type LocalizedImportFileTypeSubType,
    type LocalizedImportFileTypeSupportedEncodings,
    type LocalizedImportFileDocument,
} from '@/core/file.ts';

import type { LocaleDefaultSettings } from '@/core/setting.ts';
import type { ErrorResponse } from '@/core/api.ts';

import { UTC_TIMEZONE, ALL_TIMEZONES } from '@/consts/timezone.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { DEFAULT_EXPENSE_CATEGORIES, DEFAULT_INCOME_CATEGORIES, DEFAULT_TRANSFER_CATEGORIES } from '@/consts/category.ts';
import { KnownErrorCode, SPECIFIED_API_NOT_FOUND_ERRORS, PARAMETERIZED_ERRORS } from '@/consts/api.ts';
import { DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE, SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE, SUPPORTED_IMPORT_FILE_TYPES } from '@/consts/file.ts';

import {
    type CategorizedAccount,
    Account,
    AccountWithDisplayBalance,
    CategorizedAccountWithDisplayBalance
} from '@/models/account.ts';
import type { LatestExchangeRateResponse, LocalizedLatestExchangeRate } from '@/models/exchange_rate.ts';

import {
    isDefined,
    isObject,
    isString,
    isNumber,
    isBoolean
} from '@/lib/common.ts';

import {
    isPM,
    formatUnixTime,
    formatCurrentTime,
    formatDate,
    parseDateFromUnixTime,
    getYear,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffset,
    getBrowserTimezoneOffsetMinutes,
    getTimeDifferenceHoursAndMinutes,
    getDateTimeFormatType,
    getRecentMonthDateRanges,
    isDateRangeMatchFullYears,
    isDateRangeMatchFullMonths,
    formatMonthDay
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

import {
    getCategorizedAccountsMap,
    getAllFilteredAccountsBalance
} from '@/lib/account.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

export interface LocalizedErrorParameter {
    readonly key: string;
    readonly localized: boolean;
    readonly value: string;
}

export interface LocalizedError {
    readonly message: string;
    readonly parameters?: LocalizedErrorParameter[];
}

export function getI18nOptions(): object {
    return {
        legacy: false,
        locale: DEFAULT_LANGUAGE,
        fallbackLocale: DEFAULT_LANGUAGE,
        formatFallbackMessages: true,
        messages: (function () {
            const messages: Record<string, object> = {};

            for (const languageKey in ALL_LANGUAGES) {
                if (!Object.prototype.hasOwnProperty.call(ALL_LANGUAGES, languageKey)) {
                    continue;
                }

                const languageInfo = ALL_LANGUAGES[languageKey];
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
    const exchangeRatesStore = useExchangeRatesStore();

    // private functions
    function getLanguageDisplayName(languageName: string): string {
        return t(`language.${languageName}`);
    }

    function getDefaultLanguage(): string {
        if (!window || !window.navigator) {
            return DEFAULT_LANGUAGE;
        }

        let browserLanguage = window.navigator.browserLanguage || window.navigator.language;

        if (!browserLanguage) {
            return DEFAULT_LANGUAGE;
        }

        if (!ALL_LANGUAGES[browserLanguage]) {
            const languageKey = getLanguageKeyFromLanguageAlias(browserLanguage);

            if (languageKey) {
                browserLanguage = languageKey;
            }
        }

        if (!ALL_LANGUAGES[browserLanguage] && browserLanguage.split('-').length > 1) { // maybe language-script-region
            const languageTagParts = browserLanguage.split('-');
            browserLanguage = languageTagParts[0] + '-' + languageTagParts[1];

            if (!ALL_LANGUAGES[browserLanguage]) {
                const languageKey = getLanguageKeyFromLanguageAlias(browserLanguage);

                if (languageKey) {
                    browserLanguage = languageKey;
                }
            }

            if (!ALL_LANGUAGES[browserLanguage]) {
                browserLanguage = languageTagParts[0];
                const languageKey = getLanguageKeyFromLanguageAlias(browserLanguage);

                if (languageKey) {
                    browserLanguage = languageKey;
                }
            }
        }

        if (!ALL_LANGUAGES[browserLanguage]) {
            return DEFAULT_LANGUAGE;
        }

        return browserLanguage;
    }

    function getLanguageKeyFromLanguageAlias(alias: string): string | null {
        for (const languageKey in ALL_LANGUAGES) {
            if (!Object.prototype.hasOwnProperty.call(ALL_LANGUAGES, languageKey)) {
                continue;
            }

            if (languageKey.toLowerCase() === alias.toLowerCase()) {
                return languageKey;
            }

            const langInfo = ALL_LANGUAGES[languageKey];
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
    function translateIf(text: string | undefined, isTranslate?: boolean): string {
        if (!isDefined(text)) {
            return '';
        }

        if (isTranslate) {
            return t(text);
        }

        return text;
    }

    function translateError(message: string | { error: ErrorResponse }): string {
        let finalMessage = '';
        let parameters = {};

        if (isObject(message) && isObject(message.error)) {
            const localizedError = getLocalizedError(message.error);
            finalMessage = localizedError.message;
            parameters = getLocalizedErrorParameters(localizedError.parameters);
        } else if (isString(message)) {
            finalMessage = message;
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

    function getServerTipContent(tipConfig: Record<string, string>): string {
        if (!tipConfig) {
            return '';
        }

        const currentLanguage = getCurrentLanguageTag();

        if (isString(tipConfig[currentLanguage])) {
            return tipConfig[currentLanguage];
        }

        return tipConfig['default'] || '';
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

    function getCurrentLanguageDisplayName(): string {
        const currentLanguageInfo = getCurrentLanguageInfo();
        return currentLanguageInfo.displayName;
    }

    function getDefaultCurrency(): string {
        return t('default.currency');
    }

    function getDefaultFirstDayOfWeek(): string {
        return t('default.firstDayOfWeek');
    }

    function getDefaultFiscalYearStart(): string {
        return t('default.fiscalYearStart');
    }

    function getAllLanguageOptions(includeSystemDefault: boolean): LanguageOption[] {
        const ret: LanguageOption[] = [];

        for (const languageTag in ALL_LANGUAGES) {
            if (!Object.prototype.hasOwnProperty.call(ALL_LANGUAGES, languageTag)) {
                continue;
            }

            const languageInfo = ALL_LANGUAGES[languageTag];
            const displayName = languageInfo.displayName;
            const languageNameInCurrentLanguage = getLanguageDisplayName(languageInfo.name);

            ret.push({
                languageTag: languageTag,
                displayName: languageNameInCurrentLanguage,
                nativeDisplayName: displayName
            });
        }

        ret.sort(function (lang1, lang2) {
            return lang1.languageTag.localeCompare(lang2.languageTag);
        });

        if (includeSystemDefault) {
            ret.splice(0, 0, {
                languageTag: '',
                displayName: '',
                nativeDisplayName: t('System Default')
            });
        }

        return ret;
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

    function getAllCurrencies(): LocalizedCurrencyInfo[] {
        const allCurrencies: LocalizedCurrencyInfo[] = [];

        for (const currencyCode in ALL_CURRENCIES) {
            if (!Object.prototype.hasOwnProperty.call(ALL_CURRENCIES, currencyCode)) {
                continue;
            }

            const localizedCurrencyInfo: LocalizedCurrencyInfo = {
                currencyCode: currencyCode,
                displayName: getCurrencyName(currencyCode)
            };

            allCurrencies.push(localizedCurrencyInfo);
        }

        allCurrencies.sort(function(c1, c2) {
            return c1.displayName.localeCompare(c2.displayName);
        })

        return allCurrencies;
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

    function getAllWeekDays(firstDayOfWeek?: number): TypeAndDisplayName[] {
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

    function getLocalizedDateTimeFormats<T extends DateFormat | TimeFormat>(type: string, allFormatMap: Record<string, T>, allFormatArray: T[], languageDefaultTypeNameKey: string, systemDefaultFormatType: T): LocalizedDateTimeFormat[] {
        const defaultFormat = getLocalizedDateTimeFormat<T>(type, allFormatMap, allFormatArray, LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE, languageDefaultTypeNameKey, systemDefaultFormatType);
        const ret: LocalizedDateTimeFormat[] = [];

        ret.push({
            type: LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE,
            format: defaultFormat,
            displayName: `${t('Language Default')} (${formatCurrentTime(defaultFormat)})`
        });

        for (let i = 0; i < allFormatArray.length; i++) {
            const formatType = allFormatArray[i];
            const format = t(`format.${type}.${formatType.key}`);

            ret.push({
                type: formatType.type,
                format: format,
                displayName: formatCurrentTime(format)
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

    function getAllExpenseIncomeAmountColors(categoryType: CategoryType): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];
        let defaultAmountName = '';

        if (categoryType === CategoryType.Expense) {
            defaultAmountName = PresetAmountColor.DefaultExpenseColor.name;
        } else if (categoryType === CategoryType.Income) { // income
            defaultAmountName = PresetAmountColor.DefaultIncomeColor.name;
        }

        if (defaultAmountName) {
            defaultAmountName = t('color.amount.' + defaultAmountName);
        }

        ret.push({
            type: PresetAmountColor.SystemDefaultType,
            displayName: t('System Default') + (defaultAmountName ? ` (${defaultAmountName})` : '')
        });

        const allPresetAmountColors = PresetAmountColor.values();

        for (let i = 0; i < allPresetAmountColors.length; i++) {
            const amountColor = allPresetAmountColors[i];

            ret.push({
                type: amountColor.type,
                displayName: t('color.amount.' + amountColor.name)
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

    function getAllTransactionDefaultCategories(categoryType: 0 | CategoryType, locale: string): PartialRecord<CategoryType, LocalizedPresetCategory[]> {
        const allCategories: PartialRecord<CategoryType, LocalizedPresetCategory[]> = {};
        const categoryTypes: CategoryType[] = [];

        if (categoryType === 0) {
            categoryTypes.push(...ALL_CATEGORY_TYPES);
        } else {
            categoryTypes.push(categoryType);
        }

        for (let i = 0; i < categoryTypes.length; i++) {
            const categories: LocalizedPresetCategory[] = [];
            const categoryType = categoryTypes[i];
            let defaultCategories: PresetCategory[] = [];

            if (categoryType === CategoryType.Income) {
                defaultCategories = DEFAULT_INCOME_CATEGORIES;
            } else if (categoryType === CategoryType.Expense) {
                defaultCategories = DEFAULT_EXPENSE_CATEGORIES;
            } else if (categoryType === CategoryType.Transfer) {
                defaultCategories = DEFAULT_TRANSFER_CATEGORIES;
            }

            for (let j = 0; j < defaultCategories.length; j++) {
                const category = defaultCategories[j];

                const submitCategory: LocalizedPresetCategory = {
                    name: t('category.' + category.name, {}, { locale: locale }),
                    type: categoryType,
                    icon: category.categoryIconId,
                    color: category.color,
                    subCategories: []
                };

                for (let k = 0; k < category.subCategories.length; k++) {
                    const subCategory = category.subCategories[k];
                    const submitSubCategory: LocalizedPresetSubCategory = {
                        name: t('category.' + subCategory.name, {}, { locale: locale }),
                        type: categoryType,
                        icon: subCategory.categoryIconId,
                        color: subCategory.color
                    };

                    submitCategory.subCategories.push(submitSubCategory);
                }

                categories.push(submitCategory);
            }

            allCategories[categoryType] = categories;
        }

        return allCategories;
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

    function getAllSupportedImportFileTypes(): LocalizedImportFileType[] {
        const allSupportedImportFileTypes: LocalizedImportFileType[] = [];

        for (let i = 0; i < SUPPORTED_IMPORT_FILE_TYPES.length; i++) {
            const fileType = SUPPORTED_IMPORT_FILE_TYPES[i];
            let document: LocalizedImportFileDocument | undefined;

            if (fileType.document) {
                let documentLanguage = '';
                let documentDisplayLanguageName = '';
                let documentAnchor = '';

                if (fileType.document.supportMultiLanguages === true) {
                    documentLanguage = getCurrentLanguageTag();

                    if (SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE[documentLanguage]) {
                        documentAnchor = t(`document.anchor.export_and_import.${fileType.document.anchor}`);
                    } else {
                        documentLanguage = DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE;
                        documentAnchor = t(`document.anchor.export_and_import.${fileType.document.anchor}`, {}, { locale: documentLanguage });
                    }
                } else if (isString(fileType.document.supportMultiLanguages) && ALL_LANGUAGES[fileType.document.supportMultiLanguages]) {
                    documentLanguage = fileType.document.supportMultiLanguages;
                    documentAnchor = fileType.document.anchor;
                }

                if (documentLanguage && documentLanguage !== getCurrentLanguageTag()) {
                    documentDisplayLanguageName = getLanguageDisplayName(ALL_LANGUAGES[documentLanguage].name);
                }

                if (documentLanguage) {
                    documentLanguage = documentLanguage.replace(/-/g, '_');
                }

                if (documentAnchor) {
                    documentAnchor = documentAnchor.toLowerCase().replace(/ /g, '-');
                }

                if (documentLanguage === DEFAULT_LANGUAGE) {
                    documentLanguage = '';
                }

                document = {
                    language: documentLanguage,
                    displayLanguageName: documentDisplayLanguageName,
                    anchor: documentAnchor
                };
            } else {
                document = undefined;
            }

            const subTypes: LocalizedImportFileTypeSubType[] = [];

            if (fileType.subTypes) {
                for (let i = 0; i < fileType.subTypes.length; i++) {
                    const subType = fileType.subTypes[i];
                    const localizedSubType: LocalizedImportFileTypeSubType = {
                        type: subType.type,
                        displayName: t(subType.name),
                        extensions: subType.extensions
                    };

                    subTypes.push(localizedSubType);
                }
            }

            const supportedEncodings: LocalizedImportFileTypeSupportedEncodings[] = [];

            if (fileType.supportedEncodings) {
                for (let i = 0; i < fileType.supportedEncodings.length; i++) {
                    const encoding = fileType.supportedEncodings[i];
                    const localizedEncoding: LocalizedImportFileTypeSupportedEncodings = {
                        encoding: encoding,
                        displayName: t(`encoding.${encoding}`)
                    };

                    supportedEncodings.push(localizedEncoding);
                }
            }

            const localizedFileType: LocalizedImportFileType = {
                type: fileType.type,
                displayName: t(fileType.name),
                extensions: fileType.extensions,
                subTypes: subTypes.length ? subTypes : undefined,
                supportedEncodings: supportedEncodings.length ? supportedEncodings : undefined,
                dataFromTextbox: fileType.dataFromTextbox,
                document: document
            };
            allSupportedImportFileTypes.push(localizedFileType);
        }

        return allSupportedImportFileTypes;
    }

    function getLanguageInfo(languageKey: string): LanguageInfo | undefined {
        return ALL_LANGUAGES[languageKey];
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

        const allWeekDays = getAllWeekDays(firstDayOfWeek);
        const finalWeekdayNames = [];

        for (let i = 0; i < allWeekDays.length; i++) {
            const weekDay = allWeekDays[i];

            if (weekdayTypesMap[weekDay.type]) {
                finalWeekdayNames.push(weekDay.displayName);
            }
        }

        return joinMultiText(finalWeekdayNames);
    }

    // Returns FiscalYearStart object, to facilitate diverse uses and conversions
    function getCurrentFiscalYearStart(): FiscalYearStart {
        let fiscalYearStart = FiscalYearStart.fromNumber(userStore.currentUserFiscalYearStart);
        if ( fiscalYearStart ) {
            return fiscalYearStart;
        }
        if ( !fiscalYearStart ) {
            fiscalYearStart = FiscalYearStart.fromMonthDashDayString(getDefaultFiscalYearStart());
        }
        return FiscalYearStart.Default;
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

    function isLongDateMonthAfterYear(): boolean {
        return getLocalizedDateTimeType(LongDateFormat.all(), LongDateFormat.values(), userStore.currentUserLongDateFormat, 'longDateFormat', LongDateFormat.Default).isMonthAfterYear;
    }

    function isShortDateMonthAfterYear(): boolean {
        return getLocalizedDateTimeType(ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default).isMonthAfterYear;
    }

    function isLongTime24HourFormat(): boolean {
        return getLocalizedDateTimeType(LongTimeFormat.all(), LongTimeFormat.values(), userStore.currentUserLongTimeFormat, 'longTimeFormat', LongTimeFormat.Default).is24HourFormat;
    }

    function isLongTimeMeridiemIndicatorFirst(): boolean {
        return getLocalizedDateTimeType(LongTimeFormat.all(), LongTimeFormat.values(), userStore.currentUserLongTimeFormat, 'longTimeFormat', LongTimeFormat.Default).isMeridiemIndicatorFirst || false;
    }

    function isShortTime24HourFormat(): boolean {
        return getLocalizedDateTimeType(ShortTimeFormat.all(), ShortTimeFormat.values(), userStore.currentUserShortTimeFormat, 'shortTimeFormat', ShortTimeFormat.Default).is24HourFormat;
    }

    function isShortTimeMeridiemIndicatorFirst(): boolean {
        return getLocalizedDateTimeType(ShortTimeFormat.all(), ShortTimeFormat.values(), userStore.currentUserShortTimeFormat, 'shortTimeFormat', ShortTimeFormat.Default).isMeridiemIndicatorFirst || false;
    }

    function formatDateToLongDate(date: string): string {
        return formatDate(date, getLocalizedLongDateFormat());
    }

    function formatMonthDayToLongDate(monthDay: string): string {
        return formatMonthDay(monthDay, getLocalizedLongMonthDayFormat());
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

    function formatDateRange(dateType: number, startTime: number, endTime: number): string {
        if (dateType === DateRange.All.type) {
            return t(DateRange.All.name);
        }

        const allDateRanges = DateRange.values();

        for (let i = 0; i < allDateRanges.length; i++) {
            const dateRange = allDateRanges[i];

            if (dateRange && dateRange.type !== DateRange.Custom.type && dateRange.type === dateType && dateRange.name) {
                return t(dateRange.name);
            }
        }

        if (isDateRangeMatchFullYears(startTime, endTime)) {
            const format = getLocalizedShortYearFormat();
            const displayStartTime = formatUnixTime(startTime, format);
            const displayEndTime = formatUnixTime(endTime, format);

            return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
        }

        if (isDateRangeMatchFullMonths(startTime, endTime)) {
            const format = getLocalizedShortYearMonthFormat();
            const displayStartTime = formatUnixTime(startTime, format);
            const displayEndTime = formatUnixTime(endTime, format);

            return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
        }

        const startTimeYear = getYear(parseDateFromUnixTime(startTime));
        const endTimeYear = getYear(parseDateFromUnixTime(endTime));

        const format = getLocalizedShortDateFormat();
        const displayStartTime = formatUnixTime(startTime, format);
        const displayEndTime = formatUnixTime(endTime, format);

        if (displayStartTime === displayEndTime) {
            return displayStartTime;
        } else if (startTimeYear === endTimeYear) {
            const displayShortEndTime = formatUnixTime(endTime, getLocalizedShortMonthDayFormat());
            return `${displayStartTime} ~ ${displayShortEndTime}`;
        }

        return `${displayStartTime} ~ ${displayEndTime}`;
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

    function getFormattedAmountWithCurrency(value: number | string, currencyCode?: string | false, notConvertValue?: boolean, currencyDisplayType?: CurrencyDisplayType): string {
        if (isNumber(value)) {
            value = value.toString();
        }

        let textualValue = value;
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

        let finalCurrencyCode = '';

        if (!isBoolean(currencyCode) && !currencyCode) {
            finalCurrencyCode = userStore.currentUserDefaultCurrency;
        } else if (isBoolean(currencyCode) && !currencyCode) {
            finalCurrencyCode = '';
        } else {
            finalCurrencyCode = currencyCode;
        }

        if (!finalCurrencyCode) {
            return textualValue;
        }

        if (!currencyDisplayType) {
            currencyDisplayType = getCurrentCurrencyDisplayType();
        }

        const currencyUnit = getCurrencyUnitName(finalCurrencyCode, isPlural);
        const currencyName = getCurrencyName(finalCurrencyCode);
        return appendCurrencySymbol(textualValue, currencyDisplayType, finalCurrencyCode, currencyUnit, currencyName, isPlural);
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

    function getCategorizedAccountsWithDisplayBalance(allVisibleAccounts: Account[], showAccountBalance: boolean): CategorizedAccountWithDisplayBalance[] {
        const ret: CategorizedAccountWithDisplayBalance[] = [];
        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const allCategories = AccountCategory.values();
        const categorizedAccounts: Record<number, CategorizedAccount> = getCategorizedAccountsMap(Account.cloneAccounts(allVisibleAccounts));

        for (let i = 0; i < allCategories.length; i++) {
            const category = allCategories[i];

            if (!categorizedAccounts[category.type]) {
                continue;
            }

            const accountCategory = categorizedAccounts[category.type];
            const accountsWithDisplayBalance: AccountWithDisplayBalance[] = [];

            if (accountCategory.accounts) {
                for (let i = 0; i < accountCategory.accounts.length; i++) {
                    const account = accountCategory.accounts[i];
                    let accountWithDisplaceBalance: AccountWithDisplayBalance;

                    if (showAccountBalance && account.isAsset) {
                        accountWithDisplaceBalance = AccountWithDisplayBalance.fromAccount(account, getFormattedAmountWithCurrency(account.balance, account.currency));
                    } else if (showAccountBalance && account.isLiability) {
                        accountWithDisplaceBalance = AccountWithDisplayBalance.fromAccount(account, getFormattedAmountWithCurrency(-account.balance, account.currency));
                    } else {
                        accountWithDisplaceBalance = AccountWithDisplayBalance.fromAccount(account, '***');
                    }

                    accountsWithDisplayBalance.push(accountWithDisplaceBalance);
                }
            }

            let finalTotalBalance = '';

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

                finalTotalBalance = totalBalance.toString();

                if (hasUnCalculatedAmount) {
                    finalTotalBalance = finalTotalBalance + '+';
                }

                finalTotalBalance = getFormattedAmountWithCurrency(finalTotalBalance, defaultCurrency);
            } else {
                finalTotalBalance = '***';
            }

            const accountCategoryWithDisplayBalance = CategorizedAccountWithDisplayBalance.of(accountCategory, accountsWithDisplayBalance, finalTotalBalance);
            ret.push(accountCategoryWithDisplayBalance);
        }

        return ret;
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

    function initLocale(lastUserLanguage?: string, timezone?: string): LocaleDefaultSettings | null {
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
        getServerTipContent,
        // get current language info
        getCurrentLanguageTag,
        getCurrentLanguageInfo,
        getCurrentLanguageDisplayName,
        // get localization default type
        getDefaultCurrency,
        getDefaultFirstDayOfWeek,
        getDefaultFiscalYearStart,
        // get all localized info of specified type
        getAllLanguageOptions,
        getAllEnableDisableOptions,
        getAllCurrencies,
        getAllMeridiemIndicators,
        getAllLongMonthNames,
        getAllShortMonthNames,
        getAllLongWeekdayNames,
        getAllShortWeekdayNames,
        getAllMinWeekdayNames,
        getAllWeekDays,
        getAllLongDateFormats: () => getLocalizedDateTimeFormats<LongDateFormat>('longDate', LongDateFormat.all(), LongDateFormat.values(), 'longDateFormat', LongDateFormat.Default),
        getAllShortDateFormats: () => getLocalizedDateTimeFormats<ShortDateFormat>('shortDate', ShortDateFormat.all(), ShortDateFormat.values(), 'shortDateFormat', ShortDateFormat.Default),
        getAllLongTimeFormats: () => getLocalizedDateTimeFormats<LongTimeFormat>('longTime', LongTimeFormat.all(), LongTimeFormat.values(), 'longTimeFormat', LongTimeFormat.Default),
        getAllShortTimeFormats: () => getLocalizedDateTimeFormats<ShortTimeFormat>('shortTime', ShortTimeFormat.all(), ShortTimeFormat.values(), 'shortTimeFormat', ShortTimeFormat.Default),
        getAllDateRanges,
        getAllRecentMonthDateRanges,
        getAllTimezones,
        getAllTimezoneTypesUsedForStatistics,
        getAllDecimalSeparators: () => getLocalizedNumeralSeparatorFormats(DecimalSeparator.values(), DecimalSeparator.parse(t('default.decimalSeparator')), DecimalSeparator.Default, DecimalSeparator.LanguageDefaultType),
        getAllDigitGroupingSymbols: () => getLocalizedNumeralSeparatorFormats(DigitGroupingSymbol.values(), DigitGroupingSymbol.parse(t('default.digitGroupingSymbol')), DigitGroupingSymbol.Default, DigitGroupingSymbol.LanguageDefaultType),
        getAllDigitGroupingTypes,
        getAllCurrencyDisplayTypes,
        getAllCurrencySortingTypes: () => getLocalizedDisplayNameAndType(CurrencySortingType.values()),
        getAllExpenseAmountColors: () => getAllExpenseIncomeAmountColors(CategoryType.Expense),
        getAllIncomeAmountColors: () => getAllExpenseIncomeAmountColors(CategoryType.Income),
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
        getAllImportTransactionColumnTypes: () => getLocalizedDisplayNameAndType(ImportTransactionColumnType.values()),
        getAllTransactionDefaultCategories,
        getAllDisplayExchangeRates,
        getAllSupportedImportFileTypes,
        // get localized info
        getLanguageInfo,
        getMonthShortName,
        getMonthLongName,
        getMonthdayOrdinal,
        getMonthdayShortName,
        getWeekdayShortName,
        getWeekdayLongName,
        getMultiMonthdayShortNames,
        getMultiWeekdayLongNames,
        getCurrentFiscalYearStart,
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
        formatUnixTimeToLongMonth: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, 'MMMM', utcOffset, currentUtcOffset),
        formatUnixTimeToShortMonth: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, 'MMM', utcOffset, currentUtcOffset),
        formatUnixTimeToLongYearMonth: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongYearMonthFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortYearMonth: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortYearMonthFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongMonthDay: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongMonthDayFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortMonthDay: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortMonthDayFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToLongTime: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedLongTimeFormat(), utcOffset, currentUtcOffset),
        formatUnixTimeToShortTime: (unixTime: number, utcOffset?: number, currentUtcOffset?: number) => formatUnixTime(unixTime, getLocalizedShortTimeFormat(), utcOffset, currentUtcOffset),
        formatDateToLongDate,
        formatMonthDayToLongDate,
        formatYearQuarter,
        formatDateRange,
        getTimezoneDifferenceDisplayText,
        appendDigitGroupingSymbol: getNumberWithDigitGroupingSymbol,
        parseAmount: getParsedAmountNumber,
        formatAmount: getFormattedAmount,
        formatAmountWithCurrency: getFormattedAmountWithCurrency,
        formatExchangeRateAmount: getFormattedExchangeRateAmount,
        getAdaptiveAmountRate,
        getAmountPrependAndAppendText,
        getCategorizedAccountsWithDisplayBalance,
        // localization setting functions
        setLanguage,
        setTimeZone,
        initLocale
    };
}
