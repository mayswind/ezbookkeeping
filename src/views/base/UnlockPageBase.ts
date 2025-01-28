import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTokensStore } from '@/stores/token.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { isWebAuthnSupported } from '@/lib/webauthn.ts';
import { hasWebAuthnConfig } from '@/lib/userstate.ts';
import { getVersion } from '@/lib/version.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export function useUnlockPageBase() {
    const { setLanguage, initLocale } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const tokensStore = useTokensStore();
    const transactionsStore = useTransactionsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const version: string = `v${getVersion()}`;

    const pinCode = ref<string>('');

    const isWebAuthnAvailable = computed<boolean>(() => {
        return settingsStore.appSettings.applicationLockWebAuthn
            && hasWebAuthnConfig()
            && isWebAuthnSupported();
    });

    function isPinCodeValid(pinCode: string): boolean {
        return !!pinCode && pinCode.length === 6;
    }

    function changeLanguage(locale: string): void {
        const localeDefaultSettings = setLanguage(locale);
        settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
    }

    function doAfterUnlocked(): void {
        transactionsStore.initTransactionDraft();
        tokensStore.refreshTokenAndRevokeOldToken().then(response => {
            if (response.user) {
                const localeDefaultSettings = setLanguage(response.user.language);
                settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
            }

            if (response.notificationContent) {
                rootStore.setNotificationContent(response.notificationContent);
            }
        });

        if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
            exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
        }
    }

    function doRelogin(): void {
        rootStore.forceLogout();
        settingsStore.clearAppSettings();

        const localeDefaultSettings = initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone);
        settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

        setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);
    }

    return {
        // constants
        version,
        // states
        pinCode,
        // computed states
        isWebAuthnAvailable,
        // methods
        isPinCodeValid,
        changeLanguage,
        doAfterUnlocked,
        doRelogin
    };
}
