import Vue from 'vue';
import Vuex from 'vuex';

import VueI18n from 'vue-i18n';
import PincodeInput from 'vue-pincode-input';
import VueClipboard from 'vue-clipboard2';

import moment from 'moment-timezone';

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
import Framework7Calendar from 'framework7/components/calendar/calendar';
import Framework7Picker from 'framework7/components/picker/picker';
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

import api from './consts/api.js';
import datetime from './consts/datetime.js';
import timezone from './consts/timezone.js';
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
import { getAllLanguages, getLanguage, getDefaultLanguage, getI18nOptions, getLocalizedError, getLocalizedErrorParameters } from './lib/i18n.js';

import stores from './store/index.js';

import localizedFilter from './filters/localized.js';
import momentFilter from './filters/moment.js';
import percentFilter from './filters/percent.js';
import formatFilter from './filters/format.js';
import optionNameFilter from './filters/optionName.js';
import itemFieldContentFilter from './filters/itemFieldContent.js';
import languageNameFilter from './filters/languageName.js';
import currencyFilter from './filters/currency.js';
import exchangeRateFilter from './filters/exchangeRate.js';
import utcOffsetFilter from './filters/utcOffset.js';
import textLimitFilter from './filters/textLimit.js';
import iconFilter from './filters/icon.js';
import iconStyleFilter from './filters/iconStyle.js';
import defaultIconColorFilter from './filters/defaultIconColor.js';
import accountIconFilter from './filters/accountIcon.js';
import accountIconStyleFilter from './filters/accountIconStyle.js';
import categoryIconFilter from './filters/categoryIcon.js';
import categoryIconStyleFilter from './filters/categoryIconStyle.js';
import tokenDeviceFilter from './filters/tokenDevice.js';
import tokenIconFilter from './filters/tokenIcon.js';

import PieChart from './components/mobile/PieChart.vue';
import PasswordInputSheet from './components/mobile/PasswordInputSheet.vue';
import PasscodeInputSheet from './components/mobile/PasscodeInputSheet.vue';
import PinCodeInputSheet from './components/mobile/PinCodeInputSheet.vue';
import DateRangeSelectionSheet from './components/mobile/DateRangeSelectionSheet.vue';
import ListItemSelectionSheet from './components/mobile/ListItemSelectionSheet.vue';
import TwoColumnListItemSelectionSheet from './components/mobile/TwoColumnListItemSelectionSheet.vue';
import TreeViewSelectionSheet from './components/mobile/TreeViewSelectionSheet.vue';
import IconSelectionSheet from './components/mobile/IconSelectionSheet.vue';
import ColorSelectionSheet from './components/mobile/ColorSelectionSheet.vue';
import InformationSheet from './components/mobile/InformationSheet.vue';
import NumberPadSheet from './components/mobile/NumberPadSheet.vue';
import TransactionTagSelectionSheet from './components/mobile/TransactionTagSelectionSheet.vue';

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
Framework7.use(Framework7Calendar);
Framework7.use(Framework7Picker);
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
Vue.component('TransactionTagSelectionSheet', TransactionTagSelectionSheet);

Vue.filter('localized', (value, options) => localizedFilter({ i18n }, value, options));
Vue.filter('moment', (value, format, options) => momentFilter(value, format, options));
Vue.filter('percent', (value, precision, lowPrecisionValue) => percentFilter(value, precision, lowPrecisionValue));
Vue.filter('format', (value, format) => formatFilter(value, format));
Vue.filter('optionName', (value, options, keyField, nameField, defaultName) => optionNameFilter(value, options, keyField, nameField, defaultName));
Vue.filter('itemFieldContent', (value, fieldName, defaultValue, translate) => itemFieldContentFilter({ i18n }, value, fieldName, defaultValue, translate));
Vue.filter('languageName', (languageCode) => languageNameFilter(languageCode));
Vue.filter('currency', (value, currencyCode) => currencyFilter({ i18n }, value, currencyCode));
Vue.filter('exchangeRate', (value, currentCurrency, allExchangeRates) => exchangeRateFilter(value, currentCurrency, allExchangeRates));
Vue.filter('utcOffset', (value) => utcOffsetFilter(value));
Vue.filter('textLimit', (value, maxLength) => textLimitFilter(value, maxLength));
Vue.filter('icon', (value, iconType) => iconFilter(value, iconType));
Vue.filter('iconStyle', (value, iconType, defaultColor, additionalFieldName) => iconStyleFilter(value, iconType, defaultColor, additionalFieldName));
Vue.filter('defaultIconColor', (value, defaultColor) => defaultIconColorFilter(value, defaultColor));
Vue.filter('accountIcon', (value) => accountIconFilter(value));
Vue.filter('accountIconStyle', (value, defaultColor, additionalFieldName) => accountIconStyleFilter(value, defaultColor, additionalFieldName));
Vue.filter('categoryIcon', (value) => categoryIconFilter(value));
Vue.filter('categoryIconStyle', (value, defaultColor, additionalFieldName) => categoryIconStyleFilter(value, defaultColor, additionalFieldName));
Vue.filter('tokenDevice', (value) => tokenDeviceFilter(value));
Vue.filter('tokenIcon', (value) => tokenIconFilter(value));

const store = new Vuex.Store(stores);
const i18n = new VueI18n(getI18nOptions());

Vue.prototype.$version = version.getVersion();
Vue.prototype.$buildTime = version.getBuildTime();

Vue.prototype.$licenses = {
    license: licenses.getLicense(),
    thirdPartyLicenses: licenses.getThirdPartyLicenses()
};

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
    defaultTimezoneOffset: utils.getTimezoneOffset(),
    defaultTimezoneOffsetMinutes: utils.getTimezoneOffsetMinutes(),
    getDefaultLanguage: getDefaultLanguage,
    getAllLanguages: getAllLanguages,
    getAllLongMonthNames: function () {
        return [
            i18n.t('datetime.January.long'),
            i18n.t('datetime.February.long'),
            i18n.t('datetime.March.long'),
            i18n.t('datetime.April.long'),
            i18n.t('datetime.May.long'),
            i18n.t('datetime.June.long'),
            i18n.t('datetime.July.long'),
            i18n.t('datetime.August.long'),
            i18n.t('datetime.September.long'),
            i18n.t('datetime.October.long'),
            i18n.t('datetime.November.long'),
            i18n.t('datetime.December.long')
        ];
    },
    getAllShortMonthNames: function () {
        return [
            i18n.t('datetime.January.short'),
            i18n.t('datetime.February.short'),
            i18n.t('datetime.March.short'),
            i18n.t('datetime.April.short'),
            i18n.t('datetime.May.short'),
            i18n.t('datetime.June.short'),
            i18n.t('datetime.July.short'),
            i18n.t('datetime.August.short'),
            i18n.t('datetime.September.short'),
            i18n.t('datetime.October.short'),
            i18n.t('datetime.November.short'),
            i18n.t('datetime.December.short')
        ];
    },
    getAllLongWeekdayNames: function () {
        return [
            i18n.t('datetime.Sunday.long'),
            i18n.t('datetime.Monday.long'),
            i18n.t('datetime.Tuesday.long'),
            i18n.t('datetime.Wednesday.long'),
            i18n.t('datetime.Thursday.long'),
            i18n.t('datetime.Friday.long'),
            i18n.t('datetime.Saturday.long')
        ];
    },
    getAllShortWeekdayNames: function () {
        return [
            i18n.t('datetime.Sunday.short'),
            i18n.t('datetime.Monday.short'),
            i18n.t('datetime.Tuesday.short'),
            i18n.t('datetime.Wednesday.short'),
            i18n.t('datetime.Thursday.short'),
            i18n.t('datetime.Friday.short'),
            i18n.t('datetime.Saturday.short')
        ];
    },
    getInputTimeIntlDateTimeFormatOptions: function () {
        const hourMinuteFormat = i18n.t('input-format.datetime.long');
        const is24HourFormat = hourMinuteFormat.indexOf('H') > 0;
        const hour2Digits = (hourMinuteFormat.indexOf('HH') > 0) || (hourMinuteFormat.indexOf('hh') > 0);
        const minute2Digits = hourMinuteFormat.indexOf(':mm') > 0;

        return {
            hour12: !is24HourFormat,
            hour: hour2Digits ? '2-digit' : 'numeric',
            minute: minute2Digits ? '2-digit' : 'numeric'
        }
    },
    getLanguage: getLanguage,
    setLanguage: function (locale) {
        if (settings.getLanguage() !== locale) {
            settings.setLanguage(locale);
        }

        i18n.locale = locale;
        moment.locale(locale);
        services.setLocale(locale);
        document.querySelector('html').setAttribute('lang', locale);

        const defaultCurrency = i18n.t('default.currency');
        const defaultFirstDayOfWeekName = i18n.t('default.firstDayOfWeek');
        let defaultFirstDayOfWeek = datetime.defaultFirstDayOfWeek;

        if (datetime.allWeekDays[defaultFirstDayOfWeekName]) {
            defaultFirstDayOfWeek = datetime.allWeekDays[defaultFirstDayOfWeekName].type;
        }

        store.dispatch('updateLocalizedDefaultSettings', { defaultCurrency, defaultFirstDayOfWeek });

        return locale;
    },
    getTimezone: function () {
        return settings.getTimezone();
    },
    setTimezone: function (timezone) {
        if (timezone) {
            settings.setTimezone(timezone);
            moment.tz.setDefault(timezone);
        } else {
            settings.setTimezone('');
            moment.tz.setDefault();
        }
    },
    getAllTimezones: function (includeSystemDefault) {
        const allTimezones = timezone.all;
        const allTimezoneInfos = [];

        for (let i = 0; i < allTimezones.length; i++) {
            allTimezoneInfos.push({
                name: allTimezones[i].timezoneName,
                utcOffset: (allTimezones[i].timezoneName !== 'Etc/GMT' ? utils.getTimezoneOffset(allTimezones[i].timezoneName) : ''),
                utcOffsetMinutes: utils.getTimezoneOffsetMinutes(allTimezones[i].timezoneName),
                displayName: i18n.t(`timezone.${allTimezones[i].displayName}`)
            });
        }

        if (includeSystemDefault) {
            allTimezoneInfos.push({
                name: '',
                utcOffset: this.defaultTimezoneOffset,
                utcOffsetMinutes: this.defaultTimezoneOffsetMinutes,
                displayName: i18n.t('System Default')
            });
        }

        allTimezoneInfos.sort(function(c1, c2){
            const utcOffset1 = parseInt(c1.utcOffset.replace(':', ''));
            const utcOffset2 = parseInt(c2.utcOffset.replace(':', ''));

            if (utcOffset1 !== utcOffset2) {
                return utcOffset1 - utcOffset2;
            }

            return c1.displayName.localeCompare(c2.displayName);
        })

        return allTimezoneInfos;
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
            this.setLanguage(settings.getLanguage());
        } else {
            logger.info(`No language is set, use browser default ${getDefaultLanguage()}`);
            this.setLanguage(getDefaultLanguage());
        }

        if (settings.getTimezone()) {
            logger.info(`Current timezone is ${settings.getTimezone()}`);
            this.setTimezone(settings.getTimezone());
        } else {
            logger.info(`No timezone is set, use browser default ${utils.getTimezoneOffset()} (maybe ${moment.tz.guess(true)})`);
        }
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
Vue.prototype.$routeBackOnError = function (errorPropertyName) {
    const self = this;
    const router = self.$f7router;

    const unwatch = self.$watch(errorPropertyName, () => {
        if (self[errorPropertyName]) {
            setTimeout(() => {
                if (unwatch) {
                    unwatch();
                }

                router.back();
            }, 200);
        }
    }, {
        immediate: true
    });
};

Vue.prototype.$user = userstate;

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
