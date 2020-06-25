import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { map } from 'rxjs/operators';
import gql from 'graphql-tag';
import { ApolloQueryResult } from 'apollo-client';
import { GithubOrganization } from 'src/app/models/graphql-models/github-organization.model';

@Injectable({
  providedIn: 'root'
})
export class GithubOrganizationService {

  constructor(private apollo: Apollo) { }

  getGithubOrganizations() {
    return this.apollo.query({
      query: gql`
        {
          githubOrganizations{
            avatarUrl
            login
            name
          }
        }
      `,
    })
      .pipe(map((result: ApolloQueryResult<{ githubOrganizations: Array<GithubOrganization> }>) => result.data.githubOrganizations))
  }
}
