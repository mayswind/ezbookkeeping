<template>
    <div>
        <v-tabs show-arrows class="text-uppercase" v-model="activeTab">
            <v-tab value="basicSetting">
                <v-icon size="20" start :icon="icons.basicSetting"/>
                {{ $t('Basic') }}
            </v-tab>
            <v-tab value="securitySetting">
                <v-icon size="20" start :icon="icons.securitySetting"/>
                {{ $t('Security') }}
            </v-tab>
            <v-tab value="twoFactorSetting">
                <v-icon size="20" start :icon="icons.twoFactorSetting"/>
                {{ $t('Two-Factor Authentication') }}
            </v-tab>
            <v-tab value="dataManagementSetting">
                <v-icon size="20" start :icon="icons.dataManagementSetting"/>
                {{ $t('Data Management') }}
            </v-tab>
        </v-tabs>
        <v-divider />
        <v-window class="mt-5 disable-tab-transition" v-model="activeTab">
            <v-window-item value="basicSetting">
                <user-basic-setting-tab/>
            </v-window-item>

            <v-window-item value="securitySetting">
                <user-security-setting-tab/>
            </v-window-item>

            <v-window-item value="twoFactorSetting">
                <user-two-factor-auth-setting-tab/>
            </v-window-item>

            <v-window-item value="dataManagementSetting">
                <user-data-management-setting-tab/>
            </v-window-item>
        </v-window>
    </div>
</template>

<script>
import UserBasicSettingTab from "./settings/UserBasicSettingTab.vue";
import UserSecuritySettingTab from "./settings/UserSecuritySettingTab.vue";
import UserTwoFactorAuthSettingTab from "./settings/UserTwoFactorAuthSettingTab.vue";
import UserDataManagementSettingTab from "./settings/UserDataManagementSettingTab.vue";

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
    data() {
        return {
            activeTab: 'basicSetting',
            icons: {
                basicSetting: mdiAccountOutline,
                securitySetting: mdiLockOpenOutline,
                twoFactorSetting: mdiOnepassword,
                dataManagementSetting: mdiDatabaseCogOutline
            }
        };
    }
}
</script>
