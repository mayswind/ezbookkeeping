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
    readonly type: number;
    readonly icon: string;
    readonly color: string;
    readonly comment: string;
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
