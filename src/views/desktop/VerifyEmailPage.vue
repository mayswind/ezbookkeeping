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
                    <v-img max-width="320px" src="img/desktop/people2.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ $t('Verify your email') }}</h4>
                            <p class="mb-0" v-if="token && loading">{{ $t('Verifying...') }}</p>
                            <p class="mb-0" v-if="token && verified">{{ $t('Email address is verified') }}</p>
                            <p class="mb-0" v-if="token && !verified && errorMessage">{{ errorMessage }}</p>
                            <p class="mb-0" v-if="!token && !email">{{ $t('Parameter Invalid') }}</p>
                            <p class="mb-0" v-if="!token && email">{{ $t(hasValidEmailVerifyToken ? 'format.misc.accountActivationAndResendValidationEmailTip' : 'format.misc.resendValidationEmailTip', { email: email }) }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12" v-if="!loading && !token && email && isUserVerifyEmailEnabled">
                                        <v-text-field
                                            autocomplete="password"
                                            type="password"
                                            :disabled="loading || resending"
                                            :label="$t('Password')"
                                            :placeholder="$t('Your password')"
                                            v-model="password"
                                            @keyup.enter="resendEmail"
                                        />
                                    </v-col>

                                    <v-col cols="12" v-if="!loading && !token && email && isUserVerifyEmailEnabled">
                                        <v-btn block type="submit" :disabled="loading || resending || !password" @click="resendEmail">
                                            {{ $t('Resend Validation Email') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="resending"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12">
                                        <router-link class="d-flex align-center justify-center" :to="verified ? '/' : '/login'"
                                                     :class="{ 'disabled': loading || resending }">
                                            <v-icon :icon="icons.left"/>
                                            <span v-if="!verified">{{ $t('Back to login page') }}</span>
                                            <span v-else-if="verified">{{ $t('Back to home page') }}</span>
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
                                                   :disabled="resending"
                                                   v-bind="props">{{ currentLanguageName }}</v-btn>
                                        </template>
                                        <v-list>
                                            <v-list-item v-for="lang in allLanguages" :key="lang.languageTag">
                                                <v-list-item-title
                                                    class="cursor-pointer"
                                                    @click="changeLanguage(lang.languageTag)">
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
        <snack-bar ref="snackbar" @update:show="onSnackbarShowStateChanged" />
    </div>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import { isUserVerifyEmailEnabled } from '@/lib/server_settings.ts';

import {
    mdiChevronLeft
} from '@mdi/js';

export default {
    props: [
        'email',
        'token',
        'hasValidEmailVerifyToken'
    ],
    data() {
        return {
            password: '',
            loading: true,
            resending: false,
            verified: false,
            errorMessage: '',
            icons: {
                left: mdiChevronLeft
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore),
        ezBookkeepingLogoPath() {
            return APPLICATION_LOGO_PATH;
        },
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfoArray(false);
        },
        isDarkMode() {
            return this.globalTheme.global.name.value === ThemeType.Dark;
        },
        currentLanguageName() {
            return this.$locale.getCurrentLanguageDisplayName();
        },
        isUserVerifyEmailEnabled() {
            return isUserVerifyEmailEnabled();
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    created() {
        const self = this;

        self.verified = false;
        self.loading = true;

        if (!self.token) {
            self.loading = false;
            return;
        }

        self.rootStore.verifyEmail({
            token: self.token,
            requestNewToken: !self.$user.isUserLogined()
        }).then(() => {
            self.loading = false;
            self.verified = true;
            self.$refs.snackbar.showMessage('Email address is verified');
        }).catch(error => {
            self.loading = false;
            self.verified = false;

            if (!error.processed) {
                self.errorMessage = self.$tError(error.message || error);
                self.$refs.snackbar.showError(error);
            }
        });
    },
    methods: {
        resendEmail() {
            const self = this;

            self.resending = true;

            self.rootStore.resendVerifyEmailByUnloginUser({
                email: self.email,
                password: self.password
            }).then(() => {
                self.resending = false;
                self.$refs.snackbar.showMessage('Validation email has been sent');
            }).catch(error => {
                self.resending = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        onSnackbarShowStateChanged(newValue) {
            if (!newValue && this.verified && this.$user.isUserLogined()) {
                this.$router.replace('/');
            }
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    }
}
</script>
