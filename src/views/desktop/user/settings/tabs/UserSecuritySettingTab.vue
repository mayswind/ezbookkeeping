<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': updatingPassword }" :title="$t('Modify Password')">
                <v-form>
                    <v-card-text class="pt-0">
                        <span class="text-body-1">{{ $t('After changing the password, other devices will be logged out. Please use the new password to log in on other devices.') }}</span>
                    </v-card-text>

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    autocomplete="current-password"
                                    ref="currentPasswordInput"
                                    type="password"
                                    :disabled="updatingPassword"
                                    :label="$t('Current Password')"
                                    :placeholder="$t('Current Password')"
                                    v-model="currentPassword"
                                    @keyup.enter="$refs.newPasswordInput.focus()"
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
                                    :label="$t('New Password')"
                                    :placeholder="$t('New Password')"
                                    v-model="newPassword"
                                    @keyup.enter="$refs.confirmPasswordInput.focus()"
                                />
                            </v-col>
                        </v-row>

                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    ref="confirmPasswordInput"
                                    type="password"
                                    :disabled="updatingPassword"
                                    :label="$t('Confirm Password')"
                                    :placeholder="$t('Re-enter the password')"
                                    v-model="confirmPassword"
                                    @keyup.enter="updatePassword"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn :disabled="!currentPassword || !newPassword || !confirmPassword || updatingPassword" @click="updatePassword">
                            {{ $t('Save Changes') }}
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
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :loading="loadingSession" @click="reloadSessions">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-table class="table-striped text-no-wrap" :hover="!loadingSession">
                    <thead>
                    <tr>
                        <th>{{ $t('Type') }}</th>
                        <th>{{ $t('Device Info') }}</th>
                        <th>{{ $t('Last Activity Time') }}</th>
                        <th class="text-right">
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
                            {{ $t(session.isCurrent ? 'Current' : 'Other Device') }}
                        </td>
                        <td class="text-sm">{{ session.deviceInfo }}</td>
                        <td class="text-sm">{{ session.lastSeenDateTime }}</td>
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

import { isEquals } from '@/lib/common.ts';
import { parseSessionInfo } from '@/lib/session.ts';

import {
    mdiRefresh,
    mdiCellphone,
    mdiTablet,
    mdiWatch,
    mdiTelevision,
    mdiConsole,
    mdiDevices
} from '@mdi/js';

export default {
    data() {
        return {
            tokens: [],
            currentPassword: '',
            newPassword: '',
            confirmPassword: '',
            updatingPassword: false,
            loadingSession: true,
            icons: {
                refresh: mdiRefresh
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore),
        inputProblemMessage() {
            if (!this.currentPassword) {
                return 'Current password cannot be blank';
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
        sessions() {
            if (!this.tokens) {
                return this.tokens;
            }

            const sessions = [];

            for (let i = 0; i < this.tokens.length; i++) {
                const token = this.tokens[i];
                const sessionInfo = parseSessionInfo(token);
                sessionInfo.icon = this.getTokenIcon(sessionInfo.deviceType);
                sessionInfo.lastSeenDateTime = sessionInfo.lastSeen ? this.$locale.formatUnixTimeToLongDateTime(this.userStore, sessionInfo.lastSeen) : '-';
                sessions.push(sessionInfo);
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
        getTokenIcon(deviceType) {
            if (deviceType === 'phone') {
                return mdiCellphone;
            } else if (deviceType === 'wearable') {
                return mdiWatch;
            } else if (deviceType === 'tablet') {
                return mdiTablet;
            } else if (deviceType === 'tv') {
                return mdiTelevision;
            } else if (deviceType === 'cli') {
                return mdiConsole;
            } else {
                return mdiDevices;
            }
        }
    }
};
</script>
