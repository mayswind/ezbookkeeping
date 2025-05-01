import type { PartialRecord } from '@/core/base.ts';
import type { YearMonth, StartEndTime } from '@/core/datetime.ts';
import type { MapPosition } from '@/core/map.ts';
import { TransactionType } from '@/core/transaction.ts';

import { Account, type AccountInfoResponse } from './account.ts';
import { TransactionCategory, type TransactionCategoryInfoResponse } from './transaction_category.ts';
import { TransactionTag, type TransactionTagInfoResponse } from './transaction_tag.ts';
import { TransactionPicture, type TransactionPictureInfoBasicResponse } from './transaction_picture_info.ts';

export class Transaction implements TransactionInfoResponse {
    public id: string;
    public timeSequenceId: string;
    public type: number;
    public expenseCategoryId: string = '';
    public incomeCategoryId: string = '';
    public transferCategoryId: string = '';
    public time: number;
    public timeZone?: string; // only in new transaction
    public utcOffset: number;
    public sourceAccountId: string;
    public destinationAccountId: string;
    public sourceAmount: number;
    public destinationAmount: number;
    public hideAmount: boolean;
    public tagIds: string[];
    public comment: string;
    public editable: boolean;

    private _pictures?: TransactionPicture[];
    private _geoLocation?: TransactionGeoLocation;

    private _category?: TransactionCategory; // only for displaying transaction
    private _sourceAccount?: Account; // only for displaying transaction
    private _destinationAccount?: Account; // only for displaying transaction
    private _tags?: TransactionTag[]; // only for displaying transaction

    private _date?: string = undefined; // only for displaying transaction in transaction list
    private _day?: number = undefined; // only for displaying transaction in transaction list
    private _dayOfWeek?: string = undefined; // only for displaying transaction in transaction list

    protected constructor(id: string, timeSequenceId: string, type: number, categoryId: string, time: number, timeZone: string | undefined, utcOffset: number, sourceAccountId: string, destinationAccountId: string, sourceAmount: number, destinationAmount: number, hideAmount: boolean, tagIds: string[], comment: string, editable: boolean) {
        this.id = id;
        this.timeSequenceId = timeSequenceId;
        this.type = type;
        this.time = time;
        this.timeZone = timeZone;
        this.utcOffset = utcOffset;
        this.sourceAccountId = sourceAccountId;
        this.destinationAccountId = destinationAccountId;
        this.sourceAmount = sourceAmount;
        this.destinationAmount = destinationAmount;
        this.hideAmount = hideAmount;
        this.tagIds = tagIds;
        this.comment = comment;
        this.editable = editable;
        this.setCategoryId(categoryId);
    }

    public get pictures(): TransactionPictureInfoBasicResponse[] | undefined {
        const ret: TransactionPictureInfoBasicResponse[] = [];

        if (this._pictures) {
            for (const picture of this._pictures) {
                ret.push(picture);
            }
        }

        return ret;
    }

    public get geoLocation(): TransactionGeoLocationResponse | undefined {
        return this._geoLocation;
    }


    public set geoLocation(value: MapPosition) {
        this._geoLocation = TransactionGeoLocation.of(value);
    }

    public get categoryId(): string {
        return this.getCategoryId();
    }

    public get category(): TransactionCategoryInfoResponse | undefined {
        return this._category;
    }

    public get sourceAccount(): AccountInfoResponse | undefined {
        return this._sourceAccount;
    }

    public get destinationAccount(): AccountInfoResponse | undefined {
        return this._destinationAccount;
    }

    public get tags(): TransactionTagInfoResponse[] | undefined {
        const ret: TransactionTagInfoResponse[] = [];

        if (this._tags) {
            for (const tag of this._tags) {
                ret.push(tag);
            }
        }

        return ret;
    }

    public get date(): string | undefined {
        return this._date;
    }

    public get day(): number | undefined {
        return this._day;
    }

    public get dayOfWeek(): string | undefined {
        return this._dayOfWeek;
    }

    public getCategoryId(): string {
        if (this.type === TransactionType.Expense) {
            return this.expenseCategoryId;
        } else if (this.type === TransactionType.Income) {
            return this.incomeCategoryId;
        } else if (this.type === TransactionType.Transfer) {
            return this.transferCategoryId;
        } else {
            return '';
        }
    }

    public setCategoryId(categoryId: string): void {
        if (this.type === TransactionType.Expense) {
            this.expenseCategoryId = categoryId;
        } else if (this.type === TransactionType.Income) {
            this.incomeCategoryId = categoryId;
        } else if (this.type === TransactionType.Transfer) {
            this.transferCategoryId = categoryId;
        }
    }

    public setCategory(category: TransactionCategory): void {
        this._category = category;
    }

    public setSourceAccount(sourceAccount: Account): void {
        this._sourceAccount = sourceAccount;
    }

    public setDestinationAccount(destinationAccount: Account): void {
        this._destinationAccount = destinationAccount;
    }

    public setTags(tags: TransactionTag[]): void {
        this._tags = tags;
    }

    public getPictureIds(): string[] {
        const pictureIds: string[] = [];

        if (this._pictures) {
            for (const picture of this._pictures) {
                pictureIds.push(picture.pictureId);
            }
        }

        return pictureIds;
    }

    public setPictures(pictures: TransactionPicture[]): void {
        this._pictures = pictures;
    }

    public addPicture(pictureInfo: TransactionPictureInfoBasicResponse): void {
        if (!this._pictures) {
            this._pictures = [];
        }

        this._pictures.push(TransactionPicture.of(pictureInfo));
    }

    public removePicture(pictureInfo: TransactionPictureInfoBasicResponse): void {
        if (!this._pictures) {
            return;
        }

        for (let i = 0; i < this._pictures.length; i++) {
            if (this._pictures[i].pictureId === pictureInfo.pictureId) {
                this._pictures.splice(i, 1);
            }
        }
    }

    public clearPictures(): void {
        this._pictures = [];
    }

    public setGeoLocation(geoLocation?: TransactionGeoLocation): void {
        this._geoLocation = geoLocation;
    }

    public setLatitudeAndLongitude(latitude: number, longitude: number): void {
        this._geoLocation = TransactionGeoLocation.createNewGeoLocation(latitude, longitude);
    }

    public removeGeoLocation(): void {
        this._geoLocation = undefined;
    }

    public setDisplayDate(date: string, day: number, dayOfWeek: string): void {
        this._date = date;
        this._day = day;
        this._dayOfWeek = dayOfWeek;
    }

    public toCreateRequest(clientSessionId: string, actualTime?: number): TransactionCreateRequest {
        return {
            type: this.type,
            categoryId: this.getCategoryId(),
            time: actualTime ? actualTime : this.time,
            utcOffset: this.utcOffset,
            sourceAccountId: this.sourceAccountId,
            destinationAccountId: this.type === TransactionType.Transfer ? this.destinationAccountId : '0',
            sourceAmount: this.sourceAmount,
            destinationAmount: this.type === TransactionType.Transfer ? this.destinationAmount : 0,
            hideAmount: this.hideAmount,
            tagIds: this.tagIds,
            pictureIds: this.getPictureIds(),
            comment: this.comment,
            geoLocation: this.geoLocation,
            clientSessionId: clientSessionId
        };
    }

    public toModifyRequest(actualTime?: number): TransactionModifyRequest {
        return {
            id: this.id,
            categoryId: this.getCategoryId(),
            time: actualTime ? actualTime : this.time,
            utcOffset: this.utcOffset,
            sourceAccountId: this.sourceAccountId,
            destinationAccountId: this.type === TransactionType.Transfer ? this.destinationAccountId : '0',
            sourceAmount: this.sourceAmount,
            destinationAmount: this.type === TransactionType.Transfer ? this.destinationAmount : 0,
            hideAmount: this.hideAmount,
            tagIds: this.tagIds,
            pictureIds: this.getPictureIds(),
            comment: this.comment,
            geoLocation: this.geoLocation
        };
    }

    public toTransactionDraft(): TransactionDraft | null {
        if (this.type !== TransactionType.Expense &&
            this.type !== TransactionType.Income &&
            this.type !== TransactionType.Transfer) {
            return null;
        }

        return {
            type: this.type,
            categoryId: this.getCategoryId(),
            sourceAccountId: this.sourceAccountId,
            sourceAmount: this.sourceAmount,
            destinationAccountId: this.type === TransactionType.Transfer ? this.destinationAccountId : '0',
            destinationAmount: this.type === TransactionType.Transfer ? this.destinationAmount : 0,
            hideAmount: this.hideAmount,
            tagIds: this.tagIds,
            pictures: this.pictures,
            comment: this.comment,
        };
    }

    public static createNewTransaction(type: number, time: number, timeZone: string, utcOffset: number): Transaction {
        return new Transaction(
            '', // id
            '', // timeSequenceId
            type, // type
            '', // categoryId
            time, // time
            timeZone, // timeZone
            utcOffset, // utcOffset
            '', // sourceAccountId
            '', // destinationAccountId
            0, // sourceAmount
            0, // destinationAmount
            false, // hideAmount
            [], // tagIds
            '', // comment
            true // editable
        );
    }

    public static of(transactionResponse: TransactionInfoResponse): Transaction {
        const transaction: Transaction = new Transaction(
            transactionResponse.id,
            transactionResponse.timeSequenceId,
            transactionResponse.type,
            transactionResponse.categoryId,
            transactionResponse.time,
            undefined, // only in new transaction
            transactionResponse.utcOffset,
            transactionResponse.sourceAccountId,
            transactionResponse.destinationAccountId,
            transactionResponse.sourceAmount,
            transactionResponse.destinationAmount,
            transactionResponse.hideAmount,
            transactionResponse.tagIds,
            transactionResponse.comment,
            transactionResponse.editable
        );

        if (transactionResponse.category) {
            transaction.setCategory(TransactionCategory.of(transactionResponse.category));
        }

        if (transactionResponse.sourceAccount) {
            transaction.setSourceAccount(Account.of(transactionResponse.sourceAccount));
        }

        if (transactionResponse.destinationAccount) {
            transaction.setDestinationAccount(Account.of(transactionResponse.destinationAccount));
        }

        if (transactionResponse.tags) {
            transaction.setTags(TransactionTag.ofMulti(transactionResponse.tags));
        }

        if (transactionResponse.pictures) {
            const pictures: TransactionPicture[] = [];

            for (const picture of transactionResponse.pictures) {
                pictures.push(TransactionPicture.of(picture));
            }

            transaction.setPictures(pictures);
        }

        if (transactionResponse.geoLocation) {
            transaction.setLatitudeAndLongitude(transactionResponse.geoLocation.latitude, transactionResponse.geoLocation.longitude);
        }

        return transaction;
    }

    public static ofMulti(transactionResponses: TransactionInfoResponse[]): Transaction[] {
        const transactions: Transaction[] = [];

        for (const transactionResponse of transactionResponses) {
            transactions.push(Transaction.of(transactionResponse));
        }

        return transactions;
    }

    public static ofDraft(transactionDraft?: TransactionDraft | null): Transaction | null {
        if (!transactionDraft) {
            return null;
        }

        if (transactionDraft.type !== TransactionType.Expense &&
            transactionDraft.type !== TransactionType.Income &&
            transactionDraft.type !== TransactionType.Transfer) {
            return null;
        }

        const transaction: Transaction = new Transaction(
            '', // id
            '', // timeSequenceId
            transactionDraft.type, // type
            transactionDraft.categoryId ?? '', // categoryId
            0, // time
            undefined, // only in new transaction
            0, // utcOffset
            transactionDraft.sourceAccountId ?? '', // sourceAccountId
            transactionDraft.destinationAccountId ?? '', // destinationAccountId
            transactionDraft.sourceAmount ?? 0, // sourceAmount
            transactionDraft.destinationAmount ?? 0, // destinationAmount
            transactionDraft.hideAmount ?? false, // hideAmount
            transactionDraft.tagIds ?? [], // tagIds
            transactionDraft.comment ?? '', // comment
            true // editable
        );

        if (transactionDraft.pictures) {
            const pictures: TransactionPicture[] = [];

            for (const picture of transactionDraft.pictures) {
                pictures.push(TransactionPicture.of(picture));
            }

            transaction.setPictures(pictures);
        }

        return transaction;
    }
}

export class TransactionGeoLocation implements TransactionGeoLocationRequest {
    public latitude: number;
    public longitude: number;

    private constructor(latitude: number, longitude: number) {
        this.latitude = latitude;
        this.longitude = longitude;
    }

    public static createNewGeoLocation(latitude: number, longitude: number): TransactionGeoLocation {
        return new TransactionGeoLocation(latitude, longitude);
    }

    public static of(mapPosition: MapPosition): TransactionGeoLocation {
        return new TransactionGeoLocation(mapPosition.latitude, mapPosition.longitude);
    }
}

export interface TransactionDraft {
    readonly type?: number;
    readonly categoryId?: string;
    readonly sourceAccountId?: string;
    readonly sourceAmount?: number;
    readonly destinationAccountId?: string;
    readonly destinationAmount?: number;
    readonly hideAmount?: boolean;
    readonly tagIds?: string[];
    readonly pictures?: TransactionPictureInfoBasicResponse[];
    readonly comment?: string;
}

export interface TransactionGeoLocationRequest {
    readonly latitude: number;
    readonly longitude: number;
}

export interface TransactionCreateRequest {
    readonly type: number;
    readonly categoryId: string;
    readonly time: number;
    readonly utcOffset: number;
    readonly sourceAccountId: string;
    readonly destinationAccountId: string;
    readonly sourceAmount: number;
    readonly destinationAmount: number;
    readonly hideAmount: boolean;
    readonly tagIds: string[];
    readonly pictureIds: string[];
    readonly comment: string;
    readonly geoLocation?: TransactionGeoLocationRequest;
    readonly clientSessionId: string;
}

export interface TransactionModifyRequest {
    readonly id: string;
    readonly categoryId: string;
    readonly time: number;
    readonly utcOffset: number;
    readonly sourceAccountId: string;
    readonly destinationAccountId: string;
    readonly sourceAmount: number;
    readonly destinationAmount: number;
    readonly hideAmount: boolean;
    readonly tagIds: string[];
    readonly pictureIds: string[];
    readonly comment: string;
    readonly geoLocation?: TransactionGeoLocationRequest;
}

export interface TransactionDeleteRequest {
    readonly id: string;
}

export interface TransactionImportRequest {
    readonly transactions: TransactionCreateRequest[];
    readonly clientSessionId: string;
}

export interface TransactionListByMaxTimeRequest {
    readonly maxTime: number;
    readonly minTime: number;
    readonly count: number;
    readonly page: number;
    readonly withCount: boolean;
    readonly type: number;
    readonly categoryIds: string;
    readonly accountIds: string;
    readonly tagIds: string;
    readonly tagFilterType: number;
    readonly amountFilter: string;
    readonly keyword: string;
}

export interface TransactionListInMonthByPageRequest {
    readonly year: number;
    readonly month: number;
    readonly type: number;
    readonly categoryIds: string;
    readonly accountIds: string;
    readonly tagIds: string;
    readonly tagFilterType: number;
    readonly amountFilter: string;
    readonly keyword: string;
}

export type TransactionGeoLocationResponse = MapPosition;

export interface TransactionInfoResponse {
    readonly id: string;
    readonly timeSequenceId: string;
    readonly type: number;
    readonly categoryId: string;
    readonly category?: TransactionCategoryInfoResponse;
    readonly time: number;
    readonly utcOffset: number;
    readonly sourceAccountId: string;
    readonly sourceAccount?: AccountInfoResponse;
    readonly destinationAccountId: string;
    readonly destinationAccount?: AccountInfoResponse;
    readonly sourceAmount: number;
    readonly destinationAmount: number;
    readonly hideAmount: boolean;
    readonly tagIds: string[];
    readonly tags?: TransactionTagInfoResponse[];
    readonly pictures?: TransactionPictureInfoBasicResponse[];
    readonly comment: string;
    readonly geoLocation?: TransactionGeoLocationResponse;
    readonly editable: boolean;
}

export interface TransactionStatisticRequest {
    readonly startTime: number;
    readonly endTime: number;
    readonly tagIds: string;
    readonly tagFilterType: number;
    readonly useTransactionTimezone: boolean;
}

export interface YearMonthRangeRequest {
    readonly startYearMonth: string;
    readonly endYearMonth: string;
}

export interface TransactionStatisticTrendsRequest extends YearMonthRangeRequest {
    readonly tagIds: string;
    readonly tagFilterType: number;
    readonly useTransactionTimezone: boolean;
}

export const ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE = [
    'today',
    'thisWeek',
    'thisMonth',
    'thisYear',
    'lastMonth',
    'monthBeforeLastMonth',
    'monthBeforeLast2Months',
    'monthBeforeLast3Months',
    'monthBeforeLast4Months',
    'monthBeforeLast5Months',
    'monthBeforeLast6Months',
    'monthBeforeLast7Months',
    'monthBeforeLast8Months',
    'monthBeforeLast9Months',
    'monthBeforeLast10Months'
] as const;

export type TransactionAmountsRequestType = typeof ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE[number];

export const LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES: TransactionAmountsRequestType[] = [
    'monthBeforeLast10Months',
    'monthBeforeLast9Months',
    'monthBeforeLast8Months',
    'monthBeforeLast7Months',
    'monthBeforeLast6Months',
    'monthBeforeLast5Months',
    'monthBeforeLast4Months',
    'monthBeforeLast3Months',
    'monthBeforeLast2Months',
    'monthBeforeLastMonth',
    'lastMonth',
    'thisMonth'
];

export interface TransactionAmountsRequestParams extends PartialRecord<TransactionAmountsRequestType, StartEndTime> {
    readonly useTransactionTimezone: boolean;
    today?: StartEndTime;
    thisWeek?: StartEndTime;
    thisMonth?: StartEndTime;
    thisYear?: StartEndTime;
    lastMonth?: StartEndTime;
    monthBeforeLastMonth?: StartEndTime;
    monthBeforeLast2Months?: StartEndTime;
    monthBeforeLast3Months?: StartEndTime;
    monthBeforeLast4Months?: StartEndTime;
    monthBeforeLast5Months?: StartEndTime;
    monthBeforeLast6Months?: StartEndTime;
    monthBeforeLast7Months?: StartEndTime;
    monthBeforeLast8Months?: StartEndTime;
    monthBeforeLast9Months?: StartEndTime;
    monthBeforeLast10Months?: StartEndTime;
}

export class TransactionAmountsRequest {
    public readonly useTransactionTimezone: boolean;
    public readonly query: string;

    public constructor(useTransactionTimezone: boolean, query: string) {
        this.useTransactionTimezone = useTransactionTimezone;
        this.query = query;
    }

    public buildQuery(): string {
        return `use_transaction_timezone=${this.useTransactionTimezone}` + (this.query.length ? '&query=' + this.query : '');
    }

    public static of(params: TransactionAmountsRequestParams): TransactionAmountsRequest {
        const queryParams: string[] = [];

        ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE.forEach((type) => {
            if (params[type]) {
                queryParams.push(`${type}_${params[type].startTime}_${params[type].endTime}`);
            }
        });

        return new TransactionAmountsRequest(params.useTransactionTimezone, (queryParams.length ? queryParams.join('|') : ''));
    }
}

export interface TransactionInfoPageWrapperResponse {
    readonly items: TransactionInfoResponse[];
    readonly nextTimeSequenceId?: number;
    readonly totalCount?: number;
}

export interface TransactionInfoPageWrapperResponse2 {
    readonly items: TransactionInfoResponse[];
    readonly totalCount: number;
}

export interface TransactionPageWrapper {
    readonly items: Transaction[];
    readonly totalCount?: number;
}

export interface TransactionStatisticResponse {
    readonly startTime: number;
    readonly endTime: number;
    readonly items: TransactionStatisticResponseItem[];
}

export interface TransactionStatisticResponseItem {
    readonly categoryId: string;
    readonly accountId: string;
    readonly amount: number;
}

export interface TransactionStatisticTrendsResponseItem {
    readonly year: number;
    readonly month: number;
    readonly items: TransactionStatisticResponseItem[];
}

export interface YearMonthDataItem extends YearMonth, Record<string, unknown> {}

export interface YearMonthItems<T extends YearMonth> extends Record<string, unknown> {
    readonly items: T[];
}

export interface SortableTransactionStatisticDataItem {
    readonly name: string;
    readonly displayOrders: number[];
    readonly totalAmount: number;
}

export type TransactionStatisticDataItemType = 'category' | 'account' | 'total';

export interface TransactionStatisticDataItemBase extends SortableTransactionStatisticDataItem {
    readonly name: string;
    readonly type: TransactionStatisticDataItemType;
    readonly id: string;
    readonly icon: string;
    readonly color: string;
    readonly hidden: boolean;
    readonly displayOrders: number[];
    readonly totalAmount: number;
}

export interface TransactionCategoricalAnalysisData {
    readonly totalAmount: number;
    readonly items: TransactionCategoricalAnalysisDataItem[];
}

export interface TransactionCategoricalAnalysisDataItem extends TransactionStatisticDataItemBase {
    readonly percent: number;
}

export interface TransactionTrendsAnalysisData {
    readonly items: TransactionTrendsAnalysisDataItem[];
}

export interface TransactionTrendsAnalysisDataItem extends TransactionStatisticDataItemBase {
    readonly items: TransactionTrendsAnalysisDataAmount[];
}

export interface TransactionTrendsAnalysisDataAmount {
    readonly year: number;
    readonly month: number;
    readonly totalAmount: number;
}

export type TransactionAmountsResponse = PartialRecord<TransactionAmountsRequestType, TransactionAmountsResponseItem>;

export interface TransactionAmountsResponseItem {
    readonly startTime: number;
    readonly endTime: number;
    readonly amounts: TransactionAmountsResponseItemAmountInfo[];
}

export interface TransactionAmountsResponseItemAmountInfo {
    readonly currency: string;
    readonly incomeAmount: number;
    readonly expenseAmount: number;
}

export type TransactionOverviewResponse = PartialRecord<TransactionAmountsRequestType, TransactionOverviewResponseItem>;

export type TransactionOverviewDisplayTime = PartialRecord<TransactionAmountsRequestType, TransactionOverviewDisplayTimeItem>;

export interface TransactionOverviewDisplayTimeItem {
    readonly displayTime?: string;
    readonly startTime?: string;
    readonly endTime?: string;
}

export interface TransactionOverviewResponseItem {
    readonly valid: boolean;
    readonly incomeAmount: number;
    readonly expenseAmount: number;
    readonly incompleteIncomeAmount: boolean;
    readonly incompleteExpenseAmount: boolean;
    readonly amounts?: TransactionAmountsResponseItemAmountInfo[];
}

export interface TransactionMonthlyIncomeAndExpenseData {
    readonly monthStartTime: number;
    readonly incomeAmount: number;
    readonly expenseAmount: number;
    readonly incompleteIncomeAmount: boolean;
    readonly incompleteExpenseAmount: boolean;
}

export const EMPTY_TRANSACTION_RESULT: TransactionPageWrapper = {
    items: [],
    totalCount: 0
}
