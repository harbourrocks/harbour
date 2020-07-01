// This file can be replaced during build by using the `fileReplacements` array.
// `ng build --prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  oidcDiscoverUrl: 'https://login.microsoftonline.com/07a987a0-3eef-42a9-a7d6-53698899fcf2/v2.0/.well-known/openid-configuration',
  oidcClientId: 'd8dfd041-5a71-465d-8d15-7d2b91c0b1ba',
  oidcRedirectUrl: 'http://localhost:4200/auth/oidc/callback',
  graphQlUrl: 'http://localhost:5400/graphql'
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.
