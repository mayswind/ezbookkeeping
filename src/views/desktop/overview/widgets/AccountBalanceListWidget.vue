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

import { useAccountsStore } from '@/stores/account.ts';

import { DateRange } from '@/core/datetime.ts';
import type { Account } from '@/models/account.ts';

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

const { accountBalance } = useAccountListPageBase();

const accountsStore = useAccountsStore();

const displayAccounts = computed<Account[]>(() => accountsStore.getSortedAccounts(props.accountCategories, props.sortBy, props.itemCount));
</script>
