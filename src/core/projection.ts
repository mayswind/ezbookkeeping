import type { BigDecimal } from '@/core/numeral.ts';
import type { CategoryType } from '@/core/category.ts';

// How much of a projected month is made of transactions that already happened. A projection is split
// at the current instant, so the month containing it is the only mixed one.
export enum ProjectionMonthKind {
    Actual = 'actual',
    Mixed = 'mixed',
    Projected = 'projected'
}

export interface ProjectionMonth {
    readonly year: number;
    readonly month: number; // 1-based (1 = January, 12 = December)
    readonly yearMonth: number; // year * 100 + month, the key the api sorts by
    readonly kind: ProjectionMonthKind;
}

// A row of the projection table: one amount per month column plus the total of the period.
export interface ProjectionRow {
    readonly id: string;
    readonly name: string;
    readonly monthlyAmounts: BigDecimal[];
    readonly totalAmount: BigDecimal;
}

// A category row, collapsible into its sub-categories. Its own amounts are the subtotal of them.
export interface ProjectionCategoryRow extends ProjectionRow {
    readonly subCategories: ProjectionRow[];
}

// An income or expense section, collapsible into its categories. Its own amounts are the total of
// the section, which stays visible when the section is collapsed.
export interface ProjectionSectionRow extends ProjectionRow {
    readonly categoryType: CategoryType;
    readonly categories: ProjectionCategoryRow[];
}

export interface ProjectionTableData {
    readonly months: ProjectionMonth[];
    readonly income: ProjectionSectionRow;
    readonly expense: ProjectionSectionRow;
    // income minus expense, month by month
    readonly net: ProjectionRow;
    // running sum of net. Its total column is the value of the last month rather than the sum of the
    // row, because the row is already accumulated.
    readonly accumulated: ProjectionRow;
    // true when at least one amount was left out because its account currency could not be converted
    // to the default currency of the user, so the totals are known to be incomplete
    readonly hasUnconvertedAmounts: boolean;
}
