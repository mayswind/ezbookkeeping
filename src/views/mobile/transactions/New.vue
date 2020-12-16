<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>

            <f7-subnavbar>
                <f7-segmented strong>
                    <f7-button tab-link="#expense" :text="$t('Expense')" tab-link-active></f7-button>
                    <f7-button tab-link="#income" :text="$t('Income')"></f7-button>
                    <f7-button tab-link="#transfer" :text="$t('Transfer')"></f7-button>
                </f7-segmented>
            </f7-subnavbar>
        </f7-navbar>

        <f7-tabs>
            <f7-tab id="expense" tab-active class="page-content no-padding-top">
                <f7-card>
                    <f7-card-content class="no-safe-areas" :padding="false">
                        <f7-list form>
                            <f7-list-item
                                class="color-theme-teal transaction-edit-amount padding-top-half padding-bottom-half"
                                :header="$t('Expense Amount')"
                                :title="transaction.destinationAmount | currency"
                                @click="transaction.showDestinationAmountSheet = true"
                            >
                                <NumberPadSheet :show.sync="transaction.showDestinationAmountSheet"
                                                v-model="transaction.destinationAmount"
                                ></NumberPadSheet>
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Category')"
                            >
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Account')"
                            >
                            </f7-list-item>

                            <f7-list-input
                                :label="$t('Transaction Time')"
                                type="datetime-local"
                                class="transaction-edit-time"
                                :value="transaction.time"
                                @input="transaction.time = $event.target.value"
                            >
                            </f7-list-input>

                            <f7-list-item
                                :header="$t('Tags')"
                            >
                            </f7-list-item>

                            <f7-list-input
                                type="textarea"
                                :label="$t('Description')"
                                :placeholder="$t('Your transaction description (optional)')"
                                :value="transaction.comment"
                                @input="transaction.comment = $event.target.value"
                            ></f7-list-input>
                        </f7-list>
                    </f7-card-content>
                </f7-card>
            </f7-tab>
            <f7-tab id="income" class="page-content no-padding-top">
                <f7-card>
                    <f7-card-content class="no-safe-areas" :padding="false">
                        <f7-list form>
                            <f7-list-item
                                class="color-theme-red transaction-edit-amount padding-top-half padding-bottom-half"
                                :header="$t('Income Amount')"
                                :title="transaction.destinationAmount | currency"
                                @click="transaction.showDestinationAmountSheet = true"
                            >
                                <NumberPadSheet :show.sync="transaction.showDestinationAmountSheet"
                                                v-model="transaction.destinationAmount"
                                ></NumberPadSheet>
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Category')"
                            >
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Account')"
                            >
                            </f7-list-item>

                            <f7-list-input
                                :label="$t('Transaction Time')"
                                type="datetime-local"
                                class="transaction-edit-time"
                                :value="transaction.time"
                                @input="transaction.time = $event.target.value"
                            >
                            </f7-list-input>

                            <f7-list-item
                                :header="$t('Tags')"
                            >
                            </f7-list-item>

                            <f7-list-input
                                type="textarea"
                                :label="$t('Description')"
                                :placeholder="$t('Your transaction description (optional)')"
                                :value="transaction.comment"
                                @input="transaction.comment = $event.target.value"
                            ></f7-list-input>
                        </f7-list>
                    </f7-card-content>
                </f7-card>
            </f7-tab>
            <f7-tab id="transfer" class="page-content no-padding-top">
                <f7-card>
                    <f7-card-content class="no-safe-areas" :padding="false">
                        <f7-list form>
                            <f7-list-item
                                class="transaction-edit-amount padding-top-half padding-bottom-half"
                                :header="$t('Transfer Out Amount')"
                                :title="transaction.sourceAmount | currency"
                                @click="transaction.showSourceAmountSheet = true"
                            >
                                <NumberPadSheet :show.sync="transaction.showSourceAmountSheet"
                                                v-model="transaction.sourceAmount"
                                ></NumberPadSheet>
                            </f7-list-item>

                            <f7-list-item
                                class="transaction-edit-amount padding-top-half padding-bottom-half"
                                :header="$t('Transfer In Amount')"
                                :title="transaction.destinationAmount | currency"
                                @click="transaction.showDestinationAmountSheet = true"
                            >
                                <NumberPadSheet :show.sync="transaction.showDestinationAmountSheet"
                                                v-model="transaction.destinationAmount"
                                ></NumberPadSheet>
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Category')"
                            >
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Source Account')"
                            >
                            </f7-list-item>

                            <f7-list-item
                                :header="$t('Destination Account')"
                            >
                            </f7-list-item>

                            <f7-list-input
                                :label="$t('Transaction Time')"
                                type="datetime-local"
                                class="transaction-edit-time"
                                :value="transaction.time"
                                @input="transaction.time = $event.target.value"
                            >
                            </f7-list-input>

                            <f7-list-item
                                :header="$t('Tags')"
                            >
                            </f7-list-item>

                            <f7-list-input
                                type="textarea"
                                :label="$t('Description')"
                                :placeholder="$t('Your transaction description (optional)')"
                                :value="transaction.comment"
                                @input="transaction.comment = $event.target.value"
                            ></f7-list-input>
                        </f7-list>
                    </f7-card-content>
                </f7-card>
            </f7-tab>
        </f7-tabs>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            editTransactionId: null,
            transaction: {
                sourceAmount: 0,
                destinationAmount: 0,
                time: self.$utilities.formatDate(new Date(), 'YYYY-MM-DDTHH:mm'),
                comment: '',
                showSourceAmountSheet: false,
                showDestinationAmountSheet: false
            },
            allAccounts: [],
            allCategories: {},
            allTags: [],
            submitting: false
        };
    },
    computed: {
        title() {
            if (!this.editTransactionId) {
                return 'Add Transaction';
            } else {
                return 'Edit Transaction';
            }
        },
        saveButtonTitle() {
            if (!this.editTransactionId) {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            return null;
        }
    },
    watch: {
        'transaction.sourceAmount': function (newValue) {
            this.transaction.destinationAmount = newValue;
        },
        'transaction.destinationAmount': function (newValue) {
            this.transaction.sourceAmount = newValue;
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        const promises = [
            self.$services.getAllAccounts(),
            self.$services.getAllTransactionCategories({}),
            self.$services.getAllTransactionTags()
        ];

        Promise.all(promises).then(function (responses) {
            const accountDta = responses[0].data;
            const categoryData = responses[1].data;
            const tagData = responses[2].data;

            if (!accountDta || !accountDta.success || !accountDta.result) {
                self.$toast('Unable to get account list');
                router.back();
                return;
            }

            if (!categoryData || !categoryData.success || !categoryData.result) {
                self.$toast('Unable to get category list');
                router.back();
                return;
            }

            if (!tagData || !tagData.success || !tagData.result) {
                self.$toast('Unable to get tag list');
                router.back();
                return;
            }

            self.allAccounts = accountDta.result;
            self.allCategories = categoryData.result;
            self.allTags = tagData.result;

            self.loading = false;
        }).catch(errors => {
            self.$logger.error('failed to load essential data for editing transaction', errors);

            for (let i = 0; i < errors.length; i++) {
                const error = errors[i];

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                    router.back();
                    return;
                } else if (!error.processed) {
                    self.$toast('An error has occurred');
                    router.back();
                    return;
                }
            }
        });
    },
    methods: {
        save() {

        }
    }
};
</script>

<style>
.transaction-edit-amount {
    font-size: 40px;
    font-weight: bolder;
    color: var(--f7-theme-color);
}

.transaction-edit-time input[type="datetime-local"] {
    max-width: inherit;
}
</style>
