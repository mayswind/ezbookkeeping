<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card :title="$t('global.app.title')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Version') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <span class="text-body-1">{{ version }}</span>
                        </v-col>
                    </v-row>
                    <v-row no-gutters v-if="buildTime">
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Build Time') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <span class="text-body-1">{{ buildTime }}</span>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Official Website') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <a class="text-body-1" href="https://github.com/mayswind/ezbookkeeping" target="_blank">
                                https://github.com/mayswind/ezbookkeeping
                            </a>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Report Issue') }}</span>
                        </v-col>
                        <v-col cols="12" md="10" class="mb-6">
                            <a class="text-body-1" href="https://github.com/mayswind/ezbookkeeping/issues" target="_blank">
                                https://github.com/mayswind/ezbookkeeping/issues
                            </a>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Documents') }}</span>
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
            <v-card :title="$t('Exchange Rates Data')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Provider') }}</span>
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
            <v-card :title="$t('Map')">
                <v-card-text>
                    <v-row no-gutters>
                        <v-col cols="12" md="2">
                            <span class="text-body-1">{{ $t('Provider') }}</span>
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
            <v-card :title="$t('License')">
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

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { getMapProvider } from '@/lib/server_settings.ts';
import { getMapWebsite } from '@/lib/map/index.ts';
import { getLicense, getThirdPartyLicenses } from '@/lib/licenses.ts';

export default {
    computed: {
        ...mapStores(useUserStore, useExchangeRatesStore),
        version() {
            return 'v' + this.$version;
        },
        buildTime() {
            if (!this.$buildTime) {
                return this.$buildTime;
            }

            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.$buildTime);
        },
        exchangeRatesData() {
            return this.exchangeRatesStore.latestExchangeRates.data;
        },
        mapProviderName() {
            const provider = getMapProvider();
            return provider ? this.$t(`mapprovider.${provider}`) : '';
        },
        mapProviderWebsite() {
            return getMapWebsite();
        },
        licenseLines() {
            return getLicense().replaceAll(/\r/g, '').split('\n');
        },
        thirdPartyLicenses() {
            return getThirdPartyLicenses();
        }
    }
}
</script>
