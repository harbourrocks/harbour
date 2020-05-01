import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot, UrlTree} from '@angular/router';
import {OidcClient} from "oidc-client";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class OidcCallbackGuard implements CanActivate {
  constructor(private router: Router) {
  }

  async canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Promise<boolean | UrlTree> {

    // assemble a client
    const client = new OidcClient({
      authority: environment.oidcDiscoverUrl,
      client_id: environment.oidcClientId,
      redirect_uri: environment.oidcRedirectUrl,
      response_mode: 'query',
      scope: 'openid',
    });

    // request id_token (finish code flow)
    const response = await client.processSigninResponse();
    console.debug('OIDC response', response);

    // redirect to error page on error
    if (response.error != null && response.error != "") {
      return this.router.parseUrl(`/auth/oidc/failed?error=${response.error}&errorDescription=${response.error_description}`);
    }

    // save id_token to local storage
    console.debug('IdToken', response.id_token);
    localStorage.setItem('iam-id_token', response.id_token);

    return this.router.parseUrl('');
  }

}
