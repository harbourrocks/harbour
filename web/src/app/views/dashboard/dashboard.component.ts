import { Component, OnInit } from '@angular/core';
import { GraphQlService } from 'src/app/services/graphql.service';
import { DashboardListItem } from 'src/app/models/dashboard-list-item.model';
import { BuildStatus } from 'src/app/models/build-status.enum';

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
