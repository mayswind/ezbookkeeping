<template>
    <f7-page name="settings">
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>
        <f7-block-title>{{ userNickName }}</f7-block-title>
        <f7-list>
            <f7-list-item :title="$t('User Profile')" link="/user/profile"></f7-list-item>
            <f7-list-item :title="$t('Two-Factor Authentication')" link="/user/2fa"></f7-list-item>
            <f7-list-item :title="$t('Device & Sessions')" link="/user/sessions"></f7-list-item>
            <f7-list-button :class="{ 'disabled': logouting }" @click="logout">{{ $t('Log Out') }}</f7-list-button>
        </f7-list>
        <f7-block-title>{{ $t('Application') }}</f7-block-title>
        <f7-list>
            <f7-list-item
                :title="$t('Language')"
                smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done') }">
                <select v-model="currentLocale">
                    <option v-for="(lang, locale) in allLanguages"
                            :key="locale"
                            :value="locale">{{ lang.displayName }}</option>
                </select>
            </f7-list-item>
            <f7-list-item>
                <span>{{ $t('Enable Animate') }}</span>
                <f7-toggle :checked="isEnableAnimate" @toggle:change="isEnableAnimate = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-item>
                <span>{{ $t('Enable Auto Dark Mode') }}</span>
                <f7-toggle :checked="isEnableAutoDarkMode" @toggle:change="isEnableAutoDarkMode = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-item :title="$t('About')" link="/about" :after="version"></f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            logouting: false,
            allLanguages: self.$getAllLanguages()
        };
    },
    computed: {
        version() {
            return 'v' + this.$version();
        },
        userNickName() {
            return this.$user.getUserNickName() || this.$user.getUserName() || this.$t('User');
        },
        currentLocale: {
            get: function () {
                return this.$i18n.locale;
            },
            set: function (value) {
                this.$setLanguage(value);
            }
        },
        isEnableAnimate: {
            get: function () {
                return this.$settings.isEnableAnimate();
            },
            set: function (value) {
                if (value !== this.$settings.isEnableAnimate()) {
                    this.$settings.setEnableAnimate(value);
                    location.reload();
                }
            }
        },
        isEnableAutoDarkMode: {
            get: function () {
                return this.$settings.isEnableAutoDarkMode();
            },
            set: function (value) {
                if (value !== this.$settings.isEnableAutoDarkMode()) {
                    this.$settings.setEnableAutoDarkMode(value);
                    location.reload();
                }
            }
        }
    },
    methods: {
        logout() {
            const self = this;
            const router = self.$f7router;

            self.$confirm('Are you sure you want to log out?', () => {
                self.logouting = true;
                self.$showLoading(() => self.logouting);

                self.$services.logout().then(response => {
                    self.logouting = false;
                    self.$hideLoading();
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        self.$alert('Unable to logout');
                        return;
                    }

                    self.$user.clearToken();
                    router.navigate('/');
                }).catch(error => {
                    self.logouting = false;
                    self.$hideLoading();

                    if (error && error.processed) {
                        return;
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        self.$alert({ error: error.response.data });
                    } else if (!error.processed) {
                        self.$alert('Unable to logout');
                    }
                });
            });
        }
    }
};
</script>
