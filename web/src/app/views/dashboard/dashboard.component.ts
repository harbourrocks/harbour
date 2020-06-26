import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { ListContentModel } from 'src/app/models/list-item';
import { GraphQlService } from 'src/app/services/graphql.service';
import { DashboardListItemComponent } from 'src/app/components/dashboard-list-item/dashboard-list-item.component';
import { DashboardListItem } from 'src/app/models/dashboard-list-item.model';
import { BuildStatus } from 'src/app/models/graphql-models/repository-build.model';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public models: Array<DashboardListItem>;

  constructor(private graphQlService: GraphQlService) {
    // Dummy data
    if(!this.models){
      this.models = [
        {
          builds: [
            {buildStatus: BuildStatus.Failed, commit: "",timestamp:0},
            {buildStatus: BuildStatus.Failed, commit: "",timestamp:0},
            {buildStatus: BuildStatus.Failed, commit: "",timestamp:0},
          ],
          images: [
            {name: "name"},
            {name: "name"},
            {name: "name"},
            {name: "name"},
          ],
          name: "Frontend"
          
        }
      ]
    }
  }

  ngOnInit(): void {
    this.graphQlService.createDashboardData();
  }

}
