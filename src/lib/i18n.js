import { defaultLanguage, allLanguages } from '../locales/index.js';
import timezone from "../consts/timezone.js";
import currency from "../consts/currency.js";
import settings from "./settings";
import utils from './utils.js';

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

export function getAllLanguageInfos() {
    return allLanguages;
}

export function getLanguageInfo(locale) {
    return allLanguages[locale];
}

export function getDefaultLanguage() {
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
    }

    if (!allLanguages[browserLocale]) {
        return defaultLanguage;
    }

    return browserLocale;
}

export function transateIf(text, isTranslate, translateFn) {
    if (isTranslate) {
        return translateFn(text);
    }

    return text;
}

export function getAllLongMonthNames(translateFn) {
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

export function getAllShortMonthNames(translateFn) {
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

export function getAllLongWeekdayNames(translateFn) {
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

export function getAllShortWeekdayNames(translateFn) {
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

export function getAllMinWeekdayNames(translateFn) {
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

export function getAllTimezones(includeSystemDefault, translateFn) {
    const defaultTimezoneOffset = utils.getTimezoneOffset();
    const defaultTimezoneOffsetMinutes = utils.getTimezoneOffsetMinutes();
    const allTimezones = timezone.all;
    const allTimezoneInfos = [];

    for (let i = 0; i < allTimezones.length; i++) {
        allTimezoneInfos.push({
            name: allTimezones[i].timezoneName,
            utcOffset: (allTimezones[i].timezoneName !== 'Etc/GMT' ? utils.getTimezoneOffset(allTimezones[i].timezoneName) : ''),
            utcOffsetMinutes: utils.getTimezoneOffsetMinutes(allTimezones[i].timezoneName),
            displayName: translateFn(`timezone.${allTimezones[i].displayName}`)
        });
    }

    if (includeSystemDefault) {
        allTimezoneInfos.push({
            name: '',
            utcOffset: defaultTimezoneOffset,
            utcOffsetMinutes: defaultTimezoneOffsetMinutes,
            displayName: translateFn('System Default')
        });
    }

    allTimezoneInfos.sort(function(c1, c2){
        const utcOffset1 = parseInt(c1.utcOffset.replace(':', ''));
        const utcOffset2 = parseInt(c2.utcOffset.replace(':', ''));

        if (utcOffset1 !== utcOffset2) {
            return utcOffset1 - utcOffset2;
        }

        return c1.displayName.localeCompare(c2.displayName);
    })

    return allTimezoneInfos;
}

export function getAllCurrencies(translateFn) {
    const allCurrencyCodes = currency.all;
    const allCurrencies = [];

    for (let currencyCode in allCurrencyCodes) {
        if (!Object.prototype.hasOwnProperty.call(allCurrencyCodes, currencyCode)) {
            return;
        }

        allCurrencies.push({
            code: currencyCode,
            displayName: translateFn(`currency.${currencyCode}`)
        });
    }

    allCurrencies.sort(function(c1, c2){
        return c1.displayName.localeCompare(c2.displayName);
    })

    return allCurrencies;
}

export function getDisplayCurrency(value, currencyCode, notConvertValue, translateFn) {
    if (!utils.isNumber(value) && !utils.isString(value)) {
        return value;
    }

    if (utils.isNumber(value)) {
        value = value.toString();
    }

    if (!notConvertValue) {
        const hasIncompleteFlag = utils.isString(value) && value.charAt(value.length - 1) === '+';

        if (hasIncompleteFlag) {
            value = value.substring(0, value.length - 1);
        }

        value = utils.numericCurrencyToString(value);

        if (hasIncompleteFlag) {
            value = value + '+';
        }
    }

    const currencyDisplayMode = settings.getCurrencyDisplayMode();

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
        const currencyName = translateFn(`currency.${currencyCode}`);
        return `${value} ${currencyName}`;
    } else {
        return value;
    }
}

export function getI18nOptions() {
    return {
        legacy: false,
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

export function getLocalizedError(error) {
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

export function getLocalizedErrorParameters(parameters, i18nFunc) {
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
