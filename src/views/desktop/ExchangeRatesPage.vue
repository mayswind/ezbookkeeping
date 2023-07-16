<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <div class="d-flex flex-column flex-md-row">
                    <div>
                        <div class="mx-6 my-4">
                            <small>{{ $t('Data source') }}</small>
                            <p class="text-subtitle-1 mb-2">
                                <a tabindex="-1" target="_blank" :href="exchangeRatesData.referenceUrl" v-if="!loading && exchangeRatesData && exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</a>
                                <span v-else-if="!loading && exchangeRatesData && !exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                                <span v-else-if="!loading && !exchangeRatesData">{{ $t('None') }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="exchange-rates-summary-skeleton mt-2 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <small v-if="exchangeRatesDataUpdateTime">{{ $t('Last Updated') }}</small>
                            <p class="text-subtitle-1" v-if="exchangeRatesDataUpdateTime">
                                <span v-if="!loading">{{ exchangeRatesDataUpdateTime }}</span>
                                <span v-if="loading">
                                    <v-skeleton-loader class="exchange-rates-summary-skeleton mt-2 mb-6" type="text" :loading="true"></v-skeleton-loader>
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
                            <v-tab :key="exchangeRate.currencyCode" :value="exchangeRate.currencyCode"
                                   v-for="exchangeRate in availableExchangeRates">
                                {{ exchangeRate.currencyDisplayName }}
                                <small class="smaller ml-1">{{ exchangeRate.currencyCode }}</small>
                            </v-tab>
                        </v-tabs>
                        <div class="mx-6 mt-2 mb-4"
                             v-else-if="!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length">
                            <span v-if="!loading">{{ $t('None') }}</span>
                            <v-skeleton-loader type="paragraph" :loading="loading" v-else-if="loading"></v-skeleton-loader>
                        </div>
                    </div>

                    <v-window class="d-flex flex-grow-1 ml-md-5 disable-tab-transition w-100-window-container" v-model="activeTab">
                        <v-window-item value="exchangeRatesPage">
                            <v-card variant="flat">
                                <template #title>
                                    <div class="d-flex align-center">
                                        <span>{{ $t('Exchange Rates Data') }}</span>
                                        <v-btn density="compact" color="default" variant="text"
                                               class="ml-2" :icon="true"
                                               v-if="!loading" @click="reload">
                                            <v-icon :icon="icons.refresh" size="24" />
                                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                        </v-btn>
                                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
                                    </div>
                                </template>

                                <v-table>
                                    <thead>
                                    <tr>
                                        <th class="text-uppercase">{{ $t('Currency') }}</th>
                                        <th class="text-uppercase text-right">{{ $t('Amount') }}</th>
                                    </tr>
                                    </thead>

                                    <tbody>
                                    <tr :key="itemIdx"
                                        v-for="itemIdx in (loading && (!exchangeRatesData || !exchangeRatesData.exchangeRates || exchangeRatesData.exchangeRates.length < 1) ? [ 1, 2, 3 ] : [])">
                                        <td class="px-0" colspan="2">
                                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                        </td>
                                    </tr>

                                    <tr v-if="!loading && (!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length)">
                                        <td colspan="2">{{ $t('No exchange rates data') }}</td>
                                    </tr>

                                    <tr :key="exchangeRate.currencyCode" v-for="exchangeRate in availableExchangeRates">
                                        <td>
                                            <span>{{ exchangeRate.currencyDisplayName }}</span>
                                            <small class="smaller ml-1">{{ exchangeRate.currencyCode }}</small>
                                        </td>
                                        <td class="text-right">{{ getDisplayConvertedAmount(exchangeRate) }}</td>
                                    </tr>
                                    </tbody>
                                </v-table>
                            </v-card>
                        </v-window-item>
                    </v-window>
                </div>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { getConvertedAmount, getDisplayExchangeRateAmount } from '@/lib/currency.js';

import {
    mdiRefresh
} from '@mdi/js';

export default {
    data() {
        const userStore = useUserStore();

        return {
            activeTab: 'exchangeRatesPage',
            baseCurrency: userStore.currentUserDefaultCurrency,
            baseAmount: '1',
            loading: true,
            icons: {
                refresh: mdiRefresh
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
                return getConvertedAmount(parseFloat(this.baseAmount), fromExchangeRate, toExchangeRate);
            } catch (e) {
                return 0;
            }
        },
        getDisplayConvertedAmount(toExchangeRate) {
            if (this.baseAmount === '') {
                return '';
            }

            const rateStr = this.getConvertedAmount(toExchangeRate).toString();
            return getDisplayExchangeRateAmount(rateStr, this.isEnableThousandsSeparator);
        }
    }
}
</script>

<style>
.exchange-rates-summary-skeleton .v-skeleton-loader__text {
    margin: 0;
}
</style>
