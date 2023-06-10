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
                            v-for="chartDataType in allChartDataTypes">{{ $t(chartDataType.name) }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :title="$t('Default Date Range')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date Range'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="defaultDateRange">
                    <option :value="dateRange.type"
                            :key="dateRange.type"
                            v-for="dateRange in allDateRanges">{{ $t(dateRange.name) }}</option>
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
                            v-for="sortingType in allSortingTypes">{{ $t(sortingType.name) }}</option>
                </select>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import datetimeConstants from '@/consts/datetime.js';
import statisticsConstants from '@/consts/statistics.js';

export default {
    computed: {
        allChartTypes() {
            return statisticsConstants.allChartTypes;
        },
        allChartDataTypes() {
            return statisticsConstants.allChartDataTypes;
        },
        allSortingTypes() {
            return statisticsConstants.allSortingTypes;
        },
        allDateRanges() {
            const allDateRanges = [];

            for (let dateRangeField in datetimeConstants.allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(datetimeConstants.allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRangeType = datetimeConstants.allDateRanges[dateRangeField];

                if (dateRangeType.type !== datetimeConstants.allDateRanges.Custom.type) {
                    allDateRanges.push(dateRangeType);
                }
            }

            return allDateRanges;
        },
        defaultChartType: {
            get: function () {
                return this.$settings.getStatisticsDefaultChartType();
            },
            set: function (value) {
                this.$settings.setStatisticsDefaultChartType(value);
            }
        },
        defaultChartDataType: {
            get: function () {
                return this.$settings.getStatisticsDefaultChartDataType();
            },
            set: function (value) {
                this.$settings.setStatisticsDefaultChartDataType(value);
            }
        },
        defaultDateRange: {
            get: function () {
                return this.$settings.getStatisticsDefaultDateRange();
            },
            set: function (value) {
                this.$settings.setStatisticsDefaultDateRange(value);
            }
        },
        defaultSortingType: {
            get: function () {
                return this.$settings.getStatisticsSortingType();
            },
            set: function (value) {
                this.$settings.setStatisticsSortingType(value);
            }
        }
    }
};
</script>
