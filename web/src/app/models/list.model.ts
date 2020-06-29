import { ListItem } from './list-item.model';

export interface List {
    listItems: Array<ListItem>;
    clickHandler?: ListClickHandler;
}

export type ListClickHandler = (listItem: ListItem) => void; 