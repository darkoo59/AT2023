import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { RouterModule } from "@angular/router";
import { NgLetModule } from "ng-let";
import { MaterialModule } from "src/app/shared/material.module";
import { NavigationComponent } from "./navigation.component";

@NgModule({
  declarations: [NavigationComponent],
  imports: [
    CommonModule,
    MaterialModule,
    NgLetModule,
    RouterModule
  ],
  exports: [NavigationComponent]
})
export class NavigationModule { }