<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="allTabs" v-model="activeTab" />
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4" v-if="activeTab === 'table'">
                            <span class="text-subtitle-2">{{ tt('Transactions Per Page') }}</span>
                            <v-select class="mt-2" density="compact"
                                      item-title="name"
                                      item-value="value"
                                      :disabled="loading"
                                      :items="allPageCounts"
                                      v-model="countPerPage"
                            />
                        </div>
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="currentExploreId">
                            <v-tab class="tab-text-truncate" key="new" value="">
                                <span class="text-truncate">{{ tt('New Explore') }}</span>
                            </v-tab>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-card variant="flat" min-height="800">
                            <template #title>
                                <div class="title-and-toolbar d-flex align-center">
                                    <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                           :ripple="false" :icon="true" @click="showNav = !showNav">
                                        <v-icon :icon="mdiMenu" size="24" />
                                    </v-btn>
                                    <span>{{ tt('Insights & Explore') }}</span>
                                    <v-btn-group class="ms-4" color="default" density="comfortable" variant="outlined" divided>
                                        <v-btn class="button-icon-with-direction" :icon="mdiArrowLeft"
                                               :disabled="loading || !canShiftDateRange"
                                               @click="shiftDateRange(-1)"/>
                                        <v-menu location="bottom" max-height="500">
                                            <template #activator="{ props }">
                                                <v-btn :disabled="loading"
                                                       v-bind="props">{{ displayQueryDateRangeName }}</v-btn>
                                            </template>
                                            <v-list :selected="[query.dateRangeType]">
                                                <v-list-item :key="dateRange.type" :value="dateRange.type"
                                                             :append-icon="(query.dateRangeType === dateRange.type ? mdiCheck : undefined)"
                                                             v-for="dateRange in allDateRanges">
                                                    <v-list-item-title class="cursor-pointer"
                                                                       @click="setDateFilter(dateRange.type)">
                                                        <div class="d-flex align-center">
                                                            <span>{{ dateRange.displayName }}</span>
                                                        </div>
                                                        <div class="statistics-custom-datetime-range smaller" v-if="dateRange.isUserCustomRange && query.dateRangeType === dateRange.type && !!query.startTime && !!query.endTime">
                                                            <span>{{ displayQueryStartTime }}</span>
                                                            <span>&nbsp;-&nbsp;</span>
                                                            <br/>
                                                            <span>{{ displayQueryEndTime }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                        <v-btn class="button-icon-with-direction" :icon="mdiArrowRight"
                                               :disabled="loading || !canShiftDateRange"
                                               @click="shiftDateRange(1)"/>
                                    </v-btn-group>

                                    <v-btn density="compact" color="default" variant="text" size="24"
                                           class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20"/>
                                        </template>
                                        <v-icon :icon="mdiRefresh" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                    </v-btn>
                                    <v-spacer/>
                                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                                           :disabled="loading" :icon="true"
                                           v-if="activeTab !== 'query'">
                                        <v-icon :icon="mdiDotsVertical" />
                                        <v-menu activator="parent">
                                            <v-list>
                                                <v-list-item :prepend-icon="mdiExport"
                                                             :title="tt('Export Results')"
                                                             :disabled="loading || !filteredTransactions || filteredTransactions.length < 1"
                                                             @click="exportResults"></v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </v-btn>
                                </div>
                            </template>

                            <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                                <v-window-item value="query">
                                    <explore-query-tab :loading="loading" />
                                </v-window-item>
                                <v-window-item value="table">
                                    <explore-data-table-tab ref="exploreDataTableTab"
                                                            :loading="loading"
                                                            v-model:count-per-page="countPerPage" />
                                </v-window-item>
                                <v-window-item value="chart">

                                </v-window-item>
                            </v-window>
                        </v-card>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="tt('Custom Date Range')"
                                 :min-time="query.startTime"
                                 :max-time="query.endTime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="setCustomDateFilter"
                                 @error="onShowDateRangeError" />

    <export-dialog ref="exportDialog" />

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ExploreQueryTab from '@/views/desktop/insights/tabs/ExploreQueryTab.vue';
import ExploreDataTableTab from '@/views/desktop/insights/tabs/ExploreDataTableTab.vue';
import ExportDialog from '@/views/desktop/statistics/transaction/dialogs/ExportDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useRouter, onBeforeRouteUpdate } from 'vue-router';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { type TransactionExplorePartialFilter, type TransactionExploreFilter, useExploresStore } from '@/stores/explore.ts';

import type { NameNumeralValue } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';
import { type WeekDayValue, type LocalizedDateRange, DateRangeScene, DateRange } from '@/core/datetime.ts';

import {
    type TransactionInsightDataItem
} from '@/models/transaction.ts';

import {
    parseDateTimeFromUnixTime,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.ts';

import {
    mdiMenu,
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck,
    mdiRefresh,
    mdiDotsVertical,
    mdiExport
} from '@mdi/js';

interface InsightsExploreProps {
    initId?: string;
    initActiveTab?: string,
    initDateRangeType?: string,
    initStartTime?: string,
    initEndTime?: string,
}

const props = defineProps<InsightsExploreProps>();

type ExplorePageTabType = 'query' | 'table' | 'chart';
type SnackBarType = InstanceType<typeof SnackBar>;
type ExploreDataTableTabType = InstanceType<typeof ExploreDataTableTab>;
type ExportDialogType = InstanceType<typeof ExportDialog>;

const router = useRouter();
const display = useDisplay();

const {
    tt,
    getAllDateRanges,
    getCurrentNumeralSystemType,
    formatDateTimeToLongDateTime,
    formatDateRange
} = useI18n();

const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const exploresStore = useExploresStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const exploreDataTableTab = useTemplateRef<ExploreDataTableTabType>('exploreDataTableTab');
const exportDialog = useTemplateRef<ExportDialogType>('exportDialog');

const loading = ref<boolean>(true);
const initing = ref<boolean>(true);
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);
const activeTab = ref<ExplorePageTabType>('query');
const currentExploreId = ref<string>('');
const countPerPage = ref<number>(15);
const showCustomDateRangeDialog = ref<boolean>(false);

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);
const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

const query = computed<TransactionExploreFilter>(() => exploresStore.transactionExploreFilter);
const filteredTransactions = computed<TransactionInsightDataItem[]>(() => exploresStore.filteredTransactions);

const allDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.InsightsExplore, true));
const canShiftDateRange = computed<boolean>(() => query.value.dateRangeType !== DateRange.All.type);
const displayQueryDateRangeName = computed<string>(() => formatDateRange(query.value.dateRangeType, query.value.startTime, query.value.endTime));
const displayQueryStartTime = computed<string>(() => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(query.value.startTime)));
const displayQueryEndTime = computed<string>(() => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(query.value.endTime)));

const allTabs = computed<{ name: string, value: ExplorePageTabType }[]>(() => {
    return [
        {
            name: tt('Query'),
            value: 'query'
        },
        {
            name: tt('Data Table'),
            value: 'table'
        }
    ];
});

const allPageCounts = computed<NameNumeralValue[]>(() => {
    const pageCounts: NameNumeralValue[] = [];
    const availableCountPerPage: number[] = [ 5, 10, 15, 20, 25, 30, 50 ];

    for (const count of availableCountPerPage) {
        pageCounts.push({ value: count, name: numeralSystem.value.formatNumber(count) });
    }

    pageCounts.push({ value: -1, name: tt('All') });

    return pageCounts;
});

function getFilterLinkUrl(): string {
    return `/insights/explore?${exploresStore.getTransactionExplorePageParams(currentExploreId.value, activeTab.value)}`;
}

function init(initProps: InsightsExploreProps): void {
    const filter: TransactionExplorePartialFilter = {
        dateRangeType: initProps.initDateRangeType ? parseInt(initProps.initDateRangeType) : undefined,
        startTime: initProps.initStartTime ? parseInt(initProps.initStartTime) : undefined,
        endTime: initProps.initEndTime ? parseInt(initProps.initEndTime) : undefined
    };

    let needReload = false;

    if (filter.dateRangeType !== query.value.dateRangeType) {
        needReload = true;
    } else if (filter.dateRangeType === DateRange.Custom.type) {
        if (filter.startTime !== query.value.startTime
            || filter.endTime !== query.value.endTime) {
            needReload = true;
        }
    }

    if (initProps.initActiveTab === 'query' || initProps.initActiveTab === 'table' || initProps.initActiveTab === 'chart') {
        if (initProps.initActiveTab !== activeTab.value) {
            activeTab.value = initProps.initActiveTab;
            needReload = true;
        }
    } else {
        activeTab.value = 'query';
    }

    exploresStore.initTransactionExploreFilter(filter);

    if (!needReload && !exploresStore.transactionExploreStateInvalid) {
        loading.value = false;
        initing.value = false;
        return;
    }

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ]).then(() => {
        return exploresStore.loadAllTransactions({ force: false });
    }).then(() => {
        loading.value = false;
        initing.value = false;
    }).catch(error => {
        loading.value = false;
        initing.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function reload(force: boolean): Promise<unknown> | null {
    loading.value = true;

    return exploresStore.loadAllTransactions({
        force: force
    }).then(() => {
        loading.value = false;

        if (force) {
            snackbar.value?.showMessage('Data has been updated');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function exportResults(): void {
    if (activeTab.value === 'table' && filteredTransactions.value) {
        const results = exploreDataTableTab.value?.buildExportResults();

        if (results) {
            exportDialog.value?.open(results);
        }
    }
}

function setDateFilter(dateType: number): void {
    if (dateType === DateRange.Custom.type) { // Custom
        showCustomDateRangeDialog.value = true;
        return;
    } else if (query.value.dateRangeType === dateType) {
        return;
    }

    const dateRange = getDateRangeByDateType(dateType, firstDayOfWeek.value, fiscalYearStart.value);

    if (!dateRange) {
        return;
    }

    const changed = exploresStore.updateTransactionExploreFilter({
        dateRangeType: dateRange.dateType,
        startTime: dateRange.minTime,
        endTime: dateRange.maxTime
    });

    if (changed) {
        loading.value = true;
        exploresStore.updateTransactionExploreInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setCustomDateFilter(startTime: number, endTime: number): void {
    if (!startTime || !endTime) {
        return;
    }

    const chartDateType = getDateTypeByDateRange(startTime, endTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.InsightsExplore);

    const changed = exploresStore.updateTransactionExploreFilter({
        dateRangeType: chartDateType,
        startTime: startTime,
        endTime: endTime
    });

    showCustomDateRangeDialog.value = false;

    if (changed) {
        loading.value = true;
        exploresStore.updateTransactionExploreInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function shiftDateRange(scale: number): void {
    if (query.value.dateRangeType === DateRange.All.type) {
        return;
    }

    const newDateRange = getShiftedDateRangeAndDateType(query.value.startTime, query.value.endTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);

    const changed = exploresStore.updateTransactionExploreFilter({
        dateRangeType: newDateRange.dateType,
        startTime: newDateRange.minTime,
        endTime: newDateRange.maxTime
    });

    if (changed) {
        loading.value = true;
        exploresStore.updateTransactionExploreInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function onShowDateRangeError(message: string): void {
    snackbar.value?.showError(message);
}

onBeforeRouteUpdate((to) => {
    if (to.query) {
        init({
            initId: (to.query['id'] as string | null) || undefined,
            initActiveTab: (to.query['activeTab'] as string | null) || undefined,
            initDateRangeType: (to.query['dateRangeType'] as string | null) || undefined,
            initStartTime: (to.query['startTime'] as string | null) || undefined,
            initEndTime: (to.query['endTime'] as string | null) || undefined,
        });
    } else {
        init({});
    }
});

watch(() => display.mdAndUp.value, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

watch(activeTab, () => {
    router.push(getFilterLinkUrl());
});

init(props);
</script>
