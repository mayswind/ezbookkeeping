<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('global.app.title') }}</span>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" @click="refreshBrowserCache"
                               v-if="!clientVersionMatchServerVersion">
                            <v-icon :icon="mdiWebRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh Browser Cache') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Version') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <span class="text-body-1">{{ clientVersion }}</span>
                        </v-col>
                    </v-row>
                    <v-row no-gutters v-if="clientBuildTime">
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Build Time') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <span class="text-body-1">{{ clientBuildTime }}</span>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Official Website') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <a class="text-body-1" href="https://github.com/mayswind/ezbookkeeping" target="_blank">
                                https://github.com/mayswind/ezbookkeeping
                            </a>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Report Issue') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <a class="text-body-1" href="https://github.com/mayswind/ezbookkeeping/issues" target="_blank">
                                https://github.com/mayswind/ezbookkeeping/issues
                            </a>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Getting help') }}</span>
                        </v-col>
                        <v-col cols="12" md="10">
                            <a class="text-body-1" href="https://ezbookkeeping.mayswind.net" target="_blank">
                                https://ezbookkeeping.mayswind.net
                            </a>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" v-if="exchangeRatesData && !isUserCustomExchangeRates">
            <v-card :title="tt('Exchange Rates Data')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Provider') }}</span>
                        </v-col>
                        <v-col cols="12" md="10">
                            <a class="text-body-1" :href="exchangeRatesData.referenceUrl" target="_blank"
                               v-if="exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</a>
                            <span class="text-body-1" v-if="!exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" v-if="mapProviderName">
            <v-card :title="tt('Map')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Provider') }}</span>
                        </v-col>
                        <v-col cols="12" md="10">
                            <a class="text-body-1" :href="mapProviderWebsite" target="_blank"
                               v-if="mapProviderWebsite">{{ mapProviderName }}</a>
                            <span class="text-body-1" v-if="!mapProviderWebsite">{{ mapProviderName }}</span>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('License')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12">
                            <p>
                                <span :key="num" v-for="(line, num) in licenseLines"
                                      :style="{ 'display': line ? 'initial' : 'block', 'padding' : line ? '0' : '0 0 1em 0' }">
                                    {{ line }}
                                </span>
                            </p>
                            <v-divider/>
                            <br/>
                            <p>
                                <span>ezBookkeeping's codebase and localization translation rely on contributions from the community. The following people have contributed to ezBookkeeping:</span>
                            </p>
                            <div>
                                <strong>Project Maintainer</strong>
                                <div class="mt-2">
                                    <a target="_blank" href="https://github.com/mayswind">@mayswind</a>
                                </div>
                            </div>
                            <p class="mt-4">
                                <strong>Code Contributors</strong>
                            </p>
                            <table class="contributors-table">
                                <thead>
                                <tr>
                                    <th>Contributor</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr :key="index" v-for="(contributor, index) in contributors.code">
                                    <td>
                                        <a target="_blank" :href="`https://github.com/${contributor}`">
                                            @{{ contributor }}
                                        </a>
                                    </td>
                                </tr>
                                </tbody>
                            </table>
                            <p class="mt-4">
                                <strong>Translation Contributors</strong>
                            </p>
                            <table class="contributors-table">
                                <thead>
                                <tr>
                                    <th>Tag</th>
                                    <th>Language</th>
                                    <th>Contributors</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr :key="lang.languageTag" v-for="lang in allLanguages">
                                    <td>{{ lang.languageTag }}</td>
                                    <td>{{ lang.nativeDisplayName }}</td>
                                    <td>
                                        <template :key="contributor"
                                                  v-for="(contributor, index) in contributors.translators[lang.languageTag] ?? []">
                                            <span v-if="index > 0">, </span>
                                            <a target="_blank" :href="`https://github.com/${contributor}`">
                                                @{{ contributor }}
                                            </a>
                                        </template>
                                        <span v-if="!contributors.translators[lang.languageTag] || !contributors.translators[lang.languageTag]?.length">/</span>
                                    </td>
                                </tr>
                                </tbody>
                            </table>
                            <p class="mt-4 mb-4">
                                <span>ezBookkeeping also contains additional third party software and illustration.</span><br/>
                                <span>All the third party software / illustration included or linked is redistributed under the terms and conditions of their original licenses.</span>
                            </p>
                            <p :key="licenseInfo.name" v-for="licenseInfo in thirdPartyLicenses">
                                <strong>{{ licenseInfo.name }}</strong>
                                <br v-if="licenseInfo.copyright"/><span v-if="licenseInfo.copyright">{{ licenseInfo.copyright }}</span>
                                <br v-if="licenseInfo.url"/><a class="work-break-all" target="_blank" :href="licenseInfo.url" v-if="licenseInfo.url">{{ licenseInfo.url }}</a>
                                <br v-if="licenseInfo.licenseUrl"/><span class="work-break-all" v-if="licenseInfo.licenseUrl">{{ licenseInfo.license || 'License' }}: </span>
                                <a target="_blank" :href="licenseInfo.licenseUrl">{{ licenseInfo.licenseUrl }}</a>
                            </p>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';
import { useAboutPageBase } from '@/views/base/AboutPageBase.ts';

import {
    mdiWebRefresh
} from '@mdi/js';

const { tt, getAllLanguageOptions } = useI18n();
const {
    clientVersion,
    clientVersionMatchServerVersion,
    clientBuildTime,
    exchangeRatesData,
    isUserCustomExchangeRates,
    mapProviderName,
    mapProviderWebsite,
    contributors,
    licenseLines,
    thirdPartyLicenses,
    refreshBrowserCache,
    init
} = useAboutPageBase();


const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(false));

init();
</script>

<style>
.contributors-table {
    border-collapse: collapse;

    > thead > tr {
        > th:not(:first-child) {
            padding-inline-start: 10px;
        }

        > th:not(:last-child) {
            padding-inline-end: 10px;
        }
    }

    > thead > tr > th,
    > tbody > tr > td {
        padding: 4px 8px;
        border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
        text-align: start;
    }
}
</style>
