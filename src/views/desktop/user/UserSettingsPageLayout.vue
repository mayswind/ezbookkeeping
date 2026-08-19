<template>
    <main-page-layout>
        <template #nav-items>
            <li class="nav-section-title">
                <div class="title-wrapper">
                    <span class="title-text">{{ tt('User Settings') }}</span>
                </div>
            </li>
            <li class="nav-link">
                <router-link to="/user/settings/basic">
                    <v-icon class="nav-item-icon" :icon="mdiAccountOutline"/>
                    <span class="nav-item-title">{{ tt('Basic Settings') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/user/settings/security">
                    <v-icon class="nav-item-icon" :icon="mdiLockOpenOutline"/>
                    <span class="nav-item-title">{{ tt('Security') }}</span>
                </router-link>
            </li>
            <li class="nav-link">
                <router-link to="/user/settings/two_factor">
                    <v-icon class="nav-item-icon" :icon="mdiOnepassword"/>
                    <span class="nav-item-title">{{ tt('Two-Factor Authentication') }}</span>
                </router-link>
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
            <li class="nav-link" v-if="isUserBudgetingEnabled()"> <!-- [PLUGIN:budget] -->
                <router-link to="/budget/overview">
                    <v-icon class="nav-item-icon" :icon="mdiWalletOutline"/>
                    <span class="nav-item-title">{{ tt('Budgets') }}</span>
                </router-link>
            </li>
            <li class="nav-link" v-if="isRulesEngineEnabled()"> <!-- [PLUGIN:rules] -->
                <router-link to="/rule/list">
                    <v-icon class="nav-item-icon" :icon="mdiAutoFix"/>
                    <span class="nav-item-title">{{ tt('Rules') }}</span>
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
                <router-link to="/user/settings/data_management">
                    <v-icon class="nav-item-icon" :icon="mdiDatabaseCogOutline"/>
                    <span class="nav-item-title">{{ tt('Data Management') }}</span>
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

import { isUserCustomIconEnabled, isUserScheduledTransactionEnabled, isUserBudgetingEnabled, isRulesEngineEnabled } from '@/lib/server_settings.ts'; // [PLUGIN:budget] [PLUGIN:rules] for budgeting & rules engine

import {
    mdiAccountOutline,
    mdiViewDashboardOutline,
    mdiTagOutline,
    mdiWalletOutline, // [PLUGIN:budget]
    mdiAutoFix, // [PLUGIN:rules]
    mdiClipboardTextOutline,
    mdiClipboardTextClockOutline,
    mdiShapePlusOutline,
    mdiLockOpenOutline,
    mdiOnepassword,
    mdiDatabaseCogOutline
} from '@mdi/js';

const route = useRoute();

const { tt } = useI18n();

const currentRoutePath = computed<string>(() => route.path);
</script>
