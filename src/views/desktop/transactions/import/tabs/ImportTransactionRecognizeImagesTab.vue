<template>
    <div>
        <div class="text-center">{{ progressText }}</div>

        <v-progress-linear class="my-4" rounded color="primary" height="20"
                           :striped="!!submitting" :model-value="processedPercent" :buffer-value="progressingPercent" />

        <v-list class="recognition-failed-list" v-if="importProgress.failedCount > 0">
            <template :key="index" v-for="(item, index) in importImages">
                <v-list-item class="px-0" v-if="item.status === 'failed'">
                    <template #prepend>
                        <v-img :src="item.previewUrl" cover width="48" height="48" class="rounded me-3" />
                    </template>
                    <v-list-item-title>{{ item.file.name }}</v-list-item-title>
                    <v-list-item-subtitle class="text-error" v-if="item.failureReason">{{ item.failureReason }}</v-list-item-subtitle>
                </v-list-item>
            </template>
        </v-list>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

export type ImageStatus = 'pending' | 'recognizing' | 'success' | 'failed';

export interface BatchImportImageItem {
    file: File;
    previewUrl: string;
    status: ImageStatus;
    failureReason?: string;
}

const props = defineProps<{
    importImages: BatchImportImageItem[];
    disabled?: boolean;
    submitting?: boolean;
}>();

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const importProgress = computed<{ successCount: number, failedCount: number }>(() => {
    let successCount: number = 0;
    let failedCount: number = 0;

    for (const item of props.importImages) {
        if (item.status === 'success') {
            successCount++;
        } else if (item.status === 'failed') {
            failedCount++;
        }
    }

    return { successCount, failedCount };
});

const processedPercent = computed<number>(() => {
    if (props.importImages.length < 1) {
        return 0;
    }

    return Math.round(((importProgress.value.successCount + importProgress.value.failedCount) / props.importImages.length) * 100);
});

const progressingPercent = computed<number>(() => {
    if (props.importImages.length < 1) {
        return 0;
    }

    return Math.min(100, Math.round(((importProgress.value.successCount + importProgress.value.failedCount + 1) / props.importImages.length) * 100));
});

const progressText = computed<string>(() => {
    return tt(props.submitting ? 'format.misc.recognizingImagesProcess' : 'format.misc.recognizedImagesProcess', {
        success: formatNumberToLocalizedNumerals(importProgress.value.successCount),
        failed: formatNumberToLocalizedNumerals(importProgress.value.failedCount),
        total: formatNumberToLocalizedNumerals(props.importImages.length)
    });
});
</script>

<style>
.recognition-failed-list {
    max-height: 320px;
    overflow-y: auto;
}
</style>
