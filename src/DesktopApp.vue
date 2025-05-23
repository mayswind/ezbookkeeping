<template>
    <v-app>
        <router-view />
    </v-app>
    <v-snackbar class="cursor-pointer" color="notification-background" location="top"
                :multi-line="true" :timeout="-1" :close-on-content-click="true" v-model="showNotification">
        <v-tooltip activator="parent">{{ tt('Click to close') }}</v-tooltip>
        <div class="d-inline-flex">
            <img alt="logo" class="notification-logo" :src="APPLICATION_LOGO_PATH" />
            <span class="ml-2">{{ tt('global.app.title') }}</span>
        </div>
        <div>
            {{ currentNotificationContent }}
        </div>
    </v-snackbar>
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

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import { isProduction } from '@/lib/version.ts';
import { initMapProvider } from '@/lib/map/index.ts';
import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';
import { getSystemTheme, setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

const { tt, getCurrentLanguageInfo, setLanguage, initLocale } = useI18n();

const theme = useTheme();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const tokensStore = useTokensStore();
const exchangeRatesStore = useExchangeRatesStore();

const showNotification = ref<boolean>(false);

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
    theme.global.name.value = ThemeType.Light;
} else if (settingsStore.appSettings.theme === ThemeType.Dark) {
    theme.global.name.value = ThemeType.Dark;
} else {
    theme.global.name.value = getSystemTheme();
}

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
    if (settingsStore.appSettings.theme === 'auto') {
        if (e.matches) {
            theme.global.name.value = ThemeType.Dark;
        } else {
            theme.global.name.value = ThemeType.Light;
        }
    }
});

let localeDefaultSettings = initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone);
settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);

if (isUserLogined()) {
    if (!settingsStore.appSettings.applicationLock || isUserUnlocked()) {
        // refresh token if user is logined
        tokensStore.refreshTokenAndRevokeOldToken().then(response => {
            if (response.user) {
                localeDefaultSettings = setLanguage(response.user.language);
                settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);

                if (response.notificationContent) {
                    rootStore.setNotificationContent(response.notificationContent);
                }
            }
        });

        // auto refresh exchange rates data
        if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
            exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
        }
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
