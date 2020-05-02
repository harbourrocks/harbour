import { ListContentModel } from './list-item';
import { ListDetailsModel } from './list-detail.model';

export interface ListModel {
    text: string;
    content?: ListContentModel;
    details?: Array<ListDetailsModel>;
}