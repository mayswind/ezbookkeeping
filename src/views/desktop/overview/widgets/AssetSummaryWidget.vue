<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ tt('Asset Summary') }}</span>
        </template>

        <v-card-text class="mt-4">
            <div class="mb-16">
                <span class="text-body-medium" v-if="!loading || (allAccounts && allAccounts.length)">{{ tt('format.misc.youHaveAccounts', { count: displayAccountCount }) }}</span>
                <v-skeleton-loader class="skeleton-no-margin py-1" width="200px" type="text" v-else-if="loading && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
            </div>

            <v-row class="my-1">
                <v-col cols="12" md="4" :key="item.title" v-for="item in summaryItems">
                    <div class="d-flex align-center">
                        <div class="me-3">
                            <v-avatar rounded :color="item.color" size="42" class="elevation-1">
                                <v-icon size="24" :icon="item.icon" />
                            </v-avatar>
                        </div>
                        <div class="d-flex flex-column overflow-hidden">
                            <span class="text-body-medium">{{ tt(item.title) }}</span>
                            <span class="text-body-large text-truncate" v-if="!loading || (allAccounts && allAccounts.length)">{{ item.value }}</span>
                            <v-skeleton-loader class="skeleton-no-margin mb-1" style="margin-top: 6px" width="120px" type="text" :loading="true" v-else-if="loading && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
                        </div>
                    </div>
                </v-col>
            </v-row>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useHomePageBase } from '@/views/base/HomePageBase.ts';

import {
    mdiBankOutline,
    mdiCreditCardOutline,
    mdiPiggyBankOutline
} from '@mdi/js';

interface SummaryItem {
    title: string;
    value: string;
    icon: string;
    color: string;
}

defineProps<{
    loading: boolean
}>();

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const {
    allAccounts,
    netAssets,
    totalAssets,
    totalLiabilities
} = useHomePageBase();

const displayAccountCount = computed<string>(() => formatNumberToLocalizedNumerals(allAccounts.value?.length ?? 0));

const summaryItems = computed<SummaryItem[]>(() => [
    {
        title: 'Total assets',
        value: totalAssets.value,
        icon: mdiBankOutline,
        color: 'grey'
    },
    {
        title: 'Total liabilities',
        value: totalLiabilities.value,
        icon: mdiCreditCardOutline,
        color: 'expense'
    },
    {
        title: 'Net assets',
        value: netAssets.value,
        icon: mdiPiggyBankOutline,
        color: 'primary'
    }
]);
</script>
