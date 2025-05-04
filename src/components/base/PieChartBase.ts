import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

import { isNumber } from '@/lib/common.ts';

export interface CommonPieChartDataItem {
    id: string;
    name: string;
    displayName: string;
    value: number;
    percent: number;
    actualPercent: number;
    color: string;
    sourceItem: Record<string, unknown>;
    displayPercent?: string;
    displayValue?: string;
}

export interface CommonPieChartProps {
    skeleton?: boolean;
    items: Record<string, unknown>[];
    idField?: string;
    nameField: string;
    valueField: string;
    percentField?: string;
    colorField?: string;
    hiddenField?: string;
    minValidPercent?: number;
    defaultCurrency?: string;
    showValue?: boolean;
    enableClickItem?: boolean;
}

export function usePieChartBase(props: CommonPieChartProps) {
    const { formatAmountWithCurrency, formatPercent } = useI18n();

    const selectedIndex = ref<number>(0);

    const validItems = computed<CommonPieChartDataItem[]>(() => {
        let totalValidValue = 0;

        for (let i = 0; i < props.items.length; i++) {
            const item = props.items[i];
            const value = item[props.valueField];

            if (isNumber(value) && value > 0 && (!props.hiddenField || !item[props.hiddenField])) {
                totalValidValue += value;
            }
        }

        const validItems: CommonPieChartDataItem[] = [];

        for (let i = 0; i < props.items.length; i++) {
            const item = props.items[i];
            const value = item[props.valueField];
            const percent = props.percentField ? item[props.percentField] : -1;

            if (isNumber(value) && value > 0 &&
                (!props.hiddenField || !item[props.hiddenField]) &&
                (!props.minValidPercent || value / totalValidValue > props.minValidPercent)) {
                const finalItem: CommonPieChartDataItem = {
                    id: (props.idField && item[props.idField]) ? item[props.idField] as string : item[props.nameField] as string,
                    name: (props.idField && item[props.idField]) ? item[props.idField] as string : item[props.nameField] as string,
                    displayName: item[props.nameField] as string,
                    value: value,
                    percent: (isNumber(percent) && percent >= 0) ? percent : (value / totalValidValue * 100),
                    actualPercent: value / totalValidValue,
                    color: (props.colorField && item[props.colorField]) ? item[props.colorField] as string : DEFAULT_CHART_COLORS[validItems.length % DEFAULT_CHART_COLORS.length],
                    sourceItem: item
                };

                finalItem.displayPercent = formatPercent(finalItem.percent, 2, '&lt;0.01');
                finalItem.displayValue = formatAmountWithCurrency(finalItem.value, props.defaultCurrency);

                validItems.push(finalItem);
            }
        }

        return validItems;
    });

    watch(() => props.items, () => {
        selectedIndex.value = 0;
    });

    return {
        // states
        selectedIndex,
        // computed states
        validItems
    };
}
