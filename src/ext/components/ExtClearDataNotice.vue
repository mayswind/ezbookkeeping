<template>
    <v-alert class="mt-3" type="warning" variant="tonal" density="compact" v-if="available">
        <div>{{ tt('"Clear All Data" also removes your business records: items, stock, locations, customers, sales and repayments. Your team, receipt details and activity log are kept.') }}</div>
        <div class="mt-1">{{ tt('"Clear All Transactions" does not remove them, so sales and what customers owe stay as they are.') }}</div>
        <div class="mt-1">
            <router-link to="/settings/business">{{ tt('Download a copy of your business data first') }}</router-link>
        </div>
    </v-alert>
</template>

<script setup lang="ts">
import { useExtI18n } from '@/ext/i18n.ts';
import { useBusinessFeatures } from '@/ext/features.ts';

const { tt } = useExtI18n();
const { available } = useBusinessFeatures();
</script>
