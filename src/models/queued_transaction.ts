import { TransactionType } from '@/core/transaction.ts';

import type { TransactionGeoLocationResponse } from './transaction.ts';
import type { TransactionCategoryInfoResponse } from './transaction_category.ts';

export enum QueuedTransactionStatus {
    Pending = 1,
    Archived = 2
}

export enum QueuedTransactionSource {
    Unknown = 0,
    ScheduledTemplate = 1,
    Api = 2,
    Mcp = 3
}

/**
 * A draft transaction waiting in the queue of the user.
 *
 * It is not part of the bookkeeping data: it affects no account balance and shows up in no
 * statistic. It only becomes real once the user accepts it, which creates a regular transaction out
 * of it. Because it is a draft, every transaction field may be unset.
 */
export class QueuedTransaction implements QueuedTransactionInfoResponse {
    public readonly id: string;
    public readonly status: number;
    public readonly source: number;
    public readonly sourceTemplateId?: string;
    public readonly type: number;
    public readonly categoryId: string;
    public readonly category?: TransactionCategoryInfoResponse;
    public readonly time: number;
    public readonly utcOffset: number;
    public readonly sourceAccountId: string;
    public readonly destinationAccountId: string;
    public readonly sourceAmount: number;
    public readonly destinationAmount?: number;
    public readonly hideAmount: boolean;
    public readonly tagIds: string[];
    public readonly comment: string;
    public readonly geoLocation?: TransactionGeoLocationResponse;
    public readonly createdAt: number;
    public readonly archivedAt?: number;
    public readonly expiresAt?: number;

    private constructor(response: QueuedTransactionInfoResponse) {
        this.id = response.id;
        this.status = response.status;
        this.source = response.source;
        this.sourceTemplateId = response.sourceTemplateId;
        this.type = response.type;
        this.categoryId = response.categoryId;
        this.category = response.category;
        this.time = response.time;
        this.utcOffset = response.utcOffset;
        this.sourceAccountId = response.sourceAccountId;
        this.destinationAccountId = response.destinationAccountId;
        this.sourceAmount = response.sourceAmount;
        this.destinationAmount = response.destinationAmount;
        this.hideAmount = response.hideAmount;
        this.tagIds = response.tagIds ?? [];
        this.comment = response.comment;
        this.geoLocation = response.geoLocation;
        this.createdAt = response.createdAt;
        this.archivedAt = response.archivedAt;
        this.expiresAt = response.expiresAt;
    }

    public get isArchived(): boolean {
        return this.status === QueuedTransactionStatus.Archived;
    }

    public get isTransfer(): boolean {
        return this.type === TransactionType.Transfer;
    }

    /** Whether the draft holds everything a real transaction needs, so the user can tell at a glance what is still missing */
    public get isComplete(): boolean {
        if (!this.type || !this.time || !this.sourceAccountId || this.sourceAccountId === '0') {
            return false;
        }

        if (this.type !== TransactionType.ModifyBalance && (!this.categoryId || this.categoryId === '0')) {
            return false;
        }

        if (this.isTransfer && (!this.destinationAccountId || this.destinationAccountId === '0')) {
            return false;
        }

        return true;
    }

    public static of(response: QueuedTransactionInfoResponse): QueuedTransaction {
        return new QueuedTransaction(response);
    }

    public static ofMulti(responses: QueuedTransactionInfoResponse[]): QueuedTransaction[] {
        const queuedTransactions: QueuedTransaction[] = [];

        for (const response of responses) {
            queuedTransactions.push(QueuedTransaction.of(response));
        }

        return queuedTransactions;
    }
}

export interface QueuedTransactionListRequest {
    readonly includeArchived: boolean;
}

export interface QueuedTransactionCreateRequest {
    readonly type?: number;
    readonly categoryId?: string;
    readonly time?: number;
    readonly utcOffset?: number;
    readonly sourceAccountId?: string;
    readonly destinationAccountId?: string;
    readonly sourceAmount?: number;
    readonly destinationAmount?: number;
    readonly hideAmount?: boolean;
    readonly tagIds?: string[];
    readonly comment?: string;
    readonly clientSessionId?: string;
}

export interface QueuedTransactionArchiveRequest {
    readonly id: string;
}

export interface QueuedTransactionRestoreRequest {
    readonly id: string;
}

export interface QueuedTransactionDeleteRequest {
    readonly id: string;
}

export interface QueuedTransactionInfoResponse {
    readonly id: string;
    readonly status: number;
    readonly source: number;
    readonly sourceTemplateId?: string;
    readonly type: number;
    readonly categoryId: string;
    readonly category?: TransactionCategoryInfoResponse;
    readonly time: number;
    readonly utcOffset: number;
    readonly sourceAccountId: string;
    readonly destinationAccountId: string;
    readonly sourceAmount: number;
    readonly destinationAmount?: number;
    readonly hideAmount: boolean;
    readonly tagIds: string[];
    readonly comment: string;
    readonly geoLocation?: TransactionGeoLocationResponse;
    readonly createdAt: number;
    readonly archivedAt?: number;
    readonly expiresAt?: number;
}

export interface QueuedTransactionCountResponse {
    readonly pendingCount: number;
    readonly archivedCount: number;
}

export interface QueuedTransactionBatchOperationResponse {
    readonly affectedCount: number;
}
