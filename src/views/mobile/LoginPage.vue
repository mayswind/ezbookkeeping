<template>
    <f7-page no-navbar no-swipeback login-screen hide-toolbar-on-scroll>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
            <f7-block class="login-page-tile margin-vertical-half">{{ tt('global.app.title') }}</f7-block>
        </f7-login-screen-title>

        <f7-list inset v-if="tips">
            <f7-block-footer>{{ tips }}</f7-block-footer>
        </f7-list>

        <f7-list form dividers class="margin-bottom-half" v-if="isInternalAuthEnabled()">
            <f7-list-input
                type="text"
                autocomplete="username"
                clear-button
                :disabled="loggingInByPassword || loggingInByOAuth2"
                :label="tt('Username')"
                :placeholder="tt('Your username or email')"
                v-model:value="username"
                @input="tempToken = ''"
            ></f7-list-input>
            <f7-list-input
                type="password"
                autocomplete="current-password"
                clear-button
                :disabled="loggingInByPassword || loggingInByOAuth2"
                :label="tt('Password')"
                :placeholder="tt('Your password')"
                v-model:value="password"
                @input="tempToken = ''"
                @keyup.enter="loginByPressEnter"
            ></f7-list-input>
        </f7-list>

        <f7-list class="no-margin-vertical">
            <f7-list-item>
                <template #title>
                    <small>
                        <f7-link :class="{ 'disabled': loggingInByPassword || loggingInByOAuth2 }" @click="switchToDesktopVersion">{{ tt('Switch to Desktop Version') }}</f7-link>
                    </small>
                </template>
                <template #after>
                    <small>
                        <f7-link :class="{ 'disabled': !isUserForgetPasswordEnabled() || loggingInByPassword || loggingInByOAuth2 }" @click="forgetPasswordEmail = ''; showForgetPasswordSheet = true">{{ tt('Forget Password?') }}</f7-link>
                    </small>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list class="margin-vertical-half">
            <f7-list-button :class="{ 'disabled': inputIsEmpty || loggingInByPassword || loggingInByOAuth2 }" :text="tt('Log In')"
                            @click="login" v-if="isInternalAuthEnabled()"></f7-list-button>
            <f7-list-item class="login-divider display-flex align-items-center" v-if="isInternalAuthEnabled() && isOAuth2Enabled()">
                <hr class="margin-inline-end-half" />
                <small>{{ tt('or') }}</small>
                <hr class="margin-inline-start-half" />
            </f7-list-item>
            <f7-list-button external :class="{ 'disabled': loggingInByPassword || loggingInByOAuth2 }" :href="oauth2LoginUrl" :text="oauth2LoginDisplayName"
                            @click="loginByOAuth2" v-if="isOAuth2Enabled()"></f7-list-button>
            <f7-block-footer v-if="isInternalAuthEnabled()">
                <span>{{ tt('Don\'t have an account?') }}</span>&nbsp;
                <f7-link :class="{ 'disabled': !isUserRegistrationEnabled() || loggingInByPassword || loggingInByOAuth2 }" href="/signup" :text="tt('Create an account')"></f7-link>
            </f7-block-footer>
            <f7-block-footer class="padding-bottom">
            </f7-block-footer>
        </f7-list>

        <f7-list class="login-page-bottom">
            <f7-block-footer>
                <language-select-button :disabled="loggingInByPassword || loggingInByOAuth2" />

                <div class="login-page-powered-by margin-top-half">
                    <span>Powered by</span>
                    <f7-link @click="openExternalUrl('https://github.com/mayswind/ezbookkeeping')" target="_blank">ezBookkeeping</f7-link>
                    <span>{{ version }}</span>
                </div>
            </f7-block-footer>
        </f7-list>

        <f7-toolbar class="login-page-fixed-bottom" tabbar bottom :outline="false">
            <language-select-button :disabled="loggingInByPassword || loggingInByOAuth2" />

            <div class="login-page-powered-by margin-top-half">
                <span>Powered by</span>
                <f7-link @click="openExternalUrl('https://github.com/mayswind/ezbookkeeping')" target="_blank">ezBookkeeping</f7-link>
                <span>{{ version }}</span>
            </div>
        </f7-toolbar>

        <f7-sheet
            style="height:auto"
            :opened="show2faSheet" @sheet:closed="show2faSheet = false"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ tt('Two-Factor Authentication') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <f7-list strong class="no-margin">
                        <f7-list-input
                            type="number"
                            autocomplete="one-time-code"
                            outline
                            floating-label
                            clear-button
                            class="no-margin no-padding-bottom"
                            v-if="twoFAVerifyType === 'passcode'"
                            :label="tt('Passcode')"
                            :placeholder="tt('Passcode')"
                            v-model:value="passcode"
                            @keyup.enter="verify"
                        ></f7-list-input>
                        <f7-list-input
                            outline
                            floating-label
                            clear-button
                            class="no-margin no-padding-bottom"
                            v-if="twoFAVerifyType === 'backupcode'"
                            :label="tt('Backup Code')"
                            :placeholder="tt('Backup Code')"
                            v-model:value="backupCode"
                            @keyup.enter="verify"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': twoFAInputIsEmpty || verifying }" :text="tt('Verify')" @click="verify"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="switch2FAVerifyType" :text="tt(twoFAVerifyTypeSwitchName)"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
            style="height:auto"
            :opened="showForgetPasswordSheet" @sheet:closed="showForgetPasswordSheet = false"
        >
            <div class="swipe-handler" style="z-index: 10"></div>
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ tt('Forget Password?') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin">
                        <span>{{ tt('Please enter your email address used for registration and we\'ll send you an email with a reset password link') }}</span>
                    </p>
                    <f7-list strong class="no-margin">
                        <f7-list-input
                            type="email"
                            autocomplete="email"
                            outline
                            floating-label
                            clear-button
                            class="no-margin no-padding-bottom"
                            :label="tt('E-mail')"
                            :placeholder="tt('Your email address')"
                            v-model:value="forgetPasswordEmail"
                            @keyup.enter="requestResetPassword"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !forgetPasswordEmail || requestingForgetPassword }" :text="tt('Send Reset Link')" @click="requestResetPassword"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': requestingForgetPassword }" @click="showForgetPasswordSheet = false" :text="tt('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <password-input-sheet :title="tt('Verify your email')"
                              :hint="tt(hasValidEmailVerifyToken ? 'format.misc.accountActivationAndResendValidationEmailTip' : 'format.misc.resendValidationEmailTip', { email: resendVerifyEmail })"
                              :confirm-disabled="requestingResendVerifyEmail"
                              :cancel-disabled="requestingResendVerifyEmail"
                              v-model:show="showVerifyEmailSheet"
                              v-model="currentPasswordForResendVerifyEmail"
                              @password:confirm="requestResendVerifyEmail">
        </password-input-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useLoginPageBase } from '@/views/base/LoginPageBase.ts';

import { useRootStore } from '@/stores/index.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { KnownErrorCode } from '@/consts/api.ts';

import { generateRandomUUID } from '@/lib/misc.ts';
import {
    isUserRegistrationEnabled,
    isUserForgetPasswordEnabled,
    isUserVerifyEmailEnabled,
    isInternalAuthEnabled,
    isOAuth2Enabled
} from '@/lib/server_settings.ts';
import { getDesktopVersionPath } from '@/lib/version.ts';
import { useI18nUIComponents, showLoading, hideLoading, isModalShowing } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showAlert, showConfirm, showToast, openExternalUrl } = useI18nUIComponents();

const rootStore = useRootStore();

const {
    version,
    username,
    password,
    passcode,
    backupCode,
    tempToken,
    twoFAVerifyType,
    oauth2ClientSessionId,
    loggingInByPassword,
    loggingInByOAuth2,
    verifying,
    inputIsEmpty,
    twoFAInputIsEmpty,
    oauth2LoginUrl,
    oauth2LoginDisplayName,
    tips,
    doAfterLogin
} = useLoginPageBase('mobile');

const forgetPasswordEmail = ref<string>('');
const resendVerifyEmail = ref<string>('');
const hasValidEmailVerifyToken = ref<boolean>(false);
const currentPasswordForResendVerifyEmail = ref<string>('');
const requestingForgetPassword = ref<boolean>(false);
const requestingResendVerifyEmail = ref<boolean>(false);
const show2faSheet = ref<boolean>(false);
const showForgetPasswordSheet = ref<boolean>(false);
const showVerifyEmailSheet = ref<boolean>(false);

const twoFAVerifyTypeSwitchName = computed<string>(() => {
    if (twoFAVerifyType.value === 'backupcode') {
        return 'Use Passcode';
    } else {
        return 'Use Backup Code';
    }
});

function switchToDesktopVersion(): void {
    showConfirm('Are you sure you want to switch to desktop version?', () => {
        window.location.replace(getDesktopVersionPath());
    });
}

function login(): void {
    const router = props.f7router;

    if (!username.value) {
        showAlert('Username cannot be blank');
        return;
    }

    if (!password.value) {
        showAlert('Password cannot be blank');
        return;
    }

    if (tempToken.value) {
        show2faSheet.value = true;
        return;
    }

    loggingInByPassword.value = true;
    resendVerifyEmail.value = '';
    hasValidEmailVerifyToken.value = false;
    currentPasswordForResendVerifyEmail.value = '';
    showLoading(() => loggingInByPassword.value);

    rootStore.authorize({
        loginName: username.value,
        password: password.value
    }).then(authResponse => {
        loggingInByPassword.value = false;
        hideLoading();

        if (authResponse.need2FA) {
            tempToken.value = authResponse.token;
            show2faSheet.value = true;
            return;
        }

        doAfterLogin(authResponse);
        router.refreshPage();
    }).catch(error => {
        loggingInByPassword.value = false;
        hideLoading();

        if (isUserVerifyEmailEnabled() && error.error && error.error.errorCode === KnownErrorCode.UserEmailNotVerified && error.error.context && error.error.context.email) {
            resendVerifyEmail.value = error.error.context.email;
            hasValidEmailVerifyToken.value = error.error.context.hasValidEmailVerifyToken || false;
            currentPasswordForResendVerifyEmail.value = '';
            showVerifyEmailSheet.value = true;
            return;
        }

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function loginByPressEnter(): void {
    if (isModalShowing()) {
        return;
    }

    return login();
}

function loginByOAuth2(): void {
    loggingInByOAuth2.value = true;
    showLoading();
}

function verify(): void {
    const router = props.f7router;

    if (twoFAInputIsEmpty.value || verifying.value) {
        return;
    }

    if (twoFAVerifyType.value === 'passcode' && !passcode.value) {
        showAlert('Passcode cannot be blank');
        return;
    } else if (twoFAVerifyType.value === 'backupcode' && !backupCode.value) {
        showAlert('Backup code cannot be blank');
        return;
    }

    verifying.value = true;
    showLoading(() => verifying.value);

    rootStore.authorize2FA({
        token: tempToken.value,
        passcode: twoFAVerifyType.value === 'passcode' ? passcode.value : null,
        recoveryCode: twoFAVerifyType.value === 'backupcode' ? backupCode.value : null
    }).then(authResponse => {
        verifying.value = false;
        hideLoading();

        doAfterLogin(authResponse);
        show2faSheet.value = false;
        router.refreshPage();
    }).catch(error => {
        verifying.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function requestResetPassword(): void {
    if (!forgetPasswordEmail.value) {
        showAlert('Email address cannot be blank');
        return;
    }

    requestingForgetPassword.value = true;
    showLoading(() => requestingForgetPassword.value);

    rootStore.requestResetPassword({
        email: forgetPasswordEmail.value
    }).then(() => {
        requestingForgetPassword.value = false;
        hideLoading();

        showToast('Password reset email has been sent');
        showForgetPasswordSheet.value = false;
    }).catch(error => {
        requestingForgetPassword.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function requestResendVerifyEmail(): void {
    if (!currentPasswordForResendVerifyEmail.value) {
        showToast('Current password cannot be blank');
        return;
    }

    requestingResendVerifyEmail.value = true;
    showLoading(() => requestingResendVerifyEmail.value);

    rootStore.resendVerifyEmailByUnloginUser({
        email: resendVerifyEmail.value,
        password: currentPasswordForResendVerifyEmail.value
    }).then(() => {
        requestingResendVerifyEmail.value = false;
        hideLoading();

        showToast('Validation email has been sent');
        showVerifyEmailSheet.value = false;
    }).catch(error => {
        requestingResendVerifyEmail.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function switch2FAVerifyType(): void {
    if (twoFAVerifyType.value === 'passcode') {
        twoFAVerifyType.value = 'backupcode';
    } else {
        twoFAVerifyType.value = 'passcode';
    }
}

oauth2ClientSessionId.value = generateRandomUUID();
</script>

<style>
.login-divider > .item-content {
    width: 100%;
    min-height: 0;

    > .item-inner {
        padding-top: 0;
        padding-bottom: 0;
        min-height: 0;

        > small {
            opacity: 0.7;
        }

        > hr {display: block;
            flex: 1 1 100%;
            height: 0;
            max-height: 0;
            border-style: solid;
            border-width: thin 0 0 0;
            opacity: 0.12;
            transition: inherit;
        }
    }
}
</style>
