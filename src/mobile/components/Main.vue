<template>
    <f7-page name="home">
        <f7-navbar :title="$t('global.app.title')"></f7-navbar>
        <f7-tabs animated swipeable routable>
            <f7-tab id="main-tab-home" class="page-content" tab-active />
            <f7-tab id="main-tab-journals" class="page-content" />
            <f7-tab id="main-tab-accounts" class="page-content" />
            <f7-tab id="main-tab-charts" class="page-content" />
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
