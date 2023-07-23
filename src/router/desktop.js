import { createRouter, createWebHashHistory } from 'vue-router';

import userState from '@/lib/userstate.js';

import MainLayout from '@/views/desktop/MainLayout.vue';
import LoginPage from '@/views/desktop/LoginPage.vue';
import SignUpPage from '@/views/desktop/SignupPage.vue';
import UnlockPage from '@/views/desktop/UnlockPage.vue';

import HomePage from '@/views/desktop/HomePage.vue';
import TransactionsPage from '@/views/desktop/TransactionsPage.vue';
import StatisticsTransactionPage from '@/views/desktop/statistics/TransactionPage.vue';
import AccountsPage from '@/views/desktop/AccountsPage.vue';
import TransactionCategoriesPage from '@/views/desktop/TransactionCategoriesPage.vue';
import TransactionTagsPage from '@/views/desktop/TransactionTagsPage.vue';
import ExchangeRatesPage from '@/views/desktop/ExchangeRatesPage.vue';
import UserSettingsPage from '@/views/desktop/user/UserSettingsPage.vue';
import AppSettingsPage from '@/views/desktop/app/AppSettingsPage.vue';
import AboutPage from '@/views/desktop/AboutPage.vue';

function checkLogin() {
    if (!userState.isUserLogined()) {
        return {
            path: '/login',
            replace: true
        };
    }

    if (!userState.isUserUnlocked()) {
        return {
            path: '/unlock',
            replace: true
        };
    }
}

function checkLocked() {
    if (!userState.isUserLogined()) {
        return {
            path: '/login',
            replace: true
        };
    }

    if (userState.isUserUnlocked()) {
        return {
            path: '/',
            replace: true
        };
    }
}

function checkNotLogin() {
    if (userState.isUserLogined() && !userState.isUserUnlocked()) {
        return {
            path: '/unlock',
            replace: true
        };
    }

    if (userState.isUserLogined()) {
        return {
            path: '/',
            replace: true
        };
    }
}

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            component: MainLayout,
            beforeEnter: checkLogin,
            children: [
                {
                    path: '',
                    component: HomePage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/transactions',
                    component: TransactionsPage,
                    beforeEnter: checkLogin,
                    props: route => ({
                        initDateType: route.query.dateType,
                        initMaxTime: route.query.maxTime,
                        initMinTime: route.query.minTime,
                        initType: route.query.type,
                        initCategoryId: route.query.categoryId,
                        initAccountId: route.query.accountId
                    })
                },
                {
                    path: '/statistics/transaction',
                    component: StatisticsTransactionPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/accounts',
                    component: AccountsPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/categories',
                    component: TransactionCategoriesPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/tags',
                    component: TransactionTagsPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/exchange_rates',
                    component: ExchangeRatesPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/user/settings',
                    component: UserSettingsPage,
                    beforeEnter: checkLogin,
                    props: route => ({
                        initTab: route.query.tab
                    })
                },
                {
                    path: '/app/settings',
                    component: AppSettingsPage,
                    beforeEnter: checkLogin,
                    props: route => ({
                        initTab: route.query.tab
                    })
                },
                {
                    path: '/about',
                    component: AboutPage,
                    beforeEnter: checkLogin
                }
            ]
        },
        {
            path: '/login',
            component: LoginPage,
            beforeEnter: checkNotLogin
        },
        {
            path: '/signup',
            component: SignUpPage,
            beforeEnter: checkNotLogin
        },
        {
            path: '/unlock',
            component: UnlockPage,
            beforeEnter: checkLocked
        }
    ],
})

export default router;
