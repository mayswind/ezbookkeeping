import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';

import { type TypeAndDisplayName, itemAndIndex } from '@/core/base.ts';
import { ScheduledTemplateFrequencyType, ScheduledTransactionStatus } from '@/core/template.ts';
import { type TextualYearMonth, type Year0BasedMonth, type WeekDayValue, DateRangeScene } from '@/core/datetime.ts';

import type {
    TransactionTemplate,
    ScheduledTransactionOccurrenceResponse,
    TransactionTemplateNewDisplayOrderRequest
} from '@/models/transaction_template.ts';

import {
    isEquals,
    findDisplayNameByType
} from '@/lib/common.ts';
import {
    getYear0BasedMonthObjectFromString,
    getYear1BasedMonthObjectFromString,
    getCurrentDateTime,
    parseDateTimeFromUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange
} from '@/lib/datetime.ts';

export function useScheduledTransactionListPageBase() {
    const {
        tt,
        getMultiMonthAndDayLongNames,
        getMultiMonthdayShortNames,
        getMultiWeekdayLongNames,
        formatDateTimeToLongDate,
        formatDateRange
    } = useI18n();

    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionTemplatesStore = useTransactionTemplatesStore();

    const loading = ref<boolean>(false);
    const updating = ref<boolean>(false);
    const showHidden = ref<boolean>(false);
    const calendarView = ref<boolean>(false);
    const displayOrderModified = ref<boolean>(false);
    const minYearMonth = ref<TextualYearMonth>(getCurrentDateTime().getGregorianCalendarYearDashMonth());
    const currentYearMonth = ref<TextualYearMonth>(minYearMonth.value);
    const currentStatus = ref<ScheduledTransactionStatus>(ScheduledTransactionStatus.Active);
    const templates = ref<TransactionTemplate[]>([]);
    const occurrences = ref<ScheduledTransactionOccurrenceResponse[]>([]);

    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);

    const allStatuses = computed<TypeAndDisplayName[]>(() => [
        {
            type: ScheduledTransactionStatus.Active,
            displayName: tt('Active'),
        },
        {
            type: ScheduledTransactionStatus.Disabled,
            displayName: tt('pageSpecific.scheduledTransactionListPage.disabled')
        },
        {
            type: ScheduledTransactionStatus.Expired,
            displayName: tt('Expired')
        }
    ]);

    const supportCalendar = computed<boolean>(() => currentStatus.value === ScheduledTransactionStatus.Active);
    const visibleTemplates = computed<TransactionTemplate[]>(() => templates.value.filter(template => showHidden.value || !template.hidden));
    const queryMonth = computed<Year0BasedMonth>(() => getYear0BasedMonthObjectFromString(currentYearMonth.value)!);
    const displayStatusName = computed<string>(() => findDisplayNameByType(allStatuses.value, currentStatus.value) || '');
    const displayMonth = computed<string>(() => {
        if (!calendarView.value) {
            return tt('All');
        }

        const startTime = getYearMonthFirstUnixTime(currentYearMonth.value);
        const endTime = getYearMonthLastUnixTime(currentYearMonth.value);
        const dateType = getDateTypeByDateRange(startTime, endTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);

        return formatDateRange(dateType, startTime, endTime);
    });

    function loadScheduledTransactions(force?: boolean): Promise<void> {
        return Promise.all([
            accountsStore.loadAllAccounts({ force: false }),
            transactionTemplatesStore.loadAllScheduledTransactions({ status: currentStatus.value })
        ]).then(result => {
            if (force && isEquals(templates.value, result[1] || [])) {
                return Promise.reject({ message: 'Template list is up to date', isUpToDate: true });
            }

            templates.value = result[1] || [];
            return;
        });
    }

    function loadScheduledTransactionCalendar(): Promise<void> {
        if (!supportCalendar.value) {
            return Promise.resolve();
        }

        const yearMonth = getYear1BasedMonthObjectFromString(currentYearMonth.value);

        if (!yearMonth) {
            return Promise.reject('An error occurred');
        }

        return transactionTemplatesStore.loadScheduledTransactionsInMonth({
            year: yearMonth.year,
            month: yearMonth.month1base
        }).then(result => {
            occurrences.value = result;
        });
    }

    function moveScheduledTransaction(from: number, to: number): void {
        const template = templates.value[from];

        if (!template || !templates.value[to]) {
            return;
        }

        templates.value.splice(to, 0, templates.value.splice(from, 1)[0]!);
        displayOrderModified.value = true;
    }

    function saveScheduledTransactionDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: TransactionTemplateNewDisplayOrderRequest[] = [];

        for (const [template, index] of itemAndIndex(templates.value)) {
            newDisplayOrders.push({
                id: template.id,
                displayOrder: index + 1
            });
        }

        return transactionTemplatesStore.updateNormalTemplateDisplayOrders({
            newDisplayOrders: newDisplayOrders
        });
    }

    function getDisplayFrequency(template: TransactionTemplate): string {
        const values = (template.scheduledFrequency || '').split(',').filter(value => value).map(value => parseInt(value));

        switch (template.scheduledFrequencyType) {
            case ScheduledTemplateFrequencyType.Disabled.type:
                return tt('Disabled');
            case ScheduledTemplateFrequencyType.Daily.type:
                return tt('Daily');
            case ScheduledTemplateFrequencyType.EveryNDays.type:
                return tt('format.misc.everyNDays', {
                    n: values[0]
                });
            case ScheduledTemplateFrequencyType.Weekly.type:
                return tt('format.misc.everyMultiDaysOfWeek', {
                    days: getMultiWeekdayLongNames(values, firstDayOfWeek.value)
                });
            case ScheduledTemplateFrequencyType.Monthly.type:
                return tt('format.misc.everyMultiDaysOfMonth', {
                    days: getMultiMonthdayShortNames(values)
                });
            case ScheduledTemplateFrequencyType.Yearly.type:
                return tt('format.misc.everyMultiDaysOfYear', {
                    days: getMultiMonthAndDayLongNames(values)
                });
            default:
                return '';
        }
    }

    function getDisplayNextScheduledTime(template: TransactionTemplate): string {
        return template.nextScheduledTime ? tt('format.misc.nextExecutionTime', {
            time: formatDateTimeToLongDate(parseDateTimeFromUnixTime(template.nextScheduledTime))
        }) : '';
    }

    return {
        // states
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
        // computed states
        allStatuses,
        supportCalendar,
        visibleTemplates,
        queryMonth,
        displayStatusName,
        displayMonth,
        // functions
        loadScheduledTransactions,
        loadScheduledTransactionCalendar,
        moveScheduledTransaction,
        saveScheduledTransactionDisplayOrders,
        getDisplayFrequency,
        getDisplayNextScheduledTime
    };
}
