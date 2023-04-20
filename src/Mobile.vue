<template>
    <f7-app v-bind="f7params">
        <f7-view id="main-view" class="safe-areas" main url="/"></f7-view>
    </f7-app>
</template>

<script>
import { f7ready } from 'framework7-vue';
import routes from './router/mobile.js';

export default {
    data() {
        const self = this;

        return {
            f7params: {
                name: 'ezBookkeeping',
                theme: 'ios',
                colors: {
                    primary: '#c67e48'
                },
                routes: routes,
                darkMode: self.$settings.isEnableAutoDarkMode() ? 'auto' : false,
                touch: {
                    disableContextMenu: true,
                    tapHold: true
                },
                serviceWorker: {
                    path: self.$settings.isProduction() ? './sw.js' : undefined,
                    scope: './',
                },
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
                view: {
                    animate: self.$settings.isEnableAnimate(),
                    browserHistory: !self.isiOSHomeScreenMode(),
                    browserHistoryInitialMatch: true,
                    browserHistoryAnimate: false,
                    iosSwipeBackAnimateShadow: false,
                    mdSwipeBackAnimateShadow: false
                }
            },
            isDarkMode: undefined
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
    mounted() {
        f7ready((f7) => {
            this.isDarkMode = f7.darkMode;
            this.setThemeColorMeta(f7.darkMode);

            f7.on('pageBeforeOut',  () => {
                if (this.$ui.isModalShowing()) {
                    f7.actions.close('.actions-modal.modal-in', false);
                    f7.dialog.close('.dialog.modal-in', false);
                    f7.popover.close('.popover.modal-in', false);
                    f7.popup.close('.popup.modal-in', false);
                    f7.sheet.close('.sheet-modal.modal-in', false);
                }
            });

            f7.on('darkModeChange', (isDarkMode) => {
                this.isDarkMode = isDarkMode;
                this.setThemeColorMeta(isDarkMode);
            });
        });
    },
    methods: {
        isiOSHomeScreenMode() {
            if ((/iphone|ipod|ipad/gi).test(navigator.platform) && (/Safari/i).test(navigator.appVersion) &&
                window.matchMedia && window.matchMedia('(display-mode: standalone)').matches
            ) {
                return true;
            }

            return false;
        },
        setThemeColorMeta(isDarkMode) {
            if (isDarkMode) {
                document.querySelector('meta[name=theme-color]').setAttribute('content', '#121212');
            } else {
                document.querySelector('meta[name=theme-color]').setAttribute('content', '#f6f6f8');
            }
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

.smaller {
    font-size: 10px;
}

.readonly {
    pointer-events: none !important;
}

.skeleton-text {
    pointer-events: none !important;
}

.segmented.readonly .button:not(.button-active) > span,
.list.readonly .item-content .item-title.item-label,
.list.readonly .item-content .item-title > .item-header {
    opacity: 0.55 !important;
}

/** Replacing the default style of framework7 **/
:root {
    --default-icon-color: var(--f7-text-color);
}

:root .dark {
    --default-icon-color: var(--f7-text-color);
}

.color-gray {
    --f7-theme-color: #8e8e93;
    --f7-theme-color-rgb: 142, 142, 147;
    --f7-theme-color-shade: #79797f;
    --f7-theme-color-tint: #a3a3a7;
}

.ios .dark, .ios.dark {
    --f7-list-item-header-text-color: inherit !important;
}

i.icon.la, i.icon.las, i.icon.lab {
    font-size: 28px;
}

.chip.chip-placeholder {
    border: 0;
}

/** Replacing the default style of vue-datepicker **/
.dp__theme_light {
    --dp-primary-color: #c67e48;
}

.dp__theme_dark {
    --dp-primary-color: #c67e48;
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

.tabbar-primary-link,
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

.block-title .accordion-item-toggle .icon {
    color: var(--f7-list-chevron-icon-color);
    font-size: var(--f7-list-chevron-icon-font-size);
    font-weight: bolder;
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

.dark .list .item-content .list-item-showing {
    color: rgba(255, 255, 255, 0.2);
}

.accordion-item.list-item-checked > .item-link > .item-content .item-title {
    font-weight: bold;
}

.list.list-dividers li.list-group-title:first-child {
    border-radius: var(--f7-list-inset-border-radius) var(--f7-list-inset-border-radius) 0 0;
}

.list.list-dividers li.list-group-title:first-child:before {
    background-color: transparent;
}

.list.list-dividers li:last-child > .swipeout-content > .item-link > .item-content > .item-inner:after {
    background-color: transparent;
}

.list.inset li.list-group-title:first-child > a.button {
    border-radius: var(--f7-button-border-radius);
}

.list .item-content .list-item-checked-icon {
    font-size: 20px;
    color: var(--f7-radio-active-color, var(--f7-theme-color));
    margin-right: calc(var(--f7-list-item-media-margin) + var(--f7-checkbox-extra-margin));
}

.popover .popover-inner .list .item-content .list-item-checked-icon {
    margin-right: 0;
}

.list li.no-margin .item-content.item-input {
    margin: 0;
}

.ebk-list-item-error-info div.item-footer {
    color: var(--f7-input-error-text-color)
}

.skeleton-text .list-item-toggle .item-after {
    height: var(--f7-toggle-height);
}

.skeleton-text .list-item-toggle .item-after > span {
    line-height: var(--f7-toggle-height);
    font-size: var(--f7-toggle-height);
}

.no-sortable > .sortable-handler {
    display: none;
}

.card-header-content {
    opacity: 0.6;
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

/** Swipe handler **/
.swipe-handler {
    height: 16px;
    position: absolute;
    left: 0;
    width: 100%;
    top: 0;
    cursor: pointer;
    z-index: 10
}

.swipe-handler:after {
    content: "";
    width: 36px;
    height: 6px;
    position: absolute;
    left: 50%;
    top: 50%;
    margin-left: -18px;
    margin-top: -3px;
    border-radius: 3px;
    background: #666
}

/** list-item-with-multi-item for framework7 **/
.list-item-with-multi-item .item-content,
.list-item-with-multi-item .item-inner {
    padding: 0;
}

.list-item-with-multi-item .item-inner > div {
    width: 100%;
}

.list-item-with-multi-item > .item-content > .item-inner:after {
    background-color: transparent;
}

.list-item-with-multi-item .list-item-subitem:first-child .item-content {
    padding-left: calc(var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-left));
}

.list-item-with-multi-item .list-item-subitem .item-inner {
    display: block;
    width: 100%;
    padding-left: calc(var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-left));
    padding-top: var(--f7-list-item-padding-vertical);
    padding-bottom: var(--f7-list-item-padding-vertical);
}

.list-item-with-multi-item .list-item-subitem:first-child .item-inner {
    padding-left: 0;
}

/** Combination list for framework7 **/
.combination-list-wrapper {
    margin: 0;
    padding: 0;
}

.combination-list-wrapper .block-title {
    margin-top: 0;
    margin-bottom: 0;
}

.combination-list-wrapper .list.combination-list-header {
    margin: 0;
}

.combination-list-wrapper .list.combination-list-header .item-title {
    width: 100%;
    display: flex;
}

.combination-list-wrapper .list.combination-list-header > ul {
    background-color: var(--f7-list-group-title-bg-color);
}

.combination-list-wrapper .list.combination-list-header.combination-list-opened > ul {
    border-radius: var(--f7-list-inset-border-radius) var(--f7-list-inset-border-radius) 0 0;
}

.combination-list-wrapper .list.combination-list-header.combination-list-closed > ul {
    border-radius: var(--f7-list-inset-border-radius);
}

.combination-list-wrapper .list.combination-list-header .combination-list-chevron-icon {
    margin-left: auto;
}

.combination-list-wrapper .list.combination-list-content.inset > ul {
    border-radius: 0 0 var(--f7-list-inset-border-radius) var(--f7-list-inset-border-radius);
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

.nested-list-item.has-child-list-item .item-link.active-state {
    background-color: transparent;
}

.nested-list-item.has-child-list-item .item-link .item-inner {
    padding-right: 0;
}

.nested-list-item.has-child-list-item .item-link .item-inner:before {
    color: transparent;
}

.nested-list-item.has-child-list-item .item-link.active-state .item-inner .nested-list-item-child .item-link.active-state {
    background-color: var(--f7-list-link-pressed-bg-color);
}

.nested-list-item.has-child-list-item .item-link .item-inner .nested-list-item-child .item-link .item-inner {
    padding-right: calc(var(--f7-list-chevron-icon-area) + var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-right));
}

.nested-list-item.has-child-list-item .item-link .item-inner .nested-list-item-child .item-link .item-inner:before {
    color: var(--f7-list-chevron-icon-color);
}

.nested-list-item .nested-list-item-title {
    align-self: center;
    margin-left: var(--f7-list-item-media-margin);
    margin-right: var(--f7-list-item-media-margin);
    overflow: hidden;
    text-overflow: ellipsis;
}

.nested-list-item:last-child > .swipeout-content > .item-link > .item-content > .item-inner:after {
    background-color: transparent;
}

.sortable-enabled .nested-list-item .nested-list-item-child .item-inner {
    padding-right: var(--f7-safe-area-right) !important;
}
</style>
