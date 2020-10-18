<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>{{ $t('global.app.title') }}</f7-login-screen-title>
        <f7-list form>
            <f7-list-input
                type="text"
                clear-button
                validate
                required
                :label="$t('Username')"
                :placeholder="$t('Username or Email')"
                :value="username"
                @input="username = $event.target.value"
            ></f7-list-input>
            <f7-list-input
                type="password"
                clear-button
                validate
                required
                :label="$t('Password')"
                :placeholder="$t('Password')"
                :value="password"
                @input="password = $event.target.value"
            ></f7-list-input>
        </f7-list>
        <f7-list>
            <f7-list-button @click="login">{{ $t('Log In') }}</f7-list-button>
            <f7-block-footer>
                <span v-t="'Don\'t have an account?'"></span>&nbsp;
                <f7-link href="/signup">{{ $t('Sign up') }}</f7-link>
                <br/>
                <span v-t="'Have problem with login?'"></span>&nbsp;
                <f7-link href="/forget-pwd">{{ $t('Forget Password') }}</f7-link>
            </f7-block-footer>
        </f7-list>
    </f7-page>
</template>

<script>
import services from '../../common/services.js'
import userState from '../../common/userstate.js'

export default {
    data() {
        return {
            username: '',
            password: ''
        };
    },
    methods: {
        login() {
            const self = this
            const app = self.$f7
            const router = self.$f7router

            services.authorize({
                loginName: self.username,
                password: self.password
            }).then(response => {
                const data = response.data

                if (data && data.success && data.result && typeof(data.result) === 'string') {
                    userState.updateToken(data.result)
                    router.navigate('/')
                } else {
                    app.dialog.alert(self.$i18n.t('Unable to login'))
                }
            }).catch(error => {
                if (error.response && error.response.data && error.response.data.errorMessage) {
                    app.dialog.alert(self.$i18n.t(`error.${error.response.data.errorMessage}`))
                } else {
                    app.dialog.alert(self.$i18n.t('Unable to login'))
                }
            })
        }
    }
};
</script>
