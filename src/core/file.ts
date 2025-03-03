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
