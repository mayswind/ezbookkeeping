import Vue from 'vue';
import Vuex from 'vuex';

import VueI18n from 'vue-i18n';
import PincodeInput from 'vue-pincode-input';
import VueMoment from 'vue-moment';
import VueClipboard from 'vue-clipboard2';

import moment from 'moment';

import Framework7 from 'framework7/framework7-lite.esm.js';
import Framework7Dialog from 'framework7/components/dialog/dialog';
import Framework7Popup from 'framework7/components/popup/popup';
import Framework7LoginScreen from 'framework7/components/login-screen/login-screen';
import Framework7Popover from 'framework7/components/popover/popover';
import Framework7Actions from 'framework7/components/actions/actions';
import Framework7Sheet from 'framework7/components/sheet/sheet';
import Framework7Toast from 'framework7/components/toast/toast';
import Framework7Preloader from 'framework7/components/preloader/preloader';
import Framework7Progressbar from 'framework7/components/progressbar/progressbar';
import Framework7Sortable from 'framework7/components/sortable/sortable';
import Framework7Swipeout from 'framework7/components/swipeout/swipeout';
import Framework7Accordion from 'framework7/components/accordion/accordion';
import Framework7Card from 'framework7/components/card/card';
import Framework7Chip from 'framework7/components/chip/chip';
import Framework7Form from 'framework7/components/form/form';
import Framework7Input from 'framework7/components/input/input';
import Framework7Checkbox from 'framework7/components/checkbox/checkbox';
import Framework7Radio from 'framework7/components/radio/radio';
import Framework7Toggle from 'framework7/components/toggle/toggle';
import Framework7SmartSelect from 'framework7/components/smart-select/smart-select';
import Framework7Grid from 'framework7/components/grid/grid';
import Framework7InfiniteScroll from 'framework7/components/infinite-scroll/infinite-scroll';
import Framework7PullToRefresh from 'framework7/components/pull-to-refresh/pull-to-refresh';
import Framework7Searchbar from 'framework7/components/searchbar/searchbar';
import Framework7Tooltip from 'framework7/components/tooltip/tooltip';
import Framework7Skeleton from 'framework7/components/skeleton/skeleton';
import Framework7Menu from 'framework7/components/menu/menu';
import Framework7Treeview from 'framework7/components/treeview/treeview';
import Framework7Typography from 'framework7/components/typography/typography';
import Framework7Vue from 'framework7-vue/framework7-vue.esm.bundle.js';

import 'framework7/css/framework7.bundle.css';
import 'framework7-icons';

import 'line-awesome/dist/line-awesome/css/line-awesome.css';

import { getAllLanguages, getLanguage, getDefaultLanguage, getI18nOptions, getLocalizedError, getLocalizedErrorParameters } from './lib/i18n.js';
import api from './consts/api.js';
import datetime from './consts/datetime.js';
import currency from './consts/currency.js';
import colors from './consts/color.js';
import icons from './consts/icon.js';
import account from './consts/account.js';
import transaction from './consts/transaction.js';
import category from './consts/category.js';
import statistics from './consts/statistics.js';
import licenses from './lib/licenses.js';
import version from './lib/version.js';
import logger from './lib/logger.js';
import settings from './lib/settings.js';
import services from './lib/services.js';
import userstate from './lib/userstate.js';
import webauthn from './lib/webauthn.js';
import utils from './lib/utils.js';
import stores from './store/index.js';
import localizedFilter from './filters/localized.js';
import percentFilter from './filters/percent.js';
import itemFieldContentFilter from './filters/itemFieldContent.js';
import currencyFilter from './filters/currency.js';
import iconFilter from './filters/icon.js';
import iconStyleFilter from './filters/iconStyle.js';
import defaultIconColorFilter from './filters/defaultIconColor.js';
import accountIconFilter from './filters/accountIcon.js';
import accountIconStyleFilter from './filters/accountIconStyle.js';
import categoryIconFilter from './filters/categoryIcon.js';
import categoryIconStyleFilter from './filters/categoryIconStyle.js';
import tokenDeviceFilter from './filters/tokenDevice.js';
import tokenIconFilter from './filters/tokenIcon.js';

import PieChart from "./components/mobile/PieChart.vue";
import PasswordInputSheet from "./components/mobile/PasswordInputSheet.vue";
import PasscodeInputSheet from "./components/mobile/PasscodeInputSheet.vue";
import PinCodeInputSheet from "./components/mobile/PinCodeInputSheet.vue";
import DateRangeSelectionSheet from "./components/mobile/DateRangeSelectionSheet.vue";
import ListItemSelectionSheet from "./components/mobile/ListItemSelectionSheet.vue";
import TwoColumnListItemSelectionSheet from "./components/mobile/TwoColumnListItemSelectionSheet.vue";
import TreeViewSelectionSheet from "./components/mobile/TreeViewSelectionSheet.vue";
import IconSelectionSheet from "./components/mobile/IconSelectionSheet.vue";
import ColorSelectionSheet from "./components/mobile/ColorSelectionSheet.vue";
import InformationSheet from "./components/mobile/InformationSheet.vue";
import NumberPadSheet from "./components/mobile/NumberPadSheet.vue";

import App from './Mobile.vue';

Framework7.use(Framework7Dialog);
Framework7.use(Framework7Popup);
Framework7.use(Framework7LoginScreen);
Framework7.use(Framework7Popover);
Framework7.use(Framework7Actions);
Framework7.use(Framework7Sheet);
Framework7.use(Framework7Toast);
Framework7.use(Framework7Preloader);
Framework7.use(Framework7Progressbar);
Framework7.use(Framework7Sortable);
Framework7.use(Framework7Swipeout);
Framework7.use(Framework7Accordion);
Framework7.use(Framework7Card);
Framework7.use(Framework7Chip);
Framework7.use(Framework7Form);
Framework7.use(Framework7Input);
Framework7.use(Framework7Checkbox);
Framework7.use(Framework7Radio);
Framework7.use(Framework7Toggle);
Framework7.use(Framework7SmartSelect);
Framework7.use(Framework7Grid);
Framework7.use(Framework7InfiniteScroll);
Framework7.use(Framework7PullToRefresh);
Framework7.use(Framework7Searchbar);
Framework7.use(Framework7Tooltip);
Framework7.use(Framework7Skeleton);
Framework7.use(Framework7Menu);
Framework7.use(Framework7Treeview);
Framework7.use(Framework7Typography);
Framework7.use(Framework7Vue);

Vue.use(Vuex);
Vue.use(VueI18n);
Vue.use(VueMoment, { moment });
Vue.use(VueClipboard);

Vue.component('PincodeInput', PincodeInput);
Vue.component('PieChart', PieChart);
Vue.component('PasswordInputSheet', PasswordInputSheet);
Vue.component('PasscodeInputSheet', PasscodeInputSheet);
Vue.component('PinCodeInputSheet', PinCodeInputSheet);
Vue.component('DateRangeSelectionSheet', DateRangeSelectionSheet);
Vue.component('ListItemSelectionSheet', ListItemSelectionSheet);
Vue.component('TwoColumnListItemSelectionSheet', TwoColumnListItemSelectionSheet);
Vue.component('TreeViewSelectionSheet', TreeViewSelectionSheet);
Vue.component('IconSelectionSheet', IconSelectionSheet);
Vue.component('ColorSelectionSheet', ColorSelectionSheet);
Vue.component('InformationSheet', InformationSheet);
Vue.component('NumberPadSheet', NumberPadSheet);

const store = new Vuex.Store(stores);
const i18n = new VueI18n(getI18nOptions());

Vue.prototype.$version = version.getVersion();
Vue.prototype.$buildTime = version.getBuildTime();

Vue.prototype.$licenses = licenses.getLicenses();
Vue.prototype.$constants = {
    api: api,
    datetime: datetime,
    currency: currency,
    colors: colors,
    icons: icons,
    account: account,
    transaction: transaction,
    category: category,
    statistics: statistics,
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

Vue.prototype.$user = userstate;

Vue.filter('localized', (value, options) => localizedFilter({ i18n }, value, options));
Vue.filter('percent', (value, precision, lowPrecisionValue) => percentFilter(value, precision, lowPrecisionValue));
Vue.filter('itemFieldContent', (value, fieldName, defaultValue, translate) => itemFieldContentFilter({ i18n }, value, fieldName, defaultValue, translate));
Vue.filter('currency', (value, currencyCode) => currencyFilter({ i18n }, value, currencyCode));
Vue.filter('icon', (value, iconType) => iconFilter(value, iconType));
Vue.filter('iconStyle', (value, iconType, defaultColor) => iconStyleFilter(value, iconType, defaultColor));
Vue.filter('defaultIconColor', (value, defaultColor) => defaultIconColorFilter(value, defaultColor));
Vue.filter('accountIcon', (value) => accountIconFilter(value));
Vue.filter('accountIconStyle', (value, defaultColor) => accountIconStyleFilter(value, defaultColor));
Vue.filter('categoryIcon', (value) => categoryIconFilter(value));
Vue.filter('categoryIconStyle', (value, defaultColor) => categoryIconStyleFilter(value, defaultColor));
Vue.filter('tokenDevice', (value) => tokenDeviceFilter(value));
Vue.filter('tokenIcon', (value) => tokenIconFilter(value));

Vue.prototype.$locale.init();

new Vue({
    el: '#app',
    i18n: i18n,
    store: store,
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
