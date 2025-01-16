<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin">{{ hint }}</p>
                <f7-list class="no-margin">
                    <f7-list-item class="list-item-pincode-input padding-vertical-half">
                        <pin-code-input :secure="true" :length="6" v-model="currentPinCode" @pincode:confirm="confirm"/>
                    </f7-list-item>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentPinCodeValid || confirmDisabled }"
                           :text="tt('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link @click="cancel" :text="tt('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

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
    (e: 'pincode:confirm', value: string): void;
}>();

const { tt } = useI18n();

const currentPinCode = ref<string>('');

const currentPinCodeValid = computed<boolean>(() => {
    return currentPinCode.value?.length === 6 || false;
});

function confirm(): void {
    if (!currentPinCodeValid.value || props.confirmDisabled) {
        return;
    }

    emit('update:modelValue', currentPinCode.value);
    emit('pincode:confirm', currentPinCode.value);
}

function cancel(): void {
    emit('update:show', false);
}

function onSheetOpen(): void {
    currentPinCode.value = '';
}

function onSheetClosed(): void {
    cancel();
}

watch(() => props.modelValue, (newValue) => {
    if (newValue === '') {
        currentPinCode.value = '';
    }
});
</script>

<style>
.list-item-pincode-input .item-content {
    padding-left: 0;
    padding-right: 0;
}

.list-item-pincode-input .item-content .item-inner {
    padding-left: 0;
    padding-right: 0;
    justify-content: center;
}
</style>
