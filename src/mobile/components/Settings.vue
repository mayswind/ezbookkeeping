<template>
    <f7-page name="home">
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>
        <f7-list>
            <f7-list-button @click="logout">{{ $t('Logout') }}</f7-list-button>
        </f7-list>
    </f7-page>
</template>

<script>
import userState from '../../common/userstate.js';

export default {
    methods: {
        logout() {
            const self = this;
            const router = self.$f7router;

            userState.clearToken();
            router.navigate('/');
        }
    }
};
</script>
