<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left>
                <f7-link icon-f7="xmark" @click="cancel"></f7-link>
            </f7-nav-left>
            <f7-nav-title :title="tt('Home Page Layout')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ disabled: !isModified }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list sortable sortable-enabled sortable-tap-hold class="overview-layout-editor no-margin"
                 :sortable-move-elements="false" @sortable:sort="onSort" v-if="draftLayout.widgets.length">
            <li class="cursor-pointer" :key="widget.id" v-for="widget in draftLayout.widgets">
                <overview-widget class="overview-widget-editor-content" :widget="widget" :loading="loadingOverview" />
                <div class="overview-widget-drag-area" @click="showWidgetActions(widget)"></div>
            </li>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="!draftLayout.widgets.length">
            <f7-list-item :title="tt('No widgets')"></f7-list-item>
            <f7-list-button :title="tt('Add Widget')" @click="showAddWidgetPopup = true"></f7-list-button>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showWidgetActionSheet" @actions:closed="showWidgetActionSheet = false; selectedWidget = null">
            <f7-actions-group>
                <f7-actions-label v-if="selectedWidget">{{ tt(getMobileWidgetName(selectedWidget.type)) }}</f7-actions-label>
                <f7-actions-button @click="showConfigureWidgetPopup = true" v-if="selectedWidget && hasWidgetSettings(selectedWidget.type)">{{ tt('Settings') }}</f7-actions-button>
                <f7-actions-button color="red" @click="removeSelectedWidget">{{ tt('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="showAddWidgetPopup = true">{{ tt('Add Widget') }}</f7-actions-button>
                <f7-actions-button @click="resetLayout">{{ tt('Reset to Default') }}</f7-actions-button>
                <f7-actions-button :class="{ disabled: !draftLayout.widgets.length }" @click="clearLayout">{{ tt('Clear Layout') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button @click="openImportSheet">{{ tt('Import Layout') }}</f7-actions-button>
                <f7-actions-button @click="showExportSheet = true">{{ tt('Export Layout') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-popup push swipe-to-close :opened="showAddWidgetPopup" @popup:closed="showAddWidgetPopup = false">
            <f7-page>
                <f7-navbar>
                    <div class="swipe-handler"></div>
                    <f7-nav-title>{{ tt('Add Widget') }}</f7-nav-title>
                </f7-navbar>
                <f7-page-content class="no-padding-vertical">
                    <f7-list strong inset dividers class="margin-vertical">
                        <f7-list-item link="#" no-chevron :key="definition.type"
                                      :title="tt(definition.name)"
                                      v-for="definition in mobileWidgetDefinitions"
                                      @click="addWidget(definition.type)">
                            <template #after>
                                <f7-icon f7="plus"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-page-content>
            </f7-page>
        </f7-popup>

        <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
                  :opened="showImportSheet" @sheet:closed="showImportSheet = false">
            <div class="swipe-handler"></div>
            <f7-page-content class="margin-top no-padding-top">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ tt('Import Layout') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-list strong inset dividers class="no-margin margin-bottom">
                        <f7-list-input
                            type="textarea"
                            class="import-chart-color-scheme-textarea code-textarea"
                            :value="importText"
                            @input="importText = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill
                               :class="{ 'disabled': !importText }"
                               :text="tt('Import')"
                               @click="importLayout">
                    </f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="showImportSheet = false" :text="tt('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <information-sheet class="code-information-sheet"
                           :title="tt('Export Layout')"
                           :information="exportText"
                           :row-count="15"
                           :enable-copy="true"
                           v-model:show="showExportSheet"
                           @info:copied="onLayoutCopied"/>

        <widget-settings-popup :model-value="selectedWidget"
                               v-model:show="showConfigureWidgetPopup"
                               @update:model-value="updateWidgetSettings" />
    </f7-page>
</template>

<script setup lang="ts">
import OverviewWidget from './OverviewWidget.vue';
import WidgetSettingsPopup from './WidgetSettingsPopup.vue';

import { ref, computed, useTemplateRef } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { itemAndIndex } from '@/core/base.ts';
import {
    type OverviewWidgetType,
    type MobileOverviewLayout,
    type MobileOverviewWidgetDefinition,
    type MobileOverviewWidgetLayout,
    OverviewWidgetDataRequirement
} from '@/core/overview_layout.ts';
import { MOBILE_OVERVIEW_WIDGET_DEFINITIONS, DEFAULT_MOBILE_OVERVIEW_LAYOUT } from '@/consts/overview_layout.ts';

import {
    getOverviewDataRequirements,
    getOverviewTransactionOverviewMonths,
    getOverviewRecentTransactionCount,
    getOverviewAssetTrendMonths,
    getOverviewCalendarHeatmapMonths,
    getOverviewTransactionCategoryStatisticDateTypes,
    isDefaultMobileOverviewLayout,
    cloneMobileOverviewLayout,
    parseMobileOverviewLayout,
    serializeMobileOverviewLayout
} from '@/lib/overview_layout.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import logger from '@/lib/logger.ts';

type WidgetSettingsPopupType = InstanceType<typeof WidgetSettingsPopup>;

const props = defineProps<{
    f7router: Router.Router
}>();

const { tt } = useI18n();
const { showToast, showConfirm } = useI18nUIComponents();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();

const widgetSettingsPopup = useTemplateRef<WidgetSettingsPopupType>('widgetSettingsPopup');

const loadingOverview = ref<boolean>(true);
const showMoreActionSheet = ref<boolean>(false);
const showImportSheet = ref<boolean>(false);
const showExportSheet = ref<boolean>(false);
const showWidgetActionSheet = ref<boolean>(false);
const showAddWidgetPopup = ref<boolean>(false);
const showConfigureWidgetPopup = ref<boolean>(false);
const initialLayout = ref<MobileOverviewLayout>(getInitialLayout());
const draftLayout = ref<MobileOverviewLayout>(cloneMobileOverviewLayout(initialLayout.value));
const initialJson = ref<string>(serializeMobileOverviewLayout(initialLayout.value));
const selectedWidget = ref<MobileOverviewWidgetLayout | null>(null);
const importText = ref<string>('');

const isModified = computed<boolean>(() => serializeMobileOverviewLayout(draftLayout.value) !== initialJson.value);
const exportText = computed<string>(() => serializeMobileOverviewLayout(draftLayout.value, true));
const mobileWidgetDefinitions = computed<MobileOverviewWidgetDefinition[]>(() => Object.values(MOBILE_OVERVIEW_WIDGET_DEFINITIONS).filter(Boolean) as MobileOverviewWidgetDefinition[]);

function getInitialLayout(): MobileOverviewLayout {
    let initialLayout: MobileOverviewLayout;

    try {
        initialLayout = parseMobileOverviewLayout(settingsStore.appSettings.mobileOverviewPageLayout);
    } catch (error) {
        logger.warn('failed to parse mobile overview page layout in editor', error);
        initialLayout = cloneMobileOverviewLayout(DEFAULT_MOBILE_OVERVIEW_LAYOUT);
    }

    return initialLayout;
}

function getMobileWidgetName(type: OverviewWidgetType): string {
    return MOBILE_OVERVIEW_WIDGET_DEFINITIONS[type]?.name ?? 'Unknown';
}

function hasWidgetSettings(type: OverviewWidgetType): boolean {
    return !!MOBILE_OVERVIEW_WIDGET_DEFINITIONS[type]?.supportsSettings.length;
}

function reload(force: boolean): void {
    loadingOverview.value = true;

    const requirements: OverviewWidgetDataRequirement[] = getOverviewDataRequirements(draftLayout.value, MOBILE_OVERVIEW_WIDGET_DEFINITIONS);
    const promises: Promise<unknown>[] = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false })
    ];

    if (requirements.includes(OverviewWidgetDataRequirement.TransactionOverview)) {
        promises.push(overviewStore.loadTransactionOverview({
            force: force,
            months: getOverviewTransactionOverviewMonths(draftLayout.value)
        }));
    }

    if (requirements.includes(OverviewWidgetDataRequirement.TransactionCategoryStatistics)) {
        for (const dateType of getOverviewTransactionCategoryStatisticDateTypes(draftLayout.value)) {
            promises.push(overviewStore.loadTransactionCategoryStatistics({
                force: force,
                dateType: dateType
            }));
        }
    }

    if (requirements.includes(OverviewWidgetDataRequirement.AssetTrends)) {
        promises.push(overviewStore.loadTransactionAssetTrends({
            force: force,
            months: getOverviewAssetTrendMonths(draftLayout.value)
        }));
    }

    if (requirements.includes(OverviewWidgetDataRequirement.RecentTransactions)) {
        promises.push(overviewStore.loadRecentTransactions({
            force: force,
            count: getOverviewRecentTransactionCount(draftLayout.value)
        }));
    }

    if (requirements.includes(OverviewWidgetDataRequirement.DailyTransactionAmounts)) {
        promises.push(overviewStore.loadTransactionDailyAmounts({
            force: force,
            months: getOverviewCalendarHeatmapMonths(draftLayout.value)
        }));
    }

    Promise.all(promises).then(() => {
        loadingOverview.value = false;

        if (force) {
            showToast('Data has been updated');
        }
    }).catch(error => {
        loadingOverview.value = false;

        if (!error.processed && !error.isUpToDate) {
            showToast(error.message || error);
        }
    });
}

function addWidget(type: OverviewWidgetType): void {
    const definition = MOBILE_OVERVIEW_WIDGET_DEFINITIONS[type];

    if (!definition) {
        return;
    }

    draftLayout.value.widgets.push({
        id: generateRandomUUID(),
        type: type,
        settings: { ...definition.defaultSettings }
    });

    showAddWidgetPopup.value = false;
    reload(false);
}

function removeWidget(id: string): void {
    for (const [widget, index] of itemAndIndex(draftLayout.value.widgets)) {
        if (widget.id === id) {
            draftLayout.value.widgets.splice(index, 1);
            return;
        }
    }
}

function showWidgetActions(widget: MobileOverviewWidgetLayout): void {
    selectedWidget.value = widget;
    showWidgetActionSheet.value = true;
}

function removeSelectedWidget(): void {
    const widgetId = selectedWidget.value?.id;

    if (widgetId) {
        removeWidget(widgetId);
    }
}

function updateWidgetSettings(updatedWidget: MobileOverviewWidgetLayout): void {
    for (const currentWidget of draftLayout.value.widgets) {
        if (currentWidget.id === updatedWidget.id) {
            currentWidget.settings = updatedWidget.settings;
            break;
        }
    }

    reload(false);
}

function clearLayout(): void {
    showConfirm('Clear all widgets from the layout?', () => {
        draftLayout.value.widgets.length = 0;
    });
}

function resetLayout(): void {
    showConfirm('Reset the layout to its default value?', () => {
        draftLayout.value = cloneMobileOverviewLayout(DEFAULT_MOBILE_OVERVIEW_LAYOUT);
    });
}

function openImportSheet(): void {
    importText.value = '';
    showImportSheet.value = true;
}

function importLayout(): void {
    try {
        draftLayout.value = parseMobileOverviewLayout(importText.value);
        showImportSheet.value = false;
    } catch (error) {
        logger.error('failed to import mobile overview layout', error);
        showToast('Layout import failed. Please make sure the layout is valid and try again.');
    }
}

function save(): void {
    if (!isModified.value) {
        return;
    }

    settingsStore.setMobileOverviewPageLayout(isDefaultMobileOverviewLayout(draftLayout.value) ? '' : serializeMobileOverviewLayout(draftLayout.value));
    initialJson.value = serializeMobileOverviewLayout(draftLayout.value);
    props.f7router.back();
}

function cancel(): void {
    if (!isModified.value) {
        props.f7router.back();
        return;
    }

    showConfirm('Discard unsaved layout changes?', () => {
        props.f7router.back();
    });
}

function onSort(event: { from: number, to: number }): void {
    const widget = draftLayout.value.widgets.splice(event.from, 1)[0];

    if (!widget) {
        return;
    }

    draftLayout.value.widgets.splice(event.to, 0, widget);
}

function onLayoutCopied(): void {
    showToast('Data copied');
}

reload(false);
</script>

<style>
.overview-layout-editor {
    --f7-list-bg-color: transparent;
    --f7-sortable-sorting-item-bg-color: transparent;
    --f7-sortable-sorting-item-box-shadow: none;

    .overview-widget-drag-area {
        position: absolute;
        z-index: 1;
        inset: 0;
        width: auto;
        height: auto;
        cursor: grab;

        &::after {
            display: none;
        }
    }

    > ul > li .overview-widget-editor-content {
        transition: transform 180ms ease, filter 180ms ease;
        transform-origin: center;
    }

    > ul > li.sorting .overview-widget-editor-content {
        transform: translateY(-4px) scale(1.015);
        filter: drop-shadow(0 10px 12px rgba(0, 0, 0, 0.22));
    }

    li.sorting .overview-widget-drag-area {
        cursor: grabbing;
    }

    .overview-widget-editor-content {
        pointer-events: none;
    }

    &.sortable-enabled:not(.sortable-opposite) > ul > li .overview-widget-editor-content .item-link .item-inner::before {
        display: block;
    }
}

.overview-layout-editor > ul > li > .list > ul {
    padding-inline-start: 0;
}

html:not([dir="rtl"]) .overview-layout-editor {
    &.sortable-enabled:not(.sortable-opposite) > ul > li .overview-widget-editor-content .item-link .item-inner {
        padding-right: calc(var(--f7-list-chevron-icon-area) + var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-right));
    }
}

html[dir="rtl"] .overview-layout-editor {
    &.sortable-enabled:not(.sortable-opposite) > ul > li .overview-widget-editor-content .item-link .item-inner {
        padding-left: calc(var(--f7-list-chevron-icon-area) + var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-left));
    }
}
</style>
