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
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTransactionsStore } from '@/stores/transaction.js';

import numeralConstants from '@/consts/numeral.js';
import transactionConstants from '@/consts/transaction.js';
import { isString } from '@/lib/common.js';
import logger from '@/lib/logger.js';

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
            return numeralConstants.allAmountFilterTypeArray;
        },
        allowedMinAmount() {
            return transactionConstants.minAmountNumber;
        },
        allowedMaxAmount() {
            return transactionConstants.maxAmountNumber;
        },
        amountCount() {
            return this.getAmountFilterParameterCount(this.type);
        },
        title() {
            const amountFilterType = numeralConstants.allAmountFilterTypeMap[this.type];
            return amountFilterType ? this.$t(amountFilterType.name) : this.$t('Amount');
        },
        amount1Header() {
            if (this.type === numeralConstants.allAmountFilterType.GreaterThan.type
                || this.type === numeralConstants.allAmountFilterType.Between.type
                || this.type === numeralConstants.allAmountFilterType.NotBetween.type) {
                return this.$t('Minimum Amount');
            } else if (this.type === numeralConstants.allAmountFilterType.LessThan.type) {
                return this.$t('Maximum Amount');
            } else {
                return this.$t('Amount');
            }
        },
        amount2Header() {
            if (this.type === numeralConstants.allAmountFilterType.Between.type) {
                return this.$t('Maximum Amount');
            } else if (this.type === numeralConstants.allAmountFilterType.NotBetween.type) {
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
                logger.warn('cannot parse amount from filter value, original value is ' + query.value);
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

            this.transactionsStore.updateTransactionListFilter({
                amountFilter: amountFilter
            });
            this.transactionsStore.updateTransactionListInvalidState(true);
            router.back();
        },
        getDisplayAmount(value) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, false);
        },
        getAmountFilterParameterCount(filterType) {
            const amountFilterType = numeralConstants.allAmountFilterTypeMap[filterType];
            return amountFilterType ? amountFilterType.paramCount : 0;
        }
    }
}
</script>
