import Vue from 'vue';
import VueI18n from 'vue-i18n';
import Framework7 from 'framework7/framework7.esm.bundle.js';
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js';

import 'framework7/css/framework7.bundle.css';
import 'framework7-icons';

import { getAllLanguages, getLanguage, getDefaultLanguage, getI18nOptions } from './lib/i18n.js';
import settings from './lib/settings.js';
import services from './lib/services.js';
import userstate from './lib/userstate.js';
import App from './Mobile.vue';

Vue.use(VueI18n);
Framework7.use(Framework7Vue);

const i18n = new VueI18n(getI18nOptions());

Vue.prototype.$settings = settings;
Vue.prototype.$getDefaultLanguage = getDefaultLanguage;
Vue.prototype.$getAllLanguages = getAllLanguages;
Vue.prototype.$getLanguage = getLanguage;
Vue.prototype.$setLanguage = function (locale) {
    if (settings.getLanguage() !== locale) {
        settings.setLanguage(locale);
    }

    i18n.locale = locale;
    services.setLocale(locale);
    document.querySelector('html').setAttribute('lang', locale);
    return locale;
};
Vue.prototype.$alert = function (message, confirmCallback) {
    this.$f7.dialog.create({
        title: i18n.t('global.app.title'),
        text: i18n.t(message),
        buttons: [
            {
                text: i18n.t('OK'),
                onClick: confirmCallback
            }
        ]
    }).open();
};
Vue.prototype.$confirm = function (message, confirmCallback, cancelCallback) {
    this.$f7.dialog.create({
        title: i18n.t('global.app.title'),
        text: i18n.t(message),
        buttons: [
            {
                text: i18n.t('Cancel'),
                onClick: cancelCallback
            },
            {
                text: i18n.t('OK'),
                onClick: confirmCallback
            }
        ]
    }).open();
};
Vue.prototype.$services = services;
Vue.prototype.$user = userstate;

Vue.prototype.$setLanguage(settings.getLanguage() || getDefaultLanguage());

// refresh token if user is logined
if (userstate.isUserLogined()) {
    services.refreshToken().then(response => {
        const data = response.data;

        if (data && data.success && data.result) {
            userstate.updateToken(data.result);
        }
    });
}

new Vue({
    el: '#app',
    i18n: i18n,
    render: h => h(App),
    mounted: function () {
        const app = this.$f7;

        window.addEventListener('popstate', () => {
            if (document.querySelectorAll('.modal-in').length > 0) {
                app.dialog.close();
            }
        }, false);

        document.addEventListener('keydown', (event) => {
            if (event.key === 'Escape' || event.key === 'Esc' || event.keyCode === 27 || event.which === 27) {
                if (document.querySelectorAll('.modal-in').length > 0) {
                    app.dialog.close();
                }
            }
        }, false);
    }
});
