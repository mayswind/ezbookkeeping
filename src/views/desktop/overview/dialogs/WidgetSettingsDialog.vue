<template>
    <v-dialog width="600" v-model="showState">
        <one-column-dialog-layout :title="tt('Widget Settings')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined" @click="confirm">{{ tt('Apply') }}</v-btn>
            </template>

            <template #content>
                <div class="mt-1">
                    <v-select item-title="name" item-value="value" :label="tt('Date Range')" :items="allPeriodOptions"
                              v-model="settingValue"
                              v-if="widget?.type === OverviewWidgetType.PeriodIncomeExpense" />

                    <v-select item-title="name" item-value="value" :label="tt('Date Range')" :items="allMonthOptions"
                              v-model="settingValue"
                              v-else-if="widget?.type === OverviewWidgetType.IncomeExpenseTrend" />
                </div>

                <div class="text-body-large text-medium-emphasis text-center" v-if="widget && !hasSettings">{{ tt('This widget has no configurable parameters') }}</div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { NameValue, NameNumeralValue } from '@/core/base.ts';
import {
    type OverviewWidgetSettingValue,
    type DesktopOverviewWidgetLayout,
    OverviewWidgetType
} from '@/core/overview_layout.ts';

import { isDefined } from '@/lib/common.ts';
import { cloneWidget } from '@/lib/overview_layout.ts';

const { tt, formatNumberToLocalizedNumerals } = useI18n();

let resolveFunc: ((widget: DesktopOverviewWidgetLayout) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const widget = ref<DesktopOverviewWidgetLayout | null>(null);

const allPeriodOptions = computed<NameValue[]>(() => [
    {
        name: tt('Today'),
        value: 'today'
    },
    {
        name: tt('This Week'),
        value: 'thisWeek'
    },
    {
        name: tt('This Month'),
        value: 'thisMonth'
    },
    {
        name: tt('This Year'),
        value: 'thisYear'
    }
]);

const allMonthOptions = computed<NameNumeralValue[]>(() => [
    {
        name: tt('format.misc.nMonths', { n: formatNumberToLocalizedNumerals(6) }),
        value: 6
    },
    {
        name: tt('format.misc.nMonths', { n: formatNumberToLocalizedNumerals(12) }),
        value: 12
    }
]);


const settingValue = computed<OverviewWidgetSettingValue | undefined>({
    get: () => {
        if (widget.value?.type === OverviewWidgetType.PeriodIncomeExpense) {
            return widget.value?.settings['dateRange'];
        } else if (widget.value?.type === OverviewWidgetType.IncomeExpenseTrend) {
            return widget.value?.settings['months'];
        } else {
            return undefined;
        }
    },
    set: (value: OverviewWidgetSettingValue | undefined) => {
        if (isDefined(value) && widget.value?.type === OverviewWidgetType.PeriodIncomeExpense) {
            updateWidgetSettings('dateRange', value);
        } else if (isDefined(value) && widget.value?.type === OverviewWidgetType.IncomeExpenseTrend) {
            updateWidgetSettings('months', value);
        }
    }
});

const hasSettings = computed<boolean>(() => widget.value?.type === OverviewWidgetType.PeriodIncomeExpense || widget.value?.type === OverviewWidgetType.IncomeExpenseTrend);

function updateWidgetSettings(settingKey: string, settingValue: OverviewWidgetSettingValue): void {
    if (!widget.value) {
        return;
    }

    if (!widget.value.settings) {
        widget.value.settings = {};
    }

    widget.value.settings[settingKey] = settingValue;
}

function open(value: DesktopOverviewWidgetLayout): Promise<DesktopOverviewWidgetLayout> {
    widget.value = cloneWidget(value);
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!widget.value) {
        return;
    }

    resolveFunc?.(widget.value);
    showState.value = false;
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
