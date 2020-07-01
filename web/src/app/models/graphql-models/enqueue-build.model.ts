export interface EnqueueBuild {
    tag: string,
    scmId: string,
    commit: string,
    repository: string,
    dockerfile: string,
}

export interface EnqueueBuildReturn {
    buildId: string,
    status: string,
}