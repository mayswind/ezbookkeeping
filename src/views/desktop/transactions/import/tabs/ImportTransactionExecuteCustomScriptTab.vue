<template>
    <v-row>
        <v-col cols="12" md="6">
            <div class="d-flex w-100 mb-2">
                <v-btn density="compact" variant="tonal" :prepend-icon="mdiPlay"
                       :disabled="disabled || !sandboxLoaded || executingScript" :loading="executingScript"
                       @click="executeCustomScript()">
                    <template #loader>
                        <v-progress-circular indeterminate size="18" class="me-1"/>
                        <span>{{ tt('Execute Custom Script') }}</span>
                    </template>
                    <span>{{ tt('Execute Custom Script') }}</span>
                </v-btn>
            </div>
            <v-textarea class="w-100" style="height: 360px" :readonly="disabled"
                        v-model="customScript"></v-textarea>
        </v-col>
        <v-col cols="12" md="6">
            <div class="d-flex w-100 mb-2">
                <v-btn density="compact" color="default" variant="text"
                       :disabled="disabled || !sandboxLoaded || executingScript || !previewResult">
                    <span>{{ tt('format.misc.previewCount', { count: previewCount > 0 ? getDisplayCount(previewCount) : tt('All') }) }}</span>
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :key="count.value" :title="count.name"
                                         v-for="count in previewCounts"
                                         @click="previewCount = count.value"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </div>
            <div class="w-100 code-container">
                <v-textarea class="w-100 always-cursor-text" style="height: 360px" :readonly="true"
                            :value="displayPreviewResult"></v-textarea>
            </div>
        </v-col>
    </v-row>

    <iframe id="sandbox" ref="sandbox" sandbox="allow-scripts" style="display: none;"></iframe>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, onMounted, onUnmounted } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { NameNumeralValue } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';
import { KnownFileType } from '@/core/file.ts';

import type { ImportTransactionRequest, ImportTransactionRequestItem } from '@/models/imported_transaction.ts';

import { isDefined } from '@/lib/common.ts';
import { getBrowserTimezoneOffsetMinutes } from '@/lib/datetime.ts';
import {
    openTextFileContent,
    startDownloadFile
} from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

import {
    mdiPlay,
    mdiFolderOpenOutline,
    mdiContentSaveOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

type SandboxRequest = {
    parsedFileData: string[][];
    code: string;
};

type SandboxResponse = {
    result?: string;
    knownError?: string;
    error?: string;
};

interface ImportTransactionDefineColumnMenu {
    prependIcon: string;
    title: string;
    disabled?: boolean;
    onClick: () => void;
}

const props = defineProps<{
    parsedFileData?: string[][];
    disabled?: boolean;
}>();

const { tt, getCurrentNumeralSystemType } = useI18n();

const sandbox = useTemplateRef<HTMLIFrameElement>('sandbox');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const sandboxLoaded = ref<boolean>(false);
const customScript = ref<string>('');
const previewResult = ref<ImportTransactionRequestItem[] | undefined>(undefined);
const executingScript = ref<boolean>(false);
const previewCount = ref<number>(10);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const previewCounts = computed<NameNumeralValue[]>(() => getTablePageOptions(previewResult.value?.length));

const sampleScript = computed<string>(() => `// ${tt('sample.importTransactionCustomScript.headerComment')}
/**
 * ${tt('sample.importTransactionCustomScript.functionDescription')}
 * @param {array} row - ${tt('sample.importTransactionCustomScript.functionParamRowDescription')}
 * @param {number} index - ${tt('sample.importTransactionCustomScript.functionParamIndexDescription')}
 * @returns {object|null} ${tt('sample.importTransactionCustomScript.functionReturnDescription')}
 */
function parse(row, index) {
    if (index < 1) {
        return null;
    }

    return {
        time: row[0], // ${tt('sample.importTransactionCustomScript.fieldTimeDescription')}
        utcOffset: '${getBrowserTimezoneOffsetMinutes()}', // ${tt('sample.importTransactionCustomScript.fieldUtcOffsetDescription')}
        type: TransactionType.Expense, // ${tt('sample.importTransactionCustomScript.fieldTypeDescription')}
        categoryName: row[4], // ${tt('sample.importTransactionCustomScript.fieldCategoryNameDescription')}
        sourceAccountName: row[5], // ${tt('sample.importTransactionCustomScript.fieldSourceAccountNameDescription')}
        destinationAccountName: row[8], // ${tt('sample.importTransactionCustomScript.fieldDestinationAccountNameDescription')}
        sourceAmount: row[7], // ${tt('sample.importTransactionCustomScript.fieldSourceAmountDescription')}
        destinationAmount: row[10], // ${tt('sample.importTransactionCustomScript.fieldDestinationAmountDescription')}
        geoLocation: undefined, // ${tt('sample.importTransactionCustomScript.fieldGeoLocationDescription')}
        tagNames: '', // ${tt('sample.importTransactionCustomScript.fieldTagNamesDescription')}
        description: row[13] // ${tt('sample.importTransactionCustomScript.fieldCommentDescription')}
    };
}`);

const displayPreviewResult = computed<string>(() => {
    if (executingScript.value) {
        return tt('Executing Script...');
    } else if (previewResult.value) {
        const rows = previewResult.value.slice(0, previewCount.value);
        return JSON.stringify(rows, null, 2);
    } else {
        return tt('No Preview Result');
    }
});

const menus = computed<ImportTransactionDefineColumnMenu[]>(() => [
    {
        prependIcon: mdiFolderOpenOutline,
        title: tt('Load Script File'),
        onClick: loadScriptFile
    },
    {
        prependIcon: mdiContentSaveOutline,
        title: tt('Save Script File'),
        onClick: saveScriptFile
    }
]);

function getDisplayCount(count: number): string {
    return numeralSystem.value.formatNumber(count);
}

function getTablePageOptions(linesCount?: number): NameNumeralValue[] {
    const pageOptions: NameNumeralValue[] = [];

    if (!linesCount || linesCount < 1) {
        pageOptions.push({ value: -1, name: tt('All') });
        return pageOptions;
    }

    for (const count of [ 10, 50, 100 ]) {
        if (linesCount < count) {
            break;
        }

        pageOptions.push({ value: count, name: getDisplayCount(count) });
    }

    pageOptions.push({ value: -1, name: tt('All') });

    return pageOptions;
}

function reloadSandbox(): void {
    sandboxLoaded.value = false;

    if (sandbox.value) {
        sandbox.value.src = 'about:blank';
        sandbox.value.srcdoc = `
            <script>
                window.TransactionType = {
                    Income: 'Income',
                    Expense: 'Expense',
                    Transfer: 'Transfer'
                };

                window.addEventListener('message', function (event) {
                    try {
                        const data = JSON.parse(event.data);
                        const parsedFileData = data.parsedFileData;
                        eval(data.code);

                        if (window.parse) {
                            const result = [];

                            for (let i = 0; i < parsedFileData.length; i++) {
                                try {
                                    const row = parsedFileData[i];
                                    const transaction = window.parse(row, i);

                                    if (transaction) {
                                        result.push(transaction);
                                    }
                                } catch (error) {
                                    window.parent.postMessage({ error: error.message }, '*');
                                    return;
                                }
                            }

                            window.parent.postMessage({ result: JSON.stringify(result) }, '*');
                        } else {
                            window.parent.postMessage({ knownError: 'No parse function defined' }, '*');
                        }
                    } catch (error) {
                        window.parent.postMessage({ error: error.message }, '*');
                    }
                });
            <\/script>
        `;

        sandbox.value.onload = () => {
            sandboxLoaded.value = true;
        };
    }
}

function executeCustomScript(): void {
    if (!sandbox.value || props.disabled || executingScript.value) {
        return;
    }

    executingScript.value = true;

    const sandboxRequest: SandboxRequest = {
        parsedFileData: props.parsedFileData || [],
        code: customScript.value + '\n\n;window.parse = parse;'
    };

    sandbox.value?.contentWindow?.postMessage(JSON.stringify(sandboxRequest), '*');
}

function generateResult(): string | undefined {
    if (!previewResult.value) {
        snackbar.value?.showError('Please execute the custom script first');
        return undefined;
    }

    const result: ImportTransactionRequest = {
        transactions: previewResult.value
    };

    return JSON.stringify(result);
}

function reset(): void {
    customScript.value = sampleScript.value;
    previewResult.value = undefined;
    executingScript.value = false;
    previewCount.value = 10;
}

function loadScriptFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JS.contentType
    }).then(content => {
        customScript.value = content;
    }).catch(error => {
        logger.error('Failed to load script file', error);
        snackbar.value?.showError('Cannot load script file');
    });
}

function saveScriptFile(): void {
    const fileName = KnownFileType.JS.formatFileName(tt('dataExport.defaultImportHandlingScript'));
    startDownloadFile(fileName, KnownFileType.JS.createBlob(customScript.value));
}

function onMessage(event: MessageEvent<SandboxResponse>): void {
    if (event.source !== sandbox.value?.contentWindow) {
        return;
    }

    executingScript.value = false;

    const data = event.data;

    if (data.knownError) {
        snackbar.value?.showError(data.knownError);
        previewResult.value = undefined;
    } else if (data.error) {
        logger.error('Failed to execute custom script: ' + data.error);
        snackbar.value?.showError('Failed to execute custom script');
        previewResult.value = undefined;
    } else if (data.result) {
        const originalResult = JSON.parse(data.result) as Record<string, unknown>[];
        const finalResult: ImportTransactionRequestItem[] = [];

        for (const item of originalResult) {
            const finalItem: ImportTransactionRequestItem = {
                time: (isDefined(item['time'])) ? String(item['time']) : '',
                utcOffset: (isDefined(item['utcOffset'])) ? String(item['utcOffset']) : '',
                type: (isDefined(item['type'])) ? String(item['type']) : '',
                categoryName: (isDefined(item['categoryName']) && item['categoryName'] !== '') ? String(item['categoryName']) : undefined,
                sourceAccountName: (isDefined(item['sourceAccountName']) && item['sourceAccountName'] !== '') ? String(item['sourceAccountName']) : undefined,
                destinationAccountName: (isDefined(item['destinationAccountName']) && item['destinationAccountName'] !== '') ? String(item['destinationAccountName']) : undefined,
                sourceAmount: (isDefined(item['sourceAmount'])) ? String(item['sourceAmount']) : '',
                destinationAmount: (isDefined(item['destinationAmount']) && item['destinationAmount'] !== '') ? String(item['destinationAmount']) : undefined,
                geoLocation: (isDefined(item['geoLocation']) && item['geoLocation']) ? String(item['geoLocation']) : undefined,
                tagNames: (isDefined(item['tagNames']) && item['tagNames']) ? String(item['tagNames']) : undefined,
                comment: (isDefined(item['description']) && item['description']) ? String(item['description']) : undefined
            };
            finalResult.push(finalItem);
        }

        previewResult.value = finalResult;
    }

    reloadSandbox();
}

onMounted(() => {
    customScript.value = sampleScript.value;
    reloadSandbox();
    window.addEventListener('message', onMessage);
});

onUnmounted(() => {
    window.removeEventListener('message', onMessage);
});

defineExpose({
    menus,
    generateResult,
    reset
});
</script>
