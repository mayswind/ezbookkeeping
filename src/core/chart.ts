import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorValue } from '@/core/color.ts';

export enum ChartValueType {
    Amount = 'amount',
    Percent = 'percent',
    Number = 'number'
}

export interface CategoricalChartSourceDataItem {
    id?: string;
    name: string;
    value: BigDecimal;
    percent?: number;
    hidden?: boolean;
    color?: ColorValue;
}

export interface AxisChartSourceDataItem {
    id?: string;
    name: string;
    values: BigDecimal[];
    displayOrders: number[];
    hidden?: boolean;
    color?: ColorValue;
}

export interface TrendsChartSourceDataItem {
    id?: string;
    name: string;
    items: { value: BigDecimal }[];
    displayOrders: number[];
    percent?: number;
    hidden?: boolean;
    color?: ColorValue;
}

export interface CalendarChartSourceDataItem {
    id: string;
    value: BigDecimal;
    hidden?: boolean;
}
