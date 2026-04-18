import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useTransactionTagsStore } from './transactionTag.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { type BeforeResolveFunction, itemAndIndex, keys, values } from '@/core/base.ts';
import { NumeralSystem, AmountFilterType } from '@/core/numeral.ts';
import { type DateTime, DateRangeScene, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { AccountCategory } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    TransactionExplorerChartType,
    TransactionExplorerDataDimension,
    TransactionExplorerValueMetric,
    DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE
} from '@/core/explorer.ts';
import { AMOUNT_FACTOR } from '@/consts/numeral.ts';
import { ALL_CURRENCIES } from '@/consts/currency.ts';

import { type Account } from '@/models/account.ts';
import { type TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTag } from '@/models/transaction_tag.ts';
import {
    type TransactionInfoResponse,
    type TransactionInsightDataItem
} from '@/models/transaction.ts';
import {
    type InsightsExplorerNewDisplayOrderRequest,
    type InsightsExplorerInfoResponse,
    type InsightsExplorerMatchContext,
    InsightsExplorer,
    InsightsExplorerBasicInfo
} from '@/models/explorer.ts';

import {
    isDefined,
    isNumber,
    isInteger,
    isEquals,
    getObjectOwnFieldCount
} from '@/lib/common.ts';
import {
    mean,
    median,
    percentile,
    sumMaxN,
    cumulativePercentage,
    meanAbsoluteDeviation,
    medianAbsoluteDeviation,
    varianceAndStandardDeviation,
    coefficientOfVariation,
    skewness,
    kurtosis
} from '@/lib/math.ts';
import {
    getUtcOffsetByUtcOffsetMinutes,
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithTimezoneOffset,
    getDateRangeByDateType,
    getFiscalYearFromUnixTime
} from '@/lib/datetime.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
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
}

export interface TransactionExplorerFilter extends TransactionExplorerPartialFilter {
    dateRangeType: number;
    startTime: number;
    endTime: number;
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
    trasactions: Record<string, SeriesTransactions>;
}

export interface CategoriedTransactionExplorerData extends CategoriedInfo {
    data: CategoriedTransactionExplorerDataItem[];
}

export interface SeriesInfo {
    seriesName: string;
    seriesNameNeedI18n?: boolean;
    seriesNameI18nParameters?: Record<string, string>;
    seriesId: string;
    seriesIdType: TransactionExplorerDimensionType;
    seriesDisplayOrders: number[];
}

export interface SeriesTransactions extends SeriesInfo {
    trasactions: TransactionInsightDataItem[];
}

export interface CategoriedTransactionExplorerDataItem extends SeriesInfo {
    value: number;
}

export interface AmountRanges {
    categorySourceAmountRanges?: number[];
    categoryDestinationAmountRanges?: number[];
    seriesSourceAmountRanges?: number[];
    seriesDestinationAmountRanges?: number[];
}

export interface TransactionInsightDataItemInQuery {
    queryIndex: number;
    queryName: string;
    transaction: TransactionInsightDataItem;
}

export interface InsightsExplorerTransactionStatisticData {
    totalCount: number;
    totalAmount: number;
    totalIncome: number;
    totalExpense: number;
    netIncome: number;
    averageAmount: number;
    medianAmount: number;
    minimumAmount: number;
    maximumAmount: number;
    p90Amount: number;
    range: number;
    interquartileRange: number;
    top5AmountShare?: number;
    transactionsFor80PercentAmount?: number;
    variance?: number;
    standardDeviation?: number;
    coefficientOfVariation?: number;
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

    function buildInsightsExplorerMatchContext(insightsExplorer: InsightsExplorer, transaction: TransactionInsightDataItem): InsightsExplorerMatchContext {
        return {
            getTransactionDateTime(): DateTime {
                let transactionTimeUtfOffset: number | undefined = undefined;

                if (insightsExplorer.timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
                    transactionTimeUtfOffset = transaction.utcOffset;
                }

                return isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
            }
        };
    }

    function calculateAmountRanges(sortedAmounts: number[], dimension: TransactionExplorerDataDimension, rangeCount: number): number[] {
        const result: number[] = [];

        if (sortedAmounts.length < 1 || rangeCount <= 0) {
            return result;
        }

        const minAmount = sortedAmounts[0] as number;
        const maxAmount = sortedAmounts[sortedAmounts.length - 1] as number;
        rangeCount = Math.min(rangeCount, sortedAmounts.length);

        // [min1, max1), [min2, max2), ..., [minN, maxN]
        if (dimension === TransactionExplorerDataDimension.SourceAmountRangeEqualFrequency
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeEqualFrequency) {
            for (let i = 0; i < rangeCount; i++) {
                result.push(sortedAmounts[Math.floor(i * (sortedAmounts.length - 1) / rangeCount)] as number);
            }
            result.push(maxAmount);
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeEqualWidth
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeEqualWidth) {
            if (minAmount === maxAmount) {
                return [minAmount, maxAmount];
            }

            const width: number = (maxAmount - minAmount) / rangeCount;

            for (let i = 0; i < rangeCount; i++) {
                result.push(minAmount + i * width);
            }
            result.push(maxAmount);
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeLogScale
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeLogScale) {
            const epsilon: number = 1e-9;

            const transform = (x: number): number => {
                if (x === 0) {
                    return 0;
                }

                return Math.sign(x) * Math.log(Math.abs(x) + epsilon);
            };

            const inverse = (y: number): number => {
                if (y === 0) {
                    return 0;
                }

                return Math.sign(y) * (Math.exp(Math.abs(y)) - epsilon);
            };

            const transformed = sortedAmounts.map(transform).sort((a, b) => a - b);

            const tMin: number = transformed[0] as number;
            const tMax: number = transformed[transformed.length - 1] as number;

            if (tMin === tMax) {
                return [minAmount, maxAmount];
            }

            const width: number = (tMax - tMin) / rangeCount;

            result.push(minAmount);
            for (let i = 1; i < rangeCount; i++) {
                result.push(inverse(tMin + i * width));
            }
            result.push(maxAmount);
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeStandardDeviation
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeStandardDeviation) {
            if (minAmount === maxAmount) {
                return [minAmount, maxAmount];
            }

            const averageAmountForVarianceCalculation: number = mean(sortedAmounts, item => item) / AMOUNT_FACTOR;
            const { standardDeviation } = varianceAndStandardDeviation(sortedAmounts, averageAmountForVarianceCalculation, item => item / AMOUNT_FACTOR);

            if (standardDeviation === 0) {
                return [minAmount, maxAmount];
            }

            const rawBreaks: number[] = [];
            const halfCount = Math.floor(rangeCount / 2);

            if (rangeCount % 2 === 1) {
                for (let i = -halfCount; i <= halfCount; i++) {
                    rawBreaks.push((averageAmountForVarianceCalculation + i * standardDeviation) * AMOUNT_FACTOR);
                }
            } else {
                for (let i = -halfCount; i <= halfCount; i++) {
                    if (i === 0) {
                        continue;
                    }
                    rawBreaks.push((averageAmountForVarianceCalculation + (i - 0.5) * standardDeviation) * AMOUNT_FACTOR);
                }
                rawBreaks.sort((a, b) => a - b);
            }

            const clipped = rawBreaks.map((v) => Math.max(minAmount, Math.min(maxAmount, v)))
                .filter((v, i, arr) => i === 0 || v !== arr[i - 1]);

            clipped[0] = minAmount;

            if (clipped[clipped.length - 1] !== maxAmount) {
                clipped.push(maxAmount);
            }

            return clipped;
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeNaturalBreaks
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeNaturalBreaks) {
            if (minAmount === maxAmount) {
                return [minAmount, maxAmount];
            }

            const n = sortedAmounts.length;
            const k = Math.min(rangeCount, n);

            const lowerClassLimits: number[][] = Array.from({ length: n + 1 }, () => new Array(k + 1).fill(0));
            const varianceCombinations: number[][] = Array.from({ length: n + 1 }, () => new Array(k + 1).fill(Infinity));

            for (let i = 1; i <= k; i++) {
                lowerClassLimits[1]![i] = 1;
                varianceCombinations[1]![i] = 0;
            }

            for (let l = 2; l <= n; l++) {
                let sumZ = 0;
                let sumZ2 = 0;

                for (let m = 1; m <= l; m++) {
                    const val = sortedAmounts[l - m] as number;
                    sumZ += val;
                    sumZ2 += val * val;

                    const variance = sumZ2 - (sumZ * sumZ) / m;

                    if (m === l) {
                        for (let j = 1; j <= k; j++) {
                            if (variance < varianceCombinations[l]![j]!) {
                                lowerClassLimits[l]![j] = 1;
                                varianceCombinations[l]![j] = variance;
                            }
                        }
                    } else {
                        for (let j = 2; j <= k; j++) {
                            const combined = varianceCombinations[l - m]![j - 1]! + variance;
                            if (combined < varianceCombinations[l]![j]!) {
                                lowerClassLimits[l]![j] = l - m + 1;
                                varianceCombinations[l]![j] = combined;
                            }
                        }
                    }
                }
            }

            const breaks: number[] = new Array(k + 1);
            breaks[k] = maxAmount;

            let currentK = k;
            let currentIdx = n;

            while (currentK >= 2) {
                const lowerIdx = lowerClassLimits[currentIdx]![currentK]!;
                breaks[currentK - 1] = sortedAmounts[lowerIdx - 1] as number;
                currentIdx = lowerIdx - 1;
                currentK--;
            }

            breaks[0] = minAmount;
            return breaks;
        }

        return result;
    }

    function getDataCategoryInfo(timezoneUsedForDateRange: number, dimension: TransactionExplorerDataDimension, sourceAmountRanges: number[] | undefined, destinationAmountRanges: number[] | undefined, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): CategoriedInfo {
        const defaultCurrency = userStore.currentUserDefaultCurrency;
        let transactionTimeUtfOffset: number | undefined = undefined;

        if (timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
            transactionTimeUtfOffset = transaction.utcOffset;
        }

        if (dimension === TransactionExplorerDataDimension.None) {
            const valueMetric = TransactionExplorerValueMetric.valueOf(currentInsightsExplorer.value.valueMetric);
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
            const textualDateTime = `${dateTime.getGregorianCalendarYearDashMonthDashDay()} ${dateTime.getHour().toString(10).padStart(2, NumeralSystem.WesternArabicNumerals.digitZero)}:${dateTime.getMinute().toString(10).padStart(2, NumeralSystem.WesternArabicNumerals.digitZero)}:${dateTime.getSecond().toString(10).padStart(2, NumeralSystem.WesternArabicNumerals.digitZero)}`;

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
        } else if (dimension === TransactionExplorerDataDimension.DateTimeByHourOfDay) {
            const dateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);

            return {
                categoryName: dateTime.getHour().toString(10),
                categoryId: dateTime.getHour().toString(10),
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [dateTime.getHour()]
            };
        } else if (dimension === TransactionExplorerDataDimension.TimezoneOffset) {
            return {
                categoryName: getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset),
                categoryId: transaction.utcOffset.toString(10),
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [transaction.utcOffset]
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
                categoryDisplayOrders: [currencyDisplayOrders[transaction.sourceAccount.currency] || Number.MAX_SAFE_INTEGER]
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
                categoryDisplayOrders: transaction.type === TransactionType.Transfer ? [accountCategoryDisplayOrder] : [Number.MAX_SAFE_INTEGER]
            };
        } else if (dimension === TransactionExplorerDataDimension.DestinationAccountCurrency) {
            return {
                categoryName: transaction.type === TransactionType.Transfer ? (transaction.destinationAccount?.currency || 'Unknown') : 'None',
                categoryNameNeedI18n: transaction.type !== TransactionType.Transfer || !transaction.destinationAccount?.currency,
                categoryId: transaction.type === TransactionType.Transfer ? (transaction.destinationAccount?.currency || 'unknown') : 'none',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: transaction.type === TransactionType.Transfer ? [currencyDisplayOrders[transaction.destinationAccount?.currency ?? ''] || Number.MAX_SAFE_INTEGER] : [Number.MAX_SAFE_INTEGER]
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
        } else if (dimension === TransactionExplorerDataDimension.SourceAmount || dimension === TransactionExplorerDataDimension.DestinationAmount) {
            if (dimension === TransactionExplorerDataDimension.DestinationAmount && transaction.type !== TransactionType.Transfer) {
                return {
                    categoryName: 'None',
                    categoryNameNeedI18n: true,
                    categoryId: 'none',
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                };
            }

            const amount = dimension === TransactionExplorerDataDimension.SourceAmount ? transaction.sourceAmount : transaction.destinationAmount;
            const account = dimension === TransactionExplorerDataDimension.SourceAmount ? transaction.sourceAccount : transaction.destinationAccount;
            let amountInDefaultCurrency: number = amount;

            if (!account) {
                return {
                    categoryName: 'Unknown',
                    categoryNameNeedI18n: true,
                    categoryId: 'unknown',
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                };
            }

            if (account.currency !== defaultCurrency) {
                const exchangedAmount = exchangeRatesStore.getExchangedAmount(amount, account.currency, defaultCurrency);

                if (isNumber(exchangedAmount)) {
                    amountInDefaultCurrency = Math.trunc(exchangedAmount);
                } else {
                    return {
                        categoryName: 'Unknown',
                        categoryNameNeedI18n: true,
                        categoryId: 'unknown',
                        categoryIdType: TransactionExplorerDimensionType.Other,
                        categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                    };
                }
            }

            return {
                categoryName: amountInDefaultCurrency.toString(10),
                categoryId: amountInDefaultCurrency.toString(10),
                categoryIdType: TransactionExplorerDimensionType.Amount,
                categoryDisplayOrders: [amountInDefaultCurrency]
            };
        } else if (dimension.isSourceAmountRange || dimension.isDestinationAmountRange) {
            const isSourceAmount = dimension.isSourceAmountRange;

            if (dimension.isDestinationAmountRange && transaction.type !== TransactionType.Transfer) {
                return {
                    categoryName: 'None',
                    categoryNameNeedI18n: true,
                    categoryId: 'none',
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                };
            }

            const amount = dimension.isSourceAmountRange ? transaction.sourceAmount : transaction.destinationAmount;
            const account = dimension.isSourceAmountRange ? transaction.sourceAccount : transaction.destinationAccount;
            let amountInDefaultCurrency: number = amount;

            if (!account) {
                return {
                    categoryName: 'Unknown',
                    categoryNameNeedI18n: true,
                    categoryId: 'unknown',
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                };
            }

            if (account.currency !== defaultCurrency) {
                const exchangedAmount = exchangeRatesStore.getExchangedAmount(amount, account.currency, defaultCurrency);

                if (isNumber(exchangedAmount)) {
                    amountInDefaultCurrency = Math.trunc(exchangedAmount);
                } else {
                    return {
                        categoryName: 'Unknown',
                        categoryNameNeedI18n: true,
                        categoryId: 'unknown',
                        categoryIdType: TransactionExplorerDimensionType.Other,
                        categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                    };
                }
            }

            const amountRanges: number[] = isSourceAmount ? (sourceAmountRanges ?? []) : (destinationAmountRanges ?? []);
            let matchAmountRangeMin: number | undefined = undefined;
            let matchAmountRangeMax: number | undefined = undefined;
            let matchAmountRangeIndex: number | undefined = undefined;

            for (let i = 1; i < amountRanges.length; i++) {
                const amountRangeMin = amountRanges[i - 1] as number;
                const amountRangeMax = amountRanges[i] as number;

                if (amountInDefaultCurrency < amountRangeMin) {
                    continue;
                }

                if (amountInDefaultCurrency > amountRangeMax) {
                    continue;
                }

                if (i < amountRanges.length - 1 && amountInDefaultCurrency === amountRangeMax) {
                    continue;
                }

                matchAmountRangeMin = amountRangeMin;
                matchAmountRangeMax = amountRangeMax;
                matchAmountRangeIndex = i - 1;
            }

            if (isNumber(matchAmountRangeMin) && isNumber(matchAmountRangeMax) && isNumber(matchAmountRangeIndex)) {
                return {
                    categoryName: `${matchAmountRangeMin.toString(10)}|${matchAmountRangeMax.toString(10)}`,
                    categoryId: matchAmountRangeIndex.toString(10),
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [matchAmountRangeIndex]
                };
            } else {
                return {
                    categoryName: 'Other',
                    categoryNameNeedI18n: true,
                    categoryId: 'other',
                    categoryIdType: TransactionExplorerDimensionType.Other,
                    categoryDisplayOrders: [Number.MAX_SAFE_INTEGER]
                };
            }
        } else {
            return {
                categoryName: '',
                categoryId: '',
                categoryIdType: TransactionExplorerDimensionType.Other,
                categoryDisplayOrders: [0]
            };
        }
    }

    function addTransactionToFilteredList(filteredTransactions: TransactionInsightDataItemInQuery[], filteredTransactionSourceAmountsInDefaultCurrency: number[], filteredTransactionDestinationAmountsInDefaultCurrency: number[], defaultCurrency: string, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): void {
        filteredTransactions.push({
            queryIndex: queryIndex,
            queryName: queryName,
            transaction: transaction
        });

        let sourceAmountInDefaultCurrency: number | undefined = transaction.sourceAmount;
        let destinationAmountInDefaultCurrency: number | undefined = transaction.type === TransactionType.Transfer && transaction.destinationAccount ? transaction.destinationAmount : undefined;

        if (transaction.sourceAccount.currency !== defaultCurrency) {
            const amount = exchangeRatesStore.getExchangedAmount(transaction.sourceAmount, transaction.sourceAccount.currency, defaultCurrency);
            sourceAmountInDefaultCurrency = isNumber(amount) ? Math.trunc(amount) : undefined;
        }

        if (transaction.type === TransactionType.Transfer && transaction.destinationAccount && transaction.destinationAccount.currency !== defaultCurrency) {
            const amount = exchangeRatesStore.getExchangedAmount(transaction.destinationAmount, transaction.destinationAccount.currency, defaultCurrency);
            destinationAmountInDefaultCurrency = isNumber(amount) ? Math.trunc(amount) : undefined;
        }

        if (isNumber(sourceAmountInDefaultCurrency)) {
            filteredTransactionSourceAmountsInDefaultCurrency.push(sourceAmountInDefaultCurrency);
        }

        if (isNumber(destinationAmountInDefaultCurrency)) {
            filteredTransactionDestinationAmountsInDefaultCurrency.push(destinationAmountInDefaultCurrency);
        }
    }

    function addTransactionToCategoriedDataMap(timezoneUsedForDateRange: number, categoriedDataMap: Record<string, CategoriedTransactions>, categoryDimension: TransactionExplorerDataDimension, seriesDemension: TransactionExplorerDataDimension, allAmountRanges: AmountRanges, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): void {
        const categoriedInfo = getDataCategoryInfo(timezoneUsedForDateRange, categoryDimension, allAmountRanges.categorySourceAmountRanges, allAmountRanges.categoryDestinationAmountRanges, queryName, queryIndex, transaction);
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

        const seriesInfo = getDataCategoryInfo(timezoneUsedForDateRange, seriesDemension, allAmountRanges.seriesSourceAmountRanges, allAmountRanges.seriesDestinationAmountRanges, queryName, queryIndex, transaction);
        let seriesData = categoriedData.trasactions[seriesInfo.categoryId];

        if (!seriesData) {
            seriesData = {
                seriesName: seriesInfo.categoryName,
                seriesNameNeedI18n: seriesInfo.categoryNameNeedI18n,
                seriesNameI18nParameters: seriesInfo.categoryNameI18nParameters,
                seriesId: seriesInfo.categoryId,
                seriesIdType: seriesInfo.categoryIdType,
                seriesDisplayOrders: seriesInfo.categoryDisplayOrders,
                trasactions: []
            };
            categoriedData.trasactions[seriesInfo.categoryId] = seriesData;
        }

        seriesData.trasactions.push(transaction);
    }

    function buildAllAmountRanges(categoryDimension: TransactionExplorerDataDimension, seriesDimension: TransactionExplorerDataDimension, filteredTransactionSourceAmountsInDefaultCurrency: number[], filteredTransactionDestinationAmountsInDefaultCurrency: number[], rangeCount: number): AmountRanges {
        const allAmountRanges: AmountRanges = {};

        if (categoryDimension.isSourceAmountRange || seriesDimension.isSourceAmountRange) {
            filteredTransactionSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
            const sorteUniqueAmounts = filteredTransactionSourceAmountsInDefaultCurrency.filter((v, i, a) => i === 0 || v !== a[i - 1]);

            if (categoryDimension.isSourceAmountRange) {
                allAmountRanges.categorySourceAmountRanges = calculateAmountRanges(sorteUniqueAmounts, categoryDimension, rangeCount);
            }

            if (seriesDimension.isSourceAmountRange) {
                allAmountRanges.seriesSourceAmountRanges = calculateAmountRanges(sorteUniqueAmounts, seriesDimension, rangeCount);
            }
        }

        if (categoryDimension.isDestinationAmountRange || seriesDimension.isDestinationAmountRange) {
            filteredTransactionDestinationAmountsInDefaultCurrency.sort((a, b) => a - b);
            const sorteUniqueAmounts = filteredTransactionDestinationAmountsInDefaultCurrency.filter((v, i, a) => i === 0 || v !== a[i - 1]);
            if (categoryDimension.isDestinationAmountRange) {
                allAmountRanges.categoryDestinationAmountRanges = calculateAmountRanges(sorteUniqueAmounts, categoryDimension, rangeCount);
            }

            if (seriesDimension.isDestinationAmountRange) {
                allAmountRanges.seriesDestinationAmountRanges = calculateAmountRanges(sorteUniqueAmounts, seriesDimension, rangeCount);
            }
        }

        return allAmountRanges;
    }

    function loadInsightsExplorerList(explorers: InsightsExplorerBasicInfo[]): void {
        allInsightsExplorerBasicInfos.value = explorers;
        allInsightsExplorerBasicInfosMap.value = {};

        for (const explorer of explorers) {
            allInsightsExplorerBasicInfosMap.value[explorer.id] = explorer;
        }
    }

    function addExplorerToInsightsExplorerList(explorer: InsightsExplorerBasicInfo): void {
        allInsightsExplorerBasicInfos.value.push(explorer);
        allInsightsExplorerBasicInfosMap.value[explorer.id] = explorer;
    }

    function updateExplorerInInsightsExplorerList(currentExplorer: InsightsExplorerBasicInfo): void {
        for (const [explorer, index] of itemAndIndex(allInsightsExplorerBasicInfos.value)) {
            if (explorer.id === currentExplorer.id) {
                allInsightsExplorerBasicInfos.value.splice(index, 1, currentExplorer);
                break;
            }
        }

        allInsightsExplorerBasicInfosMap.value[currentExplorer.id] = currentExplorer;
    }

    function updateExplorerDisplayOrderInInsightsExplorerList({ from, to }: { from: number, to: number }): void {
        allInsightsExplorerBasicInfos.value.splice(to, 0, allInsightsExplorerBasicInfos.value.splice(from, 1)[0] as InsightsExplorer);
    }

    function updateExplorerVisibilityInInsightsExplorerList({ explorerId, hidden }: { explorerId: string, hidden: boolean }): void {
        if (allInsightsExplorerBasicInfosMap.value[explorerId]) {
            allInsightsExplorerBasicInfosMap.value[explorerId]!.hidden = hidden;
        }
    }

    function removeExplorerFromInsightsExplorerList(currentExplorer: InsightsExplorerBasicInfo): void {
        for (const [insightsExplorer, index] of itemAndIndex(allInsightsExplorerBasicInfos.value)) {
            if (insightsExplorer.id === currentExplorer.id) {
                allInsightsExplorerBasicInfos.value.splice(index, 1);
                break;
            }
        }

        if (allInsightsExplorerBasicInfosMap.value[currentExplorer.id]) {
            delete allInsightsExplorerBasicInfosMap.value[currentExplorer.id];
        }
    }

    const transactionExplorerFilter = ref<TransactionExplorerFilter>({
        dateRangeType: DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type,
        startTime: 0,
        endTime: 0
    });

    const transactionExplorerAllData = ref<TransactionInfoResponse[]>([]);
    const transactionExplorerStateInvalid = ref<boolean>(true);

    const allInsightsExplorerBasicInfos = ref<InsightsExplorerBasicInfo[]>([]);
    const allInsightsExplorerBasicInfosMap = ref<Record<string, InsightsExplorerBasicInfo>>({});
    const currentInsightsExplorer = ref<InsightsExplorer>(InsightsExplorer.createNewExplorer(generateRandomUUID()));
    const insightsExplorerListStateInvalid = ref<boolean>(true);

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

    const isUsingAmountRange = computed<boolean>(() => {
        const chartType = TransactionExplorerChartType.valueOf(currentInsightsExplorer.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(currentInsightsExplorer.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(currentInsightsExplorer.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;
        return categoryDimension?.isSourceAmountRange || seriesDimension?.isSourceAmountRange
            || categoryDimension?.isDestinationAmountRange || seriesDimension?.isDestinationAmountRange
            || false;
    });

    const filteredTransactionsInDataTable = computed<TransactionInsightDataItem[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        if (!currentInsightsExplorer.value.queries || currentInsightsExplorer.value.queries.length < 1) {
            return allTransactions.value;
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of allTransactions.value) {
            const matchOptions: InsightsExplorerMatchContext = buildInsightsExplorerMatchContext(currentInsightsExplorer.value, transaction);

            for (const query of currentInsightsExplorer.value.queries) {
                if (currentInsightsExplorer.value.datatableQuerySource && currentInsightsExplorer.value.datatableQuerySource !== query.id) {
                    continue;
                }

                if (query.match(transaction, matchOptions)) {
                    result.push(transaction);
                    break;
                }
            }
        }

        return result;
    });

    const filteredTransactionsInDataTableStatistic = computed<InsightsExplorerTransactionStatisticData>(() => {
        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const statisticData: InsightsExplorerTransactionStatisticData = {
            totalCount: 0,
            totalAmount: 0,
            totalIncome: 0,
            totalExpense: 0,
            netIncome: 0,
            averageAmount: 0,
            medianAmount: 0,
            minimumAmount: Number.MAX_SAFE_INTEGER,
            maximumAmount: Number.MIN_SAFE_INTEGER,
            p90Amount: 0,
            range: 0,
            interquartileRange: 0,
            top5AmountShare: undefined,
            transactionsFor80PercentAmount: undefined,
            variance: undefined,
            standardDeviation: undefined,
            coefficientOfVariation: undefined
        };

        const sourceAmounts: number[] = [];

        for (const transaction of filteredTransactionsInDataTable.value) {
            statisticData.totalCount++;

            let amountInDefaultCurrency: number = transaction.sourceAmount;

            if (transaction.sourceAccount.currency !== defaultCurrency) {
                const amount = exchangeRatesStore.getExchangedAmount(transaction.sourceAmount, transaction.sourceAccount.currency, defaultCurrency);

                if (isNumber(amount)) {
                    amountInDefaultCurrency = Math.trunc(amount);
                } else {
                    continue;
                }
            }

            statisticData.totalAmount += amountInDefaultCurrency;
            sourceAmounts.push(amountInDefaultCurrency);

            if (transaction.type === TransactionType.Income) {
                statisticData.totalIncome += amountInDefaultCurrency;
            } else if (transaction.type === TransactionType.Expense) {
                statisticData.totalExpense += amountInDefaultCurrency;
            }

            if (amountInDefaultCurrency >= 0 && amountInDefaultCurrency < statisticData.minimumAmount) {
                statisticData.minimumAmount = amountInDefaultCurrency;
            }

            if (amountInDefaultCurrency > statisticData.maximumAmount) {
                statisticData.maximumAmount = amountInDefaultCurrency;
            }
        }

        statisticData.netIncome = statisticData.totalIncome - statisticData.totalExpense;

        if (statisticData.minimumAmount === Number.MAX_SAFE_INTEGER) {
            statisticData.minimumAmount = 0;
        }

        if (statisticData.maximumAmount === Number.MIN_SAFE_INTEGER) {
            statisticData.maximumAmount = 0;
        }

        if (sourceAmounts.length > 0) {
            statisticData.averageAmount = Math.trunc(statisticData.totalAmount / sourceAmounts.length);
        }

        statisticData.range = statisticData.maximumAmount - statisticData.minimumAmount;

        if (sourceAmounts.length > 0) {
            sourceAmounts.sort((a, b) => a - b);
            statisticData.medianAmount = Math.trunc(median(sourceAmounts, item => item));
            statisticData.p90Amount = Math.trunc(percentile(sourceAmounts, 0.9, item => item));

            const q1 = percentile(sourceAmounts, 0.25, item => item);
            const q3 = percentile(sourceAmounts, 0.75, item => item);
            statisticData.interquartileRange = Math.trunc(q3 - q1);
        }

        if (sourceAmounts.length > 5) {
            const top5AmountSum = sumMaxN(sourceAmounts, 5, item => item);
            statisticData.top5AmountShare = statisticData.totalAmount > 0 ? 100.0 * top5AmountSum / statisticData.totalAmount : 0;
        }

        if (sourceAmounts.length > 0) {
            statisticData.transactionsFor80PercentAmount = cumulativePercentage(sourceAmounts, 0.8, statisticData.totalAmount, item => item);
        }

        if (sourceAmounts.length > 0) {
            const averageAmountForVarianceCalculation: number = statisticData.totalAmount / sourceAmounts.length / AMOUNT_FACTOR;
            const { variance, standardDeviation } = varianceAndStandardDeviation(sourceAmounts, averageAmountForVarianceCalculation, item => item / AMOUNT_FACTOR);
            statisticData.variance = variance;
            statisticData.standardDeviation = standardDeviation;
            statisticData.coefficientOfVariation = coefficientOfVariation(standardDeviation, averageAmountForVarianceCalculation);
        }

        return statisticData;
    });

    const categoriedTransactions = computed<Record<string, CategoriedTransactions>>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return {};
        }

        const chartType = TransactionExplorerChartType.valueOf(currentInsightsExplorer.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(currentInsightsExplorer.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(currentInsightsExplorer.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;

        if (!chartType || !categoryDimension || !seriesDimension) {
            return {};
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const filteredTransactions: TransactionInsightDataItemInQuery[] = [];
        const filteredTransactionSourceAmountsInDefaultCurrency: number[] = [];
        const filteredTransactionDestinationAmountsInDefaultCurrency: number[] = [];

        for (const transaction of allTransactions.value) {
            if (!currentInsightsExplorer.value.queries || currentInsightsExplorer.value.queries.length < 1) {
                addTransactionToFilteredList(filteredTransactions, filteredTransactionSourceAmountsInDefaultCurrency, filteredTransactionDestinationAmountsInDefaultCurrency, defaultCurrency, '', 0, transaction);
                continue;
            }

            const matchContext: InsightsExplorerMatchContext = buildInsightsExplorerMatchContext(currentInsightsExplorer.value, transaction);

            for (const [query, index] of itemAndIndex(currentInsightsExplorer.value.queries)) {
                if (query.match(transaction, matchContext)) {
                    addTransactionToFilteredList(filteredTransactions, filteredTransactionSourceAmountsInDefaultCurrency, filteredTransactionDestinationAmountsInDefaultCurrency, defaultCurrency, query.name, index, transaction);

                    if (categoryDimension !== TransactionExplorerDataDimension.Query) {
                        break;
                    }
                }
            }
        }

        const categoriedDataMap: Record<string, CategoriedTransactions> = {};
        const allAmountRanges: AmountRanges = buildAllAmountRanges(categoryDimension, seriesDimension, filteredTransactionSourceAmountsInDefaultCurrency, filteredTransactionDestinationAmountsInDefaultCurrency, currentInsightsExplorer.value.amountRangeCount);

        for (const item of filteredTransactions) {
            addTransactionToCategoriedDataMap(currentInsightsExplorer.value.timezoneUsedForDateRange, categoriedDataMap, categoryDimension, seriesDimension, allAmountRanges, item.queryName, item.queryIndex, item.transaction);
        }

        return categoriedDataMap;
    });

    const categoriedTransactionExplorerData = computed<CategoriedTransactionExplorerData[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        const chartType = TransactionExplorerChartType.valueOf(currentInsightsExplorer.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(currentInsightsExplorer.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(currentInsightsExplorer.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;
        const valueMetric = TransactionExplorerValueMetric.valueOf(currentInsightsExplorer.value.valueMetric);

        if (!chartType || !categoryDimension || !seriesDimension || !valueMetric) {
            return [];
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const result: CategoriedTransactionExplorerData[] = [];
        const categoriedDataMap = categoriedTransactions.value;
        let needCalculateDailyTransactionCount: boolean = false;

        if (valueMetric === TransactionExplorerValueMetric.ActiveTransactionDays || valueMetric === TransactionExplorerValueMetric.TransactionsPerDay) {
            needCalculateDailyTransactionCount = true;
        }

        for (const categoriedTransactions of values(categoriedDataMap)) {
            const dataItems: CategoriedTransactionExplorerDataItem[] = [];
            let allSeriesTransactions: Record<string, SeriesTransactions> = categoriedTransactions.trasactions;

            if (!chartType.seriesDimensionRequired) {
                const transactions: TransactionInsightDataItem[] = [];

                for (const seriesTransactions of values(categoriedTransactions.trasactions)) {
                    transactions.push(...seriesTransactions.trasactions);
                }

                allSeriesTransactions = {};
                allSeriesTransactions['none'] = {
                    seriesName: valueMetric?.name ?? 'Unknown',
                    seriesNameNeedI18n: true,
                    seriesId: 'none',
                    seriesIdType: TransactionExplorerDimensionType.Other,
                    seriesDisplayOrders: [0],
                    trasactions: transactions
                };
            }

            for (const seriesTransactions of values(allSeriesTransactions)) {
                const transactionDateMapCount: Record<string, number> = {};
                const allSourceAmountsInDefaultCurrency: number[] = [];
                let totalSourceAmountSumInDefaultCurrency: number = 0;
                let totalSourceIncomeAmountSumInDefaultCurrency: number = 0;
                let totalSourceExpenseAmountSumInDefaultCurrency: number = 0;
                let minimumSourceAmountInDefaultCurrency: number = Number.MAX_SAFE_INTEGER;
                let maximumSourceAmountInDefaultCurrency: number = Number.MIN_SAFE_INTEGER;

                for (const transaction of seriesTransactions.trasactions) {
                    let amountInDefaultCurrency: number = transaction.sourceAmount;

                    if (transaction.sourceAccount.currency !== defaultCurrency) {
                        const amount = exchangeRatesStore.getExchangedAmount(transaction.sourceAmount, transaction.sourceAccount.currency, defaultCurrency);

                        if (isNumber(amount)) {
                            amountInDefaultCurrency = Math.trunc(amount);
                        } else {
                            continue;
                        }
                    }

                    if (needCalculateDailyTransactionCount) {
                        let transactionTimeUtfOffset: number | undefined = undefined;

                        if (currentInsightsExplorer.value.timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
                            transactionTimeUtfOffset = transaction.utcOffset;
                        }

                        const transactionDateTime: DateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
                        const transactionYearMonthDay: string = transactionDateTime.getGregorianCalendarYearDashMonthDashDay();

                        if (transactionDateMapCount[transactionYearMonthDay]) {
                            transactionDateMapCount[transactionYearMonthDay]++;
                        } else {
                            transactionDateMapCount[transactionYearMonthDay] = 1;
                        }
                    }

                    allSourceAmountsInDefaultCurrency.push(amountInDefaultCurrency);
                    totalSourceAmountSumInDefaultCurrency += amountInDefaultCurrency;

                    if (transaction.type === TransactionType.Income) {
                        totalSourceIncomeAmountSumInDefaultCurrency += amountInDefaultCurrency;
                    } else if (transaction.type === TransactionType.Expense) {
                        totalSourceExpenseAmountSumInDefaultCurrency += amountInDefaultCurrency;
                    }

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
                } else if (valueMetric === TransactionExplorerValueMetric.ActiveTransactionDays) {
                    value = getObjectOwnFieldCount(transactionDateMapCount);
                } else if (valueMetric === TransactionExplorerValueMetric.TransactionsPerDay) {
                    const activeDays = getObjectOwnFieldCount(transactionDateMapCount);
                    value = activeDays > 0 ? allSourceAmountsInDefaultCurrency.length / activeDays : 0;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceIncomeAmountSum) {
                    value = totalSourceIncomeAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceExpenseAmountSum) {
                    value = totalSourceExpenseAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceNetIncomeAmountSum) {
                    value = totalSourceIncomeAmountSumInDefaultCurrency - totalSourceExpenseAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SrouceAmountExpenseIncomeRatio) {
                    value = totalSourceIncomeAmountSumInDefaultCurrency !== 0 ? 100.0 * totalSourceExpenseAmountSumInDefaultCurrency / totalSourceIncomeAmountSumInDefaultCurrency : 0;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSavingsRate) {
                    value = totalSourceIncomeAmountSumInDefaultCurrency !== 0 ? 100.0 * (totalSourceIncomeAmountSumInDefaultCurrency - totalSourceExpenseAmountSumInDefaultCurrency) / totalSourceIncomeAmountSumInDefaultCurrency : 0;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSum) {
                    value = totalSourceAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountAverage) {
                    value = allSourceAmountsInDefaultCurrency.length > 0 ? Math.trunc(totalSourceAmountSumInDefaultCurrency / allSourceAmountsInDefaultCurrency.length) : 0;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMedian) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        value = Math.trunc(median(allSourceAmountsInDefaultCurrency, item => item));
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountQ1Amount
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountQ3Amount
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount10thPercentile
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount90thPercentile
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount95thPercentile
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount99thPercentile) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);

                        if (valueMetric === TransactionExplorerValueMetric.SourceAmountQ1Amount) {
                            value = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.25, item => item));
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountQ3Amount) {
                            value = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.75, item => item));
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount10thPercentile) {
                            value = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.1, item => item));
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount90thPercentile) {
                            value = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.9, item => item));
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount95thPercentile) {
                            value = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.95, item => item));
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount99thPercentile) {
                            value = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.99, item => item));
                        }
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMinimum) {
                    value = minimumSourceAmountInDefaultCurrency === Number.MAX_SAFE_INTEGER ? 0 : minimumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceTop5AmountSum) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        value = sumMaxN(allSourceAmountsInDefaultCurrency, 5, item => item);
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceTop5AmountShare) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        const top5AmountSum = sumMaxN(allSourceAmountsInDefaultCurrency, 5, item => item);
                        value = totalSourceAmountSumInDefaultCurrency > 0 ? 100.0 * top5AmountSum / totalSourceAmountSumInDefaultCurrency : 0;
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.TransactionsForEightyPercentOfSourceAmount) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        value = cumulativePercentage(allSourceAmountsInDefaultCurrency, 0.8, totalSourceAmountSumInDefaultCurrency, item => item);
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMaximum) {
                    value = maximumSourceAmountInDefaultCurrency === Number.MIN_SAFE_INTEGER ? 0 : maximumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountRange) {
                    const finalMinimumSourceAmountInDefaultCurrency = minimumSourceAmountInDefaultCurrency === Number.MAX_SAFE_INTEGER ? 0 : minimumSourceAmountInDefaultCurrency;
                    const finalMaximumSourceAmountInDefaultCurrency = maximumSourceAmountInDefaultCurrency === Number.MIN_SAFE_INTEGER ? 0 : maximumSourceAmountInDefaultCurrency;
                    value = finalMaximumSourceAmountInDefaultCurrency - finalMinimumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountInterquartileRange) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        const q1 = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.25, item => item));
                        const q3 = Math.trunc(percentile(allSourceAmountsInDefaultCurrency, 0.75, item => item));
                        value = Math.trunc(q3 - q1);
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMeanAbsoluteDeviation) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        const averageSourceAmountInDefaultCurrency = totalSourceAmountSumInDefaultCurrency / allSourceAmountsInDefaultCurrency.length;
                        value = Math.trunc(meanAbsoluteDeviation(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, item => item));
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMedianAbsoluteDeviation) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a - b);
                        const medianSourceAmountInDefaultCurrency = median(allSourceAmountsInDefaultCurrency, item => item);
                        value = Math.trunc(medianAbsoluteDeviation(allSourceAmountsInDefaultCurrency, medianSourceAmountInDefaultCurrency, item => item));
                    } else {
                        value = 0;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountVariance
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountStandardDeviation
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountCoefficientOfVariation
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountSkewness
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountKurtosis) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        const averageSourceAmountInDefaultCurrency = totalSourceAmountSumInDefaultCurrency / allSourceAmountsInDefaultCurrency.length / AMOUNT_FACTOR;
                        const { variance, standardDeviation } = varianceAndStandardDeviation(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, item => item / AMOUNT_FACTOR);

                        if (valueMetric === TransactionExplorerValueMetric.SourceAmountVariance) {
                            value = variance;
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountStandardDeviation) {
                            value = standardDeviation;
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountCoefficientOfVariation) {
                            value = coefficientOfVariation(standardDeviation, averageSourceAmountInDefaultCurrency) ?? 0;
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSkewness) {
                            value = skewness(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, standardDeviation, item => item / AMOUNT_FACTOR);
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountKurtosis) {
                            value = kurtosis(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, variance, item => item / AMOUNT_FACTOR);
                        }
                    } else {
                        value = 0;
                    }
                }

                dataItems.push({
                    seriesName: seriesTransactions.seriesName,
                    seriesNameNeedI18n: seriesTransactions.seriesNameNeedI18n,
                    seriesNameI18nParameters: seriesTransactions.seriesNameI18nParameters,
                    seriesId: seriesTransactions.seriesId,
                    seriesIdType: seriesTransactions.seriesIdType,
                    seriesDisplayOrders: seriesTransactions.seriesDisplayOrders,
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

    function updateInsightsExplorerListInvalidState(invalidState: boolean): void {
        insightsExplorerListStateInvalid.value = invalidState;
    }

    function updateCurrentInsightsExplorer(explorer: InsightsExplorer): void {
        currentInsightsExplorer.value = explorer;
    }

    function resetTransactionExplorers(): void {
        transactionExplorerFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type;
        transactionExplorerFilter.value.startTime = 0;
        transactionExplorerFilter.value.endTime = 0;
        transactionExplorerAllData.value = [];
        allInsightsExplorerBasicInfos.value = [];
        allInsightsExplorerBasicInfosMap.value = {};
        currentInsightsExplorer.value = InsightsExplorer.createNewExplorer(generateRandomUUID());
        transactionExplorerStateInvalid.value = true;
        insightsExplorerListStateInvalid.value = true;
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
            currentInsightsExplorer.value = InsightsExplorer.createNewExplorer(generateRandomUUID());
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

        return changed;
    }

    function getTransactionExplorerPageParams(currentExplorerId: string, activeTab: string): string {
        const querys: string[] = [];

        if (currentExplorerId) {
            querys.push('id=' + currentExplorerId);
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

    function loadAllInsightsExplorerBasicInfos({ force }: { force?: boolean }): Promise<InsightsExplorerBasicInfo[]> {
        if (!force && !insightsExplorerListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allInsightsExplorerBasicInfos.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllInsightsExplorers().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve explorer list' });
                    return;
                }

                if (insightsExplorerListStateInvalid.value) {
                    updateInsightsExplorerListInvalidState(false);
                }

                const explorerBasicInfos = InsightsExplorerBasicInfo.ofMulti(data.result);

                if (force && data.result && isEquals(allInsightsExplorerBasicInfos.value, explorerBasicInfos)) {
                    reject({ message: 'Explorer list is up to date', isUpToDate: true });
                    return;
                }

                loadInsightsExplorerList(explorerBasicInfos);

                resolve(explorerBasicInfos);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load explorer list', error);
                } else {
                    logger.error('failed to load explorer list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve explorer list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getInsightsExplorer({ explorerId }: { explorerId: string }): Promise<InsightsExplorer> {
        return new Promise((resolve, reject) => {
            services.getInsightsExplorer({
                id: explorerId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve explorer' });
                    return;
                }

                const transactionCategory = InsightsExplorer.of(data.result);

                resolve(transactionCategory);
            }).catch(error => {
                logger.error('failed to load explorer info', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve explorer' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveInsightsExplorer({ explorer, saveAs, clientSessionId }: { explorer: InsightsExplorer, saveAs?: boolean, clientSessionId: string }): Promise<InsightsExplorer> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<InsightsExplorerInfoResponse>;

            if (!explorer.id || saveAs) {
                promise = services.addInsightsExplorer(explorer.toCreateRequest(clientSessionId));
            } else {
                promise = services.modifyInsightsExplorer(explorer.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!explorer.id) {
                        reject({ message: 'Unable to add explorer' });
                    } else {
                        reject({ message: 'Unable to save explorer' });
                    }
                    return;
                }

                const explorerBasicInfo = InsightsExplorerBasicInfo.of(data.result);

                if (!explorer.id || saveAs) {
                    addExplorerToInsightsExplorerList(explorerBasicInfo);
                } else {
                    updateExplorerInInsightsExplorerList(explorerBasicInfo);
                }

                resolve(InsightsExplorer.of(data.result));
            }).catch(error => {
                logger.error('failed to save explorer', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!explorer.id) {
                        reject({ message: 'Unable to add explorer' });
                    } else {
                        reject({ message: 'Unable to save explorer' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeInsightsExplorerDisplayOrder({ explorerId, from, to }: { explorerId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentExplorer: InsightsExplorerBasicInfo | null = null;

            for (const insightsExplorer of allInsightsExplorerBasicInfos.value) {
                if (insightsExplorer.id === explorerId) {
                    currentExplorer = insightsExplorer;
                    break;
                }
            }

            if (!currentExplorer || !allInsightsExplorerBasicInfos.value[to]) {
                reject({ message: 'Unable to move explorer' });
                return;
            }

            if (!insightsExplorerListStateInvalid.value) {
                updateInsightsExplorerListInvalidState(true);
            }

            updateExplorerDisplayOrderInInsightsExplorerList({ from, to });

            resolve();
        });
    }

    function updateInsightsExplorerDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: InsightsExplorerNewDisplayOrderRequest[] = [];

        for (const [insightsExplorer, index] of itemAndIndex(allInsightsExplorerBasicInfos.value)) {
            newDisplayOrders.push({
                id: insightsExplorer.id,
                displayOrder: index + 1
            });
        }

        return new Promise((resolve, reject) => {
            services.moveInsightsExplorer({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move explorer' });
                    return;
                }

                loadAllInsightsExplorerBasicInfos({ force: false }).finally(() => {
                    if (insightsExplorerListStateInvalid.value) {
                        updateInsightsExplorerListInvalidState(false);
                    }

                    resolve(data.result);
                });
            }).catch(error => {
                logger.error('failed to save explorers display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move explorer' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function hideInsightsExplorer({ explorer, hidden }: { explorer: InsightsExplorer | InsightsExplorerBasicInfo, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideInsightsExplorer({
                id: explorer.id,
                hidden: hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this explorer' });
                    } else {
                        reject({ message: 'Unable to unhide this explorer' });
                    }
                    return;
                }

                explorer.hidden = hidden;
                updateExplorerVisibilityInInsightsExplorerList({ explorerId: explorer.id, hidden });

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to change explorer visibility', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this explorer' });
                    } else {
                        reject({ message: 'Unable to unhide this explorer' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteInsightsExplorer({ explorer, beforeResolve }: { explorer: InsightsExplorer, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteInsightsExplorer({
                id: explorer.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this explorer' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeExplorerFromInsightsExplorerList(explorer);
                    });
                } else {
                    removeExplorerFromInsightsExplorerList(explorer);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete explorer', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this explorer' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        transactionExplorerFilter,
        transactionExplorerStateInvalid,
        allInsightsExplorerBasicInfos,
        allInsightsExplorerBasicInfosMap,
        currentInsightsExplorer,
        insightsExplorerListStateInvalid,
        // computed
        isUsingAmountRange,
        filteredTransactionsInDataTable,
        filteredTransactionsInDataTableStatistic,
        categoriedTransactionExplorerData,
        // functions
        updateTransactionExplorerInvalidState,
        updateInsightsExplorerListInvalidState,
        updateCurrentInsightsExplorer,
        resetTransactionExplorers,
        initTransactionExplorerFilter,
        updateTransactionExplorerFilter,
        getTransactionExplorerPageParams,
        getTransactionListPageParams,
        loadAllTransactions,
        loadAllInsightsExplorerBasicInfos,
        getInsightsExplorer,
        saveInsightsExplorer,
        changeInsightsExplorerDisplayOrder,
        updateInsightsExplorerDisplayOrders,
        hideInsightsExplorer,
        deleteInsightsExplorer
    };
});
