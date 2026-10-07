import { CommonModule } from '@angular/common'
import { Component, Input } from '@angular/core'
import { downloadTempExcelFile } from '../../../../tool/excel-helper'

@Component({
    selector: 'app-download-excel-template',
    standalone: true,
    template: '<div (click)="downloadTemplateExcel()">{{ buttonLabel }}</div>'
})
export class DownloadExcelTemplateComponent {
    @Input() buttonLabel: string = 'Download Excel Template'
    @Input() excelFileName: string = 'template.xlsx'
    @Input() excelFieldList: any[] = []

    downloadTemplateExcel() {
         downloadTempExcelFile(this.excelFieldList, this.excelFileName)
    }
}