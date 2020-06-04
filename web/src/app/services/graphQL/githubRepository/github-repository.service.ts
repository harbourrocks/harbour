import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';
import { GithubRpositories } from 'src/app/models/graphql-models/github-repositories.model';

@Injectable({
  providedIn: 'root'
})
export class GithubRepositoryService {

  constructor(private apollo: Apollo) { }

  getGithubRepositories(login: string) {
    return this.apollo.query({
      query: gql`
        query getGithubRepositories($login: String){
          githubRepositories(orgLogin: $login){
            id
            name
            scm_id
          }
        }
      `,
      variables: {"login": login},

    })
      .pipe(map((result: ApolloQueryResult<{ githubRepositories: Array<GithubRpositories> }>) => result.data.githubRepositories))
  }
}
