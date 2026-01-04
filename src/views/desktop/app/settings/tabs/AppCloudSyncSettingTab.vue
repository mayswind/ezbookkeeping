<template>
    <v-row>
        <v-col cols="12">
            <v-card>
                <template #title>
                    <span>{{ tt('Settings Sync') }}</span>
                    <v-progress-circular indeterminate size="20" class="ms-3" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text class="pb-0">
                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-5" type="text" style="width: 150px" :loading="true" v-if="loading"></v-skeleton-loader>
                    <p class="text-body-1" v-if="!loading && !isEnableCloudSync">
                        {{ tt('Settings sync is not enabled') }}
                    </p>
                    <p class="text-body-1" v-if="!loading && isEnableCloudSync">
                        {{ tt('Settings sync has been enabled') }}
                    </p>
                </v-card-text>

                <v-card-text>
                    <v-expansion-panels class="synchronized-settings" multiple
                                        :readonly="true" :hide-actions="true"
                                        :disabled="loading || enabling || disabling"
                                        v-model="openedPanel">
                        <v-expansion-panel class="border" value="synchronizedSettings">
                            <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                                <div class="d-flex align-center justify-center w-100">
                                    <div class="w-100">
                                        <span>{{ tt('Synchronized Settings') }}</span>
                                    </div>
                                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                                           :disabled="loading || enabling || disabling" :icon="true">
                                        <v-icon :icon="mdiDotsVertical" />
                                        <v-menu activator="parent">
                                            <v-list>
                                                <v-list-item :disable="loading || enabling || disabling"
                                                             :prepend-icon="mdiSelectAll"
                                                             :title="tt('Select All')"
                                                             @click="selectAllSettings"></v-list-item>
                                                <v-list-item :disable="loading || enabling || disabling"
                                                             :prepend-icon="mdiSelect"
                                                             :title="tt('Select None')"
                                                             @click="selectNoneSettings"></v-list-item>
                                                <v-list-item :disable="loading || enabling || disabling"
                                                             :prepend-icon="mdiSelectInverse"
                                                             :title="tt('Invert Selection')"
                                                             @click="selectInvertSettings"></v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </v-btn>
                                </div>
                            </v-expansion-panel-title>
                            <v-expansion-panel-text>
                                <v-list rounded density="comfortable" class="pa-0">
                                    <template :key="categorizedItems.categoryName"
                                              v-for="(categorizedItems, categoryIdx) in ALL_APPLICATION_CLOUD_SETTINGS">
                                        <v-divider v-if="categoryIdx > 0"/>

                                        <v-list-item>
                                            <template #prepend>
                                                <v-checkbox :disabled="loading || enabling || disabling"
                                                            :model-value="isAllSettingsSelected(categorizedItems)"
                                                            :indeterminate="hasSettingSelectedButNotAllChecked(categorizedItems)"
                                                            @update:model-value="updateSettingsSelected(categorizedItems, !!$event)">
                                                    <template #label>
                                                        <span>{{ tt(categorizedItems.categoryName) }}</span>
                                                        <span class="mx-2" v-if="categorizedItems.categorySubName">/</span>
                                                        <span v-if="categorizedItems.categorySubName">{{ tt(categorizedItems.categorySubName) }}</span>
                                                    </template>
                                                </v-checkbox>
                                            </template>
                                        </v-list-item>

                                        <v-divider/>

                                        <v-list rounded density="comfortable" class="pa-0 ms-4">
                                            <template :key="settingItem.settingKey"
                                                      v-for="(settingItem, itemIdx) in categorizedItems.items">
                                                <v-divider v-if="itemIdx > 0"/>

                                                <v-list-item>
                                                    <template #prepend>
                                                        <v-checkbox :disabled="loading || enabling || disabling"
                                                                    :model-value="enabledApplicationCloudSettings[settingItem.settingKey]"
                                                                    @update:model-value="updateSettingSelected(settingItem, !!$event)">
                                                            <template #label>
                                                                <span>{{ tt(settingItem.settingName) }}</span>
                                                                <v-icon class="ms-2 me-0" start size="16" :icon="mdiCellphone" v-if="settingItem.mobile"/>
                                                                <v-icon class="ms-2 me-0" start size="16" :icon="mdiMonitor" v-if="settingItem.desktop"/>
                                                            </template>
                                                        </v-checkbox>
                                                    </template>
                                                </v-list-item>
                                            </template>
                                        </v-list>
                                    </template>
                                </v-list>
                            </v-expansion-panel-text>
                        </v-expansion-panel>
                    </v-expansion-panels>
                </v-card-text>

                <v-card-text>
                    <v-row>
                        <v-col cols="12" class="d-flex flex-wrap gap-4">
                            <v-btn :disabled="loading || enabling || disabling || !hasEnabledApplicationCloudSettings" v-if="!isEnableCloudSync" @click="enable(false)">
                                {{ tt('Enable Settings Sync') }}
                                <v-progress-circular indeterminate size="22" class="ms-2" v-if="enabling"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="loading || enabling || disabling || !hasEnabledApplicationCloudSettings" v-if="isEnableCloudSync" @click="enable(true)">
                                {{ tt('Update Synchronized Settings') }}
                                <v-progress-circular indeterminate size="22" class="ms-2" v-if="enabling"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="loading || enabling || disabling" v-if="isEnableCloudSync" @click="disable">
                                {{ tt('Disable Settings Sync') }}
                                <v-progress-circular indeterminate size="22" class="ms-2" v-if="disabling"></v-progress-circular>
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAppCloudSyncBase } from '@/views/base/settings/AppCloudSyncPageBase.ts';

import { useUserStore } from '@/stores/user.ts';

import {
    mdiDotsVertical,
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiCellphone,
    mdiMonitor,
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const {
    ALL_APPLICATION_CLOUD_SETTINGS,
    loading,
    enabling,
    disabling,
    enabledApplicationCloudSettings,
    isEnableCloudSync,
    hasEnabledApplicationCloudSettings,
    enabledApplicationCloudSettingKeys,
    isAllSettingsSelected,
    hasSettingSelectedButNotAllChecked,
    updateSettingsSelected,
    updateSettingSelected,
    selectAllSettings,
    selectNoneSettings,
    selectInvertSettings,
    setUserApplicationCloudSettings
} = useAppCloudSyncBase();

const userStore = useUserStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const openedPanel = ref<string[]>(['synchronizedSettings']);

function init(): void {
    loading.value = true;

    userStore.getUserApplicationCloudSettings().then(response => {
        setUserApplicationCloudSettings(response);
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function enable(update: boolean): void {
    enabling.value = true;

    userStore.fullUpdateUserApplicationCloudSettings(enabledApplicationCloudSettingKeys.value).then(() => {
        enabling.value = false;

        if (!update) {
            snackbar.value?.showMessage('Settings sync has been enabled');
        } else {
            snackbar.value?.showMessage('Synchronized settings have been updated');
        }
    }).catch(error => {
        enabling.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function disable(): void {
    disabling.value = true;

    userStore.disableUserApplicationCloudSettings().then(() => {
        enabledApplicationCloudSettings.value = {};
        disabling.value = false;
        snackbar.value?.showMessage('Settings sync has been disabled');
    }).catch(error => {
        disabling.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

init();
</script>

<style>
.synchronized-settings .v-expansion-panel-title {
    cursor: inherit;
}

.synchronized-settings .v-expansion-panel-text__wrapper {
    padding: 0 0 0 0;
}
</style>
