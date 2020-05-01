import { ListContentModel } from './list-item';

export interface ListModel {
    text: string;
    content?: ListContentModel;
}