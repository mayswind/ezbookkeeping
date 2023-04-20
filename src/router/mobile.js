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
        async: ({resolve}) => resolve({component: HomePage}),
        beforeEnter: [checkLogin],
        options: {
            animate: false,
        }
    },
    {
        path: '/login',
        async: ({resolve}) => resolve({component: LoginPage}),
        beforeEnter: [checkNotLogin],
        options: {
            animate: false,
        }
    },
    {
        path: '/signup',
        async: ({resolve}) => resolve({component: SignUpPage}),
        beforeEnter: [checkNotLogin],
        options: {
            animate: false,
        }
    },
    {
        path: '/unlock',
        async: ({resolve}) => resolve({component: UnlockPage}),
        beforeEnter: [checkLocked],
        options: {
            animate: false,
        }
    },
    {
        path: '/transaction/list',
        async: ({resolve}) => resolve({component: TransactionListPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/transaction/add',
        async: ({resolve}) => resolve({component: TransactionEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/transaction/edit',
        async: ({resolve}) => resolve({component: TransactionEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/transaction/detail',
        async: ({resolve}) => resolve({component: TransactionEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/account/list',
        async: ({resolve}) => resolve({component: AccountListPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/account/add',
        async: ({resolve}) => resolve({component: AccountEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/account/edit',
        async: ({resolve}) => resolve({component: AccountEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/transaction',
        async: ({resolve}) => resolve({component: StatisticsTransactionPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/settings',
        async: ({resolve}) => resolve({component: StatisticsSettingsPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/filter/account',
        async: ({resolve}) => resolve({component: StatisticsAccountFilterSettingsPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/statistic/filter/category',
        async: ({resolve}) => resolve({component: StatisticsCategoryFilterSettingsPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings',
        async: ({resolve}) => resolve({component: SettingsPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/app_lock',
        async: ({resolve}) => resolve({component: ApplicationLockPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/exchange_rates',
        async: ({resolve}) => resolve({component: ExchangeRatesPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/about',
        async: ({resolve}) => resolve({component: AboutPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/profile',
        async: ({resolve}) => resolve({component: UserProfilePage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/data/management',
        async: ({resolve}) => resolve({component: DataManagementPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/2fa',
        async: ({resolve}) => resolve({component: TwoFactorAuthPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/user/sessions',
        async: ({resolve}) => resolve({component: SessionListPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/all',
        async: ({resolve}) => resolve({component: CategoryAllPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/list',
        async: ({resolve}) => resolve({component: CategoryListPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/add',
        async: ({resolve}) => resolve({component: CategoryEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/edit',
        async: ({resolve}) => resolve({component: CategoryEditPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/category/preset',
        async: ({resolve}) => resolve({component: CategoryPresetPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '/tag/list',
        async: ({resolve}) => resolve({component: TagListPage}),
        beforeEnter: [checkLogin]
    },
    {
        path: '(.*)',
        redirect: '/'
    }
];

export default routes;
