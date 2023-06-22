<template>
    <div class="auth-wrapper d-flex align-center justify-center pa-4">
        <v-card class="auth-card pa-4 pt-7" max-width="448">
            <v-card-item class="justify-center">
                <v-card-title class="d-grid font-weight-semibold text-2xl">
                    <v-img alt="logo" class="login-page-logo" src="/img/ezbookkeeping-192.png" :width="96" />
                    <p class="mt-4 font-weight-bold">{{ $t('global.app.title') }}</p>
                </v-card-title>
            </v-card-item>

            <v-card-text>
                <v-form>
                    <v-row>
                        <v-col cols="12">
                            <v-text-field
                                type="text"
                                autocomplete="username"
                                clearable
                                :disabled="show2faInput"
                                :label="$t('Username')"
                                :placeholder="$t('Your username or email')"
                                v-model="username"
                                @input="tempToken = ''"
                                @keyup.enter="$refs.passwordInput.focus()"
                            />
                        </v-col>

                        <v-col cols="12">
                            <v-text-field
                                autocomplete="current-password"
                                clearable
                                ref="passwordInput"
                                :type="isPasswordVisible ? 'text' : 'password'"
                                :disabled="show2faInput"
                                :label="$t('Password')"
                                :placeholder="$t('Your password')"
                                :append-inner-icon="isPasswordVisible ? icons.eyeSlash : icons.eye"
                                v-model="password"
                                @input="tempToken = ''"
                                @click:append-inner="isPasswordVisible = !isPasswordVisible"
                                @keyup.enter="login"
                            />
                        </v-col>

                        <v-col cols="12" v-show="show2faInput">
                            <v-text-field
                                type="number"
                                autocomplete="one-time-code"
                                clearable
                                ref="passcodeInput"
                                :label="$t('Passcode')"
                                :placeholder="$t('Passcode')"
                                :append-inner-icon="icons.backupCode"
                                v-model="passcode"
                                @click:append-inner="twoFAVerifyType = 'backupcode'"
                                @keyup.enter="verify"
                                v-if="twoFAVerifyType === 'passcode'"
                            />
                            <v-text-field
                                type="text"
                                clearable
                                :label="$t('Backup Code')"
                                :placeholder="$t('Backup Code')"
                                :append-inner-icon="icons.passcode"
                                v-model="backupCode"
                                @click:append-inner="twoFAVerifyType = 'passcode'"
                                @keyup.enter="verify"
                                v-if="twoFAVerifyType === 'backupcode'"
                            />
                        </v-col>

                        <v-col cols="12">
                            <v-btn block :class="{ 'disabled': inputIsEmpty || logining }"
                                   @click="login" v-if="!show2faInput">
                                {{ $t('Log In') }}
                            </v-btn>
                            <v-btn block :class="{ 'disabled': twoFAInputIsEmpty || verifying }"
                                   @click="verify" v-else-if="show2faInput">
                                {{ $t('Continue') }}
                            </v-btn>
                        </v-col>

                        <v-col cols="12" class="text-center text-base">
                            <span>{{ $t('Don\'t have an account?') }}</span>&nbsp;
                            <router-link class="text-primary ms-2" to="/signup">
                                {{ $t('Create an account') }}
                            </router-link>
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
                            <VDivider />
                        </v-col>

                        <v-col cols="12" class="text-center text-sm">
                            <span>Powered by </span>
                            <a href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</a>&nbsp;<span>{{ version }}</span>
                        </v-col>
                    </v-row>
                </v-form>
            </v-card-text>
        </v-card>

        <v-snackbar v-model="showSnackbar">
            {{ snackbarMessage }}

            <template #actions>
                <v-btn color="primary" variant="text" @click="showSnackbar = false">{{ $t('Close') }}</v-btn>
            </template>
        </v-snackbar>

        <v-overlay class="justify-center align-center" :persistent="true" v-model="logining">
            <v-progress-circular indeterminate></v-progress-circular>
        </v-overlay>

        <v-overlay class="justify-center align-center" :persistent="true" v-model="verifying">
            <v-progress-circular indeterminate></v-progress-circular>
        </v-overlay>
    </div>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import {
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiOnepassword,
    mdiHelpCircleOutline
} from '@mdi/js';

export default {
    data() {
        return {
            username: '',
            password: '',
            passcode: '',
            backupCode: '',
            tempToken: '',
            isPasswordVisible: false,
            logining: false,
            verifying: false,
            show2faInput: false,
            twoFAVerifyType: 'passcode',
            showSnackbar: false,
            snackbarMessage: '',
            icons: {
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                passcode: mdiOnepassword,
                backupCode: mdiHelpCircleOutline
            }
        };
    },

    computed: {
        ...mapStores(useRootStore, useSettingsStore, useExchangeRatesStore),
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        isUserRegistrationEnabled() {
            return this.$settings.isUserRegistrationEnabled();
        },
        inputIsEmpty() {
            return !this.username || !this.password;
        },
        twoFAInputIsEmpty() {
            if (this.twoFAVerifyType === 'backupcode') {
                return !this.backupCode;
            } else {
                return !this.passcode;
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
    methods: {
        login() {
            const self = this;

            if (!self.username) {
                self.showSnackbarMessage(self.$t('Username cannot be empty'));
                return;
            }

            if (!self.password) {
                self.showSnackbarMessage(self.$t('Password cannot be empty'));
                return;
            }

            if (self.tempToken) {
                self.show2faInput = true;
                return;
            }

            if (self.logining) {
                return;
            }

            self.isPasswordVisible = false;
            self.logining = true;

            self.rootStore.authorize({
                loginName: self.username,
                password: self.password
            }).then(authResponse => {
                self.logining = false;

                if (authResponse.need2FA) {
                    self.tempToken = authResponse.token;
                    self.show2faInput = true;

                    this.$nextTick(() => {
                        if (self.$refs.passcodeInput) {
                            self.$refs.passcodeInput.focus();
                            self.$refs.passcodeInput.select();
                        }
                    });

                    return;
                }

                if (authResponse.user && authResponse.user.language) {
                    const localeDefaultSettings = self.$locale.setLanguage(authResponse.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.exchangeRatesStore.getLatestExchangeRates({silent: true, force: false});
                }

                this.$router.replace('/');
            }).catch(error => {
                self.logining = false;

                if (!error.processed) {
                    self.showSnackbarMessage(self.$tError(error.message || error));
                }
            });
        },
        verify() {
            const self = this;

            if (self.twoFAInputIsEmpty || self.verifying) {
                return;
            }

            if (this.twoFAVerifyType === 'passcode' && !this.passcode) {
                self.showSnackbarMessage(self.$t('Passcode cannot be empty'));
                return;
            } else if (this.twoFAVerifyType === 'backupcode' && !this.backupCode) {
                self.showSnackbarMessage(self.$t('Backup code cannot be empty'));
                return;
            }

            self.verifying = true;

            self.rootStore.authorize2FA({
                token: self.tempToken,
                passcode: self.twoFAVerifyType === 'passcode' ? self.passcode : null,
                recoveryCode: self.twoFAVerifyType === 'backupcode' ? self.backupCode : null
            }).then(authResponse => {
                self.verifying = false;

                if (authResponse.user && authResponse.user.language) {
                    const localeDefaultSettings = self.$locale.setLanguage(authResponse.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                this.$router.replace('/');
            }).catch(error => {
                self.verifying = false;

                if (!error.processed) {
                    self.showSnackbarMessage(self.$tError(error.message || error));
                }
            });
        },
        switch2FAVerifyType() {
            if (this.twoFAVerifyType === 'passcode') {
                this.twoFAVerifyType = 'backupcode';
            } else {
                this.twoFAVerifyType = 'passcode';
            }
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        },
        showSnackbarMessage(message) {
            this.showSnackbar = true;
            this.snackbarMessage = message;
        }
    }
}
</script>

<style>
.auth-wrapper {
    min-block-size: calc(var(--vh, 1vh) * 100);
}

.auth-card {
    z-index: 1 !important
}

.login-page-logo {
    margin: auto;
}
</style>
