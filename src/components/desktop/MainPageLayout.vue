<template>
    <div class="layout-wrapper layout-nav-type-vertical"
         :class="{ 'layout-overlay-nav': lgAndDown }">
        <div class="layout-vertical-nav" :class="{'visible': showVerticalOverlayMenu, 'overlay-nav': lgAndDown}"
             v-if="!noNavbar">
            <div class="nav-header">
                <router-link to="/" class="app-logo d-flex align-center gap-x-3">
                    <div class="d-flex">
                        <img alt="logo" class="main-logo" :src="APPLICATION_LOGO_PATH" />
                    </div>
                    <h1 class="app-title">{{ tt('global.app.title') }}</h1>
                </router-link>
            </div>
            <ul class="nav-items" :class="navItemsClass">
                <slot name="nav-items" />
            </ul>
        </div>

        <div :class="{ 'layout-content-wrapper': !noNavbar }">
            <div class="layout-navbar navbar-blur">
                <div class="navbar-content-container">
                    <div class="navbar-content d-flex h-100 align-center">
                        <v-btn class="ms-n2 d-lg-none" color="default" variant="text"
                               :aria-label="tt('Open Menu')"
                               :icon="true" @click="showVerticalOverlayMenu = true" v-if="!noNavbar">
                            <v-icon :icon="mdiMenu" size="24" />
                        </v-btn>
                        <div class="app-logo d-flex align-center gap-x-3 me-3" v-if="noNavbar">
                            <div class="d-flex">
                                <img alt="logo" class="main-logo" :src="APPLICATION_LOGO_PATH" />
                            </div>
                            <h1 class="app-title d-none d-md-inline-flex">{{ tt('global.app.title') }}</h1>
                        </div>
                        <div class="app-top-toolbar d-inline-flex"
                             :class="{ 'app-top-toolbar-without-navbar': noNavbar }">
                            <slot name="top-toolbar">
                                <router-link to="/">
                                    <v-btn class="top-navigation-button" density="comfortable" variant="text" :icon="true"
                                           :aria-label="tt('Overview')"
                                           :active="isTopNavigationActive('/')"
                                           :color="isTopNavigationActive('/') ? 'primary' : 'default'">
                                        <v-icon :icon="isTopNavigationActive('/') ? mdiHome : mdiHomeOutline" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Overview') }}</v-tooltip>
                                    </v-btn>
                                </router-link>

                                <router-link to="/transaction/list?pageType=0&dateType=7">
                                    <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text" :icon="true"
                                           :aria-label="tt('Transaction Details')"
                                           :active="isTopNavigationActive('/transaction/list')"
                                           :color="isTopNavigationActive('/transaction/list') ? 'primary' : 'default'">
                                        <v-icon :icon="isTopNavigationActive('/transaction/list') ? mdiListBox : mdiListBoxOutline" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Transaction Details') }}</v-tooltip>
                                    </v-btn>
                                </router-link>

                                <router-link to="/statistics/transaction">
                                    <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text" :icon="true"
                                           :aria-label="tt('Statistics & Analysis')"
                                           :active="isTopNavigationActive('/statistics/transaction')"
                                           :color="isTopNavigationActive('/statistics/transaction') ? 'primary' : 'default'">
                                        <v-icon :icon="isTopNavigationActive('/statistics/transaction') ? mdiChartPie : mdiChartPieOutline" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Statistics & Analysis') }}</v-tooltip>
                                    </v-btn>
                                </router-link>

                                <router-link to="/insights/explorer">
                                    <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text" :icon="true"
                                           :aria-label="tt('Insights Explorer')"
                                           :active="isTopNavigationActive('/insights/explorer')"
                                           :color="isTopNavigationActive('/insights/explorer') ? 'primary' : 'default'">
                                        <v-icon :icon="isTopNavigationActive('/insights/explorer') ? mdiCompass : mdiCompassOutline" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Insights Explorer') }}</v-tooltip>
                                    </v-btn>
                                </router-link>

                                <router-link to="/account/list">
                                    <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text" :icon="true"
                                           :aria-label="tt('Accounts')"
                                           :active="isTopNavigationActive('/account/list')"
                                           :color="isTopNavigationActive('/account/list') ? 'primary' : 'default'">
                                        <v-icon :icon="isTopNavigationActive('/account/list') ? mdiCreditCard : mdiCreditCardOutline" size="24" />
                                        <v-tooltip activator="parent">{{ tt('Accounts') }}</v-tooltip>
                                    </v-btn>
                                </router-link>

                                <router-link class="d-inline-flex align-center" to="/transaction/list?pageType=0&dateType=7"
                                             v-if="showAddTransactionButtonInDesktopNavbar">
                                    <v-btn class="add-transaction-button ms-2" color="primary" density="comfortable" variant="flat"
                                           :aria-label="tt('Add Transaction')" :icon="true"
                                           @click="showAddDialogInTransactionListPage">
                                        <v-icon :icon="mdiPlus" size="22" />
                                        <v-tooltip activator="parent">{{ tt('Add Transaction') }}</v-tooltip>
                                    </v-btn>
                                </router-link>
                            </slot>
                        </div>
                        <v-spacer />
                        <v-btn class="ms-2" color="primary" variant="text" density="comfortable"
                               :aria-label="tt('Use on Mobile Device')"
                               :icon="true" @click="showMobileQrCode = true">
                            <v-icon :icon="mdiCellphone" size="24" />
                            <v-tooltip activator="parent">{{ tt('Use on Mobile Device') }}</v-tooltip>
                        </v-btn>
                        <v-btn class="ms-2" color="primary" variant="text" density="comfortable"
                               :aria-label="tt('Theme')"
                               :icon="true" @click="(currentTheme === 'light' ? currentTheme = 'dark' : (currentTheme === 'dark' ? currentTheme = 'auto' : currentTheme = 'light'))">
                            <v-icon :icon="(currentTheme === 'light' ? mdiWeatherSunny : (currentTheme === 'dark' ? mdiWeatherNight : mdiThemeLightDark))" size="24" />
                        </v-btn>
                        <v-avatar class="cursor-pointer ms-3" variant="tonal"
                                  :color="currentUserAvatar ? 'rgba(0,0,0,0)' : 'primary'">
                            <v-icon :color="currentUserAvatar ? 'primary' : undefined" :icon="mdiAccount"/>
                            <span class="user-avatar-image" :style="currentUserAvatarStyle" v-if="currentUserAvatar"></span>
                            <v-menu activator="parent" width="230" location="bottom end" offset="14px">
                                <v-list>
                                    <v-list-item>
                                        <template #prepend>
                                            <v-list-item-action>
                                                <v-avatar variant="tonal"
                                                          :color="currentUserAvatar ? 'rgba(0,0,0,0)' : 'primary'">
                                                    <v-img :src="currentUserAvatar" v-if="currentUserAvatar">
                                                        <template #placeholder>
                                                            <div class="d-flex align-center justify-center bg-light-primary">
                                                                <v-icon color="primary" :icon="mdiAccount"/>
                                                            </div>
                                                        </template>
                                                    </v-img>
                                                    <v-icon :icon="mdiAccount" v-else-if="!currentUserAvatar"/>
                                                </v-avatar>
                                            </v-list-item-action>
                                        </template>
                                        <v-list-item-title class="ms-1">
                                            {{ currentNickName }}
                                        </v-list-item-title>
                                    </v-list-item>
                                    <v-divider class="my-1"/>
                                    <v-list-item :prepend-icon="mdiAccountCogOutline"
                                                 :title="tt('User Settings')"
                                                 to="/user/settings/basic"></v-list-item>
                                    <v-list-item :prepend-icon="mdiCogOutline"
                                                 :title="tt('Application Settings')"
                                                 to="/app/settings/basic"></v-list-item>
                                    <v-divider class="my-1"/>
                                    <v-list-item :prepend-icon="mdiInformationOutline"
                                                 :title="tt('About')"
                                                 @click="showAboutDialog = true"></v-list-item>
                                    <v-divider class="my-1"/>
                                    <v-list-item :prepend-icon="mdiLockOutline"
                                                 :title="tt('Lock Application')"
                                                 v-if="isEnableApplicationLock"
                                                 @click="lock"></v-list-item>
                                    <v-list-item :disabled="logouting"
                                                 :prepend-icon="mdiLogout"
                                                 :title="tt('Log Out')"
                                                 @click="logout"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-avatar>
                    </div>
                </div>
            </div>
            <div class="layout-page-content">
                <div class="page-content-container">
                    <slot name="content" />
                </div>
            </div>
        </div>

        <div class="layout-overlay" :class="{ 'visible': showVerticalOverlayMenu }" @click="showVerticalOverlayMenu = false"></div>

        <v-overlay class="justify-center align-center" :persistent="true" v-model="showLoading">
            <v-progress-circular indeterminate></v-progress-circular>
        </v-overlay>

        <switch-to-mobile-dialog v-model:show="showMobileQrCode" />
        <about-dialog v-model:show="showAboutDialog" />

        <snack-bar ref="snackbar" />
    </div>
</template>

<script setup lang="ts">
import AboutDialog from '@/views/desktop/common/dialogs/AboutDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useDisplay, useTheme } from 'vuetify';
import { useRoute, useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';

import { getSystemTheme, setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

import {
    mdiMenu,
    mdiHome,
    mdiHomeOutline,
    mdiListBox,
    mdiListBoxOutline,
    mdiCreditCard,
    mdiCreditCardOutline,
    mdiChartPie,
    mdiChartPieOutline,
    mdiCompass,
    mdiCompassOutline,
    mdiPlus,
    mdiCellphone,
    mdiThemeLightDark,
    mdiWeatherSunny,
    mdiWeatherNight,
    mdiAccount,
    mdiAccountCogOutline,
    mdiCogOutline,
    mdiInformationOutline,
    mdiLockOutline,
    mdiLogout
} from '@mdi/js';

defineProps<{
    navItemsClass?: string;
    noNavbar?: boolean;
}>();

type SnackBarType = InstanceType<typeof SnackBar>;

const { lgAndDown } = useDisplay();
const theme = useTheme();
const route = useRoute();
const router = useRouter();

const { tt, initLocale } = useI18n();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const desktopPageStore = useDesktopPageStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const logouting = ref<boolean>(false);
const showVerticalOverlayMenu = ref<boolean>(false);
const showLoading = ref<boolean>(false);
const showMobileQrCode = ref<boolean>(false);
const showAboutDialog = ref<boolean>(false);

const currentNickName = computed<string>(() => userStore.currentUserNickname || tt('User'));
const currentUserAvatar = computed<string | null>(() => userStore.currentUserAvatar);
const currentUserAvatarStyle = computed<Record<string, string> | undefined>(() => currentUserAvatar.value ? {
    backgroundImage: `url("${currentUserAvatar.value}")`
} : undefined);

const currentTheme = computed<string>({
    get: () => {
        return settingsStore.appSettings.theme;
    },
    set: (value: string) => {
        if (value !== settingsStore.appSettings.theme) {
            settingsStore.setTheme(value);

            if (value === ThemeType.Light || value === ThemeType.Dark) {
                theme.change(value);
            } else {
                theme.change(getSystemTheme());
            }
        }
    }
});

const showAddTransactionButtonInDesktopNavbar = computed<boolean>(() => settingsStore.appSettings.showAddTransactionButtonInDesktopNavbar);
const isEnableApplicationLock = computed<boolean>(() => settingsStore.appSettings.applicationLock);

function isTopNavigationActive(path: string): boolean {
    return route.path === path;
}

function lock(): void {
    rootStore.lock();
    router.replace('/unlock');
}

function logout(): void {
    logouting.value = true;
    showLoading.value = true;

    rootStore.logout().then(() => {
        logouting.value = false;
        showLoading.value = false;

        settingsStore.clearAppSettings();

        const localeDefaultSettings = initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone);
        settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

        setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);

        router.replace('/login');
    }).catch(error => {
        logouting.value = false;
        showLoading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function showAddDialogInTransactionListPage(): void {
    desktopPageStore.setShowAddTransactionDialogInTransactionList();
}
</script>

<style>
.top-navigation-button.v-btn {
    width: 40px;
    min-width: 40px;
    border-radius: 10px;
}

.top-navigation-button.v-btn:not(.v-btn--active) {
    color: rgba(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity)) !important;
}

.add-transaction-button.v-btn {
    --v-btn-height: 32px;
    width: 38px;
    min-width: 38px;
    border-radius: 10px;
}

.user-avatar-image {
    position: absolute;
    z-index: 1;
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
    inset: 0;
    pointer-events: none;
}
</style>
