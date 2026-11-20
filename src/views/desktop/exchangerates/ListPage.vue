<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Exchange Rates Data') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                               :disabled="loading" @click="update"
                               v-if="isUserCustomExchangeRates">{{ tt('Update') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" class="ms-2"
                               :aria-label="tt('Refresh')" :icon="true" :loading="loading" @click="reload(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <div class="ms-2">
                            <currency-select auto-size density="compact" variant="outlined"
                                             :disabled="loading || !availableExchangeRates.length"
                                             :placeholder="tt('Base Currency')"
                                             :allowed-currencies="allCurrencyCodes"
                                             v-model="baseCurrency"/>
                            <v-tooltip activator="parent">{{ tt('Base Currency') }}</v-tooltip>
                        </div>
                        <div class="ms-3" style="width: 170px">
                            <amount-input density="compact"
                                          :placeholder="tt('Base Amount')"
                                          :currency="baseCurrency"
                                          :disabled="loading || !exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length"
                                          v-model="baseAmount"/>
                            <v-tooltip activator="parent">{{ tt('Base Amount') }}</v-tooltip>
                        </div>
                        <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                               :aria-label="tt('More')" :disabled="loading || !exchangeRatesData" :icon="true">
                            <v-icon :icon="mdiDotsVertical" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-subheader class="text-body-small" :title="tt('Data source')" />
                                    <v-list-item target="_blank" :href="exchangeRatesData.referenceUrl"
                                                 :prepend-icon="mdiOpenInNew" :title="exchangeRatesData.dataSource"
                                                 v-if="exchangeRatesData && !isUserCustomExchangeRates && exchangeRatesData.referenceUrl"></v-list-item>
                                    <v-list-item :title="exchangeRatesData.dataSource" v-else-if="exchangeRatesData && !isUserCustomExchangeRates && !exchangeRatesData.referenceUrl"></v-list-item>
                                    <v-list-item :title="tt('User Custom')" v-else-if="exchangeRatesData && isUserCustomExchangeRates"></v-list-item>
                                    <v-list-item :title="tt('None')" v-else-if="!exchangeRatesData"></v-list-item>
                                    <v-divider class="my-2" v-if="exchangeRatesDataUpdateTime" />
                                    <v-list-subheader class="text-body-small" :title="tt('Last Updated')" v-if="exchangeRatesDataUpdateTime" />
                                    <v-list-item :title="exchangeRatesDataUpdateTime" v-if="exchangeRatesDataUpdateTime"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </template>

                <v-table class="exchange-rates-table table-striped" :hover="!loading">
                    <thead>
                    <tr>
                        <th>
                            <div class="d-flex align-center">
                                <span>{{ tt('Currency') }}</span>
                                <v-spacer/>
                                <span>{{ tt('Amount') }}</span>
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
                        <td>{{ tt('No exchange rates data') }}</td>
                    </tr>

                    <tr class="exchange-rates-table-row-data" :key="exchangeRate.currencyCode"
                        v-for="exchangeRate in availableExchangeRates"
                        @mouseenter="hoveredCurrency = exchangeRate.currencyCode" @mouseleave="hoveredCurrency = ''">
                        <td>
                            <div class="d-flex align-center">
                                <span>{{ exchangeRate.currencyDisplayName }}</span>
                                <span class="text-body-small ms-1">{{ exchangeRate.currencyCode }}</span>

                                <v-spacer/>

                                <template v-if="hoveredCurrency === exchangeRate.currencyCode && !loading">
                                    <v-btn class="px-2 ms-2" color="default"
                                           density="comfortable" variant="text"
                                           v-if="exchangeRate.currencyCode !== baseCurrency"
                                           @click="setAsBaseline(exchangeRate.currencyCode, getFinalConvertedAmount(exchangeRate, false))">
                                        {{ tt('Set as Base') }}
                                    </v-btn>
                                    <v-btn class="px-2" color="default"
                                           density="comfortable" variant="text"
                                           :prepend-icon="mdiDeleteOutline"
                                           :loading="customExchangeRateRemoving[exchangeRate.currencyCode]"
                                           :disabled="loading || updating"
                                           v-if="exchangeRate.currencyCode !== defaultCurrency && isUserCustomExchangeRates"
                                           @click="remove(exchangeRate.currencyCode)">
                                        <template #loader>
                                            <v-progress-circular indeterminate size="20" width="2"/>
                                        </template>
                                        {{ tt('Delete') }}
                                    </v-btn>
                                </template>

                                <span class="text-body-large ms-3">{{ getFinalConvertedAmount(exchangeRate, true) }}</span>
                            </div>
                        </td>
                    </tr>
                    </tbody>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <update-dialog ref="updateDialog" />

    <confirm-dialog ref="confirmDialog"/>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import UpdateDialog from './list/dialogs/UpdateDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useExchangeRatesPageBase } from '@/views/base/ExchangeRatesPageBase.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import { type BigDecimal, NumeralSystem } from '@/core/numeral.ts';
import { AMOUNT_FACTOR } from '@/consts/numeral.ts';

import type { LocalizedLatestExchangeRate } from '@/models/exchange_rate.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';
import logger from '@/lib/logger.ts';

import {
    mdiRefresh,
    mdiDeleteOutline,
    mdiDotsVertical,
    mdiOpenInNew
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type UpdateDialogType = InstanceType<typeof UpdateDialog>;

const { tt, getCurrentNumeralSystemType, formatExchangeRateAmountToWesternArabicNumerals } = useI18n();
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
const desktopPageStore = useDesktopPageStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const updateDialog = useTemplateRef<UpdateDialogType>('updateDialog');

const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const hoveredCurrency = ref<string>('');
const customExchangeRateRemoving = ref<Record<string, boolean>>({});

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const allCurrencyCodes = computed<string[]>(() => availableExchangeRates.value.map(exchangeRate => exchangeRate.currencyCode));

function reload(force: boolean): void {
    loading.value = true;

    exchangeRatesStore.getLatestExchangeRates({
        silent: false,
        force: force
    }).then(() => {
        loading.value = false;

        if (exchangeRatesData.value && exchangeRatesData.value.exchangeRates) {
            const exchangeRates = exchangeRatesData.value.exchangeRates;
            let foundDefaultCurrency = false;

            for (const exchangeRate of exchangeRates) {
                if (exchangeRate.currency === baseCurrency.value) {
                    foundDefaultCurrency = true;
                    break;
                }
            }

            if (force) {
                desktopPageStore.addToast('Exchange rates data has been updated');
            } else if (!foundDefaultCurrency) {
                desktopPageStore.addToast('There is no exchange rates data for your default currency');
            }
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function update(): void {
    updateDialog.value?.open().then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
        }
    }).catch(error => {
        if (error) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function remove(currency: string): void {
    confirmDialog.value?.open('Are you sure you want to delete this user custom exchange rate?').then(() => {
        updating.value = true;
        customExchangeRateRemoving.value[currency] = true;

        exchangeRatesStore.deleteUserCustomExchangeRate({
            currency: currency
        }).then(() => {
            if (currency === baseCurrency.value) {
                baseCurrency.value = defaultCurrency.value;
            }

            updating.value = false;
            customExchangeRateRemoving.value[currency] = false;
        }).catch(error => {
            updating.value = false;
            customExchangeRateRemoving.value[currency] = false;

            if (!error.processed) {
                desktopPageStore.addErrorToast(error);
            }
        });
    });
}

function getFinalConvertedAmount(toExchangeRate: LocalizedLatestExchangeRate, displayLocalizedDigits: boolean): string {
    if (!baseCurrency.value) {
        if (displayLocalizedDigits) {
            return numeralSystem.value.digitZero;
        } else {
            return NumeralSystem.WesternArabicNumerals.digitZero;
        }
    }

    const fromExchangeRate = exchangeRatesStore.latestExchangeRateMap[baseCurrency.value];
    let exchangeRateAmount: BigDecimal | '' | null = BIG_DECIMAL_ZERO;

    try {
        exchangeRateAmount = getConvertedAmount(parseBigDecimal(baseAmount.value).divide(AMOUNT_FACTOR), fromExchangeRate, toExchangeRate);
    } catch (ex) {
        exchangeRateAmount = BIG_DECIMAL_ZERO;
        logger.warn('failed to convert amount by exchange rates, original base amount is ' + baseAmount.value, ex)
    }

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

reload(false);
</script>
