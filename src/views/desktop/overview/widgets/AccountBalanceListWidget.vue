<template>
    <v-card class="overview-widget overview-widget--list overview-widget--accounts h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Account Balance List')" :icon="mdiCreditCardOutline" />
        </template>

        <v-card-text class="overview-widget__body overview-widget__list-body">
            <v-list class="py-0" density="compact" v-if="displayAccounts.length">
                <v-list-item class="overview-widget__list-item" :key="account.id" :title="account.name"
                             :to="`/transaction/list?accountIds=${account.id}&dateType=${DateRange.All.type}`"
                             v-for="account in displayAccounts">
                    <template #prepend>
                        <span class="overview-widget__item-icon">
                            <ItemIcon size="24px" :icon-type="getAccountIconType(account.iconType)" :icon-id="account.icon" :color="account.color" />
                        </span>
                    </template>
                    <template #append>
                        <span class="overview-widget__list-amount overview-widget__amount">{{ accountBalance(account) }}</span>
                    </template>
                </v-list-item>
            </v-list>
            <div v-if="loading && !displayAccounts.length">
                <v-skeleton-loader class="skeleton-no-margin mx-2 py-3" style="margin-bottom: 1px" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="overview-widget__empty" v-else-if="!loading && !displayAccounts.length">
                <v-icon :icon="mdiCreditCardOutline" size="32" />
                <span>{{ tt('No data') }}</span>
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountListPageBase } from '@/views/base/accounts/AccountListPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';

import { DateRange } from '@/core/datetime.ts';
import type { Account } from '@/models/account.ts';

import { getAccountIconType } from '@/lib/icon.ts';

import {
    mdiCreditCardOutline
} from '@mdi/js';

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
