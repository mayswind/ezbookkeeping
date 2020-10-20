import userState from "../lib/userstate.js";

import MainPage from '../views/mobile/Main.vue';
import MainPageHomeTab from '../views/mobile/main/Home.vue';

import LoginPage from '../views/mobile/Login.vue';
import SettingsPage from '../views/mobile/Settings.vue';

function checkLogin(to, from, resolve, reject) {
    const router = this;

    if (userState.isUserLogined()) {
        resolve();
        return;
    }

    reject();
    router.navigate('/login');
}

const routes = [
    {
        path: '/',
        component: MainPage,
        tabs: [
            {
                path: '/',
                id: 'main-tab-home',
                component: MainPageHomeTab,
                beforeEnter: checkLogin
            }
        ],
        beforeEnter: checkLogin
    },
    {
        path: '/login',
        component: LoginPage
    },
    {
        path: '/settings',
        component: SettingsPage,
        beforeEnter: checkLogin
    },
    {
        path: '(.*)',
        redirect: '/'
    }
];

export default routes;
