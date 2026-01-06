<template>
    <v-dialog max-width="500" :persistent="oldExplorerName !== newExplorerName" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4 text-wrap">{{ tt('Rename Explorer') }}</h4>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <v-text-field persistent-placeholder
                              :label="tt('Explorer Name')"
                              :placeholder="tt('Explorer Name')"
                              v-model="newExplorerName"/>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="primary" :disabled="!newExplorerName || oldExplorerName === newExplorerName" @click="save">
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

const { tt } = useI18n();

let resolveFunc: ((name: string) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const oldExplorerName = ref<string>('');
const newExplorerName = ref<string>('');

function open(currentExplorerName: string): Promise<string> {
    showState.value = true;
    oldExplorerName.value = currentExplorerName;
    newExplorerName.value = currentExplorerName;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function save(): void {
    resolveFunc?.(newExplorerName.value);
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

