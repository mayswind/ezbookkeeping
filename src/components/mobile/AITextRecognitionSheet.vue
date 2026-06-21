<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar class="toolbar-with-swipe-handler">
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link icon-f7="xmark" :class="{ 'disabled': recognizing }"
                         @click="cancel"></f7-link>
            </div>
            <div class="right">
                <f7-button round fill icon-f7="checkmark_alt"
                           :class="{ 'disabled': !pastedText }"
                           @click="confirm"></f7-button>
            </div>
        </f7-toolbar>
        <f7-page-content class="no-margin-vertical no-padding-vertical">
            <div class="text-recognition-container display-flex justify-content-center align-items-center text-align-center padding-horizontal" @click="pasteFromClipboard">
                <div class="display-inline-flex flex-direction-column">
                    <textarea ref="pastedTextArea" class="text-align-center" :placeholder="tt('Click here to paste a transaction description')" v-model="pastedText"
                              @input="onPastedTextInput" @change="onPastedTextChanged"></textarea>
                    <small class="margin-top-half">{{ tt('Uploaded text and personal data will be sent to the large language model, please be aware of potential privacy risks.') }}</small>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import logger from '@/lib/logger.ts';

defineProps<{
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'text:confirm', value: string): void;
}>();

const { tt } = useI18n();

const pastedTextArea = useTemplateRef<HTMLTextAreaElement>('pastedTextArea');

const pastedText = ref<string>('');
const recognizing = ref<boolean>(false);

function autosizeTextareaHeight(textarea: HTMLTextAreaElement, resetScroll: boolean): void {
    if (resetScroll) {
        textarea.scrollTop = 0;
    }

    textarea.style.height = '';
    textarea.style.height = Math.min(textarea.scrollHeight, 400) + 'px';
}

function pasteFromClipboard(): void {
    navigator.clipboard.readText().then((text) => {
        pastedText.value = text;
    }).catch(error => {
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
    pastedText.value = '';
    recognizing.value = false;
    autosizeTextareaHeight(pastedTextArea.value as HTMLTextAreaElement, true);
}

function onSheetClosed(): void {
    close();
}

function onPastedTextInput(event: Event): void {
    autosizeTextareaHeight(event.target as HTMLTextAreaElement, false);
}

function onPastedTextChanged(event: Event): void {
    autosizeTextareaHeight(event.target as HTMLTextAreaElement, false);
}
</script>

<style>
.text-recognition-container {
    --ebk-ai-text-recognition-height: 310px;
    height: var(--ebk-ai-text-recognition-height);
    border: 1px solid var(--f7-page-master-border-color);
    font-size: var(--f7-input-font-size);

    @media (min-height: 630px) {
        --ebk-ai-text-recognition-height: 525px;
    }

    textarea::placeholder {
        color: var(--f7-text-color);
    }

    textarea:focus::placeholder {
        opacity: 0.5;
    }
}
</style>
