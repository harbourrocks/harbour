import { Component, OnInit, Input } from '@angular/core';
import { DashboardListItem } from 'src/app/models/dashboard-list-item.model';

@Component({
  selector: 'app-dashboard-list-item',
  templateUrl: './dashboard-list-item.component.html',
  styleUrls: ['./dashboard-list-item.component.scss']
})
export class DashboardListItemComponent implements OnInit {
  @Input() model: DashboardListItem;

  constructor() { }

  ngOnInit(): void {
    if(!this.model){
      this.model = {
        builds: [],
        images: [],
        name: "Frontend"
      }
    }
  }

}
