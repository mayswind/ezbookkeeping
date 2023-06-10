import userState from '@/lib/userstate.js';

import HomePage from '@/views/mobile/HomePage.vue';
import LoginPage from '@/views/mobile/LoginPage.vue';
import SignUpPage from '@/views/mobile/SignupPage.vue';
import UnlockPage from '@/views/mobile/UnlockPage.vue';

import TransactionListPage from '@/views/mobile/transactions/ListPage.vue';
import TransactionEditPage from '@/views/mobile/transactions/EditPage.vue';

import AccountListPage from '@/views/mobile/accounts/ListPage.vue';
import AccountEditPage from '@/views/mobile/accounts/EditPage.vue';

import StatisticsTransactionPage from '@/views/mobile/statistics/TransactionPage.vue';
import StatisticsSettingsPage from '@/views/mobile/statistics/SettingsPage.vue';
import StatisticsAccountFilterSettingsPage from '@/views/mobile/statistics/AccountFilterSettingsPage.vue';
import StatisticsCategoryFilterSettingsPage from '@/views/mobile/statistics/CategoryFilterSettingsPage.vue';

import SettingsPage from '@/views/mobile/SettingsPage.vue';
import ApplicationLockPage from '@/views/mobile/ApplicationLockPage.vue';
import ExchangeRatesPage from '@/views/mobile/ExchangeRatesPage.vue';
import AboutPage from '@/views/mobile/AboutPage.vue';

import UserProfilePage from '@/views/mobile/users/UserProfilePage.vue';
import DataManagementPage from '@/views/mobile/users/DataManagementPage.vue';
import TwoFactorAuthPage from '@/views/mobile/users/TwoFactorAuthPage.vue';
import SessionListPage from '@/views/mobile/users/SessionListPage.vue';

import CategoryAllPage from '@/views/mobile/categories/AllPage.vue';
import CategoryListPage from '@/views/mobile/categories/ListPage.vue';
import CategoryEditPage from '@/views/mobile/categories/EditPage.vue';
import CategoryPresetPage from '@/views/mobile/categories/PresetPage.vue';

import TagListPage from '@/views/mobile/tags/ListPage.vue';

function asyncResolve(component) {
    return function({ resolve }) {
        return resolve({
            component: component
        });
    };
}

function checkLogin({ router, resolve, reject }) {
    if (!userState.isUserLogined()) {
        reject();
        router.navigate('/login', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    if (!userState.isUserUnlocked()) {
        reject();
        router.navigate('/unlock', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    resolve();
}

function checkLocked({ router, resolve, reject }) {
    if (!userState.isUserLogined()) {
        reject();
        router.navigate('/login', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    if (userState.isUserUnlocked()) {
        reject();
        router.navigate('/', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    resolve();
}

function checkNotLogin({ router, resolve, reject }) {
    if (userState.isUserLogined() && !userState.isUserUnlocked()) {
        reject();
        router.navigate('/unlock', {
            clearPreviousHistory: true,
            pushState: false
        });
        return;
    }

    if (userState.isUserLogined()) {
        reject();
        router.navigate('/', {
            clearPreviousHistory: true,
            pushState: false
        });
        return;
    }

    resolve();
}

const routes = [
    {
        path: '/',
        async: asyncResolve(HomePage),
        beforeEnter: [checkLogin],
        options: {
            animate: false,
        }
    },
    {
        path: '/login',
        async: asyncResolve(LoginPage),
        beforeEnter: [checkNotLogin],
        options: {
            animate: false,
        }
    },
    {
        path: '/signup',
        async: asyncResolve(SignUpPage),
        beforeEnter: [checkNotLogin],
        options: {
            animate: false,
        }
    },
    {
        path: '/unlock',
        async: asyncResolve(UnlockPage),
        beforeEnter: [checkLocked],
        options: {
            animate: false,
        }
    },
    {
        path: '/transaction/list',
        async: asyncResolve(TransactionListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/transaction/add',
        async: asyncResolve(TransactionEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/transaction/edit',
        async: asyncResolve(TransactionEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/transaction/detail',
        async: asyncResolve(TransactionEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/account/list',
        async: asyncResolve(AccountListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/account/add',
        async: asyncResolve(AccountEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/account/edit',
        async: asyncResolve(AccountEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/transaction',
        async: asyncResolve(StatisticsTransactionPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/settings',
        async: asyncResolve(StatisticsSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/filter/account',
        async: asyncResolve(StatisticsAccountFilterSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/filter/category',
        async: asyncResolve(StatisticsCategoryFilterSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings',
        async: asyncResolve(SettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/app_lock',
        async: asyncResolve(ApplicationLockPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/exchange_rates',
        async: asyncResolve(ExchangeRatesPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/about',
        async: asyncResolve(AboutPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/profile',
        async: asyncResolve(UserProfilePage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/data/management',
        async: asyncResolve(DataManagementPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/2fa',
        async: asyncResolve(TwoFactorAuthPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/sessions',
        async: asyncResolve(SessionListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/all',
        async: asyncResolve(CategoryAllPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/list',
        async: asyncResolve(CategoryListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/add',
        async: asyncResolve(CategoryEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/edit',
        async: asyncResolve(CategoryEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/preset',
        async: asyncResolve(CategoryPresetPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/tag/list',
        async: asyncResolve(TagListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '(.*)',
        redirect: '/'
    }
];

export default routes;
