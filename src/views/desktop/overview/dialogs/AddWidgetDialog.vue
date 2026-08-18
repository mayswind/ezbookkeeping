<template>
    <v-dialog width="600" v-model="showState">
        <one-column-dialog-layout content-class="px-2 py-0"
                                  :title="tt('Add Widget')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #content>
                <v-list lines="two">
                    <v-list-item :key="definition.type" :title="tt(definition.name)"
                                 :subtitle="`${definition.defaultWidth} × ${definition.defaultHeight}`"
                                 v-for="definition in DESKTOP_OVERVIEW_WIDGET_DEFINITIONS"
                                 @click="select(definition.type)">
                        <template #append>
                            <v-icon :icon="mdiPlus" />
                        </template>
                    </v-list-item>
                </v-list>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { type OverviewWidgetType } from '@/core/overview_layout.ts';
import { DESKTOP_OVERVIEW_WIDGET_DEFINITIONS } from '@/consts/overview_layout.ts';

import {
    mdiPlus
} from '@mdi/js';

const { tt } = useI18n();

let resolveFunc: ((type: OverviewWidgetType) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);

function open(): Promise<OverviewWidgetType> {
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function select(type: OverviewWidgetType): void {
    resolveFunc?.(type);
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
