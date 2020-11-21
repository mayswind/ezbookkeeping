<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left back-link-force :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Application Lock')"></f7-nav-title>
        </f7-navbar>

        <f7-card v-if="isEnableApplicationLock">
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Status')" :after="$t('Enabled')"></f7-list-item>
                    <f7-list-button @click="disable">{{ $t('Disable') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!isEnableApplicationLock">
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Status')" :after="$t('Disabled')"></f7-list-item>
                    <f7-list-button @click="enable(null)">{{ $t('Enable') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-sheet
            style="height:auto;"
            :opened="showInputPinCodeSheetForEnable" @sheet:closed="showInputPinCodeSheetForEnable = false; currentPinCodeForEnable = ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('PIN Code') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('Please input a new PIN code. PIN code would encrypt your local data, so you need input this PIN code when you launch this app. If this PIN code is lost, you should re-login.') }}</p>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-item class="list-item-pincode-input">
                            <PincodeInput secure :length="6" v-model="currentPinCodeForEnable" />
                        </f7-list-item>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPinCodeValid }" :text="$t('Continue')" @click="enable(currentPinCodeForEnable)"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="showInputPinCodeSheetForEnable = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            isEnableApplicationLock: this.$settings.isEnableApplicationLock(),
            currentPinCodeForEnable: '',
            showInputPinCodeSheetForEnable: false
        };
    },
    computed: {
        currentPinCodeValid() {
            return this.currentPinCodeForEnable && this.currentPinCodeForEnable.length === 6;
        }
    },
    methods: {
        enable(pinCode) {
            if (!pinCode) {
                this.showInputPinCodeSheetForEnable = true;
                return;
            }

            if (!this.currentPinCodeValid) {
                this.$alert('PIN code is invalid');
                return;
            }

            this.$user.encryptToken(pinCode);
            this.$settings.setEnableApplicationLock(true);
            this.isEnableApplicationLock = true;

            this.showInputPinCodeSheetForEnable = false;
        },
        disable() {
            this.$user.decryptToken();
            this.$settings.setEnableApplicationLock(false);
            this.isEnableApplicationLock = false;
        }
    }
}
</script>
