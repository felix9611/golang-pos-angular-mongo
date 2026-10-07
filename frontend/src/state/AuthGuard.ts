// src/app/guards/auth.guard.ts
import { Injectable } from '@angular/core'
import {
  CanActivate,
  Router,
} from '@angular/router'
import { UserStoreService } from './user.service'

@Injectable({
  providedIn: 'root',
})
export class AuthGuard implements CanActivate {
  constructor(private userStoreService: UserStoreService, private router: Router) {}

  canActivate(): boolean {
    if (!this.userStoreService.isAuthenticated()) {
        this.router.navigate(['/login'])
        return false
    } else {
      return true
    }
    
  }
}

