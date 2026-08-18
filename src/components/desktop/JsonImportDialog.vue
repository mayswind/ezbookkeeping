<template>
    <v-dialog width="1000" :persistent="!!json && json !== sampleJson" v-model="showState">
        <one-column-dialog-layout content-class="pa-0" content-style="height: 500px"
                                  :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="!json || !onImport" @click="confirm">{{ tt('Import') }}</v-btn>
            </template>

            <template #content>
                <div class="w-100 h-100">
                    <v-textarea no-resize class="w-100 h-100 ps-3 code-textarea always-cursor-text"
                                density="compact" variant="plain" :rounded="false"
                                v-model="json"></v-textarea>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

const props = defineProps<{
    title: string;
    sampleJson?: string;
    onImport: (value: string) => boolean;
}>();

const { tt } = useI18n();

let resolveFunc: ((value: unknown) => void) | null = null;
let rejectFunc: (() => void) | null = null;

const showState = ref<boolean>(false);
const json = ref<string>('');

function open(): Promise<unknown> {
    json.value = props.sampleJson || '';
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!json.value || !props.onImport) {
        return;
    }

    const result = props.onImport(json.value);

    if (result) {
        resolveFunc?.(result);
        showState.value = false;
    }
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
