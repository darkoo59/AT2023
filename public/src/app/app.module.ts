import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgLetModule } from 'ng-let';
import { MaterialModule } from './shared/material.module';
import { NavigationModule } from './core/navigation/navigation.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { ToastrModule } from 'ngx-toastr';
import { AuthInterceptor } from './core/interceptors/auth-interceptor.service';
import { ConfirmDialogModule } from './shared/confirm-dialog/confirm-dialog.module';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    NgLetModule,
    MaterialModule,
    NavigationModule,
    HttpClientModule,
    ToastrModule.forRoot(),
    ConfirmDialogModule
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
