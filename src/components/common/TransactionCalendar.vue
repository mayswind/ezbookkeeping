<template>
    <vue-date-picker inline auto-apply
                     model-type="yyyy-MM-dd"
                     :class="`transaction-calendar ${alternateDates ? 'transaction-calendar-with-alternate-date' : ''} ${calendarClass}`"
                     :config="{ noSwipe: true, monthChangeOnArrows: false, monthChangeOnScroll: false }"
                     :time-config="{ enableTimePicker: false }"
                     :input-attrs="{ clearable: false }"
                     :readonly="readonly"
                     :dark="isDarkMode"
                     :day-names="dayNames"
                     :week-start="firstDayOfWeek"
                     :min-date="minDate"
                     :max-date="maxDate"
                     :disabled-dates="noTransactionInMonthDay"
                     :prevent-min-max-navigation="true"
                     :hide-offset-dates="true"
                     :hide-month-year-select="true"
                     v-model="dateTime">
        <template #day="{ day, date }">
            <div class="transaction-calendar-daily-amounts">
                <span :class="dayHasTransactionClass && hasVisibleAmount(day) ? dayHasTransactionClass : undefined">{{ getDisplayDay(date) }}</span>
                <span class="transaction-calendar-alternate-date" v-if="alternateDates && alternateDates[`${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`]">{{ alternateDates[`${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`] }}</span>
                <span class="transaction-calendar-daily-amount text-income" v-if="showAmount && showIncomeAmount && dailyTotalAmounts && dailyTotalAmounts[day] && dailyTotalAmounts[day].income && !dailyTotalAmounts[day].income.isZero()">{{ getDisplayMonthTotalAmount(dailyTotalAmounts[day].income, defaultCurrency, '', dailyTotalAmounts[day].incompleteIncome) }}</span>
                <span class="transaction-calendar-daily-amount text-expense" v-if="showAmount && showExpenseAmount && dailyTotalAmounts && dailyTotalAmounts[day] && dailyTotalAmounts[day].expense && !dailyTotalAmounts[day].expense.isZero()">{{ getDisplayMonthTotalAmount(dailyTotalAmounts[day].expense, defaultCurrency, '', dailyTotalAmounts[day].incompleteExpense) }}</span>
                <span class="transaction-calendar-daily-amount" v-if="!showAmount">
                    <span class="transaction-calendar-daily-amount-dot text-income" v-if="showIncomeAmount && dailyTotalAmounts && dailyTotalAmounts[day] && dailyTotalAmounts[day].income && !dailyTotalAmounts[day].income.isZero()">●</span>
                    <span class="transaction-calendar-daily-amount-dot text-expense" style="margin-inline-start: 2px" v-if="showExpenseAmount && dailyTotalAmounts && dailyTotalAmounts[day] && dailyTotalAmounts[day].expense && !dailyTotalAmounts[day].expense.isZero()">●</span>
                </span>
            </div>
        </template>
    </vue-date-picker>
</template>

<script setup lang="ts">
import { computed, } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import type { TransactionTotalAmount } from '@/stores/transaction.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import type { CalendarAlternateDate, TextualYearMonthDay, WeekDayValue } from '@/core/datetime.ts';
import { INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getYearMonthDayDateTime } from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue: TextualYearMonthDay | '';
    isDarkMode: boolean;
    defaultCurrency: string | false;
    minDate: Date;
    maxDate: Date;
    weekDayNameType?: 'long' | 'short';
    dailyTotalAmounts?: Record<string, TransactionTotalAmount>;
    showAmount?: boolean;
    showIncomeAmount?: boolean;
    showExpenseAmount?: boolean;
    showAlternateDate?: boolean;
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
    getCalendarDisplayDayOfMonthFromDateTime,
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
    if (!props.showAlternateDate) {
        return undefined;
    }

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
    return !hasVisibleAmount(date.getDate());
}

function hasVisibleAmount(day: number): boolean {
    const dailyTotalAmount = props.dailyTotalAmounts?.[day];

    if (!dailyTotalAmount) {
        return false;
    }

    return !!(props.showIncomeAmount && dailyTotalAmount.income && !dailyTotalAmount.income.isZero()) || !!(props.showExpenseAmount && dailyTotalAmount.expense && !dailyTotalAmount.expense.isZero());
}

function getDisplayMonthTotalAmount(amount: BigDecimal, currency: string | false, symbol: string, incomplete: boolean): string {
    const displayAmount = formatAmountToLocalizedNumeralsWithCurrency(amount, currency);
    return symbol + displayAmount + (incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
}

function getDisplayDay(date: Date): string {
    return getCalendarDisplayDayOfMonthFromDateTime(getYearMonthDayDateTime(date.getFullYear(), date.getMonth() + 1, date.getDate()));
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
    opacity: 0.6;
}

.dp--cell-disabled .transaction-calendar-alternate-date {
    opacity: 0.8;
}

.dp--main.transaction-calendar .dp--calendar .dp--calendar-row > .dp--calendar-item .transaction-calendar-daily-amounts > span.transaction-calendar-alternate-date,
.dp--main.transaction-calendar .dp--calendar .dp--calendar-row > .dp--calendar-item .transaction-calendar-daily-amounts > span.transaction-calendar-daily-amount {
    display: block;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>
