<template>
    <f7-page>
        <f7-navbar :title="$t('Sign Up')" :back-link="$t('Back')"></f7-navbar>
        <f7-list no-hairlines-md>
            <f7-list-input
                type="text"
                clear-button
                :label="$t('Username')"
                :placeholder="$t('Your username')"
                :value="username"
                @input="username = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password')"
                :value="password"
                @input="password = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="email"
                clear-button
                :label="$t('E-mail')"
                :placeholder="$t('Your email')"
                :value="email"
                @input="email = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="text"
                clear-button
                :label="$t('Nickname')"
                :placeholder="$t('Your nickname')"
                :value="nickname"
                @input="nickname = $event.target.value"
            ></f7-list-input>
        </f7-list>

        <f7-button large fill :class="{ 'disabled': inputIsEmpty }" :text="$t('Sign Up')" @click="signup"></f7-button>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            username: '',
            password: '',
            email: '',
            nickname: ''
        };
    },
    computed: {
        inputIsEmpty() {
            return !this.username || !this.password || !this.email || !this.nickname;
        }
    },
    methods: {
        signup() {
            const self = this;
            const app = self.$f7;
            const router = self.$f7router;

            if (!this.username) {
                self.$alert('Please input username');
                return;
            }

            if (!this.password) {
                self.$alert('Please input password');
                return;
            }

            if (!this.email) {
                self.$alert('Please input email');
                return;
            }

            if (!this.nickname) {
                self.$alert('Please input nickname');
                return;
            }

            let hasResponse = false;

            setTimeout(() => {
                if (!hasResponse) {
                    app.preloader.show();
                }
            }, 200);

            self.$services.register({
                username: self.username,
                password: self.password,
                email: self.email,
                nickname: self.nickname
            }).then(response => {
                hasResponse = true;
                app.preloader.hide();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$alert('Unable to sign up');
                    return;
                }

                if (typeof(data.result) === 'string') {
                    self.$user.updateToken(data.result);
                }

                self.$toast('You have been successfully registered');
                router.navigate('/');
            }).catch(error => {
                hasResponse = true;
                app.preloader.hide();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else {
                    self.$alert('Unable to sign up');
                }
            });
        }
    }
};
</script>
