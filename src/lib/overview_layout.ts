import { keys } from '@/core/base.ts';
import {
    type OverviewWidgetSettingValue,
    type DesktopOverviewLayout,
    type DesktopOverviewWidgetLayout,
    OverviewWidgetType,
    OverviewWidgetDataRequirement
} from '@/core/overview_layout.ts';

import {
    DESKTOP_OVERVIEW_LAYOUT_COLUMNS,
    DESKTOP_OVERVIEW_LAYOUT_MAX_WIDGETS,
    DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS,
    DESKTOP_OVERVIEW_WIDGET_DEFINITIONS,
    DEFAULT_DESKTOP_OVERVIEW_LAYOUT
} from '@/consts/overview_layout.ts';

import {
    isObject,
    isArray,
    isString,
    isInteger,
    normalizeInteger
} from '@/lib/common.ts';

function widgetsOverlap(first: DesktopOverviewWidgetLayout, second: DesktopOverviewWidgetLayout): boolean {
    return first.x < second.x + second.w && first.x + first.w > second.x && first.y < second.y + second.h && first.y + first.h > second.y;
}

function normalizeOverviewWidgetSettings(type: OverviewWidgetType, settings: unknown): Record<string, OverviewWidgetSettingValue> {
    const normalized = { ...DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[type].defaultSettings };
    const source: Record<string, unknown> = isObject(settings) ? settings as Record<string, unknown> : {};

    if (type === OverviewWidgetType.PeriodIncomeExpense) {
        const dateRange = source['dateRange'];
        normalized['dateRange'] = dateRange === 'today' || dateRange === 'thisWeek' || dateRange === 'thisMonth' || dateRange === 'thisYear' ? dateRange : 'today';
    } else if (type === OverviewWidgetType.IncomeExpenseTrend) {
        normalized['months'] = source['months'] === 6 ? 6 : 12;
    }

    return normalized;
}

export function serializeDesktopOverviewLayout(layout: DesktopOverviewLayout, pretty?: boolean): string {
    const normalized = normalizeDesktopOverviewLayout(layout);
    return JSON.stringify(normalized, null, pretty ? 2 : undefined);
}

export function isDefaultDesktopOverviewLayout(layout: DesktopOverviewLayout): boolean {
    return serializeDesktopOverviewLayout(layout) === serializeDesktopOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT);
}

export function getOverviewDataRequirements(layout: DesktopOverviewLayout): OverviewWidgetDataRequirement[] {
    const requirements: Record<string, boolean> = {};
    const result: OverviewWidgetDataRequirement[] = [];

    for (const widget of layout.widgets) {
        for (const requirement of DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.type].dataRequirements) {
            requirements[requirement] = true;
        }
    }

    if (requirements[OverviewWidgetDataRequirement.TransactionOverviewLast12Months]) {
        requirements[OverviewWidgetDataRequirement.TransactionOverview] = true;
    }

    for (const requirement of keys(requirements)) {
        if (requirements[requirement]) {
            result.push(requirement as OverviewWidgetDataRequirement);
        }
    }

    return result;
}

export function getOverviewTransactionOverviewMonths(layout: DesktopOverviewLayout): number {
    let months: number = 1;

    for (const widget of layout.widgets) {
        if (widget.type === OverviewWidgetType.IncomeExpenseTrend) {
            const monthsValue = widget.settings['months'];

            if (isInteger(monthsValue)) {
                months = Math.max(months, monthsValue);
            }
        }
    }

    return months;
}

export function resolveOverviewWidgetCollisions(widgets: DesktopOverviewWidgetLayout[], activeId?: string): DesktopOverviewWidgetLayout[] {
    const sorted: DesktopOverviewWidgetLayout[] = widgets.map(cloneWidget).sort((a, b) => {
        if (a.id === activeId) return -1;
        if (b.id === activeId) return 1;
        return a.y - b.y || a.x - b.x || a.id.localeCompare(b.id);
    });
    const placed: DesktopOverviewWidgetLayout[] = [];

    for (const widget of sorted) {
        let overlapping: DesktopOverviewWidgetLayout[] = placed.filter(other => widgetsOverlap(widget, other));

        while (overlapping.length > 0) {
            widget.y = Math.max(...overlapping.map(other => other.y + other.h));
            overlapping = placed.filter(other => widgetsOverlap(widget, other));
        }

        placed.push(widget);
    }

    return placed;
}

export function compactOverviewWidgets(widgets: DesktopOverviewWidgetLayout[]): DesktopOverviewWidgetLayout[] {
    const result: DesktopOverviewWidgetLayout[] = widgets.map(cloneWidget).sort((a, b) => a.y - b.y || a.x - b.x || a.id.localeCompare(b.id));

    for (const widget of result) {
        while (widget.y > 0) {
            const candidate = { ...widget, y: widget.y - 1 };

            if (result.some(other => other.id !== widget.id && widgetsOverlap(candidate, other))) {
                break;
            }

            widget.y--;
        }
    }

    return result.sort((a, b) => a.y - b.y || a.x - b.x || a.id.localeCompare(b.id));
}

export function normalizeDesktopOverviewLayout(input: unknown): DesktopOverviewLayout {
    if (!isObject(input)) {
        throw new Error('input is not an object');
    }

    const source = input as Record<string, unknown>;
    const sourceWidgets = source['widgets'];

    if (!isArray(sourceWidgets)) {
        throw new Error('widgets is not an array');
    }

    const finalWidgets: DesktopOverviewWidgetLayout[] = [];
    const existsIds: Record<string, boolean> = {};
    let widgetCount: number = 0;

    for (const item of sourceWidgets) {
        if (!isObject(item)) {
            continue;
        }

        if (widgetCount >= DESKTOP_OVERVIEW_LAYOUT_MAX_WIDGETS) {
            break;
        }

        const sourceWidget = item as Record<string, unknown>;
        const type = sourceWidget['type'] as OverviewWidgetType;
        const definition = DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[type];
        const widgetId = sourceWidget['id'];

        if (!definition || !widgetId || !isString(widgetId) || widgetId.length > 100 || existsIds[widgetId]) {
            continue;
        }

        const w: number = normalizeInteger(sourceWidget['w'], definition.defaultWidth, definition.minWidth, Math.min(DESKTOP_OVERVIEW_LAYOUT_COLUMNS, definition.maxWidth ?? DESKTOP_OVERVIEW_LAYOUT_COLUMNS));
        const h: number = normalizeInteger(sourceWidget['h'], definition.defaultHeight, definition.minHeight, Math.min(DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS, definition.maxHeight ?? DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS));
        const x: number = normalizeInteger(sourceWidget['x'], 0, 0, DESKTOP_OVERVIEW_LAYOUT_COLUMNS - w);
        const y: number = normalizeInteger(sourceWidget['y'], 0, 0, DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS - h);

        const finalWidget: DesktopOverviewWidgetLayout = {
            id: widgetId,
            type: type,
            x: x,
            y: y,
            w: w,
            h: h,
            settings: normalizeOverviewWidgetSettings(type, sourceWidget['settings'])
        };

        existsIds[widgetId] = true;
        finalWidgets.push(finalWidget);
        widgetCount++;
    }

    return {
        widgets: compactOverviewWidgets(resolveOverviewWidgetCollisions(finalWidgets))
    };
}

export function findOverviewWidgetPosition(widgets: DesktopOverviewWidgetLayout[], width: number, height: number): { x: number; y: number } {
    for (let y = 0; y < DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS - height; y++) {
        for (let x = 0; x <= DESKTOP_OVERVIEW_LAYOUT_COLUMNS - width; x++) {
            const candidate: DesktopOverviewWidgetLayout = {
                id: '',
                type: OverviewWidgetType.CurrentMonthOverview,
                x,
                y,
                w: width,
                h: height,
                settings: {}
            };

            if (!widgets.some(widget => widgetsOverlap(candidate, widget))) {
                return { x, y };
            }
        }
    }

    return { x: 0, y: Math.max(0, ...widgets.map(widget => widget.y + widget.h)) };
}

export function cloneWidget(widget: DesktopOverviewWidgetLayout): DesktopOverviewWidgetLayout {
    return { ...widget, settings: { ...widget.settings } };
}

export function cloneOverviewLayout(original: DesktopOverviewLayout): DesktopOverviewLayout {
    return {
        widgets: original.widgets.map(cloneWidget)
    };
}

export function parseDesktopOverviewLayout(value: string): DesktopOverviewLayout {
    if (!value) {
        return cloneOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT);
    }

    return normalizeDesktopOverviewLayout(JSON.parse(value));
}
