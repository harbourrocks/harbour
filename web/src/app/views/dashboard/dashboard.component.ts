import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { ListContentModel } from 'src/app/models/list-item';
import { SearchbarComponent } from 'src/app/components/searchbar/searchbar.component';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public groupArr: Array<GroupModel>;

  constructor() {
    this.groupArr = [
      {
        listItems: [
          {
            text: "Web Systems Frontend",
          },
          {
            text: "Web Systems Backend",
          },
          {
            text: "Custom Docker Registry",
            content: new ListContentModel(SearchbarComponent, ""),
          },
          {
            text: "Search Controller",
          },
        ]
      },

    ]
  }

  ngOnInit(): void {
  }

}
