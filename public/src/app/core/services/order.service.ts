import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

export interface CreateOrderDto {
  quantity: number;
  itemId: string;
}

@Injectable({
  providedIn: 'root',
})
export class OrderService {
  constructor(private http: HttpClient) {}

  createOrder(data: CreateOrderDto): Observable<any> {
    return this.http.post(`${environment.apiUrl}/customer/order`, {});
  }
}
