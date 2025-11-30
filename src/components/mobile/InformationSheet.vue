<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:closed="onSheetClosed">
        <div class="swipe-handler"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">
                    <span>{{ hint }}</span>
                    <f7-link class="icon-after-text"
                             icon-only icon-f7="doc_on_doc"
                             @click="copyBackupCodes"
                             v-if="enableCopy"
                    ></f7-link>
                </p>
                <textarea class="information-content full-line" :readonly="true" :rows="rowCount" :value="information"></textarea>
                <div class="margin-top text-align-center">
                    <f7-link @click="close" :text="tt('Close')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';

import { copyTextToClipboard } from '@/lib/ui/common.ts';

const props = defineProps<{
    title?: string;
    hint?: string;
    information: string;
    rowCount: number;
    enableCopy?: boolean;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'info:copied'): void;
}>();

const { tt } = useI18n();

function copyBackupCodes(): void {
    copyTextToClipboard(props.information);
    emit('info:copied');
}

function close(): void {
    emit('update:show', false);
}

function onSheetClosed(): void {
    close();
}
</script>
