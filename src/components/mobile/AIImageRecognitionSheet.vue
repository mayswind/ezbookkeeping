<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :class="{ 'disabled': loading || recognizing }" :text="tt('Cancel')" @click="cancel"></f7-link>
            </div>
            <div class="right">
                <f7-link :class="{ 'disabled': loading || recognizing || !imageFile }" :text="tt('Recognize')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="no-margin-vertical no-padding-vertical">
            <div class="image-container display-flex justify-content-center"
                 style="height: 400px" @click="showOpenImage">
                <img height="400px" :src="imageSrc" v-if="imageSrc" />
                <div class="image-container-background display-flex justify-content-center align-items-center padding-horizontal" v-if="!imageSrc">
                    <span v-if=!loading>{{ tt('Click here to select a receipt or transaction image') }}</span>
                    <span v-else-if="loading">{{ tt('Loading image...') }}</span>
                </div>
            </div>
        </f7-page-content>

        <input ref="imageInput" type="file" style="display: none" :accept="`${SUPPORTED_IMAGE_EXTENSIONS};capture=camera`" @change="openImage($event)" />
    </f7-sheet>
</template>


<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, closeAllDialog } from '@/lib/ui/mobile.ts';

import { useTransactionsStore } from '@/stores/transaction.ts';

import { KnownFileType } from '@/core/file.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import type { RecognizedReceiptImageResponse } from '@/models/large_language_model.ts';

import { generateRandomUUID } from '@/lib/misc.ts';
import { compressJpgImage } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

defineProps<{
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'recognition:change', value: RecognizedReceiptImageResponse): void;
}>();

const { tt } = useI18n();
const { showCancelableLoading, showToast } = useI18nUIComponents();

const transactionsStore = useTransactionsStore();

const imageInput = useTemplateRef<HTMLInputElement>('imageInput');

const loading = ref<boolean>(false);
const recognizing = ref<boolean>(false);
const cancelRecognizingUuid = ref<string | undefined>(undefined);
const imageFile = ref<File | null>(null);
const imageSrc = ref<string | undefined>(undefined);

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
        showToast('Unable to load image');
    });
}

function showOpenImage(): void {
    if (loading.value || recognizing.value) {
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

function confirm(): void {
    if (loading.value || recognizing.value || !imageFile.value) {
        return;
    }

    cancelRecognizingUuid.value = generateRandomUUID();
    recognizing.value = true;
    showCancelableLoading('Recognizing', 'AI can make mistakes. Check important info.', 'Cancel Recognition', cancelRecognize);

    transactionsStore.recognizeReceiptImage({
        imageFile: imageFile.value,
        cancelableUuid: cancelRecognizingUuid.value
    }).then(response => {
        recognizing.value = false;
        cancelRecognizingUuid.value = undefined;
        closeAllDialog();
        emit('update:show', false);
        emit('recognition:change', response);
    }).catch(error => {
        if (error.canceled) {
            return;
        }

        recognizing.value = false;
        cancelRecognizingUuid.value = undefined;
        closeAllDialog();

        if (!error.processed) {
            showToast(error.message || error);
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
    closeAllDialog();

    showToast('User Canceled');
}

function cancel(): void {
    close();
}

function close(): void {
    emit('update:show', false);
    loading.value = false;
    recognizing.value = false;
    cancelRecognizingUuid.value = undefined;
    imageFile.value = null;
    imageSrc.value = undefined;
}

function onSheetOpen(): void {
    loading.value = false;
    recognizing.value = false;
    cancelRecognizingUuid.value = undefined;
    imageFile.value = null;
    imageSrc.value = undefined;
}

function onSheetClosed(): void {
    close();
}
</script>

<style>
.image-container {
    border: 1px solid var(--f7-page-master-border-color);
}

.image-container-background {
    width: 100%;
    height: 100%;

    > span {
        font-size: var(--f7-input-font-size);
    }
}
</style>
