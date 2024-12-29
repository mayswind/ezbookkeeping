<template>
    <v-chart autoresize class="trends-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import { DEFAULT_ICON_COLOR, DEFAULT_CHART_COLORS } from '@/consts/color.ts';
import datetimeConstants from '@/consts/datetime.js';
import statisticsConstants from '@/consts/statistics.js';
import { ThemeType } from '@/core/theme.ts';
import {
    isArray,
    isNumber
} from '@/lib/common.ts';
import {
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange
} from '@/lib/datetime.js';
import {
    sortStatisticsItems,
    getAllDateRanges
} from '@/lib/statistics.js';

export default {
    props: [
        'skeleton',
        'type',
        'items',
        'startYearMonth',
        'endYearMonth',
        'sortingType',
        'dateAggregationType',
        'idField',
        'nameField',
        'valueField',
        'colorField',
        'hiddenField',
        'displayOrdersField',
        'translateName',
        'defaultCurrency',
        'showValue',
        'showTotalAmountInTooltip',
        'enableClickItem'
    ],
    emits: [
        'click'
    ],
    data() {
        return {
            selectedLegends: null
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        isDarkMode() {
            return this.globalTheme.global.name.value === ThemeType.Dark;
        },
        itemsMap: function () {
            const map = {};

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];
                let id = '';

                if (this.idField && item[this.idField]) {
                    id = item[this.idField];
                } else {
                    id = this.getItemName(item[this.nameField]);
                }

                map[id] = {
                    [this.idField || 'id']: id,
                    [this.nameField || 'name']: item[this.nameField],
                    [this.hiddenField || 'hidden']: item[this.hiddenField],
                    [this.displayOrdersField || 'displayOrders']: item[this.displayOrdersField]
                };
            }

            return map;
        },
        allDateRanges: function () {
            return getAllDateRanges(this.items, this.startYearMonth, this.endYearMonth, this.dateAggregationType);
        },
        allDisplayDateRanges: function () {
            const allDisplayDateRanges = [];

            for (let i = 0; i < this.allDateRanges.length; i++) {
                const dateRange = this.allDateRanges[i];

                if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Year.type) {
                    allDisplayDateRanges.push(this.$locale.formatUnixTimeToShortYear(this.userStore, dateRange.minUnixTime));
                } else if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Quarter.type) {
                    allDisplayDateRanges.push(this.$locale.formatYearQuarter(dateRange.year, dateRange.quarter));
                } else { // if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Month.type) {
                    allDisplayDateRanges.push(this.$locale.formatUnixTimeToShortYearMonth(this.userStore, dateRange.minUnixTime));
                }
            }

            return allDisplayDateRanges;
        },
        allSeries: function () {
            const allSeries = [];

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (!this.hiddenField || item[this.hiddenField]) {
                    continue;
                }

                const allAmounts = [];
                const dateRangeAmountMap = {};

                for (let j = 0; j < item.items.length; j++) {
                    const dataItem = item.items[j];
                    let dateRangeKey = '';

                    if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Year.type) {
                        dateRangeKey = dataItem.year;
                    } else if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Quarter.type) {
                        dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month - 1) / 3) + 1}`;
                    } else { // if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Month.type) {
                        dateRangeKey = `${dataItem.year}-${dataItem.month}`;
                    }

                    const dataItems = dateRangeAmountMap[dateRangeKey] || [];
                    dataItems.push(dataItem);

                    dateRangeAmountMap[dateRangeKey] = dataItems;
                }

                for (let j = 0; j < this.allDateRanges.length; j++) {
                    const dateRange = this.allDateRanges[j];
                    let dateRangeKey = '';

                    if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Year.type) {
                        dateRangeKey = dateRange.year;
                    } else if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Quarter.type) {
                        dateRangeKey = `${dateRange.year}-${dateRange.quarter}`;
                    } else { // if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Month.type) {
                        dateRangeKey = `${dateRange.year}-${dateRange.month + 1}`;
                    }

                    let amount = 0;
                    const dataItems = dateRangeAmountMap[dateRangeKey];

                    if (isArray(dataItems)) {
                        for (let i = 0; i < dataItems.length; i++) {
                            const dataItem = dataItems[i];

                            if (isNumber(dataItem[this.valueField])) {
                                amount += dataItem[this.valueField];
                            }
                        }
                    }

                    allAmounts.push(amount);
                }

                const finalItem = {
                    id: (this.idField && item[this.idField]) ? item[this.idField] : this.getItemName(item[this.nameField]),
                    name: (this.idField && item[this.idField]) ? item[this.idField] : this.getItemName(item[this.nameField]),
                    itemStyle: {
                        color: this.getColor(item[this.colorField] ? item[this.colorField] : DEFAULT_CHART_COLORS[i % DEFAULT_CHART_COLORS.length]),
                    },
                    selected: true,
                    type: 'line',
                    stack: 'a',
                    animation: !self.skeleton,
                    data: allAmounts
                };

                if (this.type === statisticsConstants.allTrendChartTypes.Area) {
                    finalItem.areaStyle = {};
                } else if (this.type === statisticsConstants.allTrendChartTypes.Column) {
                    finalItem.type = 'bar';
                }

                allSeries.push(finalItem);
            }

            return allSeries;
        },
        yAxisWidth: function () {
            let maxValue = Number.MIN_SAFE_INTEGER;
            let minValue = Number.MAX_SAFE_INTEGER;
            let width = 90;

            if (!this.allSeries || !this.allSeries.length) {
                return width;
            }

            for (let i = 0; i < this.allSeries.length; i++) {
                for (let j = 0; j < this.allSeries[i].data.length; j++) {
                    const value = this.allSeries[i].data[j];

                    if (value > maxValue) {
                        maxValue = value;
                    }

                    if (value < minValue) {
                        minValue = value;
                    }
                }
            }

            const maxValueText = this.getDisplayCurrency(maxValue, this.defaultCurrency);
            const minValueText = this.getDisplayCurrency(minValue, this.defaultCurrency);
            const maxLengthText = maxValueText.length > minValueText.length ? maxValueText : minValueText;

            const canvas = document.createElement('canvas');
            const context = canvas.getContext('2d');
            context.font = '12px Arial';

            const textMetrics = context.measureText(maxLengthText);
            const actualWidth = Math.round(textMetrics.width) + 20;

            if (actualWidth >= 200) {
                width = 200;
            } if (actualWidth > 90) {
                width = actualWidth;
            }

            return width;
        },
        chartOptions: function () {
            const self = this;

            return {
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        label: {
                            backgroundColor: self.isDarkMode ? '#333' : '#fff',
                            color: self.isDarkMode ? '#eee' : '#333'
                        },
                    },
                    backgroundColor: self.isDarkMode ? '#333' : '#fff',
                    borderColor: self.isDarkMode ? '#333' : '#fff',
                    textStyle: {
                        color: self.isDarkMode ? '#eee' : '#333'
                    },
                    formatter: params => {
                        let tooltip = '';
                        let totalAmount = 0;
                        const displayItems = [];

                        for (let i = 0; i < params.length; i++) {
                            const id = params[i].seriesId;
                            const name = self.itemsMap[id] && self.nameField && self.itemsMap[id][self.nameField] ? self.getItemName(self.itemsMap[id][self.nameField]) : id;
                            const color = params[i].color;
                            const displayOrders = self.itemsMap[id] && self.displayOrdersField && self.itemsMap[id][self.displayOrdersField] ? self.itemsMap[id][self.displayOrdersField] : [0];
                            const amount = params[i].data;

                            displayItems.push({
                                name: name,
                                color: color,
                                displayOrders: displayOrders,
                                totalAmount: amount
                            });

                            totalAmount += amount;
                        }

                        sortStatisticsItems(displayItems, self.sortingType);

                        for (let i = 0; i < displayItems.length; i++) {
                            const item = displayItems[i];

                            if (displayItems.length === 1 || item.totalAmount !== 0) {
                                const value = self.getDisplayCurrency(item.totalAmount, self.defaultCurrency);
                                tooltip += '<div><span class="chart-pointer" style="background-color: ' + item.color + '"></span>';
                                tooltip += `<span>${item.name}</span><span style="margin-left: 20px; float: right">${value}</span><br/>`;
                                tooltip += '</div>';
                            }
                        }

                        if (self.showTotalAmountInTooltip) {
                            const displayTotalAmount = self.getDisplayCurrency(totalAmount, self.defaultCurrency);
                            tooltip = '<div style="border-bottom: ' + (self.isDarkMode ? '#eee' : '#333') + ' dashed 1px">'
                                + '<span class="chart-pointer" style="background-color: ' + (self.isDarkMode ? '#eee' : '#333') + '"></span>'
                                + `<span>${self.$t('Total Amount')}</span><span style="margin-left: 20px; float: right">${displayTotalAmount}</span><br/>`
                                + '</div>' + tooltip;
                        }

                        if (params.length && params[0].name) {
                            tooltip = `${params[0].name}<br/>` + tooltip;
                        }

                        return tooltip;
                    }
                },
                legend: {
                    orient: 'horizontal',
                    data: self.allSeries.map(item => item.name),
                    selected: self.selectedLegends,
                    textStyle: {
                        color: self.isDarkMode ? '#eee' : '#333'
                    },
                    formatter: id => {
                        return self.itemsMap[id] && self.nameField && self.itemsMap[id][self.nameField] ? self.getItemName(self.itemsMap[id][self.nameField]) : id;
                    }
                },
                grid: {
                    left: self.yAxisWidth,
                    right: 20
                },
                xAxis: [
                    {
                        type: 'category',
                        data: self.allDisplayDateRanges
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        axisLabel: {
                            formatter: function (value) {
                                return self.getDisplayCurrency(value, self.defaultCurrency);
                            }
                        },
                        axisPointer: {
                            label: {
                                formatter: function (params) {
                                    return self.getDisplayCurrency(parseInt(params.value), self.defaultCurrency);
                                }
                            }
                        }
                    }
                ],
                series: self.allSeries
            }
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        clickItem: function (e) {
            if (!this.enableClickItem || e.componentType !== 'series') {
                return;
            }

            const id = e.seriesId;
            const item = this.itemsMap[id];
            const itemId = this.idField ? item[this.idField] : '';
            const dateRange = this.allDateRanges[e.dataIndex];
            let minUnixTime = dateRange.minUnixTime;
            let maxUnixTime = dateRange.maxUnixTime;

            if (this.startYearMonth) {
                const startMinUnixTime = getYearMonthFirstUnixTime(this.startYearMonth);

                if (startMinUnixTime > minUnixTime) {
                    minUnixTime = startMinUnixTime;
                }
            }

            if (this.endYearMonth) {
                const endMaxUnixTime = getYearMonthLastUnixTime(this.endYearMonth);

                if (endMaxUnixTime < maxUnixTime) {
                    maxUnixTime = endMaxUnixTime;
                }
            }

            const dateRangeType = getDateTypeByDateRange(minUnixTime, maxUnixTime, this.userStore.currentUserFirstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

            this.$emit('click', {
                itemId: itemId,
                dateRange: {
                    minTime: minUnixTime,
                    maxTime: maxUnixTime,
                    type: dateRangeType
                }
            });
        },
        getColor: function (color) {
            if (color && color !== DEFAULT_ICON_COLOR) {
                color = '#' + color;
            }

            return color;
        },
        getItemName(name) {
            return this.translateName ? this.$t(name) : name;
        },
        onLegendSelectChanged: function (e) {
            this.selectedLegends = e.selected;
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        }
    }
}
</script>

<style scoped>
.trends-chart-container {
    width: 100%;
    height: 560px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .pie-chart-container {
        height: 500px;
    }
}
</style>
