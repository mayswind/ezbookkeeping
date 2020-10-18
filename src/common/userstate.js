const tokenLocalStorageKey = 'UserToken';

function getToken() {
    return localStorage.getItem(tokenLocalStorageKey);
}

function isUserLogined() {
    return !!getToken();
}

function updateToken(token) {
    return localStorage.setItem(tokenLocalStorageKey, token);
}

function clearToken() {
    return localStorage.removeItem(tokenLocalStorageKey);
}

export default {
    getToken,
    isUserLogined,
    updateToken,
    clearToken
};
