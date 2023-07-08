import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { NgLetModule } from "ng-let";
import { BalanceRoutingModule } from "./balance-routing.module";
import { BalanceComponent } from "./balance.component";

@NgModule({
  declarations: [BalanceComponent],
  imports: [
    CommonModule, 
    NgLetModule,
    BalanceRoutingModule
  ],
  exports: [] 
})
export class BalanceModule {}