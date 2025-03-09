<template>
    <f7-page>
        <f7-navbar :title="tt('Statistics Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('Common Settings') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                link="#"
                :title="tt('Default Chart Data Type')"
                :after="findDisplayNameByType(allChartDataTypes, defaultChartDataType)"
                @click="showDefaultChartDataTypePopup = true"
            >
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
                link="#"
                :title="tt('Timezone Used for Date Range')"
                :after="findDisplayNameByType(allTimezoneTypesUsedForStatistics, defaultTimezoneType)"
                @click="showDefaultTimezoneTypePopup = true"
            >
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
                link="#"
                :title="tt('Default Sort Order')"
                :after="findDisplayNameByType(allSortingTypes, defaultSortingType)"
                @click="showDefaultSortingTypePopup = true"
            >
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
        <f7-list strong inset dividers>
            <f7-list-item
                link="#"
                :title="tt('Default Chart Type')"
                :after="findDisplayNameByType(allCategoricalChartTypes, defaultCategoricalChartType)"
                @click="showDefaultCategoricalChartTypePopup = true"
            >
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
                link="#"
                :title="tt('Default Date Range')"
                :after="findDisplayNameByType(allCategoricalChartDateRanges, defaultCategoricalChartDateRange)"
                @click="showDefaultCategoricalChartDateRangePopup = true"
            >
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
        <f7-list strong inset dividers>
            <f7-list-item
                link="#"
                :title="tt('Default Date Range')"
                :after="findDisplayNameByType(allTrendChartDateRanges, defaultTrendChartDateRange)"
                @click="showDefaultTrendChartDateRangePopup = true"
            >
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
    defaultChartDataType,
    defaultTimezoneType,
    defaultSortingType,
    defaultCategoricalChartType,
    defaultCategoricalChartDateRange,
    defaultTrendChartDateRange
} = useStatisticsSettingPageBase();

import { findDisplayNameByType } from '@/lib/common.ts';

const showDefaultChartDataTypePopup = ref<boolean>(false);
const showDefaultTimezoneTypePopup = ref<boolean>(false);
const showDefaultSortingTypePopup = ref<boolean>(false);
const showDefaultCategoricalChartTypePopup = ref<boolean>(false);
const showDefaultCategoricalChartDateRangePopup = ref<boolean>(false);
const showDefaultTrendChartDateRangePopup = ref<boolean>(false);
</script>
