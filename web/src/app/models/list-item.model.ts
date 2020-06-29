import { BuildStatus } from './build-status.enum';

export interface ListItem {
    label: string,
    preLabel?: string,
    status?: BuildStatus,
    sufLabel?: string,
    clickable?: boolean,
}