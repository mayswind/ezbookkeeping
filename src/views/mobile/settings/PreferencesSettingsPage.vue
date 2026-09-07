<template>
    <f7-page>
        <f7-navbar :title="tt('Preferences')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('General Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item>
                <template #after-title>
                    {{ tt('Show Account Balance') }}
                </template>
                <template #after>
                    <f7-toggle :checked="showAccountBalance" @toggle:change="showAccountBalance = $event"></f7-toggle>
                </template>
            </f7-list-item>
            <f7-list-item
                class="item-truncate-after-text"
                link="/settings/account_category_display_order">
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Account Category Order') }}</span>
                    </div>
                </template>
                <template #after>
                    <div>{{ accountCategorysDisplayOrderContent }}</div>
                </template>
            </f7-list-item>
            <f7-list-item
                class="item-truncate-after-text"
                link="/settings/chart_color_scheme">
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Chart Color Scheme') }}</span>
                    </div>
                </template>
                <template #after>
                    <div>{{ chartColorSchemeContent }}</div>
                </template>
            </f7-list-item>
            <f7-list-item>
                <template #after-title>
                    {{ tt('Auto-update Exchange Rates Data') }}
                </template>
                <template #after>
                    <f7-toggle :checked="isAutoUpdateExchangeRatesData" @toggle:change="isAutoUpdateExchangeRatesData = $event"></f7-toggle>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Overview Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="/settings/overview_layout">
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Home Page Layout') }}</span>
                    </div>
                </template>
                <template #after>
                    <div>{{ overviewPageLayoutDisplayContent }}</div>
                </template>
            </f7-list-item>

            <f7-list-item>
                <template #after-title>
                    {{ tt('Show Amount') }}
                </template>
                <template #after>
                    <f7-toggle :checked="showAmountInHomePage" @toggle:change="showAmountInHomePage = $event"></f7-toggle>
                </template>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showTimezoneUsedForStatisticsInHomePagePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Timezone Used for Statistics') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTimezoneTypesUsedForStatistics, timezoneUsedForStatisticsInHomePage) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Timezone Used for Statistics')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Timezone Type')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTimezoneTypesUsedForStatistics"
                                           v-model:show="showTimezoneUsedForStatisticsInHomePagePopup"
                                           v-model="timezoneUsedForStatisticsInHomePage">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="/settings/filter/account?type=homePageOverview"
                :disabled="!hasAnyAccount">
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Accounts Included in Overview Statistics') }}</span>
                    </div>
                </template>
                <template #after>
                    <f7-preloader v-if="loadingAccounts" />
                    <div v-else-if="!loadingAccounts">{{ accountsIncludedInHomePageOverviewDisplayContent }}</div>
                </template>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                :disabled="!hasAnyTransactionCategory"
                :link="`/settings/filter/category?type=homePageOverview&allowCategoryTypes=${CategoryType.Income},${CategoryType.Expense}`">
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Transaction Categories Included in Overview Statistics') }}</span>
                    </div>
                </template>
                <template #after>
                    <f7-preloader v-if="loadingTransactionCategories" />
                    <div v-else-if="!loadingTransactionCategories">{{ transactionCategoriesIncludedInHomePageOverviewDisplayContent }}</div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Transaction List Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item>
                <template #after-title>
                    {{ tt('Show Monthly Total Amount') }}
                </template>
                <template #after>
                    <f7-toggle :checked="showTotalAmountInTransactionListPage" @toggle:change="showTotalAmountInTransactionListPage = $event"></f7-toggle>
                </template>
            </f7-list-item>
            <f7-list-item>
                <template #after-title>
                    {{ tt('Show Transaction Tags') }}
                </template>
                <template #after>
                    <f7-toggle :checked="showTagInTransactionListPage" @toggle:change="showTagInTransactionListPage = $event"></f7-toggle>
                </template>
            </f7-list-item>
            <f7-list-item
                link="#"
                class="item-truncate-after-text"
                popover-open=".default-keyword-search-matching-mode-popover-menu"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Keyword Search Matching Mode') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allKeywordMatchModes, defaultKeywordMatchModeInTransactionListPage) }}
                </template>
                <f7-popover class="default-keyword-search-matching-mode-popover-menu">
                    <f7-list dividers>
                        <f7-list-item link="#" no-chevron popover-close
                                      :title="option.displayName"
                                      :class="{ 'list-item-selected': defaultKeywordMatchModeInTransactionListPage === option.type }"
                                      :key="option.type"
                                      v-for="option in allKeywordMatchModes"
                                      @click="defaultKeywordMatchModeInTransactionListPage = option.type">
                            <template #after>
                                <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="defaultKeywordMatchModeInTransactionListPage === option.type"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-popover>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Transaction Edit Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                link="#"
                class="item-truncate-after-text"
                popover-open=".quick-add-button-style-popover-menu"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Quick Save Button Style') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTransactionQuickSaveButtonStyles, quickSaveButtonStyleInMobileTransactionListPage) }}
                </template>
                <f7-popover class="quick-add-button-style-popover-menu">
                    <f7-list dividers>
                        <f7-list-item link="#" no-chevron popover-close
                                      :title="option.displayName"
                                      :class="{ 'list-item-selected': quickSaveButtonStyleInMobileTransactionListPage === option.type }"
                                      :key="option.type"
                                      v-for="option in allTransactionQuickSaveButtonStyles"
                                      @click="quickSaveButtonStyleInMobileTransactionListPage = option.type">
                            <template #after>
                                <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="quickSaveButtonStyleInMobileTransactionListPage === option.type"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-popover>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="item-truncate-after-text"
                :disabled="quickSaveButtonStyleInMobileTransactionListPage === TransactionQuickSaveButtonStyle.Disabled.type"
                popover-open=".quick-add-button-action-popover-menu"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Quick Add Button Action') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTransactionQuickAddButtonActionTypes, quickAddButtonActionInMobileTransactionEditPage) }}
                </template>
                <f7-popover class="quick-add-button-action-popover-menu">
                    <f7-list dividers>
                        <f7-list-item link="#" no-chevron popover-close
                                      :title="option.displayName"
                                      :class="{ 'list-item-selected': quickAddButtonActionInMobileTransactionEditPage === option.type }"
                                      :key="option.type"
                                      v-for="option in allTransactionQuickAddButtonActionTypes"
                                      @click="quickAddButtonActionInMobileTransactionEditPage = option.type">
                            <template #after>
                                <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="quickAddButtonActionInMobileTransactionEditPage === option.type"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-popover>
            </f7-list-item>

            <f7-list-item
                link="#"
                class="item-truncate-after-text"
                popover-open=".auto-save-draft-popover-menu"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Automatically Save Draft') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findNameByValue(allAutoSaveTransactionDraftTypes, autoSaveTransactionDraft) }}
                </template>
                <f7-popover class="auto-save-draft-popover-menu">
                    <f7-list dividers>
                        <f7-list-item link="#" no-chevron popover-close
                                      :title="option.name"
                                      :class="{ 'list-item-selected': autoSaveTransactionDraft === option.value }"
                                      :key="option.value"
                                      v-for="option in allAutoSaveTransactionDraftTypes"
                                      @click="autoSaveTransactionDraft = option.value">
                            <template #after>
                                <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="autoSaveTransactionDraft === option.value"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-popover>
            </f7-list-item>

            <f7-list-item>
                <template #after-title>
                    {{ tt('Automatically Add Geolocation') }}
                </template>
                <template #after>
                    <f7-toggle :checked="isAutoGetCurrentGeoLocation" @toggle:change="isAutoGetCurrentGeoLocation = $event"></f7-toggle>
                </template>
            </f7-list-item>

            <f7-list-item>
                <template #after-title>
                    {{ tt('Always Show Transaction Pictures') }}
                </template>
                <template #after>
                    <f7-toggle :checked="alwaysShowTransactionPicturesInMobileTransactionEditPage" @toggle:change="alwaysShowTransactionPicturesInMobileTransactionEditPage = $event"></f7-toggle>
                </template>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showTransactionPictureQualityPopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Transaction Picture Upload Quality') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allImageUploadQualityTypes, transactionPictureQuality) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Transaction Picture Upload Quality')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Transaction Picture Upload Quality')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allImageUploadQualityTypes"
                                           v-model:show="showTransactionPictureQualityPopup"
                                           v-model="transactionPictureQuality">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('AI Clipboard Text Recognition') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item>
                <template #after-title>
                    {{ tt('Always Require Confirmation of Clipboard Content Before Submission') }}
                </template>
                <template #after>
                    <f7-toggle :checked="isAlwaysRequireConfirmationOfClipboardContentBeforeSubmission" @toggle:change="isAlwaysRequireConfirmationOfClipboardContentBeforeSubmission = $event"></f7-toggle>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('AI Image Recognition') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item>
                <template #after-title>
                    {{ tt('Auto Upload AI Recognition Image as Transaction Picture') }}
                </template>
                <template #after>
                    <f7-toggle :checked="isAutoUploadTransactionPictureForAIRecognition" @toggle:change="isAutoUploadTransactionPictureForAIRecognition = $event"></f7-toggle>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Account List Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="/settings/filter/account?type=accountListTotalAmount"
                :disabled="!hasAnyVisibleAccount">
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Accounts Included in Total') }}</span>
                    </div>
                </template>
                <template #after>
                    <f7-preloader v-if="loadingAccounts" />
                    <div v-else-if="!loadingAccounts">{{ accountsIncludedInTotalDisplayContent }}</div>
                </template>
            </f7-list-item>
            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showReconciliationStatementDefaultDateRangePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Date Range for Reconciliation Statement Page') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allReconciliationStatementDateRanges, reconciliationStatementPageDefaultDateRangeTypeInMobile) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Date Range')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Date Range')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allReconciliationStatementDateRanges"
                                           v-model:show="showReconciliationStatementDefaultDateRangePopup"
                                           v-model="reconciliationStatementPageDefaultDateRangeTypeInMobile">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Exchange Rates Data Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                link="#"
                class="item-truncate-after-text"
                popover-open=".exchange-rates-data-sort-by-popover-menu"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Sort by') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allCurrencySortingTypes, currencySortByInExchangeRatesPage) }}
                </template>
                <f7-popover class="exchange-rates-data-sort-by-popover-menu">
                    <f7-list dividers>
                        <f7-list-item link="#" no-chevron popover-close
                                      :title="option.displayName"
                                      :class="{ 'list-item-selected': currencySortByInExchangeRatesPage === option.type }"
                                      :key="option.type"
                                      v-for="option in allCurrencySortingTypes"
                                      @click="currencySortByInExchangeRatesPage = option.type">
                            <template #after>
                                <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="currencySortByInExchangeRatesPage === option.type"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-popover>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useAppSettingPageBase } from '@/views/base/settings/AppSettingsPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionQuickSaveButtonStyle } from '@/core/transaction.ts';
import { DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_MOBILE } from '@/core/statistics.ts';

import { findNameByValue, findDisplayNameByType } from '@/lib/common.ts';
import { isDefaultMobileOverviewLayout, parseMobileOverviewLayout } from '@/lib/overview_layout.ts';

const {
    tt,
    getAllTransactionQuickSaveButtonStyles,
    getAllTransactionQuickAddButtonActionTypes
} = useI18n();
const { showToast } = useI18nUIComponents();
const {
    loadingAccounts,
    loadingTransactionCategories,
    hasAnyAccount,
    hasAnyVisibleAccount,
    hasAnyTransactionCategory,
    allTimezoneTypesUsedForStatistics,
    allCurrencySortingTypes,
    allKeywordMatchModes,
    allAutoSaveTransactionDraftTypes,
    allImageUploadQualityTypes,
    allReconciliationStatementDateRanges,
    isAutoUpdateExchangeRatesData,
    showAccountBalance,
    showAmountInHomePage,
    timezoneUsedForStatisticsInHomePage,
    showTotalAmountInTransactionListPage,
    showTagInTransactionListPage,
    defaultKeywordMatchModeInTransactionListPage,
    autoSaveTransactionDraft,
    isAutoGetCurrentGeoLocation,
    transactionPictureQuality,
    isAlwaysRequireConfirmationOfClipboardContentBeforeSubmission,
    isAutoUploadTransactionPictureForAIRecognition,
    currencySortByInExchangeRatesPage,
    chartColorSchemeContent,
    accountsIncludedInHomePageOverviewDisplayContent,
    accountsIncludedInTotalDisplayContent,
    accountCategorysDisplayOrderContent,
    transactionCategoriesIncludedInHomePageOverviewDisplayContent,
    getValidReconciliationStatementPageDefaultDateRangeType
} = useAppSettingPageBase();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();

const showTimezoneUsedForStatisticsInHomePagePopup = ref<boolean>(false);
const showTransactionPictureQualityPopup = ref<boolean>(false);
const showReconciliationStatementDefaultDateRangePopup = ref<boolean>(false);

const allTransactionQuickSaveButtonStyles = computed<TypeAndDisplayName[]>(() => getAllTransactionQuickSaveButtonStyles());
const allTransactionQuickAddButtonActionTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionQuickAddButtonActionTypes());

const overviewPageLayoutDisplayContent = computed<string>(() => {
    try {
        return tt(isDefaultMobileOverviewLayout(parseMobileOverviewLayout(settingsStore.appSettings.mobileOverviewPageLayout)) ? 'Default' : 'Custom');
    } catch {
        return tt('Custom');
    }
});

const quickSaveButtonStyleInMobileTransactionListPage = computed<number>({
    get: () => settingsStore.appSettings.quickSaveButtonStyleInMobileTransactionListPage,
    set: (value) => settingsStore.setQuickSaveButtonStyleInMobileTransactionListPage(value)
});

const quickAddButtonActionInMobileTransactionEditPage = computed<number>({
    get: () => settingsStore.appSettings.quickAddButtonActionInMobileTransactionEditPage,
    set: (value) => settingsStore.setQuickAddButtonActionInMobileTransactionEditPage(value)
});

const alwaysShowTransactionPicturesInMobileTransactionEditPage = computed<boolean>({
    get: () => settingsStore.appSettings.alwaysShowTransactionPicturesInMobileTransactionEditPage,
    set: (value) => settingsStore.setAlwaysShowTransactionPicturesInMobileTransactionEditPage(value)
});

const reconciliationStatementPageDefaultDateRangeTypeInMobile = computed<number>({
    get: () => getValidReconciliationStatementPageDefaultDateRangeType(settingsStore.appSettings.reconciliationStatementPageDefaultDateRangeTypeInMobile, DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_MOBILE.type),
    set: (value: number) => settingsStore.setReconciliationStatementPageDefaultDateRangeTypeInMobile(value)
});

function init(): void {
    loadingAccounts.value = true;
    loadingTransactionCategories.value = true;

    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loadingAccounts.value = false;
    }).catch(error => {
        loadingAccounts.value = false;

        if (!error.processed) {
            showToast(error.message || error);
        }
    });

    transactionCategoriesStore.loadAllCategories({
        force: false
    }).then(() => {
        loadingTransactionCategories.value = false;
    }).catch(error => {
        loadingTransactionCategories.value = false;

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

init();
</script>
