<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <f7-list media-list class="skeleton-text">
            <f7-list-item title="Placeholder"></f7-list-item>
        </f7-list>

        <f7-toolbar tabbar labels bottom>
            <f7-link href="/transaction/list">
                <f7-icon f7="square_list"></f7-icon>
                <span class="tabbar-label">{{ $t('Details') }}</span>
            </f7-link>
            <f7-link href="/account/list">
                <f7-icon f7="creditcard"></f7-icon>
                <span class="tabbar-label">{{ $t('Accounts') }}</span>
            </f7-link>
            <f7-link href="/transaction/add">
                <f7-icon f7="plus_square" size="42px" class="lab-tarbar-big-icon"></f7-icon>
            </f7-link>
            <f7-link href="/statistic/overview">
                <f7-icon f7="chart_pie"></f7-icon>
                <span class="tabbar-label">{{ $t('Statistics') }}</span>
            </f7-link>
            <f7-link href="/settings">
                <f7-icon f7="gear_alt"></f7-icon>
                <span class="tabbar-label">{{ $t('Settings') }}</span>
            </f7-link>
        </f7-toolbar>
    </f7-page>
</template>

<script>
export default {}
</script>

<style>
.lab-tarbar-big-icon {
    margin-top: 7px;
    margin-left: 0;
}
</style>
