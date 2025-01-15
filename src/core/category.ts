import type { ColorValue } from '@/core/color.ts';

export enum CategoryType {
    Income = 1,
    Expense = 2,
    Transfer = 3
}

export const ALL_CATEGORY_TYPES: CategoryType[] = [
    CategoryType.Income,
    CategoryType.Expense,
    CategoryType.Transfer
];

export interface PresetCategory {
    readonly name: string;
    readonly categoryIconId: string;
    readonly color: ColorValue;
    readonly subCategories: PresetSubCategory[];
}

export interface PresetSubCategory {
    readonly name: string;
    readonly categoryIconId: string;
    readonly color: ColorValue;
}

export interface LocalizedPresetCategory {
    readonly name: string;
    readonly type: CategoryType;
    readonly icon: string;
    readonly color: ColorValue;
    readonly subCategories: LocalizedPresetSubCategory[];
}

export interface LocalizedPresetSubCategory {
    readonly name: string;
    readonly type: CategoryType;
    readonly icon: string;
    readonly color: ColorValue;
}
