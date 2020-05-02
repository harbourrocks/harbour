import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { RepositoryComponent } from './views/repository/repository.component';
import { WizardComponent } from './views/wizard/wizard.component';
import { SettingsComponent } from './views/settings/settings.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { SearchbarComponent } from './components/searchbar/searchbar.component';
import { GroupsComponent } from './components/groups/groups.component';
import { ListComponent } from './components/list/list.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { ListDirective } from './components/list/list.directive';
import { AddgroupComponent } from './components/addgroup/addgroup.component';
import { RepositoryDetailComponent } from './views/repository-detail/repository-detail.component';
import { ListDetailsComponent } from './components/list-details/list-details.component';

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    RepositoryComponent,
    WizardComponent,
    SettingsComponent,
    NavbarComponent,
    SearchbarComponent,
    GroupsComponent,
    ListComponent,
    ListDirective,
    AddgroupComponent,
    RepositoryDetailComponent,
    ListDetailsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
