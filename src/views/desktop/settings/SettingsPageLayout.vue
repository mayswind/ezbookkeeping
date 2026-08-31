<template>
    <main-page-layout>
        <template #nav-items>
            <li class="nav-section-title">
                <div class="title-wrapper">
                    <span class="title-text">{{ tt('User Settings') }}</span>
                </div>
            </li>
            <li class="nav-link">
                <router-link to="/settings/user/basic">
                    <v-icon class="nav-item-icon" :icon="mdiAccountOutline"/>
                    <span class="nav-item-title">{{ tt('User Profile') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/user/security">
                    <v-icon class="nav-item-icon" :icon="mdiLockOpenOutline"/>
                    <span class="nav-item-title">{{ tt('Security') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/user/two_factor">
                    <v-icon class="nav-item-icon" :icon="mdiOnepassword"/>
                    <span class="nav-item-title">{{ tt('Two-Factor Authentication') }}</span>
                </router-link>
            </li>

            <li class="nav-section-title">
                <div class="title-wrapper">
                    <span class="title-text">{{ tt('Basis Data') }}</span>
                </div>
            </li>
            <li class="nav-link">
                <router-link to="/category/list">
                    <v-icon class="nav-item-icon" :icon="mdiViewDashboardOutline"/>
                    <span class="nav-item-title">{{ tt('Transaction Categories') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/tag/list">
                    <v-icon class="nav-item-icon" :icon="mdiTagOutline"/>
                    <span class="nav-item-title">{{ tt('Transaction Tags') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/template/list">
                    <v-icon class="nav-item-icon" :icon="mdiClipboardTextOutline"/>
                    <span class="nav-item-title">{{ tt('Transaction Templates') }}</span>
                </router-link>
            </li>
            <li class="nav-link" v-if="isUserScheduledTransactionEnabled()">
                <router-link to="/schedule/list">
                    <v-icon class="nav-item-icon" :icon="mdiClipboardTextClockOutline"/>
                    <span class="nav-item-title">{{ tt('Scheduled Transactions') }}</span>
                </router-link>
            </li>
            <li class="nav-link" v-if="isUserCustomIconEnabled()">
                <router-link to="/custom_icon/list">
                    <v-icon class="nav-item-icon" :icon="mdiShapePlusOutline"/>
                    <span class="nav-item-title">{{ tt('Custom Icons') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/user/data_management">
                    <v-icon class="nav-item-icon" :icon="mdiDatabaseCogOutline"/>
                    <span class="nav-item-title">{{ tt('Data Management') }}</span>
                </router-link>
            </li>

            <li class="nav-section-title">
                <div class="title-wrapper">
                    <span class="title-text">{{ tt('Application Settings') }}</span>
                </div>
            </li>
            <li class="nav-link">
                <router-link to="/settings/preferences">
                    <v-icon class="nav-item-icon" :icon="mdiCogOutline"/>
                    <span class="nav-item-title">{{ tt('Preferences') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/application_lock">
                    <v-icon class="nav-item-icon" :icon="mdiLockOpenOutline"/>
                    <span class="nav-item-title">{{ tt('Application Lock') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/statistics">
                    <v-icon class="nav-item-icon" :icon="mdiChartPieOutline"/>
                    <span class="nav-item-title">{{ tt('Statistics Settings') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/cloud_sync">
                    <v-icon class="nav-item-icon" :icon="mdiCloudOutline"/>
                    <span class="nav-item-title">{{ tt('Settings Sync') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/exchange_rate">
                    <v-icon class="nav-item-icon" :icon="mdiSwapHorizontal"/>
                    <span class="nav-item-title">{{ tt('Exchange Rates Data') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/settings/browser_cache">
                    <v-icon class="nav-item-icon" :icon="mdiDatabaseClockOutline"/>
                    <span class="nav-item-title">{{ tt('Browser Cache Management') }}</span>
                </router-link>
            </li>
        </template>

        <template #content>
            <router-view :key="currentRoutePath" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { isUserCustomIconEnabled, isUserScheduledTransactionEnabled } from '@/lib/server_settings.ts';

import {
    mdiAccountOutline,
    mdiViewDashboardOutline,
    mdiTagOutline,
    mdiClipboardTextOutline,
    mdiClipboardTextClockOutline,
    mdiShapePlusOutline,
    mdiLockOpenOutline,
    mdiOnepassword,
    mdiDatabaseCogOutline,
    mdiCogOutline,
    mdiChartPieOutline,
    mdiCloudOutline,
    mdiSwapHorizontal,
    mdiDatabaseClockOutline
} from '@mdi/js';

const route = useRoute();

const { tt } = useI18n();

const currentRoutePath = computed<string>(() => route.path);
</script>
