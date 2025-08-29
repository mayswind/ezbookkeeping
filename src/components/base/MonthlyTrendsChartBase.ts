import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type {
    TextualYearMonth,
    Year1BasedMonth,
    TimeRangeAndDateType,
    YearUnixTime,
    YearQuarterUnixTime,
    YearMonthUnixTime
} from '@/core/datetime.ts';
import type { FiscalYearUnixTime } from '@/core/fiscalyear.ts';
import type { YearMonthItems } from '@/models/transaction.ts';

import { getAllDateRangesFromItems } from '@/lib/statistics.ts';

export interface CommonMonthlyTrendsChartProps<T extends Year1BasedMonth> {
    items: YearMonthItems<T>[];
    startYearMonth: TextualYearMonth | '';
    endYearMonth: TextualYearMonth | '';
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

export interface MonthlyTrendsBarChartClickEvent {
    itemId: string;
    dateRange: TimeRangeAndDateType;
}

export function useMonthlyTrendsChartBase<T extends Year1BasedMonth>(props: CommonMonthlyTrendsChartProps<T>) {
    const { tt } = useI18n();

    const allDateRanges = computed<YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[]>(() => getAllDateRangesFromItems(props.items, props.startYearMonth, props.endYearMonth, props.fiscalYearStart, props.dateAggregationType));

    function getItemName(name: string): string {
        return props.translateName ? tt(name) : name;
    }

    return {
        // computed states
        allDateRanges,
        // functions
        getItemName
    }
}
