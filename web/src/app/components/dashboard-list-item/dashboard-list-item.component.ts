import { Component, OnInit, Input } from '@angular/core';
import { DashboardListItem } from 'src/app/models/dashboard-list-item.model';
import { BuildStatus } from 'src/app/models/build-status.enum';
import { min } from 'rxjs/operators';

@Component({
  selector: 'app-dashboard-list-item',
  templateUrl: './dashboard-list-item.component.html',
  styleUrls: ['./dashboard-list-item.component.scss']
})
export class DashboardListItemComponent implements OnInit {
  @Input() model: DashboardListItem;
  buildDetails: Array<{
    buildTime: number;
    color: any;
    failed: boolean;
  }>;

  private failedHeight = 17;
  private minHeight = 25;
  private maxHeight = 70;

  constructor() { }

  ngOnInit(): void {
    const details = this.model.builds.map(build => {
      const buildTime = build.endTime - build.startTime;
      const color = BuildStatus[build.buildStatus];
      const failed = BuildStatus[build.buildStatus] === BuildStatus.Failed
      return { buildTime, color, failed }
    })
    this.buildDetails = this.resolveBuildHeight(details);
  }

  resolveBuildHeight(buildDetails: Array<{
    buildTime: number;
    color: any;
    failed: boolean;
  }>) {
    const maxVal = buildDetails.reduce((a, b) => a.buildTime < b.buildTime ? b : a);
    const minVal = buildDetails.reduce((a, b) => a.buildTime > b.buildTime ? b : a);
    const diff = maxVal.buildTime - minVal.buildTime;

    const heightRoom = this.maxHeight - this.minHeight;


    if (diff === 0)
      return buildDetails.map(build => ({ ...build, buildTime: this.minHeight }))
    const height_per_num = heightRoom / diff;

    return buildDetails.map(build => ({ ...build, buildTime: Math.round(build.buildTime * height_per_num) + this.minHeight }))
  }




}
