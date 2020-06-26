export interface RepositoryBuild {
    buildStatus: BuildStatus,
    commit: string,
    timestamp: number,
}

export enum BuildStatus {
    "Pending",
    "Running",
    "Success",
    "Failed"
}