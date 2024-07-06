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
                                    :label="$t('Timezone Used for Date Range')"
                                    :placeholder="$t('Timezone Used for Date Range')"
                                    :items="allTimezoneTypesUsedForStatistics"
                                    v-model="defaultTimezoneType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Sort Order')"
                                    :placeholder="$t('Default Sort Order')"
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
            <v-card :title="$t('Categorical Analysis Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Chart Type')"
                                    :placeholder="$t('Default Chart Type')"
                                    :items="allCategoricalChartTypes"
                                    v-model="defaultCategoricalChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Date Range')"
                                    :placeholder="$t('Default Date Range')"
                                    :items="allCategoricalChartDateRanges"
                                    v-model="defaultCategoricalChartDateRange"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="$t('Trend Analysis Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Chart Type')"
                                    :placeholder="$t('Default Chart Type')"
                                    :items="allTrendChartTypes"
                                    v-model="defaultTrendChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="$t('Default Date Range')"
                                    :placeholder="$t('Default Date Range')"
                                    :items="allTrendChartDateRanges"
                                    v-model="defaultTrendChartDateRange"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <account-filter-settings-card type="statisticsDefault" :auto-save="true" />
        </v-col>

        <v-col cols="12">
            <category-filter-settings-card type="statisticsDefault" :auto-save="true" />
        </v-col>
    </v-row>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';

import datetimeConstants from '@/consts/datetime.js';

import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';

export default {
    components: {
        AccountFilterSettingsCard,
        CategoryFilterSettingsCard
    },
    computed: {
        ...mapStores(useSettingsStore),
        allChartDataTypes() {
            return this.$locale.getAllStatisticsChartDataTypes();
        },
        allTimezoneTypesUsedForStatistics() {
            return this.$locale.getAllTimezoneTypesUsedForStatistics();
        },
        allSortingTypes() {
            return this.$locale.getAllStatisticsSortingTypes();
        },
        allCategoricalChartTypes() {
            return this.$locale.getAllCategoricalChartTypes();
        },
        allCategoricalChartDateRanges() {
            return this.$locale.getAllDateRanges(datetimeConstants.allDateRangeScenes.Normal, false);
        },
        allTrendChartTypes() {
            return this.$locale.getAllTrendChartTypes();
        },
        allTrendChartDateRanges() {
            return this.$locale.getAllDateRanges(datetimeConstants.allDateRangeScenes.TrendAnalysis, false);
        },
        defaultChartDataType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultChartDataType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultChartDataType(value);
            }
        },
        defaultTimezoneType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultTimezoneType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultTimezoneType(value);
            }
        },
        defaultSortingType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultSortingType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsSortingType(value);
            }
        },
        defaultCategoricalChartType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultCategoricalChartType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultCategoricalChartType(value);
            }
        },
        defaultCategoricalChartDateRange: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultCategoricalChartDataRangeType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultCategoricalChartDateRange(value);
            }
        },
        defaultTrendChartType: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultTrendChartType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultTrendChartType(value);
            }
        },
        defaultTrendChartDateRange: {
            get: function () {
                return this.settingsStore.appSettings.statistics.defaultTrendChartDataRangeType;
            },
            set: function (value) {
                this.settingsStore.setStatisticsDefaultTrendChartDateRange(value);
            }
        },
    }
};
</script>

