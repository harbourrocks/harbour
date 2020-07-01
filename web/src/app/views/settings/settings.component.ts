import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { GraphQlService } from 'src/app/services/graphql.service';
import { map } from 'rxjs/operators';
import { SimpleListItem } from 'src/app/models/simple-list-item.model';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {

  @ViewChild('passwordInput') passwordInput: ElementRef;

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

  passwordSet(){
    const value = this.passwordInput.nativeElement.value;
    if(value.length === 0) return;
    this.graphQlService.setPassword(value).subscribe(console.log);
  }

}
