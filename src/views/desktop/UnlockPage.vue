<template>
    <div class="layout-wrapper">
        <div class="auth-logo d-flex align-start gap-x-3">
            <v-img alt="logo" class="login-page-logo" src="/img/ezbookkeeping-192.png" />
            <h1 class="font-weight-medium leading-normal text-2xl">{{ $t('global.app.title') }}</h1>
        </div>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer" v-if="currentTheme !== 'dark'">
                    <v-img src="/img/desktop/background.svg"/>
                </div>
                <div class="d-flex auth-img-footer" v-if="currentTheme === 'dark'">
                    <v-img src="/img/desktop/background-dark.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img max-width="600px" src="/img/desktop/people3.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex align-center justify-center">
                <v-card variant="flat" class="mt-12 mt-sm-0 pa-4" max-width="500">
                    <v-card-text>
                        <h5 class="text-h5 mb-1">{{ $t('Unlock Application') }}</h5>
                        <p class="mb-0" v-if="isWebAuthnAvailable">{{ $t('Please input your PIN code or use WebAuthn to unlock application') }}</p>
                        <p class="mb-0" v-else-if="!isWebAuthnAvailable">{{ $t('Please input your PIN code to unlock application') }}</p>
                    </v-card-text>

                    <v-card-text>
                        <v-form>
                            <v-row>
                                <v-col cols="12">
                                    <pin-code-input :autofocus="true" :secure="true" :length="6"
                                                    v-model="pinCode" @pincode:confirm="unlockByPin" />
                                </v-col>

                                <v-col cols="12">
                                    <v-btn block :class="{ 'disabled': !isPinCodeValid(pinCode) }"
                                           @click="unlockByPin(pinCode)">
                                        {{ $t('Unlock By PIN Code') }}
                                    </v-btn>
                                </v-col>

                                <v-col cols="12" v-if="isWebAuthnAvailable">
                                    <v-btn block variant="tonal"
                                           @click="unlockByWebAuthn">
                                        {{ $t('Unlock By WebAuthn') }}
                                    </v-btn>
                                </v-col>

                                <v-col cols="12" class="text-center">
                                    <span class="me-1">{{ $t('Can\'t Unlock?') }}</span>
                                    <a class="text-primary" href="javascript:void(0);" @click="relogin">
                                        {{ $t('Re-login') }}
                                    </a>
                                </v-col>

                                <v-col cols="12" class="text-center">
                                    <v-menu location="bottom">
                                        <template #activator="{ props }">
                                            <v-btn variant="text" v-bind="props">{{ currentLanguageName }}</v-btn>
                                        </template>
                                        <v-list>
                                            <v-list-item v-for="(lang, locale) in allLanguages" :key="locale">
                                                <v-list-item-title
                                                    class="cursor-pointer"
                                                    @click="changeLanguage(locale)">
                                                    {{ lang.displayName }}
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-col>

                                <v-col cols="12" class="d-flex align-center">
                                    <v-divider />
                                </v-col>

                                <v-col cols="12" class="text-center text-sm">
                                    <span>Powered by </span>
                                    <a href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</a>&nbsp;<span>{{ version }}</span>
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>

        <confirm-dialog ref="confirmDialog"/>
        <snack-bar ref="snackbar" />

        <v-overlay class="justify-center align-center" :persistent="true" v-model="verifying">
            <v-progress-circular indeterminate></v-progress-circular>
        </v-overlay>
    </div>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import logger from '@/lib/logger.js';
import webauthn from '@/lib/webauthn.js';

export default {
    data() {
        return {
            pinCode: '',
            verifying: false
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        isWebAuthnAvailable() {
            return this.settingsStore.appSettings.applicationLockWebAuthn
                && this.$user.getWebAuthnCredentialId()
                && webauthn.isSupported();
        },
        currentTheme: {
            get: function () {
                return this.globalTheme.global.name.value;
            }
        },
        currentLanguageName() {
            const currentLocale = this.$i18n.locale;
            let lang = this.$locale.getLanguageInfo(currentLocale);

            if (!lang) {
                lang = this.$locale.getLanguageInfo(this.$locale.getDefaultLanguage());
            }

            return lang.displayName;
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        unlockByWebAuthn() {
            const self = this;

            if (!self.settingsStore.appSettings.applicationLockWebAuthn || !self.$user.getWebAuthnCredentialId()) {
                self.$refs.snackbar.showMessage('WebAuthn is not enabled');
                return;
            }

            if (!webauthn.isSupported()) {
                self.$refs.snackbar.showMessage('This device does not support WebAuthn');
                return;
            }

            this.verifying = true;

            webauthn.verifyCredential(
                self.userStore.currentUserInfo,
                self.$user.getWebAuthnCredentialId()
            ).then(({ id, userName, userSecret }) => {
                this.verifying = false;

                self.$user.unlockTokenByWebAuthn(id, userName, userSecret);
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                    }
                });

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                self.$router.replace('/');
            }).catch(error => {
                this.verifying = false;
                logger.error('failed to use webauthn to verify', error);

                if (error.notSupported) {
                    self.$refs.snackbar.showMessage('This device does not support WebAuthn');
                } else if (error.name === 'NotAllowedError') {
                    self.$refs.snackbar.showMessage('User has canceled authentication');
                } else if (error.invalid) {
                    self.$refs.snackbar.showMessage('Failed to authenticate by WebAuthn');
                } else {
                    self.$refs.snackbar.showMessage('User has canceled or this device does not support WebAuthn');
                }
            });
        },
        unlockByPin(pinCode) {
            const self = this;

            if (!self.isPinCodeValid(pinCode)) {
                return;
            }

            const user = self.userStore.currentUserInfo;

            if (!user || !user.username) {
                self.$refs.snackbar.showMessage('An error has occurred');
                return;
            }

            try {
                self.$user.unlockTokenByPinCode(user.username, pinCode);
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                    }
                });

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                self.$router.replace('/');
            } catch (ex) {
                logger.error('failed to unlock by pin code', ex);
                self.$refs.snackbar.showMessage('PIN code is wrong');
            }
        },
        relogin() {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to re-login?').then(() => {
                self.rootStore.forceLogout();
                self.settingsStore.clearAppSettings();

                const localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
                self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                self.$router.replace('/login');
            });
        },
        isPinCodeValid(pinCode) {
            return pinCode && pinCode.length === 6;
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    }
}
</script>
