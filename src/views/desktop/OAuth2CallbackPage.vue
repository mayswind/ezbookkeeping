<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ tt('global.app.title') }}</h1>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="auth-image-background d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer" v-if="!isDarkMode">
                    <v-img class="img-with-direction" src="img/desktop/background.svg"/>
                </div>
                <div class="d-flex auth-img-footer" v-if="isDarkMode">
                    <v-img class="img-with-direction" src="img/desktop/background-dark.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img class="img-with-direction" max-width="300px" src="img/desktop/people2.svg" v-if="!isDarkMode"/>
                    <v-img class="img-with-direction" max-width="300px" src="img/desktop/people2-dark.svg" v-else-if="isDarkMode"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ oauth2LoginDisplayName }}</h4>
                            <p class="mb-0" v-if="!error && platform && token && !userName">{{ tt('Logging in...') }}</p>
                            <p class="mb-0" v-else-if="!error && userName">{{ tt('format.misc.oauth2bindTip', { providerName: oauth2ProviderDisplayName, userName: userName }) }}</p>
                            <p class="mb-0" v-else-if="error">{{ tt(error) }}</p>
                            <p class="mb-0" v-else>{{ tt('An error occurred') }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6" v-if="!error && userName">
                            <v-form>
                                <v-row>
                                    <v-col cols="12">
                                        <v-text-field
                                            type="password"
                                            autocomplete="password"
                                            :autofocus="true"
                                            :disabled="logining"
                                            :label="tt('Password')"
                                            :placeholder="tt('Your password')"
                                            v-model="password"
                                            @keyup.enter="verify"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block type="submit" :disabled="!password || logining" @click="verify">
                                            {{ tt('Continue') }}
                                            <v-progress-circular indeterminate size="22" class="ms-2" v-if="logining"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12">
                                        <router-link class="d-flex align-center justify-center" to="/login"
                                                     :class="{ 'disabled': logining }">
                                            <v-icon class="icon-with-direction" :icon="mdiChevronLeft"/>
                                            <span>{{ tt('Back to login page') }}</span>
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
                                    <language-select-button :disabled="logining" />
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

        <snack-bar ref="snackbar" />
    </div>
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useLoginPageBase } from '@/views/base/LoginPageBase.ts';

import { useRootStore } from '@/stores/index.ts';

import { ThemeType } from '@/core/theme.ts';
import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { KnownErrorCode } from '@/consts/api.ts';

import {
    isUserVerifyEmailEnabled,
    getOAuth2Provider,
    getOIDCCustomDisplayNames
} from '@/lib/server_settings.ts';

import {
    mdiChevronLeft
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    token?: string;
    provider?: string;
    platform?: string;
    userName?: string;
    error?: string;
}>();

const router = useRouter();
const theme = useTheme();

const { tt, getLocalizedOAuth2ProviderName, getLocalizedOAuth2LoginText } = useI18n();

const rootStore = useRootStore();

const {
    version,
    password,
    logining,
    doAfterLogin
} = useLoginPageBase('desktop');

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const oauth2ProviderDisplayName = computed<string>(() => getLocalizedOAuth2ProviderName(getOAuth2Provider(), getOIDCCustomDisplayNames()));
const oauth2LoginDisplayName = computed<string>(() => getLocalizedOAuth2LoginText(getOAuth2Provider(), getOIDCCustomDisplayNames()));

const inputProblemMessage = computed<string | null>(() => {
    if (!password.value) {
        return 'Password cannot be blank';
    } else {
        return null;
    }
});

function verify(): void  {
    const problemMessage = inputProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    logining.value = true;

    rootStore.authorizeOAuth2({
        provider: props.provider || '',
        password: password.value,
        token: props.token || ''
    }).then(authResponse => {
        logining.value = false;
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

if (!props.error && props.platform && props.token && !props.userName) {
    logining.value = true;

    rootStore.authorizeOAuth2({
        provider: props.provider || '',
        token: props.token || ''
    }).then(authResponse => {
        logining.value = false;
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
</script>
