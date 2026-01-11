<template>
    <v-dialog persistent min-width="320" width="auto" v-model="showState">
        <v-card>
            <v-toolbar :color="finalColor">
                <v-toolbar-title>{{ titleContent }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text v-if="textContent" class="pa-4 pb-6">{{ textContent }}</v-card-text>
            <v-card-text>
                <v-form>
                    <v-row>
                        <v-col cols="12">
                            <amount-input :persistent-placeholder="true"
                                          :autofocus="true"
                                          :currency="dialogOptions?.currency"
                                          :show-currency="!!dialogOptions?.currency"
                                          :label="inputLabelContent"
                                          :placeholder="inputPlaceholderContent"
                                          v-model="amount"
                                          @enter="confirm" />
                        </v-col>
                    </v-row>
                </v-form>
            </v-card-text>
            <v-card-actions class="px-4 pb-4">
                <v-spacer></v-spacer>
                <v-btn color="gray" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-btn :color="finalColor" @click="confirm">{{ tt('OK') }}</v-btn>
            </v-card-actions>
        </v-card>
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
