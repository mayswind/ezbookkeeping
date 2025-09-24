import de from './de.json';
import en from './en.json';
import es from './es.json';
import fr from './fr.json';
import it from './it.json';
import ja from './ja.json';
import nl from './nl.json';
import ru from './ru.json';
import th from './th.json';
import uk from './uk.json';
import vi from './vi.json';
import zhHans from './zh_Hans.json';
import zhHant from './zh_Hant.json';
import ptBR from './pt_BR.json';

export interface LanguageInfo {
    readonly name: string;
    readonly displayName: string;
    readonly alternativeLanguageTag: string;
    readonly aliases?: string[];
    readonly textDirection: 'ltr' | 'rtl';
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
        textDirection: 'ltr',
        content: de
    },
    'en': {
        name: 'English',
        displayName: 'English',
        alternativeLanguageTag: 'en-US',
        textDirection: 'ltr',
        content: en
    },
    'es': {
        name: 'Spanish',
        displayName: 'Español',
        alternativeLanguageTag: 'es-ES',
        textDirection: 'ltr',
        content: es
    },
    'fr': {
        name: "French",
        displayName: "Français",
        alternativeLanguageTag: "fr-FR",
        textDirection: "ltr",
        content: fr,
    },
    'it': {
        name: 'Italian',
        displayName: 'Italiano',
        alternativeLanguageTag: 'it-IT',
        textDirection: 'ltr',
        content: it
    },
    'ja': {
        name: 'Japanese',
        displayName: '日本語',
        alternativeLanguageTag: 'ja-JP',
        textDirection: 'ltr',
        content: ja
    },
    'nl': {
        name: 'Dutch',
        displayName: 'Nederlands',
        alternativeLanguageTag: 'nl-NL',
        textDirection: 'ltr',
        content: nl
    },
    'pt-BR': {
        name: 'Portuguese (Brazil)',
        displayName: 'Português (Brasil)',
        alternativeLanguageTag: 'pt-BR',
        textDirection: 'ltr',
        content: ptBR
    },
    'ru': {
        name: 'Russian',
        displayName: 'Русский',
        alternativeLanguageTag: 'ru-RU',
        textDirection: 'ltr',
        content: ru
    },
    'th': {
        name: 'Thai',
        displayName: 'ภาษาไทย',
        alternativeLanguageTag: 'th-TH',
        textDirection: 'ltr',
        content: th
    },
    'uk': {
        name: 'Ukrainian',
        displayName: 'Українська',
        alternativeLanguageTag: 'uk-UA',
        textDirection: 'ltr',
        content: uk
    },
    'vi': {
        name: 'Vietnamese',
        displayName: 'Tiếng Việt',
        alternativeLanguageTag: 'vi-VN',
        textDirection: 'ltr',
        content: vi
    },
    'zh-Hans': {
        name: 'Chinese (Simplified)',
        displayName: '中文 (简体)',
        alternativeLanguageTag: 'zh-CN',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        textDirection: 'ltr',
        content: zhHans
    },
    'zh-Hant': {
        name: 'Chinese (Traditional)',
        displayName: '中文 (繁體)',
        alternativeLanguageTag: 'zh-TW',
        aliases: ['zh-CHT', 'zh-TW', 'zh-HK', 'zh-MO'],
        textDirection: 'ltr',
        content: zhHant
    },
};
