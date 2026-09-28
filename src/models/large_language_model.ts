export enum CodingAssistantType {
    InsightsExplorerCustomChart = 'insights_explorer_custom_chart'
}

export interface RecognizedTransactionResponse {
    readonly type: number;
    readonly time?: number;
    readonly categoryId?: string;
    readonly sourceAccountId?: string;
    readonly destinationAccountId?: string;
    readonly sourceAmount?: number;
    readonly destinationAmount?: number;
    readonly tagIds?: string[];
    readonly comment?: string;
}

export interface CodingAssistantRequest {
    readonly type: CodingAssistantType;
    readonly userPrompt: string;
    readonly code?: string;
}

export interface CodingAssistantResponse {
    readonly code: string;
}
