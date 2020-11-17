<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('User Profile')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsNotChanged || inputIsInvalid || saving }" :text="$t('Save')" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-input label="Password" placeholder="Your password"></f7-list-input>
                    <f7-list-input label="Confirmation Password" placeholder="Re-enter the password"></f7-list-input>
                    <f7-list-input label="E-mail" placeholder="Your email address"></f7-list-input>
                    <f7-list-input label="Nickname" placeholder="Your nickname"></f7-list-input>
                    <f7-list-input label="Default Currency" placeholder="Default Currency"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-input
                        type="password"
                        clear-button
                        :label="$t('Password')"
                        :placeholder="$t('Your password')"
                        :value="newProfile.password"
                        @input="newProfile.password = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="password"
                        clear-button
                        :label="$t('Confirmation Password')"
                        :placeholder="$t('Re-enter the password')"
                        :value="newProfile.confirmPassword"
                        @input="newProfile.confirmPassword = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="email"
                        clear-button
                        :label="$t('E-mail')"
                        :placeholder="$t('Your email address')"
                        :value="newProfile.email"
                        @input="newProfile.email = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Nickname')"
                        :placeholder="$t('Your nickname')"
                        :value="newProfile.nickname"
                        @input="newProfile.nickname = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="select"
                        :label="$t('Default Currency')"
                        :value="newProfile.defaultCurrency"
                        @input="newProfile.defaultCurrency = $event.target.value"
                    >
                        <option v-for="currency in allCurrencies"
                                :key="currency.code"
                                :value="currency.code">{{ currency.displayName }}</option>
                    </f7-list-input>

                    <f7-list-item class="lab-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-sheet
            style="height:auto"
            :opened="showInputPasswordSheet" @sheet:closed="showInputPasswordSheet = false"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Current Password') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top">{{ $t('Please enter your current password when modifying your password') }}</p>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-input
                            type="password"
                            outline
                            clear-button
                            :placeholder="$t('Password')"
                            :value="currentPassword"
                            @input="currentPassword = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPassword || saving }" :text="$t('Continue')" @click="save"></f7-button>
                </div>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            newProfile: {
                password: '',
                confirmPassword: '',
                email: '',
                nickname: '',
                defaultCurrency: ''
            },
            oldProfile: {
                email: '',
                nickname: '',
                defaultCurrency: ''
            },
            currentPassword: '',
            loading: true,
            saving: false,
            showInputPasswordSheet: false
        };
    },
    computed: {
        allCurrencies() {
            return this.$getAllCurrencies();
        },
        inputIsNotChanged() {
            return !!this.inputIsNotChangedProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputIsNotChangedProblemMessage() {
            if (!this.newProfile.password && !this.newProfile.confirmPassword && !this.newProfile.email && !this.newProfile.nickname) {
                return 'Nothing has been modified';
            } else if (!this.newProfile.password && !this.newProfile.confirmPassword &&
                this.newProfile.email === this.oldProfile.email &&
                this.newProfile.nickname === this.oldProfile.nickname &&
                this.newProfile.defaultCurrency === this.oldProfile.defaultCurrency) {
                return 'Nothing has been modified';
            } else if (!this.newProfile.password && this.newProfile.confirmPassword) {
                return 'Password cannot be empty';
            } else if (this.newProfile.password && !this.newProfile.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.newProfile.password && this.newProfile.confirmPassword && this.newProfile.password !== this.newProfile.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else if (!this.newProfile.email) {
                return 'Email address cannot be empty';
            } else if (!this.newProfile.nickname) {
                return 'Nickname cannot be empty';
            } else if (!this.newProfile.defaultCurrency) {
                return 'Default currency cannot be empty';
            } else {
                return null;
            }
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        self.$services.getProfile().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$alert('Unable to get user profile', () => {
                    router.back();
                });
                return;
            }

            self.oldProfile.email = data.result.email;
            self.oldProfile.nickname = data.result.nickname;
            self.oldProfile.defaultCurrency = data.result.defaultCurrency;

            self.newProfile.email = self.oldProfile.email
            self.newProfile.nickname = self.oldProfile.nickname;
            self.newProfile.defaultCurrency = self.oldProfile.defaultCurrency;
            self.loading = false;
        }).catch(error => {
            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$alert({ error: error.response.data }, () => {
                    router.back();
                });
            } else if (!error.processed) {
                self.$alert('Unable to get user profile', () => {
                    router.back();
                });
            }
        });
    },
    methods: {
        save() {
            const self = this;
            const router = self.$f7router;

            self.showInputPasswordSheet = false;

            let problemMessage = self.inputIsNotChangedProblemMessage || self.inputInvalidProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            if (self.newProfile.password && !self.currentPassword) {
                self.showInputPasswordSheet = true;
                return;
            }

            self.saving = true;
            self.$showLoading(() => self.saving);

            self.$services.updateProfile({
                password: self.newProfile.password,
                oldPassword: self.currentPassword,
                email: self.newProfile.email,
                nickname: self.newProfile.nickname,
                defaultCurrency: self.newProfile.defaultCurrency
            }).then(response => {
                self.saving = false;
                self.$hideLoading();
                self.currentPassword = '';

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$alert('Unable to update user profile');
                    return;
                }

                if (self.$utilities.isString(data.result.newToken)) {
                    self.$user.updateToken(data.result.newToken);
                }

                if (self.$utilities.isObject(data.result.user)) {
                    self.$user.updateUserInfo(data.result.user);
                }

                self.$toast('Your profile has been successfully updated');
                router.back('/settings', { force: true });
            }).catch(error => {
                self.saving = false;
                self.$hideLoading();
                self.currentPassword = '';

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else if (!error.processed) {
                    self.$alert('Unable to update user profile');
                }
            });
        }
    }
};
</script>
