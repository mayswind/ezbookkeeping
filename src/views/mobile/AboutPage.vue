<template>
    <f7-page>
        <f7-navbar :title="$t('About')" :back-link="$t('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top">
            <f7-list-item :title="$t('Version')" :after="version"></f7-list-item>
            <f7-list-item :title="$t('Build Time')" :after="buildTime" v-if="buildTime"></f7-list-item>
            <f7-list-item external :title="$t('Official Website')" link="https://github.com/mayswind/ezbookkeeping" target="_blank"></f7-list-item>
            <f7-list-item :title="$t('License')" link="#" popup-open=".license-popup"></f7-list-item>
        </f7-list>

        <f7-popup push with-subnavbar swipe-to-close swipe-handler=".swipe-handler" class="license-popup">
            <f7-page>
                <f7-navbar>
                    <div class="swipe-handler" style="z-index: 10"></div>
                    <f7-subnavbar :title="$t('License') "></f7-subnavbar>
                </f7-navbar>
                <f7-block strong outline>
                    <p>
                        <span :key="num" v-for="(line, num) in licenseLines"
                              :style="{ 'display': line ? 'initial' : 'block', 'padding' : line ? '0' : '0 0 1em 0' }">
                            {{ line }}
                        </span>
                    </p>
                    <hr/>
                    <p>
                        <span>ezBookkeeping also contains additional third party software.</span><br/>
                        <span>All the third party software included or linked is redistributed under the terms and conditions of their original licenses.</span>
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

<script>
export default {
    computed: {
        version() {
            return 'v' + this.$version;
        },
        buildTime() {
            if (!this.$buildTime) {
                return this.$buildTime;
            }

            return this.$utilities.formatUnixTime(this.$buildTime, this.$locale.getLongDateTimeFormat());
        },
        licenseLines() {
            return this.$licenses.license.replaceAll(/\r/g, '').split('\n');
        },
        thirdPartyLicenses() {
            return this.$licenses.thirdPartyLicenses;
        }
    }
}
</script>

<style>
.license-popup .navbar-bg {
    background-color: rgb(var(--f7-navbar-bg-color-rgb, var(--f7-bars-bg-color-rgb)));
}

.license-popup .subnavbar {
    background-color: rgb(var(--f7-subnavbar-bg-color-rgb, var(--f7-bars-bg-color-rgb)));
}

.license-popup .subnavbar-title {
    --f7-subnavbar-title-font-size: 30px;
}
</style>
