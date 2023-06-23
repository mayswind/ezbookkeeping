<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Device & Sessions')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': sessions.length < 2 }" :text="$t('Logout All')" @click="revokeAll"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers media-list class="margin-top skeleton-text" v-if="loading">
            <f7-list-item class="list-item-media-valign-middle"
                          title="Current"
                          text="Device Name (Browser xx.x.xxxx.xx)">
                <template #media>
                    <f7-icon f7="device_phone_portrait"></f7-icon>
                </template>
                <template #after>
                    <small>MM/DD/YYYY HH:mm:ss</small>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers media-list class="margin-top" v-else-if="!loading">
            <f7-list-item class="list-item-media-valign-middle" swipeout
                          :id="session.domId"
                          :title="session.deviceType"
                          :text="session.deviceInfo"
                          :key="session.tokenId"
                          v-for="session in sessions">
                <template #media>
                    <f7-icon :f7="session.icon"></f7-icon>
                </template>
                <template #after>
                    <small>{{ session.createdAt }}</small>
                </template>
                <f7-swipeout-actions right v-if="!session.isCurrent">
                    <f7-swipeout-button color="red" :text="$t('Log Out')" @click="revoke(session)"></f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';

import iconConstants from '@/consts/icon.js';
import { isEquals } from '@/lib/common.js';
import { parseDeviceInfo, parseUserAgent } from '@/lib/misc.js';

import { onSwipeoutDeleted } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            tokens: [],
            loading: true,
            loadingError: null
        };
    },
    computed: {
        ...mapStores(useUserStore, useTokensStore),
        sessions() {
            if (!this.tokens) {
                return this.tokens;
            }

            const sessions = [];

            for (let i = 0; i < this.tokens.length; i++) {
                const token = this.tokens[i];

                sessions.push({
                    tokenId: token.tokenId,
                    domId: this.getTokenDomId(token.tokenId),
                    isCurrent: token.isCurrent,
                    deviceType: this.$t(token.isCurrent ? 'Current' : 'Other Device'),
                    deviceInfo: parseDeviceInfo(token.userAgent),
                    icon: this.getTokenIcon(token),
                    createdAt: this.$locale.formatUnixTimeToLongDateTime(this.userStore, token.createdAt)
                });
            }

            return sessions;
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.tokensStore.getAllTokens().then(tokens => {
            self.tokens = tokens;
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            const self = this;

            self.tokensStore.getAllTokens().then(tokens => {
                if (done) {
                    done();
                }

                if (isEquals(self.tokens, tokens)) {
                    self.$toast('Session list is up to date');
                } else {
                    self.$toast('Session list has been updated');
                }

                self.tokens = tokens;
            }).catch(error => {
                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        revoke(session) {
            const self = this;

            self.$confirm('Are you sure you want to logout from this session?', () => {
                self.$showLoading();

                self.tokensStore.revokeToken({
                    tokenId: session.tokenId
                }).then(() => {
                    self.$hideLoading();

                    onSwipeoutDeleted(self.getTokenDomId(session.tokenId), () => {
                        for (let i = 0; i < self.tokens.length; i++) {
                            if (self.tokens[i].tokenId === session.tokenId) {
                                self.tokens.splice(i, 1);
                            }
                        }
                    });
                }).catch(error => {
                    self.$hideLoading();

                    if (!error.processed) {
                        self.$toast(error.message || error);
                    }
                });
            });
        },
        revokeAll() {
            const self = this;

            if (self.tokens.length < 2) {
                return;
            }

            self.$confirm('Are you sure you want to logout all other sessions?', () => {
                self.$showLoading();

                self.tokensStore.revokeAllTokens().then(() => {
                    self.$hideLoading();

                    for (let i = self.tokens.length - 1; i >= 0; i--) {
                        if (!self.tokens[i].isCurrent) {
                            self.tokens.splice(i, 1);
                        }
                    }

                    self.$toast('You have logged out all other sessions');
                }).catch(error => {
                    self.$hideLoading();

                    if (!error.processed) {
                        self.$toast(error.message || error);
                    }
                });
            });
        },
        getTokenIcon(token) {
            const ua = parseUserAgent(token.userAgent);

            if (!ua || !ua.device) {
                return 'device_desktop';
            }

            if (ua.device.type === 'mobile') {
                return 'device_phone_portrait';
            } else if (ua.device.type === 'wearable') {
                return 'device_phone_portrait';
            } else if (ua.device.type === 'tablet') {
                return 'device_tablet_portrait';
            } else if (ua.device.type === 'smarttv') {
                return 'tv';
            } else {
                return 'device_desktop';
            }
        },
        getTokenDomId(tokenId) {
            return 'token_' + tokenId.replace(/:/g, '_');
        }
    }
};
</script>
