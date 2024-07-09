<template>
    <f7-page>
        <f7-navbar :title="$t('Page Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ $t('Overview Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ $t('Show Amount') }}</span>
                <f7-toggle :checked="showAmountInHomePage" @toggle:change="showAmountInHomePage = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item
                :title="$t('Timezone Used for Statistics')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Timezone Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="timezoneUsedForStatisticsInHomePage">
                    <option :value="timezoneType.type"
                            :key="timezoneType.type"
                            v-for="timezoneType in allTimezoneTypesUsedForStatistics">{{ timezoneType.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ $t('Transaction List Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ $t('Show Monthly Total Amount') }}</span>
                <f7-toggle :checked="showTotalAmountInTransactionListPage" @toggle:change="showTotalAmountInTransactionListPage = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-item>
                <span>{{ $t('Show Transaction Tag') }}</span>
                <f7-toggle :checked="showTagInTransactionListPage" @toggle:change="showTagInTransactionListPage = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ $t('Transaction Edit Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ $t('Automatically Add Geolocation') }}</span>
                <f7-toggle :checked="isAutoGetCurrentGeoLocation" @toggle:change="isAutoGetCurrentGeoLocation = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useOverviewStore } from '@/stores/overview.js';

export default {
    computed: {
        ...mapStores(useSettingsStore, useOverviewStore),
        allTimezoneTypesUsedForStatistics() {
            return this.$locale.getAllTimezoneTypesUsedForStatistics();
        },
        showAmountInHomePage: {
            get: function () {
                return this.settingsStore.appSettings.showAmountInHomePage;
            },
            set: function (value) {
                this.settingsStore.setShowAmountInHomePage(value);
            }
        },
        timezoneUsedForStatisticsInHomePage: {
            get: function () {
                return this.settingsStore.appSettings.timezoneUsedForStatisticsInHomePage;
            },
            set: function (value) {
                this.settingsStore.setTimezoneUsedForStatisticsInHomePage(value);
                this.overviewStore.updateTransactionOverviewInvalidState(true);
            }
        },
        showTotalAmountInTransactionListPage: {
            get: function () {
                return this.settingsStore.appSettings.showTotalAmountInTransactionListPage;
            },
            set: function (value) {
                this.settingsStore.setShowTotalAmountInTransactionListPage(value);
            }
        },
        showTagInTransactionListPage: {
            get: function () {
                return this.settingsStore.appSettings.showTagInTransactionListPage;
            },
            set: function (value) {
                this.settingsStore.setShowTagInTransactionListPage(value);
            }
        },
        isAutoGetCurrentGeoLocation: {
            get: function () {
                return this.settingsStore.appSettings.autoGetCurrentGeoLocation;
            },
            set: function (value) {
                this.settingsStore.setAutoGetCurrentGeoLocation(value);
            }
        }
    }
};
</script>
