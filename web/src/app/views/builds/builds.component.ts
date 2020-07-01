import { Component, OnInit } from '@angular/core';
import { GraphQlService } from 'src/app/services/graphql.service';
import { Observable } from 'rxjs';
import { List } from 'src/app/models/list.model';
import { map, mergeAll } from 'rxjs/operators';
import { BuildStatus } from 'src/app/models/build-status.enum';

@Component({
  selector: 'app-builds',
  templateUrl: './builds.component.html',
  styleUrls: ['./builds.component.scss']
})
export class BuildsComponent implements OnInit {
  builds: Observable<List>;

  constructor(private graphQlService: GraphQlService) { }

  ngOnInit(): void {
    this.graphQlService.getAllBuilds().subscribe(console.log)

    this.builds = this.graphQlService.getAllBuilds()
      .pipe(map(builds =>
        ({
          listItems: builds.map(build => ({
            label: `${build.repository}:${build.tag}`,
            preLabel: `#${build.commit}`,
            sufLabel: new Date(1593611348 * 1000).toISOString().substring(0, 10),
            color: BuildStatus[build.buildStatus]

          }))
        })))
  }

}
