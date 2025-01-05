import type { TransactionGeoLocationResponse } from './transaction.ts';

export interface ImportTransactionResponse {
    readonly type: number;
    readonly categoryId: string;
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
