import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';
import { Repository } from 'src/app/models/graphql-models/repository.model';

@Injectable({
  providedIn: 'root'
})
export class RepositoryService {

  constructor(private apollo: Apollo) { }

  getRepositories() {
    return this.apollo.query({
      query: gql`
        {
        repositories {
          name
          }
        }
      `,
    })
      .pipe(map((result: ApolloQueryResult<{ repositories: Array<Repository> }>) => result.data.repositories))
  }
}
