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

import UserSettingsPageLayout from '@/views/desktop/user/UserSettingsPageLayout.vue';
import UserBasicSettingPage from '@/views/desktop/user/UserBasicSettingPage.vue';
import UserSecuritySettingPage from '@/views/desktop/user/UserSecuritySettingPage.vue';
import UserTwoFactorAuthSettingPage from '@/views/desktop/user/UserTwoFactorAuthSettingPage.vue';
import UserDataManagementSettingPage from '@/views/desktop/user/UserDataManagementSettingPage.vue';

import AppSettingsPageLayout from '@/views/desktop/app/AppSettingsPageLayout.vue';
import AppBasicSettingPage from '@/views/desktop/app/AppBasicSettingPage.vue';
import AppLockSettingPage from '@/views/desktop/app/AppLockSettingPage.vue';
import AppStatisticsSettingPage from '@/views/desktop/app/AppStatisticsSettingPage.vue';
import AppCloudSyncSettingPage from '@/views/desktop/app/AppCloudSyncSettingPage.vue';
import AppBrowserCacheSettingPage from '@/views/desktop/app/AppBrowserCacheSettingPage.vue';

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
            component: UserSettingsPageLayout,
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
                    path: '/user/settings/basic',
                    component: UserBasicSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/user/settings/security',
                    component: UserSecuritySettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/user/settings/two_factor',
                    component: UserTwoFactorAuthSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/user/settings/data_management',
                    component: UserDataManagementSettingPage,
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
            path: '/',
            component: AppSettingsPageLayout,
            beforeEnter: checkLogin,
            children: [
                {
                    path: '/app/settings/basic',
                    component: AppBasicSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/app/settings/application_lock',
                    component: AppLockSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/app/settings/statistics',
                    component: AppStatisticsSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/app/settings/cloud_sync',
                    component: AppCloudSyncSettingPage,
                    beforeEnter: checkLogin
                },
                {
                    path: '/app/settings/browser_cache',
                    component: AppBrowserCacheSettingPage,
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
