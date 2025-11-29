<template>
    <f7-page>
        <f7-navbar :title="tt('Statistics Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('Common Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultChartDataTypePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Chart Data Type') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allChartDataTypes, defaultChartDataType) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Chart Data Type')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Chart Data Type')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allChartDataTypes"
                                           v-model:show="showDefaultChartDataTypePopup"
                                           v-model="defaultChartDataType">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultTimezoneTypePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Timezone Used for Date Range') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTimezoneTypesUsedForStatistics, defaultTimezoneType) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Timezone Used for Date Range')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Timezone Type')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTimezoneTypesUsedForStatistics"
                                           v-model:show="showDefaultTimezoneTypePopup"
                                           v-model="defaultTimezoneType">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item :title="tt('Default Account Filter')" link="/settings/filter/account?type=statisticsDefault"></f7-list-item>

            <f7-list-item :title="tt('Default Transaction Category Filter')" link="/settings/filter/category?type=statisticsDefault"></f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultSortingTypePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Sort Order') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allSortingTypes, defaultSortingType) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Sort Order')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Sort Order')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allSortingTypes"
                                           v-model:show="showDefaultSortingTypePopup"
                                           v-model="defaultSortingType">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Categorical Analysis Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultCategoricalChartTypePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Chart Type') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allCategoricalChartTypes, defaultCategoricalChartType) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Chart Type')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Chart Type')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCategoricalChartTypes"
                                           v-model:show="showDefaultCategoricalChartTypePopup"
                                           v-model="defaultCategoricalChartType">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultCategoricalChartDateRangePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Date Range') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allCategoricalChartDateRanges, defaultCategoricalChartDateRange) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Date Range')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Date Range')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCategoricalChartDateRanges"
                                           v-model:show="showDefaultCategoricalChartDateRangePopup"
                                           v-model="defaultCategoricalChartDateRange">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Trend Analysis Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultTrendChartDateRangePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Date Range') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allTrendChartDateRanges, defaultTrendChartDateRange) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Date Range')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Date Range')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTrendChartDateRanges"
                                           v-model:show="showDefaultTrendChartDateRangePopup"
                                           v-model="defaultTrendChartDateRange">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Asset Trends Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                class="item-truncate-after-text"
                link="#"
                @click="showDefaultAssetTrendsChartDateRangePopup = true"
            >
                <template #after-title>
                    <div class="item-actual-title">
                        <span>{{ tt('Default Date Range') }}</span>
                    </div>
                </template>
                <template #after>
                    {{ findDisplayNameByType(allAssetTrendsChartDateRanges, defaultAssetTrendsChartDateRange) }}
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="type" value-field="type"
                                           title-field="displayName"
                                           :title="tt('Default Date Range')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Date Range')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allAssetTrendsChartDateRanges"
                                           v-model:show="showDefaultAssetTrendsChartDateRangePopup"
                                           v-model="defaultAssetTrendsChartDateRange">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useStatisticsSettingPageBase } from '@/views/base/statistics/StatisticsSettingPageBase.ts';

const { tt } = useI18n();
const {
    allChartDataTypes,
    allTimezoneTypesUsedForStatistics,
    allSortingTypes,
    allCategoricalChartTypes,
    allCategoricalChartDateRanges,
    allTrendChartDateRanges,
    allAssetTrendsChartDateRanges,
    defaultChartDataType,
    defaultTimezoneType,
    defaultSortingType,
    defaultCategoricalChartType,
    defaultCategoricalChartDateRange,
    defaultTrendChartDateRange,
    defaultAssetTrendsChartDateRange
} = useStatisticsSettingPageBase();

import { findDisplayNameByType } from '@/lib/common.ts';

const showDefaultChartDataTypePopup = ref<boolean>(false);
const showDefaultTimezoneTypePopup = ref<boolean>(false);
const showDefaultSortingTypePopup = ref<boolean>(false);
const showDefaultCategoricalChartTypePopup = ref<boolean>(false);
const showDefaultCategoricalChartDateRangePopup = ref<boolean>(false);
const showDefaultTrendChartDateRangePopup = ref<boolean>(false);
const showDefaultAssetTrendsChartDateRangePopup = ref<boolean>(false);
</script>
