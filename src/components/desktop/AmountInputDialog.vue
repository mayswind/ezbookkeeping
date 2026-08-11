<template>
    <v-dialog persistent width="500" v-model="showState">
        <one-column-dialog-layout :title="titleContent" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #content>
                <div class="text-body-large text-wrap mb-4" v-if="textContent">{{ textContent }}</div>
                <div class="d-flex justify-center w-100">
                    <v-form class="w-100">
                        <v-row>
                            <v-col cols="12">
                                <amount-input :persistent-placeholder="true"
                                              :autofocus="true"
                                              :currency="dialogOptions?.currency"
                                              :show-currency="!!dialogOptions?.currency"
                                              :enable-formula="true"
                                              :label="inputLabelContent"
                                              :placeholder="inputPlaceholderContent"
                                              v-model="amount"
                                              @enter="confirm" />
                            </v-col>
                        </v-row>
                    </v-form>
                </div>
            </template>

            <template #footer>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-spacer />
                <v-btn :color="finalColor" @click="confirm">{{ tt('OK') }}</v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

interface AmountInputDialogOptions {
    title?: string;
    text?: string;
    textI18nOptions?: Record<string, unknown>;
    inputLabel?: string;
    inputPlaceholder?: string;
    color?: string;
    currency?: string;
    initAmount?: number;
}

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const showState = ref<boolean>(false);
const dialogOptions = ref<AmountInputDialogOptions | undefined>(undefined);
const amount = ref<number>(0);

let resolveFunc: ((value: number) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const titleContent = computed<string>(() => dialogOptions.value?.title ? tt(dialogOptions.value.title) : tt('global.app.title'));
const textContent = computed<string>(() => dialogOptions.value?.text ? tt(dialogOptions.value.text, dialogOptions.value.textI18nOptions ?? {}) : '');
const inputLabelContent = computed<string | undefined>(() => dialogOptions.value?.inputLabel ? tt(dialogOptions.value.inputLabel) : undefined);
const inputPlaceholderContent = computed<string | undefined>(() => dialogOptions.value?.inputPlaceholder ? tt(dialogOptions.value.inputPlaceholder) : undefined);
const finalColor = computed<string>(() => dialogOptions.value?.color || 'primary');

function open(options: AmountInputDialogOptions): Promise<number> {
    showState.value = true;
    dialogOptions.value = options;
    amount.value = options.initAmount ?? 0;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (resolveFunc) {
        resolveFunc(amount.value);
    }

    showState.value = false;
    emit('update:show', false);
}

function cancel(): void {
    if (rejectFunc) {
        rejectFunc();
    }

    showState.value = false;
    emit('update:show', false);
}

watch(showState, newValue => {
    emit('update:show', newValue);
});

defineExpose({
    open
});
</script>
