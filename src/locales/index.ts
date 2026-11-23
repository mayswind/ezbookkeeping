import en from './en.json';

import { TextDirection } from '@/core/text.ts';
import type { LanguageInfo, LanguageMessages } from '@/core/language.ts';

export const DEFAULT_LANGUAGE: string = 'en';
export const DEFAULT_LANGUAGE_CONTENT: LanguageMessages = en;

// To add new languages, please refer to https://ezbookkeeping.mayswind.net/translating
export const ALL_LANGUAGES: Record<string, LanguageInfo> = {
    'de': {
        name: 'German',
        displayName: 'Deutsch',
        alternativeLanguageTag: 'de-DE',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./de.json').then(module => module.default)
    },
    'el': {
        name: 'Greek',
        displayName: 'Ελληνικά',
        alternativeLanguageTag: 'el-GR',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./el.json').then(module => module.default)
    },
    'en': {
        name: 'English',
        displayName: 'English',
        alternativeLanguageTag: 'en-US',
        textDirection: TextDirection.LTR,
        loadContent: () => Promise.resolve(DEFAULT_LANGUAGE_CONTENT)
    },
    'es': {
        name: 'Spanish',
        displayName: 'Español',
        alternativeLanguageTag: 'es-ES',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./es.json').then(module => module.default)
    },
    'fr': {
        name: "French",
        displayName: "Français",
        alternativeLanguageTag: "fr-FR",
        textDirection: TextDirection.LTR,
        loadContent: () => import('./fr.json').then(module => module.default),
    },
    'hi': {
        name: 'Hindi',
        displayName: 'हिन्दी',
        alternativeLanguageTag: 'hi-IN',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./hi.json').then(module => module.default)
    },
    'id': {
        name: 'Indonesian',
        displayName: 'Bahasa Indonesia',
        alternativeLanguageTag: 'id-ID',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./id.json').then(module => module.default)
    },
    'it': {
        name: 'Italian',
        displayName: 'Italiano',
        alternativeLanguageTag: 'it-IT',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./it.json').then(module => module.default)
    },
    'ja': {
        name: 'Japanese',
        displayName: '日本語',
        alternativeLanguageTag: 'ja-JP',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./ja.json').then(module => module.default)
    },
    'kn': {
        name: 'Kannada',
        displayName: 'ಕನ್ನಡ',
        alternativeLanguageTag: 'kn-IN',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./kn.json').then(module => module.default)
    },
    'ko': {
        name: 'Korean',
        displayName: '한국어',
        alternativeLanguageTag: 'ko-KR',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./ko.json').then(module => module.default)
    },
    'nl': {
        name: 'Dutch',
        displayName: 'Nederlands',
        alternativeLanguageTag: 'nl-NL',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./nl.json').then(module => module.default)
    },
    'pl': {
        name: 'Polish',
        displayName: 'Polski',
        alternativeLanguageTag: 'pl-PL',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./pl.json').then(module => module.default)
    },
    'pt-BR': {
        name: 'Portuguese (Brazil)',
        displayName: 'Português (Brasil)',
        alternativeLanguageTag: 'pt-BR',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./pt_BR.json').then(module => module.default)
    },
    'ro': {
        name: 'Romanian',
        displayName: 'Română',
        alternativeLanguageTag: 'ro-RO',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./ro.json').then(module => module.default)
    },
    'ru': {
        name: 'Russian',
        displayName: 'Русский',
        alternativeLanguageTag: 'ru-RU',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./ru.json').then(module => module.default)
    },
    'sk': {
        name: 'Slovak',
        displayName: 'Slovenčina',
        alternativeLanguageTag: 'sk-SK',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./sk.json').then(module => module.default)
    },
    'sl': {
        name: 'Slovenian',
        displayName: 'Slovenščina',
        alternativeLanguageTag: 'sl-SI',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./sl.json').then(module => module.default)
    },
    'ta': {
        name: 'Tamil',
        displayName: 'தமிழ்',
        alternativeLanguageTag: 'ta-IN',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./ta.json').then(module => module.default)
    },
    'th': {
        name: 'Thai',
        displayName: 'ภาษาไทย',
        alternativeLanguageTag: 'th-TH',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./th.json').then(module => module.default)
    },
    'tr': {
        name: 'Turkish',
        displayName: 'Türkçe',
        alternativeLanguageTag: 'tr-TR',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./tr.json').then(module => module.default)
    },
    'uk': {
        name: 'Ukrainian',
        displayName: 'Українська',
        alternativeLanguageTag: 'uk-UA',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./uk.json').then(module => module.default)
    },
    'vi': {
        name: 'Vietnamese',
        displayName: 'Tiếng Việt',
        alternativeLanguageTag: 'vi-VN',
        textDirection: TextDirection.LTR,
        loadContent: () => import('./vi.json').then(module => module.default)
    },
    'zh-Hans': {
        name: 'Chinese (Simplified)',
        displayName: '中文 (简体)',
        alternativeLanguageTag: 'zh-CN',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        textDirection: TextDirection.LTR,
        loadContent: () => import('./zh_Hans.json').then(module => module.default)
    },
    'zh-Hant': {
        name: 'Chinese (Traditional)',
        displayName: '中文 (繁體)',
        alternativeLanguageTag: 'zh-TW',
        aliases: ['zh-CHT', 'zh-TW', 'zh-HK', 'zh-MO'],
        textDirection: TextDirection.LTR,
        loadContent: () => import('./zh_Hant.json').then(module => module.default)
    },
};
