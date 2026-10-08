// Types of the ext module API (see pkg/ext/api/types.go). Ids are strings, amounts are integers in minor
// currency units and quantities are integers scaled by QTY_SCALE.

export type BusinessRole = 'owner' | 'manager' | 'staff';
export type MembershipStatus = 'owner' | 'active' | 'pending';

export interface BusinessInfo {
    readonly ownerUid: string;
    readonly name: string;
    readonly role: BusinessRole;
    readonly status: MembershipStatus;
}

export interface StaffInfo {
    readonly staffUid: string;
    readonly username: string;
    readonly nickname: string;
    readonly email: string;
    readonly role: BusinessRole;
    readonly status: 'active' | 'pending';
}

export interface AuditEntry {
    readonly id: string;
    readonly actorUid: string;
    readonly role: BusinessRole;
    readonly method: string;
    readonly path: string;
    readonly status: number;
    readonly time: number;
}

export interface LocationInfo {
    readonly id: string;
    readonly name: string;
    readonly isDefault: boolean;
}

export interface ItemInfo {
    readonly id: string;
    readonly sku: string;
    readonly name: string;
    readonly unit: string;
    readonly costPrice: number;
    readonly salePrice: number;
    readonly reorderLevel: number;
    readonly trackStock: boolean;
}

export interface ItemRequest {
    readonly id?: string;
    readonly sku: string;
    readonly name: string;
    readonly unit: string;
    readonly costPrice: number;
    readonly salePrice: number;
    readonly reorderLevel: number;
    readonly trackStock: boolean;
}

export interface StockLevel {
    readonly itemId: string;
    readonly locationId: string;
    readonly qty: number;
}

export interface StockChangeRequest {
    readonly itemId: string;
    readonly locationId: string;
    readonly qty: number;
    readonly unitCost?: number;
    readonly note?: string;
    readonly opening?: boolean;
}

export interface StockTransferRequest {
    readonly itemId: string;
    readonly fromLocationId: string;
    readonly toLocationId: string;
    readonly qty: number;
    readonly note?: string;
}

export interface StockMovement {
    readonly id: string;
    readonly itemId: string;
    readonly locationId: string;
    readonly qtyChange: number;
    readonly reason: number;
    readonly note: string;
    readonly actorUid: string;
    readonly time: number;
}

export interface ExtApiError {
    readonly errorCode?: number;
    readonly errorMessage?: string;
}
