import userState from "../lib/userstate.js";

import HomePage from '../views/mobile/Home.vue';
import LoginPage from '../views/mobile/Login.vue';
import SignUpPage from '../views/mobile/Signup.vue';
import UnlockPage from '../views/mobile/Unlock.vue';

import TransactionDetailPage from '../views/mobile/transactions/Detail.vue';
import TransactionNewPage from '../views/mobile/transactions/New.vue';

import AccountListPage from '../views/mobile/accounts/AccountList.vue';
import AccountEditPage from '../views/mobile/accounts/AccountEdit.vue';

import StatisticsOverviewPage from '../views/mobile/statistics/Overview.vue';

import SettingsPage from '../views/mobile/Settings.vue';
import ApplicationLockPage from '../views/mobile/ApplicationLock.vue';
import ExchangeRatesPage from "../views/mobile/ExchangeRates.vue";
import AboutPage from "../views/mobile/About.vue";

import UserProfilePage from "../views/mobile/users/UserProfile.vue";
import TwoFactorAuthPage from "../views/mobile/users/TwoFactorAuth.vue";
import SessionListPage from "../views/mobile/users/SessionList.vue";

import CategoryAllPage from "../views/mobile/categories/CategoryAll.vue";
import CategoryListPage from "../views/mobile/categories/CategoryList.vue";
import CategoryEditPage from "../views/mobile/categories/CategoryEdit.vue";
import CategoryDefaultPage from "../views/mobile/categories/CategoryDefault.vue";

function checkLogin(to, from, resolve, reject) {
    const router = this;

    if (!userState.isUserLogined()) {
        reject();
        router.navigate('/login', {
            clearPreviousHistory: true,
            pushState: false
        });
        return;
    }

    if (!userState.isUserUnlocked()) {
        reject();
        router.navigate('/unlock', {
            clearPreviousHistory: true,
            pushState: false
        });
        return;
    }

    resolve();
}

function checkLocked(to, from, resolve, reject) {
    const router = this;

    if (!userState.isUserLogined()) {
        reject();
        router.navigate('/login', {
            clearPreviousHistory: true,
            pushState: false
        });
        return;
    }

    if (userState.isUserUnlocked()) {
        reject();
        router.navigate('/', {
            clearPreviousHistory: true,
            pushState: false
        });
        return;
    }

    resolve();
}

function checkNotLogin(to, from, resolve, reject) {
    const router = this;

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
        component: HomePage,
        beforeEnter: checkLogin
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
    },
    {
        path: '/transaction/details',
        component: TransactionDetailPage,
        beforeEnter: checkLogin
    },
    {
        path: '/transaction/new',
        component: TransactionNewPage,
        beforeEnter: checkLogin
    },
    {
        path: '/account/list',
        component: AccountListPage,
        beforeEnter: checkLogin
    },
    {
        path: '/account/add',
        component: AccountEditPage,
        beforeEnter: checkLogin
    },
    {
        path: '/account/edit',
        component: AccountEditPage,
        beforeEnter: checkLogin
    },
    {
        path: '/statistic/overview',
        component: StatisticsOverviewPage,
        beforeEnter: checkLogin
    },
    {
        path: '/settings',
        component: SettingsPage,
        beforeEnter: checkLogin
    },
    {
        path: '/app_lock',
        component: ApplicationLockPage,
        beforeEnter: checkLogin
    },
    {
        path: '/exchange_rates',
        component: ExchangeRatesPage,
        beforeEnter: checkLogin
    },
    {
        path: '/about',
        component: AboutPage,
        beforeEnter: checkLogin
    },
    {
        path: '/user/profile',
        component: UserProfilePage,
        beforeEnter: checkLogin
    },
    {
        path: '/user/2fa',
        component: TwoFactorAuthPage,
        beforeEnter: checkLogin
    },
    {
        path: '/user/sessions',
        component: SessionListPage,
        beforeEnter: checkLogin
    },
    {
        path: '/category/all',
        component: CategoryAllPage,
        beforeEnter: checkLogin
    },
    {
        path: '/category/list',
        component: CategoryListPage,
        beforeEnter: checkLogin
    },
    {
        path: '/category/add',
        component: CategoryEditPage,
        beforeEnter: checkLogin
    },
    {
        path: '/category/edit',
        component: CategoryEditPage,
        beforeEnter: checkLogin
    },
    {
        path: '/category/default',
        component: CategoryDefaultPage,
        beforeEnter: checkLogin
    },
    {
        path: '(.*)',
        redirect: '/'
    }
];

export default routes;
