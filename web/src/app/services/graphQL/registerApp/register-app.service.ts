import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';
import { Repository } from 'src/app/models/graphql-models/repository.model';
import { RegisterApp } from 'src/app/models/graphql-models/register-app.model';

@Injectable({
  providedIn: 'root'
})
export class RegisterAppService {

  constructor(private apollo: Apollo) { }

  registerApp(registerData: RegisterApp) {
    return this.apollo.mutate({
      mutation: gql`
        mutation manuallyRegisterApp($clientId: String, $clientSecret: String, $privateKey: String, $appId: String, $installationId: String){
          manuallyRegisterApp(
            clientId: $clientId
            clientSecret: $clientSecret
            privateKey: $privateKey
            appId: $appId
            installationId: $installationId
          ) {
            status
          }
        }
      `,
      variables: { ...registerData }
    })
      .pipe(map((result: ApolloQueryResult<{ manuallyRegisterApp: string  }>) => result.data.manuallyRegisterApp))
  }
}
