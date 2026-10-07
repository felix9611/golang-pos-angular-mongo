import { Component, OnInit, ViewChild, AfterViewInit } from '@angular/core'
import { CommonModule } from '@angular/common'
import { RouterLink, RouterLinkActive, RouterOutlet, Router } from '@angular/router'
import { FormsModule } from '@angular/forms'
import { deleteApiWithAuth, getApiWithAuth, postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
import { RoleForm } from './interface'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { buildTreeForUI } from './function'
import { NzFormatEmitEvent, NzTreeComponent, NzTreeModule, NzTreeNodeOptions } from 'ng-zorro-antd/tree';
import { UserStoreService } from '../../../../state/user.service'
import { read, write } from 'fs'
import { update } from 'plotly.js-dist-min'
import { findMenuItem } from '../../tool-function'
import { Subscription } from 'rxjs'

@Component({
    // selector: 'app-footer',
    standalone: true,
    imports: [
        NzTreeModule, 
        NzCheckboxModule, 
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzModalModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule
    ],
    templateUrl: './role.component.html',
    styleUrl: './role.component.css',
})
export class RoleComponent implements OnInit{
    @ViewChild('nzTreeComponent', { static: false }) nzTreeComponent!: NzTreeComponent
     private rightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Role', 'role')
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false
                 // keep default value
            }
        })   
    }
    
    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            //    this.rightSubscription.unsubscribe()
        }
    }
    

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false
    }

    searchForm: any = {
        page: 1,
        limit: 10
    }

    editForm: RoleForm = {
        _id: '',
        code: '',
        name: '',
        remark: '',
        menuIds: [],
        read: false,
        write: false,
        delete: false,
        update: false,
        upload: false
    }

    okText: string = 'Create'

    dataLists: any[] = []
    totals: number = 0
    editFormDialog: boolean = false
    removeDialog: boolean = false
    handleRemoveId: string = ''

    ngOnInit() {
        this.loadSysRoleLists()
        this.loadAllMenuItems()
       
    }

    async submitForm() {
        const url = this.editForm._id === '' ? '/sys/role/create' : `/sys/role/update`

        const res = await postApiWithAuth(url, {
            
            ...this.editForm,
            ...this.editForm._id ? { _id: this.editForm._id} : {},
        })

        if (res.msg) {
            this.message.error(res.msg)
        } else if (res.matchedCount === 1 || !res.msg) {
            this.editForm = {
                _id: '',
                code: '',
                name: '',
                remark: '',
                menuIds: [],
                read: false,
                write: false,
                delete: false,
                update: false,
                upload: false
            }

            this.message.success('Save successful!')
            this.closeDialog()
            this.loadSysRoleLists()
        }
    }

    async loadAllMenuItems() {
        const data = await getApiWithAuth('/sys/menu/all-menu')
        this.menuItems = buildTreeForUI(data)
        console.log(this.menuItems, 'test')
    }

    async loadSysRoleLists() {
        const res = await postApiWithAuth('/sys/role/list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    showDialog() {
        this.editFormDialog = true
    }

    closeDialog() {
        this.editFormDialog = false
    }

    handleRomeve(id: string) {
        this.handleRemoveId = id
        this.removeDialog = true
    }

    closeRemoveDialog() {
        this.removeDialog = false
    }

    async handleRemove() {
        const url = `/sys/role/remove/${this.handleRemoveId}`

        const res: any = await getApiWithAuth(url)

        if (res.msg) {
            this.message.info(res.msg)
        }
        this.loadSysRoleLists()
        this.closeRemoveDialog()
        
    }


    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    async getOneData(id:string) {
        const res = await getApiWithAuth(`/sys/role/one/${id}`)
        this.editForm = res
        this.okText = 'Update'
        this.showDialog()
    }

    selectedMenusIds: any = []
    defaultMenusIds: any = []
    menusIdsExpandedKeys: any = []

    menuItems: any[] = []
    menuDialog: boolean = false
    handleId: string = ''

    handleMenuItemsIds: any = []



    async openMenuDialog(id: string) {
        this.menuDialog = true
        this.handleId = id

        const url = `/sys/role/one/${id}`

        const res: any = await getApiWithAuth(url)

        this.selectedMenusIds = res.menuIds
     //   this.defaultMenusIds = res.menuIds
    }

    closeMenuDialog() {
        this.menuDialog = false
    }

    async handleMeunPermission() {
        const res = await postApiWithAuth('/sys/role/update-permission', {
            id: this.handleId,
            menuIds: this.selectedMenusIds
        })
        if (res.msg) {
            this.message.error(res.msg)
        } else if (res.matchedCount === 1 || !res.msg) {
            this.handleMenuItemsIds = []
            this.handleId = ''

            this.message.success('Save successful!')
            this.closeMenuDialog()
            this.loadSysRoleLists()
        }
    }


    menuItemClick(event: NzFormatEmitEvent) {
        console.log(event)
    }

    menuItemCheck(event: NzFormatEmitEvent) {
        console.log(event)
        if (event.node?.isChecked === false) {
            this.selectedMenusIds = this.selectedMenusIds.filter((item: string) => item !== event.node?.key)
        } else {
            this.selectedMenusIds.push(event.node?.key)
        }
    }

    menuItemSelect(keys: string[]): void {
        console.log(keys, this.nzTreeComponent.getSelectedNodeList());
    }

}