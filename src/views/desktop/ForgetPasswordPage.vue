<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ tt('global.app.title') }}</h1>
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
                            <h4 class="text-h4 mb-2">{{ tt('Forget Password?') }}</h4>
                        </v-card-text>
                        
                        <v-card-text class="pb-0 mb-6">
                            <template v-if="isUserForgetPasswordEnabled()">
                                <v-form>
                                    <v-row>
                                        <v-col cols="12">
                                            <p class="mb-0">{{ tt('Please enter your email address used for registration and we\'ll send you an email with a reset password link') }}</p>
                                        </v-col>
                                        <v-col cols="12">
                                            <v-text-field
                                                type="email"
                                                autocomplete="email"
                                                autofocus="autofocus"
                                                :disabled="requesting"
                                                :label="tt('E-mail')"
                                                :placeholder="tt('Your email address')"
                                                v-model="email"
                                                @keyup.enter="requestResetPassword"
                                            />
                                        </v-col>

                                        <v-col cols="12">
                                            <v-btn block type="submit" :disabled="!email || requesting" @click="requestResetPassword">
                                                {{ tt('Send Reset Link') }}
                                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="requesting"></v-progress-circular>
                                            </v-btn>
                                        </v-col>

                                        <v-col cols="12">
                                            <router-link class="d-flex align-center justify-center" to="/login"
                                                        :class="{ 'disabled': requesting }">
                                                <v-icon :icon="icons.left"/>
                                                <span>{{ tt('Back to login page') }}</span>
                                            </router-link>
                                        </v-col>
                                    </v-row>
                                </v-form>
                            </template>
                            <template v-else>
                                <v-row>
                                    <v-col cols="12">
                                        <p class="mb-0">{{ tt('Resetting passwords is not enabled. Please contact the administrator for assistance.') }}</p>
                                    </v-col>
                                    <v-col cols="12">
                                        <router-link class="d-flex align-center justify-center" to="/login">
                                            <v-icon :icon="icons.left"/>
                                            <span>{{ tt('Back to login page') }}</span>
                                        </router-link>
                                    </v-col>
                                </v-row>
                            </template>

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
        <snack-bar ref="snackbar" />
    </div>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useTheme } from 'vuetify';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';

import { isUserForgetPasswordEnabled } from '@/lib/server_settings.ts';
import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import { getVersion } from '@/lib/version.ts';

import {
    mdiChevronLeft,
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const theme = useTheme();

const { tt, getCurrentLanguageDisplayName, getAllLanguageOptions, setLanguage } = useI18n();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();

const icons = {
    left: mdiChevronLeft
};

const version = `v${getVersion()}`;

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const email = ref<string>('');
const requesting = ref<boolean>(false);

const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(false));
const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const currentLanguageName = computed<string>(() => getCurrentLanguageDisplayName());

function requestResetPassword(): void {
    if (!email.value) {
        snackbar.value?.showMessage('Email address cannot be blank');
        return;
    }

    requesting.value = true;

    rootStore.requestResetPassword({
        email: email.value
    }).then(() => {
        requesting.value = false;
        snackbar.value?.showMessage('Password reset email has been sent');
    }).catch(error => {
        requesting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function changeLanguage(locale: string): void {
    const localeDefaultSettings = setLanguage(locale);
    settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
}
</script>
