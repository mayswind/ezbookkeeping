import { describe, expect, test } from 'vitest';

import { OverviewWidgetType, OverviewWidgetDataRequirement } from '@/core/overview_layout.ts';
import { DateRange } from '@/core/datetime.ts';

import { DEFAULT_DESKTOP_OVERVIEW_LAYOUT } from '@/consts/overview_layout.ts';

import {
    serializeDesktopOverviewLayout,
    isDefaultDesktopOverviewLayout,
    getOverviewDataRequirements,
    getOverviewTransactionOverviewMonths,
    getOverviewRecentTransactionCount,
    getOverviewTransactionCategoryStatisticDateTypes,
    getOverviewAssetTrendMonths,
    getOverviewCalendarHeatmapMonths,
    resolveOverviewWidgetCollisions,
    compactOverviewWidgets,
    normalizeDesktopOverviewLayout,
    findOverviewWidgetPosition,
    parseDesktopOverviewLayout
} from '../overview_layout.ts';

describe('desktop overview layout', () => {
    test('empty setting uses default layout', () => {
        const layout = parseDesktopOverviewLayout('');
        expect(isDefaultDesktopOverviewLayout(layout)).toBe(true);
        expect(layout).not.toBe(DEFAULT_DESKTOP_OVERVIEW_LAYOUT);
    });

    test('normalizes sizes and widget settings', () => {
        const layout = normalizeDesktopOverviewLayout({
            widgets: [
                { id: 'one', type: OverviewWidgetType.PeriodIncomeExpense, x: -2, y: -1, w: 1, h: 1, settings: { dateRange: 'invalid' } },
                { id: 'ranking', type: OverviewWidgetType.ExpenseCategoryRanking, x: 0, y: 0, w: 4, h: 6, settings: { dateRange: 'invalid', categoryLevel: 'invalid', itemCount: 100 } },
                { id: 'accounts', type: OverviewWidgetType.AccountBalanceList, x: 4, y: 0, w: 4, h: 6, settings: { accountCategories: [1, 2, 3, 999, 2], itemCount: 10, sortBy: 'balance' } },
                { id: 'income-trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6, settings: { months: 6, showXAxisLabels: false, showLegend: false } },
                { id: 'assets-trend', type: OverviewWidgetType.NetAssetsTrend, x: 6, y: 0, w: 6, h: 6, settings: { months: 12, showXAxisLabels: false, showLegend: false } },
                { id: 'default-trend', type: OverviewWidgetType.NetAssetsTrend, x: 0, y: 6, w: 6, h: 6, settings: { months: 12 } }
            ]
        });

        expect(layout.widgets).toHaveLength(6);
        expect(layout.widgets[0]).toMatchObject({ id: 'income-trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6, settings: { months: 6, showXAxisLabels: false, showLegend: false } });
        expect(layout.widgets[1]).toMatchObject({ id: 'assets-trend', type: OverviewWidgetType.NetAssetsTrend, x: 6, y: 0, w: 6, h: 6, settings: { months: 12, showXAxisLabels: false, showLegend: false } });
        expect(layout.widgets[2]).toMatchObject({ id: 'one', type: OverviewWidgetType.PeriodIncomeExpense, x: 0, y: 6, w: 2, h: 3, settings: { dateRange: DateRange.Today.type } });
        expect(layout.widgets[3]).toMatchObject({ id: 'accounts', type: OverviewWidgetType.AccountBalanceList, x: 4, y: 6, w: 4, h: 6, settings: { accountCategories: [1, 2, 3], itemCount: 10, sortBy: 'balance' } });
        expect(layout.widgets[4]).toMatchObject({ id: 'ranking', type: OverviewWidgetType.ExpenseCategoryRanking, x: 0, y: 9, w: 4, h: 6, settings: { dateRange: DateRange.ThisMonth.type, categoryLevel: 'primary', itemCount: 3 } });
        expect(layout.widgets[5]).toMatchObject({ id: 'default-trend', type: OverviewWidgetType.NetAssetsTrend, x: 0, y: 15, w: 6, h: 6, settings: { months: 12, showXAxisLabels: true, showLegend: true } });
    });

    test('allows repeated widget types and removes duplicate ids', () => {
        const layout = normalizeDesktopOverviewLayout({
            widgets: [
                { id: 'one', type: OverviewWidgetType.PeriodIncomeExpense, x: 0, y: 0, w: 3, h: 3 },
                { id: 'two', type: OverviewWidgetType.PeriodIncomeExpense, x: 3, y: 0, w: 3, h: 3 },
                { id: 'two', type: OverviewWidgetType.AssetSummary, x: 6, y: 0, w: 6, h: 3 }
            ]
        });

        expect(layout.widgets).toHaveLength(2);
        expect(layout.widgets.every(widget => widget.type === OverviewWidgetType.PeriodIncomeExpense)).toBe(true);
    });

    test('pushes colliding widget down and compacts it back up', () => {
        const widgets = [
            { id: 'active', type: OverviewWidgetType.CurrentMonthOverview, x: 0, y: 0, w: 4, h: 3, settings: {} },
            { id: 'other', type: OverviewWidgetType.AssetSummary, x: 3, y: 1, w: 4, h: 3, settings: {} }
        ];
        const resolved = resolveOverviewWidgetCollisions(widgets, 'active');
        expect(resolved.find(widget => widget.id === 'other')?.y).toBe(3);
        expect(compactOverviewWidgets(resolved).find(widget => widget.id === 'other')?.y).toBe(3);
    });

    test('moves another widget into available space above the active widget', () => {
        const widgets = [
            { id: 'active', type: OverviewWidgetType.CurrentMonthOverview, x: 0, y: 3, w: 4, h: 3, settings: {} },
            { id: 'other', type: OverviewWidgetType.AssetSummary, x: 0, y: 3, w: 4, h: 3, settings: {} }
        ];
        const resolved = resolveOverviewWidgetCollisions(widgets, 'active');

        expect(resolved.find(widget => widget.id === 'active')?.y).toBe(3);
        expect(resolved.find(widget => widget.id === 'other')?.y).toBe(6);

        const compacted = compactOverviewWidgets(resolved, 'active');
        expect(compacted.find(widget => widget.id === 'other')?.y).toBe(0);
        expect(compacted.find(widget => widget.id === 'active')?.y).toBe(3);
    });

    test('keeps another widget below when space above the active widget is too small', () => {
        const widgets = [
            { id: 'active', type: OverviewWidgetType.CurrentMonthOverview, x: 0, y: 2, w: 4, h: 3, settings: {} },
            { id: 'other', type: OverviewWidgetType.AssetSummary, x: 0, y: 2, w: 4, h: 3, settings: {} }
        ];
        const resolved = resolveOverviewWidgetCollisions(widgets, 'active');
        const compacted = compactOverviewWidgets(resolved, 'active');

        expect(compacted.find(widget => widget.id === 'active')?.y).toBe(2);
        expect(compacted.find(widget => widget.id === 'other')?.y).toBe(5);
    });

    test('resolves full-height collisions without looping', () => {
        const resolved = resolveOverviewWidgetCollisions([
            { id: 'one', type: OverviewWidgetType.AssetSummary, x: 0, y: 0, w: 12, h: 1000, settings: {} },
            { id: 'two', type: OverviewWidgetType.AssetSummary, x: 0, y: 0, w: 12, h: 1000, settings: {} }
        ]);
        expect(resolved[1]?.y).toBe(1000);
    });

    test('finds first available position', () => {
        const position = findOverviewWidgetPosition([
            { id: 'one', type: OverviewWidgetType.AssetSummary, x: 0, y: 0, w: 12, h: 3, settings: {} }
        ], 4, 3);
        expect(position).toEqual({ x: 0, y: 3 });
    });

    test('merges data requirements', () => {
        const requirements = getOverviewDataRequirements(normalizeDesktopOverviewLayout({
            widgets: [{ id: 'trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6 }]
        }));
        expect(requirements.includes(OverviewWidgetDataRequirement.TransactionOverview)).toBe(true);
        expect(requirements.includes(OverviewWidgetDataRequirement.TransactionOverviewLast12Months)).toBe(true);
        expect(requirements.includes(OverviewWidgetDataRequirement.Accounts)).toBe(false);
    });

    test('gets the maximum transaction overview months required by widgets', () => {
        const sixMonthLayout = normalizeDesktopOverviewLayout({
            widgets: [{ id: 'trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6, settings: { months: 6 } }]
        });
        expect(getOverviewTransactionOverviewMonths(sixMonthLayout)).toBe(6);

        const mixedLayout = normalizeDesktopOverviewLayout({
            widgets: [
                { id: 'six-month-trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6, settings: { months: 6 } },
                { id: 'twelve-month-trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 6, y: 0, w: 6, h: 6, settings: { months: 12 } }
            ]
        });
        expect(getOverviewTransactionOverviewMonths(mixedLayout)).toBe(12);
    });

    test('gets maximum widget data ranges', () => {
        const layout = normalizeDesktopOverviewLayout({
            widgets: [
                { id: 'progress', type: OverviewWidgetType.CurrentMonthExpenseProgress, x: 0, y: 0, w: 4, h: 4 },
                { id: 'recent', type: OverviewWidgetType.RecentTransactions, x: 4, y: 0, w: 4, h: 6, settings: { itemCount: 10 } },
                { id: 'worth', type: OverviewWidgetType.NetAssetsTrend, x: 0, y: 6, w: 6, h: 6, settings: { months: 12 } },
                { id: 'calendar', type: OverviewWidgetType.TransactionCalendarHeatmap, x: 0, y: 12, w: 8, h: 6, settings: { months: 6 } }
            ]
        });

        expect(getOverviewTransactionOverviewMonths(layout)).toBe(2);
        expect(getOverviewRecentTransactionCount(layout)).toBe(10);
        expect(getOverviewAssetTrendMonths(layout)).toBe(12);
        expect(getOverviewCalendarHeatmapMonths(layout)).toBe(6);
    });

    test('gets only the category statistic date types used by ranking widgets', () => {
        const layout = normalizeDesktopOverviewLayout({
            widgets: [
                { id: 'month-ranking', type: OverviewWidgetType.ExpenseCategoryRanking, x: 0, y: 0, w: 3, h: 4, settings: { dateRange: DateRange.ThisMonth.type } },
                { id: 'year-ranking', type: OverviewWidgetType.ExpenseCategoryRanking, x: 3, y: 0, w: 3, h: 4, settings: { dateRange: DateRange.ThisYear.type } },
                { id: 'month-ranking-2', type: OverviewWidgetType.ExpenseCategoryRanking, x: 6, y: 0, w: 3, h: 4, settings: { dateRange: DateRange.ThisMonth.type } }
            ]
        });

        expect(getOverviewTransactionCategoryStatisticDateTypes(layout)).toEqual([DateRange.ThisMonth.type, DateRange.ThisYear.type]);
    });

    test('round trips serialized layout', () => {
        const json = serializeDesktopOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT, true);
        expect(serializeDesktopOverviewLayout(parseDesktopOverviewLayout(json))).toBe(serializeDesktopOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT));
    });
});
