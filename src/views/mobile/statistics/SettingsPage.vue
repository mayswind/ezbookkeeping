<template>
    <f7-page>
        <f7-navbar :title="$t('Statistics Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top">
            <f7-list-item
                :title="$t('Default Chart Type')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Chart Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultChartType">
                    <option :value="allChartTypes.Pie">{{ $t('Pie Chart') }}</option>
                    <option :value="allChartTypes.Bar">{{ $t('Bar Chart') }}</option>
                </select>
            </f7-list-item>

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
                :title="$t('Default Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date Range'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultDateRange">
                    <option :value="dateRange.type"
                            :key="dateRange.type"
                            v-for="dateRange in allDateRanges">{{ dateRange.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item :title="$t('Default Account Filter')" link="/statistic/filter/account?modifyDefault=1"></f7-list-item>

            <f7-list-item :title="$t('Default Transaction Category Filter')" link="/statistic/filter/category?modifyDefault=1"></f7-list-item>

            <f7-list-item
                :title="$t('Default Sort By')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Sort By'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultSortingType">
                    <option :value="sortingType.type"
                            :key="sortingType.type"
                            v-for="sortingType in allSortingTypes">{{ sortingType.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>
    </f7-page>
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
