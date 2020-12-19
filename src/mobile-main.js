import Vue from 'vue';
import VueI18n from 'vue-i18n';
import VueI18nFilter from 'vue-i18n-filter';
import Framework7 from 'framework7/framework7.esm.bundle.js';
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js';
import PincodeInput from 'vue-pincode-input';
import VueMoment from 'vue-moment';
import VueClipboard from 'vue-clipboard2';

import moment from 'moment';
import 'moment/min/locales';

import 'framework7/css/framework7.bundle.css';
import 'framework7-icons';

import 'line-awesome/dist/line-awesome/css/line-awesome.css';

import { getAllLanguages, getLanguage, getDefaultLanguage, getI18nOptions, getLocalizedError, getLocalizedErrorParameters } from './lib/i18n.js';
import currency from './consts/currency.js';
import colors from './consts/color.js';
import icons from './consts/icon.js';
import account from './consts/account.js';
import transaction from './consts/transaction.js';
import category from './consts/category.js';
import licenses from './consts/licenses.js';
import version from './lib/version.js';
import logger from './lib/logger.js';
import settings from './lib/settings.js';
import services from './lib/services.js';
import userstate from './lib/userstate.js';
import exchangeRates from './lib/exchangeRates.js';
import webauthn from './lib/webauthn.js';
import utils from './lib/utils.js';
import currencyFilter from './filters/currency.js';
import iconFilter from './filters/icon.js';
import accountIconFilter from './filters/accountIcon.js';
import categoryIconFilter from './filters/categoryIcon.js';
import tokenDeviceFilter from './filters/tokenDevice.js';
import tokenIconFilter from './filters/tokenIcon.js';

import PasswordInputSheet from "./components/mobile/PasswordInputSheet.vue";
import PasscodeInputSheet from "./components/mobile/PasscodeInputSheet.vue";
import PinCodeInputSheet from "./components/mobile/PinCodeInputSheet.vue";
import ListItemSelectionSheet from "./components/mobile/ListItemSelectionSheet.vue";
import TwoColumnListItemSelectionSheet from "./components/mobile/TwoColumnListItemSelectionSheet.vue";
import IconSelectionSheet from "./components/mobile/IconSelectionSheet.vue";
import ColorSelectionSheet from "./components/mobile/ColorSelectionSheet.vue";
import InformationSheet from "./components/mobile/InformationSheet.vue";
import NumberPadSheet from "./components/mobile/NumberPadSheet.vue";
import App from './Mobile.vue';

Vue.use(VueI18n);
Vue.use(VueI18nFilter);
Vue.use(VueMoment, { moment });
Vue.use(VueClipboard);
Vue.component('PincodeInput', PincodeInput);
Vue.component('PasswordInputSheet', PasswordInputSheet);
Vue.component('PasscodeInputSheet', PasscodeInputSheet);
Vue.component('PinCodeInputSheet', PinCodeInputSheet);
Vue.component('ListItemSelectionSheet', ListItemSelectionSheet);
Vue.component('TwoColumnListItemSelectionSheet', TwoColumnListItemSelectionSheet);
Vue.component('IconSelectionSheet', IconSelectionSheet);
Vue.component('ColorSelectionSheet', ColorSelectionSheet);
Vue.component('InformationSheet', InformationSheet);
Vue.component('NumberPadSheet', NumberPadSheet);
Framework7.use(Framework7Vue);

const i18n = new VueI18n(getI18nOptions());

Vue.prototype.$version = version.getVersion;
Vue.prototype.$buildTime = version.getBuildTime;

Vue.prototype.$licenses = licenses;
Vue.prototype.$constants = {
    currency: currency,
    colors: colors,
    icons: icons,
    account: account,
    transaction: transaction,
    category: category,
};

Vue.prototype.$utilities = utils;
Vue.prototype.$logger = logger;
Vue.prototype.$webauthn = webauthn;
Vue.prototype.$settings = settings;
Vue.prototype.$locale = {
    getDefaultLanguage: getDefaultLanguage,
    getAllLanguages: getAllLanguages,
    getLanguage: getLanguage,
    setLanguage: function (locale) {
        if (settings.getLanguage() !== locale) {
            settings.setLanguage(locale);
        }

        i18n.locale = locale;
        moment.locale(locale);
        services.setLocale(locale);
        document.querySelector('html').setAttribute('lang', locale);
        return locale;
    },
    getAllCurrencies: function () {
        const allCurrencyCodes = currency.all;
        const allCurrencies = [];

        for (let currencyCode in allCurrencyCodes) {
            if (!Object.prototype.hasOwnProperty.call(allCurrencyCodes, currencyCode)) {
                return;
            }

            allCurrencies.push({
                code: currencyCode,
                displayName: i18n.t(`currency.${currencyCode}`)
            });
        }

        allCurrencies.sort(function(c1, c2){
            return c1.displayName.localeCompare(c2.displayName);
        })

        return allCurrencies;
    },
    init: function () {
        if (settings.getLanguage()) {
            logger.info(`Current language is ${settings.getLanguage()}`);
        } else {
            logger.info(`No language is set, use browser default ${getDefaultLanguage()}`);
        }

        this.setLanguage(settings.getLanguage() || getDefaultLanguage());
    }
};

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
        animate: settings.isEnableAnimate(),
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
        animate: settings.isEnableAnimate(),
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
        parameters = getLocalizedErrorParameters(localizedError.parameters, s => i18n.t(s));
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
Vue.prototype.$exchangeRates = exchangeRates;
Vue.prototype.$user = userstate;

Vue.filter('currency', (value, currencyCode) => currencyFilter({ i18n }, value, currencyCode));
Vue.filter('icon', (value, iconType) => iconFilter(value, iconType));
Vue.filter('accountIcon', (value) => accountIconFilter(value));
Vue.filter('categoryIcon', (value) => categoryIconFilter(value));
Vue.filter('tokenDevice', (value) => tokenDeviceFilter(value));
Vue.filter('tokenIcon', (value) => tokenIconFilter(value));

Vue.prototype.$locale.init();

if (userstate.isUserLogined()) {
    if (!settings.isEnableApplicationLock()) {
        // refresh token if user is logined
        services.refreshToken();

        // auto refresh exchange rates data
        if (settings.isAutoUpdateExchangeRatesData()) {
            services.autoRefreshLatestExchangeRates();
        }
    }
}

new Vue({
    el: '#app',
    i18n: i18n,
    render: h => h(App),
    mounted: function () {
        const app = this.$f7;
        const $$ = app.$;

        app.on('pageBeforeOut',  () => {
            if ($$('.modal-in').length) {
                app.actions.close('.actions-modal.modal-in', false);
                app.dialog.close('.dialog.modal-in', false);
                app.popover.close('.popover.modal-in', false);
                app.popup.close('.popup.modal-in', false);
                app.sheet.close('.sheet-modal.modal-in', false);
            }
        });
    }
});
