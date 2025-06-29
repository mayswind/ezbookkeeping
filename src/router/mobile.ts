import type { Router } from 'framework7/types';

import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';

import HomePage from '@/views/mobile/HomePage.vue';
import LoginPage from '@/views/mobile/LoginPage.vue';
import SignUpPage from '@/views/mobile/SignupPage.vue';
import UnlockPage from '@/views/mobile/UnlockPage.vue';

import TransactionListPage from '@/views/mobile/transactions/ListPage.vue';
import TransactionEditPage from '@/views/mobile/transactions/EditPage.vue';
import TransactionAmountFilterPage from '@/views/mobile/transactions/AmountFilterPage.vue';

import AccountListPage from '@/views/mobile/accounts/ListPage.vue';
import AccountEditPage from '@/views/mobile/accounts/EditPage.vue';

import StatisticsTransactionPage from '@/views/mobile/statistics/TransactionPage.vue';
import StatisticsSettingsPage from '@/views/mobile/statistics/SettingsPage.vue';

import TextSizeSettingsPage from '@/views/mobile/settings/TextSizeSettingsPage.vue';
import PageSettingsPage from '@/views/mobile/settings/PageSettingsPage.vue';
import ApplicationCloudSyncSettingsPage from '@/views/mobile/settings/ApplicationCloudSyncSettingsPage.vue';
import AccountFilterSettingsPage from '@/views/mobile/settings/AccountFilterSettingsPage.vue';
import CategoryFilterSettingsPage from '@/views/mobile/settings/CategoryFilterSettingsPage.vue';
import TransactionTagFilterSettingsPage from '@/views/mobile/settings/TransactionTagFilterSettingsPage.vue';

import SettingsPage from '@/views/mobile/SettingsPage.vue';
import ApplicationLockPage from '@/views/mobile/ApplicationLockPage.vue';
import ExchangeRatesListPage from '@/views/mobile/exchangerates/ListPage.vue';
import ExchangeRatesUpdatePage from '@/views/mobile/exchangerates/UpdatePage.vue';
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

import TemplateListPage from '@/views/mobile/templates/ListPage.vue';

function asyncResolve(component: unknown): (ctx: Router.RouteCallbackCtx) => void {
    return function({ resolve }: { resolve: ({ component }: { component: unknown }) => void }): void {
        return resolve({
            component: component
        });
    } as unknown as (ctx: Router.RouteCallbackCtx) => void;
}

function checkLogin({ router, resolve, reject }: { router: Router.Router, resolve: () => void, reject: () => void }): void {
    if (!isUserLogined()) {
        reject();
        router.navigate('/login', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    if (!isUserUnlocked()) {
        reject();
        router.navigate('/unlock', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    resolve();
}

function checkLocked({ router, resolve, reject }: { router: Router.Router, resolve: () => void, reject: () => void }): void {
    if (!isUserLogined()) {
        reject();
        router.navigate('/login', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    if (isUserUnlocked()) {
        reject();
        router.navigate('/', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    resolve();
}

function checkNotLogin({ router, resolve, reject }: { router: Router.Router, resolve: () => void, reject: () => void }): void {
    if (isUserLogined() && !isUserUnlocked()) {
        reject();
        router.navigate('/unlock', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    if (isUserLogined()) {
        reject();
        router.navigate('/', {
            clearPreviousHistory: true,
            browserHistory: false
        });
        return;
    }

    resolve();
}

const routes: Router.RouteParameters[] = [
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
        path: '/transaction/filter/amount',
        async: asyncResolve(TransactionAmountFilterPage),
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
        path: '/settings/textsize',
        async: asyncResolve(TextSizeSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings/filter/account',
        async: asyncResolve(AccountFilterSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings/filter/category',
        async: asyncResolve(CategoryFilterSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings/filter/tag',
        async: asyncResolve(TransactionTagFilterSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings/page',
        async: asyncResolve(PageSettingsPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/settings/sync',
        async: asyncResolve(ApplicationCloudSyncSettingsPage),
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
        async: asyncResolve(ExchangeRatesListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/exchange_rates/update',
        async: asyncResolve(ExchangeRatesUpdatePage),
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
        path: '/template/list',
        async: asyncResolve(TemplateListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/schedule/list',
        async: asyncResolve(TemplateListPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/template/add',
        async: asyncResolve(TransactionEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '/template/edit',
        async: asyncResolve(TransactionEditPage),
        beforeEnter: [checkLogin]
    },
    {
        path: '(.*)',
        redirect: '/'
    }
];

export default routes;
