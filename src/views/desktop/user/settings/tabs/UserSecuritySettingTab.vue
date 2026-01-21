<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': updatingPassword }" :title="tt('Modify Password')">
                <v-form>
                    <v-card-text class="pt-0">
                        <span class="text-body-1">{{ tt('After changing the password, other devices will be logged out. Please use the new password to log in on other devices.') }}</span>
                    </v-card-text>

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    autocomplete="current-password"
                                    ref="currentPasswordInput"
                                    type="password"
                                    :disabled="updatingPassword"
                                    :label="tt('Current Password')"
                                    :placeholder="tt('Current Password')"
                                    v-model="currentPassword"
                                    @keyup.enter="newPasswordInput?.focus()"
                                />
                            </v-col>
                        </v-row>

                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    autocomplete="new-password"
                                    ref="newPasswordInput"
                                    type="password"
                                    :disabled="updatingPassword"
                                    :label="tt('New Password')"
                                    :placeholder="tt('New Password')"
                                    v-model="newPassword"
                                    @keyup.enter="confirmPasswordInput?.focus()"
                                />
                            </v-col>
                        </v-row>

                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    ref="confirmPasswordInput"
                                    type="password"
                                    :disabled="updatingPassword"
                                    :label="tt('Confirm Password')"
                                    :placeholder="tt('Re-enter the password')"
                                    v-model="confirmPassword"
                                    @keyup.enter="updatePassword"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn :disabled="!newPassword || !confirmPassword || updatingPassword" @click="updatePassword">
                            {{ tt('Save Changes') }}
                            <v-progress-circular indeterminate size="22" class="ms-2" v-if="updatingPassword"></v-progress-circular>
                        </v-btn>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12" v-if="isOAuth2Enabled() && (loadingExternalAuth || (thirdPartyLoginList && thirdPartyLoginList.length))">
            <v-card :class="{ 'disabled': loadingExternalAuth }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('Third-Party Login') }}</span>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :loading="loadingExternalAuth" @click="reloadExternalAuth(false)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-table class="table-striped text-no-wrap" :hover="!loadingExternalAuth">
                    <thead>
                    <tr>
                        <th>{{ tt('Type') }}</th>
                        <th>{{ tt('Username') }}</th>
                        <th>{{ tt('Linked Time') }}</th>
                        <th class="text-right">{{ tt('Operation') }}</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr :key="itemIdx"
                        v-for="itemIdx in (loadingExternalAuth && (!externalAuths || externalAuths.length < 1) ? [ 1 ] : [])">
                        <td class="px-0" colspan="4">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>

                    <tr :key="thirdPartyLogin.externalAuthType"
                        v-for="thirdPartyLogin in thirdPartyLoginList">
                        <td class="text-sm">
                            <v-icon start :icon="thirdPartyLogin.icon"/>
                            {{ thirdPartyLogin.displayName }}
                        </td>
                        <td class="text-sm">{{ thirdPartyLogin.externalUsername }}</td>
                        <td class="text-sm">{{ thirdPartyLogin.createdAt }}</td>
                        <td class="text-sm text-right">
                            <v-btn density="comfortable" variant="tonal"
                                   :disabled="loggingInByOAuth2"
                                   :href="oauth2LinkUrl"
                                   @click="loggingInByOAuth2 = true"
                                   v-if="!thirdPartyLogin.linked && isOAuth2Enabled() && getOAuth2Provider() === thirdPartyLogin.externalAuthType">
                                {{ tt('Link') }}
                                <v-progress-circular indeterminate size="22" class="ms-2" v-if="loggingInByOAuth2"></v-progress-circular>
                            </v-btn>
                            <v-btn density="comfortable" color="error" variant="tonal"
                                   :disabled="loadingExternalAuth"
                                   @click="unlinkExternalAuth(thirdPartyLogin)"
                                   v-if="thirdPartyLogin.linked">
                                {{ tt('Unlink') }}
                            </v-btn>
                        </td>
                    </tr>
                    </tbody>
                </v-table>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :class="{ 'disabled': loadingSession }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('Device & Sessions') }}</span>
                        <v-btn class="ms-3" density="compact" color="default" variant="outlined"
                               @click="generateToken" v-if="isAPITokenEnabled() || isMCPServerEnabled()">{{ tt('Generate Token') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :loading="loadingSession" @click="reloadSessions(false)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-table class="table-striped text-no-wrap" :hover="!loadingSession">
                    <thead>
                    <tr>
                        <th>{{ tt('Type') }}</th>
                        <th>{{ tt('Device Info') }}</th>
                        <th>{{ tt('Last Activity Time') }}</th>
                        <th class="text-right">
                            <v-btn density="comfortable" color="error" variant="tonal"
                                   :disabled="sessions.length < 2 || loadingSession"
                                   @click="revokeAllSessions">
                                {{ tt('Logout All') }}
                            </v-btn>
                        </th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr :key="itemIdx"
                        v-for="itemIdx in (loadingSession && (!sessions || sessions.length < 1) ? [ 1, 2, 3 ] : [])">
                        <td class="px-0" colspan="4">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>

                    <tr :key="session.tokenId"
                        v-for="session in sessions">
                        <td class="text-sm">
                            <v-icon start :icon="session.icon"/>
                            {{ tt(session.deviceName) }}
                        </td>
                        <td class="text-sm">{{ session.deviceInfo }}</td>
                        <td class="text-sm">{{ session.lastSeenDateTime }}</td>
                        <td class="text-sm text-right">
                            <v-btn density="comfortable" color="error" variant="tonal"
                                   :disabled="session.isCurrent || loadingSession"
                                   @click="revokeSession(session)">
                                {{ tt('Log Out') }}
                            </v-btn>
                        </td>
                    </tr>
                    </tbody>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <unlink-third-party-login-dialog ref="unlinkThirdPartyLoginDialog" />
    <user-generate-token-dialog ref="generateTokenDialog" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import { VTextField } from 'vuetify/components/VTextField';
import UnlinkThirdPartyLoginDialog from '@/views/desktop/user/settings/dialogs/UnlinkThirdPartyLoginDialog.vue';
import UserGenerateTokenDialog from '@/views/desktop/user/settings/dialogs/UserGenerateTokenDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserExternalAuthStore } from '@/stores/userExternalAuth.ts';
import { useTokensStore } from '@/stores/token.ts';

import { itemAndIndex, reversedItemAndIndex } from '@/core/base.ts';
import { KnownErrorCode } from '@/consts/api.ts';

import { type UserExternalAuthInfoResponse } from '@/models/user_external_auth.ts';
import { type TokenInfoResponse, SessionDeviceType, SessionInfo } from '@/models/token.ts';

import { isEquals } from '@/lib/common.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { parseSessionInfo } from '@/lib/session.ts';
import {
    isAPITokenEnabled,
    isOAuth2Enabled,
    getOAuth2Provider,
    getOIDCCustomDisplayNames,
    isMCPServerEnabled
} from '@/lib/server_settings.ts';
import { generateRandomUUID } from '@/lib/misc.ts';

import {
    mdiRefresh,
    mdiLinkVariant,
    mdiGithub,
    mdiCellphone,
    mdiTablet,
    mdiWatch,
    mdiTelevision,
    mdiConsole,
    mdiCreationOutline,
    mdiDevices
} from '@mdi/js';

class DesktopPageLinkedThirdPartyLogin {
    public readonly externalAuthType: string;
    public readonly icon: string;
    public readonly displayName: string;
    public readonly linked: boolean;
    public readonly externalUsername: string;
    public readonly createdAt: string;

    public constructor(externalAuthInfoResponse: UserExternalAuthInfoResponse) {
        this.externalAuthType = externalAuthInfoResponse.externalAuthType;
        this.linked = externalAuthInfoResponse.linked;
        this.externalUsername = externalAuthInfoResponse.externalUsername ? externalAuthInfoResponse.externalUsername : '-';
        this.createdAt = externalAuthInfoResponse.createdAt ? formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(externalAuthInfoResponse.createdAt)) : '-';

        if (externalAuthInfoResponse.externalAuthCategory === 'oauth2') {
            this.displayName = getLocalizedOAuth2ProviderName(externalAuthInfoResponse.externalAuthType, getOIDCCustomDisplayNames());

            if (externalAuthInfoResponse.externalAuthType === 'github') {
                this.icon = mdiGithub;
            } else {
                this.icon = mdiLinkVariant;
            }
        } else {
            this.displayName = externalAuthInfoResponse.externalAuthType;
            this.icon = mdiLinkVariant;
        }
    }
}

class DesktopPageSessionInfo extends SessionInfo {
    public readonly icon: string;
    public readonly lastSeenDateTime: string;

    public constructor(sessionInfo: SessionInfo) {
        super(sessionInfo.tokenId, sessionInfo.isCurrent, sessionInfo.deviceType, sessionInfo.deviceInfo, sessionInfo.deviceName, sessionInfo.lastSeen);
        this.icon = this.getTokenIcon(sessionInfo.deviceType);
        this.lastSeenDateTime = sessionInfo.lastSeen ? formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(sessionInfo.lastSeen)) : '-';
    }

    private getTokenIcon(deviceType: SessionDeviceType): string {
        if (deviceType === SessionDeviceType.Phone) {
            return mdiCellphone;
        } else if (deviceType === SessionDeviceType.Wearable) {
            return mdiWatch;
        } else if (deviceType === SessionDeviceType.Tablet) {
            return mdiTablet;
        } else if (deviceType === SessionDeviceType.TV) {
            return mdiTelevision;
        } else if (deviceType === SessionDeviceType.Api) {
            return mdiConsole;
        } else if (deviceType === SessionDeviceType.MCP) {
            return mdiCreationOutline;
        } else {
            return mdiDevices;
        }
    }
}

type UnlinkThirdPartyLoginDialogType = InstanceType<typeof UnlinkThirdPartyLoginDialog>;
type UserGenerateTokenDialogType = InstanceType<typeof UserGenerateTokenDialog>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    formatDateTimeToLongDateTime,
    getLocalizedOAuth2ProviderName,
    setLanguage
} = useI18n();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const userExternalAuthStore = useUserExternalAuthStore();
const tokensStore = useTokensStore();

const newPasswordInput = useTemplateRef<VTextField>('newPasswordInput');
const confirmPasswordInput = useTemplateRef<VTextField>('confirmPasswordInput');
const unlinkThirdPartyLoginDialog = useTemplateRef<UnlinkThirdPartyLoginDialogType>('unlinkThirdPartyLoginDialog');
const generateTokenDialog = useTemplateRef<UserGenerateTokenDialogType>('generateTokenDialog');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const externalAuths = ref<UserExternalAuthInfoResponse[]>([]);
const tokens = ref<TokenInfoResponse[]>([]);
const currentPassword = ref<string>('');
const newPassword = ref<string>('');
const confirmPassword = ref<string>('');
const updatingPassword = ref<boolean>(false);
const loadingExternalAuth = ref<boolean>(true);
const loadingSession = ref<boolean>(true);
const loggingInByOAuth2 = ref<boolean>(false);
const oauth2ClientSessionId = ref<string>(generateRandomUUID());

const oauth2LinkUrl = computed<string>(() => rootStore.generateOAuth2LinkUrl('desktop', oauth2ClientSessionId.value));

const thirdPartyLoginList = computed<DesktopPageLinkedThirdPartyLogin[]>(() => {
    const ret: DesktopPageLinkedThirdPartyLogin[] = [];

    if (!externalAuths.value) {
        return ret;
    }

    for (const externalAuth of externalAuths.value) {
        ret.push(new DesktopPageLinkedThirdPartyLogin(externalAuth));
    }

    return ret;
});

const sessions = computed<DesktopPageSessionInfo[]>(() => {
    const sessions: DesktopPageSessionInfo[] = [];

    if (!tokens.value) {
        return sessions;
    }

    for (const token of tokens.value) {
        const sessionInfo = parseSessionInfo(token);
        sessions.push(new DesktopPageSessionInfo(sessionInfo));
    }

    return sessions;
});

const inputProblemMessage = computed<string | null>(() => {
    if (!newPassword.value && !confirmPassword.value) {
        return 'Nothing has been modified';
    } else if (!newPassword.value && confirmPassword.value) {
        return 'New password cannot be blank';
    } else if (newPassword.value && !confirmPassword.value) {
        return 'Password confirmation cannot be blank';
    } else if (newPassword.value && confirmPassword.value && newPassword.value !== confirmPassword.value) {
        return 'Password and password confirmation do not match';
    } else {
        return null;
    }
});

function init(): void {
    reloadExternalAuth(true);
    reloadSessions(true);
}

function updatePassword(): void {
    const problemMessage = inputProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    updatingPassword.value = true;

    rootStore.updateUserProfile({
        password: newPassword.value,
        oldPassword: currentPassword.value
    }).then(response => {
        updatingPassword.value = false;
        currentPassword.value = '';
        newPassword.value = '';
        confirmPassword.value = '';

        if (response.user) {
            const localeDefaultSettings = setLanguage(response.user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }

        reloadSessions(true);

        snackbar.value?.showMessage('Your profile has been successfully updated');
    }).catch(error => {
        updatingPassword.value = false;
        currentPassword.value = '';

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function reloadExternalAuth(silent?: boolean): void {
    if (!isOAuth2Enabled()) {
        return;
    }

    loadingExternalAuth.value = true;

    userExternalAuthStore.getExternalAuths().then(response => {
        if (!silent) {
            if (isEquals(externalAuths.value, response)) {
                snackbar.value?.showMessage('Third-party login list is up to date');
            } else {
                snackbar.value?.showMessage('Third-party login list has been updated');
            }
        }

        externalAuths.value = response;
        loadingExternalAuth.value = false;
    }).catch(error => {
        loadingExternalAuth.value = false;

        if (error.error && error.error.errorCode === KnownErrorCode.ApiNotFound) {
            externalAuths.value = [];
        } else if (!error.processed) {
            externalAuths.value = [];
            snackbar.value?.showError(error);
        }
    });
}

function unlinkExternalAuth(thirdPartyLogin: DesktopPageLinkedThirdPartyLogin): void {
    if (!isOAuth2Enabled()) {
        return;
    }

    unlinkThirdPartyLoginDialog.value?.open(thirdPartyLogin.externalAuthType).then(() => {
        reloadExternalAuth(true);
    });
}

function generateToken(): void {
    generateTokenDialog.value?.open().then(() => {
        reloadSessions(true);
    });
}

function reloadSessions(silent?: boolean): void {
    loadingSession.value = true;

    tokensStore.getAllTokens().then(response => {
        if (!silent) {
            if (isEquals(tokens.value, response)) {
                snackbar.value?.showMessage('Session list is up to date');
            } else {
                snackbar.value?.showMessage('Session list has been updated');
            }
        }

        tokens.value = response;
        loadingSession.value = false;
    }).catch(error => {
        loadingSession.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function revokeSession(session: SessionInfo): void {
    confirmDialog.value?.open('Are you sure you want to logout from this session?').then(() => {
        loadingSession.value = true;

        tokensStore.revokeToken({
            tokenId: session.tokenId
        }).then(() => {
            loadingSession.value = false;

            for (const [token, index] of itemAndIndex(tokens.value)) {
                if (token.tokenId === session.tokenId) {
                    tokens.value.splice(index, 1);
                }
            }
        }).catch(error => {
            loadingSession.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function revokeAllSessions(): void {
    if (sessions.value.length < 2) {
        return;
    }

    confirmDialog.value?.open('Are you sure you want to logout all other sessions?').then(() => {
        loadingSession.value = true;

        tokensStore.revokeAllTokens().then(() => {
            loadingSession.value = false;

            for (const [token, index] of reversedItemAndIndex(tokens.value)) {
                if (!token.isCurrent) {
                    tokens.value.splice(index, 1);
                }
            }

            snackbar.value?.showMessage('You have logged out all other sessions');
        }).catch(error => {
            loadingSession.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

init();
</script>
