import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { NgLetModule } from "ng-let";
import { MaterialModule } from "../material.module";
import { ConfirmDialogComponent } from "./confirm-dialog.component";

@NgModule({
  declarations: [ConfirmDialogComponent],
  imports: [
    CommonModule,
    MaterialModule,
    NgLetModule
  ],
  exports: []
})
export class ConfirmDialogModule { }