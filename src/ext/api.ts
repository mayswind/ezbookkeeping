import axios from 'axios';

import type { ApiResponse } from '@/core/api.ts';

import type {
    BusinessInfo,
    StaffInfo,
    AuditEntry,
    LocationInfo,
    ItemInfo,
    ItemRequest,
    StockLevel,
    StockChangeRequest,
    StockTransferRequest,
    StockMovement,
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
    listStaff: () => get<StaffInfo[]>('staff/list.json'),
    inviteStaff: (email: string, role: BusinessRole) => post<StaffInfo>('staff/invite.json', { email, role }),
    setStaffRole: (staffUid: string, role: BusinessRole) => post<boolean>('staff/set_role.json', { staffUid, role }),
    removeStaff: (staffUid: string) => post<boolean>('staff/remove.json', { staffUid }),
    respondToInvitation: (ownerUid: string, accept: boolean) => post<boolean>('staff/respond.json', { ownerUid, accept }),
    leaveBusiness: (ownerUid: string) => post<boolean>('staff/leave.json', { ownerUid }),
    listAudit: (beforeId?: string) => get<AuditEntry[]>('audit/list.json', { limit: 50, beforeId }),

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
    listStockMovements: (itemId?: string, beforeId?: string) => get<StockMovement[]>('stock/movements/list.json', { itemId, beforeId, limit: 50 })
};
