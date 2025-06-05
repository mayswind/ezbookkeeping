import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type {
    YearMonth,
    TimeRangeAndDateType,
    YearUnixTime,
    YearQuarterUnixTime,
    YearMonthUnixTime
} from '@/core/datetime.ts';
import type { FiscalYearUnixTime } from '@/core/fiscalyear.ts';
import type { ColorValue } from '@/core/color.ts';
import { DEFAULT_ICON_COLOR } from '@/consts/color.ts';
import type { YearMonthItems } from '@/models/transaction.ts';

import { getAllDateRanges } from '@/lib/statistics.ts';

export interface CommonTrendsChartProps<T extends YearMonth> {
    items: YearMonthItems<T>[];
    startYearMonth: string;
    endYearMonth: string;
    fiscalYearStart: number;
    sortingType: number;
    dateAggregationType: number;
    idField?: string;
    nameField: string;
    valueField: string;
    colorField?: string;
    hiddenField?: string;
    displayOrdersField?: string;
    translateName?: boolean;
    defaultCurrency?: string;
    enableClickItem?: boolean;
}

export interface TrendsBarChartClickEvent {
    itemId: string;
    dateRange: TimeRangeAndDateType;
}

export function useTrendsChartBase<T extends YearMonth>(props: CommonTrendsChartProps<T>) {
    const { tt } = useI18n();

    const allDateRanges = computed<YearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] | FiscalYearUnixTime[]>(() => getAllDateRanges(props.items, props.startYearMonth, props.endYearMonth, props.fiscalYearStart, props.dateAggregationType));

    function getItemName(name: string): string {
        return props.translateName ? tt(name) : name;
    }

    function getColor(color: string): ColorValue {
        if (color && color !== DEFAULT_ICON_COLOR) {
            color = '#' + color;
        }

        return color;
    }

    return {
        // computed states
        allDateRanges,
        // functions
        getItemName,
        getColor
    }
}
