import type { TypeAndName, TypeAndDisplayName } from './base.ts';
import { KnownAmountFormat } from './numeral.ts';
import { KnownDateTimeFormat } from './datetime.ts';
import { KnownDateTimezoneFormat } from './timezone.ts';
import { TransactionType } from './transaction.ts';

export class ImportTransactionColumnType implements TypeAndName {
    private static readonly allInstances: ImportTransactionColumnType[] = [];

    public static readonly TransactionTime = new ImportTransactionColumnType(1, 'Transaction Time');
    public static readonly TransactionTimezone = new ImportTransactionColumnType(2, 'Transaction Timezone');
    public static readonly TransactionType = new ImportTransactionColumnType(3, 'Transaction Type');
    public static readonly Category = new ImportTransactionColumnType(4, 'Category');
    public static readonly SubCategory = new ImportTransactionColumnType(5, 'Secondary Category');
    public static readonly AccountName = new ImportTransactionColumnType(6, 'Account Name');
    public static readonly AccountCurrency = new ImportTransactionColumnType(7, 'Currency');
    public static readonly Amount = new ImportTransactionColumnType(8, 'Amount');
    public static readonly RelatedAccountName = new ImportTransactionColumnType(9, 'Transfer In Account Name');
    public static readonly RelatedAccountCurrency = new ImportTransactionColumnType(10, 'Transfer In Currency');
    public static readonly RelatedAmount = new ImportTransactionColumnType(11, 'Transfer In Amount');
    public static readonly GeographicLocation = new ImportTransactionColumnType(12, 'Geographic Location');
    public static readonly Tags = new ImportTransactionColumnType(13, 'Tags');
    public static readonly Description = new ImportTransactionColumnType(14, 'Description');

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        ImportTransactionColumnType.allInstances.push(this);
    }

    public static values(): ImportTransactionColumnType[] {
        return ImportTransactionColumnType.allInstances;
    }
}

export class ImportTransactionDataMapping {
    private static readonly JSON_ROOT_FIELD = 'ezBookkeepingImportTransactionDataMapping';
    private static readonly DEFAULT_INCLUDE_HEADER = true;
    private static readonly DEFAULT_TIME_FORMAT = '';
    private static readonly DEFAULT_TIMEZONE_FORMAT = '';
    private static readonly DEFAULT_AMOUNT_FORMAT = '';
    private static readonly DEFAULT_GEO_LOCATION_SEPARATOR = ' ';
    private static readonly DEFAULT_GEO_LOCATION_ORDER = 'lonlat';
    private static readonly DEFAULT_TAG_SEPARATOR = ';';

    public includeHeader: boolean;
    public dataColumnMapping: Record<number, number>;
    public transactionTypeMapping: Record<string, TransactionType>;
    public timeFormat: string;
    public timezoneFormat: string;
    public amountFormat: string;
    public geoLocationSeparator: string;
    public geoLocationOrder: string;
    public tagSeparator: string;

    private constructor(includeHeader: boolean,
                        dataColumnMapping: Record<number, number>,
                        transactionTypeMapping: Record<string, TransactionType>,
                        timeFormat: string,
                        timezoneFormat: string,
                        amountFormat: string,
                        geoLocationSeparator: string,
                        geoLocationOrder: string,
                        tagSeparator: string) {
        this.includeHeader = includeHeader;
        this.dataColumnMapping = dataColumnMapping;
        this.transactionTypeMapping = transactionTypeMapping;
        this.timeFormat = timeFormat;
        this.timezoneFormat = timezoneFormat;
        this.amountFormat = amountFormat;
        this.geoLocationSeparator = geoLocationSeparator;
        this.geoLocationOrder = geoLocationOrder;
        this.tagSeparator = tagSeparator;
    }

    public isColumnMappingSet(column: ImportTransactionColumnType | TypeAndDisplayName): boolean {
        return this.dataColumnMapping.hasOwnProperty(column.type) && typeof(this.dataColumnMapping[column.type]) === 'number';
    }

    public toggleIncludeHeader(): void {
        this.includeHeader = !this.includeHeader;
    }

    public toggleDataMappingColumn(columnIndex: number, columnType: number): void {
        if (this.dataColumnMapping[columnType] === columnIndex) {
            delete this.dataColumnMapping[columnType];
        } else {
            this.dataColumnMapping[columnType] = columnIndex;
        }

        for (const otherColumnType in this.dataColumnMapping) {
            if (otherColumnType !== columnType.toString() && this.dataColumnMapping[otherColumnType] === columnIndex) {
                delete this.dataColumnMapping[otherColumnType];
            }
        }
    }

    public setGeoLocationFormat(separator: string, order: string): void {
        this.geoLocationSeparator = separator;
        this.geoLocationOrder = order;
    }

    public formatGeoLocation(latitude: string, longitude: string): string {
        if (this.geoLocationOrder === 'latlon') {
            return `${latitude}${this.geoLocationSeparator}${longitude}`;
        } else {
            return `${longitude}${this.geoLocationSeparator}${latitude}`;
        }
    }

    public parseFileAllTransactionTypes(fileData: string[][] | undefined): string[] {
        if (!fileData || !fileData.length || !this.isColumnMappingSet(ImportTransactionColumnType.TransactionType)) {
            return [];
        }

        const allTypeMap: Record<string, boolean> = {};
        const allTypes: string[] = [];
        const typeColumnIndex = this.dataColumnMapping[ImportTransactionColumnType.TransactionType.type];

        const startIndex = this.includeHeader ? 1 : 0;

        for (let i = startIndex; i < fileData.length; i++) {
            if (fileData[i].length <= typeColumnIndex) {
                continue;
            }

            const type = fileData[i][typeColumnIndex];

            if (type && !allTypeMap[type]) {
                allTypes.push(type);
                allTypeMap[type] = true;
            }
        }

        return allTypes;
    }

    public parseFileValidMappedTransactionTypes(fileData: string[][] | undefined): Record<string, TransactionType> {
        if (!fileData || !fileData.length || !this.isColumnMappingSet(ImportTransactionColumnType.TransactionType)) {
            return {};
        }

        const result: Record<string, TransactionType> = {};

        if (!this.transactionTypeMapping) {
            return result;
        }

        for (const name in this.transactionTypeMapping) {
            if (!Object.prototype.hasOwnProperty.call(this.transactionTypeMapping, name)) {
                continue;
            }

            const type = this.transactionTypeMapping[name];

            if (TransactionType.ModifyBalance <= type && type <= TransactionType.Transfer) {
                result[name] = type;
            }
        }

        return result;
    }

    public parseFileAutoDetectedTimeFormat(fileData: string[][] | undefined): string | undefined {
        if (!fileData || !fileData.length || !this.isColumnMappingSet(ImportTransactionColumnType.TransactionTime)) {
            return undefined;
        }

        const allDateTimes: string[] = [];
        const dateTimeColumnIndex = this.dataColumnMapping[ImportTransactionColumnType.TransactionTime.type];

        const startIndex = this.includeHeader ? 1 : 0;

        for (let i = startIndex; i < fileData.length; i++) {
            if (fileData[i].length <= dateTimeColumnIndex) {
                continue;
            }

            const dateTime = fileData[i][dateTimeColumnIndex];

            if (dateTime) {
                allDateTimes.push(dateTime);
            }
        }

        const detectedFormats = KnownDateTimeFormat.detectMulti(allDateTimes);

        if (!detectedFormats || !detectedFormats.length || detectedFormats.length > 1) {
            return undefined;
        }

        return detectedFormats[0].format;
    }

    public parseFileAutoDetectedTimezoneFormat(fileData: string[][] | undefined): string | undefined {
        if (!fileData || !fileData.length || !this.isColumnMappingSet(ImportTransactionColumnType.TransactionTimezone)) {
            return undefined;
        }

        const allTimezones: string[] = [];
        const timezoneColumnIndex = this.dataColumnMapping[ImportTransactionColumnType.TransactionTimezone.type];

        const startIndex = this.includeHeader ? 1 : 0;

        for (let i = startIndex; i < fileData.length; i++) {
            if (fileData[i].length <= timezoneColumnIndex) {
                continue;
            }

            const timezone = fileData[i][timezoneColumnIndex];

            if (timezone) {
                allTimezones.push(timezone);
            }
        }

        const detectedFormats = KnownDateTimezoneFormat.detectMulti(allTimezones);

        if (!detectedFormats || !detectedFormats.length || detectedFormats.length > 1) {
            return undefined;
        }

        return detectedFormats[0].value;
    }

    public parseFileAutoDetectedAmountFormat(fileData: string[][] | undefined): string | undefined {
        if (!fileData || !fileData.length || !this.isColumnMappingSet(ImportTransactionColumnType.TransactionTimezone)) {
            return undefined;
        }

        const allAmounts: string[] = [];
        const amountColumnIndex = this.dataColumnMapping[ImportTransactionColumnType.Amount.type];

        const startIndex = this.includeHeader ? 1 : 0;

        for (let i = startIndex; i < fileData.length; i++) {
            if (fileData[i].length <= amountColumnIndex) {
                continue;
            }

            const amount = fileData[i][amountColumnIndex];

            if (amount) {
                allAmounts.push(amount);
            }
        }

        const detectedFormats = KnownAmountFormat.detectMulti(allAmounts);

        if (!detectedFormats || !detectedFormats.length) {
            return undefined;
        }

        return detectedFormats[0].type;
    }

    public reset(): void {
        this.includeHeader = ImportTransactionDataMapping.DEFAULT_INCLUDE_HEADER;
        this.dataColumnMapping = {};
        this.transactionTypeMapping = {};
        this.timeFormat = ImportTransactionDataMapping.DEFAULT_TIME_FORMAT;
        this.timezoneFormat = ImportTransactionDataMapping.DEFAULT_TIMEZONE_FORMAT;
        this.amountFormat = ImportTransactionDataMapping.DEFAULT_AMOUNT_FORMAT;
        this.geoLocationSeparator = ImportTransactionDataMapping.DEFAULT_GEO_LOCATION_SEPARATOR;
        this.geoLocationOrder = ImportTransactionDataMapping.DEFAULT_GEO_LOCATION_ORDER;
        this.tagSeparator = ImportTransactionDataMapping.DEFAULT_TAG_SEPARATOR;
    }

    public toJson(): string {
        return JSON.stringify({
            [ImportTransactionDataMapping.JSON_ROOT_FIELD]: {
                includeHeader: this.includeHeader,
                dataColumnMapping: this.dataColumnMapping,
                transactionTypeMapping: this.transactionTypeMapping,
                timeFormat: this.timeFormat,
                timezoneFormat: this.timezoneFormat,
                amountFormat: this.amountFormat,
                geoLocationSeparator: this.geoLocationSeparator,
                geoLocationOrder: this.geoLocationOrder,
                tagSeparator: this.tagSeparator
            }
        });
    }

    public static createEmpty(): ImportTransactionDataMapping {
        return new ImportTransactionDataMapping(
            ImportTransactionDataMapping.DEFAULT_INCLUDE_HEADER,
            {},
            {},
            ImportTransactionDataMapping.DEFAULT_TIME_FORMAT,
            ImportTransactionDataMapping.DEFAULT_TIMEZONE_FORMAT,
            ImportTransactionDataMapping.DEFAULT_AMOUNT_FORMAT,
            ImportTransactionDataMapping.DEFAULT_GEO_LOCATION_SEPARATOR,
            ImportTransactionDataMapping.DEFAULT_GEO_LOCATION_ORDER,
            ImportTransactionDataMapping.DEFAULT_TAG_SEPARATOR
        );
    }

    public static parseFromJson(json: string): ImportTransactionDataMapping | null {
        try {
            const parsed = JSON.parse(json);
            const root = parsed[ImportTransactionDataMapping.JSON_ROOT_FIELD];

            if (!root) {
                return null;
            }

            return new ImportTransactionDataMapping(
                root.includeHeader ?? ImportTransactionDataMapping.DEFAULT_INCLUDE_HEADER,
                root.dataColumnMapping ?? {},
                root.transactionTypeMapping ?? {},
                root.timeFormat ?? ImportTransactionDataMapping.DEFAULT_TIME_FORMAT,
                root.timezoneFormat ?? ImportTransactionDataMapping.DEFAULT_TIMEZONE_FORMAT,
                root.amountFormat ?? ImportTransactionDataMapping.DEFAULT_AMOUNT_FORMAT,
                root.geoLocationSeparator ?? ImportTransactionDataMapping.DEFAULT_GEO_LOCATION_SEPARATOR,
                root.geoLocationOrder ?? ImportTransactionDataMapping.DEFAULT_GEO_LOCATION_ORDER,
                root.tagSeparator ?? ImportTransactionDataMapping.DEFAULT_TAG_SEPARATOR
            );
        } catch {
            return null;
        }
    }
}

export type ImportTransactionReplaceRuleDataType = 'expenseCategory' | 'incomeCategory' | 'transferCategory' | 'account' | 'tag';

export class ImportTransactionReplaceRule {
    public dataType: ImportTransactionReplaceRuleDataType;
    public sourceValue: string;
    public targetId: string;

    private constructor(dataType: ImportTransactionReplaceRuleDataType, sourceValue: string, targetId: string) {
        this.dataType = dataType;
        this.sourceValue = sourceValue;
        this.targetId = targetId;
    }

    public toJsonObject(): unknown {
        return {
            type: this.dataType,
            sourceValue: this.sourceValue,
            targetId: this.targetId
        };
    }

    public static of(dataType: ImportTransactionReplaceRuleDataType, sourceValue: string, targetId: string): ImportTransactionReplaceRule {
        return new ImportTransactionReplaceRule(dataType, sourceValue, targetId);
    }

    public static parse(data: unknown): ImportTransactionReplaceRule | null {
        if (!data
            || typeof(data) !== 'object' || !('type' in data) || !('sourceValue' in data) || !('targetId' in data)
            || typeof(data.type) !== 'string' || typeof(data.sourceValue) !== 'string' || typeof(data.targetId) !== 'string') {
            return null;
        }

        if (data.type !== 'expenseCategory' && data.type !== 'incomeCategory' && data.type !== 'transferCategory' && data.type !== 'account' && data.type !== 'tag') {
            return null;
        }

        return new ImportTransactionReplaceRule(data.type as ImportTransactionReplaceRuleDataType, data.sourceValue as string, data.targetId as string);
    }
}

export class ImportTransactionReplaceRules {
    private static readonly JSON_ROOT_FIELD = 'ezBookkeepingImportTransactionReplaceRules';

    private readonly rules: ImportTransactionReplaceRule[];

    private constructor(rules: ImportTransactionReplaceRule[]) {
        this.rules = rules;
    }

    public getRules(): ImportTransactionReplaceRule[] {
        return this.rules;
    }

    public toJson(): string {
        const result: unknown[] = [];

        for (let i = 0; i < this.rules.length; i++) {
            const rule = this.rules[i];
            result.push(rule.toJsonObject());
        }

        return JSON.stringify({
            [ImportTransactionReplaceRules.JSON_ROOT_FIELD]: result
        });
    }

    public static of(rules: ImportTransactionReplaceRule[]): ImportTransactionReplaceRules {
        return new ImportTransactionReplaceRules(rules);
    }

    public static parseFromJson(json: string): ImportTransactionReplaceRules | null {
        try {
            const parsed = JSON.parse(json);
            const root = parsed[ImportTransactionReplaceRules.JSON_ROOT_FIELD];

            if (!root || !('length' in root)) {
                return null;
            }

            const result = new ImportTransactionReplaceRules([]);

            for (let i = 0; i < root.length; i++) {
                const rule = root[i];
                const replaceRule = ImportTransactionReplaceRule.parse(rule);

                if (replaceRule) {
                    result.rules.push(replaceRule);
                }
            }

            return result;
        } catch {
            return null;
        }
    }
}
