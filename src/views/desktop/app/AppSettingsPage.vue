<template>
    <div>
        <v-tabs show-arrows v-model="activeTab">
            <v-tab value="basicSetting" @click="pushRouter('basicSetting')">
                <v-icon size="20" start :icon="icons.basicSetting"/>
                {{ $t('Basic') }}
            </v-tab>
            <v-tab value="applicationLockSetting" @click="pushRouter('applicationLockSetting')">
                <v-icon size="20" start :icon="icons.applicationLockSetting"/>
                {{ $t('Application Lock') }}
            </v-tab>
            <v-tab value="statisticsSetting" @click="pushRouter('statisticsSetting')">
                <v-icon size="20" start :icon="icons.statisticsSetting"/>
                {{ $t('Statistics') }}
            </v-tab>
        </v-tabs>

        <v-window class="mt-4 disable-tab-transition" v-model="activeTab">
            <v-window-item value="basicSetting">
                <app-basic-setting-tab/>
            </v-window-item>

            <v-window-item value="applicationLockSetting">
                <app-lock-setting-tab/>
            </v-window-item>

            <v-window-item value="statisticsSetting">
                <app-statistics-setting-tab/>
            </v-window-item>
        </v-window>
    </div>
</template>

<script>
import AppBasicSettingTab from './settings/tabs/AppBasicSettingTab.vue';
import AppLockSettingTab from './settings/tabs/AppLockSettingTab.vue';
import AppStatisticsSettingTab from './settings/tabs/AppStatisticsSettingTab.vue';

import {
    mdiCogOutline,
    mdiLockOpenOutline,
    mdiChartPieOutline
} from '@mdi/js';

export default {
    components: {
        AppBasicSettingTab,
        AppLockSettingTab,
        AppStatisticsSettingTab
    },
    props: [
        'initTab'
    ],
    data() {
        let queryActiveTab = this.initTab || 'basicSetting';

        if ([
            'basicSetting',
            'applicationLockSetting',
            'statisticsSetting'
        ].indexOf(queryActiveTab) === -1) {
            queryActiveTab = 'basicSetting';
        }

        return {
            activeTab: queryActiveTab,
            icons: {
                basicSetting: mdiCogOutline,
                applicationLockSetting: mdiLockOpenOutline,
                statisticsSetting: mdiChartPieOutline
            }
        };
    },
    beforeRouteUpdate(to) {
        if (to.query && to.query.tab && [
            'basicSetting',
            'applicationLockSetting',
            'statisticsSetting'
        ].indexOf(to.query.tab) >= 0) {
            this.activeTab = to.query.tab;
        } else {
            this.activeTab = 'basicSetting';
        }
    },
    methods: {
        pushRouter(tab) {
            this.$router.push(`/app/settings?tab=${tab}`);
        }
    }
}
</script>
