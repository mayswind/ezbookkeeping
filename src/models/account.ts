export interface AccountCreateRequest {
    readonly name: string;
    readonly category: number;
    readonly type: number;
    readonly icon: string;
    readonly color: string;
    readonly currency: string;
    readonly balance: number;
    readonly balanceTime: number;
    readonly comment: string;
    readonly creditCardStatementDate: number;
    readonly subAccounts?: AccountCreateRequest[];
    readonly clientSessionId: string;
}

export interface AccountModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly category: number;
    readonly icon: string;
    readonly color: string;
    readonly comment: string;
    readonly creditCardStatementDate?: number;
    readonly hidden: boolean;
    readonly subAccounts?: AccountModifyRequest[];
}

export interface AccountInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly parentId: string;
    readonly category: number;
    readonly type: number;
    readonly icon: string;
    readonly color: string;
    readonly currency: string;
    readonly balance: number;
    readonly comment: string;
    readonly creditCardStatementDate?: number;
    readonly displayOrder: number;
    readonly isAsset?: boolean;
    readonly isLiability?: boolean;
    readonly hidden: boolean;
    readonly subAccounts?: AccountInfoResponse[];
}

export interface AccountHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface AccountMoveRequest {
    readonly newDisplayOrders: AccountNewDisplayOrderRequest[];
}

export interface AccountNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface AccountDeleteRequest {
    readonly id: string;
}
