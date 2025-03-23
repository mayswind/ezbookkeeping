import { TransactionType } from '@/core/transaction.ts';
import { TemplateType } from '@/core/template.ts';

import { Transaction, type TransactionInfoResponse } from './transaction.ts';

export class TransactionTemplate extends Transaction implements TransactionTemplateInfoResponse {
    public templateType: number;
    public name: string;
    public scheduledFrequencyType?: number;
    public scheduledFrequency?: string;
    public scheduledStartDate?: string;
    public scheduledEndDate?: string;
    public scheduledAt?: number;
    public displayOrder: number;
    public hidden: boolean;

    private constructor(id: string, templateType: number, name: string, type: number, categoryId: string, timeZone: string | undefined, utcOffset: number, sourceAccountId: string, destinationAccountId: string, sourceAmount: number, destinationAmount: number, hideAmount: boolean, scheduledFrequencyType: number | undefined, scheduledFrequency: string | undefined, scheduledStartDate: string | undefined, scheduledEndDate: string | undefined, scheduledAt: number | undefined, tagIds: string[], comment: string, editable: boolean, displayOrder: number, hidden: boolean) {
        super(id, '', type, categoryId, 0, timeZone, utcOffset, sourceAccountId, destinationAccountId, sourceAmount, destinationAmount, hideAmount, tagIds, comment, editable);
        this.templateType = templateType;
        this.name = name;
        this.scheduledFrequencyType = scheduledFrequencyType;
        this.scheduledFrequency = scheduledFrequency;
        this.scheduledStartDate = scheduledStartDate;
        this.scheduledEndDate = scheduledEndDate;
        this.scheduledAt = scheduledAt;
        this.displayOrder = displayOrder;
        this.hidden = hidden;
    }

    public fillFrom(other: TransactionTemplate): void {
        this.templateType = other.templateType;
        this.name = other.name;

        if (this.templateType === TemplateType.Schedule.type) {
            this.scheduledFrequencyType = other.scheduledFrequencyType;
            this.scheduledFrequency = other.scheduledFrequency;
            this.scheduledStartDate = other.scheduledStartDate;
            this.scheduledEndDate = other.scheduledEndDate;
            this.utcOffset = other.utcOffset;
            this.timeZone = undefined;
        }
    }

    public toTemplateCreateRequest(clientSessionId: string): TransactionTemplateCreateRequest {
        return {
            templateType: this.templateType,
            name: this.name,
            type: this.type,
            categoryId: this.getCategoryId(),
            sourceAccountId: this.sourceAccountId,
            destinationAccountId: this.type === TransactionType.Transfer ? this.destinationAccountId : '0',
            sourceAmount: this.sourceAmount,
            destinationAmount: this.type === TransactionType.Transfer ? this.destinationAmount : 0,
            hideAmount: this.hideAmount,
            tagIds: this.tagIds,
            comment: this.comment,
            scheduledFrequencyType: this.templateType === TemplateType.Schedule.type ? this.scheduledFrequencyType : undefined,
            scheduledFrequency: this.templateType === TemplateType.Schedule.type ? this.scheduledFrequency : undefined,
            scheduledStartDate: this.templateType === TemplateType.Schedule.type && this.scheduledStartDate ? this.scheduledStartDate : undefined,
            scheduledEndDate: this.templateType === TemplateType.Schedule.type && this.scheduledEndDate ? this.scheduledEndDate : undefined,
            utcOffset: this.templateType === TemplateType.Schedule.type ? this.utcOffset : undefined,
            clientSessionId: clientSessionId
        };
    }

    public toTemplateModifyRequest(): TransactionTemplateModifyRequest {
        return {
            id: this.id,
            name: this.name,
            type: this.type,
            categoryId: this.getCategoryId(),
            sourceAccountId: this.sourceAccountId,
            destinationAccountId: this.type === TransactionType.Transfer ? this.destinationAccountId : '0',
            sourceAmount: this.sourceAmount,
            destinationAmount: this.type === TransactionType.Transfer ? this.destinationAmount : 0,
            hideAmount: this.hideAmount,
            tagIds: this.tagIds,
            comment: this.comment,
            scheduledFrequencyType: this.templateType === TemplateType.Schedule.type ? this.scheduledFrequencyType : undefined,
            scheduledFrequency: this.templateType === TemplateType.Schedule.type ? this.scheduledFrequency : undefined,
            scheduledStartDate: this.templateType === TemplateType.Schedule.type && this.scheduledStartDate ? this.scheduledStartDate : undefined,
            scheduledEndDate: this.templateType === TemplateType.Schedule.type && this.scheduledEndDate ? this.scheduledEndDate : undefined,
            utcOffset: this.templateType === TemplateType.Schedule.type ? this.utcOffset : undefined
        };
    }

    public static createNewTransactionTemplate(transaction: Transaction): TransactionTemplate {
        return new TransactionTemplate(
            transaction.id,
            0, // templateType
            '', // name
            transaction.type,
            transaction.categoryId,
            transaction.timeZone,
            transaction.utcOffset,
            transaction.sourceAccountId,
            transaction.destinationAccountId,
            transaction.sourceAmount,
            transaction.destinationAmount,
            transaction.hideAmount,
            undefined, // scheduledFrequencyType
            undefined, // scheduledFrequency
            undefined, // scheduledStartDate
            undefined, // scheduledEndDate
            undefined, // scheduledAt
            transaction.tagIds,
            transaction.comment,
            true,
            0,
            false
        );
    }

    public static ofTemplate(templateResponse: TransactionTemplateInfoResponse): TransactionTemplate {
        return new TransactionTemplate(
            templateResponse.id,
            templateResponse.templateType,
            templateResponse.name,
            templateResponse.type,
            templateResponse.categoryId,
            undefined, // only in new transaction template
            templateResponse.utcOffset ?? 0,
            templateResponse.sourceAccountId,
            templateResponse.destinationAccountId,
            templateResponse.sourceAmount,
            templateResponse.destinationAmount,
            templateResponse.hideAmount,
            templateResponse.scheduledFrequencyType,
            templateResponse.scheduledFrequency,
            templateResponse.scheduledStartDate ?? undefined,
            templateResponse.scheduledEndDate ?? undefined,
            templateResponse.scheduledAt,
            templateResponse.tagIds,
            templateResponse.comment,
            true, // editable
            templateResponse.displayOrder,
            templateResponse.hidden
        );
    }

    public static ofManyTemplates(templateResponses: TransactionTemplateInfoResponse[]): TransactionTemplate[] {
        const templates: TransactionTemplate[] = [];

        for (const templateResponse of templateResponses) {
            templates.push(TransactionTemplate.ofTemplate(templateResponse));
        }

        return templates;
    }
}

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
    readonly scheduledStartDate?: string;
    readonly scheduledEndDate?: string;
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
    readonly scheduledStartDate?: string;
    readonly scheduledEndDate?: string;
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
    readonly scheduledStartDate?: string;
    readonly scheduledEndDate?: string;
    readonly scheduledAt?: number;
    readonly displayOrder: number;
    readonly hidden: boolean;
}
