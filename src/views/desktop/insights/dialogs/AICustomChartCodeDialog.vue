<template>
    <v-dialog width="600" :persistent="generating || !!userPrompt" v-model="showState">
        <one-column-dialog-layout content-class="pa-0"
                                  :disabled="generating" :loading="generating"
                                  :title="tt('AI Generate Chart Code')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="!userPrompt || !userPrompt.trim() || generating"
                       @click="generate" v-if="!generating">{{ tt('Generate') }}</v-btn>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       @click="cancelGenerate" v-if="generating && cancelGeneratingUuid">{{ tt('Cancel Generation') }}</v-btn>
            </template>

            <template #content>
                <v-textarea no-resize persistent-placeholder
                            class="w-100 h-100 ps-4 always-cursor-text"
                            rows="10" autocomplete="off" density="compact" variant="plain" :rounded="false"
                            :disabled="generating"
                            :placeholder="tt('Describe the chart you want')"
                            v-model="userPrompt"></v-textarea>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useExplorersStore } from '@/stores/explorer.ts';

import { generateRandomUUID } from '@/lib/misc.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const explorersStore = useExplorersStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: ((code: string) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const generating = ref<boolean>(false);
const cancelGeneratingUuid = ref<string | undefined>(undefined);
const userPrompt = ref<string>('');
const currentCode = ref<string>('');

function open(code: string): Promise<string> {
    userPrompt.value = '';
    currentCode.value = code;
    showState.value = true;

    return new Promise<string>((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function generate(): void {
    if (generating.value || !userPrompt.value || !userPrompt.value.trim()) {
        return;
    }

    cancelGeneratingUuid.value = generateRandomUUID();
    generating.value = true;

    explorersStore.generateCustomChartCode({
        userPrompt: userPrompt.value.trim(),
        currentCode: currentCode.value,
        cancelableUuid: cancelGeneratingUuid.value
    }).then(response => {
        resolveFunc?.(response.code);
        showState.value = false;
        generating.value = false;
        cancelGeneratingUuid.value = undefined;
    }).catch(error => {
        if (error.canceled) {
            return;
        }

        generating.value = false;
        cancelGeneratingUuid.value = undefined;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancelGenerate(): void {
    if (!cancelGeneratingUuid.value) {
        return;
    }

    explorersStore.cancelGenerateCustomChartCode(cancelGeneratingUuid.value);
    generating.value = false;
    cancelGeneratingUuid.value = undefined;

    snackbar.value?.showMessage('User Canceled');
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
    generating.value = false;
    cancelGeneratingUuid.value = undefined;
}

defineExpose({
    open
});
</script>
