<template>
    <v-dialog :persistent="!!persistent" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center">
                        <h4 class="text-h4">{{ tt('Import Transactions') }}</h4>
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :icon="true" :disabled="loading || submitting"
                           v-if="currentStep === 'defineColumn' && importTransactionDefineColumnTab?.menus">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-item :key="index"
                                             :prepend-icon="menu.prependIcon"
                                             :title="menu.title"
                                             :disabled="menu.disabled"
                                             @click="menu.onClick()"
                                             v-for="(menu, index) in importTransactionDefineColumnTab.menus"/>
                            </v-list>
                        </v-menu>
                    </v-btn>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :icon="true" :disabled="loading || submitting"
                           v-if="currentStep === 'executeCustomScript' && importTransactionExecuteCustomScriptTab?.menus">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-item :key="index"
                                             :prepend-icon="menu.prependIcon"
                                             :title="menu.title"
                                             :disabled="menu.disabled"
                                             @click="menu.onClick()"
                                             v-for="(menu, index) in importTransactionExecuteCustomScriptTab.menus"/>
                            </v-list>
                        </v-menu>
                    </v-btn>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :icon="true" :disabled="loading || submitting"
                           v-if="currentStep === 'checkData' && importTransactionCheckDataTab?.filterMenus">
                        <v-icon :icon="mdiFilterOutline" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <template :key="groupIndex" v-for="(group, groupIndex) in importTransactionCheckDataTab.filterMenus">
                                    <v-divider class="my-2" v-if="groupIndex > 0" />
                                    <v-list-subheader :title="group.title" v-if="group.title" />
                                    <v-list-item :key="`menu_${groupIndex}_${index}`"
                                                 :prepend-icon="menu.prependIcon"
                                                 :title="menu.title"
                                                 :subtitle="menu.subTitle"
                                                 :append-icon="menu.appendIcon"
                                                 :disabled="menu.disabled"
                                                 @click="menu.onClick()"
                                                 v-for="(menu, index) in group.items" />
                                </template>
                            </v-list>
                        </v-menu>
                    </v-btn>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :icon="true" :disabled="loading || submitting"
                           v-if="currentStep === 'checkData' && importTransactionCheckDataTab?.toolMenus">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <template :key="index" v-for="(menu, index) in importTransactionCheckDataTab.toolMenus">
                                    <v-divider class="my-2" v-if="menu.divider" />
                                    <v-list-item :prepend-icon="menu.prependIcon"
                                                 :title="menu.title"
                                                 :subtitle="menu.subTitle"
                                                 :append-icon="menu.appendIcon"
                                                 :disabled="menu.disabled"
                                                 @click="menu.onClick()" />
                                </template>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <v-card-text>
                <div class="cursor-default">
                    <steps-bar min-width="700" :clickable="false" :steps="allSteps" :current-step="currentStep" />
                </div>

                <v-window class="disable-tab-transition" v-model="currentStep">
                    <v-window-item value="uploadFile">
                        <v-row class="pt-2">
                            <v-col cols="12" md="12">
                                <two-column-select primary-key-field="displayCategoryName"
                                                   primary-value-field="displayCategoryName"
                                                   primary-title-field="displayCategoryName"
                                                   primary-sub-items-field="fileTypes"
                                                   secondary-key-field="type"
                                                   secondary-value-field="type"
                                                   secondary-title-field="displayName"
                                                   :disabled="submitting"
                                                   :enable-filter="true"
                                                   :filter-placeholder="tt('Find file type')"
                                                   :filter-no-items-text="tt('No available file type')"
                                                   :label="tt('File Type')"
                                                   :placeholder="tt('File Type')"
                                                   :items="allSupportedImportFileCategoryAndTypes"
                                                   :auto-update-menu-position="true"
                                                   v-model="fileType">
                                </two-column-select>
                            </v-col>

                            <v-col cols="12" md="12" v-if="allFileSubTypes">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    :disabled="submitting"
                                    :label="tt('Format')"
                                    :placeholder="tt('Format')"
                                    :items="allFileSubTypes"
                                    v-model="fileSubType"
                                />
                            </v-col>

                            <v-col cols="12" md="12" v-if="!isImportDataFromTextbox && allSupportedEncodings">
                                <v-select
                                    item-title="displayName"
                                    item-value="encoding"
                                    :disabled="submitting"
                                    :label="tt('File Encoding')"
                                    :placeholder="tt('File Encoding')"
                                    :items="allSupportedEncodings"
                                    v-model="fileEncoding"
                                />
                            </v-col>

                            <v-col cols="12" md="12" v-if="fileType === 'dsv' || fileType === 'dsv_data'">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    :disabled="submitting"
                                    :label="tt('Handling Method')"
                                    :placeholder="tt('Handling Method')"
                                    :items="[
                                        { displayName: tt('Column Mapping'), type: ImportDSVProcessMethod.ColumnMapping },
                                        { displayName: tt('Custom Script'), type: ImportDSVProcessMethod.CustomScript }
                                     ]"
                                    v-model="processDSVMethod"
                                />
                            </v-col>

                            <v-col cols="12" md="12" v-if="supportedAdditionalOptions">
                                <v-select
                                    :disabled="submitting"
                                    :label="tt('Additional Options')"
                                    :placeholder="tt('Additional Options')"
                                    v-model="fileType"
                                    v-model:menu="additionalOptionsMenuState"
                                >
                                    <template #selection>
                                        <span class="cursor-pointer">{{ displaySelectedAdditionalOptions }}</span>
                                    </template>

                                    <template #no-data>
                                        <v-list class="py-0">
                                            <template v-for="item in allSupportedAdditionalOptions">
                                                <v-list-item :key="item.key"
                                                             :append-icon="importAdditionalOptions[item.key] ? mdiCheck : undefined"
                                                             @click="importAdditionalOptions[item.key] = !importAdditionalOptions[item.key]"
                                                             v-if="isDefined(supportedAdditionalOptions[item.key])">{{ tt(item.name) }}</v-list-item>
                                            </template>
                                        </v-list>
                                    </template>
                                </v-select>
                            </v-col>

                            <v-col cols="12" md="12" v-if="!isImportDataFromTextbox">
                                <v-text-field
                                    readonly
                                    persistent-placeholder
                                    type="text"
                                    class="always-cursor-pointer"
                                    :disabled="submitting"
                                    :label="tt('Data File')"
                                    :placeholder="tt('format.misc.clickToSelectedFile', { extensions: supportedImportFileExtensions })"
                                    v-model="fileName"
                                    @click="showOpenFileDialog"
                                />
                            </v-col>

                            <v-col cols="12" md="12" v-if="isImportDataFromTextbox">
                                <v-textarea
                                    type="text"
                                    persistent-placeholder
                                    rows="5"
                                    :disabled="submitting"
                                    :placeholder="tt('Data to import')"
                                    v-model="importData"
                                />
                            </v-col>

                            <v-col cols="12" md="12" v-if="exportFileGuideDocumentUrl">
                                <a :href="exportFileGuideDocumentUrl" :class="{ 'disabled': submitting }" target="_blank">
                                    <v-icon :icon="mdiHelpCircleOutline" size="16" />
                                    <span class="ms-1" v-if="fileType === 'dsv' || fileType === 'dsv_data'">{{ tt('How to import this file?') }}</span>
                                    <span class="ms-1" v-if="fileType !== 'dsv' && fileType !== 'dsv_data'">{{ tt('How to export this file?') }}</span>
                                    <span class="ms-1" v-if="exportFileGuideDocumentLanguageName">[{{ exportFileGuideDocumentLanguageName }}]</span>
                                </a>
                            </v-col>
                        </v-row>
                    </v-window-item>
                    <v-window-item value="defineColumn">
                        <import-transaction-define-column-tab
                            ref="importTransactionDefineColumnTab"
                            :parsed-file-data="parsedFileData"
                            :disabled="loading || submitting"
                        />
                    </v-window-item>
                    <v-window-item value="executeCustomScript">
                        <import-transaction-execute-custom-script-tab
                            ref="importTransactionExecuteCustomScriptTab"
                            :parsed-file-data="parsedFileData"
                            :disabled="loading || submitting"
                        />
                    </v-window-item>
                    <v-window-item value="checkData">
                        <import-transaction-check-data-tab
                            ref="importTransactionCheckDataTab"
                            :import-transactions="importTransactions"
                            :disabled="loading || submitting"
                        />
                    </v-window-item>
                    <v-window-item value="finalResult">
                        <h4 class="text-h4 mb-1">{{ tt('Data Import Completed') }}</h4>
                        <p class="my-5">{{ tt('format.misc.importTransactionResult', { count: getDisplayCount(importedCount || 0) }) }}</p>
                    </v-window-item>
                </v-window>
            </v-card-text>
            <v-card-text>
                <div class="d-flex justify-center justify-sm-space-between flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="secondary" variant="tonal" :disabled="loading || submitting"
                           :prepend-icon="mdiClose" @click="close(false)"
                           v-if="currentStep !== 'finalResult'">{{ tt('Cancel') }}</v-btn>
                    <v-btn class="button-icon-with-direction" color="primary"
                           :disabled="loading || submitting || (!isImportDataFromTextbox && !importFile) || (isImportDataFromTextbox && !importData) || (!isImportDataFromTextbox && allSupportedEncodings && fileEncoding === 'auto' && !autoDetectedFileEncoding)"
                           :append-icon="!submitting ? mdiArrowRight : undefined" @click="parseData"
                           v-if="currentStep === 'defineColumn' || currentStep === 'executeCustomScript' || currentStep === 'uploadFile'">
                        {{ tt('Next') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn class="button-icon-with-direction" color="teal"
                           :disabled="submitting || importTransactionCheckDataTab?.isEditing || !importTransactionCheckDataTab?.canImport"
                           :append-icon="!submitting ? mdiArrowRight : undefined" @click="submit"
                           v-if="currentStep === 'checkData'">
                        {{ (submitting && importProcess > 0 ? tt('format.misc.importingTransactions', { process: formatNumberToLocalizedNumerals(importProcess, 2) }) : tt('Import')) }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal"
                           :append-icon="mdiCheck"
                           @click="close(true)"
                           v-if="currentStep === 'finalResult'">{{ tt('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="fileInput" type="file" style="display: none" :accept="supportedImportFileExtensions" @change="setImportFile($event)" />
</template>

<script setup lang="ts">
import type { StepBarItem } from '@/components/desktop/StepsBar.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import ImportTransactionDefineColumnTab from './tabs/ImportTransactionDefineColumnTab.vue';
import ImportTransactionExecuteCustomScriptTab from './tabs/ImportTransactionExecuteCustomScriptTab.vue';
import ImportTransactionCheckDataTab from './tabs/ImportTransactionCheckDataTab.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { type KeyAndName, itemAndIndex } from '@/core/base.ts';
import { type NumeralSystem } from '@/core/numeral.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    type ImportFileTypeSupportedAdditionalOptions,
    type LocalizedImportFileCategoryAndTypes,
    type LocalizedImportFileType,
    type LocalizedImportFileTypeSubType,
    type LocalizedImportFileTypeSupportedEncodings,
    KnownFileType
} from '@/core/file.ts';
import { UTF_8 } from '@/consts/file.ts';

import { ImportTransaction } from '@/models/imported_transaction.ts';

import { isDefined, isNumber } from '@/lib/common.ts';
import { findExtensionByType, isFileExtensionSupported, detectFileEncoding } from '@/lib/file.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import logger from '@/lib/logger.ts';

import {
    mdiFilterOutline,
    mdiCheck,
    mdiDotsVertical,
    mdiHelpCircleOutline,
    mdiClose,
    mdiArrowRight
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type ImportTransactionDefineColumnTabType = InstanceType<typeof ImportTransactionDefineColumnTab>;
type ImportTransactionExecuteCustomScriptTabType = InstanceType<typeof ImportTransactionExecuteCustomScriptTab>;
type ImportTransactionCheckDataTabType = InstanceType<typeof ImportTransactionCheckDataTab>;

type ImportTransactionDialogStep = 'uploadFile' | 'defineColumn' | 'executeCustomScript' | 'checkData' | 'finalResult';
enum ImportDSVProcessMethod {
    ColumnMapping,
    CustomScript
};

defineProps<{
    persistent?: boolean;
}>();

const {
    tt,
    joinMultiText,
    getCurrentNumeralSystemType,
    getAllSupportedImportFileCagtegoryAndTypes,
    formatNumberToLocalizedNumerals,
    getLocalizedFileEncodingName
} = useI18n();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();
const overviewStore = useOverviewStore();
const statisticsStore = useStatisticsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const importTransactionDefineColumnTab = useTemplateRef<ImportTransactionDefineColumnTabType>('importTransactionDefineColumnTab');
const importTransactionExecuteCustomScriptTab = useTemplateRef<ImportTransactionExecuteCustomScriptTabType>('importTransactionExecuteCustomScriptTab');
const importTransactionCheckDataTab = useTemplateRef<ImportTransactionCheckDataTabType>('importTransactionCheckDataTab');
const fileInput = useTemplateRef<HTMLInputElement>('fileInput');

const allSupportedAdditionalOptions: KeyAndName[] = [
    {
        key: 'payeeAsTag',
        name: 'Parse Payee as Tag'
    },
    {
        key: 'payeeAsDescription',
        name: 'Parse Payee as Description'
    },
    {
        key: 'memberAsTag',
        name: 'Parse Member as Tag'
    },
    {
        key: 'projectAsTag',
        name: 'Parse Project as Tag'
    },
    {
        key: 'merchantAsTag',
        name: 'Parse Merchant as Tag'
    }
];

const showState = ref<boolean>(false);
const additionalOptionsMenuState = ref<boolean>(false);
const clientSessionId = ref<string>('');
const currentStep = ref<ImportTransactionDialogStep>('uploadFile');
const importProcess = ref<number>(0);
const fileType = ref<string>('ezbookkeeping');
const fileSubType = ref<string>('ezbookkeeping_csv');
const fileEncoding = ref<string>('auto');
const detectingFileEncoding = ref<boolean>(false);
const autoDetectedFileEncoding = ref<string | undefined>(undefined);
const processDSVMethod = ref<ImportDSVProcessMethod>(ImportDSVProcessMethod.ColumnMapping);
const importFile = ref<File | null>(null);
const importData = ref<string>('');
const importAdditionalOptions = ref<ImportFileTypeSupportedAdditionalOptions>({});
const parsedFileData = ref<string[][] | undefined>(undefined);
const importTransactions = ref<ImportTransaction[] | undefined>(undefined);

const importedCount = ref<number | null>(null);
const loading = ref<boolean>(true);
const submitting = ref<boolean>(false);

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

const allSupportedImportFileCategoryAndTypes = computed<LocalizedImportFileCategoryAndTypes[]>(() => getAllSupportedImportFileCagtegoryAndTypes());
const allFileSubTypes = computed<LocalizedImportFileTypeSubType[] | undefined>(() => allSupportedImportFileTypesMap.value[fileType.value]?.subTypes);
const allSupportedEncodings = computed<LocalizedImportFileTypeSupportedEncodings[] | undefined>(() => {
    const supportedEncodings = allSupportedImportFileTypesMap.value[fileType.value]?.supportedEncodings;

    if (!supportedEncodings) {
        return undefined;
    }

    const ret: LocalizedImportFileTypeSupportedEncodings[] = [];
    let autoDetectDisplayName = tt('Auto detect');

    if (importFile.value) {
        if (detectingFileEncoding.value) {
            autoDetectDisplayName += ` [${tt('Detecting...')}]`;
        } else if (autoDetectedFileEncoding.value) {
            autoDetectDisplayName += ` [${getLocalizedFileEncodingName(autoDetectedFileEncoding.value)}]`;
        } else {
            autoDetectDisplayName += ` [${tt('Unknown')}]`;
        }
    }

    const autoDetectEncoding: LocalizedImportFileTypeSupportedEncodings = {
        displayName: autoDetectDisplayName,
        encoding: 'auto'
    };

    ret.push(autoDetectEncoding);

    if (supportedEncodings && supportedEncodings.length) {
        ret.push(...supportedEncodings);
    }

    return ret;
});
const isImportDataFromTextbox = computed<boolean>(() => allSupportedImportFileTypesMap.value[fileType.value]?.dataFromTextbox ?? false);
const supportedAdditionalOptions = computed<ImportFileTypeSupportedAdditionalOptions | undefined>(() => allSupportedImportFileTypesMap.value[fileType.value]?.supportedAdditionalOptions);

const allSteps = computed<StepBarItem[]>(() => {
    const steps: StepBarItem[] = [
        {
            name: 'uploadFile',
            title: tt('Upload File'),
            subTitle: tt('Upload Transaction Data File')
        }
    ];

    if (fileType.value === 'dsv' || fileType.value === 'dsv_data') {
        if (processDSVMethod.value === ImportDSVProcessMethod.CustomScript) {
            steps.push({
                name: 'executeCustomScript',
                title: tt('Execute Custom Script'),
                subTitle: tt('Execute Custom Script to Parse Data')
            });
        } else {
            steps.push({
                name: 'defineColumn',
                title: tt('Define Column'),
                subTitle: tt('Define and Check Column Mapping')
            });
        }
    }

    steps.push(...[
        {
            name: 'checkData',
            title: tt('Check & Modify'),
            subTitle: tt('Check and Modify Your Data')
        },
        {
            name: 'finalResult',
            title: tt('Complete'),
            subTitle: tt('Data Import Completed')
        }
    ]);

    return steps;
});

const allSupportedImportFileTypesMap = computed<Record<string, LocalizedImportFileType>>(() => {
    const ret: Record<string, LocalizedImportFileType> = {};

    for (const importFileCategoryAndTypes of allSupportedImportFileCategoryAndTypes.value) {
        for (const importFileType of importFileCategoryAndTypes.fileTypes) {
            ret[importFileType.type] = importFileType;
        }
    }

    return ret;
});

const supportedImportFileExtensions = computed<string | undefined>(() => {
    if (allFileSubTypes.value && allFileSubTypes.value.length) {
        const subTypeExtensions = findExtensionByType(allFileSubTypes.value, fileSubType.value);

        if (subTypeExtensions) {
            return subTypeExtensions;
        }
    }

    return allSupportedImportFileTypesMap.value[fileType.value]?.extensions;
});

const displaySelectedAdditionalOptions = computed<string>(() => {
    if (!supportedAdditionalOptions.value) {
        return tt('None');
    }

    const selectedOptions: string[] = [];

    for (const option of allSupportedAdditionalOptions) {
        if (isDefined(supportedAdditionalOptions.value[option.key]) && importAdditionalOptions.value[option.key]) {
            selectedOptions.push(tt(option.name));
        }
    }

    if (selectedOptions.length < 1) {
        return tt('None');
    }

    return joinMultiText(selectedOptions);
});

const exportFileGuideDocumentUrl = computed<string | undefined>(() => {
    const document = allSupportedImportFileTypesMap.value[fileType.value]?.document;

    if (!document) {
        return undefined;
    }

    const language = document.language ? document.language + '/' : '';
    const anchor = document.anchor ? '#' + document.anchor : '';
    return `https://ezbookkeeping.mayswind.net/${language}export_and_import${anchor}`;
});

const exportFileGuideDocumentLanguageName = computed<string | undefined>(() => allSupportedImportFileTypesMap.value[fileType.value]?.document?.displayLanguageName);

const fileName = computed<string>(() => importFile.value?.name || '');

function getDisplayCount(count: number): string {
    return numeralSystem.value.formatNumber(count);
}

function open(): Promise<void> {
    fileType.value = 'ezbookkeeping';
    fileSubType.value = 'ezbookkeeping_csv';
    fileEncoding.value = 'auto';
    detectingFileEncoding.value = false;
    autoDetectedFileEncoding.value = undefined;
    processDSVMethod.value = ImportDSVProcessMethod.ColumnMapping;
    currentStep.value = 'uploadFile';
    importProcess.value = 0;
    importFile.value = null;
    importData.value = '';
    importAdditionalOptions.value = Object.assign({}, supportedAdditionalOptions.value ?? {});
    parsedFileData.value = undefined;
    importTransactionDefineColumnTab.value?.reset();
    importTransactionExecuteCustomScriptTab.value?.reset();
    importTransactions.value = undefined;
    importTransactionCheckDataTab.value?.reset();
    showState.value = true;
    clientSessionId.value = generateRandomUUID();

    const promises = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ];

    Promise.all(promises).then(() => {
        loading.value = false;
    }).catch(error => {
        logger.error('failed to load essential data for importing transaction', error);

        loading.value = false;
        showState.value = false;

        if (!error.processed) {
            if (rejectFunc) {
                rejectFunc(error);
            }
        }
    });

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function showOpenFileDialog(): void {
    if (submitting.value) {
        return;
    }

    fileInput.value?.click();
}

function setImportFile(event: Event): void {
    if (!event || !event.target) {
        return;
    }

    const el = event.target as HTMLInputElement;

    if (!el.files || !el.files.length || !el.files[0]) {
        return;
    }

    importFile.value = el.files[0] as File;
    detectingFileEncoding.value = false;
    autoDetectedFileEncoding.value = undefined;
    el.value = '';

    if (allSupportedEncodings.value) {
        detectingFileEncoding.value = true;

        detectFileEncoding(importFile.value).then(detectedEncoding => {
            detectingFileEncoding.value = false;
            autoDetectedFileEncoding.value = detectedEncoding;
        }).catch(() => {
            detectingFileEncoding.value = false;
            autoDetectedFileEncoding.value = undefined;
        });
    }
}

function parseData(): void {
    let uploadFile: File;
    let type: string = fileType.value;
    let encoding: string | undefined = undefined;

    if (allFileSubTypes.value) {
        type = fileSubType.value;
    }

    if (allSupportedEncodings.value) {
        if (fileEncoding.value === 'auto') {
            encoding = autoDetectedFileEncoding.value;
        } else {
            encoding = fileEncoding.value;
        }
    }

    if (!isImportDataFromTextbox.value) {
        if (!importFile.value) {
            snackbar.value?.showError('Please select a file to import');
            return;
        }

        if (allSupportedEncodings.value) {
            if (fileEncoding.value === 'auto' && !autoDetectedFileEncoding.value) {
                snackbar.value?.showError('Unable to detect the file encoding automatically. Please select the actual encoding.');
                return;
            }
        }

        uploadFile = importFile.value;
    } else if (isImportDataFromTextbox.value) {
        if (!importData.value) {
            snackbar.value?.showError('No data to import');
            return;
        }

        if (type === 'custom_csv') {
            uploadFile = KnownFileType.CSV.createFile(importData.value, 'import');
        } else if (type === 'custom_tsv') {
            uploadFile = KnownFileType.TSV.createFile(importData.value, 'import');
        } else {
            snackbar.value?.showError('Parameter Invalid');
            return;
        }

        encoding = UTF_8;
    } else { // should not happen, but ts would check whether uploadFile has been assigned a value
        snackbar.value?.showMessage('An error occurred');
        return;
    }

    const isDsvFileType: boolean = fileType.value === 'dsv' || fileType.value === 'dsv_data';

    if (isDsvFileType && currentStep.value === 'uploadFile') {
        submitting.value = true;

        transactionsStore.parseImportDsvFile({
            fileType: type,
            fileEncoding: encoding,
            importFile: uploadFile
        }).then(response => {
            if (response && response.length) {
                if (processDSVMethod.value === ImportDSVProcessMethod.CustomScript) {
                    importTransactionExecuteCustomScriptTab.value?.reset();
                    parsedFileData.value = response;
                    currentStep.value = 'executeCustomScript';
                } else {
                    importTransactionDefineColumnTab.value?.reset();
                    parsedFileData.value = response;
                    currentStep.value = 'defineColumn';
                }
            } else {
                parsedFileData.value = undefined;
                snackbar.value?.showError('No data to import');
            }

            submitting.value = false;
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else {
        let columnMapping: Record<number, number> | undefined = undefined;
        let transactionTypeMapping: Record<string, TransactionType> | undefined = undefined;
        let hasHeaderLine: boolean | undefined = undefined;
        let timeFormat: string | undefined = undefined;
        let timezoneFormat: string | undefined = undefined;
        let amountDecimalSeparator: string | undefined = undefined;
        let amountDigitGroupingSymbol: string | undefined = undefined;
        let geoLocationSeparator: string | undefined = undefined;
        let geoLocationOrder: string | undefined = undefined;
        let tagSeparator: string | undefined = undefined;

        if (isDsvFileType && processDSVMethod.value === ImportDSVProcessMethod.ColumnMapping) {
            const defineColumnResult = importTransactionDefineColumnTab.value?.generateResult();

            if (!defineColumnResult) {
                return;
            }

            columnMapping = defineColumnResult.columnMapping;
            transactionTypeMapping = defineColumnResult.transactionTypeMapping;
            hasHeaderLine = defineColumnResult.includeHeader;
            timeFormat = defineColumnResult.timeFormat;
            timezoneFormat = defineColumnResult.timezoneFormat;
            amountDecimalSeparator = defineColumnResult.amountDecimalSeparator;
            amountDigitGroupingSymbol = defineColumnResult.amountDigitGroupingSymbol;
            geoLocationSeparator = defineColumnResult.geoLocationSeparator;
            geoLocationOrder = defineColumnResult.geoLocationOrder;
            tagSeparator = defineColumnResult.tagSeparator;
        } else if (isDsvFileType && processDSVMethod.value === ImportDSVProcessMethod.CustomScript) {
            const executeCustomScriptResult = importTransactionExecuteCustomScriptTab.value?.generateResult();

            if (!executeCustomScriptResult) {
                return;
            }

            type = 'ezbookkeeping_json';
            encoding = undefined;
            uploadFile = KnownFileType.JSON.createFile(executeCustomScriptResult, 'import');
        }

        submitting.value = true;

        transactionsStore.parseImportTransaction({
            fileType: type,
            additionalOptions: importAdditionalOptions.value,
            fileEncoding: encoding,
            importFile: uploadFile,
            columnMapping: columnMapping,
            transactionTypeMapping: transactionTypeMapping,
            hasHeaderLine: hasHeaderLine,
            timeFormat: timeFormat,
            timezoneFormat: timezoneFormat,
            amountDecimalSeparator: amountDecimalSeparator,
            amountDigitGroupingSymbol: amountDigitGroupingSymbol,
            geoSeparator: geoLocationSeparator,
            geoOrder: geoLocationOrder,
            tagSeparator: tagSeparator
        }).then(response => {
            const parsedTransactions: ImportTransaction[] = [];

            if (response.items) {
                for (const [importTransaction, index] of itemAndIndex(response.items)) {
                    const parsedTransaction = ImportTransaction.of(importTransaction, index);
                    parsedTransactions.push(parsedTransaction);
                }
            }

            importTransactionCheckDataTab.value?.reset();

            if (parsedTransactions && parsedTransactions.length >= 0 && parsedTransactions.length < 10) {
                importTransactionCheckDataTab.value?.setCountPerPage(-1);
            } else {
                importTransactionCheckDataTab.value?.setCountPerPage(10);
            }

            importTransactions.value = parsedTransactions;
            currentStep.value = 'checkData';
            submitting.value = false;
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }
}

function submit(): void {
    if (importTransactionCheckDataTab.value?.isEditing) {
        return;
    }

    const transactions: ImportTransaction[] = [];

    if (importTransactions.value) {
        for (const importTransaction of importTransactions.value) {
            if (importTransaction.valid && importTransaction.selected) {
                transactions.push(importTransaction);
            } else if (!importTransaction.valid && importTransaction.selected) {
                snackbar.value?.showError('Cannot import invalid transactions');
                return;
            }
        }
    }

    if (transactions.length < 1) {
        snackbar.value?.showError('No data to import');
        return;
    }

    confirmDialog.value?.open('format.misc.confirmImportTransactions', {
        count: getDisplayCount(transactions.length)
    }).then(() => {
        submitting.value = true;

        let showProcessTimer : number | undefined = undefined;

        if (transactions.length > 100) {
            setTimeout(() => {
                if (!submitting.value) {
                    logger.warn('transaction import is not submitting');
                    return;
                }

                // @ts-expect-error the return value of setInterval is number, but lint shows it as NodeJS.Timer
                showProcessTimer = setInterval(() => {
                    if (submitting.value) {
                        transactionsStore.getImportTransactionsProcess({
                            clientSessionId: clientSessionId.value
                        }).then(response => {
                            if (isNumber(response) && 0 <= response && response < 100) {
                                importProcess.value = response;
                            } else {
                                importProcess.value = 0;
                                clearInterval(showProcessTimer);
                                showProcessTimer = undefined;
                            }
                        }).catch(() => {
                            importProcess.value = 0;
                            clearInterval(showProcessTimer);
                            showProcessTimer = undefined;
                        });
                    }
                }, 2000);
            }, 2000);
        }

        transactionsStore.importTransactions({
            transactions: transactions,
            clientSessionId: clientSessionId.value
        }).then(response => {
            if (showProcessTimer) {
                importProcess.value = 0;
                clearInterval(showProcessTimer);
                showProcessTimer = undefined;
            }

            importedCount.value = response;
            currentStep.value = 'finalResult';

            accountsStore.updateAccountListInvalidState(true);
            transactionsStore.updateTransactionListInvalidState(true);
            overviewStore.updateTransactionOverviewInvalidState(true);
            statisticsStore.updateTransactionStatisticsInvalidState(true);

            submitting.value = false;
        }).catch(error => {
            if (showProcessTimer) {
                importProcess.value = 0;
                clearInterval(showProcessTimer);
                showProcessTimer = undefined;
            }

            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function close(completed: boolean): void {
    if (completed) {
        if (resolveFunc) {
            resolveFunc();
        }
    } else {
        if (rejectFunc) {
            rejectFunc();
        }
    }

    showState.value = false;
}

watch(fileType, () => {
    if (allFileSubTypes.value && allFileSubTypes.value.length) {
        fileSubType.value = allFileSubTypes.value[0]!.type;
    }

    importFile.value = null;
    parsedFileData.value = undefined;
    importAdditionalOptions.value = Object.assign({}, supportedAdditionalOptions.value ?? {});
    importTransactions.value = undefined;
});

watch(fileSubType, (newValue) => {
    let supportedExtensions: string | undefined = findExtensionByType(allFileSubTypes.value, newValue);

    if (!supportedExtensions) {
        supportedExtensions = allSupportedImportFileTypesMap.value[fileType.value]?.extensions;
    }

    if (importFile.value && importFile.value.name && !isFileExtensionSupported(importFile.value.name, supportedExtensions || '')) {
        importFile.value = null;
    }
});

defineExpose({
    open
});
</script>
