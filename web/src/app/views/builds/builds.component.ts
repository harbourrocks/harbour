import { Component, OnInit } from '@angular/core';
import { GraphQlService } from 'src/app/services/graphql.service';
import { Observable } from 'rxjs';
import { List } from 'src/app/models/list.model';
import { map, timestamp } from 'rxjs/operators';
import { BuildStatus } from 'src/app/models/build-status.enum';
import { convertToDDMMYYY } from 'src/app/helper/date';

@Component({
  selector: 'app-builds',
  templateUrl: './builds.component.html',
  styleUrls: ['./builds.component.scss']
})
export class BuildsComponent implements OnInit {
  builds: Observable<List>;

  constructor(private graphQlService: GraphQlService) {
  }

  ngOnInit(): void {
    this.graphQlService.getAllBuilds().subscribe(console.log)
    this.builds = this.graphQlService.getAllBuilds()
      .pipe(
        map(builds =>
          ({
            listItems: builds
              .sort((a, b) => a?.startTime - b?.startTime)
              .map(build => ( {
                  label: `${build.repository}:${build.tag}`,
                  preLabel: `#${build.buildId.substr(0, 18)}`,
                  sufLabel: convertToDDMMYYY(build.timestamp),
                  color: BuildStatus[build.buildStatus]
                }
              ))
          })))
  }
}
