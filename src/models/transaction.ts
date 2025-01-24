import type { PartialRecord } from '@/core/base.ts';
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
    readonly amount: number;
}

export interface TransactionStatisticResponseWithInfo {
    readonly startTime: number;
    readonly endTime: number;
    readonly items: TransactionStatisticResponseItemWithInfo[];
}

export interface TransactionStatisticResponseItemWithInfo extends TransactionStatisticResponseItem {
    readonly account?: AccountInfoResponse;
    readonly primaryAccount?: AccountInfoResponse;
    readonly category?: TransactionCategoryInfoResponse;
    readonly primaryCategory?: TransactionCategoryInfoResponse;
    readonly amountInDefaultCurrency: number | null;
}

export interface TransactionStatisticTrendsResponseItem {
    readonly year: number;
    readonly month: number;
    readonly items: TransactionStatisticResponseItem[];
}

export interface TransactionStatisticTrendsResponseItemWithInfo {
    readonly year: number;
    readonly month: number;
    readonly items: TransactionStatisticResponseItemWithInfo[];
}

export type TransactionStatisticDataItemType = 'category' | 'account' | 'total';

export interface TransactionStatisticDataItemBase {
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
