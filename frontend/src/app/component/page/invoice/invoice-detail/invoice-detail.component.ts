import { Component, OnInit } from '@angular/core'
import { ActivatedRoute, Router } from '@angular/router'
import { getApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzTableModule } from 'ng-zorro-antd/table'
import { MatIconModule } from '@angular/material/icon'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import moment from 'moment'


@Component({
    imports: [
        CommonModule,
        FormsModule,
        NzFormModule,
        NzSelectModule,
        NzInputModule,
        NzInputNumberModule,
        NzTableModule,
        MatIconModule,
        NzDatePickerModule
    ],
    templateUrl: './invoice-detail.component.html'
})
export class InvoiceDetailComponent implements OnInit {
    constructor(
        private route: ActivatedRoute,
        private routeTo: Router
    ) {}
    ngOnInit(): void {
        this.loadLocationList()
        this.route.queryParams.subscribe((x: any) => {
            if (x.id) {
                this.getDetail(x.id)
            }
        })
    }

    async getDetail(id: string) {
        const data = await getApiWithAuth(`/invoice/one/${id}`)
        this.invoiceDetail = data
        this.invoiceDetail.discount = this.invoiceDetail.discountType === '%' ? this.invoiceDetail.discount * 100 : this.invoiceDetail.discount
    }

    locationList: any[] = []
    async loadLocationList() {
        const data = await getApiWithAuth('/base/location/getAll')
        this.locationList = data
    }

    invoiceDetail: any = {}
    
    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    goBackList() {
        this.routeTo.navigate(['invoice-list'])
    }
}