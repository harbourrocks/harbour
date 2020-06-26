import { FaSymbol } from '@fortawesome/fontawesome-svg-core';

export interface SimpleListItem {
    label: string,
    afterLabel?: string,
    icon?: string | FaSymbol,
    id?: any,
}