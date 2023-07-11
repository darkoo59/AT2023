import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable, tap } from "rxjs";
import { GenericDataService } from "src/app/shared/generic-data.service";
import { environment } from "src/environments/environment";

@Injectable()
export class BalanceService extends GenericDataService<number> {

  constructor(private http: HttpClient) { super() }

  fetchBalance(): Observable<any> {
    return this.http.get(`${environment.apiUrl}/customer/balance`).pipe(
      tap(res => this.setData = res as number)
    )
  }

  patchBalance(data:{balance: number}): Observable<any> {
    return this.http.patch(`${environment.apiUrl}/customer/balance`, data);
  }
}