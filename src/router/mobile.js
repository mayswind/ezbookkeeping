import userState from "../lib/userstate.js";

import HomePage from '../views/mobile/Home.vue';
import LoginPage from '../views/mobile/Login.vue';
import SignUpPage from '../views/mobile/Signup.vue';

import SettingsPage from '../views/mobile/Settings.vue';
import AboutPage from "../views/mobile/About.vue";
import UserProfilePage from "../views/mobile/users/UserProfile.vue";
import TwoFactorAuthPage from "../views/mobile/users/TwoFactorAuth.vue";
import SessionListPage from "../views/mobile/users/SessionList.vue";

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
        path: '/settings',
        component: SettingsPage,
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
        path: '(.*)',
        redirect: '/'
    }
];

export default routes;
