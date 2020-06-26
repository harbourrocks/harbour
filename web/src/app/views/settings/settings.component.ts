import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { faGithub, faGitlab } from '@fortawesome/free-brands-svg-icons';
import { GraphQlService } from 'src/app/services/graphql.service';
import { map, flatMap } from 'rxjs/operators';
import { SimpleListItem } from 'src/app/models/simple-list-item.model';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  navigationLabels = ["Accounts", "Passwords"];
  listItems: Observable<SimpleListItem[]>;

  currentPageIndex: number = 0;

  constructor(private graphQlService: GraphQlService) {  }

  ngOnInit(): void {
    this.listItems = this.graphQlService.getGithubOrganizations()
      .pipe(
        map(orga => orga.map(org=> ({icon: org.avatarUrl, label: org.name})))
      )
  }

  pageChange(newPageIndex: number) {
    this.currentPageIndex = newPageIndex;
  }

}
