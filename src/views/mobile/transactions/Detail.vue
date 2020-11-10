<template>
    <f7-page>
        <f7-navbar :title="$t('Transaction Details')" :back-link="$t('Back')">
            <f7-subnavbar :inner="false">
                <f7-searchbar
                    :custom-search="true"
                    :placeholder="$t('Search')"
                    :disable-button-text="$t('Cancel')"
                ></f7-searchbar>
            </f7-subnavbar>
        </f7-navbar>

        <f7-list media-list class="skeleton-text">
            <f7-list-item title="Placeholder"></f7-list-item>
        </f7-list>
    </f7-page>
</template>
