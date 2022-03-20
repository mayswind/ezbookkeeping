import datetimeConstants from '../consts/datetime.js';
import currencyConstants from '../consts/currency.js';
import statisticsConstants from '../consts/statistics.js';
import userState from '../lib/userstate.js';
import settings from '../lib/settings.js';
import utils from '../lib/utils.js';

import {
    RESET_STATE,

    UPDATE_DEFAULT_SETTING,

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
    INIT_TRANSACTION_LIST_FILTER,
    UPDATE_TRANSACTION_LIST_FILTER,
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

    LOAD_TRANSACTION_OVERVIEW,
    UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE,

    LOAD_TRANSACTION_STATISTICS,
    INIT_TRANSACTION_STATISTICS_FILTER,
    UPDATE_TRANSACTION_STATISTICS_FILTER,
    UPDATE_TRANSACTION_STATISTICS_INVALID_STATE,
} from './mutations.js';

import {
    updateLocalizedDefaultSettings
} from './setting.js';

import {
    authorize,
    authorize2FA,
    register,
    logout,
    getCurrentUserProfile,
    updateUserProfile,
    clearUserData,
    clearUserInfoState,
    resetState,
    currentUserNickname,
    currentUserDefaultCurrency,
    currentUserFirstDayOfWeek,
} from './user.js';

import {
    get2FAStatus,
    enable2FA,
    confirmEnable2FA,
    disable2FA,
    regenerate2FARecoveryCode,
} from './twoFactorAuth.js';

import {
    getAllTokens,
    refreshTokenAndRevokeOldToken,
    revokeToken,
    revokeAllTokens,
} from './token.js';

import {
    getLatestExchangeRates,
    exchangeRatesLastUpdateTime,
    getExchangedAmount,
    getExchangeRatesFromLocalStorage,
    setExchangeRatesToLocalStorage,
    clearExchangeRatesFromLocalStorage,
} from './exchangeRates.js';

import {
    loadTransactionOverview
} from './overview.js';

import {
    loadTransactionStatistics,
    initTransactionStatisticsFilter,
    updateTransactionStatisticsFilter,
    statisticsItemsByTransactionStatisticsData,
    statisticsItemsByAccountsData,
} from './statistics.js';

import {
    loadAllAccounts,
    getAccount,
    saveAccount,
    changeAccountDisplayOrder,
    updateAccountDisplayOrders,
    hideAccount,
    deleteAccount,
    allPlainAccounts,
    allVisiblePlainAccounts,
    allAvailableAccountsCount,
    allVisibleAccountsCount,
} from './account.js';

import {
    initTransactionListFilter,
    updateTransactionListFilter,
    loadTransactions,
    getTransaction,
    saveTransaction,
    deleteTransaction,
    collapseMonthInTransactionList,
    noTransaction,
    hasMoreTransaction,
    fillTransactionObject,
    calculateMonthTotalAmount,
} from './transaction.js';

import {
    loadAllCategories,
    getCategory,
    saveCategory,
    addCategories,
    changeCategoryDisplayOrder,
    updateCategoryDisplayOrders,
    hideCategory,
    deleteCategory,
} from './transactionCategory.js';

import {
    loadAllTags,
    saveTag,
    changeTagDisplayOrder,
    updateTagDisplayOrders,
    hideTag,
    deleteTag,
} from './transactionTag.js';

const stores = {
    strict: !settings.isProduction(),
    state: {
        defaultSetting: {
            currency: currencyConstants.defaultCurrency,
            firstDayOfWeek: datetimeConstants.defaultFirstDayOfWeek
        },
        currentUserInfo: userState.getUserInfo(),
        latestExchangeRates: getExchangeRatesFromLocalStorage(),
        allAccounts: [],
        allAccountsMap: {},
        allCategorizedAccounts: {},
        accountListStateInvalid: true,
        transactionsFilter: {
            dateType: datetimeConstants.allDateRanges.All.type,
            maxTime: 0,
            minTime: 0,
            type: 0,
            categoryId: '0',
            accountId: '0',
            keyword: ''
        },
        transactions: [],
        transactionsNextTimeId: 0,
        transactionListStateInvalid: true,
        allTransactionCategories: {},
        allTransactionCategoriesMap: {},
        transactionCategoryListStateInvalid: true,
        allTransactionTags: [],
        allTransactionTagsMap: {},
        transactionTagListStateInvalid: true,
        transactionOverview: {},
        transactionOverviewStateInvalid: true,
        transactionStatisticsFilter: {
            dateType: statisticsConstants.defaultDataRangeType,
            startTime: 0,
            endTime: 0,
            chartType: statisticsConstants.defaultChartType,
            chartDataType: statisticsConstants.defaultChartDataType,
            filterAccountIds: {},
            filterCategoryIds: {}
        },
        transactionStatistics: [],
        transactionStatisticsStateInvalid: true,
    },
    getters: {
        // user
        currentUserNickname,
        currentUserDefaultCurrency,
        currentUserFirstDayOfWeek,

        // exchange rates
        exchangeRatesLastUpdateTime,
        getExchangedAmount,

        // statistics
        statisticsItemsByTransactionStatisticsData,
        statisticsItemsByAccountsData,

        // account
        allPlainAccounts,
        allVisiblePlainAccounts,
        allAvailableAccountsCount,
        allVisibleAccountsCount,

        // transaction
        noTransaction,
        hasMoreTransaction,
    },
    mutations: {
        [RESET_STATE] (state) {
            state.latestExchangeRates = {};

            state.allAccounts = [];
            state.allAccountsMap = {};
            state.allCategorizedAccounts = {};
            state.accountListStateInvalid = true;

            state.transactionsFilter.dateType = datetimeConstants.allDateRanges.All.type;
            state.transactionsFilter.maxTime = 0;
            state.transactionsFilter.minTime = 0;
            state.transactionsFilter.type = 0;
            state.transactionsFilter.categoryId = '0';
            state.transactionsFilter.accountId = '0';
            state.transactionsFilter.keyword = '';
            state.transactions = [];
            state.transactionsNextTimeId = 0;
            state.transactionListStateInvalid = true;

            state.allTransactionCategories = {};
            state.allTransactionCategoriesMap = {};
            state.transactionCategoryListStateInvalid = true;

            state.allTransactionTags = [];
            state.allTransactionTagsMap = {};
            state.transactionTagListStateInvalid = true;

            state.transactionOverview = {};
            state.transactionOverviewStateInvalid = true;

            state.transactionStatisticsFilter.dateType = statisticsConstants.defaultDataRangeType;
            state.transactionStatisticsFilter.startTime = 0;
            state.transactionStatisticsFilter.endTime = 0;
            state.transactionStatisticsFilter.chartType = statisticsConstants.defaultChartType;
            state.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            state.transactionStatisticsFilter.filterAccountIds = {};
            state.transactionStatisticsFilter.filterCategoryIds = {};
            state.transactionStatistics = {};
            state.transactionStatisticsStateInvalid = true;

            clearExchangeRatesFromLocalStorage();
        },
        [UPDATE_DEFAULT_SETTING] (state, { defaultCurrency, defaultFirstDayOfWeek }) {
            state.defaultSetting.currency = defaultCurrency;
            state.defaultSetting.firstDayOfWeek = defaultFirstDayOfWeek;
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
            setExchangeRatesToLocalStorage(latestExchangeRates);
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
        [LOAD_TRANSACTION_LIST] (state, { transactions, reload, autoExpand, defaultCurrency }) {
            if (reload) {
                state.transactions = [];
            }

            if (transactions.items && transactions.items.length) {
                const currentUtcOffset = utils.getTimezoneOffsetMinutes();
                let currentMonthListIndex = -1;
                let currentMonthList = null;

                for (let i = 0; i < transactions.items.length; i++) {
                    const item = transactions.items[i];
                    fillTransactionObject(state, item, currentUtcOffset);

                    const transactionTime = utils.parseDateFromUnixTime(item.time, item.utcOffset, currentUtcOffset);
                    const transactionYear = utils.getYear(transactionTime);
                    const transactionMonth = utils.getMonth(transactionTime);
                    const transactionYearMonth = utils.getYearAndMonth(transactionTime);

                    if (currentMonthList && currentMonthList.year === transactionYear && currentMonthList.month === transactionMonth) {
                        currentMonthList.items.push(Object.freeze(item));
                        continue;
                    }

                    for (let j = currentMonthListIndex + 1; j < state.transactions.length; j++) {
                        if (state.transactions[j].year === transactionYear && state.transactions[j].month === transactionMonth) {
                            currentMonthListIndex = j;
                            currentMonthList = state.transactions[j];
                            break;
                        }
                    }

                    if (!currentMonthList || currentMonthList.year !== transactionYear || currentMonthList.month !== transactionMonth) {
                        calculateMonthTotalAmount(state, currentMonthList, defaultCurrency, state.transactionsFilter.accountId, false);

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

                    currentMonthList.items.push(Object.freeze(item));
                    calculateMonthTotalAmount(state, currentMonthList, defaultCurrency, state.transactionsFilter.accountId, true);
                }
            }

            if (transactions.nextTimeSequenceId) {
                state.transactionsNextTimeId = transactions.nextTimeSequenceId;
            } else {
                calculateMonthTotalAmount(state, state.transactions[state.transactions.length - 1], defaultCurrency, state.transactionsFilter.accountId, false);
                state.transactionsNextTimeId = -1;
            }
        },
        [INIT_TRANSACTION_LIST_FILTER] (state, filter) {
            if (filter && utils.isNumber(filter.dateType)) {
                state.transactionsFilter.dateType = filter.dateType;
            } else {
                state.transactionsFilter.dateType = datetimeConstants.allDateRanges.All.type;
            }

            if (filter && utils.isNumber(filter.maxTime)) {
                state.transactionsFilter.maxTime = filter.maxTime;
            } else {
                state.transactionsFilter.maxTime = 0;
            }

            if (filter && utils.isNumber(filter.minTime)) {
                state.transactionsFilter.minTime = filter.minTime;
            } else {
                state.transactionsFilter.minTime = 0;
            }

            if (filter && utils.isNumber(filter.type)) {
                state.transactionsFilter.type = filter.type;
            } else {
                state.transactionsFilter.type = 0;
            }

            if (filter && utils.isString(filter.categoryId)) {
                state.transactionsFilter.categoryId = filter.categoryId;
            } else {
                state.transactionsFilter.categoryId = '0';
            }

            if (filter && utils.isString(filter.accountId)) {
                state.transactionsFilter.accountId = filter.accountId;
            } else {
                state.transactionsFilter.accountId = '0';
            }

            if (filter && utils.isString(filter.keyword)) {
                state.transactionsFilter.keyword = filter.keyword;
            } else {
                state.transactionsFilter.keyword = '';
            }
        },
        [UPDATE_TRANSACTION_LIST_FILTER] (state, filter) {
            if (filter && utils.isNumber(filter.dateType)) {
                state.transactionsFilter.dateType = filter.dateType;
            }

            if (filter && utils.isNumber(filter.maxTime)) {
                state.transactionsFilter.maxTime = filter.maxTime;
            }

            if (filter && utils.isNumber(filter.minTime)) {
                state.transactionsFilter.minTime = filter.minTime;
            }

            if (filter && utils.isNumber(filter.type)) {
                state.transactionsFilter.type = filter.type;
            }

            if (filter && utils.isString(filter.categoryId)) {
                state.transactionsFilter.categoryId = filter.categoryId;
            }

            if (filter && utils.isString(filter.accountId)) {
                state.transactionsFilter.accountId = filter.accountId;
            }

            if (filter && utils.isString(filter.keyword)) {
                state.transactionsFilter.keyword = filter.keyword;
            }
        },
        [COLLAPSE_MONTH_IN_TRANSACTION_LIST] (state, { month, collapse }) {
            if (month) {
                month.opened = !collapse;
            }
        },
        [SAVE_TRANSACTION_IN_TRANSACTION_LIST] (state, { transaction, defaultCurrency }) {
            const currentUtcOffset = utils.getTimezoneOffsetMinutes();
            const transactionTime = utils.parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentUtcOffset);
            const transactionYear = utils.getYear(transactionTime);
            const transactionMonth = utils.getMonth(transactionTime);

            for (let i = 0; i < state.transactions.length; i++) {
                const transactionMonthList = state.transactions[i];

                if (!transactionMonthList.items) {
                    continue;
                }

                for (let j = 0; j < transactionMonthList.items.length; j++) {
                    if (transactionMonthList.items[j].id === transaction.id) {
                        fillTransactionObject(state, transaction, currentUtcOffset);

                        if (transactionYear !== transactionMonthList.year ||
                            transactionMonth !== transactionMonthList.month ||
                            transaction.day !== transactionMonthList.items[j].day) {
                            state.transactionListStateInvalid = true;
                            return;
                        }

                        if ((state.transactionsFilter.categoryId && state.transactionsFilter.categoryId !== '0' && state.transactionsFilter.categoryId !== transaction.categoryId) ||
                            (state.transactionsFilter.accountId && state.transactionsFilter.accountId !== '0' &&
                                state.transactionsFilter.accountId !== transaction.sourceAccountId &&
                                state.transactionsFilter.accountId !== transaction.destinationAccountId &&
                                (!transaction.sourceAccount || state.transactionsFilter.accountId !== transaction.sourceAccount.parentId) &&
                                (!transaction.destinationAccount || state.transactionsFilter.accountId !== transaction.destinationAccount.parentId)
                            )
                        ) {
                            transactionMonthList.items.splice(j, 1);
                        } else {
                            transactionMonthList.items.splice(j, 1, transaction);
                        }

                        if (transactionMonthList.items.length < 1) {
                            state.transactions.splice(i, 1);
                        } else {
                            calculateMonthTotalAmount(state, transactionMonthList, defaultCurrency, state.transactionsFilter.accountId, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
                        }

                        return;
                    }
                }
            }
        },
        [REMOVE_TRANSACTION_FROM_TRANSACTION_LIST] (state, { transaction, defaultCurrency }) {
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
                    calculateMonthTotalAmount(state, transactionMonthList, defaultCurrency, state.transactionsFilter.accountId, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
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
        [LOAD_TRANSACTION_OVERVIEW] (state, transactionOverview) {
            state.transactionOverview = transactionOverview;
        },
        [UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE] (state, invalidState) {
            state.transactionOverviewStateInvalid = invalidState;
        },
        [LOAD_TRANSACTION_STATISTICS] (state, { statistics, defaultCurrency }) {
            if (statistics && statistics.items && statistics.items.length) {
                for (let i = 0; i < statistics.items.length; i++) {
                    const item = statistics.items[i];

                    if (item.accountId) {
                        item.account = state.allAccountsMap[item.accountId];
                    }

                    if (item.account && item.account.parentId !== '0') {
                        item.primaryAccount = state.allAccountsMap[item.account.parentId];
                    } else {
                        item.primaryAccount = item.account;
                    }

                    if (item.categoryId) {
                        item.category = state.allTransactionCategoriesMap[item.categoryId];
                    }

                    if (item.category && item.category.parentId !== '0') {
                        item.primaryCategory = state.allTransactionCategoriesMap[item.category.parentId];
                    } else {
                        item.primaryCategory = item.category;
                    }

                    if (item.account && item.account.currency !== defaultCurrency) {
                        const amount = getExchangedAmount(state)(item.amount, item.account.currency, defaultCurrency);

                        if (utils.isNumber(amount)) {
                            item.amountInDefaultCurrency = Math.floor(amount);
                        }
                    } else if (item.account && item.account.currency === defaultCurrency) {
                        item.amountInDefaultCurrency = item.amount;
                    } else {
                        item.amountInDefaultCurrency = null;
                    }
                }
            }

            state.transactionStatistics = statistics;
        },
        [INIT_TRANSACTION_STATISTICS_FILTER] (state, filter) {
            if (filter && utils.isNumber(filter.dateType)) {
                state.transactionStatisticsFilter.dateType = filter.dateType;
            } else {
                state.transactionStatisticsFilter.dateType = statisticsConstants.defaultDataRangeType;
            }

            if (filter && utils.isNumber(filter.startTime)) {
                state.transactionStatisticsFilter.startTime = filter.startTime;
            } else {
                state.transactionStatisticsFilter.startTime = 0;
            }

            if (filter && utils.isNumber(filter.endTime)) {
                state.transactionStatisticsFilter.endTime = filter.endTime;
            } else {
                state.transactionStatisticsFilter.endTime = 0;
            }

            if (filter && utils.isNumber(filter.chartType)) {
                state.transactionStatisticsFilter.chartType = filter.chartType;
            } else {
                state.transactionStatisticsFilter.chartType = statisticsConstants.defaultChartType;
            }

            if (filter && utils.isNumber(filter.chartDataType)) {
                state.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            } else {
                state.transactionStatisticsFilter.chartDataType = statisticsConstants.defaultChartDataType;
            }

            if (filter && utils.isObject(filter.filterAccountIds)) {
                state.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            } else {
                state.transactionStatisticsFilter.filterAccountIds = {};
            }

            if (filter && utils.isObject(filter.filterCategoryIds)) {
                state.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            } else {
                state.transactionStatisticsFilter.filterCategoryIds = {};
            }
        },
        [UPDATE_TRANSACTION_STATISTICS_FILTER] (state, filter) {
            if (filter && utils.isNumber(filter.dateType)) {
                state.transactionStatisticsFilter.dateType = filter.dateType;
            }

            if (filter && utils.isNumber(filter.startTime)) {
                state.transactionStatisticsFilter.startTime = filter.startTime;
            }

            if (filter && utils.isNumber(filter.endTime)) {
                state.transactionStatisticsFilter.endTime = filter.endTime;
            }

            if (filter && utils.isNumber(filter.chartType)) {
                state.transactionStatisticsFilter.chartType = filter.chartType;
            }

            if (filter && utils.isNumber(filter.chartDataType)) {
                state.transactionStatisticsFilter.chartDataType = filter.chartDataType;
            }

            if (filter && utils.isObject(filter.filterAccountIds)) {
                state.transactionStatisticsFilter.filterAccountIds = filter.filterAccountIds;
            }

            if (filter && utils.isObject(filter.filterCategoryIds)) {
                state.transactionStatisticsFilter.filterCategoryIds = filter.filterCategoryIds;
            }
        },
        [UPDATE_TRANSACTION_STATISTICS_INVALID_STATE] (state, invalidState) {
            state.transactionStatisticsStateInvalid = invalidState;
        },
    },
    actions: {
        // setting
        updateLocalizedDefaultSettings,

        // user
        authorize,
        authorize2FA,
        register,
        logout,
        getCurrentUserProfile,
        updateUserProfile,
        clearUserData,
        clearUserInfoState,
        resetState,

        // 2fa
        get2FAStatus,
        enable2FA,
        confirmEnable2FA,
        disable2FA,
        regenerate2FARecoveryCode,

        // token
        getAllTokens,
        refreshTokenAndRevokeOldToken,
        revokeToken,
        revokeAllTokens,

        // exchange rates
        getLatestExchangeRates,

        // overview
        loadTransactionOverview,

        // statistics
        loadTransactionStatistics,
        initTransactionStatisticsFilter,
        updateTransactionStatisticsFilter,

        // account
        loadAllAccounts,
        saveAccount,
        getAccount,
        changeAccountDisplayOrder,
        updateAccountDisplayOrders,
        hideAccount,
        deleteAccount,

        // transaction
        initTransactionListFilter,
        updateTransactionListFilter,
        loadTransactions,
        getTransaction,
        saveTransaction,
        deleteTransaction,
        collapseMonthInTransactionList,

        // transaction category
        loadAllCategories,
        getCategory,
        saveCategory,
        addCategories,
        changeCategoryDisplayOrder,
        updateCategoryDisplayOrders,
        hideCategory,
        deleteCategory,

        // transaction tag
        loadAllTags,
        saveTag,
        changeTagDisplayOrder,
        updateTagDisplayOrders,
        hideTag,
        deleteTag,
    }
};

export default stores;
