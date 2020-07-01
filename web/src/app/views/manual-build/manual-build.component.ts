import {Component, OnInit} from '@angular/core';
import {FormModel} from 'src/app/models/form.model';
import {Location} from '@angular/common';
import {EnqueueBuild} from 'src/app/models/graphql-models/enqueue-build.model';
import {GraphQlService} from 'src/app/services/graphql.service';
import {GithubRpositories} from 'src/app/models/graphql-models/github-repositories.model';

@Component({
  selector: 'app-manual-build',
  templateUrl: './manual-build.component.html',
  styleUrls: ['./manual-build.component.scss']
})
export class ManualBuildComponent implements OnInit {
  githubRepositories: Array<GithubRpositories> = [];

  public formModel: FormModel = {
    header: "Manual Build",
    items: [
      { name: "scmId", placeholder: "Source Control Repo", selections: [] },
      { name: "commit", placeholder: "Commit Hash" },
      { name: "dockerfile", placeholder: "Dockerfile Path" },
      { name: "repository", placeholder: "Repository", selections: [] },
      { name: "tag", placeholder: "Image Tag" },
    ],
  }

  constructor(private _location: Location, private graphQlService: GraphQlService) {
  }

  async ngOnInit(): Promise<void> {
    this.graphQlService.getRepositories()
      .subscribe(repos => this.formModel.items[3].selections = repos.map(repo => repo.name));
    this.githubRepositories = await this.graphQlService.getAllGithubRepositories().toPromise();
    this.formModel.items[0].selections = this.githubRepositories.map(repo => repo.name);
  }

  onCancel() {
    this._location.back();
  }

  onSubmit(data: EnqueueBuild) {
    this.graphQlService.enqueueBuild(data).subscribe(console.log);

    this._location.back();
  }
}
