import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { LocalizedTimezoneInfo } from '@/core/timezone.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import { TemplateType } from '@/core/template.ts';
import { TRANSACTION_MAX_PICTURE_COUNT } from '@/consts/transaction.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import type { TransactionPictureInfoBasicResponse } from '@/models/transaction_picture_info.ts';
import { Transaction } from '@/models/transaction.ts';
import { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    isArray
} from '@/lib/common.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    getCurrentUnixTime
} from '@/lib/datetime.ts';

import {
    getFirstAvailableCategoryId
} from '@/lib/category.ts';

export enum TransactionEditPageType {
    Transaction = 'transaction',
    Template = 'template'
}

export enum TransactionEditPageMode {
    Add = 'add',
    Edit = 'edit',
    View = 'view'
}

export enum GeoLocationStatus {
    Getting = 'getting',
    Success = 'success',
    Error = 'error'
}

export function useTransactionEditPageBase(type: TransactionEditPageType, initMode?: TransactionEditPageMode, transactionDefaultType?: number) {
    const {
        tt,
        getAllTimezones,
        getTimezoneDifferenceDisplayText,
        formatAmountWithCurrency,
        getAdaptiveAmountRate,
        getCategorizedAccountsWithDisplayBalance
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();
    const transactionsStore = useTransactionsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const isSupportGeoLocation: boolean = !!navigator.geolocation;

    const mode = ref<TransactionEditPageMode>(initMode ?? TransactionEditPageMode.Add);
    const editId = ref<string | null>(null);
    const addByTemplateId = ref<string | null>(null);
    const duplicateFromId = ref<string | null>(null);

    const clientSessionId = ref<string>('');
    const loading = ref<boolean>(true);
    const submitting = ref<boolean>(false);
    const uploadingPicture = ref<boolean>(false);
    const geoLocationStatus = ref<GeoLocationStatus | null>(null);

    const transaction = ref<Transaction | TransactionTemplate>(createNewTransactionModel(transactionDefaultType));

    const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone));
    const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const defaultAccountId = computed<string>(() => userStore.currentUserDefaultAccountId);
    const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);

    const allTimezones = computed<LocalizedTimezoneInfo[]>(() => getAllTimezones(true));
    const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
    const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value));
    const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
    const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
    const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
    const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);

    const canAddTransactionPicture = computed<boolean>(() => {
        if (type !== TransactionEditPageType.Transaction || (mode.value !== TransactionEditPageMode.Add && mode.value !== TransactionEditPageMode.Edit)) {
            return false;
        }

        return !isArray(transaction.value.pictures) || transaction.value.pictures.length < TRANSACTION_MAX_PICTURE_COUNT;
    });

    const title = computed<string>(() => {
        if (type === TransactionEditPageType.Transaction) {
            if (mode.value === TransactionEditPageMode.Add) {
                return 'Add Transaction';
            } else if (mode.value === TransactionEditPageMode.Edit) {
                return 'Edit Transaction';
            } else {
                return 'Transaction Detail';
            }
        } else if (type === TransactionEditPageType.Template && (transaction.value as TransactionTemplate).templateType === TemplateType.Normal.type) {
            if (mode.value === TransactionEditPageMode.Add) {
                return 'Add Transaction Template';
            } else if (mode.value === TransactionEditPageMode.Edit) {
                return 'Edit Transaction Template';
            }
        } else if (type === TransactionEditPageType.Template && (transaction.value as TransactionTemplate).templateType === TemplateType.Schedule.type) {
            if (mode.value === TransactionEditPageMode.Add) {
                return 'Add Scheduled Transaction';
            } else if (mode.value === TransactionEditPageMode.Edit) {
                return 'Edit Scheduled Transaction';
            }
        }

        return '';
    });

    const saveButtonTitle = computed<string>(() => {
        if (mode.value === TransactionEditPageMode.Add) {
            return 'Add';
        } else {
            return 'Save';
        }
    });

    const cancelButtonTitle = computed<string>(() => {
        if (mode.value === TransactionEditPageMode.View) {
            return 'Close';
        } else {
            return 'Cancel';
        }
    });

    const sourceAmountName = computed<string>(() => {
        if (transaction.value.type === TransactionType.Expense) {
            return 'Expense Amount';
        } else if (transaction.value.type === TransactionType.Income) {
            return 'Income Amount';
        } else if (transaction.value.type === TransactionType.Transfer) {
            return 'Transfer Out Amount';
        } else {
            return 'Amount';
        }
    });

    const sourceAccountTitle = computed<string>(() => {
        if (transaction.value.type === TransactionType.Expense || transaction.value.type === TransactionType.Income) {
            return 'Account';
        } else if (transaction.value.type === TransactionType.Transfer) {
            return 'Source Account';
        } else {
            return 'Account';
        }
    });

    const transferInAmountTitle = computed<string>(() => {
        const sourceAccount = allAccountsMap.value[transaction.value.sourceAccountId];
        const destinationAccount = allAccountsMap.value[transaction.value.destinationAccountId];

        if (!sourceAccount || !destinationAccount || sourceAccount.currency === destinationAccount.currency) {
            return tt('Transfer In Amount');
        }

        const fromExchangeRate = exchangeRatesStore.latestExchangeRateMap[sourceAccount.currency];
        const toExchangeRate = exchangeRatesStore.latestExchangeRateMap[destinationAccount.currency];
        const amountRate = getAdaptiveAmountRate(transaction.value.sourceAmount, transaction.value.destinationAmount, fromExchangeRate, toExchangeRate);

        if (!amountRate) {
            return tt('Transfer In Amount');
        }

        return tt('Transfer In Amount') + ` (${amountRate})`;
    });

    const hasAvailableExpenseCategories = computed<boolean>(() => {
        if (!allCategories.value || !allCategories.value[CategoryType.Expense] || !allCategories.value[CategoryType.Expense].length) {
            return false;
        }

        const firstAvailableCategoryId = getFirstAvailableCategoryId(allCategories.value[CategoryType.Expense]);
        return firstAvailableCategoryId !== '';
    });

    const hasAvailableIncomeCategories = computed<boolean>(() => {
        if (!allCategories.value || !allCategories.value[CategoryType.Income] || !allCategories.value[CategoryType.Income].length) {
            return false;
        }

        const firstAvailableCategoryId = getFirstAvailableCategoryId(allCategories.value[CategoryType.Income]);
        return firstAvailableCategoryId !== '';
    });

    const hasAvailableTransferCategories = computed<boolean>(() => {
        if (!allCategories.value || !allCategories.value[CategoryType.Transfer] || !allCategories.value[CategoryType.Transfer].length) {
            return false;
        }

        const firstAvailableCategoryId = getFirstAvailableCategoryId(allCategories.value[CategoryType.Transfer]);
        return firstAvailableCategoryId !== '';
    });

    const sourceAccountName = computed<string>(() => {
        if (transaction.value.sourceAccountId) {
            return Account.findAccountNameById(allAccounts.value, transaction.value.sourceAccountId) || '';
        } else {
            return tt('None');
        }
    });

    const destinationAccountName = computed<string>(() => {
        if (transaction.value.destinationAccountId) {
            return Account.findAccountNameById(allAccounts.value, transaction.value.destinationAccountId) || '';
        } else {
            return tt('None');
        }
    });

    const sourceAccountCurrency = computed<string>(() => {
        const sourceAccount = allAccountsMap.value[transaction.value.sourceAccountId];

        if (sourceAccount) {
            return sourceAccount.currency;
        }

        return defaultCurrency.value;
    });

    const destinationAccountCurrency = computed<string>(() => {
        const destinationAccount = allAccountsMap.value[transaction.value.destinationAccountId];

        if (destinationAccount) {
            return destinationAccount.currency;
        }

        return defaultCurrency.value;
    });

    const transactionDisplayTimezone = computed<string>(() => {
        return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.value.utcOffset)}`;
    });

    const transactionTimezoneTimeDifference = computed<string>(() => {
        return getTimezoneDifferenceDisplayText(transaction.value.utcOffset);
    });

    const geoLocationStatusInfo = computed<string>(() => {
        if (geoLocationStatus.value === GeoLocationStatus.Success) {
            return '';
        } else if (geoLocationStatus.value === GeoLocationStatus.Getting) {
            return tt('Getting Location...');
        } else {
            return tt('No Location');
        }
    });

    const inputEmptyProblemMessage = computed<string | null>(() => {
        if (transaction.value.type === TransactionType.Expense) {
            if (!transaction.value.expenseCategoryId || transaction.value.expenseCategoryId === '') {
                return 'Transaction category cannot be blank';
            }

            if (!transaction.value.sourceAccountId || transaction.value.sourceAccountId === '') {
                return 'Transaction account cannot be blank';
            }
        } else if (transaction.value.type === TransactionType.Income) {
            if (!transaction.value.incomeCategoryId || transaction.value.incomeCategoryId === '') {
                return 'Transaction category cannot be blank';
            }

            if (!transaction.value.sourceAccountId || transaction.value.sourceAccountId === '') {
                return 'Transaction account cannot be blank';
            }
        } else if (transaction.value.type === TransactionType.Transfer) {
            if (!transaction.value.transferCategoryId || transaction.value.transferCategoryId === '') {
                return 'Transaction category cannot be blank';
            }

            if (!transaction.value.sourceAccountId || transaction.value.sourceAccountId === '') {
                return 'Source account cannot be blank';
            }

            if (!transaction.value.destinationAccountId || transaction.value.destinationAccountId === '') {
                return 'Destination account cannot be blank';
            }
        }

        if (type === 'template' && transaction.value instanceof TransactionTemplate) {
            if (!transaction.value.name) {
                return 'Template name cannot be blank';
            }
        }

        return null;
    });

    const inputIsEmpty = computed<boolean>(() => {
        return !!inputEmptyProblemMessage.value;
    });

    function createNewTransactionModel(transactionType?: number): Transaction | TransactionTemplate {
        const now: number = getCurrentUnixTime();
        const currentTimezone: string = settingsStore.appSettings.timeZone;

        let defaultType: TransactionType = TransactionType.Expense;

        if (transactionType === TransactionType.Income) {
            defaultType = TransactionType.Income;
        } else if (transactionType === TransactionType.Transfer) {
            defaultType = TransactionType.Transfer;
        }

        let newTransaction: Transaction | TransactionTemplate = Transaction.createNewTransaction(defaultType, now, currentTimezone, getTimezoneOffsetMinutes(currentTimezone));

        if (type === TransactionEditPageType.Template) {
            newTransaction = TransactionTemplate.createNewTransactionTemplate(newTransaction);
        }

        return newTransaction;
    }

    function swapTransactionData(swapAccount: boolean, swapAmount: boolean): void {
        if (swapAccount) {
            const oldSourceAccountId = transaction.value.sourceAccountId;
            transaction.value.sourceAccountId = transaction.value.destinationAccountId;
            transaction.value.destinationAccountId = oldSourceAccountId;
        }

        if (swapAmount) {
            const oldSourceAmount = transaction.value.sourceAmount;
            transaction.value.sourceAmount = transaction.value.destinationAmount;
            transaction.value.destinationAmount = oldSourceAmount;
        }
    }

    function getDisplayAmount(amount: number | string, hideAmount: boolean, currencyCode: string): string {
        if (hideAmount) {
            return formatAmountWithCurrency('***', currencyCode);
        }

        return formatAmountWithCurrency(amount, currencyCode);
    }

    function getTransactionPictureUrl(pictureInfo?: TransactionPictureInfoBasicResponse | null): string | undefined {
        return transactionsStore.getTransactionPictureUrl(pictureInfo);
    }

    watch(() => transaction.value.sourceAmount, (newValue, oldValue) => {
        if (mode.value === TransactionEditPageMode.View || loading.value) {
            return;
        }

        transactionsStore.setTransactionSuitableDestinationAmount(transaction.value, oldValue, newValue);
    });

    watch(() => transaction.value.destinationAmount, (newValue) => {
        if (mode.value === TransactionEditPageMode.View || loading.value) {
            return;
        }

        if (transaction.value.type === TransactionType.Expense || transaction.value.type === TransactionType.Income) {
            transaction.value.sourceAmount = newValue;
        }
    });

    watch(() => transaction.value.timeZone, (newValue) => {
        for (let i = 0; i < allTimezones.value.length; i++) {
            if (allTimezones.value[i].name === newValue) {
                transaction.value.utcOffset = allTimezones.value[i].utcOffsetMinutes;
                break;
            }
        }
    });

    return {
        // constants
        isSupportGeoLocation,
        // states
        mode,
        editId,
        addByTemplateId,
        duplicateFromId,
        clientSessionId,
        loading,
        submitting,
        uploadingPicture,
        geoLocationStatus,
        transaction,
        // computed states
        currentTimezoneOffsetMinutes,
        showAccountBalance,
        defaultCurrency,
        firstDayOfWeek,
        defaultAccountId,
        allTimezones,
        allAccounts,
        allVisibleAccounts,
        allAccountsMap,
        allVisibleCategorizedAccounts,
        allCategories,
        allCategoriesMap,
        allTags,
        allTagsMap,
        canAddTransactionPicture,
        title,
        saveButtonTitle,
        cancelButtonTitle,
        sourceAmountName,
        sourceAccountTitle,
        transferInAmountTitle,
        hasAvailableExpenseCategories,
        hasAvailableIncomeCategories,
        hasAvailableTransferCategories,
        sourceAccountName,
        destinationAccountName,
        sourceAccountCurrency,
        destinationAccountCurrency,
        transactionDisplayTimezone,
        transactionTimezoneTimeDifference,
        geoLocationStatusInfo,
        inputEmptyProblemMessage,
        inputIsEmpty,
        // functions
        createNewTransactionModel,
        swapTransactionData,
        getDisplayAmount,
        getTransactionPictureUrl
    }
}
