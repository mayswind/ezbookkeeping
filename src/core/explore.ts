import { type NameValue } from '@/core/base.ts';
import { DateRange } from '@/core/datetime.ts';

export enum TransactionExploreConditionRelation {
    First = 'first',
    And = 'and',
    Or = 'or'
}

export const TransactionExploreConditionRelationPriority: Record<TransactionExploreConditionRelation, number> = {
    [TransactionExploreConditionRelation.First]: 0,
    [TransactionExploreConditionRelation.Or]: 1,
    [TransactionExploreConditionRelation.And]: 2
};


export enum TransactionExploreConditionFieldType {
    TransactionType = 'transactionType',
    TransactionCategory = 'transactionCategory',
    SourceAccount = 'sourceAccount',
    DestinationAccount = 'destinationAccount',
    SourceAmount = 'sourceAmount',
    DestinationAmount = 'destinationAmount',
    TransactionTag = 'transactionTag',
    Description = 'description'
}

export class TransactionExploreConditionField implements NameValue {
    private static readonly allInstances: TransactionExploreConditionField[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreConditionField> = {};

    public static readonly TransactionType = new TransactionExploreConditionField('Transaction Type', TransactionExploreConditionFieldType.TransactionType);
    public static readonly TransactionCategory = new TransactionExploreConditionField('Category', TransactionExploreConditionFieldType.TransactionCategory);
    public static readonly SourceAccount = new TransactionExploreConditionField('Source Account', TransactionExploreConditionFieldType.SourceAccount);
    public static readonly DestinationAccount = new TransactionExploreConditionField('Destination Account', TransactionExploreConditionFieldType.DestinationAccount);
    public static readonly SourceAmount = new TransactionExploreConditionField('Amount', TransactionExploreConditionFieldType.SourceAmount);
    public static readonly DestinationAmount = new TransactionExploreConditionField('Transfer In Amount', TransactionExploreConditionFieldType.DestinationAmount);
    public static readonly TransactionTag = new TransactionExploreConditionField('Tags', TransactionExploreConditionFieldType.TransactionTag);
    public static readonly Description = new TransactionExploreConditionField('Description', TransactionExploreConditionFieldType.Description);

    public readonly name: string;
    public readonly value: TransactionExploreConditionFieldType;

    private constructor(name: string, value: TransactionExploreConditionFieldType) {
        this.name = name;
        this.value = value;

        TransactionExploreConditionField.allInstances.push(this);
        TransactionExploreConditionField.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreConditionField[] {
        return TransactionExploreConditionField.allInstances;
    }

    public static valueOf(value: string): TransactionExploreConditionField | undefined {
        return TransactionExploreConditionField.allInstancesByValue[value];
    }
}

export enum TransactionExploreConditionOperatorType {
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

export class TransactionExploreConditionOperator implements NameValue {
    private static readonly allInstances: TransactionExploreConditionOperator[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreConditionOperator> = {};

    public static readonly In = new TransactionExploreConditionOperator('In', TransactionExploreConditionOperatorType.In);
    public static readonly GreaterThan = new TransactionExploreConditionOperator('Greater than', TransactionExploreConditionOperatorType.GreaterThan);
    public static readonly LessThan = new TransactionExploreConditionOperator('Less than', TransactionExploreConditionOperatorType.LessThan);
    public static readonly Equals = new TransactionExploreConditionOperator('Equal to', TransactionExploreConditionOperatorType.Equals);
    public static readonly NotEquals = new TransactionExploreConditionOperator('Not equal to', TransactionExploreConditionOperatorType.NotEquals);
    public static readonly Between = new TransactionExploreConditionOperator('Between', TransactionExploreConditionOperatorType.Between);
    public static readonly NotBetween = new TransactionExploreConditionOperator('Not between', TransactionExploreConditionOperatorType.NotBetween);
    public static readonly HasAny = new TransactionExploreConditionOperator('Has any', TransactionExploreConditionOperatorType.HasAny);
    public static readonly HasAll = new TransactionExploreConditionOperator('Has all', TransactionExploreConditionOperatorType.HasAll);
    public static readonly NotHasAny = new TransactionExploreConditionOperator('Not has any', TransactionExploreConditionOperatorType.NotHasAny);
    public static readonly NotHasAll = new TransactionExploreConditionOperator('Not has all', TransactionExploreConditionOperatorType.NotHasAll);
    public static readonly IsEmpty = new TransactionExploreConditionOperator('Is empty', TransactionExploreConditionOperatorType.IsEmpty);
    public static readonly IsNotEmpty = new TransactionExploreConditionOperator('Is not empty', TransactionExploreConditionOperatorType.IsNotEmpty);
    public static readonly Contains = new TransactionExploreConditionOperator('Contains', TransactionExploreConditionOperatorType.Contains);
    public static readonly NotContains = new TransactionExploreConditionOperator('Not contains', TransactionExploreConditionOperatorType.NotContains);
    public static readonly StartsWith = new TransactionExploreConditionOperator('Starts with', TransactionExploreConditionOperatorType.StartsWith);
    public static readonly NotStartsWith = new TransactionExploreConditionOperator('Not starts with', TransactionExploreConditionOperatorType.NotStartsWith);
    public static readonly EndsWith = new TransactionExploreConditionOperator('Ends with', TransactionExploreConditionOperatorType.EndsWith);
    public static readonly NotEndsWith = new TransactionExploreConditionOperator('Not ends with', TransactionExploreConditionOperatorType.NotEndsWith);

    public readonly name: string;
    public readonly value: TransactionExploreConditionOperatorType;

    private constructor(name: string, value: TransactionExploreConditionOperatorType) {
        this.name = name;
        this.value = value;

        TransactionExploreConditionOperator.allInstances.push(this);
        TransactionExploreConditionOperator.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreConditionOperator[] {
        return TransactionExploreConditionOperator.allInstances;
    }

    public static valueOf(value: string): TransactionExploreConditionOperator | undefined {
        return TransactionExploreConditionOperator.allInstancesByValue[value];
    }
}

export enum TransactionExploreChartTypeValue {
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

export class TransactionExploreChartType implements NameValue {
    private static readonly allInstances: TransactionExploreChartType[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreChartType> = {};

    public static readonly Pie = new TransactionExploreChartType('Pie Chart', TransactionExploreChartTypeValue.Pie, false);
    public static readonly Radar = new TransactionExploreChartType('Radar Chart', TransactionExploreChartTypeValue.Radar, false);

    public static readonly Default = TransactionExploreChartType.Pie;

    public readonly name: string;
    public readonly value: TransactionExploreChartTypeValue;
    public readonly seriesDimensionRequired: boolean;

    private constructor(name: string, value: TransactionExploreChartTypeValue, seriesDimensionRequired: boolean) {
        this.name = name;
        this.value = value;
        this.seriesDimensionRequired = seriesDimensionRequired;

        TransactionExploreChartType.allInstances.push(this);
        TransactionExploreChartType.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreChartType[] {
        return TransactionExploreChartType.allInstances;
    }

    public static valueOf(value: string): TransactionExploreChartType | undefined {
        return TransactionExploreChartType.allInstancesByValue[value];
    }
}

export enum TransactionExploreDataDimensionType {
    None = 'none',
    Query = 'query',
    DateTime = 'dateTime',
    DateTimeByDay = 'dateTimeByDay',
    DateTimeByMonth = 'dateTimeByMonth',
    DateTimeByQuarter = 'dateTimeByQuarter',
    DateTimeByYear = 'dateTimeByYear',
    DateTimeByFiscalYear = 'dateTimeByFiscalYear',
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

export class TransactionExploreDataDimension implements NameValue {
    private static readonly allInstances: TransactionExploreDataDimension[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreDataDimension> = {};

    public static readonly None = new TransactionExploreDataDimension('None', TransactionExploreDataDimensionType.None);
    public static readonly Query = new TransactionExploreDataDimension('Query', TransactionExploreDataDimensionType.Query);
    public static readonly DateTime = new TransactionExploreDataDimension('Transaction Time', TransactionExploreDataDimensionType.DateTime);
    public static readonly DateTimeByDay = new TransactionExploreDataDimension('Transaction Date', TransactionExploreDataDimensionType.DateTimeByDay);
    public static readonly DateTimeByMonth = new TransactionExploreDataDimension('Transaction Month', TransactionExploreDataDimensionType.DateTimeByMonth);
    public static readonly DateTimeByQuarter = new TransactionExploreDataDimension('Transaction Quarter', TransactionExploreDataDimensionType.DateTimeByQuarter);
    public static readonly DateTimeByYear = new TransactionExploreDataDimension('Transaction Year', TransactionExploreDataDimensionType.DateTimeByYear);
    public static readonly DateTimeByFiscalYear = new TransactionExploreDataDimension('Transaction Fiscal Year', TransactionExploreDataDimensionType.DateTimeByFiscalYear);
    public static readonly TransactionType = new TransactionExploreDataDimension('Transaction Type', TransactionExploreDataDimensionType.TransactionType);
    public static readonly SourceAccount = new TransactionExploreDataDimension('Source Account', TransactionExploreDataDimensionType.SourceAccount);
    public static readonly SourceAccountCategory = new TransactionExploreDataDimension('Source Account Category', TransactionExploreDataDimensionType.SourceAccountCategory);
    public static readonly SourceAccountCurrency = new TransactionExploreDataDimension('Source Account Currency', TransactionExploreDataDimensionType.SourceAccountCurrency);
    public static readonly DestinationAccount = new TransactionExploreDataDimension('Destination Account', TransactionExploreDataDimensionType.DestinationAccount);
    public static readonly DestinationAccountCategory = new TransactionExploreDataDimension('Destination Account Category', TransactionExploreDataDimensionType.DestinationAccountCategory);
    public static readonly DestinationAccountCurrency = new TransactionExploreDataDimension('Destination Account Currency', TransactionExploreDataDimensionType.DestinationAccountCurrency);
    public static readonly PrimaryCategory = new TransactionExploreDataDimension('Primary Category', TransactionExploreDataDimensionType.PrimaryCategory);
    public static readonly SecondaryCategory = new TransactionExploreDataDimension('Secondary Category', TransactionExploreDataDimensionType.SecondaryCategory);
    public static readonly SourceAmount = new TransactionExploreDataDimension('Amount', TransactionExploreDataDimensionType.SourceAmount);
    public static readonly DestinationAmount = new TransactionExploreDataDimension('Transfer In Amount', TransactionExploreDataDimensionType.DestinationAmount);

    public static readonly CategoryDimensionDefault = TransactionExploreDataDimension.Query;
    public static readonly SeriesDimensionDefault = TransactionExploreDataDimension.None;

    public readonly name: string;
    public readonly value: TransactionExploreDataDimensionType;

    private constructor(name: string, value: TransactionExploreDataDimensionType) {
        this.name = name;
        this.value = value;

        TransactionExploreDataDimension.allInstances.push(this);
        TransactionExploreDataDimension.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreDataDimension[] {
        return TransactionExploreDataDimension.allInstances;
    }

    public static valueOf(value: string): TransactionExploreDataDimension | undefined {
        return TransactionExploreDataDimension.allInstancesByValue[value];
    }
}

export enum TransactionExploreValueMetricType {
    TransactionCount = 'transactionCount',
    SourceAmountSum = 'sourceAmountSum',
    SourceAmountAverage = 'sourceAmountAverage',
    SourceAmountMedian = 'sourceAmountMedian',
    SourceAmountMinimum = 'sourceAmountMinimum',
    SourceAmountMaximum = 'sourceAmountMaximum'
}

export class TransactionExploreValueMetric implements NameValue {
    private static readonly allInstances: TransactionExploreValueMetric[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreValueMetric> = {};

    public static readonly TransactionCount = new TransactionExploreValueMetric('Transaction Count', TransactionExploreValueMetricType.TransactionCount, false);
    public static readonly SourceAmountSum = new TransactionExploreValueMetric('Total Amount', TransactionExploreValueMetricType.SourceAmountSum, true);
    public static readonly SourceAmountAverage = new TransactionExploreValueMetric('Average Amount', TransactionExploreValueMetricType.SourceAmountAverage, true);
    public static readonly SourceAmountMedian = new TransactionExploreValueMetric('Median Amount', TransactionExploreValueMetricType.SourceAmountMedian, true);
    public static readonly SourceAmountMinimum = new TransactionExploreValueMetric('Minimum Amount', TransactionExploreValueMetricType.SourceAmountMinimum, true);
    public static readonly SourceAmountMaximum = new TransactionExploreValueMetric('Maximum Amount', TransactionExploreValueMetricType.SourceAmountMaximum, true);

    public static readonly Default = TransactionExploreValueMetric.SourceAmountSum;

    public readonly name: string;
    public readonly value: TransactionExploreValueMetricType;
    public readonly isAmount: boolean;

    private constructor(name: string, value: TransactionExploreValueMetricType, isAmount: boolean) {
        this.name = name;
        this.value = value;
        this.isAmount = isAmount;

        TransactionExploreValueMetric.allInstances.push(this);
        TransactionExploreValueMetric.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreValueMetric[] {
        return TransactionExploreValueMetric.allInstances;
    }

    public static valueOf(value: string): TransactionExploreValueMetric | undefined {
        return TransactionExploreValueMetric.allInstancesByValue[value];
    }
}

export const DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE: DateRange = DateRange.ThisMonth;
