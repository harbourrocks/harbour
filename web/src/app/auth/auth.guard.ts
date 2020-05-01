import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree} from '@angular/router';
import {environment} from "../../environments/environment";
import {OidcClient} from "oidc-client";

function MakeRandomId(length): string {
  let result = '';
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;

  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }

  return result;
}

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {
  async canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Promise<boolean | UrlTree> {

    const id_token: string = localStorage.getItem("iam-id_token");
    console.debug(`id_token: ${id_token}`);

    if (id_token != null && id_token.length > 0) return true;

    // assemble a client
    const client = new OidcClient({
      authority: environment.oidcDiscoverUrl,
      client_id: environment.oidcClientId,
      redirect_uri: environment.oidcRedirectUrl,
      response_mode: 'query',
      scope: 'openid',
      response_type: 'code',
      extraQueryParams: {
        login_hint: 'oliver.seitz@harbour.rocks'
      }
    });

    const signInRequest = await client.createSigninRequest();
    console.debug('OIDC request', signInRequest);

    // redirect to OpenId Provider
    window.location.assign(signInRequest.url);

    // unreachable
    return false;
  }

}
