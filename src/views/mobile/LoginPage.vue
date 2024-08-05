<template>
    <f7-page no-navbar no-swipeback login-screen hide-toolbar-on-scroll>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" :src="ezBookkeepingLogoPath" />
            <f7-block class="login-page-tile margin-vertical-half">{{ $t('global.app.title') }}</f7-block>
        </f7-login-screen-title>

        <f7-list form dividers class="margin-bottom-half">
            <f7-list-input
                type="text"
                autocomplete="username"
                clear-button
                :label="$t('Username')"
                :placeholder="$t('Your username or email')"
                v-model:value="username"
                @input="tempToken = ''"
            ></f7-list-input>
            <f7-list-input
                type="password"
                autocomplete="current-password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password')"
                v-model:value="password"
                @input="tempToken = ''"
                @keyup.enter="loginByPressEnter"
            ></f7-list-input>
        </f7-list>

        <f7-list class="no-margin-vertical">
            <f7-list-item>
                <template #title>
                    <small>
                        <f7-link external :href="desktopVersionPath">{{ $t('Switch to Desktop Version') }}</f7-link>
                    </small>
                </template>
                <template #after>
                    <small>
                        <f7-link :class="{'disabled': !isUserForgetPasswordEnabled}" @click="forgetPasswordEmail = ''; showForgetPasswordSheet = true">{{ $t('Forget Password?') }}</f7-link>
                    </small>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list class="margin-vertical-half">
            <f7-list-button :class="{ 'disabled': inputIsEmpty || logining }" :text="$t('Log In')" @click="login"></f7-list-button>
            <f7-block-footer>
                <span>{{ $t('Don\'t have an account?') }}</span>&nbsp;
                <f7-link :class="{'disabled': !isUserRegistrationEnabled}" href="/signup" :text="$t('Create an account')"></f7-link>
            </f7-block-footer>
            <f7-block-footer class="padding-bottom">
            </f7-block-footer>
        </f7-list>

        <f7-button small popover-open=".lang-popover-menu" :text="currentLanguageName"></f7-button>

        <f7-list class="login-page-bottom">
            <f7-block-footer>
                <div class="login-page-powered-by">
                    <span>Powered by</span>
                    <f7-link external href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</f7-link>
                    <span>{{ version }}</span>
                </div>
            </f7-block-footer>
        </f7-list>

        <f7-toolbar class="login-page-fixed-bottom" tabbar bottom :outline="false">
            <div class="login-page-powered-by">
                <span>Powered by</span>
                <f7-link external href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</f7-link>
                <span>{{ version }}</span>
            </div>
        </f7-toolbar>

        <f7-popover class="lang-popover-menu">
            <f7-list dividers>
                <f7-list-item
                    link="#" no-chevron popover-close
                    :title="lang.displayName"
                    :key="locale"
                    v-for="(lang, locale) in allLanguages"
                    @click="changeLanguage(locale)"
                >
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="currentLanguageCode === locale"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-sheet
            style="height:auto"
            :opened="show2faSheet" @sheet:closed="show2faSheet = false"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ $t('Two-Factor Authentication') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-list strong class="no-margin">
                        <f7-list-input
                            type="number"
                            autocomplete="one-time-code"
                            outline
                            floating-label
                            clear-button
                            class="no-margin no-padding-bottom"
                            v-if="twoFAVerifyType === 'passcode'"
                            :label="$t('Passcode')"
                            :placeholder="$t('Passcode')"
                            v-model:value="passcode"
                            @keyup.enter="verify"
                        ></f7-list-input>
                        <f7-list-input
                            outline
                            floating-label
                            clear-button
                            class="no-margin no-padding-bottom"
                            v-if="twoFAVerifyType === 'backupcode'"
                            :label="$t('Backup Code')"
                            :placeholder="$t('Backup Code')"
                            v-model:value="backupCode"
                            @keyup.enter="verify"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': twoFAInputIsEmpty || verifying }" :text="$t('Verify')" @click="verify"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="switch2FAVerifyType" :text="$t(twoFAVerifyTypeSwitchName)"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
            style="height:auto"
            :opened="showForgetPasswordSheet" @sheet:closed="showForgetPasswordSheet = false"
        >
            <div class="swipe-handler" style="z-index: 10"></div>
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ $t('Forget Password?') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin">
                        <span>{{ $t('Please enter your email address used for registration and we\'ll send you an email with a reset password link') }}</span>
                    </p>
                    <f7-list strong class="no-margin">
                        <f7-list-input
                            type="email"
                            autocomplete="email"
                            outline
                            floating-label
                            clear-button
                            class="no-margin no-padding-bottom"
                            :label="$t('E-mail')"
                            :placeholder="$t('Your email address')"
                            v-model:value="forgetPasswordEmail"
                            @keyup.enter="requestResetPassword"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !forgetPasswordEmail || requestingForgetPassword }" :text="$t('Send Reset Link')" @click="requestResetPassword"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': requestingForgetPassword }" @click="showForgetPasswordSheet = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <password-input-sheet :title="$t('Verify your email')"
                              :hint="$t(hasValidEmailVerifyToken ? 'format.misc.accountActivationAndResendValidationEmailTip' : 'format.misc.resendValidationEmailTip', { email: resendVerifyEmail })"
                              :confirm-disabled="requestingResendVerifyEmail"
                              :cancel-disabled="requestingResendVerifyEmail"
                              v-model:show="showVerifyEmailSheet"
                              v-model="currentPasswordForResendVerifyEmail"
                              @password:confirm="requestResendVerifyEmail">
        </password-input-sheet>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import assetConstants from '@/consts/asset.js';
import {
    isUserRegistrationEnabled,
    isUserForgetPasswordEnabled,
    isUserVerifyEmailEnabled
} from '@/lib/server_settings.js';
import { getDesktopVersionPath } from '@/lib/version.js';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui.js';
import { isModalShowing } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            username: '',
            password: '',
            passcode: '',
            backupCode: '',
            tempToken: '',
            forgetPasswordEmail: '',
            resendVerifyEmail: '',
            hasValidEmailVerifyToken: false,
            currentPasswordForResendVerifyEmail: '',
            logining: false,
            verifying: false,
            requestingForgetPassword: false,
            requestingResendVerifyEmail: false,
            show2faSheet: false,
            showForgetPasswordSheet: false,
            showVerifyEmailSheet: false,
            twoFAVerifyType: 'passcode'
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useExchangeRatesStore),
        ezBookkeepingLogoPath() {
            return assetConstants.ezBookkeepingLogoPath;
        },
        version() {
            return 'v' + this.$version;
        },
        desktopVersionPath() {
            return getDesktopVersionPath();
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        isUserRegistrationEnabled() {
            return isUserRegistrationEnabled();
        },
        isUserForgetPasswordEnabled() {
            return isUserForgetPasswordEnabled();
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
        twoFAVerifyTypeSwitchName() {
            if (this.twoFAVerifyType === 'backupcode') {
                return 'Use Passcode';
            } else {
                return 'Use Backup Code';
            }
        },
        currentLanguageCode() {
            return this.$locale.getCurrentLanguageTag();
        },
        currentLanguageName() {
            return this.$locale.getCurrentLanguageDisplayName();
        },
        isUserVerifyEmailEnabled() {
            return isUserVerifyEmailEnabled();
        }
    },
    methods: {
        login() {
            const self = this;
            const router = self.f7router;

            if (!this.username) {
                self.$alert('Username cannot be blank');
                return;
            }

            if (!this.password) {
                self.$alert('Password cannot be blank');
                return;
            }

            if (self.tempToken) {
                self.show2faSheet = true;
                return;
            }

            self.logining = true;
            self.resendVerifyEmail = '';
            self.hasValidEmailVerifyToken = false;
            self.currentPasswordForResendVerifyEmail = '';
            self.$showLoading(() => self.logining);

            self.rootStore.authorize({
                loginName: self.username,
                password: self.password
            }).then(authResponse => {
                self.logining = false;
                self.$hideLoading();

                if (authResponse.need2FA) {
                    self.tempToken = authResponse.token;
                    self.show2faSheet = true;
                    return;
                }

                if (authResponse.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(authResponse.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                    setExpenseAndIncomeAmountColor(authResponse.user.expenseAmountColor, authResponse.user.incomeAmountColor);
                }

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                if (authResponse.notificationContent) {
                    self.rootStore.setNotificationContent(authResponse.notificationContent);
                }

                router.refreshPage();
            }).catch(error => {
                self.logining = false;
                self.$hideLoading();

                if (self.isUserVerifyEmailEnabled && error.error && error.error.errorCode === 201020 && error.error.context && error.error.context.email) {
                    self.resendVerifyEmail = error.error.context.email;
                    self.hasValidEmailVerifyToken = error.error.context.hasValidEmailVerifyToken || false;
                    self.currentPasswordForResendVerifyEmail = '';
                    self.showVerifyEmailSheet = true;
                    return;
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        loginByPressEnter() {
            if (isModalShowing()) {
                return;
            }

            return this.login();
        },
        verify() {
            const self = this;
            const router = self.f7router;

            if (self.twoFAInputIsEmpty || self.verifying) {
                return;
            }

            if (this.twoFAVerifyType === 'passcode' && !this.passcode) {
                self.$alert('Passcode cannot be blank');
                return;
            } else if (this.twoFAVerifyType === 'backupcode' && !this.backupCode) {
                self.$alert('Backup code cannot be blank');
                return;
            }

            self.verifying = true;
            self.$showLoading(() => self.verifying);

            self.rootStore.authorize2FA({
                token: self.tempToken,
                passcode: self.twoFAVerifyType === 'passcode' ? self.passcode : null,
                recoveryCode: self.twoFAVerifyType === 'backupcode' ? self.backupCode : null
            }).then(authResponse => {
                self.verifying = false;
                self.$hideLoading();

                if (authResponse.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(authResponse.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                    setExpenseAndIncomeAmountColor(authResponse.user.expenseAmountColor, authResponse.user.incomeAmountColor);
                }

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                if (authResponse.notificationContent) {
                    self.rootStore.setNotificationContent(authResponse.notificationContent);
                }

                self.show2faSheet = false;
                router.refreshPage();
            }).catch(error => {
                self.verifying = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        requestResetPassword() {
            const self = this;

            if (!self.forgetPasswordEmail) {
                self.$alert('Email address cannot be blank');
                return;
            }

            self.requestingForgetPassword = true;
            self.$showLoading(() => self.requestingForgetPassword);

            self.rootStore.requestResetPassword({
                email: self.forgetPasswordEmail
            }).then(() => {
                self.requestingForgetPassword = false;
                self.$hideLoading();

                self.$toast('Password reset email has been sent');
                self.showForgetPasswordSheet = false;
            }).catch(error => {
                self.requestingForgetPassword = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        requestResendVerifyEmail() {
            const self = this;

            if (!self.currentPasswordForResendVerifyEmail) {
                self.$toast('Current password cannot be blank');
                return;
            }

            self.requestingResendVerifyEmail = true;
            self.$showLoading(() => self.requestingResendVerifyEmail);

            self.rootStore.resendVerifyEmailByUnloginUser({
                email: self.resendVerifyEmail,
                password: self.currentPasswordForResendVerifyEmail
            }).then(() => {
                self.requestingResendVerifyEmail = false;
                self.$hideLoading();

                self.$toast('Validation email has been sent');
                self.showVerifyEmailSheet = false;
            }).catch(error => {
                self.requestingResendVerifyEmail = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
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
        }
    }
};
</script>
