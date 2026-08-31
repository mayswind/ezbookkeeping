<template>
    <f7-list strong inset dividers class="overview-widget-list no-margin-top margin-bottom" :class="{ 'skeleton-text': loading }">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || tt('Account Balance List') }}</small>
        </f7-list-item>
        <f7-list-item :key="account.id" :title="account.name"
                      :after="accountBalance(account, undefined, showAmount, accountIds) || ''"
                      :link="`/transaction/list?accountIds=${account.id}&dateType=${DateRange.All.type}`"
                      v-for="account in displayAccounts">
            <template #media>
                <ItemIcon :icon-type="getAccountIconType(account.iconType)" :icon-id="account.icon" :color="account.color" />
            </template>
        </f7-list-item>
        <template v-if="loading && !displayAccounts.length">
            <f7-list-item link="#" title="Account" after="0.00 USD" :key="idx" v-for="idx in itemCount">
                <template #media>
                    <f7-icon f7="app_fill"></f7-icon>
                </template>
            </f7-list-item>
        </template>
        <f7-list-item :title="tt('No data')" v-else-if="!loading && !displayAccounts.length"></f7-list-item>
    </f7-list>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountListPageBase } from '@/views/base/accounts/AccountListPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { DateRange } from '@/core/datetime.ts';
import type { Account } from '@/models/account.ts';

import { getAccountIconType } from '@/lib/icon.ts';

const props = defineProps<{
    loading: boolean;
    title?: string;
    showTitle: boolean;
    accountIds: string[];
    itemCount: number;
    sortBy: string;
    alwaysShowAmount: boolean;
}>();

const { tt } = useI18n();

const { accountBalance } = useAccountListPageBase();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();

const showAmount = computed<boolean>(() => !!props.alwaysShowAmount || settingsStore.appSettings.showAmountInHomePage);
const displayAccounts = computed<Account[]>(() => accountsStore.getSortedAccounts(props.accountIds, props.sortBy, props.itemCount));
</script>
