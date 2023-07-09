<template>
    <div class="layout-wrapper layout-nav-type-vertical layout-navbar-static layout-footer-static layout-content-width-fluid"
         :class="{ 'layout-overlay-nav': mdAndDown }">
        <div class="layout-vertical-nav" :class="{'visible': showVerticalOverlayMenu, 'scrolled': isVerticalNavScrolled, 'overlay-nav': mdAndDown}">
            <div class="nav-header">
                <router-link to="/" class="app-logo d-flex align-center gap-x-3 app-title-wrapper">
                    <div class="d-flex">
                        <img alt="logo" class="main-logo" src="/img/ezbookkeeping-192.png" />
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
                    <router-link to="/">
                        <v-icon class="nav-item-icon" :icon="icons.overview"/>
                        <span class="nav-item-title">{{ $t('Overview') }}</span>
                    </router-link>
                </li>
                <li class="nav-section-title">
                    <div class="title-wrapper">
                        <span class="title-text">{{ $t('Transaction Data') }}</span>
                    </div>
                </li>
                <li class="nav-link">
                    <router-link to="/transactions">
                        <v-icon class="nav-item-icon" :icon="icons.transactions"/>
                        <span class="nav-item-title">{{ $t('Transaction List') }}</span>
                    </router-link>
                </li>
                <li class="nav-link">
                    <router-link to="/statistics/transaction">
                        <v-icon class="nav-item-icon" :icon="icons.statistics"/>
                        <span class="nav-item-title">{{ $t('Statistics Data') }}</span>
                    </router-link>
                </li>
                <li class="nav-section-title">
                    <div class="title-wrapper">
                        <span class="title-text">{{ $t('Data Management') }}</span>
                    </div>
                </li>
                <li class="nav-link">
                    <router-link to="/accounts">
                        <v-icon class="nav-item-icon" :icon="icons.accounts"/>
                        <span class="nav-item-title">{{ $t('Account List') }}</span>
                    </router-link>
                </li>
                <li class="nav-link">
                    <router-link to="/categories">
                        <v-icon class="nav-item-icon" :icon="icons.categories"/>
                        <span class="nav-item-title">{{ $t('Transaction Categories') }}</span>
                    </router-link>
                </li>
                <li class="nav-link">
                    <router-link to="/tags">
                        <v-icon class="nav-item-icon" :icon="icons.tags"/>
                        <span class="nav-item-title">{{ $t('Transaction Tags') }}</span>
                    </router-link>
                </li>
                <li class="nav-section-title">
                    <div class="title-wrapper">
                        <span class="title-text">{{ $t('Other') }}</span>
                    </div>
                </li>
                <li class="nav-link">
                    <router-link to="/exchange_rates">
                        <v-icon class="nav-item-icon" :icon="icons.exchangeRates"/>
                        <span class="nav-item-title">{{ $t('Exchange Rates Data') }}</span>
                    </router-link>
                </li>
                <li class="nav-link">
                    <a href="javascript:void(0);" @click="showMobileQrCode = true">
                        <v-icon class="nav-item-icon" :icon="icons.mobile"/>
                        <span class="nav-item-title">{{ $t('Use on Mobile Device') }}</span>
                    </a>
                </li>
                <li class="nav-link">
                    <router-link to="/about">
                        <v-icon class="nav-item-icon" :icon="icons.about"/>
                        <span class="nav-item-title">{{ $t('About') }}</span>
                    </router-link>
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
                                <img alt="logo" class="main-logo" src="/img/ezbookkeeping-192.png" />
                            </div>
                            <h1 class="font-weight-medium text-xl">{{ $t('global.app.title') }}</h1>
                        </div>
                        <v-spacer />
                        <v-btn color="primary" variant="text" class="me-2"
                               :icon="true" @click="(theme === 'light' ? theme = 'dark' : (theme === 'dark' ? theme = 'auto' : theme = 'light'))">
                            <v-icon :icon="(theme === 'light' ? icons.themeLight : (theme === 'dark' ? icons.themeDark : icons.themeAuto))" size="24" />
                        </v-btn>
                        <v-avatar class="cursor-pointer" color="primary" variant="tonal">
                            <v-img :src="currentUserAvatar" v-if="currentUserAvatar">
                                <template v-slot:placeholder>
                                    <div class="d-flex align-center justify-center fill-height">
                                        <v-icon :icon="icons.user"/>
                                    </div>
                                </template>
                            </v-img>
                            <v-icon :icon="icons.user" v-else-if="!currentUserAvatar"/>
                            <v-menu activator="parent" width="230" location="bottom end" offset="14px">
                                <v-list>
                                    <v-list-item>
                                        <template #prepend>
                                            <v-list-item-action start>
                                                <v-avatar color="primary" variant="tonal">
                                                    <v-img :src="currentUserAvatar" v-if="currentUserAvatar">
                                                        <template v-slot:placeholder>
                                                            <div class="d-flex align-center justify-center fill-height">
                                                                <v-icon :icon="icons.user"/>
                                                            </div>
                                                        </template>
                                                    </v-img>
                                                    <v-icon :icon="icons.user" v-else-if="!currentUserAvatar"/>
                                                </v-avatar>
                                            </v-list-item-action>
                                        </template>
                                        <v-list-item-title class="font-weight-semibold">
                                            {{ currentNickName }}
                                        </v-list-item-title>
                                    </v-list-item>
                                    <v-divider class="my-2"/>
                                    <v-list-item :prepend-icon="icons.profile"
                                                 :title="$t('User Settings')"
                                                 to="/user/settings"></v-list-item>
                                    <v-list-item :prepend-icon="icons.settings"
                                                 :title="$t('Application Settings')"
                                                 to="/app/settings"></v-list-item>
                                    <v-divider class="my-2"/>
                                    <v-list-item :disabled="logouting"
                                                 :prepend-icon="icons.logout"
                                                 :title="$t('Log Out')"
                                                 @click="logout"></v-list-item>
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

        <v-dialog width="400" v-model="showMobileQrCode">
            <v-card>
                <v-toolbar color="primary">
                    <v-toolbar-title>{{ $t('global.app.title') }}</v-toolbar-title>
                </v-toolbar>
                <v-card-text class="pa-4">
                    <p>{{ $t('You can scan the below QR code on your mobile device.') }}</p>
                </v-card-text>
                <v-card-text class="pa-4 w-100 d-flex justify-center">
                    <img alt="qrcode" class="img-url-qrcode" :src="mobileUrlQrCodePath" />
                </v-card-text>
                <v-card-actions>
                    <v-btn :href="mobileVersionPath">{{$t('Switch to Mobile Version') }}</v-btn>
                    <v-spacer></v-spacer>
                    <v-btn @click="showMobileQrCode = false">{{ $t('Close') }}</v-btn>
                </v-card-actions>
            </v-card>
        </v-dialog>

        <v-overlay class="justify-center align-center" :persistent="true" v-model="showLoading">
            <v-progress-circular indeterminate></v-progress-circular>
        </v-overlay>

        <snack-bar ref="snackbar" />
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
import { getMobileUrlQrCodePath } from '@/lib/qrcode.js';
import { getMobileVersionPath } from '@/lib/version.js';

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
    mdiCellphone,
    mdiInformationOutline,
    mdiThemeLightDark,
    mdiWeatherSunny,
    mdiWeatherNight,
    mdiAccount,
    mdiAccountCogOutline,
    mdiLogout
} from '@mdi/js';

export default {
    data() {
        return {
            logouting: false,
            isVerticalNavScrolled: false,
            showVerticalOverlayMenu: false,
            showLoading: false,
            showMobileQrCode: false,
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
                mobile: mdiCellphone,
                about: mdiInformationOutline,
                themeAuto: mdiThemeLightDark,
                themeLight: mdiWeatherSunny,
                themeDark: mdiWeatherNight,
                user: mdiAccount,
                profile: mdiAccountCogOutline,
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
        mobileUrlQrCodePath() {
            return getMobileUrlQrCodePath();
        },
        mobileVersionPath() {
            return getMobileVersionPath();
        },
        currentNickName() {
            return this.userStore.currentUserNickname || this.$t('User');
        },
        currentUserAvatar() {
            return this.userStore.currentUserAvatar;
        },
        theme: {
            get: function () {
                return this.settingsStore.appSettings.theme;
            },
            set: function (value) {
                if (value !== this.settingsStore.appSettings.theme) {
                    this.settingsStore.setTheme(value);

                    if (value === 'light' || value === 'dark') {
                        this.globalTheme.global.name.value = value;
                    } else {
                        this.globalTheme.global.name.value = getSystemTheme();
                    }
                }
            }
        }
    },
    setup() {
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

                self.settingsStore.clearAppSettings();

                const localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
                self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                this.$router.replace('/login');
            }).catch(error => {
                self.logouting = false;
                self.showLoading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        }
    }
}
</script>

<style>
.img-url-qrcode {
    width: 320px;
    height: 320px
}

.main-logo {
    width: 1.8em;
    height: 1.8em;
}
</style>
