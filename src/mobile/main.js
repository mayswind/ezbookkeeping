import Vue from 'vue';
import VueI18n from 'vue-i18n';
import axios from 'axios';
import Framework7 from 'framework7/framework7.esm.bundle.js';
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js';

import 'framework7/css/framework7.bundle.css';
import 'framework7-icons';

import './assets/css/custom.css';

import i18nOptions from '../common/i18n.js';
import App from './App.vue';

Vue.use(VueI18n);
Framework7.use(Framework7Vue);

const i18n = new VueI18n(i18nOptions);

Vue.prototype.languages = {
    setI18nLanguage: lang => {
        i18n.locale = lang;
        axios.defaults.headers.common['Accept-Language'] = lang;
        document.querySelector('html').setAttribute('lang', lang);
        return lang;
    }
}

new Vue({
    el: '#app',
    i18n: i18n,
    render: h => h(App),
});
