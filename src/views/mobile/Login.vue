<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>
            <img class="login-page-logo" src="img/ezbookkeeping-192.png" />
            <f7-block class="margin-vertical-half">{{ $t('global.app.title') }}</f7-block>
        </f7-login-screen-title>
        <f7-list form>
            <f7-list-input
                type="text"
                autocomplete="username"
                clear-button
                :label="$t('Username')"
                :placeholder="$t('Your username or email')"
                :value="username"
                @input="username = $event.target.value; tempToken = ''"
            ></f7-list-input>
            <f7-list-input
                type="password"
                autocomplete="current-password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password')"
                :value="password"
                @input="password = $event.target.value; tempToken = ''"
                @keyup.enter.native="loginByPressEnter"
            ></f7-list-input>
        </f7-list>
        <f7-list>
            <f7-list-button :class="{ 'disabled': inputIsEmpty || logining }" :text="$t('Log In')" @click="login"></f7-list-button>
            <f7-block-footer>
                <span>{{ $t('Don\'t have an account?') }}</span>&nbsp;
                <f7-link :class="{'disabled': !isUserRegistrationEnabled}" href="/signup" :text="$t('Create an account')"></f7-link>
                <br/>
                <f7-link class="disabled" href="/forget-pwd" :text="$t('Forget Password?')"></f7-link>
            </f7-block-footer>
            <f7-block-footer>
            </f7-block-footer>
        </f7-list>
        <f7-button small popover-open=".lang-popover-menu" :text="currentLanguageName"></f7-button>
        <f7-list>
            <f7-block-footer>
                <span>Powered by </span>
                <f7-link external href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</f7-link>&nbsp;
                <span>{{ version }}</span>
            </f7-block-footer>
            <f7-block-footer>
            </f7-block-footer>
        </f7-list>

        <f7-popover class="lang-popover-menu">
            <f7-list>
                <f7-list-item
                    link="#" no-chevron popover-close
                    v-for="(lang, locale) in allLanguages"
                    :key="locale"
                    :title="lang.displayName"
                    @click="changeLanguage(locale)"
                >
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="$i18n.locale === locale"></f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-sheet
            style="height:auto"
            :opened="show2faSheet" @sheet:closed="show2faSheet = false"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Two-Factor Authentication') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-input
                            type="number"
                            autocomplete="one-time-code"
                            outline
                            clear-button
                            v-if="twoFAVerifyType === 'passcode'"
                            :placeholder="$t('Passcode')"
                            :value="passcode"
                            @input="passcode = $event.target.value"
                            @keyup.enter.native="verify"
                        ></f7-list-input>
                        <f7-list-input
                            outline
                            clear-button
                            v-if="twoFAVerifyType === 'backupcode'"
                            :placeholder="$t('Backup Code')"
                            :value="backupCode"
                            @input="backupCode = $event.target.value"
                            @keyup.enter.native="verify"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': twoFAInputIsEmpty || verifying }" :text="$t('Verify')" @click="verify"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="switch2FAVerifyType" :text="$t(twoFAVerifyTypeSwitchName)"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            username: '',
            password: '',
            passcode: '',
            backupCode: '',
            tempToken: '',
            logining: false,
            verifying: false,
            show2faSheet: false,
            twoFAVerifyType: 'passcode'
        };
    },
    computed: {
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguages();
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
        twoFAVerifyTypeSwitchName() {
            if (this.twoFAVerifyType === 'backupcode') {
                return 'Use a passcode';
            } else {
                return 'Use a backup code';
            }
        },
        currentLanguageName() {
            const currentLocale = this.$i18n.locale;
            let lang = this.$locale.getLanguage(currentLocale);

            if (!lang) {
                lang = this.$locale.getLanguage(this.$locale.getDefaultLanguage());
            }

            return lang.displayName;
        }
    },
    methods: {
        login() {
            const self = this;
            const router = self.$f7router;

            if (!this.username) {
                self.$alert('Username cannot be empty');
                return;
            }

            if (!this.password) {
                self.$alert('Password cannot be empty');
                return;
            }

            if (self.tempToken) {
                self.show2faSheet = true;
                return;
            }

            self.logining = true;
            self.$showLoading(() => self.logining);

            self.$store.dispatch('authorize', {
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

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
                }

                router.refreshPage();
            }).catch(error => {
                self.logining = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        loginByPressEnter() {
            const app = this.$f7;
            const $$ = app.$;

            if ($$('.modal-in').length) {
                return;
            }

            return this.login();
        },
        verify() {
            const self = this;
            const router = self.$f7router;

            if (self.twoFAInputIsEmpty || self.verifying) {
                return;
            }

            if (this.twoFAVerifyType === 'passcode' && !this.passcode) {
                self.$alert('Passcode cannot be empty');
                return;
            } else if (this.twoFAVerifyType === 'backupcode' && !this.backupCode) {
                self.$alert('Backup code cannot be empty');
                return;
            }

            self.verifying = true;
            self.$showLoading(() => self.verifying);

            self.$store.dispatch('authorize2FA', {
                token: self.tempToken,
                passcode: self.twoFAVerifyType === 'passcode' ? self.passcode : null,
                recoveryCode: self.twoFAVerifyType === 'backupcode' ? self.backupCode : null
            }).then(() => {
                self.verifying = false;
                self.$hideLoading();

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
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
        switch2FAVerifyType() {
            if (this.twoFAVerifyType === 'passcode') {
                this.twoFAVerifyType = 'backupcode';
            } else {
                this.twoFAVerifyType = 'passcode';
            }
        },
        changeLanguage(locale) {
            this.$locale.setLanguage(locale);
        }
    }
};
</script>

<style>
.login-page-logo {
    width: 80px;
    height: 80px;
}

.lang-popover-menu .popover-inner {
    max-height: 300px;
    overflow-Y: auto;
}
</style>
