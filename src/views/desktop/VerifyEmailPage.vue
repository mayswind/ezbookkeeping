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
                    <v-img max-width="320px" src="img/desktop/people2.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ tt('Verify your email') }}</h4>
                            <p class="mb-0" v-if="token && loading">{{ tt('Verifying...') }}</p>
                            <p class="mb-0" v-if="token && verified">{{ tt('Email address is verified') }}</p>
                            <p class="mb-0" v-if="token && !verified && errorMessage">{{ errorMessage }}</p>
                            <p class="mb-0" v-if="!token && !email">{{ tt('Parameter Invalid') }}</p>
                            <p class="mb-0" v-if="!token && email">{{ tt(hasValidEmailVerifyToken ? 'format.misc.accountActivationAndResendValidationEmailTip' : 'format.misc.resendValidationEmailTip', { email: email }) }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12" v-if="!loading && !token && email && isUserVerifyEmailEnabled()">
                                        <v-text-field
                                            autocomplete="password"
                                            type="password"
                                            :disabled="loading || resending"
                                            :label="tt('Password')"
                                            :placeholder="tt('Your password')"
                                            v-model="password"
                                            @keyup.enter="resendEmail"
                                        />
                                    </v-col>

                                    <v-col cols="12" v-if="!loading && !token && email && isUserVerifyEmailEnabled()">
                                        <v-btn block type="submit" :disabled="loading || resending || !password" @click="resendEmail">
                                            {{ tt('Resend Validation Email') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="resending"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12">
                                        <router-link class="d-flex align-center justify-center" :to="verified ? '/' : '/login'"
                                                     :class="{ 'disabled': loading || resending }">
                                            <v-icon :icon="mdiChevronLeft"/>
                                            <span v-if="!verified">{{ tt('Back to login page') }}</span>
                                            <span v-else-if="verified">{{ tt('Back to home page') }}</span>
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
                                    <language-select-button :disabled="resending" />
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

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import { isUserVerifyEmailEnabled } from '@/lib/server_settings.ts';
import { isUserLogined } from '@/lib/userstate.ts';
import { getClientDisplayVersion } from '@/lib/version.ts';

import {
    mdiChevronLeft
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    email: string;
    token?: string;
    hasValidEmailVerifyToken: boolean;
}>();

const router = useRouter();
const theme = useTheme();

const { tt, te } = useI18n();

const rootStore = useRootStore();

const version = `${getClientDisplayVersion()}`;

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const password = ref<string>('');
const loading = ref<boolean>(true);
const resending = ref<boolean>(false);
const verified = ref<boolean>(false);
const errorMessage = ref<string>('');

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

function init(): void {
    verified.value = false;
    loading.value = true;

    if (!props.token) {
        loading.value = false;
        return;
    }

    rootStore.verifyEmail({
        token: props.token,
        requestNewToken: !isUserLogined()
    }).then(() => {
        loading.value = false;
        verified.value = true;
        snackbar.value?.showMessage('Email address is verified');
    }).catch(error => {
        loading.value = false;
        verified.value = false;

        if (!error.processed) {
            errorMessage.value = te(error.message || error);
            snackbar.value?.showError(error);
        }
    });
}

function resendEmail(): void {
    resending.value = true;

    rootStore.resendVerifyEmailByUnloginUser({
        email: props.email,
        password: password.value
    }).then(() => {
        resending.value = false;
        snackbar.value?.showMessage('Validation email has been sent');
    }).catch(error => {
        resending.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function onSnackbarShowStateChanged(newValue: boolean): void {
    if (!newValue && verified.value && isUserLogined()) {
        router.replace('/');
    }
}

init();
</script>
