<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': account.type !== AccountType.MultiSubAccounts.type }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :class="{ 'disabled': isInputEmpty() || submitting }" :text="tt(saveButtonTitle)" @click="save"></f7-link>
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
                :header="tt('Account Category')"
                :title="findDisplayNameByType(allAccountCategories, account.category)"
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
                :header="tt('Account Type')"
                :title="findDisplayNameByType(allAccountTypes, account.type)"
                @click="showAccountTypeSheet = true"
            >
                <list-item-selection-sheet value-type="item"
                                           key-field="type" value-field="type" title-field="displayName"
                                           :items="allAccountTypes"
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

        <f7-list form strong inset dividers class="margin-vertical" v-else-if="!loading && account.type === AccountType.SingleAccount.type">
            <f7-list-input
                type="text"
                clear-button
                :label="tt('Account Name')"
                :placeholder="tt('Your account name')"
                v-model:value="account.name"
            ></f7-list-input>

            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="accountContext.showIconSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ tt('Account Icon') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <icon-selection-sheet :all-icon-infos="ALL_ACCOUNT_ICONS"
                                                  :color="account.color"
                                                  v-model:show="accountContext.showIconSelectionSheet"
                                                  v-model="account.icon"
                            ></icon-selection-sheet>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="accountContext.showColorSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ tt('Account Color') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <color-selection-sheet :all-color-infos="ALL_ACCOUNT_COLORS"
                                                   v-model:show="accountContext.showColorSelectionSheet"
                                                   v-model="account.color"
                            ></color-selection-sheet>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :class="{ 'disabled': editAccountId }"
                :header="tt('Currency')"
                :no-chevron="!!editAccountId"
                @click="accountContext.showCurrencyPopup = true"
            >
                <template #title>
                    <div class="no-padding no-margin">
                        <span>{{ getCurrencyName(account.currency) }}&nbsp;</span>
                        <small class="smaller">{{ account.currency }}</small>
                    </div>
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="currencyCode" value-field="currencyCode"
                                           title-field="displayName" after-field="currencyCode"
                                           :title="tt('Currency Name')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Currency')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCurrencies"
                                           v-model:show="accountContext.showCurrencyPopup"
                                           v-model="account.currency">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Statement Date')"
                :title="getAccountCreditCardStatementDate(account.creditCardStatementDate)"
                v-if="isAccountSupportCreditCardStatementDate"
                @click="accountContext.showCreditCardStatementDatePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="day" value-field="day"
                                           title-field="displayName"
                                           :title="tt('Statement Date')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Statement Date')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allAvailableMonthDays"
                                           v-model:show="accountContext.showCreditCardStatementDatePopup"
                                           v-model="account.creditCardStatementDate">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title"
                :class="{ 'disabled': editAccountId }"
                :header="account.isLiability ? tt('Account Outstanding Balance') : tt('Account Balance')"
                :title="formatAccountDisplayBalance(account)"
                @click="accountContext.showBalanceSheet = true"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  :currency="account.currency"
                                  :flip-negative="account.isLiability"
                                  v-model:show="accountContext.showBalanceSheet"
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
                    <div class="account-edit-balancetime-header" @click="showDateTimeDialog(accountContext, 'time')">{{ tt('Balance Time') }}</div>
                </template>
                <template #title>
                    <div class="account-edit-balancetime-title">
                        <div @click="showDateTimeDialog(accountContext, 'date')">{{ formatAccountBalanceDate(account) }}</div>&nbsp;<div class="account-edit-balancetime-time" @click="showDateTimeDialog(accountContext, 'time')">{{ formatAccountBalanceTime(account) }}</div>
                    </div>
                </template>
                <date-time-selection-sheet :init-mode="accountContext.balanceDateTimeSheetMode"
                                           v-model:show="accountContext.showBalanceDateTimeSheet"
                                           v-model="account.balanceTime">
                </date-time-selection-sheet>
            </f7-list-item>

            <f7-list-item :title="tt('Visible')" v-if="editAccountId">
                <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                style="height: auto"
                :label="tt('Description')"
                :placeholder="tt('Your account description (optional)')"
                v-textarea-auto-size
                v-model:value="account.comment"
            ></f7-list-input>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-else-if="!loading && account.type === AccountType.MultiSubAccounts.type">
            <f7-list-input
                type="text"
                clear-button
                :label="tt('Account Name')"
                :placeholder="tt('Your account name')"
                v-model:value="account.name"
            ></f7-list-input>

            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="accountContext.showIconSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ tt('Account Icon') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <icon-selection-sheet :all-icon-infos="ALL_ACCOUNT_ICONS"
                                                  :color="account.color"
                                                  v-model:show="accountContext.showIconSelectionSheet"
                                                  v-model="account.icon"
                            ></icon-selection-sheet>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="accountContext.showColorSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ tt('Account Color') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="account.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <color-selection-sheet :all-color-infos="ALL_ACCOUNT_COLORS"
                                                   v-model:show="accountContext.showColorSelectionSheet"
                                                   v-model="account.color"
                            ></color-selection-sheet>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Statement Date')"
                :title="getAccountCreditCardStatementDate(account.creditCardStatementDate)"
                v-if="isAccountSupportCreditCardStatementDate"
                @click="accountContext.showCreditCardStatementDatePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="day" value-field="day"
                                           title-field="displayName"
                                           :title="tt('Statement Date')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Statement Date')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allAvailableMonthDays"
                                           v-model:show="accountContext.showCreditCardStatementDatePopup"
                                           v-model="account.creditCardStatementDate">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item :title="tt('Visible')" v-if="editAccountId">
                <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                style="height: auto"
                :label="tt('Description')"
                :placeholder="tt('Your account description (optional)')"
                v-textarea-auto-size
                v-model:value="account.comment"
            ></f7-list-input>
        </f7-list>

        <f7-block class="no-padding no-margin" v-if="!loading && account.type === AccountType.MultiSubAccounts.type">
            <f7-list strong inset dividers class="subaccount-edit-list margin-vertical"
                     :key="idx"
                     v-for="(subAccount, idx) in subAccounts">
                <f7-list-item group-title>
                    <small>{{ tt('Sub Account') + ' #' + (idx + 1) }}</small>
                    <f7-button rasied fill class="subaccount-delete-button" color="red" icon-f7="trash" icon-size="16px"
                               :tooltip="tt('Remove Sub-account')"
                               @click="removeSubAccount(subAccount, false)">
                    </f7-button>
                </f7-list-item>

                <f7-list-input
                    type="text"
                    clear-button
                    :label="tt('Sub-account Name')"
                    :placeholder="tt('Your sub-account name')"
                    v-model:value="subAccount.name"
                ></f7-list-input>

                <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                    <template #default>
                        <div class="grid grid-cols-2">
                            <div class="list-item-subitem no-chevron">
                                <a class="item-link" href="#" @click="subAccountContexts[idx].showIconSelectionSheet = true">
                                    <div class="item-content">
                                        <div class="item-inner">
                                            <div class="item-header">
                                                <span>{{ tt('Sub-account Icon') }}</span>
                                            </div>
                                            <div class="item-title">
                                                <div class="list-item-custom-title no-padding">
                                                    <ItemIcon icon-type="account" :icon-id="subAccount.icon" :color="subAccount.color"></ItemIcon>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </a>

                                <icon-selection-sheet :all-icon-infos="ALL_ACCOUNT_ICONS"
                                                      :color="subAccount.color"
                                                      v-model:show="subAccountContexts[idx].showIconSelectionSheet"
                                                      v-model="subAccount.icon"
                                ></icon-selection-sheet>
                            </div>
                            <div class="list-item-subitem no-chevron">
                                <a class="item-link" href="#" @click="subAccountContexts[idx].showColorSelectionSheet = true">
                                    <div class="item-content">
                                        <div class="item-inner">
                                            <div class="item-header">
                                                <span>{{ tt('Sub-account Color') }}</span>
                                            </div>
                                            <div class="item-title">
                                                <div class="list-item-custom-title no-padding">
                                                    <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="subAccount.color"></ItemIcon>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </a>

                                <color-selection-sheet :all-color-infos="ALL_ACCOUNT_COLORS"
                                                       v-model:show="subAccountContexts[idx].showColorSelectionSheet"
                                                       v-model="subAccount.color"
                                ></color-selection-sheet>
                            </div>
                        </div>
                    </template>
                </f7-list-item>

                <f7-list-item
                    class="list-item-with-header-and-title list-item-no-item-after"
                    link="#"
                    :class="{ 'disabled': editAccountId && !isNewAccount(subAccount) }"
                    :header="tt('Currency')"
                    :no-chevron="!!editAccountId && !isNewAccount(subAccount)"
                    @click="subAccountContexts[idx].showCurrencyPopup = true"
                >
                    <template #title>
                        <div class="no-padding no-margin">
                            <span>{{ getCurrencyName(subAccount.currency) }}&nbsp;</span>
                            <small class="smaller">{{ subAccount.currency }}</small>
                        </div>
                    </template>
                    <list-item-selection-popup value-type="item"
                                               key-field="currencyCode" value-field="currencyCode"
                                               title-field="displayName" after-field="currencyCode"
                                               :title="tt('Currency Name')"
                                               :enable-filter="true"
                                               :filter-placeholder="tt('Currency')"
                                               :filter-no-items-text="tt('No results')"
                                               :items="allCurrencies"
                                               v-model:show="subAccountContexts[idx].showCurrencyPopup"
                                               v-model="subAccount.currency">
                    </list-item-selection-popup>
                </f7-list-item>

                <f7-list-item
                    link="#" no-chevron
                    class="list-item-with-header-and-title"
                    :class="{ 'disabled': editAccountId && !isNewAccount(subAccount) }"
                    :header="account.isLiability ? tt('Sub-account Outstanding Balance') : tt('Sub-account Balance')"
                    :title="formatAccountDisplayBalance(subAccount)"
                    @click="subAccountContexts[idx].showBalanceSheet = true"
                >
                    <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                      :max-value="TRANSACTION_MAX_AMOUNT"
                                      :currency="subAccount.currency"
                                      :flip-negative="account.isLiability"
                                      v-model:show="subAccountContexts[idx].showBalanceSheet"
                                      v-model="subAccount.balance"
                    ></number-pad-sheet>
                </f7-list-item>

                <f7-list-item
                    class="account-edit-balancetime list-item-with-header-and-title"
                    link="#" no-chevron
                    v-show="subAccount.balance"
                    v-if="!editAccountId || isNewAccount(subAccount)"
                >
                    <template #header>
                        <div class="account-edit-balancetime-header" @click="showDateTimeDialog(subAccountContexts[idx], 'time')">{{ tt('Sub-account Balance Time') }}</div>
                    </template>
                    <template #title>
                        <div class="account-edit-balancetime-title">
                            <div @click="showDateTimeDialog(subAccountContexts[idx], 'date')">{{ formatAccountBalanceDate(subAccount) }}</div>&nbsp;<div class="account-edit-balancetime-time" @click="showDateTimeDialog(subAccountContexts[idx], 'time')">{{ formatAccountBalanceTime(subAccount) }}</div>
                        </div>
                    </template>
                    <date-time-selection-sheet :init-mode="subAccountContexts[idx].balanceDateTimeSheetMode"
                                               v-model:show="subAccountContexts[idx].showBalanceDateTimeSheet"
                                               v-model="subAccount.balanceTime">
                    </date-time-selection-sheet>
                </f7-list-item>

                <f7-list-item :title="tt('Visible')" v-if="editAccountId && !isNewAccount(subAccount)">
                    <f7-toggle :checked="subAccount.visible" @toggle:change="subAccount.visible = $event"></f7-toggle>
                </f7-list-item>

                <f7-list-input
                    type="textarea"
                    style="height: auto"
                    :label="tt('Description')"
                    :placeholder="tt('Your sub-account description (optional)')"
                    v-textarea-auto-size
                    v-model:value="subAccount.comment"
                ></f7-list-input>
            </f7-list>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="addSubAccountAndContext">{{ tt('Add Sub-account') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to remove this sub-account?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="removeSubAccount(subAccountToDelete, true)">{{ tt('Remove') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useAccountEditPageBaseBase } from '@/views/base/accounts/AccountEditPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';

import type { LocalizedCurrencyInfo } from '@/core/currency.ts';
import { AccountType } from '@/core/account.ts';
import { ALL_ACCOUNT_ICONS } from '@/consts/icon.ts';
import { ALL_ACCOUNT_COLORS } from '@/consts/color.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import type { Account } from '@/models/account.ts';

import { isDefined, findDisplayNameByType } from '@/lib/common.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import {
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore
} from '@/lib/datetime.ts';

interface AccountContext {
    showIconSelectionSheet: boolean;
    showColorSelectionSheet: boolean;
    showCurrencyPopup: boolean;
    showCreditCardStatementDatePopup: boolean;
    showBalanceSheet: boolean;
    showBalanceDateTimeSheet: boolean;
    balanceDateTimeSheetMode: string;
}

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const { tt, getAllCurrencies, getCurrencyName, formatUnixTimeToLongDate, formatUnixTimeToLongTime, formatAmountWithCurrency } = useI18n();
const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();
const {
    editAccountId,
    clientSessionId,
    loading,
    submitting,
    account,
    subAccounts,
    title,
    saveButtonTitle,
    allAccountCategories,
    allAccountTypes,
    allAvailableMonthDays,
    isAccountSupportCreditCardStatementDate,
    getAccountCreditCardStatementDate,
    isNewAccount,
    isInputEmpty,
    getAccountOrSubAccountProblemMessage,
    addSubAccount,
    setAccount
} = useAccountEditPageBaseBase();

const accountsStore = useAccountsStore();

const DEFAULT_ACCOUNT_CONTEXT: AccountContext = {
    showIconSelectionSheet: false,
    showColorSelectionSheet: false,
    showCurrencyPopup: false,
    showCreditCardStatementDatePopup: false,
    showBalanceSheet: false,
    showBalanceDateTimeSheet: false,
    balanceDateTimeSheetMode: 'time'
};

const accountContext = ref<AccountContext>(Object.assign({}, DEFAULT_ACCOUNT_CONTEXT));
const subAccountContexts = ref<AccountContext[]>([]);
const subAccountToDelete = ref<Account | null>(null);
const loadingError = ref<unknown | null>(null);
const showAccountCategorySheet = ref<boolean>(false);
const showAccountTypeSheet = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);

const allCurrencies = computed<LocalizedCurrencyInfo[]>(() => getAllCurrencies());

function formatAccountDisplayBalance(selectedAccount: Account): string {
    const balance = account.value.isLiability ? -selectedAccount.balance : selectedAccount.balance;
    return formatAmountWithCurrency(balance, selectedAccount.currency);
}

function formatAccountBalanceDate(account: Account): string {
    if (!isDefined(account.balanceTime)) {
        return '';
    }

    return formatUnixTimeToLongDate(getActualUnixTimeForStore(account.balanceTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
}

function formatAccountBalanceTime(account: Account): string {
    if (!isDefined(account.balanceTime)) {
        return '';
    }

    return formatUnixTimeToLongTime(getActualUnixTimeForStore(account.balanceTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
}

function init(): void {
    const query = props.f7route.query;
    clientSessionId.value = generateRandomUUID();

    if (query['id']) {
        loading.value = true;

        editAccountId.value = query['id'];

        accountsStore.getAccount({
            accountId: editAccountId.value
        }).then(response => {
            setAccount(response);
            subAccountContexts.value = [];

            for (let i = 0; i < subAccounts.value.length; i++) {
                subAccountContexts.value.push(Object.assign({}, DEFAULT_ACCOUNT_CONTEXT));
            }

            loading.value = false;
        }).catch(error => {
            if (error.processed) {
                loading.value = false;
            } else {
                loadingError.value = error;
                showToast(error.message || error);
            }
        });
    } else {
        loading.value = false;
    }
}

function save(): void {
    const router = props.f7router;
    const problemMessage = getAccountOrSubAccountProblemMessage();

    if (problemMessage) {
        showAlert(problemMessage);
        return;
    }

    submitting.value = true;
    showLoading(() => submitting.value);

    accountsStore.saveAccount({
        account: account.value,
        subAccounts: subAccounts.value,
        isEdit: !!editAccountId.value,
        clientSessionId: clientSessionId.value
    }).then(() => {
        submitting.value = false;
        hideLoading();

        if (!editAccountId.value) {
            showToast('You have added a new account');
        } else {
            showToast('You have saved this account');
        }

        router.back();
    }).catch(error => {
        submitting.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function addSubAccountAndContext(): void {
    if (addSubAccount()) {
        subAccountContexts.value.push(Object.assign({}, DEFAULT_ACCOUNT_CONTEXT));
    }
}

function removeSubAccount(subAccount: Account | null, confirm: boolean): void {
    if (!subAccount) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        subAccountToDelete.value = subAccount;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    subAccountToDelete.value = null;

    for (let i = 0; i < subAccounts.value.length; i++) {
        if (subAccounts.value[i] === subAccount) {
            subAccounts.value.splice(i, 1);
            subAccountContexts.value.splice(i, 1);
        }
    }
}

function showDateTimeDialog(accountContext: AccountContext, sheetMode: string): void {
    accountContext.balanceDateTimeSheetMode = sheetMode;
    accountContext.showBalanceDateTimeSheet = true;
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

watch(() => account.value.type, () => {
    if (subAccounts.value.length < 1) {
        addSubAccountAndContext();
    }
});

init();
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
