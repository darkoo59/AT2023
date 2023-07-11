import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { NgLetModule } from "ng-let";
import { BalanceRoutingModule } from "./balance-routing.module";
import { BalanceComponent } from "./balance.component";
import { ReactiveFormsModule } from "@angular/forms";
import { MaterialModule } from "src/app/shared/material.module";
import { BalanceService } from "src/app/core/services/balance.service";

@NgModule({
  declarations: [BalanceComponent],
  providers: [BalanceService],
  imports: [
    CommonModule, 
    NgLetModule,
    BalanceRoutingModule,
    ReactiveFormsModule,
    MaterialModule
  ],
  exports: [] 
})
export class BalanceModule {}