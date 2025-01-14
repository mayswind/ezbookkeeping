import en from './en.json';
import vi from './vi.json';
import zhHans from './zh_Hans.json';

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
}

export const DEFAULT_LANGUAGE: string = 'en';

// To add new languages, please refer to https://ezbookkeeping.mayswind.net/translating
export const ALL_LANGUAGES: Record<string, LanguageInfo> = {
    'en': {
        name: 'English',
        displayName: 'English',
        alternativeLanguageTag: 'en',
        content: en
    },
    'vi': {
        name: 'Vietnamese',
        displayName: 'Tiếng Việt',
        alternativeLanguageTag: 'vi-VN',
        content: vi
    },
    'zh-Hans': {
        name: 'Simplified Chinese',
        displayName: '简体中文',
        alternativeLanguageTag: 'zh-CN',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        content: zhHans
    }
};
