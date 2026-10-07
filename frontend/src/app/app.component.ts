import { Component, OnInit } from '@angular/core'
import { Router, RouterModule, RouterOutlet } from '@angular/router'
import { CommonModule } from '@angular/common'
import { NglModule } from 'ng-lightning'
import { UserStoreService } from '../state/user.service'
import { CanvasJSAngularChartsModule } from '@canvasjs/angular-charts'

@Component({
  selector: 'app-root',
  imports: [CommonModule, RouterOutlet, RouterModule, NglModule, CanvasJSAngularChartsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit{
  constructor(
    private userStoreService: UserStoreService,
    private router: Router
  ) {}

  title = 'felix-store-app';

  ngOnInit(): void {
    this.checkAuth()
    // this.promotionStoreService.loadPromotionData()
  }

  async checkAuth() {
    const checkpoint = await this.userStoreService.isAuthenticated()
    if (!checkpoint) {
      this.router.navigate(['/login'])
    } else {
      await this.userStoreService.loadUserInfo()
    }
    
  }
}
