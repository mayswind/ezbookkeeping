import en from './langs/en.js'
import zhHans from './langs/zh_Hans.js'

const defaultLanguage = 'en';

const allLanguages = {
    'en': {
        name: 'English',
        displayName: 'English',
        content: en
    },
    'zh-Hans': {
        name: 'Simplified Chinese',
        displayName: '简体中文',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        content: zhHans
    }
};

const i18nOptions = {
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

function getAllLanguages() {
    return allLanguages;
}

function getLanguage(locale) {
    return allLanguages[locale];
}

function getLocaleFromLanguageAlias(alias) {
    for (let locale in allLanguages) {
        if (!Object.prototype.hasOwnProperty.call(allLanguages, locale)) {
            continue;
        }

        const lang = allLanguages[locale];
        const aliases = lang.aliases;

        if (!aliases || aliases.length < 1) {
            continue;
        }

        for (let i = 0; i < aliases.length; i++) {
            if (aliases[i] === alias) {
                return locale;
            }
        }
    }

    return null;
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

    if (!allLanguages[browserLocale]) {
        return defaultLanguage;
    }

    return browserLocale;
}

export default {
    i18nOptions,
    getAllLanguages,
    getLanguage,
    getDefaultLanguage
};
