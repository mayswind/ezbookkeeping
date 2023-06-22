<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <span>{{ $t('Exchange Rates Data') }}</span>
                    <v-btn density="compact" color="default" variant="text"
                           class="ml-2" :class="{ 'disabled': updating }"
                           :icon="true" @click="update">
                        <v-icon :icon="icons.refresh" size="24" />
                    </v-btn>
                </template>

                <v-card-text>
                    <span class="text-sm">
                        {{ $t('Data source') }}
                        <a target="_blank" :href="exchangeRatesData.referenceUrl" v-if="exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</a>
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
                            <v-select
                                density="compact"
                                single-line
                                item-title="currencyDisplayName"
                                item-value="currencyCode"
                                :items="availableExchangeRates"
                                v-model="baseCurrency"
                            ></v-select>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-subtitle-1">{{ $t('Base Amount') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <amount-input density="compact" v-model="baseAmount"/>
                        </v-col>
                    </v-row>
                </v-card-text>
                <v-table v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                    <thead>
                    <tr>
                        <th style="width: 50%">{{ $t('Currency') }}</th>
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

    <v-snackbar v-model="showSnackbar">
        {{ snackbarMessage }}

        <template #actions>
            <v-btn color="primary" variant="text" @click="showSnackbar = false">{{ $t('Close') }}</v-btn>
        </template>
    </v-snackbar>

    <v-overlay class="justify-center align-center" :persistent="true" v-model="updating">
        <v-progress-circular indeterminate></v-progress-circular>
    </v-overlay>
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
            updating: false,
            showSnackbar: false,
            snackbarMessage: '',
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

        this.showSnackbarMessage(this.$t('There is no exchange rates data for your default currency'));
    },
    methods: {
        update() {
            const self = this;

            if (self.updating) {
                return;
            }

            self.updating = true;
            self.exchangeRatesStore.getLatestExchangeRates({
                silent: false,
                force: true
            }).then(() => {
                self.updating = false;
                self.showSnackbarMessage(self.$t('Exchange rates data has been updated'));
            }).catch(error => {
                self.updating = false;

                if (!error.processed) {
                    self.showSnackbarMessage(self.$tError(error.message || error));
                }
            });
        },
        getConvertedAmount(toExchangeRate) {
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
        },
        showSnackbarMessage(message) {
            this.showSnackbar = true;
            this.snackbarMessage = message;
        }
    }
}
</script>
