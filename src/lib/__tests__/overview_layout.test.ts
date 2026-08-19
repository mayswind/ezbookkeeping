import { describe, expect, test } from 'vitest';

import { OverviewWidgetType, OverviewWidgetDataRequirement } from '@/core/overview_layout.ts';

import { DEFAULT_DESKTOP_OVERVIEW_LAYOUT } from '@/consts/overview_layout.ts';

import {
    serializeDesktopOverviewLayout,
    isDefaultDesktopOverviewLayout,
    getOverviewDataRequirements,
    getOverviewTransactionOverviewMonths,
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
            version: 1,
            widgets: [{ id: 'one', type: OverviewWidgetType.PeriodIncomeExpense, x: -2, y: -1, w: 1, h: 1, settings: { dateRange: 'invalid' } }]
        });

        expect(layout.widgets[0]).toMatchObject({ x: 0, y: 0, w: 2, h: 3, settings: { dateRange: 'today' } });
    });

    test('allows repeated widget types and removes duplicate ids', () => {
        const layout = normalizeDesktopOverviewLayout({
            version: 1,
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
            version: 1,
            widgets: [{ id: 'trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6 }]
        }));
        expect(requirements.includes(OverviewWidgetDataRequirement.TransactionOverview)).toBe(true);
        expect(requirements.includes(OverviewWidgetDataRequirement.TransactionOverviewLast12Months)).toBe(true);
        expect(requirements.includes(OverviewWidgetDataRequirement.Accounts)).toBe(false);
    });

    test('gets the maximum transaction overview months required by widgets', () => {
        const sixMonthLayout = normalizeDesktopOverviewLayout({
            version: 1,
            widgets: [{ id: 'trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6, settings: { months: 6 } }]
        });
        expect(getOverviewTransactionOverviewMonths(sixMonthLayout)).toBe(6);

        const mixedLayout = normalizeDesktopOverviewLayout({
            version: 1,
            widgets: [
                { id: 'six-month-trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 0, y: 0, w: 6, h: 6, settings: { months: 6 } },
                { id: 'twelve-month-trend', type: OverviewWidgetType.IncomeExpenseTrend, x: 6, y: 0, w: 6, h: 6, settings: { months: 12 } }
            ]
        });
        expect(getOverviewTransactionOverviewMonths(mixedLayout)).toBe(12);
    });

    test('round trips serialized layout', () => {
        const json = serializeDesktopOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT, true);
        expect(serializeDesktopOverviewLayout(parseDesktopOverviewLayout(json))).toBe(serializeDesktopOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT));
    });
});
