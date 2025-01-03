<template>
    <f7-list class="statistics-list-item skeleton-text" v-if="loading">
        <f7-list-item link="#" :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ]">
            <template #media>
                <div class="display-flex no-padding-horizontal">
                    <div class="display-flex align-items-center statistics-icon">
                        <f7-icon f7="app_fill"></f7-icon>
                    </div>
                </div>
            </template>
            <template #title>
                <div class="statistics-list-item-text">
                    <span>Date Range</span>
                </div>
            </template>
            <template #after>
                <span>0.00 USD</span>
            </template>
            <template #inner-end>
                <div class="statistics-item-end">
                    <div class="statistics-percent-line">
                        <f7-progressbar></f7-progressbar>
                    </div>
                </div>
            </template>
        </f7-list-item>
    </f7-list>

    <f7-list v-else-if="!loading && (!allDisplayDataItems || !allDisplayDataItems.data || !allDisplayDataItems.data.length)">
        <f7-list-item :title="$t('No transaction data')"></f7-list-item>
    </f7-list>

    <f7-list v-else-if="!loading && allDisplayDataItems && allDisplayDataItems.data && allDisplayDataItems.data.length">
        <f7-list-item v-if="allDisplayDataItems.legends && allDisplayDataItems.legends.length > 1">
            <div class="display-flex" style="flex-wrap: wrap">
                <div class="trends-bar-chart-legend display-flex align-items-center"
                     :class="{ 'trends-bar-chart-legend-unselected': !!unselectedLegends[legend.id] }"
                     :key="idx"
                     v-for="(legend, idx) in allDisplayDataItems.legends"
                     @click="toggleLegend(legend)">
                    <f7-icon f7="app_fill" class="trends-bar-chart-legend-icon" :style="{ 'color': unselectedLegends[legend.id] ? '' : legend.color }"></f7-icon>
                    <span class="trends-bar-chart-legend-text">{{ legend.name }}</span>
                </div>
            </div>
        </f7-list-item>
        <f7-list-item class="statistics-list-item"
                      link="#"
                      :key="idx"
                      v-for="(item, idx) in allDisplayDataItems.data"
                      v-show="!item.hidden"
                      @click="clickItem(item)"
        >
            <template #media>
                <div class="display-flex no-padding-horizontal">
                    <div class="display-flex align-items-center statistics-icon">
                        <f7-icon f7="calendar"></f7-icon>
                    </div>
                </div>
            </template>

            <template #title>
                <div class="statistics-list-item-text">
                    <span>{{ item.displayDateRange }}</span>
                </div>
            </template>

            <template #after>
                <span>{{ getDisplayCurrency(item.totalAmount, defaultCurrency) }}</span>
            </template>

            <template #inner-end>
                <div class="statistics-item-end">
                    <div class="statistics-percent-line statistics-multi-percent-line display-flex">
                        <div class="display-inline-flex" :style="{ 'width': (item.percent * data.totalAmount / item.totalPositiveAmount) + '%' }"
                             :key="dataIdx"
                             v-for="(data, dataIdx) in item.items"
                             v-show="data.totalAmount > 0">
                            <f7-progressbar :progress="100" :style="{ '--f7-progressbar-progress-color': (data.color ? data.color : '') } "></f7-progressbar>
                        </div>
                        <div class="display-inline-flex" :style="{ 'width': (100.0 - item.percent) + '%' }"
                             v-if="item.percent < 100.0">
                            <f7-progressbar :progress="0"></f7-progressbar>
                        </div>
                    </div>
                </div>
            </template>
        </f7-list-item>
    </f7-list>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import { DateRangeScene } from '@/core/datetime.ts';
import { DEFAULT_ICON_COLOR, DEFAULT_CHART_COLORS } from '@/consts/color.ts';
import statisticsConstants from '@/consts/statistics.js';
import { isNumber } from '@/lib/common.ts';
import {
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange
} from '@/lib/datetime.ts';
import {
    sortStatisticsItems,
    getAllDateRanges
} from '@/lib/statistics.js';

export default {
    props: [
        'loading',
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
        'enableClickItem'
    ],
    emits: [
        'click'
    ],
    data() {
        return {
            unselectedLegends: {}
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        allDateRanges: function () {
            return getAllDateRanges(this.items, this.startYearMonth, this.endYearMonth, this.dateAggregationType);
        },
        allDisplayDataItems: function () {
            const allDateRangeItemsMap = {};
            const legends = [];

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (!this.hiddenField || item[this.hiddenField]) {
                    continue;
                }

                const id = (this.idField && item[this.idField]) ? item[this.idField] : this.getItemName(item[this.nameField]);

                const legend = {
                    id: id,
                    name: (this.nameField && item[this.nameField]) ? this.getItemName(item[this.nameField]) : id,
                    color: this.getColor(item[this.colorField] ? item[this.colorField] : DEFAULT_CHART_COLORS[i % DEFAULT_CHART_COLORS.length]),
                    displayOrders: (this.displayOrdersField && item[this.displayOrdersField]) ? item[this.displayOrdersField] : [0]
                };

                legends.push(legend);

                if (this.unselectedLegends[id]) {
                    continue;
                }

                const dateRangeItemMap = {};

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

                    if (dateRangeItemMap[dateRangeKey]) {
                        dateRangeItemMap[dateRangeKey].totalAmount += (this.valueField && isNumber(dataItem[this.valueField])) ? dataItem[this.valueField] : 0;
                    } else {
                        const allDataItems = allDateRangeItemsMap[dateRangeKey] || [];
                        const finalDataItem = Object.assign({}, legend, {
                            totalAmount: (this.valueField && isNumber(dataItem[this.valueField])) ? dataItem[this.valueField] : 0
                        });

                        allDataItems.push(finalDataItem);
                        dateRangeItemMap[dateRangeKey] = finalDataItem;
                        allDateRangeItemsMap[dateRangeKey] = allDataItems;
                    }
                }
            }

            const finalDataItems = [];
            let maxTotalAmount = 0;

            for (let i = 0; i < this.allDateRanges.length; i++) {
                const dateRange = this.allDateRanges[i];
                let dateRangeKey = '';

                if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Year.type) {
                    dateRangeKey = dateRange.year;
                } else if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Quarter.type) {
                    dateRangeKey = `${dateRange.year}-${dateRange.quarter}`;
                } else { // if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Month.type) {
                    dateRangeKey = `${dateRange.year}-${dateRange.month + 1}`;
                }

                let displayDateRange = '';

                if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Year.type) {
                    displayDateRange = this.$locale.formatUnixTimeToShortYear(this.userStore, dateRange.minUnixTime);
                } else if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Quarter.type) {
                    displayDateRange = this.$locale.formatYearQuarter(dateRange.year, dateRange.quarter);
                } else { // if (this.dateAggregationType === statisticsConstants.allDateAggregationTypes.Month.type) {
                    displayDateRange = this.$locale.formatUnixTimeToShortYearMonth(this.userStore, dateRange.minUnixTime);
                }

                const dataItems = allDateRangeItemsMap[dateRangeKey] || [];
                let totalAmount = 0;
                let totalPositiveAmount = 0;

                sortStatisticsItems(dataItems, this.sortingType);

                for (let j = 0; j < dataItems.length; j++) {
                    if (dataItems[j].totalAmount > 0) {
                        totalPositiveAmount += dataItems[j].totalAmount;
                    }

                    totalAmount += dataItems[j].totalAmount;
                }

                if (totalAmount > maxTotalAmount) {
                    maxTotalAmount = totalAmount;
                }

                finalDataItems.push({
                    dateRange: dateRange,
                    displayDateRange: displayDateRange,
                    items: dataItems,
                    totalAmount: totalAmount,
                    totalPositiveAmount: totalPositiveAmount
                });
            }

            for (let i = 0; i < finalDataItems.length; i++) {
                if (maxTotalAmount > 0 && finalDataItems[i].totalAmount > 0) {
                    finalDataItems[i].percent = 100.0 * finalDataItems[i].totalAmount / maxTotalAmount;
                } else {
                    finalDataItems[i].percent = 0.0;
                }
            }

            return {
                data: finalDataItems,
                legends: legends
            };
        }
    },
    methods: {
        clickItem: function (item) {
            let itemId = '';

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (!this.hiddenField || item[this.hiddenField]) {
                    continue;
                }

                const id = (this.idField && item[this.idField]) ? item[this.idField] : this.getItemName(item[this.nameField]);

                if (this.unselectedLegends[id]) {
                    continue;
                }

                if (itemId.length) {
                    itemId += ',';
                }

                itemId += id;
            }

            const dateRange = item.dateRange;
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

            const dateRangeType = getDateTypeByDateRange(minUnixTime, maxUnixTime, this.userStore.currentUserFirstDayOfWeek, DateRangeScene.Normal);

            this.$emit('click', {
                itemId: itemId,
                dateRange: {
                    minTime: minUnixTime,
                    maxTime: maxUnixTime,
                    type: dateRangeType
                }
            });
        },
        toggleLegend(legend) {
            if (this.unselectedLegends[legend.id]) {
                delete this.unselectedLegends[legend.id];
            } else {
                this.unselectedLegends[legend.id] = true;
            }
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
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        }
    }
}
</script>

<style>
.trends-bar-chart-legend {
    margin-right: 4px;
    cursor: pointer;
}

.trends-bar-chart-legend-icon.f7-icons {
    font-size: var(--ebk-trends-bar-chart-legend-icon-font-size);
    margin-right: 2px;
}

.trends-bar-chart-legend-unselected .trends-bar-chart-legend-icon.f7-icons {
    color: #cccccc;
}

.trends-bar-chart-legend-text {
    font-size: var(--ebk-trends-bar-chart-legend-text-font-size);
}

.trends-bar-chart-legend-unselected .trends-bar-chart-legend-text {
    color: #cccccc;
}
</style>
