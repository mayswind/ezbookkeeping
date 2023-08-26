import { createRouter, createWebHashHistory } from 'vue-router';

import userState from '@/lib/userstate.js';

import MainLayout from '@/views/desktop/MainLayout.vue';
import LoginPage from '@/views/desktop/LoginPage.vue';
import SignUpPage from '@/views/desktop/SignupPage.vue';
import ForgetPasswordPage from '@/views/desktop/ForgetPasswordPage.vue';
import ResetPasswordPage from '@/views/desktop/ResetPasswordPage.vue';
import UnlockPage from '@/views/desktop/UnlockPage.vue';

import HomePage from '@/views/desktop/HomePage.vue';

import TransactionListPage from '@/views/desktop/transactions/ListPage.vue';

import StatisticsTransactionPage from '@/views/desktop/statistics/TransactionPage.vue';

import AccountListPage from '@/views/desktop/accounts/ListPage.vue';

import TransactionCategoryListPage from '@/views/desktop/categories/ListPage.vue';

import TransactionTagListPage from '@/views/desktop/tags/ListPage.vue';

import UserSettingsPage from '@/views/desktop/user/UserSettingsPage.vue';
import AppSettingsPage from '@/views/desktop/app/AppSettingsPage.vue';

import ExchangeRatesPage from '@/views/desktop/ExchangeRatesPage.vue';
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
                    path: '/transaction/list',
                    component: TransactionListPage,
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
                    path: '/account/list',
                    component: AccountListPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/category/list',
                    component: TransactionCategoryListPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/tag/list',
                    component: TransactionTagListPage,
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
            path: '/forgetpassword',
            component: ForgetPasswordPage,
            beforeEnter: checkNotLogin
        },
        {
            path: '/resetpassword',
            component: ResetPasswordPage,
            props: route => ({
                token: route.query.token
            })
        },
        {
            path: '/unlock',
            component: UnlockPage,
            beforeEnter: checkLocked
        }
    ],
})

export default router;
