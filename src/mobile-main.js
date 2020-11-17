import Vue from 'vue';
import VueI18n from 'vue-i18n';
import VueI18nFilter from 'vue-i18n-filter'
import Framework7 from 'framework7/framework7.esm.bundle.js';
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js';
import VueMoment from 'vue-moment';
import VueClipboard from 'vue-clipboard2';

import moment from 'moment';
import 'moment/min/locales';

import 'framework7/css/framework7.bundle.css';
import 'framework7-icons';

import { getAllLanguages, getLanguage, getDefaultLanguage, getI18nOptions, getLocalizedError, getLocalizedErrorParameters } from './lib/i18n.js';
import currency from './consts/currency.js';
import colors from './consts/color.js';
import icons from './consts/icon.js';
import account from './consts/account.js';
import licenses from './consts/licenses.js';
import version from './lib/version.js';
import settings from './lib/settings.js';
import services from './lib/services.js';
import userstate from './lib/userstate.js';
import utils from './lib/utils.js';
import currencyFilter from './filters/currency.js';
import accountIconFilter from './filters/accountIcon.js';
import App from './Mobile.vue';

Vue.use(VueI18n);
Vue.use(VueI18nFilter);
Vue.use(VueMoment, { moment });
Vue.use(VueClipboard);
Framework7.use(Framework7Vue);

const i18n = new VueI18n(getI18nOptions());

Vue.prototype.$version = version.getVersion;
Vue.prototype.$licenses = licenses;
Vue.prototype.$constants = {
    currency: currency,
    colors: colors,
    icons: icons,
    account: account,
};
Vue.prototype.$utilities = utils;
Vue.prototype.$settings = settings;
Vue.prototype.$getDefaultLanguage = getDefaultLanguage;
Vue.prototype.$getAllLanguages = getAllLanguages;
Vue.prototype.$getLanguage = getLanguage;
Vue.prototype.$setLanguage = function (locale) {
    if (settings.getLanguage() !== locale) {
        settings.setLanguage(locale);
    }

    i18n.locale = locale;
    moment.locale(locale);
    services.setLocale(locale);
    document.querySelector('html').setAttribute('lang', locale);
    return locale;
};
Vue.prototype.$getAllCurrencies = function () {
    const allCurrencyCodes = currency.all;
    const allCurrencies = [];

    for (let i = 0; i < allCurrencyCodes.length; i++) {
        const code = allCurrencyCodes[i];

        allCurrencies.push({
            code: code,
            displayName: i18n.t(`currency.${code}`)
        });
    }

    allCurrencies.sort(function(c1, c2){
        return c1.displayName.localeCompare(c2.displayName);
    })

    return allCurrencies;
};
Vue.prototype.$isUserRegistrationEnabled = settings.isUserRegistrationEnabled;

Vue.prototype.$alert = function (message, confirmCallback) {
    let parameters = {};

    if (message && message.error) {
        const localizedError = getLocalizedError(message.error);
        message = localizedError.message;
        parameters = getLocalizedErrorParameters(localizedError.parameters, s => i18n.t(s));
    }

    this.$f7.dialog.create({
        title: i18n.t('global.app.title'),
        text: i18n.t(message, parameters),
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
Vue.prototype.$toast = function (message, timeout) {
    let parameters = {};

    if (message && message.error) {
        const localizedError = getLocalizedError(message.error);
        message = localizedError.message;
        parameters = getLocalizedErrorParameters(localizedError.parameters, i18n.t);
    }

    this.$f7.toast.create({
        text: i18n.t(message, parameters),
        position: 'center',
        closeTimeout: timeout || 1500
    }).open();
};
Vue.prototype.$showLoading = function (delayConditionFunc, delayMills) {
    if (!delayConditionFunc) {
        return this.$f7.preloader.show();
    }

    setTimeout(() => {
        if (delayConditionFunc()) {
            this.$f7.preloader.show();
        }
    }, delayMills || 200);
};
Vue.prototype.$hideLoading = function () {
    return this.$f7.preloader.hide();
};

Vue.prototype.$services = services;
Vue.prototype.$user = userstate;

Vue.filter('currency', (value, currencyCode) => currencyFilter({ i18n }, value, currencyCode));
Vue.filter('accountIcon', (value) => accountIconFilter(value));

Vue.prototype.$setLanguage(settings.getLanguage() || getDefaultLanguage());

// refresh token if user is logined
if (userstate.isUserLogined()) {
    services.refreshToken();
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
                app.sheet.close();
                app.popup.close();
                return false;
            }
        }, false);

        document.addEventListener('keydown', (event) => {
            if (event.key === 'Escape' || event.key === 'Esc' || event.keyCode === 27 || event.which === 27) {
                if (document.querySelectorAll('.modal-in').length > 0) {
                    app.dialog.close();
                    return false;
                }
            }
        }, false);
    }
});
