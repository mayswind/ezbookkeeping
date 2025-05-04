<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading || saving }">
                <template #title>
                    <span>{{ tt('Basic Settings') }}</span>
                    <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text class="d-flex">
                    <v-avatar rounded="lg" variant="tonal" size="100" class="me-4 user-profile-avatar-icon"
                              :class="{ 'cursor-pointer': avatarProvider === 'internal', 'user-profile-avatar-icon-modifiable': avatarProvider === 'internal' }"
                              :color="currentUserAvatar ? 'rgba(0,0,0,0)' : 'primary'">
                        <v-img :src="currentUserAvatar" v-if="currentUserAvatar">
                            <template #placeholder>
                                <div class="d-flex align-center justify-center fill-height bg-light-primary">
                                    <v-icon color="primary" size="48" class="user-profile-avatar-placeholder" :icon="mdiAccount"/>
                                </div>
                            </template>
                        </v-img>
                        <v-icon size="48" class="user-profile-avatar-placeholder" :icon="mdiAccount" v-else-if="!currentUserAvatar"/>
                        <div class="avatar-edit-icon" v-if="avatarProvider === 'internal'">
                            <v-icon size="48" :icon="mdiAccountEditOutline"/>
                        </div>
                        <v-menu activator="parent" width="200" location="bottom" offset="14px" v-if="avatarProvider === 'internal'">
                            <v-list>
                                <v-list-item :disabled="saving" :title="tt('Update Avatar')" @click="showOpenAvatarDialog"></v-list-item>
                                <v-list-item :disabled="!currentUserAvatar || saving" :title="tt('Remove Avatar')" @click="removeAvatar"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-avatar>
                    <div class="d-flex flex-column justify-center gap-3">
                        <div class="d-flex text-body-1">
                            <span class="me-1">{{ tt('Username:') }}</span>
                            <v-skeleton-loader class="skeleton-no-margin" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                            <span v-if="!loading">{{ oldProfile.username }}</span>
                        </div>
                        <div class="d-flex text-body-1 align-center" style="height: 40px;">
                            <span v-if="!loading && emailVerified">{{ tt('Email address is verified') }}</span>
                            <span v-if="!loading && !emailVerified">{{ tt('Email address is not verified') }}</span>
                            <v-btn class="ml-2 px-2" size="small" variant="text" :disabled="loading || resending"
                                   @click="resendVerifyEmail" v-if="isUserVerifyEmailEnabled() && !loading && !emailVerified">
                                {{ tt('Resend Validation Email') }}
                                <v-progress-circular indeterminate size="18" class="ml-2" v-if="resending"></v-progress-circular>
                            </v-btn>
                            <v-skeleton-loader class="skeleton-no-margin mt-2 mb-1" type="text" style="width: 160px" :loading="true" v-if="loading"></v-skeleton-loader>
                        </div>
                    </div>
                </v-card-text>

                <v-divider />

                <v-form class="mt-6">
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    autocomplete="nickname"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Nickname')"
                                    :placeholder="tt('Your nickname')"
                                    v-model="newProfile.nickname"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="email"
                                    autocomplete="email"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('E-mail')"
                                    :placeholder="tt('Your email address')"
                                    v-model="newProfile.email"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <two-column-select primary-key-field="id" primary-value-field="category"
                                                   primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="account"
                                                   primary-sub-items-field="accounts"
                                                   :primary-title-i18n="true"
                                                   secondary-key-field="id" secondary-value-field="id"
                                                   secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                   :disabled="loading || saving || !allVisibleAccounts.length"
                                                   :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                   :label="tt('Default Account')"
                                                   :placeholder="tt('Default Account')"
                                                   :items="allVisibleCategorizedAccounts"
                                                   :no-item-text="tt('Unspecified')"
                                                   v-model="newProfile.defaultAccountId">
                                </two-column-select>
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Editable Transaction Range')"
                                    :placeholder="tt('Editable Transaction Range')"
                                    :items="allTransactionEditScopeTypes"
                                    v-model="newProfile.transactionEditScope"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <language-select :disabled="loading || saving"
                                                 :label="languageTitle"
                                                 :placeholder="languageTitle"
                                                 :include-system-default="true"
                                                 :use-model-value="true" v-model="newProfile.language" />
                            </v-col>

                            <v-col cols="12" md="6">
                                <currency-select :disabled="loading || saving"
                                                 :label="tt('Default Currency')"
                                                 :placeholder="tt('Default Currency')"
                                                 v-model="newProfile.defaultCurrency" />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('First Day of Week')"
                                    :placeholder="tt('First Day of Week')"
                                    :items="allWeekDays"
                                    v-model="newProfile.firstDayOfWeek"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <fiscal-year-start-select
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Fiscal Year Start Date')"
                                    :placeholder="tt('Fiscal Year Start Date')"
                                    v-model="newProfile.fiscalYearStart"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Long Date Format')"
                                    :placeholder="tt('Long Date Format')"
                                    :items="allLongDateFormats"
                                    v-model="newProfile.longDateFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Short Date Format')"
                                    :placeholder="tt('Short Date Format')"
                                    :items="allShortDateFormats"
                                    v-model="newProfile.shortDateFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Long Time Format')"
                                    :placeholder="tt('Long Time Format')"
                                    :items="allLongTimeFormats"
                                    v-model="newProfile.longTimeFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Short Time Format')"
                                    :placeholder="tt('Short Time Format')"
                                    :items="allShortTimeFormats"
                                    v-model="newProfile.shortTimeFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Fiscal Year Format')"
                                    :placeholder="tt('Fiscal Year Format')"
                                    :items="allFiscalYearFormats"
                                    v-model="newProfile.fiscalYearFormat"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Currency Display Mode')"
                                    :placeholder="tt('Currency Display Mode')"
                                    :items="allCurrencyDisplayTypes"
                                    v-model="newProfile.currencyDisplayType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Digit Grouping')"
                                    :placeholder="tt('Digit Grouping')"
                                    :items="allDigitGroupingTypes"
                                    v-model="newProfile.digitGrouping"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving || !supportDigitGroupingSymbol"
                                    :label="tt('Digit Grouping Symbol')"
                                    :placeholder="tt('Digit Grouping Symbol')"
                                    :items="allDigitGroupingSymbols"
                                    v-model="newProfile.digitGroupingSymbol"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Decimal Separator')"
                                    :placeholder="tt('Decimal Separator')"
                                    :items="allDecimalSeparators"
                                    v-model="newProfile.decimalSeparator"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Expense Amount Color')"
                                    :placeholder="tt('Expense Amount Color')"
                                    :items="allExpenseAmountColorTypes"
                                    v-model="newProfile.expenseAmountColor"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="tt('Income Amount Color')"
                                    :placeholder="tt('Income Amount Color')"
                                    :items="allIncomeAmountColorTypes"
                                    v-model="newProfile.incomeAmountColor"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn :disabled="inputIsNotChanged || inputIsInvalid || saving" @click="save">
                            {{ tt('Save Changes') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="saving"></v-progress-circular>
                        </v-btn>

                        <v-btn color="default" variant="tonal" @click="reset">
                            {{ tt('Reset') }}
                        </v-btn>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="avatarInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="updateAvatar($event)" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useUserProfilePageBase } from '@/views/base/users/UserProfilePageBase.ts';

import { useRootStore } from '@/stores/index.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';
import type { UserProfileResponse } from '@/models/user.ts';

import { generateRandomUUID } from '@/lib/misc.ts';
import { isUserVerifyEmailEnabled } from '@/lib/server_settings.ts';

import {
    mdiAccount,
    mdiAccountEditOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const {
    newProfile,
    oldProfile,
    emailVerified,
    loading,
    resending,
    saving,
    allVisibleAccounts,
    allVisibleCategorizedAccounts,
    allWeekDays,
    allLongDateFormats,
    allShortDateFormats,
    allLongTimeFormats,
    allShortTimeFormats,
    allFiscalYearFormats,
    allDecimalSeparators,
    allDigitGroupingSymbols,
    allDigitGroupingTypes,
    allCurrencyDisplayTypes,
    allExpenseAmountColorTypes,
    allIncomeAmountColorTypes,
    allTransactionEditScopeTypes,
    languageTitle,
    supportDigitGroupingSymbol,
    inputIsNotChangedProblemMessage,
    inputInvalidProblemMessage,
    langAndRegionInputInvalidProblemMessage,
    extendInputInvalidProblemMessage,
    inputIsNotChanged,
    inputIsInvalid,
    setCurrentUserProfile,
    reset,
    doAfterProfileUpdate
} = useUserProfilePageBase();

const rootStore = useRootStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const avatarInput = useTemplateRef<HTMLInputElement>('avatarInput');

const avatarUrl = ref<string>('');
const avatarProvider = ref<string | undefined>('');
const avatarNoCacheId = ref<string>('');

const currentUserAvatar = computed<string | null>(() => userStore.getUserAvatarUrl(avatarUrl.value, avatarNoCacheId.value));

function init(): void {
    loading.value = true;

    const promises = [
        accountsStore.loadAllAccounts({ force: false }),
        userStore.getCurrentUserProfile()
    ];

    Promise.all(promises).then(responses => {
        const profile = responses[1] as UserProfileResponse;
        setCurrentUserProfile(profile);
        avatarUrl.value = profile.avatar;
        avatarProvider.value = profile.avatarProvider;
        loading.value = false;
    }).catch(error => {
        oldProfile.value.nickname = '';
        oldProfile.value.email = '';
        newProfile.value.nickname = '';
        newProfile.value.email = '';
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function save(): void {
    const problemMessage = inputIsNotChangedProblemMessage.value || inputInvalidProblemMessage.value || extendInputInvalidProblemMessage.value || langAndRegionInputInvalidProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    saving.value = true;

    rootStore.updateUserProfile(newProfile.value.toProfileUpdateRequest()).then(response => {
        saving.value = false;

        doAfterProfileUpdate(response.user);
        snackbar.value?.showMessage('Your profile has been successfully updated');
    }).catch(error => {
        saving.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function updateAvatar(event: Event): void {
    if (!event || !event.target) {
        return;
    }

    const el = event.target as HTMLInputElement;

    if (!el.files || !el.files.length) {
        return;
    }

    const avatarFile = el.files[0];

    el.value = '';

    saving.value = true;

    userStore.updateUserAvatar({ avatarFile }).then(response => {
        saving.value = false;

        if (response) {
            avatarUrl.value = response.avatar;
            avatarProvider.value = response.avatarProvider;
            avatarNoCacheId.value = generateRandomUUID();
            setCurrentUserProfile(response);
        }

        snackbar.value?.showMessage('Your avatar has been successfully updated');
    }).catch(error => {
        saving.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function removeAvatar(): void {
    confirmDialog.value?.open('Are you sure you want to remove avatar?').then(() => {
        saving.value = true;

        userStore.removeUserAvatar().then(response => {
            saving.value = false;

            if (response) {
                avatarUrl.value = response.avatar;
                avatarProvider.value = response.avatarProvider;
                setCurrentUserProfile(response);
            }

            snackbar.value?.showMessage('Your profile has been successfully updated');
        }).catch(error => {
            saving.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function resendVerifyEmail(): void {
    resending.value = true;

    rootStore.resendVerifyEmailByLoginedUser().then(() => {
        resending.value = false;
        snackbar.value?.showMessage('Validation email has been sent');
    }).catch(error => {
        resending.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function showOpenAvatarDialog(): void {
    avatarInput.value?.click();
}

init();
</script>

<style>
.user-profile-avatar-icon .avatar-edit-icon {
    display: none;
    position: absolute;
    width: 100% !important;
    height: 100% !important;
    background-color: rgba(0, 0, 0, 0.4);
}

.user-profile-avatar-icon .avatar-edit-icon > i.v-icon {
    background-color: transparent;
    color: rgba(255, 255, 255, 0.8);
}

.user-profile-avatar-icon:hover .avatar-edit-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    vertical-align: middle;
}

.user-profile-avatar-icon-modifiable:hover .user-profile-avatar-placeholder {
    display: none;
}
</style>
