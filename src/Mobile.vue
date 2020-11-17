<template>
    <f7-app :params="f7params">
        <f7-view id="main-view" main url="/" :push-state="true"></f7-view>
    </f7-app>
</template>

<script>
import routes from './router/mobile.js';

export default {
    data() {
        const self = this;

        return {
            f7params: {
                name: 'lab',
                id: 'net.mayswind.lab',
                theme: 'ios',
                autoDarkTheme: self.$settings.isEnableAutoDarkMode(),
                routes: routes,
                touch: {
                    tapHold: true
                },
                sheet: {
                    backdrop: true,
                    closeOnEscape: true
                },
                popup: {
                    swipeToClose: true,
                    closeOnEscape: true
                },
                view: {
                    animate : self.$settings.isEnableAnimate()
                }
            }
        }
    }
}
</script>

<style>
body {
    -ms-user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    user-select: none;
}

:root {
    --f7-theme-color: #c67e48;
    --f7-theme-color-rgb: 198, 126, 72;
    --f7-theme-color-shade: #af6a36;
    --f7-theme-color-tint: #d09467;
}

.work-break-all {
    word-break: break-all;
}

.full-line {
    width: 100%;
}

.icon-after-text {
    margin-left: 6px;
}

.badge.right-bottom-icon {
    margin-left: -12px;
    margin-top: 14px;
    width: 16px;
    height: 16px;
}

.badge.right-bottom-icon > .icon {
    font-size: 13px;
    width: 13px;
    height: 13px;
}

.nested-list-item .item-title {
    width: 100%;
}

.nested-list-item .item-inner {
    padding-right: 0;
}

.nested-list-item.has-child-list-item .item-inner {
    padding-bottom: 0;
}

.nested-list-item.has-child-list-item .nested-list-item-child .item-inner {
    padding-bottom: var(--f7-list-item-padding-vertical);
}

.nested-list-item .nested-list-item-title {
    align-self: center;
    margin-left: var(--f7-list-item-media-margin);
}

.sortable-enabled .nested-list-item .nested-list-item-child .item-inner {
    padding-right: var(--f7-safe-area-right) !important;
}

.no-sortable > .sortable-handler {
    display: none;
}

.list-item-media-valign-middle .item-media {
    align-self: normal !important;
}

.lab-list-item-error-info div.item-footer {
    color: var(--f7-input-error-text-color)
}
</style>
