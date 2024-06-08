<template>
    <v-chart autoresize class="trends-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" />
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
        'currencyField',
        'colorField',
        'hiddenField',
        'defaultCurrency',
        'showValue',
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
                    id = item[this.nameField];
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
                    id: (this.idField && item[this.idField]) ? item[this.idField] : item[this.nameField],
                    name: (this.idField && item[this.idField]) ? item[this.idField] : item[this.nameField],
                    currency: item[this.currencyField],
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

                        if (params.length && params[0].name) {
                            tooltip += `${params[0].name}<br/>`;
                        }

                        for (let i = 0; i < params.length; i++) {
                            const id = params[i].seriesId;
                            const name = self.itemsMap[id] && self.nameField && self.itemsMap[id][self.nameField] ? self.itemsMap[id][self.nameField] : id;

                            if (params[i].data !== 0) {
                                const currency = self.itemsMap[id] && self.currencyField && self.itemsMap[id][self.currencyField] ? self.itemsMap[id][self.currencyField] : self.defaultCurrency;
                                const value = self.getDisplayCurrency(params[i].data, currency);
                                tooltip += '<div><span class="chart-pointer" style="background-color: ' + params[i].color + '"></span>';
                                tooltip += `<span>${name}</span><span style="margin-left: 20px; float: right">${value}</span><br/>`;
                                tooltip += '</div>';
                            }
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
                        return self.itemsMap[id] && self.nameField && self.itemsMap[id][self.nameField] ? self.itemsMap[id][self.nameField] : id;
                    }
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
    height: 500px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .pie-chart-container {
        height: 500px;
    }
}
</style>
