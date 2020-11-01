<template>
    <f7-page name="home">
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
            <f7-nav-right>
                <f7-link href="/new-record" icon-f7="plus_square" :text="$t('New Record')"></f7-link>
            </f7-nav-right>
        </f7-navbar>
        <f7-tabs animated routable>
            <f7-tab id="main-tab-home" class="page-content" tab-active></f7-tab>
            <f7-tab id="main-tab-journals" class="page-content"></f7-tab>
            <f7-tab id="main-tab-accounts" class="page-content"></f7-tab>
            <f7-tab id="main-tab-charts" class="page-content"></f7-tab>
        </f7-tabs>
        <f7-toolbar tabbar bottom>
            <f7-link route-tab-id="main-tab-home" href="/" class="lab-toolbar-link-with-icon" icon-f7="house" icon-size="24px" :text="$t('Home')" tab-link-active></f7-link>
            <f7-link route-tab-id="main-tab-journals" href="/journals" class="lab-toolbar-link-with-icon" icon-f7="square_list" icon-size="24px" :text="$t('Journals')"></f7-link>
            <f7-link route-tab-id="main-tab-accounts" href="/accounts" class="lab-toolbar-link-with-icon" icon-f7="creditcard" icon-size="24px" :text="$t('Accounts')"></f7-link>
            <f7-link route-tab-id="main-tab-charts" href="/charts" class="lab-toolbar-link-with-icon" icon-f7="chart_pie" icon-size="24px" :text="$t('Charts')"></f7-link>
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

.lab-list-item-error-info div.item-footer {
    color: var(--f7-input-error-text-color)
}
</style>
