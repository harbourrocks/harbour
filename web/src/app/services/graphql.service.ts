import { Injectable } from '@angular/core';
import { GithubOrganizationService } from './graphQL/githubOrganization/github-organization.service';
import { GithubRepositoryService } from './graphQL/githubRepository/github-repository.service';
import { RepositoryBuildsService } from './graphQL/repositoryBuilds/repository-builds.service';
import { RepositoryService } from './graphQL/repositoryService/repository.service';
import { TagService } from './graphQL/tagService/tag.service';
import { mergeMap, flatMap, map, mergeAll, tap, concatAll, concatMap, toArray, mapTo } from 'rxjs/operators';
import { Observable, forkJoin } from 'rxjs';
import { GithubRpositories } from '../models/graphql-models/github-repositories.model';
import { EnqueueBuild, EnqueueBuildReturn } from '../models/graphql-models/enqueue-build.model';
import { RegisterApp } from '../models/graphql-models/register-app.model';
import { RegisterAppService } from './graphQL/registerApp/register-app.service';
import { RepositoryBuild } from '../models/graphql-models/repository-build.model';
import { RegistryPasswordService } from './graphQL/registryPassword/registry-password.service';
import { subscribe } from 'graphql';

@Injectable({
  providedIn: 'root'
})
export class GraphQlService {
  constructor(
    private githubOrganizationService: GithubOrganizationService,
    private githubRepositoryService: GithubRepositoryService,
    private repositoryBuildsService: RepositoryBuildsService,
    private repositoryService: RepositoryService,
    private tagService: TagService,
    private registerAppService: RegisterAppService,
    private registryPasswordService: RegistryPasswordService
  ) { }

  getGithubOrganizations() {
    return this.githubOrganizationService.getGithubOrganizations();
  }

  getGithubRepositories(login: string) {
    return this.githubRepositoryService.getGithubRepositories(login);
  }

  getRepositoryBuildsWithScmId(scmdId: string, repositoryName: string) {
    return this.repositoryBuildsService.getRepositoryBuildsWithScmId(scmdId, repositoryName);
  }

  getRepositoryBuilds(repositoryName: string) {
    return this.repositoryBuildsService.getRepositoryBuilds(repositoryName).pipe(map(builds => builds.sort((a, b) => a.timestamp - b.timestamp)));
  }

  getRepositories() {
    return this.repositoryService.getRepositories();
  }

  getTags(repositoryName: string) {
    return this.tagService.getTags(repositoryName);
  }

  createDashboardData() {
    return this.getRepositories().pipe(
      map(repos => repos.map(repo =>
        forkJoin(
          this.getRepositoryBuilds(repo.name),
          this.getTags(repo.name)
        )
          .pipe(
            map(([builds, tags]) => ({
              builds, images: tags, name: repo.name
            })
            ))
      )),
      map(obs => forkJoin(obs)),
      mergeAll(1)
    )
  }

  getAllGithubRepositories(): Observable<Array<GithubRpositories>> {
    return this.getGithubOrganizations().pipe(
      mergeMap(orgas => orgas.map(orga => this.getGithubRepositories(orga.login))),
      mergeAll(1)
    )
  }

  addGithubAccount(data: RegisterApp): Observable<string> {
    return this.registerAppService.registerApp(data);
  }

  enqueueBuild(enqueueData: EnqueueBuild): Observable<EnqueueBuildReturn> {
    return this.repositoryBuildsService.enqueueBuild(enqueueData);
  }

  getAllBuilds(): Observable<RepositoryBuild[]> {
    return this.getRepositories().pipe(
      map(repos => repos.map(repo => this.getRepositoryBuilds(repo.name))),
      map(builds => forkJoin(builds)),
      mergeAll(1),
      map(builds => builds.reduce((a, b) => a.concat(b).sort((a, b) => a.timestamp - b.timestamp)))
    )
  }

  setPassword(password: string) {
    return this.registryPasswordService.setRegistryPassword(password);
  }

}
