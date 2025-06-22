export class KnownFileType {
    private static readonly allInstancesByExtension: Record<string, KnownFileType> = {};

    public static readonly CSV = new KnownFileType('csv', 'text/csv');
    public static readonly TSV = new KnownFileType('tsv', 'text/tab-separated-values');
    public static readonly MARKDOWN = new KnownFileType('md', 'text/markdown');

    public readonly extension: string;
    public readonly contentType: string;

    private constructor(extension: string, contentType: string) {
        this.extension = extension;
        this.contentType = contentType;

        KnownFileType.allInstancesByExtension[extension] = this;
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

    public static parse(extension: string): KnownFileType | undefined {
        return KnownFileType.allInstancesByExtension[extension];
    }
}

export interface ImportFileTypeAndExtensions {
    readonly type: string;
    readonly extensions?: string;
}

export interface ImportFileType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly name: string;
    readonly extensions: string;
    readonly subTypes?: ImportFileTypeSubType[];
    readonly supportedEncodings?: string[];
    readonly dataFromTextbox?: boolean;
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

export interface LocalizedImportFileType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly displayName: string;
    readonly extensions: string;
    readonly subTypes?: LocalizedImportFileTypeSubType[];
    readonly supportedEncodings?: LocalizedImportFileTypeSupportedEncodings[];
    readonly dataFromTextbox?: boolean;
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
