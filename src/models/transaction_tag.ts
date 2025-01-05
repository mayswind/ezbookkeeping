export interface TransactionTagCreateRequest {
    readonly name: string;
}

export interface TransactionTagModifyRequest {
    readonly id: string;
    readonly name: string;
}

export interface TransactionTagHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface TransactionTagMoveRequest {
    readonly newDisplayOrders: TransactionTagNewDisplayOrderRequest[];
}

export interface TransactionTagNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface TransactionTagDeleteRequest {
    readonly id: string;
}

export interface TransactionTagInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly displayOrder: number;
    readonly hidden: boolean;
}
