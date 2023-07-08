import { Injectable } from "@angular/core";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { ConfirmDialogComponent } from "src/app/shared/confirm-dialog/confirm-dialog.component";

@Injectable({
  providedIn: 'root'
})
export class ModalService { 

  constructor(private dialog: MatDialog){}

  openConfirmDialog(title: string, content: string, action: () => void): MatDialogRef<any, any> {
    return this.dialog.open(ConfirmDialogComponent, {
      width: '400px',
      data: {
        title: title,
        text: content,
        action: action
      }
    });
  }
}