<template>
    <f7-page>
        <f7-navbar :title="tt('About')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('global.app.title') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item :title="tt('Version')" :after="version"></f7-list-item>
            <f7-list-item :title="tt('Build Time')" :after="buildTime" v-if="buildTime"></f7-list-item>
            <f7-list-item external :title="tt('Official Website')" link="https://github.com/mayswind/ezbookkeeping" target="_blank"></f7-list-item>
            <f7-list-item external :title="tt('Report Issue')" link="https://github.com/mayswind/ezbookkeeping/issues" target="_blank"></f7-list-item>
            <f7-list-item external :title="tt('Getting help')" link="https://ezbookkeeping.mayswind.net" target="_blank"></f7-list-item>
            <f7-list-item :title="tt('License')" link="#" popup-open=".license-popup"></f7-list-item>
        </f7-list>

        <f7-block-title class="margin-top" v-if="exchangeRatesData">{{ tt('Exchange Rates Data') }}</f7-block-title>
        <f7-list strong inset dividers v-if="exchangeRatesData">
            <f7-list-item external :title="tt('Provider')" :after="exchangeRatesData.dataSource"
                          :link="exchangeRatesData.referenceUrl" target="_blank" v-if="exchangeRatesData.referenceUrl"></f7-list-item>
            <f7-list-item :title="tt('Provider')" :after="exchangeRatesData.dataSource" v-if="!exchangeRatesData.referenceUrl"></f7-list-item>
        </f7-list>

        <f7-block-title class="margin-top" v-if="mapProviderName">{{ tt('Map') }}</f7-block-title>
        <f7-list strong inset dividers v-if="mapProviderName">
            <f7-list-item external :title="tt('Provider')" :after="mapProviderName"
                          :link="mapProviderWebsite" target="_blank" v-if="mapProviderWebsite"></f7-list-item>
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
    </f7-page>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { useAboutPageBase } from '@/views/base/AboutPageBase.ts';

const { tt } = useI18n();
const { version, buildTime, exchangeRatesData, mapProviderName, mapProviderWebsite, licenseLines, thirdPartyLicenses } = useAboutPageBase();
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
