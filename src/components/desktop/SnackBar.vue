<template>
    <v-snackbar v-model="showState">
        {{ messageContent }}

        <template #actions>
            <v-btn color="primary" variant="text" @click="showState = false">{{ tt('Close') }}</v-btn>
        </template>
    </v-snackbar>
</template>

<script setup lang="ts">
import { type Ref, ref, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { isObject } from '@/lib/common.ts';

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
}>();

const { tt, te } = useI18n();

const showState: Ref<boolean> = ref(false);
const messageContent: Ref<string> = ref('');

function showMessage(message: string, options: Record<string, unknown>): void {
    showState.value = true;
    messageContent.value = tt(message, options);
}

function showError(error: string | { message: string }): void {
    showState.value = true;

    if (isObject(error) && (error as { message: string }).message) {
        messageContent.value = te((error as { message: string }).message);
    } else {
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
