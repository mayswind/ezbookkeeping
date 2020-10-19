<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>{{ $t('global.app.title') }}</f7-login-screen-title>
        <f7-list form>
            <f7-list-input
                type="text"
                clear-button
                validate
                required
                :label="$t('Username')"
                :placeholder="$t('Username or Email')"
                :value="username"
                @input="username = $event.target.value"
            ></f7-list-input>
            <f7-list-input
                type="password"
                clear-button
                validate
                required
                :label="$t('Password')"
                :placeholder="$t('Password')"
                :value="password"
                @input="password = $event.target.value"
            ></f7-list-input>
        </f7-list>
        <f7-list>
            <f7-list-button @click="login">{{ $t('Log In') }}</f7-list-button>
            <f7-block-footer>
                <span v-t="'Don\'t have an account?'"></span>&nbsp;
                <f7-link href="/signup">{{ $t('Create an account') }}</f7-link>
                <br/>
                <f7-link href="/forget-pwd">{{ $t('Forget Password?') }}</f7-link>
            </f7-block-footer>
            <f7-block-footer>
            </f7-block-footer>
        </f7-list>
        <f7-button small popover-open=".popover-menu">{{ currentLanguageName }}</f7-button>
        <f7-popover class="popover-menu">
            <f7-list>
                <f7-list-item
                    link="#" popover-close
                    v-for="(lang, locale) in allLanguages" :key="locale"
                    :title="lang.displayName"
                    @click="changeLanguage(locale)"
                ></f7-list-item>
            </f7-list>
        </f7-popover>
    </f7-page>
</template>

<script>
import services from '../../common/services.js'
import userState from '../../common/userstate.js'
import i18n from '../../common/i18n.js';

export default {
    data() {
        return {
            username: '',
            password: '',
            allLanguages: i18n.getAllLanguages()
        };
    },
    computed: {
        currentLanguageName() {
            const currentLocale = this.$i18n.locale;
            let lang = i18n.getLanguage(currentLocale);

            if (!lang) {
                lang = i18n.getLanguage(i18n.getDefaultLanguage());
            }

            return lang.displayName;
        }
    },
    methods: {
        login() {
            const self = this;
            const app = self.$f7;
            const router = self.$f7router;

            if (!this.username) {
                app.dialog.alert(self.$i18n.t('Please input username'));
                return;
            }

            if (!this.password) {
                app.dialog.alert(self.$i18n.t('Please input password'));
                return;
            }

            let hasResponse = false;

            setTimeout(() => {
                if (!hasResponse) {
                    self.$f7.preloader.show();
                }
            }, 200);

            services.authorize({
                loginName: self.username,
                password: self.password
            }).then(response => {
                hasResponse = true;
                self.$f7.preloader.hide();
                const data = response.data;

                if (data && data.success && data.result && typeof(data.result) === 'string') {
                    userState.updateToken(data.result);
                    router.navigate('/');
                } else {
                    app.dialog.alert(self.$i18n.t('Unable to login'));
                }
            }).catch(error => {
                hasResponse = true;
                self.$f7.preloader.hide();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    app.dialog.alert(self.$i18n.t(`error.${error.response.data.errorMessage}`));
                } else {
                    app.dialog.alert(self.$i18n.t('Unable to login'));
                }
            })
        },
        changeLanguage(locale) {
            this.$setLanguage(locale);
        }
    }
};
</script>
