<template>
    <f7-page>
        <f7-navbar :title="$t('Statistics Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item
                        :title="$t('Default Chart Type')"
                        smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done'), scrollToSelectedItem: true }">
                        <select v-model="defaultChartType">
                            <option :value="$constants.statistics.allChartTypes.Pie">{{ $t('Pie Chart') }}</option>
                            <option :value="$constants.statistics.allChartTypes.Bar">{{ $t('Bar Chart') }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        :title="$t('Default Chart Data Type')"
                        smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done'), scrollToSelectedItem: true }">
                        <select v-model="defaultChartDataType">
                            <option v-for="chartDataType in allChartDataTypes"
                                    :key="chartDataType.type"
                                    :value="chartDataType.type">{{ chartDataType.name | localized }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        :title="$t('Default Date Range')"
                        smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done'), scrollToSelectedItem: true }">
                        <select v-model="defaultDateRange">
                            <option v-for="dateRange in allDateRanges"
                                    :key="dateRange.type"
                                    :value="dateRange.type">{{ dateRange.name | localized }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        :title="$t('Sort By')"
                        smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done'), scrollToSelectedItem: true }">
                        <select v-model="sortBy">
                            <option :value="$constants.statistics.allSortingTypes.ByAmount">{{ $t('By Amount') }}</option>
                            <option :value="$constants.statistics.allSortingTypes.ByDisplayOrder">{{ $t('By Display Order') }}</option>
                            <option :value="$constants.statistics.allSortingTypes.ByName">{{ $t('By Name') }}</option>
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
        sortBy: {
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
