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

    LOAD_TRANSACTION_CATEGORY_LIST,
    UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE,

    LOAD_TRANSACTION_TAG_LIST,
    ADD_TAG_TO_TRANSACTION_TAG_LIST,
    SAVE_TAG_IN_TRANSACTION_TAG_LIST,
    CHANGE_TAG_DISPLAY_ORDER_IN_TRANSACTION_TAG_LIST,
    UPDATE_TAG_VISIBILITY_IN_TRANSACTION_TAG_LIST,
    REMOVE_TAG_FROM_TRANSACTION_TAG_LIST,
    UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE
} from './mutations.js';

import user from './user.js';
import twoFactorAuth from './twoFactorAuth.js';
import token from './token.js';
import exchangeRates from './exchangeRates.js';
import account from './account.js';
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
        allTransactionCategories: [],
        allTransactionCategoriesMap: {},
        transactionCategoryListStateInvalid: true,
        allTransactionTags: [],
        allTransactionTagsMap: {},
        transactionTagListStateInvalid: true,
        transactions: [],
    },
    getters: {
        currentUserNickname: user.currentUserNickname,
        currentUserDefaultCurrency: user.currentUserDefaultCurrency,
        exchangeRatesLastUpdateDate: exchangeRates.exchangeRatesLastUpdateDate,
        getExchangedAmount: exchangeRates.getExchangedAmount,
        allAvailableAccountsCount: account.allAvailableAccountsCount,
        allVisibleAccountsCount: account.allVisibleAccountsCount,
    },
    mutations: {
        [RESET_STATE] (state) {
            state.latestExchangeRates = {};
            state.allAccounts = [];
            state.allAccountsMap = {};
            state.allCategorizedAccounts = {};
            state.accountListStateInvalid = true;
            state.allTransactionCategories = [];
            state.allTransactionCategoriesMap = {};
            state.allTransactionTags = [];
            state.allTransactionTagsMap = {};
            state.transactionTagListStateInvalid = true;
            state.transactions = [];

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
        [LOAD_TRANSACTION_CATEGORY_LIST] (state, categories) {
            state.allTransactionCategories = categories;
            state.allTransactionCategoriesMap = {};

            for (let i = 0; i < categories.length; i++) {
                const category = categories[i];
                state.allTransactionCategoriesMap[category.id] = category;
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

        addTransactionCategoryBatch: transactionCategory.addTransactionCategoryBatch,

        loadAllTags: transactionTag.loadAllTags,
        saveTag: transactionTag.saveTag,
        changeTagDisplayOrder: transactionTag.changeTagDisplayOrder,
        updateTagDisplayOrders: transactionTag.updateTagDisplayOrders,
        hideTag: transactionTag.hideTag,
        deleteTag: transactionTag.deleteTag,
    }
};

export default stores;
