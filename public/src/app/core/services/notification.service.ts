import { Injectable } from "@angular/core";
import { ToastrService } from "ngx-toastr";

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  constructor(private toastr: ToastrService) { }

  showError(message: string): void {
    this.toastr.error(message, 'Error', {
      timeOut: 3000,
      positionClass: 'toast-bottom-right',
      progressBar: true
    });
  }

  showSuccess(message: string): void {
    this.toastr.success(message, 'Success', {
      timeOut: 3000,
      positionClass: 'toast-bottom-right',
      progressBar: true
    });
  }
}