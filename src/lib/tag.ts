import { reversed } from '@/core/base.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';

export function isNoAvailableTag(tags: TransactionTag[], showHidden: boolean): boolean {
    for (const tag of tags) {
        if (showHidden || !tag.hidden) {
            return false;
        }
    }

    return true;
}

export function getAvailableTagCount(tags: TransactionTag[], showHidden: boolean): number {
    let count = 0;

    for (const tag of tags) {
        if (showHidden || !tag.hidden) {
            count++;
        }
    }

    return count;
}

export function getFirstShowingId(tags: TransactionTag[], showHidden: boolean): string | null {
    for (const tag of tags) {
        if (showHidden || !tag.hidden) {
            return tag.id;
        }
    }

    return null;
}

export function getLastShowingId(tags: TransactionTag[], showHidden: boolean): string | null {
    for (const tag of reversed(tags)) {
        if (showHidden || !tag.hidden) {
            return tag.id;
        }
    }

    return null;
}
