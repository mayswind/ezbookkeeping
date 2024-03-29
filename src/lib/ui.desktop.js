export function getOuterHeight(element) {
    if (!element) {
        return 0;
    }

    const computedStyle = window.getComputedStyle(element);

    return ['height', 'padding-top', 'padding-bottom', 'margin-top', 'margin-bottom']
        .map((key) => parseInt(computedStyle.getPropertyValue(key), 10))
        .reduce((prev, cur) => prev + cur);
}

export function getNavSideBarOuterHeight(element) {
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
        totalHeight += getOuterHeight(children[i]);
    }

    return totalHeight;
}

export function getCssValue(element, name) {
    if (!element) {
        return 0;
    }

    const computedStyle = window.getComputedStyle(element);
    return computedStyle.getPropertyValue(name);
}

export function scrollToSelectedItem(parentEl, containerSelector, selectedItemSelector) {
    if (!parentEl) {
        return;
    }

    let container = parentEl;

    if (containerSelector) {
        const lists = parentEl.querySelectorAll(containerSelector);

        if (!lists.length || !lists[0]) {
            return;
        }

        container = lists[0];
    }

    const selectedItems = container.querySelectorAll(selectedItemSelector);

    if (!selectedItems.length || !selectedItems[0]) {
        return;
    }

    const selectedItem = selectedItems[0];
    const containerOuterHeight = getOuterHeight(container);
    const selectedItemOuterHeight = getOuterHeight(selectedItem);

    const targetPos = selectedItem.offsetTop - container.offsetTop - parseInt(getCssValue(container, 'padding-top'), 10)
        - (containerOuterHeight - selectedItemOuterHeight) / 2;

    if (targetPos <= 0) {
        return;
    }

    container.scrollTop = targetPos;
}
