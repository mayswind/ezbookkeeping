<template>
    <v-chart autoresize class="trends-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import colorConstants from '@/consts/color.js';
import statisticsConstants from '@/consts/statistics.js';
import { isNumber } from '@/lib/common.js';
import {
    getYearMonthStringFromObject,
    getAllYearMonthUnixTimesBetweenStartYearMonthAndEndYearMonth
} from '@/lib/datetime.js';

export default {
    props: [
        'skeleton',
        'type',
        'items',
        'startYearMonth',
        'endYearMonth',
        'idField',
        'nameField',
        'valueField',
        'colorField',
        'hiddenField',
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
            return this.globalTheme.global.name.value === 'dark';
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

                map[id] = item;
            }

            return map;
        },
        allYearMonthTimes: function () {
            if (this.startYearMonth && this.endYearMonth) {
                return getAllYearMonthUnixTimesBetweenStartYearMonthAndEndYearMonth(this.startYearMonth, this.endYearMonth);
            } else if (this.items && this.items.length) {
                let minYear = Number.MAX_SAFE_INTEGER, minMonth = Number.MAX_SAFE_INTEGER, maxYear = 0, maxMonth = 0;

                for (let i = 0; i < this.items.length; i++) {
                    const item = this.items[i];

                    for (let j = 0; j < item.items.length; j++) {
                        const dataItem = item.items[j];

                        if (dataItem.year < minYear || (dataItem.year === minYear && dataItem.month < minMonth)) {
                            minYear = dataItem.year;
                            minMonth = dataItem.month;
                        }

                        if (dataItem.year > maxYear || (dataItem.year === maxYear && dataItem.month > maxMonth)) {
                            maxYear = dataItem.year;
                            maxMonth = dataItem.month;
                        }
                    }
                }

                return getAllYearMonthUnixTimesBetweenStartYearMonthAndEndYearMonth(`${minYear}-${minMonth}`, `${maxYear}-${maxMonth}`);
            }

            return [];
        },
        allDisplayMonths: function () {
            const allDisplayMonths = [];

            for (let i = 0; i < this.allYearMonthTimes.length; i++) {
                const yearMonthTime = this.allYearMonthTimes[i];
                allDisplayMonths.push(this.$locale.formatUnixTimeToShortYearMonth(this.userStore, yearMonthTime.minUnixTime));
            }

            return allDisplayMonths;
        },
        allSeries: function () {
            const allSeries = [];

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (!this.hiddenField || item[this.hiddenField]) {
                    continue;
                }

                const allAmounts = [];
                const yearMonthDataMap = {};

                for (let j = 0; j < item.items.length; j++) {
                    const dataItem = item.items[j];
                    yearMonthDataMap[`${dataItem.year}-${dataItem.month}`] = dataItem;
                }

                for (let j = 0; j < this.allYearMonthTimes.length; j++) {
                    const yearMonth = getYearMonthStringFromObject(this.allYearMonthTimes[j]);
                    const dataItem = yearMonthDataMap[yearMonth];
                    let amount = 0;

                    if (dataItem && isNumber(dataItem[this.valueField])) {
                        amount = dataItem[this.valueField];
                    }

                    allAmounts.push(amount);
                }

                const finalItem = {
                    id: (this.idField && item[this.idField]) ? item[this.idField] : this.getItemName(item[this.nameField]),
                    name: (this.idField && item[this.idField]) ? item[this.idField] : this.getItemName(item[this.nameField]),
                    itemStyle: {
                        color: this.getColor(item[this.colorField] ? item[this.colorField] : colorConstants.defaultChartColors[i % colorConstants.defaultChartColors.length]),
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
            let maxLengthText = maxValueText.length > minValueText.length ? maxValueText : minValueText;

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

                        for (let i = 0; i < params.length; i++) {
                            const id = params[i].seriesId;
                            const name = self.itemsMap[id] && self.nameField && self.itemsMap[id][self.nameField] ? self.getItemName(self.itemsMap[id][self.nameField]) : id;

                            if (params.length === 1 || params[i].data !== 0) {
                                const value = self.getDisplayCurrency(params[i].data, self.defaultCurrency);
                                tooltip += '<div><span class="chart-pointer" style="background-color: ' + params[i].color + '"></span>';
                                tooltip += `<span>${name}</span><span style="margin-left: 20px; float: right">${value}</span><br/>`;
                                tooltip += '</div>';
                                totalAmount += params[i].data;
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
                        data: self.allDisplayMonths
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
            const yearMonthTime = this.allYearMonthTimes[e.dataIndex];

            this.$emit('click', {
                yearMonth: getYearMonthStringFromObject(yearMonthTime),
                item: item
            });
        },
        getColor: function (color) {
            if (color && color !== colorConstants.defaultColor) {
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
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
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
