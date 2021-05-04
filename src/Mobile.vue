<template>
    <f7-app :params="f7params">
        <f7-view id="main-view" class="safe-areas" main url="/"></f7-view>
    </f7-app>
</template>

<script>
import routes from './router/mobile.js';

export default {
    data() {
        const self = this;

        return {
            f7params: {
                name: 'ezBookkeeping',
                id: 'net.mayswind.ezbookkeeping',
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
                    animate: self.$settings.isEnableAnimate(),
                    pushState: !self.isiOSHomeScreenMode(),
                    pushStateAnimate: false,
                    iosSwipeBackAnimateShadow: false,
                    mdSwipeBackAnimateShadow: false
                },
                calendar: {
                    locale: 'en',
                    openIn: 'customModal',
                    backdrop: true
                },
                serviceWorker: {
                    path: self.$settings.isProduction() ? './sw.js' : undefined,
                    scope: './',
                }
            }
        }
    },
    created() {
        if (this.$user.isUserLogined()) {
            if (!this.$settings.isEnableApplicationLock()) {
                // refresh token if user is logined
                this.$store.dispatch('refreshTokenAndRevokeOldToken');

                // auto refresh exchange rates data
                if (this.$settings.isAutoUpdateExchangeRatesData()) {
                    this.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
                }
            }
        }
    },
    methods: {
        isiOSHomeScreenMode() {
            if ((/iphone|ipod|ipad/gi).test(navigator.platform) && (/Safari/i).test(navigator.appVersion) &&
                window.matchMedia && window.matchMedia('(display-mode: standalone)').matches
            ) {
                return true;
            }

            return false;
        }
    }
}
</script>

<style>
/** Global style **/
html, body {
    position: fixed;
}

body {
    -ms-user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    user-select: none;
}

/** Common class **/
.no-right-border {
    border-right: 0;
}

.no-bottom-border {
    border-bottom: 0;
}

.work-break-all {
    word-break: break-all;
}

.full-line {
    width: 100%;
}

/** Replacing the default style of framework7 **/
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

i.icon.la, i.icon.las, i.icon.lab {
    font-size: 28px;
}

.chip.chip-placeholder {
    border: 0;
}

/** Common class for replacing the default style of framework7 **/
.navbar .navbar-compact-icons.right a + a {
    margin-left: 0;
}

.toolbar-item-auto-size .toolbar-inner {
    padding-left: 16px;
    padding-right: 16px;
}

.toolbar-item-auto-size .toolbar-inner > .link {
    width: auto;
}

.tabbar-item-changed {
    color: var(--f7-theme-color);
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

.list-item-media-valign-middle .item-media {
    align-self: normal !important;
}

.list-item-no-item-after .item-after {
    display: none;
}

.list-item-with-header-and-title .item-content .item-header {
    margin-bottom: 11px;
}

.list-item-with-header-and-title.list-item-title-hide-overflow .item-content .list-item-custom-title {
    overflow: hidden;
    text-overflow: ellipsis;
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

.accordion-item.list-item-checked > .item-link > .item-content .item-title {
    font-weight: bold;
}

.list .item-content .list-item-checked-icon {
    font-size: 20px;
    color: var(--f7-radio-active-color, var(--f7-theme-color));
}

.ebk-list-item-error-info div.item-footer {
    color: var(--f7-input-error-text-color)
}

.no-sortable > .sortable-handler {
    display: none;
}

.card-header-content {
    opacity: 0.6;
}

.card-chevron-icon {
    color: var(--f7-list-chevron-icon-color);
    font-size: var(--f7-list-chevron-icon-font-size);
    font-weight: bolder;
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

/** Nested List item for framework7 **/
.nested-list-item .item-title {
    width: 100%;
}

.nested-list-item .item-inner {
    padding-right: 0;
}

.nested-list-item.has-child-list-item .nested-list-item-child .item-inner {
    padding-bottom: var(--f7-list-item-padding-vertical);
}

.nested-list-item .nested-list-item-title {
    align-self: center;
    margin-left: var(--f7-list-item-media-margin);
    margin-right: var(--f7-list-item-media-margin);
    overflow: hidden;
    text-overflow: ellipsis;
}

.sortable-enabled .nested-list-item .nested-list-item-child .item-inner {
    padding-right: var(--f7-safe-area-right) !important;
}

/** Replacing the default style of Vue-pincode-input **/
.vue-pincode-input {
    margin: 3px !important;
    padding: 5px !important;
    box-shadow: 0 0 2px rgba(0,0,0,.5) !important;
}

.theme-dark .vue-pincode-input {
    box-shadow: 0 0 2px rgba(255,255,255,.5) !important;
}
</style>
