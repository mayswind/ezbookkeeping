import userState from "../lib/userstate.js";

import MainPage from '../views/mobile/Main.vue';
import MainPageHomeTab from '../views/mobile/main/Home.vue';

import LoginPage from '../views/mobile/Login.vue';
import SignUpPage from '../views/mobile/Signup.vue';
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

function checkNotLogin(to, from, resolve, reject) {
    const router = this;

    if (!userState.isUserLogined()) {
        resolve();
        return;
    }

    reject();
    router.navigate('/');
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
        component: LoginPage,
        beforeEnter: checkNotLogin
    },
    {
        path: '/signup',
        component: SignUpPage,
        beforeEnter: checkNotLogin
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
