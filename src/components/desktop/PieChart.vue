<template>
    <v-chart class="pie-chart-container" autoresize :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';

import colorConstants from '@/consts/color.js';
import { formatPercent } from '@/lib/common.js';

export default {
    props: [
        'items',
        'idField',
        'nameField',
        'valueField',
        'percentField',
        'currencyField',
        'colorField',
        'hiddenField',
        'minValidPercent',
        'defaultCurrency',
        'showValue',
        'enableClickItem'
    ],
    emits: [
        'click'
    ],
    data() {
        return {
            selectedLegends: null,
            selectedIndex: 0
        };
    },
    computed: {
        ...mapStores(useSettingsStore),
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
        validItems: function () {
            let totalValidValue = 0;

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (item[this.valueField] && item[this.valueField] > 0 && (!this.hiddenField || !item[this.hiddenField])) {
                    totalValidValue += item[this.valueField];
                }
            }

            const validItems = [];

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (item[this.valueField] && item[this.valueField] > 0 &&
                    (!this.hiddenField || !item[this.hiddenField]) &&
                    (!this.minValidPercent || item[this.valueField] / totalValidValue > this.minValidPercent)) {
                    const finalItem = {
                        id: (this.idField && item[this.idField]) ? item[this.idField] : item[this.nameField],
                        name: (this.idField && item[this.idField]) ? item[this.idField] : item[this.nameField],
                        displayName: item[this.nameField],
                        value: item[this.valueField],
                        percent: (item[this.percentField] > 0 || item[this.percentField] === 0 || item[this.percentField] === '0') ? item[this.percentField] : (item[this.valueField] / totalValidValue * 100),
                        actualPercent: item[this.valueField] / totalValidValue,
                        currency: item[this.currencyField],
                        itemStyle: {
                            color: this.getColor(item[this.colorField] ? item[this.colorField] : colorConstants.defaultChartColors[validItems.length % colorConstants.defaultChartColors.length]),
                        },
                        selected: true,
                        sourceItem: item
                    };

                    finalItem.displayPercent = formatPercent(finalItem.percent, 2, '&lt;0.01');
                    finalItem.displayValue = this.getDisplayCurrency(finalItem.value, (finalItem.currency || this.defaultCurrency));

                    validItems.push(finalItem);
                }
            }

            return validItems;
        },
        hasUnselectedItem: function () {
            for (let i = 0; i < this.validItems.length; i++) {
                const item = this.validItems[i];

                if (this.selectedLegends && !this.selectedLegends[item.id]) {
                    return true;
                }
            }

            return false;
        },
        firstItemAndHalfCurrentItemTotalPercent: function () {
            let totalValue = 0;
            let firstValue = null;
            let firstToCurrentTotalValue = 0;

            for (let i = 0; i < this.validItems.length; i++) {
                const item = this.validItems[i];

                if (this.selectedLegends && !this.selectedLegends[item.id]) {
                    continue;
                }

                if (firstValue === null) {
                    firstValue = item.value;
                }

                if (firstValue !== null) {
                    if (i < this.selectedIndex) {
                        firstToCurrentTotalValue += item.value;
                    } else if (i === this.selectedIndex) {
                        firstToCurrentTotalValue += item.value / 2;
                    }
                }

                totalValue += item.value;
            }

            if (firstToCurrentTotalValue && totalValue > 0) {
                return firstToCurrentTotalValue / totalValue;
            } else {
                return 0;
            }
        },
        chartOptions: function () {
            const self = this;

            return {
                tooltip: {
                    trigger: 'item',
                    backgroundColor: self.isDarkMode ? '#333' : '#fff',
                    borderColor: self.isDarkMode ? '#333' : '#fff',
                    textStyle: {
                        color: self.isDarkMode ? '#eee' : '#333'
                    },
                    formatter: params => {
                        const name = params.data ? params.data.displayName : '';
                        const value = params.data ? params.data.displayValue : self.getDisplayCurrency(params.value);
                        let percent = params.data ? params.data.displayPercent : (params.percent + '%');

                        if (self.hasUnselectedItem) {
                            percent = params.percent + '%';
                        }

                        if (name) {
                            return `${name}<br/>${value} (${percent})`;
                        } else {
                            return `${value} (${percent})`;
                        }
                    }
                },
                legend: {
                    orient: 'horizontal',
                    data: self.validItems.map(item => item.name),
                    selected: self.selectedLegends,
                    textStyle: {
                        color: self.isDarkMode ? '#eee' : '#333'
                    },
                    formatter: id => {
                        return self.itemsMap[id] && self.nameField && self.itemsMap[id][self.nameField] ? self.itemsMap[id][self.nameField] : id;
                    }
                },
                series: [
                    {
                        type: 'pie',
                        data: self.validItems,
                        top: 50,
                        startAngle: -90 + self.firstItemAndHalfCurrentItemTotalPercent * 360,
                        emphasis: {
                            itemStyle: {
                                shadowBlur: 10,
                                shadowOffsetX: 0,
                                shadowColor: 'rgba(0, 0, 0, 0.5)',
                            }
                        },
                        label: {
                            color: self.isDarkMode ? '#eee' : '#333',
                            formatter: params => {
                                return params.data ? params.data.displayName : '';
                            }
                        }
                    }
                ],
                media: [
                    {
                        query: {
                            minWidth: 600
                        },
                        option: {
                            legend: {
                                orient: 'vertical',
                                left: 'left'
                            },
                            series: [
                                {
                                    type: 'pie',
                                    top: 0
                                }
                            ]
                        }
                    }
                ]
            }
        }
    },
    watch: {
        'items': function () {
            this.selectedIndex = 0;
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
            if (!this.enableClickItem || e.componentType !== 'series' || e.seriesType !=='pie') {
                return;
            }

            if (e.event && e.event.target && e.event.target.currentStates && e.event.target.currentStates[0] && e.event.target.currentStates[0] === 'emphasis') {
                this.selectedIndex = e.dataIndex;
                return;
            }

            if (!e.data || !e.data.sourceItem) {
                return;
            }

            this.$emit('click', e.data.sourceItem);
        },
        onLegendSelectChanged: function (e) {
            this.selectedLegends = e.selected;
            const selectedItem = this.validItems[this.selectedIndex];

            if (!selectedItem || !this.selectedLegends[selectedItem.id]) {
                let newSelectedIndex = 0;

                for (let i = 0; i < this.validItems.length; i++) {
                    const item = this.validItems[i];

                    if (this.selectedLegends[item.id]) {
                        newSelectedIndex = i;
                        break;
                    }
                }

                this.selectedIndex = newSelectedIndex;
            }
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
.pie-chart-container {
    width: 100%;
    height: 400px;
}

@media (min-width: 600px) {
    .pie-chart-container {
        height: 500px;
    }
}
</style>
