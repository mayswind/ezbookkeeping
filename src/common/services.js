import axios from 'axios';
import userState from "./userstate.js";

axios.defaults.baseURL = '/api';
axios.interceptors.request.use(config => {
    const token = userState.getToken();

    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
}, error => {
    return Promise.reject(error);
});

axios.interceptors.response.use(response => {
    return response;
}, error => {
    return Promise.reject(error);
});

export default {
    authorize: (params) => {
        return axios.post('authorize.json', params);
    }
};
