<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>{{ $t('global.app.title') }}</f7-login-screen-title>
        <f7-list form>
            <f7-list-input
                type="text"
                clear-button
                :label="$t('Username')"
                :placeholder="$t('Your username or Email')"
                :value="username"
                @input="username = $event.target.value; tempToken = ''"
            ></f7-list-input>
            <f7-list-input
                type="password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password')"
                :value="password"
                @input="password = $event.target.value; tempToken = ''"
            ></f7-list-input>
        </f7-list>
        <f7-list>
            <f7-list-button :class="{ 'disabled': inputIsEmpty }" :text="$t('Log In')" @click="login"></f7-list-button>
            <f7-block-footer>
                <span v-t="'Don\'t have an account?'"></span>&nbsp;
                <f7-link href="/signup" :text="$t('Create an account')"></f7-link>
                <br/>
                <f7-link class="disabled" href="/forget-pwd" :text="$t('Forget Password?')"></f7-link>
            </f7-block-footer>
            <f7-block-footer>
            </f7-block-footer>
        </f7-list>
        <f7-button small popover-open=".popover-menu" :text="currentLanguageName"></f7-button>
        <f7-popover class="popover-menu">
            <f7-list>
                <f7-list-item
                    link="#" popover-close
                    v-for="(lang, locale) in allLanguages" :key="locale"
                    :title="lang.displayName"
                    @click="changeLanguage(locale)"
                ></f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-sheet
            id="2fa-auth-sheet"
            style="height:auto; --f7-sheet-bg-color: #fff;"
            backdrop
        >
            <div class="sheet-modal-swipe-step">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b v-t="'Two-Factor Authentication'"></b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-list no-hairlines class="twofa-auth-form">
                        <f7-list-input
                            type="number"
                            outline
                            clear-button
                            v-if="twoFAVerifyType === 'passcode'"
                            :placeholder="$t('Passcode')"
                            :value="passcode"
                            @input="passcode = $event.target.value"
                        ></f7-list-input>
                        <f7-list-input
                            outline
                            clear-button
                            v-if="twoFAVerifyType === 'backupcode'"
                            :placeholder="$t('Backup Code')"
                            :value="backupCode"
                            @input="backupCode = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': twoFAInputIsEmpty }" :text="$t('Verify')" @click="verify"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="switch2FAVerifyType" :text="$t(twoFAVerifyTypeSwitchName)"></f7-link>
                    </div>
                </div>
            </div>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            username: '',
            password: '',
            passcode: '',
            backupCode: '',
            tempToken: '',
            twoFAVerifyType: 'passcode',
            twoFAVerifyTypeSwitchName: 'Use a backup code',
            allLanguages: self.$getAllLanguages()
        };
    },
    computed: {
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
            let lang = this.$getLanguage(currentLocale);

            if (!lang) {
                lang = this.$getLanguage(this.$getDefaultLanguage());
            }

            return lang.displayName;
        }
    },
    methods: {
        login() {
            const self = this;
            const app = self.$f7;
            const router = self.$f7router;

            if (!this.username) {
                self.$alert('Please input username');
                return;
            }

            if (!this.password) {
                self.$alert('Please input password');
                return;
            }

            if (self.tempToken) {
                app.sheet.open('#2fa-auth-sheet');
                return;
            }

            let hasResponse = false;

            setTimeout(() => {
                if (!hasResponse) {
                    app.preloader.show();
                }
            }, 200);

            self.$services.authorize({
                loginName: self.username,
                password: self.password
            }).then(response => {
                hasResponse = true;
                app.preloader.hide();
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.token) {
                    self.$alert('Unable to login');
                    return;
                }

                if (data.result.need2FA) {
                    self.tempToken = data.result.token;
                    app.sheet.open('#2fa-auth-sheet');
                    return;
                }

                self.$user.updateToken(data.result.token);
                router.navigate('/');
            }).catch(error => {
                hasResponse = true;
                app.preloader.hide();

                if (error && error.processed) {
                    return;
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert(`error.${error.response.data.errorMessage}`);
                } else {
                    self.$alert('Unable to login');
                }
            });
        },
        verify() {
            const self = this;
            const app = self.$f7;
            const router = self.$f7router;

            if (this.twoFAVerifyType === 'passcode' && !this.passcode) {
                self.$alert('Please input passcode');
                return;
            } else if (this.twoFAVerifyType === 'backupcode' && !this.backupCode) {
                self.$alert('Please input backup code');
                return;
            }

            app.preloader.show();

            let promise = null;

            if (self.twoFAVerifyType === 'backupcode') {
                promise = self.$services.authorize2FAByBackupCode({
                    recoveryCode: self.backupCode,
                    token: self.tempToken
                });
            } else {
                promise = self.$services.authorize2FA({
                    passcode: self.passcode,
                    token: self.tempToken
                });
            }

            promise.then(response => {
                app.preloader.hide();
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.token) {
                    self.$alert('Unable to verify');
                    return;
                }

                self.$user.updateToken(data.result.token);
                app.sheet.close('#2fa-auth-sheet');
                router.navigate('/');
            }).catch(error => {
                app.preloader.hide();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert(`error.${error.response.data.errorMessage}`);
                } else {
                    self.$alert('Unable to verify');
                }
            })
        },
        switch2FAVerifyType() {
            if (this.twoFAVerifyType === 'passcode') {
                this.twoFAVerifyType = 'backupcode';
                this.twoFAVerifyTypeSwitchName = 'Use a passcode';
            } else {
                this.twoFAVerifyType = 'passcode';
                this.twoFAVerifyTypeSwitchName = 'Use a backup code';
            }
        },
        changeLanguage(locale) {
            this.$setLanguage(locale);
        }
    }
};
</script>

<style scoped>
.twofa-auth-form {
    margin-top: 0;
    margin-bottom: 10px;
}
</style>
