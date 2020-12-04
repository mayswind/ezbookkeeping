<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': isInputEmpty() || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-input label="Account Category" placeholder="Category"></f7-list-input>
                    <f7-list-input label="Account Type" placeholder="Account Type"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
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
                        :class="{ 'disabled': editAccountId }"
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

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-input label="Account Name" placeholder="Your account name"></f7-list-input>
                    <f7-list-item header="Account Icon" after="Icon"></f7-list-item>
                    <f7-list-item header="Account Color" after="Color"></f7-list-item>
                    <f7-list-input label="Currency" placeholder="Currency"></f7-list-input>
                    <f7-list-input type="textarea" label="Description" placeholder="Your account description (optional)"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading && account.type === $constants.account.allAccountTypes.SingleAccount.toString()">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Account Name')"
                        :placeholder="$t('Your account name')"
                        :value="account.name"
                        @input="account.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Account Icon')" key="singleTypeAccountIconSelection" link="#"
                                  @click="showIconSelectionSheet(account)">
                        <f7-icon slot="after" :icon="account.icon | accountIcon" :style="{ color: '#' + account.color }"></f7-icon>
                    </f7-list-item>

                    <f7-list-item :header="$t('Account Color')" key="singleTypeAccountColorSelection" link="#"
                                  @click="showColorSelectionSheet(account)">
                        <f7-icon slot="after" f7="app_fill" :style="{ color: '#' + account.color }"></f7-icon>
                    </f7-list-item>

                    <f7-list-input
                        type="select"
                        :class="{ 'disabled': editAccountId }"
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

                    <f7-list-item :header="$t('Visible')" v-if="editAccountId">
                        <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading && account.type === $constants.account.allAccountTypes.MultiSubAccounts.toString()">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Account Name')"
                        :placeholder="$t('Your account name')"
                        :value="account.name"
                        @input="account.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Account Icon')" key="multiTypeAccountIconSelection" link="#"
                                  @click="showIconSelectionSheet(account)">
                        <f7-icon slot="after" :icon="account.icon | accountIcon" :style="{ color: '#' + account.color }"></f7-icon>
                    </f7-list-item>

                    <f7-list-item :header="$t('Account Color')" key="multiTypeAccountColorSelection" link="#"
                                  @click="showColorSelectionSheet(account)">
                        <f7-icon slot="after" f7="app_fill" :style="{ color: '#' + account.color }"></f7-icon>
                    </f7-list-item>

                    <f7-list-input
                        type="textarea"
                        :label="$t('Description')"
                        :placeholder="$t('Your account description (optional)')"
                        :value="account.comment"
                        @input="account.comment = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Visible')" v-if="editAccountId">
                        <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
            <f7-card-footer v-if="!editAccountId">
                <f7-button large fill
                           :text="$t('Add Sub Account')" @click="addSubAccount"></f7-button>
            </f7-card-footer>
        </f7-card>

        <f7-block class="no-padding no-margin" v-if="!loading && account.type === $constants.account.allAccountTypes.MultiSubAccounts.toString()">
            <f7-card v-for="(subAccount, idx) in subAccounts" :key="idx">
                <f7-card-content class="no-safe-areas" :padding="false">
                    <f7-list>
                        <f7-list-input
                            type="text"
                            clear-button
                            :label="$t('Sub Account Name')"
                            :placeholder="$t('Your sub account name')"
                            :value="subAccount.name"
                            @input="subAccount.name = $event.target.value"
                        ></f7-list-input>

                        <f7-list-item :header="$t('Sub Account Icon')" key="subAccountIconSelection" link="#"
                                      @click="showIconSelectionSheet(subAccount)">
                            <f7-icon slot="after" :icon="subAccount.icon | accountIcon" :style="{ color: '#' + subAccount.color }"></f7-icon>
                        </f7-list-item>

                        <f7-list-item :header="$t('Sub Account Color')" key="subAccountColorSelection" link="#"
                                      @click="showColorSelectionSheet(subAccount)">
                            <f7-icon slot="after" f7="app_fill" :style="{ color: '#' + subAccount.color }"></f7-icon>
                        </f7-list-item>

                        <f7-list-input
                            type="select"
                            :class="{ 'disabled': editAccountId }"
                            :label="$t('Currency')"
                            :value="subAccount.currency"
                            @input="subAccount.currency = $event.target.value"
                        >
                            <option v-for="currency in allCurrencies"
                                    :key="currency.code"
                                    :value="currency.code">{{ currency.displayName }}</option>
                        </f7-list-input>

                        <f7-list-input
                            type="textarea"
                            :label="$t('Description')"
                            :placeholder="$t('Your sub account description (optional)')"
                            :value="subAccount.comment"
                            @input="subAccount.comment = $event.target.value"
                        ></f7-list-input>

                        <f7-list-item :header="$t('Visible')" v-if="editAccountId">
                            <f7-toggle :checked="subAccount.visible" @toggle:change="subAccount.visible = $event"></f7-toggle>
                        </f7-list-item>
                    </f7-list>
                </f7-card-content>
                <f7-card-footer v-if="!editAccountId">
                    <f7-button large fill
                               color="red" :text="$t('Remove Sub Account')" @click="removeSubAccount(subAccount)"></f7-button>
                </f7-card-footer>
            </f7-card>
        </f7-block>

        <f7-sheet :opened="showIconSelection" @sheet:closed="hideIconSelectionSheet">
            <f7-toolbar>
                <div class="left"></div>
                <div class="right">
                    <f7-link sheet-close :text="$t('Done')"></f7-link>
                </div>
            </f7-toolbar>
            <f7-page-content>
                <f7-block class="margin-vertical">
                    <f7-row class="padding-vertical-half padding-horizontal-half" v-for="(row, idx) in allAccountIconRows" :key="idx">
                        <f7-col class="text-align-center" v-for="accountIcon in row" :key="accountIcon.id">
                            <f7-icon :icon="accountIcon.icon" :style="{ color: '#' + (accountChoosingIcon ? accountChoosingIcon.color : '000000') }" @click.native="setSelectedIcon(accountIcon)">
                                <f7-badge color="default" class="right-bottom-icon" v-if="accountChoosingIcon && accountChoosingIcon.icon === accountIcon.id">
                                    <f7-icon f7="checkmark_alt"></f7-icon>
                                </f7-badge>
                            </f7-icon>
                        </f7-col>
                        <f7-col v-for="idx in (iconCountPerRow - row.length)" :key="idx"></f7-col>
                    </f7-row>
                </f7-block>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet :opened="showColorSelection" @sheet:closed="hideColorSelectionSheet">
            <f7-toolbar>
                <div class="left"></div>
                <div class="right">
                    <f7-link sheet-close :text="$t('Done')"></f7-link>
                </div>
            </f7-toolbar>
            <f7-page-content>
                <f7-block class="margin-vertical">
                    <f7-row class="padding-vertical padding-horizontal-half" v-for="(row, idx) in allAccountColorRows" :key="idx">
                        <f7-col class="text-align-center" v-for="accountColor in row" :key="accountColor.color">
                            <f7-icon f7="app_fill" :style="{ color: '#' + accountColor.color }" @click.native="setSelectedColor(accountColor.color)">
                                <f7-badge color="default" class="right-bottom-icon" v-if="accountChoosingColor && accountChoosingColor.color === accountColor.color">
                                    <f7-icon f7="checkmark_alt"></f7-icon>
                                </f7-badge>
                            </f7-icon>
                        </f7-col>
                        <f7-col v-for="idx in (iconCountPerRow - row.length)" :key="idx"></f7-col>
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
            editAccountId: null,
            loading: false,
            account: {
                category: '1',
                type: self.$constants.account.allAccountTypes.SingleAccount.toString(),
                name: '',
                icon: self.$constants.icons.defaultAccountIconId,
                color: self.$constants.colors.defaultAccountColor,
                currency: self.$user.getUserInfo() ? self.$user.getUserInfo().defaultCurrency : self.$t('default.currency'),
                comment: '',
                visible: true
            },
            subAccounts: [],
            iconCountPerRow: 7,
            accountChoosingIcon: null,
            accountChoosingColor: null,
            submitting: false,
            showIconSelection: false,
            showColorSelection: false
        };
    },
    computed: {
        title() {
            if (!this.editAccountId) {
                return 'Add Account';
            } else {
                return 'Edit Account';
            }
        },
        saveButtonTitle() {
            if (!this.editAccountId) {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        allAccountIconRows() {
            const allAccountIcons = this.$constants.icons.allAccountIcons;
            const ret = [];
            let rowCount = 0;

            for (let accountIconId in allAccountIcons) {
                if (!Object.prototype.hasOwnProperty.call(allAccountIcons, accountIconId)) {
                    continue;
                }

                const accountIcon = allAccountIcons[accountIconId];

                if (!ret[rowCount]) {
                    ret[rowCount] = [];
                } else if (ret[rowCount] && ret[rowCount].length >= this.iconCountPerRow) {
                    rowCount++;
                    ret[rowCount] = [];
                }

                ret[rowCount].push({
                    id: accountIconId,
                    icon: accountIcon.icon
                });
            }

            return ret;
        },
        allAccountColorRows() {
            const allAccountColors = this.$constants.colors.allAccountColors;
            const ret = [];
            let rowCount = -1;

            for (let i = 0; i < allAccountColors.length; i++) {
                if (i % this.iconCountPerRow === 0) {
                    ret[++rowCount] = [];
                }

                ret[rowCount].push({
                    color: allAccountColors[i]
                });
            }

            return ret;
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        if (query.id) {
            self.loading = true;

            self.editAccountId = query.id;
            self.$services.getAccount({
                id: self.editAccountId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get account');
                    router.back();
                    return;
                }

                const account = data.result;
                self.account.id = account.id;
                self.account.category = account.category.toString();
                self.account.type = account.type.toString();
                self.account.name = account.name;
                self.account.icon = account.icon;
                self.account.color = account.color;
                self.account.currency = account.currency;
                self.account.comment = account.comment;
                self.account.visible = !account.hidden;

                if (account.subAccounts && account.subAccounts.length > 0) {
                    for (let i = 0; i < account.subAccounts.length; i++) {
                        const subAccount = account.subAccounts[i];

                        self.subAccounts.push({
                            id: subAccount.id,
                            category: subAccount.category.toString(),
                            type: subAccount.type.toString(),
                            name: subAccount.name,
                            icon: subAccount.icon,
                            color: subAccount.color,
                            currency: subAccount.currency,
                            comment: subAccount.comment,
                            visible: !subAccount.hidden
                        });
                    }
                }

                self.loading = false;
            }).catch(error => {
                self.$logger.error('failed to load account info', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                    router.back();
                } else if (!error.processed) {
                    self.$toast('Unable to get account');
                    router.back();
                }
            });
        } else {
            self.loading = false;
        }
    },
    methods: {
        addSubAccount() {
            const self = this;

            if (self.account.type !== self.$constants.account.allAccountTypes.MultiSubAccounts.toString()) {
                return;
            }

            this.subAccounts.push({
                category: null,
                type: null,
                name: '',
                icon: this.account.icon,
                color: this.account.color,
                currency: self.$user.getUserInfo() ? self.$user.getUserInfo().defaultCurrency : self.$t('default.currency'),
                comment: ''
            });
        },
        removeSubAccount(subAccount) {
            for (let i = 0; i < this.subAccounts.length; i++) {
                if (this.subAccounts[i] === subAccount) {
                    this.subAccounts.splice(i, 1);
                }
            }
        },
        showIconSelectionSheet(account) {
            this.accountChoosingIcon = account;
            this.showIconSelection = true;
        },
        setSelectedIcon(accountIcon) {
            if (!this.accountChoosingIcon) {
                return;
            }

            this.accountChoosingIcon.icon = accountIcon.id;
            this.accountChoosingIcon = null;
            this.showIconSelection = false;
        },
        hideIconSelectionSheet() {
            this.accountChoosingIcon = null;
            this.showIconSelection = false;
        },
        showColorSelectionSheet(account) {
            this.accountChoosingColor = account;
            this.showColorSelection = true;
        },
        setSelectedColor(color) {
            if (!this.accountChoosingColor) {
                return;
            }

            this.accountChoosingColor.color = color;
            this.accountChoosingColor = null;
            this.showColorSelection = false;
        },
        hideColorSelectionSheet() {
            this.accountChoosingColor = null;
            this.showColorSelection = false;
        },
        save() {
            const self = this;
            const router = self.$f7router;

            let problemMessage = self.getInputEmptyProblemMessage(self.account, false);

            if (!problemMessage && self.account.type === self.$constants.account.allAccountTypes.MultiSubAccounts.toString()) {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    problemMessage = self.getInputEmptyProblemMessage(self.subAccounts[i], true);

                    if (problemMessage) {
                        break;
                    }
                }
            }

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            const subAccounts = [];

            if (self.account.type === self.$constants.account.allAccountTypes.MultiSubAccounts.toString()) {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    const subAccount = self.subAccounts[i];
                    const submitAccount = {
                        category: parseInt(self.account.category),
                        type: self.$constants.account.allAccountTypes.SingleAccount,
                        name: subAccount.name,
                        icon: subAccount.icon,
                        color: subAccount.color,
                        currency: subAccount.currency,
                        comment: subAccount.comment
                    };

                    if (self.editAccountId) {
                        submitAccount.id = subAccount.id;
                        submitAccount.hidden = !subAccount.visible;
                    }

                    subAccounts.push(submitAccount);
                }
            }

            const submitAccount = {
                category: parseInt(self.account.category),
                type: parseInt(self.account.type),
                name: self.account.name,
                icon: self.account.icon,
                color: self.account.color,
                currency: self.account.type === self.$constants.account.allAccountTypes.SingleAccount.toString() ? self.account.currency : self.$constants.currency.parentAccountCurrencyPlaceholder,
                comment: self.account.comment,
                subAccounts: self.account.type === self.$constants.account.allAccountTypes.SingleAccount.toString() ? null : subAccounts,
            };

            let promise = null;

            if (!self.editAccountId) {
                promise = self.$services.addAccount(submitAccount);
            } else {
                submitAccount.id = self.account.id;
                submitAccount.hidden = !self.account.visible;
                promise = self.$services.modifyAccount(submitAccount);
            }

            promise.then(response => {
                self.submitting = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!self.editAccountId) {
                        self.$toast('Unable to add account');
                    } else {
                        self.$toast('Unable to save account');
                    }
                    return;
                }

                if (!self.editAccountId) {
                    self.$toast('You have added a new account');
                } else {
                    self.$toast('You have saved this account');
                }

                router.back();
            }).catch(error => {
                self.$logger.error('failed to save account', error);

                self.submitting = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    if (!self.editAccountId) {
                        self.$toast('Unable to add account');
                    } else {
                        self.$toast('Unable to save account');
                    }
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
        },
        isInputEmpty() {
            const isAccountEmpty = !!this.getInputEmptyProblemMessage(this.account, false);

            if (isAccountEmpty) {
                return true;
            }

            if (this.account.type === this.$constants.account.allAccountTypes.MultiSubAccounts.toString()) {
                for (let i = 0; i < this.subAccounts.length; i++) {
                    const isSubAccountEmpty = !!this.getInputEmptyProblemMessage(this.subAccounts[i], true);

                    if (isSubAccountEmpty) {
                        return true;
                    }
                }
            }

            return false;
        },
        getInputEmptyProblemMessage(account, isSubAccount) {
            if (!isSubAccount && !account.category) {
                return 'Account category cannot be empty';
            } else if (!isSubAccount && !account.type) {
                return 'Account type cannot be empty';
            } else if (!account.name) {
                return 'Account name cannot be empty';
            } else if (account.type === this.$constants.account.allAccountTypes.SingleAccount.toString() && !account.currency) {
                return 'Account currency cannot be empty';
            } else {
                return null;
            }
        }
    }
}
</script>
