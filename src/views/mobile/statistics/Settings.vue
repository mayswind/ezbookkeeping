<template>
    <f7-page>
        <f7-navbar :title="$t('Statistics Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list dividers>
                    <f7-list-item
                        :title="$t('Default Chart Type')"
                        smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, popupSwipeToClose: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Chart Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                        <select v-model="defaultChartType">
                            <option :value="$constants.statistics.allChartTypes.Pie">{{ $t('Pie Chart') }}</option>
                            <option :value="$constants.statistics.allChartTypes.Bar">{{ $t('Bar Chart') }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        :title="$t('Default Chart Data Type')"
                        smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, popupSwipeToClose: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Chart Data Type'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                        <select v-model="defaultChartDataType">
                            <option v-for="chartDataType in allChartDataTypes"
                                    :key="chartDataType.type"
                                    :value="chartDataType.type">{{ $t(chartDataType.name) }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        :title="$t('Default Date Range')"
                        smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, popupSwipeToClose: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date Range'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                        <select v-model="defaultDateRange">
                            <option v-for="dateRange in allDateRanges"
                                    :key="dateRange.type"
                                    :value="dateRange.type">{{ $t(dateRange.name) }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item :title="$t('Default Account Filter')" link="/statistic/filter/account?modifyDefault=1"></f7-list-item>

                    <f7-list-item :title="$t('Default Transaction Category Filter')" link="/statistic/filter/category?modifyDefault=1"></f7-list-item>

                    <f7-list-item
                        :title="$t('Default Sort By')"
                        smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, popupSwipeToClose: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Sort By'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                        <select v-model="defaultSortingType">
                            <option v-for="sortingType in allSortingTypes"
                                    :key="sortingType.type"
                                    :value="sortingType.type">{{ $t(sortingType.name) }}</option>
                        </select>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    computed: {
        allChartDataTypes() {
            return this.$constants.statistics.allChartDataTypes;
        },
        allSortingTypes() {
            return this.$constants.statistics.allSortingTypes;
        },
        allDateRanges() {
            const allDateRanges = [];

            for (let dateRangeField in this.$constants.datetime.allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(this.$constants.datetime.allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRangeType = this.$constants.datetime.allDateRanges[dateRangeField];

                if (dateRangeType.type !== this.$constants.datetime.allDateRanges.Custom.type) {
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
