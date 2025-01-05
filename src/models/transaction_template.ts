import type { TransactionInfoResponse } from './transaction.ts';

export interface TransactionTemplateCreateRequest {
    readonly templateType: number;
    readonly name: string;
    readonly type: number;
    readonly categoryId: string;
    readonly sourceAccountId: string;
    readonly destinationAccountId: string;
    readonly sourceAmount: number;
    readonly destinationAmount: number;
    readonly hideAmount: boolean;
    readonly tagIds: string[];
    readonly comment: string;
    readonly scheduledFrequencyType?: number;
    readonly scheduledFrequency?: string;
    readonly utcOffset?: number;
    readonly clientSessionId: string;
}

export interface TransactionTemplateModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly type: number;
    readonly categoryId: string;
    readonly sourceAccountId: string;
    readonly destinationAccountId: string;
    readonly sourceAmount: number;
    readonly destinationAmount: number;
    readonly hideAmount: boolean;
    readonly tagIds: string[];
    readonly comment: string;
    readonly scheduledFrequencyType?: number;
    readonly scheduledFrequency?: string;
    readonly utcOffset?: number;
}

export interface TransactionTemplateHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface TransactionTemplateMoveRequest {
    readonly newDisplayOrders: TransactionTemplateNewDisplayOrderRequest[];
}

export interface TransactionTemplateNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface TransactionTemplateDeleteRequest {
    readonly id: string;
}

export interface TransactionTemplateInfoResponse extends TransactionInfoResponse {
    readonly templateType: number;
    readonly name: string;
    readonly scheduledFrequencyType?: number;
    readonly scheduledFrequency?: string;
    readonly scheduledAt?: number;
    readonly displayOrder: number;
    readonly hidden: boolean;
}
