export class KnownFileType {
    private static readonly allInstancesByExtension: Record<string, KnownFileType> = {};

    public static readonly JSON = new KnownFileType('json', 'application/json');
    public static readonly CSV = new KnownFileType('csv', 'text/csv');
    public static readonly TSV = new KnownFileType('tsv', 'text/tab-separated-values');
    public static readonly SSV = new KnownFileType('txt', 'text/plain');
    public static readonly TXT = new KnownFileType('txt', 'text/plain');
    public static readonly MARKDOWN = new KnownFileType('md', 'text/markdown');
    public static readonly JS = new KnownFileType('js', 'application/javascript');
    public static readonly JPG = new KnownFileType('jpg', 'image/jpeg');

    public readonly extension: string;
    public readonly contentType: string;

    private constructor(extension: string, contentType: string) {
        this.extension = extension;
        this.contentType = contentType;

        KnownFileType.allInstancesByExtension[extension] = this;
    }

    public isSameType(contentType: string): boolean {
        if (!contentType) {
            return false;
        }

        return this.contentType === contentType || contentType.indexOf(this.contentType) === 0;
    }

    public formatFileName(fileName: string): string {
        if (fileName.endsWith(`.${this.extension}`)) {
            return fileName;
        }

        return `${fileName}.${this.extension}`;
    }

    public createBlob(content: string): Blob {
        return new Blob([content], {
            type: this.contentType,
        });
    }

    public createFile(content: string, fileName: string): File {
        return new File([content], this.formatFileName(fileName), {
            type: this.contentType,
        });
    }

    public createFileFromBlob(blob: Blob, fileName: string): File {
        return new File([blob], this.formatFileName(fileName), {
            type: this.contentType,
        });
    }

    public static parse(extension: string): KnownFileType | undefined {
        return KnownFileType.allInstancesByExtension[extension];
    }
}

export interface ImportFileTypeAndExtensions {
    readonly type: string;
    readonly extensions?: string;
}

export interface ImportFileCategoryAndTypes {
    readonly categoryName: string;
    readonly fileTypes: ImportFileType[];
}

export interface ImportFileTypeSupportedAdditionalOptions extends Record<string, boolean | undefined> {
    readonly payeeAsTag?: boolean;
    readonly payeeAsDescription?: boolean;
    readonly memberAsTag?: boolean;
    readonly projectAsTag?: boolean;
    readonly merchantAsTag?: boolean;
}

export interface ImportFileType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly name: string;
    readonly extensions: string;
    readonly subTypes?: ImportFileTypeSubType[];
    readonly supportedEncodings?: string[];
    readonly dataFromTextbox?: boolean;
    readonly supportedAdditionalOptions?: ImportFileTypeSupportedAdditionalOptions;
    readonly document?: {
        readonly supportMultiLanguages: boolean | string;
        readonly anchor: string;
    };
}

export interface ImportFileTypeSubType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly name: string;
    readonly extensions?: string;
}

export interface LocalizedImportFileCategoryAndTypes {
    readonly displayCategoryName: string;
    readonly fileTypes: LocalizedImportFileType[];
}

export interface LocalizedImportFileType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly displayName: string;
    readonly extensions: string;
    readonly subTypes?: LocalizedImportFileTypeSubType[];
    readonly supportedEncodings?: LocalizedImportFileTypeSupportedEncodings[];
    readonly dataFromTextbox?: boolean;
    readonly supportedAdditionalOptions?: ImportFileTypeSupportedAdditionalOptions;
    readonly document?: LocalizedImportFileDocument;
}

export interface LocalizedImportFileTypeSubType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly displayName: string;
    readonly extensions?: string;
}

export interface LocalizedImportFileTypeSupportedEncodings {
    readonly encoding: string;
    readonly displayName: string;
}

export interface LocalizedImportFileDocument {
    readonly language: string;
    readonly displayLanguageName: string;
    readonly anchor: string;
}
