import type { BusinessProfileInfo, CustomerInfo, ItemInfo, RepaymentInfo, SaleInfo } from './types.ts';
import { formatQty } from './qty.ts';

// Receipts are built as plain data first, then shown on screen, printed, or turned into text for messaging apps.
// Keeping the arithmetic and wording here (and not in a component) means it can be tested.

export interface ReceiptRow {
    readonly label: string;
    readonly value: string;
    readonly strong?: boolean;
}

export interface ReceiptLine {
    readonly name: string;
    readonly detail: string; // "2.5 kg x 1,500.00"
    readonly total: string;
}

export interface Receipt {
    readonly kind: 'sale' | 'repayment';
    readonly business: { readonly name: string, readonly address: string, readonly phone: string, readonly footer: string };
    readonly title: string;
    readonly number: string;
    readonly date: string;
    readonly facts: ReceiptRow[];
    readonly lines: ReceiptLine[];
    readonly totals: ReceiptRow[];
    readonly notes: string[];
    readonly voided: boolean;
}

/** The translation and formatting the builders need; the pages pass in the app's own. */
export interface ReceiptFormat {
    tt: (key: string, params?: Record<string, string | number>) => string;
    money: (minorUnits: number) => string;
    date: (unixTime: number) => string;
}

function business(profile: BusinessProfileInfo): Receipt['business'] {
    return { name: profile.name, address: profile.address, phone: profile.phone, footer: profile.footer };
}

export interface SaleReceiptInput {
    readonly sale: SaleInfo; // with its lines
    readonly items: ItemInfo[];
    readonly customer?: CustomerInfo;
    readonly servedBy: string;
    readonly locationName?: string; // only for businesses with several locations
    readonly paidIntoName?: string;
    readonly profile: BusinessProfileInfo;
}

export function buildSaleReceipt(input: SaleReceiptInput, fmt: ReceiptFormat): Receipt {
    const { sale, items, customer, profile } = input;
    const { tt, money } = fmt;

    const facts: ReceiptRow[] = [];

    if (customer) {
        facts.push({ label: tt('Customer'), value: customer.name });
    }

    if (input.locationName) {
        facts.push({ label: tt('Location'), value: input.locationName });
    }

    if (input.servedBy) {
        facts.push({ label: tt('Served by'), value: input.servedBy });
    }

    const lines: ReceiptLine[] = (sale.lines ?? []).map(line => {
        const item = items.find(i => i.id === line.itemId);
        const unit = item?.unit ? ` ${item.unit}` : '';

        return {
            name: item?.name ?? tt('Item no longer listed'),
            detail: `${formatQty(line.qty)}${unit} x ${money(line.unitPrice)}`,
            total: money(line.lineTotal)
        };
    });

    const totals: ReceiptRow[] = [];

    if (sale.discount > 0) {
        totals.push({ label: tt('Subtotal'), value: money(sale.subtotal) });
        totals.push({ label: tt('Discount'), value: `-${money(sale.discount)}` });
    }

    totals.push({ label: tt('Total'), value: money(sale.total), strong: true });
    totals.push({ label: tt('Paid'), value: money(sale.paid) });

    if (sale.outstanding > 0) {
        totals.push({ label: tt('Balance owed'), value: money(sale.outstanding), strong: true });
    }

    if (sale.paid > 0 && input.paidIntoName) {
        facts.push({ label: tt('Paid into'), value: input.paidIntoName });
    }

    const notes: string[] = [];

    if (sale.voided) {
        notes.push(tt('This sale was cancelled.'));
    } else if (sale.outstanding > 0) {
        notes.push(tt('Part or all of this sale is on credit.'));
    }

    return {
        kind: 'sale', business: business(profile), title: tt('Receipt'), number: `#${sale.id}`, date: fmt.date(sale.time),
        facts, lines, totals, notes, voided: sale.voided
    };
}

export interface RepaymentReceiptInput {
    readonly repayment: RepaymentInfo; // with its allocations
    readonly customer?: CustomerInfo;
    readonly servedBy: string;
    readonly paidIntoName?: string;
    readonly profile: BusinessProfileInfo;
}

export function buildRepaymentReceipt(input: RepaymentReceiptInput, fmt: ReceiptFormat): Receipt {
    const { repayment, customer, profile } = input;
    const { tt, money } = fmt;

    const facts: ReceiptRow[] = [];

    if (customer) {
        facts.push({ label: tt('Customer'), value: customer.name });
    }

    if (input.servedBy) {
        facts.push({ label: tt('Received by'), value: input.servedBy });
    }

    if (input.paidIntoName) {
        facts.push({ label: tt('Paid into'), value: input.paidIntoName });
    }

    const lines: ReceiptLine[] = (repayment.allocations ?? []).map(allocation => ({
        name: tt('Payment towards sale #{id}', { id: allocation.saleId }),
        detail: '',
        total: money(allocation.amount)
    }));

    const totals: ReceiptRow[] = [{ label: tt('Amount received'), value: money(repayment.amount), strong: true }];

    if (customer) {
        totals.push({ label: tt('Still owed'), value: money(customer.outstanding) });
    }

    return {
        kind: 'repayment', business: business(profile), title: tt('Payment receipt'), number: `#${repayment.id}`, date: fmt.date(repayment.time),
        facts, lines, totals, notes: repayment.note ? [repayment.note] : [], voided: false
    };
}

// ---- plain text, for WhatsApp, SMS or email

function wrap(text: string, width: number): string[] {
    const out: string[] = [];
    let line = '';

    for (const word of text.split(/\s+/).filter(Boolean)) {
        if (word.length > width) {
            if (line) {
                out.push(line);
                line = '';
            }

            for (let i = 0; i < word.length; i += width) {
                out.push(word.substring(i, i + width));
            }
        } else if (!line) {
            line = word;
        } else if (line.length + 1 + word.length <= width) {
            line += ' ' + word;
        } else {
            out.push(line);
            line = word;
        }
    }

    if (line) {
        out.push(line);
    }

    return out;
}

function pair(left: string, right: string, width: number): string[] {
    if (left.length + 1 + right.length <= width) {
        return [left + ' '.repeat(width - left.length - right.length) + right];
    }

    return [...wrap(left, width), ' '.repeat(Math.max(0, width - right.length)) + right];
}

/** A receipt as text lines of at most `width` characters, with amounts lined up on the right. */
export function receiptToText(receipt: Receipt, width: number = 32): string {
    const rule = '-'.repeat(width);
    const out: string[] = [];
    const centered = (text: string): void => {
        for (const line of wrap(text, width)) {
            out.push(' '.repeat(Math.max(0, Math.floor((width - line.length) / 2))) + line);
        }
    };

    centered(receipt.business.name);

    if (receipt.business.address) {
        centered(receipt.business.address);
    }

    if (receipt.business.phone) {
        centered(receipt.business.phone);
    }

    out.push(rule);

    if (receipt.voided) {
        centered('*** VOID ***');
    }

    out.push(...pair(receipt.title, receipt.number, width));
    out.push(receipt.date);

    for (const fact of receipt.facts) {
        out.push(...wrap(`${fact.label}: ${fact.value}`, width));
    }

    if (receipt.lines.length > 0) {
        out.push(rule);

        for (const line of receipt.lines) {
            out.push(...wrap(line.name, width));
            out.push(...pair(line.detail ? '  ' + line.detail : '', line.total, width));
        }
    }

    out.push(rule);

    for (const row of receipt.totals) {
        out.push(...pair(row.label, row.value, width));
    }

    for (const note of receipt.notes) {
        out.push(rule);
        out.push(...wrap(note, width));
    }

    if (receipt.business.footer) {
        out.push(rule);
        centered(receipt.business.footer);
    }

    return out.join('\n');
}
