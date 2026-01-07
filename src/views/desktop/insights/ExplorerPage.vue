<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading || updating" :buttons="allTabs" v-model="activeTab" />
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading || updating" :model-value="currentExplorer.id">
                            <v-tab class="tab-text-truncate" key="new" value="" @click="createNewExplorer">
                                <span class="text-truncate">{{ tt('New Explorer') }}</span>
                            </v-tab>
                            <v-tab class="tab-text-truncate" :key="explorer.id" :value="explorer.id"
                                   v-for="explorer in allExplorers"
                                   @click="loadExplorer(explorer.id)">
                                <span class="text-truncate">{{ explorer.name || tt('Untitled Explorer') }}</span>
                            </v-tab>
<!--                            <v-btn class="text-left justify-start" variant="text" color="default" :rounded="false">-->
<!--                                <span class="ps-2">{{ tt('More Explorer') }}</span>-->
<!--                            </v-btn>-->
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
                                    <span>{{ tt('Insights Explorer') }}</span>
                                    <v-btn-group class="ms-4" color="default" density="comfortable" variant="outlined" divided>
                                        <v-btn class="button-icon-with-direction" :icon="mdiArrowLeft"
                                               :disabled="loading || updating || !canShiftDateRange"
                                               @click="shiftDateRange(-1)"/>
                                        <v-menu location="bottom" max-height="500">
                                            <template #activator="{ props }">
                                                <v-btn :disabled="loading || updating"
                                                       v-bind="props">{{ displayQueryDateRangeName }}</v-btn>
                                            </template>
                                            <v-list :selected="[currentFilter.dateRangeType]">
                                                <v-list-item :key="dateRange.type" :value="dateRange.type"
                                                             :append-icon="(currentFilter.dateRangeType === dateRange.type ? mdiCheck : undefined)"
                                                             v-for="dateRange in allDateRanges">
                                                    <v-list-item-title class="cursor-pointer"
                                                                       @click="setDateFilter(dateRange.type)">
                                                        <div class="d-flex align-center">
                                                            <span>{{ dateRange.displayName }}</span>
                                                        </div>
                                                        <div class="statistics-custom-datetime-range smaller" v-if="dateRange.isUserCustomRange && currentFilter.dateRangeType === dateRange.type && !!currentFilter.startTime && !!currentFilter.endTime">
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
                                               :disabled="loading || updating || !canShiftDateRange"
                                               @click="shiftDateRange(1)"/>
                                    </v-btn-group>

                                    <v-btn density="compact" color="default" variant="text" size="24"
                                           class="ms-2" :icon="true" :loading="loading" :disabled="updating" @click="reload(true)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20"/>
                                        </template>
                                        <v-icon :icon="mdiRefresh" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                    </v-btn>
                                    <v-spacer/>
                                    <v-btn class="ms-3" color="default" variant="outlined"
                                           :disabled="loading || updating" @click="saveExplorer(false)">
                                        {{ tt('Save Explorer') }}
                                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="updating"></v-progress-circular>
                                        <v-menu activator="parent" :open-on-hover="true">
                                            <v-list>
                                                <v-list-item :prepend-icon="mdiContentSaveOutline" @click="saveExplorer(true)">
                                                    <v-list-item-title>{{ tt('Save As New Explorer') }}</v-list-item-title>
                                                </v-list-item>
                                                <v-divider class="my-2" v-if="currentExplorer.id" />
                                                <v-list-item :prepend-icon="mdiPencilOutline" @click="setExplorerName" v-if="currentExplorer.id">
                                                    <v-list-item-title>{{ tt('Rename Explorer') }}</v-list-item-title>
                                                </v-list-item>
                                                <v-divider class="my-2" v-if="currentExplorer.id" />
                                                <v-list-item :prepend-icon="mdiDeleteOutline" @click="removeExplorer" v-if="currentExplorer.id">
                                                    <v-list-item-title>{{ tt('Delete Explorer') }}</v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </v-btn>
                                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                                           :disabled="loading || updating" :icon="true">
                                        <v-icon :icon="mdiDotsVertical" />
                                        <v-menu activator="parent">
                                            <v-list>
                                                <v-list-subheader :title="tt('Timezone Used for Date Range')"
                                                                  v-if="activeTab === 'query'"/>
                                                <template v-if="activeTab === 'query'">
                                                    <v-list-item :key="timezoneType.type" :value="timezoneType.type"
                                                                 :prepend-icon="timezoneTypeIconMap[timezoneType.type]"
                                                                 :append-icon="(currentExplorer.timezoneUsedForDateRange === timezoneType.type ? mdiCheck : undefined)"
                                                                 :title="timezoneType.displayName"
                                                                 v-for="timezoneType in allTimezoneTypesUsedForDateRange"
                                                                 @click="currentExplorer.timezoneUsedForDateRange = timezoneType.type"></v-list-item>
                                                </template>
                                                <v-list-item :prepend-icon="mdiExport"
                                                             :title="tt('Export Results')"
                                                             :disabled="loading || updating || !filteredTransactions || filteredTransactions.length < 1"
                                                             @click="exportResults"
                                                             v-if="activeTab === 'table' || activeTab === 'chart'"></v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </v-btn>
                                </div>
                            </template>

                            <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                                <v-window-item value="query">
                                    <explorer-query-tab :loading="loading" :disabled="loading || updating" />
                                </v-window-item>
                                <v-window-item value="table">
                                    <explorer-data-table-tab ref="explorerDataTableTab"
                                                             :loading="loading" :disabled="loading || updating"
                                                             @click:transaction="onShowTransaction" />
                                </v-window-item>
                                <v-window-item value="chart">
                                    <explorer-chart-tab ref="explorerChartTab"
                                                       :loading="loading" :disabled="loading || updating" />
                                </v-window-item>
                            </v-window>
                        </v-card>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="tt('Custom Date Range')"
                                 :min-time="currentFilter.startTime"
                                 :max-time="currentFilter.endTime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="setCustomDateFilter"
                                 @error="onShowDateRangeError" />

    <explorer-rename-dialog ref="explorerRenameDialog" />
    <edit-dialog ref="editDialog" :type="TransactionEditPageType.Transaction" />
    <export-dialog ref="exportDialog" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import ExplorerQueryTab from '@/views/desktop/insights/tabs/ExplorerQueryTab.vue';
import ExplorerDataTableTab from '@/views/desktop/insights/tabs/ExplorerDataTableTab.vue';
import ExplorerChartTab from '@/views/desktop/insights/tabs/ExplorerChartTab.vue';
import ExplorerRenameDialog from '@/views/desktop/insights/dialogs/ExplorerRenameDialog.vue';
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';
import ExportDialog from '@/views/desktop/statistics/transaction/dialogs/ExportDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useRouter, onBeforeRouteUpdate } from 'vue-router';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { type TransactionExplorerPartialFilter, type TransactionExplorerFilter, useExplorersStore } from '@/stores/explorer.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type WeekDayValue, type LocalizedDateRange, DateRangeScene, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';

import { type TransactionInsightDataItem, Transaction } from '@/models/transaction.ts';
import { type InsightsExplorerBasicInfo, InsightsExplorer } from '@/models/explorer.ts';

import {
    parseDateTimeFromUnixTime,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.ts';

import { generateRandomUUID } from '@/lib/misc.ts';

import {
    mdiMenu,
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck,
    mdiRefresh,
    mdiDotsVertical,
    mdiContentSaveOutline,
    mdiPencilOutline,
    mdiDeleteOutline,
    mdiHomeClockOutline,
    mdiInvoiceTextClockOutline,
    mdiExport
} from '@mdi/js';

interface InsightsExplorerProps {
    initId?: string;
    initActiveTab?: string,
    initDateRangeType?: string,
    initStartTime?: string,
    initEndTime?: string,
}

const props = defineProps<InsightsExplorerProps>();

type ExplorerPageTabType = 'query' | 'table' | 'chart';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type ExplorerDataTableTabType = InstanceType<typeof ExplorerDataTableTab>;
type ExplorerChartTabType = InstanceType<typeof ExplorerChartTab>;
type ExplorerRenameDialogType = InstanceType<typeof ExplorerRenameDialog>;
type EditDialogType = InstanceType<typeof EditDialog>;
type ExportDialogType = InstanceType<typeof ExportDialog>;

const router = useRouter();
const display = useDisplay();

const {
    tt,
    getAllDateRanges,
    getAllTimezoneTypesUsedForStatistics,
    formatDateTimeToLongDateTime,
    formatDateRange
} = useI18n();

const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const explorersStore = useExplorersStore();

const timezoneTypeIconMap = {
    [TimezoneTypeForStatistics.ApplicationTimezone.type]: mdiHomeClockOutline,
    [TimezoneTypeForStatistics.TransactionTimezone.type]: mdiInvoiceTextClockOutline
};

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const explorerDataTableTab = useTemplateRef<ExplorerDataTableTabType>('explorerDataTableTab');
const explorerChartTab = useTemplateRef<ExplorerChartTabType>('explorerChartTab');
const explorerRenameDialog = useTemplateRef<ExplorerRenameDialogType>('explorerRenameDialog');
const exportDialog = useTemplateRef<ExportDialogType>('exportDialog');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const loading = ref<boolean>(true);
const initing = ref<boolean>(true);
const updating = ref<boolean>(false);
const clientSessionId = ref<string>('');
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);
const activeTab = ref<ExplorerPageTabType>('query');
const showCustomDateRangeDialog = ref<boolean>(false);

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);

const allExplorers = computed<InsightsExplorerBasicInfo[]>(() => explorersStore.allInsightsExplorerBasicInfos.slice(0, 15));
const currentFilter = computed<TransactionExplorerFilter>(() => explorersStore.transactionExplorerFilter);
const currentExplorer = computed<InsightsExplorer>(() => explorersStore.currentInsightsExplorer);
const filteredTransactions = computed<TransactionInsightDataItem[]>(() => explorersStore.filteredTransactions);

const allDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.InsightsExplorer, true));
const allTimezoneTypesUsedForDateRange = computed<TypeAndDisplayName[]>(() => getAllTimezoneTypesUsedForStatistics());
const canShiftDateRange = computed<boolean>(() => currentFilter.value.dateRangeType !== DateRange.All.type);
const displayQueryDateRangeName = computed<string>(() => formatDateRange(currentFilter.value.dateRangeType, currentFilter.value.startTime, currentFilter.value.endTime));
const displayQueryStartTime = computed<string>(() => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(currentFilter.value.startTime)));
const displayQueryEndTime = computed<string>(() => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(currentFilter.value.endTime)));

const allTabs = computed<{ name: string, value: ExplorerPageTabType }[]>(() => {
    return [
        {
            name: tt('Query'),
            value: 'query'
        },
        {
            name: tt('Data Table'),
            value: 'table'
        },
        {
            name: tt('Chart'),
            value: 'chart'
        }
    ];
});

function getFilterLinkUrl(): string {
    return `/insights/explorer?${explorersStore.getTransactionExplorerPageParams(currentExplorer.value.id, activeTab.value)}`;
}

function init(initProps: InsightsExplorerProps): void {
    clientSessionId.value = generateRandomUUID();

    const filter: TransactionExplorerPartialFilter = {
        dateRangeType: initProps.initDateRangeType ? parseInt(initProps.initDateRangeType) : undefined,
        startTime: initProps.initStartTime ? parseInt(initProps.initStartTime) : undefined,
        endTime: initProps.initEndTime ? parseInt(initProps.initEndTime) : undefined
    };

    let needReload = false;

    if (filter.dateRangeType !== currentFilter.value.dateRangeType) {
        needReload = true;
    } else if (filter.dateRangeType === DateRange.Custom.type) {
        if (filter.startTime !== currentFilter.value.startTime
            || filter.endTime !== currentFilter.value.endTime) {
            needReload = true;
        }
    }

    if (initProps.initActiveTab === 'query' || initProps.initActiveTab === 'table' || initProps.initActiveTab === 'chart') {
        if (initProps.initActiveTab !== activeTab.value) {
            activeTab.value = initProps.initActiveTab;
        }
    } else {
        activeTab.value = 'query';
    }

    explorersStore.initTransactionExplorerFilter(filter);

    if (initProps.initId) {
        if (explorersStore.currentInsightsExplorer.id !== initProps.initId) {
            loadExplorer(initProps.initId);
        }
    } else {
        explorersStore.updateCurrentInsightsExplorer(InsightsExplorer.createNewExplorer(generateRandomUUID()));
    }

    if (!needReload && !explorersStore.transactionExplorerStateInvalid && !explorersStore.insightsExplorerListStateInvalid) {
        loading.value = false;
        initing.value = false;
        return;
    }

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false }),
        explorersStore.loadAllInsightsExplorerBasicInfos({ force: false })
    ]).then(() => {
        return explorersStore.loadAllTransactions({ force: false });
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

    return explorersStore.loadAllTransactions({
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

function createNewExplorer(): void {
    if (!currentExplorer.value.id) {
        return;
    }

    explorersStore.updateCurrentInsightsExplorer(InsightsExplorer.createNewExplorer(generateRandomUUID()));
    router.push(getFilterLinkUrl());
}

function loadExplorer(explorerId: string): void {
    if (currentExplorer.value.id === explorerId) {
        return;
    }

    loading.value = true;

    explorersStore.getInsightsExplorer({
        explorerId: explorerId
    }).then(explorer => {
        explorersStore.updateCurrentInsightsExplorer(explorer);
        loading.value = false;
        router.push(getFilterLinkUrl());
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function saveExplorer(saveAs?: boolean): void {
    if (saveAs || !currentExplorer.value.name) {
        explorerRenameDialog.value?.open(currentExplorer.value.name || '').then((newName: string) => {
            currentExplorer.value.name = newName;
            doSaveExplorer(saveAs);
        })
    } else {
        doSaveExplorer(saveAs);
    }
}

function doSaveExplorer(saveAs?: boolean): Promise<unknown> {
    const oldExplorerId = currentExplorer.value.id;

    updating.value = true;

    return explorersStore.saveInsightsExplorer({
        explorer: currentExplorer.value,
        saveAs: saveAs,
        clientSessionId: clientSessionId.value
    }).then(newExplorer => {
        updating.value = false;
        clientSessionId.value = generateRandomUUID();
        explorersStore.updateCurrentInsightsExplorer(newExplorer);

        if (oldExplorerId !== newExplorer.id) {
            router.push(getFilterLinkUrl());
        }
    }).catch(error => {
        updating.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function setExplorerName(): void {
    explorerRenameDialog.value?.open(currentExplorer.value.name || '').then((newName: string) => {
        currentExplorer.value.name = newName;
    });
}

function removeExplorer(): void {
    if (!currentExplorer.value.id) {
        return;
    }

    confirmDialog.value?.open('Are you sure you want to delete this explorer?').then(() => {
        updating.value = true;

        explorersStore.deleteInsightsExplorer({
            explorer: currentExplorer.value
        }).then(() => {
            updating.value = false;
            createNewExplorer();
        }).catch(error => {
            updating.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function exportResults(): void {
    if (activeTab.value === 'table' && filteredTransactions.value) {
        const results = explorerDataTableTab.value?.buildExportResults();

        if (results) {
            exportDialog.value?.open(results);
        }
    } else if (activeTab.value === 'chart') {
        const results = explorerChartTab.value?.buildExportResults();

        if (results) {
            exportDialog.value?.open(results);
        }
    }
}

function setDateFilter(dateType: number): void {
    if (dateType === DateRange.Custom.type) { // Custom
        showCustomDateRangeDialog.value = true;
        return;
    } else if (currentFilter.value.dateRangeType === dateType) {
        return;
    }

    const dateRange = getDateRangeByDateType(dateType, firstDayOfWeek.value, fiscalYearStart.value);

    if (!dateRange) {
        return;
    }

    const changed = explorersStore.updateTransactionExplorerFilter({
        dateRangeType: dateRange.dateType,
        startTime: dateRange.minTime,
        endTime: dateRange.maxTime
    });

    if (changed) {
        loading.value = true;
        explorersStore.updateTransactionExplorerInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setCustomDateFilter(startTime: number, endTime: number): void {
    if (!startTime || !endTime) {
        return;
    }

    const chartDateType = getDateTypeByDateRange(startTime, endTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.InsightsExplorer);

    const changed = explorersStore.updateTransactionExplorerFilter({
        dateRangeType: chartDateType,
        startTime: startTime,
        endTime: endTime
    });

    showCustomDateRangeDialog.value = false;

    if (changed) {
        loading.value = true;
        explorersStore.updateTransactionExplorerInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function shiftDateRange(scale: number): void {
    if (currentFilter.value.dateRangeType === DateRange.All.type) {
        return;
    }

    const newDateRange = getShiftedDateRangeAndDateType(currentFilter.value.startTime, currentFilter.value.endTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);

    const changed = explorersStore.updateTransactionExplorerFilter({
        dateRangeType: newDateRange.dateType,
        startTime: newDateRange.minTime,
        endTime: newDateRange.maxTime
    });

    if (changed) {
        loading.value = true;
        explorersStore.updateTransactionExplorerInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function onShowTransaction(transaction: TransactionInsightDataItem): void {
    editDialog.value?.open({
        id: transaction.id,
        currentTransaction: Transaction.of(transaction)
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        reload(false);
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
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
