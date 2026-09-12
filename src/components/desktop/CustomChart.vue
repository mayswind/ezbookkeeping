<template>
    <v-row :class="{ 'readonly': disabled }">
        <v-col cols="12" :md="displayLayout.showChart || displayLayout.showChartData ? 6 : 12" v-if="displayLayout.showCode">
            <div class="title-and-toolbar d-flex w-100 mb-1">
                <v-btn density="compact" variant="tonal" :prepend-icon="mdiPlay"
                       :disabled="disabled || !sandboxLoaded || executingScript" :loading="executingScript"
                       @click="executeCustomScript">
                    <template #loader>
                        <v-progress-circular indeterminate size="18" class="me-1"/>
                        <span>{{ tt('Execute Custom Script') }}</span>
                    </template>
                    <span>{{ tt('Execute Custom Script') }}</span>
                </v-btn>
            </div>
            <code-editor class="w-100" style="height: 620px" language="javascript"
                         :readonly="disabled" :line-numbers="true"
                         :extra-libs="customChartEditorExtraLibs" v-model="customScript" />
        </v-col>
        <v-col cols="12" :md="displayLayout.showCode ? 6 : 12" v-if="displayLayout.showChart">
            <div class="w-100 custom-chart-container">
                <div class="w-100 h-100 d-flex align-center justify-content-center" v-if="executionError">
                    <v-alert class="mb-3" density="compact" type="error" variant="tonal" :text="executionError" />
                </div>
                <v-chart autoresize :option="chartOptions" :update-options="{ notMerge: true }"
                         v-if="!executionError && chartOptions" />
            </div>
        </v-col>
        <v-col cols="12" :md="displayLayout.showCode ? 6 : 12" v-else-if="displayLayout.showChartData">
            <div class="title-and-toolbar d-flex align-center w-100 mb-1">
                <span class="text-body-large">{{ tt('Chart Data') }}</span>
            </div>
            <div class="w-100">
                <code-editor class="w-100" style="height: 620px" language="json"
                             :readonly="true" :model-value="displayChartData" />
            </div>
        </v-col>
    </v-row>

    <iframe id="sandbox" ref="sandbox" sandbox="allow-scripts" style="display: none;"></iframe>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import { type CodeEditorExtraLib } from '@/components/desktop/CodeEditor.vue';

import { ref, computed, useTemplateRef, onMounted, onUnmounted, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import { TransactionType } from '@/core/transaction.ts';
import { TransactionExplorerCustomChartDisplayLayout } from '@/core/explorer.ts';
import type { TransactionInsightDataItem } from '@/models/transaction.ts';

import { parseDateTimeFromUnixTimeWithTimezoneOffset } from '@/lib/datetime.ts';
import logger from '@/lib/logger.ts';

import {
    mdiPlay
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

type SandboxRequest = {
    source: string;
    transactions: CustomChartTransaction[];
    code: string;
};

type SandboxResponse = {
    result?: string;
    knownError?: string;
    error?: string;
};

interface CustomChartTransaction {
    id: string;
    unixTime: number;
    year: number;
    month: number;
    dayOfMonth: number;
    dayOfWeek: number;
    hour: number;
    minute: number;
    second: number;
    utcOffset: number;
    type: string;
    primaryCategoryName: string;
    secondaryCategoryName: string;
    sourceAccountName: string;
    destinationAccountName?: string;
    sourceCurrency: string;
    destinationCurrency?: string;
    sourceAmount: number;
    destinationAmount: number;
    tagNames: string[];
    geoLocation?: {
        latitude: number;
        longitude: number
    };
    comment: string;
}

const props = defineProps<{
    disabled?: boolean;
    displayLayout: number;
    transactions: TransactionInsightDataItem[];
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const {
    tt
} = useI18n();

const settingsStore = useSettingsStore();

const sandboxMessageSignature: string = '#ezBookkeeping-sandbox-message#';
const sandboxBuildinScripts: string = `
<script>
window.TransactionType = {
    Income: '${tt('Income')}',
    Expense: '${tt('Expense')}',
    Transfer: '${tt('Transfer')}'
};

window.ChartColors = ${JSON.stringify(settingsStore.chartColorList.map(color => '#' + color))};

window.addEventListener('message', function (event) {
    if (!event.data || typeof event.data !== 'string' || event.data[0] !== '{' || event.data.indexOf('${sandboxMessageSignature}') <= 0) {
        return;
    }

    try {
        const data = JSON.parse(event.data);

        if (!data || !data.source || data.source !== '${sandboxMessageSignature}') {
            return;
        }

        const transactions = data.transactions;
        eval(data.code);

        if (window.buildChartOptions) {
            const result = window.buildChartOptions(transactions);
            window.parent.postMessage({ result: JSON.stringify(result) }, '*');
        } else {
            window.parent.postMessage({ knownError: 'No buildChartOptions function defined' }, '*');
        }
    } catch (error) {
        window.parent.postMessage({ error: error.message }, '*');
    }
});
<\/script>
`;

const sandbox = useTemplateRef<HTMLIFrameElement>('sandbox');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const sandboxLoaded = ref<boolean>(false);
const customScript = ref<string>('');
const chartOptions = ref<Record<string, unknown>>({});
const executingScript = ref<boolean>(false);
const executionError = ref<string>('');

const displayLayout = computed<TransactionExplorerCustomChartDisplayLayout>(() => TransactionExplorerCustomChartDisplayLayout.valueOf(props.displayLayout) ?? TransactionExplorerCustomChartDisplayLayout.Default);

const customChartEditorExtraLibs = computed<CodeEditorExtraLib[]>(() => [{
    filePath: 'inmemory://model/ezbookkeeping-custom-chart.d.ts',
    content: `
/** ${tt('sample.insightsExplorerCustomChart.customChartTransactionDescription')} */
interface CustomChartTransaction {
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.id')} */
    id: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.unixTime')} */
    unixTime: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.year')} */
    year: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.month')} */
    month: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.dayOfMonth')} */
    dayOfMonth: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.dayOfWeek')} */
    dayOfWeek: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.hour')} */
    hour: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.minute')} */
    minute: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.second')} */
    second: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.utcOffset')} */
    utcOffset: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.type')} */
    type: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.primaryCategoryName')} */
    primaryCategoryName: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.secondaryCategoryName')} */
    secondaryCategoryName: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.sourceAccountName')} */
    sourceAccountName: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.destinationAccountName')} */
    destinationAccountName?: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.sourceCurrency')} */
    sourceCurrency: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.destinationCurrency')} */
    destinationCurrency?: string;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.sourceAmount')} */
    sourceAmount: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.destinationAmount')} */
    destinationAmount: number;
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.tagNames')} */
    tagNames: string[];
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.geoLocation')} */
    geoLocation?: {
        latitude: number;
        longitude: number;
    };
    /** ${tt('sample.insightsExplorerCustomChart.transactionField.comment')} */
    comment: string;
}

/** ${tt('sample.insightsExplorerCustomChart.transactionTypeDescription')} */
declare const TransactionType: Readonly<{
    /** ${tt('Income')} */
    Income: string;
    /** ${tt('Expense')} */
    Expense: string;
    /** ${tt('Transfer')} */
    Transfer: string;
}>;

/** ${tt('sample.insightsExplorerCustomChart.chartColorsDescription')} */
declare const ChartColors: readonly string[];`
}]);

const sampleScript = computed<string>(() => `// ${tt('sample.insightsExplorerCustomChart.headerComment')}
/**
 * ${tt('sample.insightsExplorerCustomChart.functionDescription')}
 *
 * @param {CustomChartTransaction[]} transactions ${tt('sample.insightsExplorerCustomChart.transactionsFieldsDescription')}
 * @returns {object} ${tt('sample.insightsExplorerCustomChart.functionReturnDescription')}
 */
function buildChartOptions(transactions) {
    const transactionsByDate = new Map();

    transactions.forEach(transaction => {
        const date = transaction.year + '-'
            + String(transaction.month).padStart(2, '0') + '-'
            + String(transaction.dayOfMonth).padStart(2, '0');
        const dailyTransactions = transactionsByDate.get(date) || [];
        dailyTransactions.push(transaction);
        transactionsByDate.set(date, dailyTransactions);
    });

    const dailyData = Array.from(transactionsByDate.entries())
        .map(([date, dailyTransactions]) => ({
            date: date,
            year: dailyTransactions[0].year,
            month: dailyTransactions[0].month,
            dayOfMonth: dailyTransactions[0].dayOfMonth,
            totalAmount: sumTransactionAmounts(dailyTransactions)
        }))
        .sort((a, b) => a.date.localeCompare(b.date));

    const dailyDataMap = new Map(dailyData.map(item => [item.date, item]));
    const dates = dailyData.map(item => item.date);
    const totalAmounts = dailyData.map(item => item.totalAmount);

    const averageAmount = totalAmounts.length > 0
        ? totalAmounts.reduce((sum, amount) => sum + amount, 0) / totalAmounts.length
        : 0;

    const dayOverDayRates = dailyData.map(item => {
        const previousDate = new Date(Date.UTC(item.year, item.month - 1, item.dayOfMonth));
        previousDate.setUTCDate(previousDate.getUTCDate() - 1);

        const previousDateKey = previousDate.getUTCFullYear() + '-'
            + String(previousDate.getUTCMonth() + 1).padStart(2, '0') + '-'
            + String(previousDate.getUTCDate()).padStart(2, '0');
        const previousDayData = dailyDataMap.get(previousDateKey);

        if (!previousDayData || previousDayData.totalAmount === 0) {
            return null;
        }

        return (item.totalAmount - previousDayData.totalAmount) / previousDayData.totalAmount * 100;
    });

    const hasDayOverDayRates = dayOverDayRates.some(rate => rate !== null);
    const valueDimensions = [
        { name: 'date', type: 'ordinal', tooltip: false },
        { name: 'value', type: 'float', tooltip: false },
        { name: 'formattedValue', type: 'ordinal' }
    ];

    const series = [
        {
            name: ${JSON.stringify(tt('Total Amount'))},
            type: 'bar',
            yAxisIndex: 0,
            barMaxWidth: 40,
            itemStyle: { color: ChartColors[0] },
            dimensions: valueDimensions,
            encode: { x: 'date', y: 'formattedValue', tooltip: ['formattedValue'] },
            data: dailyData.map(item => [item.date, item.totalAmount, (item.totalAmount / 100).toFixed(2)])
        },
        {
            name: ${JSON.stringify(tt('Average Amount'))},
            type: 'line',
            yAxisIndex: 0,
            symbol: 'none',
            lineStyle: { type: 'dashed' },
            itemStyle: { color: ChartColors[1] },
            dimensions: valueDimensions,
            encode: { x: 'date', y: 'formattedValue', tooltip: ['formattedValue'] },
            data: dailyData.map(item => [item.date, averageAmount, (averageAmount / 100).toFixed(2)])
        }
    ];

    if (hasDayOverDayRates) {
        series.push({
            name: ${JSON.stringify(tt('Day-over-Day'))},
            type: 'line',
            yAxisIndex: 1,
            smooth: true,
            connectNulls: false,
            itemStyle: { color: ChartColors[2] },
            dimensions: valueDimensions,
            encode: { x: 'date', y: 'value', tooltip: ['formattedValue'] },
            data: dailyData.flatMap((item, index) => {
                const rate = dayOverDayRates[index];
                return rate !== null ? [[item.date, rate, rate.toFixed(2) + '%']] : [];
            })
        });
    }

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'cross' }
        },
        legend: {
            top: 0,
            type: 'scroll'
        },
        grid: {
            top: 50,
            right: 50,
            bottom: 30,
            left: 50
        },
        xAxis: {
            type: 'category',
            data: dates
        },
        yAxis: [
            {
                type: 'value',
                name: ${JSON.stringify(tt('Amount'))},
                splitLine: {
                    lineStyle: { color: '#e1e6f2' }
                }
            },
            {
                type: 'value',
                name: ${JSON.stringify(tt('Day-over-Day'))},
                axisLabel: { formatter: '{value}%' },
                splitLine: { show: false }
            }
        ],
        series: series
    };
}

/**
 * ${tt('sample.insightsExplorerCustomChart.sumTransactionAmountsDescription')}
 *
 * @param {CustomChartTransaction[]} transactions ${tt('sample.insightsExplorerCustomChart.transactionsFieldsDescription')}
 * @returns {number} ${tt('sample.insightsExplorerCustomChart.sumTransactionAmountsReturnDescription')}
 */
function sumTransactionAmounts(transactions) {
    if (!Array.isArray(transactions) || transactions.length < 1) {
        return 0;
    }

    return transactions.reduce((sum, transaction) => {
        if (typeof transaction.sourceAmount === 'number') {
            return sum + transaction.sourceAmount;
        }

        return sum;
    }, 0);
}

/**
 * ${tt('sample.insightsExplorerCustomChart.allTransactionFieldsDescription')}
 * {string} id - ${tt('sample.insightsExplorerCustomChart.transactionField.id')}
 * {number} unixTime - ${tt('sample.insightsExplorerCustomChart.transactionField.unixTime')}
 * {number} year - ${tt('sample.insightsExplorerCustomChart.transactionField.year')}
 * {number} month - ${tt('sample.insightsExplorerCustomChart.transactionField.month')}
 * {number} dayOfMonth - ${tt('sample.insightsExplorerCustomChart.transactionField.dayOfMonth')}
 * {number} dayOfWeek - ${tt('sample.insightsExplorerCustomChart.transactionField.dayOfWeek')}
 * {number} hour - ${tt('sample.insightsExplorerCustomChart.transactionField.hour')}
 * {number} minute - ${tt('sample.insightsExplorerCustomChart.transactionField.minute')}
 * {number} second - ${tt('sample.insightsExplorerCustomChart.transactionField.second')}
 * {number} utcOffset - ${tt('sample.insightsExplorerCustomChart.transactionField.utcOffset')}
 * {string} type - ${tt('sample.insightsExplorerCustomChart.transactionField.type')}
 * {string} primaryCategoryName - ${tt('sample.insightsExplorerCustomChart.transactionField.primaryCategoryName')}
 * {string} secondaryCategoryName - ${tt('sample.insightsExplorerCustomChart.transactionField.secondaryCategoryName')}
 * {string} sourceAccountName - ${tt('sample.insightsExplorerCustomChart.transactionField.sourceAccountName')}
 * {string} destinationAccountName - ${tt('sample.insightsExplorerCustomChart.transactionField.destinationAccountName')}
 * {string} sourceCurrency - ${tt('sample.insightsExplorerCustomChart.transactionField.sourceCurrency')}
 * {string} destinationCurrency - ${tt('sample.insightsExplorerCustomChart.transactionField.destinationCurrency')}
 * {number} sourceAmount - ${tt('sample.insightsExplorerCustomChart.transactionField.sourceAmount')}
 * {number} destinationAmount - ${tt('sample.insightsExplorerCustomChart.transactionField.destinationAmount')}
 * {string[]} tagNames - ${tt('sample.insightsExplorerCustomChart.transactionField.tagNames')}
 * {string} geoLocation - ${tt('sample.insightsExplorerCustomChart.transactionField.geoLocation')}
 * {string} comment - ${tt('sample.insightsExplorerCustomChart.transactionField.comment')}
 */

/**
 * ${tt('sample.insightsExplorerCustomChart.availableGlobalConstants')}
 * {TransactionType} TransactionType - ${tt('sample.insightsExplorerCustomChart.transactionTypeDescription')}
 * {string[]} ChartColors - ${tt('sample.insightsExplorerCustomChart.chartColorsDescription')}
 */`);

const displayChartData = computed<string>(() => {
    if (executingScript.value) {
        return tt('Executing Script...');
    } else if (executionError.value) {
        return `// ${executionError.value}`;
    } else if (chartOptions.value) {
        return JSON.stringify(chartOptions.value, null, 2);
    } else {
        return tt('No Preview Result');
    }
});

const customChartTransactions = computed<CustomChartTransaction[]>(() => {
    return props.transactions.map(transaction => {
        const transactionTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        let transactionType: string = '';

        if (transaction.type === TransactionType.Income) {
            transactionType = tt('Income');
        } else if (transaction.type === TransactionType.Expense) {
            transactionType = tt('Expense');
        } else if (transaction.type === TransactionType.Transfer) {
            transactionType = tt('Transfer');
        } else if (transaction.type === TransactionType.ModifyBalance) {
            transactionType = tt('Modify Balance');
        } else {
            transactionType = tt('Unknown');
        }

        const finalTransaction: CustomChartTransaction = {
            id: transaction.id,
            unixTime: transaction.time,
            year: transactionTime.getGregorianCalendarYear(),
            month: transactionTime.getGregorianCalendarMonth(),
            dayOfMonth: transactionTime.getGregorianCalendarDay(),
            dayOfWeek: transactionTime.getWeekDay().type,
            hour: transactionTime.getHour(),
            minute: transactionTime.getMinute(),
            second: transactionTime.getSecond(),
            utcOffset: transaction.utcOffset,
            type: transactionType,
            primaryCategoryName: transaction.primaryCategoryName,
            secondaryCategoryName: transaction.secondaryCategoryName,
            sourceAccountName: transaction.sourceAccountName,
            destinationAccountName: transaction.destinationAccountName,
            sourceCurrency: transaction.sourceAccount.currency,
            destinationCurrency: transaction.destinationAccount?.currency,
            sourceAmount: transaction.sourceAmount,
            destinationAmount: transaction.destinationAmount,
            tagNames: transaction.tags?.map(tag => tag.name) ?? [],
            geoLocation: transaction.geoLocation ? {
                latitude: transaction.geoLocation.latitude,
                longitude: transaction.geoLocation.longitude
            } : undefined,
            comment: transaction.comment || ''
        };

        return finalTransaction;
    });
});

function reloadSandbox(executeAfterLoaded: boolean): void {
    sandboxLoaded.value = false;

    if (sandbox.value) {
        sandbox.value.src = 'about:blank';
        sandbox.value.srcdoc = sandboxBuildinScripts;

        sandbox.value.onload = () => {
            sandboxLoaded.value = true;

            if (executeAfterLoaded && (displayLayout.value.showChart || displayLayout.value.showChartData)) {
                executeCustomScript();
            }
        };
    }
}

function executeCustomScript(): void {
    if (!sandbox.value || props.disabled || executingScript.value) {
        return;
    }

    executingScript.value = true;

    const sandboxRequest: SandboxRequest = {
        source: sandboxMessageSignature,
        transactions: customChartTransactions.value || [],
        code: customScript.value + `\n\n;if (typeof buildChartOptions !== 'undefined') { window.buildChartOptions = buildChartOptions; }`
    };

    sandbox.value?.contentWindow?.postMessage(JSON.stringify(sandboxRequest), '*');
}

function onMessage(event: MessageEvent<SandboxResponse>): void {
    if (event.source !== sandbox.value?.contentWindow) {
        return;
    }

    executingScript.value = false;

    const data = event.data;

    if (data.knownError) {
        snackbar.value?.showError(data.knownError);
        executionError.value = tt(data.knownError);
    } else if (data.error) {
        logger.error('Failed to execute custom script: ' + data.error);
        snackbar.value?.showError('Failed to execute custom script');
        executionError.value = data.error;
    } else if (data.result) {
        try {
            chartOptions.value = JSON.parse(data.result) as Record<string, unknown>;
            executionError.value = '';
        } catch (error) {
            logger.error('Failed to build custom chart options', error);
            executionError.value = tt('Failed to execute custom script');
        }
    }

    reloadSandbox(false);
}

onMounted(() => {
    customScript.value = props.modelValue || sampleScript.value;
    reloadSandbox(true);
    window.addEventListener('message', onMessage);
});

onUnmounted(() => {
    window.removeEventListener('message', onMessage);
});

watch(() => props.displayLayout, (newValue, oldValue) => {
    if ((TransactionExplorerCustomChartDisplayLayout.valueOf(newValue)?.showChart || TransactionExplorerCustomChartDisplayLayout.valueOf(newValue)?.showChartData)
        && !(TransactionExplorerCustomChartDisplayLayout.valueOf(oldValue)?.showChart || TransactionExplorerCustomChartDisplayLayout.valueOf(oldValue)?.showChartData)) {
        reloadSandbox(true);
    }
});

watch(customChartTransactions, () => {
    reloadSandbox(true);
});

watch(customScript, value => {
    emit('update:modelValue', value);
});

watch(() => props.modelValue, value => {
    if (value !== customScript.value) {
        customScript.value = value;
    }
});
</script>

<style scoped>
.custom-chart-container {
    width: 100%;
    height: 660px;
}
</style>
