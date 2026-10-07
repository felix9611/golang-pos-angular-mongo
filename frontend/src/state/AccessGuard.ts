import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { UserStoreService } from './user.service';
import { findMenuItemByPath } from '../app/component/tool-function';
import { Observable } from 'rxjs';
import { map, take } from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class AccessGuard implements CanActivate {
    constructor(
        private userStoreService: UserStoreService, 
        private router: Router
    ) {}

    canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ) {
        if (!this.userStoreService.isAuthenticated()) {
            this.router.navigate(['/login'])
            return false
        }

        const urlPath = route.url[0]?.path || ''
        
        return this.userStoreService.menuRole$.pipe(
            take(1), // Only take one
            map((menuData: any) => {
                const hasAccess = !!findMenuItemByPath(menuData, urlPath);
                if (!hasAccess) {
                    this.router.navigate(['/'])
                }
                return hasAccess
            })
        )
    }
}