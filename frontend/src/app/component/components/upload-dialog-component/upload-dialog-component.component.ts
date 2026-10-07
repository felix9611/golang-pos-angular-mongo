import { CommonModule } from '@angular/common'
import { Component, Input, OnInit } from '@angular/core'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { formatJson, readExcelFile } from '../../../../tool/excel-helper'
import { postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzIconModule } from 'ng-zorro-antd/icon'

@Component({
    selector: 'app-upload-dialog',
    standalone: true,
    imports: [CommonModule, NzUploadModule, NzModalModule, NzIconModule],
    templateUrl: './upload-dialog-component.component.html'
})
export class UploadDialogComponent implements OnInit {
    @Input() uploadButtonLabel: string = ''
    @Input() uploadApiUrl: string = ''
    @Input() excelFieldList: any[] = []
    @Input() dbFieldList: any[] = []
    constructor(private message: NzMessageService) {}

    ngOnInit(){
    }

    upLoadDialog: boolean = false

    openUploadDialog() {
        this.upLoadDialog = true
    }
    
    closeUploadDialog() {
        this.upLoadDialog = false
    }

    async uploadAction(file: any) {
        if (file.type === 'start') {
            const data = await readExcelFile(file.file)
            const reData = formatJson(this.excelFieldList, this.dbFieldList, data)
                    
            if (reData.length > 0 ) {
                await postApiWithAuth(this.uploadApiUrl, reData)
                this.message.success('In Uploading')
                this.closeUploadDialog()
                        
            } else {
                    this.message.error('Ooooops, may data is wrong, please check again.')
            }
        }
    }
}