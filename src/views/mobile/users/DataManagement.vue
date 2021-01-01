<template>
    <f7-page>
        <f7-navbar :title="$t('Data Management')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-button external no-chevron target="_blank" :link="$services.getDataExportUrl()">{{ $t('Export Data') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {};
    }
};
</script>
