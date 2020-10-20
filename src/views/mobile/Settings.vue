<template>
    <f7-page name="home">
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>
        <f7-list>
            <f7-list-item
                :title="$t('Language')"
                smart-select :smart-select-params="{ openIn: 'sheet', sheetCloseLinkText: $t('Done') }">
                <select v-model="currentLocale">
                    <option v-for="(lang, locale) in allLanguages"
                            :key="locale"
                            :value="locale">{{ lang.displayName }}</option>
                </select>
            </f7-list-item>
            <f7-list-button @click="logout">{{ $t('Log Out') }}</f7-list-button>
        </f7-list>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            allLanguages: self.$getAllLanguages()
        };
    },
    computed: {
        currentLocale: {
            get: function () {
                return this.$i18n.locale
            },
            set: function (value) {
                this.$setLanguage(value);
            }
        }
    },
    methods: {
        logout() {
            const self = this;
            const app = self.$f7;
            const router = self.$f7router;

            app.dialog.create({
                title: self.$i18n.t('global.app.title'),
                text: self.$i18n.t('Are you sure you want to log out?'),
                buttons: [
                    {
                        text: self.$i18n.t('Cancel'),
                    },
                    {
                        text: self.$i18n.t('OK'),
                        onClick: () => {
                            this.$user.clearToken();
                            router.navigate('/');
                        }
                    }
                ]
            }).open();
        }
    }
};
</script>
