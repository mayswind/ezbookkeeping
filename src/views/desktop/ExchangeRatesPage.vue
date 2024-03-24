<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <small>{{ $t('Data source') }}</small>
                            <p class="text-body-1 mt-1 mb-3">
                                <a tabindex="-1" target="_blank" :href="exchangeRatesData.referenceUrl" v-if="!loading && exchangeRatesData && exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</a>
                                <span v-else-if="!loading && exchangeRatesData && !exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                                <span v-else-if="!loading && !exchangeRatesData">{{ $t('None') }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <small v-if="exchangeRatesDataUpdateTime || loading">{{ $t('Last Updated') }}</small>
                            <p class="text-body-1 mt-1" v-if="exchangeRatesDataUpdateTime || loading">
                                <span v-if="!loading">{{ exchangeRatesDataUpdateTime }}</span>
                                <span v-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-5" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <small>{{ $t('Base Amount') }}</small>
                            <amount-input class="mt-2" density="compact"
                                          :disabled="loading || !exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length"
                                          v-model="baseAmount"/>
                        </div>
                        <div class="mx-6 mt-4">
                            <small>{{ $t('Base Currency') }}</small>
                        </div>
                        <v-tabs show-arrows class="mb-4" direction="vertical"
                                :disabled="loading" v-model="baseCurrency"
                                v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                            <v-tab class="tab-text-truncate" :key="exchangeRate.currencyCode" :value="exchangeRate.currencyCode"
                                   v-for="exchangeRate in availableExchangeRates">
                                <div class="text-truncate">
                                    <span>{{ exchangeRate.currencyDisplayName }}</span>
                                    <small class="smaller ml-1">{{ exchangeRate.currencyCode }}</small>
                                </div>
                            </v-tab>
                        </v-tabs>
                        <div class="mx-6 mt-2 mb-4"
                             v-else-if="!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length">
                            <span v-if="!loading">{{ $t('None') }}</span>
                            <span v-else-if="loading">
                                <v-skeleton-loader class="skeleton-no-margin pt-2 pb-5" type="text"
                                                   :key="itemIdx" :loading="loading"
                                                   v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]"></v-skeleton-loader>
                            </span>
                        </div>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="exchangeRatesPage">
                                <v-card variant="flat">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ $t('Exchange Rates Data') }}</span>
                                            <v-btn density="compact" color="default" variant="text"
                                                   class="ml-2" :icon="true"
                                                   v-if="!loading" @click="reload">
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                                        </div>
                                    </template>

                                    <v-table class="exchange-rates-table table-striped" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th class="text-uppercase">
                                                <div class="d-flex align-center">
                                                    <span>{{ $t('Currency') }}</span>
                                                    <v-spacer/>
                                                    <span>{{ $t('Amount') }}</span>
                                                </div>
                                            </th>
                                        </tr>
                                        </thead>

                                        <tbody>
                                        <tr :key="itemIdx"
                                            v-for="itemIdx in (loading && (!exchangeRatesData || !exchangeRatesData.exchangeRates || exchangeRatesData.exchangeRates.length < 1) ? [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ] : [])">
                                            <td class="px-0">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>

                                        <tr v-if="!loading && (!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length)">
                                            <td>{{ $t('No exchange rates data') }}</td>
                                        </tr>

                                        <tr class="exchange-rates-table-row-data" :key="exchangeRate.currencyCode"
                                            v-for="exchangeRate in availableExchangeRates">
                                            <td>
                                                <div class="d-flex align-center">
                                                    <span class="text-sm">{{ exchangeRate.currencyDisplayName }}</span>
                                                    <span class="text-caption ml-1">{{ exchangeRate.currencyCode }}</span>
                                                    <v-spacer/>
                                                    <v-btn class="px-2 ml-2 mr-3" color="default"
                                                           density="comfortable" variant="text"
                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                           v-if="exchangeRate.currencyCode !== baseCurrency"
                                                           @click="setAsBaseline(exchangeRate.currencyCode, getConvertedAmount(exchangeRate))">
                                                        {{ $t('Set As Baseline') }}
                                                    </v-btn>
                                                    <span>{{ getDisplayConvertedAmount(exchangeRate, isEnableThousandsSeparator) }}</span>
                                                </div>
                                            </td>
                                        </tr>
                                        </tbody>
                                    </v-table>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script>
import { useDisplay } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { isNumber } from '@/lib/common.js';
import {
    stringCurrencyToNumeric,
    getConvertedAmount,
    getDisplayExchangeRateAmount
} from '@/lib/currency.js';

import {
    mdiRefresh,
    mdiMenu
} from '@mdi/js';

export default {
    data() {
        const { mdAndUp } = useDisplay();
        const userStore = useUserStore();

        return {
            activeTab: 'exchangeRatesPage',
            baseCurrency: userStore.currentUserDefaultCurrency,
            baseAmount: 100,
            loading: true,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            icons: {
                refresh: mdiRefresh,
                menu: mdiMenu
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useExchangeRatesStore),
        isEnableThousandsSeparator() {
            return this.settingsStore.appSettings.thousandsSeparator;
        },
        exchangeRatesData() {
            return this.exchangeRatesStore.latestExchangeRates.data;
        },
        exchangeRatesDataUpdateTime() {
            const exchangeRatesLastUpdateTime = this.exchangeRatesStore.exchangeRatesLastUpdateTime;
            return exchangeRatesLastUpdateTime ? this.$locale.formatUnixTimeToLongDate(this.userStore, exchangeRatesLastUpdateTime) : '';
        },
        availableExchangeRates() {
            return this.$locale.getAllDisplayExchangeRates(this.exchangeRatesData);
        }
    },
    created() {
        this.reload(false);
    },
    setup() {
        const display = useDisplay();

        return {
            display: display
        };
    },
    watch: {
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    methods: {
        reload(force) {
            const self = this;

            self.loading = true;

            self.exchangeRatesStore.getLatestExchangeRates({
                silent: false,
                force: force
            }).then(() => {
                self.loading = false;

                if (this.exchangeRatesData && this.exchangeRatesData.exchangeRates) {
                    let foundDefaultCurrency = false;

                    for (let i = 0; i < this.exchangeRatesData.exchangeRates.length; i++) {
                        const exchangeRate = this.exchangeRatesData.exchangeRates[i];
                        if (exchangeRate.currency === this.baseCurrency) {
                            foundDefaultCurrency = true;
                            break;
                        }
                    }

                    if (force) {
                        self.$refs.snackbar.showMessage('Exchange rates data has been updated');
                    } else if (!foundDefaultCurrency) {
                        this.$refs.snackbar.showMessage('There is no exchange rates data for your default currency');
                    }
                }
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        getConvertedAmount(toExchangeRate) {
            if (!this.baseCurrency) {
                return 0;
            }

            const fromExchangeRate = this.exchangeRatesStore.latestExchangeRateMap[this.baseCurrency];

            try {
                return getConvertedAmount(this.baseAmount / 100, fromExchangeRate, toExchangeRate);
            } catch (e) {
                return 0;
            }
        },
        getDisplayConvertedAmount(toExchangeRate, isEnableThousandsSeparator) {
            const rateStr = this.getConvertedAmount(toExchangeRate).toString();
            return getDisplayExchangeRateAmount(rateStr, isEnableThousandsSeparator);
        },
        setAsBaseline(currency, amount) {
            if (!isNumber(amount)) {
                amount = '';
            }

            this.baseCurrency = currency;
            this.baseAmount = stringCurrencyToNumeric(amount.toString());
        }
    }
}
</script>

<style>
.exchange-rates-table tr.exchange-rates-table-row-data .hover-display {
    display: none;
}

.exchange-rates-table tr.exchange-rates-table-row-data:hover .hover-display {
    display: grid;
}

</style>
