import {computed, ref} from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';

import type { ColorValue } from '@/core/color.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

import { isHextualColor, isEquals } from '@/lib/common.ts';

export function useChartColorSchemeSettingsPageBase() {
    const settingsStore = useSettingsStore();

    const chartColors = ref<ColorValue[]>([...settingsStore.chartColorList]);
    const chartColorsModified = computed<boolean>(() => !isEquals(chartColors.value, settingsStore.chartColorList));
    const textualChartColors = computed<string>({
        get: () => chartColors.value.map(c => `#${c}`).join('\n'),
        set: (value: string) => chartColors.value = filterValidColors(value)
    });

    const canSaveColorScheme = computed<boolean>(() => chartColorsModified.value && chartColors.value.length > 0);

    function filterValidColors(str: string): ColorValue[] {
        return str.split(/[\n,]/).map(c => c.trim().toLowerCase().replace(/^#/, '')).filter(c => isHextualColor(c));
    }

    function addChartColor(color?: ColorValue): void {
        chartColors.value.push(color?.toLowerCase() ?? '000000');
    }

    function removeChartColor(index: number): void {
        chartColors.value.splice(index, 1);
    }

    function updateChartColor(index: number, color: ColorValue): void {
        chartColors.value[index] = color.toLowerCase();
    }

    function resetChartColorsToDefault(): void {
        chartColors.value = [...DEFAULT_CHART_COLORS];
    }

    function loadChartColorsFromSettings(): void {
        chartColors.value = [...settingsStore.chartColorList];
    }

    function saveChartColorsToSettings(): void {
        const value = chartColors.value.join(',');
        const defaultColors = DEFAULT_CHART_COLORS.join(',');

        if (value === defaultColors) {
            settingsStore.setChartColors('');
        } else {
            settingsStore.setChartColors(value);
        }
    }

    function onColorInput(index: number, event: Event): void {
        const target = event.target as HTMLInputElement;
        const hex = target.value.replace(/^#/, '');

        if (isHextualColor(hex)) {
            updateChartColor(index, hex as ColorValue);
        }
    }

    return {
        // states
        chartColors,
        // computed states
        chartColorsModified,
        textualChartColors,
        canSaveColorScheme,
        // functions
        filterValidColors,
        addChartColor,
        removeChartColor,
        updateChartColor,
        resetChartColorsToDefault,
        loadChartColorsFromSettings,
        saveChartColorsToSettings,
        onColorInput
    };
}
