<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Transaction Categories')" :back-link="$t('Back')"></f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list dividers>
                    <f7-list-item title="Expense" link="#"></f7-list-item>
                    <f7-list-item title="Income" link="#"></f7-list-item>
                    <f7-list-item title="Transfer" link="#"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list dividers>
                    <f7-list-item :title="$t('Expense')" link="/category/list?type=2"></f7-list-item>
                    <f7-list-item :title="$t('Income')" link="/category/list?type=1"></f7-list-item>
                    <f7-list-item :title="$t('Transfer')" link="/category/list?type=3"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>


<script>
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
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('loadAllCategories', {
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

            self.$store.dispatch('loadAllCategories', {
                force: true
            }).then(() => {
                if (done) {
                    done();
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
