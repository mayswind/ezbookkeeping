<template>
    <v-chart autoresize class="account-category-sankey-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import {
    type TransactionCategoricalOverviewAnalysisDataItem,
    type TransactionCategoricalOverviewAnalysisDataItemOutflowItem,
    TransactionCategoricalOverviewAnalysisDataItemType
} from '@/models/transaction.ts';

import { values } from '@/core/base.ts';
import { ThemeType } from '@/core/theme.ts';

import { isNumber } from '@/lib/common.ts';
import { getExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

enum SankeyChartDepth {
    PrimaryIncomeCategory = 0,
    SecondaryIncomeCategory = 1,
    AccountForIncome = 2,
    AccountForExpense = 3,
    SecondaryExpenseCategory = 4,
    PrimaryExpenseCategory = 5
}

enum SankeyChartNodeItemType {
    Account = 'account',
    Category = 'category',
    NetCashFlow = 'netCashFlow'
}

interface SankeyChartData {
    nodes: SankeyChartNodeItem[];
    links: SankeyChartLinkItem[];
}

interface SankeyChartNodeItem {
    dateItemType: TransactionCategoricalOverviewAnalysisDataItemType;
    itemType: SankeyChartNodeItemType;
    itemId: string;
    name: string;
    displayName: string;
    totalAmount: number;
    accountNetCashFlow?: number;
    percent?: number;
    depth: number;
    itemStyle?: {
        color?: string;
        opacity?: number;
    }
}

interface SankeyChartLinkItem {
    sourceItemType: SankeyChartNodeItemType;
    sourceItemId: string;
    source: string;
    sourceDisplayName: string;
    targetItemType: SankeyChartNodeItemType;
    targetItemId: string;
    target: string;
    targetDisplayName: string;
    value: number;
}

const props = defineProps<{
    skeleton?: boolean;
    items: TransactionCategoricalOverviewAnalysisDataItem[];
    defaultCurrency?: string;
    enableClickItem?: boolean;
}>();

const emit = defineEmits<{
    (e: 'click', sourceItemType: 'account' | 'category', sourceItemId: string, targetItemType?: 'account' | 'category', targetItemId?: string): void;
}>();

const theme = useTheme();

const {
    tt,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatPercentToLocalizedNumerals
} = useI18n();

const userStore = useUserStore();

const overviewDataItemTypeSankeyChartNodeItemTypeMap: Record<TransactionCategoricalOverviewAnalysisDataItemType, SankeyChartNodeItemType> = {
    [TransactionCategoricalOverviewAnalysisDataItemType.IncomeByPrimaryCategory]: SankeyChartNodeItemType.Category,
    [TransactionCategoricalOverviewAnalysisDataItemType.IncomeBySecondaryCategory]: SankeyChartNodeItemType.Category,
    [TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount]: SankeyChartNodeItemType.Account,
    [TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount]: SankeyChartNodeItemType.Account,
    [TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow]: SankeyChartNodeItemType.NetCashFlow,
    [TransactionCategoricalOverviewAnalysisDataItemType.ExpenseBySecondaryCategory]: SankeyChartNodeItemType.Category,
    [TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByPrimaryCategory]: SankeyChartNodeItemType.Category
};

const overviewDataItemTypeSankeyChartNodeItemDepthMap: Record<TransactionCategoricalOverviewAnalysisDataItemType, number> = {
    [TransactionCategoricalOverviewAnalysisDataItemType.IncomeByPrimaryCategory]: SankeyChartDepth.PrimaryIncomeCategory,
    [TransactionCategoricalOverviewAnalysisDataItemType.IncomeBySecondaryCategory]: SankeyChartDepth.SecondaryIncomeCategory,
    [TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount]: SankeyChartDepth.AccountForIncome,
    [TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount]: SankeyChartDepth.AccountForExpense,
    [TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow]: SankeyChartDepth.SecondaryExpenseCategory,
    [TransactionCategoricalOverviewAnalysisDataItemType.ExpenseBySecondaryCategory]: SankeyChartDepth.SecondaryExpenseCategory,
    [TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByPrimaryCategory]: SankeyChartDepth.PrimaryExpenseCategory
};

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const sankeyData = computed<SankeyChartData>(() => {
    const nodes: SankeyChartNodeItem[] = [];
    const links: SankeyChartLinkItem[] = [];

    for (const item of props.items) {
        if (item.hidden) {
            continue;
        }

        const itemType = overviewDataItemTypeSankeyChartNodeItemTypeMap[item.type];
        const depth = overviewDataItemTypeSankeyChartNodeItemDepthMap[item.type];

        if (!itemType || itemType === SankeyChartNodeItemType.NetCashFlow || depth === undefined) {
            continue;
        }

        if (item.totalAmount === 0 && item.outflows.length === 0) {
            continue;
        }

        const nodeItem: SankeyChartNodeItem = {
            dateItemType: item.type,
            itemType: itemType,
            itemId: item.id,
            name: `${item.type}:${item.id}`,
            displayName: item.name,
            totalAmount: item.totalAmount,
            percent: item.percent,
            depth: depth
        };

        if (!isNumber(nodeItem.percent) && nodeItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount) {
            nodeItem.itemStyle = {
                color: '#aaa',
                opacity: 0.5
            };
        }

        if (nodeItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount) {
            for (const outflowItem of item.outflows) {
                if (outflowItem.relatedItem.type !== TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow) {
                    continue;
                }

                nodeItem.accountNetCashFlow = (nodeItem.accountNetCashFlow ?? 0) + outflowItem.amount;
            }
        }

        nodes.push(nodeItem);

        const combinedOutflows: Record<string, TransactionCategoricalOverviewAnalysisDataItemOutflowItem> = {};

        for (const outflowItem of item.outflows) {
            const relatedItem = outflowItem.relatedItem;

            if (!relatedItem) {
                continue;
            }

            if (outflowItem.relatedItem) {
                const key = `${item.type}:${item.id}-${outflowItem.relatedItem.type}:${outflowItem.relatedItem.id}`;
                let combinedOutflow: TransactionCategoricalOverviewAnalysisDataItemOutflowItem | undefined = combinedOutflows[key];

                if (!combinedOutflow) {
                    combinedOutflow = {
                        relatedItem: outflowItem.relatedItem,
                        amount: 0
                    };
                    combinedOutflows[key] = combinedOutflow;
                }

                combinedOutflow.amount += outflowItem.amount;
            }
        }

        for (const outflowItem of values(combinedOutflows)) {
            const relatedItem = outflowItem.relatedItem;
            const transferItemType = overviewDataItemTypeSankeyChartNodeItemTypeMap[relatedItem.type];

            if (!transferItemType) {
                continue;
            }

            const linkItem: SankeyChartLinkItem = {
                sourceItemType: itemType,
                sourceItemId: item.id,
                source: `${item.type}:${item.id}`,
                sourceDisplayName: item.name,
                targetItemType: transferItemType,
                targetItemId: relatedItem.id,
                target: `${relatedItem.type}:${relatedItem.id}`,
                targetDisplayName: relatedItem.name,
                value: outflowItem.amount
            };

            links.push(linkItem);
        }
    }

    const ret: SankeyChartData = {
        nodes: nodes,
        links: links
    };

    return ret;
});

const chartOptions = computed<object>(() => {
    const expenseIncomeAmountColor = getExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor, isDarkMode.value);

    return {
        tooltip: {
            trigger: 'item',
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams) => {
                if (params.dataType === 'node') {
                    const dataItem = params.data as SankeyChartNodeItem;
                    const value = dataItem.totalAmount;
                    const displayValue = formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
                    let displayTypeName = '';

                    if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.IncomeByPrimaryCategory) {
                        displayTypeName = tt('Income By Primary Category');
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.IncomeBySecondaryCategory) {
                        displayTypeName = tt('Income By Secondary Category');
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount) {
                        displayTypeName = tt('Income By Account');
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByAccount) {
                        displayTypeName = tt('Expense By Account');
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.NetCashFlow) {
                        displayTypeName = tt('Net Cash Flow');
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.ExpenseBySecondaryCategory) {
                        displayTypeName = tt('Expense By Secondary Category');
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.ExpenseByPrimaryCategory) {
                        displayTypeName = tt('Expense By Primary Category');
                    }

                    let tooltip = `<div><span>${dataItem.displayName}</span>`;

                    if (displayTypeName && (dataItem.dateItemType !== TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount || isNumber(dataItem.percent))) {
                        tooltip = `<div class="mb-1">${displayTypeName}</div>` + tooltip;
                    } else if (dataItem.dateItemType === TransactionCategoricalOverviewAnalysisDataItemType.IncomeByAccount) {
                        tooltip = `<div class="mb-1">${tt('Account Balance')}</div>` + tooltip;
                    }

                    if (isNumber(dataItem.percent) && dataItem.percent > 0) {
                        const displayPercent = formatPercentToLocalizedNumerals(dataItem.percent, 2, '<0.01');
                        tooltip += `<span class="ms-1" style="float: inline-end">(${displayPercent})</span>`;
                    }

                    tooltip += `<span class="ms-5" style="float: inline-end">${displayValue}</span></div>`;

                    if (isNumber(dataItem.accountNetCashFlow) && dataItem.accountNetCashFlow !== 0) {
                        const displayAccountNetCashFlow = formatAmountToLocalizedNumeralsWithCurrency(dataItem.accountNetCashFlow, props.defaultCurrency);
                        tooltip += `<div class="mt-1"><span>${tt('Net Cash Flow')}</span><span class="ms-5" style="float: inline-end">${displayAccountNetCashFlow}</span></div>`;
                    }

                    return tooltip;
                } else if (params.dataType === 'edge') {
                    const dataItem = params.data as SankeyChartLinkItem;
                    const value = isNumber(params.value) ? params.value as number : 0;
                    const displayValue = formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
                    return `<div><span>${dataItem.sourceDisplayName} → ${dataItem.targetDisplayName}</span><span class="ms-5" style="float: inline-end">${displayValue}</span></div>`;
                } else {
                    return '';
                }
            }
        },
        series: [
            {
                type: 'sankey',
                left: 10,
                top: 0,
                bottom: 10,
                roam: true,
                layoutIterations: 0,
                label: {
                    formatter: (params: CallbackDataParams) => {
                        const dataItem = params.data as SankeyChartNodeItem;
                        return dataItem.displayName;
                    }
                },
                levels: [
                    {
                        depth: SankeyChartDepth.PrimaryIncomeCategory,
                        itemStyle: {
                            color: expenseIncomeAmountColor.incomeAmountColor,
                            opacity: 0.6
                        },
                        lineStyle: {
                            color: 'source',
                            opacity: 0.3
                        }
                    },
                    {
                        depth: SankeyChartDepth.SecondaryIncomeCategory,
                        itemStyle: {
                            color: expenseIncomeAmountColor.incomeAmountColor,
                            opacity: 0.4
                        },
                        lineStyle: {
                            color: 'source',
                            opacity: 0.2
                        }
                    },
                    {
                        depth: SankeyChartDepth.AccountForIncome,
                        itemStyle: {
                            color: '#c07d43',
                            opacity: 0.5
                        },
                        lineStyle: {
                            color: 'source',
                            opacity: 0.2
                        }
                    },
                    {
                        depth: SankeyChartDepth.AccountForExpense,
                        itemStyle: {
                            color: '#c07d43',
                            opacity: 0.5
                        },
                        lineStyle: {
                            color: 'source',
                            opacity: 0.2
                        }
                    },
                    {
                        depth: SankeyChartDepth.SecondaryExpenseCategory,
                        itemStyle: {
                            color: expenseIncomeAmountColor.expenseAmountColor,
                            opacity: 0.4
                        },
                        lineStyle: {
                            color: 'source',
                            opacity: 0.2
                        }
                    },
                    {
                        depth: SankeyChartDepth.PrimaryExpenseCategory,
                        itemStyle: {
                            color: expenseIncomeAmountColor.expenseAmountColor,
                            opacity: 0.6
                        },
                        lineStyle: {
                            color: 'source',
                            opacity: 0.3
                        }
                    }
                ],
                emphasis: {
                    focus: 'adjacency'
                },
                data: sankeyData.value.nodes,
                links: sankeyData.value.links,
                animation: !props.skeleton
            }
        ]
    };
});

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series' || e.seriesType !=='sankey') {
        return;
    }

    if (!e.data) {
        return;
    }

    if (e.dataType === 'node') {
        const dataItem = e.data as SankeyChartNodeItem;

        if (dataItem.itemType === SankeyChartNodeItemType.NetCashFlow) {
            return;
        }

        emit('click', dataItem.itemType, dataItem.itemId);
    } else if (e.dataType === 'edge') {
        const dataItem = e.data as SankeyChartLinkItem;

        if (dataItem.sourceItemType === SankeyChartNodeItemType.NetCashFlow) {
            return;
        }

        if (dataItem.sourceItemType === dataItem.targetItemType && dataItem.sourceItemId === dataItem.targetItemId) {
            emit('click', dataItem.sourceItemType, dataItem.sourceItemId);
        } else if (dataItem.targetItemType !== SankeyChartNodeItemType.NetCashFlow) {
            emit('click', dataItem.sourceItemType, dataItem.sourceItemId, dataItem.targetItemType, dataItem.targetItemId);
        }
    }
}
</script>

<style scoped>
.account-category-sankey-chart-container {
    width: 100%;
    height: 460px;
}

@media (min-width: 600px) {
    .account-category-sankey-chart-container {
        height: 660px;
    }
}

.account-category-sankey-chart-container.transition-in {
    animation: radar-chart-skeleton-fade-in 2s 1;
}

@keyframes radar-chart-skeleton-fade-in {
    0% {
        opacity: 0;
    }
    20% {
        opacity: 0;
    }
    100% {
        opacity: 1;
    }
}
</style>
