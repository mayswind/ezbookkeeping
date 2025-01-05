export interface ClearDataRequest {
    readonly password: string;
}

export interface DataStatisticsResponse {
    readonly totalAccountCount: number;
    readonly totalTransactionCategoryCount: number;
    readonly totalTransactionTagCount: number;
    readonly totalTransactionCount: number;
    readonly totalTransactionPictureCount: number;
    readonly totalTransactionTemplateCount: number;
    readonly totalScheduledTransactionCount: number;
}
