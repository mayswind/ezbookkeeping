<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">
                    <span>{{ hint }}</span>
                    <f7-link id="copy-to-clipboard-icon" ref="copyToClipboardIcon"
                             class="icon-after-text"
                             icon-only icon-f7="doc_on_doc"
                             v-if="enableCopy"
                    ></f7-link>
                </p>
                <textarea class="information-content full-line" readonly="readonly" :rows="rowCount" :value="information"></textarea>
                <div class="margin-top text-align-center">
                    <f7-link @click="close" :text="tt('Close')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { useTemplateRef, watch, onMounted, onUpdated } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { ClipboardHolder } from '@/lib/clipboard.ts';

const props = defineProps<{
    title?: string
    hint?: string
    information: string
    rowCount: number
    enableCopy?: boolean
    show: boolean
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
    (e: 'info:copied'): void
}>();

const { tt } = useI18n();

const iconCopyToClipboard = useTemplateRef('copyToClipboardIcon');

let clipboardHolder: ClipboardHolder | null = null;

function makeCopyToClipboardClickable() {
    if (clipboardHolder) {
        return;
    }

    if (iconCopyToClipboard.value) {
        clipboardHolder = ClipboardHolder.create({
            el: '#copy-to-clipboard-icon',
            text: props.information,
            successCallback: function () {
                emit('info:copied');
            }
        });
    }
}

function close() {
    emit('update:show', false);
}

function onSheetClosed() {
    close();
}

onMounted(() => {
    makeCopyToClipboardClickable();
});

onUpdated(() => {
    makeCopyToClipboardClickable();
});

watch(() => props.information, (newValue) => {
    if (clipboardHolder) {
        clipboardHolder.setClipboardText(newValue);
    }
});
</script>
