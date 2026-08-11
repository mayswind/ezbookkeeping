<template>
    <v-card :class="{ 'disabled': disabled }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="grey" size="32">
                <v-icon size="22" :icon="icon" />
            </v-avatar>
            <span class="text-title-small font-weight-bold ms-2">{{ title }}</span>
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
        <v-card-text class="py-3">
            <div class="text-truncate text-headline-small text-income me-2 mb-2" v-if="!loading || incomeAmount">{{ incomeAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-3 mb-6" type="text" width="120px" :loading="true" v-else-if="loading && !incomeAmount"></v-skeleton-loader>
            <div class="text-truncate text-title-medium text-expense pt-1" v-if="!loading || expenseAmount">{{ expenseAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin pt-1" style="padding-bottom: 7px" type="text" width="120px" :loading="true" v-else-if="loading && !expenseAmount"></v-skeleton-loader>
            <div class="text-truncate text-title-medium mt-3 mb-5" v-if="!loading && !incomeAmount && !expenseAmount">{{ tt('No data') }}</div>
        </v-card-text>
        <v-card-text>
            <span class="text-body-medium">{{ datetime }}</span>
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
