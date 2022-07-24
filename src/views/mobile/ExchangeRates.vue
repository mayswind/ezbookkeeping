<template>
    <f7-page ptr @ptr:refresh="update">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Exchange Rates Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false" v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                <f7-list>
                    <f7-list-item
                        class="list-item-with-header-and-title list-item-no-item-after"
                        :header="$t('Base Currency')"
                        smart-select :smart-select-params="{ openIn: 'popup', pageTitle: $t('Base Currency'), searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <f7-block slot="title" class="no-padding no-margin">
                            <span>{{ $t(`currency.${baseCurrency}`) }}&nbsp;</span>
                            <small class="smaller">{{ baseCurrency }}</small>
                        </f7-block>
                        <select v-model="baseCurrency">
                            <option v-for="exchangeRate in availableExchangeRates"
                                    :key="exchangeRate.currencyCode"
                                    :value="exchangeRate.currencyCode">{{ exchangeRate.currencyDisplayName }}</option>
                        </select>
                    </f7-list-item>
                    <f7-list-item
                        class="currency-base-amount"
                        link="#" no-chevron
                        :style="{ fontSize: baseAmountFontSize + 'px' }"
                        :header="$t('Base Amount')"
                        :title="baseAmount | currency"
                        @click="showBaseAmountSheet = true"
                    >
                        <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                          :max-value="$constants.transaction.maxAmount"
                                          :show.sync="showBaseAmountSheet"
                                          v-model="baseAmount"
                        ></number-pad-sheet>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false" v-if="!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length">
                <f7-list>
                    <f7-list-item :title="$t('No exchange rates data')"></f7-list-item>
                </f7-list>
            </f7-card-content>
            <f7-card-content class="no-safe-areas" :padding="false" v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                <f7-list>
                    <f7-list-item v-for="exchangeRate in availableExchangeRates" :key="exchangeRate.currencyCode"
                                  :after="getConvertedAmount(exchangeRate) | exchangeRate"
                                  swipeout>
                        <f7-block slot="title" class="no-padding no-margin">
                            <span style="margin-right: 5px">{{ exchangeRate.currencyDisplayName }}</span>
                            <small class="smaller">{{ exchangeRate.currencyCode }}</small>
                        </f7-block>
                        <f7-swipeout-actions right>
                            <f7-swipeout-button color="primary" close :text="$t('Set As Baseline')" @click="setAsBaseline(exchangeRate.currencyCode, getConvertedAmount(exchangeRate))"></f7-swipeout-button>
                        </f7-swipeout-actions>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
            <f7-card-footer v-if="exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                <span>{{ $t('Last Updated') }}</span>
                <span>{{ exchangeRatesData.updateTime | moment($t('format.date.long')) }}</span>
            </f7-card-footer>
            <f7-card-footer v-if="exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                <span>{{ $t('Data source') }}</span>
                <f7-link external target="_blank" :href="exchangeRatesData.referenceUrl" v-if="exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</f7-link>
                <span v-else-if="!exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
            </f7-card-footer>
        </f7-card>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': updating }" @click="update(null)">
                    <span>{{ $t('Update') }}</span>
                </f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            baseCurrency: self.$store.getters.currentUserDefaultCurrency,
            baseAmount: 100,
            updating: false,
            showMoreActionSheet: false,
            showBaseAmountSheet: false
        };
    },
    computed: {
        exchangeRatesData() {
            return this.$store.state.latestExchangeRates.data;
        },
        exchangeRateMap() {
            const exchangeRateMap = {};

            if (!this.exchangeRatesData.exchangeRates) {
                return exchangeRateMap;
            }

            for (let i = 0; i < this.exchangeRatesData.exchangeRates.length; i++) {
                const exchangeRate = this.exchangeRatesData.exchangeRates[i];
                exchangeRateMap[exchangeRate.currency] = exchangeRate;
            }

            return exchangeRateMap;
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

            availableExchangeRates.sort(function(c1, c2){
                return c1.currencyDisplayName.localeCompare(c2.currencyDisplayName);
            })

            return availableExchangeRates;
        },
        baseAmountFontSize() {
            return this.getFontSizeByAmount(this.baseAmount);
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

        this.$toast('There is no exchange rates data for your default currency');
    },
    methods: {
        update(done) {
            const self = this;

            if (self.updating) {
                if (done) {
                    done();
                }

                return;
            }

            self.updating = true;

            if (!done) {
                self.$showLoading();
            }

            self.$store.dispatch('getLatestExchangeRates', {
                silent: false,
                force: true
            }).then(() => {
                if (done) {
                    done();
                }

                self.updating = false;
                self.$hideLoading();

                self.$toast('Exchange rates data has been updated');
            }).catch(error => {
                if (done) {
                    done();
                }

                self.updating = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        getConvertedAmount(toExchangeRate) {
            const fromExchangeRate = this.exchangeRateMap[this.baseCurrency];

            if (!fromExchangeRate) {
                return '';
            }

            return this.$utilities.getExchangedAmount(this.baseAmount / 100, fromExchangeRate.rate, toExchangeRate.rate);
        },
        setAsBaseline(currency, amount) {
            if (!this.$utilities.isNumber(amount)) {
                amount = '';
            }

            this.baseCurrency = currency;
            this.baseAmount = this.$utilities.stringCurrencyToNumeric(amount.toString());
        },
        getFontSizeByAmount(amount) {
            if (amount >= 100000000 || amount <= -100000000) {
                return 32;
            } else if (amount >= 1000000 || amount <= -1000000) {
                return 36;
            } else {
                return 40;
            }
        }
    }
}
</script>

<style>
.currency-base-amount {
    line-height: 53px;
}

.currency-base-amount .item-header {
    padding-top: calc(var(--f7-typography-padding) / 2);
}
</style>
