import transactionConstants from '../consts/transaction.js';
import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utils from '../lib/utils.js';

import { getExchangedAmount } from "./exchangeRates.js";

import {
    LOAD_TRANSACTION_LIST,
    COLLAPSE_MONTH_IN_TRANSACTION_LIST,
    REMOVE_TRANSACTION_FROM_TRANSACTION_LIST,
    UPDATE_TRANSACTION_LIST_INVALID_STATE,
    UPDATE_ACCOUNT_LIST_INVALID_STATE,
} from './mutations.js';

const emptyTransactionResult = {
    items: [],
    transactionsNextTimeId: 0
};

export function getTransactions(context, { reload, autoExpand, defaultCurrency, maxTime, minTime, type, categoryId, accountId, keyword }) {
    let actualMaxTime = context.state.transactionsNextTimeId;

    if (reload && maxTime > 0) {
        actualMaxTime = maxTime;
    } else if (reload && maxTime <= 0) {
        actualMaxTime = 0;
    }

    return new Promise((resolve, reject) => {
        services.getTransactions({
            maxTime: actualMaxTime,
            minTime: minTime,
            type: type,
            categoryId: categoryId,
            accountId: accountId,
            keyword: keyword
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                if (reload) {
                    context.commit(LOAD_TRANSACTION_LIST, {
                        transactions: emptyTransactionResult,
                        reload: reload,
                        autoExpand: autoExpand,
                        defaultCurrency: defaultCurrency,
                        accountId: accountId
                    });
                    context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, true);
                }

                reject({ message: 'Unable to get transaction list' });
                return;
            }

            context.commit(LOAD_TRANSACTION_LIST, {
                transactions: data.result,
                reload: reload,
                autoExpand: autoExpand,
                defaultCurrency: defaultCurrency,
                accountId: accountId
            });

            if (reload) {
                context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, false);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to load transaction list', error);

            if (reload) {
                context.commit(LOAD_TRANSACTION_LIST, {
                    transactions: emptyTransactionResult,
                    reload: reload,
                    autoExpand: autoExpand,
                    defaultCurrency: defaultCurrency,
                    accountId: accountId
                });
                context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, true);
            }

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get transaction list' });
            } else {
                reject(error);
            }
        });
    });
}

export function getTransaction(context, { transactionId }) {
    return new Promise((resolve, reject) => {
        services.getTransaction({
            id: transactionId
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get transaction' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to load transaction info', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get transaction' });
            } else {
                reject(error);
            }
        });
    });
}

export function saveTransaction(context, { transaction }) {
    return new Promise((resolve, reject) => {
        let promise = null;

        if (!transaction.id) {
            promise = services.addTransaction(transaction);
        } else {
            promise = services.modifyTransaction(transaction);
        }

        promise.then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                if (!transaction.id) {
                    reject({ message: 'Unable to add transaction' });
                } else {
                    reject({ message: 'Unable to save transaction' });
                }
                return;
            }

            if (!transaction.id) {
                context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, true);
            } else {
                context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, true);
            }

            context.commit(UPDATE_ACCOUNT_LIST_INVALID_STATE, true);

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to save transaction', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                if (!transaction.id) {
                    reject({ message: 'Unable to add transaction' });
                } else {
                    reject({ message: 'Unable to save transaction' });
                }
            } else {
                reject(error);
            }
        });
    });
}

export function deleteTransaction(context, { transaction, defaultCurrency, accountId, beforeResolve }) {
    return new Promise((resolve, reject) => {
        services.deleteTransaction({
            id: transaction.id
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to delete this transaction' });
                return;
            }

            if (beforeResolve) {
                beforeResolve(() => {
                    context.commit(REMOVE_TRANSACTION_FROM_TRANSACTION_LIST, {
                        transaction: transaction,
                        defaultCurrency: defaultCurrency,
                        accountId: accountId
                    });
                });
            } else {
                context.commit(REMOVE_TRANSACTION_FROM_TRANSACTION_LIST, {
                    transaction: transaction,
                    defaultCurrency: defaultCurrency,
                    accountId: accountId
                });
            }

            context.commit(UPDATE_ACCOUNT_LIST_INVALID_STATE, true);

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to delete transaction', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to delete this transaction' });
            } else {
                reject(error);
            }
        });
    });
}

export function collapseMonthInTransactionList(context, { month, collapse }) {
    context.commit(COLLAPSE_MONTH_IN_TRANSACTION_LIST, {
        month: month,
        collapse: collapse
    });
}

export function noTransaction(state) {
    for (let i = 0; i < state.transactions.length; i++) {
        const transactionMonthList = state.transactions[i];

        for (let j = 0; j < transactionMonthList.items.length; j++) {
            if (transactionMonthList.items[j]) {
                return false;
            }
        }
    }

    return true;
}

export function hasMoreTransaction(state) {
    return state.transactionsNextTimeId > 0;
}

export function calculateMonthTotalAmount(state, transactionMonthList, defaultCurrency, accountId, incomplete) {
    if (!transactionMonthList) {
        return;
    }

    let totalExpense = 0;
    let totalIncome = 0;
    let hasUnCalculatedTotalExpense = false;
    let hasUnCalculatedTotalIncome = false;

    for (let i = 0; i < transactionMonthList.items.length; i++) {
        const transaction = transactionMonthList.items[i];

        if (!transaction.sourceAccount) {
            continue;
        }

        let amount = transaction.sourceAmount;

        if (transaction.sourceAccount.currency !== defaultCurrency) {
            const balance = getExchangedAmount(state)(amount, transaction.sourceAccount.currency, defaultCurrency);

            if (!utils.isNumber(balance)) {
                if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
                    hasUnCalculatedTotalExpense = true;
                } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
                    hasUnCalculatedTotalIncome = true;
                }

                continue;
            }

            amount = Math.floor(balance);
        }

        if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
            totalExpense += amount;
        } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
            totalIncome += amount;
        } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer && accountId) {
            if (accountId === transaction.sourceAccountId) {
                totalExpense += amount;
            } else if (accountId === transaction.destinationAccountId) {
                totalIncome += amount;
            }
        }
    }

    transactionMonthList.totalAmount = {
        expense: totalExpense,
        incompleteExpense: incomplete || hasUnCalculatedTotalExpense,
        income: totalIncome,
        incompleteIncome: incomplete || hasUnCalculatedTotalIncome
    };
}
