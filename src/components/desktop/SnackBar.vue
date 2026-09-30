<template>
    <v-snackbar :color="color" :variant="variant" :location="location"
                :timeout="notAutoHide ? -1 : undefined" v-model="showState">
        {{ messageContent }}

        <template #actions>
            <v-btn variant="text" :color="color ?? 'primary'" @click="showState = false">{{ tt('Close') }}</v-btn>
        </template>
    </v-snackbar>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import type { Anchor } from 'vuetify/lib/util/anchor.d.ts';

import { useI18n } from '@/locales/helpers.ts';

import type { ErrorResponse } from '@/core/api.ts';
import { isObject, isString } from '@/lib/common.ts';
import type { SnackBarVariant } from '@/lib/ui/desktop.ts';

defineProps<{
    notAutoHide?: boolean;
    color?: string;
    variant?: SnackBarVariant;
    location?: Anchor;
}>();

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

function showError(error: string | { message: string } | { error: ErrorResponse }): void {
    showState.value = true;

    if (isObject(error) && 'message' in error && error.message) {
        messageContent.value = te(error.message);
    } else if (isObject(error) && 'error' in error) {
        messageContent.value = te(error);
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
