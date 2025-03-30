import type { ColorValue } from '@/core/color.ts';
import { CategoryType } from '@/core/category.ts';
import { DEFAULT_CATEGORY_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';

export class TransactionCategory implements TransactionCategoryInfoResponse {
    public id: string;
    public name: string;
    public parentId: string;
    public type: CategoryType;
    public icon: string;
    public color: ColorValue;
    public comment: string;
    public displayOrder: number;
    public visible: boolean;
    public subCategories?: TransactionCategory[];

    private constructor(id: string, name: string, parentId: string, type: CategoryType, icon: string, color: ColorValue, comment: string, displayOrder: number, visible: boolean, subCategories?: TransactionCategory[]) {
        this.id = id;
        this.name = name;
        this.parentId = parentId;
        this.type = type;
        this.icon = icon;
        this.color = color;
        this.comment = comment;
        this.displayOrder = displayOrder;
        this.visible = visible;

        if (subCategories) {
            this.subCategories = subCategories;
        } else if (!subCategories && (!parentId || parentId === '0')) {
            this.subCategories = [];
        }
    }

    public get hidden(): boolean {
        return !this.visible;
    }

    public equals(other: TransactionCategory): boolean {
        const isEqual = this.id === other.id &&
            this.name === other.name &&
            this.parentId === other.parentId &&
            this.type === other.type &&
            this.icon === other.icon &&
            this.color === other.color &&
            this.comment === other.comment &&
            this.displayOrder === other.displayOrder &&
            this.visible === other.visible;

        if (!isEqual) {
            return false;
        }

        if (this.subCategories && other.subCategories) {
            if (this.subCategories.length !== other.subCategories.length) {
                return false;
            }

            for (let i = 0; i < this.subCategories.length; i++) {
                if (!this.subCategories[i].equals(other.subCategories[i])) {
                    return false;
                }
            }
        } else if ((this.subCategories && this.subCategories.length) || (other.subCategories && other.subCategories.length)) {
            return false;
        }

        return true;
    }

    public fillFrom(other: TransactionCategory): void {
        this.id = other.id;
        this.name = other.name;
        this.parentId = other.parentId;
        this.type = other.type;
        this.icon = other.icon;
        this.color = other.color;
        this.comment = other.comment;
        this.visible = other.visible;
    }

    public toCreateRequest(clientSessionId: string): TransactionCategoryCreateRequest {
        return {
            name: this.name,
            type: this.type,
            parentId: this.parentId,
            icon: this.icon,
            color: this.color,
            comment: this.comment,
            clientSessionId: clientSessionId
        };
    }

    public toModifyRequest(): TransactionCategoryModifyRequest {
        return {
            id: this.id,
            name: this.name,
            parentId: this.parentId,
            icon: this.icon,
            color: this.color,
            comment: this.comment,
            hidden: !this.visible
        };
    }

    public static of(categoryResponse: TransactionCategoryInfoResponse): TransactionCategory {
        return new TransactionCategory(
            categoryResponse.id,
            categoryResponse.name,
            categoryResponse.parentId,
            categoryResponse.type,
            categoryResponse.icon,
            categoryResponse.color,
            categoryResponse.comment,
            categoryResponse.displayOrder,
            !categoryResponse.hidden,
            categoryResponse.subCategories ? TransactionCategory.ofMulti(categoryResponse.subCategories) : undefined
        );
    }

    public static ofMulti(categoryResponses: TransactionCategoryInfoResponse[]): TransactionCategory[] {
        const categories: TransactionCategory[] = [];

        for (const categoryResponse of categoryResponses) {
            categories.push(TransactionCategory.of(categoryResponse));
        }

        return categories;
    }

    public static ofMap(categoriesByType: Record<number, TransactionCategoryInfoResponse[]>): Record<number, TransactionCategory[]> {
        const ret: Record<number, TransactionCategory[]> = {};

        for (const categoryType in categoriesByType) {
            if (!Object.prototype.hasOwnProperty.call(categoriesByType, categoryType)) {
                continue;
            }

            ret[categoryType] = TransactionCategory.ofMulti(categoriesByType[categoryType]);
        }

        return ret;
    }

    public static findNameById(categories: TransactionCategory[], id: string): string | null {
        for (const category of categories) {
            if (category.id === id) {
                return category.name;
            }
        }

        return null;
    }

    public static createNewCategory(type?: CategoryType, parentId?: string): TransactionCategory {
        return new TransactionCategory('', '', parentId || '0', type || CategoryType.Income, DEFAULT_CATEGORY_ICON_ID, DEFAULT_CATEGORY_COLOR, '', 0, true);
    }
}

export interface TransactionCategoryCreateRequest {
    readonly name: string;
    readonly type: number;
    readonly parentId: string;
    readonly icon: string;
    readonly color: string;
    readonly comment: string;
    readonly clientSessionId: string;
}

export interface TransactionCategoryCreateBatchRequest {
    readonly categories: TransactionCategoryCreateWithSubCategories[];
}

export interface TransactionCategoryCreateWithSubCategories {
    readonly name: string;
    readonly type: CategoryType;
    readonly icon: string;
    readonly color: ColorValue;
    readonly subCategories: TransactionCategoryCreateRequest[];
}

export interface TransactionCategoryModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly parentId: string;
    readonly icon: string;
    readonly color: string;
    readonly comment: string;
    readonly hidden: boolean;
}

export interface TransactionCategoryHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface TransactionCategoryMoveRequest {
    readonly newDisplayOrders: TransactionCategoryNewDisplayOrderRequest[];
}

export interface TransactionCategoryNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface TransactionCategoryDeleteRequest {
    readonly id: string;
}

export interface TransactionCategoryInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly parentId: string;
    readonly type: number;
    readonly icon: string;
    readonly color: string;
    readonly comment: string;
    readonly displayOrder: number;
    readonly hidden: boolean;
    readonly subCategories?: TransactionCategoryInfoResponse[];
}

export interface TransactionCategoriesWithVisibleCount {
    readonly type: number;
    readonly allCategories: TransactionCategory[];
    readonly allVisibleCategoryCount: number;
    readonly firstVisibleCategoryIndex: number;
    readonly allSubCategories: Record<string, TransactionCategory[]>;
    readonly allVisibleSubCategoryCounts: Record<string, number>;
    readonly allFirstVisibleSubCategoryIndexes: Record<string, number>;
}
