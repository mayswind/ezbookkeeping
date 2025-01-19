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
                    <v-img max-width="600px" src="img/desktop/people1.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ tt('Welcome to ezBookkeeping') }}</h4>
                            <p class="mb-0">{{ tt('Please log in with your ezBookkeeping account') }}</p>
                            <p class="mt-1 mb-0" v-if="tips">{{ tips }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12">
                                        <v-text-field
                                            type="text"
                                            autocomplete="username"
                                            autofocus="autofocus"
                                            :disabled="show2faInput || logining || verifying"
                                            :label="tt('Username')"
                                            :placeholder="tt('Your username or email')"
                                            v-model="username"
                                            @input="tempToken = ''"
                                            @keyup.enter="passwordInput?.focus()"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-text-field
                                            autocomplete="current-password"
                                            ref="passwordInput"
                                            type="password"
                                            :disabled="show2faInput || logining || verifying"
                                            :label="tt('Password')"
                                            :placeholder="tt('Your password')"
                                            v-model="password"
                                            @input="tempToken = ''"
                                            @keyup.enter="login"
                                        />
                                    </v-col>

                                    <v-col cols="12" v-show="show2faInput">
                                        <v-text-field
                                            type="number"
                                            autocomplete="one-time-code"
                                            ref="passcodeInput"
                                            :disabled="logining || verifying"
                                            :label="tt('Passcode')"
                                            :placeholder="tt('Passcode')"
                                            :append-inner-icon="icons.backupCode"
                                            v-model="passcode"
                                            @click:append-inner="twoFAVerifyType = 'backupcode'"
                                            @keyup.enter="verify"
                                            v-if="twoFAVerifyType === 'passcode'"
                                        />
                                        <v-text-field
                                            type="text"
                                            :disabled="logining || verifying"
                                            :label="tt('Backup Code')"
                                            :placeholder="tt('Backup Code')"
                                            :append-inner-icon="icons.passcode"
                                            v-model="backupCode"
                                            @click:append-inner="twoFAVerifyType = 'passcode'"
                                            @keyup.enter="verify"
                                            v-if="twoFAVerifyType === 'backupcode'"
                                        />
                                    </v-col>

                                    <v-col cols="12" class="py-0 mt-1 mb-4">
                                        <div class="d-flex align-center justify-space-between flex-wrap">
                                            <a href="javascript:void(0);" @click="showMobileQrCode = true">
                                                <span class="nav-item-title">{{ tt('Use on Mobile Device') }}</span>
                                            </a>
                                            <v-spacer/>
                                            <router-link class="text-primary" to="/forgetpassword"
                                                         :class="{'disabled': !isUserForgetPasswordEnabled()}">
                                                {{ tt('Forget Password?') }}
                                            </router-link>
                                        </div>
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block :disabled="inputIsEmpty || logining || verifying"
                                               @click="login" v-if="!show2faInput">
                                            {{ tt('Log In') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="logining"></v-progress-circular>
                                        </v-btn>
                                        <v-btn block :disabled="twoFAInputIsEmpty || logining || verifying"
                                               @click="verify" v-else-if="show2faInput">
                                            {{ tt('Continue') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="verifying"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" class="text-center text-base">
                                        <span class="me-1">{{ tt('Don\'t have an account?') }}</span>
                                        <router-link class="text-primary" to="/signup"
                                                     :class="{'disabled': !isUserRegistrationEnabled()}">
                                            {{ tt('Create an account') }}
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
                                                   :disabled="logining || verifying"
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

        <switch-to-mobile-dialog v-model:show="showMobileQrCode" />
        <snack-bar ref="snackbar" />
    </div>
</template>

<script setup lang="ts">
import { VTextField } from 'vuetify/components/VTextField';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';
import { useLoginPageBase } from '@/views/base/LoginPageBase.ts';

import { useRootStore } from '@/stores/index.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { KnownErrorCode } from '@/consts/api.ts';
import { ThemeType } from '@/core/theme.ts';
import { isUserRegistrationEnabled, isUserForgetPasswordEnabled, isUserVerifyEmailEnabled } from '@/lib/server_settings.ts';

import {
    mdiOnepassword,
    mdiHelpCircleOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();
const theme = useTheme();

const { tt, getCurrentLanguageDisplayName, getAllLanguageOptions } = useI18n();

const rootStore = useRootStore();

const {
    version,
    username,
    password,
    passcode,
    backupCode,
    tempToken,
    twoFAVerifyType,
    logining,
    verifying,
    inputIsEmpty,
    twoFAInputIsEmpty,
    tips,
    changeLanguage,
    doAfterLogin
} = useLoginPageBase();

const icons = {
    passcode: mdiOnepassword,
    backupCode: mdiHelpCircleOutline
};

const passwordInput = useTemplateRef<VTextField>('passwordInput');
const passcodeInput = useTemplateRef<VTextField>('passcodeInput');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const show2faInput = ref<boolean>(false);
const showMobileQrCode = ref<boolean>(false);

const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(false));
const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const currentLanguageName = computed<string>(() => getCurrentLanguageDisplayName());

function login(): void {
    if (!username.value) {
        snackbar.value?.showMessage('Username cannot be blank');
        return;
    }

    if (!password.value) {
        snackbar.value?.showMessage('Password cannot be blank');
        return;
    }

    if (tempToken.value) {
        show2faInput.value = true;
        return;
    }

    if (logining.value) {
        return;
    }

    logining.value = true;

    rootStore.authorize({
        loginName: username.value,
        password: password.value
    }).then(authResponse => {
        logining.value = false;

        if (authResponse.need2FA) {
            tempToken.value = authResponse.token;
            show2faInput.value = true;

            nextTick(() => {
                if (passcodeInput.value) {
                    passcodeInput.value.focus();
                    passcodeInput.value.select();
                }
            });

            return;
        }

        doAfterLogin(authResponse);
        router.replace('/');
    }).catch(error => {
        logining.value = false;

        if (isUserVerifyEmailEnabled() && error.error && error.error.errorCode === KnownErrorCode.UserEmailNotVerified && error.error.context && error.error.context.email) {
            router.push(`/verify_email?email=${encodeURIComponent(error.error.context.email)}&emailSent=${error.error.context.hasValidEmailVerifyToken || false}`);
            return;
        }

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function verify(): void {
    if (twoFAInputIsEmpty.value || verifying.value) {
        return;
    }

    if (twoFAVerifyType.value === 'passcode' && !passcode.value) {
        snackbar.value?.showMessage('Passcode cannot be blank');
        return;
    } else if (twoFAVerifyType.value === 'backupcode' && !backupCode.value) {
        snackbar.value?.showMessage('Backup code cannot be blank');
        return;
    }

    verifying.value = true;

    rootStore.authorize2FA({
        token: tempToken.value,
        passcode: twoFAVerifyType.value === 'passcode' ? passcode.value : null,
        recoveryCode: twoFAVerifyType.value === 'backupcode' ? backupCode.value : null
    }).then(authResponse => {
        verifying.value = false;

        doAfterLogin(authResponse);
        router.replace('/');
    }).catch(error => {
        verifying.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}
</script>
