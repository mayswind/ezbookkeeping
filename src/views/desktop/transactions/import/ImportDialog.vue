<template>
    <v-dialog :persistent="!!persistent || loading || loadingPreset || submitting || currentStep !== 'uploadFile' || isModifiedInUploadStep" v-model="showState">
        <one-column-dialog-layout :class="{ 'import-transaction-dialog-layout': currentStep === 'defineColumn' || currentStep === 'checkData' }"
                                  :content-class="currentStep === 'defineColumn' || currentStep === 'checkData' ? 'pa-0 import-transaction-dialog-content' : 'pa-0'"
                                  :title="tt('Import Transactions')" :cancel-button-title="currentStep === 'finalResult' ? tt('Close') : tt('Cancel')"
                                  :disabled="loading || loadingPreset || submitting"
                                  :loading="currentStep !== 'checkData' && (loading || loadingPreset)"
                                  @cancel="close(currentStep === 'finalResult')">
            <template #after-title>
                <div class="ms-2 text-body-large text-truncate" v-if="currentStep === 'uploadFile' && exportFileGuideDocumentUrl">
                    <a :class="{ 'disabled': submitting }" :href="exportFileGuideDocumentUrl" target="_blank">
                        <v-icon :icon="mdiHelpCircleOutline" size="16" />
                        <span class="vertical-align-middle ms-1" v-if="isCustomFileFormat">{{ tt('How to import this file?') }}</span>
                        <span class="vertical-align-middle ms-1" v-if="!isCustomFileFormat">{{ tt('How to export this file?') }}</span>
                        <span class="vertical-align-middle ms-1" v-if="exportFileGuideDocumentLanguageName">[{{ exportFileGuideDocumentLanguageName }}]</span>
                    </a>
                </div>

                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('Refresh Accounts, Categories and Tags')" :icon="true"
                       :disabled="loading || loadingPreset"
                       :loading="loading || loadingPreset"
                       @click="reloadBasisData"
                       v-if="currentStep === 'checkData'">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="22" />
                    <v-tooltip activator="parent">{{ tt('Refresh Accounts, Categories and Tags') }}</v-tooltip>
                </v-btn>

                <v-btn class="ms-2" density="comfortable" color="primary" variant="outlined"
                       :disabled="loading || loadingPreset || submitting || importTransactionCheckDataTab?.isEditing"
                       @click="batchApplyRules"
                       v-if="currentStep === 'checkData'">
                    {{ tt('Batch Apply Rules') }}
                </v-btn>
            </template>

            <template #toolbar>
                <v-btn class="ms-2 me-1" density="comfortable" variant="outlined" color="primary"
                       :disabled="loading || loadingPreset || submitting || (!isImportDataFromTextbox && !isAIImageImport && !importFile) || (isImportDataFromTextbox && !importData) || (!isImportDataFromTextbox && !isAIImageImport && allSupportedEncodings && fileEncoding === 'auto' && !autoDetectedFileEncoding) || (!isImportDataFromTextbox && isAIImageImport && importImageFiles.length === 0)"
                       @click="parseData"
                       v-if="currentStep === 'defineColumn' || currentStep === 'executeCustomScript' || (currentStep === 'recognizeImages' && !submitting) || currentStep === 'uploadFile'">
                    {{ tt('Next') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>

                <v-btn class="ms-2 me-1" density="comfortable" variant="outlined" color="secondary"
                       :disabled="loading || loadingPreset"
                       @click="cancelBatchRecognizeImages()"
                       v-if="currentStep === 'recognizeImages' && submitting">{{ tt('Cancel Recognition') }}</v-btn>

                <v-btn class="ms-2" density="comfortable" variant="outlined"
                       :disabled="submitting || importTransactionCheckDataTab?.isEditing || !importTransactionCheckDataTab?.canImport"
                       @click="submit"
                       v-if="currentStep === 'checkData'">
                    {{ (submitting && importProcess > 0 ? tt('format.misc.importingTransactions', { process: formatNumberToLocalizedNumerals(importProcess, 2) }) : tt('Import')) }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>
                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('More')" :disabled="loading || loadingPreset || submitting" :icon="true"
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
                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('More')" :disabled="loading || loadingPreset || submitting" :icon="true"
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
                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('Filter')" :icon="true" :disabled="loading || loadingPreset || submitting || importTransactionCheckDataTab?.isEditing"
                       v-if="currentStep === 'checkData' && importTransactionCheckDataTab?.filterMenus">
                    <v-icon :icon="mdiFilterOutline" />
                    <v-menu activator="parent" max-height="500">
                        <v-list>
                            <template :key="groupIndex" v-for="(group, groupIndex) in importTransactionCheckDataTab.filterMenus">
                                <v-divider class="my-2" v-if="groupIndex > 0" />
                                <v-list-subheader class="text-body-small" :title="group.title" v-if="group.title" />
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
                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('More')" :icon="true" :disabled="loading || loadingPreset || submitting"
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
            </template>

            <template #subtitle>
                <v-divider class="mt-2" />
                <div class="cursor-default mt-3 mx-3 mb-md-2">
                    <steps-bar min-width="700" :always-horizontal="true" :clickable="false"
                               :steps="allSteps" :current-step="currentStep" />
                </div>
            </template>

            <template #content>
                <v-window class="disable-tab-transition"
                          :class="{ 'import-transaction-dialog-window': currentStep === 'defineColumn' || currentStep === 'checkData' }"
                          v-model="currentStep">
                    <v-window-item value="uploadFile">
                        <div class="pa-4">
                            <v-row>
                                <v-col cols="12" md="12">
                                    <two-column-select primary-key-field="displayCategoryName"
                                                       primary-value-field="displayCategoryName"
                                                       primary-title-field="displayCategoryName"
                                                       primary-sub-items-field="fileTypes"
                                                       secondary-key-field="type"
                                                       secondary-value-field="type"
                                                       secondary-title-field="displayName"
                                                       :disabled="loading || loadingPreset || submitting"
                                                       :enable-filter="true"
                                                       :filter-placeholder="tt('Find file type')"
                                                       :filter-no-items-text="tt('No available file type')"
                                                       :label="tt('File Type')"
                                                       :placeholder="tt('File Type')"
                                                       :items="allPresetsAndSupportedImportFileCategoryAndTypes"
                                                       :auto-update-menu-position="true"
                                                       :disable-auto-select-first-secondary-item="true"
                                                       v-model="selectedImportType">
                                    </two-column-select>
                                </v-col>

                                <import-transaction-options
                                    :file-type-info="allSupportedImportFileTypesMap[fileType]"
                                    :supported-encodings="allSupportedEncodings"
                                    :disabled="loading || loadingPreset || submitting"
                                    :disabled-sub-type="!!selectedPreset"
                                    :disabled-handling-method="!!selectedPreset"
                                    v-model:file-type="fileType"
                                    v-model:file-sub-type="fileSubType"
                                    v-model:file-encoding="fileEncoding"
                                    v-model:handling-method="processCustomFileFormatMethod"
                                    v-model:additional-options="importAdditionalOptions"
                                    v-model:additional-prompt="importAIAdditionalPrompt"
                                />

                                <v-col cols="12" md="12" v-if="!isImportDataFromTextbox && !isAIImageImport">
                                    <v-text-field
                                        readonly
                                        persistent-placeholder
                                        type="text"
                                        class="always-cursor-pointer"
                                        :disabled="loading || loadingPreset || submitting"
                                        :label="tt('Data File')"
                                        :placeholder="tt('format.misc.clickToSelectedFile', { extensions: supportedImportFileExtensions })"
                                        v-model="fileName"
                                        @click="showOpenFileDialog"
                                    />
                                </v-col>

                                <v-col cols="12" md="12" v-if="isImportDataFromTextbox">
                                    <v-textarea
                                        type="text"
                                        autocomplete="off"
                                        persistent-placeholder
                                        rows="5"
                                        :disabled="loading || loadingPreset || submitting"
                                        :placeholder="tt('Data to import')"
                                        v-model="importData"
                                    />
                                </v-col>

                                <v-col cols="12" md="12" v-if="isAIImageImport">
                                    <div class="text-body-large mb-2">{{ tt('Image Files') }}</div>
                                    <div class="import-transaction-images d-flex gap-2 overflow-x-auto">
                                        <div :key="picIdx" v-for="(imageItem, picIdx) in importImageFiles">
                                            <v-avatar rounded="lg" variant="tonal" size="120"
                                                      class="cursor-pointer import-image"
                                                      color="rgba(0,0,0,0)" @click="removeImportImageFile(picIdx)">
                                                <v-img :src="imageItem.previewUrl"></v-img>
                                                <div class="picture-control-icon">
                                                    <v-icon size="48" :aria-label="tt('Remove Picture')" :icon="mdiTrashCanOutline" />
                                                </div>
                                            </v-avatar>
                                        </div>
                                        <div>
                                            <v-avatar rounded="lg" variant="tonal" size="120"
                                                      class="import-image import-image-add"
                                                      :class="{ 'enabled': !submitting, 'cursor-pointer': !submitting }"
                                                      color="rgba(0,0,0,0)" @click="showOpenFileDialog">
                                                <v-tooltip activator="parent" v-if="!submitting">{{ tt('Add Picture') }}</v-tooltip>
                                                <v-icon class="import-image-add-icon" size="48" :aria-label="tt('Add Picture')" :icon="mdiImagePlusOutline" />
                                            </v-avatar>
                                        </div>
                                    </div>
                                </v-col>
                            </v-row>
                        </div>
                    </v-window-item>
                    <v-window-item value="defineColumn">
                        <import-transaction-define-column-tab
                            ref="importTransactionDefineColumnTab"
                            :parsed-file-data="parsedFileData"
                            :parsed-file-column-separator="parsedFileColumnSeparator"
                            :disabled="loading || loadingPreset || submitting"
                        />
                    </v-window-item>
                    <v-window-item value="executeCustomScript">
                        <import-transaction-execute-custom-script-tab
                            ref="importTransactionExecuteCustomScriptTab"
                            :parsed-file-data="parsedFileData"
                            :disabled="loading || loadingPreset || submitting"
                        />
                    </v-window-item>
                    <v-window-item value="recognizeImages">
                        <import-transaction-recognize-images-tab
                            ref="importTransactionRecognizeImagesTab"
                            :disabled="loading || loadingPreset || submitting"
                            :submitting="submitting"
                            :import-images="importImageFiles"
                        />
                    </v-window-item>
                    <v-window-item value="checkData">
                        <import-transaction-check-data-tab
                            ref="importTransactionCheckDataTab"
                            :import-transactions="importTransactions"
                            :disabled="loading || loadingPreset || submitting"
                        />
                    </v-window-item>
                    <v-window-item value="finalResult">
                        <div class="mx-4 my-4">
                            <v-alert class="my-4" type="warning" variant="tonal"
                                     v-if="importResultUnknown">{{ tt('The import result could not be confirmed. Check the result before importing again to avoid duplicate transactions.') }}</v-alert>
                            <v-alert class="my-4" type="error" variant="tonal"
                                     v-else-if="importErrorMessage">{{ tt('format.misc.importTransactionStoppedWithCause', { cause: importErrorMessage }) }}</v-alert>
                            <v-alert class="my-4" type="success" color="success-darken-1" variant="tonal"
                                     v-else-if="!importResultUnknown && !importErrorMessage">{{ tt('Data Import Completed') }}</v-alert>

                            <div class="text-body-large"
                                 v-if="!importResultUnknown">{{ tt('format.misc.importTransactionBatchResult', { count: formatNumberToLocalizedNumerals(importedCount || 0), totalCount: formatNumberToLocalizedNumerals(totalImportCount) }) }}</div>
                            <div class="text-body-large"
                                 v-if="!importResultUnknown && remainingTransactions.length">{{ tt('Would you like to view the transactions that have not been imported? Transactions already imported will be kept.') }}</div>

                            <div class="mt-4">
                                <v-btn class="me-2" color="primary" variant="tonal"
                                       :loading="submitting" :disabled="loading || loadingPreset || submitting"
                                       @click="checkImportResult(true)"
                                       v-if="importResultUnknown">
                                    <template #loader>
                                        <v-progress-circular indeterminate size="20"/>
                                    </template>
                                    {{ tt('Check Import Result') }}
                                </v-btn>
                                <v-btn class="me-2" color="primary" variant="tonal"
                                       @click="showRemainingTransactions"
                                       v-else-if="remainingTransactions.length">{{ tt('View Transactions Not Imported') }}</v-btn>
                                <v-btn color="secondary" variant="tonal"
                                       :disabled="loading || loadingPreset || submitting"
                                       @click="close(true)">{{ tt('Close') }}</v-btn>
                            </div>
                        </div>
                    </v-window-item>
                </v-window>
            </template>

            <template #footer v-if="currentStep === 'uploadFile' && (needAITextRecognition || needAIImageRecognition)">
                <div class="w-100 text-body-large text-truncate">
                    <span v-if="needAITextRecognition">
                        {{ tt('Uploaded text and personal data will be sent to the large language model, please be aware of potential privacy risks.') }}
                    </span>
                    <span v-if="needAIImageRecognition">
                        {{ tt('Uploaded image and personal data will be sent to the large language model, please be aware of potential privacy risks.') }}
                    </span>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="fileInput" type="file" style="display: none" :accept="supportedImportFileExtensions" :multiple="isAIImageImport" @change="setImportFile($event)" />
</template>

<script setup lang="ts">
import type { StepBarItem } from '@/components/desktop/StepsBar.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import ImportTransactionOptions from '@/components/desktop/ImportTransactionOptions.vue';
import ImportTransactionDefineColumnTab from './tabs/ImportTransactionDefineColumnTab.vue';
import ImportTransactionExecuteCustomScriptTab from './tabs/ImportTransactionExecuteCustomScriptTab.vue';
import ImportTransactionRecognizeImagesTab, { type BatchImportImageItem } from './tabs/ImportTransactionRecognizeImagesTab.vue';
import ImportTransactionCheckDataTab from './tabs/ImportTransactionCheckDataTab.vue';

import { ref, computed, useTemplateRef, watch, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useImportPresetsStore } from '@/stores/importPreset.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { itemAndIndex } from '@/core/base.ts';
import { TransactionType } from '@/core/transaction.ts';
import { KnownFileType } from '@/core/file.ts';
import {
    type ImportFileTypeSupportedAdditionalOptions,
    type LocalizedImportFileCategoryAndTypes,
    type LocalizedImportFileType,
    type LocalizedImportFileTypeSubType,
    type LocalizedImportFileTypeSupportedEncodings,
    ImportCustomFileFormatProcessMethod
} from '@/core/import_transaction.ts';
import { ImageUploadQualityType } from '@/core/image.ts';
import type { ErrorResponse } from '@/core/api.ts';
import { IMPORT_FILE_PRESET_TYPE } from '@/consts/import_transaction.ts';
import { UTF_8 } from '@/consts/file.ts';

import type { ImportPresetBasicInfo, ImportPreset } from '@/models/import_preset.ts';
import { type ImportTransactionReplaceRule, ImportTransactionReplaceRules } from '@/models/rule.ts';
import type { TransactionImportResponse } from '@/models/transaction.ts';
import { type ImportTransactionResponse, ImportTransaction } from '@/models/imported_transaction.ts';

import { isDefined, isInteger, isBoolean } from '@/lib/common.ts';
import { findExtensionByType, isFileExtensionSupported, detectFileEncoding } from '@/lib/file.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import { isTransactionFromAITextRecognitionEnabled, isTransactionFromAIImageRecognitionEnabled } from '@/lib/server_settings.ts';
import { compressJpgImageByQuality } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

import {
    mdiRefresh,
    mdiFilterOutline,
    mdiDotsVertical,
    mdiHelpCircleOutline,
    mdiImagePlusOutline,
    mdiTrashCanOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type ImportTransactionDefineColumnTabType = InstanceType<typeof ImportTransactionDefineColumnTab>;
type ImportTransactionExecuteCustomScriptTabType = InstanceType<typeof ImportTransactionExecuteCustomScriptTab>;
type ImportTransactionRecognizeImagesTabType = InstanceType<typeof ImportTransactionRecognizeImagesTab>;
type ImportTransactionCheckDataTabType = InstanceType<typeof ImportTransactionCheckDataTab>;
type ImportTransactionDialogStep = 'uploadFile' | 'defineColumn' | 'executeCustomScript' | 'recognizeImages' | 'checkData' | 'finalResult';

defineProps<{
    persistent?: boolean;
}>();

const {
    tt,
    te,
    getAllSupportedImportFileCagtegoryAndTypes,
    formatNumberToLocalizedNumerals,
    getLocalizedFileEncodingName
} = useI18n();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();
const importPresetsStore = useImportPresetsStore();
const overviewStore = useOverviewStore();
const statisticsStore = useStatisticsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const importTransactionDefineColumnTab = useTemplateRef<ImportTransactionDefineColumnTabType>('importTransactionDefineColumnTab');
const importTransactionExecuteCustomScriptTab = useTemplateRef<ImportTransactionExecuteCustomScriptTabType>('importTransactionExecuteCustomScriptTab');
const importTransactionRecognizeImagesTab = useTemplateRef<ImportTransactionRecognizeImagesTabType>('importTransactionRecognizeImagesTab');
const importTransactionCheckDataTab = useTemplateRef<ImportTransactionCheckDataTabType>('importTransactionCheckDataTab');
const fileInput = useTemplateRef<HTMLInputElement>('fileInput');

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const clientSessionId = ref<string>('');
const currentStep = ref<ImportTransactionDialogStep>('uploadFile');
const importProcess = ref<number>(0);
const fileType = ref<string>('ezbookkeeping');
const fileSubType = ref<string>('ezbookkeeping_csv');
const fileEncoding = ref<string>('auto');
const detectingFileEncoding = ref<boolean>(false);
const autoDetectedFileEncoding = ref<string | undefined>(undefined);
const processCustomFileFormatMethod = ref<ImportCustomFileFormatProcessMethod>(ImportCustomFileFormatProcessMethod.ColumnMapping);

const selectedPreset = ref<ImportPreset | undefined>(undefined);
const presetDataMapping = ref<string>('');
const presetCustomScript = ref<string>('');
const presetAutoExecuteScript = ref<boolean>(false);
const presetReplaceRules = ref<string>('');
const presetRulesApplied = ref<boolean>(false);

const importFile = ref<File | null>(null);
const importImageFiles = ref<BatchImportImageItem[]>([]);
const importData = ref<string>('');
const importAdditionalOptions = ref<ImportFileTypeSupportedAdditionalOptions>({});
const importAIAdditionalPrompt = ref<string>('');
const importImageCancelRecognizingUuid = ref<string | undefined>(undefined);
const parsedFileData = ref<string[][] | undefined>(undefined);
const parsedFileColumnSeparator = ref<string | undefined>(undefined);
const importTransactions = ref<ImportTransaction[] | undefined>(undefined);

const importedCount = ref<number | null>(null);
const totalImportCount = ref<number>(0);
const importErrorMessage = ref<string>('');
const importResultUnknown = ref<boolean>(false);
const submittedTransactions = ref<ImportTransaction[]>([]);
const remainingTransactions = ref<ImportTransaction[]>([]);
const hasImportedTransactions = ref<boolean>(false);
const loading = ref<boolean>(true);
const loadingPreset = ref<boolean>(false);
const submitting = ref<boolean>(false);

const allPresets = computed<ImportPresetBasicInfo[]>(() => importPresetsStore.allImportPresetBasicInfos);
const allSupportedImportFileCategoryAndTypes = computed<LocalizedImportFileCategoryAndTypes[]>(() => getAllSupportedImportFileCagtegoryAndTypes(isTransactionFromAITextRecognitionEnabled(), isTransactionFromAIImageRecognitionEnabled()));
const allPresetsAndSupportedImportFileCategoryAndTypes = computed<LocalizedImportFileCategoryAndTypes[]>(() => {
    if (!allPresets.value || !allPresets.value.length) {
        return allSupportedImportFileCategoryAndTypes.value;
    }

    const ret: LocalizedImportFileCategoryAndTypes[] = [];
    const presetCategoryAndTypes: LocalizedImportFileCategoryAndTypes = {
        displayCategoryName: tt('User Presets'),
        fileTypes: []
    };

    for (const preset of allPresets.value) {
        presetCategoryAndTypes.fileTypes.push({
            type: `${IMPORT_FILE_PRESET_TYPE}:${preset.id}`,
            displayName: preset.name,
            extensions: ''
        });
    }

    ret.push(presetCategoryAndTypes);
    ret.push(...allSupportedImportFileCategoryAndTypes.value);

    return ret;
});

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

const isCustomFileFormat = computed<boolean>(() => fileType.value === 'dsv' || fileType.value === 'dsv_data' || fileType.value === 'excel');
const isImportDataFromTextbox = computed<boolean>(() => allSupportedImportFileTypesMap.value[fileType.value]?.dataFromTextbox ?? false);
const needAITextRecognition = computed<boolean>(() => allSupportedImportFileTypesMap.value[fileType.value]?.needAITextRecognition ?? false);
const isAIImageImport = computed<boolean>(() => fileType.value === 'ai_image');
const needAIImageRecognition = computed<boolean>(() => allSupportedImportFileTypesMap.value[fileType.value]?.needAIImageRecognition ?? false);
const supportedAdditionalOptions = computed<ImportFileTypeSupportedAdditionalOptions | undefined>(() => allSupportedImportFileTypesMap.value[fileType.value]?.supportedAdditionalOptions);
const supportedAIAdditionalPrompt = computed<boolean>(() => !!allSupportedImportFileTypesMap.value[fileType.value]?.supportedAIAdditionalPrompt);

const isModifiedInUploadStep = computed<boolean>(() => {
    if (importFile.value || importImageFiles.value.length > 0 || importData.value) {
        return true;
    }

    if (selectedPreset.value) {
        if (importAIAdditionalPrompt.value !== selectedPreset.value.additionalPrompt) {
            return true;
        }
    } else {
        if (importAIAdditionalPrompt.value) {
            return true;
        }
    }

    return false;
});

const selectedImportType = computed<string>({
    get: () => {
        if (selectedPreset.value) {
            return `${IMPORT_FILE_PRESET_TYPE}:${selectedPreset.value.id}`;
        } else {
            return fileType.value;
        }
    },
    set: (value) => {
        const preset: ImportPresetBasicInfo | undefined = allPresets.value.find(item => `${IMPORT_FILE_PRESET_TYPE}:${item.id}` === value);

        if (preset) {
            loadPreset(preset.id, false);
        } else {
            loadingPreset.value = false;
            fileType.value = value;
            importAdditionalOptions.value = Object.assign({}, supportedAdditionalOptions.value ?? {});
            importAIAdditionalPrompt.value = '';

            selectedPreset.value = undefined;
            presetDataMapping.value = '';
            presetCustomScript.value = '';
            presetAutoExecuteScript.value = false;
            presetReplaceRules.value = '';
        }
    }
});

const allSteps = computed<StepBarItem[]>(() => {
    const steps: StepBarItem[] = [
        {
            name: 'uploadFile',
            title: tt('Upload File'),
            subTitle: tt('Upload Transaction Data File')
        }
    ];

    if (isCustomFileFormat.value) {
        if (processCustomFileFormatMethod.value === ImportCustomFileFormatProcessMethod.CustomScript) {
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

    if (isAIImageImport.value) {
        steps.push({
            name: 'recognizeImages',
            title: tt('Recognize Images'),
            subTitle: tt('Recognize Transactions from Images')
        });
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

function loadInitFileTypeFromSettings(): void {
    if (!settingsStore.appSettings.lastSelectedFileTypeInImportTransactionDialog) {
        return;
    }

    if (settingsStore.appSettings.lastSelectedFileTypeInImportTransactionDialog.startsWith(IMPORT_FILE_PRESET_TYPE)) {
        const items: string[] = settingsStore.appSettings.lastSelectedFileTypeInImportTransactionDialog.split(':');

        if (items.length === 2 && items[1]) {
            loadPreset(items[1], true);
        }

        return;
    }

    const lastSelectedFileTypes = settingsStore.appSettings.lastSelectedFileTypeInImportTransactionDialog.split('|');
    const lastSelectedFileType = lastSelectedFileTypes[0];

    if (!lastSelectedFileType || !allSupportedImportFileTypesMap.value[lastSelectedFileType]) {
        return;
    }

    fileType.value = lastSelectedFileType;

    const fileSubTypes = allSupportedImportFileTypesMap.value[lastSelectedFileType].subTypes;

    if (!fileSubTypes || fileSubTypes.length < 1) {
        return;
    }

    const lastSelectedFileSubType = lastSelectedFileTypes[1];

    if (lastSelectedFileSubType) {
        for (const subType of fileSubTypes) {
            if (subType.type === lastSelectedFileSubType) {
                fileSubType.value = lastSelectedFileSubType;
                return;
            }
        }
    }

    const firstFileSubType = fileSubTypes[0];

    if (firstFileSubType) {
        fileSubType.value = firstFileSubType.type;
    }
}

function open(): Promise<void> {
    loading.value = true;
    loadingPreset.value = false;
    fileType.value = 'ezbookkeeping';
    fileSubType.value = 'ezbookkeeping_csv';

    if (settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog && settingsStore.appSettings.lastSelectedFileTypeInImportTransactionDialog) {
        loadInitFileTypeFromSettings();
    }

    fileEncoding.value = 'auto';
    detectingFileEncoding.value = false;
    autoDetectedFileEncoding.value = undefined;
    processCustomFileFormatMethod.value = ImportCustomFileFormatProcessMethod.ColumnMapping;
    currentStep.value = 'uploadFile';

    selectedPreset.value = undefined;
    presetDataMapping.value = '';
    presetCustomScript.value = '';
    presetAutoExecuteScript.value = false;
    presetReplaceRules.value = '';
    presetRulesApplied.value = false;

    importProcess.value = 0;
    importedCount.value = null;
    totalImportCount.value = 0;
    importErrorMessage.value = '';
    importResultUnknown.value = false;
    submittedTransactions.value = [];
    remainingTransactions.value = [];
    hasImportedTransactions.value = false;

    importFile.value = null;
    importData.value = '';
    importAdditionalOptions.value = Object.assign({}, supportedAdditionalOptions.value ?? {});
    importAIAdditionalPrompt.value = '';
    parsedFileData.value = undefined;
    parsedFileColumnSeparator.value = undefined;
    importTransactionDefineColumnTab.value?.reset();
    importTransactionExecuteCustomScriptTab.value?.reset();
    importTransactions.value = undefined;
    importTransactionCheckDataTab.value?.reset();
    showState.value = true;
    clientSessionId.value = generateRandomUUID();
    clearImportImageFiles();

    const promises: Promise<unknown>[] = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false }),
        importPresetsStore.loadAllPresetBasicInfos({ force: false })
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

function loadPreset(presetId: string, silent: boolean): void {
    loadingPreset.value = true;

    importPresetsStore.getPreset({
        presetId: presetId
    }).then(preset => {
        loadingPreset.value = false;

        const originalFileType = allSupportedImportFileTypesMap.value[preset.fileType];

        if (!originalFileType) {
            logger.error(`failed to load preset ${preset.name} because the original file type ${preset.fileType} does not exist`);
            if (!silent) {
                snackbar.value?.showError('An error occurred');
            }
            return;
        }

        if (preset.fileSubType && !originalFileType.subTypes?.find(subType => subType.type === preset.fileSubType)) {
            logger.error(`failed to load preset ${preset.name} because the original file type ${preset.fileType} does not support the sub type ${preset.fileSubType}`);
            if (!silent) {
                snackbar.value?.showError('An error occurred');
            }
            return;
        }

        selectedPreset.value = preset;
        presetRulesApplied.value = false;
        fileType.value = preset.fileType;
        fileSubType.value = preset.fileSubType;

        if (preset.fileEncoding && preset.fileEncoding !== 'auto') {
            if (originalFileType.supportedEncodings?.find(encoding => encoding.encoding === preset.fileEncoding)) {
                fileEncoding.value = preset.fileEncoding;
            } else {
                logger.warn(`failed to load preset ${preset.name} because the original file type ${preset.fileType} does not support the encoding ${preset.fileEncoding}`);
                fileEncoding.value = 'auto';
            }
        } else {
            fileEncoding.value = 'auto';
        }

        processCustomFileFormatMethod.value = ImportCustomFileFormatProcessMethod.ColumnMapping;

        if (isCustomFileFormat.value) {
            if (ImportCustomFileFormatProcessMethod[preset.handlingMethod]) {
                processCustomFileFormatMethod.value = preset.handlingMethod;
            } else {
                logger.warn(`failed to load preset ${preset.name} because the handling method ${preset.handlingMethod} is not supported`);
                processCustomFileFormatMethod.value = ImportCustomFileFormatProcessMethod.ColumnMapping;
            }
        }

        importAdditionalOptions.value = Object.assign({}, supportedAdditionalOptions.value ?? {}, preset.additionalOptions);
        importAIAdditionalPrompt.value = preset.additionalPrompt;
        presetDataMapping.value = preset.dataMapping;
        presetCustomScript.value = preset.customScript;
        presetAutoExecuteScript.value = isBoolean(preset.autoExecuteScript) ? preset.autoExecuteScript : false;
        presetReplaceRules.value = preset.replaceRules;
    }).catch(error => {
        loadingPreset.value = false;

        if (!error.processed && !silent) {
            snackbar.value?.showError(error);
        }
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
        el.value = '';
        return;
    }

    if (isAIImageImport.value) {
        for (const file of el.files) {
            importImageFiles.value.push({
                file: file,
                previewUrl: URL.createObjectURL(file),
                status: 'pending'
            });
        }

        el.value = '';
    } else {
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
}

function invalidateImportedData(): void {
    accountsStore.updateAccountListInvalidState(true);
    transactionsStore.updateTransactionListInvalidState(true);
    overviewStore.updateTransactionOverviewInvalidState(true);
    statisticsStore.updateTransactionStatisticsInvalidState(true);
}

function removeImportImageFile(index: number): void {
    const targetImageFile = importImageFiles.value[index];

    if (targetImageFile && targetImageFile.previewUrl) {
        URL.revokeObjectURL(targetImageFile.previewUrl);
    }

    importImageFiles.value.splice(index, 1);
}

function clearImportImageFiles(): void {
    for (const item of importImageFiles.value) {
        if (item.previewUrl) {
            URL.revokeObjectURL(item.previewUrl);
        }
    }

    importImageFiles.value.length = 0;
}

function reloadBasisData(): void {
    loading.value = true;

    Promise.allSettled([
        accountsStore.loadAllAccounts({ force: true }),
        transactionCategoriesStore.loadAllCategories({ force: true }),
        transactionTagsStore.loadAllTags({ force: true })
    ]).then(results => {
        loading.value = false;

        const isAllUpToDate = results.length === 3
            && results[0].status === 'rejected' && results[0].reason?.isUpToDate
            && results[1].status === 'rejected' && results[1].reason?.isUpToDate
            && results[2].status === 'rejected' && results[2].reason?.isUpToDate;

        // show info if all up to date
        if (isAllUpToDate) {
            snackbar.value?.showMessage('Data is up to date');
            return;
        }

        // show error if any
        for (const result of results) {
            if (result.status === 'rejected' && !result.reason?.isUpToDate) {
                snackbar.value?.showError(result.reason);
                return;
            }
        }

        // show info if one of them updated
        for (const result of results) {
            if (result.status === 'fulfilled') {
                snackbar.value?.showMessage('Data has been updated');
                importTransactionCheckDataTab.value?.updateAllTransactionsIsValid();
                return;
            }
        }
    });
}

function recognizeImage(item: BatchImportImageItem, additionalPrompt?: string): Promise<ImportTransactionResponse[]> {
    return new Promise<ImportTransactionResponse[]>((resolve, reject) => {
        compressJpgImageByQuality(item.file, ImageUploadQualityType.HD720P).then(blob => {
            const compressedFile = KnownFileType.JPG.createFileFromBlob(blob, "image");
            importImageCancelRecognizingUuid.value = generateRandomUUID();

            transactionsStore.parseImportTransaction({
                fileType: 'ai_image',
                aiAdditionalPrompt: additionalPrompt,
                importFile: compressedFile,
                cancelableUuid: importImageCancelRecognizingUuid.value
            }).then(response => {
                resolve(response.items || []);
            }).catch(error => {
                reject(error);
            });
        }).catch(error => {
            logger.error('failed to compress image', error);
            reject('Unable to load image');
        });
    });
}

function batchRecognizeImages(): Promise<void> {
    for (const item of importImageFiles.value) {
        item.status = 'pending';
        item.failureReason = undefined;
    }

    return new Promise<void>((resolve, reject) => {
        const finish = function (): void {
            const allFailed = importImageFiles.value.every(item => item.status === 'failed');
            importImageCancelRecognizingUuid.value = '';

            if (allFailed) {
                snackbar.value?.showError('Unable to recognize image');
                reject();
            } else {
                resolve();
            }
        }

        const recurseRecognizeImage = function (index: number): void {
            if (!submitting.value) {
                finish();
                return;
            }

            if (index >= importImageFiles.value.length) {
                finish();
                return;
            }

            const item = importImageFiles.value[index];

            if (!item) {
                finish();
                return;
            }

            item.status = 'recognizing';

            recognizeImage(item, supportedAIAdditionalPrompt.value ? importAIAdditionalPrompt.value : undefined).then(results => {
                if (submitting.value) {
                    if (!importTransactions.value) {
                        importTransactions.value = [];
                    }

                    for (const importTransactionResp of results) {
                        const index = importTransactions.value.length;
                        importTransactions.value.push(ImportTransaction.of(importTransactionResp, index));
                    }

                    item.status = 'success';
                }

                recurseRecognizeImage(index + 1);
            }).catch(error => {
                if (submitting.value) {
                    item.status = 'failed';
                    item.failureReason = te(error.message || error || 'An error occurred');
                }

                recurseRecognizeImage(index + 1);
            });
        }

        recurseRecognizeImage(0);
    });
}

function cancelBatchRecognizeImages(): void {
    submitting.value = false;

    if (importImageCancelRecognizingUuid.value) {
        transactionsStore.cancelRecognizeReceiptImage(importImageCancelRecognizingUuid.value);
    }

    snackbar.value?.showMessage('User Canceled');
}

function parseData(): void {
    let uploadFile: File;
    let type: string = fileType.value;
    let encoding: string | undefined = undefined;

    if (allFileSubTypes.value) {
        type = fileSubType.value;
    }

    presetRulesApplied.value = false;

    if (isAIImageImport.value && currentStep.value === 'uploadFile') {
        if (importImageFiles.value.length < 1) {
            snackbar.value?.showError('Please select a file to import');
            return;
        }

        currentStep.value = 'recognizeImages';
        submitting.value = true;
        importTransactions.value = undefined;

        batchRecognizeImages().then(() => {
            if (!importTransactions.value || importTransactions.value.length < 1) {
                if (submitting.value) {
                    currentStep.value = 'uploadFile';
                    snackbar.value?.showMessage('No data to import');
                }
                submitting.value = false;
                return;
            }

            importTransactionCheckDataTab.value?.reset();

            if (importTransactions.value.length > 0 && importTransactions.value.length < 10) {
                importTransactionCheckDataTab.value?.setCountPerPage(-1);
            } else {
                importTransactionCheckDataTab.value?.setCountPerPage(10);
            }

            const anyFailed = importImageFiles.value.some(item => item.status === 'failed');

            if (submitting.value && !anyFailed) {
                currentStep.value = 'checkData';
            }

            submitting.value = false;
        }).catch(() => {
            submitting.value = false;
        });

        return;
    } else if (isAIImageImport.value && currentStep.value === 'recognizeImages') {
        if (!importTransactions.value || importTransactions.value.length < 1) {
            snackbar.value?.showError('No data to import');
            return;
        }

        currentStep.value = 'checkData';
        return;
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
        } else if (type === 'custom_ssv') {
            uploadFile = KnownFileType.TXT.createFile(importData.value, 'import');
        } else if (type === 'ai_txt') {
            uploadFile = KnownFileType.TXT.createFile(importData.value, 'import');
        } else {
            snackbar.value?.showError('Parameter Invalid');
            return;
        }

        encoding = UTF_8;
    } else { // should not happen, but ts would check whether uploadFile has been assigned a value
        snackbar.value?.showMessage('An error occurred');
        return;
    }

    if (isCustomFileFormat.value && currentStep.value === 'uploadFile') {
        submitting.value = true;
        let columnSeparator: string | undefined = undefined;

        if (type === 'custom_csv') {
            columnSeparator = ',';
        } else if (type === 'custom_tsv') {
            columnSeparator = '\t';
        } else if (type === 'custom_ssv') {
            columnSeparator = ';';
        }

        transactionsStore.parseImportCustomFile({
            fileType: type,
            fileEncoding: encoding,
            importFile: uploadFile
        }).then(response => {
            if (response && response.length) {
                if (processCustomFileFormatMethod.value === ImportCustomFileFormatProcessMethod.CustomScript) {
                    importTransactionExecuteCustomScriptTab.value?.reset();
                    parsedFileData.value = response;
                    parsedFileColumnSeparator.value = columnSeparator;
                    currentStep.value = 'executeCustomScript';
                    nextTick(() => {
                        if (presetCustomScript.value) {
                            importTransactionExecuteCustomScriptTab.value?.loadScript(presetCustomScript.value, presetAutoExecuteScript.value);
                        }
                    });
                } else {
                    importTransactionDefineColumnTab.value?.reset();
                    parsedFileData.value = response;
                    parsedFileColumnSeparator.value = columnSeparator;
                    currentStep.value = 'defineColumn';
                    nextTick(() => {
                        if (presetDataMapping.value) {
                            importTransactionDefineColumnTab.value?.loadColumnMapping(presetDataMapping.value);
                        }
                    });
                }
            } else {
                parsedFileData.value = undefined;
                parsedFileColumnSeparator.value = undefined;
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

        if (isCustomFileFormat.value && processCustomFileFormatMethod.value === ImportCustomFileFormatProcessMethod.ColumnMapping) {
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
        } else if (isCustomFileFormat.value && processCustomFileFormatMethod.value === ImportCustomFileFormatProcessMethod.CustomScript) {
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
            aiAdditionalPrompt: supportedAIAdditionalPrompt.value ? importAIAdditionalPrompt.value : undefined,
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

function batchApplyRules(): void {
    if (importTransactionCheckDataTab.value?.isEditing) {
        return;
    }

    let initialRules: ImportTransactionReplaceRule[] = [];

    if (presetReplaceRules.value) {
        initialRules = ImportTransactionReplaceRules.parseFromJson(presetReplaceRules.value, generateRandomUUID)?.getRules() || [];
    }

    importTransactionCheckDataTab.value?.showBatchApplyRulesDialog(initialRules);
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
        count: formatNumberToLocalizedNumerals(transactions.length)
    }).then(() => {
        submitting.value = true;
        clientSessionId.value = generateRandomUUID();
        submittedTransactions.value = transactions;
        totalImportCount.value = transactions.length;
        importedCount.value = null;
        importErrorMessage.value = '';
        importResultUnknown.value = false;
        importProcess.value = 0;
        remainingTransactions.value = [];

        let showProcessTimer : number | undefined = undefined;

        if (transactions.length > 100) {
            const sessionId = clientSessionId.value;

            showProcessTimer = window.setInterval(() => {
                transactionsStore.getImportTransactionsProcess({
                    clientSessionId: sessionId
                }).then(response => {
                    if (submitting.value && clientSessionId.value === sessionId && response && isInteger(response.importedCount)
                        && response.totalCount === transactions.length
                        && 0 <= response.importedCount && response.importedCount < response.totalCount) {
                        importProcess.value = response.importedCount / response.totalCount * 100;
                    }
                }).catch(() => {
                    importProcess.value = 0;
                    clearInterval(showProcessTimer);
                    showProcessTimer = undefined;
                });
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

            showImportResult(response);
            submitting.value = false;
        }).catch(error => {
            if (showProcessTimer) {
                importProcess.value = 0;
                clearInterval(showProcessTimer);
                showProcessTimer = undefined;
            }

            submitting.value = false;

            if (error.resultUnknown === false && error.error) {
                showImportResult({
                    finished: true,
                    importedCount: error.error.context?.importedCount ?? 0,
                    totalCount: error.error.context?.totalCount ?? submittedTransactions.value.length
                }, error.error);
            } else {
                importResultUnknown.value = true;
                currentStep.value = 'finalResult';
                invalidateImportedData();
                checkImportResult(false);
            }
        });
    });
}

function showImportResult(result: TransactionImportResponse, error?: ErrorResponse): void {
    currentStep.value = 'finalResult';

    if (!result.finished || result.totalCount !== submittedTransactions.value.length || !isInteger(result.importedCount) || result.importedCount < 0 || result.importedCount > result.totalCount || (!error && result.importedCount !== result.totalCount)) {
        importResultUnknown.value = true;
        invalidateImportedData();
        return;
    }

    importResultUnknown.value = false;
    importedCount.value = result.importedCount;
    totalImportCount.value = result.totalCount;
    remainingTransactions.value = submittedTransactions.value.slice(result.importedCount);
    importErrorMessage.value = error ? te({ error: error }) : '';

    if (result.importedCount > 0) {
        hasImportedTransactions.value = true;
        invalidateImportedData();
    }
}

function checkImportResult(force: boolean): void {
    submitting.value = true;

    transactionsStore.getImportTransactionsProcess({
        clientSessionId: clientSessionId.value,
        force: force
    }).then(response => {
        if (response && response.finished) {
            showImportResult(response);
        }
    }).catch(error => {
        logger.error('failed to confirm transaction import result', error);

        if (error.processed) {
            return;
        }

        if (isDefined(error.error?.context?.importedCount) && isDefined(error.error?.context?.totalCount)) {
            showImportResult({
                finished: true,
                importedCount: error.error.context?.importedCount ?? 0,
                totalCount: error.error.context?.totalCount ?? submittedTransactions.value.length
            }, error.error);
        } else if (force) {
            snackbar.value?.showError(error);
        }
    }).finally(() => {
        submitting.value = false;
    });
}

function showRemainingTransactions(): void {
    importTransactions.value = remainingTransactions.value;
    importTransactionCheckDataTab.value?.reset();
    currentStep.value = 'checkData';
}

function close(completed: boolean): void {
    if (importResultUnknown.value) {
        invalidateImportedData();
    }

    if (completed || hasImportedTransactions.value || importResultUnknown.value) {
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

watch(currentStep, (newValue) => {
    if (newValue === 'checkData' && !presetRulesApplied.value) {
        nextTick().then(() => {
            if (presetReplaceRules.value) {
                const initialRules = ImportTransactionReplaceRules.parseFromJson(presetReplaceRules.value, generateRandomUUID)?.getRules();

                if (initialRules && initialRules.length > 0) {
                    importTransactionCheckDataTab.value?.applyPresetRules(initialRules);
                    presetRulesApplied.value = true;
                }
            }
        });
    }
});

watch(selectedPreset, (newValue) => {
    if (newValue) {
        if (settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog && !loading.value) {
            settingsStore.setLastSelectedFileTypeInImportTransactionDialog(`${IMPORT_FILE_PRESET_TYPE}:${newValue.id}`);
        }
    } else {
        const subFileTypes = allSupportedImportFileTypesMap.value[fileType.value]?.subTypes;

        if (settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog && !loading.value) {
            if (subFileTypes && subFileTypes.length) {
                settingsStore.setLastSelectedFileTypeInImportTransactionDialog(`${fileType.value}|${fileSubType.value}`);
            } else {
                settingsStore.setLastSelectedFileTypeInImportTransactionDialog(`${fileType.value}|`);
            }
        }
    }
});

watch(fileType, (newValue) => {
    if (!selectedPreset.value) {
        const subFileTypes = allSupportedImportFileTypesMap.value[newValue]?.subTypes;

        if (subFileTypes && subFileTypes.length) {
            if (fileSubType.value !== subFileTypes[0]!.type && !subFileTypes.find(subType => subType.type === fileSubType.value)) {
                fileSubType.value = subFileTypes[0]!.type;
            } else if (settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog && !loading.value) {
                settingsStore.setLastSelectedFileTypeInImportTransactionDialog(`${newValue}|${fileSubType.value}`);
            }
        } else {
            if (settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog && !loading.value) {
                settingsStore.setLastSelectedFileTypeInImportTransactionDialog(`${newValue}|`);
            }
        }

        importAdditionalOptions.value = Object.assign({}, supportedAdditionalOptions.value ?? {});
    }

    importFile.value = null;
    parsedFileData.value = undefined;
    parsedFileColumnSeparator.value = undefined;
    importTransactions.value = undefined;
    clearImportImageFiles();
});

watch(fileSubType, (newValue) => {
    if (!selectedPreset.value) {
        if (settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog && !loading.value) {
            settingsStore.setLastSelectedFileTypeInImportTransactionDialog(`${fileType.value}|${newValue}`);
        }
    }

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

<style>
.import-transaction-dialog-layout {
    max-height: calc(100dvh - 48px);

    .import-transaction-dialog-content {
        min-height: 0;
        overflow-y: hidden !important;

        .import-transaction-dialog-window {
            min-height: 0;

            > .v-window__container {
                height: auto;
            }

            .v-window-item--active {
                min-height: 0;
            }

            .import-transaction-table-container {
                max-height: calc(100dvh - 164px);
                min-height: 0;
                overflow-y: hidden !important;

                .import-transaction-table {
                    flex: 1 1 auto;
                    min-height: 0;
                    overflow: hidden;

                    > .v-table__wrapper {
                        height: 100% !important;
                        min-height: 0;
                        overflow-y: auto;
                    }
                }

                .import-transaction-table-footer {
                    flex: 0 0 auto;
                }
            }
        }
    }
}

.import-transaction-images {
    .import-image {
        .picture-control-icon {
            display: none;
            position: absolute;
            width: 100% !important;
            height: 100% !important;
            background-color: rgba(0, 0, 0, 0.4);
        }

        .picture-control-icon > i.v-icon {
            background-color: transparent;
            color: rgba(255, 255, 255, 0.8);
        }
    }

    .import-image:hover {
        .picture-control-icon {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            vertical-align: middle;
        }
    }

    .import-image-add {
        border: 2px dashed rgba(var(--v-theme-grey-500));
    }

    .import-image-add .import-image-add-icon {
        color: rgba(var(--v-theme-grey-500));
    }

    .import-image-add.enabled:hover {
        border: 2px dashed rgba(var(--v-theme-grey-700));
    }

    .import-image-add.enabled:hover .import-image-add-icon {
        color: rgba(var(--v-theme-grey-700));
    }
}
</style>
