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
                        <v-btn :disabled="!currentPassword || !newPassword || !confirmPassword || updatingPassword" @click="updatePassword">
                            {{ tt('Save Changes') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="updatingPassword"></v-progress-circular>
                        </v-btn>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :class="{ 'disabled': loadingSession }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ tt('Device & Sessions') }}</span>
                        <v-btn class="ml-3" density="compact" color="default" variant="outlined"
                               @click="generateMCPToken" v-if="isMCPServerEnabled()">{{ tt('Generate MCP token') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :loading="loadingSession" @click="reloadSessions(false)">
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
                            {{ session.deviceType === 'mcp' ? 'MCP' : (tt(session.isCurrent ? 'Current' : 'Other Device')) }}
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

    <user-generate-m-c-p-token-dialog ref="generateMCPTokenDialog" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import { VTextField } from 'vuetify/components/VTextField';
import UserGenerateMCPTokenDialog from '@/views/desktop/user/settings/dialogs/UserGenerateMCPTokenDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useTokensStore } from '@/stores/token.ts';

import { type TokenInfoResponse, SessionInfo } from '@/models/token.ts';

import { isEquals } from '@/lib/common.ts';
import { parseSessionInfo } from '@/lib/session.ts';
import { isMCPServerEnabled } from '@/lib/server_settings.ts';

import {
    mdiRefresh,
    mdiCellphone,
    mdiTablet,
    mdiWatch,
    mdiTelevision,
    mdiCreationOutline,
    mdiConsole,
    mdiDevices
} from '@mdi/js';

class DesktopPageSessionInfo extends SessionInfo {
    public readonly icon: string;
    public readonly lastSeenDateTime: string;

    public constructor(sessionInfo: SessionInfo) {
        super(sessionInfo.tokenId, sessionInfo.isCurrent, sessionInfo.deviceType, sessionInfo.deviceInfo, sessionInfo.createdByCli, sessionInfo.lastSeen);
        this.icon = getTokenIcon(sessionInfo.deviceType);
        this.lastSeenDateTime = sessionInfo.lastSeen ? formatUnixTimeToLongDateTime(sessionInfo.lastSeen) : '-';
    }
}

type UserGenerateMCPTokenDialogType = InstanceType<typeof UserGenerateMCPTokenDialog>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt, formatUnixTimeToLongDateTime, setLanguage } = useI18n();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const tokensStore = useTokensStore();

const newPasswordInput = useTemplateRef<VTextField>('newPasswordInput');
const confirmPasswordInput = useTemplateRef<VTextField>('confirmPasswordInput');
const generateMCPTokenDialog = useTemplateRef<UserGenerateMCPTokenDialogType>('generateMCPTokenDialog');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const tokens = ref<TokenInfoResponse[]>([]);
const currentPassword = ref<string>('');
const newPassword = ref<string>('');
const confirmPassword = ref<string>('');
const updatingPassword = ref<boolean>(false);
const loadingSession = ref<boolean>(true);

const sessions = computed<DesktopPageSessionInfo[]>(() => {
    const sessions: DesktopPageSessionInfo[] = [];

    if (!tokens.value) {
        return sessions;
    }

    for (let i = 0; i < tokens.value.length; i++) {
        const token = tokens.value[i];
        const sessionInfo = parseSessionInfo(token);
        sessions.push(new DesktopPageSessionInfo(sessionInfo));
    }

    return sessions;
});

const inputProblemMessage = computed<string | null>(() => {
    if (!currentPassword.value) {
        return 'Current password cannot be blank';
    } else if (!newPassword.value && !confirmPassword.value) {
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

function getTokenIcon(deviceType: string): string {
    if (deviceType === 'phone') {
        return mdiCellphone;
    } else if (deviceType === 'wearable') {
        return mdiWatch;
    } else if (deviceType === 'tablet') {
        return mdiTablet;
    } else if (deviceType === 'tv') {
        return mdiTelevision;
    } else if (deviceType === 'mcp') {
        return mdiCreationOutline;
    } else if (deviceType === 'cli') {
        return mdiConsole;
    } else {
        return mdiDevices;
    }
}

function init(): void {
    loadingSession.value = true;

    tokensStore.getAllTokens().then(response => {
        tokens.value = response;
        loadingSession.value = false;
    }).catch(error => {
        loadingSession.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
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

function generateMCPToken(): void {
    generateMCPTokenDialog.value?.open().then(() => {
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

            for (let i = 0; i < tokens.value.length; i++) {
                if (tokens.value[i].tokenId === session.tokenId) {
                    tokens.value.splice(i, 1);
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

            for (let i = tokens.value.length - 1; i >= 0; i--) {
                if (!tokens.value[i].isCurrent) {
                    tokens.value.splice(i, 1);
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
