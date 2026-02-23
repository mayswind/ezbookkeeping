import type { TypeAndName, TypeAndNameWithAlternativeName } from './base.ts';
import { DateRange } from '@/core/datetime.ts';

export enum StatisticsAnalysisType {
    CategoricalAnalysis = 0,
    TrendAnalysis = 1,
    AssetTrends = 2
}

export enum ChartDataAggregationType {
    Sum = 0,
    Last = 1
}

export class CategoricalChartType implements TypeAndName {
    private static readonly allInstancesForAll: CategoricalChartType[] = [];
    private static readonly allInstancesForDesktop: CategoricalChartType[] = [];
    private static readonly allInstancesByType: Record<number, CategoricalChartType> = {};

    public static readonly Pie = new CategoricalChartType(0, 'Pie Chart', false);
    public static readonly Bar = new CategoricalChartType(1, 'Bar Chart', false);
    public static readonly Radar = new CategoricalChartType(2, 'Radar Chart', true);

    public static readonly Default = CategoricalChartType.Pie;

    public readonly type: number;
    public readonly name: string;
    public readonly desktopOnly: boolean = false;

    private constructor(type: number, name: string, desktopOnly: boolean) {
        this.type = type;
        this.name = name;
        this.desktopOnly = desktopOnly;

        if (!desktopOnly) {
            CategoricalChartType.allInstancesForAll.push(this);
        }

        CategoricalChartType.allInstancesForDesktop.push(this);
        CategoricalChartType.allInstancesByType[type] = this;
    }

    public static isValidType(type: number): boolean {
        return !!CategoricalChartType.allInstancesByType[type];
    }

    public static values(withDesktopOnlyChart: boolean): CategoricalChartType[] {
        if (withDesktopOnlyChart) {
            return CategoricalChartType.allInstancesForDesktop;
        } else {
            return CategoricalChartType.allInstancesForAll;
        }
    }
}

export class TrendChartType implements TypeAndName {
    private static readonly allInstances: TrendChartType[] = [];
    private static readonly allInstancesByType: Record<number, TrendChartType> = {};

    public static readonly Area = new TrendChartType(0, 'Area Chart');
    public static readonly Column = new TrendChartType(1, 'Column Chart');
    public static readonly Bubble = new TrendChartType(2, 'Bubble Chart');

    public static readonly Default = TrendChartType.Column;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        TrendChartType.allInstances.push(this);
        TrendChartType.allInstancesByType[type] = this;
    }

    public static isValidType(type: number): boolean {
        return !!TrendChartType.allInstancesByType[type];
    }

    public static values(): TrendChartType[] {
        return TrendChartType.allInstances;
    }
}

export class AccountBalanceTrendChartType implements TypeAndName {
    private static readonly allInstances: AccountBalanceTrendChartType[] = [];

    public static readonly Area = new AccountBalanceTrendChartType(0, 'Area Chart');
    public static readonly Column = new AccountBalanceTrendChartType(1, 'Column Chart');
    public static readonly Candlestick = new AccountBalanceTrendChartType(2, 'Candlestick Chart');

    public static readonly Default = TrendChartType.Column;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        AccountBalanceTrendChartType.allInstances.push(this);
    }

    public static values(): TrendChartType[] {
        return AccountBalanceTrendChartType.allInstances;
    }
}

export class ChartDataType implements TypeAndName {
    private static readonly allInstancesForAll: ChartDataType[] = [];
    private static readonly allInstancesForDesktop: ChartDataType[] = [];
    private static readonly allInstancesByType: Record<number, ChartDataType> = {};

    public static readonly Overview = new ChartDataType(16, 'Overview', true, true, StatisticsAnalysisType.CategoricalAnalysis);
    public static readonly OutflowsByAccount = new ChartDataType(11, 'Outflows By Account', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly ExpenseByAccount = new ChartDataType(0, 'Expense By Account', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly ExpenseByPrimaryCategory = new ChartDataType(1, 'Expense By Primary Category', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly ExpenseBySecondaryCategory = new ChartDataType(2, 'Expense By Secondary Category', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly InflowsByAccount = new ChartDataType(12, 'Inflows By Account', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly IncomeByAccount = new ChartDataType(3, 'Income By Account', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly IncomeByPrimaryCategory = new ChartDataType(4, 'Income By Primary Category', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly IncomeBySecondaryCategory = new ChartDataType(5, 'Income By Secondary Category', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.TrendAnalysis);
    public static readonly AccountTotalAssets = new ChartDataType(6, 'Account Total Assets', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.AssetTrends);
    public static readonly AccountTotalLiabilities = new ChartDataType(7, 'Account Total Liabilities', false, false, StatisticsAnalysisType.CategoricalAnalysis, StatisticsAnalysisType.AssetTrends);
    public static readonly TotalOutflows = new ChartDataType(13, 'Total Outflows', false, false, StatisticsAnalysisType.TrendAnalysis);
    public static readonly TotalExpense = new ChartDataType(8, 'Total Expense', false, false, StatisticsAnalysisType.TrendAnalysis);
    public static readonly TotalInflows = new ChartDataType(14, 'Total Inflows', false, false, StatisticsAnalysisType.TrendAnalysis);
    public static readonly TotalIncome = new ChartDataType(9, 'Total Income', false, false, StatisticsAnalysisType.TrendAnalysis);
    public static readonly NetCashFlow = new ChartDataType(15, 'Net Cash Flow', false, false, StatisticsAnalysisType.TrendAnalysis);
    public static readonly NetIncome = new ChartDataType(10, 'Net Income', false, false, StatisticsAnalysisType.TrendAnalysis);
    public static readonly NetWorth = new ChartDataType(17, 'Net Worth', false, false, StatisticsAnalysisType.AssetTrends);

    public static readonly Default = ChartDataType.ExpenseByPrimaryCategory;
    public static readonly DefaultForAssetTrends = ChartDataType.NetWorth;

    public readonly type: number;
    public readonly name: string;
    public readonly desktopOnly: boolean = false;
    public readonly specialChart: boolean = false;
    private readonly availableAnalysisTypes: Record<number, boolean>;

    private constructor(type: number, name: string, desktopOnly: boolean, specialChart: boolean, ...availableAnalysisTypes: StatisticsAnalysisType[]) {
        this.type = type;
        this.name = name;
        this.desktopOnly = desktopOnly;
        this.specialChart = specialChart;
        this.availableAnalysisTypes = {};

        if (availableAnalysisTypes) {
            for (const analysisType of availableAnalysisTypes) {
                this.availableAnalysisTypes[analysisType] = true;
            }
        }

        if (!desktopOnly) {
            ChartDataType.allInstancesForAll.push(this);
        }

        ChartDataType.allInstancesForDesktop.push(this);
        ChartDataType.allInstancesByType[type] = this;
    }

    public isAvailableAnalysisType(analysisType: StatisticsAnalysisType): boolean {
        return this.availableAnalysisTypes[analysisType] || false;
    }

    public static values(analysisType?: StatisticsAnalysisType, withDesktopOnlyChart?: boolean): ChartDataType[] {
        const availableInstances: ChartDataType[] = withDesktopOnlyChart ? ChartDataType.allInstancesForDesktop : ChartDataType.allInstancesForAll;

        if (analysisType === undefined) {
            return availableInstances;
        }

        const ret: ChartDataType[] = [];

        for (const chartDataType of availableInstances) {
            if (chartDataType.isAvailableAnalysisType(analysisType)) {
                ret.push(chartDataType);
            }
        }

        return ret;
    }

    public static valueOf(type: number): ChartDataType | undefined {
        return ChartDataType.allInstancesByType[type];
    }

    public static isAvailableForAnalysisType(type: number, analysisType: StatisticsAnalysisType): boolean {
        const chartDataType = ChartDataType.allInstancesByType[type];
        return chartDataType?.isAvailableAnalysisType(analysisType) || false;
    }
}

export class ChartSortingType implements TypeAndNameWithAlternativeName {
    private static readonly allInstances: ChartSortingType[] = [];
    private static readonly allInstancesByType: Record<number, ChartSortingType> = {};

    public static readonly Amount = new ChartSortingType(0, 'Amount', 'Value');
    public static readonly DisplayOrder = new ChartSortingType(1, 'Display Order');
    public static readonly Name = new ChartSortingType(2, 'Name');

    public static readonly Default = ChartSortingType.Amount;

    public readonly type: number;
    public readonly name: string;
    public readonly alternativeName?: string;

    private constructor(type: number, name: string, alternativeName?: string) {
        this.type = type;
        this.name = name;
        this.alternativeName = alternativeName;

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

export class ChartDateAggregationType {
    private static readonly allInstances: ChartDateAggregationType[] = [];
    private static readonly allInstancesByType: Record<number, ChartDateAggregationType> = {};

    public static readonly Day = new ChartDateAggregationType(4, 'Daily', 'Aggregate by Day', StatisticsAnalysisType.AssetTrends);
    public static readonly Month = new ChartDateAggregationType(0, 'Monthly', 'Aggregate by Month', StatisticsAnalysisType.TrendAnalysis, StatisticsAnalysisType.AssetTrends);
    public static readonly Quarter = new ChartDateAggregationType(1, 'Quarterly', 'Aggregate by Quarter', StatisticsAnalysisType.TrendAnalysis, StatisticsAnalysisType.AssetTrends);
    public static readonly Year = new ChartDateAggregationType(2, 'Yearly', 'Aggregate by Year', StatisticsAnalysisType.TrendAnalysis, StatisticsAnalysisType.AssetTrends);
    public static readonly FiscalYear = new ChartDateAggregationType(3, 'FiscalYearly', 'Aggregate by Fiscal Year', StatisticsAnalysisType.TrendAnalysis, StatisticsAnalysisType.AssetTrends);

    public static readonly Default = ChartDateAggregationType.Month;

    public readonly type: number;
    public readonly shortName: string;
    public readonly fullName: string;
    private readonly availableAnalysisTypes: Record<number, boolean>;

    private constructor(type: number, shortName: string, fullName: string, ...availableAnalysisTypes: StatisticsAnalysisType[]) {
        this.type = type;
        this.shortName = shortName;
        this.fullName = fullName;
        this.availableAnalysisTypes = {};

        if (availableAnalysisTypes) {
            for (const analysisType of availableAnalysisTypes) {
                this.availableAnalysisTypes[analysisType] = true;
            }
        }

        ChartDateAggregationType.allInstances.push(this);
        ChartDateAggregationType.allInstancesByType[type] = this;
    }

    public isAvailableAnalysisType(analysisType: StatisticsAnalysisType): boolean {
        return this.availableAnalysisTypes[analysisType] || false;
    }

    public static values(analysisType?: StatisticsAnalysisType): ChartDateAggregationType[] {
        const availableInstances: ChartDateAggregationType[] = ChartDateAggregationType.allInstances;

        if (analysisType === undefined) {
            return availableInstances;
        }

        const ret: ChartDateAggregationType[] = [];

        for (const chartDataType of availableInstances) {
            if (chartDataType.isAvailableAnalysisType(analysisType)) {
                ret.push(chartDataType);
            }
        }

        return ret;
    }

    public static valueOf(type: number): ChartDateAggregationType | undefined {
        return ChartDateAggregationType.allInstancesByType[type];
    }
}

export const DEFAULT_CATEGORICAL_CHART_DATA_RANGE: DateRange = DateRange.ThisMonth;
export const DEFAULT_TREND_CHART_DATA_RANGE: DateRange = DateRange.ThisYear;
export const DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE: DateRange = DateRange.ThisYear;
