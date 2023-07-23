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
                                    :label="$t('Enable Thousands Separator')"
                                    :placeholder="$t('Enable Thousands Separator')"
                                    :items="enableDisableOptions"
                                    v-model="isEnableThousandsSeparator"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Currency Display Mode')"
                                    :placeholder="$t('Currency Display Mode')"
                                    :items="[
                                        { value: allCurrencyDisplayModes.None, displayName: $t('None') },
                                        { value: allCurrencyDisplayModes.Symbol, displayName: $t('Currency Symbol') },
                                        { value: allCurrencyDisplayModes.Code, displayName: $t('Currency Code') },
                                        { value: allCurrencyDisplayModes.Name, displayName: $t('Currency Name') }
                                    ]"
                                    v-model="currencyDisplayMode"
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
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="$t('Show Monthly Total Amount')"
                                    :placeholder="$t('Show Monthly Total Amount')"
                                    :items="enableDisableOptions"
                                    v-model="showTotalAmountInTransactionListPage"
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

import currencyConstants from '@/consts/currency.js';
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
        allCurrencyDisplayModes() {
            return currencyConstants.allCurrencyDisplayModes;
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
        isEnableThousandsSeparator: {
            get: function () {
                return this.settingsStore.appSettings.thousandsSeparator;
            },
            set: function (value) {
                this.settingsStore.setEnableThousandsSeparator(value);
            }
        },
        currencyDisplayMode: {
            get: function () {
                return this.settingsStore.appSettings.currencyDisplayMode;
            },
            set: function (value) {
                this.settingsStore.setCurrencyDisplayMode(value);
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
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    }
};
</script>
