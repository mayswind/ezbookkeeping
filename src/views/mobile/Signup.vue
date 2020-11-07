<template>
    <f7-page>
        <f7-navbar :title="$t('Sign Up')" :back-link="$t('Back')"></f7-navbar>
        <f7-list no-hairlines-md>
            <f7-list-input
                type="text"
                clear-button
                :label="$t('Username')"
                :placeholder="$t('Your username')"
                :value="username"
                @input="username = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password, at least 6 characters')"
                :value="password"
                @input="password = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="password"
                clear-button
                :label="$t('Confirmation Password')"
                :placeholder="$t('Re-enter the password')"
                :value="confirmPassword"
                @input="confirmPassword = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="email"
                clear-button
                :label="$t('E-mail')"
                :placeholder="$t('Your email address')"
                :value="email"
                @input="email = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="text"
                clear-button
                :label="$t('Nickname')"
                :placeholder="$t('Your nickname')"
                :value="nickname"
                @input="nickname = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="select"
                :label="$t('Default Currency')"
                :value="defaultCurrency"
                @input="defaultCurrency = $event.target.value"
            >
                <option v-for="currency in allCurrencies"
                        :key="currency.code"
                        :value="currency.code">{{ currency.displayName }}</option>
            </f7-list-input>

            <f7-list-item class="lab-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <f7-button large fill :class="{ 'disabled': inputIsEmpty || signuping }" :text="$t('Sign Up')" @click="signup"></f7-button>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            username: '',
            password: '',
            confirmPassword: '',
            email: '',
            nickname: '',
            defaultCurrency: self.$t('default.currency'),
            signuping: false,
            allCurrencies: self.$getAllCurrencies()
        };
    },
    computed: {
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.username) {
                return 'Username cannot be empty';
            } else if (!this.password) {
                return 'Password cannot be empty';
            } else if (!this.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else if (!this.email) {
                return 'Email address cannot be empty';
            } else if (!this.nickname) {
                return 'Nickname cannot be empty';
            } else if (!this.defaultCurrency) {
                return 'Default currency cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.password && this.confirmPassword && this.password !== this.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else {
                return null;
            }
        }
    },
    methods: {
        signup() {
            const self = this;
            const router = self.$f7router;

            let problemMessage = self.inputEmptyProblemMessage || self.inputInvalidProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.signuping = true;
            self.$showLoading(() => self.signuping);

            self.$services.register({
                username: self.username,
                password: self.password,
                email: self.email,
                nickname: self.nickname,
                defaultCurrency: self.defaultCurrency
            }).then(response => {
                self.signuping = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$alert('Unable to sign up');
                    return;
                }

                if (typeof(data.result) === 'object') {
                    self.$user.updateToken(data.result);
                }

                self.$toast('You have been successfully registered');
                router.navigate('/');
            }).catch(error => {
                self.signuping = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else if (!error.processed) {
                    self.$alert('Unable to sign up');
                }
            });
        }
    }
};
</script>
