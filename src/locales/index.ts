import en from './en.json';

export interface LanguageMessages {
    [key: string]: string | LanguageMessages;
}

export interface LanguageInfo {
    readonly name: string;
    readonly displayName: string;
    readonly alternativeLanguageTag: string;
    readonly aliases?: string[];
    readonly textDirection: 'ltr' | 'rtl';
    readonly loadContent: () => Promise<LanguageMessages>;
}

export interface LanguageOption {
    readonly languageTag: string;
    readonly displayName: string;
    readonly nativeDisplayName: string;
}

export const DEFAULT_LANGUAGE: string = 'en';
export const DEFAULT_LANGUAGE_CONTENT: LanguageMessages = en;

// To add new languages, please refer to https://ezbookkeeping.mayswind.net/translating
export const ALL_LANGUAGES: Record<string, LanguageInfo> = {
    'de': {
        name: 'German',
        displayName: 'Deutsch',
        alternativeLanguageTag: 'de-DE',
        textDirection: 'ltr',
        loadContent: () => import('./de.json').then(module => module.default)
    },
    'el': {
        name: 'Greek',
        displayName: 'Ελληνικά',
        alternativeLanguageTag: 'el-GR',
        textDirection: 'ltr',
        loadContent: () => import('./el.json').then(module => module.default)
    },
    'en': {
        name: 'English',
        displayName: 'English',
        alternativeLanguageTag: 'en-US',
        textDirection: 'ltr',
        loadContent: () => Promise.resolve(DEFAULT_LANGUAGE_CONTENT)
    },
    'es': {
        name: 'Spanish',
        displayName: 'Español',
        alternativeLanguageTag: 'es-ES',
        textDirection: 'ltr',
        loadContent: () => import('./es.json').then(module => module.default)
    },
    'fr': {
        name: "French",
        displayName: "Français",
        alternativeLanguageTag: "fr-FR",
        textDirection: "ltr",
        loadContent: () => import('./fr.json').then(module => module.default),
    },
    'hi': {
        name: 'Hindi',
        displayName: 'हिन्दी',
        alternativeLanguageTag: 'hi-IN',
        textDirection: 'ltr',
        loadContent: () => import('./hi.json').then(module => module.default)
    },
    'it': {
        name: 'Italian',
        displayName: 'Italiano',
        alternativeLanguageTag: 'it-IT',
        textDirection: 'ltr',
        loadContent: () => import('./it.json').then(module => module.default)
    },
    'ja': {
        name: 'Japanese',
        displayName: '日本語',
        alternativeLanguageTag: 'ja-JP',
        textDirection: 'ltr',
        loadContent: () => import('./ja.json').then(module => module.default)
    },
    'kn': {
        name: 'Kannada',
        displayName: 'ಕನ್ನಡ',
        alternativeLanguageTag: 'kn-IN',
        textDirection: 'ltr',
        loadContent: () => import('./kn.json').then(module => module.default)
    },
    'ko': {
        name: 'Korean',
        displayName: '한국어',
        alternativeLanguageTag: 'ko-KR',
        textDirection: 'ltr',
        loadContent: () => import('./ko.json').then(module => module.default)
    },
    'nl': {
        name: 'Dutch',
        displayName: 'Nederlands',
        alternativeLanguageTag: 'nl-NL',
        textDirection: 'ltr',
        loadContent: () => import('./nl.json').then(module => module.default)
    },
    'pl': {
        name: 'Polish',
        displayName: 'Polski',
        alternativeLanguageTag: 'pl-PL',
        textDirection: 'ltr',
        loadContent: () => import('./pl.json').then(module => module.default)
    },
    'pt-BR': {
        name: 'Portuguese (Brazil)',
        displayName: 'Português (Brasil)',
        alternativeLanguageTag: 'pt-BR',
        textDirection: 'ltr',
        loadContent: () => import('./pt_BR.json').then(module => module.default)
    },
    'ro': {
        name: 'Romanian',
        displayName: 'Română',
        alternativeLanguageTag: 'ro-RO',
        textDirection: 'ltr',
        loadContent: () => import('./ro.json').then(module => module.default)
    },
    'ru': {
        name: 'Russian',
        displayName: 'Русский',
        alternativeLanguageTag: 'ru-RU',
        textDirection: 'ltr',
        loadContent: () => import('./ru.json').then(module => module.default)
    },
    'sk': {
        name: 'Slovak',
        displayName: 'Slovenčina',
        alternativeLanguageTag: 'sk-SK',
        textDirection: 'ltr',
        loadContent: () => import('./sk.json').then(module => module.default)
    },
    'sl': {
        name: 'Slovenian',
        displayName: 'Slovenščina',
        alternativeLanguageTag: 'sl-SI',
        textDirection: 'ltr',
        loadContent: () => import('./sl.json').then(module => module.default)
    },
    'ta': {
        name: 'Tamil',
        displayName: 'தமிழ்',
        alternativeLanguageTag: 'ta-IN',
        textDirection: 'ltr',
        loadContent: () => import('./ta.json').then(module => module.default)
    },
    'th': {
        name: 'Thai',
        displayName: 'ภาษาไทย',
        alternativeLanguageTag: 'th-TH',
        textDirection: 'ltr',
        loadContent: () => import('./th.json').then(module => module.default)
    },
    'tr': {
        name: 'Turkish',
        displayName: 'Türkçe',
        alternativeLanguageTag: 'tr-TR',
        textDirection: 'ltr',
        loadContent: () => import('./tr.json').then(module => module.default)
    },
    'uk': {
        name: 'Ukrainian',
        displayName: 'Українська',
        alternativeLanguageTag: 'uk-UA',
        textDirection: 'ltr',
        loadContent: () => import('./uk.json').then(module => module.default)
    },
    'vi': {
        name: 'Vietnamese',
        displayName: 'Tiếng Việt',
        alternativeLanguageTag: 'vi-VN',
        textDirection: 'ltr',
        loadContent: () => import('./vi.json').then(module => module.default)
    },
    'zh-Hans': {
        name: 'Chinese (Simplified)',
        displayName: '中文 (简体)',
        alternativeLanguageTag: 'zh-CN',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        textDirection: 'ltr',
        loadContent: () => import('./zh_Hans.json').then(module => module.default)
    },
    'zh-Hant': {
        name: 'Chinese (Traditional)',
        displayName: '中文 (繁體)',
        alternativeLanguageTag: 'zh-TW',
        aliases: ['zh-CHT', 'zh-TW', 'zh-HK', 'zh-MO'],
        textDirection: 'ltr',
        loadContent: () => import('./zh_Hant.json').then(module => module.default)
    },
};
