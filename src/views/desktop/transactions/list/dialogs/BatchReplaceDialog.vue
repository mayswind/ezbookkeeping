<template>
    <v-dialog width="600" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'expenseCategory'">{{ $t('Batch Replace Selected Expense Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'incomeCategory'">{{ $t('Batch Replace Selected Income Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'transferCategory'">{{ $t('Batch Replace Selected Transfer Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'account'">{{ $t('Batch Replace Selected Accounts') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'destinationAccount'">{{ $t('Batch Replace Selected Destination Accounts') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'expenseCategory'">{{ $t('Replace Invalid Expense Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'incomeCategory'">{{ $t('Replace Invalid Income Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'transferCategory'">{{ $t('Replace Invalid Transfer Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'account'">{{ $t('Replace Invalid Accounts') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'tag'">{{ $t('Replace Invalid Transaction Tags') }}</h4>
                </div>
            </template>
            <v-card-text class="my-md-4 w-100 d-flex justify-center" v-if="type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory'">
                <v-row>
                    <v-col cols="12" v-if="mode === 'replaceInvalidItems'">
                        <v-autocomplete
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            :label="$t('Invalid Category')"
                            :placeholder="$t('Invalid Category')"
                            :items="invalidItems"
                            :no-data-text="$t('No available category')"
                            v-model="sourceItem">
                        </v-autocomplete>
                    </v-col>
                    <v-col cols="12">
                        <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :disabled="!hasAvailableExpenseCategories"
                                           :show-selection-primary-text="true"
                                           :custom-selection-primary-text="getPrimaryCategoryName(targetItem, allCategories[allCategoryTypes.Expense])"
                                           :custom-selection-secondary-text="getSecondaryCategoryName(targetItem, allCategories[allCategoryTypes.Expense])"
                                           :label="$t('Target Category')"
                                           :placeholder="$t('Target Category')"
                                           :items="allCategories[allCategoryTypes.Expense]"
                                           v-model="targetItem"
                                           v-if="type === 'expenseCategory'">
                        </two-column-select>
                        <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :disabled="!hasAvailableIncomeCategories"
                                           :show-selection-primary-text="true"
                                           :custom-selection-primary-text="getPrimaryCategoryName(targetItem, allCategories[allCategoryTypes.Income])"
                                           :custom-selection-secondary-text="getSecondaryCategoryName(targetItem, allCategories[allCategoryTypes.Income])"
                                           :label="$t('Target Category')"
                                           :placeholder="$t('Target Category')"
                                           :items="allCategories[allCategoryTypes.Income]"
                                           v-model="targetItem"
                                           v-if="type === 'incomeCategory'">
                        </two-column-select>
                        <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :disabled="!hasAvailableTransferCategories"
                                           :show-selection-primary-text="true"
                                           :custom-selection-primary-text="getPrimaryCategoryName(targetItem, allCategories[allCategoryTypes.Transfer])"
                                           :custom-selection-secondary-text="getSecondaryCategoryName(targetItem, allCategories[allCategoryTypes.Transfer])"
                                           :label="$t('Target Category')"
                                           :placeholder="$t('Target Category')"
                                           :items="allCategories[allCategoryTypes.Transfer]"
                                           v-model="targetItem"
                                           v-if="type === 'transferCategory'">
                        </two-column-select>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="my-md-4 w-100 d-flex justify-center" v-if="type === 'account' || type === 'destinationAccount'">
                <v-row>
                    <v-col cols="12" v-if="mode === 'replaceInvalidItems'">
                        <v-autocomplete
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            :label="$t('Invalid Account')"
                            :placeholder="$t('Invalid Account')"
                            :items="invalidItems"
                            :no-data-text="$t('No available account')"
                            v-model="sourceItem">
                        </v-autocomplete>
                    </v-col>
                    <v-col cols="12">
                        <two-column-select primary-key-field="id" primary-value-field="category"
                                           primary-title-field="name" primary-footer-field="displayBalance"
                                           primary-icon-field="icon" primary-icon-type="account"
                                           primary-sub-items-field="accounts"
                                           :primary-title-i18n="true"
                                           secondary-key-field="id" secondary-value-field="id"
                                           secondary-title-field="name" secondary-footer-field="displayBalance"
                                           secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                           :disabled="!allVisibleAccounts.length"
                                           :custom-selection-primary-text="getAccountDisplayName(targetItem)"
                                           :label="$t('Target Account')"
                                           :placeholder="$t('Target Account')"
                                           :items="allVisibleCategorizedAccounts"
                                           v-model="targetItem">
                        </two-column-select>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="my-md-4 w-100 d-flex justify-center" v-if="type === 'tag'">
                <v-row>
                    <v-col cols="12" v-if="mode === 'replaceInvalidItems'">
                        <v-autocomplete
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            :label="$t('Invalid Tag')"
                            :placeholder="$t('Invalid Tag')"
                            :items="invalidItems"
                            :no-data-text="$t('No available tag')"
                            v-model="sourceItem">
                        </v-autocomplete>
                    </v-col>
                    <v-col cols="12">
                        <v-autocomplete
                            item-title="name"
                            item-value="id"
                            persistent-placeholder
                            chips
                            :label="$t('Target Tag')"
                            :placeholder="$t('Target Tag')"
                            :items="allTags"
                            :no-data-text="$t('No available tag')"
                            v-model="targetItem"
                        >
                            <template #chip="{ props, item }">
                                <v-chip :prepend-icon="icons.tag" :text="item.title" v-bind="props"/>
                            </template>

                            <template #item="{ props, item }">
                                <v-list-item :value="item.value" v-bind="props" v-if="!item.raw.hidden">
                                    <template #title>
                                        <v-list-item-title>
                                            <div class="d-flex align-center">
                                                <v-icon size="20" start :icon="icons.tag"/>
                                                <span>{{ item.title }}</span>
                                            </div>
                                        </v-list-item-title>
                                    </template>
                                </v-list-item>
                            </template>
                        </v-autocomplete>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="(mode === 'replaceInvalidItems' && !sourceItem && sourceItem !== '') || (!targetItem && targetItem !== '')" @click="confirm">{{ $t('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ $t('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { CategoryType } from '@/core/category.ts';
import {
    getNameByKeyValue
} from '@/lib/common.ts';
import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName,
    getFirstAvailableCategoryId
} from '@/lib/category.js';

import {
    mdiPound
} from '@mdi/js';

export default {
    props: [
        'persistent'
    ],
    expose: [
        'open'
    ],
    data() {
        return {
            showState: false,
            mode: '',
            type: '',
            invalidItems: [],
            sourceItem: null,
            targetItem: null,
            icons: {
                tag: mdiPound
            }
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useExchangeRatesStore),
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        allCategoryTypes() {
            return CategoryType;
        },
        allAccounts() {
            return this.accountsStore.allPlainAccounts;
        },
        allVisibleCategorizedAccounts() {
            return this.$locale.getCategorizedAccountsWithDisplayBalance(this.allVisibleAccounts, this.showAccountBalance, this.defaultCurrency, this.settingsStore, this.userStore, this.exchangeRatesStore);
        },
        allVisibleAccounts() {
            return this.accountsStore.allVisiblePlainAccounts;
        },
        allCategories() {
            return this.transactionCategoriesStore.allTransactionCategories;
        },
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        allTagsMap() {
            return this.transactionTagsStore.allTransactionTagsMap;
        },
        hasAvailableExpenseCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Expense] || !this.allCategories[this.allCategoryTypes.Expense].length) {
                return false;
            }

            const firstAvailableCategoryId = getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Expense]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableIncomeCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Income] || !this.allCategories[this.allCategoryTypes.Income].length) {
                return false;
            }

            const firstAvailableCategoryId = getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Income]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableTransferCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Transfer] || !this.allCategories[this.allCategoryTypes.Transfer].length) {
                return false;
            }

            const firstAvailableCategoryId = getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Transfer]);
            return firstAvailableCategoryId !== '';
        },
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
        },
    },
    methods: {
        open(options) {
            const self = this;
            self.mode = options.mode;
            self.type = options.type;
            self.sourceItem = null;

            if (self.mode === 'batchReplace') {
                self.invalidItems = null;
            } else if (self.mode === 'replaceInvalidItems') {
                self.invalidItems = options.invalidItems;
            }

            self.targetItem = null;
            self.showState = true;

            return new Promise((resolve, reject) => {
                self.resolve = resolve;
                self.reject = reject;
            });
        },
        confirm() {
            if (this.resolve) {
                if (this.mode === 'batchReplace') {
                    this.resolve({
                        targetItem: this.targetItem
                    });
                } else if (this.mode === 'replaceInvalidItems') {
                    this.resolve({
                        sourceItem: this.sourceItem,
                        targetItem: this.targetItem
                    });
                }
            }

            this.showState = false;
        },
        cancel() {
            if (this.reject) {
                this.reject();
            }

            this.showState = false;
        },
        getPrimaryCategoryName(categoryId, allCategories) {
            return getTransactionPrimaryCategoryName(categoryId, allCategories);
        },
        getSecondaryCategoryName(categoryId, allCategories) {
            return getTransactionSecondaryCategoryName(categoryId, allCategories);
        },
        getAccountDisplayName(accountId) {
            if (accountId) {
                return getNameByKeyValue(this.allAccounts, accountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        }
    }
}
</script>
