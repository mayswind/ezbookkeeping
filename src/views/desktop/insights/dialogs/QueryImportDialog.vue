<template>
    <v-dialog width="1000" :persistent="!!queriesJson" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4">{{ tt('Import Queries') }}</h4>
            </template>

            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto" style="height: 485px">
                <div class="w-100 h-100">
                    <v-textarea class="w-100 h-100 always-cursor-text" v-model="queriesJson"></v-textarea>
                </div>
            </v-card-text>

            <v-card-text>
                <div ref="buttonContainer" class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="primary" :disabled="!queriesJson" @click="confirm">{{ tt('Import') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { TransactionExplorerQuery } from '@/models/explorer.ts';

import logger from '@/lib/logger.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const queriesJson = ref<string>('');

let resolveFunc: ((queries: TransactionExplorerQuery[]) => void) | null = null;
let rejectFunc: (() => void) | null = null;

function open(): Promise<TransactionExplorerQuery[]> {
    queriesJson.value = '';
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!queriesJson.value) {
        return;
    }

    try {
        const queries: TransactionExplorerQuery[] = TransactionExplorerQuery.parseFromQueryiesJson(queriesJson.value);

        if (!queries || queries.length < 1) {
            snackbar.value?.showError('No valid queries found in the input. Please make sure the queries are valid and try again.');
            return;
        }

        resolveFunc?.(queries);
        showState.value = false;
    } catch (error) {
        logger.error('Failed to import queries', error);
        snackbar.value?.showError('Queries import failed. Please make sure the queries are valid and try again.');
    }
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
