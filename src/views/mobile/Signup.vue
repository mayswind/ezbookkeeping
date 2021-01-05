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
                        autocomplete="username"
                        clear-button
                        :label="$t('Username')"
                        :placeholder="$t('Your username')"
                        :value="user.username"
                        @input="user.username = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="password"
                        autocomplete="new-password"
                        clear-button
                        :label="$t('Password')"
                        :placeholder="$t('Your password, at least 6 characters')"
                        :value="user.password"
                        @input="user.password = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="password"
                        autocomplete="new-password"
                        clear-button
                        :label="$t('Confirmation Password')"
                        :placeholder="$t('Re-enter the password')"
                        :value="user.confirmPassword"
                        @input="user.confirmPassword = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="email"
                        autocomplete="email"
                        clear-button
                        :label="$t('E-mail')"
                        :placeholder="$t('Your email address')"
                        :value="user.email"
                        @input="user.email = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="text"
                        autocomplete="nickname"
                        clear-button
                        :label="$t('Nickname')"
                        :placeholder="$t('Your nickname')"
                        :value="user.nickname"
                        @input="user.nickname = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item class="lab-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-item
                        :header="$t('Default Currency')"
                        smart-select :smart-select-params="{ openIn: 'popup', searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Close'), scrollToSelectedItem: true }"
                    >
                        <select autocomplete="transaction-currency" v-model="user.defaultCurrency">
                            <option v-for="currency in allCurrencies"
                                    :key="currency.code"
                                    :value="currency.code">{{ currency.displayName }}</option>
                        </select>
                    </f7-list-item>
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

            self.$store.dispatch('register', {
                user: self.user
            }).then(() => {
                self.submitting = false;
                self.$hideLoading();

                if (self.$user.isUserLogined()) {
                    if (self.$settings.isAutoUpdateExchangeRatesData()) {
                        self.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
                    }
                }

                self.$toast('You have been successfully registered');
                router.navigate('/');

                self.$confirm('Do you want to initialize transaction categories now?', () => {
                    router.navigate('/category/default?type=0');
                });
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    }
};
</script>
