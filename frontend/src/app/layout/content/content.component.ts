import { Component, HostListener } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
    selector: 'app-content',
    imports: [RouterOutlet],
    templateUrl: './content.component.html',
    styleUrls: ['./content.component.scss']
})
export class ContentComponent {
    isMenuExpanded: boolean = window.innerWidth > 200; // Expand only on large screens
    screenWidth: number = window.innerWidth;

    @HostListener('window:resize', ['$event'])
    onResize(event: any) {
        this.screenWidth = window.innerWidth;
        this.isMenuExpanded = this.screenWidth > 200
    }
}
