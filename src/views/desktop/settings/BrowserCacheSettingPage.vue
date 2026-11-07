<template>
    <v-row>
        <v-col cols="12" v-if="isSupportedFileCache">
            <v-card :class="{ 'disabled': loading }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('File Cache') }}</span>
                        <v-btn density="compact" color="default" variant="text" class="ms-2"
                               :aria-label="tt('Refresh')" :icon="true" :loading="loading" @click="loadCacheStatistics(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text class="d-flex align-end" style="height: 3rem">
                    <span class="text-body-large">{{ tt('Used storage') }}</span>
                    <v-skeleton-loader class="d-inline-block skeleton-no-margin ms-2 pt-1 pb-1" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                    <span class="text-title-medium ms-2" v-if="!loading">{{ fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.totalCacheSize, 2) : '-' }}</span>
                </v-card-text>

                <v-card-text class="py-0">
                    <v-row>
                        <v-col cols="6" sm="3" :key="idx" v-for="(item, idx) in [
                            {
                                title: 'Application Code',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.codeCacheSize, 2) : '-',
                                icon: mdiFileCodeOutline,
                                color: 'info-darken-1'
                            },
                            {
                                title: 'Resource Files',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.assetsCacheSize, 2) : '-',
                                icon: mdiFileImageOutline,
                                color: 'teal'
                            },
                            {
                                title: 'Map Data',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.mapCacheSize, 2) : '-',
                                icon: mdiFileImageMarkerOutline,
                                color: 'warning'
                            },
                            {
                                title: 'User Avatars',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.userAvatarCacheSize, 2) : '-',
                                icon: mdiAccountCircleOutline,
                                color: 'primary'
                            },
                            {
                                title: 'Transaction Pictures',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.transactionPictureCacheSize, 2) : '-',
                                icon: mdiImage,
                                color: 'error-darken-1'
                            },
                            {
                                title: 'Custom Icons',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.customIconCacheSize, 2) : '-',
                                icon: mdiShapePlusOutline,
                                color: 'purple'
                            }
                        ]">
                            <div class="d-flex align-center">
                                <div class="me-3">
                                    <v-avatar rounded :color="item.color" size="42" class="elevation-1">
                                        <v-icon size="24" :icon="item.icon"/>
                                    </v-avatar>
                                </div>

                                <div class="d-flex flex-column">
                                    <span class="text-body-small">{{ tt(item.title) }}</span>
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-2" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                                    <span class="text-title-medium" v-if="!loading">{{ item.count }}</span>
                                </div>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text>
                    <v-btn color="secondary" variant="tonal"
                           :disabled="!isSupportedFileCache || loading || clearingCache || !fileCacheStatistics" @click="openClearFileCache">
                        {{ tt('Clear File Cache') }}
                    </v-btn>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('Exchange Rates Data') }}</span>
                    </div>
                </template>

                <v-card-text class="d-flex align-end py-0">
                    <span class="text-body-large">{{ tt('Used storage') }}</span>
                    <v-skeleton-loader class="d-inline-block skeleton-no-margin ms-2 mt-1 pt-1 pb-1" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                    <span class="text-title-medium ms-2" v-if="!loading">{{ formatVolumeToLocalizedNumerals(exchangeRatesCacheSize ?? 0, 2) }}</span>
                </v-card-text>

                <v-card-text>
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || !exchangeRatesCacheSize" @click="clearExchangeRatesCache()">
                        {{ tt('Clear Exchange Rates Data Cache') }}
                    </v-btn>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Cache Expiration Time')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" sm="6" v-if="getMapProvider()">
                                <v-select
                                    item-title="name"
                                    item-value="value"
                                    persistent-placeholder
                                    :disabled="loading || !isSupportedFileCache || !fileCacheStatistics || isMapProviderUseExternalSDK() || !isMapDataFetchProxyEnabled()"
                                    :label="tt('Cache Expiration for Map Data')"
                                    :placeholder="tt('Cache Expiration for Map Data')"
                                    :items="allMapCacheExpirationOptions"
                                    v-model="mapCacheExpiration"
                                />
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-select
                                    item-title="name"
                                    item-value="value"
                                    persistent-placeholder
                                    :disabled="loading || !isSupportedFileCache || !fileCacheStatistics"
                                    :label="tt('Cache Expiration for Transaction Pictures')"
                                    :placeholder="tt('Cache Expiration for Transaction Pictures')"
                                    :items="allTransactionPictureCacheExpirationOptions"
                                    v-model="transactionPictureCacheExpiration"
                                />
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-select
                                    item-title="name"
                                    item-value="value"
                                    persistent-placeholder
                                    :disabled="loading"
                                    :label="tt('Cache Expiration for Exchange Rates Data')"
                                    :placeholder="tt('Cache Expiration for Exchange Rates Data')"
                                    :items="allExchangeRatesDataCacheExpirationOptions"
                                    v-model="exchangeRatesDataCacheExpiration"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <v-dialog width="600" :persistent="clearingCache" v-model="showClearFileCache">
        <one-column-dialog-layout content-class="pa-0"
                                  :title="tt('Clear File Cache')" :cancel-button-title="tt('Cancel')"
                                  :disabled="clearingCache" @cancel="showClearFileCache = false">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="loading || clearingCache || !selectedFileCacheTypes.length" @click="clearFileCache">
                    {{ tt('Clear') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="clearingCache"></v-progress-circular>
                </v-btn>
            </template>

            <template #content>
                <v-list class="mx-3 d-flex flex-column gap-1" density="comfortable" select-strategy="classic"
                        :disabled="loading || clearingCache" v-model:selected="selectedFileCacheTypes">
                    <v-list-item class="py-0" density="compact"
                                 :key="item.key" :value="item" :title="tt(item.name)"
                                 v-for="item in FileCacheType.values()">
                        <template #prepend="{ isActive }">
                            <v-list-item-action start>
                                <v-checkbox-btn :model-value="isActive"
                                                @update:model-value="setFileCacheSelected(item, !!$event)"></v-checkbox-btn>
                            </v-list-item-action>
                        </template>

                        <template #append>
                            <span class="text-body-medium">{{ getDisplayFileCacheSize(item) }}</span>
                        </template>
                    </v-list-item>
                </v-list>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useAppBrowserCacheSettingPageBase } from '@/views/base/settings/AppBrowserCacheSettingPageBase.ts';

import { FileCacheType } from '@/core/cache.ts';

import { isMapProviderUseExternalSDK } from '@/lib/map/index.ts';
import { getMapProvider, isMapDataFetchProxyEnabled } from '@/lib/server_settings.ts';
import logger from '@/lib/logger.ts';

import {
    mdiRefresh,
    mdiFileCodeOutline,
    mdiFileImageOutline,
    mdiFileImageMarkerOutline,
    mdiAccountCircleOutline,
    mdiImage,
    mdiShapePlusOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    formatVolumeToLocalizedNumerals
} = useI18n();

const {
    isSupportedFileCache,
    loading,
    clearingCache,
    fileCacheStatistics,
    exchangeRatesCacheSize,
    selectedFileCacheTypes,
    showClearFileCache,
    allMapCacheExpirationOptions,
    allTransactionPictureCacheExpirationOptions,
    allExchangeRatesDataCacheExpirationOptions,
    mapCacheExpiration,
    transactionPictureCacheExpiration,
    exchangeRatesDataCacheExpiration,
    loadCacheStatistics,
    getDisplayFileCacheSize,
    openClearFileCache,
    setFileCacheSelected,
    clearSelectedFileCaches,
    clearExchangeRatesDataCache
} = useAppBrowserCacheSettingPageBase();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

function clearFileCache(): void {
    clearSelectedFileCaches().then(() => {
        snackbar.value?.showMessage('Cache has been cleared');
    }).catch(error => {
        logger.error('failed to clear browser cache', error);
        snackbar.value?.showMessage('Unable to clear cache');
    });
}

function clearExchangeRatesCache(): void {
    confirmDialog.value?.open('Are you sure you want to clear exchange rates data cache?').then(() => {
        clearExchangeRatesDataCache();
    });
}

loadCacheStatistics(true);
</script>
