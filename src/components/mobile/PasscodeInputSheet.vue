<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin" v-if="hint">{{ hint }}</p>
                <slot></slot>
                <f7-list class="no-margin">
                    <f7-list-input
                        type="number"
                        autocomplete="one-time-code"
                        outline
                        clear-button
                        class="no-margin no-padding-bottom"
                        :placeholder="tt('Passcode')"
                        v-model:value="currentPasscode"
                        @keyup.enter="confirm()"
                    ></f7-list-input>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentPasscode || confirmDisabled }"
                           :text="tt('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link :class="{ 'disabled': cancelDisabled }" @click="cancel" :text="tt('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

const props = defineProps<{
    modelValue: string;
    title?: string;
    hint?: string;
    confirmDisabled?: boolean;
    cancelDisabled?: boolean;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'update:show', value: boolean): void;
    (e: 'passcode:confirm', value: string): void;
}>();

const { tt } = useI18n();

const currentPasscode = ref<string>('');

function confirm(): void {
    if (!currentPasscode.value || props.confirmDisabled) {
        return;
    }

    emit('update:modelValue', currentPasscode.value);
    emit('passcode:confirm', currentPasscode.value);
}

function cancel(): void {
    close();
}

function close(): void {
    emit('update:show', false);
}

function onSheetOpen(): void {
    currentPasscode.value = '';
}

function onSheetClosed(): void {
    close();
}
</script>
