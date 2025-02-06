<template>
    <div>
        <v-tabs show-arrows v-model="activeTab">
            <v-tab value="basicSetting" @click="pushRouter('basicSetting')">
                <v-icon size="20" start :icon="mdiAccountOutline"/>
                {{ tt('Basic') }}
            </v-tab>
            <v-tab value="securitySetting" @click="pushRouter('securitySetting')">
                <v-icon size="20" start :icon="mdiLockOpenOutline"/>
                {{ tt('Security') }}
            </v-tab>
            <v-tab value="twoFactorSetting" @click="pushRouter('twoFactorSetting')">
                <v-icon size="20" start :icon="mdiOnepassword"/>
                {{ tt('Two-Factor Authentication') }}
            </v-tab>
            <v-tab value="dataManagementSetting" @click="pushRouter('dataManagementSetting')">
                <v-icon size="20" start :icon="mdiDatabaseCogOutline"/>
                {{ tt('Data Management') }}
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

<script setup lang="ts">
import UserBasicSettingTab from './settings/tabs/UserBasicSettingTab.vue';
import UserSecuritySettingTab from './settings/tabs/UserSecuritySettingTab.vue';
import UserTwoFactorAuthSettingTab from './settings/tabs/UserTwoFactorAuthSettingTab.vue';
import UserDataManagementSettingTab from './settings/tabs/UserDataManagementSettingTab.vue';

import { ref, useTemplateRef, watch } from 'vue';
import { useRouter, onBeforeRouteUpdate } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import {
    mdiAccountOutline,
    mdiLockOpenOutline,
    mdiOnepassword,
    mdiDatabaseCogOutline
} from '@mdi/js';

type TwoFactorSettingTabType = InstanceType<typeof UserTwoFactorAuthSettingTab>;

const props = defineProps<{
    initTab?: string;
}>();

const router = useRouter();

const { tt } = useI18n();

const ALL_TABS: string[] = [
    'basicSetting',
    'securitySetting',
    'twoFactorSetting',
    'dataManagementSetting'
];

const twoFactorSettingTab = useTemplateRef<TwoFactorSettingTabType>('twoFactorSettingTab');

const activeTab = ref<string>((() => {
    let queryActiveTab = props.initTab || 'basicSetting';

    if (ALL_TABS.indexOf(queryActiveTab) < 0) {
        queryActiveTab = 'basicSetting';
    }

    return queryActiveTab;
})());

const pushRouter = (tab: string) => {
    router.push(`/user/settings?tab=${tab}`);
};

onBeforeRouteUpdate((to) => {
    if (to.query && to.query['tab'] && ALL_TABS.indexOf(to.query['tab'] as string) >= 0) {
        activeTab.value = to.query['tab'] as string;
    } else {
        activeTab.value = 'basicSetting';
    }
});

watch(activeTab, (newValue, oldValue) => {
    if (oldValue === 'twoFactorSetting' && newValue !== 'twoFactorSetting') {
        twoFactorSettingTab.value?.reset();
    }
});
</script>
