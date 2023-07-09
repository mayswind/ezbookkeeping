<template>
    <div>
        <v-tabs show-arrows class="text-uppercase" v-model="activeTab">
            <v-tab value="basicSetting">
                <v-icon size="20" start :icon="icons.basicSetting"/>
                {{ $t('Basic') }}
            </v-tab>
            <v-tab value="applicationLockSetting">
                <v-icon size="20" start :icon="icons.applicationLockSetting"/>
                {{ $t('Application Lock') }}
            </v-tab>
            <v-tab value="statisticsSetting">
                <v-icon size="20" start :icon="icons.statisticsSetting"/>
                {{ $t('Statistics') }}
            </v-tab>
        </v-tabs>
        <v-divider />
        <v-window class="mt-5 disable-tab-transition" v-model="activeTab">
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
import AppBasicSettingTab from './settings/AppBasicSettingTab.vue';
import AppLockSettingTab from './settings/AppLockSettingTab.vue';
import AppStatisticsSettingTab from './settings/AppStatisticsSettingTab.vue';

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
        'tab'
    ],
    data() {
        let queryActiveTab = this.tab || 'basicSetting';

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
    }
}
</script>
