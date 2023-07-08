import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { MainComponent } from './main.component';
import { NgLetModule } from 'ng-let';
import { NavigationModule } from '../core/navigation/navigation.module';
import { MaterialModule } from '../shared/material.module';
import { MainRoutingModule } from './main-routing.module';
import { CommonModule } from '@angular/common';

@NgModule({
  declarations: [
    MainComponent
  ],
  imports: [
    CommonModule,
    MainRoutingModule,
    NgLetModule,
    MaterialModule,
    NavigationModule,
    HttpClientModule,
  ],
  bootstrap: [MainComponent]
})
export class MainModule { }
