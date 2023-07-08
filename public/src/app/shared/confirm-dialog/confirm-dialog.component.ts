import { Component, Inject } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { Subject, tap } from "rxjs";

@Component({
  templateUrl: './confirm-dialog.component.html'
})
export class ConfirmDialogComponent {
  constructor(
    private dialogRef: MatDialogRef<ConfirmDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { action: () => void, text: string, title: string }
  ) { }

  confirm$: Subject<any> = new Subject<any>().pipe(
    tap(() => {
      this.data.action()
      this.dialogRef.close(true)
    })
  ) as Subject<any>;
}