<template>
    <div>
        <v-tabs show-arrows class="v-tabs-pill text-uppercase" v-model="activeTab">
            <v-tab value="basicSetting" @click="pushRouter('basicSetting')">
                <v-icon size="20" start :icon="icons.basicSetting"/>
                {{ $t('Basic') }}
            </v-tab>
            <v-tab value="securitySetting" @click="pushRouter('securitySetting')">
                <v-icon size="20" start :icon="icons.securitySetting"/>
                {{ $t('Security') }}
            </v-tab>
            <v-tab value="twoFactorSetting" @click="pushRouter('twoFactorSetting')">
                <v-icon size="20" start :icon="icons.twoFactorSetting"/>
                {{ $t('Two-Factor Authentication') }}
            </v-tab>
            <v-tab value="dataManagementSetting" @click="pushRouter('dataManagementSetting')">
                <v-icon size="20" start :icon="icons.dataManagementSetting"/>
                {{ $t('Data Management') }}
            </v-tab>
        </v-tabs>

        <v-window class="mt-4 disable-tab-transition" v-model="activeTab">
            <v-window-item value="basicSetting">
                <user-basic-setting-tab/>
            </v-window-item>

            <v-window-item value="securitySetting">
                <user-security-setting-tab/>
            </v-window-item>

            <v-window-item value="twoFactorSetting">
                <user-two-factor-auth-setting-tab ref="twoFactorSettingTab"/>
            </v-window-item>

            <v-window-item value="dataManagementSetting">
                <user-data-management-setting-tab/>
            </v-window-item>
        </v-window>
    </div>
</template>

<script>
import UserBasicSettingTab from './settings/tabs/UserBasicSettingTab.vue';
import UserSecuritySettingTab from './settings/tabs/UserSecuritySettingTab.vue';
import UserTwoFactorAuthSettingTab from './settings/tabs/UserTwoFactorAuthSettingTab.vue';
import UserDataManagementSettingTab from './settings/tabs/UserDataManagementSettingTab.vue';

import {
    mdiAccountOutline,
    mdiLockOpenOutline,
    mdiOnepassword,
    mdiDatabaseCogOutline
} from '@mdi/js';

export default {
    components: {
        UserBasicSettingTab,
        UserSecuritySettingTab,
        UserTwoFactorAuthSettingTab,
        UserDataManagementSettingTab
    },
    props: [
        'initTab'
    ],
    data() {
        let queryActiveTab = this.initTab || 'basicSetting';

        if ([
            'basicSetting',
            'securitySetting',
            'twoFactorSetting',
            'dataManagementSetting'
        ].indexOf(queryActiveTab) === -1) {
            queryActiveTab = 'basicSetting';
        }

        return {
            activeTab: queryActiveTab,
            icons: {
                basicSetting: mdiAccountOutline,
                securitySetting: mdiLockOpenOutline,
                twoFactorSetting: mdiOnepassword,
                dataManagementSetting: mdiDatabaseCogOutline
            }
        };
    },
    watch: {
        'activeTab': function (newValue, oldValue) {
            if (oldValue === 'twoFactorSetting' && newValue !== 'twoFactorSetting') {
                this.$refs.twoFactorSettingTab.reset();
            }
        }
    },
    beforeRouteUpdate(to) {
        if (to.query && to.query.tab && [
            'basicSetting',
            'securitySetting',
            'twoFactorSetting',
            'dataManagementSetting'
        ].indexOf(to.query.tab) >= 0) {
            this.activeTab = to.query.tab;
        } else {
            this.activeTab = 'basicSetting';
        }
    },
    methods: {
        pushRouter(tab) {
            this.$router.push(`/user/settings?tab=${tab}`);
        }
    }
}
</script>
