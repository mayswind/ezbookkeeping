<template>
    <v-dialog persistent min-width="320" width="auto" v-model="showState">
        <v-card>
            <v-toolbar :color="finalColor">
                <v-toolbar-title>{{ titleContent }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text v-if="textContent" class="pa-4 pb-6">{{ textContent }}</v-card-text>
            <v-card-actions class="px-4 pb-4">
                <v-spacer></v-spacer>
                <v-btn color="gray" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-btn :color="finalColor" @click="confirm">{{ tt('OK') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { type Ref, ref, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { isString } from '@/lib/common.ts';

const props = defineProps<{
    show?: boolean
    color?: string
    title?: string
    text?: string
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
}>();

const { tt } = useI18n();

const showState: Ref<boolean> = ref(false);
const titleContent: Ref<string> = ref(props.title || tt('global.app.title'));
const textContent: Ref<string> = ref(props.text || '');
const finalColor: Ref<string> = ref(props.color || 'primary');

let resolveFunc: ((value?: unknown) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

function open(titleOrText: string, textOrOptions: string | Record<string, unknown>, options: Record<string, unknown>) {
    showState.value = true;

    if (isString(textOrOptions)) { // second parameter is text
        titleContent.value = tt(titleOrText, options);
        textContent.value = tt(textOrOptions as string, options);
    } else { // second parameter is options
        const actualOptions = textOrOptions as Record<string, unknown>;
        titleContent.value = tt('global.app.title');
        textContent.value = tt(titleOrText, actualOptions);
    }

    if (options && isString(options.color)) {
        finalColor.value = (options.color as string) || 'primary';
    }

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (resolveFunc) {
        resolveFunc();
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
