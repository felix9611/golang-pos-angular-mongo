import { CommonModule } from '@angular/common'
import { Component, ComponentFactoryResolver, EventEmitter, HostListener, Input, OnInit, Output, ViewChild, ViewContainerRef } from '@angular/core'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { getApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzMessageService } from 'ng-zorro-antd/message'
import { MatIconModule } from '@angular/material/icon'

@Component({
    selector: 'app-fileview-dialog',
    standalone: true,
    imports: [NzButtonModule, NzModalModule, CommonModule, MatIconModule],
    templateUrl: './file-view-dialog.component.html',
    styleUrls: ['./file-view-dialog.component.scss']
})
export class FileViewComponent implements OnInit{
    constructor(
        private resolver: ComponentFactoryResolver,
        private message: NzMessageService
    ) {}
    ngOnInit(): void {
        this.loadingFile()
    }

    @Input() fileList!: any
    @Input() fileNameKey!: string
    @Input() filePathKey!: string
    @Input() removeFileApi!: string
    @Input() loadFileApi!: string

    @Output() fileDialogOpenChange = new EventEmitter<boolean>()
    @ViewChild('container', { read: ViewContainerRef }) container!: ViewContainerRef

    @Input() fileDialogOpen!: boolean

    closeFileDialog() {
        this.fileDialogOpen = false
        this.fileDialogOpenChange.emit(false)
    }

    async removeFile(id: string) {
        const finalApi = this.removeFileApi.replace(':id', id)
        const res = await getApiWithAuth(finalApi)

        if (res.finished) {
            this.message.success('Remove file success')
            await this.loadingFile()
        } else {
            this.message.error(res.msg)
        }
    }

    async loadingFile() {
        this.fileList = await getApiWithAuth(this.loadFileApi)
    }
}
