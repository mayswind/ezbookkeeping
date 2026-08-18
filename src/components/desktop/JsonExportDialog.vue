<template>
    <v-dialog width="1000" v-model="showState">
        <one-column-dialog-layout content-class="pa-0" content-style="height: 500px"
                                  :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <div ref="buttonContainer">
                    <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true"
                           :disabled="!json" @click="copy">
                        <v-icon :icon="mdiContentCopy" size="20" />
                        <v-tooltip activator="parent">{{ tt('Copy') }}</v-tooltip>
                    </v-btn>
                    <v-btn density="compact" color="default" variant="text" class="ms-1" :icon="true"
                           :disabled="!json" @click="save">
                        <v-icon :icon="mdiContentSaveOutline" size="22" />
                        <v-tooltip activator="parent">{{ tt('Save') }}</v-tooltip>
                    </v-btn>
                </div>
            </template>

            <template #content>
                <div class="w-100 h-100 code-container">
                    <v-textarea no-resize class="w-100 h-100 ps-3 always-cursor-text"
                                density="compact" variant="plain"
                                :readonly="true" :rounded="false"
                                :value="json"></v-textarea>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { KnownFileType } from '@/core/file.ts';

import { copyTextToClipboard, startDownloadFile } from '@/lib/ui/common.ts';

import {
    mdiContentCopy,
    mdiContentSaveOutline
} from '@mdi/js';

export interface JsonExportDialogOptions {
    json: string;
}

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    title: string;
    fileName: string;
}>();

const { tt } = useI18n();

const buttonContainer = useTemplateRef<HTMLElement>('buttonContainer');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const json = ref<string>('');

function open(options: JsonExportDialogOptions): void {
    json.value = options.json;
    showState.value = true;
}

function copy(): void {
    copyTextToClipboard(json.value, buttonContainer.value);
    snackbar.value?.showMessage('Data copied');
}

function save(): void {
    const fileType = KnownFileType.JSON;
    startDownloadFile(fileType.formatFileName(props.fileName), fileType.createBlob(json.value));
}

function cancel(): void {
    showState.value = false;
}

defineExpose({
    open
});
</script>
