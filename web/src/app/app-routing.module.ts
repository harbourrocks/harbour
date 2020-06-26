import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {DashboardComponent} from './views/dashboard/dashboard.component';
import {RepositoryComponent} from './views/repository/repository.component';
import {SettingsComponent} from './views/settings/settings.component';
import {WizardComponent} from './views/wizard/wizard.component';
import {AuthGuard} from "./auth/auth.guard";
import {OidcCallbackGuard} from "./auth/oidc-callback.guard";
import { RepositoryDetailComponent } from './views/repository-detail/repository-detail.component';
import { BuildsComponent } from './views/builds/builds.component';


const routes: Routes = [
  { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
  { path: 'repository', component: RepositoryComponent, canActivate: [AuthGuard] },
  { path: 'repository/:id', component: RepositoryDetailComponent, canActivate: [AuthGuard] },
  { path: 'settings', component: SettingsComponent, canActivate: [AuthGuard] },
  { path: 'builds', component: BuildsComponent, canActivate: [AuthGuard] },
  { path: 'auth/oidc/callback', component: DashboardComponent, canActivate: [OidcCallbackGuard] },
  { path: '**', redirectTo: 'dashboard' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
