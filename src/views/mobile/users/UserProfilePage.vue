<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('User Profile')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !isUserVerifyEmailEnabled() || loading || emailVerified }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :class="{ 'disabled': inputIsNotChanged || inputIsInvalid || saving }" :text="tt('Save')" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-input label="Password" placeholder="Your password"></f7-list-input>
            <f7-list-input label="Confirm Password" placeholder="Re-enter the password"></f7-list-input>
            <f7-list-input label="E-mail" placeholder="Your email address"></f7-list-input>
            <f7-list-input label="Nickname" placeholder="Your nickname"></f7-list-input>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Account" title="Unspecified"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Editable Transaction Range" title="All" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Language" title="Language" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Currency" title="Currency" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="First Day of Week" title="Week Day" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Long Date Format" title="YYYY-MM-DD" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Short Date Format" title="YYYY-MM-DD" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Long Time Format" title="HH:mm:ss" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Short Time Format" title="HH:mm" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Currency Display Mode" title="Currency Symbol" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Digit Grouping" title="Thousands Separator" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Digit Grouping Symbol" title="Comma (,)" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Decimal Separator" title="Dot (.)" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Expense Amount Color" title="Amount Color" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Income Amount Color" title="Amount Color" link="#"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="tt('Password')"
                :placeholder="tt('Your password')"
                v-model:value="newProfile.password"
            ></f7-list-input>

            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="tt('Confirm Password')"
                :placeholder="tt('Re-enter the password')"
                v-model:value="newProfile.confirmPassword"
            ></f7-list-input>

            <f7-list-input
                type="email"
                autocomplete="email"
                clear-button
                :label="tt('E-mail') + ' ' + (emailVerified ? tt('(Verified)') : tt('(Not Verified)'))"
                :placeholder="tt('Your email address')"
                v-model:value="newProfile.email"
            ></f7-list-input>

            <f7-list-input
                type="text"
                autocomplete="nickname"
                clear-button
                :label="tt('Nickname')"
                :placeholder="tt('Your nickname')"
                v-model:value="newProfile.nickname"
            ></f7-list-input>

            <f7-list-item class="ebk-list-item-error-info" v-if="inputIsInvalid" :footer="tt(inputInvalidProblemMessage || '')"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': !allVisibleAccounts.length }"
                :header="tt('Default Account')"
                :title="Account.findAccountNameById(allAccounts, newProfile.defaultAccountId, tt('Unspecified'))"
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
                                                      :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                      :items="allVisibleCategorizedAccounts"
                                                      v-model:show="showAccountSheet"
                                                      v-model="newProfile.defaultAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Editable Transaction Range')"
                :title="findDisplayNameByType(allTransactionEditScopeTypes, newProfile.transactionEditScope)"
                @click="showEditableTransactionRangePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Editable Transaction Range')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Date Range')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTransactionEditScopeTypes"
                                           v-model:show="showEditableTransactionRangePopup"
                                           v-model="newProfile.transactionEditScope">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item class="ebk-list-item-error-info" v-if="extendInputIsInvalid" :footer="tt(extendInputInvalidProblemMessage || '')"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :header="languageTitle"
                :title="currentLanguageName"
                @click="showLanguagePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="languageTag" value-field="languageTag"
                                           title-field="nativeDisplayName" after-field="displayName"
                                           :title="languageTitle"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Language')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allLanguages"
                                           v-model:show="showLanguagePopup"
                                           v-model="newProfile.language">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :header="tt('Default Currency')"
                @click="showDefaultCurrencyPopup = true"
            >
                <template #title>
                    <f7-block class="no-padding no-margin">
                        <span>{{ getCurrencyName(newProfile.defaultCurrency) }}&nbsp;</span>
                        <small class="smaller">{{ newProfile.defaultCurrency }}</small>
                    </f7-block>
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="currencyCode" value-field="currencyCode"
                                           title-field="displayName" after-field="currencyCode"
                                           :title="tt('Default Currency')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Currency')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCurrencies"
                                           v-model:show="showDefaultCurrencyPopup"
                                           v-model="newProfile.defaultCurrency">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('First Day of Week')"
                :title="currentDayOfWeekName"
                @click="showFirstDayOfWeekPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('First Day of Week')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Date')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allWeekDays"
                                           v-model:show="showFirstDayOfWeekPopup"
                                           v-model="newProfile.firstDayOfWeek">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Fiscal Year Start Date')"
                :title="currentFiscalYearStartDate"
                @click="showFiscalYearStartSheet = true"
            >
                <fiscal-year-start-selection-sheet
                    v-model:show="showFiscalYearStartSheet"
                    v-model="newProfile.fiscalYearStart"
                    v-model:title="currentFiscalYearStartDate">
                </fiscal-year-start-selection-sheet>
            </f7-list-item>

        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Long Date Format')"
                :title="findDisplayNameByType(allLongDateFormats, newProfile.longDateFormat)"
                @click="showLongDateFormatPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Long Date Format')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Long Date Format')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allLongDateFormats"
                                           v-model:show="showLongDateFormatPopup"
                                           v-model="newProfile.longDateFormat">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Short Date Format')"
                :title="findDisplayNameByType(allShortDateFormats, newProfile.shortDateFormat)"
                @click="showShortDateFormatPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Short Date Format')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Short Date Format')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allShortDateFormats"
                                           v-model:show="showShortDateFormatPopup"
                                           v-model="newProfile.shortDateFormat">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Long Time Format')"
                :title="findDisplayNameByType(allLongTimeFormats, newProfile.longTimeFormat)"
                @click="showLongTimeFormatPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Long Time Format')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Long Time Format')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allLongTimeFormats"
                                           v-model:show="showLongTimeFormatPopup"
                                           v-model="newProfile.longTimeFormat">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Short Time Format')"
                :title="findDisplayNameByType(allShortTimeFormats, newProfile.shortTimeFormat)"
                @click="showShortTimeFormatPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Short Time Format')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Short Time Format')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allShortTimeFormats"
                                           v-model:show="showShortTimeFormatPopup"
                                           v-model="newProfile.shortTimeFormat">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Fiscal Year Format')"
                :title="findDisplayNameByType(allFiscalYearFormats, newProfile.fiscalYearFormat)"
                @click="showFiscalYearFormatPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Fiscal Year Format')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Fiscal Year Format')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allFiscalYearFormats"
                                           v-model:show="showFiscalYearFormatPopup"
                                           v-model="newProfile.fiscalYearFormat">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Currency Display Mode')"
                :title="findDisplayNameByType(allCurrencyDisplayTypes, newProfile.currencyDisplayType)"
                @click="showCurrencyDisplayTypePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Currency Display Mode')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Currency Display Mode')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCurrencyDisplayTypes"
                                           v-model:show="showCurrencyDisplayTypePopup"
                                           v-model="newProfile.currencyDisplayType">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Digit Grouping')"
                :title="findDisplayNameByType(allDigitGroupingTypes, newProfile.digitGrouping)"
                @click="showDigitGroupingPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Digit Grouping')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Digit Grouping')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allDigitGroupingTypes"
                                           v-model:show="showDigitGroupingPopup"
                                           v-model="newProfile.digitGrouping">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :disabled="!supportDigitGroupingSymbol"
                :header="tt('Digit Grouping Symbol')"
                :title="findDisplayNameByType(allDigitGroupingSymbols, newProfile.digitGroupingSymbol)"
                @click="showDigitGroupingSymbolPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Digit Grouping Symbol')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Digit Grouping Symbol')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allDigitGroupingSymbols"
                                           v-model:show="showDigitGroupingSymbolPopup"
                                           v-model="newProfile.digitGroupingSymbol">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Decimal Separator')"
                :title="findDisplayNameByType(allDecimalSeparators, newProfile.decimalSeparator)"
                @click="showDecimalSeparatorPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Decimal Separator')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Decimal Separator')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allDecimalSeparators"
                                           v-model:show="showDecimalSeparatorPopup"
                                           v-model="newProfile.decimalSeparator">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Expense Amount Color')"
                :title="findDisplayNameByType(allExpenseAmountColorTypes, newProfile.expenseAmountColor)"
                @click="showExpenseAmountColorPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Expense Amount Color')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Color')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allExpenseAmountColorTypes"
                                           v-model:show="showExpenseAmountColorPopup"
                                           v-model="newProfile.expenseAmountColor">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="tt('Income Amount Color')"
                :title="findDisplayNameByType(allIncomeAmountColorTypes, newProfile.incomeAmountColor)"
                @click="showIncomeAmountColorPopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Income Amount Color')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Color')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allIncomeAmountColorTypes"
                                           v-model:show="showIncomeAmountColorPopup"
                                           v-model="newProfile.incomeAmountColor">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item class="ebk-list-item-error-info" v-if="langAndRegionInputIsInvalid" :footer="tt(langAndRegionInputInvalidProblemMessage || '')"></f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading || resending }" @click="resendVerifyEmail"
                                   v-if="isUserVerifyEmailEnabled() && !loading && !emailVerified"
                >{{ tt('Resend Validation Email') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <password-input-sheet :title="tt('Modify Password')"
                              :hint="tt('Please enter your current password when modifying your password')"
                              :confirm-disabled="saving"
                              :cancel-disabled="saving"
                              v-model:show="showInputPasswordSheet"
                              v-model="currentPassword"
                              @password:confirm="save()">
        </password-input-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useUserProfilePageBase } from '@/views/base/users/UserProfilePageBase.ts';

import { useRootStore } from '@/stores/index.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import type { LocalizedCurrencyInfo } from '@/core/currency.ts';

import type { UserProfileResponse } from '@/models/user.ts';
import { Account } from '@/models/account.ts';

import { findDisplayNameByType } from '@/lib/common.ts';
import { isUserVerifyEmailEnabled } from '@/lib/server_settings.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, getAllLanguageOptions, getAllCurrencies, getCurrencyName, getCurrentFiscalYearStartFormatted } = useI18n();
const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();

const {
    newProfile,
    emailVerified,
    loading,
    resending,
    saving,
    allAccounts,
    allVisibleAccounts,
    allVisibleCategorizedAccounts,
    allWeekDays,
    allLongDateFormats,
    allShortDateFormats,
    allLongTimeFormats,
    allShortTimeFormats,
    allFiscalYearFormats,
    allDecimalSeparators,
    allDigitGroupingSymbols,
    allDigitGroupingTypes,
    allCurrencyDisplayTypes,
    allExpenseAmountColorTypes,
    allIncomeAmountColorTypes,
    allTransactionEditScopeTypes,
    languageTitle,
    supportDigitGroupingSymbol,
    inputIsNotChangedProblemMessage,
    inputInvalidProblemMessage,
    langAndRegionInputInvalidProblemMessage,
    extendInputInvalidProblemMessage,
    inputIsNotChanged,
    inputIsInvalid,
    langAndRegionInputIsInvalid,
    extendInputIsInvalid,
    setCurrentUserProfile,
    doAfterProfileUpdate
} = useUserProfilePageBase();

const rootStore = useRootStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();

const currentPassword = ref<string>('');
const loadingError = ref<unknown | null>(null);
const showInputPasswordSheet = ref<boolean>(false);
const showAccountSheet = ref<boolean>(false);
const showEditableTransactionRangePopup = ref<boolean>(false);
const showLanguagePopup = ref<boolean>(false);
const showDefaultCurrencyPopup = ref<boolean>(false);
const showFirstDayOfWeekPopup = ref<boolean>(false);
const showFiscalYearStartSheet = ref<boolean>(false);
const showLongDateFormatPopup = ref<boolean>(false);
const showShortDateFormatPopup = ref<boolean>(false);
const showLongTimeFormatPopup = ref<boolean>(false);
const showShortTimeFormatPopup = ref<boolean>(false);
const showFiscalYearFormatPopup = ref<boolean>(false);
const showCurrencyDisplayTypePopup = ref<boolean>(false);
const showDigitGroupingPopup = ref<boolean>(false);
const showDigitGroupingSymbolPopup = ref<boolean>(false);
const showDecimalSeparatorPopup = ref<boolean>(false);
const showExpenseAmountColorPopup = ref<boolean>(false);
const showIncomeAmountColorPopup = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);

const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(true));
const allCurrencies = computed<LocalizedCurrencyInfo[]>(() => getAllCurrencies());

const currentLanguageName = computed<string>(() => {
    for (let i = 0; i < allLanguages.value.length; i++) {
        if (allLanguages.value[i].languageTag === newProfile.value.language) {
            return allLanguages.value[i].nativeDisplayName;
        }
    }

    return tt('Unknown');
});

const currentDayOfWeekName = computed<string | null>(() => findDisplayNameByType(allWeekDays.value, newProfile.value.firstDayOfWeek));
const currentFiscalYearStartDate = ref<string>(getCurrentFiscalYearStartFormatted());

function init(): void {
    loading.value = true;

    const promises = [
        accountsStore.loadAllAccounts({ force: false }),
        userStore.getCurrentUserProfile()
    ];

    Promise.all(promises).then(responses => {
        const profile = responses[1] as UserProfileResponse;
        setCurrentUserProfile(profile);
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function save(): void {
    const router = props.f7router;

    showInputPasswordSheet.value = false;

    const problemMessage = inputIsNotChangedProblemMessage.value || inputInvalidProblemMessage.value || extendInputInvalidProblemMessage.value || langAndRegionInputInvalidProblemMessage.value;

    if (problemMessage) {
        showAlert(problemMessage);
        return;
    }

    if (newProfile.value.password && !currentPassword.value) {
        showInputPasswordSheet.value = true;
        return;
    }

    saving.value = true;
    showLoading(() => saving.value);

    rootStore.updateUserProfile(newProfile.value.toProfileUpdateRequest(currentPassword.value)).then(response => {
        saving.value = false;
        hideLoading();
        currentPassword.value = '';

        doAfterProfileUpdate(response.user);
        showToast('Your profile has been successfully updated');
        router.back();
    }).catch(error => {
        saving.value = false;
        hideLoading();
        currentPassword.value = '';

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function resendVerifyEmail(): void {
    resending.value = true;
    showLoading(() => resending.value);

    rootStore.resendVerifyEmailByLoginedUser().then(() => {
        resending.value = false;
        hideLoading();

        showToast('Validation email has been sent');
    }).catch(error => {
        resending.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>
