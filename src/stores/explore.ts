import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useTransactionTagsStore } from './transactionTag.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { itemAndIndex, values } from '@/core/base.ts';
import { AmountFilterType } from '@/core/numeral.ts';
import { DateRangeScene, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { AccountCategory } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    TransactionExploreChartTypeValue,
    TransactionExploreChartType,
    TransactionExploreDataDimensionType,
    TransactionExploreDataDimension,
    TransactionExploreValueMetricType,
    TransactionExploreValueMetric,
    DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE
} from '@/core/explore.ts';

import { type Account } from '@/models/account.ts';
import { type TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTag } from '@/models/transaction_tag.ts';
import {
    type TransactionInfoResponse,
    type TransactionInsightDataItem
} from '@/models/transaction.ts';
import {
    TransactionExploreQuery
} from '@/models/explore.ts';

import {
    isDefined,
    isNumber,
    isInteger,
    isEquals,
} from '@/lib/common.ts';
import {
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithTimezoneOffset,
    getYearFirstUnixTimeBySpecifiedUnixTime,
    getQuarterFirstUnixTimeBySpecifiedUnixTime,
    getMonthFirstUnixTimeBySpecifiedUnixTime,
    getDayFirstUnixTimeBySpecifiedUnixTime,
    getDateRangeByDateType,
    getFiscalYearStartUnixTime
} from '@/lib/datetime.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export enum TransactionExploreDimensionType {
    TransactionType = 'transactionType',
    Category = 'category',
    Account = 'account',
    Amount = 'amount',
    Other = 'other'
}

export interface TransactionExplorePartialFilter {
    dateRangeType?: number;
    startTime?: number;
    endTime?: number;
    queryId?: string;
    chartType?: TransactionExploreChartTypeValue;
    categoryDimension?: TransactionExploreDataDimensionType;
    seriesDimension?: TransactionExploreDataDimensionType;
    valueMetric?: TransactionExploreValueMetricType;
}

export interface TransactionExploreFilter extends TransactionExplorePartialFilter {
    dateRangeType: number;
    startTime: number;
    endTime: number;
    query: TransactionExploreQuery[];
    chartType: TransactionExploreChartTypeValue;
    categoryDimension: TransactionExploreDataDimensionType;
    seriesDimension: TransactionExploreDataDimensionType;
    valueMetric: TransactionExploreValueMetricType;
}

export interface CategoriedInfo {
    categoryName: string;
    categoryNameNeedI18n?: boolean;
    categoryNameI18nParameters?: Record<string, string>;
    categoryId: string;
    categoryIdType: TransactionExploreDimensionType;
}

export interface CategoriedTransactions extends CategoriedInfo {
    trasactions: Record<string, SeriesedTransactions>;
}

export interface CategoriedTransactionExploreData extends CategoriedInfo {
    data: CategoriedTransactionExploreDataItem[];
}

export interface SeriesedInfo {
    seriesName: string;
    seriesNameNeedI18n?: boolean;
    seriesNameI18nParameters?: Record<string, string>;
    seriesId: string;
    seriesIdType: TransactionExploreDimensionType;
}

export interface SeriesedTransactions extends SeriesedInfo {
    trasactions: TransactionInsightDataItem[];
}

export interface CategoriedTransactionExploreDataItem extends SeriesedInfo {
    value: number;
}

export const useExploresStore = defineStore('explores', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    function getDataCategoryInfo(dimension: TransactionExploreDataDimension, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): CategoriedInfo {
        let transactionTimeUtfOffset: number | undefined = undefined;

        if (settingsStore.appSettings.timezoneUsedForInsightsExplorePage === TimezoneTypeForStatistics.TransactionTimezone.type) {
            transactionTimeUtfOffset = transaction.utcOffset;
        }

        if (dimension === TransactionExploreDataDimension.None) {
            const valueMetric = TransactionExploreValueMetric.valueOf(transactionExploreFilter.value.valueMetric);
            return {
                categoryName: valueMetric?.name ?? 'Unknown',
                categoryNameNeedI18n: true,
                categoryId: 'none',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.Query) {
            if (queryName) {
                return {
                    categoryName: queryName,
                    categoryId: (queryIndex + 1).toString(10),
                    categoryIdType: TransactionExploreDimensionType.Other
                };
            } else {
                return {
                    categoryName: `format.misc.queryIndex`,
                    categoryNameNeedI18n: true,
                    categoryNameI18nParameters: {
                        index: (queryIndex + 1).toString(10)
                    },
                    categoryId: (queryIndex + 1).toString(10),
                    categoryIdType: TransactionExploreDimensionType.Other
                };
            }
        } else if (dimension === TransactionExploreDataDimension.DateTime) {
            const unixTime = transaction.time.toString(10);

            return {
                categoryName: unixTime,
                categoryId: unixTime,
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByYearMonthDay) {
            const unixTime = getDayFirstUnixTimeBySpecifiedUnixTime(transaction.time, transactionTimeUtfOffset).toString(10);

            return {
                categoryName: unixTime,
                categoryId: unixTime,
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByYearMonth) {
            const unixTime = getMonthFirstUnixTimeBySpecifiedUnixTime(transaction.time, transactionTimeUtfOffset).toString(10);

            return {
                categoryName: unixTime,
                categoryId: unixTime,
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByYearQuarter) {
            const unixTime = getQuarterFirstUnixTimeBySpecifiedUnixTime(transaction.time, transactionTimeUtfOffset).toString(10);

            return {
                categoryName: unixTime,
                categoryId: unixTime,
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByYear) {
            const unixTime = getYearFirstUnixTimeBySpecifiedUnixTime(transaction.time, transactionTimeUtfOffset).toString(10);

            return {
                categoryName: unixTime,
                categoryId: unixTime,
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByFiscalYear) {
            const unixTime = getFiscalYearStartUnixTime(transaction.time, userStore.currentUserFiscalYearStart, transactionTimeUtfOffset).toString(10);

            return {
                categoryName: unixTime,
                categoryId: unixTime,
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByDayOfWeek) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getWeekDay().name,
                categoryId: dateTime.getWeekDay().type.toString(10),
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByDayOfMonth) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarDay().toString(10),
                categoryId: dateTime.getGregorianCalendarDay().toString(10),
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByMonthOfYear) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarMonth().toString(10),
                categoryId: dateTime.getGregorianCalendarMonth().toString(10),
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DateTimeByQuarterOfYear) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getGregorianCalendarQuarter().toString(10),
                categoryId: dateTime.getGregorianCalendarQuarter().toString(10),
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.TransactionType) {
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
                categoryIdType: TransactionExploreDimensionType.TransactionType
            };
        } else if (dimension === TransactionExploreDataDimension.SourceAccount) {
            return {
                categoryName: transaction.sourceAccountName || 'Unknown',
                categoryNameNeedI18n: !transaction.sourceAccountName,
                categoryId: transaction.sourceAccountId || 'unknown',
                categoryIdType: TransactionExploreDimensionType.Account
            };
        } else if (dimension === TransactionExploreDataDimension.SourceAccountCategory) {
            const accountCategory = AccountCategory.valueOf(transaction.sourceAccount.category);

            return {
                categoryName: accountCategory?.name || 'Unknown',
                categoryNameNeedI18n: true,
                categoryId: accountCategory?.type.toString(10) || 'unknown',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.SourceAccountCurrency) {
            return {
                categoryName: transaction.sourceAccount.currency || 'Unknown',
                categoryNameNeedI18n: !transaction.sourceAccount.currency,
                categoryId: transaction.sourceAccount.currency || 'unknown',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        }  else if (dimension === TransactionExploreDataDimension.DestinationAccount) {
            return {
                categoryName: transaction.type === TransactionType.Transfer ? (transaction.destinationAccountName || 'Unknown') : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer || !transaction.destinationAccountName,
                categoryId: transaction.type === TransactionType.Transfer ? (transaction.destinationAccountId || 'unknown') : 'none',
                categoryIdType: TransactionExploreDimensionType.Account
            };
        } else if (dimension === TransactionExploreDataDimension.DestinationAccountCategory) {
            const accountCategory = transaction.type === TransactionType.Transfer && transaction.destinationAccount ? AccountCategory.valueOf(transaction.destinationAccount.category) : undefined;

            return {
                categoryName: transaction.type === TransactionType.Transfer ? (accountCategory?.name || 'Unknown') : 'None',
                categoryNameNeedI18n: true,
                categoryId: transaction.type === TransactionType.Transfer ? (accountCategory?.name || 'unknown') : 'none',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.DestinationAccountCurrency) {
            return {
                categoryName: transaction.type === TransactionType.Transfer ? (transaction.destinationAccount?.currency || 'Unknown') : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer || !transaction.destinationAccount?.currency,
                categoryId: transaction.type === TransactionType.Transfer ? (transaction.destinationAccount?.currency || 'unknown') : 'none',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else if (dimension === TransactionExploreDataDimension.PrimaryCategory) {
            return {
                categoryName: transaction.primaryCategory.name,
                categoryId: transaction.primaryCategory.id,
                categoryIdType: TransactionExploreDimensionType.Category
            };
        } else if (dimension === TransactionExploreDataDimension.SecondaryCategory) {
            return {
                categoryName: transaction.secondaryCategory.name,
                categoryId: transaction.categoryId,
                categoryIdType: TransactionExploreDimensionType.Category
            };
        } else if (dimension === TransactionExploreDataDimension.SourceAmount) {
            return {
                categoryName: transaction.sourceAmount.toString(10),
                categoryId: transaction.sourceAmount.toString(10),
                categoryIdType: TransactionExploreDimensionType.Amount
            };
        } else if (dimension === TransactionExploreDataDimension.DestinationAmount) {
            return {
                categoryName: transaction.type === TransactionType.Transfer ? transaction.destinationAmount.toString(10) : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer,
                categoryId: transaction.type === TransactionType.Transfer ? transaction.destinationAmount.toString(10) : 'none',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        } else {
            return {
                categoryName: '',
                categoryId: '',
                categoryIdType: TransactionExploreDimensionType.Other
            };
        }
    }

    function addTransactionToCategoriedDataMap(categoriedDataMap: Record<string, CategoriedTransactions>, categoryDimension: TransactionExploreDataDimension, seriesDemension: TransactionExploreDataDimension, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): void {
        const categoriedInfo = getDataCategoryInfo(categoryDimension, queryName, queryIndex, transaction);
        let categoriedData = categoriedDataMap[categoriedInfo.categoryId];

        if (!categoriedData) {
            categoriedData = {
                categoryName: categoriedInfo.categoryName,
                categoryNameNeedI18n: categoriedInfo.categoryNameNeedI18n,
                categoryNameI18nParameters: categoriedInfo.categoryNameI18nParameters,
                categoryId: categoriedInfo.categoryId,
                categoryIdType: categoriedInfo.categoryIdType,
                trasactions: {}
            };
            categoriedDataMap[categoriedInfo.categoryId] = categoriedData;
        }

        const seriesedInfo = getDataCategoryInfo(seriesDemension, queryName, queryIndex, transaction);
        let seriesedData = categoriedData.trasactions[seriesedInfo.categoryId];

        if (!seriesedData) {
            seriesedData = {
                seriesName: seriesedInfo.categoryName,
                seriesNameNeedI18n: seriesedInfo.categoryNameNeedI18n,
                seriesNameI18nParameters: seriesedInfo.categoryNameI18nParameters,
                seriesId: seriesedInfo.categoryId,
                seriesIdType: seriesedInfo.categoryIdType,
                trasactions: []
            };
            categoriedData.trasactions[seriesedInfo.categoryId] = seriesedData;
        }

        seriesedData.trasactions.push(transaction);
    }

    const transactionExploreFilter = ref<TransactionExploreFilter>({
        dateRangeType: DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE.type,
        startTime: 0,
        endTime: 0,
        query: [],
        categoryDimension: TransactionExploreDataDimension.CategoryDimensionDefault.value,
        seriesDimension: TransactionExploreDataDimension.SeriesDimensionDefault.value,
        valueMetric: TransactionExploreValueMetric.Default.value,
        chartType: TransactionExploreChartType.Default.value
    });

    const transactionExploreAllData = ref<TransactionInfoResponse[]>([]);
    const transactionExploreStateInvalid = ref<boolean>(true);

    const allTransactions = computed<TransactionInsightDataItem[]>(() => {
        if (!transactionExploreAllData.value || transactionExploreAllData.value.length < 1) {
            return [];
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of transactionExploreAllData.value) {
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

        if (!transactionExploreFilter.value.query || transactionExploreFilter.value.query.length < 1) {
            return allTransactions.value;
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of allTransactions.value) {
            for (const query of transactionExploreFilter.value.query) {
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

        const chartType = TransactionExploreChartType.valueOf(transactionExploreFilter.value.chartType);
        const categoryDimension = TransactionExploreDataDimension.valueOf(transactionExploreFilter.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExploreDataDimension.valueOf(transactionExploreFilter.value.seriesDimension) : TransactionExploreDataDimension.SeriesDimensionDefault;

        if (!chartType || !categoryDimension || !seriesDimension) {
            return {};
        }

        const categoriedDataMap: Record<string, CategoriedTransactions> = {};

        for (const transaction of allTransactions.value) {
            if (!transactionExploreFilter.value.query || transactionExploreFilter.value.query.length < 1) {
                addTransactionToCategoriedDataMap(categoriedDataMap, categoryDimension, seriesDimension, '', 0, transaction);
                continue;
            }

            for (const [query, index] of itemAndIndex(transactionExploreFilter.value.query)) {
                if (query.match(transaction)) {
                    addTransactionToCategoriedDataMap(categoriedDataMap, categoryDimension, seriesDimension, query.name, index, transaction);

                    if (categoryDimension !== TransactionExploreDataDimension.Query) {
                        break;
                    }
                }
            }
        }

        return categoriedDataMap;
    });

    const categoriedTransactionExploreData = computed<CategoriedTransactionExploreData[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        const chartType = TransactionExploreChartType.valueOf(transactionExploreFilter.value.chartType);
        const categoryDimension = TransactionExploreDataDimension.valueOf(transactionExploreFilter.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExploreDataDimension.valueOf(transactionExploreFilter.value.seriesDimension) : TransactionExploreDataDimension.SeriesDimensionDefault;
        const valueMetric = TransactionExploreValueMetric.valueOf(transactionExploreFilter.value.valueMetric);

        if (!chartType || !categoryDimension || !seriesDimension || !valueMetric) {
            return [];
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const result: CategoriedTransactionExploreData[] = [];
        const categoriedDataMap = categoriedTransactions.value;

        for (const categoriedTransactions of values(categoriedDataMap)) {
            const dataItems: CategoriedTransactionExploreDataItem[] = [];
            let allSeriesedTransactions: Record<string, SeriesedTransactions> = categoriedTransactions.trasactions;

            // merge all series into one for pie/radar chart
            if (chartType === TransactionExploreChartType.Pie || chartType === TransactionExploreChartType.Radar) {
                const transactions: TransactionInsightDataItem[] = [];

                for (const seriesedTransactions of values(categoriedTransactions.trasactions)) {
                    transactions.push(...seriesedTransactions.trasactions);
                }

                allSeriesedTransactions = {};
                allSeriesedTransactions['none'] = {
                    seriesName: valueMetric?.name ?? 'Unknown',
                    seriesNameNeedI18n: true,
                    seriesId: 'none',
                    seriesIdType: TransactionExploreDimensionType.Other,
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

                if (valueMetric === TransactionExploreValueMetric.TransactionCount) {
                    value = allSourceAmountsInDefaultCurrency.length;
                } else if (valueMetric === TransactionExploreValueMetric.SourceAmountSum) {
                    value = totalSourceAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExploreValueMetric.SourceAmountAverage) {
                    value = allSourceAmountsInDefaultCurrency.length > 0 ? Math.trunc(totalSourceAmountSumInDefaultCurrency / allSourceAmountsInDefaultCurrency.length) : 0;
                } else if (valueMetric === TransactionExploreValueMetric.SourceAmountMedian) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        value = allSourceAmountsInDefaultCurrency[Math.floor(allSourceAmountsInDefaultCurrency.length / 2)] as number;
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExploreValueMetric.SourceAmountMinimum) {
                    value = minimumSourceAmountInDefaultCurrency === Number.MAX_SAFE_INTEGER ? 0 : minimumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExploreValueMetric.SourceAmountMaximum) {
                    value = maximumSourceAmountInDefaultCurrency === Number.MIN_SAFE_INTEGER ? 0 : maximumSourceAmountInDefaultCurrency;
                }

                dataItems.push({
                    seriesName: seriesedTransactions.seriesName,
                    seriesNameNeedI18n: seriesedTransactions.seriesNameNeedI18n,
                    seriesNameI18nParameters: seriesedTransactions.seriesNameI18nParameters,
                    seriesId: seriesedTransactions.seriesId,
                    seriesIdType: seriesedTransactions.seriesIdType,
                    value: value
                });
            }

            result.push({
                categoryName: categoriedTransactions.categoryName,
                categoryNameNeedI18n: categoriedTransactions.categoryNameNeedI18n,
                categoryNameI18nParameters: categoriedTransactions.categoryNameI18nParameters,
                categoryId: categoriedTransactions.categoryId,
                categoryIdType: categoriedTransactions.categoryIdType,
                data: dataItems
            });
        }

        return result;
    });

    function updateTransactionExploreInvalidState(invalidState: boolean): void {
        transactionExploreStateInvalid.value = invalidState;
    }

    function resetTransactionExplores(): void {
        transactionExploreFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE.type;
        transactionExploreFilter.value.startTime = 0;
        transactionExploreFilter.value.endTime = 0;
        transactionExploreFilter.value.query = [];
        transactionExploreFilter.value.chartType = TransactionExploreChartType.Default.value;
        transactionExploreFilter.value.categoryDimension = TransactionExploreDataDimension.CategoryDimensionDefault.value;
        transactionExploreFilter.value.seriesDimension = TransactionExploreDataDimension.SeriesDimensionDefault.value;
        transactionExploreFilter.value.valueMetric = TransactionExploreValueMetric.Default.value;
        transactionExploreAllData.value = [];
        transactionExploreStateInvalid.value = true;
    }

    function initTransactionExploreFilter(filter?: TransactionExplorePartialFilter, resetQuery?: boolean): void {
        if (filter && isInteger(filter.dateRangeType)) {
            transactionExploreFilter.value.dateRangeType = filter.dateRangeType;
        } else {
            transactionExploreFilter.value.dateRangeType = settingsStore.appSettings.insightsExploreDefaultDateRangeType;
        }

        let dateRangeTypeValid = true;

        if (!DateRange.isAvailableForScene(transactionExploreFilter.value.dateRangeType, DateRangeScene.InsightsExplore)) {
            transactionExploreFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE.type;
            dateRangeTypeValid = false;
        }

        if (dateRangeTypeValid && transactionExploreFilter.value.dateRangeType === DateRange.Custom.type) {
            if (filter && isInteger(filter.startTime)) {
                transactionExploreFilter.value.startTime = filter.startTime;
            } else {
                transactionExploreFilter.value.startTime = 0;
            }

            if (filter && isInteger(filter.endTime)) {
                transactionExploreFilter.value.endTime = filter.endTime;
            } else {
                transactionExploreFilter.value.endTime = 0;
            }
        } else {
            const dateRange = getDateRangeByDateType(transactionExploreFilter.value.dateRangeType, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart);

            if (dateRange) {
                transactionExploreFilter.value.dateRangeType = dateRange.dateType;
                transactionExploreFilter.value.startTime = dateRange.minTime;
                transactionExploreFilter.value.endTime = dateRange.maxTime;
            }
        }

        if (resetQuery) {
            transactionExploreFilter.value.query = [];
            transactionExploreFilter.value.chartType = TransactionExploreChartType.Default.value;
            transactionExploreFilter.value.categoryDimension = TransactionExploreDataDimension.CategoryDimensionDefault.value;
            transactionExploreFilter.value.seriesDimension = TransactionExploreDataDimension.SeriesDimensionDefault.value;
            transactionExploreFilter.value.valueMetric = TransactionExploreValueMetric.Default.value;
        }
    }

    function updateTransactionExploreFilter(filter: TransactionExplorePartialFilter): boolean {
        let changed = false;

        if (filter && isInteger(filter.dateRangeType) && transactionExploreFilter.value.dateRangeType !== filter.dateRangeType) {
            transactionExploreFilter.value.dateRangeType = filter.dateRangeType;
            changed = true;
        }

        if (filter && isInteger(filter.startTime) && transactionExploreFilter.value.startTime !== filter.startTime) {
            transactionExploreFilter.value.startTime = filter.startTime;
            changed = true;
        }

        if (filter && isInteger(filter.endTime) && transactionExploreFilter.value.endTime !== filter.endTime) {
            transactionExploreFilter.value.endTime = filter.endTime;
            changed = true;
        }

        if (filter && isDefined(filter.chartType) && transactionExploreFilter.value.chartType !== filter.chartType) {
            transactionExploreFilter.value.chartType = filter.chartType;
            changed = true;
        }

        if (filter && isDefined(filter.categoryDimension) && transactionExploreFilter.value.categoryDimension !== filter.categoryDimension) {
            transactionExploreFilter.value.categoryDimension = filter.categoryDimension;
            changed = true;
        }

        if (filter && isDefined(filter.seriesDimension) && transactionExploreFilter.value.seriesDimension !== filter.seriesDimension) {
            transactionExploreFilter.value.seriesDimension = filter.seriesDimension;
            changed = true;
        }

        if (filter && isDefined(filter.valueMetric) && transactionExploreFilter.value.valueMetric !== filter.valueMetric) {
            transactionExploreFilter.value.valueMetric = filter.valueMetric;
            changed = true;
        }

        if (transactionExploreFilter.value.seriesDimension === transactionExploreFilter.value.categoryDimension && transactionExploreFilter.value.seriesDimension !== TransactionExploreDataDimension.SeriesDimensionDefault.value) {
            transactionExploreFilter.value.seriesDimension = TransactionExploreDataDimension.SeriesDimensionDefault.value;
            changed = true;
        }

        return changed;
    }

    function getTransactionExplorePageParams(currentExploreId: string, activeTab: string): string {
        const querys: string[] = [];

        if (currentExploreId) {
            querys.push('id=' + currentExploreId);
        }

        if (activeTab) {
            querys.push('activeTab=' + activeTab);
        }

        querys.push('dateRangeType=' + transactionExploreFilter.value.dateRangeType);
        querys.push('startTime=' + transactionExploreFilter.value.startTime);
        querys.push('endTime=' + transactionExploreFilter.value.endTime);

        return querys.join('&');
    }

    function getTransactionListPageParams(dimensionType: TransactionExploreDimensionType, itemId: string): string {
        const querys: string[] = [];

        if (dimensionType === TransactionExploreDimensionType.TransactionType) {
            querys.push(`type=${itemId}`);
        } else if (dimensionType === TransactionExploreDimensionType.Account) {
            querys.push(`accountIds=${itemId}`);
        } else if (dimensionType === TransactionExploreDimensionType.Category) {
            querys.push(`categoryIds=${itemId}`);
        } else if (dimensionType === TransactionExploreDimensionType.Amount) {
            querys.push(`amountFilter=${encodeURIComponent(AmountFilterType.EqualTo.toTextualFilter(parseInt(itemId)))}`);
        } else {
            return '';
        }

        querys.push('dateType=' + transactionExploreFilter.value.dateRangeType);
        querys.push('minTime=' + transactionExploreFilter.value.startTime);
        querys.push('maxTime=' + transactionExploreFilter.value.endTime);

        return querys.join('&');
    }

    function loadAllTransactions({ force }: { force: boolean }): Promise<TransactionInfoResponse[]> {
        return new Promise((resolve, reject) => {
            services.getAllTransactions({
                startTime: transactionExploreFilter.value.startTime,
                endTime: transactionExploreFilter.value.endTime
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve all transactions' });
                    return;
                }

                if (transactionExploreStateInvalid.value) {
                    updateTransactionExploreInvalidState(false);
                }

                if (force && data.result && isEquals(transactionExploreAllData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionExploreAllData.value = data.result;

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
        transactionExploreFilter,
        transactionExploreStateInvalid,
        // computed
        filteredTransactions,
        categoriedTransactionExploreData,
        // functions
        updateTransactionExploreInvalidState,
        resetTransactionExplores,
        initTransactionExploreFilter,
        updateTransactionExploreFilter,
        getTransactionExplorePageParams,
        getTransactionListPageParams,
        loadAllTransactions
    };
});
