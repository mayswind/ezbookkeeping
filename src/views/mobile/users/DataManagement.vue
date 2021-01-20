<template>
    <f7-page>
        <f7-navbar :title="$t('Data Management')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-button external no-chevron target="_blank" :link="`${$constants.api.baseUrlPath}/data/export.csv?token=${$user.getToken()}`">{{ $t('Export Data') }}</f7-list-button>
                    <f7-list-button color="red" @click="clearData(null)">{{ $t('Clear User Data') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <password-input-sheet :title="$t('Are you sure you want to clear all data?')"
                              :hint="$t('You CANNOT undo this action. This will clear your accounts, categories, tags and transactions data. Please input your current password to confirm.')"
                              :show.sync="showInputPasswordSheetForClearData"
                              :confirm-disabled="clearingData"
                              :cancel-disabled="clearingData"
                              v-model="currentPasswordForClearData"
                              @password:confirm="clearData">
        </password-input-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            currentPasswordForClearData: '',
            clearingData: false,
            showInputPasswordSheetForClearData: false,
        };
    },
    methods: {
        clearData(password) {
            const self = this;

            if (!password) {
                self.currentPasswordForClearData = '';
                self.showInputPasswordSheetForClearData = true;
                return;
            }

            self.clearingData = true;
            self.$showLoading(() => self.clearingData);

            self.$store.dispatch('clearUserData', {
                password: password
            }).then(() => {
                self.clearingData = false;
                self.$hideLoading();

                self.showInputPasswordSheetForClearData = false;
                self.$toast('All user data has been cleared');
            }).catch(error => {
                self.clearingData = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
    }
};
</script>
