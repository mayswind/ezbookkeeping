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
            <f7-card-content :padding="false">
                <f7-list media-list>
                    <f7-list-item title="Token Name" after="MM/DD/YYYY HH:mm:ss"
                        text="Mozilla/5.0 (OS Name & Version or Device Name; CPU Architecture ...) AppleWebKit/Version (KHTML, like Gecko) Software Name/Version Software Name/Version ..."></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content :padding="false">
                <f7-list media-list>
                    <f7-list-item swipeout v-for="token in tokens" :key="token.tokenId" :id="token | tokenDomId" :title="token | tokenTitle | t" :after="token.createdAt | moment($t('format.datetime.long'))" :text="token.userAgent">
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

        self.$services.getTokens().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$alert('Unable to get session list', () => {
                    router.back();
                });
                return;
            }

            self.tokens = data.result;
            self.loading = false;
        }).catch(error => {
            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$alert({ error: error.response.data }, () => {
                    router.back();
                });
            } else if (!error.processed) {
                self.$alert('Unable to get session list', () => {
                    router.back();
                });
            }
        });
    },
    methods: {
        reload(done) {
            const self = this;

            self.$services.getTokens().then(response => {
                done();

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get session list');
                    return;
                }

                self.tokens = data.result;
            }).catch(error => {
                done();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to get session list');
                }
            });
        },
        revoke(token) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            self.$confirm('Are you sure you want to logout from this session?', () => {
                self.$showLoading();

                self.$services.revokeToken({
                    tokenId: token.tokenId
                }).then(response => {
                    self.$hideLoading();
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        self.$alert('Unable to logout from this session');
                        return;
                    }

                    app.swipeout.delete($$(`#${self.$options.filters.tokenDomId(token)}`), () => {
                        for (let i = 0; i < self.tokens.length; i++) {
                            if (self.tokens[i].tokenId === token.tokenId) {
                                self.tokens.splice(i, 1);
                            }
                        }
                    });
                }).catch(error => {
                    self.$hideLoading();

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        self.$alert({error: error.response.data});
                    } else if (!error.processed) {
                        self.$alert('Unable to logout from this session');
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

                self.$services.revokeAllTokens().then(response => {
                    self.$hideLoading();
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        self.$alert('Unable to logout all other sessions');
                        return;
                    }

                    for (let i = self.tokens.length - 1; i >= 0; i--) {
                        if (!self.tokens[i].isCurrent) {
                            self.tokens.splice(i, 1);
                        }
                    }

                    self.$toast('You have logged out all other sessions');
                }).catch(error => {
                    self.$hideLoading();

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        self.$alert({error: error.response.data});
                    } else if (!error.processed) {
                        self.$alert('Unable to logout all other sessions');
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
