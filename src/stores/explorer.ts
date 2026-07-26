import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useTransactionTagsStore } from './transactionTag.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { type BeforeResolveFunction, itemAndIndex, keys, values } from '@/core/base.ts';
import { NormalizedText } from '@/core/text.ts';
import { type BigDecimal, NumeralSystem } from '@/core/numeral.ts';
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
    BIG_DECIMAL_ZERO,
    BIG_DECIMAL_POSITIVE_INFINITY,
    BIG_DECIMAL_NEGATIVE_INFINITY,
    parseBigDecimal
} from '@/lib/numeral.ts';
import {
    min,
    max,
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
    kurtosis,
    giniCoefficient,
    herfindahlHirschmanIndex
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
    value: BigDecimal;
}

export interface AmountRanges {
    categorySourceAmountRanges?: BigDecimal[];
    categoryDestinationAmountRanges?: BigDecimal[];
    seriesSourceAmountRanges?: BigDecimal[];
    seriesDestinationAmountRanges?: BigDecimal[];
}

export interface TransactionInsightDataItemInQuery {
    queryIndex: number;
    queryName: string;
    transaction: TransactionInsightDataItem;
}

export interface InsightsExplorerTransactionStatisticData {
    totalCount: number;
    totalAmount: BigDecimal;
    totalIncome: BigDecimal;
    totalExpense: BigDecimal;
    netIncome: BigDecimal;
    averageAmount: BigDecimal;
    medianAmount: BigDecimal;
    minimumAmount: BigDecimal;
    maximumAmount: BigDecimal;
    p90Amount: BigDecimal;
    range: BigDecimal;
    interquartileRange: BigDecimal;
    medianToMeanRatio?: BigDecimal;
    top5AmountShare?: BigDecimal;
    transactionsFor80PercentAmount?: BigDecimal;
    variance?: BigDecimal;
    standardDeviation?: BigDecimal;
    coefficientOfVariation?: BigDecimal;
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

    function buildInsightsExplorerMatchContext(exploration: InsightsExplorer, transaction: TransactionInsightDataItem): InsightsExplorerMatchContext {
        let cachedTransactionDateTime: DateTime | undefined = undefined;
        let cachedNormalizedDescription: NormalizedText | undefined = undefined;

        return {
            getTransactionDateTime(): DateTime {
                if (!cachedTransactionDateTime) {
                    let transactionTimeUtfOffset: number | undefined = undefined;

                    if (exploration.timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
                        transactionTimeUtfOffset = transaction.utcOffset;
                    }

                    cachedTransactionDateTime = isDefined(transactionTimeUtfOffset) ? parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transactionTimeUtfOffset) : parseDateTimeFromUnixTime(transaction.time);
                }

                return cachedTransactionDateTime;
            },
            getNormalizedDescription(): NormalizedText {
                if (!cachedNormalizedDescription) {
                    cachedNormalizedDescription = NormalizedText.of(transaction.comment ? transaction.comment : '');
                }

                return cachedNormalizedDescription;
            }
        };
    }

    function calculateAmountRanges(sortedAmounts: BigDecimal[], dimension: TransactionExplorerDataDimension, rangeCount: number): BigDecimal[] {
        const result: BigDecimal[] = [];

        if (sortedAmounts.length < 1 || rangeCount <= 0) {
            return result;
        }

        const minAmount = sortedAmounts[0] as BigDecimal;
        const maxAmount = sortedAmounts[sortedAmounts.length - 1] as BigDecimal;
        rangeCount = Math.min(rangeCount, sortedAmounts.length);

        // [min1, max1), [min2, max2), ..., [minN, maxN]
        if (dimension === TransactionExplorerDataDimension.SourceAmountRangeEqualFrequency
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeEqualFrequency) {
            for (let i = 0; i < rangeCount; i++) {
                result.push(sortedAmounts[Math.floor(i * (sortedAmounts.length - 1) / rangeCount)] as BigDecimal);
            }
            result.push(maxAmount);
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeEqualWidth
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeEqualWidth) {
            if (minAmount === maxAmount) {
                return [minAmount, maxAmount];
            }

            const width: BigDecimal = maxAmount.subtract(minAmount).divide(rangeCount); // (maxAmount - minAmount) / rangeCount

            for (let i = 0; i < rangeCount; i++) {
                result.push(minAmount.add(width.multiply(i))); // minAmount + i * width
            }
            result.push(maxAmount);
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeLogScale
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeLogScale) {
            const epsilon: number = 1e-9;

            const transform = (x: BigDecimal): BigDecimal => {
                if (x.isZero()) {
                    return x;
                }

                return x.sign().multiply(x.abs().add(epsilon).log()); // sign(x) * log(abs(x) + epsilon)
            };

            const inverse = (y: BigDecimal): BigDecimal => {
                if (y.isZero()) {
                    return y;
                }

                return y.sign().multiply(y.abs().exp().subtract(epsilon)); // sign(y) * (exp(abs(y)) - epsilon)
            };

            const transformed: BigDecimal[] = sortedAmounts.map(transform).sort((a, b) => a.compareTo(b));

            const tMin: BigDecimal = transformed[0] as BigDecimal;
            const tMax: BigDecimal = transformed[transformed.length - 1] as BigDecimal;

            if (tMin === tMax) {
                return [minAmount, maxAmount];
            }

            const width: BigDecimal = tMax.subtract(tMin).divide(rangeCount); // (tMax - tMin) / rangeCount

            result.push(minAmount);
            for (let i = 1; i < rangeCount; i++) {
                result.push(inverse(tMin.add(width.multiply(i)))); // inverse(tMin + i * width)
            }
            result.push(maxAmount);
        } else if (dimension === TransactionExplorerDataDimension.SourceAmountRangeStandardDeviation
            || dimension === TransactionExplorerDataDimension.DestinationAmountRangeStandardDeviation) {
            if (minAmount === maxAmount) {
                return [minAmount, maxAmount];
            }

            const averageAmountForVarianceCalculation: BigDecimal = mean(sortedAmounts, item => item).divide(AMOUNT_FACTOR);
            const { standardDeviation } = varianceAndStandardDeviation(sortedAmounts, averageAmountForVarianceCalculation, item => item.divide(AMOUNT_FACTOR));

            if (standardDeviation.isZero()) {
                return [minAmount, maxAmount];
            }

            const rawBreaks: BigDecimal[] = [];
            const halfCount = Math.floor(rangeCount / 2);

            if (rangeCount % 2 === 1) {
                for (let i = -halfCount; i <= halfCount; i++) {
                    rawBreaks.push(averageAmountForVarianceCalculation.add(standardDeviation.multiply(i)).multiply(AMOUNT_FACTOR));
                }
            } else {
                for (let i = -halfCount; i <= halfCount; i++) {
                    if (i === 0) {
                        continue;
                    }
                    rawBreaks.push(averageAmountForVarianceCalculation.add(standardDeviation.multiply(i - 0.5)).multiply(AMOUNT_FACTOR));
                }
                rawBreaks.sort((a, b) => a.compareTo(b));
            }

            const clipped = rawBreaks.map((v) => max(minAmount, min(maxAmount, v)))
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
            const varianceCombinations: BigDecimal[][] = Array.from({ length: n + 1 }, () => new Array(k + 1).fill(BIG_DECIMAL_POSITIVE_INFINITY));

            for (let i = 1; i <= k; i++) {
                lowerClassLimits[1]![i] = 1;
                varianceCombinations[1]![i] = BIG_DECIMAL_ZERO;
            }

            for (let l = 2; l <= n; l++) {
                let sumZ: BigDecimal = BIG_DECIMAL_ZERO;
                let sumZ2: BigDecimal = BIG_DECIMAL_ZERO;

                for (let m = 1; m <= l; m++) {
                    const val = sortedAmounts[l - m] as BigDecimal;
                    sumZ = sumZ.add(val);
                    sumZ2 = sumZ2.add(val.pow(2));

                    const variance = sumZ2.subtract(sumZ.pow(2).divide(m));


                    if (m === l) {
                        for (let j = 1; j <= k; j++) {
                            if (variance.lessThan(varianceCombinations[l]![j]!)) {
                                lowerClassLimits[l]![j] = 1;
                                varianceCombinations[l]![j] = variance;
                            }
                        }
                    } else {
                        for (let j = 2; j <= k; j++) {
                            const combined = varianceCombinations[l - m]![j - 1]!.add(variance);
                            if (combined.lessThan(varianceCombinations[l]![j]!)) {
                                lowerClassLimits[l]![j] = l - m + 1;
                                varianceCombinations[l]![j] = combined;
                            }
                        }
                    }
                }
            }

            const breaks: BigDecimal[] = new Array(k + 1);
            breaks[k] = maxAmount;

            let currentK = k;
            let currentIdx = n;

            while (currentK >= 2) {
                const lowerIdx = lowerClassLimits[currentIdx]![currentK]!;
                breaks[currentK - 1] = sortedAmounts[lowerIdx - 1] as BigDecimal;
                currentIdx = lowerIdx - 1;
                currentK--;
            }

            breaks[0] = minAmount;
            return breaks;
        }

        return result;
    }

    function getDataCategoryInfo(timezoneUsedForDateRange: number, dimension: TransactionExplorerDataDimension, sourceAmountRanges: BigDecimal[] | undefined, destinationAmountRanges: BigDecimal[] | undefined, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): CategoriedInfo {
        const defaultCurrency = userStore.currentUserDefaultCurrency;
        let transactionTimeUtfOffset: number | undefined = undefined;

        if (timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
            transactionTimeUtfOffset = transaction.utcOffset;
        }

        if (dimension === TransactionExplorerDataDimension.None) {
            const valueMetric = TransactionExplorerValueMetric.valueOf(currentExploration.value.valueMetric);
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
                categoryDisplayOrders: [transaction.primaryCategory.type, transaction.primaryCategory.displayOrder]
            };
        } else if (dimension === TransactionExplorerDataDimension.SecondaryCategory) {
            return {
                categoryName: transaction.secondaryCategory.name,
                categoryId: transaction.categoryId,
                categoryIdType: TransactionExplorerDimensionType.Category,
                categoryDisplayOrders: [transaction.primaryCategory.type, transaction.primaryCategory.displayOrder, transaction.secondaryCategory.displayOrder]
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

            const amount: BigDecimal = dimension === TransactionExplorerDataDimension.SourceAmount ? parseBigDecimal(transaction.sourceAmount) : parseBigDecimal(transaction.destinationAmount);
            const account = dimension === TransactionExplorerDataDimension.SourceAmount ? transaction.sourceAccount : transaction.destinationAccount;
            let amountInDefaultCurrency: BigDecimal = amount;

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

                if (exchangedAmount) {
                    amountInDefaultCurrency = exchangedAmount.truncate();
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
                categoryName: amountInDefaultCurrency.toString(),
                categoryId: amountInDefaultCurrency.toString(),
                categoryIdType: TransactionExplorerDimensionType.Amount,
                categoryDisplayOrders: [amountInDefaultCurrency.toDoubleNumber()]
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

            const amount: BigDecimal = dimension.isSourceAmountRange ? parseBigDecimal(transaction.sourceAmount) : parseBigDecimal(transaction.destinationAmount);
            const account = dimension.isSourceAmountRange ? transaction.sourceAccount : transaction.destinationAccount;
            let amountInDefaultCurrency: BigDecimal = amount;

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

                if (exchangedAmount) {
                    amountInDefaultCurrency = exchangedAmount.truncate();
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

            const amountRanges: BigDecimal[] = isSourceAmount ? (sourceAmountRanges ?? []) : (destinationAmountRanges ?? []);
            let matchAmountRangeMin: BigDecimal | undefined = undefined;
            let matchAmountRangeMax: BigDecimal | undefined = undefined;
            let matchAmountRangeIndex: number | undefined = undefined;

            for (let i = 1; i < amountRanges.length; i++) {
                const amountRangeMin = amountRanges[i - 1] as BigDecimal;
                const amountRangeMax = amountRanges[i] as BigDecimal;

                if (amountInDefaultCurrency.lessThan(amountRangeMin)) {
                    continue;
                }

                if (amountInDefaultCurrency.greaterThan(amountRangeMax)) {
                    continue;
                }

                if (i < amountRanges.length - 1 && amountInDefaultCurrency.equals(amountRangeMax)) {
                    continue;
                }

                matchAmountRangeMin = amountRangeMin;
                matchAmountRangeMax = amountRangeMax;
                matchAmountRangeIndex = i - 1;
            }

            if (matchAmountRangeMin && matchAmountRangeMax && isNumber(matchAmountRangeIndex)) {
                return {
                    categoryName: `${matchAmountRangeMin.toString()}|${matchAmountRangeMax.toString()}`,
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

    function addTransactionToFilteredList(filteredTransactions: TransactionInsightDataItemInQuery[], filteredTransactionSourceAmountsInDefaultCurrency: BigDecimal[], filteredTransactionDestinationAmountsInDefaultCurrency: BigDecimal[], defaultCurrency: string, queryName: string, queryIndex: number, transaction: TransactionInsightDataItem): void {
        filteredTransactions.push({
            queryIndex: queryIndex,
            queryName: queryName,
            transaction: transaction
        });

        let sourceAmountInDefaultCurrency: BigDecimal | undefined = parseBigDecimal(transaction.sourceAmount);
        let destinationAmountInDefaultCurrency: BigDecimal | undefined = transaction.type === TransactionType.Transfer && transaction.destinationAccount ? parseBigDecimal(transaction.destinationAmount) : undefined;

        if (transaction.sourceAccount.currency !== defaultCurrency) {
            const amount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(transaction.sourceAmount), transaction.sourceAccount.currency, defaultCurrency);
            sourceAmountInDefaultCurrency = amount?.truncate();
        }

        if (transaction.type === TransactionType.Transfer && transaction.destinationAccount && transaction.destinationAccount.currency !== defaultCurrency) {
            const amount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(transaction.destinationAmount), transaction.destinationAccount.currency, defaultCurrency);
            destinationAmountInDefaultCurrency = amount?.truncate();
        }

        if (sourceAmountInDefaultCurrency) {
            filteredTransactionSourceAmountsInDefaultCurrency.push(sourceAmountInDefaultCurrency);
        }

        if (destinationAmountInDefaultCurrency) {
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

    function buildAllAmountRanges(categoryDimension: TransactionExplorerDataDimension, seriesDimension: TransactionExplorerDataDimension, filteredTransactionSourceAmountsInDefaultCurrency: BigDecimal[], filteredTransactionDestinationAmountsInDefaultCurrency: BigDecimal[], rangeCount: number): AmountRanges {
        const allAmountRanges: AmountRanges = {};

        if (categoryDimension.isSourceAmountRange || seriesDimension.isSourceAmountRange) {
            filteredTransactionSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
            const sorteUniqueAmounts = filteredTransactionSourceAmountsInDefaultCurrency.filter((v, i, a) => i === 0 || v.notEquals(a[i - 1]));

            if (categoryDimension.isSourceAmountRange) {
                allAmountRanges.categorySourceAmountRanges = calculateAmountRanges(sorteUniqueAmounts, categoryDimension, rangeCount);
            }

            if (seriesDimension.isSourceAmountRange) {
                allAmountRanges.seriesSourceAmountRanges = calculateAmountRanges(sorteUniqueAmounts, seriesDimension, rangeCount);
            }
        }

        if (categoryDimension.isDestinationAmountRange || seriesDimension.isDestinationAmountRange) {
            filteredTransactionDestinationAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
            const sorteUniqueAmounts = filteredTransactionDestinationAmountsInDefaultCurrency.filter((v, i, a) => i === 0 || v.notEquals(a[i - 1]));
            if (categoryDimension.isDestinationAmountRange) {
                allAmountRanges.categoryDestinationAmountRanges = calculateAmountRanges(sorteUniqueAmounts, categoryDimension, rangeCount);
            }

            if (seriesDimension.isDestinationAmountRange) {
                allAmountRanges.seriesDestinationAmountRanges = calculateAmountRanges(sorteUniqueAmounts, seriesDimension, rangeCount);
            }
        }

        return allAmountRanges;
    }

    function loadInsightsExplorerList(explorations: InsightsExplorerBasicInfo[]): void {
        allExplorationBasicInfos.value = explorations;
        allExplorationBasicInfosMap.value = {};

        for (const exploration of explorations) {
            allExplorationBasicInfosMap.value[exploration.id] = exploration;
        }
    }

    function addExplorationToInsightsExplorerList(exploration: InsightsExplorerBasicInfo): void {
        allExplorationBasicInfos.value.push(exploration);
        allExplorationBasicInfosMap.value[exploration.id] = exploration;
    }

    function updateExplorationInInsightsExplorerList(currentExploration: InsightsExplorerBasicInfo): void {
        for (const [explorer, index] of itemAndIndex(allExplorationBasicInfos.value)) {
            if (explorer.id === currentExploration.id) {
                allExplorationBasicInfos.value.splice(index, 1, currentExploration);
                break;
            }
        }

        allExplorationBasicInfosMap.value[currentExploration.id] = currentExploration;
    }

    function updateExplorerDisplayOrderInInsightsExplorerList({ from, to }: { from: number, to: number }): void {
        allExplorationBasicInfos.value.splice(to, 0, allExplorationBasicInfos.value.splice(from, 1)[0] as InsightsExplorer);
    }

    function updateExplorationVisibilityInInsightsExplorerList({ explorationId, hidden }: { explorationId: string, hidden: boolean }): void {
        if (allExplorationBasicInfosMap.value[explorationId]) {
            allExplorationBasicInfosMap.value[explorationId]!.hidden = hidden;
        }
    }

    function removeExplorationFromInsightsExplorerList(currentExploration: InsightsExplorerBasicInfo): void {
        for (const [insightsExplorer, index] of itemAndIndex(allExplorationBasicInfos.value)) {
            if (insightsExplorer.id === currentExploration.id) {
                allExplorationBasicInfos.value.splice(index, 1);
                break;
            }
        }

        if (allExplorationBasicInfosMap.value[currentExploration.id]) {
            delete allExplorationBasicInfosMap.value[currentExploration.id];
        }
    }

    const transactionExplorerFilter = ref<TransactionExplorerFilter>({
        dateRangeType: DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type,
        startTime: 0,
        endTime: 0
    });

    const transactionExplorerAllData = ref<TransactionInfoResponse[]>([]);
    const transactionExplorerStateInvalid = ref<boolean>(true);

    const allExplorationBasicInfos = ref<InsightsExplorerBasicInfo[]>([]);
    const allExplorationBasicInfosMap = ref<Record<string, InsightsExplorerBasicInfo>>({});
    const currentExploration = ref<InsightsExplorer>(InsightsExplorer.createNewExplorer(generateRandomUUID()));
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
        const chartType = TransactionExplorerChartType.valueOf(currentExploration.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(currentExploration.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(currentExploration.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;
        return categoryDimension?.isSourceAmountRange || seriesDimension?.isSourceAmountRange
            || categoryDimension?.isDestinationAmountRange || seriesDimension?.isDestinationAmountRange
            || false;
    });

    const filteredTransactionsInDataTable = computed<TransactionInsightDataItem[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        if (!currentExploration.value.queries || currentExploration.value.queries.length < 1) {
            return allTransactions.value;
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of allTransactions.value) {
            const matchOptions: InsightsExplorerMatchContext = buildInsightsExplorerMatchContext(currentExploration.value, transaction);

            for (const query of currentExploration.value.queries) {
                if (currentExploration.value.datatableQuerySource && currentExploration.value.datatableQuerySource !== query.id) {
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
            totalAmount: BIG_DECIMAL_ZERO,
            totalIncome: BIG_DECIMAL_ZERO,
            totalExpense: BIG_DECIMAL_ZERO,
            netIncome: BIG_DECIMAL_ZERO,
            averageAmount: BIG_DECIMAL_ZERO,
            medianAmount: BIG_DECIMAL_ZERO,
            minimumAmount: BIG_DECIMAL_POSITIVE_INFINITY,
            maximumAmount: BIG_DECIMAL_NEGATIVE_INFINITY,
            p90Amount: BIG_DECIMAL_ZERO,
            range: BIG_DECIMAL_ZERO,
            interquartileRange: BIG_DECIMAL_ZERO,
            medianToMeanRatio: undefined,
            top5AmountShare: undefined,
            transactionsFor80PercentAmount: undefined,
            variance: undefined,
            standardDeviation: undefined,
            coefficientOfVariation: undefined
        };

        const sourceAmounts: BigDecimal[] = [];

        for (const transaction of filteredTransactionsInDataTable.value) {
            statisticData.totalCount++;

            let amountInDefaultCurrency: BigDecimal = parseBigDecimal(transaction.sourceAmount);

            if (transaction.sourceAccount.currency !== defaultCurrency) {
                const amount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(transaction.sourceAmount), transaction.sourceAccount.currency, defaultCurrency);

                if (amount) {
                    amountInDefaultCurrency = amount.truncate();
                } else {
                    continue;
                }
            }

            statisticData.totalAmount = statisticData.totalAmount.add(amountInDefaultCurrency);
            sourceAmounts.push(amountInDefaultCurrency);

            if (transaction.type === TransactionType.Income) {
                statisticData.totalIncome = statisticData.totalIncome.add(amountInDefaultCurrency);
            } else if (transaction.type === TransactionType.Expense) {
                statisticData.totalExpense = statisticData.totalExpense.add(amountInDefaultCurrency);
            }

            if (amountInDefaultCurrency.isPositiveOrZero() && amountInDefaultCurrency.lessThan(statisticData.minimumAmount)) {
                statisticData.minimumAmount = amountInDefaultCurrency;
            }

            if (amountInDefaultCurrency.greaterThan(statisticData.maximumAmount)) {
                statisticData.maximumAmount = amountInDefaultCurrency;
            }
        }

        statisticData.netIncome = statisticData.totalIncome.subtract(statisticData.totalExpense);

        if (statisticData.minimumAmount.isPositiveInfinity()) {
            statisticData.minimumAmount = BIG_DECIMAL_ZERO;
        }

        if (statisticData.maximumAmount.isNegativeInfinity()) {
            statisticData.maximumAmount = BIG_DECIMAL_ZERO;
        }

        if (sourceAmounts.length > 0) {
            statisticData.averageAmount = statisticData.totalAmount.divide(sourceAmounts.length).truncate();
        }

        statisticData.range = statisticData.maximumAmount.subtract(statisticData.minimumAmount);

        if (sourceAmounts.length > 0) {
            sourceAmounts.sort((a, b) => a.compareTo(b));
            statisticData.medianAmount = median(sourceAmounts, item => item).truncate();
            statisticData.p90Amount = percentile(sourceAmounts, 0.9, item => item).truncate();

            const q1 = percentile(sourceAmounts, 0.25, item => item);
            const q3 = percentile(sourceAmounts, 0.75, item => item);
            statisticData.interquartileRange = q3.subtract(q1).truncate();
            statisticData.medianToMeanRatio = !statisticData.averageAmount.isZero() ? statisticData.medianAmount.divide(statisticData.averageAmount) : undefined;
        }

        if (sourceAmounts.length > 5) {
            const top5AmountSum = sumMaxN(sourceAmounts, 5, item => item);
            statisticData.top5AmountShare = statisticData.totalAmount.isPositive() ? top5AmountSum.divide(statisticData.totalAmount).multiply(100) : BIG_DECIMAL_ZERO;
        }

        if (sourceAmounts.length > 0) {
            statisticData.transactionsFor80PercentAmount = cumulativePercentage(sourceAmounts, 0.8, statisticData.totalAmount, item => item);
        }

        if (sourceAmounts.length > 0) {
            const averageAmountForVarianceCalculation: BigDecimal = statisticData.totalAmount.divide(sourceAmounts.length).divide(AMOUNT_FACTOR);
            const { variance, standardDeviation } = varianceAndStandardDeviation(sourceAmounts, averageAmountForVarianceCalculation, item => item.divide(AMOUNT_FACTOR));
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

        const chartType = TransactionExplorerChartType.valueOf(currentExploration.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(currentExploration.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(currentExploration.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;

        if (!chartType || !categoryDimension || !seriesDimension) {
            return {};
        }

        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const filteredTransactions: TransactionInsightDataItemInQuery[] = [];
        const filteredTransactionSourceAmountsInDefaultCurrency: BigDecimal[] = [];
        const filteredTransactionDestinationAmountsInDefaultCurrency: BigDecimal[] = [];

        for (const transaction of allTransactions.value) {
            if (!currentExploration.value.queries || currentExploration.value.queries.length < 1) {
                addTransactionToFilteredList(filteredTransactions, filteredTransactionSourceAmountsInDefaultCurrency, filteredTransactionDestinationAmountsInDefaultCurrency, defaultCurrency, '', 0, transaction);
                continue;
            }

            const matchContext: InsightsExplorerMatchContext = buildInsightsExplorerMatchContext(currentExploration.value, transaction);

            for (const [query, index] of itemAndIndex(currentExploration.value.queries)) {
                if (query.match(transaction, matchContext)) {
                    addTransactionToFilteredList(filteredTransactions, filteredTransactionSourceAmountsInDefaultCurrency, filteredTransactionDestinationAmountsInDefaultCurrency, defaultCurrency, query.name, index, transaction);

                    if (categoryDimension !== TransactionExplorerDataDimension.Query) {
                        break;
                    }
                }
            }
        }

        const categoriedDataMap: Record<string, CategoriedTransactions> = {};
        const allAmountRanges: AmountRanges = buildAllAmountRanges(categoryDimension, seriesDimension, filteredTransactionSourceAmountsInDefaultCurrency, filteredTransactionDestinationAmountsInDefaultCurrency, currentExploration.value.amountRangeCount);

        for (const item of filteredTransactions) {
            addTransactionToCategoriedDataMap(currentExploration.value.timezoneUsedForDateRange, categoriedDataMap, categoryDimension, seriesDimension, allAmountRanges, item.queryName, item.queryIndex, item.transaction);
        }

        return categoriedDataMap;
    });

    const categoriedTransactionExplorerData = computed<CategoriedTransactionExplorerData[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        const chartType = TransactionExplorerChartType.valueOf(currentExploration.value.chartType);
        const categoryDimension = TransactionExplorerDataDimension.valueOf(currentExploration.value.categoryDimension);
        const seriesDimension = chartType?.seriesDimensionRequired ? TransactionExplorerDataDimension.valueOf(currentExploration.value.seriesDimension) : TransactionExplorerDataDimension.SeriesDimensionDefault;
        const valueMetric = TransactionExplorerValueMetric.valueOf(currentExploration.value.valueMetric);

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
                const allSourceAmountsInDefaultCurrency: BigDecimal[] = [];
                let totalSourceAmountSumInDefaultCurrency: BigDecimal = BIG_DECIMAL_ZERO;
                let totalSourceIncomeAmountSumInDefaultCurrency: BigDecimal = BIG_DECIMAL_ZERO;
                let totalSourceExpenseAmountSumInDefaultCurrency: BigDecimal = BIG_DECIMAL_ZERO;
                let minimumSourceAmountInDefaultCurrency: BigDecimal = BIG_DECIMAL_POSITIVE_INFINITY;
                let maximumSourceAmountInDefaultCurrency: BigDecimal = BIG_DECIMAL_NEGATIVE_INFINITY;

                for (const transaction of seriesTransactions.trasactions) {
                    let amountInDefaultCurrency: BigDecimal = parseBigDecimal(transaction.sourceAmount);

                    if (transaction.sourceAccount.currency !== defaultCurrency) {
                        const amount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(transaction.sourceAmount), transaction.sourceAccount.currency, defaultCurrency);

                        if (amount) {
                            amountInDefaultCurrency = amount.truncate();
                        } else {
                            continue;
                        }
                    }

                    if (needCalculateDailyTransactionCount) {
                        let transactionTimeUtfOffset: number | undefined = undefined;

                        if (currentExploration.value.timezoneUsedForDateRange === TimezoneTypeForStatistics.TransactionTimezone.type) {
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
                    totalSourceAmountSumInDefaultCurrency = totalSourceAmountSumInDefaultCurrency.add(amountInDefaultCurrency);

                    if (transaction.type === TransactionType.Income) {
                        totalSourceIncomeAmountSumInDefaultCurrency = totalSourceIncomeAmountSumInDefaultCurrency.add(amountInDefaultCurrency);
                    } else if (transaction.type === TransactionType.Expense) {
                        totalSourceExpenseAmountSumInDefaultCurrency = totalSourceExpenseAmountSumInDefaultCurrency.add(amountInDefaultCurrency);
                    }

                    if (amountInDefaultCurrency.isPositiveOrZero() && amountInDefaultCurrency.lessThan(minimumSourceAmountInDefaultCurrency)) {
                        minimumSourceAmountInDefaultCurrency = amountInDefaultCurrency;
                    }

                    if (amountInDefaultCurrency.greaterThan(maximumSourceAmountInDefaultCurrency)) {
                        maximumSourceAmountInDefaultCurrency = amountInDefaultCurrency;
                    }
                }

                let value: BigDecimal = BIG_DECIMAL_ZERO;

                if (valueMetric === TransactionExplorerValueMetric.TransactionCount) {
                    value = parseBigDecimal(allSourceAmountsInDefaultCurrency.length);
                } else if (valueMetric === TransactionExplorerValueMetric.ActiveTransactionDays) {
                    value = parseBigDecimal(getObjectOwnFieldCount(transactionDateMapCount));
                } else if (valueMetric === TransactionExplorerValueMetric.TransactionsPerDay) {
                    const activeDays = getObjectOwnFieldCount(transactionDateMapCount);
                    value = activeDays > 0 ? parseBigDecimal(allSourceAmountsInDefaultCurrency.length).divide(activeDays) : BIG_DECIMAL_ZERO;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceIncomeAmountSum) {
                    value = totalSourceIncomeAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceExpenseAmountSum) {
                    value = totalSourceExpenseAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceNetIncomeAmountSum) {
                    value = totalSourceIncomeAmountSumInDefaultCurrency.subtract(totalSourceExpenseAmountSumInDefaultCurrency);
                } else if (valueMetric === TransactionExplorerValueMetric.SrouceAmountExpenseIncomeRatio) {
                    value = !totalSourceIncomeAmountSumInDefaultCurrency.isZero() ? totalSourceExpenseAmountSumInDefaultCurrency.divide(totalSourceIncomeAmountSumInDefaultCurrency).multiply(100) : BIG_DECIMAL_ZERO;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSavingsRate) {
                    value = !totalSourceIncomeAmountSumInDefaultCurrency.isZero() ? totalSourceIncomeAmountSumInDefaultCurrency.subtract(totalSourceExpenseAmountSumInDefaultCurrency).divide(totalSourceIncomeAmountSumInDefaultCurrency).multiply(100) : BIG_DECIMAL_ZERO;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSum) {
                    value = totalSourceAmountSumInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountAverage) {
                    value = allSourceAmountsInDefaultCurrency.length > 0 ? totalSourceAmountSumInDefaultCurrency.divide(allSourceAmountsInDefaultCurrency.length).truncate() : BIG_DECIMAL_ZERO;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMedian) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        value = median(allSourceAmountsInDefaultCurrency, item => item).truncate();
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMinimum) {
                    value = minimumSourceAmountInDefaultCurrency.isPositiveInfinity() ? BIG_DECIMAL_ZERO : minimumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMaximum) {
                    value = maximumSourceAmountInDefaultCurrency.isNegativeInfinity() ? BIG_DECIMAL_ZERO : maximumSourceAmountInDefaultCurrency;
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountQ1Amount
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountQ3Amount
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount10thPercentile
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount90thPercentile
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount95thPercentile
                    || valueMetric === TransactionExplorerValueMetric.SourceAmount99thPercentile) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));

                        if (valueMetric === TransactionExplorerValueMetric.SourceAmountQ1Amount) {
                            value = percentile(allSourceAmountsInDefaultCurrency, 0.25, item => item).truncate();
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountQ3Amount) {
                            value = percentile(allSourceAmountsInDefaultCurrency, 0.75, item => item).truncate();
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount10thPercentile) {
                            value = percentile(allSourceAmountsInDefaultCurrency, 0.1, item => item).truncate();
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount90thPercentile) {
                            value = percentile(allSourceAmountsInDefaultCurrency, 0.9, item => item).truncate();
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount95thPercentile) {
                            value = percentile(allSourceAmountsInDefaultCurrency, 0.95, item => item).truncate();
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmount99thPercentile) {
                            value = percentile(allSourceAmountsInDefaultCurrency, 0.99, item => item).truncate();
                        }
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountRange) {
                    const finalMinimumSourceAmountInDefaultCurrency = minimumSourceAmountInDefaultCurrency.isPositiveInfinity() ? BIG_DECIMAL_ZERO : minimumSourceAmountInDefaultCurrency;
                    const finalMaximumSourceAmountInDefaultCurrency = maximumSourceAmountInDefaultCurrency.isNegativeInfinity() ? BIG_DECIMAL_ZERO : maximumSourceAmountInDefaultCurrency;
                    value = finalMaximumSourceAmountInDefaultCurrency.subtract(finalMinimumSourceAmountInDefaultCurrency);
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountInterquartileRange) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        const q1 = percentile(allSourceAmountsInDefaultCurrency, 0.25, item => item);
                        const q3 = percentile(allSourceAmountsInDefaultCurrency, 0.75, item => item);
                        value = q3.subtract(q1).truncate();
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMeanAbsoluteDeviation) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        const averageSourceAmountInDefaultCurrency = totalSourceAmountSumInDefaultCurrency.divide(allSourceAmountsInDefaultCurrency.length);
                        value = meanAbsoluteDeviation(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, item => item).truncate();
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMedianAbsoluteDeviation) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        const medianSourceAmountInDefaultCurrency = median(allSourceAmountsInDefaultCurrency, item => item);
                        value = medianAbsoluteDeviation(allSourceAmountsInDefaultCurrency, medianSourceAmountInDefaultCurrency, item => item).truncate();
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountMedianToMeanRatio) {
                    if (allSourceAmountsInDefaultCurrency.length > 0 && !totalSourceAmountSumInDefaultCurrency.isZero()) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        const medianSourceAmountInDefaultCurrency = median(allSourceAmountsInDefaultCurrency, item => item);
                        const averageSourceAmountInDefaultCurrency = totalSourceAmountSumInDefaultCurrency.divide(allSourceAmountsInDefaultCurrency.length);
                        value = !averageSourceAmountInDefaultCurrency.isZero() ? medianSourceAmountInDefaultCurrency.divide(averageSourceAmountInDefaultCurrency) : BIG_DECIMAL_ZERO;
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceMaximumAmountShare) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        value = !maximumSourceAmountInDefaultCurrency.isNegativeInfinity() ? maximumSourceAmountInDefaultCurrency.divide(totalSourceAmountSumInDefaultCurrency).multiply(100) : BIG_DECIMAL_ZERO;
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceTop5AmountSum) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        value = sumMaxN(allSourceAmountsInDefaultCurrency, 5, item => item);
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceTop5AmountShare) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        const top5AmountSum = sumMaxN(allSourceAmountsInDefaultCurrency, 5, item => item);
                        value = totalSourceAmountSumInDefaultCurrency.isPositive() ? top5AmountSum.divide(totalSourceAmountSumInDefaultCurrency).multiply(100) : BIG_DECIMAL_ZERO;
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.TransactionsForEightyPercentOfSourceAmount) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        value = cumulativePercentage(allSourceAmountsInDefaultCurrency, 0.8, totalSourceAmountSumInDefaultCurrency, item => item);
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountOutlierCount
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountOutlierRatio) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        const q1 = percentile(allSourceAmountsInDefaultCurrency, 0.25, item => item);
                        const q3 = percentile(allSourceAmountsInDefaultCurrency, 0.75, item => item);
                        const iqr = q3.subtract(q1);
                        const lowerBound = q1.subtract(iqr.multiply(1.5));
                        const upperBound = q3.add(iqr.multiply(1.5));

                        let outlierCount = 0;
                        for (const amount of allSourceAmountsInDefaultCurrency) {
                            if (amount.lessThan(lowerBound) || amount.greaterThan(upperBound)) {
                                outlierCount++;
                            }
                        }

                        if (valueMetric === TransactionExplorerValueMetric.SourceAmountOutlierCount) {
                            value = parseBigDecimal(outlierCount);
                        } else {
                            value = allSourceAmountsInDefaultCurrency.length > 0 ? parseBigDecimal(outlierCount).divide(allSourceAmountsInDefaultCurrency.length).multiply(100) : BIG_DECIMAL_ZERO;
                        }
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountVariance
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountStandardDeviation
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountCoefficientOfVariation
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountSkewness
                    || valueMetric === TransactionExplorerValueMetric.SourceAmountKurtosis) {
                    if (allSourceAmountsInDefaultCurrency.length > 0) {
                        const averageSourceAmountInDefaultCurrency = totalSourceAmountSumInDefaultCurrency.divide(allSourceAmountsInDefaultCurrency.length).divide(AMOUNT_FACTOR);
                        const { variance, standardDeviation } = varianceAndStandardDeviation(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, item => item.divide(AMOUNT_FACTOR));

                        if (valueMetric === TransactionExplorerValueMetric.SourceAmountVariance) {
                            value = variance;
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountStandardDeviation) {
                            value = standardDeviation;
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountCoefficientOfVariation) {
                            value = coefficientOfVariation(standardDeviation, averageSourceAmountInDefaultCurrency) ?? BIG_DECIMAL_ZERO;
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountSkewness) {
                            value = skewness(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, standardDeviation, item => item.divide(AMOUNT_FACTOR));
                        } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountKurtosis) {
                            value = kurtosis(allSourceAmountsInDefaultCurrency, averageSourceAmountInDefaultCurrency, variance, item => item.divide(AMOUNT_FACTOR));
                        }
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountGiniCoefficient) {
                    if (allSourceAmountsInDefaultCurrency.length > 0 && !totalSourceAmountSumInDefaultCurrency.isZero()) {
                        allSourceAmountsInDefaultCurrency.sort((a, b) => a.compareTo(b));
                        value = giniCoefficient(allSourceAmountsInDefaultCurrency, totalSourceAmountSumInDefaultCurrency, item => item);
                    } else {
                        value = BIG_DECIMAL_ZERO;
                    }
                } else if (valueMetric === TransactionExplorerValueMetric.SourceAmountHerfindahlHirschmanIndex) {
                    if (allSourceAmountsInDefaultCurrency.length > 0 && !totalSourceAmountSumInDefaultCurrency.isZero()) {
                        value = herfindahlHirschmanIndex(allSourceAmountsInDefaultCurrency, totalSourceAmountSumInDefaultCurrency, item => item);
                    } else {
                        value = BIG_DECIMAL_ZERO;
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

    function updateCurrentExploration(exploration: InsightsExplorer): void {
        currentExploration.value = exploration;
    }

    function resetTransactionExplorers(): void {
        transactionExplorerFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type;
        transactionExplorerFilter.value.startTime = 0;
        transactionExplorerFilter.value.endTime = 0;
        transactionExplorerAllData.value = [];
        allExplorationBasicInfos.value = [];
        allExplorationBasicInfosMap.value = {};
        currentExploration.value = InsightsExplorer.createNewExplorer(generateRandomUUID());
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
            currentExploration.value = InsightsExplorer.createNewExplorer(generateRandomUUID());
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

    function loadAllExplorationBasicInfos({ force }: { force?: boolean }): Promise<InsightsExplorerBasicInfo[]> {
        if (!force && !insightsExplorerListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allExplorationBasicInfos.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllExplorations().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve exploration list' });
                    return;
                }

                if (insightsExplorerListStateInvalid.value) {
                    updateInsightsExplorerListInvalidState(false);
                }

                const explorationBasicInfos = InsightsExplorerBasicInfo.ofMulti(data.result);

                if (force && data.result && isEquals(allExplorationBasicInfos.value, explorationBasicInfos)) {
                    reject({ message: 'Exploration list is up to date', isUpToDate: true });
                    return;
                }

                loadInsightsExplorerList(explorationBasicInfos);

                resolve(explorationBasicInfos);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load exploration list', error);
                } else {
                    logger.error('failed to load exploration list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve exploration list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getExploration({ explorationId }: { explorationId: string }): Promise<InsightsExplorer> {
        return new Promise((resolve, reject) => {
            services.getExploration({
                id: explorationId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve exploration' });
                    return;
                }

                const transactionCategory = InsightsExplorer.of(data.result);

                resolve(transactionCategory);
            }).catch(error => {
                logger.error('failed to load exploration info', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve exploration' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveExploration({ exploration, saveAs, clientSessionId }: { exploration: InsightsExplorer, saveAs?: boolean, clientSessionId: string }): Promise<InsightsExplorer> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<InsightsExplorerInfoResponse>;

            if (!exploration.id || saveAs) {
                promise = services.addExploration(exploration.toCreateRequest(clientSessionId));
            } else {
                promise = services.modifyExploration(exploration.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!exploration.id) {
                        reject({ message: 'Unable to add exploration' });
                    } else {
                        reject({ message: 'Unable to save exploration' });
                    }
                    return;
                }

                const explorationBasicInfo = InsightsExplorerBasicInfo.of(data.result);

                if (!exploration.id || saveAs) {
                    addExplorationToInsightsExplorerList(explorationBasicInfo);
                } else {
                    updateExplorationInInsightsExplorerList(explorationBasicInfo);
                }

                resolve(InsightsExplorer.of(data.result));
            }).catch(error => {
                logger.error('failed to save exploration', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!exploration.id) {
                        reject({ message: 'Unable to add exploration' });
                    } else {
                        reject({ message: 'Unable to save exploration' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeExplorationDisplayOrder({ explorationId, from, to }: { explorationId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentExploration: InsightsExplorerBasicInfo | null = null;

            for (const exploration of allExplorationBasicInfos.value) {
                if (exploration.id === explorationId) {
                    currentExploration = exploration;
                    break;
                }
            }

            if (!currentExploration || !allExplorationBasicInfos.value[to]) {
                reject({ message: 'Unable to move exploration' });
                return;
            }

            if (!insightsExplorerListStateInvalid.value) {
                updateInsightsExplorerListInvalidState(true);
            }

            updateExplorerDisplayOrderInInsightsExplorerList({ from, to });

            resolve();
        });
    }

    function updateExplorationDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: InsightsExplorerNewDisplayOrderRequest[] = [];

        for (const [exploration, index] of itemAndIndex(allExplorationBasicInfos.value)) {
            newDisplayOrders.push({
                id: exploration.id,
                displayOrder: index + 1
            });
        }

        return new Promise((resolve, reject) => {
            services.moveExploration({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move exploration' });
                    return;
                }

                loadAllExplorationBasicInfos({ force: false }).finally(() => {
                    if (insightsExplorerListStateInvalid.value) {
                        updateInsightsExplorerListInvalidState(false);
                    }

                    resolve(data.result);
                });
            }).catch(error => {
                logger.error('failed to save explorations display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move exploration' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function hideExploration({ exploration, hidden }: { exploration: InsightsExplorer | InsightsExplorerBasicInfo, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideExploration({
                id: exploration.id,
                hidden: hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this exploration' });
                    } else {
                        reject({ message: 'Unable to unhide this exploration' });
                    }
                    return;
                }

                exploration.hidden = hidden;
                updateExplorationVisibilityInInsightsExplorerList({ explorationId: exploration.id, hidden });

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to change exploration visibility', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this exploration' });
                    } else {
                        reject({ message: 'Unable to unhide this exploration' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteExploration({ exploration, beforeResolve }: { exploration: InsightsExplorer, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteExploration({
                id: exploration.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this exploration' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeExplorationFromInsightsExplorerList(exploration);
                    });
                } else {
                    removeExplorationFromInsightsExplorerList(exploration);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete exploration', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this exploration' });
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
        allExplorationBasicInfos,
        allExplorationBasicInfosMap,
        currentExploration,
        insightsExplorerListStateInvalid,
        // computed
        isUsingAmountRange,
        filteredTransactionsInDataTable,
        filteredTransactionsInDataTableStatistic,
        categoriedTransactionExplorerData,
        categoriedTransactions,
        // functions
        updateTransactionExplorerInvalidState,
        updateInsightsExplorerListInvalidState,
        updateCurrentExploration,
        resetTransactionExplorers,
        initTransactionExplorerFilter,
        updateTransactionExplorerFilter,
        getTransactionExplorerPageParams,
        loadAllTransactions,
        loadAllExplorationBasicInfos,
        getExploration,
        saveExploration,
        changeExplorationDisplayOrder,
        updateExplorationDisplayOrders,
        hideExploration,
        deleteExploration
    };
});
