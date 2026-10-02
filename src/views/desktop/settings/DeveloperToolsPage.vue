<template>
    <v-row>
        <v-col cols="12">
            <v-card v-if="isLanguagePreviewEnabled()">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('Preview Language Configuration') }}</span>
                    </div>
                </template>

                <v-card-text class="mt-2 text-body-medium">
                    <div>
                        <v-alert variant="tonal">
                            {{ tt('This tool temporarily replaces the language configuration for the current language. Changes only affect the current session and are discarded when the page is refreshed.') }}
                        </v-alert>
                    </div>

                    <div class="d-flex mt-4">
                        <span>{{ tt('Current Language') }}</span>
                        <v-spacer/>
                        <v-text-field readonly disabled density="compact" max-width="400px"
                                      :model-value="`${currentLanguageDisplayName} (${currentLanguage})`"/>
                    </div>

                    <div class="d-flex mt-2">
                        <span>{{ tt('Text Direction') }}</span>
                        <v-spacer/>
                        <v-select
                            density="compact"
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            max-width="400px"
                            :placeholder="tt('Text Direction')"
                            :items="textDirectionOptions"
                            v-model="currentLanguageTextDirection" />
                    </div>


                    <div class="mt-2">
                        <div class="d-flex align-center">
                            <span>{{ tt('Language Configuration') }}</span>
                            <v-btn density="compact" color="default" variant="text" class="ms-2"
                                   :aria-label="tt('Load Language Configuration File')" :icon="true" @click="loadLanguageConfigurationFile">
                                <v-icon :icon="mdiFolderOpenOutline" size="20" />
                                <v-tooltip activator="parent">{{ tt('Load Language Configuration File') }}</v-tooltip>
                            </v-btn>
                            <v-btn density="compact" color="default" variant="text" class="ms-1"
                                   :aria-label="tt('Save to File')" :icon="true"
                                   :disabled="!currentLanguageConfigJson || !currentLanguageConfigJson.trim()" @click="saveLanguageConfigurationFile">
                                <v-icon :icon="mdiContentSaveOutline" size="20" />
                                <v-tooltip activator="parent">{{ tt('Save to File') }}</v-tooltip>
                            </v-btn>
                        </div>

                        <code-editor class="w-100 mt-2" style="height: 400px" language="json"
                                     :rounded="true" :line-numbers="true"
                                     v-model="currentLanguageConfigJson" />
                    </div>

                    <div class="mt-4">
                        <span>{{ tt('Paste a language JSON object or load a local file, then click Apply Preview. After editing the file, load it again to test the updated content. The JSON replaces the entire current language configuration.') }}</span>
                    </div>
                </v-card-text>

                <v-divider />

                <v-card-text class="d-flex flex-wrap gap-4">
                    <v-btn :disabled="!currentLanguageConfigJson || !currentLanguageConfigJson.trim()" @click="applyPreview">
                        {{ tt('Apply Preview') }}
                    </v-btn>

                    <v-btn color="default" variant="tonal" :disabled="!hasPreview" @click="restorePreview">
                        {{ tt('Restore Built-in Language') }}
                    </v-btn>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { NameValue } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import { KnownFileType } from '@/core/file.ts';

import { isLanguagePreviewEnabled } from '@/lib/server_settings.ts';
import { openTextFileContent, startDownloadFile } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

import {
    mdiFolderOpenOutline,
    mdiContentSaveOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    getCurrentLanguageTag,
    getCurrentLanguageDisplayName,
    getCurrentLanguageTextDirection,
    getCurrentLanguageMessagesJson,
    hasPreviewedConfiguration,
    previewLanguage,
    restoreLanguagePreview
} = useI18n();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const currentLanguageConfigJson = ref<string>('');
const currentLanguageTextDirection = ref<TextDirection>(TextDirection.LTR);

const currentLanguage = computed(() => getCurrentLanguageTag());
const currentLanguageDisplayName = computed(() => getCurrentLanguageDisplayName());
const hasPreview = computed(() => hasPreviewedConfiguration(currentLanguage.value));

const textDirectionOptions = computed<NameValue[]>(() => [
    { name: tt('Left to right (LTR)'), value: TextDirection.LTR },
    { name: tt('Right to left (RTL)'), value: TextDirection.RTL }
]);

function loadLanguageConfigurationFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JSON.contentType
    }).then(content => {
        currentLanguageConfigJson.value = content;
    }).catch(error => {
        logger.error('Failed to load language configuration file', error);
        snackbar.value?.showError('Unable to load language configuration file');
    });
}

function saveLanguageConfigurationFile(): void {
    startDownloadFile('lang.json', KnownFileType.JSON.createBlob(currentLanguageConfigJson.value));
}

function applyPreview(): void {
    try {
        previewLanguage(currentLanguageConfigJson.value, currentLanguageTextDirection.value);
        snackbar.value?.showMessage('Language preview applied. You can navigate to other pages to check the result.');
    } catch (ex) {
        logger.error('failed to apply language preview', ex);
        snackbar.value?.showError( 'Unable to apply language preview');
    }
}

function restorePreview(): void {
    try {
        restoreLanguagePreview();
        currentLanguageConfigJson.value = getCurrentLanguageMessagesJson();
        currentLanguageTextDirection.value = getCurrentLanguageTextDirection();
        snackbar.value?.showMessage('Built-in language configuration restored');
    } catch (ex) {
        logger.error('failed to restore language configuration', ex);
        snackbar.value?.showError( 'Unable to restore the language configuration');
    }
}

watch(currentLanguage, () => {
    currentLanguageConfigJson.value = getCurrentLanguageMessagesJson();
    currentLanguageTextDirection.value = getCurrentLanguageTextDirection();
}, { immediate: true });
</script>
