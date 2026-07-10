<template>
    <f7-page>
        <f7-navbar :title="tt('Page Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('General Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
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
        </f7-list>

        <f7-block-title>{{ tt('Overview Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
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
                class="item-truncate-after-text"
                link="#"
                @click="showKeywordMatchModeInTransactionListPagePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Keyword Search Matching Mode') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allKeywordMatchModes, defaultKeywordMatchModeInTransactionListPage) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Keyword Search Matching Mode')"
                                           :enable-filter="false"
                                           :items="allKeywordMatchModes"
                                           v-model:show="showKeywordMatchModeInTransactionListPagePopup"
                                           v-model="defaultKeywordMatchModeInTransactionListPage">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Transaction Edit Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showQuickSaveButtonStyleInMobileTransactionListPagePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Quick Save Button Style') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTransactionQuickSaveButtonStyles, quickSaveButtonStyleInMobileTransactionListPage) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Quick Save Button Style')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Quick Save Button Style')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTransactionQuickSaveButtonStyles"
                                           v-model:show="showQuickSaveButtonStyleInMobileTransactionListPagePopup"
                                           v-model="quickSaveButtonStyleInMobileTransactionListPage">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                :disabled="quickSaveButtonStyleInMobileTransactionListPage === TransactionQuickSaveButtonStyle.Disabled.type"
                @click="showQuickAddButtonActionInMobileTransactionEditPagePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Quick Add Button Action') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTransactionQuickAddButtonActionTypes, quickAddButtonActionInMobileTransactionEditPage) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Quick Add Button Action')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Quick Add Button Action')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTransactionQuickAddButtonActionTypes"
                                           v-model:show="showQuickAddButtonActionInMobileTransactionEditPagePopup"
                                           v-model="quickAddButtonActionInMobileTransactionEditPage">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showAutoSaveTransactionDraftPopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Automatically Save Draft') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findNameByValue(allAutoSaveTransactionDraftTypes, autoSaveTransactionDraft) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="value" value-field="value"
                                           title-field="name"
                                           :title="tt('Automatically Save Draft')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Automatically Save Draft')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allAutoSaveTransactionDraftTypes"
                                           v-model:show="showAutoSaveTransactionDraftPopup"
                                           v-model="autoSaveTransactionDraft">
                </list-item-selection-popup>
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
                class="item-truncate-after-text"
                link="#"
                @click="showCurrencySortByInExchangeRatesPagePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Sort by') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allCurrencySortingTypes, currencySortByInExchangeRatesPage) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Sort by')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Sort by')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCurrencySortingTypes"
                                           v-model:show="showCurrencySortByInExchangeRatesPagePopup"
                                           v-model="currencySortByInExchangeRatesPage">
                </list-item-selection-popup>
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
const showKeywordMatchModeInTransactionListPagePopup = ref<boolean>(false);
const showQuickSaveButtonStyleInMobileTransactionListPagePopup = ref<boolean>(false);
const showQuickAddButtonActionInMobileTransactionEditPagePopup = ref<boolean>(false);
const showAutoSaveTransactionDraftPopup = ref<boolean>(false);
const showTransactionPictureQualityPopup = ref<boolean>(false);
const showReconciliationStatementDefaultDateRangePopup = ref<boolean>(false);
const showCurrencySortByInExchangeRatesPagePopup = ref<boolean>(false);

const allTransactionQuickSaveButtonStyles = computed<TypeAndDisplayName[]>(() => getAllTransactionQuickSaveButtonStyles());
const allTransactionQuickAddButtonActionTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionQuickAddButtonActionTypes());

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
