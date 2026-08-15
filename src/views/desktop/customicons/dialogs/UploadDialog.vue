<template>
    <v-dialog width="600" :persistent="uploading || !!imageSrc" v-model="showState" @paste="onPaste">
        <one-column-dialog-layout content-class="pa-0" content-style="height: 420px"
                                  :disabled="uploading" :loading="uploading"
                                  :title="tt('Upload Custom Icon')"
                                  :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="ms-2 me-1" density="comfortable" variant="outlined" color="secondary"
                       :disabled="uploading" @click="showOpenImageDialog" v-if="imageSrc">
                    {{ tt('Select Another Image') }}
                </v-btn>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="!imageSrc || uploading" @click="upload">
                    {{ tt('Upload') }}
                </v-btn>
            </template>

            <template #content>
                <div class="w-100 h-100 border position-relative cursor-pointer"
                     @dragenter.prevent="onDragEnter"
                     @dragover.prevent
                     @dragleave.prevent="onDragLeave"
                     @drop.prevent="onDrop"
                     @click="showOpenImageDialog"
                     v-if="!imageSrc">
                    <div class="d-flex w-100 h-100 justify-center align-center justify-content-center text-center px-4"
                         :class="{ 'dropzone': true, 'dropzone-dark': isDarkMode, 'dropzone-blurry-bg': isDragOver, 'dropzone-dragover': isDragOver }">
                        <div class="d-inline-flex flex-column" v-if="!isDragOver">
                            <span class="text-title-medium font-weight-bold pa-2">{{ tt('You can drag and drop, paste or click to select an icon') }}</span>
                        </div>
                        <span class="text-title-medium font-weight-bold pa-2" v-else-if="isDragOver">{{ tt('Release to load image') }}</span>
                    </div>
                </div>
                <div class="d-flex w-100 h-100 flex-column justify-center align-center" v-else-if="imageSrc">
                    <div ref="cropStage" class="crop-stage" :class="{ panning: isPanning }"
                         @wheel.prevent="onWheel" @pointerdown.prevent="onPointerDown"
                         @pointermove.prevent="onPointerMove" @pointerup="onPointerUp" @pointercancel="onPointerUp">
                        <img ref="previewImage" :src="imageSrc" alt="" draggable="false"
                             :style="previewStyle" @load="onImageLoaded" />
                        <div class="crop-grid"></div>
                    </div>
                    <div class="text-body-large mt-4 mx-4">{{ tt('Use the mouse wheel to zoom and drag the image to adjust the crop') }}</div>
                </div>
            </template>
        </one-column-dialog-layout>

        <input ref="imageInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="openImage($event)" />
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import OneColumnDialogLayout from '@/components/desktop/OneColumnDialogLayout.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useUserCustomIconsStore } from '@/stores/userCustomIcon.ts';

import { KnownFileType } from '@/core/file.ts';
import { ThemeType } from '@/core/theme.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import { generateRandomUUID } from '@/lib/misc.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

interface CustomIconUploadResponse {
    message: string;
}

const theme = useTheme();

const { tt } = useI18n();

const customIconsStore = useUserCustomIconsStore();

const CROP_STAGE_SIZE: number = 320;
const MIN_ZOOM: number = 1;
const MAX_ZOOM: number = 5;

const imageInput = useTemplateRef<HTMLInputElement>('imageInput');
const previewImage = useTemplateRef<HTMLImageElement>('previewImage');
const cropStage = useTemplateRef<HTMLElement>('cropStage');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: ((response?: CustomIconUploadResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const uploading = ref<boolean>(false);
const isDragOver = ref<boolean>(false);
const clientSessionId = ref<string>('');
const imageSrc = ref<string | undefined>(undefined);
const zoom = ref<number>(1);
const imageWidth = ref<number>(0);
const imageHeight = ref<number>(0);
const offsetX = ref<number>(0);
const offsetY = ref<number>(0);
const isPanning = ref<boolean>(false);
const activePointerId = ref<number>();
const lastPointerX = ref<number>(0);
const lastPointerY = ref<number>(0);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const baseScale = computed(() => imageWidth.value && imageHeight.value ? Math.max(CROP_STAGE_SIZE / imageWidth.value, CROP_STAGE_SIZE / imageHeight.value) : 1);
const previewStyle = computed(() => ({
    width: `${imageWidth.value * baseScale.value * zoom.value}px`,
    height: `${imageHeight.value * baseScale.value * zoom.value}px`,
    transform: `translate(calc(-50% + ${offsetX.value}px), calc(-50% + ${offsetY.value}px))`
}));

function loadImage(file: File): void {
    imageSrc.value = URL.createObjectURL(file);
    resetCropState();
}

function resetCropState(): void {
    zoom.value = 1;
    imageWidth.value = 0;
    imageHeight.value = 0;
    offsetX.value = 0;
    offsetY.value = 0;
}

function clampOffsets(zoomValue: number = zoom.value): void {
    const displayedWidth = imageWidth.value * baseScale.value * zoomValue;
    const displayedHeight = imageHeight.value * baseScale.value * zoomValue;
    const maxX = Math.max(0, (displayedWidth - CROP_STAGE_SIZE) / 2);
    const maxY = Math.max(0, (displayedHeight - CROP_STAGE_SIZE) / 2);

    offsetX.value = Math.max(-maxX, Math.min(maxX, offsetX.value));
    offsetY.value = Math.max(-maxY, Math.min(maxY, offsetY.value));
}

function open(): Promise<CustomIconUploadResponse | undefined> {
    showState.value = true;
    clientSessionId.value = generateRandomUUID();
    imageSrc.value = undefined;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function showOpenImageDialog(): void {
    if (uploading.value || isDragOver.value) {
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

function upload(): void {
    if (uploading.value || !imageSrc.value) {
        return;
    }

    const image = previewImage.value;

    if (!image || !image.naturalWidth || !image.naturalHeight) {
        return;
    }

    uploading.value = true;

    const renderedScale = baseScale.value * zoom.value;
    const cropSize = CROP_STAGE_SIZE / renderedScale;
    const centerX = image.naturalWidth / 2 - offsetX.value / renderedScale;
    const centerY = image.naturalHeight / 2 - offsetY.value / renderedScale;
    const sx = Math.max(0, Math.min(image.naturalWidth - cropSize, centerX - cropSize / 2));
    const sy = Math.max(0, Math.min(image.naturalHeight - cropSize, centerY - cropSize / 2));
    const canvas = document.createElement('canvas');
    canvas.width = 256;
    canvas.height = 256;
    canvas.getContext('2d')?.drawImage(image, sx, sy, cropSize, cropSize, 0, 0, 256, 256);

    canvas.toBlob((blob) => {
        if (!blob) {
            uploading.value = false;
            snackbar.value?.showError('Unable to crop image');
            return;
        }

        customIconsStore.uploadCustomIcon({
            iconFile: KnownFileType.PNG.createFileFromBlob(blob, 'icon'),
            clientSessionId: clientSessionId.value
        }).then(() => {
            resolveFunc?.({ message: 'You have added a new custom icon' });
            uploading.value = false;
            showState.value = false;
        }).catch(error => {
            uploading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }, KnownFileType.PNG.contentType);
}

function cancel(): void {
    if (uploading.value) {
        return;
    }

    rejectFunc?.();
    showState.value = false;
    imageSrc.value = undefined;
}

function onDragEnter(): void {
    isDragOver.value = true;
}

function onDragLeave(): void {
    isDragOver.value = false;
}

function onDrop(event: DragEvent): void {
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

function onImageLoaded(): void {
    const image = previewImage.value;

    if (!image) {
        return;
    }

    imageWidth.value = image.naturalWidth;
    imageHeight.value = image.naturalHeight;
    zoom.value = 1;
    offsetX.value = 0;
    offsetY.value = 0;
}

function onWheel(event: WheelEvent): void {
    if (!imageWidth.value || !imageHeight.value) {
        return;
    }

    const previousZoom = zoom.value;
    const nextZoom = Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, previousZoom * (event.deltaY < 0 ? 1.1 : 1 / 1.1)));

    if (nextZoom === previousZoom) {
        return;
    }

    const bounds = cropStage.value?.getBoundingClientRect();
    const pointerX = bounds ? event.clientX - bounds.left - CROP_STAGE_SIZE / 2 : 0;
    const pointerY = bounds ? event.clientY - bounds.top - CROP_STAGE_SIZE / 2 : 0;
    const zoomRatio = nextZoom / previousZoom;

    offsetX.value = pointerX - (pointerX - offsetX.value) * zoomRatio;
    offsetY.value = pointerY - (pointerY - offsetY.value) * zoomRatio;
    zoom.value = nextZoom;
    clampOffsets(nextZoom);
}

function onPointerMove(event: PointerEvent): void {
    if (!isPanning.value || activePointerId.value !== event.pointerId) {
        return;
    }

    offsetX.value += event.clientX - lastPointerX.value;
    offsetY.value += event.clientY - lastPointerY.value;
    lastPointerX.value = event.clientX;
    lastPointerY.value = event.clientY;
    clampOffsets();
}

function onPointerUp(event: PointerEvent): void {
    if (activePointerId.value !== event.pointerId) {
        return;
    }

    isPanning.value = false;
    activePointerId.value = undefined;
    const target = event.currentTarget as HTMLElement;
    if (target.hasPointerCapture(event.pointerId)) {
        target.releasePointerCapture(event.pointerId);
    }
}

function onPointerDown(event: PointerEvent): void {
    if (!imageWidth.value) {
        return;
    }

    activePointerId.value = event.pointerId;
    lastPointerX.value = event.clientX;
    lastPointerY.value = event.clientY;
    isPanning.value = true;
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
}

defineExpose({
    open
});
</script>

<style>
.crop-stage {
    width: 320px;
    height: 320px;
    position: relative;
    overflow: hidden;
    cursor: grab;
    touch-action: none;
    user-select: none;
    background: repeating-conic-gradient(#ddd 0 25%, #fff 0 50%) 50% / 20px 20px;
}

.crop-stage.panning {
    cursor: grabbing;
}

.crop-stage img {
    max-width: none;
    position: absolute;
    left: 50%;
    top: 50%;
    pointer-events: none;
}

.crop-grid {
    position: absolute;
    inset: 0;
    pointer-events: none;
    border: 2px solid rgb(var(--v-theme-primary));
    background: linear-gradient(90deg, transparent 33%, rgba(255, 255, 255, 0.7) 33% 33.5%, transparent 33.5% 66%, rgba(255, 255, 255, 0.7) 66% 66.5%, transparent 66.5%), linear-gradient(transparent 33%, rgba(255, 255, 255, 0.7) 33% 33.5%, transparent 33.5% 66%, rgba(255, 255, 255, 0.7) 66% 66.5%, transparent 66.5%);
}
</style>
