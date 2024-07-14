<template>
    <img style="display: none;" :src="devCookiePath" v-if="!isProduction" />
    <v-app>
        <router-view />
    </v-app>
</template>

<script>
import { useTheme } from 'vuetify';
import { register } from 'register-service-worker';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { isProduction } from '@/lib/version.js';
import { loadMapAssets } from '@/lib/map/index.js';
import { getSystemTheme, setExpenseAndIncomeAmountColor } from '@/lib/ui.js';

export default {
    data() {
        return {
            isProduction: isProduction(),
            devCookiePath: isProduction() ? '' : '/dev/cookies'
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
    },
    created() {
        const self = this;
        const theme = useTheme();

        if (self.settingsStore.appSettings.theme === 'light') {
            theme.global.name.value = 'light';
        } else if (self.settingsStore.appSettings.theme === 'dark') {
            theme.global.name.value = 'dark';
        } else {
            theme.global.name.value = getSystemTheme();
        }

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
            if (self.settingsStore.appSettings.theme === 'auto') {
                if (e.matches) {
                    theme.global.name.value = 'dark';
                } else {
                    theme.global.name.value = 'light';
                }
            }
        });

        let localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

        setExpenseAndIncomeAmountColor(self.userStore.currentUserExpenseAmountColor, self.userStore.currentUserIncomeAmountColor);

        if (self.$user.isUserLogined()) {
            if (!self.settingsStore.appSettings.applicationLock) {
                // refresh token if user is logined
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                        setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
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
            loadMapAssets(languageInfo ? languageInfo.code : null);
        });
    }
}
</script>
