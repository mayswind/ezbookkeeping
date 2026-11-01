import {
    type ImportFileTypeSupportedAdditionalOptions,
    ImportCustomFileFormatProcessMethod
} from '@/core/import_transaction.ts';

export class ImportPresetBasicInfo {
    public id: string = '';
    public name: string = '';
    public displayOrder: number;

    public constructor(id: string, name: string, displayOrder: number) {
        this.id = id;
        this.name = name;
        this.displayOrder = displayOrder;
    }

    public static of(response: ImportPresetInfoResponse): ImportPresetBasicInfo {
        return new ImportPresetBasicInfo(
            response.id,
            response.name,
            response.displayOrder
        );
    }

    public static ofMulti(presetResponses: ImportPresetInfoResponse[]): ImportPresetBasicInfo[] {
        const presets: ImportPresetBasicInfo[] = [];

        for (const presetResponse of presetResponses) {
            presets.push(ImportPresetBasicInfo.of(presetResponse));
        }

        return presets;
    }
}

export class ImportPreset implements ImportPresetInfoResponse {
    public id: string = '';
    public name: string = '';
    public displayOrder: number;
    public fileType: string = 'ezbookkeeping';
    public fileSubType: string = 'ezbookkeeping_csv';
    public fileEncoding: string = 'auto';
    public handlingMethod: ImportCustomFileFormatProcessMethod = ImportCustomFileFormatProcessMethod.ColumnMapping;
    public additionalOptions: ImportFileTypeSupportedAdditionalOptions = {};
    public additionalPrompt: string = '';
    public dataMapping: string = '';
    public customScript: string = '';
    public autoExecuteScript: boolean = false;
    public replaceRules: string = '';

    private constructor(id: string, name: string, displayOrder: number, fileType: string, fileSubType: string, fileEncoding: string, handlingMethod: ImportCustomFileFormatProcessMethod, additionalOptions: ImportFileTypeSupportedAdditionalOptions, additionalPrompt: string, dataMapping: string, customScript: string, autoExecuteScript: boolean, replaceRules: string) {
        this.id = id;
        this.name = name;
        this.displayOrder = displayOrder;
        this.fileType = fileType;
        this.fileSubType = fileSubType;
        this.fileEncoding = fileEncoding;
        this.handlingMethod = handlingMethod;
        this.additionalOptions = { ...additionalOptions };
        this.additionalPrompt = additionalPrompt;
        this.dataMapping = dataMapping;
        this.customScript = customScript;
        this.autoExecuteScript = autoExecuteScript;
        this.replaceRules = replaceRules;
    }

    public fillFrom(other: ImportPresetInfoResponse): void {
        this.id = other.id;
        this.name = other.name;
        this.displayOrder = other.displayOrder;
        this.fileType = other.fileType;
        this.fileSubType = other.fileSubType;
        this.fileEncoding = other.fileEncoding;
        this.handlingMethod = other.handlingMethod;
        this.additionalOptions = { ...other.additionalOptions };
        this.additionalPrompt = other.additionalPrompt;
        this.dataMapping = other.dataMapping;
        this.customScript = other.customScript;
        this.autoExecuteScript = other.autoExecuteScript;
        this.replaceRules = other.replaceRules;
    }

    public toCreateRequest(clientSessionId: string): ImportPresetCreateRequest {
        return {
            name: this.name,
            fileType: this.fileType,
            fileSubType: this.fileSubType,
            fileEncoding: this.fileEncoding,
            handlingMethod: this.handlingMethod,
            additionalOptions: this.additionalOptions,
            additionalPrompt: this.additionalPrompt,
            dataMapping: this.dataMapping,
            customScript: this.customScript,
            autoExecuteScript: this.autoExecuteScript,
            replaceRules: this.replaceRules,
            clientSessionId: clientSessionId
        };
    }

    public toModifyRequest(): ImportPresetModifyRequest {
        return {
            id: this.id,
            name: this.name,
            fileType: this.fileType,
            fileSubType: this.fileSubType,
            fileEncoding: this.fileEncoding,
            handlingMethod: this.handlingMethod,
            additionalOptions: this.additionalOptions,
            additionalPrompt: this.additionalPrompt,
            dataMapping: this.dataMapping,
            customScript: this.customScript,
            autoExecuteScript: this.autoExecuteScript,
            replaceRules: this.replaceRules
        }
    }

    public static of(response: ImportPresetInfoResponse): ImportPreset {
        return new ImportPreset(
            response.id,
            response.name,
            response.displayOrder,
            response.fileType,
            response.fileSubType,
            response.fileEncoding,
            response.handlingMethod,
            { ...response.additionalOptions },
            response.additionalPrompt,
            response.dataMapping,
            response.customScript,
            response.autoExecuteScript,
            response.replaceRules
        );
    }

    public static createNewPreset(): ImportPreset {
        return new ImportPreset(
            '',
            '',
            0,
            'ezbookkeeping',
            'ezbookkeeping_csv',
            'auto',
            ImportCustomFileFormatProcessMethod.ColumnMapping,
            {},
            '',
            '',
            '',
            false,
            ''
        );
    }
}

export interface ImportPresetCreateRequest {
    readonly name: string;
    readonly fileType: string;
    readonly fileSubType: string;
    readonly fileEncoding: string;
    readonly handlingMethod: ImportCustomFileFormatProcessMethod;
    readonly additionalOptions: ImportFileTypeSupportedAdditionalOptions;
    readonly additionalPrompt: string;
    readonly dataMapping: string;
    readonly customScript: string;
    readonly autoExecuteScript: boolean;
    readonly replaceRules: string;
    readonly clientSessionId: string;
}

export interface ImportPresetModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly fileType: string;
    readonly fileSubType: string;
    readonly fileEncoding: string;
    readonly handlingMethod: ImportCustomFileFormatProcessMethod;
    readonly additionalOptions: ImportFileTypeSupportedAdditionalOptions;
    readonly additionalPrompt: string;
    readonly dataMapping: string;
    readonly customScript: string;
    readonly autoExecuteScript: boolean;
    readonly replaceRules: string;
}

export interface ImportPresetInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly displayOrder: number;
    readonly fileType: string;
    readonly fileSubType: string;
    readonly fileEncoding: string;
    readonly handlingMethod: ImportCustomFileFormatProcessMethod;
    readonly additionalOptions: ImportFileTypeSupportedAdditionalOptions;
    readonly additionalPrompt: string;
    readonly dataMapping: string;
    readonly customScript: string;
    readonly autoExecuteScript: boolean;
    readonly replaceRules: string;
}

export interface ImportPresetMoveRequest {
    readonly newDisplayOrders: ImportPresetNewDisplayOrderRequest[];
}

export interface ImportPresetNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface ImportPresetDeleteRequest {
    readonly id: string;
}
