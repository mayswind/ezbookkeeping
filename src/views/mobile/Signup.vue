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

                    <f7-list-item class="ebk-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :key="currentLocale + '_lang'"
                        :header="$t('Language')"
                        :title="currentLocale | languageName"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('Language'), searchbar: true, searchbarPlaceholder: $t('Language'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <select v-model="currentLocale">
                            <option v-for="(lang, locale) in allLanguages"
                                    :key="locale"
                                    :value="locale">{{ lang.displayName }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :key="currentLocale + '_currency'"
                        :header="$t('Default Currency')"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('Default Currency'), searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <f7-block slot="title" class="no-padding no-margin">
                            <span>{{ $t(`currency.${user.defaultCurrency}`) }}&nbsp;</span>
                            <small class="smaller">{{ user.defaultCurrency }}</small>
                        </f7-block>
                        <select autocomplete="transaction-currency" v-model="user.defaultCurrency">
                            <option v-for="currency in allCurrencies"
                                    :key="currency.code"
                                    :value="currency.code">{{ currency.displayName }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :key="currentLocale + '_firstDayOfWeek'"
                        :header="$t('First Day of Week')"
                        :title="user.firstDayOfWeek | optionName(allWeekDays, 'type', 'name') | format('datetime.#{value}.long') | localized"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('First Day of Week'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <select v-model="user.firstDayOfWeek">
                            <option v-for="weekDay in allWeekDays"
                                    :key="weekDay.type"
                                    :value="weekDay.type">{{ $t(`datetime.${weekDay.name}.long`) }}</option>
                        </select>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-item :title="$t('Use preset transaction categories')" link="#" @click="showPresetCategories = true">
                        <f7-toggle :checked="usePresetCategories" @toggle:change="usePresetCategories = $event"></f7-toggle>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-popup :opened="showPresetCategories" @popup:closed="showPresetCategories = false">
            <f7-page>
                <f7-navbar>
                    <f7-nav-left>
                        <f7-link popup-close :text="$t('Cancel')"></f7-link>
                    </f7-nav-left>
                    <f7-nav-title :title="$t('Preset Categories')"></f7-nav-title>
                    <f7-nav-right>
                        <f7-link icon-f7="ellipsis" @click="showPresetCategoriesMoreActionSheet = true"></f7-link>
                        <f7-link close @click="usePresetCategories = true; showPresetCategories = false" v-if="!usePresetCategories">{{ $t('Enable') }}</f7-link>
                        <f7-link close @click="usePresetCategories = false; showPresetCategories = false" v-if="usePresetCategories">{{ $t('Disable') }}</f7-link>
                    </f7-nav-right>
                </f7-navbar>
                <f7-card v-for="(categories, categoryType) in presetCategories" :key="categoryType">
                    <f7-card-header>
                        <small class="card-header-content">
                            <span>{{ categoryType | categoryTypeName($constants.category.allCategoryTypes) | localized }}</span>
                        </small>
                    </f7-card-header>
                    <f7-card-content class="no-safe-areas" :padding="false">
                        <f7-list v-if="showPresetCategories">
                            <f7-list-item v-for="(category, idx) in categories"
                                          :key="idx"
                                          :accordion-item="!!category.subCategories.length"
                                          :title="$t('category.' + category.name, currentLocale)">
                                <f7-icon slot="media"
                                         :icon="category.categoryIconId | categoryIcon"
                                         :style="category.color | categoryIconStyle('var(--default-icon-color)')">
                                </f7-icon>

                                <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                                    <f7-list>
                                        <f7-list-item v-for="(subCategory, subIdx) in category.subCategories"
                                                      :key="subIdx"
                                                      :title="$t('category.' + subCategory.name, currentLocale)">
                                            <f7-icon slot="media"
                                                     :icon="subCategory.categoryIconId | categoryIcon"
                                                     :style="subCategory.color | categoryIconStyle('var(--default-icon-color)')">
                                            </f7-icon>
                                        </f7-list-item>
                                    </f7-list>
                                </f7-accordion-content>
                            </f7-list-item>
                        </f7-list>
                    </f7-card-content>
                </f7-card>
            </f7-page>

            <f7-actions close-by-outside-click close-on-escape :opened="showPresetCategoriesMoreActionSheet" @actions:closed="showPresetCategoriesMoreActionSheet = false">
                <f7-actions-group>
                    <f7-actions-button @click="showPresetCategoriesChangeLocaleSheet = true">{{ $t('Change Language') }}</f7-actions-button>
                </f7-actions-group>
                <f7-actions-group>
                    <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
                </f7-actions-group>
            </f7-actions>

            <list-item-selection-sheet value-type="index"
                                       title-field="displayName"
                                       :items="allLanguages"
                                       :show.sync="showPresetCategoriesChangeLocaleSheet"
                                       v-model="currentLocale">
            </list-item-selection-sheet>
        </f7-popup>
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
                defaultCurrency: self.$store.state.defaultSetting.currency,
                firstDayOfWeek: self.$constants.datetime.allWeekDays[self.$t('default.firstDayOfWeek')] ? self.$constants.datetime.allWeekDays[self.$t('default.firstDayOfWeek')].type : 0
            },
            submitting: false,
            presetCategories: {
                [self.$constants.category.allCategoryTypes.Income]: self.$utilities.copyArrayTo(self.$constants.category.defaultIncomeCategories, []),
                [self.$constants.category.allCategoryTypes.Expense]: self.$utilities.copyArrayTo(self.$constants.category.defaultExpenseCategories, []),
                [self.$constants.category.allCategoryTypes.Transfer]: self.$utilities.copyArrayTo(self.$constants.category.defaultTransferCategories, [])
            },
            usePresetCategories: false,
            showPresetCategories: false,
            showPresetCategoriesMoreActionSheet: false,
            showPresetCategoriesChangeLocaleSheet: false
        };
    },
    computed: {
        allLanguages() {
            return this.$locale.getAllLanguages();
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allWeekDays() {
            return this.$constants.datetime.allWeekDays;
        },
        currentLocale: {
            get: function () {
                return this.$i18n.locale;
            },
            set: function (value) {
                const isCurrencyDefault = this.user.defaultCurrency === this.$store.state.defaultSetting.currency;
                const isFirstWeekDayDefault = this.user.firstDayOfWeek === (this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')] ? this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')].type : 0);

                this.$locale.setLanguage(value);

                if (isCurrencyDefault) {
                    this.user.defaultCurrency = this.$store.state.defaultSetting.currency;
                }

                if (isFirstWeekDayDefault) {
                    this.user.firstDayOfWeek = this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')] ? this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')].type : 0;
                }
            }
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

            const allCategories = [];

            if (self.usePresetCategories) {
                for (let categoryType in self.presetCategories) {
                    if (!Object.prototype.hasOwnProperty.call(self.presetCategories, categoryType)) {
                        continue;
                    }

                    const categories = self.presetCategories[categoryType];

                    for (let j = 0; j < categories.length; j++) {
                        const category = categories[j];
                        const submitCategory = {
                            name: self.$t('category.' + category.name, self.currentLocale),
                            type: parseInt(categoryType),
                            icon: category.categoryIconId,
                            color: category.color,
                            subCategories: []
                        }

                        for (let k = 0; k < category.subCategories.length; k++) {
                            const subCategory = category.subCategories[k];
                            submitCategory.subCategories.push({
                                name: self.$t('category.' + subCategory.name, self.currentLocale),
                                type: parseInt(categoryType),
                                icon: subCategory.categoryIconId,
                                color: subCategory.color
                            });
                        }

                        allCategories.push(submitCategory);
                    }
                }
            }

            self.$store.dispatch('register', {
                user: self.user
            }).then(() => {
                if (!self.$user.isUserLogined()) {
                    self.submitting = false;
                    self.$hideLoading();

                    if (self.usePresetCategories) {
                        self.$toast('You have been successfully registered, but something wrong with adding preset categories. You can re-add preset categories in settings page anytime.');
                    } else {
                        self.$toast('You have been successfully registered');
                    }

                    router.navigate('/');
                    return;
                }

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
                }

                if (!self.usePresetCategories) {
                    self.submitting = false;
                    self.$hideLoading();

                    self.$toast('You have been successfully registered');
                    router.navigate('/');
                    return;
                }

                self.$store.dispatch('addCategories', {
                    categories: allCategories
                }).then(() => {
                    self.submitting = false;
                    self.$hideLoading();

                    self.$toast('You have been successfully registered');
                    router.navigate('/');
                }).catch(() => {
                    self.submitting = false;
                    self.$hideLoading();

                    self.$toast('You have been successfully registered, but something wrong with adding preset categories. You can re-add preset categories in settings page anytime.');
                    router.navigate('/');
                });
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    },
    filters: {
        categoryTypeName(categoryType, allCategoryTypes) {
            switch (categoryType) {
                case allCategoryTypes.Income.toString():
                    return 'Income Categories';
                case allCategoryTypes.Expense.toString():
                    return 'Expense Categories';
                case allCategoryTypes.Transfer.toString():
                    return 'Transfer Categories';
                default:
                    return 'Transaction Categories';
            }
        }
    }
};
</script>
