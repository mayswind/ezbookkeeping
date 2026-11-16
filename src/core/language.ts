import { TextDirection } from './text.ts';

export interface LanguageInfo {
    readonly name: string;
    readonly displayName: string;
    readonly alternativeLanguageTag: string;
    readonly aliases?: string[];
    readonly textDirection: TextDirection;
    readonly loadContent: () => Promise<LanguageMessages>;
}

export interface LanguageMessages {
    [key: string]: string | LanguageMessages;
}

export interface LanguageMessagesCacheEntry {
    readonly languageKey: string;
    readonly updatedAt: number;
    readonly data: LanguageMessages;
}

export interface LanguageOption {
    readonly languageTag: string;
    readonly displayName: string;
    readonly nativeDisplayName: string;
}
