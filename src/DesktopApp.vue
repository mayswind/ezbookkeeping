<template>
    <img style="display: none;" :src="devCookiePath" v-if="!isProduction" />
    <v-app>
        <router-view />
    </v-app>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { loadMapAssets } from '@/lib/map/index.js';
import { getSystemTheme } from '@/lib/ui.js';

export default {
    data() {
        const self = this;

        return {
            isProduction: self.$settings.isProduction(),
            devCookiePath: self.$settings.isProduction() ? '' : '/dev/cookies'
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
    },
    created() {
        const self = this;
        const theme = useTheme();

        if (self.$settings.getTheme() === 'light') {
            theme.global.name.value = 'light';
        } else if (self.$settings.getTheme() === 'dark') {
            theme.global.name.value = 'dark';
        } else {
            theme.global.name.value = getSystemTheme();
        }

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
            if (self.$settings.getTheme() === 'auto') {
                if (e.matches) {
                    theme.global.name.value = 'dark';
                } else {
                    theme.global.name.value = 'light';
                }
            }
        });

        let localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage);
        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

        if (self.$user.isUserLogined()) {
            if (!self.$settings.isEnableApplicationLock()) {
                // refresh token if user is logined
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user && response.user.language) {
                        localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                    }
                });

                // auto refresh exchange rates data
                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }
            }
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

<style>
/** Global style **/
* {
    padding: 0;
    margin: 0
}
</style>
