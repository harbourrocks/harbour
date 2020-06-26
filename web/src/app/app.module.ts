import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { RepositoryComponent } from './views/repository/repository.component';
import { WizardComponent } from './views/wizard/wizard.component';
import { SettingsComponent } from './views/settings/settings.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { GroupsComponent } from './components/groups/groups.component';
import { ListComponent } from './components/list/list.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { ListDirective } from './components/list/list.directive';
import { AddgroupComponent } from './components/addgroup/addgroup.component';
import { RepositoryDetailComponent } from './views/repository-detail/repository-detail.component';
import { ListDetailsComponent } from './components/list-details/list-details.component';
import { GraphQLModule } from './services/graphQL/graphql.module';
import { HttpClientModule } from '@angular/common/http';
import { ButtonComponent } from './components/button/button.component';
import { RouterIconComponent } from './components/router-icon/router-icon.component';
import { DashboardListItemComponent } from './components/dashboard-list-item/dashboard-list-item.component';
import { SimpleListComponent } from './components/simple-list/simple-list.component';
import { VerticalNavbarComponent } from './components/vertical-navbar/vertical-navbar.component';
import { BuildsComponent } from './views/builds/builds.component';

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    RepositoryComponent,
    WizardComponent,
    SettingsComponent,
    NavbarComponent,
    GroupsComponent,
    ListComponent,
    ListDirective,
    AddgroupComponent,
    RepositoryDetailComponent,
    ListDetailsComponent,
    ButtonComponent,
    RouterIconComponent,
    DashboardListItemComponent,
    SimpleListComponent,
    VerticalNavbarComponent,
    BuildsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule,
    GraphQLModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
