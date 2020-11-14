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
                        @input="chooseSuitableIcon(account.category, $event.target.value); account.category = $event.target.value"
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
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-if="account.type === '1'">
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Account Name')"
                        :placeholder="$t('Your account name')"
                        :value="account.name"
                        @input="account.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Account Icon')" link="#"
                                  @click="showIconSelection = true">
                        <f7-icon slot="after" :f7="account.icon | accountIcon"></f7-icon>
                    </f7-list-item>

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


        <f7-sheet :opened="showIconSelection" @sheet:closed="showIconSelection = false">
            <f7-toolbar>
                <div class="left"></div>
                <div class="right">
                    <f7-link sheet-close :text="$t('Done')"></f7-link>
                </div>
            </f7-toolbar>
            <f7-page-content>
                <f7-block>
                    <f7-row class="padding-vertical-half padding-horizontal-half" v-for="(row, idx) in allAccountIconRows" :key="idx">
                        <f7-col v-for="accountIcon in row" :key="accountIcon.id">
                            <f7-icon :f7="accountIcon.f7Icon" @click.native="account.icon = accountIcon.id; showIconSelection = false">
                                <f7-badge color="default" class="right-bottom-icon" v-if="account.icon === accountIcon.id">
                                    <f7-icon f7="checkmark_alt"></f7-icon>
                                </f7-badge>
                            </f7-icon>
                        </f7-col>
                    </f7-row>
                </f7-block>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            account: {
                category: '1',
                type: '1',
                name: '',
                icon: "1",
                currency: self.$user.getUserInfo() ? self.$user.getUserInfo().defaultCurrency : self.$t('default.currency'),
                comment: ''
            },
            submitting: false,
            showIconSelection: false
        };
    },
    computed: {
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        allAccountIconRows() {
            const allAccountIcons = this.$constants.icons.allAccountIcons;
            const iconPerRow = 7;
            const ret = [];
            let rowCount = 0;

            for (let accountIconId in allAccountIcons) {
                if (!Object.prototype.hasOwnProperty.call(allAccountIcons, accountIconId)) {
                    continue;
                }

                const accountIcon = allAccountIcons[accountIconId];

                if (!ret[rowCount]) {
                    ret[rowCount] = [];
                } else if (ret[rowCount] && ret[rowCount].length >= iconPerRow) {
                    rowCount++;
                    ret[rowCount] = [];
                }

                ret[rowCount].push({
                    id: accountIconId,
                    f7Icon: accountIcon.f7Icon
                });
            }

            return ret;
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
        },
        chooseSuitableIcon(oldCategory, newCategory) {
            const allCategories = this.$constants.account.allCategories;

            for (let i = 0; i < allCategories.length; i++) {
                if (allCategories[i].id.toString() === oldCategory) {
                    if (this.account.icon !== allCategories[i].defaultAccountIconId) {
                        return;
                    } else {
                        break;
                    }
                }
            }

            for (let i = 0; i < allCategories.length; i++) {
                if (allCategories[i].id.toString() === newCategory) {
                    this.account.icon = allCategories[i].defaultAccountIconId;
                }
            }
        }
    }
}
</script>
