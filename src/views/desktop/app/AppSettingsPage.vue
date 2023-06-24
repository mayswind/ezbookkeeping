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
                <app-lock-seting-tab/>
            </v-window-item>

            <v-window-item value="statisticsSetting">
                <app-statistics-setting-tab/>
            </v-window-item>
        </v-window>
    </div>
</template>

<script>
import AppBasicSettingTab from './settings/AppBasicSettingTab.vue';
import AppLockSetingTab from './settings/AppLockSetingTab.vue';
import AppStatisticsSettingTab from './settings/AppStatisticsSettingTab.vue';

import {
    mdiCogOutline,
    mdiLockOpenOutline,
    mdiChartPieOutline
} from '@mdi/js';

export default {
    components: {
        AppBasicSettingTab,
        AppLockSetingTab,
        AppStatisticsSettingTab
    },
    data() {
        return {
            activeTab: 'basicSetting',
            icons: {
                basicSetting: mdiCogOutline,
                applicationLockSetting: mdiLockOpenOutline,
                statisticsSetting: mdiChartPieOutline
            }
        };
    }
}
</script>
