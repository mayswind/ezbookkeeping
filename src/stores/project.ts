import { Project, type ProjectCreateRequest, type ProjectModifyRequest, type ProjectHideRequest, type ProjectMoveRequest, type ProjectDeleteRequest } from '@/models/project.ts';
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

import { useRootStore } from '@/stores/index.ts';
import { useUserStore } from '@/stores/user.ts';

import { isDefined } from '@/lib/common.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useProjectsStore = defineStore('projects', () => {
    const rootStore = useRootStore();
    const userStore = useUserStore();

    const allProjects = ref<Project[]>([]);
    const projectsMap = ref<Record<string, Project>>({});
    const projectListLoaded = ref<boolean>(false);

    const allVisibleProjects = computed<Project[]>(() => {
        return allProjects.value.filter(project => !project.hidden);
    });

    function getProject(projectId: string): Project | undefined {
        return projectsMap.value[projectId];
    }

    function loadAllProjects({ force }: { force: boolean }): Promise<Project[]> {
        if (!force && projectListLoaded.value) {
            return Promise.resolve(allProjects.value);
        }

        return new Promise((resolve, reject) => {
            services.getProjects().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve project list' });
                    return;
                }

                const projects: Project[] = [];
                const projectMap: Record<string, Project> = {};

                if (data.result) {
                    for (const project of data.result) {
                        const projectInstance = Project.of(project);
                        projects.push(projectInstance);
                        projectMap[projectInstance.id] = projectInstance;
                    }
                }

                projects.sort((a, b) => {
                    if (a.displayOrder !== b.displayOrder) {
                        return a.displayOrder - b.displayOrder;
                    }

                    return a.id.localeCompare(b.id);
                });

                allProjects.value = projects;
                projectsMap.value = projectMap;
                projectListLoaded.value = true;

                resolve(projects);
            }).catch(error => {
                logger.error('failed to load project list', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to retrieve project list' });
                }
            });
        });
    }

    function saveProject({ project }: { project: ProjectCreateRequest | ProjectModifyRequest }): Promise<Project> {
        return new Promise((resolve, reject) => {
            services.saveProject(project).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to save project' });
                    return;
                }

                const savedProject = Project.of(data.result);

                if ('id' in project && project.id) { // Modify
                    for (let i = 0; i < allProjects.value.length; i++) {
                        if (allProjects.value[i].id === savedProject.id) {
                            allProjects.value[i] = savedProject;
                            break;
                        }
                    }

                    projectsMap.value[savedProject.id] = savedProject;
                } else { // Create
                    allProjects.value.push(savedProject);
                    projectsMap.value[savedProject.id] = savedProject;
                }

                allProjects.value.sort((a, b) => {
                    if (a.displayOrder !== b.displayOrder) {
                        return a.displayOrder - b.displayOrder;
                    }

                    return a.id.localeCompare(b.id);
                });

                resolve(savedProject);
            }).catch(error => {
                logger.error('failed to save project', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to save project' });
                }
            });
        });
    }

    function hideProject({ project, hidden }: { project: ProjectHideRequest, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideProject({
                id: project.id,
                hidden: hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to hide project' });
                    return;
                }

                const targetProject = projectsMap.value[project.id];

                if (targetProject) {
                    targetProject.hidden = hidden;
                }

                resolve(true);
            }).catch(error => {
                logger.error('failed to hide project', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to hide project' });
                }
            });
        });
    }

    function moveProject({ project }: { project: ProjectMoveRequest }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.moveProject(project).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move project' });
                    return;
                }

                if (project.newDisplayOrders) {
                    for (const newOrder of project.newDisplayOrders) {
                         const targetProject = projectsMap.value[newOrder.id];
                         if (targetProject) {
                             targetProject.displayOrder = newOrder.displayOrder;
                         }
                    }

                    allProjects.value.sort((a, b) => {
                        if (a.displayOrder !== b.displayOrder) {
                            return a.displayOrder - b.displayOrder;
                        }

                        return a.id.localeCompare(b.id);
                    });
                }

                resolve(true);
            }).catch(error => {
                logger.error('failed to move project', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to move project' });
                }
            });
        });
    }

    function deleteProject({ project }: { project: ProjectDeleteRequest }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteProject(project).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete project' });
                    return;
                }

                const index = allProjects.value.findIndex(p => p.id === project.id);
                if (index >= 0) {
                    allProjects.value.splice(index, 1);
                }

                delete projectsMap.value[project.id];

                resolve(true);
            }).catch(error => {
                logger.error('failed to delete project', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else {
                    reject({ message: 'Unable to delete project' });
                }
            });
        });
    }

    return {
        allProjects,
        projectsMap,
        projectListLoaded,
        allVisibleProjects,
        getProject,
        loadAllProjects,
        saveProject,
        hideProject,
        moveProject,
        deleteProject
    };
});
