import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTokensStore } from '@/stores/token.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { isObject, isString } from '@/lib/common.ts';
import { isWebAuthnSupported } from '@/lib/webauthn.ts';
import { hasWebAuthnConfig } from '@/lib/userstate.ts';
import { updateMapCacheExpiration, updateTransactionPictureCacheExpiration } from '@/lib/cache.ts';
import { getClientDisplayVersion } from '@/lib/version.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

export function useUnlockPageBase() {
    const { setLanguage, initLocale } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const tokensStore = useTokensStore();
    const transactionsStore = useTransactionsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const version: string = `${getClientDisplayVersion()}`;

    const pinCode = ref<string>('');

    const isWebAuthnAvailable = computed<boolean>(() => {
        return settingsStore.appSettings.applicationLockWebAuthn
            && hasWebAuthnConfig()
            && isWebAuthnSupported();
    });

    function isPinCodeValid(pinCode: string): boolean {
        return !!pinCode && pinCode.length === 6;
    }

    function doAfterUnlocked(): Promise<string> {
        transactionsStore.initTransactionDraft();

        return new Promise<string>((resolve, reject) => {
            let unlockHint: string = '';

            tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                const languagePromise = response.user ? setLanguage(response.user.language) : Promise.reject();

                languagePromise.then(localeDefaultSettings => {
                    settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }).catch(error => {
                    if (isObject(error) && 'message' in error && isString(error.message)) {
                        unlockHint = error.message;
                    } else if (isString(error)) {
                        unlockHint = error;
                    }
                }).finally(() => {
                    if (response.user) {
                        setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
                    }

                    updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
                    updateTransactionPictureCacheExpiration(settingsStore.appSettings.transactionPictureCacheExpiration);
                    exchangeRatesStore.removeExpiredExchangeRates(true);
                    exchangeRatesStore.autoUpdateExchangeRatesData();

                    if (response.notificationContent) {
                        rootStore.setNotificationContent(response.notificationContent);
                    }

                    resolve(unlockHint);
                });
            }).catch(error => {
                logger.error('failed to refresh token after unlocking', error);
                reject(error);
            });
        });
    }

    function doRelogin(): Promise<void> {
        rootStore.forceLogout();
        settingsStore.clearAppSettings();

        return initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone).then(localeDefaultSettings => {
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }).finally(() => {
            setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);
        });
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
        doAfterUnlocked,
        doRelogin
    };
}
