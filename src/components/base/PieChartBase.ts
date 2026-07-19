import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import type { ColorValue, ColorStyleValue } from '@/core/color.ts';

import { isNumber } from '@/lib/common.ts';
import { getDisplayColor } from '@/lib/color.ts';

export interface CommonPieChartDataItem {
    id: string;
    name: string;
    displayName: string;
    value: number;
    actualValue: number;
    percent: number;
    paintPercent: number;
    color: ColorStyleValue;
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
    amountValue?: boolean;
    percentValue?: boolean;
    defaultCurrency?: string;
    showValue?: boolean;
    showPercent?: boolean;
    enableClickItem?: boolean;
}

export function usePieChartBase(props: CommonPieChartProps) {
    const {
        formatAmountToLocalizedNumeralsWithCurrency,
        formatNumberToLocalizedNumerals,
        formatPercentToLocalizedNumerals
    } = useI18n();

    const settingsStore = useSettingsStore();

    const selectedIndex = ref<number>(0);

    const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);

    const validItems = computed<CommonPieChartDataItem[]>(() => {
        let totalValidValue = 0;

        for (const item of props.items) {
            const value = item[props.valueField];

            if (isNumber(value) && value > 0 && (!props.hiddenField || !item[props.hiddenField])) {
                totalValidValue += value;
            }
        }

        const validItems: CommonPieChartDataItem[] = [];
        let accumulatedPaintPercent: number = 0;

        for (const item of props.items) {
            const value = item[props.valueField];
            const percent = props.percentField ? item[props.percentField] : -1;

            if (isNumber(value) &&
                (!props.hiddenField || !item[props.hiddenField])) {
                const finalItem: CommonPieChartDataItem = {
                    id: (props.idField && item[props.idField]) ? item[props.idField] as string : item[props.nameField] as string,
                    name: (props.idField && item[props.idField]) ? item[props.idField] as string : item[props.nameField] as string,
                    displayName: item[props.nameField] as string,
                    value: value > 0 ? value : 0,
                    actualValue: value,
                    percent: (isNumber(percent) && percent >= 0) ? percent : (value > 0 ? value / totalValidValue * 100 : 0),
                    paintPercent: value > 0 ? value / totalValidValue : 0,
                    color: getDisplayColor((props.colorField && item[props.colorField]) ? item[props.colorField] as ColorValue : chartColors.value[validItems.length % chartColors.value.length]),
                    sourceItem: item
                };

                accumulatedPaintPercent += finalItem.paintPercent;
                finalItem.displayPercent = formatPercentToLocalizedNumerals(finalItem.percent, 2, '<0.01');
                finalItem.displayValue = getDisplayValue(value);

                validItems.push(finalItem);
            }
        }

        if (validItems.length > 0) {
            validItems[validItems.length - 1]!.paintPercent += 1 - accumulatedPaintPercent;
        }

        return validItems;
    });

    function getDisplayValue(value: number): string {
        if (props.percentValue) {
            return formatPercentToLocalizedNumerals(value, 2, '<0.01');
        }

        if (props.amountValue) {
            return formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
        }

        return formatNumberToLocalizedNumerals(value, 4);
    }

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
