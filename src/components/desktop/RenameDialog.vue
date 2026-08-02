<template>
    <v-dialog width="500" :persistent="oldName !== newName" v-model="showState">
        <one-column-dialog-layout :title="dialogTitle || defaultTitle" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #content>
                <div class="mt-4">
                    <v-text-field persistent-placeholder
                                  :autofocus="true"
                                  :label="label"
                                  :placeholder="placeholder"
                                  v-model="newName"
                                  @keyup.enter="save" />
                </div>
            </template>

            <template #footer>
                <v-spacer />
                <v-btn color="primary" :disabled="!newName || oldName === newName" @click="save">{{ tt('Save') }}</v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

defineProps<{
    label?: string;
    placeholder?: string;
    defaultTitle?: string;
}>();

const { tt } = useI18n();

let resolveFunc: ((name: string) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const dialogTitle = ref<string | undefined>(undefined);
const oldName = ref<string>('');
const newName = ref<string>('');

function open(currentName: string, title?: string): Promise<string> {
    showState.value = true;
    dialogTitle.value = title;
    oldName.value = currentName;
    newName.value = currentName;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function save(): void {
    if (!newName.value || oldName.value === newName.value) {
        return;
    }

    resolveFunc?.(newName.value);
    showState.value = false;
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>

