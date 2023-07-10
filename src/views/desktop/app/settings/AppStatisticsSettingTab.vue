<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="$t('Statistics Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Chart Type')"
                                    :placeholder="$t('Chart Type')"
                                    :items="[
                                        { type: allChartTypes.Pie, displayName: $t('Pie Chart') },
                                        { type: allChartTypes.Bar, displayName: $t('Bar Chart') }
                                    ]"
                                    v-model="defaultChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Chart Data Type')"
                                    :placeholder="$t('Default Chart Data Type')"
                                    :items="allChartDataTypes"
                                    v-model="defaultChartDataType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Date Range')"
                                    :placeholder="$t('Default Date Range')"
                                    :items="allDateRanges"
                                    v-model="defaultDateRange"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Sort By')"
                                    :placeholder="$t('Default Sort By')"
                                    :items="allSortingTypes"
                                    v-model="defaultSortingType"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="$t('Default Account Filter')">

            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="$t('Default Transaction Category Filter')">

            </v-card>
        </v-col>
    </v-row>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';

import statisticsConstants from '@/consts/statistics.js';

export default {
    computed: {
        ...mapStores(useSettingsStore),
        allChartTypes() {
            return statisticsConstants.allChartTypes;
        },
        allChartDataTypes() {
            return this.$locale.getAllStatisticsChartDataTypes();
        },
        allSortingTypes() {
            return this.$locale.getAllStatisticsSortingTypes();
        },
        allDateRanges() {
            return this.$locale.getAllDateRanges(false);
        },
        defaultChartType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultChartType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultChartType(value);
            }
        },
        defaultChartDataType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultChartDataType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultChartDataType(value);
            }
        },
        defaultDateRange: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultDataRangeType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultDateRange(value);
            }
        },
        defaultSortingType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultSortingType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsSortingType(value);
            }
        }
    }
};
</script>

