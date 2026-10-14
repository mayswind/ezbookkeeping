import { nextTick } from 'vue';

export type ComponentDensity = 'default' | 'comfortable' | 'compact';
export type InputVariant = 'filled' | 'underlined' | 'outlined' | 'plain' | 'solo' | 'solo-inverted' | 'solo-filled';
export type SnackBarVariant = 'text' | 'flat' | 'elevated' | 'tonal' | 'outlined' | 'plain';

export type TwoColumnListColumn = 'left' | 'right';

export interface TwoColumnListKeyboardOptions {
    getList: (column: TwoColumnListColumn) => HTMLElement | undefined;
    getSelectedIndex: (column: TwoColumnListColumn) => number;
    onSelect: (column: TwoColumnListColumn, index: number) => void;
    onLeave: () => void;
}

function getListItems(list: HTMLElement | undefined): HTMLElement[] {
    return Array.from(list?.children ?? [])
        .filter((item): item is HTMLElement => item instanceof HTMLElement && item.hasAttribute('tabindex'));
}

export function getOuterHeight(element: HTMLElement | null): number {
    if (!element) {
        return 0;
    }

    const computedStyle = window.getComputedStyle(element);

    return ['height', 'padding-top', 'padding-bottom', 'margin-top', 'margin-bottom']
        .map((key) => parseInt(computedStyle.getPropertyValue(key), 10))
        .reduce((prev, cur) => prev + cur);
}

export function getNavSideBarOuterHeight(element: HTMLElement | null): number {
    if (!element) {
        return 0;
    }

    const contentEl = element.querySelectorAll('.v-navigation-drawer__content');

    if (!contentEl || !contentEl[0]) {
        return 0;
    }

    const children = contentEl[0].children;

    if (!children || children.length < 1) {
        return 0;
    }

    let totalHeight = 0;

    for (let i = 0; i < children.length; i++) {
        totalHeight += getOuterHeight(children[i] as HTMLElement);
    }

    return totalHeight;
}

export function setChildInputFocus(parentEl: HTMLElement | undefined, childSelector: string): void {
    if (!parentEl) {
        return;
    }

    const childElement = parentEl.querySelector(childSelector);

    if (!childElement || !(childElement as HTMLInputElement)) {
        return;
    }

    const childInput = (childElement as HTMLInputElement);
    childInput.focus();
    childInput.select();
}

export function focusParentWhenClicked(event: MouseEvent, currentClass: string, parentSelector: string): void {
    const activeElement = document.activeElement as HTMLElement | null;
    if (!activeElement?.classList.contains(currentClass)) {
        return;
    }

    const currentTarget = event.currentTarget as HTMLElement | null;
    currentTarget?.closest<HTMLElement>(parentSelector)?.focus({ preventScroll: true });
}

export function focusTableScrollContainer(event: MouseEvent): void {
    const target = event.target as HTMLElement | null;
    const tableRoot = event.currentTarget as HTMLElement | null;

    if (!target || !tableRoot) {
        return;
    }

    if (target.closest([
        'button',
        'input',
        'textarea',
        'select',
        'label',
        '[contenteditable="true"]',
        '[role="button"]',
        '[role="checkbox"]',
        '[role="radio"]',
        '[role="switch"]',
        '[role="slider"]'
    ].join(','))) {
        return;
    }

    const wrapper = tableRoot.querySelector<HTMLElement>(':scope > .v-table__wrapper');

    if (!wrapper) {
        return;
    }

    const canScroll = wrapper.scrollHeight > wrapper.clientHeight || wrapper.scrollWidth > wrapper.clientWidth;

    if (!canScroll) {
        return;
    }

    wrapper.tabIndex = -1;
    wrapper.focus({ preventScroll: true });
}

export function setTwoColumnItemFocusWhenKeydown(event: KeyboardEvent, currentColumn: TwoColumnListColumn, options: TwoColumnListKeyboardOptions): void {
    let targetColumn = currentColumn;
    let direction: 'next' | 'prev' | 'first' | 'last' = 'first';

    if (event.key === 'Tab') {
        targetColumn = event.shiftKey ? 'left' : 'right';
    } else if (event.key === 'ArrowRight') {
        if (currentColumn !== 'left') {
            return;
        }
        targetColumn = 'right';
    } else if (event.key === 'ArrowLeft') {
        if (currentColumn !== 'right') {
            return;
        }
        targetColumn = 'left';
    } else if (event.key === 'ArrowDown') {
        direction = 'next';
    } else if (event.key === 'ArrowUp') {
        direction = 'prev';
    } else if (event.key === 'Home') {
        direction = 'first';
    } else if (event.key === 'End') {
        direction = 'last';
    } else {
        return;
    }

    event.stopPropagation();

    const focusTargetItem = function (): void {
        const items = getListItems(options.getList(targetColumn));

        if (!items.length) {
            if (event.key === 'Tab') {
                options.onLeave();
            }
            return;
        }

        let index = items.findIndex(item => item.contains(document.activeElement));

        if (targetColumn !== currentColumn) {
            index = Math.max(options.getSelectedIndex(targetColumn), 0);
        } else if (direction === 'first') {
            index = 0;
        } else if (direction === 'last') {
            index = items.length - 1;
        } else if (index < 0) {
            index = direction === 'next' ? 0 : items.length - 1;
        } else {
            index = (index + (direction === 'next' ? 1 : -1) + items.length) % items.length;
        }

        const item = items[Math.min(index, items.length - 1)];
        item?.focus({ preventScroll: true });
        item?.scrollIntoView({ block: 'nearest', inline: 'nearest' });
    }

    if (event.key === 'Tab') {
        const focusedIndex = getListItems(options.getList(currentColumn)).findIndex(item => item.contains(document.activeElement));

        if (focusedIndex >= 0) {
            options.onSelect(currentColumn, focusedIndex);
        }

        if (targetColumn === currentColumn) {
            options.onLeave();
            return;
        }

        event.preventDefault();
        nextTick(focusTargetItem);
        return;
    }

    event.preventDefault();
    focusTargetItem();
}
