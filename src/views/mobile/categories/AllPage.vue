<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Transaction Categories')" :back-link="$t('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Expense" link="#"></f7-list-item>
            <f7-list-item title="Income" link="#"></f7-list-item>
            <f7-list-item title="Transfer" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-else-if="!loading">
            <f7-list-item :title="$t('Expense')" link="/category/list?type=2"></f7-list-item>
            <f7-list-item :title="$t('Income')" link="/category/list?type=1"></f7-list-item>
            <f7-list-item :title="$t('Transfer')" link="/category/list?type=3"></f7-list-item>
        </f7-list>
    </f7-page>
</template>


<script>
import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore)
    },
    created() {
        const self = this;

        self.loading = true;

        self.transactionCategoriesStore.loadAllCategories({
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            const self = this;
            const force = !!done;

            self.transactionCategoriesStore.loadAllCategories({
                force: force
            }).then(() => {
                if (done) {
                    done();
                }

                if (force) {
                    self.$toast('Category list has been updated');
                }
            }).catch(error => {
                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    }
}
</script>
