<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('Projections') }}</span>
                        <v-btn class="ml-3" density="comfortable" variant="tonal"
                               :disabled="loading" :prepend-icon="mdiCalendarMonthOutline"
                               @click="showMonthRangeDialog = true">
                            {{ displayPeriod }}
                        </v-btn>
                        <v-btn class="ml-2" density="comfortable" color="default" variant="text"
                               :disabled="loading" :icon="true" @click="reload">
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-progress-circular indeterminate size="20" width="2" class="ml-2" v-if="loading" />
                    </div>
                </template>

                <v-card-text>
                    <p class="text-body-2 mb-4">{{ tt('Amounts up to today come from your transactions, later ones are estimated from your scheduled transactions') }}</p>
                    <projection-table :data="projectionTableData" />
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <month-range-selection-dialog :title="tt('Custom Date Range')"
                                  :min-time="projectionsStore.projectionStartYearMonth || undefined"
                                  :max-time="projectionsStore.projectionEndYearMonth || undefined"
                                  v-model:show="showMonthRangeDialog"
                                  @dateRange:change="setPeriod"
                                  @error="onPeriodError" />

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, onMounted, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useProjectionsStore } from '@/stores/projection.ts';

import type { TextualYearMonth, Year0BasedMonth } from '@/core/datetime.ts';
import type { ProjectionTableData } from '@/core/projection.ts';

import {
    getCurrentUnixTime,
    getYear0BasedMonthObjectFromUnixTime,
    getYearMonthStringFromYear0BasedMonthObject,
    getYearMonthFirstUnixTime,
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

import {
    mdiCalendarMonthOutline,
    mdiRefresh
} from '@mdi/js';

// How far ahead the table looks the first time the page is opened. The period starts at the current
// month rather than in the past, so the table reads as a projection instead of as a history.
const DEFAULT_PROJECTION_MONTH_COUNT = 12;

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    formatDateTimeToGregorianLikeShortYearMonth
} = useI18n();

const projectionsStore = useProjectionsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const showMonthRangeDialog = ref<boolean>(false);

const projectionTableData = computed<ProjectionTableData>(() => projectionsStore.projectionTableData);

const displayPeriod = computed<string>(() => {
    const months = projectionTableData.value.months;

    if (!months.length) {
        return tt('Custom Date Range');
    }

    const firstMonth = months[0]!;
    const lastMonth = months[months.length - 1]!;

    return `${formatYearMonth(firstMonth.year, firstMonth.month)} ~ ${formatYearMonth(lastMonth.year, lastMonth.month)}`;
});

function formatYearMonth(year: number, month1base: number): string {
    return formatDateTimeToGregorianLikeShortYearMonth(parseDateTimeFromUnixTime(getYearMonthFirstUnixTime({ year: year, month1base: month1base })));
}

function getDefaultPeriod(): { startYearMonth: TextualYearMonth | '', endYearMonth: TextualYearMonth | '' } {
    const startYearMonth = getYear0BasedMonthObjectFromUnixTime(getCurrentUnixTime());
    const endMonth0base = startYearMonth.month0base + DEFAULT_PROJECTION_MONTH_COUNT - 1;
    const endYearMonth: Year0BasedMonth = {
        year: startYearMonth.year + Math.floor(endMonth0base / 12),
        month0base: endMonth0base % 12
    };

    return {
        startYearMonth: getYearMonthStringFromYear0BasedMonthObject(startYearMonth),
        endYearMonth: getYearMonthStringFromYear0BasedMonthObject(endYearMonth)
    };
}

function load(force: boolean): void {
    loading.value = true;

    projectionsStore.loadProjections({ force: force }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function reload(): void {
    load(true);
}

function setPeriod(minYearMonth: TextualYearMonth | '', maxYearMonth: TextualYearMonth | ''): void {
    if (!minYearMonth || !maxYearMonth) {
        return;
    }

    showMonthRangeDialog.value = false;
    projectionsStore.setProjectionPeriod(minYearMonth, maxYearMonth);
    load(false);
}

function onPeriodError(message: string): void {
    snackbar.value?.showMessage(message);
}

onMounted(() => {
    if (!projectionsStore.projectionStartYearMonth || !projectionsStore.projectionEndYearMonth) {
        const defaultPeriod = getDefaultPeriod();
        projectionsStore.setProjectionPeriod(defaultPeriod.startYearMonth, defaultPeriod.endYearMonth);
    }

    load(false);
});
</script>
