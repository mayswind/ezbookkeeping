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
                    <f7-button :text="$t('Expense')" :active="transaction.type === allTransactionTypes.Expense"
                               v-if="transaction.type !== allTransactionTypes.ModifyBalance"
                               @click="transaction.type = allTransactionTypes.Expense"></f7-button>
                    <f7-button :text="$t('Income')" :active="transaction.type === allTransactionTypes.Income"
                               v-if="transaction.type !== allTransactionTypes.ModifyBalance"
                               @click="transaction.type = allTransactionTypes.Income"></f7-button>
                    <f7-button :text="$t('Transfer')" :active="transaction.type === allTransactionTypes.Transfer"
                               v-if="transaction.type !== allTransactionTypes.ModifyBalance"
                               @click="transaction.type = allTransactionTypes.Transfer"></f7-button>
                    <f7-button :text="$t('Modify Balance')" :active="transaction.type === allTransactionTypes.ModifyBalance"
                               v-if="transaction.type === allTransactionTypes.ModifyBalance"></f7-button>
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
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow list-item-no-item-after" header="Transaction Timezone" title="(UTC XX:XX) System Default" link="#" :no-chevron="mode === 'view'"></f7-list-item>
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
                class="transaction-edit-amount text-color-primary"
                link="#" no-chevron
                :class="destinationAmountClass"
                :header="transferInAmountTitle"
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
                :class="{ 'readonly': mode === 'view' && transaction.utcOffset === currentTimezoneOffsetMinutes }"
                :header="$t('Transaction Time')"
                :title="transactionDisplayTime"
                @click="mode !== 'view' ? showTransactionDateTimeSheet = true : showTimeInDefaultTimezone = !showTimeInDefaultTimezone"
            >
                <date-time-selection-sheet v-model:show="showTransactionDateTimeSheet"
                                           v-model="transaction.time">
                </date-time-selection-sheet>
            </f7-list-item>

            <f7-list-item
                :no-chevron="mode === 'view'"
                class="list-item-with-header-and-title list-item-title-hide-overflow list-item-no-item-after"
                :class="{ 'readonly': mode === 'view' }"
                :header="$t('Transaction Timezone')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Timezone'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Transaction Timezone'), popupCloseLinkText: $t('Done') }">
                <select v-model="transaction.timeZone">
                    <option :value="timezone.name" :key="timezone.name"
                            v-for="timezone in allTimezones">{{ timezone.displayNameWithUtcOffset }}</option>
                </select>
                <template #title>
                    <f7-block class="list-item-custom-title no-padding no-margin">
                        <span>{{ `(${transactionDisplayTimezone})` }}</span>
                        <span class="transaction-edit-timezone-name" v-if="transaction.timeZone || transaction.timeZone === ''">{{ transactionDisplayTimezoneName }}</span>
                        <span class="transaction-edit-timezone-name" v-else-if="!transaction.timeZone && transaction.timeZone !== ''">{{ transactionTimezoneTimeDifference }}</span>
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
            <f7-actions-group v-if="mapProvider">
                <f7-actions-button :class="{ 'disabled': !transaction.geoLocation }" @click="showGeoLocationMapSheet = true">{{ $t('Show on the map') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group v-if="transaction.type === allTransactionTypes.Transfer">
                <f7-actions-button @click="swapTransactionData(true, false)">{{ $t('Swap Account') }}</f7-actions-button>
                <f7-actions-button @click="swapTransactionData(false, true)">{{ $t('Swap Amount') }}</f7-actions-button>
                <f7-actions-button @click="swapTransactionData(true, true)">{{ $t('Swap Account and Amount') }}</f7-actions-button>
            </f7-actions-group>
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
import { useSettingsStore } from '@/stores/setting.js';
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
    getNameByKeyValue
} from '@/lib/common.js';
import {
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getUtcOffsetByUtcOffsetMinutes,
    getActualUnixTimeForStore
} from '@/lib/datetime.js';
import {
    getAdaptiveDisplayAmountRate
} from '@/lib/currency.js';
import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName,
    getFirstAvailableCategoryId
} from '@/lib/category.js';
import { getMapProvider } from '@/lib/server_settings.js';
import { setTransactionModelByTransaction } from '@/lib/transaction.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const query = this.f7route.query;
        const transactionsStore = useTransactionsStore();
        const newTransaction = transactionsStore.generateNewTransactionModel(query.type);

        return {
            mode: 'add',
            editTransactionId: null,
            transaction: newTransaction,
            loading: true,
            loadingError: null,
            geoLocationStatus: null,
            submitting: false,
            isSupportGeoLocation: !!navigator.geolocation,
            showTimeInDefaultTimezone: false,
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
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useTransactionsStore, useExchangeRatesStore),
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
                return 'Amount';
            }
        },
        sourceAccountTitle() {
            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                return 'Account';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Source Account';
            } else {
                return 'Account';
            }
        },
        transferInAmountTitle() {
            const sourceAccount = this.allAccountsMap[this.transaction.sourceAccountId];
            const destinationAccount = this.allAccountsMap[this.transaction.destinationAccountId];

            if (!sourceAccount || !destinationAccount || sourceAccount.currency === destinationAccount.currency) {
                return this.$t('Transfer In Amount');
            }

            const fromExchangeRate = this.exchangeRatesStore.latestExchangeRateMap[sourceAccount.currency];
            const toExchangeRate = this.exchangeRatesStore.latestExchangeRateMap[destinationAccount.currency];
            const amountRate = getAdaptiveDisplayAmountRate(this.transaction.sourceAmount, this.transaction.destinationAmount, fromExchangeRate, toExchangeRate, this.settingsStore.appSettings.thousandsSeparator);

            if (!amountRate) {
                return this.$t('Transfer In Amount');
            }

            return this.$t('Transfer In Amount') + ` (${amountRate})`;
        },
        currentTimezoneOffsetMinutes() {
            return getTimezoneOffsetMinutes(this.settingsStore.appSettings.timeZone);
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
            return this.$locale.getCategorizedAccountsWithDisplayBalance(this.exchangeRatesStore, this.allVisibleAccounts, this.showAccountBalance, this.defaultCurrency, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
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
            if (this.mode !== 'view' || !this.showTimeInDefaultTimezone) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, getActualUnixTimeForStore(this.transaction.time, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
            }

            return `${this.$locale.formatUnixTimeToLongDateTime(this.userStore, getActualUnixTimeForStore(this.transaction.time, this.transaction.utcOffset, getBrowserTimezoneOffsetMinutes()))} (UTC${getTimezoneOffset(this.settingsStore.appSettings.timeZone)})`;
        },
        transactionDisplayTimezone() {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(this.transaction.utcOffset)}`;
        },
        transactionDisplayTimezoneName() {
            return getNameByKeyValue(this.allTimezones, this.transaction.timeZone, 'name', 'displayName');
        },
        transactionTimezoneTimeDifference() {
            return this.$locale.getTimezoneDifferenceDisplayText(this.transaction.utcOffset);
        },
        sourceAmountClass() {
            const classes = {
                'readonly': this.mode === 'view',
                'text-color-teal': this.transaction.type === this.allTransactionTypes.Expense,
                'text-color-red': this.transaction.type === this.allTransactionTypes.Income,
                'text-color-primary': this.transaction.type === this.allTransactionTypes.Transfer
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
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
        },
        mapProvider() {
            return getMapProvider();
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

            this.transactionsStore.setTransactionSuitableDestinationAmount(this.transaction, oldValue, newValue);
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
                self.$toast('Unable to retrieve transaction');
                self.loadingError = 'Unable to retrieve transaction';
                return;
            }

            setTransactionModelByTransaction(
                self.transaction,
                query.id ? responses[3] : null,
                self.allCategories,
                self.allCategoriesMap,
                self.allVisibleAccounts,
                self.allAccountsMap,
                self.defaultAccountId,
                {
                    type: query.type,
                    categoryId: query.categoryId,
                    accountId: query.accountId
                },
                (self.mode === 'edit' || self.mode === 'view'),
                (self.mode === 'edit' || self.mode === 'view')
            );

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

            if (this.settingsStore.appSettings.autoGetCurrentGeoLocation && this.mode === 'add'
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

            const doSubmit = function () {
                self.submitting = true;
                self.$showLoading(() => self.submitting);

                self.transactionsStore.saveTransaction({
                    transaction: self.transaction,
                    defaultCurrency: self.defaultCurrency,
                    isEdit: self.mode === 'edit'
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

            if (self.transaction.sourceAmount === 0) {
                self.$confirm('Are you sure you want to save this transaction with a zero amount?', () => {
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
                    self.$toast('Unable to retrieve current position');
                }
                return;
            }

            navigator.geolocation.getCurrentPosition(function (position) {
                if (!position || !position.coords) {
                    logger.error('current position is null');
                    self.geoLocationStatus = 'error';

                    if (forceUpdate) {
                        self.$toast('Unable to retrieve current position');
                    }

                    return;
                }

                self.geoLocationStatus = 'success';

                self.transaction.geoLocation = {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude
                };
            }, function (err) {
                logger.error('cannot retrieve current position', err);
                self.geoLocationStatus = 'error';

                if (forceUpdate) {
                    self.$toast('Unable to retrieve current position');
                }
            });

            self.geoLocationStatus = 'getting';
        },
        clearGeoLocation() {
            this.geoLocationStatus = null;
            this.transaction.geoLocation = null;
        },
        swapTransactionData(swapAccount, swapAmount) {
            if (swapAccount) {
                const oldSourceAccountId = this.transaction.sourceAccountId;
                this.transaction.sourceAccountId = this.transaction.destinationAccountId;
                this.transaction.destinationAccountId = oldSourceAccountId;
            }

            if (swapAmount) {
                const oldSourceAmount = this.transaction.sourceAmount;
                this.transaction.sourceAmount = this.transaction.destinationAmount;
                this.transaction.destinationAmount = oldSourceAmount;
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
                return this.getDisplayCurrency('***');
            }

            return this.getDisplayCurrency(amount);
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
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
