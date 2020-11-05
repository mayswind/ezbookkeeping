<template>
    <f7-page name="home">
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <f7-list media-list class="skeleton-text">
            <f7-list-item title="Placehodler"></f7-list-item>
        </f7-list>

        <f7-toolbar tabbar bottom>
            <f7-link href="/details" class="lab-toolbar-link-with-icon" icon-f7="square_list" icon-size="24px" :text="$t('Details')"></f7-link>
            <f7-link href="/accounts" class="lab-toolbar-link-with-icon" icon-f7="creditcard" icon-size="24px" :text="$t('Accounts')"></f7-link>
            <f7-link href="/new" class="lab-toolbar-icon-only" icon-f7="plus_square" icon-size="36px"></f7-link>
            <f7-link href="/statistics" class="lab-toolbar-link-with-icon" icon-f7="chart_pie" icon-size="24px" :text="$t('Statistics')"></f7-link>
            <f7-link href="/settings" class="lab-toolbar-link-with-icon" icon-f7="gear_alt" icon-size="24px" :text="$t('Settings')"></f7-link>
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
    margin-top: 8px;
    margin-left: 0;
}
</style>
