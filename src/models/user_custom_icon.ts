export interface UserCustomIconInfoResponse {
    readonly id: string;
    readonly displayOrder: number;
    readonly originalUrl: string;
}

export interface UserCustomIconListResponse {
    readonly items: UserCustomIconInfoResponse[];
    readonly expiresAt: number;
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
