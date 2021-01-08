import userState from "../lib/userstate.js";
import utils from "../lib/utils.js";

import {
    RESET_STATE,

    STORE_USER_INFO,
    CLEAR_USER_INFO,

    STORE_LATEST_EXCHANGE_RATES,

    LOAD_ACCOUNT_LIST,
    ADD_ACCOUNT_TO_ACCOUNT_LIST,
    SAVE_ACCOUNT_IN_ACCOUNT_LIST,
    CHANGE_ACCOUNT_DISPLAY_ORDER_IN_ACCOUNT_LIST,
    UPDATE_ACCOUNT_VISIBILITY_IN_ACCOUNT_LIST,
    REMOVE_ACCOUNT_FROM_ACCOUNT_LIST,
    UPDATE_ACCOUNT_LIST_INVALID_STATE,

    LOAD_TRANSACTION_LIST,
    COLLAPSE_MONTH_IN_TRANSACTION_LIST,
    SAVE_TRANSACTION_IN_TRANSACTION_LIST,
    REMOVE_TRANSACTION_FROM_TRANSACTION_LIST,
    UPDATE_TRANSACTION_LIST_INVALID_STATE,

    LOAD_TRANSACTION_CATEGORY_LIST,
    ADD_CATEGORY_TO_TRANSACTION_CATEGORY_LIST,
    SAVE_CATEGORY_IN_TRANSACTION_CATEGORY_LIST,
    CHANGE_CATEGORY_DISPLAY_ORDER_IN_CATEGORY_LIST,
    UPDATE_CATEGORY_VISIBILITY_IN_TRANSACTION_CATEGORY_LIST,
    REMOVE_CATEGORY_FROM_TRANSACTION_CATEGORYLIST,
    UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE,

    LOAD_TRANSACTION_TAG_LIST,
    ADD_TAG_TO_TRANSACTION_TAG_LIST,
    SAVE_TAG_IN_TRANSACTION_TAG_LIST,
    CHANGE_TAG_DISPLAY_ORDER_IN_TRANSACTION_TAG_LIST,
    UPDATE_TAG_VISIBILITY_IN_TRANSACTION_TAG_LIST,
    REMOVE_TAG_FROM_TRANSACTION_TAG_LIST,
    UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE,
} from './mutations.js';

import user from './user.js';
import twoFactorAuth from './twoFactorAuth.js';
import token from './token.js';
import exchangeRates from './exchangeRates.js';
import account from './account.js';
import transaction from './transaction.js';
import transactionCategory from './transactionCategory.js';
import transactionTag from './transactionTag.js';

const stores = {
    strict: process.env.NODE_ENV !== 'production',
    state: {
        currentUserInfo: userState.getUserInfo(),
        latestExchangeRates: exchangeRates.getExchangeRatesFromLocalStorage(),
        allAccounts: [],
        allAccountsMap: {},
        allCategorizedAccounts: {},
        accountListStateInvalid: true,
        transactions: [],
        transactionsNextTimeId: 0,
        transactionListStateInvalid: true,
        allTransactionCategories: {},
        allTransactionCategoriesMap: {},
        transactionCategoryListStateInvalid: true,
        allTransactionTags: [],
        allTransactionTagsMap: {},
        transactionTagListStateInvalid: true,
    },
    getters: {
        currentUserNickname: user.currentUserNickname,
        currentUserDefaultCurrency: user.currentUserDefaultCurrency,
        exchangeRatesLastUpdateDate: exchangeRates.exchangeRatesLastUpdateDate,
        getExchangedAmount: exchangeRates.getExchangedAmount,
        allPlainAccounts: account.allPlainAccounts,
        allVisiblePlainAccounts: account.allVisiblePlainAccounts,
        allAvailableAccountsCount: account.allAvailableAccountsCount,
        allVisibleAccountsCount: account.allVisibleAccountsCount,
        noTransaction: transaction.noTransaction,
        hasMoreTransaction: transaction.hasMoreTransaction,
    },
    mutations: {
        [RESET_STATE] (state) {
            state.latestExchangeRates = {};

            state.allAccounts = [];
            state.allAccountsMap = {};
            state.allCategorizedAccounts = {};
            state.accountListStateInvalid = true;

            state.transactions = [];
            state.transactionsNextTimeId = 0;
            state.transactionListStateInvalid = true;

            state.allTransactionCategories = {};
            state.allTransactionCategoriesMap = {};
            state.transactionCategoryListStateInvalid = true;

            state.allTransactionTags = [];
            state.allTransactionTagsMap = {};
            state.transactionTagListStateInvalid = true;

            exchangeRates.clearExchangeRatesFromLocalStorage();
        },
        [STORE_USER_INFO] (state, userInfo) {
            state.currentUserInfo = userInfo;
            userState.updateUserInfo(userInfo);
        },
        [CLEAR_USER_INFO] (state) {
            state.currentUserInfo = null;
            userState.clearUserInfo();
        },
        [STORE_LATEST_EXCHANGE_RATES] (state, latestExchangeRates) {
            state.latestExchangeRates = latestExchangeRates;
            exchangeRates.setExchangeRatesToLocalStorage(latestExchangeRates);
        },
        [LOAD_ACCOUNT_LIST] (state, accounts) {
            state.allAccounts = accounts;
            state.allAccountsMap = {};

            for (let i = 0; i < accounts.length; i++) {
                const account = accounts[i];
                state.allAccountsMap[account.id] = account;

                if (account.subAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];
                        state.allAccountsMap[subAccount.id] = subAccount;
                    }
                }
            }

            state.allCategorizedAccounts = utils.getCategorizedAccounts(accounts);
        },
        [ADD_ACCOUNT_TO_ACCOUNT_LIST] (state, account) {
            let insertIndexToAllList = 0;

            for (let i = 0; i < state.allAccounts.length; i++) {
                if (state.allAccounts[i].category > account.category) {
                    insertIndexToAllList = i;
                    break;
                }
            }

            state.allAccounts.splice(insertIndexToAllList, 0, account);

            state.allAccountsMap[account.id] = account;

            if (account.subAccounts) {
                for (let i = 0; i < account.subAccounts.length; i++) {
                    const subAccount = account.subAccounts[i];
                    state.allAccountsMap[subAccount.id] = subAccount;
                }
            }

            if (state.allCategorizedAccounts[account.category]) {
                const accountList = state.allCategorizedAccounts[account.category].accounts;
                accountList.push(account);
            } else {
                state.allCategorizedAccounts = utils.getCategorizedAccounts(state.allAccounts);
            }
        },
        [SAVE_ACCOUNT_IN_ACCOUNT_LIST] (state, account) {
            for (let i = 0; i < state.allAccounts.length; i++) {
                if (state.allAccounts[i].id === account.id) {
                    state.allAccounts.splice(i, 1, account);
                    break;
                }
            }

            state.allAccountsMap[account.id] = account;

            if (account.subAccounts) {
                for (let i = 0; i < account.subAccounts.length; i++) {
                    const subAccount = account.subAccounts[i];
                    state.allAccountsMap[subAccount.id] = subAccount;
                }
            }

            if (state.allCategorizedAccounts[account.category]) {
                const accountList = state.allCategorizedAccounts[account.category].accounts;

                for (let i = 0; i < accountList.length; i++) {
                    if (accountList[i].id === account.id) {
                        accountList.splice(i, 1, account);
                        break;
                    }
                }
            }
        },
        [CHANGE_ACCOUNT_DISPLAY_ORDER_IN_ACCOUNT_LIST] (state, { account, from, to }) {
            let fromAccount = null;
            let toAccount = null;

            if (state.allCategorizedAccounts[account.category]) {
                const accountList = state.allCategorizedAccounts[account.category].accounts;
                fromAccount = accountList[from];
                toAccount = accountList[to];

                accountList.splice(to, 0, accountList.splice(from, 1)[0]);
            }

            if (fromAccount && toAccount) {
                let globalFromIndex = -1;
                let globalToIndex = -1;

                for (let i = 0; i < state.allAccounts.length; i++) {
                    if (state.allAccounts[i].id === fromAccount.id) {
                        globalFromIndex = i;
                    } else if (state.allAccounts[i].id === toAccount.id) {
                        globalToIndex = i;
                    }
                }

                if (globalFromIndex >= 0 && globalToIndex >= 0) {
                    state.allAccounts.splice(globalToIndex, 0, state.allAccounts.splice(globalFromIndex, 1)[0]);
                }
            }
        },
        [UPDATE_ACCOUNT_VISIBILITY_IN_ACCOUNT_LIST] (state, { account, hidden }) {
            if (state.allAccountsMap[account.id]) {
                state.allAccountsMap[account.id].hidden = hidden;
            }
        },
        [REMOVE_ACCOUNT_FROM_ACCOUNT_LIST] (state, account) {
            for (let i = 0; i < state.allAccounts.length; i++) {
                if (state.allAccounts[i].id === account.id) {
                    state.allAccounts.splice(i, 1);
                    break;
                }
            }

            if (state.allAccountsMap[account.id] && state.allAccountsMap[account.id].subAccounts) {
                const subAccounts = state.allAccountsMap[account.id].subAccounts;

                for (let i = 0; i < subAccounts.length; i++) {
                    const subAccount = subAccounts[i];
                    if (state.allAccountsMap[subAccount.id]) {
                        delete state.allAccountsMap[subAccount.id];
                    }
                }
            }

            if (state.allAccountsMap[account.id]) {
                delete state.allAccountsMap[account.id];
            }

            if (state.allCategorizedAccounts[account.category]) {
                const accountList = state.allCategorizedAccounts[account.category].accounts;

                for (let i = 0; i < accountList.length; i++) {
                    if (accountList[i].id === account.id) {
                        accountList.splice(i, 1);
                        break;
                    }
                }
            }
        },
        [UPDATE_ACCOUNT_LIST_INVALID_STATE] (state, invalidState) {
            state.accountListStateInvalid = invalidState;
        },
        [LOAD_TRANSACTION_LIST] (state, { transactions, reload, autoExpand, defaultCurrency, accountId }) {
            if (reload) {
                state.transactions = [];
            }

            if (transactions.items && transactions.items.length) {
                let currentMonthListIndex = -1;
                let currentMonthList = null;

                for (let i = 0; i < transactions.items.length; i++) {
                    const item = transactions.items[i];
                    const transactionTime = utils.parseDateFromUnixTime(item.time);

                    item.day = utils.getDay(transactionTime);
                    item.dayOfWeek = utils.getDayOfWeek(transactionTime);
                    item.sourceAccount = state.allAccountsMap[item.sourceAccountId];
                    item.destinationAccount = state.allAccountsMap[item.destinationAccountId];
                    item.category = state.allTransactionCategoriesMap[item.categoryId];
                    item.tags = [];

                    if (item.tagIds && item.tagIds.length) {
                        for (let j = 0; j < item.tagIds.length; j++) {
                            const tag = state.allTransactionTagsMap[item.tagIds[j]];

                            if (tag) {
                                item.tags.push(tag);
                            }
                        }
                    }

                    const transactionYear = utils.getYear(transactionTime);
                    const transactionMonth = utils.getMonth(transactionTime);
                    const transactionYearMonth = utils.getYearAndMonth(transactionTime);

                    if (currentMonthList && currentMonthList.year === transactionYear && currentMonthList.month === transactionMonth) {
                        currentMonthList.items.push(item);
                        transaction.calculateMonthTotalAmount(state, currentMonthList, defaultCurrency, accountId, true);
                        continue;
                    }

                    for (let j = currentMonthListIndex + 1; j < state.transactions.length; j++) {
                        if (state.transactions[j].year === transactionYear && state.transactions[j].month === transactionMonth) {
                            currentMonthListIndex = j;
                            currentMonthList = state.transactions[j];

                            if (j > 0) {
                                transaction.calculateMonthTotalAmount(state, state.transactions[j - 1], defaultCurrency, accountId, false);
                            }

                            break;
                        }
                    }

                    if (!currentMonthList && state.transactions.length > 0) {
                        transaction.calculateMonthTotalAmount(state, state.transactions[state.transactions.length - 1], defaultCurrency, accountId, false);
                    }

                    if (!currentMonthList || currentMonthList.year !== transactionYear || currentMonthList.month !== transactionMonth) {
                        transaction.calculateMonthTotalAmount(state, currentMonthList, defaultCurrency, accountId, false);

                        state.transactions.push({
                            year: transactionYear,
                            month: transactionMonth,
                            yearMonth: transactionYearMonth,
                            opened: autoExpand,
                            items: []
                        });

                        currentMonthListIndex = state.transactions.length - 1;
                        currentMonthList = state.transactions[state.transactions.length - 1];
                    }

                    currentMonthList.items.push(item);
                    transaction.calculateMonthTotalAmount(state, currentMonthList, defaultCurrency, accountId, true);
                }
            }

            if (transactions.nextTimeSequenceId) {
                state.transactionsNextTimeId = transactions.nextTimeSequenceId;
            } else {
                transaction.calculateMonthTotalAmount(state, state.transactions[state.transactions.length - 1], defaultCurrency, accountId, false);
                state.transactionsNextTimeId = -1;
            }
        },
        [COLLAPSE_MONTH_IN_TRANSACTION_LIST] (state, { month, collapse }) {
            if (month) {
                month.opened = !collapse;
            }
        },
        [SAVE_TRANSACTION_IN_TRANSACTION_LIST] (state, { transaction, defaultCurrency, accountId }) {
            for (let i = 0; i < state.transactions.length; i++) {
                const transactionMonthList = state.transactions[i];

                if (!transactionMonthList.items ||
                    transactionMonthList.items[0].time < transaction.time ||
                    transactionMonthList.items[transactionMonthList.items.length - 1].time > transaction.time) {
                    continue;
                }

                for (let j = 0; j < transactionMonthList.items.length; j++) {
                    if (transactionMonthList.items[j].id === transaction.id) {
                        transactionMonthList.items.splice(j, 1, transaction);
                        transaction.calculateMonthTotalAmount(state, transactionMonthList, defaultCurrency, accountId, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
                        return;
                    }
                }
            }
        },
        [REMOVE_TRANSACTION_FROM_TRANSACTION_LIST] (state, { transaction, defaultCurrency, accountId }) {
            for (let i = 0; i < state.transactions.length; i++) {
                const transactionMonthList = state.transactions[i];

                if (!transactionMonthList.items ||
                    transactionMonthList.items[0].time < transaction.time ||
                    transactionMonthList.items[transactionMonthList.items.length - 1].time > transaction.time) {
                    continue;
                }

                for (let j = 0; j < transactionMonthList.items.length; j++) {
                    if (transactionMonthList.items[j].id === transaction.id) {
                        transactionMonthList.items.splice(j, 1);
                    }
                }

                if (transactionMonthList.items.length < 1) {
                    state.transactions.splice(i, 1);
                } else {
                    transaction.calculateMonthTotalAmount(state, transactionMonthList, defaultCurrency, accountId, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
                }
            }
        },
        [UPDATE_TRANSACTION_LIST_INVALID_STATE] (state, invalidState) {
            state.transactionListStateInvalid = invalidState;
        },
        [LOAD_TRANSACTION_CATEGORY_LIST] (state, allCategories) {
            state.allTransactionCategories = allCategories;
            state.allTransactionCategoriesMap = {};

            for (let categoryType in allCategories) {
                if (!Object.prototype.hasOwnProperty.call(allCategories, categoryType)) {
                    continue;
                }

                const categories = allCategories[categoryType];

                for (let i = 0; i < categories.length; i++) {
                    const category = categories[i];
                    state.allTransactionCategoriesMap[category.id] = category;

                    for (let j = 0; j < category.subCategories.length; j++) {
                        const subCategory = category.subCategories[j];
                        state.allTransactionCategoriesMap[subCategory.id] = subCategory;
                    }
                }
            }
        },
        [ADD_CATEGORY_TO_TRANSACTION_CATEGORY_LIST] (state, category) {
            let categoryList = null;

            if (!category.parentId || category.parentId === '0') {
                categoryList = state.allTransactionCategories[category.type];
            } else if (state.allTransactionCategoriesMap[category.parentId]) {
                categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
            }

            if (categoryList) {
                categoryList.push(category);
            }

            state.allTransactionCategoriesMap[category.id] = category;
        },
        [SAVE_CATEGORY_IN_TRANSACTION_CATEGORY_LIST] (state, category) {
            let categoryList = null;

            if (!category.parentId || category.parentId === '0') {
                categoryList = state.allTransactionCategories[category.type];
            } else if (state.allTransactionCategoriesMap[category.parentId]) {
                categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
            }

            if (categoryList) {
                for (let i = 0; i < categoryList.length; i++) {
                    if (categoryList[i].id === category.id) {
                        categoryList.splice(i, 1, category);
                        break;
                    }
                }
            }

            state.allTransactionCategoriesMap[category.id] = category;
        },
        [CHANGE_CATEGORY_DISPLAY_ORDER_IN_CATEGORY_LIST] (state, { category, from, to }) {
            let categoryList = null;

            if (!category.parentId || category.parentId === '0') {
                categoryList = state.allTransactionCategories[category.type];
            } else if (state.allTransactionCategoriesMap[category.parentId]) {
                categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
            }

            if (categoryList) {
                categoryList.splice(to, 0, categoryList.splice(from, 1)[0]);
            }
        },
        [UPDATE_CATEGORY_VISIBILITY_IN_TRANSACTION_CATEGORY_LIST] (state, { category, hidden }) {
            if (state.allTransactionCategoriesMap[category.id]) {
                state.allTransactionCategoriesMap[category.id].hidden = hidden;
            }
        },
        [REMOVE_CATEGORY_FROM_TRANSACTION_CATEGORYLIST] (state, category) {
            let categoryList = null;

            if (!category.parentId || category.parentId === '0') {
                categoryList = state.allTransactionCategories[category.type];
            } else if (state.allTransactionCategoriesMap[category.parentId]) {
                categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
            }

            if (categoryList) {
                for (let i = 0; i < categoryList.length; i++) {
                    if (categoryList[i].id === category.id) {
                        categoryList.splice(i, 1);
                        break;
                    }
                }
            }

            if (state.allTransactionCategoriesMap[category.id] && state.allTransactionCategoriesMap[category.id].subCategories) {
                const subCategories = state.allTransactionCategoriesMap[category.id].subCategories;

                for (let i = 0; i < subCategories.length; i++) {
                    const subCategory = subCategories[i];
                    if (state.allTransactionCategoriesMap[subCategory.id]) {
                        delete state.allTransactionCategoriesMap[subCategory.id];
                    }
                }
            }

            if (state.allTransactionCategoriesMap[category.id]) {
                delete state.allTransactionCategoriesMap[category.id];
            }
        },
        [UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE] (state, invalidState) {
            state.transactionCategoryListStateInvalid = invalidState;
        },
        [LOAD_TRANSACTION_TAG_LIST] (state, tags) {
            state.allTransactionTags = tags;
            state.allTransactionTagsMap = {};

            for (let i = 0; i < tags.length; i++) {
                const tag = tags[i];
                state.allTransactionTagsMap[tag.id] = tag;
            }
        },
        [ADD_TAG_TO_TRANSACTION_TAG_LIST] (state, tag) {
            state.allTransactionTags.push(tag);
            state.allTransactionTagsMap[tag.id] = tag;
        },
        [SAVE_TAG_IN_TRANSACTION_TAG_LIST] (state, tag) {
            for (let i = 0; i < state.allTransactionTags.length; i++) {
                if (state.allTransactionTags[i].id === tag.id) {
                    state.allTransactionTags.splice(i, 1, tag);
                    break;
                }
            }

            state.allTransactionTagsMap[tag.id] = tag;
        },
        [CHANGE_TAG_DISPLAY_ORDER_IN_TRANSACTION_TAG_LIST] (state, { from, to }) {
            state.allTransactionTags.splice(to, 0, state.allTransactionTags.splice(from, 1)[0]);
        },
        [UPDATE_TAG_VISIBILITY_IN_TRANSACTION_TAG_LIST] (state, { tag, hidden }) {
            if (state.allTransactionTagsMap[tag.id]) {
                state.allTransactionTagsMap[tag.id].hidden = hidden;
            }
        },
        [REMOVE_TAG_FROM_TRANSACTION_TAG_LIST] (state, tag) {
            for (let i = 0; i < state.allTransactionTags.length; i++) {
                if (state.allTransactionTags[i].id === tag.id) {
                    state.allTransactionTags.splice(i, 1);
                    break;
                }
            }

            if (state.allTransactionTagsMap[tag.id]) {
                delete state.allTransactionTagsMap[tag.id];
            }
        },
        [UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE] (state, invalidState) {
            state.transactionTagListStateInvalid = invalidState;
        },
    },
    actions: {
        authorize: user.authorize,
        authorize2FA: user.authorize2FA,
        register: user.register,
        logout: user.logout,
        getCurrentUserProfile: user.getCurrentUserProfile,
        updateUserProfile: user.updateUserProfile,
        clearUserInfoState: user.clearUserInfoState,
        resetState: user.resetState,

        get2FAStatus: twoFactorAuth.get2FAStatus,
        enable2FA: twoFactorAuth.enable2FA,
        confirmEnable2FA: twoFactorAuth.confirmEnable2FA,
        disable2FA: twoFactorAuth.disable2FA,
        regenerate2FARecoveryCode: twoFactorAuth.regenerate2FARecoveryCode,

        getAllTokens: token.getAllTokens,
        refreshTokenAndRevokeOldToken: token.refreshTokenAndRevokeOldToken,
        revokeToken: token.revokeToken,
        revokeAllTokens: token.revokeAllTokens,

        getLatestExchangeRates: exchangeRates.getLatestExchangeRates,

        loadAllAccounts: account.loadAllAccounts,
        saveAccount: account.saveAccount,
        getAccount: account.getAccount,
        changeAccountDisplayOrder: account.changeAccountDisplayOrder,
        updateAccountDisplayOrders: account.updateAccountDisplayOrders,
        hideAccount: account.hideAccount,
        deleteAccount: account.deleteAccount,

        getTransactions: transaction.getTransactions,
        getTransaction: transaction.getTransaction,
        saveTransaction: transaction.saveTransaction,
        deleteTransaction: transaction.deleteTransaction,
        collapseMonthInTransactionList: transaction.collapseMonthInTransactionList,

        loadAllCategories: transactionCategory.loadAllCategories,
        getCategory: transactionCategory.getCategory,
        saveCategory: transactionCategory.saveCategory,
        addCategories: transactionCategory.addCategories,
        changeCategoryDisplayOrder: transactionCategory.changeCategoryDisplayOrder,
        updateCategoryDisplayOrders: transactionCategory.updateCategoryDisplayOrders,
        hideCategory: transactionCategory.hideCategory,
        deleteCategory: transactionCategory.deleteCategory,

        loadAllTags: transactionTag.loadAllTags,
        saveTag: transactionTag.saveTag,
        changeTagDisplayOrder: transactionTag.changeTagDisplayOrder,
        updateTagDisplayOrders: transactionTag.updateTagDisplayOrders,
        hideTag: transactionTag.hideTag,
        deleteTag: transactionTag.deleteTag,
    }
};

export default stores;
