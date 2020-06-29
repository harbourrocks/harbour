import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { SettingsComponent } from './views/settings/settings.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { GraphQLModule } from './services/graphQL/graphql.module';
import { HttpClientModule } from '@angular/common/http';
import { ButtonComponent } from './components/button/button.component';
import { RouterIconComponent } from './components/router-icon/router-icon.component';
import { DashboardListItemComponent } from './components/dashboard-list-item/dashboard-list-item.component';
import { SimpleListComponent } from './components/simple-list/simple-list.component';
import { VerticalNavbarComponent } from './components/vertical-navbar/vertical-navbar.component';
import { BuildsComponent } from './views/builds/builds.component';
import { RepositoryComponent } from './views/repository/repository.component';
import { ListComponent } from './components/list/list.component';
import { ListItemComponent } from './components/list-item/list-item.component';
import { RepositoryDetailsComponent } from './views/repository-details/repository-details.component';
import { AbstractFormComponent } from './components/abstract-form/abstract-form.component';
import { ManualBuildComponent } from './views/manual-build/manual-build.component';
import { AddAccountComponent } from './views/add-account/add-account.component';

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    SettingsComponent,
    NavbarComponent,
    ButtonComponent,
    RouterIconComponent,
    DashboardListItemComponent,
    SimpleListComponent,
    VerticalNavbarComponent,
    BuildsComponent,
    RepositoryComponent,
    ListComponent,
    ListItemComponent,
    RepositoryDetailsComponent,
    AbstractFormComponent,
    ManualBuildComponent,
    AddAccountComponent
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
