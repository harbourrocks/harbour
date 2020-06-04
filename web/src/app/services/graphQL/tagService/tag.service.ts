import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { map } from 'rxjs/operators';
import gql from 'graphql-tag';
import { ApolloQueryResult } from 'apollo-client';
import { Tag } from 'src/app/models/graphql-models/tags.model';

@Injectable({
  providedIn: 'root'
})
export class TagService {

  constructor(private apollo: Apollo) { }

  // getTags(repositoryName: string) {
  //   return this.apollo.query({
  //     query: gql`
  //       {
  //         tags(repository: $repositoryName){
  //           name
  //         }
  //       }
  //     `,
  //     variables: {repositoryName: repositoryName},
      
  //   })
  //   .pipe(map((result: ApolloQueryResult<{tags: Array<Tag>}>) => result.data.tags))
  getTags(repositoryName: string) {
    return this.apollo.query({
      query: gql`
        query getTags($repository: String){
          tags(repository: $repository){
            name
          }
        }
      `,
      variables: {repository: repositoryName},
      
    })
    .pipe(map((result: ApolloQueryResult<{tags: Array<Tag>}>) => result.data.tags))
  }
}
