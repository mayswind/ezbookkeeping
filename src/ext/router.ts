import type { NavigationGuardReturn, RouteRecordRaw } from 'vue-router';

import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';

import { canUseBusinessFeatures } from '@/ext/features.ts';

import TeamPage from '@/ext/views/TeamPage.vue';
import InventoryPage from '@/ext/views/InventoryPage.vue';
import SalesPage from '@/ext/views/SalesPage.vue';
import CustomersPage from '@/ext/views/CustomersPage.vue';
import ReportsPage from '@/ext/views/ReportsPage.vue';
import BusinessSettingsPage from '@/ext/views/BusinessSettingsPage.vue';

// Same rule as the upstream pages: be logged in and unlocked
function checkLogin(): NavigationGuardReturn {
    if (!isUserLogined()) {
        return { path: '/login', replace: true };
    }

    if (!isUserUnlocked()) {
        return { path: '/unlock', replace: true };
    }

    return true;
}

// Opt-in pages: logged in, and the business features switched on in Settings (or invited to a business)
async function checkBusinessFeatures(): Promise<NavigationGuardReturn> {
    const login = checkLogin();

    if (login !== true) {
        return login;
    }

    return await canUseBusinessFeatures() ? true : { path: '/settings/business', replace: true };
}

// The team page is also where invitations are answered, so an invitation is enough to open it
async function checkTeamAccess(): Promise<NavigationGuardReturn> {
    const login = checkLogin();

    if (login !== true) {
        return login;
    }

    return await canUseBusinessFeatures(true) ? true : { path: '/settings/business', replace: true };
}

/** Routes of the ext module, spread into the desktop router (one line in src/router/desktop.ts). */
export const extRoutes: RouteRecordRaw[] = [
    { path: '/ext/team', component: TeamPage, beforeEnter: checkTeamAccess },
    { path: '/ext/sales', component: SalesPage, beforeEnter: checkBusinessFeatures },
    { path: '/ext/customers', component: CustomersPage, beforeEnter: checkBusinessFeatures },
    { path: '/ext/inventory', component: InventoryPage, beforeEnter: checkBusinessFeatures },
    { path: '/ext/reports', component: ReportsPage, beforeEnter: checkBusinessFeatures }
];

/** Settings page of the ext module, spread into the children of the settings route (one line in src/router/desktop.ts). */
export const extSettingsRoutes: RouteRecordRaw[] = [
    { path: '/settings/business', component: BusinessSettingsPage, beforeEnter: checkLogin }
];
