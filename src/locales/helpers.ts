import { useI18n as useVueI18n } from 'vue-i18n';
import moment from 'moment-timezone';

import { type LanguageInfo, allLanguages, DEFAULT_LANGUAGE } from '@/locales/index.ts';

import {
    type LocalizedMeridiemIndicator,
    Month,
    WeekDay,
    MeridiemIndicator
} from '@/core/datetime.ts';

import type { LocaleDefaultSettings } from '@/core/setting.ts';
import type { ErrorResponse } from '@/core/api.ts';

import { KnownErrorCode, SPECIFIED_API_NOT_FOUND_ERRORS, PARAMETERIZED_ERRORS } from '@/consts/api.ts';

import {
    isString,
    isObject
} from '@/lib/common.ts';

import {
    isPM,
    getTimezoneOffset
} from '@/lib/datetime.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

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

    function getCurrencyName(currencyCode: string): string {
        return t(`currency.name.${currencyCode}`);
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
            defaultFirstDayOfWeek = WeekDay.parse(defaultFirstDayOfWeekName).type;
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
        tt: t,
        ti: translateIf,
        te: translateError,
        getCurrentLanguageTag,
        getCurrentLanguageInfo,
        getCurrentLanguageDisplayName,
        getDefaultCurrency,
        getDefaultFirstDayOfWeek,
        getCurrencyName,
        getAllMeridiemIndicators,
        getAllLongMonthNames,
        getAllShortMonthNames,
        getAllLongWeekdayNames,
        getAllShortWeekdayNames,
        getAllMinWeekdayNames,
        setLanguage,
        setTimeZone,
        initLocale
    };
}
