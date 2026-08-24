<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-title :title="tt('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <overview-dashboard :layout="layout" :loading="loading" />

        <f7-toolbar tabbar icons bottom class="main-tabbar">
            <f7-link class="link" href="/transaction/list">
                <f7-icon f7="square_list"></f7-icon>
                <span class="tabbar-label">{{ tt('Details') }}</span>
            </f7-link>
            <f7-link class="link" href="/account/list">
                <f7-icon f7="creditcard"></f7-icon>
                <span class="tabbar-label">{{ tt('Accounts') }}</span>
            </f7-link>
            <!-- "homepage-add-button" must have the "dragenabled" class, otherwise the popover disappears immediately after the second long press -->
            <f7-link id="homepage-add-button" class="link dragenabled"
                     href="/transaction/add" @taphold="openTransactionTemplatePopover">
                <f7-icon f7="plus_square" class="ebk-tarbar-big-icon"></f7-icon>
            </f7-link>
            <f7-link class="link" href="/statistic/transaction">
                <f7-icon f7="chart_pie"></f7-icon>
                <span class="tabbar-label">{{ tt('Statistics') }}</span>
            </f7-link>
            <f7-link class="link" href="/settings">
                <f7-icon f7="gear_alt"></f7-icon>
                <span class="tabbar-label">{{ tt('Settings') }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-popover class="template-popover-menu" target-el="#homepage-add-button"
                    v-model:opened="showTransactionTemplatePopover">
            <f7-list dividers v-if="isTransactionFromAITextRecognitionEnabled() || isTransactionFromAIImageRecognitionEnabled() || (allTransactionTemplates && allTransactionTemplates.length)">
                <f7-list-item key="AIClipboardTextRecognition" link="#" no-chevron popover-close
                              :title="tt('AI Clipboard Text Recognition')"
                              @click="addByRecognizingClipboardText"
                              v-if="isTransactionFromAITextRecognitionEnabled()">
                    <template #media>
                        <f7-icon f7="wand_stars"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item key="AIImageRecognition" link="#" no-chevron popover-close
                              :title="tt('AI Image Recognition')"
                              @click="showAIReceiptImageRecognitionSheet = true"
                              v-if="isTransactionFromAIImageRecognitionEnabled()">
                    <template #media>
                        <f7-icon f7="wand_stars"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item popover-close :key="template.id" :title="template.name"
                              :link="'/transaction/add?templateId=' + template.id"
                              v-for="template in allTransactionTemplates">
                    <template #media>
                        <f7-icon f7="doc_plaintext"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <a-i-image-recognition-sheet ref="aiImageRecognitionSheet"
                                     v-model:show="showAIReceiptImageRecognitionSheet"
                                     @recognition:change="onReceiptRecognitionChanged"/>
    </f7-page>
</template>

<script setup lang="ts">
import AIImageRecognitionSheet, { type AIImageRecognitionResult } from '@/components/mobile/AIImageRecognitionSheet.vue';
import OverviewDashboard from './overview/OverviewDashboard.vue';

import { ref, computed, useTemplateRef } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, isiOS } from '@/lib/ui/mobile.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { type MobileOverviewLayout, OverviewWidgetDataRequirement } from '@/core/overview_layout.ts';
import { TemplateType } from '@/core/template.ts';
import { MOBILE_OVERVIEW_WIDGET_DEFINITIONS, DEFAULT_MOBILE_OVERVIEW_LAYOUT } from '@/consts/overview_layout.ts';

import { TransactionTemplate } from '@/models/transaction_template.ts';

import { isFunction } from '@/lib/common.ts';
import {
    getOverviewDataRequirements,
    getOverviewTransactionOverviewMonths,
    getOverviewRecentTransactionCount,
    getOverviewAssetTrendMonths,
    getOverviewCalendarHeatmapMonths,
    getOverviewTransactionCategoryStatisticDateTypes,
    parseMobileOverviewLayout
} from '@/lib/overview_layout.ts';
import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';
import { getShareCacheImageBlob } from '@/lib/cache.ts';
import {
    isTransactionFromAITextRecognitionEnabled,
    isTransactionFromAIImageRecognitionEnabled
} from '@/lib/server_settings.ts';
import logger from '@/lib/logger.ts';

type AIImageRecognitionSheetType = InstanceType<typeof AIImageRecognitionSheet>;

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTemplatesStore = useTransactionTemplatesStore();
const overviewStore = useOverviewStore();

const aiImageRecognitionSheet = useTemplateRef<AIImageRecognitionSheetType>('aiImageRecognitionSheet');

const loading = ref<boolean>(true);
const showTransactionTemplatePopover = ref<boolean>(false);
const showAIReceiptImageRecognitionSheet = ref<boolean>(false);

const layout = computed<MobileOverviewLayout>(() => {
    try {
        return parseMobileOverviewLayout(settingsStore.appSettings.mobileOverviewPageLayout);
    } catch (error) {
        logger.warn('failed to parse mobile overview page layout, fallback to default layout', error);
        return DEFAULT_MOBILE_OVERVIEW_LAYOUT;
    }
});

const allTransactionTemplates = computed<TransactionTemplate[]>(() => {
    const allTemplates = transactionTemplatesStore.allVisibleTemplates;
    return allTemplates[TemplateType.Normal.type] || [];
});

function openTransactionTemplatePopover(): void {
    if (isTransactionFromAIImageRecognitionEnabled() || (allTransactionTemplates.value && allTransactionTemplates.value.length)) {
        showTransactionTemplatePopover.value = true;
    }
}

function init(): void {
    if (isUserLogined() && isUserUnlocked()) {
        loading.value = true;

        const promises: Promise<unknown>[] = [
            getShareCacheImageBlob(),
            accountsStore.loadAllAccounts({ force: false }),
            transactionCategoriesStore.loadAllCategories({ force: false }),
            transactionTemplatesStore.loadAllTemplates({ templateType: TemplateType.Normal.type,  force: false }),
            ...reloadOverviewData(false)
        ];

        Promise.all(promises).then(responses => {
            if (responses[0] && responses[0] instanceof Blob) {
                aiImageRecognitionSheet.value?.loadImage(responses[0]);
                showAIReceiptImageRecognitionSheet.value = true;
            }

            loading.value = false;
        }).catch(error => {
            loading.value = false;

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    }
}

function reload(done?: () => void): void {
    const force = !!done;
    const promises: Promise<unknown>[] = reloadOverviewData(force);

    if (promises.length < 1) {
        done?.();
        return;
    }

    Promise.all(promises).then(() => {
        done?.();

        if (force) {
            showToast('Data has been updated');
        }
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function reloadOverviewData(force: boolean): Promise<unknown>[] {
    const requirements: OverviewWidgetDataRequirement[] = getOverviewDataRequirements(layout.value, MOBILE_OVERVIEW_WIDGET_DEFINITIONS);
    const promises: Promise<unknown>[] = [];

    if (requirements.includes(OverviewWidgetDataRequirement.TransactionOverview)) {
        promises.push(overviewStore.loadTransactionOverview({
            force: force,
            months: getOverviewTransactionOverviewMonths(layout.value)
        }));
    }

    if (requirements.includes(OverviewWidgetDataRequirement.TransactionCategoryStatistics)) {
        for (const dateType of getOverviewTransactionCategoryStatisticDateTypes(layout.value)) {
            promises.push(overviewStore.loadTransactionCategoryStatistics({
                force: force,
                dateType: dateType
            }));
        }
    }

    if (requirements.includes(OverviewWidgetDataRequirement.AssetTrends)) {
        promises.push(overviewStore.loadTransactionAssetTrends({
            force: force,
            months: getOverviewAssetTrendMonths(layout.value)
        }));
    }

    if (requirements.includes(OverviewWidgetDataRequirement.RecentTransactions)) {
        promises.push(overviewStore.loadRecentTransactions({
            force: force,
            count: getOverviewRecentTransactionCount(layout.value)
        }));
    }

    if (requirements.includes(OverviewWidgetDataRequirement.DailyTransactionAmounts)) {
        promises.push(overviewStore.loadTransactionDailyAmounts({
            force: force,
            months: getOverviewCalendarHeatmapMonths(layout.value)
        }));
    }

    return promises;
}

function addByRecognizingClipboardText(): void {
    if (navigator.clipboard && isFunction(navigator.clipboard.readText) && !isiOS()) {
        navigator.clipboard.readText().then(text => {
            const clipboardText = text && text.trim() ? text.trim() : '';
            props.f7router.navigate('/transaction/add', {
                props: {
                    autoRecognizeClipboardText: clipboardText,
                }
            });
        }).catch(error => {
            logger.error('failed to read clipboard', error);
            props.f7router.navigate('/transaction/add', {
                props: {
                    autoRecognizeClipboardText: '',
                }
            });
        });
    } else {
        props.f7router.navigate('/transaction/add', {
            props: {
                autoRecognizeClipboardText: '',
            }
        });
    }
}

function onReceiptRecognitionChanged(result: AIImageRecognitionResult): void {
    const recognizedResponse = result.response;
    const autoUploadRecognizedImage = settingsStore.appSettings.autoUploadTransactionPictureForAIRecognition;
    const params: string[] = [];

    if (recognizedResponse.type) {
        params.push(`type=${recognizedResponse.type}`);
    }

    if (recognizedResponse.time) {
        params.push(`time=${recognizedResponse.time}`);
    }

    if (recognizedResponse.categoryId) {
        params.push(`categoryId=${recognizedResponse.categoryId}`);
    }

    if (recognizedResponse.sourceAccountId) {
        params.push(`accountId=${recognizedResponse.sourceAccountId}`);
    }

    if (recognizedResponse.destinationAccountId) {
        params.push(`destinationAccountId=${recognizedResponse.destinationAccountId}`);
    }

    if (recognizedResponse.sourceAmount) {
        params.push(`amount=${recognizedResponse.sourceAmount}`);
    }

    if (recognizedResponse.destinationAmount) {
        params.push(`destinationAmount=${recognizedResponse.destinationAmount}`);
    }

    if (recognizedResponse.tagIds) {
        params.push(`tagIds=${recognizedResponse.tagIds.join(',')}`);
    }

    if (recognizedResponse.comment) {
        params.push(`comment=${encodeURIComponent(recognizedResponse.comment)}`);
    }

    params.push(`noTransactionDraft=true`);

    props.f7router.navigate(`/transaction/add?${params.join('&')}`, {
        props: {
            autoUploadPicture: autoUploadRecognizedImage ? result.imageFile : undefined,
        }
    });
}

function onPageAfterIn(): void {
    if (!loading.value) {
        reload();
    }
}

init();
</script>

<style>
.home-summary-card {
    background-color: var(--f7-color-yellow);
}

.home-summary-card .home-summary-month {
    font-size: 1.3em;
}

.home-summary-card .month-expense {
    font-size: 1.5em;
}

.home-summary-card .home-summary-misc {
    opacity: 0.6;
}

.home-summary-misc > span {
    margin-inline-end: 4px;
}

.home-summary-misc > span:last-child {
    margin-inline-end: 0;
}

.dark .home-summary-card {
    background-color: var(--f7-theme-color);
}

.dark .home-summary-card a {
    color: var(--f7-text-color);
    opacity: 0.6;
}

.overview-transaction-list .item-title > div {
    overflow: hidden;
    text-overflow: ellipsis;
}

.overview-transaction-list .item-after {
    max-width: 100%;
}

.overview-transaction-list .overview-transaction-footer {
    padding-top: 6px;
    font-size: var(--ebk-large-footer-font-size);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.overview-transaction-list .overview-transaction-footer > span {
    margin-inline-end: 4px;
}

.overview-transaction-list .overview-transaction-amount {
    max-width: 100%;
}

.overview-transaction-list .overview-transaction-amount > div {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tabbar.main-tabbar .link i + span.tabbar-label {
    margin-top: var(--ebk-icon-text-margin);
}

.tabbar.main-tabbar .link i.ebk-tarbar-big-icon {
    font-size: var(--ebk-big-icon-button-size);
    width: var(--ebk-big-icon-button-size);
    height: var(--ebk-big-icon-button-size);
    line-height: var(--ebk-big-icon-button-size);
}

.template-popover-menu .popover-inner {
    max-height: 400px;
    overflow-y: auto;
}
</style>
