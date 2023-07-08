import { Injectable } from "@angular/core";
import { GenericDataService } from "../../shared/generic-data.service";

@Injectable({
  providedIn: 'root',
})
export class LoadingService extends GenericDataService<boolean> {
  
}