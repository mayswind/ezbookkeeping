import utils from './utils.js'

const tokenLocalStorageKey = 'lab_user_token';
const userInfoLocalStorageKey = 'lab_user_info';

function getToken() {
    return localStorage.getItem(tokenLocalStorageKey);
}

function getUserInfo() {
    const data = localStorage.getItem(userInfoLocalStorageKey);
    return JSON.parse(data);
}

function isUserLogined() {
    return !!getToken();
}

function updateToken(token) {
    if (utils.isString(token)) {
        localStorage.setItem(tokenLocalStorageKey, token);
    }
}

function updateUserInfo(user) {
    if (utils.isObject(user)) {
        localStorage.setItem(userInfoLocalStorageKey, JSON.stringify(user));
    }
}

function updateTokenAndUserInfo(item) {
    if (utils.isObject(item)) {
        if (item.token) {
            updateToken(item.token);
        }

        if (item.user) {
            updateUserInfo(item.user);
        }
    }
}

function clearTokenAndUserInfo() {
    localStorage.removeItem(tokenLocalStorageKey);
    localStorage.removeItem(userInfoLocalStorageKey);
}

export default {
    getToken,
    getUserInfo,
    isUserLogined,
    updateToken,
    updateUserInfo,
    updateTokenAndUserInfo,
    clearTokenAndUserInfo
};
