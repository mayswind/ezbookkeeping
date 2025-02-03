export interface ImportFileTypeAndExtensions {
    readonly type: string;
    readonly extensions?: string;
}

export interface ImportFileType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly name: string;
    readonly extensions: string;
    readonly subTypes?: ImportFileTypeSubType[];
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
    readonly document?: LocalizedImportFileDocument;
}

export interface LocalizedImportFileTypeSubType extends ImportFileTypeAndExtensions {
    readonly type: string;
    readonly displayName: string;
    readonly extensions?: string;
}

export interface LocalizedImportFileDocument {
    readonly language: string;
    readonly displayLanguageName: string;
    readonly anchor: string;
}
