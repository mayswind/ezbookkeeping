<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Statistics Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Data Type')"
                                    :placeholder="tt('Default Chart Data Type')"
                                    :items="allChartDataTypes"
                                    v-model="defaultChartDataType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Timezone Used for Date Range')"
                                    :placeholder="tt('Timezone Used for Date Range')"
                                    :items="allTimezoneTypesUsedForStatistics"
                                    v-model="defaultTimezoneType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Sort Order')"
                                    :placeholder="tt('Default Sort Order')"
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
            <v-card :title="tt('Categorical Analysis Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Type')"
                                    :placeholder="tt('Default Chart Type')"
                                    :items="allCategoricalChartTypes"
                                    v-model="defaultCategoricalChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Date Range')"
                                    :placeholder="tt('Default Date Range')"
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
            <v-card :title="tt('Trend Analysis Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Type')"
                                    :placeholder="tt('Default Chart Type')"
                                    :items="allTrendChartTypes"
                                    v-model="defaultTrendChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Date Range')"
                                    :placeholder="tt('Default Date Range')"
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
            <v-card :title="tt('Asset Trends Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Type')"
                                    :placeholder="tt('Default Chart Type')"
                                    :items="allTrendChartTypes"
                                    v-model="defaultAssetTrendsChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Date Range')"
                                    :placeholder="tt('Default Date Range')"
                                    :items="allAssetTrendsChartDateRanges"
                                    v-model="defaultAssetTrendsChartDateRange"
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

<script setup lang="ts">
import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';

import { useI18n } from '@/locales/helpers.ts';
import { useStatisticsSettingPageBase } from '@/views/base/statistics/StatisticsSettingPageBase.ts';

const { tt } = useI18n();
const {
    allChartDataTypes,
    allTimezoneTypesUsedForStatistics,
    allSortingTypes,
    allCategoricalChartTypes,
    allCategoricalChartDateRanges,
    allTrendChartTypes,
    allTrendChartDateRanges,
    allAssetTrendsChartDateRanges,
    defaultChartDataType,
    defaultTimezoneType,
    defaultSortingType,
    defaultCategoricalChartType,
    defaultCategoricalChartDateRange,
    defaultTrendChartType,
    defaultTrendChartDateRange,
    defaultAssetTrendsChartType,
    defaultAssetTrendsChartDateRange
} = useStatisticsSettingPageBase();
</script>

