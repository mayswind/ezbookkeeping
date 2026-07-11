<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Chart Color Scheme')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !canSaveColorScheme }" @click="saveChartColors"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list ref="colorSchemeList" strong inset dividers sortable sortable-enabled class="margin-top-half chart-color-list"
                 @sortable:sort="onSort">
            <f7-list-item :id="getColorIndexDomId(index)"
                          :key="index"
                          :title="'#' + color.toLowerCase()"
                          v-for="(color, index) in chartColors">
                <template #media>
                    <div class="color-preview-wrapper">
                        <input type="color" class="color-picker-input"
                               :value="'#' + color" @input="onColorInput(index, $event)" />
                        <div class="color-preview-box" :style="{ backgroundColor: '#' + color }"></div>
                    </div>
                </template>
                <template #after>
                    <div class="display-flex align-items-center">
                        <f7-link icon-f7="minus_circle_fill" color="red"
                                 @click="removeChartColor(index)"></f7-link>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="addNewColor()">{{ tt('Add') }}</f7-actions-button>
                <f7-actions-button @click="importText = ''; showImportSheet = true">{{ tt('Import') }}</f7-actions-button>
                <f7-actions-button @click="showExportSheet = true">{{ tt('Export') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button @click="resetToDefault()">{{ tt('Reset to Default') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
                  :opened="showImportSheet" @sheet:closed="showImportSheet = false">
            <div class="swipe-handler"></div>
            <f7-page-content class="margin-top no-padding-top">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ tt('Import') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-list strong inset dividers class="no-margin margin-bottom">
                        <f7-list-input
                            type="textarea"
                            class="import-chart-color-scheme-textarea"
                            :placeholder="tt('Each line should be a hex color value (e.g. c67e48 or #c67e48)')"
                            :value="importText"
                            @input="importText = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill
                               :class="{ 'disabled': !importText }"
                               :text="tt('Import')"
                               @click="doImport">
                    </f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="showImportSheet = false" :text="tt('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <information-sheet class="backup-code-sheet"
                           :title="tt('Export')"
                           :information="textualChartColors"
                           :row-count="15"
                           :enable-copy="true"
                           v-model:show="showExportSheet"
                           @info:copied="onColorsCopied">
        </information-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, useTemplateRef, nextTick } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useChartColorSchemeSettingsPageBase } from '@/views/base/settings/ChartColorSchemeSettingsPageBase.ts';

import type { ColorValue } from '@/core/color.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const {
    chartColors,
    textualChartColors,
    canSaveColorScheme,
    filterValidColors,
    addChartColor,
    removeChartColor,
    resetChartColorsToDefault,
    loadChartColorsFromSettings,
    saveChartColorsToSettings,
    onColorInput
} = useChartColorSchemeSettingsPageBase();

const colorSchemeList = useTemplateRef<{ $el: HTMLElement }>('colorSchemeList');

const showMoreActionSheet = ref<boolean>(false);
const showImportSheet = ref<boolean>(false);
const showExportSheet = ref<boolean>(false);
const importText = ref<string>('');

function getColorIndexDomId(index: number): string {
    return 'chart_color_' + index;
}

function parseColorIndexFromDomId(domId: string): string | null {
    if (!domId || domId.indexOf('chart_color_') !== 0) {
        return null;
    }

    return domId.substring(12); // chart_color_
}

function init(): void {
    loadChartColorsFromSettings();
}

function addNewColor(): void {
    addChartColor();

    nextTick(() => {
        console.log(colorSchemeList.value?.$el)

        colorSchemeList.value?.$el?.querySelector('ul > li:last-child')?.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        });
    });
}

function saveChartColors(): void {
    saveChartColorsToSettings();
    showToast('Chart color scheme saved');
    props.f7router.back();
}

function resetToDefault(): void {
    resetChartColorsToDefault();
    showMoreActionSheet.value = false;
}

function doImport(): void {
    const validColors = filterValidColors(importText.value);

    if (validColors.length < 1) {
        showToast('No valid colors found');
        return;
    }

    textualChartColors.value = importText.value;
    showImportSheet.value = false;
    importText.value = '';
    showToast('Chart color scheme imported');
}

function onSort(event: { el: { id: string }, from: number, to: number }): void {
    if (!event || !event.el || !event.el.id) {
        showToast('Unable to move color');
        return;
    }

    const indexStr = parseColorIndexFromDomId(event.el.id);

    if (!indexStr) {
        showToast('Unable to move color');
        return;
    }

    const fromIndex = parseInt(indexStr);

    if (isNaN(fromIndex) || fromIndex < 0 || fromIndex >= chartColors.value.length) {
        showToast('Unable to move color');
        return;
    }

    chartColors.value.splice(event.to, 0, chartColors.value.splice(event.from, 1)[0] as ColorValue);
}

function onColorsCopied(): void {
    showToast('Data copied');
}

init();
</script>

<style>
.color-preview-wrapper {
    position: relative;
    width: 28px;
    height: 28px;
}

.color-picker-input {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0;
    cursor: pointer;
}

.color-preview-box {
    width: 28px;
    height: 28px;
    border-radius: 4px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    pointer-events: none;
}

.chart-color-list .item-content > .item-inner > .item-title {
    font-family: monospace;
}

.import-chart-color-scheme-textarea textarea {
    height: 200px;

    @media (min-height: 630px) {
        height: 240px;
    }
}
</style>
