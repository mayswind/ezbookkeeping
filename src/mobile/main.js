import Vue from 'vue';
import VueI18n from 'vue-i18n';
import axios from 'axios';
import Framework7 from 'framework7/framework7.esm.bundle.js';
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js';

import 'framework7/css/framework7.bundle.css';
import 'framework7-icons';

import './assets/css/custom.css';

import i18n from '../common/i18n.js';
import settings from '../common/settings.js';
import App from './App.vue';

Vue.use(VueI18n);
Framework7.use(Framework7Vue);

const i18nInstance = new VueI18n(i18n.i18nOptions);

Vue.prototype.$setLanguage = function (locale) {
    if (settings.getLanguage() !== locale) {
        settings.setLanguage(locale);
    }

    i18nInstance.locale = locale;
    axios.defaults.headers.common['Accept-Language'] = locale;
    document.querySelector('html').setAttribute('lang', locale);
    return locale;
}

Vue.prototype.$setLanguage(settings.getLanguage() || i18n.getDefaultLanguage());

new Vue({
    el: '#app',
    i18n: i18nInstance,
    render: h => h(App),
});
