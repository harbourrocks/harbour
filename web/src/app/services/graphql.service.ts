import { Injectable } from '@angular/core';
import { GithubOrganizationService } from './graphQL/githubOrganization/github-organization.service';
import { GithubRepositoryService } from './graphQL/githubRepository/github-repository.service';
import { RepositoryBuildsService } from './graphQL/repositoryBuilds/repository-builds.service';
import { RepositoryService } from './graphQL/repositoryService/repository.service';
import { TagService } from './graphQL/tagService/tag.service';

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
    ) { }

  getGithubOrganizations() {
    return this.githubOrganizationService.getGithubOrganizations();
  }

  getGithubRepositories(login: string) {
    return this.githubRepositoryService.getGithubRepositories(login);
  }
  
  getRepositoryBuilds(scmdId: string, repositoryName: string) {
    return this.repositoryBuildsService.getRepositoryBuilds(scmdId,repositoryName);
  }

  getRepositories() {
    return this.repositoryService.getRepositories();
  }

  getTags(repositoryName: string) {
    return this.tagService.getTags(repositoryName);
  }



}