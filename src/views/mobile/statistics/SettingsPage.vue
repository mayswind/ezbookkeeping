<template>
    <f7-page>
        <f7-navbar :title="tt('Statistics Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ tt('Common Settings') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :title="tt('Default Chart Data Type')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Chart Data Type'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="defaultChartDataType">
                    <option :value="chartDataType.type"
                            :key="chartDataType.type"
                            v-for="chartDataType in allChartDataTypes">{{ chartDataType.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :title="tt('Timezone Used for Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Timezone Type'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="defaultTimezoneType">
                    <option :value="timezoneType.type"
                            :key="timezoneType.type"
                            v-for="timezoneType in allTimezoneTypesUsedForStatistics">{{ timezoneType.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item :title="tt('Default Account Filter')" link="/settings/filter/account?type=statisticsDefault"></f7-list-item>

            <f7-list-item :title="tt('Default Transaction Category Filter')" link="/settings/filter/category?type=statisticsDefault"></f7-list-item>

            <f7-list-item
                :title="tt('Default Sort Order')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Sort Order'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="defaultSortingType">
                    <option :value="sortingType.type"
                            :key="sortingType.type"
                            v-for="sortingType in allSortingTypes">{{ sortingType.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Categorical Analysis Settings') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :title="tt('Default Chart Type')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Chart Type'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="defaultCategoricalChartType">
                    <option :value="chartType.type"
                            :key="chartType.type"
                            v-for="chartType in allCategoricalChartTypes">{{ chartType.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :title="tt('Default Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Date Range'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="defaultCategoricalChartDateRange">
                    <option :value="dateRange.type"
                            :key="dateRange.type"
                            v-for="dateRange in allCategoricalChartDateRanges">{{ dateRange.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ tt('Trend Analysis Settings') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :title="tt('Default Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Date Range'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), popupCloseLinkText: tt('Done') }">
                <select v-model="defaultTrendChartDateRange">
                    <option :value="dateRange.type"
                            :key="dateRange.type"
                            v-for="dateRange in allTrendChartDateRanges">{{ dateRange.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
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
</script>
