import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { PageLoaderComponent } from "./page-loader.component";

@NgModule({
  declarations: [PageLoaderComponent],
  imports: [
    CommonModule
  ],
  exports: [PageLoaderComponent]
})
export class PageLoaderModule {}