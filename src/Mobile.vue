<template>
    <f7-app :params="f7params">
        <f7-view id="main-view" class="safe-areas" main url="/" :push-state="true" :push-state-animate="false"></f7-view>
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
                actions: {
                    animate: self.$settings.isEnableAnimate(),
                    backdrop: true,
                    closeOnEscape: true
                },
                dialog: {
                    animate: self.$settings.isEnableAnimate(),
                    backdrop: true
                },
                popover: {
                    animate: self.$settings.isEnableAnimate(),
                    backdrop: true,
                    closeOnEscape: true
                },
                popup: {
                    animate: self.$settings.isEnableAnimate(),
                    backdrop: true,
                    closeOnEscape: true,
                    swipeToClose: true
                },
                sheet: {
                    animate: self.$settings.isEnableAnimate(),
                    backdrop: true,
                    closeOnEscape: true
                },
                smartSelect: {
                    routableModals: false
                },
                touch: {
                    tapHold: true,
                    disableContextMenu: true
                },
                view: {
                    animate: self.$settings.isEnableAnimate()
                }
            }
        }
    }
}
</script>

<style>
html, body {
    position: fixed;
}

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

    --default-icon-color: var(--f7-text-color);
}

:root .theme-dark {
    --default-icon-color: var(--f7-text-color);
}

.no-right-border {
    border-right: 0;
}

.no-bottom-border {
    border-bottom: 0;
}

i.icon.la, i.icon.las, i.icon.lab {
    font-size: 28px;
}

.navbar .navbar-compact-icons.right a + a {
    margin-left: 0;
}

.vue-pincode-input {
    margin: 3px !important;
    padding: 5px !important;
    box-shadow: 0 0 3px rgba(0,0,0,.5) !important;
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

.tabbar-text-with-ellipsis > span {
    display: block;
    width: 100%;
    overflow: hidden;
    text-align: center;
    text-overflow: ellipsis;
    word-break: break-all;
    white-space: nowrap;
}

.tabbar-item-changed {
    color: var(--f7-theme-color);
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

.list .item-content .input.list-title-input {
    margin-top: calc(-1 * var(--f7-list-item-padding-vertical));
    margin-bottom: calc(-1 * var(--f7-list-item-padding-vertical));
}

.list .item-content .list-item-valign-middle {
    align-self: center;
}

.list .item-content .list-item-showing {
    color: rgba(0, 0, 0, 0.2);
    font-size: 16px;
    font-weight: bold;
}

.theme-dark .list .item-content .list-item-showing {
    color: rgba(255, 255, 255, 0.2);
}

.list .item-content .list-item-checked {
    font-size: 20px;
    color: var(--f7-radio-active-color, var(--f7-theme-color));
}

.lab-list-item-error-info div.item-footer {
    color: var(--f7-input-error-text-color)
}
</style>
