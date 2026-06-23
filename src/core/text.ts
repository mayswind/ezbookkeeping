import { FOLDED_CHARACTER_MAP } from '@/consts/text.ts';

export enum TextDirection {
    LTR = 'ltr',
    RTL = 'rtl'
}

export class NormalizedText {
    private readonly text: string;
    private cachedLowerCase: string | undefined;
    private cachedNormalized: string | undefined;

    private constructor(text: string) {
        this.text = text;
    }

    public get originalText(): string {
        return this.text;
    }

    public get lowerCaseText(): string {
        if (this.cachedLowerCase === undefined) {
            this.cachedLowerCase = NormalizedText.caseInsensitiveText(this.text);
        }

        return this.cachedLowerCase;
    }

    public get normalizedText(): string {
        if (this.cachedNormalized === undefined) {
            this.cachedNormalized = NormalizedText.normalizeForSearch(this.text);
        }

        return this.cachedNormalized;
    }

    public static caseInsensitiveText(text: string): string {
        return text.toLowerCase();
    }

    public static normalizeForSearch(text: string): string {
        if (text.length === 0) {
            return text;
        }

        let hasNonAscii: boolean = false;
        let needsLower: boolean = false;

        for (let i = 0; i < text.length; i++) {
            const c = text.charCodeAt(i);

            if (c > 127) {
                hasNonAscii = true;
                break;
            }

            if (c >= 65 && c <= 90) {
                needsLower = true;
            }
        }

        if (!hasNonAscii) {
            return needsLower ? text.toLowerCase() : text;
        }

        const nfkdText = text.normalize("NFKD");
        const finalChars: string[] = [];

        for (const char of nfkdText) {
            // case folding
            const foldedText = FOLDED_CHARACTER_MAP[char] ?? char;

            for (const foldedChar of foldedText) {
                const cp = foldedChar.codePointAt(0) as number;

                // strip diacritics
                if ((cp >= 0x0300 && cp <= 0x036F) ||
                    (cp >= 0x1AB0 && cp <= 0x1AFF) ||
                    (cp >= 0x1DC0 && cp <= 0x1DFF) ||
                    (cp >= 0x20D0 && cp <= 0x20FF) ||
                    (cp >= 0xFE20 && cp <= 0xFE2F)) {
                    continue;
                }

                // for Japanese, replace hiragana with katakana
                if (cp >= 0x3041 && cp <= 0x3096) {
                    finalChars.push(String.fromCharCode(cp + 0x60));
                    continue;
                }

                finalChars.push(foldedChar);
            }
        }

        return finalChars.join('');
    }

    public static of(text: string): NormalizedText {
        return new NormalizedText(text);
    }
}
