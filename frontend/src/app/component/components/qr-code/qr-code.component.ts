import { CommonModule } from '@angular/common'
import { Component, ComponentFactoryResolver, EventEmitter, HostListener, Input, OnInit, Output, ViewChild, ViewContainerRef } from '@angular/core'
import { QRCodeComponent } from 'angularx-qrcode'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule } from 'ng-zorro-antd/modal'

@Component({
    selector: 'app-qrcode',
    standalone: true,
    imports: [NzButtonModule, QRCodeComponent, NzModalModule, CommonModule],
    templateUrl: './qr-code.component.html',
    styleUrls: ['./qr-code.component.scss']
})
export class QRcodeComponent implements OnInit{
    constructor(
        private resolver: ComponentFactoryResolver
    ) {}
    ngOnInit(): void {
        // this.createHtml()
    }

    @Input() qrCodeString!: string
    @Input() isOpen!: boolean
    @Input() randomHtml: string = '<div>123</div>'
    @Input() widthInBox: string = '400px'
    @Output() isOpenChange = new EventEmitter<boolean>()
    @ViewChild('container', { read: ViewContainerRef }) container!: ViewContainerRef

    closeDialog() {
        this.isOpen = false
        this.isOpenChange.emit(false)
    }
}
