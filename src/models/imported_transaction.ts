import { TransactionType } from '@/core/transaction.ts';

import type { TransactionCreateRequest, TransactionGeoLocationResponse } from './transaction.ts';

// 组合键分隔符：用于在批量创建对话框中唯一标识“一级分类名 + 二级分类名”或“账户名 + 币种”。
// 选用单元分隔符 U+001F（ASCII 0x1F），它不会出现在用户填写的分类名 / 账户名 / 币种中，可避免歧义与碰撞。
const IMPORT_COMPOSITE_KEY_SEPARATOR = '\u001f';

// 构造"一级分类名 + 二级分类名"的组合键
// primaryCategoryName: 一级分类名（可能为空字符串）
// subCategoryName: 二级分类名
export function getImportCategoryCompositeKey(primaryCategoryName: string, subCategoryName: string): string {
    return `${primaryCategoryName}${IMPORT_COMPOSITE_KEY_SEPARATOR}${subCategoryName}`;
}

// 解析"一级分类名 + 二级分类名"的组合键
// key: 由 getImportCategoryCompositeKey 生成的组合键
export function parseImportCategoryCompositeKey(key: string): { primaryCategoryName: string, subCategoryName: string } {
    const separatorIndex = key.indexOf(IMPORT_COMPOSITE_KEY_SEPARATOR);

    if (separatorIndex < 0) {
        return { primaryCategoryName: '', subCategoryName: key };
    }

    return {
        primaryCategoryName: key.substring(0, separatorIndex),
        subCategoryName: key.substring(separatorIndex + IMPORT_COMPOSITE_KEY_SEPARATOR.length)
    };
}

// 构造"账户名 + 币种"的组合键
// accountName: 账户名
// accountCurrency: 账户币种（可能为空字符串）
export function getImportAccountCompositeKey(accountName: string, accountCurrency: string): string {
    return `${accountName}${IMPORT_COMPOSITE_KEY_SEPARATOR}${accountCurrency}`;
}

// 解析"账户名 + 币种"的组合键
// key: 由 getImportAccountCompositeKey 生成的组合键
export function parseImportAccountCompositeKey(key: string): { accountName: string, accountCurrency: string } {
    const separatorIndex = key.indexOf(IMPORT_COMPOSITE_KEY_SEPARATOR);

    if (separatorIndex < 0) {
        return { accountName: key, accountCurrency: '' };
    }

    return {
        accountName: key.substring(0, separatorIndex),
        accountCurrency: key.substring(separatorIndex + IMPORT_COMPOSITE_KEY_SEPARATOR.length)
    };
}

export class ImportTransaction implements ImportTransactionResponse {
    public type: number;
    public categoryId: string;
    public originalPrimaryCategoryName: string;
    public originalCategoryName: string;
    public time: number;
    public utcOffset: number;
    public sourceAccountId: string;
    public originalSourceAccountName: string;
    public originalSourceAccountCurrency: string;
    public destinationAccountId: string;
    public originalDestinationAccountName?: string;
    public originalDestinationAccountCurrency?: string;
    public sourceAmount: number;
    public destinationAmount: number;
    public tagIds: string[];
    public originalTagNames: string[];
    public comment: string;
    public geoLocation?: TransactionGeoLocationResponse;

    public actualCategoryName: string;
    public actualSourceAccountName: string;
    public actualDestinationAccountName?: string;
    public index: number;
    public selected: boolean;
    public valid: boolean;

    private constructor(response: ImportTransactionResponse, index: number) {
        this.type = response.type;
        this.categoryId = response.categoryId;
        this.originalPrimaryCategoryName = response.originalPrimaryCategoryName || '';
        this.originalCategoryName = response.originalCategoryName;
        this.time = response.time;
        this.utcOffset = response.utcOffset;
        this.sourceAccountId = response.sourceAccountId;
        this.originalSourceAccountName = response.originalSourceAccountName;
        this.originalSourceAccountCurrency = response.originalSourceAccountCurrency;
        this.destinationAccountId = response.destinationAccountId || '';
        this.originalDestinationAccountName = response.originalDestinationAccountName;
        this.originalDestinationAccountCurrency = response.originalDestinationAccountCurrency;
        this.sourceAmount = response.sourceAmount;
        this.destinationAmount = response.destinationAmount || 0;
        this.tagIds = response.tagIds || [];
        this.originalTagNames = response.originalTagNames || [];
        this.comment = response.comment;
        this.geoLocation = response.geoLocation;

        this.actualCategoryName = response.originalCategoryName;
        this.actualSourceAccountName = response.originalSourceAccountName;
        this.actualDestinationAccountName = response.originalDestinationAccountName;
        this.index = index;
        this.selected = false;
        this.valid = this.isTransactionValid();
    }

    public toCreateRequest(): TransactionCreateRequest {
        return {
            type: this.type,
            categoryId: this.categoryId,
            time: this.time,
            utcOffset: this.utcOffset,
            sourceAccountId: this.sourceAccountId,
            destinationAccountId: this.type === TransactionType.Transfer ? this.destinationAccountId : '0',
            sourceAmount: this.sourceAmount,
            destinationAmount: this.type === TransactionType.Transfer ? this.destinationAmount : 0,
            hideAmount: false,
            tagIds: this.tagIds,
            pictureIds: [],
            comment: this.comment,
            geoLocation: this.geoLocation,
            clientSessionId: ''
        };
    }

    public isTransactionValid(): boolean {
        if (this.type !== TransactionType.ModifyBalance && (!this.categoryId || this.categoryId === '0')) {
            return false;
        }

        if (!this.sourceAccountId || this.sourceAccountId === '0') {
            return false;
        }

        if (this.type === TransactionType.Transfer && (!this.destinationAccountId || this.destinationAccountId === '0')) {
            return false;
        }

        if (this.tagIds && this.tagIds.length) {
            for (const tagId of this.tagIds) {
                if (!tagId || tagId === '0') {
                    return false;
                }
            }
        }

        return true;
    }

    public static of(response: ImportTransactionResponse, index: number): ImportTransaction {
        return new ImportTransaction(response, index);
    }
}

export interface ImportTransactionRequest {
    readonly transactions: ImportTransactionRequestItem[];
}

export interface ImportTransactionRequestItem {
    readonly time: string;
    readonly utcOffset: string;
    readonly type: string;
    readonly categoryName?: string;
    readonly sourceAccountName?: string;
    readonly destinationAccountName?: string;
    readonly sourceAmount: string;
    readonly destinationAmount?: string;
    readonly geoLocation?: string;
    readonly tagNames?: string;
    readonly comment?: string;
}

export interface ImportTransactionResponse {
    readonly type: number;
    readonly categoryId: string;
    readonly originalPrimaryCategoryName?: string;
    readonly originalCategoryName: string;
    readonly time: number;
    readonly utcOffset: number;
    readonly sourceAccountId: string;
    readonly originalSourceAccountName: string;
    readonly originalSourceAccountCurrency: string;
    readonly destinationAccountId?: string;
    readonly originalDestinationAccountName?: string;
    readonly originalDestinationAccountCurrency?: string;
    readonly sourceAmount: number;
    readonly destinationAmount?: number;
    readonly tagIds: string[];
    readonly originalTagNames: string[];
    readonly comment: string;
    readonly geoLocation?: TransactionGeoLocationResponse;
}

export interface ImportTransactionResponsePageWrapper {
    readonly items: ImportTransactionResponse[];
    readonly totalCount: number;
}
