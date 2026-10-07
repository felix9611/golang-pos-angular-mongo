import { Component, ViewEncapsulation } from '@angular/core'
import { ContentComponent } from './content/content.component'
import { MenuComponent } from './menu/menu.component'
import { RouterLinkActive, RouterModule, RouterOutlet } from '@angular/router'
import { HeaderComponent } from './header/header.component'
import { UserStoreService } from '../../state/user.service'

@Component({
    imports: [ContentComponent, MenuComponent, HeaderComponent, RouterModule],
    selector: 'app-main',
    templateUrl: './main.component.html',
    styleUrl: './main.component.css',
    encapsulation: ViewEncapsulation.None
})
export class MainComponent {
    constructor(private userStoreService: UserStoreService) {
        this.userStoreService.loadUserInfo()
        this.userStoreService.loadMenus()
        this.userStoreService.loadMenuRoles()
    }
    isCollapsed = false;

    toggleSidebar() {
        this.isCollapsed = !this.isCollapsed;
    }

}