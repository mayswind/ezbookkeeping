import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import { ChartValueType, type CategoricalChartSourceDataItem } from '@/core/chart.ts';

import { isNumber } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, isBigDecimal } from '@/lib/numeral.ts';
import { getDisplayColor } from '@/lib/color.ts';

export interface CommonPieChartDataItem {
    id: string;
    name: string;
    displayName: string;
    value: number; // only used for echarts rendering, the actual value is in originalValue
    originalValue: string;
    displayValue: string;
    percent: number;
    paintPercent: number;
    displayPercent: string;
    color: ColorStyleValue;
    sourceItem: CategoricalChartSourceDataItem;
}

export interface CommonPieChartProps {
    skeleton?: boolean;
    items: CategoricalChartSourceDataItem[];
    valueType: ChartValueType;
    defaultCurrency?: string;
    showValue?: boolean;
    showPercent?: boolean;
    useCustomColor?: boolean;
    enableClickItem?: boolean;
}

export function usePieChartBase(props: CommonPieChartProps) {
    const {
        formatPercentToLocalizedNumerals,
        formatChartValueToLocalizedNumerals
    } = useI18n();

    const settingsStore = useSettingsStore();

    const selectedIndex = ref<number>(0);

    const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);

    const validItems = computed<CommonPieChartDataItem[]>(() => {
        let totalValidValue: BigDecimal = BIG_DECIMAL_ZERO;

        for (const item of props.items) {
            if (isBigDecimal(item.value) && item.value.isPositive() && !item.hidden) {
                totalValidValue = totalValidValue.add(item.value);
            }
        }

        const validItems: CommonPieChartDataItem[] = [];
        let accumulatedPaintPercent: number = 0;

        for (const item of props.items) {
            if (isBigDecimal(item.value) && !item.hidden) {
                let percent: number = isNumber(item.percent) ? item.percent : -1;

                if (percent < 0) {
                    if (item.value.isPositive()) {
                        percent = item.value.divide(totalValidValue).multiply(100).toDoubleNumber();
                    } else {
                        percent = 0;
                    }
                }

                const finalItem: CommonPieChartDataItem = {
                    id: item.id ?? item.name,
                    name: item.id ?? item.name,
                    displayName: item.name,
                    value: item.value.isPositive() ? item.value.toDoubleNumber() : 0,
                    originalValue: item.value.toString(),
                    displayValue: formatChartValueToLocalizedNumerals(item.value, props.valueType, props.defaultCurrency),
                    percent: percent,
                    paintPercent: item.value.isPositive() ? item.value.divide(totalValidValue).toDoubleNumber() : 0,
                    displayPercent: formatPercentToLocalizedNumerals(percent, 2, '<0.01'),
                    color: getDisplayColor(props.useCustomColor && item.color ? item.color : chartColors.value[validItems.length % chartColors.value.length]),
                    sourceItem: item
                };

                accumulatedPaintPercent += finalItem.paintPercent;
                validItems.push(finalItem);
            }
        }

        if (validItems.length > 0) {
            validItems[validItems.length - 1]!.paintPercent += 1 - accumulatedPaintPercent;
        }

        return validItems;
    });

    const allItemsMap = computed<Record<string, CategoricalChartSourceDataItem>>(() => {
        const map: Record<string, CategoricalChartSourceDataItem> = {};

        for (const item of props.items) {
            map[item.id ?? item.name] = item;
        }

        return map;
    });

    watch(() => props.items, () => {
        selectedIndex.value = 0;
    });

    return {
        // states
        selectedIndex,
        // computed states
        validItems,
        allItemsMap
    };
}
