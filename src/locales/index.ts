import en from './en.json';
import ru from './ru.json';
import es from './es.json';
import vi from './vi.json';
import zhHans from './zh_Hans.json';
import de from './de.json';
import ja from './ja.json';

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
        alternativeLanguageTag: 'en',
        content: en
    },
    'es': {
        name: 'Spanish',
        displayName: 'Español',
        alternativeLanguageTag: 'es',
        content: es
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
    }
};
