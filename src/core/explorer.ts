import { type NameValue } from '@/core/base.ts';
import { DateRange } from '@/core/datetime.ts';
import { ChartSortingType } from '@/core/statistics.ts';

export enum TransactionExplorerConditionRelation {
    First = 'first',
    And = 'and',
    Or = 'or',
    AndSub = 'and(',
    OrSub = 'or(',
    SubEnd = ')'
}

export type TransactionExplorerSubConditionStartRelation = '(';
export const TransactionExplorerSubConditionStartRelationPlaceholder: TransactionExplorerSubConditionStartRelation = '(';

export const TransactionExplorerConditionRelationPriority: Record<TransactionExplorerConditionRelation, number> = {
    [TransactionExplorerConditionRelation.First]: 0,
    [TransactionExplorerConditionRelation.Or]: 1,
    [TransactionExplorerConditionRelation.And]: 2,
    [TransactionExplorerConditionRelation.AndSub]: 0,
    [TransactionExplorerConditionRelation.OrSub]: 0,
    [TransactionExplorerConditionRelation.SubEnd]: 0
};


export enum TransactionExplorerConditionFieldType {
    Undefined = 'undefined',
    TransactionTimeDayOfWeek = 'transactionTimeDayOfWeek',
    TransactionTimeDayOfMonth = 'transactionTimeDayOfMonth',
    TransactionTimeMonthOfYear = 'transactionTimeMonthOfYear',
    TransactionTimeHourOfDay = 'transactionTimeHourOfDay',
    TransactionTimezone = 'transactionTimezone',
    TransactionType = 'transactionType',
    TransactionCategory = 'transactionCategory',
    SourceAccount = 'sourceAccount',
    DestinationAccount = 'destinationAccount',
    SourceAmount = 'sourceAmount',
    DestinationAmount = 'destinationAmount',
    GeoLocation = 'geoLocation',
    TransactionTag = 'transactionTag',
    Pictures = 'pictures',
    Description = 'description'
}

export class TransactionExplorerConditionField implements NameValue {
    private static readonly allInstances: TransactionExplorerConditionField[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerConditionField> = {};

    public static readonly TransactionTimeDayOfWeek = new TransactionExplorerConditionField('Transaction Day of Week', TransactionExplorerConditionFieldType.TransactionTimeDayOfWeek);
    public static readonly TransactionTimeDayOfMonth = new TransactionExplorerConditionField('Transaction Day of Month', TransactionExplorerConditionFieldType.TransactionTimeDayOfMonth)
    public static readonly TransactionTimeMonthOfYear = new TransactionExplorerConditionField('Transaction Month of Year', TransactionExplorerConditionFieldType.TransactionTimeMonthOfYear);
    public static readonly TransactionTimeHourOfDay = new TransactionExplorerConditionField('Transaction Hour of Day', TransactionExplorerConditionFieldType.TransactionTimeHourOfDay);
    public static readonly TransactionTimezone = new TransactionExplorerConditionField('Transaction Timezone', TransactionExplorerConditionFieldType.TransactionTimezone);
    public static readonly TransactionType = new TransactionExplorerConditionField('Transaction Type', TransactionExplorerConditionFieldType.TransactionType);
    public static readonly TransactionCategory = new TransactionExplorerConditionField('Category', TransactionExplorerConditionFieldType.TransactionCategory);
    public static readonly SourceAccount = new TransactionExplorerConditionField('Source Account', TransactionExplorerConditionFieldType.SourceAccount);
    public static readonly DestinationAccount = new TransactionExplorerConditionField('Destination Account', TransactionExplorerConditionFieldType.DestinationAccount);
    public static readonly SourceAmount = new TransactionExplorerConditionField('Amount', TransactionExplorerConditionFieldType.SourceAmount);
    public static readonly DestinationAmount = new TransactionExplorerConditionField('Transfer In Amount', TransactionExplorerConditionFieldType.DestinationAmount);
    public static readonly GeoLocation = new TransactionExplorerConditionField('Geographic Location', TransactionExplorerConditionFieldType.GeoLocation);
    public static readonly TransactionTag = new TransactionExplorerConditionField('Tags', TransactionExplorerConditionFieldType.TransactionTag);
    public static readonly Pictures = new TransactionExplorerConditionField('Pictures', TransactionExplorerConditionFieldType.Pictures);
    public static readonly Description = new TransactionExplorerConditionField('Description', TransactionExplorerConditionFieldType.Description);

    public readonly name: string;
    public readonly value: TransactionExplorerConditionFieldType;

    private constructor(name: string, value: TransactionExplorerConditionFieldType) {
        this.name = name;
        this.value = value;

        TransactionExplorerConditionField.allInstances.push(this);
        TransactionExplorerConditionField.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExplorerConditionField[] {
        return TransactionExplorerConditionField.allInstances;
    }

    public static valueOf(value: string): TransactionExplorerConditionField | undefined {
        return TransactionExplorerConditionField.allInstancesByValue[value];
    }
}

export enum TransactionExplorerConditionOperatorType {
    In = 'in',
    NotIn = 'notIn',
    GreaterThan = 'greaterThan',
    LessThan = 'lessThan',
    Equals = 'equals',
    NotEquals = 'notEquals',
    Between = 'between',
    NotBetween = 'notBetween',
    HasAny = 'hasAny',
    HasAll = 'hasAll',
    NotHasAny = 'notHasAny',
    NotHasAll = 'notHasAll',
    IsEmpty = 'isEmpty',
    IsNotEmpty = 'isNotEmpty',
    Contains = 'contains',
    NotContains = 'notContains',
    StartsWith = 'startsWith',
    NotStartsWith = 'notStartsWith',
    EndsWith = 'endsWith',
    NotEndsWith = 'notEndsWith',
    RegexMatch = 'regexMatch',
    NotRegexMatch = 'notRegexMatch',
    MinuteOffsetBetween = 'minuteOffsetBetween',
    MinuteOffsetNotBetween = 'minuteOffsetNotBetween',
    LatitudeBetween = 'latitudeBetween',
    LatitudeNotBetween = 'latitudeNotBetween',
    LongitudeBetween = 'longitudeBetween',
    LongitudeNotBetween = 'longitudeNotBetween'
}

export class TransactionExplorerConditionOperator implements NameValue {
    private static readonly allInstances: TransactionExplorerConditionOperator[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerConditionOperator> = {};

    public static readonly In = new TransactionExplorerConditionOperator('In', TransactionExplorerConditionOperatorType.In);
    public static readonly NotIn = new TransactionExplorerConditionOperator('Not in', TransactionExplorerConditionOperatorType.NotIn);
    public static readonly GreaterThan = new TransactionExplorerConditionOperator('Greater than', TransactionExplorerConditionOperatorType.GreaterThan);
    public static readonly LessThan = new TransactionExplorerConditionOperator('Less than', TransactionExplorerConditionOperatorType.LessThan);
    public static readonly Equals = new TransactionExplorerConditionOperator('Equal to', TransactionExplorerConditionOperatorType.Equals);
    public static readonly NotEquals = new TransactionExplorerConditionOperator('Not equal to', TransactionExplorerConditionOperatorType.NotEquals);
    public static readonly Between = new TransactionExplorerConditionOperator('Between', TransactionExplorerConditionOperatorType.Between);
    public static readonly NotBetween = new TransactionExplorerConditionOperator('Not between', TransactionExplorerConditionOperatorType.NotBetween);
    public static readonly HasAny = new TransactionExplorerConditionOperator('Has any', TransactionExplorerConditionOperatorType.HasAny);
    public static readonly HasAll = new TransactionExplorerConditionOperator('Has all', TransactionExplorerConditionOperatorType.HasAll);
    public static readonly NotHasAny = new TransactionExplorerConditionOperator('Does not have any', TransactionExplorerConditionOperatorType.NotHasAny);
    public static readonly NotHasAll = new TransactionExplorerConditionOperator('Does not have all', TransactionExplorerConditionOperatorType.NotHasAll);
    public static readonly IsEmpty = new TransactionExplorerConditionOperator('Is empty', TransactionExplorerConditionOperatorType.IsEmpty);
    public static readonly IsNotEmpty = new TransactionExplorerConditionOperator('Is not empty', TransactionExplorerConditionOperatorType.IsNotEmpty);
    public static readonly Contains = new TransactionExplorerConditionOperator('Contains', TransactionExplorerConditionOperatorType.Contains);
    public static readonly NotContains = new TransactionExplorerConditionOperator('Does not contain', TransactionExplorerConditionOperatorType.NotContains);
    public static readonly StartsWith = new TransactionExplorerConditionOperator('Starts with', TransactionExplorerConditionOperatorType.StartsWith);
    public static readonly NotStartsWith = new TransactionExplorerConditionOperator('Does not start with', TransactionExplorerConditionOperatorType.NotStartsWith);
    public static readonly EndsWith = new TransactionExplorerConditionOperator('Ends with', TransactionExplorerConditionOperatorType.EndsWith);
    public static readonly NotEndsWith = new TransactionExplorerConditionOperator('Does not end with', TransactionExplorerConditionOperatorType.NotEndsWith);
    public static readonly RegexMatch = new TransactionExplorerConditionOperator('Matches regex', TransactionExplorerConditionOperatorType.RegexMatch);
    public static readonly NotRegexMatch = new TransactionExplorerConditionOperator('Does not match regex', TransactionExplorerConditionOperatorType.NotRegexMatch);
    public static readonly MinuteOffsetBetween = new TransactionExplorerConditionOperator('Minute offset is between', TransactionExplorerConditionOperatorType.MinuteOffsetBetween);
    public static readonly MinuteOffsetNotBetween = new TransactionExplorerConditionOperator('Minute offset is not between', TransactionExplorerConditionOperatorType.MinuteOffsetNotBetween);
    public static readonly LatitudeBetween = new TransactionExplorerConditionOperator('Latitude is between', TransactionExplorerConditionOperatorType.LatitudeBetween);
    public static readonly LatitudeNotBetween = new TransactionExplorerConditionOperator('Latitude is not between', TransactionExplorerConditionOperatorType.LatitudeNotBetween);
    public static readonly LongitudeBetween = new TransactionExplorerConditionOperator('Longitude is between', TransactionExplorerConditionOperatorType.LongitudeBetween);
    public static readonly LongitudeNotBetween = new TransactionExplorerConditionOperator('Longitude is not between', TransactionExplorerConditionOperatorType.LongitudeNotBetween);

    public readonly name: string;
    public readonly value: TransactionExplorerConditionOperatorType;

    private constructor(name: string, value: TransactionExplorerConditionOperatorType) {
        this.name = name;
        this.value = value;

        TransactionExplorerConditionOperator.allInstances.push(this);
        TransactionExplorerConditionOperator.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExplorerConditionOperator[] {
        return TransactionExplorerConditionOperator.allInstances;
    }

    public static valueOf(value: string): TransactionExplorerConditionOperator | undefined {
        return TransactionExplorerConditionOperator.allInstancesByValue[value];
    }
}

export enum TransactionExplorerDataDimensionType {
    None = 'none',
    Query = 'query',
    DateTime = 'dateTime',
    DateTimeByYearMonthDay = 'dateTimeByYearMonthDay',
    DateTimeByYearMonth = 'dateTimeByYearMonth',
    DateTimeByYearQuarter = 'dateTimeByYearQuarter',
    DateTimeByYear = 'dateTimeByYear',
    DateTimeByFiscalYear = 'dateTimeByFiscalYear',
    DateTimeByDayOfWeek = 'dateTimeByDayOfWeek',
    DateTimeByDayOfMonth = 'dateTimeByDayOfMonth',
    DateTimeByMonthOfYear = 'dateTimeByMonthOfYear',
    DateTimeByQuarterOfYear = 'dateTimeByQuarterOfYear',
    DateTimeByHourOfDay = 'dateTimeByHourOfDay',
    TimezoneOffset = 'timezoneOffset',
    TransactionType = 'transactionType',
    SourceAccount = 'sourceAccount',
    SourceAccountCategory = 'sourceAccountCategory',
    SourceAccountCurrency = 'sourceAccountCurrency',
    DestinationAccount = 'destinationAccount',
    DestinationAccountCategory = 'destinationAccountCategory',
    DestinationAccountCurrency = 'destinationAccountCurrency',
    PrimaryCategory = 'primaryCategory',
    SecondaryCategory = 'secondaryCategory',
    SourceAmount = 'sourceAmount',
    DestinationAmount = 'destinationAmount',
    SourceAmountRangeEqualFrequency = 'sourceAmountRangeEqualFrequency',
    SourceAmountRangeEqualWidth = 'sourceAmountRangeEqualWidth',
    SourceAmountRangeLogScale = 'sourceAmountRangeLogScale',
    SourceAmountRangeStandardDeviation = 'sourceAmountRangeStandardDeviation',
    SourceAmountRangeNaturalBreaks = 'sourceAmountRangeNaturalBreaks',
    DestinationAmountRangeEqualFrequency = 'destinationAmountRangeEqualFrequency',
    DestinationAmountRangeEqualWidth = 'destinationAmountRangeEqualWidth',
    DestinationAmountRangeLogScale = 'destinationAmountRangeLogScale',
    DestinationAmountRangeStandardDeviation = 'destinationAmountRangeStandardDeviation',
    DestinationAmountRangeNaturalBreaks = 'destinationAmountRangeNaturalBreaks'
}

export class TransactionExplorerDataDimension implements NameValue {
    private static readonly allInstances: TransactionExplorerDataDimension[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerDataDimension> = {};

    public static readonly None = new TransactionExplorerDataDimension('None', TransactionExplorerDataDimensionType.None, false, false);
    public static readonly Query = new TransactionExplorerDataDimension('Query', TransactionExplorerDataDimensionType.Query, false, false);
    public static readonly DateTime = new TransactionExplorerDataDimension('Transaction Time', TransactionExplorerDataDimensionType.DateTime, false, false);
    public static readonly DateTimeByYearMonthDay = new TransactionExplorerDataDimension('Transaction Date', TransactionExplorerDataDimensionType.DateTimeByYearMonthDay, false, false);
    public static readonly DateTimeByYearMonth = new TransactionExplorerDataDimension('Transaction Year-Month', TransactionExplorerDataDimensionType.DateTimeByYearMonth, false, false);
    public static readonly DateTimeByYearQuarter = new TransactionExplorerDataDimension('Transaction Year-Quarter', TransactionExplorerDataDimensionType.DateTimeByYearQuarter, false, false);
    public static readonly DateTimeByYear = new TransactionExplorerDataDimension('Transaction Year', TransactionExplorerDataDimensionType.DateTimeByYear, false, false);
    public static readonly DateTimeByFiscalYear = new TransactionExplorerDataDimension('Transaction Fiscal Year', TransactionExplorerDataDimensionType.DateTimeByFiscalYear, false, false);
    public static readonly DateTimeByDayOfWeek = new TransactionExplorerDataDimension('Transaction Day of Week', TransactionExplorerDataDimensionType.DateTimeByDayOfWeek, false, false);
    public static readonly DateTimeByDayOfMonth = new TransactionExplorerDataDimension('Transaction Day of Month', TransactionExplorerDataDimensionType.DateTimeByDayOfMonth, false, false);
    public static readonly DateTimeByMonthOfYear = new TransactionExplorerDataDimension('Transaction Month of Year', TransactionExplorerDataDimensionType.DateTimeByMonthOfYear, false, false);
    public static readonly DateTimeByQuarterOfYear = new TransactionExplorerDataDimension('Transaction Quarter of Year', TransactionExplorerDataDimensionType.DateTimeByQuarterOfYear, false, false);
    public static readonly DateTimeByHourOfDay = new TransactionExplorerDataDimension('Transaction Hour of Day', TransactionExplorerDataDimensionType.DateTimeByHourOfDay, false, false);
    public static readonly TimezoneOffset = new TransactionExplorerDataDimension('Transaction Timezone', TransactionExplorerDataDimensionType.TimezoneOffset, false, false);
    public static readonly TransactionType = new TransactionExplorerDataDimension('Transaction Type', TransactionExplorerDataDimensionType.TransactionType, false, false);
    public static readonly SourceAccount = new TransactionExplorerDataDimension('Source Account', TransactionExplorerDataDimensionType.SourceAccount, false, false);
    public static readonly SourceAccountCategory = new TransactionExplorerDataDimension('Source Account Category', TransactionExplorerDataDimensionType.SourceAccountCategory, false, false);
    public static readonly SourceAccountCurrency = new TransactionExplorerDataDimension('Source Account Currency', TransactionExplorerDataDimensionType.SourceAccountCurrency, false, false);
    public static readonly DestinationAccount = new TransactionExplorerDataDimension('Destination Account', TransactionExplorerDataDimensionType.DestinationAccount, false, false);
    public static readonly DestinationAccountCategory = new TransactionExplorerDataDimension('Destination Account Category', TransactionExplorerDataDimensionType.DestinationAccountCategory, false, false);
    public static readonly DestinationAccountCurrency = new TransactionExplorerDataDimension('Destination Account Currency', TransactionExplorerDataDimensionType.DestinationAccountCurrency, false, false);
    public static readonly PrimaryCategory = new TransactionExplorerDataDimension('Primary Category', TransactionExplorerDataDimensionType.PrimaryCategory, false, false);
    public static readonly SecondaryCategory = new TransactionExplorerDataDimension('Secondary Category', TransactionExplorerDataDimensionType.SecondaryCategory, false, false);
    public static readonly SourceAmount = new TransactionExplorerDataDimension('Amount', TransactionExplorerDataDimensionType.SourceAmount, false, false);
    public static readonly DestinationAmount = new TransactionExplorerDataDimension('Transfer In Amount', TransactionExplorerDataDimensionType.DestinationAmount, false, false);
    public static readonly SourceAmountRangeEqualFrequency = new TransactionExplorerDataDimension('Amount Range (Equal Frequency)', TransactionExplorerDataDimensionType.SourceAmountRangeEqualFrequency, true, false);
    public static readonly SourceAmountRangeEqualWidth = new TransactionExplorerDataDimension('Amount Range (Equal Width)', TransactionExplorerDataDimensionType.SourceAmountRangeEqualWidth, true, false);
    public static readonly SourceAmountRangeLogScale = new TransactionExplorerDataDimension('Amount Range (Log Scale)', TransactionExplorerDataDimensionType.SourceAmountRangeLogScale, true, false);
    public static readonly SourceAmountRangeStandardDeviation = new TransactionExplorerDataDimension('Amount Range (Standard Deviation)', TransactionExplorerDataDimensionType.SourceAmountRangeStandardDeviation, true, false);
    public static readonly SourceAmountRangeNaturalBreaks = new TransactionExplorerDataDimension('Amount Range (Natural Breaks)', TransactionExplorerDataDimensionType.SourceAmountRangeNaturalBreaks, true, false);
    public static readonly DestinationAmountRangeEqualFrequency = new TransactionExplorerDataDimension('Transfer In Amount Range (Equal Frequency)', TransactionExplorerDataDimensionType.DestinationAmountRangeEqualFrequency, false, true);
    public static readonly DestinationAmountRangeEqualWidth = new TransactionExplorerDataDimension('Transfer In Amount Range (Equal Width)', TransactionExplorerDataDimensionType.DestinationAmountRangeEqualWidth, false, true);
    public static readonly DestinationAmountRangeLogScale = new TransactionExplorerDataDimension('Transfer In Amount Range (Log Scale)', TransactionExplorerDataDimensionType.DestinationAmountRangeLogScale, false, true);
    public static readonly DestinationAmountRangeStandardDeviation = new TransactionExplorerDataDimension('Transfer In Amount Range (Standard Deviation)', TransactionExplorerDataDimensionType.DestinationAmountRangeStandardDeviation, false, true);
    public static readonly DestinationAmountRangeNaturalBreaks = new TransactionExplorerDataDimension('Transfer In Amount Range (Natural Breaks)', TransactionExplorerDataDimensionType.DestinationAmountRangeNaturalBreaks, false, true);

    public static readonly CategoryDimensionDefault = TransactionExplorerDataDimension.Query;
    public static readonly SeriesDimensionDefault = TransactionExplorerDataDimension.None;

    public readonly name: string;
    public readonly value: TransactionExplorerDataDimensionType;
    public readonly isSourceAmountRange: boolean;
    public readonly isDestinationAmountRange: boolean;

    private constructor(name: string, value: TransactionExplorerDataDimensionType, isSourceAmountRange: boolean, isDestinationAmountRange: boolean) {
        this.name = name;
        this.value = value;
        this.isSourceAmountRange = isSourceAmountRange;
        this.isDestinationAmountRange = isDestinationAmountRange;

        TransactionExplorerDataDimension.allInstances.push(this);
        TransactionExplorerDataDimension.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExplorerDataDimension[] {
        return TransactionExplorerDataDimension.allInstances;
    }

    public static valueOf(value: string): TransactionExplorerDataDimension | undefined {
        return TransactionExplorerDataDimension.allInstancesByValue[value];
    }
}

export enum TransactionExplorerValueMetricType {
    TransactionCount = 'transactionCount',
    ActiveTransactionDays = 'activeTransactionDays',
    TransactionsPerActiveDay = 'transactionsPerActiveDay',
    SourceAmountSum = 'sourceAmountSum',
    SourceIncomeAmountSum = 'sourceIncomeAmountSum',
    SourceExpenseAmountSum = 'sourceExpenseAmountSum',
    SourceNetIncomeAmountSum = 'sourceNetIncomeAmountSum',
    SrouceAmountExpenseIncomeRatio = 'sourceExpenseIncomeRatio',
    SourceAmountSavingsRate = 'sourceAmountSavingsRate',
    SourceAmountAverage = 'sourceAmountAverage',
    SourceAmountMedian = 'sourceAmountMedian',
    SourceAmountMinimum = 'sourceAmountMinimum',
    SourceAmountMaximum = 'sourceAmountMaximum',
    SourceAmountQ1Amount = 'sourceQ1Amount',
    SourceAmountQ3Amount = 'sourceQ3Amount',
    SourceAmount10thPercentile = 'source10thPercentileAmount',
    SourceAmount90thPercentile = 'source90thPercentileAmount',
    SourceAmount95thPercentile = 'source95thPercentileAmount',
    SourceAmount99thPercentile = 'source99thPercentileAmount',
    SourceAmountRange = 'sourceAmountRange',
    SourceAmountInterquartileRange = 'sourceAmountInterquartileRange',
    SourceAmountMeanAbsoluteDeviation = 'sourceAmountMeanAbsoluteDeviation',
    SourceAmountMedianAbsoluteDeviation = 'sourceAmountMedianAbsoluteDeviation',
    SourceMaximumAmountShare = 'sourceMaximumAmountShare',
    SourceTop5AmountSum = 'sourceTop5AmountSum',
    SourceTop5AmountShare = 'sourceTop5AmountShare',
    TransactionsForEightyPercentOfSourceAmount = 'transactionsForEightyPercentOfSourceAmount',
    SourceAmountVariance = 'sourceAmountVariance',
    SourceAmountStandardDeviation = 'sourceAmountStandardDeviation',
    SourceAmountCoefficientOfVariation = 'sourceAmountCoefficientOfVariation',
    SourceAmountSkewness = 'sourceAmountSkewness',
    SourceAmountKurtosis = 'sourceAmountKurtosis'
}

export class TransactionExplorerValueMetric implements NameValue {
    private static readonly allInstances: TransactionExplorerValueMetric[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerValueMetric> = {};

    public static readonly TransactionCount = new TransactionExplorerValueMetric('Transaction Count', TransactionExplorerValueMetricType.TransactionCount, false, false, true);
    public static readonly ActiveTransactionDays = new TransactionExplorerValueMetric('Active Transaction Days', TransactionExplorerValueMetricType.ActiveTransactionDays, false, false, true);
    public static readonly TransactionsPerDay = new TransactionExplorerValueMetric('Transactions per Active Day', TransactionExplorerValueMetricType.TransactionsPerActiveDay, false, false, true);
    public static readonly SourceAmountSum = new TransactionExplorerValueMetric('Total Amount', TransactionExplorerValueMetricType.SourceAmountSum, true, false, true);
    public static readonly SourceIncomeAmountSum = new TransactionExplorerValueMetric('Total Income', TransactionExplorerValueMetricType.SourceIncomeAmountSum, true, false, true);
    public static readonly SourceExpenseAmountSum = new TransactionExplorerValueMetric('Total Expense', TransactionExplorerValueMetricType.SourceExpenseAmountSum, true, false, true);
    public static readonly SourceNetIncomeAmountSum = new TransactionExplorerValueMetric('Net Income', TransactionExplorerValueMetricType.SourceNetIncomeAmountSum, true, false, true);
    public static readonly SrouceAmountExpenseIncomeRatio = new TransactionExplorerValueMetric('Expense / Income Ratio', TransactionExplorerValueMetricType.SrouceAmountExpenseIncomeRatio, false, true, false);
    public static readonly SourceAmountSavingsRate = new TransactionExplorerValueMetric('Savings Rate', TransactionExplorerValueMetricType.SourceAmountSavingsRate, false, true, false);
    public static readonly SourceAmountAverage = new TransactionExplorerValueMetric('Average Amount', TransactionExplorerValueMetricType.SourceAmountAverage, true, false, true);
    public static readonly SourceAmountMedian = new TransactionExplorerValueMetric('Median Amount', TransactionExplorerValueMetricType.SourceAmountMedian, true, false, true);
    public static readonly SourceAmountMinimum = new TransactionExplorerValueMetric('Minimum Amount', TransactionExplorerValueMetricType.SourceAmountMinimum, true, false, true);
    public static readonly SourceAmountMaximum = new TransactionExplorerValueMetric('Maximum Amount', TransactionExplorerValueMetricType.SourceAmountMaximum, true, false, true);
    public static readonly SourceAmountQ1Amount = new TransactionExplorerValueMetric('Q1 Amount (First Quartile)', TransactionExplorerValueMetricType.SourceAmountQ1Amount, true, false, true);
    public static readonly SourceAmountQ3Amount = new TransactionExplorerValueMetric('Q3 Amount (Third Quartile)', TransactionExplorerValueMetricType.SourceAmountQ3Amount, true, false, true);
    public static readonly SourceAmount10thPercentile = new TransactionExplorerValueMetric('10th Percentile Amount', TransactionExplorerValueMetricType.SourceAmount10thPercentile, true, false, true);
    public static readonly SourceAmount90thPercentile = new TransactionExplorerValueMetric('90th Percentile Amount', TransactionExplorerValueMetricType.SourceAmount90thPercentile, true, false, true);
    public static readonly SourceAmount95thPercentile = new TransactionExplorerValueMetric('95th Percentile Amount', TransactionExplorerValueMetricType.SourceAmount95thPercentile, true, false, true);
    public static readonly SourceAmount99thPercentile = new TransactionExplorerValueMetric('99th Percentile Amount', TransactionExplorerValueMetricType.SourceAmount99thPercentile, true, false, true);
    public static readonly SourceAmountRange = new TransactionExplorerValueMetric('Range (Max - Min)', TransactionExplorerValueMetricType.SourceAmountRange, true, false, true);
    public static readonly SourceAmountInterquartileRange = new TransactionExplorerValueMetric('Interquartile Range (Q3 - Q1)', TransactionExplorerValueMetricType.SourceAmountInterquartileRange, true, false, true);
    public static readonly SourceAmountMeanAbsoluteDeviation = new TransactionExplorerValueMetric('Mean Absolute Deviation', TransactionExplorerValueMetricType.SourceAmountMeanAbsoluteDeviation, true, false, false);
    public static readonly SourceAmountMedianAbsoluteDeviation = new TransactionExplorerValueMetric('Median Absolute Deviation', TransactionExplorerValueMetricType.SourceAmountMedianAbsoluteDeviation, true, false, false);
    public static readonly SourceMaximumAmountShare = new TransactionExplorerValueMetric('Maximum Amount Share', TransactionExplorerValueMetricType.SourceMaximumAmountShare, false, true, false);
    public static readonly SourceTop5AmountSum = new TransactionExplorerValueMetric('Top 5 Amount Sum', TransactionExplorerValueMetricType.SourceTop5AmountSum, true, false, true);
    public static readonly SourceTop5AmountShare = new TransactionExplorerValueMetric('Top 5 Amount Share', TransactionExplorerValueMetricType.SourceTop5AmountShare, false, true, false);
    public static readonly TransactionsForEightyPercentOfSourceAmount = new TransactionExplorerValueMetric('Transactions for 80% of Amount', TransactionExplorerValueMetricType.TransactionsForEightyPercentOfSourceAmount, false, true, false);
    public static readonly SourceAmountVariance = new TransactionExplorerValueMetric('Variance', TransactionExplorerValueMetricType.SourceAmountVariance, false, false, false);
    public static readonly SourceAmountStandardDeviation = new TransactionExplorerValueMetric('Standard Deviation', TransactionExplorerValueMetricType.SourceAmountStandardDeviation, false, false, false);
    public static readonly SourceAmountCoefficientOfVariation = new TransactionExplorerValueMetric('Coefficient of Variation', TransactionExplorerValueMetricType.SourceAmountCoefficientOfVariation, false, false, false);
    public static readonly SourceAmountSkewness = new TransactionExplorerValueMetric('Skewness', TransactionExplorerValueMetricType.SourceAmountSkewness, false, false, false);
    public static readonly SourceAmountKurtosis = new TransactionExplorerValueMetric('Kurtosis', TransactionExplorerValueMetricType.SourceAmountKurtosis, false, false, false);

    public static readonly Default = TransactionExplorerValueMetric.SourceAmountSum;

    public readonly name: string;
    public readonly value: TransactionExplorerValueMetricType;
    public readonly isAmount: boolean;
    public readonly isPercent: boolean;
    public readonly supportSum: boolean;

    private constructor(name: string, value: TransactionExplorerValueMetricType, isAmount: boolean, isPercent: boolean, supportSum: boolean) {
        this.name = name;
        this.value = value;
        this.isAmount = isAmount;
        this.isPercent = isPercent;
        this.supportSum = supportSum;

        TransactionExplorerValueMetric.allInstances.push(this);
        TransactionExplorerValueMetric.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExplorerValueMetric[] {
        return TransactionExplorerValueMetric.allInstances;
    }

    public static valueOf(value: string): TransactionExplorerValueMetric | undefined {
        return TransactionExplorerValueMetric.allInstancesByValue[value];
    }
}

export enum TransactionExplorerChartTypeValue {
    Pie = 'pie',
    ColumnStacked = 'columnStacked',
    Column100PercentStacked = 'column100%Stacked',
    ColumnGrouped = 'columnGrouped',
    LineGrouped = 'lineGrouped',
    AreaStacked = 'areaStacked',
    Area100PercentStacked = 'area100%Stacked',
    BubbleGrouped = 'bubbleGrouped',
    Radar = 'radar',
    Heatmap = 'heatmap',
    CalendarHeatmap = 'calendarHeatmap'
}

export class TransactionExplorerChartType implements NameValue {
    private static readonly allInstances: TransactionExplorerChartType[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerChartType> = {};

    public static readonly Pie = new TransactionExplorerChartType('Pie Chart', TransactionExplorerChartTypeValue.Pie, undefined, false, undefined);
    public static readonly Radar = new TransactionExplorerChartType('Radar Chart', TransactionExplorerChartTypeValue.Radar, undefined, false, undefined);
    public static readonly ColumnStacked = new TransactionExplorerChartType('Column Chart (Stacked)', TransactionExplorerChartTypeValue.ColumnStacked, undefined, true, undefined);
    public static readonly Column100PercentStacked = new TransactionExplorerChartType('Column Chart (100% Stacked)', TransactionExplorerChartTypeValue.Column100PercentStacked, undefined, true, undefined);
    public static readonly ColumnGrouped = new TransactionExplorerChartType('Column Chart (Grouped)', TransactionExplorerChartTypeValue.ColumnGrouped, undefined, true, undefined);
    public static readonly LineGrouped = new TransactionExplorerChartType('Line Chart (Grouped)', TransactionExplorerChartTypeValue.LineGrouped, undefined, true, undefined);
    public static readonly AreaStacked = new TransactionExplorerChartType('Area Chart (Stacked)', TransactionExplorerChartTypeValue.AreaStacked, undefined, true, undefined);
    public static readonly Area100PercentStacked = new TransactionExplorerChartType('Area Chart (100% Stacked)', TransactionExplorerChartTypeValue.Area100PercentStacked, undefined, true, undefined);
    public static readonly BubbleGrouped = new TransactionExplorerChartType('Bubble Chart (Grouped)', TransactionExplorerChartTypeValue.BubbleGrouped, undefined, true, undefined);
    public static readonly Heatmap = new TransactionExplorerChartType('Heatmap Chart', TransactionExplorerChartTypeValue.Heatmap, undefined, true, undefined);
    public static readonly CalendarHeatmap = new TransactionExplorerChartType('Calendar Heatmap Chart', TransactionExplorerChartTypeValue.CalendarHeatmap, TransactionExplorerDataDimensionType.DateTimeByYearMonthDay, false, ChartSortingType.DisplayOrder.type);

    public static readonly Default = TransactionExplorerChartType.Pie;

    public readonly name: string;
    public readonly value: TransactionExplorerChartTypeValue;
    public readonly fixedCategoryDimension: TransactionExplorerDataDimensionType | undefined;
    public readonly seriesDimensionRequired: boolean;
    public readonly fixedSortingType: number | undefined;

    private constructor(name: string, value: TransactionExplorerChartTypeValue, fixedCategoryDimension: TransactionExplorerDataDimensionType | undefined, seriesDimensionRequired: boolean, fixedSortingType: number | undefined) {
        this.name = name;
        this.value = value;
        this.fixedCategoryDimension = fixedCategoryDimension;
        this.seriesDimensionRequired = seriesDimensionRequired;
        this.fixedSortingType = fixedSortingType;

        TransactionExplorerChartType.allInstances.push(this);
        TransactionExplorerChartType.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExplorerChartType[] {
        return TransactionExplorerChartType.allInstances;
    }

    public static valueOf(value: string): TransactionExplorerChartType | undefined {
        return TransactionExplorerChartType.allInstancesByValue[value];
    }
}

export const DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE: DateRange = DateRange.ThisMonth;
