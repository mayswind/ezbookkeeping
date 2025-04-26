<template>
    <v-dialog :width="account.type === AccountType.MultiSubAccounts.type ? 1000 : 800" :persistent="isAccountModified" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ tt(title) }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           :disabled="loading || submitting || account.type !== AccountType.MultiSubAccounts.type">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiCreditCardPlusOutline"
                                             :title="tt('Add Sub-account')"
                                             @click="addSubAccount"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-md-4 pt-0">
                <div class="mb-4" v-if="account.type === AccountType.MultiSubAccounts.type">
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="currentAccountIndex">
                        <v-tab :value="-1">
                            <span>{{ tt('Main Account') }}</span>
                        </v-tab>
                        <template v-if="account.type === AccountType.MultiSubAccounts.type">
                            <v-tab :key="idx" :value="idx" v-for="(subAccount, idx) in subAccounts">
                                <span>{{ tt('Sub Account') + ' #' + (idx + 1) }}</span>
                                <v-btn class="ml-2" color="error" size="24" variant="text"
                                       :icon="mdiDeleteOutline"
                                       @click="removeSubAccount(subAccount)"></v-btn>
                            </v-tab>
                        </template>
                    </v-tabs>
                </div>

                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container"
                          :class="{ 'ml-md-5': account.type === AccountType.MultiSubAccounts.type }"
                          v-model="activeTab">
                    <v-window-item value="account">
                        <v-form class="mt-2">
                            <v-row>
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
                                                            <span class="ml-2">{{ item.title }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </template>
                                            </v-list-item>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="12" v-if="account.type === AccountType.SingleAccount.type || currentAccountIndex < 0">
                                    <v-select
                                        item-title="displayName"
                                        item-value="type"
                                        persistent-placeholder
                                        :disabled="loading || submitting || !!editAccountId"
                                        :label="tt('Account Type')"
                                        :placeholder="tt('Account Type')"
                                        :items="allAccountTypes"
                                        :no-data-text="tt('No results')"
                                        v-model="selectedAccount.type"
                                    />
                                </v-col>
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
                                        item-value="day"
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
                                <v-col cols="12" :md="(!editAccountId || isNewAccount(selectedAccount)) && selectedAccount.balance ? 6 : 12"
                                       v-if="account.type === AccountType.SingleAccount.type || currentAccountIndex >= 0">
                                    <amount-input :disabled="loading || submitting || (!!editAccountId && !isNewAccount(selectedAccount))"
                                                  :persistent-placeholder="true"
                                                  :currency="selectedAccount.currency"
                                                  :show-currency="true"
                                                  :flip-negative="account.isLiability"
                                                  :label="accountAmountTitle"
                                                  :placeholder="accountAmountTitle"
                                                  v-model="selectedAccount.balance"/>
                                </v-col>
                                <v-col cols="12" md="6" v-show="selectedAccount.balance"
                                       v-if="(!editAccountId || isNewAccount(selectedAccount)) && (account.type === AccountType.SingleAccount.type || currentAccountIndex >= 0)">
                                    <date-time-select
                                        :disabled="loading || submitting"
                                        :label="tt('Balance Time')"
                                        v-model="selectedAccount.balanceTime"
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
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="isInputEmpty() || loading || submitting" @click="save">
                        {{ tt(saveButtonTitle) }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountEditPageBaseBase } from '@/views/base/accounts/AccountEditPageBase.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { AccountType } from '@/core/account.ts';
import { ALL_ACCOUNT_ICONS } from '@/consts/icon.ts';
import { ALL_ACCOUNT_COLORS } from '@/consts/color.ts';
import { Account } from '@/models/account.ts';

import { isNumber } from '@/lib/common.ts';
import { getCurrentUnixTime } from '@/lib/datetime.ts';
import { generateRandomUUID } from '@/lib/misc.ts';

import {
    mdiDotsVertical,
    mdiCreditCardPlusOutline,
    mdiDeleteOutline
} from '@mdi/js';

interface AccountEditResponse {
    message: string;
}

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

defineProps<{
    show?: boolean;
}>();

const { tt } = useI18n();
const {
    editAccountId,
    clientSessionId,
    loading,
    submitting,
    account,
    subAccounts,
    title,
    saveButtonTitle,
    allAccountCategories,
    allAccountTypes,
    allAvailableMonthDays,
    isAccountSupportCreditCardStatementDate,
    isNewAccount,
    isInputEmpty,
    getAccountOrSubAccountProblemMessage,
    addSubAccount,
    setAccount
} = useAccountEditPageBaseBase();

const userStore = useUserStore();
const accountsStore = useAccountsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const activeTab = ref<string>('account');
const currentAccountIndex = ref<number>(-1);

const selectedAccount = computed<Account>(() => {
    if (currentAccountIndex.value < 0) {
        return account.value;
    }

    return subAccounts.value[currentAccountIndex.value];
});

const accountAmountTitle = computed<string>(() => {
    if (currentAccountIndex.value < 0) {
        return account.value.isLiability ? tt('Account Outstanding Balance') : tt('Account Balance');
    } else {
        return account.value.isLiability ? tt('Sub-account Outstanding Balance') : tt('Sub-account Balance');
    }
});

const isAccountModified = computed<boolean>(() => {
    if (!editAccountId.value) {
        return !account.value.equals(Account.createNewAccount(userStore.currentUserDefaultCurrency, account.value.balanceTime ?? getCurrentUnixTime()));
    } else {
        return true;
    }
});

let resolveFunc: ((value: AccountEditResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

function open(options?: { id?: string, currentAccount?: Account, category?: number }): Promise<AccountEditResponse> {
    showState.value = true;
    loading.value = true;
    submitting.value = false;

    const newAccount = Account.createNewAccount(userStore.currentUserDefaultCurrency, getCurrentUnixTime());
    account.value.fillFrom(newAccount);
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
            account.value.category = options.category;
            account.value.setSuitableIcon(1, options.category);
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
    const problemMessage = getAccountOrSubAccountProblemMessage();

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

function removeSubAccount(subAccount: Account): void {
    confirmDialog.value?.open('Are you sure you want to remove this sub-account?').then(() => {
        for (let i = 0; i < subAccounts.value.length; i++) {
            if (subAccounts.value[i] === subAccount) {
                subAccounts.value.splice(i, 1);

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
