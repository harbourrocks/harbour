import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { faGithub, faGitlab } from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  styleUrls: ['./repository.component.scss']
})
export class RepositoryComponent implements OnInit {
  public groupArr: Array<GroupModel>;

  constructor() {
    this.groupArr = [
      {
        smallColorbox: true,
        title: "Github",
        icon: faGithub,
        collapsable: true,
        listItems: [
          {
            text: "github.com",
          },
        ]
      },
      {
        smallColorbox: true,
        title: "Gitlab",
        icon: faGitlab,
        collapsable: true,
        listItems: [
          {
            text: "gitlab.com",
          },
          {
            text: "r-n-d.informatik.hs-augsburg.de",
          }
        ]
      },

    ]
  }

  ngOnInit(): void {
  }

}
