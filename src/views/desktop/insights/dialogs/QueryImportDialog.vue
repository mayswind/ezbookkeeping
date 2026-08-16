<template>
    <v-dialog width="1000" :persistent="!!queriesJson && queriesJson !== sampleJson" v-model="showState">
        <one-column-dialog-layout content-class="pa-0" content-style="height: 500px"
                                  :title="tt('Import Queries')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="!queriesJson" @click="confirm">{{ tt('Import') }}</v-btn>
            </template>

            <template #content>
                <div class="w-100 h-100">
                    <v-textarea no-resize class="w-100 h-100 ps-3 code-textarea always-cursor-text"
                                density="compact" variant="plain" :rounded="false"
                                v-model="queriesJson"></v-textarea>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef} from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { TransactionExplorerQuery } from '@/models/explorer.ts';

import { isTextualUUID } from '@/lib/common.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import logger from '@/lib/logger.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

let resolveFunc: ((queries: TransactionExplorerQuery[]) => void) | null = null;
let rejectFunc: (() => void) | null = null;

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const queriesJson = ref<string>('');

const sampleJson = computed<string>(() => `[
    {
        "id": "", // ${tt('sample.importInsightsExplorerQueries.queryIdDescription')}
        "name": "", // ${tt('sample.importInsightsExplorerQueries.queryNameDescription')}
        "conditions": [
            {
                "condition": { // ${tt('sample.importInsightsExplorerQueries.conditionDescription')}
                    "field": "", // ${tt('sample.importInsightsExplorerQueries.conditionFieldDescription')}
                    "operator": "", // ${tt('sample.importInsightsExplorerQueries.conditionOperatorDescription')}
                    "value": "" // ${tt('sample.importInsightsExplorerQueries.conditionValueDescription')}
                },
                "relation": "" // ${tt('sample.importInsightsExplorerQueries.conditionRelationDescription')}
            }
            // ${tt('sample.importInsightsExplorerQueries.additionalQueryDescription')}
        ]
    }
]`);

function open(): Promise<TransactionExplorerQuery[]> {
    queriesJson.value = sampleJson.value;
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
        const queryItems = JSON.parse(queriesJson.value);

        if (!Array.isArray(queryItems)) {
            snackbar.value?.showError('Queries import failed. Please make sure the queries are valid and try again.');
            return;
        }

        if (!queryItems || queryItems.length < 1) {
            snackbar.value?.showError('No valid queries found in the input. Please make sure the queries are valid and try again.');
            return;
        }

        const queries: TransactionExplorerQuery[] = [];
        const queryIds: Record<string, boolean> = {};

        for (const queryItem of queryItems) {
            let originalId: string = '';

            if (!('id' in queryItem) || !queryItem['id']) {
                queryItem['id'] = generateRandomUUID();
            } else {
                const queryId = queryItem['id'];

                if (!isTextualUUID(queryId)) {
                    snackbar.value?.showMessage('format.misc.queryIdInvalidTip', { id: queryId });
                    return;
                }

                originalId = queryId;
            }

            const query = TransactionExplorerQuery.parse(queryItem);

            if (!query) {
                snackbar.value?.showMessage('format.misc.queryInvalidTip', { id: originalId });
                return;
            }

            if (queryIds[query.id]) {
                snackbar.value?.showMessage('format.misc.queryIdDuplicatedTip', { id: query.id });
                return;
            }

            queries.push(query);
            queryIds[query.id] = true;
        }

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
