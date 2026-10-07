import { Component, OnInit } from '@angular/core'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { MemberForm, MemberSpecialDay } from './interface'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import moment from 'moment'
import { NzMessageService } from 'ng-zorro-antd/message'
import { ActivatedRoute, Router } from '@angular/router'
import { MatIconModule } from '@angular/material/icon'

@Component({
    templateUrl: './member-form.component.html',
    imports: [
        CommonModule,
        FormsModule,
        NzFormModule,
        NzInputModule,
        NzSelectModule,
        NzInputNumberModule,
        NzTableModule, 
        NzButtonModule,
        NzDatePickerModule,
        MatIconModule
    ]
})
export class MemberFormComponent implements OnInit {
    constructor(
        private routeTo: Router,
        private route: ActivatedRoute,
        private message: NzMessageService
    ) {}
    ngOnInit(): void {
        this.route.queryParams.subscribe((x: any) => {
            if (x.id) {
                this.theId = x.id
                this.getOne()
            }
        })
        this.loadMemberClassList()
    }

    editForm: MemberForm = {
        _id: '',
        memberCode: '',
        name: '',
        address: '',
        phone: '',
        email: '',
        fax: '',
        classId: '',
        remark: '',
        memberSpecialDays: []
    }
    
    specialDayForm: MemberSpecialDay = {
        name: '',
        date: '',
        remark: ''
    }

    theId: string = ''

    memberClassList: any[] = []
    async loadMemberClassList() {
        this.memberClassList = await getApiWithAuth('/member/member-class/getAll')
    }


    resetForm() {
        this.editForm = {
            _id: '',
            memberCode: '',
            name: '',
            address: '',
            phone: '',
            email: '',
            fax: '',
            classId: '',
            remark: '',
            memberSpecialDays: []
        }
        this.specialDayForm = {
            name: '',
            date: '',
            remark: ''
        }
    }

    async getOne() {
        this.editForm = await getApiWithAuth(`/member/member-list/one/${this.theId}`)
    }

    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY') : null
    }

    putItemSpecialDay() {
        this.editForm.memberSpecialDays.push(this.specialDayForm)
        this.specialDayForm = {
            name: '',
            date: '',
            remark: ''
        }
    }

    async removeRowSpecialDay(index: number, id?: string) {
        if (id) {
            const res = await getApiWithAuth(`/member/member-list/special-day/remove/${id}`)
            if (res) {
                this.message.info(res.msg)
            }
            this.editForm.memberSpecialDays.splice(index, 1)
        } else {
            this.editForm.memberSpecialDays.splice(index, 1)
        }
    }

    async submitForm() {
        const res = await postApiWithAuth(
            this.editForm._id ? '/member/member-list/update' : '/member/member-list/create', 
            this.editForm
        )

        if (res) {
            this.message.info(res.msg)
        }
    }

    backToList() {
        this.routeTo.navigate([`/member-list`])
    }
}