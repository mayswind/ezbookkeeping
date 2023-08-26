<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="ezBookkeepingLogoPath" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ $t('global.app.title') }}</h1>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer" v-if="!isDarkMode">
                    <v-img src="img/desktop/background.svg"/>
                </div>
                <div class="d-flex auth-img-footer" v-if="isDarkMode">
                    <v-img src="img/desktop/background-dark.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img max-width="600px" src="img/desktop/people4.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h5 class="text-h5 mb-3">{{ $t('Forget Password?') }}</h5>
                            <p class="mb-0">{{ $t('Please input your email address used for registration and we\'ll send you an email with reset password link') }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12">
                                        <v-text-field
                                            type="email"
                                            autocomplete="email"
                                            autofocus="autofocus"
                                            clearable
                                            :disabled="requesting"
                                            :label="$t('E-mail')"
                                            :placeholder="$t('Your email address')"
                                            v-model="email"
                                            @keyup.enter="requestResetPassword"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block type="submit" :disabled="!email || requesting" @click="requestResetPassword">
                                            {{ $t('Send Reset Link') }}
                                            <v-progress-circular indeterminate size="24" class="ml-2" v-if="requesting"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12">
                                        <router-link class="d-flex align-center justify-center" to="/login"
                                                     :class="{ 'disabled': requesting }">
                                            <v-icon :icon="icons.left"/>
                                            <span>{{ $t('Back to log in') }}</span>
                                        </router-link>
                                    </v-col>
                                </v-row>
                            </v-form>
                        </v-card-text>
                    </v-card>
                </div>
                <v-spacer/>
                <div class="d-flex align-center justify-center">
                    <v-card variant="flat" class="w-100 px-4 pb-4" max-width="500">
                        <v-card-text class="pt-0">
                            <v-row>
                                <v-col cols="12" class="text-center">
                                    <v-menu location="bottom">
                                        <template #activator="{ props }">
                                            <v-btn variant="text"
                                                   :disabled="requesting"
                                                   v-bind="props">{{ currentLanguageName }}</v-btn>
                                        </template>
                                        <v-list>
                                            <v-list-item v-for="(lang, locale) in allLanguages" :key="locale">
                                                <v-list-item-title
                                                    class="cursor-pointer"
                                                    @click="changeLanguage(locale)">
                                                    {{ lang.displayName }}
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-col>

                                <v-col cols="12" class="d-flex align-center pt-0">
                                    <v-divider />
                                </v-col>

                                <v-col cols="12" class="text-center text-sm">
                                    <span>Powered by </span>
                                    <a href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</a>&nbsp;<span>{{ version }}</span>
                                </v-col>
                            </v-row>
                        </v-card-text>
                    </v-card>
                </div>
            </v-col>
        </v-row>

        <confirm-dialog ref="confirmDialog"/>
        <snack-bar ref="snackbar" />
    </div>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';

import assetConstants from '@/consts/asset.js';

import {
    mdiChevronLeft,
} from '@mdi/js';

export default {
    data() {
        return {
            email: '',
            requesting: false,
            icons: {
                left: mdiChevronLeft
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore),
        ezBookkeepingLogoPath() {
            return assetConstants.ezBookkeepingLogoPath;
        },
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        isDarkMode() {
            return this.globalTheme.global.name.value === 'dark';
        },
        currentLanguageName() {
            return this.$locale.getCurrentLanguageDisplayName();
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        requestResetPassword() {
            const self = this;

            if (!self.email) {
                self.$refs.snackbar.showMessage('Email address cannot be empty');
                return;
            }

            self.requesting = true;

            self.rootStore.requestResetPassword({
                email: self.email
            }).then(() => {
                self.requesting = false;
                self.$refs.snackbar.showMessage('Password reset email has been sent');
            }).catch(error => {
                self.requesting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    }
}
</script>
