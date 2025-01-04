<template>
    <v-dialog persistent min-width="320" width="auto" v-model="showState">
        <v-card>
            <v-toolbar :color="finalColor">
                <v-toolbar-title>{{ titleContent }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text v-if="textContent" class="pa-4 pb-6">{{ textContent }}</v-card-text>
            <v-card-actions class="px-4 pb-4">
                <v-spacer></v-spacer>
                <v-btn color="gray" @click="cancel">{{ $t('Cancel') }}</v-btn>
                <v-btn :color="finalColor" @click="confirm">{{ $t('OK') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { type Ref, ref, watch } from 'vue';

import { useI18n } from '@/lib/i18n.js';
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

let resolveFunc: (value: T | PromiseLike<T>) => void = null;
let rejectFunc: (reason?: unknown) => void = null;

function open(title: string, text: string, options: Record<string, unknown>) {
    showState.value = true;

    if (isString(text)) {
        titleContent.value = tt(title, options);
        textContent.value = tt(text, options);
    } else {
        options = text;
        titleContent.value = tt('global.app.title');
        textContent.value = tt(title, options);
    }

    if (options && options.color) {
        finalColor.value = options.color || 'primary';
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

watch(() => showState, (newValue) => {
    emit('update:show', newValue);
});

defineExpose({
    open
});
</script>
