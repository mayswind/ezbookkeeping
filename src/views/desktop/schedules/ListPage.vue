<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Scheduled Transactions') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn class="ms-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" class="ms-2"
                               :aria-label="tt('Refresh')" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-btn-group class="ms-2" color="default" density="comfortable" variant="outlined" divided
                                     v-if="supportCalendar">
                            <v-btn class="button-icon-with-direction" :icon="mdiArrowLeft"
                                   :aria-label="tt('Previous Period')"
                                   :disabled="loading || !calendarView || currentYearMonth <= minYearMonth"
                                   @click="shiftMonth(-1)"/>
                            <v-menu location="bottom" max-height="500">
                                <template #activator="{ props }">
                                    <v-btn :disabled="loading || !calendarView" v-bind="props">{{ displayMonth }}</v-btn>
                                </template>
                                <v-list>
                                    <v-list-item :title="tt('This Month')"
                                                 @click="changeCustomMonthDateFilter(getCurrentDateTime().toGregorianCalendarYear0BasedMonth())"></v-list-item>
                                    <v-list-item :title="tt('Custom Date')"
                                                 @click="showCustomMonthDialog = true"></v-list-item>
                                </v-list>
                            </v-menu>
                            <v-btn class="button-icon-with-direction" :icon="mdiArrowRight"
                                   :aria-label="tt('Next Period')"
                                   :disabled="loading || !calendarView"
                                   @click="shiftMonth(1)"/>
                        </v-btn-group>

                        <v-btn-toggle class="toggle-buttons ms-3" style="min-height: 36px"
                                      density="compact" color="default" variant="outlined"
                                      mandatory="force"
                                      :disabled="loading || updating || displayOrderModified"
                                      :model-value="calendarView"
                                      v-if="supportCalendar"
                                      @update:model-value="switchView">
                            <v-btn :value="false">{{ tt('List View') }}</v-btn>
                            <v-btn :value="true" :disabled="!supportCalendar">{{ tt('Calendar View') }}</v-btn>
                        </v-btn-toggle>

                        <v-btn-toggle class="toggle-buttons ms-3" style="min-height: 36px"
                                      density="compact" color="default" variant="outlined"
                                      mandatory="force"
                                      :disabled="loading || updating || displayOrderModified"
                                      :model-value="currentStatus"
                                      @update:model-value="switchStatus">
                            <v-btn :key="item.type" :value="item.type"
                                   v-for="item in allStatuses">{{ item.displayName }}</v-btn>
                        </v-btn-toggle>

                        <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                               :aria-label="tt('More')" :disabled="loading || updating" :icon="true">
                            <v-icon :icon="mdiDotsVertical" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="mdiEyeOutline"
                                                 :title="tt('Show Hidden Scheduled Transactions')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="mdiEyeOffOutline"
                                                 :title="tt('Hide Hidden Scheduled Transactions')"
                                                 v-if="showHidden" @click="showHidden = false"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </template>

                <v-card-text class="scheduled-transaction-calendar-container" v-if="supportCalendar && calendarView">
                    <scheduled-transaction-calendar show-alternate-date
                                                    week-day-name-type="long"
                                                    date-align="left"
                                                    :loading="loading && !occurrences.length"
                                                    :show-hidden="showHidden"
                                                    :templates="visibleTemplates"
                                                    :occurrences="occurrences"
                                                    :model-value="currentYearMonth"
                                                    @edit="edit" />
                </v-card-text>

                <v-table class="scheduled-transactions-table table-striped" :hover="!loading" v-if="!supportCalendar || !calendarView">
                    <thead>
                    <tr>
                        <th>
                            <div class="d-flex align-center">
                                <span>{{ tt('Template Name') }}</span>
                                <v-spacer/>
                                <span>{{ tt('Operation') }}</span>
                            </div>
                        </th>
                    </tr>
                    </thead>

                    <tbody v-if="loading && !visibleTemplates.length">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && !visibleTemplates.length">
                    <tr>
                        <td>{{ tt('No available scheduled transactions') }}</td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    :disabled="loading || updating || visibleTemplates.length < 2"
                                    v-model="templates"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr class="scheduled-transactions-table-row" v-if="showHidden || !element.hidden"
                                @mouseenter="hoveredTemplateId = element.id" @mouseleave="hoveredTemplateId = ''">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <v-badge class="right-bottom-icon" color="secondary"
                                                     location="bottom right" offset-x="8"
                                                     :aria-label="tt('Hidden')" :icon="mdiEyeOffOutline"
                                                     v-if="element.hidden">
                                                <v-icon size="20" start :icon="mdiClockTimeNineOutline"/>
                                            </v-badge>
                                            <v-icon size="20" start :icon="mdiClockTimeNineOutline" v-else-if="!element.hidden"/>
                                            <div>
                                                <div class="scheduled-transaction-name">{{ element.name }}</div>
                                                <div class="text-body-small text-medium-emphasis">
                                                    <div>
                                                        <span>{{ getDisplayFrequency(element) }}</span>
                                                        <span class="mx-1" v-if="element.nextScheduledTime">·</span>
                                                        <span v-if="element.nextScheduledTime">{{ getDisplayNextScheduledTime(element) }}</span>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>

                                        <v-spacer/>

                                        <template v-if="hoveredTemplateId === element.id && !loading">
                                            <v-btn class="px-2 ms-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="element.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                                   :loading="templateHiding[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="hide(element, !element.hidden)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ element.hidden ? tt('Show') : tt('Hide') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiPencilOutline"
                                                   :disabled="loading || updating"
                                                   @click="edit(element)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ tt('Edit') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiDeleteOutline"
                                                   :loading="templateRemoving[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="remove(element)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ tt('Delete') }}
                                            </v-btn>
                                        </template>

                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && visibleTemplates.length > 1 ? 'drag-handle' : 'disabled'"
                                                    :aria-label="tt('Drag to Reorder')" :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && visibleTemplates.length > 1 && hoveredTemplateId === element.id">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                    </div>
                                </td>
                            </tr>
                        </template>
                    </draggable-list>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <month-selection-dialog :title="tt('Select Month')"
                            :model-value="queryMonth"
                            v-model:show="showCustomMonthDialog"
                            @update:modelValue="changeCustomMonthDateFilter"
                            @error="desktopPageStore.addErrorToast" />

    <edit-dialog ref="editDialog" :type="TransactionEditPageType.Template" />

    <confirm-dialog ref="confirmDialog"/>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import ScheduledTransactionCalendar from '@/components/common/ScheduledTransactionCalendar.vue';
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';
import { useScheduledTransactionListPageBase } from '@/views/base/schedules/ScheduledTransactionListPageBase.ts';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import type { Year0BasedMonth } from '@/core/datetime.ts';
import { TemplateType, ScheduledTransactionStatus } from '@/core/template.ts';
import type { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    getCurrentDateTime,
    getYearMonthFirstUnixTime,
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical,
    mdiClockTimeNineOutline,
    mdiArrowLeft,
    mdiArrowRight
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type EditDialogType = InstanceType<typeof EditDialog>;

const { tt } = useI18n();

const transactionTemplatesStore = useTransactionTemplatesStore();
const desktopPageStore = useDesktopPageStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const {
    loading,
    updating,
    showHidden,
    calendarView,
    displayOrderModified,
    minYearMonth,
    currentYearMonth,
    currentStatus,
    templates,
    occurrences,
    allStatuses,
    supportCalendar,
    visibleTemplates,
    queryMonth,
    displayMonth,
    loadScheduledTransactions,
    loadScheduledTransactionCalendar,
    saveScheduledTransactionDisplayOrders,
    getDisplayFrequency,
    getDisplayNextScheduledTime
} = useScheduledTransactionListPageBase();

const hoveredTemplateId = ref<string>('');
const templateHiding = ref<Record<string, boolean>>({});
const templateRemoving = ref<Record<string, boolean>>({});
const showCustomMonthDialog = ref<boolean>(false);

function switchStatus(value: ScheduledTransactionStatus): void {
    if (currentStatus.value === value || updating.value || displayOrderModified.value) {
        return;
    }

    currentStatus.value = value;
    calendarView.value = false;
    templates.value = [];
    occurrences.value = [];
    reload();
}

function switchView(value: boolean): void {
    if (calendarView.value === value || loading.value || updating.value || displayOrderModified.value || (value && !supportCalendar.value)) {
        return;
    }

    calendarView.value = value;

    if (calendarView.value) {
        loadCalendar();
    }
}

function changeCustomMonthDateFilter(yearMonth: Year0BasedMonth): void {
    if (!yearMonth) {
        return;
    }

    const yearMonthDateTime = parseDateTimeFromUnixTime(getYearMonthFirstUnixTime(yearMonth));
    const yearDashMonth = yearMonthDateTime.getGregorianCalendarYearDashMonth();

    if (yearDashMonth < minYearMonth.value) {
        desktopPageStore.addErrorToast('Date is too early');
        return;
    }

    showCustomMonthDialog.value = false;

    if (currentYearMonth.value === yearDashMonth) {
        return;
    }

    currentYearMonth.value = yearDashMonth;
    occurrences.value = [];
    loadCalendar();
}

function shiftMonth(scale: number): void {
    const currentYearMonthDateTime = parseDateTimeFromUnixTime(getYearMonthFirstUnixTime(currentYearMonth.value));
    const newYearMonthDateTime = currentYearMonthDateTime.add(scale, 'months');
    const yearDashMonth = newYearMonthDateTime.toGregorianCalendarYear0BasedMonth();

    changeCustomMonthDateFilter(yearDashMonth);
}

function reload(force?: boolean): Promise<void> {
    if (currentYearMonth.value < minYearMonth.value) {
        currentYearMonth.value = minYearMonth.value;
    }

    loading.value = true;

    return loadScheduledTransactions(force).then(() => {
        displayOrderModified.value = false;

        if (supportCalendar.value && calendarView.value) {
            return loadScheduledTransactionCalendar();
        }

        return;
    }).then(() => {
        loading.value = false;

        if (force) {
            desktopPageStore.addToast('Template list has been updated');
        }
    }).catch(error => {
        loading.value = false;

        if (error && error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (error && !error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function loadCalendar(): Promise<void> {
    loading.value = true;

    return loadScheduledTransactionCalendar().then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function add(): void {
    editDialog.value?.open({
        templateType: TemplateType.Schedule.type
    }).then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
            reload();
        }
    }).catch(error => {
        if (error && !error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function edit(template: TransactionTemplate): void {
    editDialog.value?.open({
        id: template.id,
        currentTemplate: template
    }).then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
            reload();
        }
    }).catch(error => {
        if (error && !error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function hide(template: TransactionTemplate, hidden: boolean): void {
    updating.value = true;
    templateHiding.value[template.id] = true;

    transactionTemplatesStore.hideTemplate({
        template: template,
        hidden: hidden
    }).then(() => {
        template.hidden = hidden;
        updating.value = false;
        templateHiding.value[template.id] = false;
    }).catch(error => {
        updating.value = false;
        templateHiding.value[template.id] = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function remove(template: TransactionTemplate): void {
    confirmDialog.value?.open('Are you sure you want to delete this template?').then(() => {
        updating.value = true;
        templateRemoving.value[template.id] = true;

        transactionTemplatesStore.deleteTemplate({
            template: template
        }).then(() => {
            return reload();
        }).then(() => {
            updating.value = false;
            templateRemoving.value[template.id] = false;
        }).catch(error => {
            updating.value = false;
            templateRemoving.value[template.id] = false;

            if (!error.processed) {
                desktopPageStore.addErrorToast(error);
            }
        });
    });
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    updating.value = true;

    saveScheduledTransactionDisplayOrders().then(() => {
        return reload();
    }).then(() => {
        updating.value = false;
    }).catch(error => {
        updating.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function onMove(event: { moved?: { element: { id: string } } }): void {
    if (event?.moved?.element.id) {
        displayOrderModified.value = true;
    }
}

reload();
</script>

<style>
.scheduled-transaction-calendar-container {
    .scheduled-transaction-calendar {
        font-size: var(--dp-font-size);

        .scheduled-calendar-grid {
            .scheduled-calendar-day {
                min-height: 120px;

                .scheduled-calendar-alternate-date {
                    font-size: 0.9rem;
                }

                .scheduled-calendar-event {
                    font-size: 0.8125rem;
                    border-inline-start: 2px solid rgb(var(--ebk-primary-color));

                    &.scheduled-calendar-event-income {
                        border-inline-start-color: rgb(var(--v-theme-income));
                        background: rgba(var(--v-theme-income), 0.12);
                    }

                    &.scheduled-calendar-event-expense {
                        border-inline-start-color: rgb(var(--v-theme-expense));
                        background: rgba(var(--v-theme-expense), 0.12);
                    }
                }
            }
        }
    }
}
</style>
