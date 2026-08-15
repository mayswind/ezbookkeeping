export interface UserCustomIconInfoResponse {
    readonly id: string;
    readonly displayOrder: number;
}

export interface UserCustomIconNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface UserCustomIconMoveRequest {
    readonly newDisplayOrders: UserCustomIconNewDisplayOrderRequest[];
}

export interface UserCustomIconDeleteRequest {
    readonly id: string;
}
