<template>
    <v-dialog width="700" :persistent="chartColorsModified" v-model="showState">
        <one-column-dialog-layout content-class="pa-0"
                                  :title="tt('Chart Color Scheme')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true"
                       @click="addNewColor()">
                    <v-icon :icon="mdiPlus" size="22" />
                    <v-tooltip activator="parent">{{ tt('Add') }}</v-tooltip>
                </v-btn>
                <v-btn density="compact" color="primary" variant="text" class="ms-1" :icon="true"
                       :disabled="!canSaveColorScheme" @click="saveChartColors()">
                    <v-icon :icon="mdiCheck" size="22" />
                    <v-tooltip activator="parent">{{ tt('Save') }}</v-tooltip>
                </v-btn>
            </template>

            <template #toolbar>
                <toggle-button class="ms-2" :false-name="tt('List')" :true-name="tt('Raw Data')"
                               v-model="showRawData"/>

                <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true">
                    <v-icon :icon="mdiDotsVertical" size="22" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiRestore"
                                         :title="tt('Reset to Default')"
                                         @click="resetChartColorsToDefault"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </template>

            <template #content>
                <div class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto" style="height: 420px">
                    <v-table ref="colorSchemeTable" hover density="comfortable" class="w-100 table-striped" v-if="!showRawData">
                        <draggable-list tag="tbody"
                                        item-key="index"
                                        handle=".drag-handle"
                                        ghost-class="dragging-item"
                                        v-model="chartColors">
                            <template #item="{ element, index }">
                                <tr @mouseenter="hoveredIndex = index" @mouseleave="hoveredIndex = ''">
                                    <td>
                                        <div class="d-flex align-center">
                                            <div class="d-flex align-center">
                                                <div class="color-preview-wrapper">
                                                    <input type="color" class="color-picker-input"
                                                           :value="'#' + element" @input="onColorInput(index, $event)" />
                                                    <div class="color-preview-box" :style="{ backgroundColor: '#' + element }"></div>
                                                </div>
                                                <div class="ms-3 hextual-color">
                                                    <span class="always-ltr">{{ '#' + element.toLowerCase() }}</span>
                                                </div>
                                            </div>

                                            <v-spacer/>

                                            <template v-if="hoveredIndex === index">
                                                <v-btn class="px-2" color="default"
                                                       density="comfortable" variant="text"
                                                       :prepend-icon="mdiDeleteOutline"
                                                       @click="removeChartColor(index)">
                                                    {{ tt('Delete') }}
                                                </v-btn>
                                            </template>

                                            <span class="ms-1">
                                            <v-icon class="drag-handle" :icon="mdiDrag"/>
                                            <v-tooltip activator="parent">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                        </div>
                                    </td>
                                </tr>
                            </template>
                        </draggable-list>
                    </v-table>
                    <div class="w-100 h-100" v-if="showRawData">
                        <v-textarea no-resize class="w-100 h-100 ps-4 code-textarea hextual-color always-cursor-text"
                                    density="compact" variant="plain" :rounded="false"
                                    :placeholder="tt('Each line should be a hex color value (e.g. c67e48 or #c67e48)')"
                                    v-model="textualChartColors"></v-textarea>
                    </div>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { VTable } from 'vuetify/components/VTable';
import { ref, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useChartColorSchemeSettingsPageBase } from '@/views/base/settings/ChartColorSchemeSettingsPageBase.ts';

import {
    mdiPlus,
    mdiCheck,
    mdiDotsVertical,
    mdiRestore,
    mdiDeleteOutline,
    mdiDrag
} from '@mdi/js';

const { tt } = useI18n();

const {
    chartColors,
    chartColorsModified,
    textualChartColors,
    canSaveColorScheme,
    addChartColor,
    removeChartColor,
    resetChartColorsToDefault,
    loadChartColorsFromSettings,
    saveChartColorsToSettings,
    onColorInput
} = useChartColorSchemeSettingsPageBase();

const colorSchemeTable = useTemplateRef<VTable>('colorSchemeTable');

let resolveFunc: (() => void) | null = null;
let rejectFunc: (() => void) | null = null;

const showState = ref<boolean>(false);
const showRawData = ref<boolean>(false);
const hoveredIndex = ref<string>('');

function open(): Promise<void> {
    loadChartColorsFromSettings();
    showState.value = true;
    showRawData.value = false;

    return new Promise<void>((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function addNewColor(): void {
    addChartColor();

    nextTick(() => {
        colorSchemeTable.value?.$el?.querySelector('tbody > tr:last-child')?.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        });
    });
}

function saveChartColors(): void {
    saveChartColorsToSettings();
    resolveFunc?.();
    showState.value = false;
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>

<style>
.color-preview-wrapper {
    position: relative;
    width: 32px;
    height: 32px;
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
    width: 32px;
    height: 32px;
    border-radius: 4px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    pointer-events: none;
}

.hextual-color {
    font-family: monospace;
}
</style>
