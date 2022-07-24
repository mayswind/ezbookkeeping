<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('User Profile')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsNotChanged || inputIsInvalid || saving }" :text="$t('Save')" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-input label="Password" placeholder="Your password"></f7-list-input>
                    <f7-list-input label="Confirmation Password" placeholder="Re-enter the password"></f7-list-input>
                    <f7-list-input label="E-mail" placeholder="Your email address"></f7-list-input>
                    <f7-list-input label="Nickname" placeholder="Your nickname"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Currency" title="Currency"></f7-list-item>
                    <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="First Day of Week" title="Week Day"></f7-list-item>
                    <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Editable Transaction Scope" title="All"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-input
                        type="password"
                        autocomplete="new-password"
                        clear-button
                        :label="$t('Password')"
                        :placeholder="$t('Your password')"
                        :value="newProfile.password"
                        @input="newProfile.password = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="password"
                        autocomplete="new-password"
                        clear-button
                        :label="$t('Confirmation Password')"
                        :placeholder="$t('Re-enter the password')"
                        :value="newProfile.confirmPassword"
                        @input="newProfile.confirmPassword = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="email"
                        autocomplete="email"
                        clear-button
                        :label="$t('E-mail')"
                        :placeholder="$t('Your email address')"
                        :value="newProfile.email"
                        @input="newProfile.email = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="text"
                        autocomplete="nickname"
                        clear-button
                        :label="$t('Nickname')"
                        :placeholder="$t('Your nickname')"
                        :value="newProfile.nickname"
                        @input="newProfile.nickname = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item class="ebk-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :header="$t('Default Currency')"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('Default Currency'), searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <f7-block slot="title" class="no-padding no-margin">
                            <span>{{ $t(`currency.${newProfile.defaultCurrency}`) }}&nbsp;</span>
                            <small class="smaller">{{ newProfile.defaultCurrency }}</small>
                        </f7-block>
                        <select autocomplete="transaction-currency" v-model="newProfile.defaultCurrency">
                            <option v-for="currency in allCurrencies"
                                    :key="currency.code"
                                    :value="currency.code">{{ currency.displayName }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        class="list-item-with-header-and-title"
                        link="#"
                        :class="{ 'disabled': !allVisibleAccounts.length }"
                        :header="$t('Default Account')"
                        :title="newProfile.defaultAccountId | optionName(allAccounts, 'id', 'name', $t('Not Specified'))"
                        @click="showAccountSheet = true"
                    >
                        <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                              primary-title-field="name"
                                                              primary-icon-field="icon" primary-icon-type="account"
                                                              primary-sub-items-field="accounts"
                                                              :primary-title-i18n="true"
                                                              secondary-key-field="id" secondary-value-field="id"
                                                              secondary-title-field="name"
                                                              secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                              :items="allCategorizedAccounts"
                                                              :show.sync="showAccountSheet"
                                                              v-model="newProfile.defaultAccountId">
                        </two-column-list-item-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :header="$t('First Day of Week')"
                        :title="newProfile.firstDayOfWeek | optionName(allWeekDays, 'type', 'name') | format('datetime.#{value}.long') | localized"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('First Day of Week'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <select v-model="newProfile.firstDayOfWeek">
                            <option v-for="weekDay in allWeekDays"
                                    :key="weekDay.type"
                                    :value="weekDay.type">{{ $t(`datetime.${weekDay.name}.long`) }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :header="$t('Editable Transaction Scope')"
                        :title="newProfile.transactionEditScope | optionName(allTransactionEditScopeTypes, 'value', 'name') | localized"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('Editable Transaction Scope'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <select v-model="newProfile.transactionEditScope">
                            <option v-for="option in allTransactionEditScopeTypes"
                                    :key="option.value"
                                    :value="option.value">{{ $t(option.name) }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item class="ebk-list-item-error-info" v-if="extendInputIsInvalid" :footer="$t(extendInputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <password-input-sheet :title="$t('Current Password')"
                              :hint="$t('Please enter your current password when modifying your password')"
                              :show.sync="showInputPasswordSheet"
                              :confirm-disabled="saving"
                              :cancel-disabled="saving"
                              v-model="currentPassword"
                              @password:confirm="save()">
        </password-input-sheet>
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
                defaultCurrency: '',
                defaultAccountId: '',
                firstDayOfWeek: 0,
                transactionEditScope: 1
            },
            oldProfile: {
                email: '',
                nickname: '',
                defaultCurrency: '',
                defaultAccountId: '',
                firstDayOfWeek: 0,
                transactionEditScope: 1
            },
            currentPassword: '',
            loading: true,
            loadingError: null,
            saving: false,
            showInputPasswordSheet: false,
            showAccountSheet: false
        };
    },
    computed: {
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allAccounts() {
            return this.$store.getters.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.$store.getters.allVisiblePlainAccounts;
        },
        allCategorizedAccounts() {
            return this.$utilities.getCategorizedAccounts(this.allVisibleAccounts);
        },
        allWeekDays() {
            return this.$constants.datetime.allWeekDays;
        },
        allTransactionEditScopeTypes() {
            return [{
                value: 0,
                name: 'None'
            }, {
                value: 1,
                name: 'All'
            }, {
                value: 2,
                name: 'Today or later'
            }, {
                value: 3,
                name: 'Recent 24 hours or later'
            }, {
                value: 4,
                name: 'This week or later'
            }, {
                value: 5,
                name: 'This month or later'
            }, {
                value: 6,
                name: 'This year or later'
            }];
        },
        inputIsNotChanged() {
            return !!this.inputIsNotChangedProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        extendInputIsInvalid() {
            return !!this.extendInputInvalidProblemMessage;
        },
        inputIsNotChangedProblemMessage() {
            if (!this.newProfile.password && !this.newProfile.confirmPassword && !this.newProfile.email && !this.newProfile.nickname) {
                return 'Nothing has been modified';
            } else if (!this.newProfile.password && !this.newProfile.confirmPassword &&
                this.newProfile.email === this.oldProfile.email &&
                this.newProfile.nickname === this.oldProfile.nickname &&
                this.newProfile.defaultCurrency === this.oldProfile.defaultCurrency &&
                this.newProfile.defaultAccountId === this.oldProfile.defaultAccountId &&
                this.newProfile.firstDayOfWeek === this.oldProfile.firstDayOfWeek &&
                this.newProfile.transactionEditScope === this.oldProfile.transactionEditScope) {
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
        },
        extendInputInvalidProblemMessage() {
            if (!this.newProfile.defaultCurrency) {
                return 'Default currency cannot be empty';
            } else {
                return null;
            }
        }
    },
    created() {
        const self = this;

        self.loading = true;

        const promises = [
            self.$store.dispatch('loadAllAccounts', { force: false }),
            self.$store.dispatch('getCurrentUserProfile')
        ];

        Promise.all(promises).then(responses => {
            const profile = responses[1];

            self.oldProfile.email = profile.email;
            self.oldProfile.nickname = profile.nickname;
            self.oldProfile.defaultCurrency = profile.defaultCurrency;
            self.oldProfile.defaultAccountId = profile.defaultAccountId;
            self.oldProfile.firstDayOfWeek = profile.firstDayOfWeek;
            self.oldProfile.transactionEditScope = profile.transactionEditScope;

            self.newProfile.email = self.oldProfile.email
            self.newProfile.nickname = self.oldProfile.nickname;
            self.newProfile.defaultCurrency = self.oldProfile.defaultCurrency;
            self.newProfile.defaultAccountId = self.oldProfile.defaultAccountId;
            self.newProfile.firstDayOfWeek = self.oldProfile.firstDayOfWeek;
            self.newProfile.transactionEditScope = self.oldProfile.transactionEditScope;

            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError('loadingError');
        },
        save() {
            const self = this;
            const router = self.$f7router;

            self.showInputPasswordSheet = false;

            let problemMessage = self.inputIsNotChangedProblemMessage || self.inputInvalidProblemMessage || self.extendInputInvalidProblemMessage;

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

            self.$store.dispatch('updateUserProfile', {
                profile: self.newProfile,
                currentPassword: self.currentPassword
            }).then(() => {
                self.saving = false;
                self.$hideLoading();
                self.currentPassword = '';

                self.$toast('Your profile has been successfully updated');
                router.back();
            }).catch(error => {
                self.saving = false;
                self.$hideLoading();
                self.currentPassword = '';

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    }
};
</script>
