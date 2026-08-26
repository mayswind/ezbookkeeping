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

        <f7-list ref="colorSchemeList" strong inset dividers sortable sortable-enabled
                 class="margin-top-half chart-color-list" @sortable:sort="onSort">
            <f7-list-input readonly type="colorpicker"
                           :id="`chart_color_${chartColorIds[index]}`"
                           :key="chartColorIds[index]"
                           :color-picker-params="{
                               modules: ['sb-spectrum', 'hue-slider', 'hex'],
                               hexLabel: false,
                               hexValueEditable: true
                           }"
                           :value="{ hex: '#' + color }"
                           @colorpicker:change="updateChartColor(index, $event.hex ?? '')"
                           v-for="(color, index) in chartColors">
                <template #inner-start>
                    <ItemIcon icon-type="fixed-f7" icon-id="app_fill" size="36px" :color="color" @click="openColorPicker"></ItemIcon>
                </template>
                <template #inner-end>
                    <div class="chart-color-actions">
                        <f7-link icon-f7="minus_circle_fill" color="red"
                                 @click="removeColor(index)"></f7-link>
                    </div>
                </template>
            </f7-list-input>
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
                            class="import-chart-color-scheme-textarea code-textarea"
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

        <information-sheet class="code-information-sheet"
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

import { isInteger } from '@/lib/common.ts';

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
    updateChartColor,
    resetChartColorsToDefault,
    loadChartColorsFromSettings,
    saveChartColorsToSettings
} = useChartColorSchemeSettingsPageBase();

const colorSchemeList = useTemplateRef<{ $el: HTMLElement }>('colorSchemeList');

const showMoreActionSheet = ref<boolean>(false);
const showImportSheet = ref<boolean>(false);
const showExportSheet = ref<boolean>(false);
const chartColorIds = ref<number[]>([]);
const nextChartColorId = ref<number>(0);
const importText = ref<string>('');

function resetChartColorIds(): void {
    chartColorIds.value = chartColors.value.map(() => nextChartColorId.value++);
}

function init(): void {
    loadChartColorsFromSettings();
    resetChartColorIds();
}

function openColorPicker(event: Event): void {
    const colorInput = (event.currentTarget as HTMLElement | null)?.closest('.item-inner')?.querySelector<HTMLInputElement>('input');
    colorInput?.click();
}

function addNewColor(): void {
    addChartColor();
    chartColorIds.value.push(nextChartColorId.value++);

    nextTick(() => {
        colorSchemeList.value?.$el?.querySelector('ul > li:last-child')?.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        });
    });
}

function removeColor(index: number): void {
    removeChartColor(index);
    chartColorIds.value.splice(index, 1);
}

function saveChartColors(): void {
    saveChartColorsToSettings();
    showToast('Chart color scheme saved');
    props.f7router.back();
}

function resetToDefault(): void {
    resetChartColorsToDefault();
    resetChartColorIds();
    showMoreActionSheet.value = false;
}

function doImport(): void {
    const validColors = filterValidColors(importText.value);

    if (validColors.length < 1) {
        showToast('No valid colors found');
        return;
    }

    textualChartColors.value = importText.value;
    resetChartColorIds();
    showImportSheet.value = false;
    importText.value = '';
    showToast('Chart color scheme imported');
}

function onSort(event: { from: number, to: number }): void {
    if (!event || !isInteger(event.from) || !isInteger(event.to) ||
        event.from < 0 || event.from >= chartColors.value.length ||
        event.to < 0 || event.to >= chartColors.value.length) {
        showToast('Unable to move color');
        return;
    }

    if (event.from === event.to) {
        return;
    }

    chartColors.value.splice(event.to, 0, chartColors.value.splice(event.from, 1)[0] as ColorValue);
    chartColorIds.value.splice(event.to, 0, chartColorIds.value.splice(event.from, 1)[0] as number);
}

function onColorsCopied(): void {
    showToast('Data copied');
}

init();
</script>

<style>
.color-preview-box {
    width: 28px;
    height: 28px;
    border-radius: 4px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    pointer-events: none;
}

.chart-color-list {
    > ul {
        > li {
            > .item-input {
                > .item-inner {
                    --f7-list-item-padding-vertical: 0px;
                    flex-direction: row;
                    align-items: center;

                    input.input-with-value {
                        margin-inline: 4px;
                        font-family: monospace;
                        font-size: var(--f7-list-item-title-font-size);
                    }
                }
            }
        }
    }
}

html[dir="rtl"] .chart-color-list {
    > ul {
        > li {
            > .item-input {
                > .item-inner {
                    input.input-with-value {
                        direction: ltr;
                        text-align: right;
                    }
                }
            }
        }
    }
}

.import-chart-color-scheme-textarea textarea {
    height: 200px;

    @media (min-height: 630px) {
        height: 240px;
    }
}
</style>
