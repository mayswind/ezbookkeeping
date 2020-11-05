<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <f7-list media-list class="skeleton-text">
            <f7-list-item title="Placeholder"></f7-list-item>
        </f7-list>

        <f7-toolbar tabbar labels bottom>
            <f7-link href="/transaction/details" class="lab-toolbar-link-with-icon" icon-f7="square_list" :text="$t('Details')"></f7-link>
            <f7-link href="/account/list" class="lab-toolbar-link-with-icon" icon-f7="creditcard" :text="$t('Accounts')"></f7-link>
            <f7-link href="/transaction/new" class="lab-toolbar-icon-only" icon-f7="plus_square" icon-size="42px"></f7-link>
            <f7-link href="/statistic/overview" class="lab-toolbar-link-with-icon" icon-f7="chart_pie" :text="$t('Statistics')"></f7-link>
            <f7-link href="/settings" class="lab-toolbar-link-with-icon" icon-f7="gear_alt" :text="$t('Settings')"></f7-link>
        </f7-toolbar>
    </f7-page>
</template>

<script>
export default {}
</script>

<style>
.lab-toolbar-link-with-icon i+span {
    font-size: 11px;
    margin-top: 3px;
    margin-left: 0;
}

.lab-toolbar-icon-only i {
    margin-top: 7px;
    margin-left: 0;
}
</style>
