const tokenLocalStorageKey = 'lab_user_token';
const userNameLocalStorageKey = 'lab_user_name';
const userNickNameLocalStorageKey = 'lab_user_nickname';

function getToken() {
    return localStorage.getItem(tokenLocalStorageKey);
}

function getUserName() {
    return localStorage.getItem(userNameLocalStorageKey);
}

function getUserNickName() {
    return localStorage.getItem(userNickNameLocalStorageKey);
}

function isUserLogined() {
    return !!getToken();
}

function updateToken(item) {
    if (typeof(item) === 'string') {
        return localStorage.setItem(tokenLocalStorageKey, item);
    } else if (typeof(item) === 'object') {
        localStorage.setItem(tokenLocalStorageKey, item.token);
        localStorage.setItem(userNameLocalStorageKey, item.username);
        localStorage.setItem(userNickNameLocalStorageKey, item.nickname);

        return true;
    } else {
        return false;
    }
}

function updateUsername(value) {
    localStorage.setItem(userNameLocalStorageKey, value);
}

function updateUserNickname(value) {
    localStorage.setItem(userNickNameLocalStorageKey, value);
}

function clearToken() {
    localStorage.removeItem(tokenLocalStorageKey);
    localStorage.removeItem(userNameLocalStorageKey);
    localStorage.removeItem(userNickNameLocalStorageKey);

    return true;
}

export default {
    getToken,
    getUserName,
    getUserNickName,
    isUserLogined,
    updateToken,
    updateUsername,
    updateUserNickname,
    clearToken
};
