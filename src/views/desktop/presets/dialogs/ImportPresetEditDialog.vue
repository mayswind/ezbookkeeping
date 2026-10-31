<template>
    <v-dialog width="1200" :persistent="loading || submitting || isPresetModified" v-model="showState">
        <one-column-dialog-layout :disabled="loading || submitting" :loading="loading"
                                  :title="tt(preset.id ? 'Edit Import Preset' : 'Add Import Preset')"
                                  :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #content>
                <v-form class="mt-2">
                    <v-row>
                        <v-col cols="12" md="12">
                            <v-text-field
                                persistent-placeholder
                                type="text"
                                autocomplete="off"
                                :disabled="loading || submitting"
                                :label="tt('Preset Name')"
                                :placeholder="tt('Preset Name')"
                                v-model="preset.name" />
                        </v-col>

                        <v-col cols="12" md="12">
                            <two-column-select primary-key-field="displayCategoryName"
                                               primary-value-field="displayCategoryName"
                                               primary-title-field="displayCategoryName"
                                               primary-sub-items-field="fileTypes"
                                               secondary-key-field="type"
                                               secondary-value-field="type"
                                               secondary-title-field="displayName"
                                               :disabled="loading || submitting"
                                               :enable-filter="true"
                                               :filter-placeholder="tt('Find file type')"
                                               :filter-no-items-text="tt('No available file type')"
                                               :label="tt('File Type')"
                                               :placeholder="tt('File Type')"
                                               :items="allSupportedImportFileCategoryAndTypes"
                                               :auto-update-menu-position="true"
                                               v-model="preset.fileType" />
                        </v-col>

                        <import-transaction-options :file-type-info="selectedFileTypeInfo"
                                                    :supported-encodings="allSupportedEncodings"
                                                    :disabled="loading || submitting"
                                                    v-model:file-type="preset.fileType"
                                                    v-model:file-sub-type="preset.fileSubType"
                                                    v-model:file-encoding="preset.fileEncoding"
                                                    v-model:handling-method="preset.handlingMethod"
                                                    v-model:additional-options="preset.additionalOptions"
                                                    v-model:additional-prompt="preset.additionalPrompt" />

                        <v-col cols="12" md="12" v-if="isCustomFileFormat && preset.handlingMethod === ImportCustomFileFormatProcessMethod.ColumnMapping">
                            <div class="d-flex align-center">
                                <v-text-field
                                    readonly
                                    persistent-placeholder
                                    class="always-cursor-pointer"
                                    type="text"
                                    :disabled="loading || submitting"
                                    :clearable="!!preset.dataMapping"
                                    :label="tt('Data Mapping File')"
                                    :placeholder="tt('Data Mapping File')"
                                    :model-value="preset.dataMapping ? tt('Selected') : tt('None')"
                                    @pointerdown.prevent
                                    @click="loadDataMappingFile"
                                    @click:clear="preset.dataMapping = ''" />
                            </div>
                        </v-col>

                        <v-col cols="12" md="12" v-if="isCustomFileFormat && preset.handlingMethod === ImportCustomFileFormatProcessMethod.CustomScript">
                            <div class="d-flex align-center mb-2 mt-n2">
                                <span>{{ tt('Custom Script') }}</span>
                                <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true"
                                       :aria-label="tt('Load Script File')" :disabled="loading || submitting"
                                       @click="loadScriptFile">
                                    <v-icon :icon="mdiFolderOpenOutline" size="20" />
                                    <v-tooltip activator="parent">{{ tt('Load Script File') }}</v-tooltip>
                                </v-btn>
                                <v-spacer/>
                                <v-switch :disabled="loading || submitting"
                                          :label="tt('Automatically Execute Script')"
                                          v-model="preset.autoExecuteScript" />
                            </div>
                            <import-transaction-script-editor height="300px"
                                                              :readonly="loading || submitting"
                                                              v-model="preset.customScript" />
                        </v-col>

                        <v-col cols="12" md="12">
                            <v-text-field
                                readonly
                                persistent-placeholder
                                class="always-cursor-pointer"
                                type="text"
                                autocomplete="off"
                                :disabled="loading || submitting"
                                :clearable="!!preset.replaceRules && preset.replaceRules.length > 0"
                                :label="tt('Replace Rules')"
                                :placeholder="tt('Replace Rules')"
                                :model-value="replaceRuleCount > 0 ? tt('format.misc.ruleCount', { count: formatNumberToLocalizedNumerals(replaceRuleCount) }) : tt('None')"
                                @pointerdown.prevent
                                @click="editReplaceRules"
                                @click:clear="preset.replaceRules = ''" />
                        </v-col>
                    </v-row>
                </v-form>
            </template>

            <template #footer>
                <v-spacer/>
                <v-tooltip :disabled="!inputIsEmpty" :text="inputEmptyProblemMessage ? tt(inputEmptyProblemMessage) : ''">
                    <template v-slot:activator="{ props }">
                        <div v-bind="props" class="d-inline-block">
                            <v-btn :disabled="inputIsEmpty || loading || submitting || !isPresetModified || !selectedFileTypeInfo" @click="save">
                                {{ preset.id ? tt('Save') : tt('Add') }}
                                <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                            </v-btn>
                        </div>
                    </template>
                </v-tooltip>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <batch-apply-rules-dialog ref="replaceRulesDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import ImportTransactionOptions from '@/components/desktop/ImportTransactionOptions.vue';
import ImportTransactionScriptEditor from '@/components/desktop/ImportTransactionScriptEditor.vue';
import BatchApplyRulesDialog from '@/views/desktop/transactions/import/dialogs/BatchApplyRulesDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useImportPresetsStore } from '@/stores/importPreset.ts';

import { KnownFileType } from '@/core/file.ts';
import {
    type LocalizedImportFileCategoryAndTypes,
    type LocalizedImportFileType,
    type LocalizedImportFileTypeSupportedEncodings,
    ImportCustomFileFormatProcessMethod,
    ImportTransactionDataMapping
} from '@/core/import_transaction.ts';

import { ImportPreset } from '@/models/import_preset.ts';
import { ImportTransactionReplaceRules } from '@/models/rule.ts';

import { isEquals } from '@/lib/common.ts';
import { getBrowserTimezoneName } from '@/lib/datetime.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import { isTransactionFromAITextRecognitionEnabled, isTransactionFromAIImageRecognitionEnabled } from '@/lib/server_settings.ts';
import { openTextFileContent } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

import {
    mdiFolderOpenOutline
} from '@mdi/js';

interface ImportPresetEditResponse {
    message: string;
}

type SnackBarType = InstanceType<typeof SnackBar>;
type ReplaceRulesDialogType = InstanceType<typeof BatchApplyRulesDialog>;

const {
    tt,
    getAllSupportedImportFileCagtegoryAndTypes,
    formatNumberToLocalizedNumerals,
    getLocalizedImportTransactionCustomSampleScript
} = useI18n();

const settingsStore = useSettingsStore();
const importPresetsStore = useImportPresetsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const replaceRulesDialog = useTemplateRef<ReplaceRulesDialogType>('replaceRulesDialog');

let resolveFunc: ((value: ImportPresetEditResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const editPresetId = ref<string | null>(null);
const clientSessionId = ref<string>('');
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const preset = ref<ImportPreset>(ImportPreset.createNewPreset());
const initPreset = ref<ImportPreset | null>(null);

const currentTimezoneName = computed<string>(() => settingsStore.appSettings.timeZone || getBrowserTimezoneName());
const allSupportedImportFileCategoryAndTypes = computed<LocalizedImportFileCategoryAndTypes[]>(() => getAllSupportedImportFileCagtegoryAndTypes(isTransactionFromAITextRecognitionEnabled(), isTransactionFromAIImageRecognitionEnabled()));
const isCustomFileFormat = computed<boolean>(() => isFileTypeCustomFileFormat(preset.value.fileType));

const sampleScript = computed<string>(() => getLocalizedImportTransactionCustomSampleScript(currentTimezoneName.value));

const allSupportedImportFileTypesMap = computed<Record<string, LocalizedImportFileType>>(() => {
    const ret: Record<string, LocalizedImportFileType> = {};

    for (const importFileCategoryAndTypes of allSupportedImportFileCategoryAndTypes.value) {
        for (const importFileType of importFileCategoryAndTypes.fileTypes) {
            ret[importFileType.type] = importFileType;
        }
    }

    return ret;
});

const allSupportedEncodings = computed<LocalizedImportFileTypeSupportedEncodings[] | undefined>(() => {
    if (!selectedFileTypeInfo.value?.supportedEncodings) {
        return undefined;
    }

    const ret: LocalizedImportFileTypeSupportedEncodings[] = [];
    const autoDetectEncoding: LocalizedImportFileTypeSupportedEncodings = {
        displayName: tt('Auto detect'),
        encoding: 'auto'
    };

    ret.push(autoDetectEncoding);

    if (selectedFileTypeInfo.value.supportedEncodings.length) {
        ret.push(...selectedFileTypeInfo.value.supportedEncodings);
    }

    return ret;
});

const isPresetModified = computed<boolean>(() => {
    if (!editPresetId.value) { // Add
        return !!initPreset.value && !isEquals(preset.value.toCreateRequest(clientSessionId.value), initPreset.value.toCreateRequest(clientSessionId.value));
    } else { // Edit
        return !!initPreset.value && !isEquals(preset.value.toModifyRequest(), initPreset.value.toModifyRequest());
    }
});

const inputEmptyProblemMessage = computed<string | null>(() => {
    if (!preset.value.name) {
        return 'Preset name cannot be blank';
    } else {
        return null;
    }
});

const inputIsEmpty = computed<boolean>(() => !!inputEmptyProblemMessage.value);

const selectedFileTypeInfo = computed<LocalizedImportFileType | undefined>(() => {
    if (!preset.value.fileType) {
        return undefined;
    }

    return allSupportedImportFileTypesMap.value[preset.value.fileType];
});

const replaceRuleCount = computed<number>(() => ImportTransactionReplaceRules.parseFromJson(preset.value.replaceRules, generateRandomUUID)?.getRules().length || 0);

function isFileTypeCustomFileFormat(fileType: string): boolean {
    return fileType === 'dsv' || fileType === 'dsv_data' || fileType === 'excel';
}

function open(options: { id?: string }): Promise<ImportPresetEditResponse> {
    showState.value = true;
    loading.value = true;
    submitting.value = false;

    initPreset.value = ImportPreset.createNewPreset();
    preset.value.fillFrom(initPreset.value);

    if (options.id) {
        editPresetId.value = options.id;
        importPresetsStore.getPreset({
            presetId: options.id
        }).then(response => {
            preset.value.fillFrom(response);
            initPreset.value = ImportPreset.of(response);
            loading.value = false;
        }).catch(error => {
            loading.value = false;
            showState.value = false;

            if (!error.processed) {
                if (rejectFunc) {
                    rejectFunc(error);
                }
            }
        });
    } else {
        editPresetId.value = null;
        clientSessionId.value = generateRandomUUID();
        loading.value = false;
    }

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function save(): void {
    const problemMessage = inputEmptyProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    submitting.value = true;

    importPresetsStore.savePreset({
        preset: preset.value,
        clientSessionId: clientSessionId.value
    }).then(() => {
        submitting.value = false;

        let message = 'You have saved this import preset';

        if (!editPresetId.value) {
            message = 'You have added a new import preset';
        }

        resolveFunc?.({ message });
        showState.value = false;
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function loadDataMappingFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JSON.contentType
    }).then(content => {
        const mapping = ImportTransactionDataMapping.parseFromJson(content);

        if (!mapping) {
            snackbar.value?.showError('Data mapping file is invalid');
            return;
        }

        preset.value.dataMapping = mapping.toJson();
    }).catch(error => {
        logger.error('failed to load data mapping file', error);
        snackbar.value?.showError('Data mapping file is invalid');
    });
}

function loadScriptFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JS.contentType
    }).then(content => {
        preset.value.customScript = content;
    }).catch(error => {
        logger.error('failed to load script file', error);
        snackbar.value?.showError('Cannot load script file');
    });
}

function editReplaceRules(): void {
    replaceRulesDialog.value?.open({
        expenseCategoryNames: [],
        incomeCategoryNames: [],
        transferCategoryNames: [],
        sourceAccountNames: [],
        destinationAccountNames: [],
        tagNames: [],
        selectedTransactionCount: 0,
        initialRules: ImportTransactionReplaceRules.parseFromJson(preset.value.replaceRules, generateRandomUUID)?.getRules() || [],
        presetMode: true
    }).then(result => {
        preset.value.replaceRules = ImportTransactionReplaceRules.of(result.rules).toJson();
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

watch(() => preset.value.fileType, (newValue) => {
    if (loading.value) {
        return;
    }

    const fileTypeInfo = allSupportedImportFileTypesMap.value[newValue];

    if (fileTypeInfo?.subTypes && !preset.value.fileSubType || !fileTypeInfo?.subTypes?.find(subType => subType.type === preset.value.fileSubType)) {
        preset.value.fileSubType = fileTypeInfo?.subTypes?.[0]?.type || '';
    }

    if (!fileTypeInfo?.supportedEncodings) {
        preset.value.fileEncoding = 'auto';
    }

    if (!isFileTypeCustomFileFormat(newValue)) {
        preset.value.handlingMethod = ImportCustomFileFormatProcessMethod.ColumnMapping;
    }

    if (!fileTypeInfo?.supportedAdditionalOptions) {
        preset.value.additionalOptions = {};
    }

    if (!fileTypeInfo?.supportedAIAdditionalPrompt) {
        preset.value.additionalPrompt = '';
    }

    if (preset.value.handlingMethod !== ImportCustomFileFormatProcessMethod.ColumnMapping) {
        preset.value.dataMapping = '';
    }

    if (preset.value.handlingMethod !== ImportCustomFileFormatProcessMethod.CustomScript) {
        preset.value.customScript = '';
        preset.value.autoExecuteScript = false;
    }
});

watch(() => preset.value.handlingMethod, (newValue) => {
    if (loading.value) {
        return;
    }

    if (newValue !== ImportCustomFileFormatProcessMethod.ColumnMapping) {
        preset.value.dataMapping = '';
    }

    if (newValue !== ImportCustomFileFormatProcessMethod.CustomScript) {
        preset.value.customScript = '';
        preset.value.autoExecuteScript = false;
    }

    if (newValue === ImportCustomFileFormatProcessMethod.CustomScript && !preset.value.customScript) {
        preset.value.customScript = sampleScript.value;
    }
});

defineExpose({
    open
});
</script>
