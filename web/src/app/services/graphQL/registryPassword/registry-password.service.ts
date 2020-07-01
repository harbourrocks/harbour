import { Injectable } from '@angular/core';
import { Apollo } from 'apollo-angular';
import { RegisterApp } from 'src/app/models/graphql-models/register-app.model';
import gql from 'graphql-tag';
import { map } from 'rxjs/operators';
import { ApolloQueryResult } from 'apollo-client';
import { SetPassword } from 'src/app/models/graphql-models/set-password.model';

@Injectable({
  providedIn: 'root'
})
export class RegistryPasswordService {

  constructor(private apollo: Apollo) { }

  setRegistryPassword(password: string) {
    return this.apollo.mutate({
      mutation: gql`
        mutation setRegistryPassword($password: String){
          setRegistryPassword(password: $password){
            passwordSet
            username
          }
        }
      `,
      variables: { password }
    })
      .pipe(map((result: ApolloQueryResult<{ setRegistryPassword: SetPassword  }>) => result.data.setRegistryPassword))
  }
}
