import { Component, OnInit } from '@angular/core';
import { faTachometerAlt, faBoxOpen, faHatWizard, faCog } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  tachometerAlt = faTachometerAlt;
  public navbarBuild = [
    {
      routerLink: "/dashboard",
      icon: faTachometerAlt,
      text: "Dashboard",
    },
    {
      routerLink: "/repository",
      icon: faBoxOpen,
      text: "Repository",
    },
    {
      routerLink: "/wizard",
      icon: faHatWizard,
      text: "Wizard",
    },
    {
      routerLink: "/settings",
      icon: faCog,
      text: "Settings",
    },

  ]

  constructor() { }

  ngOnInit(): void {
  }

}
