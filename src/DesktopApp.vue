<template>
    <v-app>
        <router-view />
    </v-app>

    <v-snackbar class="cursor-pointer" color="notification-background" location="top"
                :multi-line="true" :timeout="-1" :close-on-content-click="true" v-model="showNotification">
        <v-tooltip activator="parent">{{ tt('Click to close') }}</v-tooltip>
        <div class="d-inline-flex">
            <img alt="logo" class="notification-logo" :src="APPLICATION_LOGO_PATH" />
            <span class="ms-2">{{ tt('global.app.title') }}</span>
        </div>
        <div>
            {{ currentNotificationContent }}
        </div>
    </v-snackbar>

    <v-snackbar-queue :total-visible="3" display-strategy="overflow" :close-text="tt('Close')" v-model="toastQueue">
        <template #actions="{ props }">
            <v-btn variant="text" color="primary" v-bind="props">{{ tt('Close') }}</v-btn>
        </template>
    </v-snackbar-queue>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';

import { useTheme } from 'vuetify';
import { register } from 'register-service-worker';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTokensStore } from '@/stores/token.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';

import { isObject, isString } from '@/lib/common.ts';
import { isProduction } from '@/lib/version.ts';
import { initMapProvider } from '@/lib/map/index.ts';
import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';
import { updateMapCacheExpiration, updateTransactionPictureCacheExpiration } from '@/lib/cache.ts';
import { getSystemTheme, setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

interface SnackbarMessage {
    text: string;
    style: string;
    onDismiss: (reason: string) => void;
}

const { tt, te, getCurrentLanguageInfo, setLanguage, initLocale } = useI18n();

const theme = useTheme();

const rootStore = useRootStore();
const desktopPageStore = useDesktopPageStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const tokensStore = useTokensStore();
const exchangeRatesStore = useExchangeRatesStore();

const initialRoutePath: string = (() => {
    if (!window.location.hash) {
        return '/';
    }

    const hash = window.location.hash;
    const hashIndex = hash.indexOf('#/');

    if (hashIndex < 0) {
        return '/';
    }

    const routePath = hash.substring(hashIndex + 1);
    const queryIndex = routePath.indexOf('?');

    if (queryIndex < 0) {
        return routePath;
    }

    return routePath.substring(0, queryIndex);
})();

const showNotification = ref<boolean>(false);

const toastQueue = computed<SnackbarMessage[]>({
    get: () => {
        const localizedToasts: SnackbarMessage[] = [];

        for (const toast of desktopPageStore.pendingToastQueue) {
            const snackbarMessage: SnackbarMessage = {
                text: '',
                style: 'z-index: 10000',
                onDismiss: () => {
                    desktopPageStore.removeShowingToast(toast.id);
                }
            };

            if (toast.type === 'message') {
                if (toast.options) {
                    snackbarMessage.text = tt(toast.message, toast.options);
                } else {
                    snackbarMessage.text = tt(toast.message);
                }
            } else if (toast.type === 'error') {
                if (isObject(toast.error) && 'message' in toast.error && toast.error.message) {
                    snackbarMessage.text = te(toast.error.message);
                } else if (isObject(toast.error) && 'error' in toast.error) {
                    snackbarMessage.text = te(toast.error);
                } else if (isString(toast.error)) {
                    snackbarMessage.text = te(toast.error);
                } else {
                    snackbarMessage.text = tt('An error occurred');
                }
            } else {
                snackbarMessage.text = tt('An error occurred');
            }

            localizedToasts.push(snackbarMessage);
        }

        return localizedToasts;
    },
    set: (newValue: SnackbarMessage[]) => {
        const oldCount = desktopPageStore.pendingToastQueue.length;
        const newCount = newValue.length;

        if (newCount < oldCount) {
            desktopPageStore.removePendingToast(oldCount - newCount);
        }
    }
});

const currentNotificationContent = computed<string | null>(() => rootStore.currentNotification);

onMounted(() => {
    document.addEventListener('DOMContentLoaded', () => {
        const languageInfo = getCurrentLanguageInfo();
        initMapProvider(languageInfo?.alternativeLanguageTag);
    });
});

watch(currentNotificationContent, (newValue) => {
    showNotification.value = !!newValue;
});

if (settingsStore.appSettings.theme === ThemeType.Light) {
    theme.change(ThemeType.Light);
} else if (settingsStore.appSettings.theme === ThemeType.Dark) {
    theme.change(ThemeType.Dark);
} else {
    theme.change(getSystemTheme());
}

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
    if (settingsStore.appSettings.theme === 'auto') {
        if (e.matches) {
            theme.change(ThemeType.Dark);
        } else {
            theme.change(ThemeType.Light);
        }
    }
});

initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone).then(localeDefaultSettings => {
    settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
}).catch(error => {
    logger.error('failed to initialize locale', error);

    if (!error.processed) {
        desktopPageStore.addErrorToast(error);
    }
});

setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);

if (isUserLogined() && initialRoutePath !== '/verify_email' && initialRoutePath !== '/resetpassword' && initialRoutePath !== '/oauth2_callback') {
    if (!settingsStore.appSettings.applicationLock || isUserUnlocked()) {
        // refresh token if user is logined
        tokensStore.refreshTokenAndRevokeOldToken().then(response => {
            const user = response.user;
            const languageSettings = user ? setLanguage(user.language).catch(error => {
                if (!error.processed) {
                    desktopPageStore.addErrorToast(error);
                }

                return Promise.reject(error);
            }) : Promise.reject();

            languageSettings.then(localeDefaultSettings => {
                settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
            }).finally(() => {
                if (user) {
                    setExpenseAndIncomeAmountColor(user.expenseAmountColor, user.incomeAmountColor);

                    if (response.notificationContent) {
                        rootStore.setNotificationContent(response.notificationContent);
                    }
                }

                updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
                updateTransactionPictureCacheExpiration(settingsStore.appSettings.transactionPictureCacheExpiration);
                exchangeRatesStore.removeExpiredExchangeRates(true);
                exchangeRatesStore.autoUpdateExchangeRatesData();
            });
        }).catch(error => {
            logger.error('failed to refresh token', error);
        });
    }
}

if (isProduction()) {
    register('./sw.js', {
        registrationOptions: {
            scope: './'
        }
    });
}
</script>

<style>
.notification-logo {
    width: 1.2rem;
    height: 1.2rem;
}
</style>
