import { FormItem } from './form-items.model';

export interface FormModel {
    header: string;
    items: Array<FormItem>;
}