import { describe, expect, it } from 'vitest';

import type { BigDecimal } from '@/core/numeral.ts';
import { CategoryType } from '@/core/category.ts';
import { ProjectionMonthKind } from '@/core/projection.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionStatisticTrendsResponseItem } from '@/models/transaction.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { buildProjectionTableData, getProjectionMonths, type ProjectionTableOptions } from '@/lib/projection.ts';

const ARS_ACCOUNT_ID = '3001';
const USD_ACCOUNT_ID = '3002';

// income: Trabajo > Salario, Bonos ; expense: Casa > Alquiler, Entretenimiento > Suscripciones
const CATEGORIES: Record<string, TransactionCategory> = {};

function defineCategory(id: string, name: string, parentId: string, type: CategoryType, displayOrder: number): void {
    CATEGORIES[id] = TransactionCategory.of({
        id: id,
        name: name,
        parentId: parentId,
        type: type,
        icon: '1',
        iconType: 0,
        color: '000000',
        comment: '',
        displayOrder: displayOrder,
        hidden: false
    });
}

defineCategory('10', 'Trabajo', '0', CategoryType.Income, 1);
defineCategory('11', 'Salario', '10', CategoryType.Income, 1);
defineCategory('12', 'Bonos', '10', CategoryType.Income, 2);
defineCategory('20', 'Casa', '0', CategoryType.Expense, 1);
defineCategory('21', 'Alquiler', '20', CategoryType.Expense, 1);
defineCategory('30', 'Entretenimiento', '0', CategoryType.Expense, 2);
defineCategory('31', 'Suscripciones', '30', CategoryType.Expense, 1);
defineCategory('40', 'Movimientos', '0', CategoryType.Transfer, 1);
defineCategory('41', 'Ahorro', '40', CategoryType.Transfer, 1);

// converts nothing: every account is treated as being in the default currency
const IDENTITY_OPTIONS: ProjectionTableOptions = {
    categoriesMap: CATEGORIES,
    convertAmount: (amount: string) => parseBigDecimal(amount),
    currentYearMonth: 202608
};

function month(year: number, month1base: number, items: [string, string, string][]): TransactionStatisticTrendsResponseItem {
    return {
        year: year,
        month: month1base,
        items: items.map(([categoryId, accountId, amount]) => ({
            categoryId: categoryId,
            accountId: accountId,
            amount: amount
        }))
    };
}

function amounts(values: BigDecimal[]): string[] {
    return values.map(value => value.toString());
}

const MONTHS_2026_08_TO_12 = getProjectionMonths('2026-08', '2026-12', 202608);

describe('getProjectionMonths', () => {
    it('should return one column per month of the period', () => {
        const months = getProjectionMonths('2026-11', '2027-02', 202611);

        expect(months.map(m => m.yearMonth)).toEqual([202611, 202612, 202701, 202702]);
        expect(months.map(m => m.year)).toEqual([2026, 2026, 2027, 2027]);
        expect(months.map(m => m.month)).toEqual([11, 12, 1, 2]);
    });

    it('should mark the current month as mixed and split the rest around it', () => {
        const months = getProjectionMonths('2026-06', '2026-10', 202608);

        expect(months.map(m => m.kind)).toEqual([
            ProjectionMonthKind.Actual,
            ProjectionMonthKind.Actual,
            ProjectionMonthKind.Mixed,
            ProjectionMonthKind.Projected,
            ProjectionMonthKind.Projected
        ]);
    });

    it('should return no columns for an empty period', () => {
        expect(getProjectionMonths('', '', 202608)).toEqual([]);
    });
});

describe('buildProjectionTableData', () => {
    it('should group amounts into sections, categories and sub-categories', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1000'], ['21', ARS_ACCOUNT_ID, '600']]),
            month(2026, 9, [['11', ARS_ACCOUNT_ID, '1000'], ['12', ARS_ACCOUNT_ID, '500'], ['21', ARS_ACCOUNT_ID, '600']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.income.categories.map(c => c.name)).toEqual(['Trabajo']);
        expect(data.income.categories[0]!.subCategories.map(c => c.name)).toEqual(['Salario', 'Bonos']);

        expect(data.expense.categories.map(c => c.name)).toEqual(['Casa']);
        expect(data.expense.categories[0]!.subCategories.map(c => c.name)).toEqual(['Alquiler']);

        expect(amounts(data.income.categories[0]!.subCategories[0]!.monthlyAmounts)).toEqual(['1000', '1000', '0', '0', '0']);
        expect(amounts(data.income.categories[0]!.subCategories[1]!.monthlyAmounts)).toEqual(['0', '500', '0', '0', '0']);
    });

    it('should subtotal every category and total every section', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1000'], ['12', ARS_ACCOUNT_ID, '500'], ['21', ARS_ACCOUNT_ID, '600'], ['31', ARS_ACCOUNT_ID, '200']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        // Trabajo = Salario + Bonos
        expect(amounts(data.income.categories[0]!.monthlyAmounts)).toEqual(['1500', '0', '0', '0', '0']);
        expect(data.income.categories[0]!.totalAmount.toString()).toEqual('1500');

        expect(amounts(data.income.monthlyAmounts)).toEqual(['1500', '0', '0', '0', '0']);
        expect(amounts(data.expense.monthlyAmounts)).toEqual(['800', '0', '0', '0', '0']);
        expect(data.expense.totalAmount.toString()).toEqual('800');
    });

    it('should sort categories and sub-categories by display order', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['31', ARS_ACCOUNT_ID, '200'], ['21', ARS_ACCOUNT_ID, '600']]),
            month(2026, 9, [['12', ARS_ACCOUNT_ID, '500'], ['11', ARS_ACCOUNT_ID, '1000']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.expense.categories.map(c => c.name)).toEqual(['Casa', 'Entretenimiento']);
        expect(data.income.categories[0]!.subCategories.map(c => c.name)).toEqual(['Salario', 'Bonos']);
    });

    it('should compute the net of every month as income minus expense', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 9, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '2000']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(amounts(data.net.monthlyAmounts)).toEqual(['500', '-500', '0', '0', '0']);
        expect(data.net.totalAmount.toString()).toEqual('0');
    });

    it('should accumulate the net month after month', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 9, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 10, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 11, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(amounts(data.net.monthlyAmounts)).toEqual(['500', '500', '500', '500', '0']);
        expect(amounts(data.accumulated.monthlyAmounts)).toEqual(['500', '1000', '1500', '2000', '2000']);
    });

    // the accumulated row is already a running total, so summing it across the columns would be
    // meaningless: its total column is the value it reached at the end of the period
    it('should use the value of the last month as the total of the accumulated row', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 9, [['11', ARS_ACCOUNT_ID, '1500'], ['21', ARS_ACCOUNT_ID, '1000']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.accumulated.totalAmount.toString()).toEqual('1000');
        expect(data.net.totalAmount.toString()).toEqual('1000');
    });

    it('should total every row across the months of the period', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 9, [['11', ARS_ACCOUNT_ID, '1000']]),
            month(2026, 10, [['11', ARS_ACCOUNT_ID, '1000']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.income.categories[0]!.subCategories[0]!.totalAmount.toString()).toEqual('3000');
        expect(data.income.totalAmount.toString()).toEqual('3000');
    });

    it('should reproduce the example of the technical spec', () => {
        const months = getProjectionMonths('2026-08', '2026-11', 202608);
        const monthlyItems: [string, string, string][] = [
            ['11', ARS_ACCOUNT_ID, '1000'],
            ['12', ARS_ACCOUNT_ID, '500'],
            ['31', ARS_ACCOUNT_ID, '200'],
            ['21', ARS_ACCOUNT_ID, '600']
        ];

        const data = buildProjectionTableData([
            month(2026, 8, monthlyItems),
            month(2026, 9, monthlyItems),
            month(2026, 10, monthlyItems),
            month(2026, 11, monthlyItems)
        ], months, IDENTITY_OPTIONS);

        expect(amounts(data.income.monthlyAmounts)).toEqual(['1500', '1500', '1500', '1500']);
        expect(data.income.totalAmount.toString()).toEqual('6000');

        expect(amounts(data.expense.monthlyAmounts)).toEqual(['800', '800', '800', '800']);
        expect(data.expense.totalAmount.toString()).toEqual('3200');

        expect(amounts(data.net.monthlyAmounts)).toEqual(['700', '700', '700', '700']);
        expect(amounts(data.accumulated.monthlyAmounts)).toEqual(['700', '1400', '2100', '2800']);
        expect(data.accumulated.totalAmount.toString()).toEqual('2800');
    });

    it('should convert amounts of accounts held in another currency', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['21', ARS_ACCOUNT_ID, '600'], ['21', USD_ACCOUNT_ID, '10']])
        ], MONTHS_2026_08_TO_12, {
            ...IDENTITY_OPTIONS,
            convertAmount: (amount: string, accountId: string) => {
                const value = parseBigDecimal(amount);
                return accountId === USD_ACCOUNT_ID ? value.multiply(1000) : value;
            }
        });

        // 600 ARS + 10 USD at 1000 = 10600 in the default currency
        expect(amounts(data.expense.categories[0]!.subCategories[0]!.monthlyAmounts)).toEqual(['10600', '0', '0', '0', '0']);
        expect(data.hasUnconvertedAmounts).toEqual(false);
    });

    // an amount that cannot be converted is left out rather than added raw, which would mix
    // currencies, and the flag lets the table say so
    it('should leave out amounts it cannot convert and report it', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['21', ARS_ACCOUNT_ID, '600'], ['21', USD_ACCOUNT_ID, '10']])
        ], MONTHS_2026_08_TO_12, {
            ...IDENTITY_OPTIONS,
            convertAmount: (amount: string, accountId: string) => accountId === USD_ACCOUNT_ID ? null : parseBigDecimal(amount)
        });

        expect(amounts(data.expense.categories[0]!.subCategories[0]!.monthlyAmounts)).toEqual(['600', '0', '0', '0', '0']);
        expect(data.hasUnconvertedAmounts).toEqual(true);
    });

    it('should ignore transfer categories', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1000'], ['41', ARS_ACCOUNT_ID, '9999']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.income.totalAmount.toString()).toEqual('1000');
        expect(data.expense.totalAmount.toString()).toEqual('0');
        expect(data.income.categories.length).toEqual(1);
    });

    it('should ignore unknown categories and months outside the period', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['11', ARS_ACCOUNT_ID, '1000'], ['999', ARS_ACCOUNT_ID, '9999']]),
            month(2026, 7, [['11', ARS_ACCOUNT_ID, '7777']]),
            month(2027, 1, [['11', ARS_ACCOUNT_ID, '8888']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.income.totalAmount.toString()).toEqual('1000');
    });

    it('should add up several accounts of the same sub-category in the same month', () => {
        const data = buildProjectionTableData([
            month(2026, 8, [['21', ARS_ACCOUNT_ID, '600'], ['21', USD_ACCOUNT_ID, '400']])
        ], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(amounts(data.expense.categories[0]!.subCategories[0]!.monthlyAmounts)).toEqual(['1000', '0', '0', '0', '0']);
    });

    it('should return empty sections and zeroed rows when there is no data', () => {
        const data = buildProjectionTableData([], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.income.categories).toEqual([]);
        expect(data.expense.categories).toEqual([]);
        expect(amounts(data.net.monthlyAmounts)).toEqual(['0', '0', '0', '0', '0']);
        expect(amounts(data.accumulated.monthlyAmounts)).toEqual(['0', '0', '0', '0', '0']);
        expect(data.accumulated.totalAmount.toString()).toEqual('0');
        expect(data.hasUnconvertedAmounts).toEqual(false);
    });

    it('should carry the month columns through to the table data', () => {
        const data = buildProjectionTableData([], MONTHS_2026_08_TO_12, IDENTITY_OPTIONS);

        expect(data.months.map(m => m.yearMonth)).toEqual([202608, 202609, 202610, 202611, 202612]);
        expect(data.months[0]!.kind).toEqual(ProjectionMonthKind.Mixed);
        expect(data.months[1]!.kind).toEqual(ProjectionMonthKind.Projected);
    });
});
