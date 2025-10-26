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

import type {
    SortableTransactionStatisticDataItem,
    TransactionStatisticResponseItemWithInfo
} from '@/models/transaction.ts';

import { values } from '@/core/base.ts';
import { ThemeType } from '@/core/theme.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionRelatedAccountType } from '@/core/transaction.ts';

import { isNumber } from '@/lib/common.ts';
import { sortStatisticsItems } from '@/lib/statistics.ts';
import { getExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

enum SankeyChartNodeItemType {
    Account = 'account',
    Category = 'category'
}

interface SankeyChartData {
    nodes: SankeyChartNodeItem[];
    links: SankeyChartLinkItem[];
}

interface SankeyChartNodeItem extends SortableTransactionStatisticDataItem {
    itemType: SankeyChartNodeItemType;
    itemId: string;
    name: string;
    nameId: string;
    displayName: string;
    displayOrders: number[];
    totalAmount: number;
    depth: number;
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
    items: TransactionStatisticResponseItemWithInfo[];
    sortingType: number;
    defaultCurrency?: string;
    enableClickItem?: boolean;
}>();

const emit = defineEmits<{
    (e: 'click', sourceItemType: 'account' | 'category', sourceItemId: string, targetItemType?: 'account' | 'category', targetItemId?: string): void;
}>();

const theme = useTheme();

const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

const userStore = useUserStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const sankeyData = computed<SankeyChartData>(() => {
    const primaryIncomeCategoryNodesMap: Record<string, SankeyChartNodeItem> = {};
    const secondaryIncomeCategoryNodesMap: Record<string, SankeyChartNodeItem> = {};
    const incomeAccountNodesMap: Record<string, SankeyChartNodeItem> = {};
    const expenseAccountNodesMap: Record<string, SankeyChartNodeItem> = {};
    const secondaryExpenseCategoryNodesMap: Record<string, SankeyChartNodeItem> = {};
    const primaryExpenseCategoryNodesMap: Record<string, SankeyChartNodeItem> = {};
    const linksMap: Record<string, SankeyChartLinkItem> = {};

    for (const item of props.items) {
        if (!item.primaryAccount || !item.account || !item.primaryCategory || !item.category || !item.amountInDefaultCurrency) {
            continue;
        }

        if (item.account.hidden || item.primaryAccount.hidden || item.category.hidden || item.primaryCategory.hidden) {
            continue;
        }

        if (item.relatedAccount && (item.relatedAccountType === TransactionRelatedAccountType.TransferFrom || item.relatedAccount.hidden || !item.relatedPrimaryAccount || item.relatedPrimaryAccount.hidden)) {
            continue;
        }

        const incomeAccountNameId = `income_account:${item.account.id}`;
        const expenseAccountNameId = `expense_account:${item.account.id}`;

        updateNodeItem(incomeAccountNodesMap, {
            itemType: SankeyChartNodeItemType.Account,
            id: item.account.id,
            name: item.account.name,
            nameId: incomeAccountNameId,
            displayOrders: [item.account.displayOrder],
            amount: (item.primaryCategory.type == CategoryType.Income || item.primaryCategory.type == CategoryType.Transfer) ? item.amountInDefaultCurrency : 0,
            depth: 2
        });

        updateNodeItem(expenseAccountNodesMap, {
            itemType: SankeyChartNodeItemType.Account,
            id: item.account.id,
            name: item.account.name,
            nameId: expenseAccountNameId,
            displayOrders: [item.account.displayOrder],
            amount: item.primaryCategory.type == CategoryType.Expense ? item.amountInDefaultCurrency : 0,
            depth: 3
        });

        if (item.primaryCategory.type == CategoryType.Income) {
            updateNodeItem(primaryIncomeCategoryNodesMap, {
                itemType: SankeyChartNodeItemType.Category,
                id: item.primaryCategory.id,
                name: item.primaryCategory.name,
                nameId: item.primaryCategory.id,
                displayOrders: [item.primaryCategory.displayOrder],
                amount: item.amountInDefaultCurrency,
                depth: 0
            });

            updateNodeItem(secondaryIncomeCategoryNodesMap, {
                itemType: SankeyChartNodeItemType.Category,
                id: item.category.id,
                name: item.category.name,
                nameId: item.category.id,
                displayOrders: [item.primaryCategory.displayOrder, item.category.displayOrder],
                amount: item.amountInDefaultCurrency,
                depth: 1
            });

            updateLinkItem(linksMap, {
                sourceItemType: SankeyChartNodeItemType.Category,
                sourceItemId: item.primaryCategory.id,
                source: item.primaryCategory.id,
                sourceName: item.primaryCategory.name,
                targetItemType: SankeyChartNodeItemType.Category,
                targetItemId: item.category.id,
                target: item.category.id,
                targetName: item.category.name,
                value: item.amountInDefaultCurrency
            });

            updateLinkItem(linksMap, {
                sourceItemType: SankeyChartNodeItemType.Category,
                sourceItemId: item.category.id,
                source: item.category.id,
                sourceName: item.category.name,
                targetItemType: SankeyChartNodeItemType.Account,
                targetItemId: item.account.id,
                target: incomeAccountNameId,
                targetName: item.account.name,
                value: item.amountInDefaultCurrency
            });
        } else if (item.primaryCategory.type == CategoryType.Expense) {
            updateNodeItem(secondaryExpenseCategoryNodesMap, {
                itemType: SankeyChartNodeItemType.Category,
                id: item.category.id,
                name: item.category.name,
                nameId: item.category.id,
                displayOrders: [item.primaryCategory.displayOrder, item.category.displayOrder],
                amount: item.amountInDefaultCurrency,
                depth: 4
            });

            updateNodeItem(primaryExpenseCategoryNodesMap, {
                itemType: SankeyChartNodeItemType.Category,
                id: item.primaryCategory.id,
                name: item.primaryCategory.name,
                nameId: item.primaryCategory.id,
                displayOrders: [item.primaryCategory.displayOrder],
                amount: item.amountInDefaultCurrency,
                depth: 5
            });

            updateLinkItem(linksMap, {
                sourceItemType: SankeyChartNodeItemType.Account,
                sourceItemId: item.account.id,
                source: expenseAccountNameId,
                sourceName: item.account.name,
                targetItemType: SankeyChartNodeItemType.Category,
                targetItemId: item.category.id,
                target: item.category.id,
                targetName: item.category.name,
                value: item.amountInDefaultCurrency
            });

            updateLinkItem(linksMap, {
                sourceItemType: SankeyChartNodeItemType.Category,
                sourceItemId: item.category.id,
                source: item.category.id,
                sourceName: item.category.name,
                targetItemType: SankeyChartNodeItemType.Category,
                targetItemId: item.primaryCategory.id,
                target: item.primaryCategory.id,
                targetName: item.primaryCategory.name,
                value: item.amountInDefaultCurrency
            });
        } else if (item.primaryCategory.type == CategoryType.Transfer && item.relatedAccount) {
            const relatedAccountNameId = `expense_account:${item.relatedAccount.id}`;
            updateLinkItem(linksMap, {
                sourceItemType: SankeyChartNodeItemType.Account,
                sourceItemId: item.account.id,
                source: incomeAccountNameId,
                sourceName: item.account.name,
                targetItemType: SankeyChartNodeItemType.Account,
                targetItemId: item.relatedAccount.id,
                target: relatedAccountNameId,
                targetName: item.relatedAccount.name,
                value: item.amountInDefaultCurrency
            });
        }
    }

    for (const expenseAccountNode of values(expenseAccountNodesMap)) {
        const incomeAccountNode = incomeAccountNodesMap[`income_account:${expenseAccountNode.itemId}`];

        if (!incomeAccountNode) {
            continue;
        }

        let totalOutflowAmount = 0;
        let totalInflowAmount = 0;

        for (const link of values(linksMap)) {
            if (link.sourceItemType === SankeyChartNodeItemType.Account && link.sourceItemId === expenseAccountNode.itemId) {
                totalOutflowAmount += link.value;
            } else if (link.targetItemType === SankeyChartNodeItemType.Account && link.targetItemId === expenseAccountNode.itemId) {
                totalInflowAmount += link.value;
            }
        }

        const amountDifference = totalOutflowAmount - totalInflowAmount;

        if (amountDifference > 0) {
            updateLinkItem(linksMap, {
                sourceItemType: SankeyChartNodeItemType.Account,
                sourceItemId: incomeAccountNode.itemId,
                source: incomeAccountNode.nameId,
                sourceName: incomeAccountNode.displayName,
                targetItemType: SankeyChartNodeItemType.Account,
                targetItemId: expenseAccountNode.itemId,
                target: expenseAccountNode.nameId,
                targetName: expenseAccountNode.displayName,
                value: amountDifference
            });
        }
    }

    const nodes: SankeyChartNodeItem[] = [];
    const links: SankeyChartLinkItem[] = [];
    addFinalSortedNodeItems(primaryIncomeCategoryNodesMap, nodes);
    addFinalSortedNodeItems(secondaryIncomeCategoryNodesMap, nodes);
    addFinalSortedNodeItems(incomeAccountNodesMap, nodes);
    addFinalSortedNodeItems(expenseAccountNodesMap, nodes);
    addFinalSortedNodeItems(secondaryExpenseCategoryNodesMap, nodes);
    addFinalSortedNodeItems(primaryExpenseCategoryNodesMap, nodes);
    addFinalLinkItems(linksMap, links);

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
                const value = isNumber(params.value) ? params.value as number : 0;
                const displayValue = formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);

                if (params.dataType === 'node') {
                    const dataItem = params.data as SankeyChartNodeItem;
                    return `<div><span>${dataItem.displayName}</span><span class="ms-5" style="float: inline-end">${displayValue}</span></div>`;
                } else if (params.dataType === 'edge') {
                    const dataItem = params.data as SankeyChartLinkItem;
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
                top: 10,
                bottom: 10,
                label: {
                    formatter: (params: CallbackDataParams) => {
                        const dataItem = params.data as SankeyChartNodeItem;
                        return dataItem.displayName;
                    }
                },
                levels: [
                    {
                        depth: 0,
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
                        depth: 1,
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
                        depth: 2,
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
                        depth: 3,
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
                        depth: 4,
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
                        depth: 5,
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

function updateNodeItem(nodesMap: Record<string, SankeyChartNodeItem>, { itemType, id, name, nameId, displayOrders, amount, depth }: { itemType: SankeyChartNodeItemType, id: string, name: string, nameId: string, displayOrders: number[], amount: number, depth: number }) {
    const node: SankeyChartNodeItem | undefined = nodesMap[nameId];

    if (!node) {
        nodesMap[nameId] = {
            itemType: itemType,
            itemId: id,
            name: name,
            nameId: nameId,
            displayName: name,
            displayOrders: displayOrders,
            totalAmount: amount,
            depth: depth
        };
    } else {
        node.totalAmount += amount;
    }
}

function updateLinkItem(linksMap: Record<string, SankeyChartLinkItem>, { sourceItemType, sourceItemId, source, sourceName, targetItemType, targetItemId, target, targetName, value }: { sourceItemType: SankeyChartNodeItemType, sourceItemId: string, source: string, sourceName: string, targetItemType: SankeyChartNodeItemType, targetItemId: string, target: string, targetName: string, value: number }) {
    const key = `${source}:${target}`;
    const link: SankeyChartLinkItem | undefined = linksMap[key];

    if (!link) {
        linksMap[key] = {
            sourceItemType: sourceItemType,
            sourceItemId: sourceItemId,
            source: source,
            sourceDisplayName: sourceName,
            targetItemType: targetItemType,
            targetItemId: targetItemId,
            target: target,
            targetDisplayName: targetName,
            value: value
        };
    } else {
        link.value += value;
    }
}

function addFinalSortedNodeItems(nodesMap: Record<string, SankeyChartNodeItem>, allNodesArray: SankeyChartNodeItem[]): void {
    const nodesArray: SankeyChartNodeItem[] = [];

    for (const node of values(nodesMap)) {
        nodesArray.push(node);
    }

    sortStatisticsItems(nodesArray, props.sortingType);

    for (const node of nodesArray) {
        node.name = node.nameId;
    }

    allNodesArray.push(...nodesArray);
}

function addFinalLinkItems(linksMap: Record<string, SankeyChartLinkItem>, allLinksArray: SankeyChartLinkItem[]): void {
    const linksArray: SankeyChartLinkItem[] = [];

    for (const link of values(linksMap)) {
        linksArray.push(link);
    }

    allLinksArray.push(...linksArray);
}

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series' || e.seriesType !=='sankey') {
        return;
    }

    if (!e.data) {
        return;
    }

    if (e.dataType === 'node') {
        const dataItem = e.data as SankeyChartNodeItem;
        emit('click', dataItem.itemType, dataItem.itemId);
    } else if (e.dataType === 'edge') {
        const dataItem = e.data as SankeyChartLinkItem;

        if (dataItem.sourceItemType === dataItem.targetItemType && dataItem.sourceItemId === dataItem.targetItemId) {
            emit('click', dataItem.sourceItemType, dataItem.sourceItemId);
        } else {
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
