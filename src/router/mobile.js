import userState from '../lib/userstate.js';

import HomePage from '../views/mobile/Home.vue';
import LoginPage from '../views/mobile/Login.vue';
import SignUpPage from '../views/mobile/Signup.vue';
import UnlockPage from '../views/mobile/Unlock.vue';

import TransactionListPage from '../views/mobile/transactions/List.vue';
import TransactionEditPage from '../views/mobile/transactions/Edit.vue';

import AccountListPage from '../views/mobile/accounts/List.vue';
import AccountEditPage from '../views/mobile/accounts/Edit.vue';

import StatisticsTransactionPage from '../views/mobile/statistics/Transaction.vue';
import StatisticsSettingsPage from '../views/mobile/statistics/Settings.vue';
import StatisticsAccountFilterSettingsPage from '../views/mobile/statistics/AccountFilterSettings.vue';
import StatisticsCategoryFilterSettingsPage from '../views/mobile/statistics/CategoryFilterSettings.vue';

import SettingsPage from '../views/mobile/Settings.vue';
import ApplicationLockPage from '../views/mobile/ApplicationLock.vue';
import ExchangeRatesPage from '../views/mobile/ExchangeRates.vue';
import AboutPage from '../views/mobile/About.vue';

import UserProfilePage from '../views/mobile/users/UserProfile.vue';
import DataManagementPage from '../views/mobile/users/DataManagement.vue';
import TwoFactorAuthPage from '../views/mobile/users/TwoFactorAuth.vue';
import SessionListPage from '../views/mobile/users/SessionList.vue';

import CategoryAllPage from '../views/mobile/categories/All.vue';
import CategoryListPage from '../views/mobile/categories/List.vue';
import CategoryEditPage from '../views/mobile/categories/Edit.vue';
import CategoryPresetPage from '../views/mobile/categories/Preset.vue';

import TagListPage from '../views/mobile/tags/List.vue';

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
