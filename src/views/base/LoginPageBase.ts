import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { AuthResponse } from '@/models/auth_response.ts';

import { getLoginPageTips } from '@/lib/server_settings.ts';
import { getVersion } from '@/lib/version.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export function useLoginPageBase() {
    const { getServerTipContent, setLanguage } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const version = `v${getVersion()}`;

    const username = ref<string>('');
    const password = ref<string>('');
    const passcode = ref<string>('');
    const backupCode = ref<string>('');
    const tempToken = ref<string>('');
    const twoFAVerifyType = ref<string>('passcode');

    const logining = ref<boolean>(false);
    const verifying = ref<boolean>(false);

    const inputIsEmpty = computed<boolean>(() => !username.value || !password.value);
    const twoFAInputIsEmpty = computed<boolean>(() => {
        if (twoFAVerifyType.value === 'backupcode') {
            return !backupCode.value;
        } else {
            return !passcode.value;
        }
    });

    const tips = computed<string>(() => getServerTipContent(getLoginPageTips()));

    function changeLanguage(locale: string): void {
        const localeDefaultSettings = setLanguage(locale);
        settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
    }

    function doAfterLogin(authResponse: AuthResponse): void {
        if (authResponse.user) {
            const localeDefaultSettings = setLanguage(authResponse.user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(authResponse.user.expenseAmountColor, authResponse.user.incomeAmountColor);
        }

        if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
            exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
        }

        if (authResponse.notificationContent) {
            rootStore.setNotificationContent(authResponse.notificationContent);
        }
    }

    return {
        // constants
        version,
        // states
        username,
        password,
        passcode,
        backupCode,
        tempToken,
        twoFAVerifyType,
        logining,
        verifying,
        inputIsEmpty,
        twoFAInputIsEmpty,
        tips,
        // functions
        changeLanguage,
        doAfterLogin
    }
}
