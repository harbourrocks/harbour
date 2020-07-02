import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {GraphQlService} from 'src/app/services/graphql.service';
import {SimpleListItem} from 'src/app/models/simple-list-item.model';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {List} from 'src/app/models/list.model';
import {BuildStatus} from 'src/app/models/build-status.enum';
import { convertToDDMMYYY } from 'src/app/helper/date';

@Component({
  selector: 'app-repository-details',
  templateUrl: './repository-details.component.html',
  styleUrls: ['./repository-details.component.scss']
})
export class RepositoryDetailsComponent implements OnInit {
  navigationLabels = ["Images", "Builds"];
  repositoryName: string = "";
  tags: Observable<SimpleListItem[]>;
  builds: Observable<List>;

  currentPageIndex = 0;

  constructor(private route: ActivatedRoute, private graphQlService: GraphQlService) {
  }

  ngOnInit(): void {
    this.repositoryName = this.route.snapshot.paramMap.get('repo_name');
    this.tags = this.graphQlService.getTags(this.repositoryName)
      .pipe(map(tags => tags.map(tag => ({ label: tag.name }))));

    this.builds = this.graphQlService.getRepositoryBuilds(this.repositoryName)
      .pipe(map(builds =>
        ({
          listItems: builds
            .sort((a, b) => a?.startTime - b?.startTime)
            .map(build => ({
              label: `${build.repository}:${build.tag}`,
              preLabel: `#${build.buildId.substr(0, 18)}`,
              sufLabel: convertToDDMMYYY(build.timestamp),
              color: BuildStatus[build.buildStatus]
            }))
        })))
  }

  pageChange(newIndex) {
    this.currentPageIndex = newIndex;
  }

}
