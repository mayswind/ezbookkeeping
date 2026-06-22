<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar class="toolbar-with-swipe-handler">
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link icon-f7="xmark" :class="{ 'disabled': recognizing }"
                         @click="cancel"></f7-link>
            </div>
            <div id="clipboard-text-recognition-sheet-toolbar-space"
                 class="clipboard-text-recognition-sheet-toolbar-space" @click="onToolbarClick()"></div>
            <div class="right">
                <f7-button round fill icon-f7="checkmark_alt"
                           :class="{ 'disabled': !pastedText }"
                           @click="confirm"></f7-button>
            </div>
        </f7-toolbar>
        <f7-page-content class="margin-top">
            <div class="padding-horizontal padding-bottom clipboard-text-recognition-sheet-content">
                <f7-list strong inset class="no-margin margin-vertical">
                    <f7-list-input
                        type="textarea"
                        class="clipboard-text"
                        :placeholder="tt('Click here to paste a transaction description')"
                        v-model:value="pastedText"
                    ></f7-list-input>
                </f7-list>
                <div class="margin-top-half display-flex justify-content-center align-items-center text-align-center">
                    <small>{{ tt('Uploaded text and personal data will be sent to the large language model, please be aware of potential privacy risks.') }}</small>
                </div>
            </div>
        </f7-page-content>

        <f7-popover class="paste-context-menu-popover" target-el="#clipboard-text-recognition-sheet-toolbar-space"
                    v-model:opened="showPastePopover">
            <f7-list class="paste-context-menu-list">
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt('Paste')" @click="pasteFromClipboard"></f7-list-item>
            </f7-list>
        </f7-popover>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { isiOS } from '@/lib/ui/mobile.ts';

import logger from '@/lib/logger.ts';

const props = defineProps<{
    show: boolean;
    initialText?: string;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'text:confirm', value: string): void;
}>();

const { tt } = useI18n();

const isSupportClipboard = !!navigator.clipboard;

const pastedText = ref<string>('');
const recognizing = ref<boolean>(false);
const pastingAmount = ref<boolean>(false);
const showPastePopover = ref<boolean>(false);

function pasteFromClipboard(): void {
    if (pastingAmount.value) {
        pastingAmount.value = false;
        return;
    }

    pastingAmount.value = true;

    navigator.clipboard.readText().then((text) => {
        pastingAmount.value = false;

        if (!text.trim()) {
            return;
        }

        pastedText.value = text;
    }).catch(error => {
        // Do not set pastingAmount to false here
        // In iOS, system will show the paste context menu, if user click outside, the paste action should not be triggered again
        logger.error('failed to read clipboard', error);
    });
}

function confirm(): void {
    if (!pastedText.value) {
        return;
    }

    recognizing.value = true;
    emit('text:confirm', pastedText.value);
    emit('update:show', false);
}

function cancel(): void {
    close();
}

function close(): void {
    emit('update:show', false);
    pastedText.value = '';
    recognizing.value = false;
}

function onSheetOpen(): void {
    pastedText.value = props.initialText || '';
    recognizing.value = false;
}

function onSheetClosed(): void {
    close();
}

function onToolbarClick(): void {
    if (!isSupportClipboard) {
        return;
    }

    if (isiOS()) {
        pasteFromClipboard();
    } else {
        showPastePopover.value = true;
    }
}
</script>

<style>
.clipboard-text-recognition-sheet-toolbar-space {
    width: 100%;
    height: 100%;
}

.clipboard-text-recognition-sheet-content {
    font-size: var(--f7-input-font-size);

    .clipboard-text {
        height: 200px;

        @media (min-height: 630px) {
            height: 370px;
        }
    }
}
</style>
