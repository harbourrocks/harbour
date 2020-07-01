import { Component, OnInit } from '@angular/core';
import { GraphQlService } from 'src/app/services/graphql.service';
import { DashboardListItem } from 'src/app/models/dashboard-list-item.model';
import { BuildStatus } from 'src/app/models/build-status.enum';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public models: Observable<Array<DashboardListItem>>;

  constructor(private graphQlService: GraphQlService) {}

  ngOnInit(): void {
    this.models = this.graphQlService.createDashboardData();
  }

}
