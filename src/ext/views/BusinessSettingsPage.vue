<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Business Features')">
                <template #subtitle>
                    {{ tt('Turn on inventory, sales and customers if you run a business. They stay hidden otherwise.') }}
                </template>
                <v-card-text>
                    <v-switch color="primary" hide-details inset
                              :label="tt('Enable inventory, sales and customers')"
                              :disabled="worksForSomeoneElse || saving"
                              :loading="saving"
                              :model-value="available"
                              @update:model-value="change(!!$event)" />

                    <v-alert class="mt-4" type="info" variant="tonal" density="compact" v-if="worksForSomeoneElse">
                        {{ tt('These features are on because you were invited to work in a business. They stay on while you are a member.') }}
                    </v-alert>

                    <div class="text-body-1 mt-4">{{ tt('When enabled you get:') }}</div>
                    <ul class="ms-6 mt-1">
                        <li>{{ tt('Inventory: items, stock levels and several locations') }}</li>
                        <li>{{ tt('Sales: a cart that records the sale in your accounts, also on credit') }}</li>
                        <li>{{ tt('Customers: who owes you money, and repayments') }}</li>
                        <li>{{ tt('Team: invite managers and staff to work in your business') }}</li>
                    </ul>

                    <div class="text-body-2 text-medium-emphasis mt-4">
                        {{ tt('This setting is saved to your account, so it applies on all your devices.') }}
                    </div>
                </v-card-text>
            </v-card>
            <v-card class="mt-4" :title="tt('Receipt details')" v-if="available">
                <template #subtitle>
                    {{ tt('This is printed at the top and bottom of your receipts.') }}
                </template>
                <v-card-text>
                    <v-form class="d-flex flex-column ga-4" @submit.prevent="saveProfile">
                        <v-text-field density="compact" hide-details="auto" :label="tt('Business name on receipts')"
                                      :hint="tt('Leave empty to use your own name.')" :disabled="savingProfile" v-model="profileForm.receiptName" />
                        <v-text-field density="compact" hide-details="auto" :label="tt('Address')" :disabled="savingProfile" v-model="profileForm.address" />
                        <v-text-field density="compact" hide-details="auto" :label="tt('Phone')" :disabled="savingProfile" v-model="profileForm.phone" />
                        <v-text-field density="compact" hide-details="auto" :label="tt('Message at the bottom')"
                                      :hint="tt('For example: Thank you for your business')" :disabled="savingProfile" v-model="profileForm.footer" />
                        <div>
                            <v-btn color="primary" type="submit" :loading="savingProfile">{{ tt('Save receipt details') }}</v-btn>
                        </div>
                    </v-form>
                </v-card-text>
            </v-card>
            <v-card class="mt-4" :title="tt('Export your business data')">
                <template #subtitle>
                    {{ tt('Take a copy of your records whenever you like. It opens in Excel, Google Sheets or LibreOffice.') }}
                </template>
                <v-card-text>
                    <p class="mb-3">
                        {{ tt('The download is a ZIP file with your items, stock, locations, customers, sales, repayments, team and activity, as spreadsheet files, plus a short guide.') }}
                    </p>
                    <p class="mb-4 text-medium-emphasis">
                        {{ tt('Your accounts, categories, tags and transactions are exported from') }}
                        <router-link to="/settings/user/data_management">{{ tt('Data Management') }}</router-link>.
                    </p>
                    <v-btn color="primary" :loading="exporting" @click="exportData">{{ tt('Download my business data') }}</v-btn>
                    <div class="text-body-2 text-medium-emphasis mt-3">
                        {{ tt('Only the owner of a business can download its data.') }}
                    </div>
                </v-card-text>
            </v-card>
            <ext-snack-bar ref="snackbar" />
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';

import { ref, reactive, watch, onMounted, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import api from '@/ext/api.ts';
import { saveBlob } from '@/ext/csv.ts';
import { useBusiness } from '@/ext/business.ts';
import { useBusinessFeatures } from '@/ext/features.ts';

type SnackBarType = InstanceType<typeof ExtSnackBar>;

const { tt } = useExtI18n();
const { ensureLoaded } = useBusiness();
const { available, worksForSomeoneElse, loadFeatureSetting, setEnabled } = useBusinessFeatures();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const saving = ref<boolean>(false);

const exporting = ref<boolean>(false);

async function exportData(): Promise<void> {
    exporting.value = true;

    try {
        const { blob, fileName } = await api.downloadBusinessExport();
        saveBlob(fileName, blob);
        snackbar.value?.showMessage(tt('Your data was downloaded'));
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        exporting.value = false;
    }
}

// receipt details belong to the person's own business
const savingProfile = ref<boolean>(false);
const profileForm = reactive({ receiptName: '', address: '', phone: '', footer: '' });

async function loadProfile(): Promise<void> {
    try {
        const profile = await api.getMyBusinessProfile();
        profileForm.receiptName = profile.receiptName;
        profileForm.address = profile.address;
        profileForm.phone = profile.phone;
        profileForm.footer = profile.footer;
    } catch (error) {
        snackbar.value?.showError(error);
    }
}

async function saveProfile(): Promise<void> {
    savingProfile.value = true;

    try {
        await api.updateMyBusinessProfile({ ...profileForm });
        snackbar.value?.showMessage(tt('Receipt details saved'));
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        savingProfile.value = false;
    }
}

watch(available, isAvailable => {
    if (isAvailable) {
        loadProfile();
    }
});

async function change(value: boolean): Promise<void> {
    saving.value = true;

    try {
        await setEnabled(value);
        snackbar.value?.showMessage(value ? tt('Business features are on') : tt('Business features are off'));
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        saving.value = false;
    }
}

onMounted(() => {
    ensureLoaded().catch(() => {
        // the switch still works without the list of businesses
    });

    loadFeatureSetting().catch(error => snackbar.value?.showError(error));

    if (available.value) {
        loadProfile();
    }
});
</script>
