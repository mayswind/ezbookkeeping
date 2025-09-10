<template>
    <vue-date-picker inline auto-apply
                     model-type="yyyy-MM-dd"
                     :class="`transaction-calendar ${alternateDates ? 'transaction-calendar-with-alternate-date' : ''} ${calendarClass}`"
                     :config="{ noSwipe: true }"
                     :readonly="readonly"
                     :dark="isDarkMode"
                     :enable-time-picker="false"
                     :clearable="false"
                     :day-names="dayNames"
                     :week-start="firstDayOfWeek"
                     :min-date="minDate"
                     :max-date="maxDate"
                     :disabled-dates="noTransactionInMonthDay"
                     :prevent-min-max-navigation="true"
                     :hide-offset-dates="true"
                     :disable-month-year-select="true"
                     :month-change-on-scroll="false"
                     :month-change-on-arrows="false"
                     v-model="dateTime">
        <template #day="{ day, date }">
            <div class="transaction-calendar-daily-amounts">
                <span :class="dayHasTransactionClass && dailyTotalAmounts && dailyTotalAmounts[day] ? dayHasTransactionClass : undefined">{{ getDisplayDay(date) }}</span>
                <span class="transaction-calendar-alternate-date" v-if="alternateDates && alternateDates[`${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`]">{{ alternateDates[`${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`] }}</span>
                <span class="transaction-calendar-daily-amount text-income" v-if="dailyTotalAmounts && dailyTotalAmounts[day] && dailyTotalAmounts[day].income">{{ getDisplayMonthTotalAmount(dailyTotalAmounts[day].income, defaultCurrency, '', dailyTotalAmounts[day].incompleteIncome) }}</span>
                <span class="transaction-calendar-daily-amount text-expense" v-if="dailyTotalAmounts && dailyTotalAmounts[day] && dailyTotalAmounts[day].expense">{{ getDisplayMonthTotalAmount(dailyTotalAmounts[day].expense, defaultCurrency, '', dailyTotalAmounts[day].incompleteExpense) }}</span>
            </div>
        </template>
    </vue-date-picker>
</template>

<script setup lang="ts">
import { computed, } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import type { TransactionTotalAmount } from '@/stores/transaction.ts';

import type { CalendarAlternateDate, TextualYearMonthDay, WeekDayValue } from '@/core/datetime.ts';
import { INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import {
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getUnixTimeFromLocalDatetime,
    getActualUnixTimeForStore,
    getYearMonthDayDateTime,
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue: TextualYearMonthDay | '';
    isDarkMode: boolean;
    defaultCurrency: string | false;
    minDate: Date;
    maxDate: Date;
    weekDayNameType?: 'long' | 'short';
    dailyTotalAmounts?: Record<string, TransactionTotalAmount>;
    readonly?: boolean;
    calendarClass?: string;
    dayHasTransactionClass?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const {
    getAllLongWeekdayNames,
    getAllShortWeekdayNames,
    getCalendarDisplayDayOfMonthFromUnixTime,
    getCalendarAlternateDates,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const userStore = useUserStore();

const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(props.weekDayNameType === 'short' ? getAllShortWeekdayNames() : getAllLongWeekdayNames(), firstDayOfWeek.value));
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);

const dateTime = computed<TextualYearMonthDay | ''>({
    get: () => props.modelValue,
    set: (value: TextualYearMonthDay | '') => emit('update:modelValue', value)
});

const alternateDates = computed<Record<TextualYearMonthDay, string> | undefined>(() => {
    const yearMonthDay = props.modelValue ? props.modelValue.split('-') : null;

    if (!yearMonthDay || yearMonthDay.length !== 3) {
        return undefined;
    }

    const allDates: CalendarAlternateDate[] | undefined = getCalendarAlternateDates({ year: parseInt(yearMonthDay[0] as string), month1base: parseInt(yearMonthDay[1] as string) })

    if (!allDates) {
        return undefined;
    }

    const ret: Record<TextualYearMonthDay, string> = {};

    for (const alternateDate of allDates) {
        ret[`${alternateDate.year}-${alternateDate.month}-${alternateDate.day}`] = alternateDate.displayDate;
    }

    return ret;
});

function noTransactionInMonthDay(date: Date): boolean {
    const dateTime = parseDateTimeFromUnixTime(getActualUnixTimeForStore(getUnixTimeFromLocalDatetime(date), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    return !props.dailyTotalAmounts || !props.dailyTotalAmounts[dateTime.getGregorianCalendarDay()];
}

function getDisplayMonthTotalAmount(amount: number, currency: string | false, symbol: string, incomplete: boolean): string {
    const displayAmount = formatAmountToLocalizedNumeralsWithCurrency(amount, currency);
    return symbol + displayAmount + (incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
}

function getDisplayDay(date: Date): string {
    return getCalendarDisplayDayOfMonthFromUnixTime(getYearMonthDayDateTime(date.getFullYear(), date.getMonth() + 1, date.getDate()).getUnixTime());
}
</script>

<style>
.transaction-calendar-daily-amounts {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.transaction-calendar-alternate-date {
    margin-top: -3px;
    opacity: 0.5;
}

.dp__cell_disabled .transaction-calendar-alternate-date {
    opacity: 0.8;
}

.dp__main.transaction-calendar .dp__calendar .dp__calendar_row > .dp__calendar_item .transaction-calendar-daily-amounts > span.transaction-calendar-alternate-date,
.dp__main.transaction-calendar .dp__calendar .dp__calendar_row > .dp__calendar_item .transaction-calendar-daily-amounts > span.transaction-calendar-daily-amount {
    display: block;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>
