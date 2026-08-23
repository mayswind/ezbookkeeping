<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ title || tt('Account Balance List') }}</span>
        </template>

        <v-card-text class="pt-0 px-2">
            <v-list class="py-0" density="compact" v-if="displayAccounts.length">
                <v-list-item class="ps-2 pe-3 mb-h1" :key="account.id" :title="account.name"
                             :to="`/transaction/list?accountIds=${account.id}&dateType=${DateRange.All.type}`"
                             v-for="account in displayAccounts">
                    <template #prepend>
                        <ItemIcon class="me-2" size="24px" :icon-type="getAccountIconType(account.iconType)" :icon-id="account.icon" :color="account.color" />
                    </template>
                    <template #append>
                        <span>{{ accountBalance(account) }}</span>
                    </template>
                </v-list-item>
            </v-list>
            <div v-if="loading && !displayAccounts.length">
                <v-skeleton-loader class="mb-h1" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="text-medium-emphasis text-center pt-4" v-else-if="!loading && !displayAccounts.length">{{ tt('No data') }}</div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountListPageBase } from '@/views/base/accounts/AccountListPageBase.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { DateRange } from '@/core/datetime.ts';
import { AccountType } from '@/core/account.ts';

import type { Account } from '@/models/account.ts';

import { isArray } from '@/lib/common.ts';
import { parseBigDecimal, BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { getAccountIconType } from '@/lib/icon.ts';

const props = defineProps<{
    loading: boolean;
    title?: string;
    accountCategories: number[];
    itemCount: number;
    sortBy: string;
    editing?: boolean
}>();

const { tt } = useI18n();

const userStore = useUserStore();
const accountsStore = useAccountsStore();
const exchangeRatesStore = useExchangeRatesStore();

const { accountBalance } = useAccountListPageBase();

const displayAccounts = computed<Account[]>(() => {
    const accounts: Account[] = accountsStore.allAccounts.filter(account => {
        if (account.hidden) {
            return false;
        }

        return isArray(props.accountCategories) && (props.accountCategories.includes(0) || props.accountCategories.includes(account.category));
    });

    if (props.sortBy === 'balance') {
        const accountBalances: Record<string, BigDecimal> = {};

        for (const account of accounts) {
            let totalBalance: BigDecimal = BIG_DECIMAL_ZERO;

            if (account.type === AccountType.SingleAccount.type) {
                totalBalance = parseBigDecimal(account.balance);

                if (account.currency !== userStore.currentUserDefaultCurrency) {
                    const exchangedBalance = exchangeRatesStore.getExchangedAmount(totalBalance, account.currency, userStore.currentUserDefaultCurrency);

                    if (exchangedBalance) {
                        totalBalance = exchangedBalance.truncate();
                    }
                }
            } else if (account.type === AccountType.MultiSubAccounts.type && account.subAccounts) {
                for (const subAccount of account.subAccounts) {
                    let subAccountBalance = parseBigDecimal(subAccount.balance);

                    if (subAccount.currency !== userStore.currentUserDefaultCurrency) {
                        const exchangedBalance = exchangeRatesStore.getExchangedAmount(subAccountBalance, subAccount.currency, userStore.currentUserDefaultCurrency);

                        if (exchangedBalance) {
                            subAccountBalance = exchangedBalance.truncate();
                        }
                    }

                    totalBalance = totalBalance.add(subAccountBalance);
                }
            }

            accountBalances[account.id] = totalBalance;
        }

        accounts.sort((a, b) => {
            const balanceA = accountBalances[a.id] ?? BIG_DECIMAL_ZERO;
            const balanceB = accountBalances[b.id] ?? BIG_DECIMAL_ZERO;
            return balanceB.compareTo(balanceA);
        });
    }

    return accounts.slice(0, props.itemCount);
});
</script>
