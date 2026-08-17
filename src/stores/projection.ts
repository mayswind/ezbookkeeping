import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import type { TextualYearMonth } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import type { ProjectionTableData } from '@/core/projection.ts';
import type { TransactionStatisticTrendsResponseItem } from '@/models/transaction.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { getCurrentUnixTime, getGregorianCalendarYearAndMonthFromUnixTime } from '@/lib/datetime.ts';
import { buildProjectionTableData, getProjectionMonths } from '@/lib/projection.ts';
import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useProjectionsStore = defineStore('projections', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const projectionStartYearMonth = ref<TextualYearMonth | ''>('');
    const projectionEndYearMonth = ref<TextualYearMonth | ''>('');
    const projectionData = ref<TransactionStatisticTrendsResponseItem[]>([]);

    // A projection depends on the instant it was built at and on the set of scheduled templates, not
    // only on the period asked for, so it cannot be cached by period alone. The state is invalidated
    // when a scheduled template changes and when the day rolls over, and starts out invalid.
    const projectionStateInvalid = ref<boolean>(true);
    // the period the data in the store was loaded for, and the day it was loaded on
    const projectionDataPeriodKey = ref<string>('');
    const projectionDataDay = ref<number>(0);

    const projectionTableData = computed<ProjectionTableData>(() => {
        const months = getProjectionMonths(projectionStartYearMonth.value, projectionEndYearMonth.value, getCurrentYearMonthNumber());

        return buildProjectionTableData(projectionData.value, months, {
            categoriesMap: transactionCategoriesStore.allTransactionCategoriesMap,
            convertAmount: convertAmountToDefaultCurrency,
            currentYearMonth: getCurrentYearMonthNumber()
        });
    });

    // convertAmountToDefaultCurrency turns an amount of the specified account into the default
    // currency of the user, returning null when it cannot be converted, in the same way the
    // statistics store does it for the trends of real transactions
    function convertAmountToDefaultCurrency(amount: string, accountId: string): BigDecimal | null {
        const defaultCurrency = userStore.currentUserDefaultCurrency;
        const account = accountsStore.allAccountsMap[accountId];

        if (!account) {
            return null;
        }

        if (account.currency === defaultCurrency) {
            return parseBigDecimal(amount);
        }

        const exchangedAmount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(amount), account.currency, defaultCurrency);

        return exchangedAmount ? exchangedAmount.truncate() : null;
    }

    function getCurrentYearMonthNumber(): number {
        const yearMonth = getGregorianCalendarYearAndMonthFromUnixTime(getCurrentUnixTime());

        if (!yearMonth) {
            return 0;
        }

        const parts = yearMonth.split('-');

        return parseInt(parts[0]!, 10) * 100 + parseInt(parts[1]!, 10);
    }

    function getProjectionPeriodKey(): string {
        return `${projectionStartYearMonth.value}|${projectionEndYearMonth.value}`;
    }

    function getCurrentDayNumber(): number {
        return Math.floor(getCurrentUnixTime() / 86400);
    }

    function updateProjectionInvalidState(invalidState: boolean): void {
        projectionStateInvalid.value = invalidState;
    }

    function setProjectionPeriod(startYearMonth: TextualYearMonth | '', endYearMonth: TextualYearMonth | ''): void {
        if (projectionStartYearMonth.value === startYearMonth && projectionEndYearMonth.value === endYearMonth) {
            return;
        }

        projectionStartYearMonth.value = startYearMonth;
        projectionEndYearMonth.value = endYearMonth;
        projectionStateInvalid.value = true;
    }

    function resetProjections(): void {
        projectionStartYearMonth.value = '';
        projectionEndYearMonth.value = '';
        projectionData.value = [];
        projectionDataPeriodKey.value = '';
        projectionDataDay.value = 0;
        projectionStateInvalid.value = true;
    }

    // isProjectionDataUpToDate reports whether the data in the store can be reused: it must have been
    // loaded for the current period, on the current day, and no scheduled template may have changed
    // since
    function isProjectionDataUpToDate(): boolean {
        return !projectionStateInvalid.value
            && projectionDataPeriodKey.value === getProjectionPeriodKey()
            && projectionDataDay.value === getCurrentDayNumber();
    }

    function loadProjections({ force }: { force?: boolean }): Promise<TransactionStatisticTrendsResponseItem[]> {
        if (!force && isProjectionDataUpToDate()) {
            return Promise.resolve(projectionData.value);
        }

        return new Promise((resolve, reject) => {
            services.getTransactionProjections({
                startYearMonth: projectionStartYearMonth.value,
                endYearMonth: projectionEndYearMonth.value,
                useTransactionTimezone: settingsStore.appSettings.statistics.defaultTimezoneType === TimezoneTypeForStatistics.TransactionTimezone.type
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve transaction projections' });
                    return;
                }

                projectionData.value = data.result;
                projectionDataPeriodKey.value = getProjectionPeriodKey();
                projectionDataDay.value = getCurrentDayNumber();
                projectionStateInvalid.value = false;

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve transaction projections', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve transaction projections' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        projectionStartYearMonth,
        projectionEndYearMonth,
        projectionData,
        projectionStateInvalid,
        // computed states
        projectionTableData,
        // functions
        setProjectionPeriod,
        updateProjectionInvalidState,
        resetProjections,
        loadProjections
    };
});
