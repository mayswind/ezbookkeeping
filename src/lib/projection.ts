import type { BigDecimal } from '@/core/numeral.ts';
import { CategoryType } from '@/core/category.ts';
import type { TextualYearMonth, Year0BasedMonth, Year1BasedMonth } from '@/core/datetime.ts';
import {
    type ProjectionMonth,
    type ProjectionRow,
    type ProjectionCategoryRow,
    type ProjectionSectionRow,
    type ProjectionTableData,
    ProjectionMonthKind
} from '@/core/projection.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionStatisticTrendsResponseItem } from '@/models/transaction.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { getAllMonthsStartAndEndUnixTimes } from '@/lib/datetime.ts';

export interface ProjectionTableOptions {
    // every category of the user, by id, to resolve the name and the parent of each amount
    readonly categoriesMap: Record<string, TransactionCategory>;
    // converts an amount of the specified account to the default currency of the user, returning null
    // when there is no exchange rate available for it
    readonly convertAmount: (amount: string, accountId: string) => BigDecimal | null;
    // year * 100 + month of the month the user is currently in, which is the mixed one
    readonly currentYearMonth: number;
}

interface WritableProjectionRow {
    id: string;
    name: string;
    displayOrder: number;
    monthlyAmounts: BigDecimal[];
    subCategories?: Record<string, WritableProjectionRow>;
}

const ZERO: BigDecimal = parseBigDecimal(0);

// getProjectionMonths returns one column descriptor per month of the specified period
export function getProjectionMonths(startYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year0BasedMonth | Year1BasedMonth | TextualYearMonth | '', currentYearMonth: number): ProjectionMonth[] {
    const months: ProjectionMonth[] = [];

    for (const yearMonthTime of getAllMonthsStartAndEndUnixTimes(startYearMonth, endYearMonth)) {
        const month = yearMonthTime.month0base + 1;
        const yearMonth = yearMonthTime.year * 100 + month;

        months.push({
            year: yearMonthTime.year,
            month: month,
            yearMonth: yearMonth,
            kind: getProjectionMonthKind(yearMonth, currentYearMonth)
        });
    }

    return months;
}

function getProjectionMonthKind(yearMonth: number, currentYearMonth: number): ProjectionMonthKind {
    if (yearMonth < currentYearMonth) {
        return ProjectionMonthKind.Actual;
    } else if (yearMonth > currentYearMonth) {
        return ProjectionMonthKind.Projected;
    }

    return ProjectionMonthKind.Mixed;
}

// buildProjectionTableData turns the response of the projections api into the rows of the projection
// table: two collapsible sections with their categories and sub-categories, the net of every month
// and the running total of the period.
//
// Whether an amount is an income or an expense comes from the type of its category, because the api
// does not send it. Amounts of accounts held in a currency that cannot be converted are left out and
// reported through hasUnconvertedAmounts, so the table never mixes currencies silently.
export function buildProjectionTableData(items: TransactionStatisticTrendsResponseItem[], months: ProjectionMonth[], options: ProjectionTableOptions): ProjectionTableData {
    const monthIndexMap: Record<number, number> = {};

    for (let i = 0; i < months.length; i++) {
        monthIndexMap[months[i]!.yearMonth] = i;
    }

    const incomeCategories: Record<string, WritableProjectionRow> = {};
    const expenseCategories: Record<string, WritableProjectionRow> = {};
    let hasUnconvertedAmounts = false;

    for (const monthItem of items) {
        const monthIndex = monthIndexMap[monthItem.year * 100 + monthItem.month];

        if (monthIndex === undefined) {
            continue;
        }

        for (const item of monthItem.items) {
            const category = options.categoriesMap[item.categoryId];

            if (!category) {
                continue;
            }

            // the api excludes transfers from both halves already, this only guards against a
            // category whose type changed after its transactions were created
            if (category.type !== CategoryType.Income && category.type !== CategoryType.Expense) {
                continue;
            }

            const amount = options.convertAmount(item.amount, item.accountId);

            if (amount === null) {
                hasUnconvertedAmounts = true;
                continue;
            }

            const primaryCategory = category.parentId && category.parentId !== '0' ? options.categoriesMap[category.parentId] : category;

            if (!primaryCategory) {
                continue;
            }

            const categories = category.type === CategoryType.Income ? incomeCategories : expenseCategories;
            const categoryRow = getOrCreateRow(categories, primaryCategory.id, primaryCategory.name, primaryCategory.displayOrder, months.length);
            const subCategoryRow = getOrCreateRow(categoryRow.subCategories!, category.id, category.name, category.displayOrder, months.length);

            subCategoryRow.monthlyAmounts[monthIndex] = subCategoryRow.monthlyAmounts[monthIndex]!.add(amount);
        }
    }

    const income = toSectionRow(CategoryType.Income, incomeCategories, months.length);
    const expense = toSectionRow(CategoryType.Expense, expenseCategories, months.length);

    const netAmounts: BigDecimal[] = [];
    const accumulatedAmounts: BigDecimal[] = [];
    let runningTotal: BigDecimal = ZERO;

    for (let i = 0; i < months.length; i++) {
        const net = income.monthlyAmounts[i]!.subtract(expense.monthlyAmounts[i]!);
        runningTotal = runningTotal.add(net);

        netAmounts.push(net);
        accumulatedAmounts.push(runningTotal);
    }

    return {
        months: months,
        income: income,
        expense: expense,
        net: {
            id: 'net',
            name: 'Net',
            monthlyAmounts: netAmounts,
            totalAmount: sumAmounts(netAmounts)
        },
        accumulated: {
            id: 'accumulated',
            name: 'Accumulated',
            monthlyAmounts: accumulatedAmounts,
            // already a running total, so the total column is the value it reached, not the sum
            totalAmount: accumulatedAmounts.length ? accumulatedAmounts[accumulatedAmounts.length - 1]! : ZERO
        },
        hasUnconvertedAmounts: hasUnconvertedAmounts
    };
}

function getOrCreateRow(rows: Record<string, WritableProjectionRow>, id: string, name: string, displayOrder: number, monthCount: number): WritableProjectionRow {
    let row = rows[id];

    if (!row) {
        row = {
            id: id,
            name: name,
            displayOrder: displayOrder,
            monthlyAmounts: newZeroAmounts(monthCount),
            subCategories: {}
        };

        rows[id] = row;
    }

    return row;
}

function toSectionRow(categoryType: CategoryType, categories: Record<string, WritableProjectionRow>, monthCount: number): ProjectionSectionRow {
    const categoryRows: ProjectionCategoryRow[] = [];
    const sectionAmounts = newZeroAmounts(monthCount);

    for (const writableCategory of sortRows(categories)) {
        const subCategoryRows: ProjectionRow[] = [];
        const categoryAmounts = newZeroAmounts(monthCount);

        for (const writableSubCategory of sortRows(writableCategory.subCategories!)) {
            for (let i = 0; i < monthCount; i++) {
                categoryAmounts[i] = categoryAmounts[i]!.add(writableSubCategory.monthlyAmounts[i]!);
            }

            subCategoryRows.push({
                id: writableSubCategory.id,
                name: writableSubCategory.name,
                monthlyAmounts: writableSubCategory.monthlyAmounts,
                totalAmount: sumAmounts(writableSubCategory.monthlyAmounts)
            });
        }

        for (let i = 0; i < monthCount; i++) {
            sectionAmounts[i] = sectionAmounts[i]!.add(categoryAmounts[i]!);
        }

        categoryRows.push({
            id: writableCategory.id,
            name: writableCategory.name,
            monthlyAmounts: categoryAmounts,
            totalAmount: sumAmounts(categoryAmounts),
            subCategories: subCategoryRows
        });
    }

    return {
        id: categoryType === CategoryType.Income ? 'income' : 'expense',
        name: categoryType === CategoryType.Income ? 'Total Income' : 'Total Expense',
        categoryType: categoryType,
        monthlyAmounts: sectionAmounts,
        totalAmount: sumAmounts(sectionAmounts),
        categories: categoryRows
    };
}

function sortRows(rows: Record<string, WritableProjectionRow>): WritableProjectionRow[] {
    return Object.values(rows).sort((row1, row2) => {
        if (row1.displayOrder !== row2.displayOrder) {
            return row1.displayOrder - row2.displayOrder;
        }

        return row1.name.localeCompare(row2.name);
    });
}

function newZeroAmounts(monthCount: number): BigDecimal[] {
    const amounts: BigDecimal[] = [];

    for (let i = 0; i < monthCount; i++) {
        amounts.push(ZERO);
    }

    return amounts;
}

function sumAmounts(amounts: BigDecimal[]): BigDecimal {
    let total: BigDecimal = ZERO;

    for (const amount of amounts) {
        total = total.add(amount);
    }

    return total;
}
