<template>
    <v-dialog max-width="500" :persistent="oldName !== newName" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4 text-wrap">{{ dialogTitle || defaultTitle }}</h4>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <v-text-field persistent-placeholder
                              :autofocus="true"
                              :label="label"
                              :placeholder="placeholder"
                              v-model="newName"
                              @keyup.enter="save" />
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="primary" :disabled="!newName || oldName === newName" @click="save">
                        {{ tt('Save') }}
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">
                        {{ tt('Cancel') }}
                    </v-btn>
                </div>
            </v-card-text>
        </v-card>
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

