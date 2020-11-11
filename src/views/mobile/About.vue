<template>
    <f7-page>
        <f7-navbar :title="$t('About')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Version')" :after="version"></f7-list-item>
                    <f7-list-item external :title="$t('Official Website')" after="https://github.com/mayswind/lab"
                                  link="https://github.com/mayswind/lab" target="_blank"></f7-list-item>
                    <f7-list-item :title="$t('License')" link="#" popup-open=".license-popup"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-popup class="license-popup">
            <f7-page>
                <f7-navbar :title="$t('License')">
                    <f7-nav-right>
                        <f7-link popup-close :text="$t('Close')"></f7-link>
                    </f7-nav-right>
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
                    <p>
                        <strong>Gin</strong><br/>
                        <span>Copyright (c) 2014 Manuel Martínez-Almeida</span><br/>
                        <span class="work-break-all">https://github.com/gin-gonic/gin</span><br/>
                        <span class="work-break-all">License: https://github.com/gin-gonic/gin/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>GZIP gin's middleware</strong><br/>
                        <span>Copyright (c) 2017 Gin-Gonic</span><br/>
                        <span class="work-break-all">https://github.com/gin-contrib/gzip</span><br/>
                        <span class="work-break-all">License: https://github.com/gin-contrib/gzip/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>xorm</strong><br/>
                        <span>Copyright (c) 2013 - 2015 The Xorm Authors</span><br/>
                        <span class="work-break-all">https://xorm.io/xorm</span><br/>
                        <span class="work-break-all">License: https://gitea.com/xorm/xorm/src/branch/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Logrus</strong><br/>
                        <span>Copyright (c) 2014 Simon Eskildsen</span><br/>
                        <span class="work-break-all">https://github.com/sirupsen/logrus</span><br/>
                        <span class="work-break-all">License: https://github.com/sirupsen/logrus/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>cli</strong><br/>
                        <span>Copyright (c) 2016 Jeremy Saenz & Contributors</span><br/>
                        <span class="work-break-all">https://github.com/urfave/cli</span><br/>
                        <span class="work-break-all">License: https://github.com/urfave/cli/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>INI</strong><br/>
                        <span>Copyright 2014 Unknwon</span><br/>
                        <span class="work-break-all">https://gopkg.in/ini.v1</span><br/>
                        <span class="work-break-all">License: https://github.com/go-ini/ini/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Package validator</strong><br/>
                        <span>Copyright (c) 2015 Dean Karn</span><br/>
                        <span class="work-break-all">https://github.com/go-playground/validator</span><br/>
                        <span class="work-break-all">License: https://github.com/go-playground/validator/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>jwt-go</strong><br/>
                        <span>Copyright (c) 2012 Dave Grijalva</span><br/>
                        <span class="work-break-all">https://github.com/dgrijalva/jwt-go</span><br/>
                        <span class="work-break-all">License: https://github.com/dgrijalva/jwt-go/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>otp</strong><br/>
                        <span class="work-break-all">https://github.com/pquerna/otp</span><br/>
                        <span class="work-break-all">License: https://github.com/pquerna/otp/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Go-MySQL-Driver</strong><br/>
                        <span class="work-break-all">https://github.com/go-sql-driver/mysql</span><br/>
                        <span class="work-break-all">License: https://github.com/go-sql-driver/mysql/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>pq</strong><br/>
                        <span>Copyright (c) 2011-2013, 'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany</span><br/>
                        <span class="work-break-all">https://github.com/lib/pq</span><br/>
                        <span class="work-break-all">License: https://github.com/lib/pq/blob/master/LICENSE.md</span>
                    </p>
                    <p>
                        <strong>go-sqlite3</strong><br/>
                        <span>Copyright (c) 2014 Yasuhiro Matsumoto</span><br/>
                        <span class="work-break-all">https://github.com/mattn/go-sqlite3</span><br/>
                        <span class="work-break-all">License: https://github.com/mattn/go-sqlite3/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Go Cryptography</strong><br/>
                        <span>Copyright (c) 2009 The Go Authors. All rights reserved.</span><br/>
                        <span class="work-break-all">https://golang.org/x/crypto</span><br/>
                        <span class="work-break-all">License: https://github.com/golang/crypto/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Testify</strong><br/>
                        <span>Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors.</span><br/>
                        <span class="work-break-all">https://github.com/stretchr/testify</span><br/>
                        <span class="work-break-all">License: https://github.com/stretchr/testify/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>vue</strong><br/>
                        <span>Copyright (c) 2013-present, Yuxi (Evan) You</span><br/>
                        <span class="work-break-all">https://github.com/vuejs/vue</span><br/>
                        <span class="work-break-all">License: https://github.com/vuejs/vue/blob/dev/LICENSE</span>
                    </p>
                    <p>
                        <strong>vue-i18n</strong><br/>
                        <span>Copyright (c) 2016 kazuya kawaguchi</span><br/>
                        <span class="work-break-all">https://github.com/kazupon/vue-i18n</span><br/>
                        <span class="work-break-all">License: https://github.com/kazupon/vue-i18n/blob/v8.x/LICENSE</span>
                    </p>
                    <p>
                        <strong>vue-i18n-filter</strong><br/>
                        <span>Copyright (c) 2018 +v</span><br/>
                        <span class="work-break-all">https://github.com/chiaweilee/vue-i18n-filter</span><br/>
                        <span class="work-break-all">License: https://github.com/chiaweilee/vue-i18n-filter/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>vue-moment</strong><br/>
                        <span>Copyright (c) 2017 Brock Petrie</span><br/>
                        <span class="work-break-all">https://github.com/brockpetrie/vue-moment</span><br/>
                        <span class="work-break-all">License: https://github.com/brockpetrie/vue-moment/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>vue-clipboard2</strong><br/>
                        <span>Copyright (c) 2017 Inndy &lt;inndy \dot tw \at gmail \dot com&gt;</span><br/>
                        <span class="work-break-all">https://github.com/Inndy/vue-clipboard2</span><br/>
                        <span class="work-break-all">License: https://github.com/Inndy/vue-clipboard2/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>core-js</strong><br/>
                        <span>Copyright (c) 2014-2020 Denis Pushkarev</span><br/>
                        <span class="work-break-all">https://github.com/zloirock/core-js</span><br/>
                        <span class="work-break-all">License: https://github.com/zloirock/core-js/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Framework7</strong><br/>
                        <span>Copyright (c) 2014 Vladimir Kharlampidi</span><br/>
                        <span class="work-break-all">https://framework7.io/</span><br/>
                        <span class="work-break-all">License: https://github.com/framework7io/framework7/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Framework7-vue</strong><br/>
                        <span>Copyright (c) 2014 Vladimir Kharlampidi</span><br/>
                        <span class="work-break-all">https://framework7.io/vue/</span><br/>
                        <span class="work-break-all">License: https://github.com/framework7io/framework7/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Framework7-icons</strong><br/>
                        <span>Copyright (c) 2016 Vladimir Kharlampidi</span><br/>
                        <span class="work-break-all">https://framework7.io/icons/</span><br/>
                        <span class="work-break-all">License: https://github.com/framework7io/framework7-icons</span>
                    </p>
                    <p>
                        <strong>axios</strong><br/>
                        <span>Copyright (c) 2014-present Matt Zabriskie</span><br/>
                        <span class="work-break-all">https://github.com/axios/axios</span><br/>
                        <span class="work-break-all">License: https://github.com/axios/axios/blob/master/LICENSE</span>
                    </p>
                    <p>
                        <strong>Moment.js</strong><br/>
                        <span>Copyright (c) JS Foundation and other contributors</span><br/>
                        <span class="work-break-all">https://momentjs.com</span><br/>
                        <span class="work-break-all">License: https://github.com/moment/moment/blob/develop/LICENSE</span>
                    </p>
                    <p>
                        <strong>js-cookie</strong><br/>
                        <span>Copyright (c) 2018 Copyright 2018 Klaus Hartl, Fagner Brack, GitHub Contributors</span><br/>
                        <span class="work-break-all">https://github.com/js-cookie/js-cookie</span><br/>
                        <span class="work-break-all">License: https://github.com/js-cookie/js-cookie/blob/master/LICENSE</span>
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
            return 'v' + this.$version();
        }
    }
}
</script>
