<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Filter Amount')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :text="$t('Apply')" @click="confirm"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list form strong inset dividers class="margin-vertical">
            <f7-list-item
                class="ebk-small-amount"
                link="#" no-chevron
                :header="amount1Header"
                :title="getDisplayAmount(amount1)"
                @click="showAmount1Sheet = true"
            >
                <number-pad-sheet :min-value="allowedMinAmount"
                                  :max-value="allowedMaxAmount"
                                  v-model:show="showAmount1Sheet"
                                  v-model="amount1"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="ebk-small-amount"
                link="#" no-chevron
                :header="amount2Header"
                :title="getDisplayAmount(amount2)"
                @click="showAmount2Sheet = true"
                v-if="amountCount === 2"
            >
                <number-pad-sheet :min-value="allowedMinAmount"
                                  :max-value="allowedMaxAmount"
                                  v-model:show="showAmount2Sheet"
                                  v-model="amount2"
                ></number-pad-sheet>
            </f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical">
            <f7-list-item :key="filterType.type" :title="$t(filterType.name)"
                          v-for="filterType in allAmountFilterTypes"
                          @click="type = filterType.type">
                <template #after>
                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="type === filterType.type"></f7-icon>
                </template>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTransactionsStore } from '@/stores/transaction.js';

import { AmountFilterType } from '@/core/numeral.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { isString } from '@/lib/common.ts';
import logger from '@/lib/logger.ts';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        return {
            type: '',
            amount1: 0,
            amount2: 0,
            showAmount1Sheet: false,
            showAmount2Sheet: false
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useTransactionsStore),
        allAmountFilterTypes() {
            return AmountFilterType.values();
        },
        allowedMinAmount() {
            return TRANSACTION_MIN_AMOUNT;
        },
        allowedMaxAmount() {
            return TRANSACTION_MAX_AMOUNT;
        },
        amountCount() {
            return this.getAmountFilterParameterCount(this.type);
        },
        title() {
            const amountFilterType = AmountFilterType.valueOf(this.type);
            return amountFilterType ? this.$t(amountFilterType.name) : this.$t('Amount');
        },
        amount1Header() {
            if (this.type === AmountFilterType.GreaterThan.type
                || this.type === AmountFilterType.Between.type
                || this.type === AmountFilterType.NotBetween.type) {
                return this.$t('Minimum Amount');
            } else if (this.type === AmountFilterType.LessThan.type) {
                return this.$t('Maximum Amount');
            } else {
                return this.$t('Amount');
            }
        },
        amount2Header() {
            if (this.type === AmountFilterType.Between.type) {
                return this.$t('Maximum Amount');
            } else if (this.type === AmountFilterType.NotBetween.type) {
                return this.$t('Maximum Amount');
            } else {
                return this.$t('Amount');
            }
        }
    },
    created() {
        const query = this.f7route.query;
        this.type = query.type;

        let amount1 = 0, amount2 = 0;

        if (isString(query.value)) {
            try {
                const filterItems = query.value.split(':');
                const amountCount = this.getAmountFilterParameterCount(filterItems[0]);

                if (filterItems.length === 2 && amountCount === 1) {
                    amount1 = parseInt(filterItems[1]);
                } else if (filterItems.length === 3 && amountCount === 2) {
                    amount1 = parseInt(filterItems[1]);
                    amount2 = parseInt(filterItems[2]);
                }
            } catch (ex) {
                logger.warn('cannot parse amount from filter value, original value is ' + query.value, ex);
            }
        }

        this.amount1 = amount1;
        this.amount2 = amount2;
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        confirm() {
            const router = this.f7router;
            let amountFilter = this.type;

            if (this.amountCount === 1) {
                amountFilter += ':' + this.amount1;
            } else if (this.amountCount === 2) {
                if (this.amount2 < this.amount1) {
                    this.$toast('Incorrect amount range');
                    return;
                }

                amountFilter += ':' + this.amount1 + ':' + this.amount2;
            } else {
                router.back();
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                amountFilter: amountFilter
            });

            if (changed) {
                this.transactionsStore.updateTransactionListInvalidState(true);
            }

            router.back();
        },
        getDisplayAmount(value) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, false);
        },
        getAmountFilterParameterCount(filterType) {
            const amountFilterType = AmountFilterType.valueOf(filterType);
            return amountFilterType ? amountFilterType.paramCount : 0;
        }
    }
}
</script>
