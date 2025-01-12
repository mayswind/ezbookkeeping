<template>
    <f7-page>
        <f7-navbar :title="tt('Page Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('Overview Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ tt('Show Amount') }}</span>
                <f7-toggle :checked="showAmountInHomePage" @toggle:change="showAmountInHomePage = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item
                :title="tt('Timezone Used for Statistics')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Timezone Type'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="timezoneUsedForStatisticsInHomePage">
                    <option :value="timezoneType.type"
                            :key="timezoneType.type"
                            v-for="timezoneType in allTimezoneTypesUsedForStatistics">{{ timezoneType.displayName }}</option>
                </select>
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
                :title="tt('Automatically Save Draft')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Automatically Save Draft'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="autoSaveTransactionDraft">
                    <option value="disabled">{{ tt('Disabled') }}</option>
                    <option value="enabled">{{ tt('Enabled') }}</option>
                    <option value="confirmation">{{ tt('Show Confirmation Every Time') }}</option>
                </select>
            </f7-list-item>
            <f7-list-item>
                <span>{{ tt('Automatically Add Geolocation') }}</span>
                <f7-toggle :checked="isAutoGetCurrentGeoLocation" @toggle:change="isAutoGetCurrentGeoLocation = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Exchange Rates Data Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :title="tt('Sort by')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Sort by'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="currencySortByInExchangeRatesPage">
                    <option :value="sortingType.type"
                            :key="sortingType.type"
                            v-for="sortingType in allCurrencySortingTypes">{{ sortingType.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { useAppSettingPageBase } from '@/views/base/settings/AppSettingsPageBase.ts';

const { tt } = useI18n();
const {
    allTimezoneTypesUsedForStatistics,
    allCurrencySortingTypes,
    showAmountInHomePage,
    timezoneUsedForStatisticsInHomePage,
    showTotalAmountInTransactionListPage,
    showTagInTransactionListPage,
    autoSaveTransactionDraft,
    isAutoGetCurrentGeoLocation,
    currencySortByInExchangeRatesPage
} = useAppSettingPageBase();
</script>
