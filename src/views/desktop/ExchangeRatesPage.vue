<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading }">
                <template #title>
                    <span>{{ $t('Exchange Rates Data') }}</span>
                    <v-btn density="compact" color="default" variant="text"
                           class="ml-2" :icon="true"
                           v-if="!loading" @click="update">
                        <v-icon :icon="icons.refresh" size="24" />
                        <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                    </v-btn>
                    <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text>
                    <span class="text-sm">
                        {{ $t('Data source') }}
                        <a tabindex="-1" target="_blank" :href="exchangeRatesData.referenceUrl" v-if="exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</a>
                        <span v-else-if="!exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                        <span v-if="exchangeRatesDataUpdateTime">&nbsp;,&nbsp;{{ $t('Last Updated') }}&nbsp;{{ exchangeRatesDataUpdateTime }}</span>
                    </span>
                </v-card-text>
                <v-card-text v-if="!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length">
                    <span class="text-subtitle-1">{{ $t('No exchange rates data') }}</span>
                </v-card-text>
                <v-card-text v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-subtitle-1">{{ $t('Base Currency') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <v-autocomplete
                                density="compact"
                                single-line
                                item-title="currencyDisplayName"
                                item-value="currencyCode"
                                :disabled="loading"
                                :items="availableExchangeRates"
                                v-model="baseCurrency"
                            >
                                <template v-slot:no-data>
                                    <div class="px-4">{{ $t('No results') }}</div>
                                </template>
                            </v-autocomplete>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-subtitle-1">{{ $t('Base Amount') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <amount-input density="compact" :disabled="loading" v-model="baseAmount"/>
                        </v-col>
                    </v-row>
                </v-card-text>
                <v-table v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                    <thead>
                    <tr>
                        <th class="text-uppercase" style="width: 50%">{{ $t('Currency') }}</th>
                        <th class="text-uppercase">{{ $t('Amount') }}</th>
                    </tr>
                    </thead>

                    <tbody>
                    <tr :key="exchangeRate.currencyCode" v-for="exchangeRate in availableExchangeRates">
                        <td>
                            <span style="margin-right: 5px">{{ exchangeRate.currencyDisplayName }}</span>
                            <small class="smaller">{{ exchangeRate.currencyCode }}</small>
                        </td>
                        <td>{{ getDisplayConvertedAmount(exchangeRate) }}</td>
                    </tr>
                    </tbody>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <snackbar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { appendThousandsSeparator } from '@/lib/common.js';
import { getExchangedAmount } from '@/lib/currency.js';

import {
    mdiRefresh
} from '@mdi/js';

export default {
    data() {
        const userStore = useUserStore();

        return {
            baseCurrency: userStore.currentUserDefaultCurrency,
            baseAmount: '1',
            loading: false,
            icons: {
                refresh: mdiRefresh
            }
        };
    },
    computed: {
        ...mapStores(useUserStore, useExchangeRatesStore),
        exchangeRatesData() {
            return this.exchangeRatesStore.latestExchangeRates.data;
        },
        exchangeRatesDataUpdateTime() {
            const exchangeRatesLastUpdateTime = this.exchangeRatesStore.exchangeRatesLastUpdateTime;
            return exchangeRatesLastUpdateTime ? this.$locale.formatUnixTimeToLongDate(this.userStore, exchangeRatesLastUpdateTime) : '';
        },
        availableExchangeRates() {
            if (!this.exchangeRatesData || !this.exchangeRatesData.exchangeRates) {
                return [];
            }

            const availableExchangeRates = [];

            for (let i = 0; i < this.exchangeRatesData.exchangeRates.length; i++) {
                const exchangeRate = this.exchangeRatesData.exchangeRates[i];

                availableExchangeRates.push({
                    currencyCode: exchangeRate.currency,
                    currencyDisplayName: this.$t(`currency.${exchangeRate.currency}`),
                    rate: exchangeRate.rate
                });
            }

            availableExchangeRates.sort(function(c1, c2) {
                return c1.currencyDisplayName.localeCompare(c2.currencyDisplayName);
            })

            return availableExchangeRates;
        }
    },
    created() {
        if (!this.exchangeRatesData || !this.exchangeRatesData.exchangeRates) {
            return;
        }

        for (let i = 0; i < this.exchangeRatesData.exchangeRates.length; i++) {
            const exchangeRate = this.exchangeRatesData.exchangeRates[i];
            if (exchangeRate.currency === this.baseCurrency) {
                return;
            }
        }

        this.$refs.snackbar.showMessage('There is no exchange rates data for your default currency');
    },
    methods: {
        update() {
            const self = this;

            if (self.loading) {
                return;
            }

            self.loading = true;
            self.exchangeRatesStore.getLatestExchangeRates({
                silent: false,
                force: true
            }).then(() => {
                self.loading = false;
                self.$refs.snackbar.showMessage('Exchange rates data has been updated');
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

            if (!fromExchangeRate) {
                return '';
            }

            if (this.baseAmount === '') {
                return 0;
            }

            try {
                return getExchangedAmount(parseFloat(this.baseAmount), fromExchangeRate.rate, toExchangeRate.rate);
            } catch (e) {
                return 0;
            }
        },
        getDisplayConvertedAmount(toExchangeRate) {
            const rateStr = this.getConvertedAmount(toExchangeRate).toString();

            if (rateStr.indexOf('.') < 0) {
                return appendThousandsSeparator(rateStr);
            } else {
                let firstNonZeroPos = 0;

                for (let i = 0; i < rateStr.length; i++) {
                    if (rateStr.charAt(i) !== '.' && rateStr.charAt(i) !== '0') {
                        firstNonZeroPos = Math.min(i + 4, rateStr.length);
                        break;
                    }
                }

                const trimmedRateStr = rateStr.substring(0, Math.max(6, Math.max(firstNonZeroPos, rateStr.indexOf('.') + 2)));
                return appendThousandsSeparator(trimmedRateStr);
            }
        }
    }
}
</script>
