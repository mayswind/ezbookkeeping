import type { TypeAndName } from './base.ts';
import { DateRange } from '@/core/datetime.ts';

export enum StatisticsAnalysisType {
    CategoricalAnalysis = 0,
    TrendAnalysis = 1
}

type CategoricalChartTypeName = 'Pie' | 'Bar';

export class CategoricalChartType implements TypeAndName {
    private static readonly allInstances: CategoricalChartType[] = [];
    private static readonly allInstancesByTypeName: Record<string, CategoricalChartType> = {};

    public static readonly Pie = new CategoricalChartType(0, 'Pie', 'Pie Chart');
    public static readonly Bar = new CategoricalChartType(1, 'Bar', 'Bar Chart');

    public static readonly Default = CategoricalChartType.Pie;

    public readonly type: number;
    public readonly typeName: CategoricalChartTypeName;
    public readonly name: string;

    private constructor(type: number, typeName: CategoricalChartTypeName, name: string) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;

        CategoricalChartType.allInstances.push(this);
        CategoricalChartType.allInstancesByTypeName[typeName] = this;
    }

    public static values(): CategoricalChartType[] {
        return CategoricalChartType.allInstances;
    }

    public static all(): Record<CategoricalChartTypeName, CategoricalChartType> {
        return CategoricalChartType.allInstancesByTypeName;
    }
}

type TrendChartTypeName = 'Area' | 'Column';

export class TrendChartType implements TypeAndName {
    private static readonly allInstances: TrendChartType[] = [];
    private static readonly allInstancesByTypeName: Record<string, TrendChartType> = {};

    public static readonly Area = new TrendChartType(0, 'Area', 'Area Chart');
    public static readonly Column = new TrendChartType(1, 'Column', 'Column Chart');

    public static readonly Default = TrendChartType.Column;

    public readonly type: number;
    public readonly typeName: TrendChartTypeName;
    public readonly name: string;

    private constructor(type: number, typeName: TrendChartTypeName, name: string) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;

        TrendChartType.allInstances.push(this);
        TrendChartType.allInstancesByTypeName[typeName] = this;
    }

    public static values(): TrendChartType[] {
        return TrendChartType.allInstances;
    }

    public static all(): Record<TrendChartTypeName, TrendChartType> {
        return TrendChartType.allInstancesByTypeName;
    }
}

type ChartDataTypeName = 'ExpenseByAccount' | 'ExpenseByPrimaryCategory' | 'ExpenseBySecondaryCategory' | 'IncomeByAccount' | 'IncomeByPrimaryCategory' | 'IncomeBySecondaryCategory' | 'AccountTotalAssets' | 'AccountTotalLiabilities' | 'TotalExpense' | 'TotalIncome' | 'TotalBalance';

export class ChartDataType implements TypeAndName {
    private static readonly allInstances: ChartDataType[] = [];
    private static readonly allInstancesByType: Record<number, ChartDataType> = {};
    private static readonly allInstancesByTypeName: Record<string, ChartDataType> = {};

    public static readonly ExpenseByAccount = new ChartDataType(0, 'ExpenseByAccount', 'Expense By Account', StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly ExpenseByPrimaryCategory = new ChartDataType(1, 'ExpenseByPrimaryCategory', 'Expense By Primary Category', StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly ExpenseBySecondaryCategory = new ChartDataType(2, 'ExpenseBySecondaryCategory', 'Expense By Secondary Category', StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly IncomeByAccount = new ChartDataType(3, 'IncomeByAccount', 'Income By Account', StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly IncomeByPrimaryCategory = new ChartDataType(4, 'IncomeByPrimaryCategory', 'Income By Primary Category', StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly IncomeBySecondaryCategory = new ChartDataType(5, 'IncomeBySecondaryCategory', 'Income By Secondary Category', StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly AccountTotalAssets = new ChartDataType(6, 'AccountTotalAssets', 'Account Total Assets', StatisticsAnalysisType.CategoricalAnalysis);
    public static readonly AccountTotalLiabilities = new ChartDataType(7, 'AccountTotalLiabilities', 'Account Total Liabilities', StatisticsAnalysisType.CategoricalAnalysis);
    public static readonly TotalExpense = new ChartDataType(8, 'TotalExpense', 'Total Expense', StatisticsAnalysisType.TrendAnalysis);
    public static readonly TotalIncome = new ChartDataType(9, 'TotalIncome', 'Total Income', StatisticsAnalysisType.TrendAnalysis);
    public static readonly TotalBalance = new ChartDataType(10, 'TotalBalance', 'Total Balance', StatisticsAnalysisType.TrendAnalysis);

    public static readonly Default = ChartDataType.ExpenseByPrimaryCategory;

    public readonly type: number;
    public readonly typeName: ChartDataTypeName;
    public readonly name: string;
    private readonly availableAnalysisTypes: Record<number, boolean>;

    private constructor(type: number, typeName: ChartDataTypeName, name: string, ...availableAnalysisTypes: StatisticsAnalysisType[]) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;
        this.availableAnalysisTypes = {};

        if (availableAnalysisTypes) {
            for (const analysisType of availableAnalysisTypes) {
                this.availableAnalysisTypes[analysisType] = true;
            }
        }

        ChartDataType.allInstances.push(this);
        ChartDataType.allInstancesByType[type] = this;
        ChartDataType.allInstancesByTypeName[typeName] = this;
    }

    public isAvailableAnalysisType(analysisType: StatisticsAnalysisType): boolean {
        return this.availableAnalysisTypes[analysisType] || false;
    }

    public static values(analysisType: StatisticsAnalysisType | undefined): ChartDataType[] {
        if (analysisType === undefined) {
            return ChartDataType.allInstances;
        }

        const ret: ChartDataType[] = [];

        for (const chartDataType of ChartDataType.allInstances) {
            if (chartDataType.isAvailableAnalysisType(analysisType)) {
                ret.push(chartDataType);
            }
        }

        return ret;
    }

    public static all(): Record<ChartDataTypeName, ChartDataType> {
        return ChartDataType.allInstancesByTypeName;
    }

    public static valueOf(type: number): ChartDataType | undefined {
        return ChartDataType.allInstancesByType[type];
    }

    public static isAvailableForAnalysisType(type: number, analysisType: StatisticsAnalysisType): boolean {
        const chartDataType = ChartDataType.allInstancesByType[type];
        return chartDataType?.isAvailableAnalysisType(analysisType) || false;
    }
}

export class ChartSortingType implements TypeAndName {
    private static readonly allInstances: ChartSortingType[] = [];
    private static readonly allInstancesByType: Record<number, ChartSortingType> = {};

    public static readonly Amount = new ChartSortingType(0, 'Amount', 'Sort by Amount');
    public static readonly DisplayOrder = new ChartSortingType(1, 'Display Order', 'Sort by Display Order');
    public static readonly Name = new ChartSortingType(2, 'Name', 'Sort by Name');

    public static readonly Default = ChartSortingType.Amount;

    public readonly type: number;
    public readonly name: string;
    public readonly fullName: string;

    private constructor(type: number, name: string, fullName: string) {
        this.type = type;
        this.name = name;
        this.fullName = fullName;

        ChartSortingType.allInstances.push(this);
        ChartSortingType.allInstancesByType[type] = this;
    }

    public static values(): ChartSortingType[] {
        return ChartSortingType.allInstances;
    }

    public static valueOf(type: number): ChartSortingType | undefined {
        return ChartSortingType.allInstancesByType[type];
    }
}

export class ChartDateAggregationType implements TypeAndName {
    private static readonly allInstances: ChartDateAggregationType[] = [];
    private static readonly allInstancesByType: Record<number, ChartDateAggregationType> = {};

    public static readonly Month = new ChartDateAggregationType(0, 'Aggregate by Month');
    public static readonly Quarter = new ChartDateAggregationType(1, 'Aggregate by Quarter');
    public static readonly Year = new ChartDateAggregationType(2, 'Aggregate by Year');

    public static readonly Default = ChartDateAggregationType.Month;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        ChartDateAggregationType.allInstances.push(this);
        ChartDateAggregationType.allInstancesByType[type] = this;
    }

    public static values(): ChartDateAggregationType[] {
        return ChartDateAggregationType.allInstances;
    }

    public static valueOf(type: number): ChartDateAggregationType | undefined {
        return ChartDateAggregationType.allInstancesByType[type];
    }
}

export const DEFAULT_CATEGORICAL_CHART_DATA_RANGE: DateRange = DateRange.ThisMonth;
export const DEFAULT_TREND_CHART_DATA_RANGE: DateRange = DateRange.ThisYear;
