import { Component, Renderer2 } from '@angular/core';
import { tap } from 'rxjs';
import { ThemeService } from './core/services/theme.service';
import { LoadingService } from './core/services/loading.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  loading$ = this.loadingService.data$;
  theme$ = this.themeService.data$.pipe(
    tap((theme: 'dark' | 'light' | null) => {
      if (theme === 'dark') this.renderer.addClass(document.body, 'dark-theme');
      else if (theme === 'light')
        this.renderer.removeClass(document.body, 'dark-theme');
    })
  );

  constructor(
    private renderer: Renderer2,
    private themeService: ThemeService,
    private loadingService: LoadingService
  ) {}
}
