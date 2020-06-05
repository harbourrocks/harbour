import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { ListContentModel } from 'src/app/models/list-item';
import { TagService } from 'src/app/services/graphQL/tagService/tag.service';
import { GithubRepositoryService } from 'src/app/services/graphQL/githubRepository/github-repository.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public groupArr: Array<GroupModel>;

  constructor(private githubRepository: GithubRepositoryService) {
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
            content: new ListContentModel(DashboardComponent, ""),
          },
          {
            text: "Search Controller",
          },
        ]
      },

    ]
  }

  ngOnInit(): void {
    this.githubRepository.getGithubRepositories("harbourrocks").subscribe(console.log)
  }

}
