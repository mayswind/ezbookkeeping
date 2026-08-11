<template>
    <v-dialog width="1000" v-model="showState">
        <one-column-dialog-layout content-class="pa-0" content-style="height: 500px"
                                  :title="tt('Export Queries')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <div ref="buttonContainer">
                    <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true"
                           :disabled="!queriesJson" @click="copy">
                        <v-icon :icon="mdiContentCopy" size="20" />
                        <v-tooltip activator="parent">{{ tt('Copy') }}</v-tooltip>
                    </v-btn>
                    <v-btn density="compact" color="default" variant="text" class="ms-1" :icon="true"
                           :disabled="!queriesJson" @click="save()">
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
                                :value="queriesJson"></v-textarea>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { KnownFileType } from '@/core/file.ts';

import { copyTextToClipboard, startDownloadFile } from '@/lib/ui/common.ts';

import {
    mdiContentCopy,
    mdiContentSaveOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const userStore = useUserStore();

const buttonContainer = useTemplateRef<HTMLElement>('buttonContainer');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const queriesJson = ref<string>('');

const fileName = computed<string>(() => {
    const nickname = userStore.currentUserNickname;

    if (nickname) {
        return tt('dataExport.insightsExplorerQueryFileName', {
            nickname: nickname
        });
    }

    return tt('dataExport.defaultInsightsExplorerQueryFileName');
});

function open(options: { queriesJson: string }): void {
    queriesJson.value = options.queriesJson;
    showState.value = true;
}

function copy(): void {
    copyTextToClipboard(queriesJson.value, buttonContainer.value);
    snackbar.value?.showMessage('Data copied');
}

function save(): void {
    const fileType = KnownFileType.JSON;
    startDownloadFile(fileType.formatFileName(fileName.value), fileType.createBlob(queriesJson.value));
}

function cancel(): void {
    showState.value = false;
}

defineExpose({
    open
});
</script>
