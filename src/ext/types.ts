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
    readonly action: string;
    readonly entityType: string;
    readonly entityId: string;
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

export interface CustomerInfo {
    readonly id: string;
    readonly name: string;
    readonly phone: string;
    readonly email: string;
    readonly note: string;
    readonly outstanding: number;
}

export interface CustomerRequest {
    readonly id?: string;
    readonly name: string;
    readonly phone?: string;
    readonly email?: string;
    readonly note?: string;
}

export interface SaleLineInfo {
    readonly itemId: string;
    readonly qty: number;
    readonly unitPrice: number;
    readonly lineTotal: number;
}

export interface SaleInfo {
    readonly id: string;
    readonly locationId: string;
    readonly customerId: string;
    readonly time: number;
    readonly subtotal: number;
    readonly discount: number;
    readonly total: number;
    readonly paid: number;
    readonly outstanding: number;
    readonly voided: boolean;
    readonly paymentAccountId: string;
    readonly receivableAccountId: string;
    readonly categoryId: string;
    readonly note: string;
    readonly actorUid: string;
    readonly lines?: SaleLineInfo[];
}

export interface SaleLineRequest {
    readonly itemId: string;
    readonly qty: number;
    readonly unitPrice?: number;
}

export interface SaleRequest {
    readonly locationId: string;
    readonly customerId: string;
    readonly time: number;
    readonly utcOffset: number;
    readonly lines: SaleLineRequest[];
    readonly discount: number;
    readonly amountPaid: number;
    readonly paymentAccountId: string;
    readonly receivableAccountId: string;
    readonly categoryId: string;
    readonly note: string;
}

export interface RepaymentInfo {
    readonly id: string;
    readonly customerId: string;
    readonly amount: number;
    readonly transactionId: string;
    readonly paymentAccountId: string;
    readonly receivableAccountId: string;
    readonly time: number;
    readonly note: string;
    readonly actorUid: string;
    readonly allocations?: { readonly saleId: string, readonly amount: number }[];
}

export interface RepaymentRequest {
    readonly customerId: string;
    readonly amount: number;
    readonly saleId: string;
    readonly time: number;
    readonly utcOffset: number;
    readonly paymentAccountId: string;
    readonly receivableAccountId: string;
    readonly categoryId: string;
    readonly note: string;
}

export interface PersonInfo {
    readonly uid: string;
    readonly name: string;
    readonly role: BusinessRole;
    readonly active: boolean;
}

export interface StockValueRow {
    readonly item: ItemInfo;
    readonly qty: number;
    readonly costValue: number;
    readonly retailValue: number;
    readonly locations?: { readonly locationId: string, readonly qty: number }[];
}

export interface StockValueReport {
    readonly rows: StockValueRow[];
    readonly totalCostValue: number;
    readonly totalRetailValue: number;
}

export interface LowStockRow {
    readonly item: ItemInfo;
    readonly qty: number;
    readonly shortfall: number;
}

export interface ReceivableRow {
    readonly customer: CustomerInfo;
    readonly outstanding: number;
    readonly current: number;
    readonly days31To60: number;
    readonly days61To90: number;
    readonly over90: number;
    readonly oldestSaleTime: number;
    readonly openSales: number;
}

export interface ReceivablesReport {
    readonly rows: ReceivableRow[];
    readonly totalOutstanding: number;
    readonly current: number;
    readonly days31To60: number;
    readonly days61To90: number;
    readonly over90: number;
}

export interface BusinessProfileInfo {
    readonly receiptName: string;
    readonly name: string; // the receipt name, or the owner's name when none is set
    readonly address: string;
    readonly phone: string;
    readonly footer: string;
}

export interface BusinessProfileRequest {
    readonly receiptName: string;
    readonly address: string;
    readonly phone: string;
    readonly footer: string;
}

export interface MySettings {
    readonly businessFeatures: boolean;
    readonly configured: boolean;
}

export interface ExtApiError {
    readonly errorCode?: number;
    readonly errorMessage?: string;
}
