import { useI18n as useVueI18n } from 'vue-i18n';
import moment from 'moment-timezone';

import type { NameValue, TypeAndName, TypeAndDisplayName, LocalizedSwitchOption } from '@/core/base.ts';

import {
    type LanguageInfo,
    type LanguageOption,
    ALL_LANGUAGES,
    DEFAULT_LANGUAGE
} from '@/locales/index.ts';

import {
    ALL_LANGUAGES as CHINESE_CALENDAR_ALL_LANGUAGES,
    DEFAULT_CONTENT as CHINESE_CALENDAR_DEFAULT_CONTENT
} from '@/locales/calendar/chinese/index.ts';

import {
    ALL_LANGUAGES as PERSIAN_CALENDAR_ALL_LANGUAGES,
    DEFAULT_CONTENT as PERSIAN_CALENDAR_DEFAULT_CONTENT
} from '@/locales/calendar/persian/index.ts';

import {
    entries,
    keys
} from '@/core/base.ts';

import {
    TextDirection
} from '@/core/text.ts';

import {
    type ChineseCalendarLocaleData,
    type PersianCalendarLocaleData,
    CalendarType,
    CalendarDisplayType,
    DateDisplayType
} from '@/core/calendar.ts';

import {
    type DateTime,
    type DateTimeFormatOptions,
    type DateTimeLocaleData,
    type TextualMonthDay,
    type TextualYearMonthDay,
    type Year1BasedMonth,
    type YearMonthDay,
    type CalendarAlternateDate,
    type DateFormat,
    type TimeFormat,
    type LocalizedDateTimeFormat,
    type LocalizedDateRange,
    type LocalizedRecentMonthDateRange,
    type UnixTimeRange,
    type WeekDayValue,
    Month,
    WeekDay,
    MeridiemIndicator,
    KnownDateTimeFormat,
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
    type HiddenAmount,
    type NumberFormatOptions,
    type NumberWithSuffix,
    type NumeralSymbolType,
    type LocalizedNumeralSymbolType,
    type LocalizedDigitGroupingType,
    NumeralSystem,
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
    FiscalYearStart,
    FiscalYearFormat,
    FiscalYearUnixTime,
    LANGUAGE_DEFAULT_FISCAL_YEAR_FORMAT_VALUE,
} from '@/core/fiscalyear.ts';

import {
    CoordinateDisplayType
} from '@/core/coordinate.ts';

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
    TransactionEditScopeType
} from '@/core/transaction.ts';

import {
    ImportTransactionColumnType
} from '@/core/import_transaction.ts';

import {
    ScheduledTemplateFrequencyType
} from '@/core/template.ts';

import {
    StatisticsAnalysisType,
    CategoricalChartType,
    TrendChartType,
    AccountBalanceTrendChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import {
    TransactionExploreConditionField,
    TransactionExploreConditionOperator,
    TransactionExploreDataDimension,
    TransactionExploreValueMetric,
    TransactionExploreChartType
} from '@/core/explore.ts';

import {
    type LocalizedImportFileCategoryAndTypes,
    type LocalizedImportFileType,
    type LocalizedImportFileTypeSubType,
    type LocalizedImportFileTypeSupportedEncodings,
    type LocalizedImportFileDocument
} from '@/core/file.ts';

import type { LocaleDefaultSettings } from '@/core/setting.ts';
import type { ErrorResponse } from '@/core/api.ts';

import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';
import { UTC_TIMEZONE, ALL_TIMEZONES } from '@/consts/timezone.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';
import { DEFAULT_EXPENSE_CATEGORIES, DEFAULT_INCOME_CATEGORIES, DEFAULT_TRANSFER_CATEGORIES } from '@/consts/category.ts';
import { KnownErrorCode, SPECIFIED_API_NOT_FOUND_ERRORS, PARAMETERIZED_ERRORS } from '@/consts/api.ts';
import { OAUTH2_PROVIDER_DISPLAY_NAME } from '@/consts/oauth2.ts';
import { DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE, SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE, SUPPORTED_IMPORT_FILE_CATEGORY_AND_TYPES } from '@/consts/file.ts';

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
    formatCurrentTime,
    formatGregorianCalendarYearDashMonthDashDay,
    formatGregorianCalendarMonthDashDay,
    formatDateTime,
    formatUnixTime,
    getBrowserTimezoneOffset,
    getBrowserTimezoneOffsetMinutes,
    getCurrentUnixTime,
    getYearMonthDayDateTime,
    parseDateTimeFromUnixTime,
    parseDateTimeFromString,
    getGregorianCalendarYearMonthDays,
    getDateTimeFormatType,
    getFiscalYearTimeRangeFromUnixTime,
    getFiscalYearTimeRangeFromYear,
    getRecentMonthDateRanges,
    getTimeDifferenceHoursAndMinutes,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneName,
    isDateRangeMatchFullMonths,
    isDateRangeMatchFullYears,
    isPM
} from '@/lib/datetime.ts';

import {
    type ChineseYearMonthDayInfo,
    getChineseYearMonthAllDayInfos,
    getChineseYearMonthDayInfo,
    getChineseCalendarAlternateDisplayDate
} from '@/lib/calendar/chinese_calendar.ts';

import {
    appendDigitGroupingSymbolAndDecimalSeparator,
    parseAmount,
    formatAmount,
    formatHiddenAmount,
    formatNumber,
    formatPercent,
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

import {
    getSessionCurrentLanguageKey,
    setSessionCurrentLanguageKey
} from '@/lib/settings.ts';

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

            for (const [languageKey, languageInfo] of entries(ALL_LANGUAGES)) {
                messages[languageKey] = languageInfo.content;
            }

            return messages;
        })()
    };
}

export function getRtlLocales(): Record<string, boolean> {
    const rtlLocales: Record<string, boolean> = {};

    for (const [languageKey, languageInfo] of entries(ALL_LANGUAGES)) {
        if (languageInfo.textDirection === 'rtl') {
            rtlLocales[languageKey] = true;
        }
    }

    return rtlLocales;
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

        // try to match the full browser language tag with full language tag in i18n file
        if (ALL_LANGUAGES[browserLanguage]) {
            return browserLanguage;
        }

        // try to match the full browser language tag with language alias tags in i18n file
        let alternativeLanguage = getLanguageKeyFromLanguageAlias(browserLanguage);

        if (alternativeLanguage && ALL_LANGUAGES[alternativeLanguage]) {
            return alternativeLanguage;
        }

        const languageTagParts = browserLanguage.split('-');

        // maybe browser language is language-script-region format
        if (languageTagParts.length > 2) {
            // fallback to use language tag with language-script / language-region format
            browserLanguage = languageTagParts[0] + '-' + languageTagParts[1];

            // try to match language tag in language-script / language-region format with full language tag in i18n file
            if (ALL_LANGUAGES[browserLanguage]) {
                return browserLanguage;
            }

            // try to match language tag in language-script / language-region format with language alias tags in i18n file
            alternativeLanguage = getLanguageKeyFromLanguageAlias(browserLanguage);

            if (alternativeLanguage && ALL_LANGUAGES[alternativeLanguage]) {
                return alternativeLanguage;
            }
        }

        // fallback to use marco language tag
        if (languageTagParts.length > 1) {
            browserLanguage = languageTagParts[0] as string;

            // try to match marco language tag with full language tag in i18n file
            if (ALL_LANGUAGES[browserLanguage]) {
                return browserLanguage;
            }

            // try to match marco language tag with language alias tags in i18n file
            alternativeLanguage = getLanguageKeyFromLanguageAlias(browserLanguage);

            if (alternativeLanguage && ALL_LANGUAGES[alternativeLanguage]) {
                return alternativeLanguage;
            }
        }

        // fallback to match marco language tag with marco language tag in i18n file
        alternativeLanguage = getLanguageKeyFromMarcoLanguageTag(browserLanguage);

        if (alternativeLanguage && ALL_LANGUAGES[alternativeLanguage]) {
            return alternativeLanguage;
        }

        // fallback to use the default language
        return DEFAULT_LANGUAGE;
    }

    function getLanguageKeyFromLanguageAlias(alias: string): string | null {
        for (const [languageKey, languageInfo] of entries(ALL_LANGUAGES)) {
            if (languageKey.toLowerCase() === alias.toLowerCase()) {
                return languageKey;
            }

            const aliases = languageInfo.aliases;

            if (!aliases || aliases.length < 1) {
                continue;
            }

            for (const aliasName of aliases) {
                if (aliasName.toLowerCase() === alias.toLowerCase()) {
                    return languageKey;
                }
            }
        }

        return null;
    }

    function getLanguageKeyFromMarcoLanguageTag(languageTag: string): string | null {
        for (const languageKey of keys(ALL_LANGUAGES)) {
            if (languageKey.indexOf('-') < 0) {
                continue;
            }

            const marcoLanguageTag = languageKey.split('-')[0] as string;

            if (marcoLanguageTag.toLowerCase() === languageTag.toLowerCase()) {
                return languageKey;
            }
        }

        return null;
    }

    function getLocalizedError(error: ErrorResponse): LocalizedError {
        if (error.errorCode === KnownErrorCode.ApiNotFound && SPECIFIED_API_NOT_FOUND_ERRORS[error.path]) {
            return {
                message: `${SPECIFIED_API_NOT_FOUND_ERRORS[error.path]!.message}`
            };
        }

        if (error.errorCode !== KnownErrorCode.ValidatorError) {
            return {
                message: `error.${error.errorMessage}`
            };
        }

        for (const errorInfo of PARAMETERIZED_ERRORS) {
            const matches = error.errorMessage.match(errorInfo.regex);

            if (matches && matches.length === errorInfo.parameters.length + 1) {
                return {
                    message: `parameterizedError.${errorInfo.localeKey}`,
                    parameters: errorInfo.parameters.map((param, index) => {
                        return {
                            key: param.field,
                            localized: param.localized,
                            value: matches[index + 1] as string
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
            for (const parameter of parameters) {
                if (parameter.localized) {
                    localizedParameters[parameter.key] = t(`parameter.${parameter.value}`);
                } else {
                    localizedParameters[parameter.key] = parameter.value;
                }
            }
        }

        return localizedParameters;
    }

    function getDateTimeLocaleData(): DateTimeLocaleData {
        return moment.localeData();
    }

    function getChineseCalendarLocaleData(): ChineseCalendarLocaleData {
        const localeData = CHINESE_CALENDAR_ALL_LANGUAGES[locale.value] ?? CHINESE_CALENDAR_DEFAULT_CONTENT;
        const chineseCalendarLocaleData: ChineseCalendarLocaleData = {
            numerals: localeData['numerals'],
            monthNames: localeData['monthNames'],
            dayNames: localeData['dayNames'],
            leapMonthPrefix: localeData['leapMonthPrefix'],
            solarTermNames: localeData['solarTermNames']
        };

        return chineseCalendarLocaleData;
    }

    function getPersianCalendarLocaleData(): PersianCalendarLocaleData {
        const localeData = PERSIAN_CALENDAR_ALL_LANGUAGES[locale.value] ?? PERSIAN_CALENDAR_DEFAULT_CONTENT;
        const persianCalendarLocaleData: PersianCalendarLocaleData = {
            monthNames: localeData['monthNames'],
            monthShortNames: localeData['monthShortNames']
        };

        return persianCalendarLocaleData;
    }

    function getAllCurrencyDisplayTypes(numeralSystem: NumeralSystem, decimalSeparator: string): TypeAndDisplayName[] {
        const defaultCurrencyDisplayTypeName = t('default.currencyDisplayType');
        let defaultCurrencyDisplayType = CurrencyDisplayType.parse(defaultCurrencyDisplayTypeName);

        if (!defaultCurrencyDisplayType) {
            defaultCurrencyDisplayType = CurrencyDisplayType.Default;
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;

        const ret = [];
        const defaultSampleValue = getFormattedAmountWithCurrency(12345, defaultCurrency, defaultCurrencyDisplayType, numeralSystem, decimalSeparator);

        ret.push({
            type: CurrencyDisplayType.LanguageDefaultType,
            displayName: `${t('Language Default')} (${defaultSampleValue})`
        });

        const allCurrencyDisplayTypes = CurrencyDisplayType.values();

        for (const type of allCurrencyDisplayTypes) {
            const sampleValue = getFormattedAmountWithCurrency(12345, defaultCurrency, type, numeralSystem, decimalSeparator);
            const displayName = `${t(type.name)} (${sampleValue})`

            ret.push({
                type: type.type,
                displayName: displayName
            });
        }

        return ret;
    }

    function getAllLocalizedCalendarTypes(allCalendarDisplayTypeArray: CalendarDisplayType[] | DateDisplayType[], localeDefaultType: CalendarDisplayType | DateDisplayType | undefined, systemDefaultType: CalendarDisplayType | DateDisplayType, languageDefaultValue: number): TypeAndDisplayName[] {
        let defaultType: TypeAndName | undefined = localeDefaultType;

        if (!defaultType) {
            defaultType = systemDefaultType;
        }

        const ret: TypeAndDisplayName[] = [];

        ret.push({
            type: languageDefaultValue,
            displayName: `${t('Language Default')} (${t('calendar.' + defaultType.name)})`
        });

        for (const calendarDisplayType of allCalendarDisplayTypeArray) {
            ret.push({
                type: calendarDisplayType.type,
                displayName: t('calendar.' + calendarDisplayType.name)
            });
        }

        return ret;
    }

    function getLocalizedDisplayNameAndType(typeAndNames: TypeAndName[]): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];

        for (const typeAndName of typeAndNames) {
            ret.push({
                type: typeAndName.type,
                displayName: t(typeAndName.name)
            });
        }

        return ret;
    }

    function getLocalizedNameValue(nameValues: NameValue[]): NameValue[] {
        const ret: NameValue[] = [];

        for (const nameValue of nameValues) {
            ret.push({
                name: t(nameValue.name),
                value: nameValue.value
            });
        }

        return ret;
    }

    function getLocalizedDisplayNameAndTypeWithSystemDefault(typeAndNames: TypeAndName[], defaultValue: number, defaultType: TypeAndName): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];

        ret.push({
            type: defaultValue,
            displayName: t('System Default') + (defaultType.name ? ` (${t(defaultType.name)})` : '')
        });

        for (const typeAndName of typeAndNames) {
            ret.push({
                type: typeAndName.type,
                displayName: t(typeAndName.name)
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

        for (const separator of allSeparatorArray) {
            ret.push({
                type: separator.type,
                symbol: separator.symbol,
                displayName: `${t('numeral.' + separator.name)} (${separator.symbol})`
            });
        }

        return ret;
    }

    function getLocalizedChartDateAggregationTypeAndDisplayName(analysisType: StatisticsAnalysisType, fullName: boolean): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];
        const allTypes: ChartDateAggregationType[] = ChartDateAggregationType.values(analysisType);

        for (const type of allTypes) {
            ret.push({
                type: type.type,
                displayName: t(fullName ? type.fullName : `granularity.${type.shortName}`)
            });
        }

        return ret;
    }

    function getAllMonthNames(type: string): string[] {
        const ret = [];
        const allMonths = Month.values();

        for (const month of allMonths) {
            ret.push(t(`datetime.${month.name}.${type}`));
        }

        return ret;
    }

    function getAllWeekdayNames(type: string): string[] {
        const ret = [];
        const allWeekDays = WeekDay.values();

        for (const weekDay of allWeekDays) {
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

    function getLocalizedShortDayFormat(): string {
        return getLocalizedDateTimeFormat<ShortDateFormat>('shortDay', ShortDateFormat.all(), ShortDateFormat.values(), userStore.currentUserShortDateFormat, 'shortDateFormat', ShortDateFormat.Default);
    }

    function getLocalizedLongTimeFormat(): string {
        return getLocalizedDateTimeFormat<LongTimeFormat>('longTime', LongTimeFormat.all(), LongTimeFormat.values(), userStore.currentUserLongTimeFormat, 'longTimeFormat', LongTimeFormat.Default);
    }

    function getLocalizedShortTimeFormat(): string {
        return getLocalizedDateTimeFormat<ShortTimeFormat>('shortTime', ShortTimeFormat.all(), ShortTimeFormat.values(), userStore.currentUserShortTimeFormat, 'shortTimeFormat', ShortTimeFormat.Default);
    }

    function getDateTimeFormatOptions(options?: { calendarType?: CalendarType, numeralSystem?: NumeralSystem }): DateTimeFormatOptions {
        let numeralSystem: NumeralSystem | undefined = options?.numeralSystem;
        let calendarType: CalendarType | undefined = options?.calendarType;

        if (!isDefined(numeralSystem)) {
            numeralSystem = getCurrentNumeralSystemType();
        }

        if (!isDefined(calendarType)) {
            calendarType = getCurrentDateDisplayType().calendarType;
        }

        return {
            numeralSystem: numeralSystem,
            calendarType: calendarType,
            localeData: getDateTimeLocaleData(),
            chineseCalendarLocaleData: getChineseCalendarLocaleData(),
            persianCalendarLocaleData: getPersianCalendarLocaleData()
        };
    }

    function getNumberFormatOptions({numeralSystem, digitGrouping, decimalSeparator, currencyCode}: {
        numeralSystem?: NumeralSystem,
        digitGrouping?: DigitGroupingType,
        decimalSeparator?: string,
        currencyCode?: string
    }): NumberFormatOptions {
        if (!isDefined(numeralSystem)) {
            numeralSystem = getCurrentNumeralSystemType();
        }

        if (!isDefined(digitGrouping)) {
            digitGrouping = getCurrentDigitGroupingType();
        }

        if (!isDefined(decimalSeparator)) {
            decimalSeparator = getCurrentDecimalSeparator();
        }

        return {
            numeralSystem: numeralSystem,
            digitGrouping: digitGrouping,
            digitGroupingSymbol: getCurrentDigitGroupingSymbol(),
            decimalSeparator: decimalSeparator,
            decimalNumberCount: getCurrencyFraction(currencyCode),
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

    function isGregorianLikeCalendarType(calendarType: CalendarType): boolean {
        return calendarType === CalendarType.Gregorian || calendarType === CalendarType.Buddhist;
    }

    function getGregorianLikeCalendarType(): CalendarType {
        const currentDateDisplayType = getCurrentDateDisplayType();

        if (isGregorianLikeCalendarType(currentDateDisplayType.calendarType)) {
            return currentDateDisplayType.calendarType;
        }

        return CalendarType.Gregorian;
    }

    function formatYearQuarter(year: string, quarter: number): string {
        if (1 <= quarter && quarter <= 4) {
            return t('format.yearQuarter.q' + quarter, {
                year: year,
                quarter: quarter
            });
        } else {
            return '';
        }
    }

    function formatTimeRangeToGregorianLikeFiscalYearFormat(format: FiscalYearFormat, timeRange: FiscalYearUnixTime | UnixTimeRange, numeralSystem?: NumeralSystem, calendarType?: CalendarType): string {
        if (!format) {
            format = FiscalYearFormat.Default;
        }

        if (!isDefined(calendarType)) {
            calendarType = getGregorianLikeCalendarType();
        } else if (!isGregorianLikeCalendarType(calendarType)) {
            calendarType = CalendarType.Gregorian;
        }

        const dateTimeFormatOptions = getDateTimeFormatOptions({
            calendarType: calendarType,
            numeralSystem: numeralSystem
        });

        return t('format.fiscalYear.' + format.typeName, {
            StartYYYY: formatUnixTime(timeRange.minUnixTime, 'YYYY', dateTimeFormatOptions),
            StartYY: formatUnixTime(timeRange.minUnixTime, 'YY', dateTimeFormatOptions),
            EndYYYY: formatUnixTime(timeRange.maxUnixTime, 'YYYY', dateTimeFormatOptions),
            EndYY: formatUnixTime(timeRange.maxUnixTime, 'YY', dateTimeFormatOptions),
        });
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

    function getCalendarAlternateDisplayDate(dateTime: DateTime, dateTimeFormatOptions: DateTimeFormatOptions): CalendarAlternateDate {
        const numeralSystem = getCurrentNumeralSystemType();
        let displayDate = numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(dateTime.getLocalizedCalendarDay(dateTimeFormatOptions));

        if (dateTime.isLocalizedCalendarFirstDayOfMonth(dateTimeFormatOptions)) {
            displayDate = dateTime.getLocalizedCalendarMonthDisplayShortName(dateTimeFormatOptions);
        }

        const alternateDate: CalendarAlternateDate = {
            year: dateTime.getGregorianCalendarYear(),
            month: dateTime.getGregorianCalendarMonth(),
            day: dateTime.getGregorianCalendarDay(),
            displayDate: displayDate
        };

        return alternateDate;
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

    function getServerMultiLanguageConfigContent(multiLanguageConfig: Record<string, string>): string {
        if (!multiLanguageConfig) {
            return '';
        }

        const currentLanguage = getCurrentLanguageTag();

        if (isString(multiLanguageConfig[currentLanguage])) {
            return multiLanguageConfig[currentLanguage];
        }

        return multiLanguageConfig['default'] || '';
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

    function getCurrentLanguageTextDirection(): TextDirection {
        const currentLanguageInfo = getCurrentLanguageInfo();

        if (currentLanguageInfo.textDirection === 'rtl') {
            return TextDirection.RTL;
        } else {
            return TextDirection.LTR;
        }
    }

    function getDefaultCurrency(): string {
        return t('default.currency');
    }

    function getDefaultFirstDayOfWeek(): string {
        return t('default.firstDayOfWeek');
    }

    function getAllLanguageOptions(includeSystemDefault: boolean): LanguageOption[] {
        const ret: LanguageOption[] = [];

        for (const [languageKey, languageInfo] of entries(ALL_LANGUAGES)) {
            const displayName = languageInfo.displayName;
            const languageNameInCurrentLanguage = getLanguageDisplayName(languageInfo.name);

            ret.push({
                languageTag: languageKey,
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
        }, {
            value: false,
            displayName: t('Disable')
        }];
    }

    function getAllCurrencies(): LocalizedCurrencyInfo[] {
        const allCurrencies: LocalizedCurrencyInfo[] = [];

        for (const currencyCode of keys(ALL_CURRENCIES)) {
            const localizedCurrencyInfo: LocalizedCurrencyInfo = {
                currencyCode: currencyCode,
                displayName: getCurrencyName(currencyCode)
            };

            allCurrencies.push(localizedCurrencyInfo);
        }

        allCurrencies.sort(function (c1, c2) {
            return c1.displayName.localeCompare(c2.displayName);
        })

        return allCurrencies;
    }

    function getAllMeridiemIndicators(): NameValue[] {
        const allMeridiemIndicators = MeridiemIndicator.values();
        const localizedMeridiemIndicatorNames = [];

        for (const indicator of allMeridiemIndicators) {
            localizedMeridiemIndicatorNames.push({
                name: t(`datetime.${indicator.name}.content`),
                value: indicator.name
            });
        }

        return localizedMeridiemIndicatorNames;
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

    function getAllWeekDays(firstDayOfWeek?: WeekDayValue): TypeAndDisplayName[] {
        const ret: TypeAndDisplayName[] = [];
        const allWeekDays = WeekDay.values();

        if (!isNumber(firstDayOfWeek)) {
            firstDayOfWeek = WeekDay.DefaultFirstDay.type;
        }

        for (let i = firstDayOfWeek; i < allWeekDays.length; i++) {
            const weekDay = allWeekDays[i] as WeekDay;

            ret.push({
                type: weekDay.type,
                displayName: t(`datetime.${weekDay.name}.long`)
            });
        }

        for (let i = 0; i < firstDayOfWeek; i++) {
            const weekDay = allWeekDays[i] as WeekDay;

            ret.push({
                type: weekDay.type,
                displayName: t(`datetime.${weekDay.name}.long`)
            });
        }

        return ret;
    }

    function getLocalizedDateTimeFormats<T extends DateFormat | TimeFormat>(type: string, allFormatMap: Record<string, T>, allFormatArray: T[], languageDefaultTypeNameKey: string, systemDefaultFormatType: T, numeralSystem: NumeralSystem, calendarType?: CalendarType): LocalizedDateTimeFormat[] {
        const defaultFormat = getLocalizedDateTimeFormat<T>(type, allFormatMap, allFormatArray, LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE, languageDefaultTypeNameKey, systemDefaultFormatType);
        const ret: LocalizedDateTimeFormat[] = [];
        const dateTimeFormatOptions = getDateTimeFormatOptions({ numeralSystem, calendarType });

        ret.push({
            type: LANGUAGE_DEFAULT_DATE_TIME_FORMAT_VALUE,
            format: defaultFormat,
            displayName: `${t('Language Default')} (${formatCurrentTime(defaultFormat, dateTimeFormatOptions)})`
        });

        for (const formatType of allFormatArray) {
            const format = t(`format.${type}.${formatType.key}`);

            ret.push({
                type: formatType.type,
                format: format,
                displayName: formatCurrentTime(format, dateTimeFormatOptions)
            });
        }

        return ret;
    }

    function getAllDateRanges(scene: DateRangeScene, includeCustom?: boolean, includeBillingCycle?: boolean): LocalizedDateRange[] {
        const ret: LocalizedDateRange[] = [];
        const allDateRanges = DateRange.values();

        for (const dateRange of allDateRanges) {
            if (!dateRange.isAvailableForScene(scene)) {
                continue;
            }

            if (dateRange.isBillingCycle) {
                if (includeBillingCycle) {
                    ret.push({
                        type: dateRange.type,
                        displayName: t(dateRange.name),
                        isBillingCycle: dateRange.isBillingCycle,
                        isUserCustomRange: dateRange.isUserCustomRange
                    });
                }

                continue;
            }

            if (includeCustom || dateRange.type !== DateRange.Custom.type) {
                ret.push({
                    type: dateRange.type,
                    displayName: t(dateRange.name),
                    isBillingCycle: dateRange.isBillingCycle,
                    isUserCustomRange: dateRange.isUserCustomRange
                });
            }
        }

        return ret;
    }

    function getAllRecentMonthDateRanges(includeAll: boolean, includeCustom: boolean): LocalizedRecentMonthDateRange[] {
        const allRecentMonthDateRanges: LocalizedRecentMonthDateRange[] = [];
        const recentDateRanges = getRecentMonthDateRanges(12);
        const currentCalendarDisplayType = getCurrentCalendarDisplayType();
        const dateTimeFormatOptions = getDateTimeFormatOptions({ calendarType: currentCalendarDisplayType.primaryCalendarType });

        if (includeAll) {
            allRecentMonthDateRanges.push({
                dateType: DateRange.All.type,
                minTime: 0,
                maxTime: 0,
                displayName: t('All')
            });
        }

        for (const recentDateRange of recentDateRanges) {
            allRecentMonthDateRanges.push({
                dateType: recentDateRange.dateType,
                minTime: recentDateRange.minTime,
                maxTime: recentDateRange.maxTime,
                year: recentDateRange.year,
                month: recentDateRange.month,
                isPreset: true,
                displayName: formatUnixTime(recentDateRange.minTime, getLocalizedLongYearMonthFormat(), dateTimeFormatOptions)
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

    function getAllTimezones(unixTime: number, includeSystemDefault?: boolean): LocalizedTimezoneInfo[] {
        const numeralSystem = getCurrentNumeralSystemType();
        const defaultTimezoneOffset = numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(getBrowserTimezoneOffset(unixTime));
        const defaultTimezoneOffsetMinutes = getBrowserTimezoneOffsetMinutes(unixTime);
        const allTimezoneInfos: LocalizedTimezoneInfo[] = [];

        for (const timezoneInfo of ALL_TIMEZONES) {
            const utcOffset = (timezoneInfo.timezoneName !== UTC_TIMEZONE.timezoneName ? numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(getTimezoneOffset(unixTime, timezoneInfo.timezoneName)) : '');
            const displayName = t(`timezone.${timezoneInfo.displayName}`);

            allTimezoneInfos.push({
                name: timezoneInfo.timezoneName,
                utcOffset: utcOffset,
                utcOffsetMinutes: getTimezoneOffsetMinutes(unixTime, timezoneInfo.timezoneName),
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

        allTimezoneInfos.sort(function (c1, c2) {
            const utcOffsetMinutes1 = c1.utcOffsetMinutes;
            const utcOffsetMinutes2 = c2.utcOffsetMinutes;

            if (utcOffsetMinutes1 !== utcOffsetMinutes2) {
                return utcOffsetMinutes1 - utcOffsetMinutes2;
            }

            return c1.displayName.localeCompare(c2.displayName);
        })

        return allTimezoneInfos;
    }

    function getAllTimezoneTypesUsedForStatistics(currentTimezone?: string): TypeAndDisplayName[] {
        const numeralSystem = getCurrentNumeralSystemType();
        const currentTimezoneOffset = numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(getTimezoneOffset(getCurrentUnixTime(), currentTimezone));

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

    function getAllFiscalYearFormats(numeralSystem: NumeralSystem, calendarType: CalendarType): TypeAndDisplayName[] {
        const now = getCurrentUnixTime();
        let fiscalYearStart = userStore.currentUserFiscalYearStart;

        if (!fiscalYearStart) {
            fiscalYearStart = FiscalYearStart.Default.value;
        }

        const currentFiscalYearRange = getFiscalYearTimeRangeFromUnixTime(now, fiscalYearStart);
        let defaultFiscalYearFormat = FiscalYearFormat.parse(t('default.fiscalYearFormat'));

        if (!defaultFiscalYearFormat) {
            defaultFiscalYearFormat = FiscalYearFormat.Default;
        }

        const ret: TypeAndDisplayName[] = [];

        ret.push({
            type: LANGUAGE_DEFAULT_FISCAL_YEAR_FORMAT_VALUE,
            displayName: `${t('Language Default')} (${formatTimeRangeToGregorianLikeFiscalYearFormat(defaultFiscalYearFormat, currentFiscalYearRange, numeralSystem, calendarType)})`
        });

        const allFiscalYearFormats = FiscalYearFormat.values();

        for (const fiscalYearFormat of allFiscalYearFormats) {
            ret.push({
                type: fiscalYearFormat.type,
                displayName: formatTimeRangeToGregorianLikeFiscalYearFormat(fiscalYearFormat, currentFiscalYearRange, numeralSystem, calendarType),
            });
        }

        return ret;
    }

    function getAllNumeralSystemTypes(): TypeAndDisplayName[] {
        const defaultNumeralSystemTypeName = t('default.numeralSystem');
        let defaultNumeralSystemType = NumeralSystem.parse(defaultNumeralSystemTypeName);

        if (!defaultNumeralSystemType) {
            defaultNumeralSystemType = NumeralSystem.Default;
        }

        const ret: TypeAndDisplayName[] = [];

        ret.push({
            type: NumeralSystem.LanguageDefaultType,
            displayName: `${t('Language Default')} (${defaultNumeralSystemType.textualAllDigits})`
        });

        const allNumeralSystemTypes = NumeralSystem.values();

        for (const type of allNumeralSystemTypes) {
            ret.push({
                type: type.type,
                displayName: `${t('numeral.' + type.name)} (${type.textualAllDigits})`
            });
        }

        return ret;
    }

    function getAllDigitGroupingTypes(numeralSystem: NumeralSystem, digitGroupingSymbol: string): LocalizedDigitGroupingType[] {
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
        const numberCharacters = numeralSystem.replaceWesternArabicDigitsToLocalizedDigits('123456789').split('');

        for (const type of allDigitGroupingTypes) {
            const sampleValue = type.format(numberCharacters, digitGroupingSymbol);

            ret.push({
                type: type.type,
                enabled: type.enabled,
                displayName: `${t('numeral.' + type.name)} (${sampleValue})`
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

        for (const amountColor of allPresetAmountColors) {
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

        for (const accountCategory of allCategories) {
            ret.push({
                type: accountCategory.type,
                displayName: t(accountCategory.name),
                defaultAccountIconId: accountCategory.defaultAccountIconId
            });
        }

        return ret;
    }

    function getAllTransactionDefaultCategories(categoryType: 0 | CategoryType, locale: string): Record<string, LocalizedPresetCategory[]> {
        const allCategories: Record<string, LocalizedPresetCategory[]> = {};
        const categoryTypes: CategoryType[] = [];

        if (categoryType === 0) {
            categoryTypes.push(...ALL_CATEGORY_TYPES);
        } else {
            categoryTypes.push(categoryType);
        }

        for (const categoryType of categoryTypes) {
            const categories: LocalizedPresetCategory[] = [];
            let defaultCategories: PresetCategory[] = [];

            if (categoryType === CategoryType.Income) {
                defaultCategories = DEFAULT_INCOME_CATEGORIES;
            } else if (categoryType === CategoryType.Expense) {
                defaultCategories = DEFAULT_EXPENSE_CATEGORIES;
            } else if (categoryType === CategoryType.Transfer) {
                defaultCategories = DEFAULT_TRANSFER_CATEGORIES;
            }

            for (const category of defaultCategories) {
                const submitCategory: LocalizedPresetCategory = {
                    name: t('category.' + category.name, {}, { locale: locale }),
                    type: categoryType,
                    icon: category.categoryIconId,
                    color: category.color,
                    subCategories: []
                };

                for (const subCategory of category.subCategories) {
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

            allCategories[`${categoryType}`] = categories;
        }

        return allCategories;
    }

    function getAllDisplayExchangeRates(exchangeRatesData?: LatestExchangeRateResponse): LocalizedLatestExchangeRate[] {
        const availableExchangeRates: LocalizedLatestExchangeRate[] = [];

        if (!exchangeRatesData || !exchangeRatesData.exchangeRates) {
            return availableExchangeRates;
        }

        for (const exchangeRate of exchangeRatesData.exchangeRates) {
            availableExchangeRates.push({
                currencyCode: exchangeRate.currency,
                currencyDisplayName: getCurrencyName(exchangeRate.currency),
                rate: exchangeRate.rate
            });
        }

        if (settingsStore.appSettings.currencySortByInExchangeRatesPage === CurrencySortingType.Name.type) {
            availableExchangeRates.sort(function (c1, c2) {
                return c1.currencyDisplayName.localeCompare(c2.currencyDisplayName);
            });
        } else if (settingsStore.appSettings.currencySortByInExchangeRatesPage === CurrencySortingType.CurrencyCode.type) {
            availableExchangeRates.sort(function (c1, c2) {
                return c1.currencyCode.localeCompare(c2.currencyCode);
            });
        } else if (settingsStore.appSettings.currencySortByInExchangeRatesPage === CurrencySortingType.ExchangeRate.type) {
            availableExchangeRates.sort(function (c1, c2) {
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

    function getAllSupportedImportFileCagtegoryAndTypes(): LocalizedImportFileCategoryAndTypes[] {
        const allSupportedImportFileCategoryAndTypes: LocalizedImportFileCategoryAndTypes[] = [];

        for (const categoryAndTypes of SUPPORTED_IMPORT_FILE_CATEGORY_AND_TYPES) {
            const localizedCategoryAndTypes: LocalizedImportFileCategoryAndTypes = {
                displayCategoryName: t(categoryAndTypes.categoryName),
                fileTypes: []
            };

            for (const fileType of categoryAndTypes.fileTypes) {
                let document: LocalizedImportFileDocument | undefined;

                if (fileType.document) {
                    let documentLanguage = '';
                    let documentDisplayLanguageName = '';
                    let documentAnchor = '';

                    if (fileType.document.supportMultiLanguages === true) {
                        documentLanguage = getCurrentLanguageTag();

                        if (SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE[documentLanguage] === documentLanguage) {
                            documentAnchor = t(`document.anchor.export_and_import.${fileType.document.anchor}`);
                        } else if (SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE[documentLanguage]) {
                            documentLanguage = SUPPORTED_DOCUMENT_LANGUAGES_FOR_IMPORT_FILE[documentLanguage] as string;
                            documentAnchor = t(`document.anchor.export_and_import.${fileType.document.anchor}`, {}, { locale: documentLanguage });
                        } else {
                            documentLanguage = DEFAULT_DOCUMENT_LANGUAGE_FOR_IMPORT_FILE;
                            documentAnchor = t(`document.anchor.export_and_import.${fileType.document.anchor}`, {}, { locale: documentLanguage });
                        }
                    } else if (isString(fileType.document.supportMultiLanguages) && ALL_LANGUAGES[fileType.document.supportMultiLanguages]) {
                        documentLanguage = fileType.document.supportMultiLanguages;
                        documentAnchor = fileType.document.anchor;
                    }

                    if (documentLanguage && documentLanguage !== getCurrentLanguageTag() && ALL_LANGUAGES[documentLanguage]) {
                        documentDisplayLanguageName = getLanguageDisplayName((ALL_LANGUAGES[documentLanguage] as LanguageInfo).name);
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
                    for (const subType of fileType.subTypes) {
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
                    for (const encoding of fileType.supportedEncodings) {
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
                    supportedAdditionalOptions: fileType.supportedAdditionalOptions,
                    document: document
                };

                localizedCategoryAndTypes.fileTypes.push(localizedFileType);
            }

            allSupportedImportFileCategoryAndTypes.push(localizedCategoryAndTypes);
        }

        return allSupportedImportFileCategoryAndTypes;
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

    function getWeekdayShortName(weekDay: WeekDay): string {
        return t(`datetime.${weekDay.name}.short`);
    }

    function getWeekdayLongName(weekDay: WeekDay): string {
        return t(`datetime.${weekDay.name}.long`);
    }

    function getMultiMonthdayShortNames(monthDays: number[]): string {
        if (!monthDays) {
            return '';
        }

        if (monthDays.length === 1) {
            return t('format.misc.monthDay', {
                ordinal: getMonthdayOrdinal(monthDays[0] as number)
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

    function getMultiWeekdayLongNames(weekdayTypes: number[], firstDayOfWeek?: WeekDayValue): string {
        const weekdayTypesMap: Record<number, boolean> = {};

        if (!isNumber(firstDayOfWeek)) {
            firstDayOfWeek = WeekDay.DefaultFirstDay.type;
        }

        for (const weekdayType of weekdayTypes) {
            weekdayTypesMap[weekdayType] = true;
        }

        const allWeekDays = getAllWeekDays(firstDayOfWeek);
        const finalWeekdayNames = [];

        for (const weekDay of allWeekDays) {
            if (weekdayTypesMap[weekDay.type]) {
                finalWeekdayNames.push(weekDay.displayName);
            }
        }

        return joinMultiText(finalWeekdayNames);
    }

    function getAllLocalizedDigits(): string[] {
        const numeralSystem = getCurrentNumeralSystemType();
        return numeralSystem.getAllDigits();
    }

    function getLocaleDefaultCalendarDisplayType(): CalendarDisplayType {
        const defaultCalendarDisplayTypeName = t('default.calendarDisplayType');
        let calendarDisplayType = CalendarDisplayType.parse(defaultCalendarDisplayTypeName);

        if (!calendarDisplayType) {
            calendarDisplayType = CalendarDisplayType.Default;
        }

        return calendarDisplayType;
    }

    function getLocaleDefaultDateDisplayType(): DateDisplayType {
        const defaultDateDisplayTypeName = t('default.dateDisplayType');
        let dateDisplayType = DateDisplayType.parse(defaultDateDisplayTypeName);

        if (!dateDisplayType) {
            dateDisplayType = DateDisplayType.Default;
        }

        return dateDisplayType;
    }

    function getLocaleDefaultNumeralSystemType(): NumeralSystem {
        const defaultNumeralSystemTypeName = t('default.numeralSystem');
        let numeralSystemType = NumeralSystem.parse(defaultNumeralSystemTypeName);

        if (!numeralSystemType) {
            numeralSystemType = NumeralSystem.Default;
        }

        return numeralSystemType;
    }

    function getLocaleDefaultDecimalSeparator(): DecimalSeparator {
        const defaultDecimalSeparatorTypeName = t('default.decimalSeparator');
        let decimalSeparatorType = DecimalSeparator.parse(defaultDecimalSeparatorTypeName);

        if (!decimalSeparatorType) {
            decimalSeparatorType = DecimalSeparator.Default;
        }

        return decimalSeparatorType;
    }

    function getLocaleDefaultDigitGroupingSymbol(): DigitGroupingSymbol {
        const defaultDigitGroupingSymbolTypeName = t('default.digitGroupingSymbol');
        let digitGroupingSymbolType = DigitGroupingSymbol.parse(defaultDigitGroupingSymbolTypeName);

        if (!digitGroupingSymbolType) {
            digitGroupingSymbolType = DigitGroupingSymbol.Default;
        }

        return digitGroupingSymbolType;
    }

    function getLocaleDefaultDigitGroupingType(): DigitGroupingType {
        const defaultDigitGroupingTypeName = t('default.digitGrouping');
        let digitGroupingType = DigitGroupingType.parse(defaultDigitGroupingTypeName);

        if (!digitGroupingType) {
            digitGroupingType = DigitGroupingType.Default;
        }

        return digitGroupingType;
    }

    function getLocaleFiscalYearFormatType(): FiscalYearFormat {
        const defaultFiscalYearFormatTypeName = t('default.fiscalYearFormat');
        let fiscalYearFormat = FiscalYearFormat.parse(defaultFiscalYearFormatTypeName);

        if (!fiscalYearFormat) {
            fiscalYearFormat = FiscalYearFormat.Default;
        }

        return fiscalYearFormat;
    }

    function getCurrentCalendarDisplayType(): CalendarDisplayType {
        return CalendarDisplayType.valueOf(userStore.currentUserCalendarDisplayType) ?? getLocaleDefaultCalendarDisplayType();
    }

    function getCurrentDateDisplayType(): DateDisplayType {
        return DateDisplayType.valueOf(userStore.currentUserDateDisplayType) ?? getLocaleDefaultDateDisplayType();
    }

    function getCurrentNumeralSystemType(): NumeralSystem {
        return NumeralSystem.valueOf(userStore.currentUserNumeralSystem) ?? getLocaleDefaultNumeralSystemType();
    }

    function getCurrentDecimalSeparator(): string {
        return DecimalSeparator.valueOf(userStore.currentUserDecimalSeparator)?.symbol ?? getLocaleDefaultDecimalSeparator().symbol;
    }

    function getCurrentDigitGroupingSymbol(): string {
        return DigitGroupingSymbol.valueOf(userStore.currentUserDigitGroupingSymbol)?.symbol ?? getLocaleDefaultDigitGroupingSymbol().symbol;
    }

    function getCurrentDigitGroupingType(): DigitGroupingType {
        return DigitGroupingType.valueOf(userStore.currentUserDigitGrouping) ?? getLocaleDefaultDigitGroupingType();
    }

    function getCurrentFiscalYearFormatType(): number {
        return FiscalYearFormat.valueOf(userStore.currentUserFiscalYearFormat)?.type ?? getLocaleFiscalYearFormatType().type;
    }

    function getCurrencyName(currencyCode: string): string {
        if (!currencyCode) {
            return '';
        }

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

    function isLongTimeHourTwoDigits(): boolean {
        const longTimeFormat = getLocalizedLongTimeFormat();
        return longTimeFormat.indexOf('HH') >= 0 || longTimeFormat.indexOf('hh') >= 0;
    }

    function isLongTimeMinuteTwoDigits(): boolean {
        return getLocalizedLongTimeFormat().indexOf('mm') >= 0;
    }

    function isLongTimeSecondTwoDigits(): boolean {
        return getLocalizedLongTimeFormat().indexOf('ss') >= 0;
    }

    function formatGregorianTextualMonthDayToGregorianLikeLongMonthDay(monthDay: TextualMonthDay, numeralSystem?: NumeralSystem): string {
        const gregorianLikeCalendarType = getGregorianLikeCalendarType();
        return formatGregorianCalendarMonthDashDay(monthDay, getLocalizedLongMonthDayFormat(), getDateTimeFormatOptions({ calendarType: gregorianLikeCalendarType, numeralSystem: numeralSystem }));
    }

    function formatDateTimeToGregorianLikeYearQuarter(dateTime: DateTime): string {
        const gregorianLikeCalendarType = getGregorianLikeCalendarType();
        const dateTimeFormatOptions = getDateTimeFormatOptions({ calendarType: gregorianLikeCalendarType });
        const year = dateTime.getLocalizedCalendarYear(dateTimeFormatOptions);
        const quarter = dateTime.getLocalizedCalendarQuarter(dateTimeFormatOptions);
        return formatYearQuarter(year, quarter);
    }

    function formatYearQuarterToGregorianLikeYearQuarter(year: number, quarter: number): string {
        const gregorianLikeCalendarType = getGregorianLikeCalendarType();
        const dateTimeFormatOptions = getDateTimeFormatOptions({ calendarType: gregorianLikeCalendarType });
        const date = getYearMonthDayDateTime(year, 1, 1);
        const textualYear = date.getLocalizedCalendarYear(dateTimeFormatOptions);
        return formatYearQuarter(textualYear, quarter);
    }

    function formatUnixTimeToGregorianLikeFiscalYear(unixTime: number): string {
        let fiscalYearFormat = FiscalYearFormat.valueOf(getCurrentFiscalYearFormatType());

        if (!fiscalYearFormat) {
            fiscalYearFormat = FiscalYearFormat.Default;
        }

        const timeRange = getFiscalYearTimeRangeFromUnixTime(unixTime, userStore.currentUserFiscalYearStart);
        return formatTimeRangeToGregorianLikeFiscalYearFormat(fiscalYearFormat, timeRange);
    }

    function formatGregorianYearToGregorianLikeFiscalYear(year: number) {
        let fiscalYearFormat = FiscalYearFormat.valueOf(getCurrentFiscalYearFormatType());

        if (!fiscalYearFormat) {
            fiscalYearFormat = FiscalYearFormat.Default;
        }

        const timeRange = getFiscalYearTimeRangeFromYear(year, userStore.currentUserFiscalYearStart);
        return formatTimeRangeToGregorianLikeFiscalYearFormat(fiscalYearFormat, timeRange);
    }

    function formatFiscalYearStartToGregorianLikeLongMonth(fiscalYearStartValue: number, numeralSystem?: NumeralSystem): string {
        let fiscalYearStart = FiscalYearStart.valueOf(fiscalYearStartValue);

        if (!fiscalYearStart) {
            fiscalYearStart = FiscalYearStart.Default;
        }

        return formatGregorianTextualMonthDayToGregorianLikeLongMonthDay(fiscalYearStart.toMonthDashDayString(), numeralSystem);
    }

    function formatDateRange(dateType: number, startTime: number, endTime: number): string {
        if (dateType === DateRange.All.type) {
            return t(DateRange.All.name);
        }

        const allDateRanges = DateRange.values();
        const gregorianLikeCalendarType = getGregorianLikeCalendarType();
        const dateTimeFormatOptions = getDateTimeFormatOptions();
        const gregorianLikeDateTimeFormatOptions = getDateTimeFormatOptions({ calendarType: gregorianLikeCalendarType });

        for (const dateRange of allDateRanges) {
            if (dateRange && dateRange.type !== DateRange.Custom.type && dateRange.type === dateType && dateRange.name) {
                return t(dateRange.name);
            }
        }

        if (isDateRangeMatchFullYears(startTime, endTime)) {
            const format = getLocalizedShortYearFormat();
            const displayStartTime = formatUnixTime(startTime, format, gregorianLikeDateTimeFormatOptions);
            const displayEndTime = formatUnixTime(endTime, format, gregorianLikeDateTimeFormatOptions);

            return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
        }

        if (isDateRangeMatchFullMonths(startTime, endTime)) {
            const format = getLocalizedShortYearMonthFormat();
            const displayStartTime = formatUnixTime(startTime, format, gregorianLikeDateTimeFormatOptions);
            const displayEndTime = formatUnixTime(endTime, format, gregorianLikeDateTimeFormatOptions);

            return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
        }

        const startTimeYear = parseDateTimeFromUnixTime(startTime).getLocalizedCalendarYear(gregorianLikeDateTimeFormatOptions);
        const endTimeYear = parseDateTimeFromUnixTime(endTime).getLocalizedCalendarYear(gregorianLikeDateTimeFormatOptions);

        const format = getLocalizedShortDateFormat();
        const displayStartTime = formatUnixTime(startTime, format, dateTimeFormatOptions);
        const displayEndTime = formatUnixTime(endTime, format, dateTimeFormatOptions);

        if (displayStartTime === displayEndTime) {
            return displayStartTime;
        } else if (startTimeYear === endTimeYear) {
            const displayShortEndTime = formatUnixTime(endTime, getLocalizedShortMonthDayFormat(), gregorianLikeDateTimeFormatOptions);
            return `${displayStartTime} ~ ${displayShortEndTime}`;
        }

        return `${displayStartTime} ~ ${displayEndTime}`;
    }

    function getTimezoneDifferenceDisplayText(unixTime: number, utcOffset: number): string {
        const numeralSystem = getCurrentNumeralSystemType();
        const defaultTimezoneOffset = getTimezoneOffsetMinutes(unixTime);
        const offsetTime = getTimeDifferenceHoursAndMinutes(utcOffset - defaultTimezoneOffset);

        if (utcOffset > defaultTimezoneOffset) {
            if (offsetTime.offsetMinutes) {
                return t('format.misc.hoursMinutesAheadOfDefaultTimezone', {
                    hours: numeralSystem.formatNumber(offsetTime.offsetHours),
                    minutes: numeralSystem.formatNumber(offsetTime.offsetMinutes)
                });
            } else {
                return t('format.misc.hoursAheadOfDefaultTimezone', {
                    hours: numeralSystem.formatNumber(offsetTime.offsetHours)
                });
            }
        } else if (utcOffset < defaultTimezoneOffset) {
            if (offsetTime.offsetMinutes) {
                return t('format.misc.hoursMinutesBehindDefaultTimezone', {
                    hours: numeralSystem.formatNumber(offsetTime.offsetHours),
                    minutes: numeralSystem.formatNumber(offsetTime.offsetMinutes)
                });
            } else {
                return t('format.misc.hoursBehindDefaultTimezone', {
                    hours: numeralSystem.formatNumber(offsetTime.offsetHours)
                });
            }
        } else {
            return t('Same time as default timezone');
        }
    }

    function getCalendarAlternateDates(yearMonth: Year1BasedMonth): CalendarAlternateDate[] | undefined {
        const calendarDisplayType = getCurrentCalendarDisplayType().secondaryCalendarType;

        if (!calendarDisplayType) {
            return undefined;
        }

        const dateTimeFormatOptions = getDateTimeFormatOptions({ calendarType: calendarDisplayType });

        if (calendarDisplayType === CalendarType.Chinese) {
            const chineseCalendarLocaleData = getChineseCalendarLocaleData();
            const chineseDates: ChineseYearMonthDayInfo[] | undefined = getChineseYearMonthAllDayInfos(yearMonth, chineseCalendarLocaleData);

            if (!chineseDates) {
                return undefined;
            }

            const ret: CalendarAlternateDate[] = [];

            for (const chineseDate of chineseDates) {
                const alternateDate = getChineseCalendarAlternateDisplayDate(chineseDate);
                ret.push(alternateDate);
            }

            return ret;
        } else if (calendarDisplayType === CalendarType.Persian) {
            const monthDays: number = getGregorianCalendarYearMonthDays(yearMonth);
            const ret: CalendarAlternateDate[] = [];

            for (let i = 1; i <= monthDays; i++) {
                const dateTime = getYearMonthDayDateTime(yearMonth.year, yearMonth.month1base, i);
                ret.push(getCalendarAlternateDisplayDate(dateTime, dateTimeFormatOptions));
            }

            return ret;
        }

        return undefined;
    }

    function getCalendarAlternateDate(yearMonthDay: YearMonthDay): CalendarAlternateDate | undefined {
        const calendarDisplayType = getCurrentCalendarDisplayType().secondaryCalendarType;

        if (!calendarDisplayType) {
            return undefined;
        }

        const dateTimeFormatOptions = getDateTimeFormatOptions({ calendarType: calendarDisplayType });

        if (calendarDisplayType === CalendarType.Chinese) {
            const chineseCalendarLocaleData = getChineseCalendarLocaleData();
            const chineseDate: ChineseYearMonthDayInfo | undefined = getChineseYearMonthDayInfo(yearMonthDay, chineseCalendarLocaleData);

            if (!chineseDate) {
                return undefined;
            }

            return getChineseCalendarAlternateDisplayDate(chineseDate);
        } else if (calendarDisplayType === CalendarType.Persian) {
            const dateTime = getYearMonthDayDateTime(yearMonthDay.year, yearMonthDay.month, yearMonthDay.day);
            return getCalendarAlternateDisplayDate(dateTime, dateTimeFormatOptions);
        }

        return undefined;
    }

    function getParsedAmountNumber(value: string, numeralSystem?: NumeralSystem): number {
        const numberFormatOptions = getNumberFormatOptions({ numeralSystem });
        return parseAmount(value, numberFormatOptions);
    }

    function getFormattedAmount(value: number, numeralSystem?: NumeralSystem, digitGrouping?: DigitGroupingType, currencyCode?: string): string {
        const numberFormatOptions = getNumberFormatOptions({ numeralSystem, digitGrouping, currencyCode });
        return formatAmount(value, numberFormatOptions);
    }

    function getFormattedAmountWithCurrency(value: number | HiddenAmount | NumberWithSuffix, currencyCode?: string | false, currencyDisplayType?: CurrencyDisplayType, numeralSystem?: NumeralSystem, decimalSeparator?: string): string {
        let finalCurrencyCode = '';

        if (!isBoolean(currencyCode) && !currencyCode) {
            finalCurrencyCode = userStore.currentUserDefaultCurrency;
        } else if (isBoolean(currencyCode) && !currencyCode) {
            finalCurrencyCode = '';
        } else {
            finalCurrencyCode = currencyCode;
        }

        if (!currencyDisplayType) {
            currencyDisplayType = getCurrentCurrencyDisplayType();
        }

        if (!numeralSystem) {
            numeralSystem = getCurrentNumeralSystemType();
        }

        let suffix = '';

        if (isObject(value) && isNumber(value.value) && isString(value.suffix)) {
            suffix = value.suffix;
            value = value.value;
        }

        const numberFormatOptions = getNumberFormatOptions({ numeralSystem, decimalSeparator, currencyCode: finalCurrencyCode });
        const currencyName = getCurrencyName(finalCurrencyCode);

        if (isNumber(value)) {
            const isPlural: boolean = value !== 100 && value !== -100;
            const textualValue = formatAmount(value, numberFormatOptions);

            if (!finalCurrencyCode) {
                return textualValue;
            }

            const currencyUnit = getCurrencyUnitName(finalCurrencyCode, isPlural);
            const ret = appendCurrencySymbol(textualValue, currencyDisplayType, finalCurrencyCode, currencyUnit, currencyName, isPlural);

            if (suffix) {
                return ret + suffix;
            } else {
                return ret;
            }
        } else if (isString(value)) {
            const isPlural: boolean = true;
            const textualValue = formatHiddenAmount(value, numberFormatOptions);

            if (!finalCurrencyCode) {
                return textualValue;
            }

            const currencyUnit = getCurrencyUnitName(finalCurrencyCode, isPlural);
            return appendCurrencySymbol(textualValue, currencyDisplayType, finalCurrencyCode, currencyUnit, currencyName, isPlural);
        } else {
            return '';
        }
    }

    function getFormattedNumber(value: number, numeralSystem?: NumeralSystem, precision?: number): string {
        const numberFormatOptions = getNumberFormatOptions({ numeralSystem, digitGrouping: DigitGroupingType.None });
        return formatNumber(value, numberFormatOptions, precision);
    }

    function getFormattedPercentValue(value: number, precision: number, lowPrecisionValue: string, numeralSystem?: NumeralSystem): string {
        const numberFormatOptions = getNumberFormatOptions({ numeralSystem });
        return formatPercent(value, precision, lowPrecisionValue, numberFormatOptions);
    }

    function getFormattedExchangeRateAmount(value: number, numeralSystem?: NumeralSystem): string {
        const numberFormatOptions = getNumberFormatOptions({ numeralSystem });
        return formatExchangeRateAmount(value, numberFormatOptions);
    }

    function getAdaptiveAmountRate(amount1: number, amount2: number, fromExchangeRate?: { rate: string }, toExchangeRate?: { rate: string }): string | null {
        const numberFormatOptions = getNumberFormatOptions({});
        return getAdaptiveDisplayAmountRate(amount1, amount2, numberFormatOptions, fromExchangeRate, toExchangeRate);
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

        for (const category of allCategories) {
            const accountCategory = categorizedAccounts[category.type];

            if (!accountCategory) {
                continue;
            }

            const accountsWithDisplayBalance: AccountWithDisplayBalance[] = [];

            if (accountCategory.accounts) {
                for (const account of accountCategory.accounts) {
                    let accountWithDisplaceBalance: AccountWithDisplayBalance;

                    if (showAccountBalance && account.isAsset) {
                        accountWithDisplaceBalance = AccountWithDisplayBalance.fromAccount(account, getFormattedAmountWithCurrency(account.balance, account.currency));
                    } else if (showAccountBalance && account.isLiability) {
                        accountWithDisplaceBalance = AccountWithDisplayBalance.fromAccount(account, getFormattedAmountWithCurrency(-account.balance, account.currency));
                    } else {
                        accountWithDisplaceBalance = AccountWithDisplayBalance.fromAccount(account, DISPLAY_HIDDEN_AMOUNT);
                    }

                    accountsWithDisplayBalance.push(accountWithDisplaceBalance);
                }
            }

            let finalTotalBalance = '';

            if (showAccountBalance) {
                const accountsBalance = getAllFilteredAccountsBalance(categorizedAccounts, account => account.category === accountCategory.category);
                let totalBalance = 0;
                let hasUnCalculatedAmount = false;

                for (const accountBalance of accountsBalance) {
                    if (accountBalance.currency === defaultCurrency) {
                        if (accountBalance.isAsset) {
                            totalBalance += accountBalance.balance;
                        } else if (accountBalance.isLiability) {
                            totalBalance -= accountBalance.balance;
                        }
                    } else {
                        const balance = exchangeRatesStore.getExchangedAmount(accountBalance.balance, accountBalance.currency, defaultCurrency);

                        if (!isNumber(balance)) {
                            hasUnCalculatedAmount = true;
                            continue;
                        }

                        if (accountBalance.isAsset) {
                            totalBalance += Math.trunc(balance);
                        } else if (accountBalance.isLiability) {
                            totalBalance -= Math.trunc(balance);
                        }
                    }
                }

                finalTotalBalance = getFormattedAmountWithCurrency(totalBalance, defaultCurrency);

                if (hasUnCalculatedAmount) {
                    finalTotalBalance = finalTotalBalance + INCOMPLETE_AMOUNT_SUFFIX;
                }
            } else {
                finalTotalBalance = DISPLAY_HIDDEN_AMOUNT;
            }

            const accountCategoryWithDisplayBalance = CategorizedAccountWithDisplayBalance.of(accountCategory, accountsWithDisplayBalance, finalTotalBalance);
            ret.push(accountCategoryWithDisplayBalance);
        }

        return ret;
    }

    function getLocalizedFileEncodingName(encoding: string): string {
        return t(`encoding.${encoding}`);
    }

    function getLocalizedOAuth2ProviderName(oauth2Provider: string, oidcDisplayNames: Record<string, string>): string {
        if (oauth2Provider === 'oidc') {
            const providerDisplayName = getServerMultiLanguageConfigContent(oidcDisplayNames);

            if (providerDisplayName) {
                return providerDisplayName;
            } else {
                return 'Connect ID';
            }
        } else {
            const providerDisplayName = OAUTH2_PROVIDER_DISPLAY_NAME[oauth2Provider];

            if (providerDisplayName) {
                return providerDisplayName;
            } else {
                return 'OAuth 2.0';
            }
        }
    }

    function getLocalizedOAuth2LoginText(oauth2Provider: string, oidcDisplayNames: Record<string, string>): string {
        if (oauth2Provider === 'oidc') {
            const providerDisplayName = getServerMultiLanguageConfigContent(oidcDisplayNames);

            if (providerDisplayName) {
                return t('format.misc.loginWithCustomProvider', { name: providerDisplayName });
            } else {
                return t('Log in with Connect ID');
            }
        } else {
            const providerDisplayName = OAUTH2_PROVIDER_DISPLAY_NAME[oauth2Provider];

            if (providerDisplayName) {
                return t('format.misc.loginWithCustomProvider', { name: providerDisplayName });
            } else {
                return t('Log in with OAuth 2.0');
            }
        }
    }

    function setLanguage(languageKey: string | null, force?: boolean): LocaleDefaultSettings | null {
        if (!languageKey) {
            languageKey = getDefaultLanguage();
            logger.info(`No specified language, use browser default language ${languageKey}`);
        }

        const languageInfo = getLanguageInfo(languageKey);

        if (!languageInfo) {
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
            months: getAllLongMonthNames(),
            monthsShort: getAllShortMonthNames(),
            weekdays: getAllLongWeekdayNames(),
            weekdaysShort: getAllShortWeekdayNames(),
            weekdaysMin: getAllMinWeekdayNames(),
            meridiem: function (hours) {
                if (isPM(hours)) {
                    return t(`datetime.${MeridiemIndicator.PM.name}.content`);
                } else {
                    return t(`datetime.${MeridiemIndicator.AM.name}.content`);
                }
            }
        });

        setSessionCurrentLanguageKey(languageKey);
        services.setLocale(languageKey);
        document.querySelector('html')?.setAttribute('lang', languageKey);

        if (document.querySelector('html')?.getAttribute('data-dir-mode') === 'static') {
            if (languageInfo && languageInfo.textDirection === TextDirection.LTR) {
                if (location.search.includes('rtl')) {
                    const url = new URL(window.location.href);
                    url.search = '';
                    url.hash = '#/';
                    window.location.replace(url.toString());
                }
            } else if (languageInfo && languageInfo.textDirection === TextDirection.RTL) {
                if (!location.search.includes('rtl')) {
                    const url = new URL(window.location.href);
                    url.searchParams.set('rtl', '');
                    url.hash = '#/';
                    window.location.replace(url.toString());
                }
            }
        } else {
            if (languageInfo && languageInfo.textDirection === TextDirection.LTR) {
                document.querySelector('html')?.removeAttribute('dir');
            } else if (languageInfo && languageInfo.textDirection === TextDirection.RTL) {
                document.querySelector('html')?.setAttribute('dir', 'rtl');
            }
        }

        const defaultCurrency = getDefaultCurrency();
        const defaultFirstDayOfWeekName = getDefaultFirstDayOfWeek();
        let defaultFirstDayOfWeek = WeekDay.DefaultFirstDay.type;

        if (WeekDay.parse(defaultFirstDayOfWeekName)) {
            defaultFirstDayOfWeek = (WeekDay.parse(defaultFirstDayOfWeekName) as WeekDay).type;
        }

        const localeDefaultSettings: LocaleDefaultSettings = {
            currency: defaultCurrency,
            firstDayOfWeek: defaultFirstDayOfWeek
        };

        return localeDefaultSettings;
    }

    function setTimeZone(timezone: string): void {
        if (timezone) {
            moment.tz.setDefault(timezone);
        } else {
            moment.tz.setDefault();
        }
    }

    function initLocale(lastUserLanguage?: string, timezone?: string): LocaleDefaultSettings | null {
        const sessionLanguageKey: string = getSessionCurrentLanguageKey();
        let localeDefaultSettings: LocaleDefaultSettings | null = null;

        if (lastUserLanguage && getLanguageInfo(lastUserLanguage)) {
            logger.info(`Last user language is ${lastUserLanguage}`);
            localeDefaultSettings = setLanguage(lastUserLanguage, true);
        } else if (sessionLanguageKey && getLanguageInfo(sessionLanguageKey)) {
            logger.info(`Session language is ${sessionLanguageKey}`);
            localeDefaultSettings = setLanguage(sessionLanguageKey, true);
        } else {
            localeDefaultSettings = setLanguage(null, true);
        }

        if (timezone) {
            logger.info(`Current timezone is ${timezone}`);
            setTimeZone(timezone);
        } else {
            logger.info(`No timezone is set, use browser default ${getTimezoneOffset(getCurrentUnixTime())} (${getBrowserTimezoneName()})`);
            setTimeZone('');
        }

        return localeDefaultSettings;
    }

    return {
        // common functions
        tt: t,
        ti: translateIf,
        te: translateError,
        joinMultiText,
        getServerMultiLanguageConfigContent,
        // get current language info
        getCurrentLanguageTag,
        getCurrentLanguageInfo,
        getCurrentLanguageDisplayName,
        getCurrentLanguageTextDirection,
        // get localization default type
        getDefaultCurrency,
        getDefaultFirstDayOfWeek,
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
        getAllCalendarDisplayTypes: () => getAllLocalizedCalendarTypes(CalendarDisplayType.values(), CalendarDisplayType.parse(t('default.calendarDisplayType')), CalendarDisplayType.Default, CalendarDisplayType.LanguageDefaultType),
        getAllDateDisplayTypes: () => getAllLocalizedCalendarTypes(DateDisplayType.values(), DateDisplayType.parse(t('default.dateDisplayType')), DateDisplayType.Default, DateDisplayType.LanguageDefaultType),
        getAllLongDateFormats: (numeralSystem: NumeralSystem, calendarType: CalendarType) => getLocalizedDateTimeFormats<LongDateFormat>('longDate', LongDateFormat.all(), LongDateFormat.values(), 'longDateFormat', LongDateFormat.Default, numeralSystem, calendarType),
        getAllShortDateFormats: (numeralSystem: NumeralSystem, calendarType: CalendarType) => getLocalizedDateTimeFormats<ShortDateFormat>('shortDate', ShortDateFormat.all(), ShortDateFormat.values(), 'shortDateFormat', ShortDateFormat.Default, numeralSystem, calendarType),
        getAllLongTimeFormats: (numeralSystem: NumeralSystem) => getLocalizedDateTimeFormats<LongTimeFormat>('longTime', LongTimeFormat.all(), LongTimeFormat.values(), 'longTimeFormat', LongTimeFormat.Default, numeralSystem),
        getAllShortTimeFormats: (numeralSystem: NumeralSystem) => getLocalizedDateTimeFormats<ShortTimeFormat>('shortTime', ShortTimeFormat.all(), ShortTimeFormat.values(), 'shortTimeFormat', ShortTimeFormat.Default, numeralSystem),
        getAllFiscalYearFormats,
        getAllDateRanges,
        getAllRecentMonthDateRanges,
        getAllTimezones,
        getAllTimezoneTypesUsedForStatistics,
        getAllNumeralSystemTypes,
        getAllDecimalSeparators: () => getLocalizedNumeralSeparatorFormats(DecimalSeparator.values(), DecimalSeparator.parse(t('default.decimalSeparator')), DecimalSeparator.Default, DecimalSeparator.LanguageDefaultType),
        getAllDigitGroupingSymbols: () => getLocalizedNumeralSeparatorFormats(DigitGroupingSymbol.values(), DigitGroupingSymbol.parse(t('default.digitGroupingSymbol')), DigitGroupingSymbol.Default, DigitGroupingSymbol.LanguageDefaultType),
        getAllDigitGroupingTypes,
        getAllCurrencyDisplayTypes,
        getAllCurrencySortingTypes: () => getLocalizedDisplayNameAndType(CurrencySortingType.values()),
        getAllCoordinateDisplayTypes: () => getLocalizedDisplayNameAndTypeWithSystemDefault(CoordinateDisplayType.values(), CoordinateDisplayType.SystemDefaultType, CoordinateDisplayType.Default),
        getAllExpenseAmountColors: () => getAllExpenseIncomeAmountColors(CategoryType.Expense),
        getAllIncomeAmountColors: () => getAllExpenseIncomeAmountColors(CategoryType.Income),
        getAllAccountCategories,
        getAllAccountTypes: () => getLocalizedDisplayNameAndType(AccountType.values()),
        getAllCategoricalChartTypes: (withDesktopOnlyChart?: boolean) => getLocalizedDisplayNameAndType(CategoricalChartType.values(!!withDesktopOnlyChart)),
        getAllTrendChartTypes: () => getLocalizedDisplayNameAndType(TrendChartType.values()),
        getAllAccountBalanceTrendChartTypes: () => getLocalizedDisplayNameAndType(AccountBalanceTrendChartType.values()),
        getAllStatisticsChartDataTypes: (analysisType: StatisticsAnalysisType, withDesktopOnlyChart?: boolean) => getLocalizedDisplayNameAndType(ChartDataType.values(analysisType, withDesktopOnlyChart)),
        getAllStatisticsSortingTypes: () => getLocalizedDisplayNameAndType(ChartSortingType.values()),
        getAllStatisticsDateAggregationTypes: (analysisType: StatisticsAnalysisType) => getLocalizedChartDateAggregationTypeAndDisplayName(analysisType, true),
        getAllStatisticsDateAggregationTypesWithShortName: (analysisType: StatisticsAnalysisType) => getLocalizedChartDateAggregationTypeAndDisplayName(analysisType, false),
        getAllTransactionEditScopeTypes: () => getLocalizedDisplayNameAndType(TransactionEditScopeType.values()),
        getAllTransactionScheduledFrequencyTypes: () => getLocalizedDisplayNameAndType(ScheduledTemplateFrequencyType.values()),
        getAllImportTransactionColumnTypes: () => getLocalizedDisplayNameAndType(ImportTransactionColumnType.values()),
        getAllTransactionDefaultCategories,
        getAllDisplayExchangeRates,
        getAllSupportedImportFileCagtegoryAndTypes,
        getAllTransactionExploreConditionFields: () => getLocalizedNameValue(TransactionExploreConditionField.values()),
        getAllTransactionExploreConditionOperators: (operators?: TransactionExploreConditionOperator[]) => getLocalizedNameValue(operators ?? TransactionExploreConditionOperator.values()),
        getAllTransactionExploreDataDimensions: (operators?: TransactionExploreDataDimension[]) => getLocalizedNameValue(operators ?? TransactionExploreDataDimension.values()),
        getAllTransactionExploreValueMetrics: (operators?: TransactionExploreValueMetric[]) => getLocalizedNameValue(operators ?? TransactionExploreValueMetric.values()),
        getAllTransactionExploreChartTypes: (operators?: TransactionExploreChartType[]) => getLocalizedNameValue(operators ?? TransactionExploreChartType.values()),
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
        getAllLocalizedDigits,
        getLocaleDefaultCalendarDisplayType,
        getLocaleDefaultDateDisplayType,
        getLocaleDefaultNumeralSystemType,
        getLocaleDefaultDecimalSeparator,
        getLocaleDefaultDigitGroupingSymbol,
        getLocaleDefaultDigitGroupingType,
        getLocaleFiscalYearFormatType,
        getCurrentCalendarDisplayType,
        getCurrentDateDisplayType,
        getCurrentNumeralSystemType,
        getCurrentDecimalSeparator,
        getCurrentDigitGroupingSymbol,
        getCurrentDigitGroupingType,
        getCurrentFiscalYearFormatType,
        getCurrencyName,
        isLongDateMonthAfterYear,
        isShortDateMonthAfterYear,
        isLongTime24HourFormat,
        isLongTimeMeridiemIndicatorFirst,
        isShortTime24HourFormat,
        isShortTimeMeridiemIndicatorFirst,
        isLongTimeHourTwoDigits,
        isLongTimeMinuteTwoDigits,
        isLongTimeSecondTwoDigits,
        // format date time (by calendar display type) functions
        getCalendarDisplayShortYearFromDateTime: (dateTime: DateTime, numeralSystem?: NumeralSystem) => formatDateTime(dateTime, getLocalizedShortYearFormat(), getDateTimeFormatOptions({ calendarType: getCurrentCalendarDisplayType().primaryCalendarType, numeralSystem: numeralSystem })),
        getCalendarDisplayShortMonthFromDateTime: (dateTime: DateTime, numeralSystem?: NumeralSystem) => formatDateTime(dateTime, 'MMM', getDateTimeFormatOptions({ calendarType: getCurrentCalendarDisplayType().primaryCalendarType, numeralSystem: numeralSystem })),
        getCalendarDisplayDayOfMonthFromDateTime: (dateTime: DateTime, numeralSystem?: NumeralSystem) => formatDateTime(dateTime, getLocalizedShortDayFormat(), getDateTimeFormatOptions({ calendarType: getCurrentCalendarDisplayType().primaryCalendarType, numeralSystem: numeralSystem })),
        // format date time (by date display type) functions
        parseDateTimeFromLongDateTime: (dateTime: string) => parseDateTimeFromString(dateTime, getLocalizedLongDateFormat() + ' ' + getLocalizedLongTimeFormat()),
        parseDateTimeFromShortDateTime: (dateTime: string) => parseDateTimeFromString(dateTime, getLocalizedShortDateFormat() + ' ' + getLocalizedShortTimeFormat()),
        formatDateTimeToLongDateTime: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedLongDateFormat() + ' ' + getLocalizedLongTimeFormat(), getDateTimeFormatOptions()),
        formatDateTimeToShortDateTime: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedShortDateFormat() + ' ' + getLocalizedShortTimeFormat(), getDateTimeFormatOptions()),
        formatDateTimeToLongDate: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedLongDateFormat(), getDateTimeFormatOptions()),
        formatDateTimeToShortDate: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedShortDateFormat(), getDateTimeFormatOptions()),
        formatDateTimeToLongMonthDay: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedLongMonthDayFormat(), getDateTimeFormatOptions()),
        formatDateTimeToShortMonthDay: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedShortMonthDayFormat(), getDateTimeFormatOptions()),
        formatDateTimeToLongTime: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedLongTimeFormat(), getDateTimeFormatOptions()),
        formatDateTimeToShortTime: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedShortTimeFormat(), getDateTimeFormatOptions()),
        formatGregorianTextualYearMonthDayToLongDate: (date: TextualYearMonthDay) => formatGregorianCalendarYearDashMonthDashDay(date, getLocalizedLongDateFormat(), getDateTimeFormatOptions()),
        // format date time (Gregorian calendar and Gregorian-like calendar) functions
        formatDateTimeToGregorianLikeLongYear: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedLongYearFormat(), getDateTimeFormatOptions({ calendarType: getGregorianLikeCalendarType() })),
        formatDateTimeToGregorianLikeShortYear: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedShortYearFormat(), getDateTimeFormatOptions({ calendarType: getGregorianLikeCalendarType() })),
        formatDateTimeToGregorianLikeLongYearMonth: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedLongYearMonthFormat(), getDateTimeFormatOptions({ calendarType: getGregorianLikeCalendarType() })),
        formatDateTimeToGregorianLikeShortYearMonth: (dateTime: DateTime) => formatDateTime(dateTime, getLocalizedShortYearMonthFormat(), getDateTimeFormatOptions({ calendarType: getGregorianLikeCalendarType() })),
        formatDateTimeToGregorianLikeLongMonth: (dateTime: DateTime) => formatDateTime(dateTime, 'MMMM', getDateTimeFormatOptions({ calendarType: getGregorianLikeCalendarType() })),
        formatDateTimeToGregorianLikeShortMonth: (dateTime: DateTime) => formatDateTime(dateTime, 'MMM', getDateTimeFormatOptions({ calendarType: getGregorianLikeCalendarType() })),
        formatGregorianTextualMonthDayToGregorianLikeLongMonthDay,
        formatDateTimeToGregorianLikeYearQuarter,
        formatYearQuarterToGregorianLikeYearQuarter,
        formatDateTimeToGregorianLikeFiscalYear: (dateTime: DateTime) => formatUnixTimeToGregorianLikeFiscalYear(dateTime.getUnixTime()),
        formatGregorianYearToGregorianLikeFiscalYear,
        formatFiscalYearStartToGregorianLikeLongMonth,
        // format date time (Gregorian calendar) functions
        formatDateTimeToGregorianDefaultDateTime: (dateTime: DateTime) => formatDateTime(dateTime, KnownDateTimeFormat.DefaultDateTime.format, getDateTimeFormatOptions({ numeralSystem: NumeralSystem.WesternArabicNumerals, calendarType: CalendarType.Gregorian })),
        // other format date time functions
        formatDateRange,
        getTimezoneDifferenceDisplayText,
        getCalendarAlternateDates,
        getCalendarAlternateDate,
        // format amount/number functions
        parseAmountFromLocalizedNumerals: (value: string) => getParsedAmountNumber(value),
        parseAmountFromWesternArabicNumerals: (value: string) => getParsedAmountNumber(value, NumeralSystem.WesternArabicNumerals),
        formatAmountToLocalizedNumerals: (value: number, currencyCode?: string) => getFormattedAmount(value, undefined, undefined, currencyCode),
        formatAmountToWesternArabicNumerals: (value: number, currencyCode?: string) => getFormattedAmount(value, NumeralSystem.WesternArabicNumerals, undefined, currencyCode),
        formatAmountToLocalizedNumeralsWithoutDigitGrouping: (value: number, currencyCode?: string) => getFormattedAmount(value, undefined, DigitGroupingType.None, currencyCode),
        formatAmountToWesternArabicNumeralsWithoutDigitGrouping: (value: number, currencyCode?: string) => getFormattedAmount(value, NumeralSystem.WesternArabicNumerals, DigitGroupingType.None, currencyCode),
        formatAmountToLocalizedNumeralsWithCurrency: (value: number | HiddenAmount | NumberWithSuffix, currencyCode?: string | false, currencyDisplayType?: CurrencyDisplayType) => getFormattedAmountWithCurrency(value, currencyCode, currencyDisplayType),
        formatAmountToWesternArabicNumeralsWithCurrency: (value: number | HiddenAmount | NumberWithSuffix, currencyCode?: string | false, currencyDisplayType?: CurrencyDisplayType) => getFormattedAmountWithCurrency(value, currencyCode, currencyDisplayType, NumeralSystem.WesternArabicNumerals),
        formatNumberToLocalizedNumerals: (value: number, precision?: number) => getFormattedNumber(value, undefined, precision),
        formatNumberToWesternArabicNumerals: (value: number, precision?: number) => getFormattedNumber(value, NumeralSystem.WesternArabicNumerals, precision),
        formatPercentToLocalizedNumerals: (value: number, precision: number, lowPrecisionValue: string) => getFormattedPercentValue(value, precision, lowPrecisionValue),
        formatPercentToWesternArabicNumerals: (value: number, precision: number, lowPrecisionValue: string) => getFormattedPercentValue(value, precision, lowPrecisionValue, NumeralSystem.WesternArabicNumerals),
        formatExchangeRateAmountToWesternArabicNumerals: (value: number) => getFormattedExchangeRateAmount(value, NumeralSystem.WesternArabicNumerals),
        appendDigitGroupingSymbolAndDecimalSeparator: (value: string) => appendDigitGroupingSymbolAndDecimalSeparator(value, getNumberFormatOptions({})),
        getAdaptiveAmountRate,
        getAmountPrependAndAppendText,
        getCategorizedAccountsWithDisplayBalance,
        // other format functions
        getLocalizedFileEncodingName,
        getLocalizedOAuth2ProviderName,
        getLocalizedOAuth2LoginText,
        // localization setting functions
        setLanguage,
        setTimeZone,
        initLocale
    };
}
