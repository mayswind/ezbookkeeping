import { ref } from 'vue';
import { defineStore } from 'pinia';

import { type BeforeResolveFunction, itemAndIndex } from '@/core/base.ts';

import {
    type ImportPresetInfoResponse,
    type ImportPresetNewDisplayOrderRequest,
    ImportPresetBasicInfo,
    ImportPreset
} from '@/models/import_preset.ts';

import { isEquals } from '@/lib/common.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useImportPresetsStore = defineStore('importPresets', () => {
    const allImportPresetBasicInfos = ref<ImportPresetBasicInfo[]>([]);
    const importPresetListStateInvalid = ref<boolean>(true);

    function loadImportPresetList(presets: ImportPresetBasicInfo[]): void {
        allImportPresetBasicInfos.value = presets;
    }

    function addPresetToImportPresetList(preset: ImportPresetBasicInfo): void {
        allImportPresetBasicInfos.value.push(preset);
    }

    function updatePresetInImportPresetList(currentPreset: ImportPresetBasicInfo): void {
        for (const [preset, index] of itemAndIndex(allImportPresetBasicInfos.value)) {
            if (preset.id === currentPreset.id) {
                allImportPresetBasicInfos.value.splice(index, 1, currentPreset);
                break;
            }
        }
    }

    function updatePresetDisplayOrderInImportPresetList({ from, to }: { from: number, to: number }): void {
        allImportPresetBasicInfos.value.splice(to, 0, allImportPresetBasicInfos.value.splice(from, 1)[0] as ImportPresetBasicInfo);
    }

    function removePresetFromImportPresetList(currentPreset: ImportPresetBasicInfo): void {
        for (const [preset, index] of itemAndIndex(allImportPresetBasicInfos.value)) {
            if (preset.id === currentPreset.id) {
                allImportPresetBasicInfos.value.splice(index, 1);
                break;
            }
        }
    }

    function updateImportPresetListInvalidState(invalidState: boolean): void {
        importPresetListStateInvalid.value = invalidState;
    }

    function resetImportPresets(): void {
        allImportPresetBasicInfos.value = [];
        importPresetListStateInvalid.value = true;
    }

    function loadAllPresetBasicInfos({ force }: { force?: boolean }): Promise<ImportPresetBasicInfo[]> {
        if (!force && !importPresetListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allImportPresetBasicInfos.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getImportPresets().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve import preset list' });
                    return;
                }

                if (allImportPresetBasicInfos.value) {
                    updateImportPresetListInvalidState(false);
                }

                const importPresetBasicInfos = ImportPresetBasicInfo.ofMulti(data.result);

                if (force && data.result && isEquals(allImportPresetBasicInfos.value, importPresetBasicInfos)) {
                    reject({ message: 'Import preset list is up to date', isUpToDate: true });
                    return;
                }

                loadImportPresetList(importPresetBasicInfos);

                resolve(importPresetBasicInfos);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load import preset list', error);
                } else {
                    logger.error('failed to load import preset list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve import preset list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getPreset({ presetId }: { presetId: string }): Promise<ImportPreset> {
        return new Promise((resolve, reject) => {
            services.getImportPreset({
                id: presetId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve import preset' });
                    return;
                }

                const preset: ImportPreset = ImportPreset.of(data.result);

                resolve(preset);
            }).catch(error => {
                logger.error('failed to load import preset', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve import preset' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function savePreset({ preset, clientSessionId }: { preset: ImportPreset, clientSessionId: string }): Promise<ImportPreset> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<ImportPresetInfoResponse>;

            if (!preset.id) {
                promise = services.addImportPreset(preset.toCreateRequest(clientSessionId));
            } else {
                promise = services.modifyImportPreset(preset.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!preset.id) {
                        reject({ message: 'Unable to add import preset' });
                    } else {
                        reject({ message: 'Unable to save import preset' });
                    }
                    return;
                }

                const presetBasicInfo: ImportPresetBasicInfo = ImportPresetBasicInfo.of(data.result);

                if (!preset.id) {
                    addPresetToImportPresetList(presetBasicInfo);
                } else {
                    updatePresetInImportPresetList(presetBasicInfo);
                }

                resolve(ImportPreset.of(data.result));
            }).catch(error => {
                logger.error('failed to save import preset', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!preset.id) {
                        reject({ message: 'Unable to add import preset' });
                    } else {
                        reject({ message: 'Unable to save import preset' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function changePresetDisplayOrder({ presetId, from, to }: { presetId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentPreset: ImportPresetBasicInfo | null = null;

            for (const preset of allImportPresetBasicInfos.value) {
                if (preset.id === presetId) {
                    currentPreset = preset;
                    break;
                }
            }

            if (!currentPreset || !allImportPresetBasicInfos.value[to]) {
                reject({ message: 'Unable to move import preset' });
                return;
            }

            if (!importPresetListStateInvalid.value) {
                updateImportPresetListInvalidState(true);
            }

            updatePresetDisplayOrderInImportPresetList({ from, to });

            resolve();
        });
    }

    function updatePresetDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: ImportPresetNewDisplayOrderRequest[] = [];

        for (const [preset, index] of itemAndIndex(allImportPresetBasicInfos.value)) {
            newDisplayOrders.push({
                id: preset.id,
                displayOrder: index + 1
            });
        }

        return new Promise((resolve, reject) => {
            services.moveImportPreset({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move import preset' });
                    return;
                }

                loadAllPresetBasicInfos({ force: false }).finally(() => {
                    if (importPresetListStateInvalid.value) {
                        updateImportPresetListInvalidState(false);
                    }

                    resolve(data.result);
                });
            }).catch(error => {
                logger.error('failed to save presets display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move import preset' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deletePreset({ preset, beforeResolve }: { preset: ImportPresetBasicInfo, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteImportPreset({
                id: preset.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete import preset' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removePresetFromImportPresetList(preset);
                    });
                } else {
                    removePresetFromImportPresetList(preset);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete import preset', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete import preset' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        allImportPresetBasicInfos,
        // functions
        updateImportPresetListInvalidState,
        resetImportPresets,
        loadAllPresetBasicInfos,
        getPreset,
        savePreset,
        changePresetDisplayOrder,
        updatePresetDisplayOrders,
        deletePreset
    };
});
