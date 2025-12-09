import {
    isObject,
    isString,
    isBoolean,
    isNumber
} from '@/lib/common.ts';

// Project represents project data stored in database
export class Project {
    public id: string;
    public name: string;
    public color: string;
    public comment: string;
    public displayOrder: number;
    public hidden: boolean;

    public constructor(id: string, name: string, color: string, comment: string, displayOrder: number, hidden: boolean) {
        this.id = id;
        this.name = name;
        this.color = color;
        this.comment = comment;
        this.displayOrder = displayOrder;
        this.hidden = hidden;
    }

    public static of(data: unknown): Project {
        if (!isObject(data)) {
            throw new Error('Project data is invalid');
        }

        const d = data as Record<string, unknown>;

        const id = isString(d['id']) ? d['id'] : '';
        const name = isString(d['name']) ? d['name'] : '';
        const color = isString(d['color']) ? d['color'] : '';
        const comment = isString(d['comment']) ? d['comment'] : '';
        const displayOrder = isNumber(d['displayOrder']) ? d['displayOrder'] : 0;
        const hidden = isBoolean(d['hidden']) ? d['hidden'] : false;

        return new Project(id, name, color, comment, displayOrder, hidden);
    }
}

// ProjectCreateRequest represents all parameters of single project creation request
export interface ProjectCreateRequest {
    name: string;
    color: string;
    comment: string;
    clientSessionId?: string;
}

// ProjectModifyRequest represents all parameters of project modification request
export interface ProjectModifyRequest {
    id: string;
    name: string;
    color: string;
    comment: string;
    hidden: boolean;
}

// ProjectHideRequest represents all parameters of project hiding request
export interface ProjectHideRequest {
    id: string;
    hidden: boolean;
}

// ProjectMoveRequest represents all parameters of project moving request
export interface ProjectMoveRequest {
    newDisplayOrders: ProjectNewDisplayOrderRequest[];
}

// ProjectNewDisplayOrderRequest represents a data pair of id and display order
export interface ProjectNewDisplayOrderRequest {
    id: string;
    displayOrder: number;
}

// ProjectDeleteRequest represents all parameters of project deleting request
export interface ProjectDeleteRequest {
    id: string;
}

// ProjectInfoResponse represents a view-object of project
export interface ProjectInfoResponse {
    id: string;
    name: string;
    color: string;
    comment: string;
    displayOrder: number;
    hidden: boolean;
}
