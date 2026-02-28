<template>
    <v-row>
        <v-col cols="12" v-if="isSupportedFileCache">
            <v-card :class="{ 'disabled': loading }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('File Cache') }}</span>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :loading="loading" @click="loadCacheStatistics()">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text class="d-flex align-end" style="height: 3rem">
                    <span class="text-body-1">{{ tt('Used storage') }}</span>
                    <v-skeleton-loader class="d-inline-block skeleton-no-margin ms-1 pt-1 pb-1" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                    <span class="text-xl ms-1" v-if="!loading">{{ fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.totalCacheSize, 2) : '-' }}</span>
                </v-card-text>

                <v-card-text>
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
                                title: 'Others',
                                count: fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.othersCacheSize, 2) : '-',
                                icon: mdiFileOutline,
                                color: 'grey'
                            }
                        ]">
                            <div class="d-flex align-center">
                                <div class="me-3">
                                    <v-avatar rounded :color="item.color" size="42" class="elevation-1">
                                        <v-icon size="24" :icon="item.icon"/>
                                    </v-avatar>
                                </div>

                                <div class="d-flex flex-column">
                                    <span class="text-caption">{{ tt(item.title) }}</span>
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-2" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                                    <span class="text-xl" v-if="!loading">{{ item.count }}</span>
                                </div>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text class="mt-2">
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || !isSupportedFileCache || !fileCacheStatistics" @click="clearMapCache()">
                        {{ tt('Clear Map Data Cache') }}
                    </v-btn>
                    <v-btn class="ms-2" color="secondary" variant="tonal"
                           :disabled="loading || !isSupportedFileCache || !fileCacheStatistics" @click="clearAllFileCache()">
                        {{ tt('Clear All File Cache') }}
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

                <v-card-text>
                    <span class="text-body-1">{{ tt('Used storage') }}</span>
                    <span class="text-xl ms-1">{{ formatVolumeToLocalizedNumerals(exchangeRatesCacheSize ?? 0, 2) }}</span>
                </v-card-text>

                <v-card-text class="mt-2">
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

    <confirm-dialog ref="confirmDialog"/>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';

import { useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useAppBrowserCacheSettingPageBase } from '@/views/base/settings/AppBrowserCacheSettingPageBase.ts';

import { isMapProviderUseExternalSDK } from '@/lib/map/index.ts';
import { getMapProvider, isMapDataFetchProxyEnabled } from '@/lib/server_settings.ts';

import {
    mdiRefresh,
    mdiFileCodeOutline,
    mdiFileImageOutline,
    mdiFileImageMarkerOutline,
    mdiFileOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;

const {
    tt,
    formatVolumeToLocalizedNumerals
} = useI18n();

const {
    isSupportedFileCache,
    loading,
    fileCacheStatistics,
    exchangeRatesCacheSize,
    allMapCacheExpirationOptions,
    allExchangeRatesDataCacheExpirationOptions,
    mapCacheExpiration,
    exchangeRatesDataCacheExpiration,
    loadCacheStatistics,
    clearMapDataCache,
    clearAllBrowserCaches,
    clearExchangeRatesDataCache
} = useAppBrowserCacheSettingPageBase();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');

function clearMapCache(): void {
    confirmDialog.value?.open('Are you sure you want to clear map data cache?').then(() => {
        clearMapDataCache().then(() => {
            loadCacheStatistics();
        });
    });
}

function clearAllFileCache(): void {
    confirmDialog.value?.open('Are you sure you want to clear all file cache?').then(() => {
        clearAllBrowserCaches().then(() => {
            location.reload();
        });
    });
}

function clearExchangeRatesCache(): void {
    confirmDialog.value?.open('Are you sure you want to clear exchange rates data cache?').then(() => {
        clearExchangeRatesDataCache();
    });
}

loadCacheStatistics();
</script>
