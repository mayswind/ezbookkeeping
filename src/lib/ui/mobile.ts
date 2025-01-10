import { useI18n as useVueI18n } from 'vue-i18n';
import { f7, f7ready } from 'framework7-vue';
import type { Dialog, Picker, Router } from 'framework7/types';

import { FontSize, FONT_SIZE_PREVIEW_CLASSNAME_PREFIX } from '@/core/font.ts';
import { getNumberValue } from '../common.ts';
import { isEnableAnimate } from '../settings.ts';
// @ts-expect-error the above file is migrating to ts
import { translateError } from '@/locales/helper.js';

interface Framework7Dom {
    length: number;
    [index: number]: Element;
    find: (selector?: string) => Framework7Dom;
    offset(): { top: number; left: number };
    scrollTop(position: number, duration?: number, callback?: () => void): Framework7Dom;
    outerHeight(includeMargin?: boolean): number;
    css(property: string): string | number;
}

type TranslateFunction = (message: string) => string;

export function showAlert(message: string, confirmCallback: (dialog: Dialog.Dialog, e: Event) => void, translateFn: TranslateFunction): void {
    f7ready((f7) => {
        f7.dialog.create({
            title: translateFn('global.app.title'),
            text: translateError(message, translateFn),
            animate: isEnableAnimate(),
            buttons: [
                {
                    text: translateFn('OK'),
                    onClick: confirmCallback
                }
            ]
        }).open();
    });
}

export function showConfirm(message: string, confirmCallback: (dialog: Dialog.Dialog, e: Event) => void, cancelCallback: (dialog: Dialog.Dialog, e: Event) => void, translateFn: TranslateFunction): void {
    f7ready((f7) => {
        f7.dialog.create({
            title: translateFn('global.app.title'),
            text: translateFn(message),
            animate: isEnableAnimate(),
            buttons: [
                {
                    text: translateFn('Cancel'),
                    onClick: cancelCallback
                },
                {
                    text: translateFn('OK'),
                    onClick: confirmCallback
                }
            ]
        }).open();
    });
}

export function showToast(message: string, timeout: number | undefined, translateFn: TranslateFunction):void {
    f7ready((f7) => {
        f7.toast.create({
            text: translateError(message, translateFn),
            position: 'center',
            closeTimeout: timeout || 1500
        }).open();
    });
}

export function showLoading(delayConditionFunc: () => boolean, delayMills: number): void {
    if (!delayConditionFunc) {
        f7ready((f7) => {
            return f7.preloader.show();
        });
        return;
    }

    f7ready((f7) => {
        setTimeout(() => {
            if (delayConditionFunc()) {
                f7.preloader.show();
            }
        }, delayMills || 200);
    });
}

export function hideLoading() {
    f7ready((f7) => {
        return f7.preloader.hide();
    });
}

export function createInlinePicker(containerEl: HTMLElement, inputEl: HTMLElement, cols: Picker.ColumnParameters[], value: unknown[], events?: { change: (picker: Picker.Picker, value: unknown, displayValue: unknown) => void }): Picker.Picker {
    return f7.picker.create({
        containerEl: containerEl,
        inputEl: inputEl,
        toolbar: false,
        rotateEffect: true,
        value: value,
        cols: cols,
        on: events || {}
    });
}

export function routeBackOnError(f7router: Router.Router, errorPropertyName: string): void {
    // @ts-expect-error vue SFC would be migrated to composition API and this function would be removed in the future
    const self = this;
    const router = f7router;

    const unwatch = self.$watch(errorPropertyName, () => {
        if (self[errorPropertyName]) {
            setTimeout(() => {
                if (unwatch) {
                    unwatch();
                }

                router.back();
            }, 200);
        }
    }, {
        immediate: true
    });
}

export function isModalShowing(): number {
    return f7.$('.modal-in').length;
}

export function onSwipeoutDeleted(domId: string, callback: () => void): void {
    f7.swipeout.delete(f7.$('#' + domId), callback);
}

export function autoChangeTextareaSize(el: HTMLElement): void {
    f7.$(el).find('textarea').each((child: HTMLElement) => {
        child.scrollTop = 0;
        child.style.height = '';
        child.style.height = child.scrollHeight + 'px';
    });
}

export function setAppFontSize(type: number): void {
    const htmlElement = f7.$('html');
    const allFontSizes = FontSize.values();

    for (let i = 0; i < allFontSizes.length; i++) {
        const fontSizeType = allFontSizes[i];

        if (fontSizeType.type === type) {
            if (!htmlElement.hasClass(fontSizeType.className)) {
                htmlElement.addClass(fontSizeType.className);
            }
        } else {
            htmlElement.removeClass(fontSizeType.className);
        }
    }
}

export function getFontSizePreviewClassName(type: number): string {
    const allFontSizes = FontSize.values();

    for (let i = 0; i < allFontSizes.length; i++) {
        const fontSizeType = allFontSizes[i];

        if (fontSizeType.type === type) {
            return FONT_SIZE_PREVIEW_CLASSNAME_PREFIX + fontSizeType.className;
        }
    }

    return FONT_SIZE_PREVIEW_CLASSNAME_PREFIX + FontSize.Default.className;
}

export function scrollToSelectedItem(parentEl: Framework7Dom, containerSelector: string, selectedItemSelector: string): void {
    if (!parentEl || !parentEl.length) {
        return;
    }

    const container = parentEl.find(containerSelector);
    const selectedItem = parentEl.find(selectedItemSelector);

    if (!container.length || !selectedItem.length) {
        return;
    }

    const containerPaddingTop = getNumberValue(container.css('padding-top'), 0);

    let targetPos = selectedItem.offset().top - container.offset().top - containerPaddingTop
        - (container.outerHeight() - selectedItem.outerHeight()) / 2;

    if (selectedItem.length > 1) {
        const firstSelectedItem = f7.$(selectedItem[0]);
        const lastSelectedItem = f7.$(selectedItem[selectedItem.length - 1]);

        const firstSelectedItemInTop = firstSelectedItem.offset().top - container.offset().top - containerPaddingTop;
        const lastSelectedItemInTop = lastSelectedItem.offset().top - container.offset().top - containerPaddingTop;
        const lastSelectedItemInBottom = lastSelectedItem.offset().top - container.offset().top - containerPaddingTop
            - (container.outerHeight() - firstSelectedItem.outerHeight());

        targetPos = (firstSelectedItemInTop + lastSelectedItemInBottom) / 2;

        if (lastSelectedItemInTop - firstSelectedItemInTop > container.outerHeight()) {
            targetPos = firstSelectedItemInTop;
        }
    }

    if (targetPos <= 0) {
        return;
    }

    container.scrollTop(targetPos);
}

export function useI18nUIComponents() {
    const i18nGlobal = useVueI18n();

    return {
        showAlert: (message: string, confirmCallback: (dialog: Dialog.Dialog, e: Event) => void) => showAlert(message, confirmCallback, i18nGlobal.t),
        showConfirm: (message: string, confirmCallback: (dialog: Dialog.Dialog, e: Event) => void, cancelCallback: (dialog: Dialog.Dialog, e: Event) => void): void => showConfirm(message, confirmCallback, cancelCallback, i18nGlobal.t),
        showToast: (message: string, timeout?: number): void => showToast(message, timeout, i18nGlobal.t)
    }
}
