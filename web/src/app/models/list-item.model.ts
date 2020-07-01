import { BuildStatus } from './build-status.enum';

export interface ListItem {
    label: string,
    preLabel?: string,
    color?: string,
    sufLabel?: string,
    clickable?: boolean,
}