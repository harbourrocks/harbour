import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { faGithub, faGitlab } from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
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
