<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar :title="$t('Transaction Categories')" :back-link="$t('Back')"></f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item title="Expense" link="#"></f7-list-item>
                    <f7-list-item title="Income" link="#"></f7-list-item>
                    <f7-list-item title="Transfer" link="#"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
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
    data() {
        return {
            loading: true
        };
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        self.$store.dispatch('loadAllCategories', {
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            self.logining = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
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
