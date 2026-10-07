import { CommonModule } from '@angular/common'
import { Component, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { MatIconModule } from '@angular/material/icon'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzTableModule } from 'ng-zorro-antd/table'
import { UpdateReturnInvoiceDto } from './interface'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { ActivatedRoute, Router } from '@angular/router'
import { NzMessageService } from 'ng-zorro-antd/message'
import { UserStoreService } from '../../../../../state/user.service'
import { debounceTime, Subscription } from 'rxjs'
import { findMenuItem } from '../../../tool-function'

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
            NzDatePickerModule,
            NzCheckboxModule, 
        ],
        templateUrl: './return-invoice-detail.component.html'
})
export class ReturnInvoiceDetailComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
            private route: ActivatedRoute, 
            private routeTo: Router,
            private message: NzMessageService,
            private userStoreService: UserStoreService
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Create Return', 'create-return-invoice')
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false,
                upload: answer.upload ?? false
                     // keep default value
                }
         })
    }


    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false
    }

    ngOnInit(): void {
        this.loadLocationList()
        this.route.queryParams.subscribe((x: any) => {
            if (x.id) {
                this.getDetail(x.id)
            }
        })
    }

    returnInvoiceForm: any = {
    }

    async getDetail(id: string) {
        const data = await getApiWithAuth(`/return-invoice/one/${id}`)
        this.returnInvoiceForm = data
    }

    goBackList() {
        this.routeTo.navigate(['return-invoice-list'])
    }

    locationList: any[] = []
    async loadLocationList() {
        const data = await getApiWithAuth('/base/location/getAll')
        this.locationList = data
    }

}
