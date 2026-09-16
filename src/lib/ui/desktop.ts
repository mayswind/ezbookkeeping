export type ComponentDensity = 'default' | 'comfortable' | 'compact';
export type InputVariant = 'filled' | 'underlined' | 'outlined' | 'plain' | 'solo' | 'solo-inverted' | 'solo-filled';

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
