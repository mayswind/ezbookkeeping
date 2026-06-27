<template>
    <v-dialog width="1000" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4">{{ tt('Export Queries') }}</h4>
            </template>

            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto" style="height: 485px">
                <div class="w-100 h-100 code-container">
                    <v-textarea class="w-100 h-100 always-cursor-text" :readonly="true" :value="queriesJson"></v-textarea>
                </div>
            </v-card-text>

            <v-card-text>
                <div ref="buttonContainer" class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn-group variant="tonal" density="comfortable">
                        <v-btn color="primary" :disabled="!queriesJson" @click="copy">{{ tt('Copy') }}</v-btn>
                        <v-btn density="compact" color="primary" :disabled="!queriesJson" :icon="true">
                            <v-icon :icon="mdiMenuDown" size="24" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :title="tt('Save')" @click="save()"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </v-btn-group>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
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
    mdiMenuDown
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
