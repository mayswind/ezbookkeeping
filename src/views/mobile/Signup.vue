<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Sign Up')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t('Submit')" @click="submit"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Username')"
                        :placeholder="$t('Your username')"
                        :value="user.username"
                        @input="user.username = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="password"
                        clear-button
                        :label="$t('Password')"
                        :placeholder="$t('Your password, at least 6 characters')"
                        :value="user.password"
                        @input="user.password = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="password"
                        clear-button
                        :label="$t('Confirmation Password')"
                        :placeholder="$t('Re-enter the password')"
                        :value="user.confirmPassword"
                        @input="user.confirmPassword = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="email"
                        clear-button
                        :label="$t('E-mail')"
                        :placeholder="$t('Your email address')"
                        :value="user.email"
                        @input="user.email = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Nickname')"
                        :placeholder="$t('Your nickname')"
                        :value="user.nickname"
                        @input="user.nickname = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="select"
                        :label="$t('Default Currency')"
                        :value="user.defaultCurrency"
                        @input="user.defaultCurrency = $event.target.value"
                    >
                        <option v-for="currency in allCurrencies"
                                :key="currency.code"
                                :value="currency.code">{{ currency.displayName }}</option>
                    </f7-list-input>

                    <f7-list-item class="lab-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            user: {
                username: '',
                password: '',
                confirmPassword: '',
                email: '',
                nickname: '',
                defaultCurrency: self.$t('default.currency')
            },
            submitting: false
        };
    },
    computed: {
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.user.username) {
                return 'Username cannot be empty';
            } else if (!this.user.password) {
                return 'Password cannot be empty';
            } else if (!this.user.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else if (!this.user.email) {
                return 'Email address cannot be empty';
            } else if (!this.user.nickname) {
                return 'Nickname cannot be empty';
            } else if (!this.user.defaultCurrency) {
                return 'Default currency cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.user.password && this.user.confirmPassword && this.user.password !== this.user.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else {
                return null;
            }
        }
    },
    methods: {
        submit() {
            const self = this;
            const router = self.$f7router;

            let problemMessage = self.inputEmptyProblemMessage || self.inputInvalidProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            self.$services.register({
                username: self.user.username,
                password: self.user.password,
                email: self.user.email,
                nickname: self.user.nickname,
                defaultCurrency: self.user.defaultCurrency
            }).then(response => {
                self.submitting = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to sign up');
                    return;
                }

                if (self.$settings.isEnableApplicationLock()) {
                    self.$settings.setEnableApplicationLock(false);
                    self.$settings.setEnableApplicationLockWebAuthn(false);
                    self.$user.clearWebAuthnConfig();
                }

                if (self.$utilities.isString(data.result.token)) {
                    self.$user.updateTokenAndUserInfo(data.result);

                    if (self.$settings.isAutoUpdateExchangeRatesData()) {
                        self.$services.autoRefreshLatestExchangeRates();
                    }
                }

                self.$toast('You have been successfully registered');
                router.navigate('/');
            }).catch(error => {
                self.$logger.error('failed to sign up', error);

                self.submitting = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to sign up');
                }
            });
        }
    }
};
</script>
