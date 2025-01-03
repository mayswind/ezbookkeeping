<template>
    <v-card :class="{ 'disabled': disabled }">
        <template #title>
            <span>{{ $t('Income and Expense Trends') }}</span>
        </template>

        <v-card-text class="overview-monthly-chart-container overview-monthly-chart-overlay" v-if="loading && !hasAnyData">
            <div class="overview-monthly-chart-skeleton-container h-100" style="margin-top: -30px">
                <div class="d-flex w-100 h-100 align-center justify-center"
                     :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ]">
                    <v-skeleton-loader width="16" height="200" :loading="true"></v-skeleton-loader>
                </div>
            </div>
        </v-card-text>

        <v-card-text class="overview-monthly-chart-container overview-monthly-chart-overlay" v-else-if="!loading && !hasAnyData">
            <div class="d-flex flex-column align-center justify-center w-100 h-100">
                <h2 style="margin-top: -40px">{{ $t('No data') }}</h2>
            </div>
        </v-card-text>

        <v-chart autoresize class="overview-monthly-chart-container"
                 :class="{ 'readonly': !hasAnyData }" :option="chartOptions" @click="clickItem"/>
    </v-card>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import { TransactionType } from '@/core/transaction.ts';
import {
    parseDateFromUnixTime,
    getMonthName
} from '@/lib/datetime.ts';
import { getExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export default {
    props: [
        'loading',
        'data',
        'disabled',
        'isDarkMode',
        'enableClickItem'
    ],
    emits: [
        'click'
    ],
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        showAmountInHomePage() {
            return this.settingsStore.appSettings.showAmountInHomePage;
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        hasAnyData() {
            if (!this.data || !this.data.length || this.data.length < 1) {
                return false;
            }

            for (let i = 0; i < this.data.length; i++) {
                const item = this.data[i];

                if (item.incomeAmount > 0 || item.incomeAmount < 0 || item.expenseAmount > 0 || item.expenseAmount < 0) {
                    return true;
                }
            }

            return false;
        },
        chartOptions() {
            const self = this;
            const monthNames = [];
            const incomeAmounts = [];
            const expenseAmounts = [];
            let minAmount = 0;
            let maxAmount = 0;

            const expenseIncomeAmountColor = getExpenseAndIncomeAmountColor(this.userStore.currentUserExpenseAmountColor, this.userStore.currentUserIncomeAmountColor, this.isDarkMode);

            if (self.data) {
                for (let i = 0; i < self.data.length; i++) {
                    const item = self.data[i];
                    const month = getMonthName(parseDateFromUnixTime(item.monthStartTime));

                    monthNames.push(self.$locale.getMonthShortName(month));
                    incomeAmounts.push(item.incomeAmount);
                    expenseAmounts.push(-item.expenseAmount);

                    if (item.incomeAmount > maxAmount) {
                        maxAmount = item.incomeAmount;
                    }

                    if (-item.expenseAmount > maxAmount) {
                        maxAmount = -item.expenseAmount;
                    }

                    if (item.incomeAmount < minAmount) {
                        minAmount = item.incomeAmount;
                    }

                    if (-item.expenseAmount < minAmount) {
                        minAmount = -item.expenseAmount;
                    }
                }
            }

            const amountGap = maxAmount - minAmount;

            return {
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    },
                    backgroundColor: self.isDarkMode ? '#333' : '#fff',
                    borderColor: self.isDarkMode ? '#333' : '#fff',
                    textStyle: {
                        color: self.isDarkMode ? '#eee' : '#333'
                    },
                    formatter: params => {
                        let incomeAmount = 0;
                        let expenseAmount = 0;

                        for (let i = 0; i < params.length; i++) {
                            const param = params[i];
                            const dataIndex = param.dataIndex;
                            const data = self.data[dataIndex];

                            if (param.seriesId === 'seriesIncome') {
                                incomeAmount = self.getDisplayIncomeAmount(data);
                            } else if (param.seriesId === 'seriesExpense') {
                                expenseAmount = self.getDisplayExpenseAmount(data);
                            }
                        }

                        return `<table>` +
                            `<thead>` +
                                `<tr>` +
                                    `<td colspan="2" class="text-left">${params[0].name}</td>` +
                                `</tr>` +
                            `</thead>` +
                            `<tbody>` +
                                `<tr>` +
                                    `<td><span class="overview-monthly-chart-tooltip-indicator bg-income mr-1"></span><span class="mr-4">${self.$t('Income')}</span></td>` +
                                    `<td><strong>${incomeAmount}</strong></td>` +
                                `</tr>` +
                                `<tr>` +
                                    `<td><span class="overview-monthly-chart-tooltip-indicator bg-expense mr-1"></span><span class="mr-4">${self.$t('Expense')}</span></td>` +
                                    `<td><strong>${expenseAmount}</strong></td>` +
                                `</tr>` +
                            `</tbody>` +
                            `</table>`;
                    }
                },
                legend: {
                    bottom: 20,
                    itemWidth: 14,
                    itemHeight: 14,
                    textStyle: {
                        color: self.isDarkMode ? '#eee' : '#333'
                    },
                    icon: 'circle',
                    data: [ self.$t('Income'), self.$t('Expense') ]
                },
                grid: {
                    left: '20px',
                    right: '20px',
                    top: '10px',
                    bottom: '100px'
                },
                xAxis: [
                    {
                        type: 'category',
                        data: monthNames,
                        axisLine: {
                            show: false
                        },
                        axisTick: {
                            show: false
                        },
                        axisLabel: {
                            padding: [ 20, 0, 0, 0 ]
                        }
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        min: minAmount - amountGap / 20,
                        max: maxAmount,
                        splitNumber: 10,
                        axisLabel: {
                            show: false
                        },
                        splitLine: {
                            show: false
                        }
                    },
                    {
                        type: 'value',
                        min: minAmount,
                        max: maxAmount + amountGap / 20,
                        splitNumber: 10,
                        axisLabel: {
                            show: false
                        },
                        splitLine: {
                            show: false
                        }
                    }
                ],
                series: [
                    {
                        type: 'bar',
                        id: 'seriesIncome',
                        name: self.$t('Income'),
                        yAxisIndex: 0,
                        stack: 'Total',
                        itemStyle: {
                            color: expenseIncomeAmountColor.incomeAmountColor,
                            borderRadius: 16
                        },
                        emphasis: {
                            focus: 'series',
                            labelLine: {
                                show: false
                            }
                        },
                        barMaxWidth: 16,
                        data: incomeAmounts
                    },
                    {
                        type: 'bar',
                        id: 'seriesExpense',
                        name: self.$t('Expense'),
                        yAxisIndex: 1,
                        stack: 'Total',
                        itemStyle: {
                            color: expenseIncomeAmountColor.expenseAmountColor,
                            borderRadius: 16
                        },
                        emphasis: {
                            focus: 'series',
                            labelLine: {
                                show: false
                            }
                        },
                        barMaxWidth: 16,
                        data: expenseAmounts
                    }
                ]
            };
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
            if (!this.enableClickItem  || !this.data || e.componentType !== 'series') {
                return;
            }

            const clickData = this.data[e.dataIndex];

            if (clickData && e.seriesId === 'seriesIncome') {
                this.$emit('click', {
                    transactionType: TransactionType.Income,
                    monthStartTime: clickData.monthStartTime
                });
            } else if (clickData && e.seriesId === 'seriesExpense') {
                this.$emit('click', {
                    transactionType: TransactionType.Expense,
                    monthStartTime: clickData.monthStartTime
                });
            }
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        getDisplayAmount(amount, incomplete) {
            if (!this.showAmountInHomePage) {
                return this.getDisplayCurrency('***', this.defaultCurrency);
            }

            return this.getDisplayCurrency(amount, this.defaultCurrency) + (incomplete ? '+' : '');
        },
        getDisplayIncomeAmount(category) {
            return this.getDisplayAmount(category.incomeAmount, category.incompleteIncomeAmount);
        },
        getDisplayExpenseAmount(category) {
            return this.getDisplayAmount(category.expenseAmount, category.incompleteExpenseAmount);
        }
    }
}
</script>

<style>
.overview-monthly-chart-container {
    width: 100%;
    height: 400px;
}

.overview-monthly-chart-overlay {
    position: absolute !important;
    z-index: 10;
}

.overview-monthly-chart-skeleton-container {
    display: grid;
    grid-template-columns: repeat(12, minmax(0, 1fr));
}

.overview-monthly-chart-tooltip-indicator {
    display: inline-block;
    width: 10px;
    height: 10px;
    border-radius: 10px;
}
</style>
