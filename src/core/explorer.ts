import { type NameValue } from '@/core/base.ts';
import { DateRange } from '@/core/datetime.ts';

export enum TransactionExplorerConditionRelation {
    First = 'first',
    And = 'and',
    Or = 'or'
}

export const TransactionExplorerConditionRelationPriority: Record<TransactionExplorerConditionRelation, number> = {
    [TransactionExplorerConditionRelation.First]: 0,
    [TransactionExplorerConditionRelation.Or]: 1,
    [TransactionExplorerConditionRelation.And]: 2
};


export enum TransactionExplorerConditionFieldType {
    TransactionType = 'transactionType',
    TransactionCategory = 'transactionCategory',
    SourceAccount = 'sourceAccount',
    DestinationAccount = 'destinationAccount',
    SourceAmount = 'sourceAmount',
    DestinationAmount = 'destinationAmount',
    TransactionTag = 'transactionTag',
    Description = 'description'
}

export class TransactionExplorerConditionField implements NameValue {
    private static readonly allInstances: TransactionExplorerConditionField[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerConditionField> = {};

    public static readonly TransactionType = new TransactionExplorerConditionField('Transaction Type', TransactionExplorerConditionFieldType.TransactionType);
    public static readonly TransactionCategory = new TransactionExplorerConditionField('Category', TransactionExplorerConditionFieldType.TransactionCategory);
    public static readonly SourceAccount = new TransactionExplorerConditionField('Source Account', TransactionExplorerConditionFieldType.SourceAccount);
    public static readonly DestinationAccount = new TransactionExplorerConditionField('Destination Account', TransactionExplorerConditionFieldType.DestinationAccount);
    public static readonly SourceAmount = new TransactionExplorerConditionField('Amount', TransactionExplorerConditionFieldType.SourceAmount);
    public static readonly DestinationAmount = new TransactionExplorerConditionField('Transfer In Amount', TransactionExplorerConditionFieldType.DestinationAmount);
    public static readonly TransactionTag = new TransactionExplorerConditionField('Tags', TransactionExplorerConditionFieldType.TransactionTag);
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
    NotEndsWith = 'notEndsWith'
}

export class TransactionExplorerConditionOperator implements NameValue {
    private static readonly allInstances: TransactionExplorerConditionOperator[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerConditionOperator> = {};

    public static readonly In = new TransactionExplorerConditionOperator('In', TransactionExplorerConditionOperatorType.In);
    public static readonly GreaterThan = new TransactionExplorerConditionOperator('Greater than', TransactionExplorerConditionOperatorType.GreaterThan);
    public static readonly LessThan = new TransactionExplorerConditionOperator('Less than', TransactionExplorerConditionOperatorType.LessThan);
    public static readonly Equals = new TransactionExplorerConditionOperator('Equal to', TransactionExplorerConditionOperatorType.Equals);
    public static readonly NotEquals = new TransactionExplorerConditionOperator('Not equal to', TransactionExplorerConditionOperatorType.NotEquals);
    public static readonly Between = new TransactionExplorerConditionOperator('Between', TransactionExplorerConditionOperatorType.Between);
    public static readonly NotBetween = new TransactionExplorerConditionOperator('Not between', TransactionExplorerConditionOperatorType.NotBetween);
    public static readonly HasAny = new TransactionExplorerConditionOperator('Has any', TransactionExplorerConditionOperatorType.HasAny);
    public static readonly HasAll = new TransactionExplorerConditionOperator('Has all', TransactionExplorerConditionOperatorType.HasAll);
    public static readonly NotHasAny = new TransactionExplorerConditionOperator('Not has any', TransactionExplorerConditionOperatorType.NotHasAny);
    public static readonly NotHasAll = new TransactionExplorerConditionOperator('Not has all', TransactionExplorerConditionOperatorType.NotHasAll);
    public static readonly IsEmpty = new TransactionExplorerConditionOperator('Is empty', TransactionExplorerConditionOperatorType.IsEmpty);
    public static readonly IsNotEmpty = new TransactionExplorerConditionOperator('Is not empty', TransactionExplorerConditionOperatorType.IsNotEmpty);
    public static readonly Contains = new TransactionExplorerConditionOperator('Contains', TransactionExplorerConditionOperatorType.Contains);
    public static readonly NotContains = new TransactionExplorerConditionOperator('Not contains', TransactionExplorerConditionOperatorType.NotContains);
    public static readonly StartsWith = new TransactionExplorerConditionOperator('Starts with', TransactionExplorerConditionOperatorType.StartsWith);
    public static readonly NotStartsWith = new TransactionExplorerConditionOperator('Not starts with', TransactionExplorerConditionOperatorType.NotStartsWith);
    public static readonly EndsWith = new TransactionExplorerConditionOperator('Ends with', TransactionExplorerConditionOperatorType.EndsWith);
    public static readonly NotEndsWith = new TransactionExplorerConditionOperator('Not ends with', TransactionExplorerConditionOperatorType.NotEndsWith);

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

export enum TransactionExplorerChartTypeValue {
    Pie = 'pie',
    ColumnStacked = 'columnStacked',
    Column100PercentStacked = 'column100%Stacked',
    ColumnGrouped = 'columnGrouped',
    LineGrouped = 'lineGrouped',
    AreaStacked = 'areaStacked',
    Area100PercentStacked = 'area100%Stacked',
    BubbleGrouped = 'bubbleGrouped',
    Radar = 'radar'
}

export class TransactionExplorerChartType implements NameValue {
    private static readonly allInstances: TransactionExplorerChartType[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerChartType> = {};

    public static readonly Pie = new TransactionExplorerChartType('Pie Chart', TransactionExplorerChartTypeValue.Pie, false);
    public static readonly Radar = new TransactionExplorerChartType('Radar Chart', TransactionExplorerChartTypeValue.Radar, false);

    public static readonly Default = TransactionExplorerChartType.Pie;

    public readonly name: string;
    public readonly value: TransactionExplorerChartTypeValue;
    public readonly seriesDimensionRequired: boolean;

    private constructor(name: string, value: TransactionExplorerChartTypeValue, seriesDimensionRequired: boolean) {
        this.name = name;
        this.value = value;
        this.seriesDimensionRequired = seriesDimensionRequired;

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
    TransactionType = 'transactionType',
    SourceAccount = 'sourceAccount',
    SourceAccountCategory = 'sourceAccountCategory',
    SourceAccountCurrency = 'sourceAccountCurrency',
    DestinationAccount = 'destinationAccount',
    DestinationAccountCategory = 'destinationAccountCategory',
    DestinationAccountCurrency = 'destinationAccountCurrency',
    SourceAmount = 'sourceAmount',
    DestinationAmount = 'destinationAmount',
    PrimaryCategory = 'primaryCategory',
    SecondaryCategory = 'secondaryCategory'
}

export class TransactionExplorerDataDimension implements NameValue {
    private static readonly allInstances: TransactionExplorerDataDimension[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerDataDimension> = {};

    public static readonly None = new TransactionExplorerDataDimension('None', TransactionExplorerDataDimensionType.None);
    public static readonly Query = new TransactionExplorerDataDimension('Query', TransactionExplorerDataDimensionType.Query);
    public static readonly DateTime = new TransactionExplorerDataDimension('Transaction Time', TransactionExplorerDataDimensionType.DateTime);
    public static readonly DateTimeByYearMonthDay = new TransactionExplorerDataDimension('Transaction Date', TransactionExplorerDataDimensionType.DateTimeByYearMonthDay);
    public static readonly DateTimeByYearMonth = new TransactionExplorerDataDimension('Transaction Year-Month', TransactionExplorerDataDimensionType.DateTimeByYearMonth);
    public static readonly DateTimeByYearQuarter = new TransactionExplorerDataDimension('Transaction Year-Quarter', TransactionExplorerDataDimensionType.DateTimeByYearQuarter);
    public static readonly DateTimeByYear = new TransactionExplorerDataDimension('Transaction Year', TransactionExplorerDataDimensionType.DateTimeByYear);
    public static readonly DateTimeByFiscalYear = new TransactionExplorerDataDimension('Transaction Fiscal Year', TransactionExplorerDataDimensionType.DateTimeByFiscalYear);
    public static readonly DateTimeByDayOfWeek = new TransactionExplorerDataDimension('Transaction Day of Week', TransactionExplorerDataDimensionType.DateTimeByDayOfWeek);
    public static readonly DateTimeByDayOfMonth = new TransactionExplorerDataDimension('Transaction Day of Month', TransactionExplorerDataDimensionType.DateTimeByDayOfMonth);
    public static readonly DateTimeByMonthOfYear = new TransactionExplorerDataDimension('Transaction Month of Year', TransactionExplorerDataDimensionType.DateTimeByMonthOfYear);
    public static readonly DateTimeByQuarterOfYear = new TransactionExplorerDataDimension('Transaction Quarter of Year', TransactionExplorerDataDimensionType.DateTimeByQuarterOfYear);
    public static readonly TransactionType = new TransactionExplorerDataDimension('Transaction Type', TransactionExplorerDataDimensionType.TransactionType);
    public static readonly SourceAccount = new TransactionExplorerDataDimension('Source Account', TransactionExplorerDataDimensionType.SourceAccount);
    public static readonly SourceAccountCategory = new TransactionExplorerDataDimension('Source Account Category', TransactionExplorerDataDimensionType.SourceAccountCategory);
    public static readonly SourceAccountCurrency = new TransactionExplorerDataDimension('Source Account Currency', TransactionExplorerDataDimensionType.SourceAccountCurrency);
    public static readonly DestinationAccount = new TransactionExplorerDataDimension('Destination Account', TransactionExplorerDataDimensionType.DestinationAccount);
    public static readonly DestinationAccountCategory = new TransactionExplorerDataDimension('Destination Account Category', TransactionExplorerDataDimensionType.DestinationAccountCategory);
    public static readonly DestinationAccountCurrency = new TransactionExplorerDataDimension('Destination Account Currency', TransactionExplorerDataDimensionType.DestinationAccountCurrency);
    public static readonly PrimaryCategory = new TransactionExplorerDataDimension('Primary Category', TransactionExplorerDataDimensionType.PrimaryCategory);
    public static readonly SecondaryCategory = new TransactionExplorerDataDimension('Secondary Category', TransactionExplorerDataDimensionType.SecondaryCategory);
    public static readonly SourceAmount = new TransactionExplorerDataDimension('Amount', TransactionExplorerDataDimensionType.SourceAmount);
    public static readonly DestinationAmount = new TransactionExplorerDataDimension('Transfer In Amount', TransactionExplorerDataDimensionType.DestinationAmount);

    public static readonly CategoryDimensionDefault = TransactionExplorerDataDimension.Query;
    public static readonly SeriesDimensionDefault = TransactionExplorerDataDimension.None;

    public readonly name: string;
    public readonly value: TransactionExplorerDataDimensionType;

    private constructor(name: string, value: TransactionExplorerDataDimensionType) {
        this.name = name;
        this.value = value;

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
    SourceAmountSum = 'sourceAmountSum',
    SourceAmountAverage = 'sourceAmountAverage',
    SourceAmountMedian = 'sourceAmountMedian',
    SourceAmountMinimum = 'sourceAmountMinimum',
    SourceAmountMaximum = 'sourceAmountMaximum'
}

export class TransactionExplorerValueMetric implements NameValue {
    private static readonly allInstances: TransactionExplorerValueMetric[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExplorerValueMetric> = {};

    public static readonly TransactionCount = new TransactionExplorerValueMetric('Transaction Count', TransactionExplorerValueMetricType.TransactionCount, false);
    public static readonly SourceAmountSum = new TransactionExplorerValueMetric('Total Amount', TransactionExplorerValueMetricType.SourceAmountSum, true);
    public static readonly SourceAmountAverage = new TransactionExplorerValueMetric('Average Amount', TransactionExplorerValueMetricType.SourceAmountAverage, true);
    public static readonly SourceAmountMedian = new TransactionExplorerValueMetric('Median Amount', TransactionExplorerValueMetricType.SourceAmountMedian, true);
    public static readonly SourceAmountMinimum = new TransactionExplorerValueMetric('Minimum Amount', TransactionExplorerValueMetricType.SourceAmountMinimum, true);
    public static readonly SourceAmountMaximum = new TransactionExplorerValueMetric('Maximum Amount', TransactionExplorerValueMetricType.SourceAmountMaximum, true);

    public static readonly Default = TransactionExplorerValueMetric.SourceAmountSum;

    public readonly name: string;
    public readonly value: TransactionExplorerValueMetricType;
    public readonly isAmount: boolean;

    private constructor(name: string, value: TransactionExplorerValueMetricType, isAmount: boolean) {
        this.name = name;
        this.value = value;
        this.isAmount = isAmount;

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

export const DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE: DateRange = DateRange.ThisMonth;
