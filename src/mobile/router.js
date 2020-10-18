import userState from "../common/userstate.js";

import MainPage from './components/Main.vue';
import MainPageHomeTab from './components/main/Home.vue';

import LoginPage from './components/Login.vue';
import SettingsPage from './components/Settings.vue';

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
