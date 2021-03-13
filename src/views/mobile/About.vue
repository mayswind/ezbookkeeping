<template>
    <f7-page>
        <f7-navbar :title="$t('About')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Version')" :after="version"></f7-list-item>
                    <f7-list-item :title="$t('Build Time')" :after="buildTime | moment($t('format.datetime.long'))"></f7-list-item>
                    <f7-list-item external :title="$t('Official Website')" after="https://github.com/mayswind/lab"
                                  link="https://github.com/mayswind/lab" target="_blank"></f7-list-item>
                    <f7-list-item :title="$t('License')" link="#" popup-open=".license-popup"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-popup class="license-popup">
            <f7-page>
                <f7-navbar>
                    <f7-nav-left>
                        <f7-link popup-close :text="$t('Cancel')"></f7-link>
                    </f7-nav-left>
                    <f7-nav-title :title="$t('License')"></f7-nav-title>
                </f7-navbar>
                <f7-block>
                    <p><p>MIT License</p>
                    <p></p>
                    <p>Copyright (c) 2020 MaysWind (i@mayswind.net)</p>
                    <p></p>
                    <p>Permission is hereby granted, free of charge, to any person obtaining a copy
                        of this software and associated documentation files (the "Software"), to deal
                        in the Software without restriction, including without limitation the rights
                        to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
                        copies of the Software, and to permit persons to whom the Software is
                        furnished to do so, subject to the following conditions:</p>
                    <p></p>
                    <p>The above copyright notice and this permission notice shall be included in all
                        copies or substantial portions of the Software.</p>
                    <p></p>
                    <p>THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
                        IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
                        FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
                        AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
                        LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
                        OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
                        SOFTWARE.</p>
                    <p></p>
                    <hr/>
                    <p>
                        <span>lab also contains additional third party software.</span><br/>
                        <span>All the third party software included or linked is redistributed under the terms and conditions of their original licenses.</span>
                    </p>
                    <p></p>
                    <p v-for="license in licenses" :key="license.name">
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
            return this.$buildTime;
        },
        licenses() {
            return this.$licenses;
        }
    }
}
</script>
