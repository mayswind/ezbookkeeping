<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loadingCacheStatistics }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('File Cache') }}</span>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :loading="loadingCacheStatistics" @click="reloadCacheStatistics()">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text class="mt-1" v-if="loadingCacheStatistics">
                    <span class="text-body-1">{{ tt('Used storage') }}</span>
                    <v-skeleton-loader class="d-inline-block skeleton-no-margin ml-1 pt-1 pb-1" type="text" style="width: 100px; height: 24px" :loading="true" v-if="loadingCacheStatistics"></v-skeleton-loader>
                </v-card-text>

                <v-card-text v-else-if="!loadingCacheStatistics">
                    <span class="text-body-1">{{ tt('Used storage') }}</span>
                    <span class="text-xl ml-1" v-if="!loadingCacheStatistics">{{ fileCacheStatistics ? formatVolumeToLocalizedNumerals(fileCacheStatistics.totalCacheSize, 2) : '-' }}</span>
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
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-2" type="text" style="width: 100px" :loading="true" v-if="loadingCacheStatistics"></v-skeleton-loader>
                                    <span class="text-xl" v-if="!loadingCacheStatistics">{{ item.count }}</span>
                                </div>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text class="mt-2">
                    <v-btn color="gray" variant="tonal"
                           :disabled="loadingCacheStatistics || !fileCacheStatistics"  @click="clearFileCache()">
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

                <v-card-text>
                    <span class="text-body-1">{{ tt('Used storage') }}</span>
                    <span class="text-xl ml-1">{{ formatVolumeToLocalizedNumerals(exchangeRatesCacheSize ?? 0, 2) }}</span>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { type BrowserCacheStatistics } from '@/core/cache.ts';

import { loadBrowserCacheStatistics, clearAllBrowserCaches } from '@/lib/cache.ts';

import {
    mdiRefresh,
    mdiFileCodeOutline,
    mdiFileImageOutline,
    mdiFileImageMarkerOutline,
    mdiFileOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt, formatVolumeToLocalizedNumerals } = useI18n();

const exchangeRatesStore = useExchangeRatesStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loadingCacheStatistics = ref<boolean>(true);
const fileCacheStatistics = ref<BrowserCacheStatistics | undefined>(undefined);
const exchangeRatesCacheSize = ref<number | undefined>(undefined);

function reloadCacheStatistics(): void {
    loadingCacheStatistics.value = true;

    loadBrowserCacheStatistics().then(statistics => {
        fileCacheStatistics.value = statistics;
        exchangeRatesCacheSize.value = exchangeRatesStore.getExchangeRatesCacheSize();
        loadingCacheStatistics.value = false;
    }).catch(() => {
        loadingCacheStatistics.value = false;
        snackbar.value?.showError('Failed to load browser cache statistics');
    });
}

function clearFileCache(): void {
    confirmDialog.value?.open('Are you sure you want to clear file cache?').then(() => {
        clearAllBrowserCaches().then(() => {
            location.reload();
        });
    });
}

reloadCacheStatistics();
</script>
