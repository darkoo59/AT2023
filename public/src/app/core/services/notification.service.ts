import { Injectable } from "@angular/core";
import { ToastrService } from "ngx-toastr";
import { catchError, finalize, of, switchMap, tap } from "rxjs";
import { WebSocketSubject } from 'rxjs/webSocket';
import { environment } from "src/environments/environment";
import { AuthService } from "./auth.service";
import { getDecodedAccessToken } from "src/utils/utility";

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  socket$: WebSocketSubject<any> = this.authService.data$.pipe(switchMap((token: string | null) => {
    let payload = getDecodedAccessToken(token ?? "")
    if(payload == null) return of({})
    return new WebSocketSubject(`${environment.notificationUrl}/${payload.sub}`).pipe(
      tap((response: any) => {
        console.log('Received message from server:', response)
        this.showNotification(response.Message)
      }),
      catchError((error) => {
        console.error('An error occurred:', error)
        return of({})
      }),
      finalize(() => console.log('WebSocket connection closed'))
    );
  })) as WebSocketSubject<any>

  constructor(private toastr: ToastrService, private authService: AuthService) {}

  showError(message: string): void {
    this.toastr.error(message, 'Error', {
      timeOut: 5000,
      positionClass: 'toast-bottom-right',
      progressBar: true
    });
  }

  showSuccess(message: string): void {
    this.toastr.success(message, 'Success', {
      timeOut: 5000,
      positionClass: 'toast-bottom-right',
      progressBar: true
    });
  }

  showNotification(message: string): void {
    this.toastr.info(message, 'Notification', {
      timeOut: 5000,
      positionClass: 'toast-bottom-right',
      progressBar: true
    })
  }
}