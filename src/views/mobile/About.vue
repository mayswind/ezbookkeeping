<template>
    <f7-page>
        <f7-navbar :title="$t('About')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list dividers>
                    <f7-list-item :title="$t('Version')" :after="version"></f7-list-item>
                    <f7-list-item :title="$t('Build Time')" :after="buildTime" v-if="buildTime"></f7-list-item>
                    <f7-list-item external :title="$t('Official Website')" link="https://github.com/mayswind/ezbookkeeping" target="_blank"></f7-list-item>
                    <f7-list-item :title="$t('License')" link="#" popup-open=".license-popup"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-popup push swipe-to-close swipe-handler=".swipe-handler" class="license-popup">
            <f7-page>
                <template #fixed>
                    <div class="swipe-handler"></div>
                </template>
                <f7-block-title large>{{ $t('License') }}</f7-block-title>
                <f7-block strong outline>
                    <p>
                        <span v-for="(line, num) in licenseLines" :key="num"
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
                    <p v-for="license in thirdPartyLicenses" :key="license.name">
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

            return this.$utilities.formatUnixTime(this.$buildTime, this.$t('format.datetime.long'));
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
