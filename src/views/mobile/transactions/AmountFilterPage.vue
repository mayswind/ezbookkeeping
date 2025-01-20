<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Filter Amount')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :text="tt('Apply')" @click="confirm"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list form strong inset dividers class="margin-vertical">
            <f7-list-item
                class="ebk-small-amount"
                link="#" no-chevron
                :header="amount1Header"
                :title="formatAmountWithCurrency(amount1)"
                @click="showAmount1Sheet = true"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  v-model:show="showAmount1Sheet"
                                  v-model="amount1"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="ebk-small-amount"
                link="#" no-chevron
                :header="amount2Header"
                :title="formatAmountWithCurrency(amount2)"
                @click="showAmount2Sheet = true"
                v-if="amountCount === 2"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  v-model:show="showAmount2Sheet"
                                  v-model="amount2"
                ></number-pad-sheet>
            </f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical">
            <f7-list-item :key="filterType.type" :title="tt(filterType.name)"
                          v-for="filterType in AmountFilterType.values()"
                          @click="type = filterType.type">
                <template #after>
                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="type === filterType.type"></f7-icon>
                </template>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';

import { useTransactionsStore } from '@/stores/transaction.js';

import { AmountFilterType } from '@/core/numeral.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { isString } from '@/lib/common.ts';
import logger from '@/lib/logger.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const type = ref<string>('');
const amount1 = ref<number>(0);
const amount2 = ref<number>(0);
const showAmount1Sheet = ref<boolean>(false);
const showAmount2Sheet = ref<boolean>(false);

const { tt, formatAmountWithCurrency } = useI18n();
const { showToast } = useI18nUIComponents();

const transactionsStore = useTransactionsStore();

const amountCount = computed<number>(() => getAmountFilterParameterCount(type.value));

const amount1Header = computed<string>(() => {
    if (type.value === AmountFilterType.GreaterThan.type
        || type.value === AmountFilterType.Between.type
        || type.value === AmountFilterType.NotBetween.type) {
        return tt('Minimum Amount');
    } else if (type.value === AmountFilterType.LessThan.type) {
        return tt('Maximum Amount');
    } else {
        return tt('Amount');
    }
});

const amount2Header = computed<string>(() => {
    if (type.value === AmountFilterType.Between.type) {
        return tt('Maximum Amount');
    } else if (type.value === AmountFilterType.NotBetween.type) {
        return tt('Maximum Amount');
    } else {
        return tt('Amount');
    }
});

function getAmountFilterParameterCount(filterType: string): number {
    const amountFilterType = AmountFilterType.valueOf(filterType);
    return amountFilterType ? amountFilterType.paramCount : 0;
}

function init(): void {
    const query = props.f7route.query;
    type.value = query['type'] || '';

    let queryAmount1 = 0, queryAmount2 = 0;

    if (isString(query['value'])) {
        try {
            const filterItems = query['value'].split(':');
            const amountCount = getAmountFilterParameterCount(filterItems[0]);

            if (filterItems.length === 2 && amountCount === 1) {
                queryAmount1 = parseInt(filterItems[1]);
            } else if (filterItems.length === 3 && amountCount === 2) {
                queryAmount1 = parseInt(filterItems[1]);
                queryAmount2 = parseInt(filterItems[2]);
            }
        } catch (ex) {
            logger.warn('cannot parse amount from filter value, original value is ' + query['value'], ex);
        }
    }

    amount1.value = queryAmount1;
    amount2.value = queryAmount2;
}

function confirm(): void {
    const router = props.f7router;
    let amountFilter = type.value;

    if (amountCount.value === 1) {
        amountFilter += ':' + amount1.value;
    } else if (amountCount.value === 2) {
        if (amount2.value < amount1.value) {
            showToast('Incorrect amount range');
            return;
        }

        amountFilter += ':' + amount1.value + ':' + amount2.value;
    } else {
        router.back();
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        amountFilter: amountFilter
    });

    if (changed) {
        transactionsStore.updateTransactionListInvalidState(true);
    }

    router.back();
}

init();
</script>
