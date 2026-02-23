import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { NumeralSystem } from '@/core/numeral.ts';
import type { WeekDayValue } from '@/core/datetime.ts';
import type { LocalizedTimezoneInfo } from '@/core/timezone.ts';
import { TransactionType } from '@/core/transaction.ts';
import { TemplateType } from '@/core/template.ts';
import { DISPLAY_HIDDEN_AMOUNT } from '@/consts/numeral.ts';
import { TRANSACTION_MAX_PICTURE_COUNT } from '@/consts/transaction.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import type { TransactionPictureInfoBasicResponse } from '@/models/transaction_picture_info.ts';
import { Transaction } from '@/models/transaction.ts';
import { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    isArray,
    isDefined
} from '@/lib/common.ts';

import {
    getExchangedAmountByRate
} from '@/lib/numeral.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    getSameDateTimeWithCurrentTimezone,
    parseDateTimeFromUnixTimeWithBrowserTimezone,
    getCurrentUnixTime
} from '@/lib/datetime.ts';

import {
    type SetTransactionOptions,
    setTransactionModelByTransaction
} from '@/lib/transaction.ts';

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
        getCurrentNumeralSystemType,
        getTimezoneDifferenceDisplayText,
        formatAmountToLocalizedNumeralsWithCurrency,
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
    const setGeoLocationByClickMap = ref<boolean>(false);

    const transaction = ref<Transaction | TransactionTemplate>(createNewTransactionModel(transactionDefaultType));

    const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
    const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(transaction.value.time));
    const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
    const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const defaultAccountId = computed<string>(() => userStore.currentUserDefaultAccountId);
    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const coordinateDisplayType = computed<number>(() => userStore.currentUserCoordinateDisplayType);

    const allTimezones = computed<LocalizedTimezoneInfo[]>(() => {
        if (type === TransactionEditPageType.Template && transaction.value instanceof TransactionTemplate) {
            return getAllTimezones(getCurrentUnixTime(), true);
        } else {
            return getAllTimezones(transaction.value.time, true)
        }
    });
    const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
    const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));
    const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
    const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
    const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);
    const firstVisibleAccountId = computed<string | undefined>(() => allVisibleAccounts.value && allVisibleAccounts.value[0] ? allVisibleAccounts.value[0].id : undefined);

    const hasVisibleExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleExpenseCategories);
    const hasVisibleIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleIncomeCategories);
    const hasVisibleTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleTransferCategories);

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

    const sourceAmountTitle = computed<string>(() => {
        const sourceAccount = allAccountsMap.value[transaction.value.sourceAccountId];
        const amountName = tt(sourceAmountName.value);

        if (!sourceAccount || sourceAccount.currency === defaultCurrency.value || !transaction.value.sourceAmount || transaction.value.hideAmount) {
            return amountName;
        }

        const fromExchangeRate = exchangeRatesStore.latestExchangeRateMap[sourceAccount.currency];
        const toExchangeRate = exchangeRatesStore.latestExchangeRateMap[defaultCurrency.value];

        if (!fromExchangeRate || !fromExchangeRate.rate || !toExchangeRate || !toExchangeRate.rate) {
            return amountName;
        }

        let amountInDefaultCurrency = getExchangedAmountByRate(transaction.value.sourceAmount, fromExchangeRate.rate, toExchangeRate.rate);

        if (!amountInDefaultCurrency) {
            return amountName;
        }

        amountInDefaultCurrency = Math.trunc(amountInDefaultCurrency);

        const displayAmountInDefaultCurrency = getDisplayAmount(amountInDefaultCurrency, transaction.value.hideAmount, defaultCurrency.value);
        return amountName + ` (${displayAmountInDefaultCurrency})`;
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
        const utcOffset = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(getUtcOffsetByUtcOffsetMinutes(transaction.value.utcOffset));
        return `UTC${utcOffset}`;
    });

    const transactionTimezoneTimeDifference = computed<string>(() => {
        return getTimezoneDifferenceDisplayText(transaction.value.time, transaction.value.utcOffset);
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

        if (type === TransactionEditPageType.Template && transaction.value instanceof TransactionTemplate) {
            if (!transaction.value.name) {
                return 'Template name cannot be blank';
            }
        }

        return null;
    });

    const inputIsEmpty = computed<boolean>(() => {
        return !!inputEmptyProblemMessage.value;
    });

    function getCurrentUnixTimeForNewTransaction(): number {
        return getSameDateTimeWithCurrentTimezone(parseDateTimeFromUnixTimeWithBrowserTimezone(getCurrentUnixTime())).getUnixTime();
    }

    function createNewTransactionModel(transactionType?: number): Transaction | TransactionTemplate {
        const now: number = getCurrentUnixTimeForNewTransaction();
        const currentTimezone: string = settingsStore.appSettings.timeZone;

        let defaultType: TransactionType = TransactionType.Expense;

        if (transactionType === TransactionType.Income) {
            defaultType = TransactionType.Income;
        } else if (transactionType === TransactionType.Transfer) {
            defaultType = TransactionType.Transfer;
        }

        let newTransaction: Transaction | TransactionTemplate = Transaction.createNewTransaction(defaultType, now, currentTimezone, getTimezoneOffsetMinutes(now, currentTimezone));

        if (type === TransactionEditPageType.Template) {
            newTransaction = TransactionTemplate.createNewTransactionTemplate(newTransaction);
        }

        return newTransaction;
    }

    function setTransactionModel(newTransaction: Transaction | null, options: SetTransactionOptions | undefined, setContextData: boolean): void {
        setTransactionModelByTransaction(
            transaction.value,
            newTransaction,
            allCategories.value,
            allCategoriesMap.value,
            allVisibleAccounts.value,
            allAccountsMap.value,
            allTagsMap.value,
            defaultAccountId.value,
            {
                time: options?.time,
                type: options?.type,
                categoryId: options?.categoryId,
                accountId: options?.accountId,
                destinationAccountId: options?.destinationAccountId,
                amount: options?.amount,
                destinationAmount: options?.destinationAmount,
                tagIds: options?.tagIds,
                comment: options?.comment
            },
            setContextData
        );
    }

    function updateTransactionTime(newTime: number): void {
        transaction.value.time = newTime;
        updateTransactionTimezone(transaction.value.timeZone ?? '');
    }

    function updateTransactionTimezone(timezoneName: string): void {
        const oldUtcOffset = transaction.value.utcOffset;

        for (const timezone of allTimezones.value) {
            if (timezone.name === timezoneName) {
                transaction.value.timeZone = timezone.name;
                transaction.value.utcOffset = timezone.utcOffsetMinutes;
                break;
            }
        }

        transaction.value.time = transaction.value.time - (transaction.value.utcOffset - oldUtcOffset) * 60;
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

    function getDisplayAmount(amount: number, hideAmount: boolean, currencyCode: string): string {
        if (hideAmount) {
            return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, currencyCode);
        }

        return formatAmountToLocalizedNumeralsWithCurrency(amount, currencyCode);
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

    // Watch for account changes and recalculate destination amount for transfers
    watch(() => [transaction.value.sourceAccountId, transaction.value.destinationAccountId], ([newSourceAccountId, newDestinationAccountId], [oldSourceAccountId, oldDestinationAccountId]) => {
        if (mode.value === TransactionEditPageMode.View || loading.value) {
            return;
        }

        if (transaction.value.type !== TransactionType.Transfer) {
            return;
        }

        // Only recalculate if accounts actually changed (skip initial watch call)
        if (isDefined(oldSourceAccountId) && isDefined(oldDestinationAccountId)) {
            if (newSourceAccountId === oldSourceAccountId && newDestinationAccountId === oldDestinationAccountId) {
                return;
            }
        }

        transactionsStore.setTransactionSuitableDestinationAmount(transaction.value, transaction.value.sourceAmount, transaction.value.sourceAmount, oldSourceAccountId, oldDestinationAccountId);
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
        setGeoLocationByClickMap,
        transaction,
        // computed states
        numeralSystem,
        currentTimezoneOffsetMinutes,
        showAccountBalance,
        defaultCurrency,
        defaultAccountId,
        firstDayOfWeek,
        coordinateDisplayType,
        allTimezones,
        allAccounts,
        allVisibleAccounts,
        allAccountsMap,
        allVisibleCategorizedAccounts,
        allCategories,
        allCategoriesMap,
        allTagsMap,
        firstVisibleAccountId,
        hasVisibleExpenseCategories,
        hasVisibleIncomeCategories,
        hasVisibleTransferCategories,
        canAddTransactionPicture,
        title,
        saveButtonTitle,
        cancelButtonTitle,
        sourceAmountName,
        sourceAmountTitle,
        sourceAccountTitle,
        transferInAmountTitle,
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
        setTransactionModel,
        updateTransactionTime,
        updateTransactionTimezone,
        swapTransactionData,
        getDisplayAmount,
        getTransactionPictureUrl
    }
}
