import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useTransactionTagsStore } from './transactionTag.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { itemAndIndex, keys, values } from '@/core/base.ts';
import { AmountFilterType } from '@/core/numeral.ts';
import { DateRangeScene, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { AccountCategory } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import { ChartSortingType } from '@/core/statistics.ts';
import {
    TransactionExplorerChartTypeValue,
    TransactionExplorerChartType,
    TransactionExplorerDataDimensionType,
    TransactionExplorerDataDimension,
    TransactionExplorerValueMetricType,
    TransactionExplorerValueMetric,
    DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE
} from '@/core/explorer.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';

import { type Account } from '@/models/account.ts';
import { type TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTag } from '@/models/transaction_tag.ts';
import {
    type TransactionInfoResponse,
    type TransactionInsightDataItem
} from '@/models/transaction.ts';
import {
    TransactionExplorerQuery
} from '@/models/explorer.ts';

import {
    isDefined,
    isNumber,
    isInteger,
    isEquals,
} from '@/lib/common.ts';
import {
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithTimezoneOffset,
    getDateRangeByDateType,
    getFiscalYearFromUnixTime
} from '@/lib/datetime.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export enum TransactionExplorerDimensionType {
    DateTime = 'YYYY-MM-DD HH:mm:ss',
    YearMonthDay = 'YYYY-MM-DD',
    YearMonth = 'YYYY-MM',
    YearQuarter = 'YYYY-Q',
    Year = 'YYYY',
    TransactionType = 'transactionType',
    Category = 'category',
    Account = 'account',
    Amount = 'amount',
    Other = 'other'
}

export interface TransactionExplorerPartialFilter {
    dateRangeType?: number;
    startTime?: number;
    endTime?: number;
    queryId?: string;
    timezoneUsedForDateRange?: number;
    chartType?: TransactionExplorerChartTypeValue;
    categoryDimension?: TransactionExplorerDataDimensionType;
    seriesDimension?: TransactionExplorerDataDimensionType;
    valueMetric?: TransactionExplorerValueMetricType;
    chartSortingType?: number;
}

export interface TransactionExplorerFilter extends TransactionExplorerPartialFilter {
    dateRangeType: number;
    startTime: number;
    endTime: number;
    query: TransactionExplorerQuery[];
    timezoneUsedForDateRange: number;
    chartType: TransactionExplorerChartTypeValue;
    categoryDimension: TransactionExplorerDataDimensionType;
    seriesDimension: TransactionExplorerDataDimensionType;
    valueMetric: TransactionExplorerValueMetricType;
    chartSortingType: number;
}

export interface CategoriedInfo {
    categoryName: string;
    categoryNameNeedI18n?: boolean;
    categoryNameI18nParameters?: Record<string, string>;
    categoryId: string;
    categoryIdType: TransactionExplorerDimensionType;
    categoryDisplayOrders: number[];
}

export interface CategoriedTransactions extends CategoriedInfo {
    trasactions: Record<string, SeriesedTransactions>;
}

export interface CategoriedTransactionExplorerData extends CategoriedInfo {
    data: CategoriedTransactionExplorerDataItem[];
}

export interface SeriesedInfo {
    seriesName: string;
    seriesNameNeedI18n?: boolean;
    seriesNameI18nParameters?: Record<string, string>;
    seriesId: string;
    seriesIdType: TransactionExplorerDimensionType;
    seriesDisplayOrders: number[];
}

export interface SeriesedTransactions extends SeriesedInfo {
    trasactions: TransactionInsightDataItem[];
}

export interface CategoriedTransactionExplorerDataItem extends SeriesedInfo {
    value: number;
}

export const useExplorersStore = defineStore('explorers', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const currencyDisplayOrders: Record<string, number> = (() => {
        const result: Record<string, number> = {};
        let index: number = 0;

        for (const currency of keys(ALL_CURRENCIES)) {
            result[currency] = ++index;
        }

        return result;
    })();

    function getDataCategoryInfo(timezoneUsedForDateRange: number, dimension: TransactionExplorerDataDimension, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): CategoriedInfo {
        let transactionTimeUtfOffset: number | undefined = undefined;

        if (timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
            transactionTimeUtfOffset = transaction.utcOffset;
        }

        if (dimension === TransactionExplorerDataDimension.None) {
            const valueMetric = TransactionExplorerValueMetric.valueOf(transactionExplorerFilter.value.valueMetric);
            return {
                categoryName: valueMetric?.name ?? 'Unknown',
                categoryNameNeedI18n: true,
                categoryId: 'none',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [1]
            };
        } else if (dimension === TransactionExplorerDataDimension.Query) {
            if (queryName) {
                return {
                    categoryName: queryName,
                    categoryId: (queryIndex + 1).toString(10),
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [queryIndex + 1]
                };
            } else {
                return {
                    categoryName: `format.misc.queryIndex`,
                    categoryNameNeedI18n: true,
                    categoryNameI18nParameters: {
                        index: (queryIndex + 1).toString(10)
                    },
                    categoryId: (queryIndex + 1).toString(10),
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [queryIndex + 1]
                };
            }
        } else if (dimension === TransactionExplorerDataDimension.DateTime) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
            const textualDateTime = `${dateTime.getGregorianCalendarYearDashMonthDashDay()} ${dateTime.getHour().toString(10).padStart(2, '0')}:${dateTime.getMinute().toString(10).padStart(2, '0')}:${dateTime.getSecond().toString(10).padStart(2, '0')}`;

            return {
                categoryName: textualDateTime,
                categoryId: textualDateTime,
                categoryIdType: TransactionExplorerDimensionType.DateTime,
                categoryDisplayOrders: [dateTime.getUnixTime()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYearMonthDay) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
            const yearMonthDay = dateTime.getGregorianCalendarYearDashMonthDashDay();

            return {
                categoryName: yearMonthDay,
                categoryId: yearMonthDay,
                categoryIdType: TransactionExplorerDimensionType.YearMonthDay,
                categoryDisplayOrders: [dateTime.getUnixTime()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYearMonth) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
            const yearMonth = dateTime.getGregorianCalendarYearDashMonth();

            return {
                categoryName: yearMonth,
                categoryId: yearMonth,
                categoryIdType: TransactionExplorerDimensionType.YearMonth,
                categoryDisplayOrders: [dateTime.getUnixTime()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYearQuarter) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
            const yearQuarter = `${dateTime.getGregorianCalendarYear().toString(10)}-${dateTime.getGregorianCalendarQuarter().toString(10)}`;

            return {
                categoryName: yearQuarter,
                categoryId: yearQuarter,
                categoryIdType: TransactionExplorerDimensionType.YearQuarter,
                categoryDisplayOrders: [dateTime.getUnixTime()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByYear) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarYear().toString(10),
                categoryId: dateTime.getGregorianCalendarYear().toString(10),
                categoryIdType: TransactionExplorerDimensionType.Year,
                categoryDisplayOrders: [dateTime.getUnixTime()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByFiscalYear) {
            const fiscalYear = getFiscalYearFromUnixTime(transaction.time, userStore.currentUserFiscalYearStart, transactionTimeUtfOffset).toString(10);

            return {
                categoryName: fiscalYear,
                categoryId: fiscalYear,
                categoryIdType: TransactionExplorerDimensionType.Year,
                categoryDisplayOrders: [parseInt(fiscalYear)]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByDayOfWeek) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getWeekDay().name,
                categoryId: dateTime.getWeekDay().type.toString(10),
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [dateTime.getWeekDay().getDisplayOrder(userStore.currentUserFirstDayOfWeek)]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByDayOfMonth) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarDay().toString(10),
                categoryId: dateTime.getGregorianCalendarDay().toString(10),
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [dateTime.getGregorianCalendarDay()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByMonthOfYear) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarMonth().toString(10),
                categoryId: dateTime.getGregorianCalendarMonth().toString(10),
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [dateTime.getGregorianCalendarMonth()]
            };
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByQuarterOfYear) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarQuarter().toString(10),
                categoryId: dateTime.getGregorianCalendarQuarter().toString(10),
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [dateTime.getGregorianCalendarQuarter()]
            };
        } else if (dimension === TransactionExplorerDataDimension.TransactionType) {
            let transactionTypeName = 'Unknown';

            if (transaction.type === TransactionType.ModifyBalance) {
                transactionTypeName = 'Modify Balance';
            } else if (transaction.type === TransactionType.Income) {
                transactionTypeName = 'Income';
            } else if (transaction.type === TransactionType.Expense) {
                transactionTypeName = 'Expense';
            } else if (transaction.type === TransactionType.Transfer) {
                transactionTypeName = 'Transfer';
            }

            return {
                categoryName: transactionTypeName,
                categoryNameNeedI18n: true,
                categoryId: transaction.type.toString(10),
                categoryIdType: TransactionExplorerDimensionType.TransactionType,
                categoryDisplayOrders: [transaction.type]
            };
        } else if (dimension === TransactionExplorerDataDimension.SourceAccount) {
            const primaryAccount = accountsStore.allAccountsMap[transaction.sourceAccount.parentId] ?? transaction.sourceAccount;
            const primaryAccountCategoryDisplayOrder: number = settingsStore.accountCategoryDisplayOrders[primaryAccount.category] || Number.MAX_SAFE_INTEGER;

            return {
                categoryName: transaction.sourceAccountName || 'Unknown',
                categoryNameNeedI18n: !transaction.sourceAccountName,
                categoryId: transaction.sourceAccountId || 'unknown',
                categoryIdType: TransactionExplorerDimensionType.Account,
                categoryDisplayOrders: [primaryAccountCategoryDisplayOrder, primaryAccount.displayOrder, transaction.sourceAccount.displayOrder]
            };
        } else if (dimension === TransactionExplorerDataDimension.SourceAccountCategory) {
            const accountCategory = AccountCategory.valueOf(transaction.sourceAccount.category);
            const accountCategoryDisplayOrder: number = settingsStore.accountCategoryDisplayOrders[accountCategory?.type ?? 0] || Number.MAX_SAFE_INTEGER;

            return {
                categoryName: accountCategory?.name || 'Unknown',
                categoryNameNeedI18n: true,
                categoryId: accountCategory?.type.toString(10) || 'unknown',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [accountCategoryDisplayOrder]
            };
        } else if (dimension === TransactionExplorerDataDimension.SourceAccountCurrency) {
            return {
                categoryName: transaction.sourceAccount.currency || 'Unknown',
                categoryNameNeedI18n: !transaction.sourceAccount.currency,
                categoryId: transaction.sourceAccount.currency || 'unknown',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [currencyDisplayOrders[transaction.sourceAccount.currency] || 0]
            };
        }  else if (dimension === TransactionExplorerDataDimension.DestinationAccount) {
            const primaryAccount = accountsStore.allAccountsMap[transaction.destinationAccount?.parentId ?? ''] ?? transaction.destinationAccount;
            const primaryAccountCategoryDisplayOrder: number = settingsStore.accountCategoryDisplayOrders[primaryAccount?.category || 0] || Number.MAX_SAFE_INTEGER;

            return {
                categoryName: transaction.type === TransactionType.Transfer ? (transaction.destinationAccountName || 'Unknown') : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer || !transaction.destinationAccountName,
                categoryId: transaction.type === TransactionType.Transfer ? (transaction.destinationAccountId || 'unknown') : 'none',
                categoryIdType: TransactionExplorerDimensionType.Account,
                categoryDisplayOrders: transaction.type === TransactionType.Transfer && primaryAccount && transaction.destinationAccount ? [primaryAccountCategoryDisplayOrder, primaryAccount.displayOrder, transaction.destinationAccount.displayOrder] : [0]
            };
        } else if (dimension === TransactionExplorerDataDimension.DestinationAccountCategory) {
            const accountCategory = transaction.type === TransactionType.Transfer && transaction.destinationAccount ? AccountCategory.valueOf(transaction.destinationAccount.category) : undefined;
            const accountCategoryDisplayOrder: number = settingsStore.accountCategoryDisplayOrders[accountCategory?.type ?? 0] || Number.MAX_SAFE_INTEGER;

            return {
                categoryName: transaction.type === TransactionType.Transfer ? (accountCategory?.name || 'Unknown') : 'None',
                categoryNameNeedI18n: true,
                categoryId: transaction.type === TransactionType.Transfer ? (accountCategory?.name || 'unknown') : 'none',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: transaction.type === TransactionType.Transfer ? [accountCategoryDisplayOrder] : [0]
            };
        } else if (dimension === TransactionExplorerDataDimension.DestinationAccountCurrency) {
            return {
                categoryName: transaction.type === TransactionType.Transfer ? (transaction.destinationAccount?.currency || 'Unknown') : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer || !transaction.destinationAccount?.currency,
                categoryId: transaction.type === TransactionType.Transfer ? (transaction.destinationAccount?.currency || 'unknown') : 'none',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: transaction.type === TransactionType.Transfer ? [currencyDisplayOrders[transaction.destinationAccount?.currency ?? ''] || 0] : [0]
            };
        } else if (dimension === TransactionExplorerDataDimension.PrimaryCategory) {
            return {
                categoryName: transaction.primaryCategory.name,
                categoryId: transaction.primaryCategory.id,
                categoryIdType: TransactionExplorerDimensionType.Category,
                categoryDisplayOrders: [transaction.primaryCategory.displayOrder]
            };
        } else if (dimension === TransactionExplorerDataDimension.SecondaryCategory) {
            return {
                categoryName: transaction.secondaryCategory.name,
                categoryId: transaction.categoryId,
                categoryIdType: TransactionExplorerDimensionType.Category,
                categoryDisplayOrders: [transaction.primaryCategory.displayOrder, transaction.secondaryCategory.displayOrder]
            };
        } else if (dimension === TransactionExplorerDataDimension.SourceAmount) {
            return {
                categoryName: transaction.sourceAmount.toString(10),
                categoryId: transaction.sourceAmount.toString(10),
                categoryIdType: TransactionExplorerDimensionType.Amount,
                categoryDisplayOrders: [transaction.sourceAmount]
            };
        } else if (dimension === TransactionExplorerDataDimension.DestinationAmount) {
            return {
                categoryName: transaction.type === TransactionType.Transfer ? transaction.destinationAmount.toString(10) : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer,
                categoryId: transaction.type === TransactionType.Transfer ? transaction.destinationAmount.toString(10) : 'none',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [transaction.destinationAmount]
            };
        } else {
            return {
                categoryName: '',
                categoryId: '',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [0]
            };
        }
    }

    function addTransactionToCategoriedDataMap(timezoneUsedForDateRange: number, categoriedDataMap: Record<string, CategoriedTransactions>, categoryDimension: TransactionExplorerDataDimension, seriesDemension: TransactionExplorerDataDimension, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): void {
        const categoriedInfo = getDataCategoryInfo(timezoneUsedForDateRange, categoryDimension, queryName, queryIndex, transaction);
        let categoriedData = categoriedDataMap[categoriedInfo.categoryId];

        if (!categoriedData) {
            categoriedData = {
                categoryName: categoriedInfo.categoryName,
                categoryNameNeedI18n: categoriedInfo.categoryNameNeedI18n,
                categoryNameI18nParameters: categoriedInfo.categoryNameI18nParameters,
                categoryId: categoriedInfo.categoryId,
                categoryIdType: categoriedInfo.categoryIdType,
                categoryDisplayOrders: categoriedInfo.categoryDisplayOrders,
                trasactions: {}
            };
            categoriedDataMap[categoriedInfo.categoryId] = categoriedData;
        }

        const seriesedInfo = getDataCategoryInfo(timezoneUsedForDateRange, seriesDemension, queryName, queryIndex, transaction);
        let seriesedData = categoriedData.trasactions[seriesedInfo.categoryId];

        if (!seriesedData) {
            seriesedData = {
                seriesName: seriesedInfo.categoryName,
                seriesNameNeedI18n: seriesedInfo.categoryNameNeedI18n,
                seriesNameI18nParameters: seriesedInfo.categoryNameI18nParameters,
                seriesId: seriesedInfo.categoryId,
                seriesIdType: seriesedInfo.categoryIdType,
                seriesDisplayOrders: seriesedInfo.categoryDisplayOrders,
                trasactions: []
            };
            categoriedData.trasactions[seriesedInfo.categoryId] = seriesedData;
        }

        seriesedData.trasactions.push(transaction);
    }

    const transactionExplorerFilter = ref<TransactionExplorerFilter>({
        dateRangeType: DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type,
        startTime: 0,
        endTime: 0,
        query: [],
        timezoneUsedForDateRange: TimezoneTypeForStatistics.Default.type,
        chartType: TransactionExplorerChartType.Default.value,
        categoryDimension: TransactionExplorerDataDimension.CategoryDimensionDefault.value,
        seriesDimension: TransactionExplorerDataDimension.SeriesDimensionDefault.value,
        valueMetric: TransactionExplorerValueMetric.Default.value,
        chartSortingType: ChartSortingType.Default.type
    });

    const transactionExplorerAllData = ref<TransactionInfoResponse[]>([]);
    const transactionExplorerStateInvalid = ref<boolean>(true);

    const allTransactions = computed<TransactionInsightDataItem[]>(() => {
        if (!transactionExplorerAllData.value || transactionExplorerAllData.value.length < 1) {
            return [];
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of transactionExplorerAllData.value) {
            const sourceAccount: Account | undefined = accountsStore.allAccountsMap[transaction.sourceAccountId];

            if (!sourceAccount) {
                continue;
            }

            let destinationAccount: Account | undefined = undefined

            if (transaction.destinationAccountId && transaction.destinationAccountId !== '0') {
                destinationAccount = accountsStore.allAccountsMap[transaction.destinationAccountId];

                if (!destinationAccount) {
                    continue;
                }
            }

            const secondaryCategory: TransactionCategory | undefined = transactionCategoriesStore.allTransactionCategoriesMap[transaction.categoryId];

            if (!secondaryCategory) {
                continue;
            }

            const primaryCategory: TransactionCategory | undefined = transactionCategoriesStore.allTransactionCategoriesMap[secondaryCategory.parentId];

            if (!primaryCategory) {
                continue;
            }

            const tags: TransactionTag[] = [];

            for (const tagId of transaction.tagIds) {
                const tag: TransactionTag | undefined = transactionTagsStore.allTransactionTagsMap[tagId];

                if (tag) {
                    tags.push(tag);
                }
            }

            const item: TransactionInsightDataItem = {
                ...transaction,
                id: transaction.id,
                time: transaction.time,
                utcOffset: transaction.utcOffset,
                type: transaction.type,
                primaryCategory: primaryCategory,
                primaryCategoryName: primaryCategory.name,
                secondaryCategory: secondaryCategory,
                secondaryCategoryName: secondaryCategory.name,
                sourceAccount: sourceAccount,
                sourceAccountName: sourceAccount.name,
                destinationAccount: destinationAccount,
                destinationAccountName: destinationAccount?.name,
                sourceAmount: transaction.sourceAmount,
                destinationAmount: transaction.destinationAmount,
                hideAmount: transaction.hideAmount,
                tags: tags,
                comment: transaction.comment,
                geoLocation: transaction.geoLocation
            };

            result.push(item);
        }

        return result;
    });

    const filteredTransactions = computed<TransactionInsightDataItem[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        if (!transactionExplorerFilter.value.query || transactionExplorerFilter.value.query.length < 1) {
            return allTransactions.value;
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of allTransactions.value) {
            for (const query of transactionExplorerFilter.value.query) {
                if (query.match(transaction)) {
                    result.push(transaction);
                    break;
                }
            }
        }

        return result;
    });

    const categoriedTransactions = computed<Record<string, CategoriedTransactions>>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return {};
        }

        const chartType = TransactionExplorerChartType.valueOf(transactionExplorerFilter.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(transactionExplorerFilter.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(transactionExplorerFilter.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;

        if (!chartType || !categoryDimension || !seriesDimension) {
            return {};
        }

        const categoriedDataMap: Record<string, CategoriedTransactions> = {};

        for (const transaction of allTransactions.value) {
            if (!transactionExplorerFilter.value.query || transactionExplorerFilter.value.query.length < 1) {
                addTransactionToCategoriedDataMap(transactionExplorerFilter.value.timezoneUsedForDateRange, categoriedDataMap, categoryDimension, seriesDimension, '', 0, transaction);
                continue;
            }

            for (const [query, index] of itemAndIndex(transactionExplorerFilter.value.query)) {
                if (query.match(transaction)) {
                    addTransactionToCategoriedDataMap(transactionExplorerFilter.value.timezoneUsedForDateRange, categoriedDataMap, categoryDimension, seriesDimension, query.name, index, transaction);

                    if (categoryDimension !== TransactionExplorerDataDimension.Query) {
                        break;
                    }
                }
            }
        }

        return categoriedDataMap;
    });

    const categoriedTransactionExplorerData = computed<CategoriedTransactionExplorerData[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        const chartType = TransactionExplorerChartType.valueOf(transactionExplorerFilter.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(transactionExplorerFilter.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(transactionExplorerFilter.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;
        const valueMetric = TransactionExplorerValueMetric.valueOf(transactionExplorerFilter.value.valueMetric);

        if (!chartType || !categoryDimension || !seriesDimension || !valueMetric) {
            return [];
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const result: CategoriedTransactionExplorerData[] = [];
        const categoriedDataMap = categoriedTransactions.value;

        for (const categoriedTransactions of values(categoriedDataMap)) {
            const dataItems: CategoriedTransactionExplorerDataItem[] = [];
            let allSeriesedTransactions: Record<string, SeriesedTransactions> = categoriedTransactions.trasactions;

            // merge all series into one for pie/radar chart
            if (chartType === TransactionExplorerChartType.Pie || chartType === TransactionExplorerChartType.Radar) {
                const transactions: TransactionInsightDataItem[] = [];

                for (const seriesedTransactions of values(categoriedTransactions.trasactions)) {
                    transactions.push(...seriesedTransactions.trasactions);
                }

                allSeriesedTransactions = {};
                allSeriesedTransactions['none'] = {
                    seriesName: valueMetric?.name ?? 'Unknown',
                    seriesNameNeedI18n: true,
                    seriesId: 'none',
                    seriesIdType: TransactionExplorerDimensionType.Other,
                    seriesDisplayOrders: [0],
                    trasactions: transactions
                };
            }

            for (const seriesedTransactions of values(allSeriesedTransactions)) {
                const allSourceAmountsInDefaultCurrency: number[] = [];
                let totalSourceAmountSumInDefaultCurrency: number = 0;
                let minimumSourceAmountInDefaultCurrency: number = Number.MAX_SAFE_INTEGER;
                let maximumSourceAmountInDefaultCurrency: number = Number.MIN_SAFE_INTEGER;

                for (const transaction of seriesedTransactions.trasactions) {
                    let amountInDefaultCurrency: number = transaction.sourceAmount;

                    if (transaction.sourceAccount.currency !== defaultCurrency) {
                        const amount = exchangeRatesStore.getExchangedAmount(transaction.sourceAmount, transaction.sourceAccount.currency, defaultCurrency);

                        if (isNumber(amount)) {
                            amountInDefaultCurrency = Math.trunc(amount);
                        } else {
                            continue;
                        }
                    }

                    allSourceAmountsInDefaultCurrency.push(amountInDefaultCurrency);
                    totalSourceAmountSumInDefaultCurrency += amountInDefaultCurrency;

                    if (amountInDefaultCurrency >= 0 && amountInDefaultCurrency < minimumSourceAmountInDefaultCurrency) {
                        minimumSourceAmountInDefaultCurrency = amountInDefaultCurrency;
                    }

                    if (amountInDefaultCurrency > maximumSourceAmountInDefaultCurrency) {
                        maximumSourceAmountInDefaultCurrency = amountInDefaultCurrency;
                    }
                }

                let value: number = 0;

                if (valueMetric === TransactionExplorerValueMetric.TransactionCount) {
                    value = allSourceAmountsInDefaultCurrency.length;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSum) {
                    value = totalSourceAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountAverage) {
                    value = allSourceAmountsInDefaultCurrency.length > 0 ? Math.trunc(totalSourceAmountSumInDefaultCurrency / allSourceAmountsInDefaultCurrency.length) : 0;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMedian) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        value = allSourceAmountsInDefaultCurrency[Math.floor(allSourceAmountsInDefaultCurrency.length / 2)] as number;
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMinimum) {
                    value = minimumSourceAmountInDefaultCurrency === Number.MAX_SAFE_INTEGER ? 0 : minimumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMaximum) {
                    value = maximumSourceAmountInDefaultCurrency === Number.MIN_SAFE_INTEGER ? 0 : maximumSourceAmountInDefaultCurrency;
                }

                dataItems.push({
                    seriesName: seriesedTransactions.seriesName,
                    seriesNameNeedI18n: seriesedTransactions.seriesNameNeedI18n,
                    seriesNameI18nParameters: seriesedTransactions.seriesNameI18nParameters,
                    seriesId: seriesedTransactions.seriesId,
                    seriesIdType: seriesedTransactions.seriesIdType,
                    seriesDisplayOrders: seriesedTransactions.seriesDisplayOrders,
                    value: value
                });
            }

            result.push({
                categoryName: categoriedTransactions.categoryName,
                categoryNameNeedI18n: categoriedTransactions.categoryNameNeedI18n,
                categoryNameI18nParameters: categoriedTransactions.categoryNameI18nParameters,
                categoryId: categoriedTransactions.categoryId,
                categoryIdType: categoriedTransactions.categoryIdType,
                categoryDisplayOrders: categoriedTransactions.categoryDisplayOrders,
                data: dataItems
            });
        }

        return result;
    });

    function updateTransactionExplorerInvalidState(invalidState: boolean): void {
        transactionExplorerStateInvalid.value = invalidState;
    }

    function resetTransactionExplorers(): void {
        transactionExplorerFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type;
        transactionExplorerFilter.value.startTime = 0;
        transactionExplorerFilter.value.endTime = 0;
        transactionExplorerFilter.value.query = [];
        transactionExplorerFilter.value.timezoneUsedForDateRange = TimezoneTypeForStatistics.Default.type;
        transactionExplorerFilter.value.chartType = TransactionExplorerChartType.Default.value;
        transactionExplorerFilter.value.categoryDimension = TransactionExplorerDataDimension.CategoryDimensionDefault.value;
        transactionExplorerFilter.value.seriesDimension = TransactionExplorerDataDimension.SeriesDimensionDefault.value;
        transactionExplorerFilter.value.valueMetric = TransactionExplorerValueMetric.Default.value;
        transactionExplorerFilter.value.chartSortingType = ChartSortingType.Default.type;
        transactionExplorerAllData.value = [];
        transactionExplorerStateInvalid.value = true;
    }

    function initTransactionExplorerFilter(filter?: TransactionExplorerPartialFilter, resetQuery?: boolean): void {
        if (filter && isInteger(filter.dateRangeType)) {
            transactionExplorerFilter.value.dateRangeType = filter.dateRangeType;
        } else {
            transactionExplorerFilter.value.dateRangeType = settingsStore.appSettings.insightsExplorerDefaultDateRangeType;
        }

        let dateRangeTypeValid = true;

        if (!DateRange.isAvailableForScene(transactionExplorerFilter.value.dateRangeType, DateRangeScene.InsightsExplorer)) {
            transactionExplorerFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type;
            dateRangeTypeValid = false;
        }

        if (dateRangeTypeValid && transactionExplorerFilter.value.dateRangeType === DateRange.Custom.type) {
            if (filter && isInteger(filter.startTime)) {
                transactionExplorerFilter.value.startTime = filter.startTime;
            } else {
                transactionExplorerFilter.value.startTime = 0;
            }

            if (filter && isInteger(filter.endTime)) {
                transactionExplorerFilter.value.endTime = filter.endTime;
            } else {
                transactionExplorerFilter.value.endTime = 0;
            }
        } else {
            const dateRange = getDateRangeByDateType(transactionExplorerFilter.value.dateRangeType, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart);

            if (dateRange) {
                transactionExplorerFilter.value.dateRangeType = dateRange.dateType;
                transactionExplorerFilter.value.startTime = dateRange.minTime;
                transactionExplorerFilter.value.endTime = dateRange.maxTime;
            }
        }

        if (resetQuery) {
            transactionExplorerFilter.value.query = [];
            transactionExplorerFilter.value.timezoneUsedForDateRange = TimezoneTypeForStatistics.Default.type;
            transactionExplorerFilter.value.chartType = TransactionExplorerChartType.Default.value;
            transactionExplorerFilter.value.categoryDimension = TransactionExplorerDataDimension.CategoryDimensionDefault.value;
            transactionExplorerFilter.value.seriesDimension = TransactionExplorerDataDimension.SeriesDimensionDefault.value;
            transactionExplorerFilter.value.valueMetric = TransactionExplorerValueMetric.Default.value;
            transactionExplorerFilter.value.chartSortingType = ChartSortingType.Default.type;
        }
    }

    function updateTransactionExplorerFilter(filter: TransactionExplorerPartialFilter): boolean {
        let changed = false;

        if (filter && isInteger(filter.dateRangeType) && transactionExplorerFilter.value.dateRangeType !== filter.dateRangeType) {
            transactionExplorerFilter.value.dateRangeType = filter.dateRangeType;
            changed = true;
        }

        if (filter && isInteger(filter.startTime) && transactionExplorerFilter.value.startTime !== filter.startTime) {
            transactionExplorerFilter.value.startTime = filter.startTime;
            changed = true;
        }

        if (filter && isInteger(filter.endTime) && transactionExplorerFilter.value.endTime !== filter.endTime) {
            transactionExplorerFilter.value.endTime = filter.endTime;
            changed = true;
        }

        if (filter && isDefined(filter.timezoneUsedForDateRange) && transactionExplorerFilter.value.timezoneUsedForDateRange !== filter.timezoneUsedForDateRange) {
            transactionExplorerFilter.value.timezoneUsedForDateRange = filter.timezoneUsedForDateRange;
            changed = true;
        }

        if (filter && isDefined(filter.chartType) && transactionExplorerFilter.value.chartType !== filter.chartType) {
            transactionExplorerFilter.value.chartType = filter.chartType;
            changed = true;
        }

        if (filter && isDefined(filter.categoryDimension) && transactionExplorerFilter.value.categoryDimension !== filter.categoryDimension) {
            transactionExplorerFilter.value.categoryDimension = filter.categoryDimension;
            changed = true;
        }

        if (filter && isDefined(filter.seriesDimension) && transactionExplorerFilter.value.seriesDimension !== filter.seriesDimension) {
            transactionExplorerFilter.value.seriesDimension = filter.seriesDimension;
            changed = true;
        }

        if (filter && isDefined(filter.valueMetric) && transactionExplorerFilter.value.valueMetric !== filter.valueMetric) {
            transactionExplorerFilter.value.valueMetric = filter.valueMetric;
            changed = true;
        }

        if (filter && isDefined(filter.chartSortingType) && transactionExplorerFilter.value.chartSortingType !== filter.chartSortingType) {
            transactionExplorerFilter.value.chartSortingType = filter.chartSortingType;
            changed = true;
        }

        if (transactionExplorerFilter.value.seriesDimension === transactionExplorerFilter.value.categoryDimension && transactionExplorerFilter.value.seriesDimension !== TransactionExplorerDataDimension.SeriesDimensionDefault.value) {
            transactionExplorerFilter.value.seriesDimension = TransactionExplorerDataDimension.SeriesDimensionDefault.value;
            changed = true;
        }

        return changed;
    }

    function getTransactionExplorerPageParams(currentExplorationId: string, activeTab: string): string {
        const querys: string[] = [];

        if (currentExplorationId) {
            querys.push('id=' + currentExplorationId);
        }

        if (activeTab) {
            querys.push('activeTab=' + activeTab);
        }

        querys.push('dateRangeType=' + transactionExplorerFilter.value.dateRangeType);
        querys.push('startTime=' + transactionExplorerFilter.value.startTime);
        querys.push('endTime=' + transactionExplorerFilter.value.endTime);

        return querys.join('&');
    }

    function getTransactionListPageParams(dimensionType: TransactionExplorerDimensionType, itemId: string): string {
        const querys: string[] = [];

        if (dimensionType === TransactionExplorerDimensionType.TransactionType) {
            querys.push(`type=${itemId}`);
        } else if (dimensionType === TransactionExplorerDimensionType.Account) {
            querys.push(`accountIds=${itemId}`);
        } else if (dimensionType === TransactionExplorerDimensionType.Category) {
            querys.push(`categoryIds=${itemId}`);
        } else if (dimensionType === TransactionExplorerDimensionType.Amount) {
            querys.push(`amountFilter=${encodeURIComponent(AmountFilterType.EqualTo.toTextualFilter(parseInt(itemId)))}`);
        } else {
            return '';
        }

        querys.push('dateType=' + transactionExplorerFilter.value.dateRangeType);
        querys.push('minTime=' + transactionExplorerFilter.value.startTime);
        querys.push('maxTime=' + transactionExplorerFilter.value.endTime);

        return querys.join('&');
    }

    function loadAllTransactions({ force }: { force: boolean }): Promise<TransactionInfoResponse[]> {
        return new Promise((resolve, reject) => {
            services.getAllTransactions({
                startTime: transactionExplorerFilter.value.startTime,
                endTime: transactionExplorerFilter.value.endTime,
                withPictures: true
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve all transactions' });
                    return;
                }

                if (transactionExplorerStateInvalid.value) {
                    updateTransactionExplorerInvalidState(false);
                }

                if (force && data.result && isEquals(transactionExplorerAllData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionExplorerAllData.value = data.result;

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to load all transactions', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve all transactions' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        transactionExplorerFilter: transactionExplorerFilter,
        transactionExplorerStateInvalid,
        // computed
        filteredTransactions,
        categoriedTransactionExplorerData,
        // functions
        updateTransactionExplorerInvalidState,
        resetTransactionExplorers,
        initTransactionExplorerFilter,
        updateTransactionExplorerFilter,
        getTransactionExplorerPageParams,
        getTransactionListPageParams,
        loadAllTransactions
    };
});
