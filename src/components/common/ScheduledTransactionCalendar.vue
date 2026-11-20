<template>
    <div class="scheduled-transaction-calendar" :style="{ textAlign: dateAlign }">
        <div class="scheduled-calendar-grid" :aria-label="displayMonth">
            <div class="scheduled-calendar-weekday" :key="day" v-for="day in dayNames">{{ day }}</div>
            <div class="scheduled-calendar-day" :class="{ 'scheduled-calendar-not-available': !dataInDate.isCurrentMonth || dataInDate.date < today }"
                 :key="index" v-for="(dataInDate, index) in templatesByDates">
                <div class="scheduled-calendar-date"
                     :class="{ 'scheduled-calendar-today': dataInDate.date === today, 'scheduled-calendar-other-month': !dataInDate.isCurrentMonth }"
                     :aria-label="formatGregorianTextualYearMonthDayToLongDate(dataInDate.date)">
                    <div>{{ dataInDate.label }}</div>
                    <div class="scheduled-calendar-alternate-date" v-if="dataInDate.alternateDate">{{ dataInDate.alternateDate }}</div>
                </div>
                <div class="scheduled-calendar-skeleton" aria-hidden="true" v-if="loading && dataInDate.isCurrentMonth && dataInDate.date >= today"></div>
                <template v-else-if="!loading">
                    <button class="scheduled-calendar-event" type="button" :key="template.id"
                            :class="{ 'scheduled-calendar-event-income': template.type === TransactionType.Income, 'scheduled-calendar-event-expense': template.type === TransactionType.Expense }"
                            v-for="template in dataInDate.templates"
                            @click="emit('edit', template)">{{ template.name }}</button>
                </template>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { DateTime, TextualYearMonth, TextualYearMonthDay, WeekDayValue } from '@/core/datetime.ts';
import { TransactionType } from '@/core/transaction.ts';

import type { TransactionTemplate, ScheduledTransactionOccurrenceResponse } from '@/models/transaction_template.ts';

import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getCurrentDateTime, getYearMonthFirstUnixTime, parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

interface ScheduledTransactionCalendarDay {
    date: TextualYearMonthDay;
    label: string;
    alternateDate?: string;
    isCurrentMonth: boolean;
    templates: TransactionTemplate[];
}

const props = defineProps<{
    loading?: boolean;
    showAlternateDate?: boolean;
    showHidden?: boolean;
    weekDayNameType?: 'long' | 'short';
    dateAlign?: 'left' | 'center';
    modelValue: TextualYearMonth;
    templates: TransactionTemplate[];
    occurrences: ScheduledTransactionOccurrenceResponse[];
}>();

const emit = defineEmits<{
    (e: 'edit', template: TransactionTemplate): void;
}>();

const {
    getAllLongWeekdayNames,
    getAllShortWeekdayNames,
    getCalendarDisplayDayOfMonthFromDateTime,
    getCalendarAlternateDate,
    formatGregorianTextualYearMonthDayToLongDate,
    formatDateTimeToGregorianLikeLongYearMonth
} = useI18n();

const userStore = useUserStore();

const currentDateTime = ref<DateTime>(getCurrentDateTime());

const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(props.weekDayNameType === 'short' ? getAllShortWeekdayNames() : getAllLongWeekdayNames(), firstDayOfWeek.value));
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);

const today = computed<TextualYearMonthDay>(() => currentDateTime.value.getGregorianCalendarYearDashMonthDashDay());
const monthDate = computed<DateTime>(() => parseDateTimeFromUnixTime(getYearMonthFirstUnixTime(props.modelValue)));
const displayMonth = computed<string>(() => formatDateTimeToGregorianLikeLongYearMonth(monthDate.value));

const templatesByDateMap = computed<Record<TextualYearMonthDay, TransactionTemplate[]>>(() => {
    const templatesMap: Record<string, TransactionTemplate> = {};
    const ret: Record<TextualYearMonthDay, TransactionTemplate[]> = {};

    for (const template of props.templates) {
        templatesMap[template.id] = template;
    }

    for (const occurrence of props.occurrences) {
        const template = templatesMap[occurrence.templateId];

        if (!template || (!props.showHidden && template.hidden)) {
            continue;
        }

        const date = parseDateTimeFromUnixTime(occurrence.scheduledTime).getGregorianCalendarYearDashMonthDashDay();
        let templatesInDate = ret[date];

        if (!templatesInDate) {
            templatesInDate = [];
            ret[date] = templatesInDate;
        }

        templatesInDate.push(template);
    }

    return ret;
});

const templatesByDates = computed<ScheduledTransactionCalendarDay[]>(() => {
    const firstDay = monthDate.value;
    const offset = (firstDay.getWeekDay().type - firstDayOfWeek.value + 7) % 7;
    const count = firstDay.getMaxDayOfGregorianCalendarMonth();
    const cellCount = Math.ceil((offset + count) / 7) * 7;
    const ret: ScheduledTransactionCalendarDay[] = [];

    for (let index = 0; index < cellCount; index++) {
        const day = index - offset + 1;
        const isCurrentMonth = day >= 1 && day <= count;
        const dateTime = firstDay.add(index - offset, 'days');
        const date = dateTime.getGregorianCalendarYearDashMonthDashDay();

        ret.push({
            date: date,
            label: getCalendarDisplayDayOfMonthFromDateTime(dateTime),
            alternateDate: props.showAlternateDate ? getCalendarAlternateDate(dateTime.toGregorianCalendarYearMonthDay())?.displayDate : undefined,
            isCurrentMonth: isCurrentMonth,
            templates: isCurrentMonth ? templatesByDateMap.value[date] || [] : []
        });
    }

    return ret;
});
</script>

<style scoped>
.scheduled-transaction-calendar {
    --scheduled-calendar-border: color-mix(in srgb, currentColor 15%, transparent);

    .scheduled-calendar-grid {
        display: grid;
        grid-template-columns: repeat(7, minmax(0, 1fr));
        border-top: 1px solid var(--scheduled-calendar-border);
        border-inline-start: 1px solid var(--scheduled-calendar-border);
    }

    .scheduled-calendar-weekday,
    .scheduled-calendar-day {
        box-sizing: border-box;
        min-width: 0;
        border-bottom: 1px solid var(--scheduled-calendar-border);
        border-inline-end: 1px solid var(--scheduled-calendar-border);
        padding: 6px;
    }

    .scheduled-calendar-weekday {
        text-align: center;
        font-weight: 600;
    }

    .scheduled-calendar-day {
        &.scheduled-calendar-not-available {
            background: color-mix(in srgb, currentColor 3%, transparent);
        }
    }

    .scheduled-calendar-date {
        box-sizing: border-box;
        display: inline-block;
        max-width: 100%;
        padding: 4px 7px;
        border-radius: 4px;

        &.scheduled-calendar-other-month {
            opacity: 0.5;
        }

        &.scheduled-calendar-today {
            color: rgb(var(--ebk-primary-color));
            background: rgba(var(--ebk-primary-color), 0.15);
            font-weight: 600;
        }
    }

    .scheduled-calendar-alternate-date {
        font-weight: 400;
        opacity: 0.6;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .scheduled-calendar-event {
        box-sizing: border-box;
        display: block;
        width: 100%;
        font: inherit;
        line-height: 1.4;
        text-align: start;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        color: inherit;
        cursor: pointer;
        margin-top: 3px;
        padding: 2px 5px;
        border: 0;
        border-radius: 4px;
        background: rgba(var(--ebk-primary-color), 0.12);
    }

    .scheduled-calendar-skeleton {
        height: 22px;
        margin-top: 3px;
        border-radius: 4px;
        background: color-mix(in srgb, currentColor 12%, transparent);
        animation: scheduled-calendar-skeleton-loading 2s ease-in-out infinite;
    }
}

@keyframes scheduled-calendar-skeleton-loading {
    50% {
        opacity: 0.5;
    }
}
</style>
