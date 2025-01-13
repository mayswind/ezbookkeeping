export interface ClearDataRequest {
    readonly password: string;
}

export interface DataStatisticsResponse {
    readonly totalAccountCount: string;
    readonly totalTransactionCategoryCount: string;
    readonly totalTransactionTagCount: string;
    readonly totalTransactionCount: string;
    readonly totalTransactionPictureCount: string;
    readonly totalTransactionTemplateCount: string;
    readonly totalScheduledTransactionCount: string;
}

export interface DisplayDataStatistics {
    readonly totalAccountCount: string;
    readonly totalTransactionCategoryCount: string;
    readonly totalTransactionTagCount: string;
    readonly totalTransactionCount: string;
    readonly totalTransactionPictureCount: string;
    readonly totalTransactionTemplateCount: string;
    readonly totalScheduledTransactionCount: string;
}
