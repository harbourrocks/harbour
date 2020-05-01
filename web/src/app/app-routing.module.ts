import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AppComponent } from './app.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { RepositoryComponent } from './views/repository/repository.component';
import { SettingsComponent } from './views/settings/settings.component';
import { WizardComponent } from './views/wizard/wizard.component';
import { RepositoryDetailComponent } from './views/repository-detail/repository-detail.component';


const routes: Routes = [
  { path: 'dashboard', component: DashboardComponent },
  { path: 'repository', component: RepositoryComponent },
  { path: 'repository/:repoId', component: RepositoryDetailComponent },
  { path: 'settings', component: SettingsComponent },
  { path: 'wizard', component: WizardComponent },
  { path: '**', component: DashboardComponent },
  
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
