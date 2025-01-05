import type { StartEndTime } from '@/core/datetime.ts';

import type { AccountInfoResponse } from './account.ts';
import type { TransactionCategoryInfoResponse } from './transaction_category.ts';
import type { TransactionPictureInfoBasicResponse } from './transaction_picture_info.ts';
import type { TransactionTagInfoResponse } from './transaction_tag.ts';

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

export interface TransactionGeoLocationResponse {
    readonly latitude: number;
    readonly longitude: number;
}

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

export interface TransactionAmountsRequestParams {
    readonly useTransactionTimezone: boolean;
    readonly today: StartEndTime;
    readonly thisWeek: StartEndTime;
    readonly thisMonth: StartEndTime;
    readonly thisYear: StartEndTime;
    readonly lastMonth: StartEndTime;
    readonly monthBeforeLastMonth: StartEndTime;
    readonly monthBeforeLast2Months: StartEndTime;
    readonly monthBeforeLast3Months: StartEndTime;
    readonly monthBeforeLast4Months: StartEndTime;
    readonly monthBeforeLast5Months: StartEndTime;
    readonly monthBeforeLast6Months: StartEndTime;
    readonly monthBeforeLast7Months: StartEndTime;
    readonly monthBeforeLast8Months: StartEndTime;
    readonly monthBeforeLast9Months: StartEndTime;
    readonly monthBeforeLast10Months: StartEndTime;
}

export class TransactionAmountsRequest {
    public readonly useTransactionTimezone: boolean;
    public readonly query: string;

    constructor(useTransactionTimezone: boolean, query: string) {
        this.useTransactionTimezone = useTransactionTimezone;
        this.query = query;
    }

    public buildQuery(): string {
        return `use_transaction_timezone=${this.useTransactionTimezone}` + (this.query.length ? '&query=' + this.query : '');
    }

    public static of(params: TransactionAmountsRequestParams): TransactionAmountsRequest {
        const queryParams = [];

        if (params.today) {
            queryParams.push(`today_${params.today.startTime}_${params.today.endTime}`);
        }

        if (params.thisWeek) {
            queryParams.push(`thisWeek_${params.thisWeek.startTime}_${params.thisWeek.endTime}`);
        }

        if (params.thisMonth) {
            queryParams.push(`thisMonth_${params.thisMonth.startTime}_${params.thisMonth.endTime}`);
        }

        if (params.thisYear) {
            queryParams.push(`thisYear_${params.thisYear.startTime}_${params.thisYear.endTime}`);
        }

        if (params.lastMonth) {
            queryParams.push(`lastMonth_${params.lastMonth.startTime}_${params.lastMonth.endTime}`);
        }

        if (params.monthBeforeLastMonth) {
            queryParams.push(`monthBeforeLastMonth_${params.monthBeforeLastMonth.startTime}_${params.monthBeforeLastMonth.endTime}`);
        }

        if (params.monthBeforeLast2Months) {
            queryParams.push(`monthBeforeLast2Months_${params.monthBeforeLast2Months.startTime}_${params.monthBeforeLast2Months.endTime}`);
        }

        if (params.monthBeforeLast3Months) {
            queryParams.push(`monthBeforeLast3Months_${params.monthBeforeLast3Months.startTime}_${params.monthBeforeLast3Months.endTime}`);
        }

        if (params.monthBeforeLast4Months) {
            queryParams.push(`monthBeforeLast4Months_${params.monthBeforeLast4Months.startTime}_${params.monthBeforeLast4Months.endTime}`);
        }

        if (params.monthBeforeLast5Months) {
            queryParams.push(`monthBeforeLast5Months_${params.monthBeforeLast5Months.startTime}_${params.monthBeforeLast5Months.endTime}`);
        }

        if (params.monthBeforeLast6Months) {
            queryParams.push(`monthBeforeLast6Months_${params.monthBeforeLast6Months.startTime}_${params.monthBeforeLast6Months.endTime}`);
        }

        if (params.monthBeforeLast7Months) {
            queryParams.push(`monthBeforeLast7Months_${params.monthBeforeLast7Months.startTime}_${params.monthBeforeLast7Months.endTime}`);
        }

        if (params.monthBeforeLast8Months) {
            queryParams.push(`monthBeforeLast8Months_${params.monthBeforeLast8Months.startTime}_${params.monthBeforeLast8Months.endTime}`);
        }

        if (params.monthBeforeLast9Months) {
            queryParams.push(`monthBeforeLast9Months_${params.monthBeforeLast9Months.startTime}_${params.monthBeforeLast9Months.endTime}`);
        }

        if (params.monthBeforeLast10Months) {
            queryParams.push(`monthBeforeLast10Months_${params.monthBeforeLast10Months.startTime}_${params.monthBeforeLast10Months.endTime}`);
        }

        return new TransactionAmountsRequest(params.useTransactionTimezone, (queryParams.length ? queryParams.join('|') : ''));
    }
}

export interface TransactionInfoPageWrapperResponse {
    readonly items: TransactionInfoResponse[];
    readonly nextTimeSequenceId?: string;
    readonly totalCount?: number;
}

export interface TransactionInfoPageWrapperResponse2 {
    readonly items: TransactionInfoResponse[];
    readonly totalCount: number;
}

export interface TransactionStatisticResponse {
    readonly startTime: number;
    readonly endTime: number;
    readonly items: TransactionStatisticResponseItem[];
}

export interface TransactionStatisticResponseItem {
    readonly categoryId: string;
    readonly accountId: string;
    readonly totalAmount: number;
}

export interface TransactionStatisticTrendsItem {
    readonly year: number;
    readonly month: number;
    readonly items: TransactionStatisticResponseItem[];
}

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
