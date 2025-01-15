<template>
    <v-snackbar v-model="showState">
        {{ messageContent }}

        <template #actions>
            <v-btn color="primary" variant="text" @click="showState = false">{{ tt('Close') }}</v-btn>
        </template>
    </v-snackbar>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { isObject, isString } from '@/lib/common.ts';

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
}>();

const { tt, te } = useI18n();

const showState= ref<boolean>(false);
const messageContent = ref<string>('');

function showMessage(message: string, options?: Record<string, unknown>): void {
    showState.value = true;

    if (options) {
        messageContent.value = tt(message, options);
    } else {
        messageContent.value = tt(message);
    }
}

function showError(error: string | { message: string }): void {
    showState.value = true;

    if (isObject(error) && error.message) {
        messageContent.value = te(error.message);
    } else if (isString(error)) {
        messageContent.value = te(error);
    }
}

watch(showState, (newValue) => {
    emit('update:show', newValue);
});

defineExpose({
    showMessage,
    showError
});
</script>
