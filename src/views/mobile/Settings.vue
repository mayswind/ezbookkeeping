<template>
    <f7-page name="home">
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>
        <f7-list>
            <f7-list-button @click="logout">{{ $t('Logout') }}</f7-list-button>
        </f7-list>
    </f7-page>
</template>

<script>
export default {
    methods: {
        logout() {
            const router = this.$f7router;

            this.$user.clearToken();
            router.navigate('/');
        }
    }
};
</script>
