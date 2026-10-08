<template>
    <v-snackbar :color="isError ? 'error' : undefined" :timeout="isError ? 6000 : 3000" v-model="show">
        {{ message }}
    </v-snackbar>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { describeError } from '@/ext/api.ts';

const show = ref<boolean>(false);
const message = ref<string>('');
const isError = ref<boolean>(false);

function showMessage(text: string): void {
    message.value = text;
    isError.value = false;
    show.value = true;
}

function showError(error: unknown): void {
    const text = typeof error === 'string' ? error : describeError(error);

    if (!text) {
        return; // already handled by the app (for example a session that expired)
    }

    message.value = text;
    isError.value = true;
    show.value = true;
}

defineExpose({ showMessage, showError });
</script>
