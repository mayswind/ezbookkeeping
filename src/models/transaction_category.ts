import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';
import { DEFAULT_CATEGORY_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';

export class TransactionCategory implements TransactionCategoryInfoResponse {
    public id: string;
    public name: string;
    public parentId: string;
    public type: number;
    public icon: string;
    public color: string;
    public comment: string;
    public displayOrder: number;
    public visible: boolean;
    public secondaryCategories?: TransactionCategory[];

    private constructor(id: string, name: string, parentId: string, type: number, icon: string, color: string, comment: string, displayOrder: number, visible: boolean, secondaryCategories?: TransactionCategory[]) {
        this.id = id;
        this.name = name;
        this.parentId = parentId;
        this.type = type;
        this.icon = icon;
        this.color = color;
        this.comment = comment;
        this.displayOrder = displayOrder;
        this.visible = visible;

        if (secondaryCategories) {
            this.secondaryCategories = secondaryCategories;
        } else if (!secondaryCategories && (!parentId || parentId === '0')) {
            this.secondaryCategories = [];
        }
    }

    get hidden(): boolean {
        return !this.visible;
    }

    get subCategories(): TransactionCategoryInfoResponse[] | undefined {
        if (typeof(this.secondaryCategories) === 'undefined') {
            return undefined;
        }

        const ret: TransactionCategoryInfoResponse[] = [];

        if (this.secondaryCategories) {
            for (const subCategory of this.secondaryCategories) {
                ret.push(subCategory);
            }
        }

        return ret;
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
            categoryResponse.subCategories ? TransactionCategory.ofMany(categoryResponse.subCategories) : undefined
        );
    }

    public static ofMany(categoryResponses: TransactionCategoryInfoResponse[]): TransactionCategory[] {
        const tags: TransactionCategory[] = [];

        for (const tagResponse of categoryResponses) {
            tags.push(TransactionCategory.of(tagResponse));
        }

        return tags;
    }

    public static ofMap(categoriesByType: Record<number, TransactionCategoryInfoResponse[]>): Record<number, TransactionCategory[]> {
        const ret: Record<number, TransactionCategory[]> = {};

        for (const categoryType in categoriesByType) {
            if (!Object.prototype.hasOwnProperty.call(categoriesByType, categoryType)) {
                continue;
            }

            ret[categoryType] = TransactionCategory.ofMany(categoriesByType[categoryType]);
        }

        return ret;
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
    readonly categories: LocalizedPresetCategory[];
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
