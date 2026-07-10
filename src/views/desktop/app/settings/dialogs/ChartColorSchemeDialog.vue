<template>
    <v-dialog width="700" :persistent="chartColorsModified" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex align-center">
                        <h4 class="text-h4">{{ tt('Chart Color Scheme') }}</h4>
                        <v-btn class="ms-3" color="default" variant="outlined" density="comfortable"
                               @click="addNewColor()">{{ tt('Add') }}</v-btn>
                    </div>
                    <v-spacer/>
                    <v-switch class="bidirectional-switch ms-2 pt-1" color="secondary"
                              :label="tt('Raw Data')"
                              v-model="showRawData"
                              @click="showRawData = !showRawData">
                        <template #prepend>
                            <span>{{ tt('List') }}</span>
                        </template>
                    </v-switch>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiRestore"
                                             :title="tt('Reset to Default')"
                                             @click="resetChartColorsToDefault"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto" style="height: 440px">
                <v-table ref="colorSchemeTable" hover density="comfortable" class="w-100 table-striped" v-if="!showRawData">
                    <draggable-list tag="tbody"
                                    item-key="index"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    v-model="chartColors">
                        <template #item="{ element, index }">
                            <tr class="text-sm" @mouseenter="hoveredIndex = index" @mouseleave="hoveredIndex = ''">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <div class="color-preview-wrapper">
                                                <input type="color" class="color-picker-input"
                                                       :value="'#' + element" @input="onColorInput(index, $event)" />
                                                <div class="color-preview-box" :style="{ backgroundColor: '#' + element }"></div>
                                            </div>
                                            <span class="ms-3 hextual-color">{{ '#' + element.toLowerCase() }}</span>
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
                    <v-textarea class="w-100 h-100 hextual-color always-cursor-text"
                                :placeholder="tt('Each line should be a hex color value (e.g. c67e48 or #c67e48)')"
                                v-model="textualChartColors"></v-textarea>
                </div>
            </v-card-text>

            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="!canSaveColorScheme" @click="saveChartColors">{{ tt('Save') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { VTable } from 'vuetify/components/VTable';
import { ref, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useChartColorSchemeSettingsPageBase } from '@/views/base/settings/ChartColorSchemeSettingsPageBase.ts';

import {
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
