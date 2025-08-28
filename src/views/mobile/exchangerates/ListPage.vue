<template>
    <f7-page ptr @ptr:refresh="reload">
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
                                           :filter-placeholder="tt('Currency')"
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
                          :id="getExchangeRateDomId(exchangeRate)"
                          :after="getFinalConvertedAmount(exchangeRate, true)"
                          :key="baseCurrencyChangedTime + '_' + exchangeRate.currencyCode" v-for="exchangeRate in availableExchangeRates"
                          @swipeout:closed="onExchangeRateSwipeoutClosed()">
                <template #title>
                    <div class="no-padding no-margin">
                        <span style="margin-inline-end: 5px">{{ exchangeRate.currencyDisplayName }}</span>
                        <small class="smaller">{{ exchangeRate.currencyCode }}</small>
                    </div>
                </template>
                <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                     :right="textDirection === TextDirection.LTR"
                                     v-if="exchangeRate.currencyCode !== baseCurrency || (exchangeRate.currencyCode !== defaultCurrency && isUserCustomExchangeRates)">
                    <f7-swipeout-button color="primary" close
                                        :text="tt('Set as Base')"
                                        :class="{ 'disabled': exchangeRate.currencyCode === baseCurrency }"
                                        @click="setAsBaseline(exchangeRate.currencyCode, getFinalConvertedAmount(exchangeRate, false)); settingBaseLine = true"
                                        v-if="settingBaseLine || exchangeRate.currencyCode !== baseCurrency"></f7-swipeout-button>
                    <f7-swipeout-button color="red" class="padding-horizontal"
                                        @click="remove(exchangeRate, false)"
                                        v-if="exchangeRate.currencyCode !== defaultCurrency && isUserCustomExchangeRates">
                        <f7-icon f7="trash"></f7-icon>
                    </f7-swipeout-button>
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
                    <f7-link external target="_blank" :href="exchangeRatesData.referenceUrl" v-if="!isUserCustomExchangeRates && exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</f7-link>
                    <span v-else-if="!isUserCustomExchangeRates && !exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                    <span v-else-if="isUserCustomExchangeRates">{{ tt('User Custom') }}</span>
                </small>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group v-if="isUserCustomExchangeRates">
                <f7-actions-button :class="{ 'disabled': loading }" @click="update()">
                    <span>{{ tt('Update') }}</span>
                </f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading }" @click="reload(undefined)">
                    <span>{{ tt('Refresh') }}</span>
                </f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this user custom exchange rate?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(customExchangeRateToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted } from '@/lib/ui/mobile.ts';
import { useExchangeRatesPageBase } from '@/views/base/ExchangeRatesPageBase.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { TextDirection } from '@/core/text.ts';
import { NumeralSystem } from '@/core/numeral.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';

import type { LocalizedLatestExchangeRate } from '@/models/exchange_rate.ts';

import {
    getCurrentUnixTime
} from '@/lib/datetime.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const {
    tt,
    getCurrentLanguageTextDirection,
    getCurrentNumeralSystemType,
    getCurrencyName,
    formatAmountToLocalizedNumerals,
    formatExchangeRateAmountToWesternArabicNumerals
} = useI18n();

const { showAlert, showToast } = useI18nUIComponents();

const {
    baseCurrency,
    baseAmount,
    defaultCurrency,
    exchangeRatesData,
    isUserCustomExchangeRates,
    exchangeRatesDataUpdateTime,
    availableExchangeRates,
    getConvertedAmount,
    setAsBaseline
} = useExchangeRatesPageBase();

const exchangeRatesStore = useExchangeRatesStore();

const loading = ref<boolean>(false);
const baseCurrencyChangedTime = ref<number>(getCurrentUnixTime());
const settingBaseLine = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showBaseCurrencyPopup = ref<boolean>(false);
const showBaseAmountSheet = ref<boolean>(false);
const customExchangeRateToDelete = ref<LocalizedLatestExchangeRate | null>(null);
const showDeleteActionSheet = ref<boolean>(false);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const displayBaseAmount = computed<string>(() => formatAmountToLocalizedNumerals(baseAmount.value, baseCurrency.value));
const baseAmountFontSizeClass = computed<string>(() => {
    if (baseAmount.value >= 100000000 || baseAmount.value <= -100000000) {
        return 'ebk-small-amount';
    } else if (baseAmount.value >= 1000000 || baseAmount.value <= -1000000) {
        return 'ebk-normal-amount';
    } else {
        return 'ebk-large-amount';
    }
});

function getExchangeRateDomId(exchangeRate: LocalizedLatestExchangeRate): string {
    return 'exchangeRate_' + exchangeRate.currencyCode;
}

function reload(done?: () => void): void {
    if (loading.value) {
        done?.();
        return;
    }

    loading.value = true;

    if (!done) {
        showLoading();
    }

    exchangeRatesStore.getLatestExchangeRates({
        silent: false,
        force: true
    }).then(() => {
        done?.();

        loading.value = false;
        hideLoading();

        showToast('Exchange rates data has been updated');
    }).catch(error => {
        done?.();

        loading.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function update(): void {
    props.f7router.navigate('/exchange_rates/update');
}

function remove(customExchangeRate: LocalizedLatestExchangeRate | null, confirm: boolean): void {
    if (!customExchangeRate) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        customExchangeRateToDelete.value = customExchangeRate;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    customExchangeRateToDelete.value = null;
    showLoading();

    exchangeRatesStore.deleteUserCustomExchangeRate({
        currency: customExchangeRate.currencyCode,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getExchangeRateDomId(customExchangeRate), done);
        }
    }).then(() => {
        if (customExchangeRate.currencyCode === baseCurrency.value) {
            baseCurrency.value = defaultCurrency.value;
        }

        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function getFinalConvertedAmount(toExchangeRate: LocalizedLatestExchangeRate, displayLocalizedDigits: boolean): string {
    const fromExchangeRate = exchangeRatesStore.latestExchangeRateMap[baseCurrency.value];
    const exchangeRateAmount = getConvertedAmount(baseAmount.value / 100, fromExchangeRate, toExchangeRate);

    if (!exchangeRateAmount) {
        if (displayLocalizedDigits) {
            return numeralSystem.value.digitZero;
        } else {
            return NumeralSystem.WesternArabicNumerals.digitZero;
        }
    }

    let ret = formatExchangeRateAmountToWesternArabicNumerals(exchangeRateAmount);

    if (displayLocalizedDigits) {
        ret = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(ret);
    }

    return ret;
}

function onExchangeRateSwipeoutClosed(): void {
    baseCurrencyChangedTime.value = getCurrentUnixTime();
    settingBaseLine.value = false;
}

exchangeRatesStore.getLatestExchangeRates({
    silent: true,
    force: false
}).then(() => {
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
});
</script>

<style>
.currency-base-amount {
    line-height: 53px;
}

.currency-base-amount .item-header {
    padding-top: calc(var(--f7-typography-padding) / 2);
}
</style>
