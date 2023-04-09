<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Device & Sessions')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': sessions.length < 2 }" :text="$t('Logout All')" @click="revokeAll"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list media-list dividers>
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
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list media-list dividers>
                    <f7-list-item class="list-item-media-valign-middle" swipeout
                                  v-for="session in sessions"
                                  :key="session.tokenId"
                                  :id="session.domId"
                                  :title="session.deviceType"
                                  :text="session.deviceInfo">
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
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
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
                    deviceInfo: this.$utilities.parseDeviceInfo(token.userAgent),
                    icon: this.getTokenIcon(token),
                    createdAt: this.$utilities.formatUnixTime(token.createdAt, this.$t('format.datetime.long'))
                });
            }

            return sessions;
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('getAllTokens').then(tokens => {
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

            self.$store.dispatch('getAllTokens').then(tokens => {
                if (done) {
                    done();
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

                self.$store.dispatch('revokeToken', {
                    tokenId: session.tokenId
                }).then(() => {
                    self.$hideLoading();

                    self.$ui.onSwipeoutDeleted(self.getTokenDomId(session.tokenId), () => {
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

                self.$store.dispatch('revokeAllTokens').then(() => {
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
            const ua = this.$utilities.parseUserAgent(token.userAgent);

            if (!ua || !ua.device) {
                return this.$constants.icons.deviceIcons.desktop.f7Icon;
            }

            if (ua.device.type === 'mobile') {
                return this.$constants.icons.deviceIcons.mobile.f7Icon;
            } else if (ua.device.type === 'wearable') {
                return this.$constants.icons.deviceIcons.wearable.f7Icon;
            } else if (ua.device.type === 'tablet') {
                return this.$constants.icons.deviceIcons.tablet.f7Icon;
            } else if (ua.device.type === 'smarttv') {
                return this.$constants.icons.deviceIcons.tv.f7Icon;
            } else {
                return this.$constants.icons.deviceIcons.desktop.f7Icon;
            }
        },
        getTokenDomId(tokenId) {
            return 'token_' + tokenId.replace(/:/g, '_');
        }
    }
};
</script>
