<template>
    <f7-page with-subnavbar @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true" v-if="mode !== 'view'"></f7-link>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t(saveButtonTitle)" @click="save" v-if="mode !== 'view'"></f7-link>
            </f7-nav-right>

            <f7-subnavbar>
                <f7-segmented strong :class="{ 'readonly': mode !== 'add' }">
                    <f7-button :text="$t('Expense')" :active="transaction.type === allTransactionTypes.Expense" @click="transaction.type = allTransactionTypes.Expense"></f7-button>
                    <f7-button :text="$t('Income')" :active="transaction.type === allTransactionTypes.Income" @click="transaction.type = allTransactionTypes.Income"></f7-button>
                    <f7-button :text="$t('Transfer')" :active="transaction.type === allTransactionTypes.Transfer" @click="transaction.type = allTransactionTypes.Transfer"></f7-button>
                </f7-segmented>
            </f7-subnavbar>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item
                class="transaction-edit-amount ebk-large-amount"
                header="Expense Amount" title="0.00">
            </f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow" header="Category" title="Category Names"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Account" title="Account Name"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Transaction Time" title="YYYY/MM/DD HH:mm:ss"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow list-item-no-item-after" header="Transaction Time Zone" title="(UTC XX:XX) System Default" link="#" :no-chevron="mode === 'view'"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow" header="Geographic Location" title="No Location"></f7-list-item>
            <f7-list-item header="Tags">
                <template #footer>
                    <f7-block class="margin-top-half no-padding no-margin">
                        <f7-chip class="transaction-edit-tag" text="None"></f7-chip>
                    </f7-block>
                </template>
            </f7-list-item>
            <f7-list-input type="textarea" label="Description" placeholder="Your transaction description (optional)"></f7-list-input>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical"
                 v-else-if="!loading">
            <f7-list-item
                class="transaction-edit-amount"
                link="#" no-chevron
                :class="sourceAmountClass"
                :header="$t(sourceAmountName)"
                :title="getDisplayAmount(transaction.sourceAmount, transaction.hideAmount)"
                @click="showSourceAmountSheet = true"
            >
                <number-pad-sheet :min-value="allowedMinAmount"
                                  :max-value="allowedMaxAmount"
                                  v-model:show="showSourceAmountSheet"
                                  v-model="transaction.sourceAmount"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="transaction-edit-amount"
                link="#" no-chevron
                :class="destinationAmountClass"
                :header="$t('Transfer In Amount')"
                :title="getDisplayAmount(transaction.destinationAmount, transaction.hideAmount)"
                @click="showDestinationAmountSheet = true"
                v-if="transaction.type === allTransactionTypes.Transfer"
            >
                <number-pad-sheet :min-value="allowedMinAmount"
                                  :max-value="allowedMaxAmount"
                                  v-model:show="showDestinationAmountSheet"
                                  v-model="transaction.destinationAmount"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                key="expenseCategorySelection"
                link="#" no-chevron
                :class="{ 'disabled': !hasAvailableExpenseCategories, 'readonly': mode === 'view' }"
                :header="$t('Category')"
                @click="showCategorySheet = true"
                v-if="transaction.type === allTransactionTypes.Expense"
            >
                <template #title>
                    <div class="list-item-custom-title" v-if="hasAvailableExpenseCategories">
                        <span>{{ getPrimaryCategoryName(transaction.expenseCategory, allCategories[allCategoryTypes.Expense]) }}</span>
                        <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                        <span>{{ getSecondaryCategoryName(transaction.expenseCategory, allCategories[allCategoryTypes.Expense]) }}</span>
                    </div>
                    <div class="list-item-custom-title" v-else-if="!hasAvailableExpenseCategories">
                        <span>{{ $t('None') }}</span>
                    </div>
                </template>
                <tree-view-selection-sheet primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           :items="allCategories[allCategoryTypes.Expense]"
                                           v-model:show="showCategorySheet"
                                           v-model="transaction.expenseCategory">
                </tree-view-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                key="incomeCategorySelection"
                link="#" no-chevron
                :class="{ 'disabled': !hasAvailableIncomeCategories, 'readonly': mode === 'view' }"
                :header="$t('Category')"
                @click="showCategorySheet = true"
                v-if="transaction.type === allTransactionTypes.Income"
            >
                <template #title>
                    <div class="list-item-custom-title" v-if="hasAvailableIncomeCategories">
                        <span>{{ getPrimaryCategoryName(transaction.incomeCategory, allCategories[allCategoryTypes.Income]) }}</span>
                        <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                        <span>{{ getSecondaryCategoryName(transaction.incomeCategory, allCategories[allCategoryTypes.Income]) }}</span>
                    </div>
                    <div class="list-item-custom-title" v-else-if="!hasAvailableIncomeCategories">
                        <span>{{ $t('None') }}</span>
                    </div>
                </template>
                <tree-view-selection-sheet primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           :items="allCategories[allCategoryTypes.Income]"
                                           v-model:show="showCategorySheet"
                                           v-model="transaction.incomeCategory">
                </tree-view-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                key="transferCategorySelection"
                link="#" no-chevron
                :class="{ 'disabled': !hasAvailableTransferCategories, 'readonly': mode === 'view' }"
                :header="$t('Category')"
                @click="showCategorySheet = true"
                v-if="transaction.type === allTransactionTypes.Transfer"
            >
                <template #title>
                    <div class="list-item-custom-title" v-if="hasAvailableTransferCategories">
                        <span>{{ getPrimaryCategoryName(transaction.transferCategory, allCategories[allCategoryTypes.Transfer]) }}</span>
                        <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                        <span>{{ getSecondaryCategoryName(transaction.transferCategory, allCategories[allCategoryTypes.Transfer]) }}</span>
                    </div>
                    <div class="list-item-custom-title" v-else-if="!hasAvailableTransferCategories">
                        <span>{{ $t('None') }}</span>
                    </div>
                </template>
                <tree-view-selection-sheet primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           :items="allCategories[allCategoryTypes.Transfer]"
                                           v-model:show="showCategorySheet"
                                           v-model="transaction.transferCategory">
                </tree-view-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': !allVisibleAccounts.length, 'readonly': mode === 'view' }"
                :header="$t(sourceAccountTitle)"
                :title="sourceAccountName"
                @click="showSourceAccountSheet = true"
            >
                <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                      primary-title-field="name" primary-footer-field="displayBalance"
                                                      primary-icon-field="icon" primary-icon-type="account"
                                                      primary-sub-items-field="accounts"
                                                      :primary-title-i18n="true"
                                                      secondary-key-field="id" secondary-value-field="id"
                                                      secondary-title-field="name" secondary-footer-field="displayBalance"
                                                      secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                      :items="categorizedAccounts"
                                                      v-model:show="showSourceAccountSheet"
                                                      v-model="transaction.sourceAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': !allVisibleAccounts.length, 'readonly': mode === 'view' }"
                :header="$t('Destination Account')"
                :title="destinationAccountName"
                v-if="transaction.type === allTransactionTypes.Transfer"
                @click="showDestinationAccountSheet = true"
            >
                <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                      primary-title-field="name" primary-footer-field="displayBalance"
                                                      primary-icon-field="icon" primary-icon-type="account"
                                                      primary-sub-items-field="accounts"
                                                      :primary-title-i18n="true"
                                                      secondary-key-field="id" secondary-value-field="id"
                                                      secondary-title-field="name" secondary-footer-field="displayBalance"
                                                      secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                      :items="categorizedAccounts"
                                                      v-model:show="showDestinationAccountSheet"
                                                      v-model="transaction.destinationAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'readonly': mode === 'view' }"
                :header="$t('Transaction Time')"
                :title="transactionDisplayTime"
                @click="showTransactionDateTimeSheet = true"
            >
                <date-time-selection-sheet v-model:show="showTransactionDateTimeSheet"
                                           v-model="transaction.time">
                </date-time-selection-sheet>
            </f7-list-item>

            <f7-list-item
                :no-chevron="mode === 'view'"
                class="list-item-with-header-and-title list-item-title-hide-overflow list-item-no-item-after"
                :class="{ 'readonly': mode === 'view' }"
                :header="$t('Transaction Time Zone')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Timezone'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Transaction Time Zone'), popupCloseLinkText: $t('Done') }">
                <select v-model="transaction.timeZone">
                    <option :value="timezone.name"
                            :key="timezone.name"
                            v-for="timezone in allTimezones">{{ `(UTC${timezone.utcOffset}) ${timezone.displayName}` }}</option>
                </select>
                <template #title>
                    <f7-block class="list-item-custom-title no-padding no-margin">
                        <span>{{ `(${transactionDisplayTimezone})` }}</span>
                        <span class="transaction-edit-timezone-name" v-if="transaction.timeZone || transaction.timeZone === ''">{{ transactionDisplayTimezoneName }}</span>
                    </f7-block>
                </template>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                :class="{ 'readonly': mode === 'view' && !transaction.geoLocation }"
                :header="$t('Geographic Location')"
                @click="showGeoLocationActionSheet = true"
            >
                <template #title>
                    <f7-block class="list-item-custom-title no-padding no-margin">
                        <span v-if="transaction.geoLocation">{{ `(${transaction.geoLocation.longitude}, ${transaction.geoLocation.latitude})` }}</span>
                        <span v-else-if="!transaction.geoLocation">{{ geoLocationStatusInfo }}</span>
                    </f7-block>
                </template>

                <map-sheet v-model="transaction.geoLocation"
                           v-model:show="showGeoLocationMapSheet">
                </map-sheet>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                :class="{ 'readonly': mode === 'view' }"
                :header="$t('Tags')"
                @click="showTransactionTagSheet = true"
            >
                <transaction-tag-selection-sheet :items="allTags"
                                                 v-model:show="showTransactionTagSheet"
                                                 v-model="transaction.tagIds">
                </transaction-tag-selection-sheet>

                <template #footer>
                    <f7-block class="margin-top-half no-padding no-margin" v-if="transaction.tagIds && transaction.tagIds.length">
                        <f7-chip media-bg-color="black" class="transaction-edit-tag"
                                 :text="getTagName(tagId)"
                                 :key="tagId"
                                 v-for="tagId in transaction.tagIds">
                            <template #media>
                                <f7-icon f7="number"></f7-icon>
                            </template>
                        </f7-chip>
                    </f7-block>
                    <f7-block class="margin-top-half no-padding no-margin" v-else-if="!transaction.tagIds || !transaction.tagIds.length">
                        <f7-chip class="transaction-edit-tag" :text="$t('None')">
                        </f7-chip>
                    </f7-block>
                </template>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                class="transaction-edit-comment"
                style="height: auto"
                :class="{ 'readonly': mode === 'view' }"
                :label="$t('Description')"
                :placeholder="mode !== 'view' ? $t('Your transaction description (optional)') : ''"
                v-textarea-auto-size
                v-model:value="transaction.comment"
            ></f7-list-input>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showGeoLocationActionSheet" @actions:closed="showGeoLocationActionSheet = false">
            <f7-actions-group>
                <f7-actions-button v-if="mode !== 'view'" @click="updateGeoLocation(true)">{{ $t('Update Geographic Location') }}</f7-actions-button>
                <f7-actions-button v-if="mode !== 'view'" @click="clearGeoLocation">{{ $t('Clear Geographic Location') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="$settings.getMapProvider()">
                <f7-actions-button :class="{ 'disabled': !transaction.geoLocation }" @click="showGeoLocationMapSheet = true">{{ $t('Show on the map') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button v-if="transaction.hideAmount" @click="transaction.hideAmount = false">{{ $t('Show Amount') }}</f7-actions-button>
                <f7-actions-button v-if="!transaction.hideAmount" @click="transaction.hideAmount = true">{{ $t('Hide Amount') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-toolbar tabbar bottom v-if="mode !== 'view'">
            <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" @click="save">
                <span class="tabbar-primary-link">{{ $t(saveButtonTitle) }}</span>
            </f7-link>
        </f7-toolbar>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';
import logger from '@/lib/logger.js';
import {
    isNumber,
    copyObjectTo,
    getNameByKeyValue
} from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getUtcOffsetByUtcOffsetMinutes,
    getDummyUnixTimeForLocalUsage,
    getActualUnixTimeForStore
} from '@/lib/datetime.js';
import {
    stringCurrencyToNumeric
} from '@/lib/currency.js';
import {
    getCategorizedAccounts,
    getAllFilteredAccountsBalance
} from '@/lib/account.js';
import {
    categoryTypeToTransactionType,
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const self = this;
        const query = self.f7route.query;
        const now = getCurrentUnixTime();
        const currentTimezone = self.$locale.getTimezone();

        let defaultType = transactionConstants.allTransactionTypes.Expense;

        if (query.type === transactionConstants.allTransactionTypes.Income.toString()) {
            defaultType = transactionConstants.allTransactionTypes.Income;
        } else if (query.type === transactionConstants.allTransactionTypes.Transfer.toString()) {
            defaultType = transactionConstants.allTransactionTypes.Transfer;
        }

        return {
            mode: 'add',
            editTransactionId: null,
            transaction: {
                type: defaultType,
                time: now,
                timeZone: currentTimezone,
                utcOffset: getTimezoneOffsetMinutes(currentTimezone),
                expenseCategory: '',
                incomeCategory: '',
                transferCategory: '',
                sourceAccountId: '',
                destinationAccountId: '',
                sourceAmount: 0,
                destinationAmount: 0,
                hideAmount: false,
                tagIds: [],
                comment: '',
                geoLocation: null
            },
            loading: true,
            loadingError: null,
            geoLocationStatus: null,
            submitting: false,
            isSupportGeoLocation: !!navigator.geolocation,
            showAccountBalance: self.$settings.isShowAccountBalance(),
            showGeoLocationActionSheet: false,
            showMoreActionSheet: false,
            showSourceAmountSheet: false,
            showDestinationAmountSheet: false,
            showCategorySheet: false,
            showSourceAccountSheet: false,
            showDestinationAccountSheet: false,
            showTransactionDateTimeSheet: false,
            showGeoLocationMapSheet: false,
            showTransactionTagSheet: false
        };
    },
    computed: {
        ...mapStores(useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useTransactionsStore, useExchangeRatesStore),
        title() {
            if (this.mode === 'add') {
                return 'Add Transaction';
            } else if (this.mode === 'edit') {
                return 'Edit Transaction';
            } else {
                return 'Transaction Detail';
            }
        },
        saveButtonTitle() {
            if (this.mode === 'add') {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        sourceAmountName() {
            if (this.transaction.type === this.allTransactionTypes.Expense) {
                return 'Expense Amount';
            } else if (this.transaction.type === this.allTransactionTypes.Income) {
                return 'Income Amount';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Transfer Out Amount';
            } else {
                return '';
            }
        },
        sourceAccountTitle() {
            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                return 'Account';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Source Account';
            } else {
                return '';
            }
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        defaultAccountId() {
            return this.userStore.currentUserDefaultAccountId;
        },
        allTransactionTypes() {
            return transactionConstants.allTransactionTypes;
        },
        allCategoryTypes() {
            return categoryConstants.allCategoryTypes;
        },
        allTimezones() {
            return this.$locale.getAllTimezones(true);
        },
        allAccounts() {
            return this.accountsStore.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.accountsStore.allVisiblePlainAccounts;
        },
        allAccountsMap() {
            return this.accountsStore.allAccountsMap;
        },
        categorizedAccounts() {
            const categorizedAccounts = copyObjectTo(getCategorizedAccounts(this.allVisibleAccounts), {});

            for (let category in categorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(categorizedAccounts, category)) {
                    continue;
                }

                const accountCategory = categorizedAccounts[category];

                if (accountCategory.accounts) {
                    for (let i = 0; i < accountCategory.accounts.length; i++) {
                        const account = accountCategory.accounts[i];

                        if (this.showAccountBalance && account.isAsset) {
                            account.displayBalance = this.$locale.getDisplayCurrency(account.balance, account.currency);
                        } else if (this.showAccountBalance && account.isLiability) {
                            account.displayBalance = this.$locale.getDisplayCurrency(-account.balance, account.currency);
                        } else {
                            account.displayBalance = '***';
                        }
                    }
                }

                if (this.showAccountBalance) {
                    const accountsBalance = getAllFilteredAccountsBalance(categorizedAccounts, account => account.category === accountCategory.category);
                    let totalBalance = 0;
                    let hasUnCalculatedAmount = false;

                    for (let i = 0; i < accountsBalance.length; i++) {
                        if (accountsBalance[i].currency === this.defaultCurrency) {
                            if (accountsBalance[i].isAsset) {
                                totalBalance += accountsBalance[i].balance;
                            } else if (accountsBalance[i].isLiability) {
                                totalBalance -= accountsBalance[i].balance;
                            }
                        } else {
                            const balance = this.exchangeRatesStore.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                            if (!isNumber(balance)) {
                                hasUnCalculatedAmount = true;
                                continue;
                            }

                            if (accountsBalance[i].isAsset) {
                                totalBalance += Math.floor(balance);
                            } else if (accountsBalance[i].isLiability) {
                                totalBalance -= Math.floor(balance);
                            }
                        }
                    }

                    if (hasUnCalculatedAmount) {
                        totalBalance = totalBalance + '+';
                    }

                    accountCategory.displayBalance = this.$locale.getDisplayCurrency(totalBalance, this.defaultCurrency);
                } else {
                    accountCategory.displayBalance = '***';
                }
            }

            return categorizedAccounts;
        },
        allCategories() {
            return this.transactionCategoriesStore.allTransactionCategories;
        },
        allCategoriesMap() {
            return this.transactionCategoriesStore.allTransactionCategoriesMap;
        },
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        hasAvailableExpenseCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Expense] || !this.allCategories[this.allCategoryTypes.Expense].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Expense]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableIncomeCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Income] || !this.allCategories[this.allCategoryTypes.Income].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Income]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableTransferCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Transfer] || !this.allCategories[this.allCategoryTypes.Transfer].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Transfer]);
            return firstAvailableCategoryId !== '';
        },
        sourceAccountName() {
            if (this.transaction.sourceAccountId) {
                return getNameByKeyValue(this.allAccounts, this.transaction.sourceAccountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        },
        destinationAccountName() {
            if (this.transaction.destinationAccountId) {
                return getNameByKeyValue(this.allAccounts, this.transaction.destinationAccountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        },
        transactionDisplayTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, getActualUnixTimeForStore(this.transaction.time, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
        },
        transactionDisplayTimezone() {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(this.transaction.utcOffset)}`;
        },
        transactionDisplayTimezoneName() {
            return getNameByKeyValue(this.allTimezones, this.transaction.timeZone, 'name', 'displayName');
        },
        sourceAmountClass() {
            const classes = {
                'readonly': this.mode === 'view',
                'color-teal': this.transaction.type === this.allTransactionTypes.Expense,
                'color-red': this.transaction.type === this.allTransactionTypes.Income,
            };

            classes[this.getFontClassByAmount(this.transaction.sourceAmount)] = true;

            return classes;
        },
        destinationAmountClass() {
            const classes = {
                'readonly': this.mode === 'view'
            };

            classes[this.getFontClassByAmount(this.transaction.destinationAmount)] = true;

            return classes;
        },
        geoLocationStatusInfo() {
            if (this.geoLocationStatus === 'success') {
                return '';
            } else if (this.geoLocationStatus === 'getting') {
                return this.$t('Getting Location...');
            } else {
                return this.$t('No Location');
            }
        },
        allowedMinAmount() {
            return transactionConstants.minAmount;
        },
        allowedMaxAmount() {
            return transactionConstants.maxAmount;
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            return null;
        }
    },
    watch: {
        'transaction.sourceAmount': function (newValue, oldValue) {
            if (this.mode === 'view') {
                return;
            }

            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                this.transaction.destinationAmount = newValue;
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                const sourceAccount = this.allAccountsMap[this.transaction.sourceAccountId]
                const destinationAccount = this.allAccountsMap[this.transaction.destinationAccountId]

                if (sourceAccount && destinationAccount && sourceAccount.currency !== destinationAccount.currency) {
                    const exchangedOldValue = this.exchangeRatesStore.getExchangedAmount(oldValue, sourceAccount.currency, destinationAccount.currency);
                    const exchangedNewValue = this.exchangeRatesStore.getExchangedAmount(newValue, sourceAccount.currency, destinationAccount.currency);

                    if (isNumber(exchangedOldValue)) {
                        oldValue = Math.floor(exchangedOldValue);
                    }

                    if (isNumber(exchangedNewValue)) {
                        newValue = Math.floor(exchangedNewValue);
                    }
                }

                if ((!sourceAccount || !destinationAccount || this.transaction.destinationAmount === oldValue) &&
                    (stringCurrencyToNumeric(this.allowedMinAmount) <= newValue &&
                        newValue <= stringCurrencyToNumeric(this.allowedMaxAmount))) {
                    this.transaction.destinationAmount = newValue;
                }
            }
        },
        'transaction.destinationAmount': function (newValue) {
            if (this.mode === 'view') {
                return;
            }

            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                this.transaction.sourceAmount = newValue;
            }
        },
        'transaction.timeZone': function (newValue) {
            for (let i = 0; i < this.allTimezones.length; i++) {
                if (this.allTimezones[i].name === newValue) {
                    this.transaction.utcOffset = this.allTimezones[i].utcOffsetMinutes;
                    break;
                }
            }
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        if (self.f7route.path === '/transaction/edit') {
            self.mode = 'edit';
        } else if (self.f7route.path === '/transaction/detail') {
            self.mode = 'view';
        }

        self.loading = true;

        const promises = [
            self.accountsStore.loadAllAccounts({ force: false }),
            self.transactionCategoriesStore.loadAllCategories({ force: false }),
            self.transactionTagsStore.loadAllTags({ force: false })
        ];

        if (query.id) {
            if (self.mode === 'edit') {
                self.editTransactionId = query.id;
            }

            promises.push(self.transactionsStore.getTransaction({ transactionId: query.id }));
        }

        if (query.type && query.type !== '0' &&
            query.type >= self.allTransactionTypes.Income &&
            query.type <= self.allTransactionTypes.Transfer) {
            self.transaction.type = parseInt(query.type);
        }

        Promise.all(promises).then(function (responses) {
            if (query.id && !responses[3]) {
                self.$toast('Unable to get transaction');
                self.loadingError = 'Unable to get transaction';
                return;
            }

            if ((!query.type || query.type === '0') && query.categoryId && query.categoryId !== '0' && self.allCategoriesMap[query.categoryId]) {
                const category = self.allCategoriesMap[query.categoryId];
                const type = categoryTypeToTransactionType(category.type);

                if (isNumber(type)) {
                    self.transaction.type = type;
                }
            }

            if (self.allCategories[self.allCategoryTypes.Expense] &&
                self.allCategories[self.allCategoryTypes.Expense].length) {
                if (query.categoryId && query.categoryId !== '0' && self.isCategoryIdAvailable(self.allCategories[self.allCategoryTypes.Expense], query.categoryId)) {
                    self.transaction.expenseCategory = query.categoryId;
                }

                if (!self.transaction.expenseCategory) {
                    self.transaction.expenseCategory = self.getFirstAvailableCategoryId(self.allCategories[self.allCategoryTypes.Expense]);
                }
            }

            if (self.allCategories[self.allCategoryTypes.Income] &&
                self.allCategories[self.allCategoryTypes.Income].length) {
                if (query.categoryId && query.categoryId !== '0' && self.isCategoryIdAvailable(self.allCategories[self.allCategoryTypes.Income], query.categoryId)) {
                    self.transaction.incomeCategory = query.categoryId;
                }

                if (!self.transaction.incomeCategory) {
                    self.transaction.incomeCategory = self.getFirstAvailableCategoryId(self.allCategories[self.allCategoryTypes.Income]);
                }
            }

            if (self.allCategories[self.allCategoryTypes.Transfer] &&
                self.allCategories[self.allCategoryTypes.Transfer].length) {
                if (query.categoryId && query.categoryId !== '0' && self.isCategoryIdAvailable(self.allCategories[self.allCategoryTypes.Transfer], query.categoryId)) {
                    self.transaction.transferCategory = query.categoryId;
                }

                if (!self.transaction.transferCategory) {
                    self.transaction.transferCategory = self.getFirstAvailableCategoryId(self.allCategories[self.allCategoryTypes.Transfer]);
                }
            }

            if (self.allVisibleAccounts.length) {
                if (query.accountId && query.accountId !== '0') {
                    for (let i = 0; i < self.allVisibleAccounts.length; i++) {
                        if (self.allVisibleAccounts[i].id === query.accountId) {
                            self.transaction.sourceAccountId = query.accountId;
                            self.transaction.destinationAccountId = query.accountId;
                            break;
                        }
                    }
                }

                if (!self.transaction.sourceAccountId) {
                    if (self.defaultAccountId && self.allAccountsMap[self.defaultAccountId]) {
                        self.transaction.sourceAccountId = self.defaultAccountId;
                    } else {
                        self.transaction.sourceAccountId = self.allVisibleAccounts[0].id;
                    }
                }

                if (!self.transaction.destinationAccountId) {
                    if (self.defaultAccountId && self.allAccountsMap[self.defaultAccountId]) {
                        self.transaction.destinationAccountId = self.defaultAccountId;
                    } else {
                        self.transaction.destinationAccountId = self.allVisibleAccounts[0].id;
                    }
                }
            }

            if (query.id) {
                const transaction = responses[3];

                if (self.mode === 'edit') {
                    self.transaction.id = transaction.id;
                }

                self.transaction.type = transaction.type;

                if (self.transaction.type === self.allTransactionTypes.Expense) {
                    self.transaction.expenseCategory = transaction.categoryId;
                } else if (self.transaction.type === self.allTransactionTypes.Income) {
                    self.transaction.incomeCategory = transaction.categoryId;
                } else if (self.transaction.type === self.allTransactionTypes.Transfer) {
                    self.transaction.transferCategory = transaction.categoryId;
                }

                if (self.mode === 'edit' || self.mode === 'view') {
                    self.transaction.utcOffset = transaction.utcOffset;
                    self.transaction.timeZone = null;
                    self.transaction.time = getDummyUnixTimeForLocalUsage(transaction.time, self.transaction.utcOffset, getBrowserTimezoneOffsetMinutes());
                }

                self.transaction.sourceAccountId = transaction.sourceAccountId;

                if (transaction.destinationAccountId) {
                    self.transaction.destinationAccountId = transaction.destinationAccountId;
                }

                self.transaction.sourceAmount = transaction.sourceAmount;

                if (transaction.destinationAmount) {
                    self.transaction.destinationAmount = transaction.destinationAmount;
                }

                self.transaction.hideAmount = transaction.hideAmount;
                self.transaction.tagIds = transaction.tagIds || [];
                self.transaction.comment = transaction.comment;

                if (self.mode === 'edit' || self.mode === 'view') {
                    self.transaction.geoLocation = transaction.geoLocation;
                }
            }

            self.loading = false;
        }).catch(error => {
            logger.error('failed to load essential data for editing transaction', error);

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
            this.$routeBackOnError(this.f7router, 'loadingError');

            if (this.$settings.isAutoGetCurrentGeoLocation() && this.mode === 'add'
                && !this.geoLocationStatus && !this.transaction.geoLocation) {
                this.updateGeoLocation(false);
            }
        },
        save() {
            const self = this;
            const router = self.f7router;

            if (self.mode === 'view') {
                return;
            }

            const submitTransaction = {
                type: self.transaction.type,
                time: getActualUnixTimeForStore(self.transaction.time, self.transaction.utcOffset, getBrowserTimezoneOffsetMinutes()),
                sourceAccountId: self.transaction.sourceAccountId,
                sourceAmount: self.transaction.sourceAmount,
                destinationAccountId: '0',
                destinationAmount: 0,
                hideAmount: self.transaction.hideAmount,
                tagIds: self.transaction.tagIds,
                comment: self.transaction.comment,
                geoLocation: self.transaction.geoLocation,
                utcOffset: self.transaction.utcOffset
            };

            if (self.transaction.type === self.allTransactionTypes.Expense) {
                submitTransaction.categoryId = self.transaction.expenseCategory;
            } else if (self.transaction.type === self.allTransactionTypes.Income) {
                submitTransaction.categoryId = self.transaction.incomeCategory;
            } else if (self.transaction.type === self.allTransactionTypes.Transfer) {
                submitTransaction.categoryId = self.transaction.transferCategory;
                submitTransaction.destinationAccountId = self.transaction.destinationAccountId;
                submitTransaction.destinationAmount = self.transaction.destinationAmount;
            } else {
                self.$toast('An error has occurred');
                return;
            }

            if (self.mode === 'edit') {
                submitTransaction.id = self.transaction.id;
            }

            const doSubmit = function () {
                self.submitting = true;
                self.$showLoading(() => self.submitting);

                self.transactionsStore.saveTransaction({
                    transaction: submitTransaction,
                    defaultCurrency: self.defaultCurrency
                }).then(() => {
                    self.submitting = false;
                    self.$hideLoading();

                    if (self.mode === 'add') {
                        self.$toast('You have added a new transaction');
                    } else if (self.mode === 'edit') {
                        self.$toast('You have saved this transaction');
                    }

                    router.back();
                }).catch(error => {
                    self.submitting = false;
                    self.$hideLoading();

                    if (!error.processed) {
                        self.$toast(error.message || error);
                    }
                });
            };

            if (submitTransaction.sourceAmount === 0) {
                self.$confirm('Are you sure you want to save this transaction whose amount is 0?', () => {
                    doSubmit();
                });
            } else {
                doSubmit();
            }
        },
        updateGeoLocation(forceUpdate) {
            const self = this;

            if (!self.isSupportGeoLocation) {
                logger.warn('this browser does not support geo location');

                if (forceUpdate) {
                    self.$toast('Unable to get current position');
                }
                return;
            }

            navigator.geolocation.getCurrentPosition(function (position) {
                if (!position || !position.coords) {
                    logger.error('current position is null');
                    self.geoLocationStatus = 'error';

                    if (forceUpdate) {
                        self.$toast('Unable to get current position');
                    }

                    return;
                }

                self.geoLocationStatus = 'success';

                self.transaction.geoLocation = {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude
                };
            }, function (err) {
                logger.error('cannot get current position', err);
                self.geoLocationStatus = 'error';

                if (forceUpdate) {
                    self.$toast('Unable to get current position');
                }
            });

            self.geoLocationStatus = 'getting';
        },
        clearGeoLocation() {
            this.geoLocationStatus = null;
            this.transaction.geoLocation = null;
        },
        isCategoryIdAvailable(categories, categoryId) {
            if (!categories || !categories.length) {
                return false;
            }

            for (let i = 0; i < categories.length; i++) {
                for (let j = 0; j < categories[i].subCategories.length; j++) {
                    if (categories[i].subCategories[j].id === categoryId) {
                        return true;
                    }
                }
            }

            return false;
        },
        getFirstAvailableCategoryId(categories) {
            if (!categories || !categories.length) {
                return '';
            }

            for (let i = 0; i < categories.length; i++) {
                for (let j = 0; j < categories[i].subCategories.length; j++) {
                    return categories[i].subCategories[j].id;
                }
            }
        },
        getFontClassByAmount(amount) {
            if (amount >= 100000000 || amount <= -100000000) {
                return 'ebk-small-amount';
            } else if (amount >= 1000000 || amount <= -1000000) {
                return 'ebk-normal-amount';
            } else {
                return 'ebk-large-amount';
            }
        },
        getDisplayAmount(amount, hideAmount) {
            if (hideAmount) {
                return this.$locale.getDisplayCurrency('***');
            }

            return this.$locale.getDisplayCurrency(amount);
        },
        getPrimaryCategoryName(categoryId, allCategories) {
            return getTransactionPrimaryCategoryName(categoryId, allCategories);
        },
        getSecondaryCategoryName(categoryId, allCategories) {
            return getTransactionSecondaryCategoryName(categoryId, allCategories);
        },
        getTagName(tagId) {
            return getNameByKeyValue(this.allTags, tagId, 'id', 'name');
        }
    }
};
</script>

<style>
.category-separate-icon.icon {
    margin-left: 5px;
    margin-right: 5px;
    font-size: var(--ebk-category-separate-icon-font-size);
    line-height: 16px;
    color: var(--f7-color-gray-tint);
}

.transaction-edit-amount {
    line-height: 53px;
    color: var(--f7-theme-color);
}

.transaction-edit-amount .item-title {
    font-weight: bolder;
}

.transaction-edit-amount .item-header {
    padding-top: calc(var(--f7-typography-padding) / 2);
}

.transaction-edit-timezone-name {
    padding-left: 4px;
}

.transaction-edit-tag {
    margin-right: 4px;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
