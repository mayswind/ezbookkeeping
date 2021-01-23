<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Device & Sessions')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': tokens.length < 2 }" :text="$t('Logout All')" @click="revokeAll"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list media-list>
                    <f7-list-item class="list-item-media-valign-middle"
                                  title="Current"
                                  text="Device Name (Browser xx.x.xxxx.xx)">
                        <f7-icon slot="media" f7="device_phone_portrait"></f7-icon>
                        <small slot="after">MM/DD/YYYY HH:mm:ss</small>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list media-list>
                    <f7-list-item class="list-item-media-valign-middle" swipeout
                                  v-for="token in tokens"
                                  :key="token.tokenId"
                                  :id="token | tokenDomId"
                                  :title="token | tokenTitle | localized"
                                  :text="token | tokenDevice | localized">
                        <f7-icon slot="media" :f7="token | tokenIcon"></f7-icon>
                        <small slot="after">{{ token.createdAt | moment($t('format.datetime.long')) }}</small>
                        <f7-swipeout-actions right v-if="!token.isCurrent">
                            <f7-swipeout-button color="red" :text="$t('Log Out')" @click="revoke(token)"></f7-swipeout-button>
                        </f7-swipeout-actions>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            tokens: [],
            loading: true
        };
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        self.$store.dispatch('getAllTokens').then(tokens => {
            self.tokens = tokens;
            self.loading = false;
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
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
        revoke(token) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            self.$confirm('Are you sure you want to logout from this session?', () => {
                self.$showLoading();

                self.$store.dispatch('revokeToken', {
                    tokenId: token.tokenId
                }).then(() => {
                    self.$hideLoading();

                    app.swipeout.delete($$(`#${self.$options.filters.tokenDomId(token)}`), () => {
                        for (let i = 0; i < self.tokens.length; i++) {
                            if (self.tokens[i].tokenId === token.tokenId) {
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
        }
    },
    filters: {
        tokenTitle(token) {
            if (token.isCurrent) {
                return 'Current';
            }

            return 'Other Device';
        },
        tokenDomId(token) {
            return 'token_' + token.tokenId.replace(/:/g, '_');
        }
    }
};
</script>
