import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type {
    LatestExchangeRate,
    LatestExchangeRateResponse,
    LocalizedLatestExchangeRate
} from '@/models/exchange_rate.ts';

import { getExchangedAmountByRate } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

export function useExchangeRatesPageBase() {
    const { getAllDisplayExchangeRates, formatDateTimeToLongDate, parseAmountFromWesternArabicNumerals } = useI18n();

    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const baseCurrency = ref<string>(userStore.currentUserDefaultCurrency);
    const baseAmount = ref<number>(100);

    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const exchangeRatesData = computed<LatestExchangeRateResponse | undefined>(() => exchangeRatesStore.latestExchangeRates.data);
    const isUserCustomExchangeRates = computed<boolean>(() => exchangeRatesStore.isUserCustomExchangeRates);

    const exchangeRatesDataUpdateTime = computed<string>(() => {
        if (!exchangeRatesStore.exchangeRatesLastUpdateTime) {
            return '';
        }

        const exchangeRatesLastUpdateTime = parseDateTimeFromUnixTime(exchangeRatesStore.exchangeRatesLastUpdateTime);
        return formatDateTimeToLongDate(exchangeRatesLastUpdateTime);
    });

    const availableExchangeRates = computed<LocalizedLatestExchangeRate[]>(() => {
        return getAllDisplayExchangeRates(exchangeRatesData.value);
    });

    function getConvertedAmount(baseAmount: number | '', fromExchangeRate?: LatestExchangeRate | LocalizedLatestExchangeRate, toExchangeRate?: LatestExchangeRate | LocalizedLatestExchangeRate): number | '' | null {
        if (!fromExchangeRate || !toExchangeRate) {
            return '';
        }

        if (baseAmount === '') {
            return 0;
        }

        return getExchangedAmountByRate(baseAmount as number, fromExchangeRate.rate, toExchangeRate.rate);
    }

    function setAsBaseline(currency: string, amount: string): void {
        baseCurrency.value = currency;
        baseAmount.value = parseAmountFromWesternArabicNumerals(amount);
    }

    return {
        // states
        baseCurrency,
        baseAmount,
        // computed states
        defaultCurrency,
        exchangeRatesData,
        isUserCustomExchangeRates,
        exchangeRatesDataUpdateTime,
        availableExchangeRates,
        // functions
        getConvertedAmount,
        setAsBaseline
    };
}
