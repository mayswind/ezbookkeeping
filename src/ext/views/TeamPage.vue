<template>
    <main-page-layout>
        <template #content>
            <v-row class="match-height">
                <v-col cols="12">
                    <v-card :title="tt('Invitations')" v-if="invitations.length > 0">
                        <v-card-text>
                            <v-list>
                                <v-list-item :key="invitation.ownerUid" :title="invitation.name"
                                             :subtitle="tt('Invited as {role}', { role: tt(invitation.role) })"
                                             v-for="invitation in invitations">
                                    <template #append>
                                        <v-btn class="me-2" color="primary" size="small" :disabled="busy"
                                               @click="respond(invitation, true)">{{ tt('Accept') }}</v-btn>
                                        <v-btn variant="text" size="small" :disabled="busy"
                                               @click="respond(invitation, false)">{{ tt('Decline') }}</v-btn>
                                    </template>
                                </v-list-item>
                            </v-list>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12" v-if="otherBusinesses.length > 0">
                    <v-card :title="tt('Businesses I work in')">
                        <v-card-text>
                            <v-list>
                                <v-list-item :key="business.ownerUid" :title="business.name" :subtitle="tt(business.role)"
                                             v-for="business in otherBusinesses">
                                    <template #append>
                                        <v-btn class="me-2" size="small" variant="tonal" color="primary"
                                               @click="switchBusiness(business.ownerUid)">{{ tt('Switch to this business') }}</v-btn>
                                        <v-btn size="small" variant="text" :disabled="busy"
                                               @click="leave(business)">{{ tt('Leave') }}</v-btn>
                                    </template>
                                </v-list-item>
                            </v-list>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12">
                    <v-card :title="tt('My team')">
                        <template #subtitle>
                            {{ tt('People you invite can work in your business with the role you give them.') }}
                        </template>
                        <v-card-text>
                            <v-form class="d-flex align-center flex-wrap ga-3" @submit.prevent="invite">
                                <v-text-field class="flex-grow-1" type="email" density="compact" hide-details
                                              :label="tt('Email of an existing user')"
                                              :disabled="busy" v-model="inviteEmail" />
                                <v-select style="max-width: 180px" density="compact" hide-details
                                          :label="tt('Role')" :items="roleOptions" item-title="title" item-value="value"
                                          :disabled="busy" v-model="inviteRole" />
                                <v-btn color="primary" type="submit" :disabled="busy || !inviteEmail.trim()">{{ tt('Invite') }}</v-btn>
                            </v-form>
                            <div class="text-body-2 text-medium-emphasis mt-2">
                                {{ tt('Managers run the business day to day. Staff can record sales and repayments.') }}
                            </div>
                        </v-card-text>

                        <v-table class="mt-2" v-if="staff.length > 0">
                            <thead>
                                <tr>
                                    <th>{{ tt('Name') }}</th>
                                    <th>{{ tt('Email') }}</th>
                                    <th style="width: 180px">{{ tt('Role') }}</th>
                                    <th>{{ tt('Status') }}</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr :key="member.staffUid" v-for="member in staff">
                                    <td>{{ member.nickname || member.username }}</td>
                                    <td>{{ member.email }}</td>
                                    <td>
                                        <v-select density="compact" hide-details variant="underlined"
                                                  :items="roleOptions" item-title="title" item-value="value"
                                                  :disabled="busy" :model-value="member.role"
                                                  @update:model-value="changeRole(member, $event)" />
                                    </td>
                                    <td>
                                        <v-chip size="small" :color="member.status === 'active' ? 'success' : 'warning'">
                                            {{ member.status === 'active' ? tt('Active') : tt('Waiting for reply') }}
                                        </v-chip>
                                    </td>
                                    <td class="text-end">
                                        <v-btn size="small" variant="text" color="error" :disabled="busy"
                                               @click="remove(member)">{{ tt('Remove') }}</v-btn>
                                    </td>
                                </tr>
                            </tbody>
                        </v-table>
                        <v-card-text v-else-if="!loading">
                            {{ tt('You have no team members yet.') }}
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col cols="12" v-if="staff.length > 0">
                    <v-card :title="tt('Recent activity')">
                        <template #subtitle>
                            {{ tt('Changes your managers and staff made in your business.') }}
                        </template>
                        <v-table density="compact">
                            <thead>
                                <tr>
                                    <th>{{ tt('When') }}</th>
                                    <th>{{ tt('Who') }}</th>
                                    <th>{{ tt('What') }}</th>
                                    <th>{{ tt('Result') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr :key="entry.id" v-for="entry in audit">
                                    <td>{{ formatTime(entry.time) }}</td>
                                    <td>{{ nameOf(entry.actorUid) }} <span class="text-medium-emphasis">({{ tt(entry.role) }})</span></td>
                                    <td>{{ describeAction(entry.path) }}</td>
                                    <td>
                                        <v-chip size="x-small" :color="entry.status < 400 ? 'success' : 'error'">{{ entry.status }}</v-chip>
                                    </td>
                                </tr>
                                <tr v-if="audit.length < 1">
                                    <td colspan="4" class="text-medium-emphasis">{{ tt('Nothing yet.') }}</td>
                                </tr>
                            </tbody>
                        </v-table>
                    </v-card>
                </v-col>
            </v-row>

            <confirm-dialog ref="confirmDialog" />
            <ext-snack-bar ref="snackbar" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';

import { ref, computed, onMounted, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from '@/ext/api.ts';
import { switchBusiness, useBusiness } from '@/ext/business.ts';
import type { AuditEntry, BusinessInfo, BusinessRole, StaffInfo } from '@/ext/types.ts';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof ExtSnackBar>;

const { tt, formatDateTimeToLongDateTime } = useI18n();
const { businesses, invitations, refresh } = useBusiness();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const busy = ref<boolean>(false);
const staff = ref<StaffInfo[]>([]);
const audit = ref<AuditEntry[]>([]);
const inviteEmail = ref<string>('');
const inviteRole = ref<BusinessRole>('staff');

const roleOptions = computed(() => [
    { title: tt('Staff'), value: 'staff' },
    { title: tt('Manager'), value: 'manager' }
]);

// businesses where the person works for somebody else
const otherBusinesses = computed<BusinessInfo[]>(() => businesses.value.filter(b => b.status === 'active'));

function formatTime(unixTime: number): string {
    return formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(unixTime));
}

function nameOf(uid: string): string {
    const member = staff.value.find(s => s.staffUid === uid);
    return member ? (member.nickname || member.username) : uid;
}

// "/api/v1/ext/sales/add.json" -> "ext sales add"
function describeAction(path: string): string {
    return path.replace(/^\/api\/v1\//, '').replace(/\.json$/, '').replace(/[/_]/g, ' ');
}

async function load(): Promise<void> {
    loading.value = true;

    try {
        await refresh();
        [staff.value, audit.value] = await Promise.all([api.listStaff(), api.listAudit()]);
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        loading.value = false;
    }
}

async function run(action: () => Promise<unknown>, doneMessage?: string): Promise<void> {
    busy.value = true;

    try {
        await action();

        if (doneMessage) {
            snackbar.value?.showMessage(doneMessage);
        }

        await load();
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        busy.value = false;
    }
}

async function invite(): Promise<void> {
    const email = inviteEmail.value.trim();

    if (!email) {
        return;
    }

    await run(async () => {
        await api.inviteStaff(email, inviteRole.value);
        inviteEmail.value = '';
    }, tt('Invitation sent'));
}

async function changeRole(member: StaffInfo, role: BusinessRole): Promise<void> {
    await run(() => api.setStaffRole(member.staffUid, role), tt('Role updated'));
}

function remove(member: StaffInfo): void {
    confirmDialog.value?.open('Remove {name} from your business? They lose access immediately.', { name: member.nickname || member.username }).then(() => {
        run(() => api.removeStaff(member.staffUid), tt('Removed'));
    }).catch(() => {
        // cancelled
    });
}

async function respond(invitation: BusinessInfo, accept: boolean): Promise<void> {
    await run(() => api.respondToInvitation(invitation.ownerUid, accept), accept ? tt('Invitation accepted') : tt('Invitation declined'));
}

function leave(business: BusinessInfo): void {
    confirmDialog.value?.open('Leave {name}? You will no longer have access to their business.', { name: business.name }).then(() => {
        run(() => api.leaveBusiness(business.ownerUid), tt('You left the business'));
    }).catch(() => {
        // cancelled
    });
}

onMounted(load);
</script>
