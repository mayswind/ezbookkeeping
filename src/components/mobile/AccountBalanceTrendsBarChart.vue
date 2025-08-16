<template>
    <f7-list class="skeleton-text margin-top-half" media-list v-if="loading">
        <f7-list-item class="account-balance-trends-list-item" title="Date Range" after="0.00 USD"
                      :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ]">
            <template #media>
                <f7-icon f7="app_fill"></f7-icon>
            </template>
            <template #inner>
                <div class="display-flex padding-top-half">
                    <div class="account-balance-percent-line width-100">
                        <f7-progressbar :progress="0"></f7-progressbar>
                    </div>
                </div>
            </template>
        </f7-list-item>
    </f7-list>

    <f7-list v-else-if="!loading && (!allVirtualListItems || !allVirtualListItems.length)">
        <f7-list-item :title="tt('No transaction data')"></f7-list-item>
    </f7-list>

    <f7-list class="margin-top-half" media-list virtual-list :virtual-list-params="{ items: allVirtualListItems, renderExternal, height: 'auto' }"
             :key="`account-balance-trends-${dateAggregationType}`"
             v-else-if="!loading && allVirtualListItems && allVirtualListItems.length > 0">
        <ul>
            <f7-list-item class="account-balance-trends-list-item"
                          :key="item.index"
                          :style="`top: ${virtualDataItems.topPosition}px`"
                          :virtual-list-index="item.index"
                          :title="item.displayDate"
                          :after="formatAmountToLocalizedNumeralsWithCurrency(item.closingBalance, account.currency)"
                          v-for="item in virtualDataItems.items"
            >
                <template #media>
                    <f7-icon f7="calendar"></f7-icon>
                </template>
                <template #inner>
                    <div class="display-flex padding-top-half">
                        <div class="account-balance-percent-line" :style="{ 'width': item.percent + '%' }">
                            <f7-progressbar :progress="100" :style="{ '--f7-progressbar-progress-color': (item.color ? item.color : '') } "></f7-progressbar>
                        </div>
                        <div class="account-balance-percent-line" :style="{ 'width': (100.0 - item.percent) + '%' }"
                             v-if="item.percent < 100.0">
                            <f7-progressbar :progress="0"></f7-progressbar>
                        </div>
                    </div>
                </template>
            </f7-list-item>
        </ul>
    </f7-list>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type AccountBalanceTrendsChartItem,
    type CommonAccountBalanceTrendsChartProps,
    useAccountBalanceTrendsChartBase
} from '@/components/base/AccountBalanceTrendsChartBase.ts'

import type { ColorValue } from '@/core/color.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

interface MobileAccountBalanceTrendsChartItem extends AccountBalanceTrendsChartItem {
    index: number;
    percent: number;
    color: ColorValue;
}

interface MobileAccountBalanceTrendsChartProps extends CommonAccountBalanceTrendsChartProps {
    loading?: boolean;
}

interface MobileAccountBalanceTrendsChartVirtualListData {
    items: MobileAccountBalanceTrendsChartItem[],
    topPosition: number
}

const props = defineProps<MobileAccountBalanceTrendsChartProps>();

const { tt, formatAmountToLocalizedNumeralsWithCurrency } = useI18n();
const { allDataItems } = useAccountBalanceTrendsChartBase(props);

const virtualDataItems = ref<MobileAccountBalanceTrendsChartVirtualListData>({
    items: [],
    topPosition: 0
});

const allVirtualListItems = computed<MobileAccountBalanceTrendsChartItem[]>(() => {
    const ret: MobileAccountBalanceTrendsChartItem[] = [];
    let maxClosingBalance = 0;

    for (let i = 0; i < allDataItems.value.length; i++) {
        const dataItem = allDataItems.value[i];

        if (dataItem.closingBalance > maxClosingBalance) {
            maxClosingBalance = dataItem.closingBalance;
        }

        const finalDataItem: MobileAccountBalanceTrendsChartItem = {
            index: i,
            displayDate: dataItem.displayDate,
            openingBalance: dataItem.openingBalance,
            closingBalance: dataItem.closingBalance,
            medianBalance: dataItem.medianBalance,
            averageBalance: dataItem.averageBalance,
            minimumBalance: dataItem.minimumBalance,
            maximumBalance: dataItem.maximumBalance,
            color: `#${DEFAULT_CHART_COLORS[0]}`,
            percent: 0.0
        };

        ret.push(finalDataItem);
    }

    for (let i = 0; i < ret.length; i++) {
        if (maxClosingBalance > 0 && ret[i].closingBalance > 0) {
            ret[i].percent = 100.0 * ret[i].closingBalance / maxClosingBalance;
        } else {
            ret[i].percent = 0.0;
        }
    }

    return ret;
});

function renderExternal(vl: unknown, vlData: MobileAccountBalanceTrendsChartVirtualListData): void {
    virtualDataItems.value = vlData;
}
</script>

<style>
.account-balance-trends-list-item .account-balance-percent-line {
    --f7-progressbar-bg-color: #f8f8f8;
}
</style>
