<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('About')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="" v-if="clientVersionMatchServerVersion && !forceShowRefreshBrowserCacheMenu"/>
                <f7-link icon-f7="ellipsis" @click="showRefreshBrowserCacheSheet = true"
                         v-else-if="!clientVersionMatchServerVersion || forceShowRefreshBrowserCacheMenu"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block-title class="margin-top">{{ tt('global.app.title') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item :title="tt('Version')" :after="clientVersion" @click="showVersion"></f7-list-item>
            <f7-list-item :title="tt('Build Time')" :after="clientBuildTime" v-if="clientBuildTime"></f7-list-item>
            <f7-list-item :title="tt('Official Website')" link="#" @click="openExternalUrl('https://github.com/mayswind/ezbookkeeping')"></f7-list-item>
            <f7-list-item :title="tt('Report Issue')" link="#" @click="openExternalUrl('https://github.com/mayswind/ezbookkeeping/issues')"></f7-list-item>
            <f7-list-item :title="tt('Getting help')" link="#" @click="openExternalUrl('https://ezbookkeeping.mayswind.net')"></f7-list-item>
            <f7-list-item :title="tt('License')" link="#" popup-open=".license-popup"></f7-list-item>
        </f7-list>

        <f7-block-title class="margin-top" v-if="exchangeRatesData && !isUserCustomExchangeRates">{{ tt('Exchange Rates Data') }}</f7-block-title>
        <f7-list strong inset dividers v-if="exchangeRatesData && !isUserCustomExchangeRates">
            <f7-list-item :title="tt('Provider')" :after="exchangeRatesData.dataSource" link="#"
                          @click="openExternalUrl(exchangeRatesData.referenceUrl)" v-if="exchangeRatesData.referenceUrl"></f7-list-item>
            <f7-list-item :title="tt('Provider')" :after="exchangeRatesData.dataSource" v-if="!exchangeRatesData.referenceUrl"></f7-list-item>
        </f7-list>

        <f7-block-title class="margin-top" v-if="mapProviderName">{{ tt('Map') }}</f7-block-title>
        <f7-list strong inset dividers v-if="mapProviderName">
            <f7-list-item :title="tt('Provider')" :after="mapProviderName" link="#"
                          @click="openExternalUrl(mapProviderWebsite)" v-if="mapProviderWebsite"></f7-list-item>
            <f7-list-item :title="tt('Provider')" :after="mapProviderName" v-if="!mapProviderWebsite"></f7-list-item>
        </f7-list>

        <f7-popup push with-subnavbar swipe-to-close swipe-handler=".swipe-handler" class="license-popup">
            <f7-page>
                <f7-navbar>
                    <div class="swipe-handler" style="z-index: 10"></div>
                    <f7-subnavbar :title="tt('License') "></f7-subnavbar>
                </f7-navbar>
                <f7-block strong outline class="license-content">
                    <p>
                        <span :key="num" v-for="(line, num) in licenseLines"
                              :style="{ 'display': line ? 'initial' : 'block', 'padding' : line ? '0' : '0 0 1em 0' }">
                            {{ line }}
                        </span>
                    </p>
                    <hr/>
                    <p>
                        <span>ezBookkeeping also contains additional third party software and illustration.</span><br/>
                        <span>All the third party software/illustration included or linked is redistributed under the terms and conditions of their original licenses.</span>
                    </p>
                    <p></p>
                    <p :key="license.name" v-for="license in thirdPartyLicenses">
                        <strong>{{ license.name }}</strong>
                        <br v-if="license.copyright"/><span v-if="license.copyright">{{ license.copyright }}</span>
                        <br v-if="license.url"/><span class="work-break-all" v-if="license.url">{{ license.url }}</span>
                        <br v-if="license.licenseUrl"/><span class="work-break-all" v-if="license.licenseUrl">License: {{ license.licenseUrl }}</span>
                    </p>
                </f7-block>
            </f7-page>
        </f7-popup>

        <f7-actions close-by-outside-click close-on-escape :opened="showRefreshBrowserCacheSheet" @actions:closed="showRefreshBrowserCacheSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="refreshBrowserCache">{{ tt('Refresh Browser Cache') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useAboutPageBase } from '@/views/base/AboutPageBase.ts';

const { tt } = useI18n();
const { showAlert, openExternalUrl } = useI18nUIComponents();
const {
    clientVersion,
    clientVersionMatchServerVersion,
    serverDisplayVersion,
    clientBuildTime,
    exchangeRatesData,
    isUserCustomExchangeRates,
    mapProviderName,
    mapProviderWebsite,
    licenseLines,
    thirdPartyLicenses,
    refreshBrowserCache,
    init
} = useAboutPageBase();

const showRefreshBrowserCacheSheet = ref<boolean>(false);
const versionClickCount = ref<number>(0);
const forceShowRefreshBrowserCacheMenu = computed<boolean>(() => versionClickCount.value >= 5);

function showVersion(): void {
    let versionMessage = `${tt('Frontend Version')}: ${clientVersion}`;

    if (serverDisplayVersion.value) {
        versionMessage += `<br/>${tt('Backend Version')}: ${serverDisplayVersion.value}`;
    }

    versionClickCount.value++;

    if (serverDisplayVersion.value && serverDisplayVersion.value !== 'unknown' && serverDisplayVersion.value !== clientVersion) {
        showAlert(versionMessage);
    }
}

init();
</script>

<style>
.license-popup .navbar-bg {
    background-color: rgb(var(--f7-navbar-bg-color-rgb, var(--f7-bars-bg-color-rgb)));
}

.license-popup .subnavbar {
    background-color: rgb(var(--f7-subnavbar-bg-color-rgb, var(--f7-bars-bg-color-rgb)));
}

.license-popup .subnavbar-title {
    --f7-subnavbar-title-font-size: var(--ebk-license-popup-title-font-size);
}

.license-content {
    font-size: var(--ebk-license-content-font-size);
}
</style>
