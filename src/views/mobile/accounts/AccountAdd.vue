<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Add Account')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t('Add')" @click="add"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-input
                        type="select"
                        :label="$t('Account Category')"
                        :value="account.category"
                        @input="account.category = $event.target.value"
                    >
                        <option v-for="accountCategory in allAccountCategories"
                                :key="accountCategory.id"
                                :value="accountCategory.id">{{ $t(accountCategory.name) }}</option>
                    </f7-list-input>

                    <f7-list-input
                        type="select"
                        disabled
                        :label="$t('Account Type')"
                        :value="account.type"
                        @input="account.type = $event.target.value"
                    >
                        <option value="1">{{ $t('Single Account') }}</option>
                        <option value="2">{{ $t('Multi Sub Accounts') }}</option>
                    </f7-list-input>

                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Account Name')"
                        :placeholder="$t('Your account name')"
                        :value="account.name"
                        @input="account.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-input
                        type="select"
                        :label="$t('Currency')"
                        :value="account.currency"
                        @input="account.currency = $event.target.value"
                    >
                        <option v-for="currency in allCurrencies"
                                :key="currency.code"
                                :value="currency.code">{{ currency.displayName }}</option>
                    </f7-list-input>

                    <f7-list-input
                        type="textarea"
                        :label="$t('Description')"
                        :placeholder="$t('Your account description (optional)')"
                        :value="account.comment"
                        @input="account.comment = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item class="lab-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            account: {
                category: 1,
                type: 1,
                name: '',
                icon: "1",
                currency: self.$user.getUserInfo() ? self.$user.getUserInfo().defaultCurrency : self.$t('default.currency'),
                comment: ''
            },
            submitting: false
        };
    },
    computed: {
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        allCurrencies() {
            return this.$getAllCurrencies();
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.account.category) {
                return 'Account category cannot be empty';
            } else if (!this.account.type) {
                return 'Account type cannot be empty';
            } else if (!this.account.name) {
                return 'Account name cannot be empty';
            } else if (!this.account.currency) {
                return 'Account currency cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            return null;
        }
    },
    methods: {
        add() {
            const self = this;
            const router = self.$f7router;

            let problemMessage = self.inputEmptyProblemMessage || self.inputInvalidProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            self.$services.addAccount({
                category: parseInt(self.account.category),
                type: parseInt(self.account.type),
                name: self.account.name,
                icon: self.account.icon,
                currency: self.account.currency,
                comment: self.account.comment
            }).then(response => {
                self.submitting = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$alert('Unable to add account');
                    return;
                }

                self.$toast('You have added a new account');
                router.back('/account/list', { force: true });
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else if (!error.processed) {
                    self.$alert('Unable to add account');
                }
            });
        }
    }
}
</script>
