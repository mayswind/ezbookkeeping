import de from './de.json';
import en from './en.json';
import es from './es.json';
import it from './it.json';
import ja from './ja.json';
import ru from './ru.json';
import uk from './uk.json';
import vi from './vi.json';
import zhHans from './zh_Hans.json';
import zhHant from './zh_Hant.json';

export interface LanguageInfo {
    readonly name: string;
    readonly displayName: string;
    readonly alternativeLanguageTag: string;
    readonly aliases?: string[];
    readonly content: object;
}

export interface LanguageOption {
    readonly languageTag: string;
    readonly displayName: string;
    readonly nativeDisplayName: string;
}

export const DEFAULT_LANGUAGE: string = 'en';

// To add new languages, please refer to https://ezbookkeeping.mayswind.net/translating
export const ALL_LANGUAGES: Record<string, LanguageInfo> = {
    'de': {
        name: 'German',
        displayName: 'Deutsch',
        alternativeLanguageTag: 'de-DE',
        content: de
    },
    'en': {
        name: 'English',
        displayName: 'English',
        alternativeLanguageTag: 'en-US',
        content: en
    },
    'es': {
        name: 'Spanish',
        displayName: 'Español',
        alternativeLanguageTag: 'es-ES',
        content: es
    },
    'it': {
        name: 'Italian',
        displayName: 'Italiano',
        alternativeLanguageTag: 'it-IT',
        content: it
    },
    'ja': {
        name: 'Japanese',
        displayName: '日本語',
        alternativeLanguageTag: 'ja-JP',
        content: ja
    },
    'ru': {
        name: 'Russian',
        displayName: 'Русский',
        alternativeLanguageTag: 'ru-RU',
        content: ru
    },
    'uk': {
        name: 'Ukrainian',
        displayName: 'Українська',
        alternativeLanguageTag: 'uk-UA',
        content: uk
    },
    'vi': {
        name: 'Vietnamese',
        displayName: 'Tiếng Việt',
        alternativeLanguageTag: 'vi-VN',
        content: vi
    },
    'zh-Hans': {
        name: 'Chinese (Simplified)',
        displayName: '中文 (简体)',
        alternativeLanguageTag: 'zh-CN',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        content: zhHans
    },
    'zh-Hant': {
        name: 'Chinese (Traditional)',
        displayName: '中文 (繁體)',
        alternativeLanguageTag: 'zh-TW',
        aliases: ['zh-CHT', 'zh-TW', 'zh-HK', 'zh-MO'],
        content: zhHant
    }
};
