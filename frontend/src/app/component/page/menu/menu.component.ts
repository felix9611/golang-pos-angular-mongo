import { Component, OnInit } from '@angular/core'
import { deleteApiWithAuth, getApiWithAuth, postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { CommonModule } from '@angular/common'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzFormModule } from 'ng-zorro-antd/form'
import { FormsModule } from '@angular/forms'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzInputModule } from 'ng-zorro-antd/input'
import { MenuForm } from './interface'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzRadioModule } from 'ng-zorro-antd/radio'
import { NzMessageService } from 'ng-zorro-antd/message'
import { findMenuItem } from '../../tool-function'
import { UserStoreService } from '../../../../state/user.service'
import { Subscription } from 'rxjs'

@Component({
    standalone: true,
    imports: [
        NzCheckboxModule, 
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzModalModule, 
        NzTableModule, 
        NzInputModule,
        NzSelectModule,
        NzInputNumberModule,
        NzRadioModule
    ],
    templateUrl: './menu.component.html',
    styleUrl: './menu.component.css',
})
export class MenuListComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Menu', 'menu')
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


    ngOnInit(): void {
        this.loadSysMenuLists()
        this.loadMainItemLists()
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false
    }

    listOfData = [
        {
          id: '1',
          name: 'John Brown',
          age: 32,
          expand: false,
          address: 'New York No. 1 Lake Park',
          description: 'My name is John Brown, I am 32 years old, living in New York No. 1 Lake Park.'
        },
        {
          id:'2',
          name: 'Jim Green',
          age: 42,
          expand: false,
          address: 'London No. 1 Lake Park',
          description: 'My name is Jim Green, I am 42 years old, living in London No. 1 Lake Park.'
        },
        {
          id: '3',
          name: 'Joe Black',
          age: 32,
          expand: false,
          address: 'Sidney No. 1 Lake Park',
          description: 'My name is Joe Black, I am 32 years old, living in Sidney No. 1 Lake Park.'
        }
      ];

    searchForm = {
        name: ''
    }

    editForm: MenuForm = {
        _id: '',
        mainId: '',
        name: '',
        icon: '',
        path: '',
        sort: 0,
        type: '',
        menuIds: [],
        excelFunctionCode: '',
        excelFunctionName: ''
    }

    typeOptions = [
        { label: 'Main', value: '0' },
        { label: 'Sub Item', value: '1' },
        { label: 'Non-Menu', value: '2' }
    ]
    
    mainItem: any[] = []
    async loadMainItemLists() {
        const res = await getApiWithAuth('/sys/menu/main-item')
        this.mainItem = res
        this.mainItem.push({ _id: '', name: 'None under Main Menu', mainId: ''})
    }

    expandSet = new Set<String>()
    onExpandChange(id: string, checked: boolean): void {
        if (checked) {
          this.expandSet.add(id)
        } else {
          this.expandSet.delete(id)
        }
      }
    dataLists: any[] = []
    async loadSysMenuLists() {
        const res = await postApiWithAuth('/sys/menu/list', this.searchForm)
        this.dataLists = res
    }

    okText: string = 'Create'
    editDialogVisible: boolean = false
    showEditDialog(){
        this.editDialogVisible = true
        this.editForm = {
            _id: '',
            mainId: '',
            name: '',
            icon: '',
            path: '',
            sort: 0,
            type: '',
            menuIds: [],
            excelFunctionCode: '',
            excelFunctionName: ''
        }
    }

    closeEditDialog(){
        this.editDialogVisible = false
    }

    async submitForm() {
        const url = this.editForm._id === '' ? '/sys/menu/create' : `/sys/menu/update`
        const res = await postApiWithAuth(url, this.editForm)

        if (res.msg) {
            this.message.error(res.msg)
        } else if (res.matchedCount === 1 || !res.msg) {
            this.editForm = {
                _id: '',
                mainId: '',
                name: '',
                icon: '',
                path: '',
                sort: 0,
                type: '',
                menuIds: [],
                excelFunctionCode: '',
                excelFunctionName: ''
            }

            this.message.success('Save successful!')
            this.closeEditDialog()
            this.loadSysMenuLists()
            this.loadMainItemLists()

        }

    }

    async getDataById(id: string) {
        this.okText = 'Update'
        this.editForm = await getApiWithAuth(`/sys/menu/one/${id}`)
        this.editDialogVisible = true
    }
} 