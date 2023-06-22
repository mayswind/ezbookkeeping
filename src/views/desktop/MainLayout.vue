<template>
    <div class="layout-wrapper layout-nav-type-vertical layout-navbar-static layout-footer-static layout-content-width-fluid"
         :class="{ 'layout-overlay-nav': mdAndDown }">
        <div class="layout-vertical-nav" :class="{'visible': showVerticalOverlayMenu, 'scrolled': isVerticalNavScrolled, 'overlay-nav': mdAndDown}">
            <div class="nav-header">
                <router-link to="/" class="app-logo d-flex align-center gap-x-3 app-title-wrapper">
                    <div class="d-flex">
                        <v-img alt="logo" class="main-logo" src="/img/ezbookkeeping-192.png" />
                    </div>
                    <h1 class="font-weight-medium text-xl">{{ $t('global.app.title') }}</h1>
                </router-link>
            </div>
            <perfect-scrollbar
                tag="ul" class="nav-items"
                :options="{ wheelPropagation: false }"
                @ps-scroll-y="handleNavScroll"
            >
                <li class="nav-link">
                    <RouterLink to="/">
                        <v-icon class="nav-item-icon" :icon="icons.overview"/>
                        <span class="nav-item-title">{{ $t('Overview') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-section-title">
                    <div class="title-wrapper">
                        <span class="title-text">{{ $t('Transaction Data') }}</span>
                    </div>
                </li>
                <li class="nav-link">
                    <RouterLink to="/transactions">
                        <v-icon class="nav-item-icon" :icon="icons.transactions"/>
                        <span class="nav-item-title">{{ $t('Transaction List') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-link">
                    <RouterLink to="/statistics/transaction">
                        <v-icon class="nav-item-icon" :icon="icons.statistics"/>
                        <span class="nav-item-title">{{ $t('Statistics Data') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-section-title">
                    <div class="title-wrapper">
                        <span class="title-text">{{ $t('Data Management') }}</span>
                    </div>
                </li>
                <li class="nav-link">
                    <RouterLink to="/accounts">
                        <v-icon class="nav-item-icon" :icon="icons.accounts"/>
                        <span class="nav-item-title">{{ $t('Account List') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-link">
                    <RouterLink to="/categories">
                        <v-icon class="nav-item-icon" :icon="icons.categories"/>
                        <span class="nav-item-title">{{ $t('Transaction Categories') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-link">
                    <RouterLink to="/tags">
                        <v-icon class="nav-item-icon" :icon="icons.tags"/>
                        <span class="nav-item-title">{{ $t('Transaction Tags') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-section-title">
                    <div class="title-wrapper">
                        <span class="title-text">{{ $t('Other') }}</span>
                    </div>
                </li>
                <li class="nav-link">
                    <RouterLink to="/exchange_rates">
                        <v-icon class="nav-item-icon" :icon="icons.exchangeRates"/>
                        <span class="nav-item-title">{{ $t('Exchange Rates Data') }}</span>
                    </RouterLink>
                </li>
                <li class="nav-link">
                    <RouterLink to="/about">
                        <v-icon class="nav-item-icon" :icon="icons.about"/>
                        <span class="nav-item-title">{{ $t('About') }}</span>
                    </RouterLink>
                </li>
            </perfect-scrollbar>
        </div>

        <div class="layout-content-wrapper">
            <div class="layout-navbar navbar-blur">
                <div class="navbar-content-container">
                    <div class="d-flex h-100 align-center">
                        <v-btn class="ms-n3 mr-2 d-lg-none" color="default" variant="text"
                               :icon="true" @click="showVerticalOverlayMenu = true">
                            <v-icon :icon="icons.menu" size="24" />
                        </v-btn>
                        <div class="app-logo d-flex align-center gap-x-3 app-title-wrapper" v-if="mdAndDown">
                            <div class="d-flex">
                                <v-img alt="logo" class="main-logo" src="/img/ezbookkeeping-192.png" />
                            </div>
                            <h1 class="font-weight-medium text-xl">{{ $t('global.app.title') }}</h1>
                        </div>
                        <v-spacer />
                        <v-btn color="primary" variant="text" class="me-2"
                               :icon="true" @click="(currentTheme === 'light' ? currentTheme = 'dark' : (currentTheme === 'dark' ? currentTheme = 'auto' : currentTheme = 'light'))">
                            <v-icon :icon="(currentTheme === 'light' ? icons.themeLight : (currentTheme === 'dark' ? icons.themeDark : icons.themeAuto))" size="24" />
                        </v-btn>
                        <v-avatar class="cursor-pointer" color="primary" variant="tonal">
                            <v-icon :icon="icons.user"/>
                            <v-menu activator="parent" width="230" location="bottom end" offset="14px">
                                <v-list>
                                    <v-list-item>
                                        <template #prepend>
                                            <v-list-item-action start>
                                                <v-avatar color="primary" variant="tonal">
                                                    <v-icon :icon="icons.user"/>
                                                </v-avatar>
                                            </v-list-item-action>
                                        </template>
                                        <v-list-item-title class="font-weight-semibold">
                                            {{ currentNickName }}
                                        </v-list-item-title>
                                    </v-list-item>
                                    <v-divider class="my-2"/>
                                    <v-list-item to="/user/settings">
                                        <template #prepend>
                                            <v-icon class="me-2" :icon="icons.profile" size="22"/>
                                        </template>
                                        <v-list-item-title>{{ $t('User Profile') }}</v-list-item-title>
                                    </v-list-item>
                                    <v-list-item to="/app/settings">
                                        <template #prepend>
                                            <v-icon class="me-2" :icon="icons.settings" size="22"/>
                                        </template>
                                        <v-list-item-title>{{ $t('Application Settings') }}</v-list-item-title>
                                    </v-list-item>
                                    <v-divider class="my-2"/>
                                    <v-list-item :class="{ 'disabled': logouting }" @click="logout">
                                        <template #prepend>
                                            <v-icon class="me-2" :icon="icons.logout" size="22"/>
                                        </template>
                                        <v-list-item-title>{{ $t('Log Out') }}</v-list-item-title>
                                    </v-list-item>
                                </v-list>
                            </v-menu>
                        </v-avatar>
                    </div>
                </div>
            </div>
            <div class="layout-page-content">
                <div class="page-content-container">
                    <router-view/>
                </div>
            </div>
        </div>

        <div class="layout-overlay" :class="{ 'visible': showVerticalOverlayMenu }" @click="showVerticalOverlayMenu = false"></div>

        <v-overlay class="justify-center align-center" :persistent="true" v-model="showLoading">
            <v-progress-circular indeterminate></v-progress-circular>
        </v-overlay>

        <v-snackbar v-model="showSnackbar">
            {{ snackbarMessage }}

            <template #actions>
                <v-btn color="red" variant="text" @click="showSnackbar = false">{{ $t('Close') }}</v-btn>
            </template>
        </v-snackbar>
    </div>
</template>

<script>
import { useDisplay } from 'vuetify';
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import { getSystemTheme } from '@/lib/ui.js';

import {
    mdiMenu,
    mdiHomeOutline,
    mdiListBoxOutline,
    mdiCreditCardOutline,
    mdiViewDashboardOutline,
    mdiTagOutline,
    mdiChartPieOutline,
    mdiSwapHorizontal,
    mdiCogOutline,
    mdiInformationOutline,
    mdiThemeLightDark,
    mdiWeatherSunny,
    mdiWeatherNight,
    mdiAccount,
    mdiAccountOutline,
    mdiLogout
} from '@mdi/js';

export default {
    data() {
        const self = this;

        return {
            theme: self.$settings.getTheme(),
            logouting: false,
            isVerticalNavScrolled: false,
            showVerticalOverlayMenu: false,
            showLoading: false,
            showSnackbar: false,
            snackbarMessage: '',
            icons: {
                menu: mdiMenu,
                overview: mdiHomeOutline,
                transactions: mdiListBoxOutline,
                accounts: mdiCreditCardOutline,
                categories: mdiViewDashboardOutline,
                tags: mdiTagOutline,
                statistics: mdiChartPieOutline,
                exchangeRates: mdiSwapHorizontal,
                settings: mdiCogOutline,
                about: mdiInformationOutline,
                themeAuto: mdiThemeLightDark,
                themeLight: mdiWeatherSunny,
                themeDark: mdiWeatherNight,
                user: mdiAccount,
                profile: mdiAccountOutline,
                logout: mdiLogout
            }
        }
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore),
        mdAndDown() {
            const { mdAndDown } = useDisplay();
            return mdAndDown.value;
        },
        currentNickName() {
            return this.userStore.currentUserNickname || this.$t('User');
        },
        currentTheme: {
            get: function () {
                return this.theme;
            },
            set: function (value) {
                if (value !== this.$settings.getTheme()) {
                    this.theme = value;
                    this.$settings.setTheme(value);

                    if (value === 'light' || value === 'dark') {
                        this.globalTheme.global.name.value = value;
                    } else {
                        this.globalTheme.global.name.value = getSystemTheme();
                    }
                }
            }
        }
    },
    setup () {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        handleNavScroll(e) {
            this.isVerticalNavScrolled = e.target.scrollTop > 0;
        },
        logout() {
            const self = this;

            self.logouting = true;
            self.showLoading = true;

            self.rootStore.logout().then(() => {
                self.logouting = false;
                self.showLoading = false;

                self.$settings.clearSettings();

                const localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage);
                self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                this.$router.replace('/login');
            }).catch(error => {
                self.logouting = false;
                self.showLoading = false;

                if (!error.processed) {
                    self.showSnackbarMessage(self.$tError(error.message || error));
                }
            });
        },
        showSnackbarMessage(message) {
            this.showSnackbar = true;
            this.snackbarMessage = message;
        }
    }
}
</script>

<style>
.main-logo {
    width: 1.8em;
    height: 1.8em;
}
</style>
