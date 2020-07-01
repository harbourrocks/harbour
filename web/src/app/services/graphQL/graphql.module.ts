import { NgModule } from '@angular/core';
import { ApolloModule, APOLLO_OPTIONS } from 'apollo-angular';
import { ApolloLink } from 'apollo-link';
import { setContext } from 'apollo-link-context';
import { HttpLinkModule, HttpLink } from 'apollo-angular-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpClientModule } from '@angular/common/http';
import {environment} from "../../../environments/environment";

export function createApollo(httpLink: HttpLink) {

  const token = localStorage.getItem('iam-id_token');
  const auth = setContext((operation, context) => ({
    headers: {
      Authorization: `Bearer ${token}`
    },
  }));

  const link = ApolloLink.from([auth, httpLink.create({ uri: environment.graphQlUrl })]);
  const cache = new InMemoryCache();

  return {
    link,
    cache
  }
}


@NgModule({
  exports: [
    HttpClientModule,
    ApolloModule,
    HttpLinkModule],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink],
    },
  ],
})
export class GraphQLModule { }
