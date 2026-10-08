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
            <ext-snack-bar ref="snackbar" />
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';

import { ref, onMounted, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { useBusiness } from '@/ext/business.ts';
import { useBusinessFeatures } from '@/ext/features.ts';

type SnackBarType = InstanceType<typeof ExtSnackBar>;

const { tt } = useExtI18n();
const { ensureLoaded } = useBusiness();
const { available, worksForSomeoneElse, loadFeatureSetting, setEnabled } = useBusinessFeatures();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const saving = ref<boolean>(false);

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
});
</script>
