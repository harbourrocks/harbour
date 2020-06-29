import { Injectable } from '@angular/core';
import { GithubOrganizationService } from './graphQL/githubOrganization/github-organization.service';
import { GithubRepositoryService } from './graphQL/githubRepository/github-repository.service';
import { RepositoryBuildsService } from './graphQL/repositoryBuilds/repository-builds.service';
import { RepositoryService } from './graphQL/repositoryService/repository.service';
import { TagService } from './graphQL/tagService/tag.service';
import { DashboardListItem } from '../models/dashboard-list-item.model';
import { DashboardListItemComponent } from '../components/dashboard-list-item/dashboard-list-item.component';

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
  
  getRepositoryBuildsWithScmId(scmdId: string, repositoryName: string) {
    return this.repositoryBuildsService.getRepositoryBuildsWithScmId(scmdId,repositoryName);
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

  async createDashboardData() {
    const githubOrganizations = await this.getGithubOrganizations().toPromise();
    githubOrganizations.forEach(async orga => {
      const githubRepos = await this.getGithubRepositories(orga.login).toPromise();
      

    })


    const item : DashboardListItem = {
      builds: [] ,
      images: [], // tags
      name:"" ,

    }

  }

}