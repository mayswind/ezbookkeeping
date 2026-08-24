<template>
    <div class="mobile-overview-dashboard">
        <div class="mobile-overview-dashboard-item" :key="widget.id" v-for="widget in layout.widgets">
            <overview-widget :widget="widget" :loading="loading" />
        </div>

        <f7-block class="text-align-center text-color-gray padding-vertical" v-if="!layout.widgets.length">
            {{ tt('No widgets') }}
        </f7-block>
    </div>
</template>

<script setup lang="ts">
import OverviewWidget from './OverviewWidget.vue';

import { useI18n } from '@/locales/helpers.ts';

import type { MobileOverviewLayout } from '@/core/overview_layout.ts';

defineProps<{
    layout: MobileOverviewLayout;
    loading: boolean;
}>();

const { tt } = useI18n();
</script>

<style>
.mobile-overview-dashboard-item + .mobile-overview-dashboard-item {
    margin-block-start: 16px;
}

.mobile-overview-dashboard-item > .card,
.mobile-overview-dashboard-item > .list {
    margin-block-start: 0 !important;
    margin-block-end: 0 !important;
}
</style>
