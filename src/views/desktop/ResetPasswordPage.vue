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
                            <h4 class="text-h4 mb-2">{{ $t('Reset Password') }}</h4>
                            <p class="mb-0">{{ $t('Please enter your email address again, and input the new password.') }}</p>
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
                                            :disabled="updating"
                                            :label="$t('E-mail')"
                                            :placeholder="$t('Your email address')"
                                            v-model="email"
                                            @keyup.enter="$refs.passwordInput.focus()"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-text-field
                                            autocomplete="new-password"
                                            clearable
                                            ref="passwordInput"
                                            :type="isNewPasswordVisible ? 'text' : 'password'"
                                            :disabled="updating"
                                            :label="$t('Password')"
                                            :placeholder="$t('Your password')"
                                            :append-inner-icon="isNewPasswordVisible ? icons.eyeSlash : icons.eye"
                                            v-model="newPassword"
                                            @click:append-inner="isNewPasswordVisible = !isNewPasswordVisible"
                                            @keyup.enter="$refs.confirmPasswordInput.focus()"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-text-field
                                            clearable
                                            ref="confirmPasswordInput"
                                            :type="isConfirmPasswordVisible ? 'text' : 'password'"
                                            :disabled="updating"
                                            :label="$t('Confirm Password')"
                                            :placeholder="$t('Re-enter the password')"
                                            :append-inner-icon="isConfirmPasswordVisible ? icons.eyeSlash : icons.eye"
                                            v-model="confirmPassword"
                                            @click:append-inner="isConfirmPasswordVisible = !isConfirmPasswordVisible"
                                            @keyup.enter="resetPassword"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block :disabled="!email || !newPassword || !confirmPassword || updating" @click="resetPassword">
                                            {{ $t('Update Password') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="updating"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12">
                                        <router-link class="d-flex align-center justify-center" to="/login"
                                                     :class="{ 'disabled': updating }">
                                            <v-icon :icon="icons.left"/>
                                            <span>{{ $t('Back to login page') }}</span>
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
                                                   :disabled="updating"
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
        <snack-bar ref="snackbar" @update:show="onSnackbarShowStateChanged" />
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
    mdiEyeOffOutline,
    mdiEyeOutline
} from '@mdi/js';

export default {
    props: [
        'token'
    ],
    data() {
        return {
            email: '',
            newPassword: '',
            confirmPassword: '',
            isNewPasswordVisible: false,
            isConfirmPasswordVisible: false,
            updating: false,
            passwordChanged: false,
            icons: {
                left: mdiChevronLeft,
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore),
        inputProblemMessage() {
            if (!this.email) {
                return 'Email address cannot be blank';
            } else if (!this.newPassword && !this.confirmPassword) {
                return 'Nothing has been modified';
            } else if (!this.newPassword && this.confirmPassword) {
                return 'New password cannot be blank';
            } else if (this.newPassword && !this.confirmPassword) {
                return 'Password confirmation cannot be blank';
            } else if (this.newPassword && this.confirmPassword && this.newPassword !== this.confirmPassword) {
                return 'Password and password confirmation do not match';
            } else {
                return null;
            }
        },
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
        resetPassword() {
            const self = this;
            self.passwordChanged = false;

            const problemMessage = self.inputProblemMessage;

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            self.updating = true;

            self.rootStore.resetPassword({
                token: self.token,
                email: self.email,
                password: self.newPassword
            }).then(() => {
                self.updating = false;
                self.passwordChanged = true;
                self.$refs.snackbar.showMessage('Password has been updated');
            }).catch(error => {
                self.updating = false;
                self.passwordChanged = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        onSnackbarShowStateChanged(newValue) {
            if (!newValue && this.passwordChanged) {
                this.$router.replace('/login');
            }
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    }
}
</script>
