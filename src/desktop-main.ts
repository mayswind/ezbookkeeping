import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createI18n } from 'vue-i18n';

import { createVuetify } from 'vuetify';
import { VAlert } from 'vuetify/components/VAlert';
import { VApp } from 'vuetify/components/VApp';
import { VAvatar } from 'vuetify/components/VAvatar';
import { VAutocomplete } from 'vuetify/components/VAutocomplete';
import { VBadge } from 'vuetify/components/VBadge';
import { VBtn } from 'vuetify/components/VBtn';
import { VBtnGroup } from 'vuetify/components/VBtnGroup';
import { VBtnToggle } from 'vuetify/components/VBtnToggle';
import { VCard, VCardActions, VCardItem, VCardSubtitle, VCardText, VCardTitle } from 'vuetify/components/VCard';
import { VCheckbox, VCheckboxBtn } from 'vuetify/components/VCheckbox';
import { VChip } from 'vuetify/components/VChip';
import { VDataTable } from 'vuetify/components/VDataTable';
import { VDialog } from 'vuetify/components/VDialog';
import { VDivider } from 'vuetify/components/VDivider';
import { VExpansionPanel, VExpansionPanelText, VExpansionPanelTitle, VExpansionPanels } from 'vuetify/components/VExpansionPanel';
import { VForm } from 'vuetify/components/VForm';
import { VContainer, VCol, VRow, VSpacer } from 'vuetify/components/VGrid';
import { VIcon } from 'vuetify/components/VIcon';
import { VImg } from 'vuetify/components/VImg';
import { VInput } from 'vuetify/components/VInput';
import { VLabel } from 'vuetify/components/VLabel';
import { VLayout } from 'vuetify/components/VLayout';
import { VList, VListGroup, VListImg, VListItem, VListItemAction, VListItemMedia, VListItemSubtitle, VListItemTitle, VListSubheader } from 'vuetify/components/VList';
import { VMain } from 'vuetify/components/VMain';
import { VMenu } from 'vuetify/components/VMenu';
import { VNavigationDrawer } from 'vuetify/components/VNavigationDrawer';
import { VOverlay } from 'vuetify/components/VOverlay';
import { VPagination } from 'vuetify/components/VPagination';
import { VProgressCircular } from 'vuetify/components/VProgressCircular';
import { VProgressLinear } from 'vuetify/components/VProgressLinear';
import { VSelect } from 'vuetify/components/VSelect';
import { VSkeletonLoader } from 'vuetify/components/VSkeletonLoader';
import { VSlideGroup, VSlideGroupItem } from 'vuetify/components/VSlideGroup';
import { VSnackbar } from 'vuetify/components/VSnackbar';
import { VSwitch } from 'vuetify/components/VSwitch';
import { VTabs, VTab } from 'vuetify/components/VTabs';
import { VTable } from 'vuetify/components/VTable';
import { VTextarea } from 'vuetify/components/VTextarea';
import { VTextField } from 'vuetify/components/VTextField';
import { VToolbar } from 'vuetify/components/VToolbar';
import { VTooltip } from 'vuetify/components/VTooltip';
import { VWindow, VWindowItem } from 'vuetify/components/VWindow';

import { aliases, mdi } from 'vuetify/iconsets/mdi-svg';
import 'vuetify/styles';

import * as echarts from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart, BarChart, PieChart } from 'echarts/charts';
import {
    GridComponent,
    TooltipComponent,
    LegendComponent,
} from 'echarts/components';
import VChart from 'vue-echarts';

import 'line-awesome/dist/line-awesome/css/line-awesome.css';

import { PerfectScrollbar } from 'vue3-perfect-scrollbar';

import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

import draggable from 'vuedraggable';

import router from '@/router/desktop.ts';

import { getI18nOptions } from '@/locales/helpers.ts';

import PinCodeInput from '@/components/common/PinCodeInput.vue';
import MapView from '@/components/common/MapView.vue';

import ItemIcon from '@/components/desktop/ItemIcon.vue';
import BtnVerticalGroup from '@/components/desktop/BtnVerticalGroup.vue';
import AmountInput from '@/components/desktop/AmountInput.vue';
import LanguageSelect from '@/components/desktop/LanguageSelect.vue';
import LanguageSelectButton from '@/components/desktop/LanguageSelectButton.vue';
import CurrencySelect from '@/components/desktop/CurrencySelect.vue';
import DateTimeSelect from '@/components/desktop/DateTimeSelect.vue';
import DateSelect from '@/components/desktop/DateSelect.vue';
import FiscalYearStartSelect from '@/components/desktop/FiscalYearStartSelect.vue';
import ColorSelect from '@/components/desktop/ColorSelect.vue';
import IconSelect from '@/components/desktop/IconSelect.vue';
import TwoColumnSelect from '@/components/desktop/TwoColumnSelect.vue';
import ScheduleFrequencySelect from '@/components/desktop/ScheduleFrequencySelect.vue';
import StepsBar from '@/components/desktop/StepsBar.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import PieChartComponent from '@/components/desktop/PieChart.vue';
import TrendsChartComponent from '@/components/desktop/TrendsChart.vue';
import DateRangeSelectionDialog from '@/components/desktop/DateRangeSelectionDialog.vue';
import MonthRangeSelectionDialog from '@/components/desktop/MonthRangeSelectionDialog.vue';
import SwitchToMobileDialog from '@/components/desktop/SwitchToMobileDialog.vue';

import '@/styles/desktop/template/vuetify/index.scss';
import '@/styles/desktop/template/template/index.scss';
import '@/styles/desktop/template/layout/index.scss';
import '@/styles/desktop/template/layout/component/index.scss';
import '@/styles/desktop/template/layout/_default-layout.scss';
import '@/styles/desktop/global.scss';
import '@/styles/desktop/font-size.scss';
import '@/styles/desktop/amount-color.scss';

import App from './DesktopApp.vue';

const app = createApp(App);
const pinia = createPinia();
const i18n = createI18n(getI18nOptions());
const vuetify = createVuetify({
    components: {
        VAlert,
        VApp,
        VAvatar,
        VAutocomplete,
        VBadge,
        VBtn,
        VBtnGroup,
        VBtnToggle,
        VCard,
        VCardActions,
        VCardItem,
        VCardSubtitle,
        VCardText,
        VCardTitle,
        VCheckbox,
        VCheckboxBtn,
        VChip,
        VDataTable,
        VDialog,
        VDivider,
        VExpansionPanel,
        VExpansionPanelText,
        VExpansionPanelTitle,
        VExpansionPanels,
        VForm,
        VContainer,
        VCol,
        VRow,
        VSpacer,
        VIcon,
        VImg,
        VInput,
        VLabel,
        VLayout,
        VList,
        VListGroup,
        VListImg,
        VListItem,
        VListItemAction,
        VListItemMedia,
        VListItemSubtitle,
        VListItemTitle,
        VListSubheader,
        VMain,
        VMenu,
        VNavigationDrawer,
        VOverlay,
        VPagination,
        VProgressCircular,
        VProgressLinear,
        VSelect,
        VSkeletonLoader,
        VSlideGroup,
        VSlideGroupItem,
        VSnackbar,
        VSwitch,
        VTabs,
        VTab,
        VTable,
        VTextarea,
        VTextField,
        VToolbar,
        VTooltip,
        VWindow,
        VWindowItem
    },
    icons: {
        defaultSet: 'mdi',
        aliases,
        sets: {
            mdi
        }
    },
    defaults: {
        VAlert: {
            VBtn: {
                color: undefined
            }
        },
        VAutocomplete: {
            variant: 'outlined',
            density: 'comfortable',
            color: 'primary',
            hideDetails: 'auto'
        },
        VAvatar: {
            variant: 'flat',
            VIcon: {
                size: 24,
            },
        },
        VBadge: {
            color: 'primary'
        },
        VBtn: {
            color: 'primary'
        },
        VCheckbox: {
            color: 'primary',
            hideDetails: 'auto'
        },
        VChip: {
            elevation: 0
        },
        VList: {
            color: 'primary'
        },
        VPagination: {
            density: 'comfortable',
            activeColor: 'primary'
        },
        VRadio: {
            density: 'comfortable',
            color: 'primary',
            hideDetails: 'auto'
        },
        VSelect: {
            variant: 'outlined',
            density: 'comfortable',
            color: 'primary',
            hideDetails: 'auto'
        },
        VSlider: {
            color: 'primary',
            hideDetails: 'auto'
        },
        VSwitch: {
            inset: true,
            color: 'primary',
            hideDetails: 'auto'
        },
        VProgressCircular: {
            size: 40
        },
        VSnackbar: {
            timeout: 3000
        },
        VTable: {
            hover: true
        },
        VTabs: {
            color: 'primary',
            VSlideGroup: {
                showArrows: true
            }
        },
        VTextarea: {
            variant: 'outlined',
            density: 'comfortable',
            color: 'primary',
            hideDetails: 'auto'
        },
        VTextField: {
            variant: 'outlined',
            density: 'comfortable',
            color: 'primary',
            hideDetails: 'auto'
        },
        VToolbar: {
            color: 'primary'
        },
        VTooltip: {
            location: 'top'
        },
        VWindow: {
            touch: false
        }
    },
    theme: {
        defaultTheme: 'light',
        themes: {
            light: {
                dark: false,
                colors: {
                    'primary': '#c67e48',
                    'primary-darken-1': '#b67443',
                    'on-primary': '#ffffff',
                    'secondary': '#8a8d93',
                    'secondary-darken-1': '#545659',
                    'on-secondary': '#ffffff',
                    'success': '#4cd964',
                    'success-darken-1': '#40b654',
                    'on-success': '#ffffff',
                    'info': '#2196f3',
                    'info-darken-1': '#1e85d7',
                    'on-info': '#ffffff',
                    'warning': '#ff9500',
                    'warning-darken-1': '#de8201',
                    'on-warning': '#ffffff',
                    'error': '#ff3b30',
                    'error-darken-1': '#e1342b',
                    'on-error': '#ffffff',
                    'teal': '#009688',
                    'background': '#faf8f4',
                    'on-background': '#413935',
                    'surface': '#fff',
                    'on-surface': '#413935',
                    'notification-background': '#ffffff',
                    'on-notification-background': '#000',
                    'grey-50': '#fafafa',
                    'grey-100': '#f0f2f8',
                    'grey-200': '#eeeeee',
                    'grey-300': '#e0e0e0',
                    'grey-400': '#bdbdbd',
                    'grey-500': '#9e9e9e',
                    'grey-600': '#757575',
                    'grey-700': '#616161',
                    'grey-800': '#424242',
                    'grey-900': '#212121',
                    'perfect-scrollbar-thumb': '#dbdade',
                    'skin-bordered-background': '#fff',
                    'skin-bordered-surface': '#fff',
                    'expansion-panel-text-custom-bg': '#fafafa'
                },
                variables: {
                    'code-color': '#ff8000',
                    'overlay-scrim-background': '#413935',
                    'tooltip-background': '#212121',
                    'overlay-scrim-opacity': 0.5,
                    'hover-opacity': 0.04,
                    'focus-opacity': 0.1,
                    'selected-opacity': 0.08,
                    'activated-opacity': 0.16,
                    'pressed-opacity': 0.14,
                    'dragged-opacity': 0.1,
                    'disabled-opacity': 0.4,
                    'border-color': '#413f3b',
                    'border-opacity': 0.12,
                    'table-header-color': '#fdfcf9',
                    'high-emphasis-opacity': 0.9,
                    'medium-emphasis-opacity': 0.7,

                    // 👉 shadows
                    'shadow-key-umbra-color': '#413935',
                    'shadow-xs-opacity': '0.16',
                    'shadow-sm-opacity': '0.18',
                    'shadow-md-opacity': '0.20',
                    'shadow-lg-opacity': '0.22',
                    'shadow-xl-opacity': '0.24',
                }
            },
            dark: {
                dark: true,
                colors: {
                    'primary': '#c67e48',
                    'primary-darken-1': '#b67443',
                    'on-primary': '#ffffff',
                    'secondary': '#8a8d93',
                    'secondary-darken-1': '#545659',
                    'on-secondary': '#fff',
                    'success': '#4cd964',
                    'success-darken-1': '#40b654',
                    'on-success': '#ffffff',
                    'info': '#2196f3',
                    'info-darken-1': '#1e85d7',
                    'on-info': '#ffffff',
                    'warning': '#ff9500',
                    'warning-darken-1': '#de8201',
                    'on-warning': '#ffffff',
                    'error': '#ff3b30',
                    'error-darken-1': '#e1342b',
                    'on-error': '#ffffff',
                    'teal': '#009688',
                    'background': '#000000',
                    'on-background': '#fcf0e3',
                    'surface': '#1c1c1d',
                    'on-surface': '#fcf0e3',
                    'notification-background': '#1e1e1e',
                    'on-notification-background': '#fff',
                    'grey-50': '#212121',
                    'grey-100': '#424242',
                    'grey-200': '#616161',
                    'grey-300': '#757575',
                    'grey-400': '#909090',
                    'grey-500': '#a2a2a2',
                    'grey-600': '#b4b4b4',
                    'grey-700': '#c6c6c6',
                    'grey-800': '#d8d8d8',
                    'grey-900': '#eaeaea',
                    'perfect-scrollbar-thumb': '#4a5072',
                    'skin-bordered-background': '#312d4b',
                    'skin-bordered-surface': '#312d4b',
                    'expansion-panel-text-custom-bg': '#373350'
                },
                variables: {
                    'code-color': '#ff8000',
                    'overlay-scrim-background': '#1c1c1d',
                    'tooltip-background': '#d7d7d7',
                    'overlay-scrim-opacity': 0.6,
                    'hover-opacity': 0.04,
                    'focus-opacity': 0.1,
                    'selected-opacity': 0.08,
                    'activated-opacity': 0.16,
                    'pressed-opacity': 0.14,
                    'disabled-opacity': 0.4,
                    'dragged-opacity': 0.1,
                    'border-color': '#edece9',
                    'border-opacity': 0.12,
                    'table-header-color': '#312f2b',
                    'high-emphasis-opacity': 0.9,
                    'medium-emphasis-opacity': 0.7,

                    // 👉 Shadows
                    'shadow-key-umbra-color': '#67615d',
                    'shadow-xs-opacity': '0.20',
                    'shadow-sm-opacity': '0.22',
                    'shadow-md-opacity': '0.24',
                    'shadow-lg-opacity': '0.26',
                    'shadow-xl-opacity': '0.28',
                }
            }
        }
    }
});

echarts.use([
    CanvasRenderer,
    LineChart,
    BarChart,
    PieChart,
    GridComponent,
    TooltipComponent,
    LegendComponent
]);

app.use(pinia);
app.use(i18n);
app.use(vuetify);
app.use(router);

app.component('VChart', VChart);
app.component('PerfectScrollbar', PerfectScrollbar);
app.component('VueDatePicker', VueDatePicker);
app.component('DraggableList', draggable);

app.component('PinCodeInput', PinCodeInput);
app.component('MapView', MapView);

app.component('ItemIcon', ItemIcon);
app.component('BtnVerticalGroup', BtnVerticalGroup);
app.component('AmountInput', AmountInput);
app.component('LanguageSelect', LanguageSelect);
app.component('LanguageSelectButton', LanguageSelectButton);
app.component('CurrencySelect', CurrencySelect);
app.component('DateTimeSelect', DateTimeSelect);
app.component('DateSelect', DateSelect);
app.component('FiscalYearStartSelect', FiscalYearStartSelect);
app.component('ColorSelect', ColorSelect);
app.component('IconSelect', IconSelect);
app.component('TwoColumnSelect', TwoColumnSelect);
app.component('ScheduleFrequencySelect', ScheduleFrequencySelect);
app.component('StepsBar', StepsBar);
app.component('ConfirmDialog', ConfirmDialog);
app.component('SnackBar', SnackBar);
app.component('PieChart', PieChartComponent);
app.component('TrendsChart', TrendsChartComponent);
app.component('DateRangeSelectionDialog', DateRangeSelectionDialog);
app.component('MonthRangeSelectionDialog', MonthRangeSelectionDialog);
app.component('SwitchToMobileDialog', SwitchToMobileDialog);

app.mount('#app');
