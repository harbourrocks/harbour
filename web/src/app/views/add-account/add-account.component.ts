import { Component, OnInit, Input } from '@angular/core';
import { FormModel } from 'src/app/models/form.model';
import {Location} from '@angular/common';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-add-account',
  templateUrl: './add-account.component.html',
  styleUrls: ['./add-account.component.scss']
})
export class AddAccountComponent implements OnInit {

  public formModel: FormModel;

  constructor(private _location: Location, private route: ActivatedRoute,) { 
    
  }
  
  ngOnInit(): void {
    const accountProvider = this.route.snapshot.paramMap.get('acc_prov');
    this.formModel = {
      header:`Add ${accountProvider} Account`,
      items: [
        {name: "appId", placeholder: "AppId"},
        {name: "installationId", placeholder: "InstallationId"},
        {name: "clientId", placeholder: "ClientId"},
        {name: "clientSecret", placeholder: "ClientSecret"},
        {name: "privateKey", placeholder: "PrivateKey"},
      ],
    }
    
  }

  onCancel() {
    this._location.back();
  }

  onSubmit() {
    //CALL

    this._location.back();
  }
}
