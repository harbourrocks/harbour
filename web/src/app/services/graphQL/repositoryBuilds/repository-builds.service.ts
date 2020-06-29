import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';
import { RepositoryBuild } from 'src/app/models/graphql-models/repository-build.model';

@Injectable({
  providedIn: 'root'
})
export class RepositoryBuildsService {

  constructor(private apollo: Apollo) { }

  getRepositoryBuildsWithScmId(scmdId: string, repositoryName: string) {
    return this.apollo.query({
      query: gql`
      query getRepositoryBuildsWithScmId($scmId: String,$repository: String){
        repositoryBuilds(scmId: $scmdId,  repository: $repository)
        {
          buildStatus,
          commit,
          timestamp
        }
      }
      `,
      variables: { scmId: scmdId, repository: repositoryName },

    })
      .pipe(map((result: ApolloQueryResult<{ repositoryBuilds: Array<RepositoryBuild> }>) => result.data.repositoryBuilds))
  }

  getRepositoryBuilds(repository: string) {
    return this.apollo.query({
      query: gql`
      query getRepositoryBuilds($repository: String){
        repositoryBuilds(repository: $repository){
          buildStatus
          commit
          endTime
          repository
          scmId
          startTime
          tag
          timestamp
        }
      }`,
      variables: { repository: repository },
    })
      .pipe(map((result: ApolloQueryResult<{ repositoryBuilds: Array<RepositoryBuild> }>) => result.data.repositoryBuilds))
  }
}
