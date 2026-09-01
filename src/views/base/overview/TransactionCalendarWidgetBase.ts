import { ref, computed } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import type { TransactionTotalAmount } from '@/stores/transaction.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { DateRange, KnownDateTimeFormat } from '@/core/datetime.ts';
import type { TextualYearMonthDay } from '@/core/datetime.ts';
import { TransactionType } from '@/core/transaction.ts';

import {
    getLocalDatetimeFromUnixTime,
    getSameDateTimeWithBrowserTimezone,
    getCurrentDateTime,
    parseDateTimeFromUnixTime,
    parseDateTimeFromKnownDateTimeFormat,
    getUnixTimeBeforeUnixTime,
    getUnixTimeAfterUnixTime,
    getValidMonthDayOrCurrentDayShortDate
} from '@/lib/datetime.ts';

export interface CommonTransactionCalendarWidgetProps {
    transactionTypes: number[];
}

export function useTransactionCalendarWidgetBase(props: CommonTransactionCalendarWidgetProps) {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const overviewStore = useOverviewStore();

    const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const showIncome = computed<boolean>(() => props.transactionTypes.includes(TransactionType.Income));
    const showExpense = computed<boolean>(() => props.transactionTypes.includes(TransactionType.Expense));
    const dailyTotalAmounts = computed<Record<string, TransactionTotalAmount>>(() => overviewStore.currentMonthTransactionDailyTotalAmounts);

    const currentCalendarDate = ref<TextualYearMonthDay>(getValidMonthDayOrCurrentDayShortDate(overviewStore.transactionDataRange.thisMonth.startTime, getCurrentDateTime().getGregorianCalendarYearDashMonthDashDay()));
    const transactionCalendarMinDate = computed<Date>(() => getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTime(overviewStore.transactionDataRange.thisMonth.startTime)).getUnixTime()));
    const transactionCalendarMaxDate = computed<Date>(() => getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTime(overviewStore.transactionDataRange.thisMonth.endTime)).getUnixTime()));

    function getTransactionListUrl(date: TextualYearMonthDay): string {
        const dateTime = parseDateTimeFromKnownDateTimeFormat(date, KnownDateTimeFormat.DefaultDate);

        if (!dateTime) {
            return '';
        }

        const minTime = dateTime.getUnixTime();
        const maxTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(minTime, 1, 'days'), 1, 'seconds');
        const type = props.transactionTypes.length === 1 ? props.transactionTypes[0] as TransactionType : undefined;
        return `/transaction/list?${overviewStore.getTransactionListPageParams({
            type: type,
            dateType: DateRange.Custom.type,
            minTime,
            maxTime
        })}`;
    }

    return {
        // computed states
        currentCalendarDate,
        showAmountInHomePage,
        showIncome,
        showExpense,
        defaultCurrency,
        dailyTotalAmounts,
        transactionCalendarMinDate,
        transactionCalendarMaxDate,
        // functions
        getTransactionListUrl
    };
}
