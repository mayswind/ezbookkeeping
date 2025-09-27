import { computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

export interface CalendarHeatmapDataItem {
    date: string; // YYYY-MM-DD format
    value: number;
    count: number;
    income: number;
    expense: number;
    netAmount: number; // income - expense
    displayValue?: string;
    displayIncome?: string;
    displayExpense?: string;
    displayDate?: string;
    level: number; // 0-4 intensity levels
    isIncome: boolean; // true if net is positive (income > expense)
}

export interface CommonCalendarHeatmapProps {
    data: Record<string, { value: number; count: number; income: number; expense: number }>;
    defaultCurrency?: string;
    maxValue?: number;
    colorScheme?: string[];
    incomeColorScheme?: string[];
    expenseColorScheme?: string[];
    showWeekdays?: boolean;
    showMonthLabels?: boolean;
    endDate?: Date; // Optional end date, defaults to today
}

export function useCalendarHeatmapBase(props: CommonCalendarHeatmapProps) {
    const { formatAmountToLocalizedNumeralsWithCurrency, formatNumberToLocalizedNumerals } = useI18n();

    const defaultColors = [
        '#ebedf0', // level 0 - no activity
        '#9be9a8', // level 1 - low activity
        '#40c463', // level 2 - medium activity
        '#30a14e', // level 3 - high activity
        '#216e39'  // level 4 - very high activity
    ];

    // Income color scheme (red tones) - uses system income color
    const defaultIncomeColors = [
        '#ebedf0', // level 0 - no activity
        'rgba(212, 63, 63, 0.3)', // level 1 - low income
        'rgba(212, 63, 63, 0.5)', // level 2 - medium income
        'rgba(212, 63, 63, 0.7)', // level 3 - high income
        'rgba(212, 63, 63, 0.9)'  // level 4 - very high income
    ];

    // Expense color scheme (green tones) - uses system expense color
    const defaultExpenseColors = [
        '#ebedf0', // level 0 - no activity
        'rgba(0, 150, 136, 0.3)', // level 1 - low expense
        'rgba(0, 150, 136, 0.5)', // level 2 - medium expense
        'rgba(0, 150, 136, 0.7)', // level 3 - high expense
        'rgba(0, 150, 136, 0.9)'  // level 4 - very high expense
    ];

    const colors = computed(() => props.colorScheme || defaultColors);
    const incomeColors = computed(() => props.incomeColorScheme || defaultIncomeColors);
    const expenseColors = computed(() => props.expenseColorScheme || defaultExpenseColors);

    const maxDataValue = computed(() => {
        if (props.maxValue) return props.maxValue;

        let maxIncome = 0;
        let maxExpense = 0;
        Object.values(props.data || {}).forEach(item => {
            const income = item.income || 0;
            const expense = item.expense || 0;
            if (income > maxIncome) maxIncome = income;
            if (expense > maxExpense) maxExpense = expense;
        });
        return Math.max(maxIncome, maxExpense);
    });

    // Calculate the 365-day date range
    const dateRange = computed(() => {
        const endDate = props.endDate || new Date();
        const startDate = new Date(endDate);
        startDate.setDate(startDate.getDate() - 364);
        return { startDate, endDate };
    });

    const heatmapData = computed<CalendarHeatmapDataItem[]>(() => {
        const result: CalendarHeatmapDataItem[] = [];
        const { startDate, endDate } = dateRange.value;

        for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
            // Use local date string to match the data keys from transaction processing
            const year = d.getFullYear();
            const month = String(d.getMonth() + 1).padStart(2, '0');
            const day = String(d.getDate()).padStart(2, '0');
            const dateStr = `${year}-${month}-${day}`;
            const dayData = (props.data && props.data[dateStr]) ? props.data[dateStr] : {
                value: 0,
                count: 0,
                income: 0,
                expense: 0
            };

            const income = dayData.income || 0;
            const expense = dayData.expense || 0;
            const netAmount = income - expense;
            const absoluteNet = Math.abs(netAmount);
            const isIncome = netAmount > 0;

            // Calculate intensity level (0-4) based on the dominant transaction type amount
            let level = 0;
            if (maxDataValue.value > 0) {
                const dominantAmount = isIncome ? income : expense;
                if (dominantAmount > 0) {
                    const ratio = dominantAmount / maxDataValue.value;
                    if (ratio > 0.8) level = 4;
                    else if (ratio > 0.6) level = 3;
                    else if (ratio > 0.4) level = 2;
                    else if (ratio > 0.2) level = 1;
                    else if (ratio > 0) level = 1;
                }
            }

            result.push({
                date: dateStr,
                value: dayData.value,
                count: dayData.count,
                income,
                expense,
                netAmount,
                level,
                isIncome,
                displayValue: formatAmountToLocalizedNumeralsWithCurrency(absoluteNet, props.defaultCurrency || 'USD'),
                displayIncome: formatAmountToLocalizedNumeralsWithCurrency(income, props.defaultCurrency || 'USD'),
                displayExpense: formatAmountToLocalizedNumeralsWithCurrency(expense, props.defaultCurrency || 'USD'),
                displayDate: d.toLocaleDateString()
            });
        }

        return result;
    });

    const weekLabels = computed(() => ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']);
    const monthLabels = computed(() => [
        'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
        'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
    ]);

    // Calculate grid positions for items (CSS Grid based)
    const gridData = computed(() => {
        const data = heatmapData.value;
        const { startDate } = dateRange.value;

        // Find first Sunday
        const firstSunday = new Date(startDate);
        firstSunday.setDate(firstSunday.getDate() - startDate.getDay());

        const gridItems: Array<CalendarHeatmapDataItem & {
            gridRow: number;
            gridColumn: number;
        }> = [];

        const monthPositions: Array<{ month: number; column: number; name: string }> = [];
        let currentMonth = -1;

        data.forEach((item, index) => {
            const itemDate = new Date(item.date);
            const daysSinceFirstSunday = Math.floor((itemDate.getTime() - firstSunday.getTime()) / (24 * 60 * 60 * 1000));

            const column = Math.floor(daysSinceFirstSunday / 7) + 2; // +2 for month labels and weekday labels
            const row = (daysSinceFirstSunday % 7) + 2; // +2 for header

            gridItems.push({
                ...item,
                gridRow: row,
                gridColumn: column
            });

            // Track month changes for labels
            const month = itemDate.getMonth();
            if (month !== currentMonth && itemDate.getDay() === 0) { // Only on Sundays
                monthPositions.push({
                    month: month,
                    column: column,
                    name: monthLabels.value[month] || ''
                });
                currentMonth = month;
            }
        });

        return { gridItems, monthPositions };
    });

    function getColorForLevel(level: number, isIncome?: boolean): string {
        if (isIncome === true) {
            return incomeColors.value[level] || incomeColors.value[0] || '#ebedf0';
        } else if (isIncome === false) {
            return expenseColors.value[level] || expenseColors.value[0] || '#ebedf0';
        }
        return colors.value[level] || colors.value[0] || '#ebedf0';
    }

    function getTooltipText(item: CalendarHeatmapDataItem): string {
        const transactions = formatNumberToLocalizedNumerals(item.count);
        if (item.income > 0 && item.expense > 0) {
            return `${item.displayDate || item.date}\nIncome: ${item.displayIncome || ''}\nExpense: ${item.displayExpense || ''}\nNet: ${item.isIncome ? '+' : '-'}${item.displayValue || ''}\n${transactions} transactions`;
        } else if (item.income > 0) {
            return `${item.displayDate || item.date}\nIncome: ${item.displayIncome || ''}\n${transactions} transactions`;
        } else if (item.expense > 0) {
            return `${item.displayDate || item.date}\nExpense: ${item.displayExpense || ''}\n${transactions} transactions`;
        } else {
            // Check if it's a future date
            const itemDate = new Date(item.date);
            const today = new Date();
            today.setHours(23, 59, 59, 999); // End of today

            if (itemDate > today) {
                return `${item.displayDate || item.date}\nFuture date`;
            } else {
                return `${item.displayDate || item.date}\nNo transactions\nClick to add transaction`;
            }
        }
    }

    return {
        heatmapData,
        gridData,
        weekLabels,
        monthLabels,
        getColorForLevel,
        getTooltipText,
        maxDataValue
    };
}