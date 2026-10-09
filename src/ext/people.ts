import { ref } from 'vue';

import { useExtI18n } from './i18n.ts';
import api from './api.ts';
import { useBusiness } from './business.ts';
import type { PersonInfo } from './types.ts';

// Who works in the business shown on screen, so records can say who made them. The list is per business and is
// loaded again whenever a page asks, because people come and go.
const people = ref<PersonInfo[]>([]);

export function usePeople() {
    const { tt } = useExtI18n();
    const { own } = useBusiness();

    /**
     * The name of the person who made a record: "Sam", "Sam (you)" for the logged-in user, "Sam (no longer here)" for
     * somebody who was removed. Empty when no person is recorded (records from before this was kept).
     */
    function describe(list: PersonInfo[], uid: string): string {
        if (!uid || uid === '0') {
            return '';
        }

        const person = list.find(p => p.uid === uid);

        if (!person) {
            return tt('Someone who has left');
        }

        const name = person.name || uid;

        if (own.value && own.value.ownerUid === uid) {
            return tt('{name} (you)', { name });
        }

        return person.active ? name : tt('{name} (no longer here)', { name });
    }

    /** People of the business being worked in. */
    async function load(): Promise<void> {
        people.value = await api.listPeople();
    }

    function nameOf(uid: string): string {
        return describe(people.value, uid);
    }

    /** People of the caller's own business (the Team page), unaffected by the business being worked in. */
    function ownTeam() {
        const team = ref<PersonInfo[]>([]);

        return {
            team,
            load: async (): Promise<void> => {
                team.value = await api.listMyPeople();
            },
            nameOf: (uid: string): string => describe(team.value, uid)
        };
    }

    return { people, load, nameOf, ownTeam };
}
