<template>
    <v-app>
        <router-view />
    </v-app>
    <v-snackbar class="cursor-pointer" color="notification-background" location="top"
                :multi-line="true" :timeout="-1" :close-on-content-click="true" v-model="showNotification">
        <v-tooltip activator="parent">{{ $t('Click to close') }}</v-tooltip>
        <div class="d-inline-flex">
            <img alt="logo" class="notification-logo" :src="ezBookkeepingLogoPath" />
            <span class="ml-2">{{ $t('global.app.title') }}</span>
        </div>
        <div>
            {{ currentNotificationContent }}
        </div>
    </v-snackbar>
</template>

<script>
import { useTheme } from 'vuetify';
import { register } from 'register-service-worker';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import { isProduction } from '@/lib/version.ts';
import { initMapProvider } from '@/lib/map/index.ts';
import { getSystemTheme, setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export default {
    data() {
        return {
            showNotification: false
        }
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
        ezBookkeepingLogoPath() {
            return APPLICATION_LOGO_PATH;
        },
        currentNotificationContent() {
            return this.rootStore.currentNotification;
        }
    },
    watch: {
        currentNotificationContent: function (newValue) {
            this.showNotification = !!newValue;
        }
    },
    created() {
        const self = this;
        const theme = useTheme();

        if (self.settingsStore.appSettings.theme === ThemeType.Light) {
            theme.global.name.value = ThemeType.Light;
        } else if (self.settingsStore.appSettings.theme === ThemeType.Dark) {
            theme.global.name.value = ThemeType.Dark;
        } else {
            theme.global.name.value = getSystemTheme();
        }

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
            if (self.settingsStore.appSettings.theme === 'auto') {
                if (e.matches) {
                    theme.global.name.value = ThemeType.Dark;
                } else {
                    theme.global.name.value = ThemeType.Light;
                }
            }
        });

        let localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

        setExpenseAndIncomeAmountColor(self.userStore.currentUserExpenseAmountColor, self.userStore.currentUserIncomeAmountColor);

        if (self.$user.isUserLogined()) {
            if (!self.settingsStore.appSettings.applicationLock || self.$user.isUserUnlocked()) {
                // refresh token if user is logined
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                        setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);

                        if (response.notificationContent) {
                            self.rootStore.setNotificationContent(response.notificationContent);
                        }
                    }
                });

                // auto refresh exchange rates data
                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
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
    },
    mounted() {
        document.addEventListener('DOMContentLoaded', () => {
            const languageInfo = this.$locale.getCurrentLanguageInfo();
            initMapProvider(languageInfo ? languageInfo.alternativeLanguageTag : null);
        });
    }
}
</script>

<style>
.notification-logo {
    width: 1.2rem;
    height: 1.2rem;
}
</style>
