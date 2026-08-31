import { type NavigationGuardReturn, createRouter, createWebHashHistory } from 'vue-router';

import { TemplateType } from '@/core/template.ts';
import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';

import LoginPage from '@/views/desktop/LoginPage.vue';
import SignUpPage from '@/views/desktop/SignupPage.vue';
import VerifyEmailPage from '@/views/desktop/VerifyEmailPage.vue';
import ForgetPasswordPage from '@/views/desktop/ForgetPasswordPage.vue';
import ResetPasswordPage from '@/views/desktop/ResetPasswordPage.vue';
import OAuth2CallbackPage from '@/views/desktop/OAuth2CallbackPage.vue';
import UnlockPage from '@/views/desktop/UnlockPage.vue';

import HomePage from '@/views/desktop/HomePage.vue';

import TransactionListPage from '@/views/desktop/transactions/ListPage.vue';

import StatisticsTransactionPage from '@/views/desktop/statistics/TransactionPage.vue';

import InsightsExplorerPage from '@/views/desktop/insights/ExplorerPage.vue';

import AccountListPage from '@/views/desktop/accounts/ListPage.vue';

import TransactionCategoryListPage from '@/views/desktop/categories/ListPage.vue';

import TransactionTagListPage from '@/views/desktop/tags/ListPage.vue';

import TransactionTemplateListPage from '@/views/desktop/templates/ListPage.vue';

import UserCustomIconListPage from '@/views/desktop/customicons/ListPage.vue';

import OverviewLayoutEditorPage from '@/views/desktop/overview/LayoutEditorPage.vue';

import SettingsPageLayout from '@/views/desktop/settings/SettingsPageLayout.vue';
import UserBasicSettingPage from '@/views/desktop/users/UserBasicSettingPage.vue';
import UserSecuritySettingPage from '@/views/desktop/users/UserSecuritySettingPage.vue';
import TwoFactorAuthPage from '@/views/desktop/users/TwoFactorAuthPage.vue';
import DataManagementPage from '@/views/desktop/users/DataManagementPage.vue';
import PreferencesSettingsPage from '@/views/desktop/settings/PreferencesSettingsPage.vue';
import ApplicationLockPage from '@/views/desktop/settings/ApplicationLockPage.vue';
import StatisticsSettingPage from '@/views/desktop/settings/StatisticsSettingPage.vue';
import ApplicationCloudSyncSettingsPage from '@/views/desktop/settings/ApplicationCloudSyncSettingsPage.vue';
import BrowserCacheSettingPage from '@/views/desktop/settings/BrowserCacheSettingPage.vue';

import ExchangeRatesListPage from '@/views/desktop/exchangerates/ListPage.vue';

function checkLogin(): NavigationGuardReturn {
    if (!isUserLogined()) {
        return {
            path: '/login',
            replace: true
        };
    }

    if (!isUserUnlocked()) {
        return {
            path: '/unlock',
            replace: true
        };
    }

    return true;
}

function checkLocked(): NavigationGuardReturn {
    if (!isUserLogined()) {
        return {
            path: '/login',
            replace: true
        };
    }

    if (isUserUnlocked()) {
        return {
            path: '/',
            replace: true
        };
    }

    return true;
}

function checkNotLogin(): NavigationGuardReturn {
    if (isUserLogined() && !isUserUnlocked()) {
        return {
            path: '/unlock',
            replace: true
        };
    }

    if (isUserLogined()) {
        return {
            path: '/',
            replace: true
        };
    }

    return true;
}

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            component: HomePage,
            beforeEnter: checkLogin
        },
        {
            path: '/transaction/list',
            component: TransactionListPage,
            beforeEnter: checkLogin,
            props: route => ({
                initPageType: route.query['pageType'],
                initDateType: route.query['dateType'],
                initMaxTime: route.query['maxTime'],
                initMinTime: route.query['minTime'],
                initType: route.query['type'],
                initCategoryIds: route.query['categoryIds'],
                initAccountIds: route.query['accountIds'],
                initTagFilter: route.query['tagFilter'],
                initAmountFilter: route.query['amountFilter'],
                initKeyword: route.query['keyword'],
                initMatchMode: route.query['matchMode']
            })
        },
        {
            path: '/statistics/transaction',
            component: StatisticsTransactionPage,
            beforeEnter: checkLogin,
            props: route => ({
                initAnalysisType: route.query['analysisType'],
                initChartDataType: route.query['chartDataType'],
                initChartType: route.query['chartType'],
                initChartDateType: route.query['chartDateType'],
                initStartTime: route.query['startTime'],
                initEndTime: route.query['endTime'],
                initFilterAccountIds: route.query['filterAccountIds'],
                initFilterCategoryIds: route.query['filterCategoryIds'],
                initTagFilter: route.query['tagFilter'],
                initKeyword: route.query['keyword'],
                initMatchMode: route.query['matchMode'],
                initSortingType: route.query['sortingType'],
                initTrendDateAggregationType: route.query['trendDateAggregationType'],
                initAssetTrendsDateAggregationType: route.query['assetTrendsDateAggregationType']
            })
        },
        {
            path: '/insights/explorer',
            component: InsightsExplorerPage,
            beforeEnter: checkLogin,
            props: route => ({
                initId: route.query['id'],
                initActiveTab: route.query['activeTab'],
                initDateRangeType: route.query['dateRangeType'],
                initStartTime: route.query['startTime'],
                initEndTime: route.query['endTime']
            })
        },
        {
            path: '/account/list',
            component: AccountListPage,
            beforeEnter: checkLogin
        },
        {
            path: '/',
            component: SettingsPageLayout,
            beforeEnter: checkLogin,
            children: [
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
                    path: '/template/list',
                    component: TransactionTemplateListPage,
                    beforeEnter: checkLogin,
                    props: {
                        initType: TemplateType.Normal.type
                    }
                },
                {
                    path: '/schedule/list',
                    component: TransactionTemplateListPage,
                    beforeEnter: checkLogin,
                    props: {
                        initType: TemplateType.Schedule.type
                    }
                },
                {
                    path: '/custom_icon/list',
                    component: UserCustomIconListPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/user/basic',
                    component: UserBasicSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/user/security',
                    component: UserSecuritySettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/user/two_factor',
                    component: TwoFactorAuthPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/user/data_management',
                    component: DataManagementPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/preferences',
                    component: PreferencesSettingsPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/application_lock',
                    component: ApplicationLockPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/statistics',
                    component: StatisticsSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/cloud_sync',
                    component: ApplicationCloudSyncSettingsPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/settings/browser_cache',
                    component: BrowserCacheSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/exchange_rate',
                    component: ExchangeRatesListPage,
                    beforeEnter: checkLogin
                }
            ]
        },
        {
            path: '/overview/edit',
            component: OverviewLayoutEditorPage,
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
            path: '/verify_email',
            component: VerifyEmailPage,
            props: route => ({
                email: route.query['email'],
                token: route.query['token'],
                hasValidEmailVerifyToken: route.query['emailSent'] === 'true'
            })
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
                token: route.query['token']
            })
        },
        {
            path: '/oauth2_callback',
            component: OAuth2CallbackPage,
            props: route => ({
                token: route.query['token'],
                provider: route.query['provider'],
                platform: route.query['platform'],
                userName: route.query['userName'],
                errorCode: route.query['errorCode'],
                errorMessage: route.query['errorMessage']
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
