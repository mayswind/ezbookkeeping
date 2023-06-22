import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createI18n } from 'vue-i18n';

import { createVuetify } from 'vuetify';
import { VApp } from 'vuetify/components/VApp';
import { VAvatar } from 'vuetify/components/VAvatar';
import { VBtn } from 'vuetify/components/VBtn';
import { VCard, VCardActions, VCardItem, VCardSubtitle, VCardText, VCardTitle } from 'vuetify/components/VCard';
import { VChip } from 'vuetify/components/VChip';
import { VDialog } from 'vuetify/components/VDialog';
import { VDivider } from 'vuetify/components/VDivider';
import { VForm } from 'vuetify/components/VForm';
import { VContainer, VCol, VRow, VSpacer } from 'vuetify/components/VGrid';
import { VIcon } from 'vuetify/components/VIcon';
import { VImg } from 'vuetify/components/VImg';
import { VInput } from 'vuetify/components/VInput';
import { VList, VListGroup, VListImg, VListItem, VListItemAction, VListItemMedia, VListItemSubtitle, VListItemTitle, VListSubheader } from 'vuetify/components/VList';
import { VMenu } from 'vuetify/components/VMenu';
import { VOverlay } from 'vuetify/components/VOverlay';
import { VPagination } from 'vuetify/components/VPagination';
import { VProgressCircular } from 'vuetify/components/VProgressCircular';
import { VProgressLinear } from 'vuetify/components/VProgressLinear';
import { VSelect } from 'vuetify/components/VSelect';
import { VSheet } from 'vuetify/components/VSheet';
import { VSnackbar } from 'vuetify/components/VSnackbar';
import { VTabs } from 'vuetify/components/VTabs';
import { VTable } from 'vuetify/components/VTable';
import { VTextField } from 'vuetify/components/VTextField';
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg';
import 'vuetify/styles';

import 'line-awesome/dist/line-awesome/css/line-awesome.css';

import { PerfectScrollbar } from 'vue3-perfect-scrollbar';
import 'vue3-perfect-scrollbar/dist/vue3-perfect-scrollbar.min.css';

import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

import router from '@/router/desktop.js';

import version from '@/lib/version.js';
import settings from '@/lib/settings.js';
import userstate from '@/lib/userstate.js';
import {
    getI18nOptions,
    translateIf,
    translateError,
    i18nFunctions
} from '@/lib/i18n.js';

import AmountInput from '@/components/desktop/AmountInput.vue';

import '@/styles/desktop/base.css';
import '@/styles/desktop/layout.css';
import '@/styles/desktop/font-size.css';
import '@/styles/desktop/gap-size.css';
import '@/styles/desktop/vuetify.css';
import '@/styles/desktop/classess.css';

import App from './DesktopApp.vue';

const app = createApp(App);
const pinia = createPinia();
const i18n = createI18n(getI18nOptions());
const vuetify = createVuetify({
    components: {
        VApp,
        VAvatar,
        VBtn,
        VCard,
        VCardActions,
        VCardItem,
        VCardSubtitle,
        VCardText,
        VCardTitle,
        VChip,
        VDialog,
        VDivider,
        VForm,
        VContainer,
        VCol,
        VRow,
        VSpacer,
        VIcon,
        VImg,
        VInput,
        VList,
        VListGroup,
        VListImg,
        VListItem,
        VListItemAction,
        VListItemMedia,
        VListItemSubtitle,
        VListItemTitle,
        VListSubheader,
        VMenu,
        VOverlay,
        VPagination,
        VProgressCircular,
        VProgressLinear,
        VSelect,
        VSheet,
        VSnackbar,
        VTabs,
        VTable,
        VTextField
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
            activeColor: 'primary'
        },
        VRadio: {
            color: 'primary',
            hideDetails: 'auto'
        },
        VSelect: {
            variant: 'outlined',
            color: 'primary',
            hideDetails: 'auto'
        },
        VSlider: {
            color: 'primary',
            hideDetails: 'auto'
        },
        VSwitch: {
            color: 'primary',
            hideDetails: 'auto'
        },
        VProgressCircular: {
            size: 40
        },
        VSnackbar: {
            timeout: 3000
        },
        VTabs: {
            color: 'primary',
            VSlideGroup: {
                showArrows: true
            }
        },
        VTextarea: {
            variant: 'outlined',
            color: 'primary',
            hideDetails: 'auto'
        },
        VTextField: {
            variant: 'outlined',
            color: 'primary',
            hideDetails: 'auto'
        },
        VTooltip: {
            location: 'top'
        }
    },
    theme: {
        defaultTheme: 'light',
        themes: {
            light: {
                dark: false,
                colors: {
                    'primary': '#c67e48',
                    'secondary': '#8a8d93',
                    'on-secondary': '#fff',
                    'success': '#4cd964',
                    'info': '#2196f3',
                    'warning': '#ff9500',
                    'error': '#ff3b30',
                    'on-primary': '#ffffff',
                    'on-success': '#ffffff',
                    'on-warning': '#ffffff',
                    'background': '#faf8f4',
                    'on-background': '#413935',
                    'on-surface': '#413935',
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
                    'skin-bordered-surface': '#fff'
                },
                variables: {
                    'btn-height': '38px',
                    'code-color': '#ff8000',
                    'overlay-scrim-background': '#3a3541',
                    'overlay-scrim-opacity': 0.5,
                    'hover-opacity': 0.04,
                    'focus-opacity': 0.1,
                    'selected-opacity': 0.12,
                    'activated-opacity': 0.1,
                    'pressed-opacity': 0.14,
                    'dragged-opacity': 0.1,
                    'border-color': '#3a3541',
                    'table-header-background': '#f9fafc',
                    'custom-background': '#f9f8f9',
                    'shadow-key-umbra-opacity': 'rgba(var(--v-theme-on-surface), 0.08)',
                    'shadow-key-penumbra-opacity': 'rgba(var(--v-theme-on-surface), 0.12)',
                    'shadow-key-ambient-opacity': 'rgba(var(--v-theme-on-surface), 0.04)'
                }
            },
            dark: {
                dark: true,
                colors: {
                    'primary': '#c67e48',
                    'secondary': '#8a8d93',
                    'on-secondary': '#fff',
                    'success': '#4cd964',
                    'info': '#2196f3',
                    'warning': '#ff9500',
                    'error': '#ff3b30',
                    'on-primary': '#ffffff',
                    'on-success': '#ffffff',
                    'on-warning': '#ffffff',
                    'background': '#000000',
                    'on-background': '#fcf0e3',
                    'surface': '#1c1c1d',
                    'on-surface': '#fcf0e3',
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
                    'skin-bordered-surface': '#312d4b'
                },
                variables: {
                    'btn-height': '38px',
                    'code-color': '#ff8000',
                    'overlay-scrim-background': '#2C2942',
                    'overlay-scrim-opacity': 0.6,
                    'hover-opacity': 0.04,
                    'focus-opacity': 0.1,
                    'selected-opacity': 0.12,
                    'activated-opacity': 0.1,
                    'pressed-opacity': 0.14,
                    'dragged-opacity': 0.1,
                    'border-color': '#E7E3FC',
                    'table-header-background': '#3D3759',
                    'custom-background': '#373452',
                    'shadow-key-umbra-opacity': 'rgba(20, 18, 33, 0.08)',
                    'shadow-key-penumbra-opacity': 'rgba(20, 18, 33, 0.12)',
                    'shadow-key-ambient-opacity': 'rgba(20, 18, 33, 0.04)'
                }
            }
        }
    }
});

app.use(pinia);
app.use(i18n);
app.use(vuetify);
app.use(router);

app.component('PerfectScrollbar', PerfectScrollbar);
app.component('VueDatePicker', VueDatePicker);

app.component('AmountInput', AmountInput);

app.config.globalProperties.$version = version.getVersion();
app.config.globalProperties.$buildTime = version.getBuildTime();

app.config.globalProperties.$settings = settings;
app.config.globalProperties.$locale = i18nFunctions(i18n.global);
app.config.globalProperties.$tIf = (text, isTranslate) => translateIf(text, isTranslate, i18n.global.t);
app.config.globalProperties.$tError = (message) => translateError(message, i18n.global.t);

app.config.globalProperties.$user = userstate;

app.mount('#app');
