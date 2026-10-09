import axios from 'axios';

import type { ApiResponse } from '@/core/api.ts';

import type {
    BusinessInfo,
    MySettings,
    PersonInfo,
    StockValueReport,
    LowStockRow,
    ReceivablesReport,
    StaffInfo,
    AuditEntry,
    LocationInfo,
    ItemInfo,
    ItemRequest,
    StockLevel,
    StockChangeRequest,
    StockTransferRequest,
    StockMovement,
    CustomerInfo,
    CustomerRequest,
    SaleInfo,
    SaleRequest,
    RepaymentInfo,
    RepaymentRequest,
    BusinessRole,
    ExtApiError
} from './types.ts';

// The shared axios instance already carries the base url, the login token and the error handling of the app;
// the business header is added by business.ts.
const BASE = 'v1/ext/';

async function get<T>(path: string, params?: Record<string, string | number | boolean | undefined>): Promise<T> {
    const response = await axios.get<ApiResponse<T>>(BASE + path, { params });
    return response.data.result;
}

async function post<T>(path: string, body: object): Promise<T> {
    const response = await axios.post<ApiResponse<T>>(BASE + path, body);
    return response.data.result;
}

/** Turns whatever a failed request threw into a message a person can read. */
export function describeError(error: unknown): string {
    const data = (error as { response?: { data?: ExtApiError } } | undefined)?.response?.data;

    if (data?.errorMessage) {
        return data.errorMessage.charAt(0).toUpperCase() + data.errorMessage.substring(1);
    }

    if ((error as { processed?: boolean } | undefined)?.processed) {
        return '';
    }

    return (error as { message?: string } | undefined)?.message || 'Something went wrong';
}

export default {
    // business membership (always about the logged-in person)
    getMyBusinesses: () => get<BusinessInfo[]>('me/businesses.json'),
    getMySettings: () => get<MySettings>('me/settings.json'),
    updateMySettings: (businessFeatures: boolean) => post<MySettings>('me/settings/update.json', { businessFeatures }),
    listStaff: () => get<StaffInfo[]>('staff/list.json'),
    inviteStaff: (email: string, role: BusinessRole) => post<StaffInfo>('staff/invite.json', { email, role }),
    setStaffRole: (staffUid: string, role: BusinessRole) => post<boolean>('staff/set_role.json', { staffUid, role }),
    removeStaff: (staffUid: string) => post<boolean>('staff/remove.json', { staffUid }),
    respondToInvitation: (ownerUid: string, accept: boolean) => post<boolean>('staff/respond.json', { ownerUid, accept }),
    leaveBusiness: (ownerUid: string) => post<boolean>('staff/leave.json', { ownerUid }),
    listAudit: (beforeId?: string) => get<AuditEntry[]>('audit/list.json', { limit: 50, beforeId }),

    // who works in the business, for showing who did what
    listPeople: () => get<PersonInfo[]>('people/list.json'),
    listMyPeople: () => get<PersonInfo[]>('staff/people.json'), // always my own business, even while working in another

    // reports
    getStockValueReport: (locationId?: string) => get<StockValueReport>('reports/stock_value.json', { locationId }),
    getLowStockReport: (locationId?: string) => get<LowStockRow[]>('reports/low_stock.json', { locationId }),
    getReceivablesReport: () => get<ReceivablesReport>('reports/receivables.json'),

    // locations
    listLocations: () => get<LocationInfo[]>('locations/list.json'),
    addLocation: (name: string) => post<LocationInfo>('locations/add.json', { name }),
    renameLocation: (id: string, name: string) => post<LocationInfo>('locations/modify.json', { id, name }),
    deleteLocation: (id: string) => post<boolean>('locations/delete.json', { id }),

    // items and stock
    listItems: () => get<ItemInfo[]>('items/list.json'),
    addItem: (req: ItemRequest) => post<ItemInfo>('items/add.json', req),
    modifyItem: (req: ItemRequest) => post<ItemInfo>('items/modify.json', req),
    deleteItem: (id: string) => post<boolean>('items/delete.json', { id }),
    getStockLevels: () => get<StockLevel[]>('items/stock.json'),
    receiveStock: (req: StockChangeRequest) => post<StockMovement>('stock/receive.json', req),
    adjustStock: (req: StockChangeRequest) => post<StockMovement>('stock/adjust.json', req),
    transferStock: (req: StockTransferRequest) => post<boolean>('stock/transfer.json', req),
    // customers and sales
    listCustomers: () => get<CustomerInfo[]>('customers/list.json'),
    addCustomer: (req: CustomerRequest) => post<CustomerInfo>('customers/add.json', req),
    modifyCustomer: (req: CustomerRequest) => post<CustomerInfo>('customers/modify.json', req),
    deleteCustomer: (id: string) => post<boolean>('customers/delete.json', { id }),
    listSales: (options: { beforeId?: string, customerId?: string, onlyOpen?: boolean } = {}) =>
        get<SaleInfo[]>('sales/list.json', { beforeId: options.beforeId, customerId: options.customerId, onlyOpen: options.onlyOpen, limit: 200 }),
    listRepayments: (customerId?: string, beforeId?: string) => get<RepaymentInfo[]>('repayments/list.json', { customerId, beforeId, limit: 100 }),
    addRepayment: (req: RepaymentRequest) => post<RepaymentInfo>('repayments/add.json', req),
    getSale: (id: string) => get<SaleInfo>('sales/get.json', { id }),
    createSale: (req: SaleRequest) => post<SaleInfo>('sales/add.json', req),
    voidSale: (id: string) => post<boolean>('sales/void.json', { id }),

    listStockMovements: (itemId?: string, beforeId?: string) => get<StockMovement[]>('stock/movements/list.json', { itemId, beforeId, limit: 50 })
};
