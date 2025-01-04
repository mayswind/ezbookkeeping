import { defineStore } from 'pinia';

import { isEquals } from '@/lib/common.ts';
import services from '@/lib/services.js';
import logger from '@/lib/logger.ts';

function loadTransactionTagList(state, tags) {
    state.allTransactionTags = tags;
    state.allTransactionTagsMap = {};

    for (let i = 0; i < tags.length; i++) {
        const tag = tags[i];
        state.allTransactionTagsMap[tag.id] = tag;
    }
}

function addTagToTransactionTagList(state, tag) {
    state.allTransactionTags.push(tag);
    state.allTransactionTagsMap[tag.id] = tag;
}

function updateTagInTransactionTagList(state, tag) {
    for (let i = 0; i < state.allTransactionTags.length; i++) {
        if (state.allTransactionTags[i].id === tag.id) {
            state.allTransactionTags.splice(i, 1, tag);
            break;
        }
    }

    state.allTransactionTagsMap[tag.id] = tag;
}

function updateTagDisplayOrderInTransactionTagList(state, { from, to }) {
    state.allTransactionTags.splice(to, 0, state.allTransactionTags.splice(from, 1)[0]);
}

function updateTagVisibilityInTransactionTagList(state, { tag, hidden }) {
    if (state.allTransactionTagsMap[tag.id]) {
        state.allTransactionTagsMap[tag.id].hidden = hidden;
    }
}

function removeTagFromTransactionTagList(state, tag) {
    for (let i = 0; i < state.allTransactionTags.length; i++) {
        if (state.allTransactionTags[i].id === tag.id) {
            state.allTransactionTags.splice(i, 1);
            break;
        }
    }

    if (state.allTransactionTagsMap[tag.id]) {
        delete state.allTransactionTagsMap[tag.id];
    }
}

export const useTransactionTagsStore = defineStore('transactionTags', {
    state: () => ({
        allTransactionTags: [],
        allTransactionTagsMap: {},
        transactionTagListStateInvalid: true,
    }),
    getters: {
        allVisibleTags(state) {
            const allVisibleTags = [];

            for (let i = 0; i < state.allTransactionTags.length; i++) {
                const tag = state.allTransactionTags[i];

                if (!tag.hidden) {
                    allVisibleTags.push(tag);
                }
            }

            return allVisibleTags;
        },
        allAvailableTagsCount(state) {
            return state.allTransactionTags.length;
        },
        allVisibleTagsCount(state) {
            return state.allVisibleTags.length;
        }
    },
    actions: {
        generateNewTransactionTagModel() {
            return {
                id: '',
                name: ''
            };
        },
        updateTransactionTagListInvalidState(invalidState) {
            this.transactionTagListStateInvalid = invalidState;
        },
        resetTransactionTags() {
            this.allTransactionTags = [];
            this.allTransactionTagsMap = {};
            this.transactionTagListStateInvalid = true;
        },
        loadAllTags({ force }) {
            const self = this;

            if (!force && !self.transactionTagListStateInvalid) {
                return new Promise((resolve) => {
                    resolve(self.allTransactionTags);
                });
            }

            return new Promise((resolve, reject) => {
                services.getAllTransactionTags().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve tag list' });
                        return;
                    }

                    if (self.transactionTagListStateInvalid) {
                        self.updateTransactionTagListInvalidState(false);
                    }

                    if (force && data.result && isEquals(self.allTransactionTags, data.result)) {
                        reject({ message: 'Tag list is up to date' });
                        return;
                    }

                    loadTransactionTagList(self, data.result);

                    resolve(data.result);
                }).catch(error => {
                    if (force) {
                        logger.error('failed to force load tag list', error);
                    } else {
                        logger.error('failed to load tag list', error);
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve tag list' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        saveTag({ tag }) {
            const self = this;

            return new Promise((resolve, reject) => {
                let promise = null;

                if (!tag.id) {
                    promise = services.addTransactionTag(tag);
                } else {
                    promise = services.modifyTransactionTag(tag);
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (!tag.id) {
                            reject({ message: 'Unable to add tag' });
                        } else {
                            reject({ message: 'Unable to save tag' });
                        }
                        return;
                    }

                    if (!tag.id) {
                        addTagToTransactionTagList(self, data.result);
                    } else {
                        updateTagInTransactionTagList(self, data.result);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save tag', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (!tag.id) {
                            reject({ message: 'Unable to add tag' });
                        } else {
                            reject({ message: 'Unable to save tag' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        changeTagDisplayOrder({ tagId, from, to }) {
            const self = this;

            return new Promise((resolve, reject) => {
                let tag = null;

                for (let i = 0; i < self.allTransactionTags.length; i++) {
                    if (self.allTransactionTags[i].id === tagId) {
                        tag = self.allTransactionTags[i];
                        break;
                    }
                }

                if (!tag || !self.allTransactionTags[to]) {
                    reject({ message: 'Unable to move tag' });
                    return;
                }

                if (!self.transactionTagListStateInvalid) {
                    self.updateTransactionTagListInvalidState(true);
                }

                updateTagDisplayOrderInTransactionTagList(self, {
                    tag: tag,
                    from: from,
                    to: to
                });

                resolve();
            });
        },
        updateTagDisplayOrders() {
            const self = this;
            const newDisplayOrders = [];

            for (let i = 0; i < self.allTransactionTags.length; i++) {
                newDisplayOrders.push({
                    id: self.allTransactionTags[i].id,
                    displayOrder: i + 1
                });
            }

            return new Promise((resolve, reject) => {
                services.moveTransactionTag({
                    newDisplayOrders: newDisplayOrders
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to move tag' });
                        return;
                    }

                    if (self.transactionTagListStateInvalid) {
                        self.updateTransactionTagListInvalidState(false);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save tags display order', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to move tag' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        hideTag({ tag, hidden }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.hideTransactionTag({
                    id: tag.id,
                    hidden: hidden
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (hidden) {
                            reject({ message: 'Unable to hide this tag' });
                        } else {
                            reject({ message: 'Unable to unhide this tag' });
                        }

                        return;
                    }

                    updateTagVisibilityInTransactionTagList(self, {
                        tag: tag,
                        hidden: hidden
                    });

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to change tag visibility', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (hidden) {
                            reject({ message: 'Unable to hide this tag' });
                        } else {
                            reject({ message: 'Unable to unhide this tag' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        deleteTag({ tag, beforeResolve }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.deleteTransactionTag({
                    id: tag.id
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to delete this tag' });
                        return;
                    }

                    if (beforeResolve) {
                        beforeResolve(() => {
                            removeTagFromTransactionTagList(self, tag);
                        });
                    } else {
                        removeTagFromTransactionTagList(self, tag);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to delete tag', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to delete this tag' });
                    } else {
                        reject(error);
                    }
                });
            });
        }
    }
});
