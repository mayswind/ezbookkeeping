<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card :title="tt('global.app.title')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Version') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <span class="text-body-1">{{ version }}</span>
                        </v-col>
                    </v-row>
                    <v-row no-gutters v-if="buildTime">
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ tt('Build Time') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <span class="text-body-1">{{ buildTime }}</span>
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

        <v-col cols="12" v-if="exchangeRatesData">
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
                            <v-divider/><br/>
                            <p>
                                <span>ezBookkeeping also contains additional third party software and illustration.</span><br/>
                                <span>All the third party software/illustration included or linked is redistributed under the terms and conditions of their original licenses.</span>
                            </p>
                            <p></p>
                            <p :key="license.name" v-for="license in thirdPartyLicenses">
                                <strong>{{ license.name }}</strong>
                                <br v-if="license.copyright"/><span v-if="license.copyright">{{ license.copyright }}</span>
                                <br v-if="license.url"/><a class="work-break-all" target="_blank" :href="license.url" v-if="license.url">{{ license.url }}</a>
                                <br v-if="license.licenseUrl"/><span class="work-break-all" v-if="license.licenseUrl">License: </span>
                                <a target="_blank" :href="license.licenseUrl">{{ license.licenseUrl }}</a>
                            </p>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { useAboutPageBase } from '@/views/base/AboutPageBase.ts';

const { tt } = useI18n();
const { version, buildTime, exchangeRatesData, mapProviderName, mapProviderWebsite, licenseLines, thirdPartyLicenses } = useAboutPageBase();
</script>
