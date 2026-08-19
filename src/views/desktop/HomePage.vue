<template>
    <main-page-layout no-navbar>
        <template #content>
            <overview-dashboard :layout="layout" :loading="loadingOverview" @refresh="reload(true)" />
        </template>
    </main-page-layout>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import OverviewDashboard from './overview/OverviewDashboard.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { type DesktopOverviewLayout, OverviewWidgetDataRequirement } from '@/core/overview_layout.ts';
import { DEFAULT_DESKTOP_OVERVIEW_LAYOUT } from '@/consts/overview_layout.ts';

import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';
import { getShareCacheImageBlob } from '@/lib/cache.ts';
import {
    getOverviewDataRequirements,
    getOverviewTransactionOverviewMonths,
    parseDesktopOverviewLayout
} from '@/lib/overview_layout.ts';
import logger from '@/lib/logger.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loadingOverview = ref<boolean>(true);

const layout = computed<DesktopOverviewLayout>(() => {
    try {
        return parseDesktopOverviewLayout(settingsStore.appSettings.desktopOverviewPageLayout);
    } catch (error) {
        logger.warn('failed to parse desktop overview page layout, fallback to default layout', error);
        return DEFAULT_DESKTOP_OVERVIEW_LAYOUT;
    }
});

function clearShareImageCache(): void {
    getShareCacheImageBlob().then(blob => {
        if (blob) {
            logger.warn('desktop version does not support receving shared image, the share image cache has been cleared');
        }
    });
}

function reload(force: boolean): void {
    loadingOverview.value = true;

    const requirements: OverviewWidgetDataRequirement[] = getOverviewDataRequirements(layout.value);
    const promises: Promise<unknown>[] = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false })
    ];

    if (requirements.includes(OverviewWidgetDataRequirement.TransactionOverview)) {
        promises.push(overviewStore.loadTransactionOverview({
            force: force,
            months: getOverviewTransactionOverviewMonths(layout.value)
        }));
    }

    Promise.all(promises).then(() => {
        loadingOverview.value = false;

        if (force) {
            snackbar.value?.showMessage('Data has been updated');
        }
    }).catch(error => {
        loadingOverview.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

clearShareImageCache();

if (isUserLogined() && isUserUnlocked()) {
    reload(false);
}
</script>

<style>
.overview-card-background {
    position: absolute;
    inline-size: 9rem;
    inset-block-end: 0;
    inset-inline-end: 0;
}

.overview-card-background-image {
    position: absolute;
    inline-size: 5rem;
    inset-block-end: 0.5rem;
    inset-inline-end: 1rem;
}
</style>
