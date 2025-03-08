<template>
    <f7-page ptr @ptr:refresh="update">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Exchange Rates Data')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical" v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :header="tt('Base Currency')"
                @click="showBaseCurrencyPopup = true"
            >
                <template #title>
                    <div class="no-padding no-margin">
                        <span>{{ getCurrencyName(baseCurrency) }}&nbsp;</span>
                        <small class="smaller">{{ baseCurrency }}</small>
                    </div>
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="currencyCode" value-field="currencyCode"
                                           title-field="currencyDisplayName" after-field="currencyCode"
                                           :title="tt('Base Currency')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Currency Name')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="availableExchangeRates"
                                           v-model:show="showBaseCurrencyPopup"
                                           v-model="baseCurrency">
                </list-item-selection-popup>
            </f7-list-item>
            <f7-list-item
                class="currency-base-amount"
                link="#" no-chevron
                :class="baseAmountFontSizeClass"
                :header="tt('Base Amount')"
                :title="displayBaseAmount"
                @click="showBaseAmountSheet = true"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  :currency="baseCurrency"
                                  v-model:show="showBaseAmountSheet"
                                  v-model="baseAmount"
                ></number-pad-sheet>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length">
            <f7-list-item :title="tt('No exchange rates data')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
            <f7-list-item swipeout
                          :after="getFinalConvertedAmount(exchangeRate)"
                          :key="exchangeRate.currencyCode" v-for="exchangeRate in availableExchangeRates">
                <template #title>
                    <div class="no-padding no-margin">
                        <span style="margin-right: 5px">{{ exchangeRate.currencyDisplayName }}</span>
                        <small class="smaller">{{ exchangeRate.currencyCode }}</small>
                    </div>
                </template>
                <f7-swipeout-actions right v-if="exchangeRate.currencyCode !== baseCurrency">
                    <f7-swipeout-button color="primary" close :text="tt('Set as Base')" @click="setAsBaseline(exchangeRate.currencyCode, getFinalConvertedAmount(exchangeRate))"></f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
            <f7-list-item v-if="exchangeRatesDataUpdateTime">
                <small>{{ tt('Last Updated') }}</small>
                <small>{{ exchangeRatesDataUpdateTime }}</small>
            </f7-list-item>
            <f7-list-item>
                <small>{{ tt('Data source') }}</small>
                <small>
                    <f7-link external target="_blank" :href="exchangeRatesData.referenceUrl" v-if="exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</f7-link>
                    <span v-else-if="!exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                </small>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': updating }" @click="update(undefined)">
                    <span>{{ tt('Update') }}</span>
                </f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useExchangeRatesPageBase } from '@/views/base/ExchangeRatesPageBase.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';

import type { LocalizedLatestExchangeRate } from '@/models/exchange_rate.ts';

const { tt, getCurrencyName, formatAmount, formatExchangeRateAmount } = useI18n();
const { showToast } = useI18nUIComponents();
const { baseCurrency, baseAmount, exchangeRatesData, exchangeRatesDataUpdateTime, availableExchangeRates, getConvertedAmount, setAsBaseline } = useExchangeRatesPageBase();

const exchangeRatesStore = useExchangeRatesStore();

const updating = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showBaseCurrencyPopup = ref<boolean>(false);
const showBaseAmountSheet = ref<boolean>(false);

const displayBaseAmount = computed<string>(() => formatAmount(baseAmount.value, baseCurrency.value));
const baseAmountFontSizeClass = computed<string>(() => {
    if (baseAmount.value >= 100000000 || baseAmount.value <= -100000000) {
        return 'ebk-small-amount';
    } else if (baseAmount.value >= 1000000 || baseAmount.value <= -1000000) {
        return 'ebk-normal-amount';
    } else {
        return 'ebk-large-amount';
    }
});

function update(done?: () => void): void {
    if (updating.value) {
        done?.();
        return;
    }

    updating.value = true;

    if (!done) {
        showLoading();
    }

    exchangeRatesStore.getLatestExchangeRates({
        silent: false,
        force: true
    }).then(() => {
        done?.();

        updating.value = false;
        hideLoading();

        showToast('Exchange rates data has been updated');
    }).catch(error => {
        done?.();

        updating.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function getFinalConvertedAmount(toExchangeRate: LocalizedLatestExchangeRate): string {
    const fromExchangeRate = exchangeRatesStore.latestExchangeRateMap[baseCurrency.value];
    const exchangeRateAmount = getConvertedAmount(baseAmount.value / 100, fromExchangeRate, toExchangeRate);

    if (!exchangeRateAmount) {
        return '0';
    }

    return formatExchangeRateAmount(exchangeRateAmount);
}

if (exchangeRatesData.value && exchangeRatesData.value.exchangeRates) {
    const exchangeRates = exchangeRatesData.value.exchangeRates;
    let hasBaseCurrency = false;

    for (let i = 0; i < exchangeRates.length; i++) {
        const exchangeRate = exchangeRates[i];
        if (exchangeRate.currency === baseCurrency.value) {
            hasBaseCurrency = true;
            break;
        }
    }

    if (!hasBaseCurrency) {
        showToast('There is no exchange rates data for your default currency');
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
