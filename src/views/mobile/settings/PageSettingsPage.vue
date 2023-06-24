<template>
    <f7-page>
        <f7-navbar :title="$t('Page Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ $t('Overview Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ $t('Show Amount') }}</span>
                <f7-toggle :checked="showAmountInHomePage" @toggle:change="showAmountInHomePage = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ $t('Transaction List Page') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item>
                <span>{{ $t('Show Monthly Total Amount') }}</span>
                <f7-toggle :checked="showTotalAmountInTransactionListPage" @toggle:change="showTotalAmountInTransactionListPage = $event"></f7-toggle>
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

export default {
    computed: {
        ...mapStores(useSettingsStore),
        showAmountInHomePage: {
            get: function () {
                return this.settingsStore.appSettings.showAmountInHomePage;
            },
            set: function (value) {
                this.settingsStore.setShowAmountInHomePage(value);
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
