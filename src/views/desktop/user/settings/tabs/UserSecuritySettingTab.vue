<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': updatingPassword }" :title="$t('Modify Password')">
                <v-form>
                    <v-card-text class="pt-0">
                        <span class="text-body-1">{{ $t('After the password is changed, other devices will be logged out, please log in again on other devices by using the new password.') }}</span>
                    </v-card-text>

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    ref="currentPasswordInput"
                                    autocomplete="current-password"
                                    clearable
                                    :disabled="updatingPassword"
                                    :label="$t('Current Password')"
                                    :placeholder="$t('Current Password')"
                                    :type="isCurrentPasswordVisible ? 'text' : 'password'"
                                    :append-inner-icon="isCurrentPasswordVisible ? icons.eyeSlash : icons.eye"
                                    v-model="currentPassword"
                                    @click:append-inner="isCurrentPasswordVisible = !isCurrentPasswordVisible"
                                    @keyup.enter="$refs.newPasswordInput.focus()"
                                />
                            </v-col>
                        </v-row>

                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    ref="newPasswordInput"
                                    autocomplete="new-password"
                                    clearable
                                    :disabled="updatingPassword"
                                    :label="$t('New Password')"
                                    :placeholder="$t('New Password')"
                                    :type="isNewPasswordVisible ? 'text' : 'password'"
                                    :append-inner-icon="isNewPasswordVisible ? icons.eyeSlash : icons.eye"
                                    v-model="newPassword"
                                    @click:append-inner="isNewPasswordVisible = !isNewPasswordVisible"
                                    @keyup.enter="$refs.confirmPasswordInput.focus()"
                                />
                            </v-col>
                        </v-row>

                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    ref="confirmPasswordInput"
                                    clearable
                                    :disabled="updatingPassword"
                                    :type="isConfirmPasswordVisible ? 'text' : 'password'"
                                    :label="$t('Confirmation Password')"
                                    :placeholder="$t('Re-enter the password')"
                                    :append-inner-icon="isConfirmPasswordVisible ? icons.eyeSlash : icons.eye"
                                    v-model="confirmPassword"
                                    @click:append-inner="isConfirmPasswordVisible = !isConfirmPasswordVisible"
                                    @keyup.enter="updatePassword"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn :disabled="!currentPassword || !newPassword || !confirmPassword || updatingPassword" @click="updatePassword">
                            {{ $t('Save changes') }}
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
                        <span>{{ $t('Device & Sessions') }}</span>
                        <v-btn density="compact" color="default" variant="text"
                               class="ml-2" :icon="true"
                               v-if="!loadingSession" @click="reloadSessions">
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-progress-circular indeterminate size="20" class="ml-3" v-if="loadingSession"></v-progress-circular>
                    </div>
                </template>

                <v-table class="table-striped text-no-wrap" :hover="!loadingSession">
                    <thead>
                    <tr>
                        <th class="text-uppercase">{{ $t('Type') }}</th>
                        <th class="text-uppercase">{{ $t('Device Info') }}</th>
                        <th class="text-uppercase">{{ $t('Last Activity Time') }}</th>
                        <th class="text-uppercase text-right">
                            <v-btn density="comfortable" color="error" variant="tonal"
                                   :disabled="sessions.length < 2 || loadingSession"
                                   @click="revokeAllSessions">
                                {{ $t('Logout All') }}
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
                            {{ session.deviceType }}
                        </td>
                        <td class="text-sm">{{ session.deviceInfo }}</td>
                        <td class="text-sm">{{ session.createdAt }}</td>
                        <td class="text-sm text-right">
                            <v-btn density="comfortable" color="error" variant="tonal"
                                   :disabled="session.isCurrent || loadingSession"
                                   @click="revokeSession(session)">
                                {{ $t('Log Out') }}
                            </v-btn>
                        </td>
                    </tr>
                    </tbody>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';

import { isEquals } from '@/lib/common.js';
import { parseDeviceInfo, parseUserAgent } from '@/lib/misc.js';

import {
    mdiRefresh,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiCellphone,
    mdiTablet,
    mdiWatch,
    mdiTelevision,
    mdiDevices
} from '@mdi/js';

export default {
    data() {
        return {
            tokens: [],
            currentPassword: '',
            newPassword: '',
            confirmPassword: '',
            isCurrentPasswordVisible: false,
            isNewPasswordVisible: false,
            isConfirmPasswordVisible: false,
            updatingPassword: false,
            loadingSession: true,
            icons: {
                refresh: mdiRefresh,
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore),
        inputProblemMessage() {
            if (!this.currentPassword) {
                return 'Current password cannot be empty';
            } else if (!this.newPassword && !this.confirmPassword) {
                return 'Nothing has been modified';
            } else if (!this.newPassword && this.confirmPassword) {
                return 'New password cannot be empty';
            } else if (this.newPassword && !this.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else if (this.newPassword && this.confirmPassword && this.newPassword !== this.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else {
                return null;
            }
        },
        sessions() {
            if (!this.tokens) {
                return this.tokens;
            }

            const sessions = [];

            for (let i = 0; i < this.tokens.length; i++) {
                const token = this.tokens[i];

                sessions.push({
                    tokenId: token.tokenId,
                    isCurrent: token.isCurrent,
                    deviceType: this.$t(token.isCurrent ? 'Current' : 'Other Device'),
                    deviceInfo: parseDeviceInfo(token.userAgent),
                    icon: this.getTokenIcon(token),
                    createdAt: this.$locale.formatUnixTimeToLongDateTime(this.userStore, token.createdAt)
                });
            }

            return sessions;
        }
    },
    created() {
        const self = this;

        self.loadingSession = true;

        self.tokensStore.getAllTokens().then(tokens => {
            self.tokens = tokens;
            self.loadingSession = false;
        }).catch(error => {
            self.loadingSession = false;

            if (!error.processed) {
                self.$refs.snackbar.showError(error);
            }
        });
    },
    methods: {
        updatePassword() {
            const self = this;

            const problemMessage = self.inputProblemMessage;

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            self.updatingPassword = true;
            self.isCurrentPasswordVisible = false;
            self.isNewPasswordVisible = false;
            self.isConfirmPasswordVisible = false;

            self.rootStore.updateUserProfile({
                profile: {
                    password: self.newPassword,
                    confirmPassword: self.confirmPassword
                },
                currentPassword: self.currentPassword
            }).then(response => {
                self.updatingPassword = false;
                self.currentPassword = '';
                self.newPassword = '';
                self.confirmPassword = '';

                if (response.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                self.reloadSessions();

                self.$refs.snackbar.showMessage('Your profile has been successfully updated');
            }).catch(error => {
                self.updatingPassword = false;
                self.currentPassword = '';

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        reloadSessions() {
            const self = this;

            self.loadingSession = true;

            self.tokensStore.getAllTokens().then(tokens => {
                if (isEquals(self.tokens, tokens)) {
                    self.$refs.snackbar.showMessage('Session list is up to date');
                } else {
                    self.$refs.snackbar.showMessage('Session list has been updated');
                }

                self.tokens = tokens;
                self.loadingSession = false;
            }).catch(error => {
                self.loadingSession = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        revokeSession(session) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to logout from this session?').then(() => {
                self.loadingSession = true;

                self.tokensStore.revokeToken({
                    tokenId: session.tokenId
                }).then(() => {
                    self.loadingSession = false;

                    for (let i = 0; i < self.tokens.length; i++) {
                        if (self.tokens[i].tokenId === session.tokenId) {
                            self.tokens.splice(i, 1);
                        }
                    }
                }).catch(error => {
                    self.loadingSession = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        revokeAllSessions() {
            const self = this;

            if (self.tokens.length < 2) {
                return;
            }

            self.$refs.confirmDialog.open('Are you sure you want to logout all other sessions?').then(() => {
                self.loadingSession = true;

                self.tokensStore.revokeAllTokens().then(() => {
                    self.loadingSession = false;

                    for (let i = self.tokens.length - 1; i >= 0; i--) {
                        if (!self.tokens[i].isCurrent) {
                            self.tokens.splice(i, 1);
                        }
                    }

                    self.$refs.snackbar.showMessage('You have logged out all other sessions');
                }).catch(error => {
                    self.loadingSession = false

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        getTokenIcon(token) {
            const ua = parseUserAgent(token.userAgent);

            if (!ua || !ua.device) {
                return mdiDevices;
            }

            if (ua.device.type === 'mobile') {
                return mdiCellphone;
            } else if (ua.device.type === 'wearable') {
                return mdiWatch;
            } else if (ua.device.type === 'tablet') {
                return mdiTablet;
            } else if (ua.device.type === 'smarttv') {
                return mdiTelevision;
            } else {
                return mdiDevices;
            }
        }
    }
};
</script>
