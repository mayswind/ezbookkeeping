<template>
    <f7-page>
        <f7-navbar :title="$t('Statistics Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ $t('Common Settings') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :title="$t('Default Chart Data Type')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Chart Data Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultChartDataType">
                    <option :value="chartDataType.type"
                            :key="chartDataType.type"
                            v-for="chartDataType in allChartDataTypes">{{ chartDataType.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :title="$t('Timezone Used for Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Timezone Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultTimezoneType">
                    <option :value="timezoneType.type"
                            :key="timezoneType.type"
                            v-for="timezoneType in allTimezoneTypesUsedForStatistics">{{ timezoneType.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item :title="$t('Default Account Filter')" link="/settings/filter/account?type=statisticsDefault"></f7-list-item>

            <f7-list-item :title="$t('Default Transaction Category Filter')" link="/settings/filter/category?type=statisticsDefault"></f7-list-item>

            <f7-list-item
                :title="$t('Default Sort Order')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Sort Order'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultSortingType">
                    <option :value="sortingType.type"
                            :key="sortingType.type"
                            v-for="sortingType in allSortingTypes">{{ sortingType.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>

        <f7-block-title>{{ $t('Categorical Analysis Settings') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :title="$t('Default Chart Type')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Chart Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultCategoricalChartType">
                    <option :value="chartType.type"
                            :key="chartType.type"
                            v-for="chartType in allCategoricalChartTypes">{{ chartType.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :title="$t('Default Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date Range'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultCategoricalChartDateRange">
                    <option :value="dateRange.type"
                            :key="dateRange.type"
                            v-for="dateRange in allCategoricalChartDateRanges">{{ dateRange.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>

<!--        <f7-block-title>{{ $t('Trend Analysis Settings') }}</f7-block-title>-->
<!--        <f7-list strong inset dividers>-->
<!--            <f7-list-item-->
<!--                :title="$t('Default Chart Type')"-->
<!--                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Chart Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">-->
<!--                <select v-model="defaultTrendChartType">-->
<!--                    <option :value="chartType.type"-->
<!--                            :key="chartType.type"-->
<!--                            v-for="chartType in allTrendChartTypes">{{ chartType.displayName }}</option>-->
<!--                </select>-->
<!--            </f7-list-item>-->

<!--            <f7-list-item-->
<!--                :title="$t('Default Date Range')"-->
<!--                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date Range'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">-->
<!--                <select v-model="defaultTrendChartDateRange">-->
<!--                    <option :value="dateRange.type"-->
<!--                            :key="dateRange.type"-->
<!--                            v-for="dateRange in allTrendChartDateRanges">{{ dateRange.displayName }}</option>-->
<!--                </select>-->
<!--            </f7-list-item>-->
<!--        </f7-list>-->
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';

import datetimeConstants from '@/consts/datetime.js';
import statisticsConstants from '@/consts/statistics.js';

export default {
    computed: {
        ...mapStores(useSettingsStore),
        allChartDataTypes() {
            return this.$locale.getAllStatisticsChartDataTypes(statisticsConstants.allAnalysisTypes.CategoricalAnalysis);
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
