export interface RepositoryBuild {
    buildStatus: BuildStatus,
    commit: string,
    timestamp: number,
}

enum BuildStatus {
    "Pending",
    "Running",
    "Success",
    "Failed"
}