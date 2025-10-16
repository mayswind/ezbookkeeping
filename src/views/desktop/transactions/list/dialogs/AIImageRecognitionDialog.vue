<template>
    <v-dialog width="800" :persistent="loading || recognizing || !!imageFile" v-model="showState" @paste="onPaste">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ tt('AI Image Recognition') }}</h4>
                </div>
            </template>

            <v-card-text class="d-flex justify-center w-100 my-md-4 pt-0">
                <div class="w-100 border position-relative"
                     @dragenter.prevent="onDragEnter"
                     @dragover.prevent
                     @dragleave.prevent="onDragLeave"
                     @drop.prevent="onDrop">
                    <div class="d-flex w-100 fill-height justify-center align-center justify-content-center px-4"
                         :class="{ 'dropzone': true, 'dropzone-dragover': isDragOver }" style="height: 480px">
                        <h3 :class="{ 'pa-2': true, 'bg-grey-200': !isDarkMode, 'bg-grey-100': isDarkMode }" v-if="!loading && !imageFile && !isDragOver">{{ tt('You can drag and drop, paste or click to select a receipt or transaction image') }}</h3>
                        <h3 :class="{ 'pa-2': true, 'bg-grey-200': !isDarkMode, 'bg-grey-100': isDarkMode }" v-else-if="!loading && isDragOver">{{ tt('Release to load image') }}</h3>
                        <h3 class="pa-2" v-else-if="loading">{{ tt('Loading image...') }}</h3>
                        <h3 :class="{ 'pa-2': true, 'bg-grey-200': !isDarkMode, 'bg-grey-100': isDarkMode }" v-else-if="recognizing">{{ tt('AI can make mistakes. Check important info.') }}</h3>
                    </div>
                    <v-img height="480px" :class="{ 'cursor-pointer': !loading && !recognizing && !isDragOver }"
                           :src="imageSrc" @click="showOpenImageDialog">
                        <template #placeholder>
                            <div :class="{ 'w-100 fill-height': true, 'bg-grey-200': !isDarkMode, 'bg-grey-100': isDarkMode }"></div>
                        </template>
                    </v-img>
                </div>
            </v-card-text>

            <v-card-text class="overflow-y-visible">
                <div ref="buttonContainer" class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="loading || recognizing || !imageFile" @click="recognize">
                        {{ tt('Recognize') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="recognizing"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading"
                           @click="cancelRecognize" v-if="recognizing && cancelRecognizingUuid">{{ tt('Cancel Recognition') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading || recognizing"
                           @click="cancel" v-if="!recognizing || !cancelRecognizingUuid">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
    <input ref="imageInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="openImage($event)" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionsStore } from '@/stores/transaction.ts';

import { KnownFileType } from '@/core/file.ts';
import { ThemeType } from '@/core/theme.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import type { RecognizedReceiptImageResponse } from '@/models/large_language_model.ts';

import { generateRandomUUID } from '@/lib/misc.ts';
import { compressJpgImage } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const theme = useTheme();

const { tt } = useI18n();

const transactionsStore = useTransactionsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const imageInput = useTemplateRef<HTMLInputElement>('imageInput');

let resolveFunc: ((response: RecognizedReceiptImageResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const recognizing = ref<boolean>(false);
const cancelRecognizingUuid = ref<string | undefined>(undefined);
const imageFile = ref<File | null>(null);
const imageSrc = ref<string | undefined>(undefined);
const isDragOver = ref<boolean>(false);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

function loadImage(file: File): void {
    loading.value = true;
    imageFile.value = null;
    imageSrc.value = undefined;

    compressJpgImage(file, 1280, 1280, 0.8).then(blob => {
        imageFile.value = KnownFileType.JPG.createFileFromBlob(blob, "image");
        imageSrc.value = URL.createObjectURL(blob);
        loading.value = false;
    }).catch(error => {
        imageFile.value = null;
        imageSrc.value = undefined;
        loading.value = false;
        logger.error('failed to compress image', error);
        snackbar.value?.showError('Unable to load image');
    });
}

function open(): Promise<RecognizedReceiptImageResponse> {
    showState.value = true;
    loading.value = false;
    recognizing.value = false;
    cancelRecognizingUuid.value = undefined;
    imageFile.value = null;
    imageSrc.value = undefined;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function showOpenImageDialog(): void {
    if (loading.value || recognizing.value || isDragOver.value) {
        return;
    }

    imageInput.value?.click();
}

function openImage(event: Event): void {
    if (!event || !event.target) {
        return;
    }

    const el = event.target as HTMLInputElement;

    if (!el.files || !el.files.length || !el.files[0]) {
        return;
    }

    const image = el.files[0] as File;

    el.value = '';

    loadImage(image);
}

function recognize(): void {
    if (loading.value || recognizing.value || !imageFile.value) {
        return;
    }

    cancelRecognizingUuid.value = generateRandomUUID();
    recognizing.value = true;

    transactionsStore.recognizeReceiptImage({
        imageFile: imageFile.value,
        cancelableUuid: cancelRecognizingUuid.value
    }).then(response => {
        resolveFunc?.(response);
        showState.value = false;
        recognizing.value = false;
        cancelRecognizingUuid.value = undefined;
    }).catch(error => {
        if (error.canceled) {
            return;
        }

        recognizing.value = false;
        cancelRecognizingUuid.value = undefined;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancelRecognize(): void {
    if (!cancelRecognizingUuid.value) {
        return;
    }

    transactionsStore.cancelRecognizeReceiptImage(cancelRecognizingUuid.value);
    recognizing.value = false;
    cancelRecognizingUuid.value = undefined;

    snackbar.value?.showMessage('User Canceled');
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
    loading.value = false;
    recognizing.value = false;
    cancelRecognizingUuid.value = undefined;
    imageFile.value = null;
    imageSrc.value = undefined;
}

function onDragEnter(): void {
    if (loading.value || recognizing.value) {
        return;
    }

    isDragOver.value = true;
}

function onDragLeave(): void {
    isDragOver.value = false;
}

function onDrop(event: DragEvent): void {
    if (loading.value || recognizing.value) {
        return;
    }

    isDragOver.value = false;

    if (event.dataTransfer && event.dataTransfer.files && event.dataTransfer.files.length && event.dataTransfer.files[0]) {
        loadImage(event.dataTransfer.files[0] as File);
    }
}

function onPaste(event: ClipboardEvent) {
    if (!event.clipboardData) {
        event.preventDefault();
        return;
    }

    for (let i = 0; i < event.clipboardData.items.length; i++) {
        const item = event.clipboardData.items[i];

        if (item && item.type.startsWith('image/')) {
            const file = item.getAsFile();

            if (file) {
                loadImage(file);
                event.preventDefault();
                return;
            }
        }
    }
}

defineExpose({
    open
});
</script>

<style>
.dropzone {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    pointer-events: none;
    border-radius: 8px;
    z-index: 10;
}

.dropzone-dragover {
    border: 6px dashed rgba(var(--v-border-color),var(--v-border-opacity));
}
</style>
