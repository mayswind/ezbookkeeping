import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import moment from 'moment-timezone';

import { useRootStore } from '@/stores/index.ts';
import { useUserStore } from '@/stores/user.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { EMPTY_USER_BASIC_INFO } from '@/models/user.ts';

import services from '@/lib/services.ts';

// The stores read window.location and localStorage when they are loaded and the test environment is node
vi.hoisted(() => {
    const data = new Map<string, string>();
    const storage = {
        getItem: (key: string) => data.has(key) ? data.get(key) as string : null,
        setItem: (key: string, value: string) => { data.set(key, String(value)); },
        removeItem: (key: string) => { data.delete(key); },
        clear: () => { data.clear(); },
        key: (index: number) => Array.from(data.keys())[index] ?? null,
        get length() { return data.size; }
    };

    Object.assign(globalThis, {
        window: { location: { pathname: '/', origin: 'http://localhost', href: 'http://localhost/' }, localStorage: storage, sessionStorage: storage },
        localStorage: storage,
        sessionStorage: storage
    });
});

describe('overview store', () => {
    beforeEach(() => {
        moment.tz.setDefault('UTC');
        vi.useFakeTimers({ toFake: ['Date'] });
        vi.setSystemTime(new Date('2026-10-18T12:00:00Z')); // Sunday
        localStorage.clear();
        setActivePinia(createPinia());
    });

    afterEach(() => {
        vi.useRealTimers();
        moment.tz.setDefault();
    });

    it('uses the first day of week of a user known before the store is created', async () => {
        useUserStore().storeUserBasicInfo({ ...EMPTY_USER_BASIC_INFO, username: 'test', firstDayOfWeek: 1 });

        const overviewStore = useOverviewStore();
        const getTransactionAmounts = vi.spyOn(services, 'getTransactionAmounts').mockResolvedValue({
            data: { success: true, result: {} }
        } as never);

        await overviewStore.loadTransactionOverview({ force: false });

        const monday = moment.utc('2026-10-12T00:00:00Z').unix();
        const nextMonday = moment.utc('2026-10-19T00:00:00Z').unix();

        expect(getTransactionAmounts.mock.calls[0]?.[0].thisWeek).toEqual({ startTime: monday, endTime: nextMonday - 1 });
    });

    it('uses the first day of week of the user who logged in', async () => {
        // The login page creates the root store, and with it the overview store, before any user is known
        const rootStore = useRootStore();
        const overviewStore = useOverviewStore();

        vi.spyOn(services, 'authorize').mockResolvedValue({
            data: { success: true, result: { token: 'token', user: { username: 'test', firstDayOfWeek: 1 } } }
        } as never);
        const getTransactionAmounts = vi.spyOn(services, 'getTransactionAmounts').mockResolvedValue({
            data: { success: true, result: {} }
        } as never);

        await rootStore.authorize({ loginName: 'test', password: 'password' });
        await overviewStore.loadTransactionOverview({ force: false });

        const monday = moment.utc('2026-10-12T00:00:00Z').unix();
        const nextMonday = moment.utc('2026-10-19T00:00:00Z').unix();

        expect(getTransactionAmounts.mock.calls[0]?.[0].thisWeek).toEqual({ startTime: monday, endTime: nextMonday - 1 });
        expect(overviewStore.transactionDataRange.thisWeek).toEqual({ startTime: monday, endTime: nextMonday - 1 });
    });
});
