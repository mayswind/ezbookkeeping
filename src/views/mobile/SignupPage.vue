<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Sign Up')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t('Submit')" @click="submit"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list form strong inset dividers class="margin-vertical">
            <f7-list-input
                type="text"
                autocomplete="username"
                clear-button
                :label="$t('Username')"
                :placeholder="$t('Your username')"
                v-model:value="user.username"
            ></f7-list-input>

            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password, at least 6 characters')"
                v-model:value="user.password"
            ></f7-list-input>

            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="$t('Confirmation Password')"
                :placeholder="$t('Re-enter the password')"
                v-model:value="user.confirmPassword"
            ></f7-list-input>

            <f7-list-input
                type="email"
                autocomplete="email"
                clear-button
                :label="$t('E-mail')"
                :placeholder="$t('Your email address')"
                v-model:value="user.email"
            ></f7-list-input>

            <f7-list-input
                type="text"
                autocomplete="nickname"
                clear-button
                :label="$t('Nickname')"
                :placeholder="$t('Your nickname')"
                v-model:value="user.nickname"
            ></f7-list-input>

            <f7-list-item class="ebk-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical">
            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :key="currentLocale + '_lang'"
                :header="$t('Language')"
                :title="currentLanguageName"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Language'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Language'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="currentLocale">
                    <option :value="locale"
                            :key="locale"
                            v-for="(lang, locale) in allLanguages">{{ lang.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :key="currentLocale + '_currency'"
                :header="$t('Default Currency')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Default Currency'), popupCloseLinkText: $t('Done') }"
            >
                <template #title>
                    <f7-block class="no-padding no-margin">
                        <span>{{ $t(`currency.${user.defaultCurrency}`) }}&nbsp;</span>
                        <small class="smaller">{{ user.defaultCurrency }}</small>
                    </f7-block>
                </template>
                <select autocomplete="transaction-currency" v-model="user.defaultCurrency">
                    <option :value="currency.code"
                            :key="currency.code"
                            v-for="currency in allCurrencies">{{ currency.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :key="currentLocale + '_firstDayOfWeek'"
                :header="$t('First Day of Week')"
                :title="getDayOfWeekName(user.firstDayOfWeek)"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('First Day of Week'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="user.firstDayOfWeek">
                    <option :value="weekDay.type"
                            :key="weekDay.type"
                            v-for="weekDay in allWeekDays">{{ $t(`datetime.${weekDay.name}.long`) }}</option>
                </select>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical">
            <f7-list-item :title="$t('Use preset transaction categories')" link="#" @click="showPresetCategories = true">
                <f7-toggle :checked="usePresetCategories" @toggle:change="usePresetCategories = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-popup push :close-on-escape="false" :opened="showPresetCategories"
                  @popup:closed="showPresetCategories = false">
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
                <f7-block class="no-padding no-margin"
                          :key="categoryType" v-for="(categories, categoryType) in presetCategories">
                    <f7-block-title class="margin-top margin-horizontal">{{ getCategoryTypeName(categoryType) }}</f7-block-title>
                    <f7-list strong inset dividers v-if="showPresetCategories">
                        <f7-list-item :title="$t('category.' + category.name)"
                                      :accordion-item="!!category.subCategories.length"
                                      :key="idx"
                                      v-for="(category, idx) in categories">
                            <template #media>
                                <ItemIcon icon-type="category" :icon-id="category.categoryIconId" :color="category.color"></ItemIcon>
                            </template>

                            <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                                <f7-list>
                                    <f7-list-item :title="$t('category.' + subCategory.name)"
                                                  :key="subIdx"
                                                  v-for="(subCategory, subIdx) in category.subCategories">
                                        <template #media>
                                            <ItemIcon icon-type="category" :icon-id="subCategory.categoryIconId" :color="subCategory.color"></ItemIcon>
                                        </template>
                                    </f7-list-item>
                                </f7-list>
                            </f7-accordion-content>
                        </f7-list-item>
                    </f7-list>
                </f7-block>
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
                                       v-model:show="showPresetCategoriesChangeLocaleSheet"
                                       v-model="currentLocale">
            </list-item-selection-sheet>
        </f7-popup>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import datetimeConstants from '@/consts/datetime.js';
import categoryConstants from '@/consts/category.js';
import { getNameByKeyValue, copyArrayTo } from '@/lib/common.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        const self = this;
        const settingsStore = useSettingsStore();

        return {
            user: {
                username: '',
                password: '',
                confirmPassword: '',
                email: '',
                nickname: '',
                language: self.$i18n.locale,
                defaultCurrency: settingsStore.defaultSetting.currency,
                firstDayOfWeek: datetimeConstants.allWeekDays[self.$t('default.firstDayOfWeek')] ? datetimeConstants.allWeekDays[self.$t('default.firstDayOfWeek')].type : 0
            },
            submitting: false,
            presetCategories: {
                [categoryConstants.allCategoryTypes.Income]: copyArrayTo(categoryConstants.defaultIncomeCategories, []),
                [categoryConstants.allCategoryTypes.Expense]: copyArrayTo(categoryConstants.defaultExpenseCategories, []),
                [categoryConstants.allCategoryTypes.Transfer]: copyArrayTo(categoryConstants.defaultTransferCategories, [])
            },
            usePresetCategories: false,
            showPresetCategories: false,
            showPresetCategoriesMoreActionSheet: false,
            showPresetCategoriesChangeLocaleSheet: false
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useTransactionCategoriesStore, useExchangeRatesStore),
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allWeekDays() {
            return datetimeConstants.allWeekDays;
        },
        currentLocale: {
            get: function () {
                return this.$i18n.locale;
            },
            set: function (value) {
                const isCurrencyDefault = this.user.defaultCurrency === this.settingsStore.defaultSetting.currency;
                const isFirstWeekDayDefault = this.user.firstDayOfWeek === (datetimeConstants.allWeekDays[this.$t('default.firstDayOfWeek')] ? datetimeConstants.allWeekDays[this.$t('default.firstDayOfWeek')].type : 0);

                this.user.language = value;

                const localeDefaultSettings = this.$locale.setLanguage(value);
                this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                if (isCurrencyDefault) {
                    this.user.defaultCurrency = this.settingsStore.defaultSetting.currency;
                }

                if (isFirstWeekDayDefault) {
                    this.user.firstDayOfWeek = datetimeConstants.allWeekDays[this.$t('default.firstDayOfWeek')] ? datetimeConstants.allWeekDays[this.$t('default.firstDayOfWeek')].type : 0;
                }
            }
        },
        currentLanguageName() {
            const languageInfo = this.$locale.getLanguageInfo(this.currentLocale);

            if (!languageInfo) {
                return '';
            }

            return languageInfo.displayName;
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
            const router = self.f7router;

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
                            name: self.$t('category.' + category.name),
                            type: parseInt(categoryType),
                            icon: category.categoryIconId,
                            color: category.color,
                            subCategories: []
                        }

                        for (let k = 0; k < category.subCategories.length; k++) {
                            const subCategory = category.subCategories[k];
                            submitCategory.subCategories.push({
                                name: self.$t('category.' + subCategory.name),
                                type: parseInt(categoryType),
                                icon: subCategory.categoryIconId,
                                color: subCategory.color
                            });
                        }

                        allCategories.push(submitCategory);
                    }
                }
            }

            self.rootStore.register({
                user: self.user
            }).then(response => {
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

                if (response.user && response.user.language) {
                    const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                if (!self.usePresetCategories) {
                    self.submitting = false;
                    self.$hideLoading();

                    self.$toast('You have been successfully registered');
                    router.navigate('/');
                    return;
                }

                self.transactionCategoriesStore.addCategories({
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
        },
        getDayOfWeekName(dayOfWeek) {
            const weekName = getNameByKeyValue(datetimeConstants.allWeekDays, dayOfWeek, 'type', 'name');
            const i18nWeekNameKey = `datetime.${weekName}.long`;
            return this.$t(i18nWeekNameKey);
        },
        getCategoryTypeName(categoryType) {
            switch (categoryType) {
                case categoryConstants.allCategoryTypes.Income.toString():
                    return this.$t('Income Categories');
                case categoryConstants.allCategoryTypes.Expense.toString():
                    return this.$t('Expense Categories');
                case categoryConstants.allCategoryTypes.Transfer.toString():
                    return this.$t('Transfer Categories');
                default:
                    return this.$t('Transaction Categories');
            }
        }
    }
};
</script>

