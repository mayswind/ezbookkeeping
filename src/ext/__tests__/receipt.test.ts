import { describe, expect, test } from 'vitest';

import { buildRepaymentReceipt, buildSaleReceipt, receiptToText, type ReceiptFormat } from '@/ext/receipt.ts';
import type { BusinessProfileInfo, CustomerInfo, ItemInfo, RepaymentInfo, SaleInfo } from '@/ext/types.ts';

// plain English and whole-number money keep the expectations readable
const fmt: ReceiptFormat = {
    tt: (key, params) => Object.entries(params ?? {}).reduce((text, [name, value]) => text.replace(`{${name}}`, String(value)), key),
    money: minor => `N${(minor / 100).toFixed(2)}`,
    date: () => '9 Oct 2026 10:30'
};

const profile: BusinessProfileInfo = { receiptName: 'Ada Stores', name: 'Ada Stores', address: '12 Market Road', phone: '0800 000', footer: 'Thank you' };
const rice: ItemInfo = { id: '1', sku: 'RICE', name: 'Rice 1kg', unit: 'kg', costPrice: 70000, salePrice: 150000, reorderLevel: 0, trackStock: true };
const oil: ItemInfo = { id: '2', sku: 'OIL', name: 'Groundnut oil', unit: '', costPrice: 200000, salePrice: 400000, reorderLevel: 0, trackStock: true };
const ada: CustomerInfo = { id: '9', name: 'Ada Obi', phone: '', email: '', note: '', outstanding: 275000 };

function sale(over: Partial<SaleInfo> = {}): SaleInfo {
    return {
        id: '45', locationId: '1', customerId: '9', time: 1, subtotal: 375000, discount: 0, total: 375000, paid: 100000, outstanding: 275000,
        voided: false, paymentAccountId: '5', receivableAccountId: '6', categoryId: '7', note: '', actorUid: '3',
        lines: [{ itemId: '1', qty: 2500, unitPrice: 150000, lineTotal: 375000 }], ...over
    };
}

describe('buildSaleReceipt', () => {
    test('a sale on credit shows what was paid and what is owed', () => {
        const receipt = buildSaleReceipt({ sale: sale(), items: [rice], customer: ada, servedBy: 'Sam', paidIntoName: 'Cash', profile }, fmt);

        expect(receipt.number).toBe('#45');
        expect(receipt.title).toBe('Receipt');
        expect(receipt.lines).toEqual([{ name: 'Rice 1kg', detail: '2.5 kg x N1500.00', total: 'N3750.00' }]);
        expect(receipt.totals.map(r => `${r.label}=${r.value}`)).toEqual(['Total=N3750.00', 'Paid=N1000.00', 'Balance owed=N2750.00']);
        expect(receipt.facts.map(r => `${r.label}=${r.value}`)).toEqual(['Customer=Ada Obi', 'Served by=Sam', 'Paid into=Cash']);
        expect(receipt.notes).toEqual(['Part or all of this sale is on credit.']);
        expect(receipt.voided).toBe(false);
    });

    test('a paid sale has no balance and no credit note', () => {
        const receipt = buildSaleReceipt({ sale: sale({ paid: 375000, outstanding: 0 }), items: [rice], servedBy: '', profile }, fmt);

        expect(receipt.totals.map(r => r.label)).toEqual(['Total', 'Paid']);
        expect(receipt.notes).toEqual([]);
        expect(receipt.facts.some(f => f.label === 'Customer' || f.label === 'Served by')).toBe(false);
    });

    test('a discount shows the subtotal and the discount before the total', () => {
        const receipt = buildSaleReceipt({
            sale: sale({ subtotal: 375000, discount: 25000, total: 350000, paid: 350000, outstanding: 0 }), items: [rice], servedBy: '', profile
        }, fmt);

        expect(receipt.totals.map(r => `${r.label}=${r.value}`)).toEqual(['Subtotal=N3750.00', 'Discount=-N250.00', 'Total=N3500.00', 'Paid=N3500.00']);
    });

    test('a voided sale says so', () => {
        const receipt = buildSaleReceipt({ sale: sale({ voided: true, outstanding: 0 }), items: [rice], servedBy: '', profile }, fmt);

        expect(receipt.voided).toBe(true);
        expect(receipt.notes).toEqual(['This sale was cancelled.']);
    });

    test('items without a unit, and items that were deleted since', () => {
        const receipt = buildSaleReceipt({
            sale: sale({ lines: [{ itemId: '2', qty: 1000, unitPrice: 400000, lineTotal: 400000 }, { itemId: '404', qty: 3000, unitPrice: 100, lineTotal: 300 }] }),
            items: [oil], servedBy: '', profile
        }, fmt);

        expect(receipt.lines[0]).toEqual({ name: 'Groundnut oil', detail: '1 x N4000.00', total: 'N4000.00' });
        expect(receipt.lines[1]!.name).toBe('Item no longer listed');
    });

    test('the location is named only when asked for', () => {
        const withLocation = buildSaleReceipt({ sale: sale(), items: [rice], servedBy: '', locationName: 'Shop 2', profile }, fmt);
        const without = buildSaleReceipt({ sale: sale(), items: [rice], servedBy: '', profile }, fmt);

        expect(withLocation.facts).toContainEqual({ label: 'Location', value: 'Shop 2' });
        expect(without.facts.some(f => f.label === 'Location')).toBe(false);
    });
});

describe('buildRepaymentReceipt', () => {
    const repayment: RepaymentInfo = {
        id: '8', customerId: '9', amount: 350000, transactionId: '1', paymentAccountId: '5', receivableAccountId: '6', time: 1, note: 'cash on pickup', actorUid: '3',
        allocations: [{ saleId: '45', amount: 300000 }, { saleId: '46', amount: 50000 }]
    };

    test('lists what the money paid off and what is still owed', () => {
        const receipt = buildRepaymentReceipt({ repayment, customer: ada, servedBy: 'Sam', paidIntoName: 'Bank', profile }, fmt);

        expect(receipt.title).toBe('Payment receipt');
        expect(receipt.number).toBe('#8');
        expect(receipt.lines).toEqual([
            { name: 'Payment towards sale #45', detail: '', total: 'N3000.00' },
            { name: 'Payment towards sale #46', detail: '', total: 'N500.00' }
        ]);
        expect(receipt.totals.map(r => `${r.label}=${r.value}`)).toEqual(['Amount received=N3500.00', 'Still owed=N2750.00']);
        expect(receipt.facts).toContainEqual({ label: 'Received by', value: 'Sam' });
        expect(receipt.facts).toContainEqual({ label: 'Paid into', value: 'Bank' });
        expect(receipt.notes).toEqual(['cash on pickup']);
    });
});

describe('receiptToText', () => {
    const receipt = buildSaleReceipt({ sale: sale(), items: [rice], customer: ada, servedBy: 'Sam', paidIntoName: 'Cash', profile }, fmt);

    test('has the business, the lines, the totals and the footer in order', () => {
        const text = receiptToText(receipt, 32);
        const at = (needle: string): number => text.indexOf(needle);

        expect(text).toContain('Ada Stores');
        expect(at('Ada Stores')).toBeLessThan(at('Receipt'));
        expect(at('Receipt')).toBeLessThan(at('Rice 1kg'));
        expect(at('Rice 1kg')).toBeLessThan(at('Balance owed'));
        expect(at('Balance owed')).toBeLessThan(at('Thank you'));
    });

    test('no line is longer than the width, and amounts line up on the right edge', () => {
        for (const width of [24, 32, 42]) {
            const lines = receiptToText(receipt, width).split('\n');

            expect(Math.max(...lines.map(l => l.length)), `width ${width}`).toBeLessThanOrEqual(width);

            const total = lines.find(l => l.startsWith('Total'))!;
            expect(total.endsWith('N3750.00')).toBe(true);
            expect(total.length).toBe(width);
        }
    });

    test('long item names wrap instead of overflowing', () => {
        const long = buildSaleReceipt({
            sale: sale(), items: [{ ...rice, name: 'Extra long grain parboiled premium rice family pack' }], servedBy: '', profile
        }, fmt);

        const lines = receiptToText(long, 24).split('\n');
        expect(Math.max(...lines.map(l => l.length))).toBeLessThanOrEqual(24);
        expect(lines.join(' ')).toContain('premium rice family pack');
    });

    test('a voided sale is marked at the top', () => {
        const voided = buildSaleReceipt({ sale: sale({ voided: true }), items: [rice], servedBy: '', profile }, fmt);

        expect(receiptToText(voided)).toContain('*** VOID ***');
        expect(receiptToText(receipt)).not.toContain('VOID');
    });

    test('optional business details are left out when empty', () => {
        const bare = buildSaleReceipt({ sale: sale(), items: [rice], servedBy: '', profile: { receiptName: '', name: 'Ada', address: '', phone: '', footer: '' } }, fmt);
        const text = receiptToText(bare);

        expect(text).not.toContain('Market Road');
        expect(text.split('\n')[0]).toContain('Ada');
        expect(text.trimEnd().endsWith('credit.')).toBe(true);
    });
});
