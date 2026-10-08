import type { NavigationGuardReturn, RouteRecordRaw } from 'vue-router';

import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';

import TeamPage from '@/ext/views/TeamPage.vue';
import InventoryPage from '@/ext/views/InventoryPage.vue';
import SalesPage from '@/ext/views/SalesPage.vue';

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

/** Routes of the ext module, spread into the desktop router (one line in src/router/desktop.ts). */
export const extRoutes: RouteRecordRaw[] = [
    {
        path: '/ext/team',
        component: TeamPage,
        beforeEnter: checkLogin
    },
    {
        path: '/ext/sales',
        component: SalesPage,
        beforeEnter: checkLogin
    },
    {
        path: '/ext/inventory',
        component: InventoryPage,
        beforeEnter: checkLogin
    }
];
