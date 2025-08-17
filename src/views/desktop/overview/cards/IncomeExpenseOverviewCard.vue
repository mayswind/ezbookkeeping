<template>
    <v-card :class="{ 'disabled': disabled }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="secondary" size="38">
                <v-icon size="24" :icon="icon" />
            </v-avatar>
            <span class="font-weight-bold ms-3">{{ title }}</span>
            <v-spacer/>
            <v-btn density="comfortable" color="default" variant="text" class="ms-2" :icon="true">
                <v-icon :icon="mdiDotsVertical" />
                <v-menu activator="parent">
                    <v-list>
                        <slot name="menus"></slot>
                    </v-list>
                </v-menu>
            </v-btn>
        </v-card-text>
        <v-card-text class="mt-1 pb-1">
            <div class="font-weight-semibold text-truncate text-h4 text-income me-2 mb-2" v-if="!loading || incomeAmount">{{ incomeAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-4 mb-8" type="text" width="120px" :loading="true" v-else-if="loading && !incomeAmount"></v-skeleton-loader>
            <div class="text-truncate text-h5 text-expense" v-if="!loading || expenseAmount">{{ expenseAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mb-1" style="padding-bottom: 2px" type="text" width="120px" :loading="true" v-else-if="loading && !expenseAmount"></v-skeleton-loader>
            <div class="text-truncate text-h5 mt-2 mb-7" style="padding-bottom: 2px" v-if="!loading && !incomeAmount && !expenseAmount">{{ tt('No data') }}</div>
        </v-card-text>
        <v-card-text class="mt-6">
            <span class="text-caption">{{ datetime }}</span>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';

import {
    mdiDotsVertical
} from '@mdi/js';

defineProps<{
    loading: boolean;
    disabled: boolean;
    icon: string;
    title: string;
    expenseAmount: string;
    incomeAmount: string;
    datetime: string;
}>();

const { tt } = useI18n();
</script>
