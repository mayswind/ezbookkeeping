<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="ezBookkeepingLogoPath" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ $t('global.app.title') }}</h1>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer" v-if="!isDarkMode">
                    <v-img src="img/desktop/background.svg"/>
                </div>
                <div class="d-flex auth-img-footer" v-if="isDarkMode">
                    <v-img src="img/desktop/background-dark.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img max-width="600px" src="img/desktop/people3.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ $t('Unlock Application') }}</h4>
                            <p class="mb-0" v-if="isWebAuthnAvailable">{{ $t('Please input your PIN code or use WebAuthn to unlock application') }}</p>
                            <p class="mb-0" v-else-if="!isWebAuthnAvailable">{{ $t('Please input your PIN code to unlock application') }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12">
                                        <pin-code-input :disabled="verifyingByWebAuthn" :autofocus="true"
                                                        :secure="true" :length="6"
                                                        v-model="pinCode" @pincode:confirm="unlockByPin" />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block :disabled="!isPinCodeValid(pinCode) || verifyingByWebAuthn"
                                               @click="unlockByPin(pinCode)">
                                            {{ $t('Unlock By PIN Code') }}
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" v-if="isWebAuthnAvailable">
                                        <v-btn block variant="tonal" :disabled="verifyingByWebAuthn"
                                               @click="unlockByWebAuthn">
                                            {{ $t('Unlock By WebAuthn') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="verifyingByWebAuthn"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" class="text-center">
                                        <span class="me-1">{{ $t('Can\'t Unlock?') }}</span>
                                        <a class="text-primary" href="javascript:void(0);" @click="relogin">
                                            {{ $t('Re-login') }}
                                        </a>
                                    </v-col>
                                </v-row>
                            </v-form>
                        </v-card-text>
                    </v-card>
                </div>
                <v-spacer/>
                <div class="d-flex align-center justify-center">
                    <v-card variant="flat" class="w-100 px-4 pb-4" max-width="500">
                        <v-card-text class="pt-0">
                            <v-row>
                                <v-col cols="12" class="text-center">
                                    <v-menu location="bottom">
                                        <template #activator="{ props }">
                                            <v-btn variant="text"
                                                   :disabled="verifyingByWebAuthn"
                                                   v-bind="props">{{ currentLanguageName }}</v-btn>
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

                                <v-col cols="12" class="d-flex align-center pt-0">
                                    <v-divider />
                                </v-col>

                                <v-col cols="12" class="text-center text-sm">
                                    <span>Powered by </span>
                                    <a href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</a>&nbsp;<span>{{ version }}</span>
                                </v-col>
                            </v-row>
                        </v-card-text>
                    </v-card>
                </div>
            </v-col>
        </v-row>

        <confirm-dialog ref="confirmDialog"/>
        <snack-bar ref="snackbar" />
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

import assetConstants from '@/consts/asset.js';
import logger from '@/lib/logger.js';
import webauthn from '@/lib/webauthn.js';

export default {
    data() {
        return {
            pinCode: '',
            verifyingByWebAuthn: false
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
        ezBookkeepingLogoPath() {
            return assetConstants.ezBookkeepingLogoPath;
        },
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
        isDarkMode() {
            return this.globalTheme.global.name.value === 'dark';
        },
        currentLanguageName() {
            return this.$locale.getCurrentLanguageDisplayName();
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

            self.verifyingByWebAuthn = true;

            webauthn.verifyCredential(
                self.userStore.currentUserInfo,
                self.$user.getWebAuthnCredentialId()
            ).then(({ id, userName, userSecret }) => {
                self.verifyingByWebAuthn = false;

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
                self.verifyingByWebAuthn = false;
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
