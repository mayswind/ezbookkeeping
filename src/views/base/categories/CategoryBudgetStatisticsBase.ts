import { computed, ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { KeywordMatchMode } from '@/core/text.ts';
import type { TextualYearMonth, Year1BasedMonth } from '@/core/datetime.ts';

import type { TransactionStatisticResponse } from '@/models/transaction.ts';

import {
    getCurrentDateTime,
    getYear0BasedMonthObjectFromString,
    getYearMonthDayDateTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime
} from '@/lib/datetime.ts';
import services from '@/lib/services.ts';

export function useCategoryBudgetStatisticsBase(initialMonth?: string) {
    const { formatDateTimeToGregorianLikeLongYearMonth } = useI18n();
    const now = getCurrentDateTime();
    const parsedInitialMonth = initialMonth
        ? getYear0BasedMonthObjectFromString(initialMonth as TextualYearMonth)
        : null;

    const selectedBudgetMonth = ref<Year1BasedMonth>({
        year: parsedInitialMonth?.year ?? now.getGregorianCalendarYear(),
        month1base: parsedInitialMonth ? parsedInitialMonth.month0base + 1 : now.getGregorianCalendarMonth()
    });
    const budgetStatistics = ref<TransactionStatisticResponse | null>(null);
    const budgetStatisticsLoading = ref<boolean>(false);

    let requestSequence = 0;

    const selectedBudgetMonthText = computed<string>(() => formatDateTimeToGregorianLikeLongYearMonth(
        getYearMonthDayDateTime(selectedBudgetMonth.value.year, selectedBudgetMonth.value.month1base, 1)
    ));

    const selectedBudgetMonthValue = computed<string>(() => `${selectedBudgetMonth.value.year}-${selectedBudgetMonth.value.month1base.toString().padStart(2, '0')}`);

    function changeBudgetMonth(offset: number): Promise<TransactionStatisticResponse> {
        const date = getYearMonthDayDateTime(selectedBudgetMonth.value.year, selectedBudgetMonth.value.month1base, 1).add(offset, 'months');

        selectedBudgetMonth.value = {
            year: date.getGregorianCalendarYear(),
            month1base: date.getGregorianCalendarMonth()
        };

        return loadBudgetStatistics();
    }

    function loadBudgetStatistics(): Promise<TransactionStatisticResponse> {
        const currentRequestSequence = ++requestSequence;
        budgetStatisticsLoading.value = true;

        return new Promise((resolve, reject) => {
            services.getTransactionStatistics({
                startTime: getYearMonthFirstUnixTime(selectedBudgetMonth.value),
                endTime: getYearMonthLastUnixTime(selectedBudgetMonth.value),
                tagFilter: '',
                keyword: '',
                matchMode: KeywordMatchMode.Default.type,
                useTransactionTimezone: false
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (currentRequestSequence === requestSequence) {
                        budgetStatisticsLoading.value = false;
                    }

                    reject({ message: 'Unable to retrieve transaction statistics' });
                    return;
                }

                if (currentRequestSequence === requestSequence) {
                    budgetStatistics.value = data.result;
                    budgetStatisticsLoading.value = false;
                }

                resolve(data.result);
            }).catch(error => {
                if (currentRequestSequence === requestSequence) {
                    budgetStatisticsLoading.value = false;
                }

                reject(error);
            });
        });
    }

    return {
        selectedBudgetMonth,
        selectedBudgetMonthText,
        selectedBudgetMonthValue,
        budgetStatistics,
        budgetStatisticsLoading,
        changeBudgetMonth,
        loadBudgetStatistics
    };
}
