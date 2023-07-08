import { NgModule } from "@angular/core";
import { SafeResourcePipe } from "./safe-resource.pipe";

@NgModule({
  declarations: [
    SafeResourcePipe
  ],
  imports: [],
  exports: [
    SafeResourcePipe
  ]
})
export class CustomPipesModule { }