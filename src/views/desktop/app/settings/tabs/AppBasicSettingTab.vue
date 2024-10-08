<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="$t('Basic Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Theme')"
                                    :placeholder="$t('Theme')"
                                    :items="[
                                        { value: 'auto', displayName: $t('System Default') },
                                        { value: 'light', displayName: $t('Light') },
                                        { value: 'dark', displayName: $t('Dark') }
                                    ]"
                                    v-model="theme"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-autocomplete
                                    item-title="displayNameWithUtcOffset"
                                    item-value="name"
                                    auto-select-first
                                    persistent-placeholder
                                    :label="$t('Timezone')"
                                    :placeholder="$t('Timezone')"
                                    :items="allTimezones"
                                    :no-data-text="$t('No results')"
                                    v-model="timeZone"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Auto-update Exchange Rates Data')"
                                    :placeholder="$t('Auto-update Exchange Rates Data')"
                                    :items="enableDisableOptions"
                                    v-model="isAutoUpdateExchangeRatesData"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Show Account Balance')"
                                    :placeholder="$t('Show Account Balance')"
                                    :items="enableDisableOptions"
                                    v-model="showAccountBalance"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="$t('Overview Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Show Amount')"
                                    :placeholder="$t('Show Amount')"
                                    :items="enableDisableOptions"
                                    v-model="showAmountInHomePage"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Timezone Used for Statistics')"
                                    :placeholder="$t('Timezone Used for Statistics')"
                                    :items="allTimezoneTypesUsedForStatistics"
                                    v-model="timezoneUsedForStatisticsInHomePage"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="$t('Transaction List Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    persistent-placeholder
                                    :label="$t('Transactions Per Page')"
                                    :placeholder="$t('Transactions Per Page')"
                                    :items="[ 5, 10, 15, 20, 25, 30, 50 ]"
                                    v-model="itemsCountInTransactionListPage"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Show Monthly Total Amount')"
                                    :placeholder="$t('Show Monthly Total Amount')"
                                    :items="enableDisableOptions"
                                    v-model="showTotalAmountInTransactionListPage"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Show Transaction Tag')"
                                    :placeholder="$t('Show Transaction Tag')"
                                    :items="enableDisableOptions"
                                    v-model="showTagInTransactionListPage"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="$t('Transaction Edit Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Automatically Save Draft')"
                                    :placeholder="$t('Automatically Save Draft')"
                                    :items="[
                                        { value: 'disabled', displayName: $t('Disabled') },
                                        { value: 'enabled', displayName: $t('Enabled') },
                                        { value: 'confirmation', displayName: $t('Show Confirmation Every Time') }
                                    ]"
                                    v-model="autoSaveTransactionDraft"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Automatically Add Geolocation')"
                                    :placeholder="$t('Automatically Add Geolocation')"
                                    :items="enableDisableOptions"
                                    v-model="isAutoGetCurrentGeoLocation"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useOverviewStore } from '@/stores/overview.js';
import { useStatisticsStore } from '@/stores/statistics.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { getSystemTheme } from '@/lib/ui.js';

export default {
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTransactionsStore, useOverviewStore, useStatisticsStore, useExchangeRatesStore),
        enableDisableOptions() {
            return this.$locale.getEnableDisableOptions();
        },
        allTimezones() {
            return this.$locale.getAllTimezones(true);
        },
        allTimezoneTypesUsedForStatistics() {
            return this.$locale.getAllTimezoneTypesUsedForStatistics(this.timeZone);
        },
        theme: {
            get: function () {
                return this.settingsStore.appSettings.theme;
            },
            set: function (value) {
                if (value !== this.settingsStore.appSettings.theme) {
                    this.settingsStore.setTheme(value);

                    if (value === 'light' || value === 'dark') {
                        this.globalTheme.global.name.value = value;
                    } else {
                        this.globalTheme.global.name.value = getSystemTheme();
                    }
                }
            }
        },
        timeZone: {
            get: function () {
                return this.settingsStore.appSettings.timeZone;
            },
            set: function (value) {
                this.settingsStore.setTimeZone(value);
                this.$locale.setTimeZone(value);
                this.transactionsStore.updateTransactionListInvalidState(true);
                this.overviewStore.updateTransactionOverviewInvalidState(true);
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
            }
        },
        isAutoUpdateExchangeRatesData: {
            get: function () {
                return this.settingsStore.appSettings.autoUpdateExchangeRatesData;
            },
            set: function (value) {
                this.settingsStore.setAutoUpdateExchangeRatesData(value);
            }
        },
        showAccountBalance: {
            get: function () {
                return this.settingsStore.appSettings.showAccountBalance;
            },
            set: function (value) {
                this.settingsStore.setShowAccountBalance(value);
            }
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
        itemsCountInTransactionListPage: {
            get: function () {
                return this.settingsStore.appSettings.itemsCountInTransactionListPage;
            },
            set: function (value) {
                this.settingsStore.setItemsCountInTransactionListPage(value);
            }
        },
        autoSaveTransactionDraft: {
            get: function () {
                return this.settingsStore.appSettings.autoSaveTransactionDraft;
            },
            set: function (value) {
                this.settingsStore.setAutoSaveTransactionDraft(value);

                if (value === 'disabled') {
                    this.transactionsStore.clearTransactionDraft();
                }
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
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    }
};
</script>
