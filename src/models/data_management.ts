export interface ExportTransactionDataRequest {
    readonly maxTime: number;
    readonly minTime: number;
    readonly type: number;
    readonly categoryIds: string;
    readonly accountIds: string;
    readonly tagFilter: string;
    readonly amountFilter: string;
    readonly keyword: string;
}

export interface ClearDataRequest {
    readonly password: string;
}

export interface ClearAccountTransactionsRequest {
    readonly accountId: string;
    readonly password: string;
}

export interface DataStatisticsResponse {
    readonly totalAccountCount: string;
    readonly totalTransactionCategoryCount: string;
    readonly totalTransactionTagCount: string;
    readonly totalTransactionCount: string;
    readonly totalTransactionPictureCount: string;
    readonly totalInsightsExplorerCount: string;
    readonly totalTransactionTemplateCount: string;
    readonly totalScheduledTransactionCount: string;
}

export interface DisplayDataStatistics {
    readonly totalAccountCount: string;
    readonly totalTransactionCategoryCount: string;
    readonly totalTransactionTagCount: string;
    readonly totalTransactionCount: string;
    readonly totalTransactionPictureCount: string;
    readonly totalInsightsExplorerCount: string;
    readonly totalTransactionTemplateCount: string;
    readonly totalScheduledTransactionCount: string;
}
