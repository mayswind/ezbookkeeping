import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import type { BigDecimal, HiddenAmount, BigDecimalWithSuffix } from '@/core/numeral.ts';

import { Account } from '@/models/account.ts';

export function useAssetSummaryWidgetBase(scene: 'overview' | 'accountList') {
    const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();

    const showAmountInHomePage = computed<boolean>({
        get: () => settingsStore.appSettings.showAmountInHomePage,
        set: (value) => settingsStore.setShowAmountInHomePage(value)
    });

    const showAccountBalance = computed<boolean>({
        get: () => settingsStore.appSettings.showAccountBalance,
        set: (value) => settingsStore.setShowAccountBalance(value)
    });

    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const allAccounts = computed<Account[]>(() => accountsStore.allAccounts);

    const netAssets = computed<string>(() => {
        if (scene === 'overview') {
            const netAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getNetAssets(showAmountInHomePage.value, settingsStore.appSettings.overviewAccountFilterInHomePage);
            return formatAmountToLocalizedNumeralsWithCurrency(netAssets, defaultCurrency.value);
        } else if (scene === 'accountList') {
            const netAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getNetAssets(showAccountBalance.value, settingsStore.appSettings.totalAmountExcludeAccountIds);
            return formatAmountToLocalizedNumeralsWithCurrency(netAssets, defaultCurrency.value);
        } else {
            return '';
        }
    });

    const totalAssets = computed<string>(() => {
        if (scene === 'overview') {
            const totalAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalAssets(showAmountInHomePage.value, settingsStore.appSettings.overviewAccountFilterInHomePage);
            return formatAmountToLocalizedNumeralsWithCurrency(totalAssets, defaultCurrency.value);
        } else if (scene === 'accountList') {
            const totalAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalAssets(showAccountBalance.value, settingsStore.appSettings.totalAmountExcludeAccountIds);
            return formatAmountToLocalizedNumeralsWithCurrency(totalAssets, defaultCurrency.value);
        } else {
            return '';
        }
    });

    const totalLiabilities = computed<string>(() => {
        if (scene === 'overview') {
            const totalLiabilities: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalLiabilities(showAmountInHomePage.value, settingsStore.appSettings.overviewAccountFilterInHomePage);
            return formatAmountToLocalizedNumeralsWithCurrency(totalLiabilities, defaultCurrency.value);
        } else if (scene === 'accountList') {
            const totalLiabilities: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalLiabilities(showAccountBalance.value, settingsStore.appSettings.totalAmountExcludeAccountIds);
            return formatAmountToLocalizedNumeralsWithCurrency(totalLiabilities, defaultCurrency.value);
        } else {
            return '';
        }
    });

    return {
        // computed states
        showAmountInHomePage,
        showAccountBalance,
        defaultCurrency,
        allAccounts,
        netAssets,
        totalAssets,
        totalLiabilities
    };
}
