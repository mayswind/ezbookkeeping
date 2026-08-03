<template>
    <v-dialog width="1000" :persistent="isAccountModified" v-model="showState">
        <two-column-dialog-layout :disabled="loading || submitting" :loading="loading"
                                  :title="tt(title)" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #content-left-column>
                <div class="px-4">
                    <v-tabs class="v-tabs-pill" direction="vertical" :class="{ 'readonly': !!editAccountId }"
                            :disabled="loading || submitting" v-model="account.type">
                        <v-tab :key="accountType.type" :value="accountType.type" :disabled="!!editAccountId && accountType.type !== account.type"
                               v-for="accountType in allAccountTypes">
                            <span>{{ accountType.displayName }}</span>
                        </v-tab>
                    </v-tabs>
                </div>
                <v-divider class="my-2"/>
                <div class="px-4" v-if="account.type === AccountType.SingleAccount.type">
                    <v-tabs direction="vertical" :disabled="loading || submitting" :model-value="-1">
                        <v-tab :value="-1">
                            <span>{{ tt('Basic Information') }}</span>
                        </v-tab>
                    </v-tabs>
                </div>
                <div class="px-4" v-else-if="account.type === AccountType.MultiSubAccounts.type">
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="currentAccountIndex">
                        <v-tab :value="-1">
                            <span>{{ tt('Main Account') }}</span>
                        </v-tab>
                        <template v-if="account.type === AccountType.MultiSubAccounts.type">
                            <v-tab :key="idx" :value="idx" v-for="(subAccount, idx) in subAccounts">
                                <span>{{ tt('Sub Account') + ' #' + (idx + 1) }}</span>
                                <v-btn class="ms-2" color="error" size="24" variant="text"
                                       :icon="mdiDeleteOutline"
                                       @click="removeSubAccount(subAccount)"></v-btn>
                            </v-tab>
                        </template>
                    </v-tabs>
                    <div class="w-100">
                        <v-btn class="mt-2 w-100" color="primary" variant="text" density="comfortable"
                               :disabled="loading || submitting" :prepend-icon="mdiPlus"
                               @click="addSubAccount">{{ tt('Add Sub-account') }}</v-btn>
                    </div>
                </div>
            </template>

            <template #content-right-column>
                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container"
                          v-model="activeTab">
                    <v-window-item value="account">
                        <v-form class="my-4">
                            <v-row>
                                <v-col cols="12" md="12">
                                    <v-text-field
                                        type="text"
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="currentAccountIndex < 0 ? tt('Account Name') : tt('Sub-account Name')"
                                        :placeholder="currentAccountIndex < 0 ? tt('Your account name') : tt('Your sub-account name')"
                                        v-model="selectedAccount.name"
                                    />
                                </v-col>
                                <v-col cols="12" md="12" v-if="account.type === AccountType.SingleAccount.type || currentAccountIndex < 0">
                                    <v-select
                                        item-title="displayName"
                                        item-value="type"
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="tt('Account Category')"
                                        :placeholder="tt('Account Category')"
                                        :items="allAccountCategories"
                                        :no-data-text="tt('No results')"
                                        v-model="selectedAccount.category"
                                    >
                                        <template #item="{ props, item }">
                                            <v-list-item :value="item.value" v-bind="props">
                                                <template #title>
                                                    <v-list-item-title>
                                                        <div class="d-flex align-center">
                                                            <ItemIcon icon-type="account"
                                                                      :icon-id="item.raw.defaultAccountIconId"
                                                                      v-if="item.raw" />
                                                            <span class="ms-2">{{ item.title }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </template>
                                            </v-list-item>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="6">
                                    <icon-select icon-type="account"
                                                 :all-icon-infos="ALL_ACCOUNT_ICONS"
                                                 :label="currentAccountIndex < 0 ? tt('Account Icon') : tt('Sub-account Icon')"
                                                 :color="selectedAccount.color"
                                                 :disabled="loading || submitting"
                                                 v-model="selectedAccount.icon" />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <color-select :all-color-infos="ALL_ACCOUNT_COLORS"
                                                  :label="currentAccountIndex < 0 ? tt('Account Color') : tt('Sub-account Color')"
                                                  :disabled="loading || submitting"
                                                  v-model="selectedAccount.color" />
                                </v-col>
                                <v-col cols="12" :md="currentAccountIndex < 0 && isAccountSupportCreditCardStatementDate ? 6 : 12" v-if="account.type === AccountType.SingleAccount.type || currentAccountIndex >= 0">
                                    <currency-select :disabled="loading || submitting || (!!editAccountId && !isNewAccount(selectedAccount))"
                                                     :label="tt('Currency')"
                                                     :placeholder="tt('Currency')"
                                                     v-model="selectedAccount.currency" />
                                </v-col>
                                <v-col cols="12" :md="account.type === AccountType.SingleAccount.type || currentAccountIndex >= 0 ? 6 : 12" v-if="currentAccountIndex < 0 && isAccountSupportCreditCardStatementDate">
                                    <v-autocomplete
                                        item-title="displayName"
                                        item-value="type"
                                        auto-select-first
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="tt('Statement Date')"
                                        :placeholder="tt('Statement Date')"
                                        :items="allAvailableMonthDays"
                                        :no-data-text="tt('No results')"
                                        v-model="account.creditCardStatementDate"
                                    ></v-autocomplete>
                                </v-col>
                                <v-col cols="12" :md="((canShowBalanceTime && selectedAccount.numericBalance) || canShowLastReconciledTime) ? 6 : 12"
                                       v-if="account.type === AccountType.SingleAccount.type || currentAccountIndex >= 0">
                                    <amount-input :disabled="loading || submitting || (!!editAccountId && !isNewAccount(selectedAccount))"
                                                  :persistent-placeholder="true"
                                                  :currency="selectedAccount.currency"
                                                  :show-currency="true"
                                                  :flip-negative="account.isLiability"
                                                  :label="accountAmountTitle"
                                                  :placeholder="accountAmountTitle"
                                                  v-model="selectedAccount.numericBalance"/>
                                </v-col>
                                <v-col cols="12" md="6" v-show="selectedAccount.numericBalance" v-if="canShowBalanceTime">
                                    <date-time-select
                                        :disabled="loading || submitting"
                                        :label="tt('Balance Time')"
                                        :timezone-utc-offset="getDefaultTimezoneOffsetMinutes(selectedAccount.balanceTime)"
                                        :model-value="selectedAccount.balanceTime"
                                        @update:model-value="updateAccountBalanceTime(selectedAccount, $event)"
                                        @error="onShowDateTimeError" />
                                </v-col>
                                <v-col cols="12" md="6" v-if="canShowLastReconciledTime">
                                    <date-time-select
                                        :disabled="loading || submitting"
                                        :clearable="true"
                                        :label="tt('Last Reconciled Time')"
                                        :timezone-utc-offset="getDefaultTimezoneOffsetMinutes(selectedAccount.lastReconciledTime)"
                                        :model-value="selectedAccount.lastReconciledTime ?? getCurrentUnixTime()"
                                        :empty-value="!selectedAccount.lastReconciledTime"
                                        @update:model-value="updateAccountLastReconciledTime(selectedAccount, $event)"
                                        @clear:model-value="selectedAccount.lastReconciledTime = undefined"
                                        @error="onShowDateTimeError" />
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-textarea
                                        type="text"
                                        persistent-placeholder
                                        rows="3"
                                        :disabled="loading || submitting"
                                        :label="tt('Description')"
                                        :placeholder="currentAccountIndex < 0 ? tt('Your account description (optional)') : tt('Your sub-account description (optional)')"
                                        v-model="selectedAccount.comment"
                                    />
                                </v-col>
                                <v-col class="py-0" cols="12" md="12" v-if="editAccountId && !isNewAccount(selectedAccount)">
                                    <v-switch :disabled="loading || submitting"
                                              :label="tt('Visible')" v-model="selectedAccount.visible"/>
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-window-item>
                </v-window>
            </template>

            <template #footer>
                <v-spacer/>
                <v-tooltip :disabled="!inputIsEmpty" :text="inputEmptyProblemMessage ? tt(inputEmptyProblemMessage) : ''">
                    <template v-slot:activator="{ props }">
                        <div v-bind="props" class="d-inline-block">
                            <v-btn :disabled="inputIsEmpty || loading || submitting" @click="save">
                                {{ tt(saveButtonTitle) }}
                                <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                            </v-btn>
                        </div>
                    </template>
                </v-tooltip>
            </template>
        </two-column-dialog-layout>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountEditPageBase } from '@/views/base/accounts/AccountEditPageBase.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { itemAndIndex } from '@/core/base.ts';
import { AccountType, AccountCategory } from '@/core/account.ts';
import { ALL_ACCOUNT_ICONS } from '@/consts/icon.ts';
import { ALL_ACCOUNT_COLORS } from '@/consts/color.ts';
import { Account } from '@/models/account.ts';

import { isNumber, isEquals } from '@/lib/common.ts';
import { getCurrentUnixTime } from '@/lib/datetime.ts';
import { generateRandomUUID } from '@/lib/misc.ts';

import {
    mdiPlus,
    mdiDeleteOutline
} from '@mdi/js';

interface AccountEditResponse {
    message: string;
}

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const {
    defaultAccountCategory,
    editAccountId,
    clientSessionId,
    loading,
    submitting,
    account,
    subAccounts,
    useLastReconciledTime,
    title,
    saveButtonTitle,
    inputEmptyProblemMessage,
    inputIsEmpty,
    allAccountCategories,
    allAccountTypes,
    allAvailableMonthDays,
    isAccountSupportCreditCardStatementDate,
    getCurrentUnixTimeForNewAccount,
    getDefaultTimezoneOffsetMinutes,
    updateAccountBalanceTime,
    updateAccountLastReconciledTime,
    isNewAccount,
    addSubAccount,
    setAccount
} = useAccountEditPageBase();

const userStore = useUserStore();
const accountsStore = useAccountsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: ((value: AccountEditResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const activeTab = ref<string>('account');
const initAccountCategory = ref<AccountCategory>(defaultAccountCategory);
const initAccount = ref<Account | null>(null);
const currentAccountIndex = ref<number>(-1);

const canShowBalanceTime = computed<boolean>(() => (!editAccountId.value || isNewAccount(selectedAccount.value)) && (account.value.type === AccountType.SingleAccount.type || currentAccountIndex.value >= 0));
const canShowLastReconciledTime = computed<boolean>(() => useLastReconciledTime.value && !!editAccountId.value && !isNewAccount(selectedAccount.value) && (account.value.type === AccountType.SingleAccount.type || currentAccountIndex.value >= 0));

const selectedAccount = computed<Account>(() => {
    if (currentAccountIndex.value < 0) {
        return account.value;
    }

    return subAccounts.value[currentAccountIndex.value] as Account;
});

const accountAmountTitle = computed<string>(() => {
    if (currentAccountIndex.value < 0) {
        return account.value.isLiability ? tt('Account Outstanding Balance') : tt('Account Balance');
    } else {
        return account.value.isLiability ? tt('Sub-account Outstanding Balance') : tt('Sub-account Balance');
    }
});

const isAccountModified = computed<boolean>(() => {
    if (!editAccountId.value) { // Add
        return !!initAccount.value && !isEquals(account.value.toCreateRequest(clientSessionId.value, subAccounts.value), initAccount.value.toCreateRequest(clientSessionId.value, initAccount.value.subAccounts));
    } else { // Edit
        return !!initAccount.value && !isEquals(account.value.toModifyRequest(clientSessionId.value, subAccounts.value), initAccount.value.toModifyRequest(clientSessionId.value, initAccount.value.subAccounts));
    }
});

function open(options?: { id?: string, currentAccount?: Account, category?: number }): Promise<AccountEditResponse> {
    showState.value = true;
    loading.value = true;
    submitting.value = false;

    if (isNumber(options?.category) && AccountCategory.valueOf(options.category)) {
        initAccountCategory.value = AccountCategory.valueOf(options.category)!;
    }

    initAccount.value = Account.createNewAccount(initAccountCategory.value, userStore.currentUserDefaultCurrency, getCurrentUnixTimeForNewAccount());
    account.value.fillFrom(initAccount.value);
    subAccounts.value = [];
    currentAccountIndex.value = -1;
    clientSessionId.value = generateRandomUUID();

    if (options && options.id) {
        if (options.currentAccount) {
            setAccount(options.currentAccount);
        }

        editAccountId.value = options.id;
        accountsStore.getAccount({
            accountId: editAccountId.value
        }).then(response => {
            setAccount(response);
            initAccount.value = Account.of(response);
            loading.value = false;
        }).catch(error => {
            loading.value = false;
            showState.value = false;

            if (!error.processed) {
                if (rejectFunc) {
                    rejectFunc(error);
                }
            }
        });
    } else {
        if (options && isNumber(options.category)) {
            initAccount.value.category = options.category;
            initAccount.value.setSuitableIcon(1, options.category);

            account.value.category = initAccount.value.category;
            account.value.icon = initAccount.value.icon;
        }

        editAccountId.value = null;
        loading.value = false;
    }

    return new Promise<AccountEditResponse>((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function save(): void {
    const problemMessage = inputEmptyProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    submitting.value = true;

    accountsStore.saveAccount({
        account: account.value,
        subAccounts: subAccounts.value,
        isEdit: !!editAccountId.value,
        clientSessionId: clientSessionId.value
    }).then(() => {
        submitting.value = false;

        let message = 'You have saved this account';

        if (!editAccountId.value) {
            message = 'You have added a new account';
        }

        resolveFunc?.({ message });
        showState.value = false;
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function removeSubAccount(currentSubAccount: Account): void {
    confirmDialog.value?.open('Are you sure you want to remove this sub-account?').then(() => {
        for (const [subAccount, index] of itemAndIndex(subAccounts.value)) {
            if (subAccount === currentSubAccount) {
                subAccounts.value.splice(index, 1);

                if (currentAccountIndex.value >= subAccounts.value.length) {
                    currentAccountIndex.value = subAccounts.value.length - 1;
                }
            }
        }
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

function onShowDateTimeError(error: string): void {
    snackbar.value?.showError(error);
}

watch(() => account.value.type, () => {
    if (subAccounts.value.length < 1) {
        addSubAccount();
    }
});

defineExpose({
    open
});
</script>
