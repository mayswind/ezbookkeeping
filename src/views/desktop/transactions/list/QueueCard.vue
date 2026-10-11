<template>
    <v-card variant="flat" min-height="920">
        <template #title>
            <div class="title-and-toolbar d-flex align-center text-no-wrap">
                <v-btn class="me-3 d-lg-none" density="compact" color="default" variant="plain"
                       :ripple="false" :icon="true" @click="emit('toggle:nav')">
                    <v-icon :icon="mdiMenu" size="24" />
                </v-btn>
                <span>{{ tt('Queue') }}</span>
                <v-btn class="ms-3" color="default" variant="outlined"
                       :disabled="loading || processing || !pendingQueuedTransactions.length"
                       @click="acceptAll">
                    {{ tt('Accept All') }}
                </v-btn>
                <v-btn class="ms-3" color="default" variant="outlined"
                       :disabled="loading || processing || !pendingQueuedTransactions.length"
                       @click="rejectAll">
                    {{ tt('Reject All') }}
                </v-btn>
                <v-btn density="compact" color="default" variant="text"
                       class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="24" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
                <v-spacer/>
                <v-btn class="ms-2" color="default" variant="text" density="comfortable"
                       :disabled="loading || processing"
                       :prepend-icon="showArchived ? mdiEyeOffOutline : mdiEyeOutline"
                       @click="toggleArchived">
                    {{ showArchived ? tt('Hide Archived Drafts') : tt('Show Archived Drafts') }}
                </v-btn>
            </div>
        </template>

        <v-card-text class="pt-0">
            <v-progress-linear indeterminate v-if="loading && !displayedQueuedTransactions.length"/>

            <div class="d-flex flex-column align-center justify-center pt-10 pb-10"
                 v-if="!loading && !displayedQueuedTransactions.length">
                <v-avatar size="80" color="primary" variant="tonal">
                    <v-icon size="48" :icon="mdiTrayFull"/>
                </v-avatar>
                <span class="mt-4 text-body-1">{{ tt('The transaction queue is empty') }}</span>
                <span class="mt-1 text-body-2 text-medium-emphasis">{{ tt('Drafts created by scheduled transactions or by the API show up here') }}</span>
            </div>

            <v-table class="queue-table" v-if="displayedQueuedTransactions.length">
                <thead>
                    <tr>
                        <th class="text-left">{{ tt('Transaction Time') }}</th>
                        <th class="text-left">{{ tt('Type') }}</th>
                        <th class="text-left">{{ tt('Category') }}</th>
                        <th class="text-left">{{ tt('Account') }}</th>
                        <th class="text-right">{{ tt('Amount') }}</th>
                        <th class="text-left">{{ tt('Description') }}</th>
                        <th class="text-left">{{ tt('Source') }}</th>
                        <th class="text-right">{{ tt('Operation') }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr :key="draft.id" v-for="draft in displayedQueuedTransactions"
                        :class="{ 'queue-row-archived': draft.isArchived }">
                        <td class="text-left">
                            <span v-if="draft.time">{{ getDisplayDateTime(draft) }}</span>
                            <span class="text-medium-emphasis" v-else>{{ tt('Not set') }}</span>
                        </td>
                        <td class="text-left">{{ getDisplayType(draft) }}</td>
                        <td class="text-left">
                            <span v-if="getDisplayCategory(draft)">{{ getDisplayCategory(draft) }}</span>
                            <span class="text-medium-emphasis" v-else>{{ tt('Not set') }}</span>
                        </td>
                        <td class="text-left">
                            <span v-if="getDisplayAccount(draft)">{{ getDisplayAccount(draft) }}</span>
                            <span class="text-medium-emphasis" v-else>{{ tt('Not set') }}</span>
                        </td>
                        <td class="text-right" :class="getAmountColorClass(draft)">{{ getDisplayAmount(draft) }}</td>
                        <td class="text-left queue-comment-cell">
                            <span class="text-truncate d-inline-block">{{ draft.comment }}</span>
                        </td>
                        <td class="text-left">
                            <span class="text-body-2">{{ tt(getDisplaySourceName(draft)) }}</span>
                            <div v-if="draft.isArchived">
                                <v-chip size="x-small" label color="warning">{{ tt('Archived') }}</v-chip>
                                <div class="text-caption text-medium-emphasis" v-if="draft.expiresAt">
                                    {{ tt('format.misc.deletedAt', { time: getDisplayExpiresAt(draft) }) }}
                                </div>
                            </div>
                            <v-chip size="x-small" label color="info" v-else-if="!draft.isComplete">{{ tt('Incomplete') }}</v-chip>
                        </td>
                        <td class="text-right text-no-wrap">
                            <template v-if="!draft.isArchived">
                                <v-btn density="comfortable" color="success" variant="text" :icon="true"
                                       :disabled="loading || processing" @click="accept(draft)">
                                    <v-icon :icon="mdiCheck" size="22"/>
                                    <v-tooltip activator="parent">{{ tt('Accept') }}</v-tooltip>
                                </v-btn>
                                <v-btn density="comfortable" color="error" variant="text" :icon="true"
                                       :disabled="loading || processing" @click="reject(draft)">
                                    <v-icon :icon="mdiClose" size="22"/>
                                    <v-tooltip activator="parent">{{ tt('Reject') }}</v-tooltip>
                                </v-btn>
                            </template>
                            <v-btn density="comfortable" color="default" variant="text" :icon="true"
                                   :disabled="loading || processing" v-else @click="restore(draft)">
                                <v-icon :icon="mdiRestore" size="22"/>
                                <v-tooltip activator="parent">{{ tt('Restore') }}</v-tooltip>
                            </v-btn>
                        </td>
                    </tr>
                </tbody>
            </v-table>
        </v-card-text>

        <edit-dialog ref="editDialog" :type="TransactionEditPageType.Transaction" />
        <confirm-dialog ref="confirmDialog"/>
        <snack-bar ref="snackbar" />
    </v-card>
</template>

<script setup lang="ts">
import EditDialog from './dialogs/EditDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, onMounted } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useQueuedTransactionsStore } from '@/stores/queuedTransaction.ts';

import { TransactionType } from '@/core/transaction.ts';
import { QueuedTransaction, QueuedTransactionSource } from '@/models/queued_transaction.ts';

import { DISPLAY_HIDDEN_AMOUNT } from '@/consts/numeral.ts';

import type { ErrorResponse } from '@/core/api.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTimeWithTimezoneOffset } from '@/lib/datetime.ts';

import {
    mdiMenu,
    mdiRefresh,
    mdiCheck,
    mdiClose,
    mdiRestore,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiTrayFull
} from '@mdi/js';

type EditDialogType = InstanceType<typeof EditDialog>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const emit = defineEmits<{
    (e: 'toggle:nav'): void;
}>();

const { tt, formatDateTimeToLongDateTime, formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const queuedTransactionsStore = useQueuedTransactionsStore();

const editDialog = useTemplateRef<EditDialogType>('editDialog');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const processing = ref<boolean>(false);
const showArchived = ref<boolean>(false);

const pendingQueuedTransactions = computed<QueuedTransaction[]>(() => queuedTransactionsStore.pendingQueuedTransactions);

const displayedQueuedTransactions = computed<QueuedTransaction[]>(() => {
    if (showArchived.value) {
        return queuedTransactionsStore.allQueuedTransactions;
    }

    return pendingQueuedTransactions.value;
});

function getDisplayDateTime(draft: QueuedTransaction): string {
    return formatDateTimeToLongDateTime(parseDateTimeFromUnixTimeWithTimezoneOffset(draft.time, draft.utcOffset));
}

function getDisplayExpiresAt(draft: QueuedTransaction): string {
    if (!draft.expiresAt) {
        return '';
    }

    return formatDateTimeToLongDateTime(parseDateTimeFromUnixTimeWithTimezoneOffset(draft.expiresAt, draft.utcOffset));
}

function getDisplayType(draft: QueuedTransaction): string {
    switch (draft.type) {
        case TransactionType.ModifyBalance:
            return tt('Modify Balance');
        case TransactionType.Income:
            return tt('Income');
        case TransactionType.Expense:
            return tt('Expense');
        case TransactionType.Transfer:
            return tt('Transfer');
        default:
            return tt('Not set');
    }
}

function getDisplayCategory(draft: QueuedTransaction): string {
    if (!draft.categoryId || draft.categoryId === '0') {
        return '';
    }

    const category = transactionCategoriesStore.allTransactionCategoriesMap[draft.categoryId];
    return category ? category.name : '';
}

function getDisplayAccount(draft: QueuedTransaction): string {
    if (!draft.sourceAccountId || draft.sourceAccountId === '0') {
        return '';
    }

    const sourceAccount = accountsStore.allAccountsMap[draft.sourceAccountId];
    const sourceAccountName = sourceAccount ? sourceAccount.name : '';

    if (!draft.isTransfer) {
        return sourceAccountName;
    }

    const destinationAccount = draft.destinationAccountId ? accountsStore.allAccountsMap[draft.destinationAccountId] : undefined;
    const destinationAccountName = destinationAccount ? destinationAccount.name : tt('Not set');

    return sourceAccountName + ' → ' + destinationAccountName;
}

function getDisplayAmount(draft: QueuedTransaction): string {
    const sourceAccount = draft.sourceAccountId ? accountsStore.allAccountsMap[draft.sourceAccountId] : undefined;
    const currency = sourceAccount ? sourceAccount.currency : false;

    if (draft.hideAmount) {
        return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, currency);
    }

    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(draft.sourceAmount), currency);
}

function getAmountColorClass(draft: QueuedTransaction): string {
    if (draft.type === TransactionType.Expense) {
        return 'text-expense';
    } else if (draft.type === TransactionType.Income) {
        return 'text-income';
    }

    return '';
}

function getDisplaySourceName(draft: QueuedTransaction): string {
    switch (draft.source) {
        case QueuedTransactionSource.ScheduledTemplate:
            return 'Scheduled Transaction';
        case QueuedTransactionSource.Api:
            return 'API';
        case QueuedTransactionSource.Mcp:
            return 'MCP';
        default:
            return 'Unknown';
    }
}

function toggleArchived(): void {
    showArchived.value = !showArchived.value;
    reload(false);
}

function reload(force: boolean): void {
    loading.value = true;

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ]).then(() => {
        return queuedTransactionsStore.loadAllQueuedTransactions({
            includeArchived: showArchived.value,
            force: true
        });
    }).then(() => {
        loading.value = false;

        if (force) {
            snackbar.value?.showMessage('Data has been updated');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

/**
 * Reports a rejection that carries an actual failure, and tells the caller whether it did. A
 * cancelled edit dialog rejects with no reason, which is a normal outcome and not worth reporting.
 */
function showErrorIfAny(error: unknown): boolean {
    const typedError = error as { processed?: boolean, message?: string, error?: ErrorResponse } | undefined;

    if (!typedError || typedError.processed) {
        return false;
    }

    if (typedError.error) {
        snackbar.value?.showError({ error: typedError.error });
        return true;
    }

    if (typedError.message) {
        snackbar.value?.showError({ message: typedError.message });
        return true;
    }

    return false;
}

/**
 * Opens the transaction edit dialog filled in with the draft, so the user completes whatever the
 * draft is missing and turns it into a real transaction. The draft only leaves the queue once that
 * transaction has actually been saved.
 */
function acceptOne(draft: QueuedTransaction): Promise<void> {
    if (!editDialog.value) {
        return Promise.reject(new Error('the transaction edit dialog is not ready'));
    }

    return editDialog.value.open({
        type: draft.type || undefined,
        categoryId: draft.categoryId && draft.categoryId !== '0' ? draft.categoryId : undefined,
        accountId: draft.sourceAccountId && draft.sourceAccountId !== '0' ? draft.sourceAccountId : undefined,
        destinationAccountId: draft.destinationAccountId && draft.destinationAccountId !== '0' ? draft.destinationAccountId : undefined,
        amount: draft.sourceAmount,
        destinationAmount: draft.destinationAmount,
        tagIds: draft.tagIds && draft.tagIds.length ? draft.tagIds.join(',') : undefined,
        comment: draft.comment || undefined,
        time: draft.time || undefined,
        // accepting a draft must not overwrite the locally autosaved transaction draft of the user
        noTransactionDraft: true
    }).then(() => {
        return queuedTransactionsStore.deleteQueuedTransaction({ id: draft.id });
    }).then(() => {
        return;
    });
}

function accept(draft: QueuedTransaction): void {
    processing.value = true;

    acceptOne(draft).then(() => {
        processing.value = false;
        snackbar.value?.showMessage('You have added a new transaction');
    }).catch(error => {
        processing.value = false;
        showErrorIfAny(error);
    });
}

/**
 * Walks the pending drafts one by one, opening the edit dialog for the next one as soon as the
 * current one is done. Cancelling a dialog stops the walk and leaves the remaining drafts queued.
 */
async function acceptAll(): Promise<void> {
    const drafts = [...pendingQueuedTransactions.value];
    processing.value = true;

    let acceptedCount = 0;

    for (const draft of drafts) {
        try {
            await acceptOne(draft);
            acceptedCount++;
        } catch (error) {
            processing.value = false;

            if (showErrorIfAny(error)) {
                return;
            }

            // the dialog rejects with no reason when the user cancels, so stop the walk here and
            // leave the drafts that were not reached in the queue
            if (acceptedCount > 0) {
                snackbar.value?.showMessage('format.misc.acceptedDraftCount', { count: acceptedCount });
            }

            return;
        }
    }

    processing.value = false;

    if (acceptedCount > 0) {
        snackbar.value?.showMessage('format.misc.acceptedDraftCount', { count: acceptedCount });
    }
}

function reject(draft: QueuedTransaction): void {
    confirmDialog.value?.open('Are you sure you want to reject this draft? It will be archived and permanently deleted after 24 hours.').then(() => {
        processing.value = true;

        queuedTransactionsStore.archiveQueuedTransaction({ id: draft.id }).then(() => {
            processing.value = false;
            reload(false);
        }).catch(error => {
            processing.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function rejectAll(): void {
    confirmDialog.value?.open('Are you sure you want to reject all drafts? They will be archived and permanently deleted after 24 hours.').then(() => {
        processing.value = true;

        queuedTransactionsStore.archiveAllQueuedTransactions().then(() => {
            processing.value = false;
            reload(false);
        }).catch(error => {
            processing.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function restore(draft: QueuedTransaction): void {
    processing.value = true;

    queuedTransactionsStore.restoreQueuedTransaction({ id: draft.id }).then(() => {
        processing.value = false;
        reload(false);
    }).catch(error => {
        processing.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

onMounted(() => {
    reload(false);
});

defineExpose({
    reload
});
</script>

<style>
.queue-table .queue-row-archived {
    opacity: 0.6;
}

.queue-table .queue-comment-cell {
    max-width: 240px;
}

.queue-table .queue-comment-cell > span {
    max-width: 100%;
}
</style>
