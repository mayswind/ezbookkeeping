<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :class="{ 'disabled': loading || recognizing }" :text="tt('Choose from Library')" @click="showOpenImage"></f7-link>
            </div>
            <div class="right">
                <f7-link :class="{ 'disabled': loading || recognizing }" :text="tt('Take Photo')" @click="showCamera"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="margin-top no-padding-top">
            <div class="padding-horizontal padding-bottom">
                <div class="image-container display-flex justify-content-center width-100 margin-bottom" style="height: 240px">
                    <img height="240px" :src="imageSrc" v-if="imageSrc" />
                    <div class="image-container-background display-flex justify-content-center align-items-center" v-if="!imageSrc">
                        <span>{{ tt('Please select a receipt or transaction image first') }}</span>
                    </div>
                </div>
                <f7-button large fill color="primary"
                           :class="{ 'disabled': loading || recognizing || !imageFile }"
                           :text="tt('Recognize')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link :class="{ 'disabled': loading || recognizing }" @click="cancel" :text="tt('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>

        <input ref="imageInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="openImage($event)" />
        <input ref="cameraInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" capture="environment" @change="openImage($event)" />
    </f7-sheet>
</template>


<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import { useTransactionsStore } from '@/stores/transaction.ts';

import { KnownFileType } from '@/core/file.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import type { RecognizedReceiptImageResponse } from '@/models/large_language_model.ts';

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
const { showToast } = useI18nUIComponents();

const transactionsStore = useTransactionsStore();

const imageInput = useTemplateRef<HTMLInputElement>('imageInput');
const cameraInput = useTemplateRef<HTMLInputElement>('cameraInput');

const loading = ref<boolean>(false);
const recognizing = ref<boolean>(false);
const imageFile = ref<File | null>(null);
const imageSrc = ref<string | undefined>(undefined);

function loadImage(file: File): void {
    compressJpgImage(file, 1280, 1280, 0.8).then(blob => {
        imageFile.value = KnownFileType.JPG.createFileFromBlob(blob, "image");
        imageSrc.value = URL.createObjectURL(blob);
    }).catch(error => {
        imageFile.value = null;
        imageSrc.value = undefined;
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

function showCamera(): void {
    if (loading.value || recognizing.value) {
        return;
    }

    cameraInput.value?.click();
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

    recognizing.value = true;
    showLoading(() => recognizing.value);

    transactionsStore.recognizeReceiptImage({
        imageFile: imageFile.value
    }).then(response => {
        recognizing.value = false;
        hideLoading();
        emit('update:show', false);
        emit('recognition:change', response);
    }).catch(error => {
        recognizing.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function cancel(): void {
    close();
}

function close(): void {
    emit('update:show', false);
    loading.value = false;
    recognizing.value = false;
    imageFile.value = null;
    imageSrc.value = undefined;
}

function onSheetOpen(): void {
    loading.value = false;
    recognizing.value = false;
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
    background-color: var(--f7-page-bg-color);
}
</style>
