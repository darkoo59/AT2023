import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { NgLetModule } from "ng-let";
import { OrderComponent } from "./order.component";
import { OrderRoutingModule } from "./order-routing.module";
import { MaterialModule } from "src/app/shared/material.module";
import { ReactiveFormsModule } from "@angular/forms";

@NgModule({
  declarations: [OrderComponent],
  imports: [
    CommonModule, 
    NgLetModule,
    OrderRoutingModule,
    MaterialModule,
    ReactiveFormsModule
  ],
  exports: [] 
})
export class OrderModule {}