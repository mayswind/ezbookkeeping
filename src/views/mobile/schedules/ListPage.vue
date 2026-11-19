<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')" v-if="!sortable"></f7-nav-left>
            <f7-nav-left v-else-if="sortable">
                <f7-link icon-f7="xmark" :class="{ 'disabled': displayOrderSaving }" :aria-label="tt('Cancel')" @click="cancelSort"></f7-link>
            </f7-nav-left>
            <f7-nav-title>
                <f7-link popover-open=".scheduled-status-popover-menu" :class="{ 'disabled': loading || sortable || displayOrderModified }">
                    <span style="color: var(--f7-text-color)">{{ displayStatusName }}</span>
                    <f7-icon class="page-title-bar-icon" color="gray" style="opacity: 0.5" f7="chevron_down_circle_fill"></f7-icon>
                </f7-link>
            </f7-nav-title>
            <f7-nav-right :class="{ 'navbar-compact-icons': true, 'disabled': loading }">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !templates.length || sortable }" :aria-label="tt('More')" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="plus" :aria-label="tt('Add')" :href="'/template/add?templateType=' + TemplateType.Schedule.type" v-if="!sortable"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': displayOrderSaving || !displayOrderModified }" :aria-label="tt('Save')" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-popover class="scheduled-status-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="item.displayName"
                              :class="{ 'list-item-selected': currentStatus === item.type }"
                              :key="item.type"
                              v-for="item in allStatuses"
                              @click="switchStatus(item.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="currentStatus === item.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-block strong inset class="scheduled-transaction-calendar-container margin-vertical-half no-padding" v-if="supportCalendar && calendarView && !sortable">
            <scheduled-transaction-calendar show-alternate-date
                                            week-day-name-type="short"
                                            date-align="center"
                                            :loading="(loading || loadingCalendar) && !occurrences.length"
                                            :show-hidden="showHidden"
                                            :templates="visibleTemplates"
                                            :occurrences="occurrences"
                                            :model-value="currentYearMonth"
                                            @edit="edit"/>
        </f7-block>

        <f7-list strong inset dividers sortable class="margin-top-half scheduled-transaction-list skeleton-text"
                 v-if="loading && (!supportCalendar || !calendarView)">
            <f7-list-item title="Template Name"
                          :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                <template #media>
                    <f7-icon f7="app_fill"></f7-icon>
                </template>
                <template #footer>
                    <div class="scheduled-detail">
                        <span>Frequency</span>
                        <span v-if="currentStatus === ScheduledTransactionStatus.Active">·</span>
                        <span v-if="currentStatus === ScheduledTransactionStatus.Active">Next execution: YYYY-MM-DD</span>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-top-half" v-if="!loading && !visibleTemplates.length && (!supportCalendar || !calendarView)">
            <f7-list-item :title="tt('No available scheduled transactions')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers sortable class="margin-top-half scheduled-transaction-list"
                 :sortable-enabled="sortable"
                 v-if="!loading && (!supportCalendar || !calendarView)"
                 @sortable:sort="onSort">
            <f7-list-item swipeout
                          :class="{ 'actual-first-child': template.id === firstShowingId, 'actual-last-child': template.id === lastShowingId }"
                          :id="getTemplateDomId(template)"
                          :title="template.name"
                          :key="template.id"
                          v-for="template in templates"
                          v-show="showHidden || !template.hidden"
                          @taphold="setSortable()">
                <template #media>
                    <f7-icon f7="clock">
                        <f7-badge color="gray" class="right-bottom-icon" v-if="template.hidden">
                            <f7-icon f7="eye_slash_fill"></f7-icon>
                        </f7-badge>
                    </f7-icon>
                </template>
                <template #footer>
                    <div class="scheduled-detail">
                        <span>{{ getDisplayFrequency(template) }}</span>
                        <span v-if="template.nextScheduledTime">·</span>
                        <span v-if="template.nextScheduledTime">{{ getDisplayNextScheduledTime(template) }}</span>
                    </div>
                </template>
                <f7-swipeout-actions :left="textDirection === TextDirection.LTR"
                                     :right="textDirection === TextDirection.RTL"
                                     v-if="sortable">
                    <f7-swipeout-button class="padding-horizontal" overswipe close
                                        :aria-label="template.hidden ? tt('Show') : tt('Hide')"
                                        :color="template.hidden ? 'blue' : 'gray'"
                                        @click="hide(template, !template.hidden)">
                        <f7-icon :f7="template.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
                <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                     :right="textDirection === TextDirection.LTR"
                                     v-if="!sortable">
                    <f7-swipeout-button color="orange" close :text="tt('Edit')" @click="edit(template)"></f7-swipeout-button>
                    <f7-swipeout-button color="red" class="padding-horizontal" :aria-label="tt('Delete')" @click="remove(template, false)">
                        <f7-icon f7="trash"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>

        <f7-toolbar tabbar bottom :class="{ 'compact-tabbar': true, 'toolbar-item-auto-size': true, 'disabled': loading }">
            <f7-link :class="{ 'disabled': !supportCalendar || !calendarView || currentYearMonth <= minYearMonth }"
                     :aria-label="tt('Previous Period')" @click="shiftMonth(-1)">
                <f7-icon class="icon-with-direction" f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': !supportCalendar || !calendarView }" popover-open=".scheduled-month-popover-menu">
                <span :class="{ 'tabbar-item-changed': currentYearMonth !== minYearMonth }">{{ displayMonth }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': !supportCalendar || !calendarView }" :aria-label="tt('Next Period')" @click="shiftMonth(1)">
                <f7-icon class="icon-with-direction" f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" :class="{ 'disabled': !supportCalendar }" @click="switchView(false)">
                <span :class="{ 'tabbar-item-changed': !supportCalendar || !calendarView }">{{ tt('List View') }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" :class="{ 'disabled': !supportCalendar }" @click="switchView(true)">
                <span :class="{ 'tabbar-item-changed': supportCalendar && calendarView }">{{ tt('Calendar View') }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-popover class="scheduled-month-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt('This Month')"
                              :class="{ 'list-item-selected': currentYearMonth === minYearMonth }"
                              @click="changeCustomMonthDateFilter(getCurrentDateTime().toGregorianCalendarYear0BasedMonth())">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="currentYearMonth === minYearMonth"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt('Custom Date')"
                              @click="showCustomMonthSheet = true">
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <month-selection-sheet :title="tt('Select Month')"
                               :model-value="queryMonth"
                               v-model:show="showCustomMonthSheet"
                               @update:modelValue="changeCustomMonthDateFilter">
        </month-selection-sheet>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !templates || templates.length < 2 }" @click="setSortable()">{{ tt('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Scheduled Transactions') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Scheduled Transactions') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this template?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(templateToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import ScheduledTransactionCalendar from '@/components/common/ScheduledTransactionCalendar.vue';
import { useScheduledTransactionListPageBase } from '@/views/base/schedules/ScheduledTransactionListPageBase.ts';

import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted } from '@/lib/ui/mobile.ts';

import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';

import { TextDirection } from '@/core/text.ts';
import type { Year0BasedMonth } from '@/core/datetime.ts';
import { TemplateType, ScheduledTransactionStatus } from '@/core/template.ts';
import type { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    getCurrentDateTime,
    getYearMonthFirstUnixTime,
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

import {
    getFirstShowingId,
    getLastShowingId
} from '@/lib/template.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, getCurrentLanguageTextDirection } = useI18n();
const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();

const transactionTemplatesStore = useTransactionTemplatesStore();

const {
    loading,
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
    displayStatusName,
    displayMonth,
    loadScheduledTransactions,
    loadScheduledTransactionCalendar,
    moveScheduledTransaction,
    saveScheduledTransactionDisplayOrders,
    getDisplayFrequency,
    getDisplayNextScheduledTime
} = useScheduledTransactionListPageBase();

const loadingCalendar = ref<boolean>(false);
const loadingError = ref<unknown | null>(null);
const sortable = ref<boolean>(false);
const templateToDelete = ref<TransactionTemplate | null>(null);
const showMoreActionSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);
const showCustomMonthSheet = ref<boolean>(false);
const displayOrderSaving = ref<boolean>(false);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const firstShowingId = computed<string | null>(() => getFirstShowingId(templates.value, showHidden.value));
const lastShowingId = computed<string | null>(() => getLastShowingId(templates.value, showHidden.value));

function getTemplateDomId(template: TransactionTemplate): string {
    return 'scheduled_template_' + template.id;
}

function switchStatus(value: ScheduledTransactionStatus): void {
    if (currentStatus.value === value || displayOrderModified.value) {
        return;
    }

    currentStatus.value = value;
    calendarView.value = false;
    templates.value = [];
    occurrences.value = [];
    init(false);
}

function switchView(value: boolean): void {
    if (calendarView.value === value || loading.value || sortable.value || displayOrderModified.value || (value && !supportCalendar.value)) {
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
        showToast('Date is too early');
        return;
    }

    showCustomMonthSheet.value = false;

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

function init(firstVisit: boolean): void {
    loading.value = true;

    loadScheduledTransactions(false).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            if (firstVisit) {
                loadingError.value = error;
            }

            showToast(error.message || error);
        }
    });
}

function reload(done?: () => void): Promise<void> {
    if (sortable.value) {
        done?.();
        return Promise.resolve();
    }

    const force = !!done;

    if (currentYearMonth.value < minYearMonth.value) {
        currentYearMonth.value = minYearMonth.value;
    }

    return loadScheduledTransactions(force).then(() => {
        if (supportCalendar.value && calendarView.value) {
            return loadScheduledTransactionCalendar();
        }

        return;
    }).then(() => {
        done?.();

        if (force) {
            showToast('Template list has been updated');
        }
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function loadCalendar(): Promise<void> {
    loadingCalendar.value = true;

    return loadScheduledTransactionCalendar().then(() => {
        loadingCalendar.value = false;
    }).catch(error => {
        loadingCalendar.value = false;

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function edit(template: TransactionTemplate): void {
    props.f7router.navigate(`/template/edit?id=${template.id}&templateType=${TemplateType.Schedule.type}`);
}

function hide(template: TransactionTemplate, hidden: boolean): void {
    showLoading();

    transactionTemplatesStore.hideTemplate({
        template: template,
        hidden: hidden
    }).then(() => {
        template.hidden = hidden;
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function remove(template: TransactionTemplate | null, confirm: boolean): void {
    if (!template) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        templateToDelete.value = template;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    templateToDelete.value = null;
    showLoading();

    transactionTemplatesStore.deleteTemplate({
        template: template,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getTemplateDomId(template), done);
        }
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function setSortable(): void {
    if (sortable.value) {
        return;
    }

    calendarView.value = false;
    showHidden.value = true;
    sortable.value = true;
    displayOrderModified.value = false;
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        showHidden.value = false;
        sortable.value = false;
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    saveScheduledTransactionDisplayOrders().then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        showHidden.value = false;
        sortable.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function cancelSort(): void {
    if (!displayOrderModified.value) {
        showHidden.value = false;
        sortable.value = false;
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    loadScheduledTransactions(false).then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        showHidden.value = false;
        sortable.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onSort(event: { from: number; to: number }): void {
    moveScheduledTransaction(event.from, event.to);
}

function onPageAfterIn(): void {
    if (!loading.value) {
        reload();
    }

    routeBackOnError(props.f7router, loadingError);
}

init(true);
</script>

<style>
.scheduled-transaction-list {
    --f7-list-item-footer-font-size: var(--ebk-large-footer-font-size);

    .scheduled-detail > span {
        unicode-bidi: isolate;
        margin-inline-end: 4px;
    }
}

.scheduled-transaction-list .item-footer {
    padding-top: 4px;
}

.scheduled-transaction-calendar-container {
    .scheduled-transaction-calendar {
        font-size: var(--dp-font-size);

        .scheduled-calendar-grid {
            overflow: hidden;
            border: 1px solid var(--scheduled-calendar-border);
            border-radius: var(--f7-block-inset-border-radius);

            .scheduled-calendar-weekday,
            .scheduled-calendar-day {
                &:nth-child(7n) {
                    border-inline-end: 0;
                }

                &:nth-last-child(-n + 7) {
                    border-bottom: 0;
                }
            }

            .scheduled-calendar-day {
                min-height: 74px;

                .scheduled-calendar-alternate-date {
                    font-size: var(--ebk-transaction-calendar-alternate-date-font-size);
                }

                .scheduled-calendar-event {
                    font-size: var(--ebk-transaction-calendar-amount-font-size);
                }
            }
        }
    }
}
</style>
