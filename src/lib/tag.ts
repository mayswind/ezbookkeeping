import { TransactionTag } from '@/models/transaction_tag.ts';

export function isNoAvailableTag(tags: TransactionTag[], showHidden: boolean): boolean {
    for (let i = 0; i < tags.length; i++) {
        if (showHidden || !tags[i].hidden) {
            return false;
        }
    }

    return true;
}

export function getAvailableTagCount(tags: TransactionTag[], showHidden: boolean): number {
    let count = 0;

    for (let i = 0; i < tags.length; i++) {
        if (showHidden || !tags[i].hidden) {
            count++;
        }
    }

    return count;
}

export function getFirstShowingId(tags: TransactionTag[], showHidden: boolean): string | null {
    for (let i = 0; i < tags.length; i++) {
        if (showHidden || !tags[i].hidden) {
            return tags[i].id;
        }
    }

    return null;
}

export function getLastShowingId(tags: TransactionTag[], showHidden: boolean): string | null {
    for (let i = tags.length - 1; i >= 0; i--) {
        if (showHidden || !tags[i].hidden) {
            return tags[i].id;
        }
    }

    return null;
}
