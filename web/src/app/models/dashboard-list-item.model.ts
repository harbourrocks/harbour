import { RepositoryBuild } from './graphql-models/repository-build.model';
import { Tag } from './graphql-models/tags.model';

export interface DashboardListItem {
    name: string;
    builds: Array<RepositoryBuild>;
    images: Array<Tag>;
}