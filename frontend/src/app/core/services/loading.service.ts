import { Injectable } from "@angular/core";
import { GenericDataService } from "./generic-data.service";

@Injectable({
  providedIn: 'root',
})
export class LoadingService extends GenericDataService<boolean> {
  
}