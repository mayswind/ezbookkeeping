<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': editAccountId || account.type !== allAccountTypes.MultiSubAccounts.type }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :class="{ 'disabled': isInputEmpty() || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title" header="Account Category" title="Category"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Account Type" title="Account Type"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-else-if="!loading">
            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title"
                :header="$t('Account Category')"
                :title="getAccountCategoryName(account.category)"
                @click="showAccountCategorySheet = true"
            >
                <list-item-selection-sheet value-type="item"
                                           key-field="type" value-field="type" title-field="displayName"
                                           icon-field="defaultAccountIconId" icon-type="account"
                                           :items="allAccountCategories"
                                           v-model:show="showAccountCategorySheet"
                                           v-model="account.category">
                </list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title"
                :class="{ 'disabled': editAccountId }"
                :header="$t('Account Type')"
                :title="getAccountTypeName(account.type)"
                @click="showAccountTypeSheet = true"
            >
                <list-item-selection-sheet value-type="item"
                                           key-field="type" value-field="type" title-field="displayName"
                                           :items="allAccountTypesArray"
                                           v-model:show="showAccountTypeSheet"
                                           v-model="account.type">
                </list-item-selection-sheet>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-input label="Account Name" placeholder="Your account name"></f7-list-input>
            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>Account Icon</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <f7-icon f7="app_fill"></f7-icon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>Account Color</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <f7-icon f7="app_fill"></f7-icon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>
                        </div>
                    </div>
                </template>
            </f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Currency" title="Currency" :link="editAccountId ? null : '#'"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Account Balance" title="Balance"></f7-list-item>
            <f7-list-item class="list-item-toggle" header="Visible" after="True"></f7-list-item>
            <f7-list-input label="Description" type="textarea" placeholder="Your account description (optional)"></f7-list-input>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-else-if="!loading && account.type === allAccountTypes.SingleAccount.type">
            <f7-list-input
                type="text"
                clear-button
                :label="$t('Account Name')"
                :placeholder="$t('Your account name')"
                v-model:value="account.name"
            ></f7-list-input>

            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="account.showIconSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ $t('Account Icon') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <icon-selection-sheet :all-icon-infos="allAccountIcons"
                                                  :color="account.color"
                                                  v-model:show="account.showIconSelectionSheet"
                                                  v-model="account.icon"
                            ></icon-selection-sheet>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="account.showColorSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ $t('Account Color') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <color-selection-sheet :all-color-infos="allAccountColors"
                                                   v-model:show="account.showColorSelectionSheet"
                                                   v-model="account.color"
                            ></color-selection-sheet>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :class="{ 'disabled': editAccountId }"
                :header="$t('Currency')"
                :no-chevron="!!editAccountId"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Currency Name'), popupCloseLinkText: $t('Done') }"
            >
                <template #title>
                    <div class="no-padding no-margin">
                        <span>{{ getCurrencyName(account.currency) }}&nbsp;</span>
                        <small class="smaller">{{ account.currency }}</small>
                    </div>
                </template>
                <select autocomplete="transaction-currency" v-model="account.currency">
                    <option :value="currency.currencyCode"
                            :key="currency.currencyCode"
                            v-for="currency in allCurrencies">{{ currency.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Statement Date')"
                :title="getAccountCreditCardStatementDate(account.creditCardStatementDate)"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Statement Date'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Statement Date'), popupCloseLinkText: $t('Done') }"
                v-if="isAccountSupportCreditCardStatementDate()"
            >
                <select v-model="account.creditCardStatementDate">
                    <option :value="monthDay.day"
                            :key="monthDay.day"
                            v-for="monthDay in allAvailableMonthDays">{{ monthDay.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title"
                :class="{ 'disabled': editAccountId }"
                :header="$t('Account Balance')"
                :title="getAccountBalance(account)"
                @click="account.showBalanceSheet = true"
            >
                <number-pad-sheet :min-value="allowedMinAmount"
                                  :max-value="allowedMaxAmount"
                                  :currency="account.currency"
                                  v-model:show="account.showBalanceSheet"
                                  v-model="account.balance"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="account-edit-balancetime list-item-with-header-and-title"
                link="#" no-chevron
                v-show="account.balance"
                v-if="!editAccountId"
            >
                <template #header>
                    <div class="account-edit-balancetime-header" @click="showDateTimeDialog(account, 'time')">{{ $t('Balance Time') }}</div>
                </template>
                <template #title>
                    <div class="account-edit-balancetime-title">
                        <div @click="showDateTimeDialog(account, 'date')">{{ getAccountBalanceDate(account.balanceTime) }}</div>&nbsp;<div class="account-edit-balancetime-time" @click="showDateTimeDialog(account, 'time')">{{ getAccountBalanceTime(account.balanceTime) }}</div>
                    </div>
                </template>
                <date-time-selection-sheet :init-mode="account.balanceDateTimeSheetMode"
                                           v-model:show="account.showBalanceDateTimeSheet"
                                           v-model="account.balanceTime">
                </date-time-selection-sheet>
            </f7-list-item>

            <f7-list-item :title="$t('Visible')" v-if="editAccountId">
                <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                style="height: auto"
                :label="$t('Description')"
                :placeholder="$t('Your account description (optional)')"
                v-textarea-auto-size
                v-model:value="account.comment"
            ></f7-list-input>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-else-if="!loading && account.type === allAccountTypes.MultiSubAccounts.type">
            <f7-list-input
                type="text"
                clear-button
                :label="$t('Account Name')"
                :placeholder="$t('Your account name')"
                v-model:value="account.name"
            ></f7-list-input>

            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="account.showIconSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ $t('Account Icon') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <icon-selection-sheet :all-icon-infos="allAccountIcons"
                                                  :color="account.color"
                                                  v-model:show="account.showIconSelectionSheet"
                                                  v-model="account.icon"
                            ></icon-selection-sheet>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="account.showColorSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ $t('Account Color') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <color-selection-sheet :all-color-infos="allAccountColors"
                                                   v-model:show="account.showColorSelectionSheet"
                                                   v-model="account.color"
                            ></color-selection-sheet>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Statement Date')"
                :title="getAccountCreditCardStatementDate(account.creditCardStatementDate)"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Statement Date'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Statement Date'), popupCloseLinkText: $t('Done') }"
                v-if="isAccountSupportCreditCardStatementDate()"
            >
                <select v-model="account.creditCardStatementDate">
                    <option :value="monthDay.day"
                            :key="monthDay.day"
                            v-for="monthDay in allAvailableMonthDays">{{ monthDay.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item :title="$t('Visible')" v-if="editAccountId">
                <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                style="height: auto"
                :label="$t('Description')"
                :placeholder="$t('Your account description (optional)')"
                v-textarea-auto-size
                v-model:value="account.comment"
            ></f7-list-input>
        </f7-list>

        <f7-block class="no-padding no-margin" v-if="!loading && account.type === allAccountTypes.MultiSubAccounts.type">
            <f7-list strong inset dividers class="subaccount-edit-list margin-vertical"
                     :key="idx"
                     v-for="(subAccount, idx) in subAccounts">
                <f7-list-item group-title>
                    <small>{{ $t('Sub Account') + ' #' + (idx + 1) }}</small>
                    <f7-button rasied fill class="subaccount-delete-button" color="red" icon-f7="trash" icon-size="16px"
                               :tooltip="$t('Remove Sub-account')"
                               v-if="!editAccountId"
                               @click="removeSubAccount(subAccount, false)">
                    </f7-button>
                </f7-list-item>

                <f7-list-input
                    type="text"
                    clear-button
                    :label="$t('Sub-account Name')"
                    :placeholder="$t('Your sub-account name')"
                    v-model:value="subAccount.name"
                ></f7-list-input>

                <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                    <template #default>
                        <div class="grid grid-cols-2">
                            <div class="list-item-subitem no-chevron">
                                <a class="item-link" href="#" @click="subAccount.showIconSelectionSheet = true">
                                    <div class="item-content">
                                        <div class="item-inner">
                                            <div class="item-header">
                                                <span>{{ $t('Sub-account Icon') }}</span>
                                            </div>
                                            <div class="item-title">
                                                <div class="list-item-custom-title no-padding">
                                                    <ItemIcon icon-type="account" :icon-id="subAccount.icon" :color="subAccount.color"></ItemIcon>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </a>

                                <icon-selection-sheet :all-icon-infos="allAccountIcons"
                                                      :color="subAccount.color"
                                                      v-model:show="subAccount.showIconSelectionSheet"
                                                      v-model="subAccount.icon"
                                ></icon-selection-sheet>
                            </div>
                            <div class="list-item-subitem no-chevron">
                                <a class="item-link" href="#" @click="subAccount.showColorSelectionSheet = true">
                                    <div class="item-content">
                                        <div class="item-inner">
                                            <div class="item-header">
                                                <span>{{ $t('Sub-account Color') }}</span>
                                            </div>
                                            <div class="item-title">
                                                <div class="list-item-custom-title no-padding">
                                                    <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="subAccount.color"></ItemIcon>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </a>

                                <color-selection-sheet :all-color-infos="allAccountColors"
                                                       v-model:show="subAccount.showColorSelectionSheet"
                                                       v-model="subAccount.color"
                                ></color-selection-sheet>
                            </div>
                        </div>
                    </template>
                </f7-list-item>

                <f7-list-item
                    class="list-item-with-header-and-title list-item-no-item-after"
                    :class="{ 'disabled': editAccountId }"
                    :header="$t('Currency')"
                    :no-chevron="!!editAccountId"
                    smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Currency Name'), popupCloseLinkText: $t('Done') }"
                >
                    <template #title>
                        <div class="no-padding no-margin">
                            <span>{{ getCurrencyName(subAccount.currency) }}&nbsp;</span>
                            <small class="smaller">{{ subAccount.currency }}</small>
                        </div>
                    </template>
                    <select autocomplete="transaction-currency" v-model="subAccount.currency">
                        <option :value="currency.currencyCode"
                                :key="currency.currencyCode"
                                v-for="currency in allCurrencies">{{ currency.displayName }}</option>
                    </select>
                </f7-list-item>

                <f7-list-item
                    link="#" no-chevron
                    class="list-item-with-header-and-title"
                    :class="{ 'disabled': editAccountId }"
                    :header="$t('Sub-account Balance')"
                    :title="getAccountBalance(subAccount)"
                    @click="subAccount.showBalanceSheet = true"
                >
                    <number-pad-sheet :min-value="allowedMinAmount"
                                      :max-value="allowedMaxAmount"
                                      :currency="subAccount.currency"
                                      v-model:show="subAccount.showBalanceSheet"
                                      v-model="subAccount.balance"
                    ></number-pad-sheet>
                </f7-list-item>

                <f7-list-item
                    class="account-edit-balancetime list-item-with-header-and-title"
                    link="#" no-chevron
                    v-show="subAccount.balance"
                    v-if="!editAccountId"
                >
                    <template #header>
                        <div class="account-edit-balancetime-header" @click="showDateTimeDialog(subAccount, 'time')">{{ $t('Sub-account Balance Time') }}</div>
                    </template>
                    <template #title>
                        <div class="account-edit-balancetime-title">
                            <div @click="showDateTimeDialog(subAccount, 'date')">{{ getAccountBalanceDate(subAccount.balanceTime) }}</div>&nbsp;<div class="account-edit-balancetime-time" @click="showDateTimeDialog(subAccount, 'time')">{{ getAccountBalanceTime(subAccount.balanceTime) }}</div>
                        </div>
                    </template>
                    <date-time-selection-sheet :init-mode="subAccount.balanceDateTimeSheetMode"
                                               v-model:show="subAccount.showBalanceDateTimeSheet"
                                               v-model="subAccount.balanceTime">
                    </date-time-selection-sheet>
                </f7-list-item>

                <f7-list-item :title="$t('Visible')" v-if="editAccountId">
                    <f7-toggle :checked="subAccount.visible" @toggle:change="subAccount.visible = $event"></f7-toggle>
                </f7-list-item>

                <f7-list-input
                    type="textarea"
                    style="height: auto"
                    :label="$t('Description')"
                    :placeholder="$t('Your sub-account description (optional)')"
                    v-textarea-auto-size
                    v-model:value="subAccount.comment"
                ></f7-list-input>
            </f7-list>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="addSubAccount">{{ $t('Add Sub-account') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to remove this sub-account?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="removeSubAccount(subAccountToDelete, true)">{{ $t('Remove') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';

import { AccountType, AccountCategory } from '@/core/account.ts';
import { ALL_ACCOUNT_ICONS } from '@/consts/icon.ts';
import { ALL_ACCOUNT_COLORS } from '@/consts/color.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { getNameByKeyValue } from '@/lib/common.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import {
    setAccountModelByAnotherAccount,
    setAccountSuitableIcon
} from '@/lib/account.js';
import {
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore
} from '@/lib/datetime.ts';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const accountsStore = useAccountsStore();
        const newAccount = accountsStore.generateNewAccountModel();
        newAccount.showIconSelectionSheet = false;
        newAccount.showColorSelectionSheet = false;
        newAccount.showBalanceSheet = false;
        newAccount.showBalanceDateTimeSheet = false;
        newAccount.balanceDateTimeSheetMode = 'time';

        return {
            editAccountId: null,
            clientSessionId: '',
            loading: false,
            loadingError: null,
            account: newAccount,
            subAccounts: [],
            subAccountToDelete: null,
            submitting: false,
            showAccountCategorySheet: false,
            showAccountTypeSheet: false,
            showMoreActionSheet: false,
            showDeleteActionSheet: false
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore),
        title() {
            if (!this.editAccountId) {
                return 'Add Account';
            } else {
                return 'Edit Account';
            }
        },
        saveButtonTitle() {
            if (!this.editAccountId) {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        allAccountTypes() {
            return AccountType.all();
        },
        allAccountCategories() {
            return this.$locale.getAllAccountCategories();
        },
        allAccountTypesArray() {
            return this.$locale.getAllAccountTypes();
        },
        allAccountIcons() {
            return ALL_ACCOUNT_ICONS;
        },
        allAccountColors() {
            return ALL_ACCOUNT_COLORS;
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allAvailableMonthDays() {
            const allAvailableDays = [];

            allAvailableDays.push({
                day: 0,
                displayName: this.$t('Not set'),
            });

            for (let i = 1; i <= 28; i++) {
                allAvailableDays.push({
                    day: i,
                    displayName: this.$locale.getMonthdayShortName(i),
                });
            }

            return allAvailableDays;
        },
        allowedMinAmount() {
            return TRANSACTION_MIN_AMOUNT;
        },
        allowedMaxAmount() {
            return TRANSACTION_MAX_AMOUNT;
        }
    },
    watch: {
        'account.category': function (newValue, oldValue) {
            this.chooseSuitableIcon(oldValue, newValue);
        },
        'account.type': function () {
            if (this.subAccounts.length < 1) {
                this.addSubAccount();
            }
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        if (query.id) {
            self.loading = true;

            self.editAccountId = query.id;

            self.accountsStore.getAccount({
                accountId: self.editAccountId
            }).then(account => {
                setAccountModelByAnotherAccount(self.account, account);
                self.subAccounts = [];

                if (account.subAccounts && account.subAccounts.length > 0) {
                    for (let i = 0; i < account.subAccounts.length; i++) {
                        const subAccount = self.accountsStore.generateNewSubAccountModel(self.account);
                        setAccountModelByAnotherAccount(subAccount, account.subAccounts[i]);
                        subAccount.showIconSelectionSheet = false;
                        subAccount.showColorSelectionSheet = false;
                        subAccount.showBalanceSheet = false;
                        subAccount.showBalanceDateTimeSheet = false;
                        subAccount.balanceDateTimeSheetMode = 'time';

                        self.subAccounts.push(subAccount);
                    }
                }

                self.loading = false;
            }).catch(error => {
                if (error.processed) {
                    self.loading = false;
                } else {
                    self.loadingError = error;
                    self.$toast(error.message || error);
                }
            });
        } else {
            self.clientSessionId = generateRandomUUID();
            self.loading = false;
        }
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        addSubAccount() {
            if (this.account.type !== this.allAccountTypes.MultiSubAccounts.type) {
                return;
            }

            const subAccount = this.accountsStore.generateNewSubAccountModel(this.account);
            subAccount.showIconSelectionSheet = false;
            subAccount.showColorSelectionSheet = false;
            subAccount.showBalanceSheet = false;
            subAccount.showBalanceDateTimeSheet = false;
            subAccount.balanceDateTimeSheetMode = 'time';

            this.subAccounts.push(subAccount);
        },
        removeSubAccount(subAccount, confirm) {
            if (!subAccount) {
                this.$alert('An error occurred');
                return;
            }

            if (!confirm) {
                this.subAccountToDelete = subAccount;
                this.showDeleteActionSheet = true;
                return;
            }

            this.showDeleteActionSheet = false;
            this.subAccountToDelete = null;

            for (let i = 0; i < this.subAccounts.length; i++) {
                if (this.subAccounts[i] === subAccount) {
                    this.subAccounts.splice(i, 1);
                }
            }
        },
        save() {
            const self = this;
            const router = self.f7router;

            let problemMessage = self.getInputEmptyProblemMessage(self.account, false);

            if (!problemMessage && self.account.type === self.allAccountTypes.MultiSubAccounts.type) {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    problemMessage = self.getInputEmptyProblemMessage(self.subAccounts[i], true);

                    if (problemMessage) {
                        break;
                    }
                }
            }

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            self.accountsStore.saveAccount({
                account: self.account,
                subAccounts: self.subAccounts,
                isEdit: !!self.editAccountId,
                clientSessionId: self.clientSessionId
            }).then(() => {
                self.submitting = false;
                self.$hideLoading();

                if (!self.editAccountId) {
                    self.$toast('You have added a new account');
                } else {
                    self.$toast('You have saved this account');
                }

                router.back();
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        showDateTimeDialog(account, sheetMode) {
            account.balanceDateTimeSheetMode = sheetMode;
            account.showBalanceDateTimeSheet = true;
        },
        getCurrencyName(currencyCode) {
            return this.$locale.getCurrencyName(currencyCode);
        },
        getAccountTypeName(accountType) {
            return getNameByKeyValue(this.allAccountTypesArray, accountType, 'type', 'displayName');
        },
        getAccountCategoryName(accountCategory) {
            return getNameByKeyValue(this.allAccountCategories, accountCategory, 'type', 'displayName');
        },
        getAccountCreditCardStatementDate(statementDate) {
            return getNameByKeyValue(this.allAvailableMonthDays, statementDate, 'day', 'displayName');
        },
        getAccountBalance(account) {
            return this.getDisplayCurrency(account.balance, account.currency);
        },
        getAccountBalanceDate(balanceTime) {
            return this.$locale.formatUnixTimeToLongDate(this.userStore, getActualUnixTimeForStore(balanceTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
        },
        getAccountBalanceTime(balanceTime) {
            return this.$locale.formatUnixTimeToLongTime(this.userStore, getActualUnixTimeForStore(balanceTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        isAccountSupportCreditCardStatementDate() {
            return this.account && this.account.category === AccountCategory.CreditCard.type;
        },
        chooseSuitableIcon(oldCategory, newCategory) {
            setAccountSuitableIcon(this.account, oldCategory, newCategory);
        },
        isInputEmpty() {
            const isAccountEmpty = !!this.getInputEmptyProblemMessage(this.account, false);

            if (isAccountEmpty) {
                return true;
            }

            if (this.account.type === this.allAccountTypes.MultiSubAccounts.type) {
                for (let i = 0; i < this.subAccounts.length; i++) {
                    const isSubAccountEmpty = !!this.getInputEmptyProblemMessage(this.subAccounts[i], true);

                    if (isSubAccountEmpty) {
                        return true;
                    }
                }
            }

            return false;
        },
        getInputEmptyProblemMessage(account, isSubAccount) {
            if (!isSubAccount && !account.category) {
                return 'Account category cannot be blank';
            } else if (!isSubAccount && !account.type) {
                return 'Account type cannot be blank';
            } else if (!account.name) {
                return 'Account name cannot be blank';
            } else if (account.type === this.allAccountTypes.SingleAccount.type && !account.currency) {
                return 'Account currency cannot be blank';
            } else {
                return null;
            }
        }
    }
}
</script>

<style>

.account-edit-balancetime .item-title {
    width: 100%;
}

.account-edit-balancetime .item-title > .item-header > .account-edit-balancetime-header {
    display: block;
    width: 100%;
}

.account-edit-balancetime .item-title > .account-edit-balancetime-title {
    display: flex;
    width: 100%;
}

.account-edit-balancetime .item-title > .account-edit-balancetime-title > .account-edit-balancetime-time {
    flex-grow: 1;
    overflow: hidden;
    text-overflow: ellipsis;
}

.subaccount-edit-list {
    --f7-list-group-title-height: 40px;
}

.subaccount-delete-button {
    margin-left: auto;
}
</style>
