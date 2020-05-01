import { FontawesomeObject, IconDefinition } from '@fortawesome/fontawesome-svg-core';
import { ListModel } from './list.model';

export interface GroupModel {
    title?: string;
    icon?: FontawesomeObject | IconDefinition,
    collapsable?: boolean;
    smallColorbox?: boolean;
    listItems: Array<ListModel>;
}