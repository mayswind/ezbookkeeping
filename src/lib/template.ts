import { TransactionTemplate } from '@/models/transaction_template.ts';

export function isNoAvailableTemplate(templates: TransactionTemplate[], showHidden: boolean): boolean {
    for (let i = 0; i < templates.length; i++) {
        if (showHidden || !templates[i].hidden) {
            return false;
        }
    }

    return true;
}

export function getAvailableTemplateCount(templates: TransactionTemplate[], showHidden: boolean): number {
    let count = 0;

    for (let i = 0; i < templates.length; i++) {
        if (showHidden || !templates[i].hidden) {
            count++;
        }
    }

    return count;
}

export function getFirstShowingId(templates: TransactionTemplate[], showHidden: boolean): string | null {
    for (let i = 0; i < templates.length; i++) {
        if (showHidden || !templates[i].hidden) {
            return templates[i].id;
        }
    }

    return null;
}

export function getLastShowingId(templates: TransactionTemplate[], showHidden: boolean): string | null {
    for (let i = templates.length - 1; i >= 0; i--) {
        if (showHidden || !templates[i].hidden) {
            return templates[i].id;
        }
    }

    return null;
}
