import { Injectable } from '@angular/core';
import { GithubOrganizationService } from './graphQL/githubOrganization/github-organization.service';
import { GithubRepositoryService } from './graphQL/githubRepository/github-repository.service';
import { RepositoryBuildsService } from './graphQL/repositoryBuilds/repository-builds.service';
import { RepositoryService } from './graphQL/repositoryService/repository.service';
import { TagService } from './graphQL/tagService/tag.service';
import { DashboardListItem } from '../models/dashboard-list-item.model';
import { DashboardListItemComponent } from '../components/dashboard-list-item/dashboard-list-item.component';
import { mergeMap, flatMap, map, mergeAll, tap, concatAll } from 'rxjs/operators';
import { Observable, scheduled, asyncScheduler, forkJoin } from 'rxjs';
import { GithubRpositories } from '../models/graphql-models/github-repositories.model';
import { EnqueueBuild, EnqueueBuildReturn } from '../models/graphql-models/enqueue-build.model';
import { RegisterApp } from '../models/graphql-models/register-app.model';
import { RegisterAppService } from './graphQL/registerApp/register-app.service';
import { RepositoryBuild } from '../models/graphql-models/repository-build.model';
import { RegistryPasswordService } from './graphQL/registryPassword/registry-password.service';

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
    return this.repositoryBuildsService.getRepositoryBuilds(repositoryName);
  }

  getRepositories() {
    return this.repositoryService.getRepositories();
  }

  getTags(repositoryName: string) {
    return this.tagService.getTags(repositoryName);
  }

  createDashboardData(){
    return this.getRepositories().pipe(
      map(repos => repos.map(repo => 
          forkJoin(
            this.getRepositoryBuilds(repo.name),
            this.getTags(repo.name)
          )
          .pipe(
            map(([builds, tags]) => ({
            builds, images: tags,name: repo.name
          })
          ),)
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
    return this.registerAppService.registerApp(data).pipe(tap(_=>this.getAllGithubRepositories()));
  }

  enqueueBuild(enqueueData: EnqueueBuild): Observable<EnqueueBuildReturn> {
    return this.repositoryBuildsService.enqueueBuild(enqueueData).pipe(tap(_=> this.getAllBuilds()));
    
  }

  getAllBuilds(): Observable<RepositoryBuild[]> {
    return this.getRepositories().pipe(
      map(repos => repos.map(repo => this.getRepositoryBuilds(repo.name).pipe(concatAll()))),
      map(repos => forkJoin(repos)),
      mergeAll(1),
    )
  }

  setPassword(password: string) {
    this.registryPasswordService.setRegistryPassword(password);
  }

}