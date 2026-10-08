// Remembered form choices, per business, so a till operator does not pick the same accounts for every sale.

export interface FormDefaults {
    payment?: string; // account the money goes into
    receivable?: string; // account that tracks what customers owe
    category?: string; // income category of sales
    transfer?: string; // transfer category of repayments
}

function key(businessKey: string): string {
    return `ebk_ext_sale_defaults_${businessKey}`;
}

export function loadFormDefaults(businessKey: string): FormDefaults {
    try {
        return JSON.parse(localStorage.getItem(key(businessKey)) ?? '{}') as FormDefaults;
    } catch {
        return {};
    }
}

/** Merges the given choices into what is already remembered. */
export function saveFormDefaults(businessKey: string, patch: FormDefaults): void {
    try {
        localStorage.setItem(key(businessKey), JSON.stringify({ ...loadFormDefaults(businessKey), ...patch }));
    } catch {
        // remembering choices is a convenience only
    }
}
