<template>
    <f7-page>
        <f7-navbar :title="$t('User Profile')" :back-link="$t('Back')"></f7-navbar>
        <f7-list no-hairlines-md>
            <f7-list-input
                type="password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password')"
                :value="password"
                @input="password = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="password"
                clear-button
                :label="$t('Confirmation Password')"
                :placeholder="$t('Re-enter the password')"
                :value="confirmPassword"
                @input="confirmPassword = $event.target.value"
            ></f7-list-input>

            <f7-list-input
                type="email"
                clear-button
                :label="$t('E-mail')"
                :placeholder="$t('Your email address')"
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

            <f7-list-item class="lab-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <f7-button large fill :class="{ 'disabled': inputIsNotChanged }" :text="$t('Update')" @click="update"></f7-button>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            password: '',
            confirmPassword: '',
            oldEmail: '',
            email: '',
            oldNickname: '',
            nickname: ''
        };
    },
    computed: {
        inputIsNotChanged() {
            return !!this.inputIsNotChangedProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputIsNotChangedProblemMessage() {
            if (!this.password && !this.confirmPassword && !this.email && !this.nickname) {
                return 'Nothing has been modified';
            } else if (!this.password && !this.confirmPassword && this.email === this.oldEmail && this.nickname === this.oldNickname) {
                return 'Nothing has been modified';
            } else if (!this.password && this.confirmPassword) {
                return 'Password cannot be empty';
            } else if (this.password && !this.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.password && this.confirmPassword && this.password !== this.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else {
                return null;
            }
        }
    },
    created() {
        const self = this;
        const app = self.$f7;
        const router = self.$f7router;

        app.preloader.show();

        self.$services.getProfile().then(response => {
            app.preloader.hide();
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$alert('Unable to get user profile', () => {
                    router.back();
                });
                return;
            }

            self.oldEmail = data.result.email;
            self.oldNickname = data.result.nickname;

            self.email = self.oldEmail
            self.nickname = self.oldNickname;
        }).catch(error => {
            app.preloader.hide();

            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$alert({ error: error.response.data }, () => {
                    router.back();
                });
            } else if (!error.processed) {
                self.$alert('Unable to get user profile', () => {
                    router.back();
                });
            }
        });
    },
    methods: {
        update() {
            const self = this;
            const app = self.$f7;
            const router = self.$f7router;

            let problemMessage = self.inputIsNotChangedProblemMessage || self.inputInvalidProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            let hasResponse = false;

            setTimeout(() => {
                if (!hasResponse) {
                    app.preloader.show();
                }
            }, 200);

            self.$services.updateProfile({
                password: self.password,
                email: self.email,
                nickname: self.nickname
            }).then(response => {
                hasResponse = true;
                app.preloader.hide();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$alert('Unable to update user profile');
                    return;
                }

                if (self.nickname) {
                    self.$user.updateUserNickname(self.nickname);
                }

                if (typeof(data.result) === 'string') {
                    self.$user.updateToken(data.result);
                }

                self.$toast('Your profile has been successfully updated');
                router.back('/settings', { force: true });
            }).catch(error => {
                hasResponse = true;
                app.preloader.hide();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else if (!error.processed) {
                    self.$alert('Unable to update user profile');
                }
            });
        }
    }
};
</script>
