<template>
    <f7-page ptr @ptr:refresh="update">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Exchange Rates Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': updating }" :text="$t('Update')" @click="update(null)"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false" v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                <f7-list>
                    <f7-list-item
                        :title="$t('Base Currency')"
                        smart-select :smart-select-params="{ openIn: 'popup', searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }">
                        <select v-model="baseCurrency">
                            <option v-for="exchangeRate in availableExchangeRates"
                                    :key="exchangeRate.currencyCode"
                                    :value="exchangeRate.currencyCode">{{ exchangeRate.currencyDisplayName }}</option>
                        </select>
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
                                  :title="exchangeRate.currencyDisplayName"
                                  :after="exchangeRate.rate | exchangeRate(exchangeRatesData.exchangeRates, baseCurrency)"></f7-list-item>
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
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            baseCurrency: self.$store.getters.currentUserDefaultCurrency || self.$t('default.currency'),
            updating: false
        };
    },
    computed: {
        exchangeRatesData() {
            return this.$store.state.latestExchangeRates.data;
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
        }
    },
    filters: {
        exchangeRate(oldRate, exchangeRates, currentCurrency) {
            const exchangeRateMap = {};

            for (let i = 0; i < exchangeRates.length; i++) {
                const exchangeRate = exchangeRates[i];
                exchangeRateMap[exchangeRate.currency] = exchangeRate;
            }

            const toCurrencyExchangeRate = exchangeRateMap[currentCurrency];

            if (!toCurrencyExchangeRate) {
                return '';
            }

            const newRate = parseFloat(oldRate) / parseFloat(toCurrencyExchangeRate.rate);
            const newRateStr = newRate.toString();

            if (newRateStr.indexOf('.') < 0) {
                return newRateStr;
            } else {
                let firstNonZeroPos = 0;

                for (let i = 0; i < newRateStr.length; i++) {
                    if (newRateStr.charAt(i) !== '.' && newRateStr.charAt(i) !== '0') {
                        firstNonZeroPos = Math.min(i + 4, newRateStr.length);
                        break;
                    }
                }

                return newRateStr.substr(0, Math.max(6, Math.max(firstNonZeroPos, newRateStr.indexOf('.') + 2)));
            }
        }
    }
}
</script>
