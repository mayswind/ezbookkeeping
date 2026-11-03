import { ref, computed } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';
import { useRootStore } from '@/stores/index.ts';

import { parseQuickSwitchCurrencies } from '@/lib/currency.ts';

export function useCurrencyQuickSwitchBase() {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();
    const rootStore = useRootStore();

    const switchingCurrency = ref<string | null>(null);

    const quickSwitchCurrencies = computed<string[]>({
        get: () => parseQuickSwitchCurrencies(settingsStore.appSettings.quickSwitchCurrenciesInHomePage),
        set: (value: string[]) => settingsStore.setQuickSwitchCurrenciesInHomePage(parseQuickSwitchCurrencies(value.join(',')).join(','))
    });

    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    // the chip row only makes sense when there are at least two currencies to switch between
    const showQuickSwitch = computed<boolean>(() => quickSwitchCurrencies.value.length >= 2);

    // changes the user's default currency, every overview amount is converted from it on the client side
    function switchCurrency(currency: string): Promise<unknown> {
        if (switchingCurrency.value || currency === defaultCurrency.value) {
            return Promise.resolve();
        }

        switchingCurrency.value = currency;

        return rootStore.updateUserProfile({ defaultCurrency: currency })
            // the profile update clears the cached exchange rates when the default currency changes, so reload them
            .then(() => exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false }))
            .finally(() => {
                switchingCurrency.value = null;
            });
    }

    return {
        // states
        switchingCurrency,
        // computed states
        quickSwitchCurrencies,
        defaultCurrency,
        showQuickSwitch,
        // functions
        switchCurrency
    };
}
