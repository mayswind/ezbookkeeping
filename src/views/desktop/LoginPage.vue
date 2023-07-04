<template>
    <div class="layout-wrapper">
        <div class="auth-logo d-flex align-start gap-x-3">
            <v-img alt="logo" class="login-page-logo" src="/img/ezbookkeeping-192.png" />
            <h1 class="font-weight-medium leading-normal text-2xl">{{ $t('global.app.title') }}</h1>
        </div>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer">
                    <v-img src="/img/illustrations/background.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img max-width="600px" src="/img/illustrations/people1.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex align-center justify-center">
                <v-card variant="flat" class="mt-12 mt-sm-0 pa-4" max-width="500">
                    <v-card-text>
                        <h5 class="text-h5 mb-1">{{ $t('Welcome to ezBookkeeping') }}</h5>
                        <p class="mb-0">{{ $t('Please log in with your ezBookkeeping account') }}</p>
                    </v-card-text>

                    <v-card-text>
                        <v-form>
                            <v-row>
                                <v-col cols="12">
                                    <v-text-field
                                        type="text"
                                        autocomplete="username"
                                        autofocus="autofocus"
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

                                <v-col cols="12" class="text-center text-base" v-if="isUserRegistrationEnabled">
                                    <span class="me-1">{{ $t('Don\'t have an account?') }}</span>
                                    <router-link class="text-primary" to="/signup">
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

        <snack-bar ref="snackbar" />

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

import { isUserRegistrationEnabled } from '@/lib/server_settings.js';

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
            return isUserRegistrationEnabled();
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
                self.$refs.snackbar.showMessage('Username cannot be empty');
                return;
            }

            if (!self.password) {
                self.$refs.snackbar.showMessage('Password cannot be empty');
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

                if (authResponse.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(authResponse.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                this.$router.replace('/');
            }).catch(error => {
                self.logining = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        verify() {
            const self = this;

            if (self.twoFAInputIsEmpty || self.verifying) {
                return;
            }

            if (this.twoFAVerifyType === 'passcode' && !this.passcode) {
                self.$refs.snackbar.showMessage('Passcode cannot be empty');
                return;
            } else if (this.twoFAVerifyType === 'backupcode' && !this.backupCode) {
                self.$refs.snackbar.showMessage('Backup code cannot be empty');
                return;
            }

            self.verifying = true;

            self.rootStore.authorize2FA({
                token: self.tempToken,
                passcode: self.twoFAVerifyType === 'passcode' ? self.passcode : null,
                recoveryCode: self.twoFAVerifyType === 'backupcode' ? self.backupCode : null
            }).then(authResponse => {
                self.verifying = false;

                if (authResponse.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(authResponse.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                this.$router.replace('/');
            }).catch(error => {
                self.verifying = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    }
}
</script>
