import { BuildStatus } from '../build-status.enum';

export interface RepositoryBuild {
    buildStatus?: BuildStatus,
    commit?: string,
    timestamp?: number,
    endTime?: number,
    repository?: string,
    scmId?: string,
    startTime?: number,
    tag?: string,
}

