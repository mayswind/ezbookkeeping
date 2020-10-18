import Vue from 'vue'
import VueI18n from 'vue-i18n'
import Framework7 from 'framework7/framework7.esm.bundle.js'
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js'
import 'framework7/css/framework7.bundle.css'
import 'framework7-icons';
import './assets/css/custom.css'

Vue.use(VueI18n)
Framework7.use(Framework7Vue);

import i18nOptions from '../common/i18n.js';
import App from './App.vue'

const i18n = new VueI18n(i18nOptions)

new Vue({
    el: '#app',
    i18n: i18n,
    render: h => h(App),
})
