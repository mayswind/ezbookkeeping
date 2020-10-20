import Vue from 'vue';
import VueI18n from 'vue-i18n';
import axios from 'axios';
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
    axios.defaults.headers.common['Accept-Language'] = locale;
    document.querySelector('html').setAttribute('lang', locale);
    return locale;
}
Vue.prototype.$services = services;
Vue.prototype.$user = userstate;

Vue.prototype.$setLanguage(settings.getLanguage() || getDefaultLanguage());

new Vue({
    el: '#app',
    i18n: i18n,
    render: h => h(App),
});
