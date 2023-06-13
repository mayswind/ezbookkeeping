import { f7, f7ready } from 'framework7-vue';

import fontConstants from '@/consts/font.js';
import settings from './settings.js';
import {
    getLocalizedError,
    getLocalizedErrorParameters
} from './i18n.js';

export function showAlert(message, confirmCallback, translateFn) {
    let parameters = {};

    if (message && message.error) {
        const localizedError = getLocalizedError(message.error);
        message = localizedError.message;
        parameters = getLocalizedErrorParameters(localizedError.parameters, s => translateFn(s));
    }

    f7ready((f7) => {
        f7.dialog.create({
            title: translateFn('global.app.title'),
            text: translateFn(message, parameters),
            animate: settings.isEnableAnimate(),
            buttons: [
                {
                    text: translateFn('OK'),
                    onClick: confirmCallback
                }
            ]
        }).open();
    });
}

export function showConfirm(message, confirmCallback, cancelCallback, translateFn) {
    f7ready((f7) => {
        f7.dialog.create({
            title: translateFn('global.app.title'),
            text: translateFn(message),
            animate: settings.isEnableAnimate(),
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

export function showToast(message, timeout, translateFn) {
    let parameters = {};

    if (message && message.error) {
        const localizedError = getLocalizedError(message.error);
        message = localizedError.message;
        parameters = getLocalizedErrorParameters(localizedError.parameters, s => translateFn(s));
    }

    f7ready((f7) => {
        f7.toast.create({
            text: translateFn(message, parameters),
            position: 'center',
            closeTimeout: timeout || 1500
        }).open();
    });
}

export function showLoading(delayConditionFunc, delayMills) {
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

export function routeBackOnError(f7router, errorPropertyName) {
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

export function elements(selector) {
    return f7.$(selector);
}

export function isModalShowing() {
    return f7.$('.modal-in').length;
}

export function onSwipeoutDeleted(domId, callback) {
    f7.swipeout.delete(f7.$('#' + domId), callback);
}

export function autoChangeTextareaSize(el) {
    f7.$(el).find('textarea').each(el => {
        el.scrollTop = 0;
        el.style.height = '';
        el.style.height = el.scrollHeight + 'px';
    });
}

export function setAppFontSize(type) {
    const htmlElement = elements('html');

    for (let i = 0; i < fontConstants.allFontSizeArray.length; i++) {
        const fontSizeType = fontConstants.allFontSizeArray[i];

        if (fontSizeType.type === type) {
            if (!htmlElement.hasClass(fontSizeType.className)) {
                htmlElement.addClass(fontSizeType.className);
            }
        } else {
            htmlElement.removeClass(fontSizeType.className);
        }
    }
}

export function getFontSizePreviewClassName(type) {
    for (let i = 0; i < fontConstants.allFontSizeArray.length; i++) {
        const fontSizeType = fontConstants.allFontSizeArray[i];

        if (fontSizeType.type === type) {
            return fontConstants.fontSizePreviewClassNamePrefix + fontSizeType.className;
        }
    }

    return fontConstants.fontSizePreviewClassNamePrefix + fontConstants.defaultFontSize.className;
}
