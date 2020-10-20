<template>
    <f7-app :params="f7params">
        <f7-view id="main-view" main url="/" :push-state="true"></f7-view>
    </f7-app>
</template>

<script>
import routes from './router/mobile.js';

export default {
    data() {
        return {
            f7params: {
                name: 'lab',
                id: 'net.mayswind.lab',
                theme: 'ios',
                routes: routes
            }
        }
    }
}
</script>

<style>
body {
    -ms-user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    user-select: none;
}
</style>
