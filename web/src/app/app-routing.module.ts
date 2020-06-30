import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {DashboardComponent} from './views/dashboard/dashboard.component';
import {SettingsComponent} from './views/settings/settings.component';
import {AuthGuard} from "./auth/auth.guard";
import {OidcCallbackGuard} from "./auth/oidc-callback.guard";
import { BuildsComponent } from './views/builds/builds.component';
import { RepositoryComponent } from './views/repository/repository.component';
import { RepositoryDetailsComponent } from './views/repository-details/repository-details.component';
import { ManualBuildComponent } from './views/manual-build/manual-build.component';
import { AddAccountComponent } from './views/add-account/add-account.component';


const routes: Routes = [
  { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
  { path: 'repository', component: RepositoryComponent, canActivate: [AuthGuard] },
  { path: 'repository/:repo_name', component: RepositoryDetailsComponent, canActivate: [AuthGuard] },
  { path: 'settings', component: SettingsComponent, canActivate: [AuthGuard] },
  { path: 'builds', component: BuildsComponent, canActivate: [AuthGuard] },
  { path: 'manual-build', component: ManualBuildComponent, canActivate: [AuthGuard] },
  { path: 'add-account/:acc_prov', component: AddAccountComponent, canActivate: [AuthGuard] },
  { path: 'auth/oidc/callback', component: DashboardComponent, canActivate: [OidcCallbackGuard] },
  { path: '**', redirectTo: 'dashboard' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
