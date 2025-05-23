import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createI18n } from 'vue-i18n';

import Framework7 from 'framework7/lite';
import Framework7Dialog from 'framework7/components/dialog';
import Framework7Popup from 'framework7/components/popup';
import Framework7LoginScreen from 'framework7/components/login-screen';
import Framework7Popover from 'framework7/components/popover';
import Framework7Actions from 'framework7/components/actions';
import Framework7Sheet from 'framework7/components/sheet';
import Framework7Notification from 'framework7/components/notification';
import Framework7Toast from 'framework7/components/toast';
import Framework7Preloader from 'framework7/components/preloader';
import Framework7Progressbar from 'framework7/components/progressbar';
import Framework7Sortable from 'framework7/components/sortable';
import Framework7Swipeout from 'framework7/components/swipeout';
import Framework7Accordion from 'framework7/components/accordion';
import Framework7Card from 'framework7/components/card';
import Framework7Chip from 'framework7/components/chip';
import Framework7Form from 'framework7/components/form';
import Framework7Input from 'framework7/components/input';
import Framework7Checkbox from 'framework7/components/checkbox';
import Framework7Radio from 'framework7/components/radio';
import Framework7Toggle from 'framework7/components/toggle';
import Framework7Range from 'framework7/components/range';
import Framework7Grid from 'framework7/components/grid';
import Framework7Picker from 'framework7/components/picker';
import Framework7InfiniteScroll from 'framework7/components/infinite-scroll';
import Framework7PullToRefresh from 'framework7/components/pull-to-refresh';
import Framework7Searchbar from 'framework7/components/searchbar';
import Framework7Tooltip from 'framework7/components/tooltip';
import Framework7Skeleton from 'framework7/components/skeleton';
import Framework7Treeview from 'framework7/components/treeview';
import Framework7Typography from 'framework7/components/typography';
import Framework7Swiper from 'framework7/components/swiper';
import Framework7PhotoBrowser from 'framework7/components/photo-browser';
// @ts-expect-error there is a function called "registerComponents" in the framework7-vue package, but it is not declared in the type definition file
import Framework7Vue, { registerComponents } from 'framework7-vue/bundle';

import 'framework7/css';
import 'framework7/components/dialog/css';
import 'framework7/components/popup/css';
import 'framework7/components/login-screen/css';
import 'framework7/components/popover/css';
import 'framework7/components/actions/css';
import 'framework7/components/sheet/css';
import 'framework7/components/notification/css';
import 'framework7/components/toast/css';
import 'framework7/components/preloader/css';
import 'framework7/components/progressbar/css';
import 'framework7/components/sortable/css';
import 'framework7/components/swipeout/css';
import 'framework7/components/accordion/css';
import 'framework7/components/card/css';
import 'framework7/components/chip/css';
import 'framework7/components/form/css';
import 'framework7/components/input/css';
import 'framework7/components/checkbox/css';
import 'framework7/components/radio/css';
import 'framework7/components/toggle/css';
import 'framework7/components/range/css';
import 'framework7/components/grid/css';
import 'framework7/components/picker/css';
import 'framework7/components/infinite-scroll/css';
import 'framework7/components/pull-to-refresh/css';
import 'framework7/components/searchbar/css';
import 'framework7/components/tooltip/css';
import 'framework7/components/skeleton/css';
import 'framework7/components/treeview/css';
import 'framework7/components/typography/css';
import 'framework7/components/swiper/css';
import 'framework7/components/photo-browser/css';

import 'framework7-icons';
import 'line-awesome/dist/line-awesome/css/line-awesome.css';

import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

import { getI18nOptions } from '@/locales/helpers.ts';

import PinCodeInput from '@/components/common/PinCodeInput.vue';
import MapView from '@/components/common/MapView.vue';

import ItemIcon from '@/components/mobile/ItemIcon.vue';
import LanguageSelectButton from '@/components/mobile/LanguageSelectButton.vue';
import PieChart from '@/components/mobile/PieChart.vue';
import TrendsBarChart from '@/components/mobile/TrendsBarChart.vue';
import PinCodeInputSheet from '@/components/mobile/PinCodeInputSheet.vue';
import PasswordInputSheet from '@/components/mobile/PasswordInputSheet.vue';
import PasscodeInputSheet from '@/components/mobile/PasscodeInputSheet.vue';
import DateTimeSelectionSheet from '@/components/mobile/DateTimeSelectionSheet.vue';
import DateSelectionSheet from '@/components/mobile/DateSelectionSheet.vue';
import FiscalYearStartSelectionSheet from '@/components/mobile/FiscalYearStartSelectionSheet.vue';
import DateRangeSelectionSheet from '@/components/mobile/DateRangeSelectionSheet.vue';
import MonthSelectionSheet from '@/components/mobile/MonthSelectionSheet.vue';
import MonthRangeSelectionSheet from '@/components/mobile/MonthRangeSelectionSheet.vue';
import ListItemSelectionSheet from '@/components/mobile/ListItemSelectionSheet.vue';
import ListItemSelectionPopup from '@/components/mobile/ListItemSelectionPopup.vue';
import TwoColumnListItemSelectionSheet from '@/components/mobile/TwoColumnListItemSelectionSheet.vue';
import TreeViewSelectionSheet from '@/components/mobile/TreeViewSelectionSheet.vue';
import IconSelectionSheet from '@/components/mobile/IconSelectionSheet.vue';
import ColorSelectionSheet from '@/components/mobile/ColorSelectionSheet.vue';
import InformationSheet from '@/components/mobile/InformationSheet.vue';
import NumberPadSheet from '@/components/mobile/NumberPadSheet.vue';
import MapSheet from '@/components/mobile/MapSheet.vue';
import TransactionTagSelectionSheet from '@/components/mobile/TransactionTagSelectionSheet.vue';
import ScheduleFrequencySheet from '@/components/mobile/ScheduleFrequencySheet.vue';

import TextareaAutoSize from '@/directives/mobile/textareaAutoSize.ts';

import '@/styles/mobile/global.css';
import '@/styles/mobile/font-size-default.css';
import '@/styles/mobile/font-size-small.css';
import '@/styles/mobile/font-size-large.css';
import '@/styles/mobile/font-size-x-large.css';
import '@/styles/mobile/font-size-xx-large.css';
import '@/styles/mobile/font-size-xxx-large.css';
import '@/styles/mobile/font-size-xxxx-large.css';
import '@/styles/mobile/amount-color.css';

import App from '@/MobileApp.vue';

Framework7.use([
    Framework7Dialog,
    Framework7Popup,
    Framework7LoginScreen,
    Framework7Popover,
    Framework7Actions,
    Framework7Sheet,
    Framework7Notification,
    Framework7Toast,
    Framework7Preloader,
    Framework7Progressbar,
    Framework7Sortable,
    Framework7Swipeout,
    Framework7Accordion,
    Framework7Card,
    Framework7Chip,
    Framework7Form,
    Framework7Input,
    Framework7Checkbox,
    Framework7Radio,
    Framework7Toggle,
    Framework7Range,
    Framework7Grid,
    Framework7Picker,
    Framework7InfiniteScroll,
    Framework7PullToRefresh,
    Framework7Searchbar,
    Framework7Tooltip,
    Framework7Skeleton,
    Framework7Treeview,
    Framework7Typography,
    Framework7Swiper,
    Framework7PhotoBrowser,
    Framework7Vue
]);

const app = createApp(App);
const pinia = createPinia();
const i18n = createI18n(getI18nOptions());
registerComponents(app);
app.use(pinia);
app.use(i18n);

app.component('VueDatePicker', VueDatePicker);

app.component('PinCodeInput', PinCodeInput);
app.component('MapView', MapView);

app.component('ItemIcon', ItemIcon);
app.component('LanguageSelectButton', LanguageSelectButton);
app.component('PieChart', PieChart);
app.component('TrendsBarChart', TrendsBarChart);
app.component('PinCodeInputSheet', PinCodeInputSheet);
app.component('PasswordInputSheet', PasswordInputSheet);
app.component('PasscodeInputSheet', PasscodeInputSheet);
app.component('DateTimeSelectionSheet', DateTimeSelectionSheet);
app.component('DateSelectionSheet', DateSelectionSheet);
app.component('FiscalYearStartSelectionSheet', FiscalYearStartSelectionSheet);
app.component('DateRangeSelectionSheet', DateRangeSelectionSheet);
app.component('MonthSelectionSheet', MonthSelectionSheet);
app.component('MonthRangeSelectionSheet', MonthRangeSelectionSheet);
app.component('ListItemSelectionSheet', ListItemSelectionSheet);
app.component('ListItemSelectionPopup', ListItemSelectionPopup);
app.component('TwoColumnListItemSelectionSheet', TwoColumnListItemSelectionSheet);
app.component('TreeViewSelectionSheet', TreeViewSelectionSheet);
app.component('IconSelectionSheet', IconSelectionSheet);
app.component('ColorSelectionSheet', ColorSelectionSheet);
app.component('InformationSheet', InformationSheet);
app.component('NumberPadSheet', NumberPadSheet);
app.component('MapSheet', MapSheet);
app.component('TransactionTagSelectionSheet', TransactionTagSelectionSheet);
app.component('ScheduleFrequencySheet', ScheduleFrequencySheet);

app.directive('TextareaAutoSize', TextareaAutoSize);

app.mount('#app');
