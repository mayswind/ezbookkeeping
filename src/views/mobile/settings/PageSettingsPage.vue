<template>
    <f7-page>
        <f7-navbar :title="tt('Page Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('Overview Page') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item>
                <span>{{ tt('Show Amount') }}</span>
                <f7-toggle :checked="showAmountInHomePage" @toggle:change="showAmountInHomePage = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item
                link="#"
                :title="tt('Timezone Used for Statistics')"
                :after="findDisplayNameByType(allTimezoneTypesUsedForStatistics, timezoneUsedForStatisticsInHomePage)"
                @click="showTimezoneUsedForStatisticsInHomePagePopup = true"
            >
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
        </f7-list>

        <f7-block-title>{{ tt('Transaction List Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ tt('Show Monthly Total Amount') }}</span>
                <f7-toggle :checked="showTotalAmountInTransactionListPage" @toggle:change="showTotalAmountInTransactionListPage = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-item>
                <span>{{ tt('Show Transaction Tag') }}</span>
                <f7-toggle :checked="showTagInTransactionListPage" @toggle:change="showTagInTransactionListPage = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Transaction Edit Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                link="#"
                :title="tt('Automatically Save Draft')"
                :after="findNameByValue(allAutoSaveTransactionDraftTypes, autoSaveTransactionDraft)"
                @click="showAutoSaveTransactionDraftPopup = true"
            >
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
                <span>{{ tt('Automatically Add Geolocation') }}</span>
                <f7-toggle :checked="isAutoGetCurrentGeoLocation" @toggle:change="isAutoGetCurrentGeoLocation = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ tt('Always Show Transaction Pictures') }}</span>
                <f7-toggle :checked="alwaysShowTransactionPicturesInMobileTransactionEditPage" @toggle:change="alwaysShowTransactionPicturesInMobileTransactionEditPage = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Account List Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item :disabled="!hasAnyVisibleAccount"
                          :title="tt('Accounts Included in Total')"
                          link="/settings/filter/account?type=accountListTotalAmount">
                <template #after>
                    <f7-preloader v-if="loadingAccounts" />
                    <span v-else-if="!loadingAccounts">{{ accountsIncludedInTotalDisplayContent }}</span>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Exchange Rates Data Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                link="#"
                :title="tt('Sort by')"
                :after="findDisplayNameByType(allCurrencySortingTypes, currencySortByInExchangeRatesPage)"
                @click="showCurrencySortByInExchangeRatesPagePopup = true"
            >
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

import { findNameByValue, findDisplayNameByType } from '@/lib/common.ts';

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();
const {
    loadingAccounts,
    hasAnyVisibleAccount,
    allTimezoneTypesUsedForStatistics,
    allCurrencySortingTypes,
    allAutoSaveTransactionDraftTypes,
    showAmountInHomePage,
    timezoneUsedForStatisticsInHomePage,
    showTotalAmountInTransactionListPage,
    showTagInTransactionListPage,
    autoSaveTransactionDraft,
    isAutoGetCurrentGeoLocation,
    currencySortByInExchangeRatesPage,
    accountsIncludedInTotalDisplayContent
} = useAppSettingPageBase();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();

const showTimezoneUsedForStatisticsInHomePagePopup = ref<boolean>(false);
const showAutoSaveTransactionDraftPopup = ref<boolean>(false);
const showCurrencySortByInExchangeRatesPagePopup = ref<boolean>(false);

const alwaysShowTransactionPicturesInMobileTransactionEditPage = computed<boolean>({
    get: () => settingsStore.appSettings.alwaysShowTransactionPicturesInMobileTransactionEditPage,
    set: (value) => settingsStore.setAlwaysShowTransactionPicturesInMobileTransactionEditPage(value)
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
}

init();
</script>
