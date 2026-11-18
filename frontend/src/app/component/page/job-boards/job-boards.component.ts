import {Component, computed, signal} from '@angular/core'
import {bootstrapApplication} from '@angular/platform-browser'
import {CommonModule} from '@angular/common'
import {FormsModule} from '@angular/forms'
import {CdkDragDrop, DragDropModule, moveItemInArray} from '@angular/cdk/drag-drop'
import {Job, newJob, demoJobs} from './job'
import {FilterListApi, ListOption} from '../../components/filter-list/filter-list.component'
import { deleteApiWithAuth, getApiWithAuth, postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { Subscription } from 'rxjs'
import { UserStoreService } from '../../../../state/user.service'


@Component({
    imports: [CommonModule, FormsModule, DragDropModule, FilterListApi],
    standalone: true,
    templateUrl: './job-boards.component.html',
    styleUrls: ['./job-boards.component.css']
})
export class JobBoardsComponent  {
  private userSubscription: Subscription
  constructor(
    private userStoreService: UserStoreService
  ) {
    this.userSubscription = this.userStoreService.user$.subscribe((data: any) => {
      this.yourUserData._id = data._id
      this.yourUserData.avatar = data.avatarBase64
      this.yourUserData.username = data.username
      this.yourUserData.superAdmin = data.superAdmin
    })
  }

  ngOnDestroy() {
    if (this.userStoreService.user$) {
      this.userSubscription.unsubscribe()
    }
    this.destroyed = true
    clearInterval(this.reloadTimer)
  }

  yourUserData: any = {
    _id: '',
    avatar: '',
    username: '',
    superAdmin: false
  }
 
  jobs: Job[] = [];
  search = signal('');
  filter = signal('all');
  view = signal('board');
  columns = [
    { id: 1, name: 'To Do List', en: 'To do', color: '#8b91a1', key: 'toDoList' },
    { id: 2, name: 'In Progress', en: 'In progress', color: '#6c5ce7', key: 'inProgressList' },
    { id: 3, name: 'Pending Review', en: 'Pending review', color: '#e5a345', key: 'inReviewList' },
    { id: 4, name: 'In Review', en: 'In review', color: '#ff5e00', key: 'inReviewList' },
    { id: 5, name: 'Completed', en: 'Done', color: '#39a982', key: 'completedList' },
  ]
  assigneeOptions(): ListOption[] {
    const options = [...this.memberOptions];
    for (const id of this.selected?.userIds || [])
      if (!options.some((o) => o.value === id)) options.push({ value: id, label: id });
    return options;
  }
  selected: Job | null = null
  tab = 'details'
  tagText = ''
  labelText = ''
  comment = ''
  checklist = ''
  error = ''
  toast = signal('')
  failed = signal('')

  private reloadTimer?: ReturnType<typeof setInterval>
  private loadingJobs = false
  private destroyed = false

  ngOnInit() {
    void this.loadJobLists()
    this.loadSomeUser()

    this.reloadTimer = setInterval(() => {
      // When user in edit task, this will stopped reloading
      if (!this.selected) {
        void this.loadJobLists()
      }
    }, 60_000);
  }

  dataSearchForm: any = {
  }

  async loadJobLists(): Promise<void> {
    if (this.loadingJobs || this.destroyed) return
    this.loadingJobs = true

     try {
      const data = await postApiWithAuth('/tasks/list', this.dataSearchForm)

      if (this.destroyed) return

      this.jobs = [
        ...(data.toDoList ?? []),
        ...(data.inProgressList ?? []),
        ...(data.pendingReviewList ?? []),
        ...(data.inReviewList ?? []),
        ...(data.completedList ?? []),
      ]
    } catch (error) {
      console.error(error)
    } finally {
      this.loadingJobs = false
    }
  }

  userLists: any[] = []
  async loadSomeUser() {
    const allUser = await postApiWithAuth('/sys/user/list', {page: 1, limit: 10})
    this.userLists = allUser.lists
  }
  postApiWithAuth(arg0: string, arg1: { page: number; limit: number }) {
    throw new Error('Method not implemented.')
  }

  // members data
  memberTotal: number = 0
  members: any = []
  memberOptions: ListOption[] = [];
  systemUserSearchForm: any = {
    username: '',
    page: 1,
    limit: 10
  }
  loadUserLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.systemUserSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.members = res.data.lists ?? [];
      this.memberTotal = res.data.total ?? 0;

      this.memberOptions = this.members.map((m: any) => ({
        value: m.id,
        label: m.username,
      }));

      return this.memberOptions;
    } catch (error) {
      console.error('API request failed', error);
      throw error;
    }
  }

  // Product data
  productTotal: number = 0
  products: any = []
  productOptions: ListOption[] = [];
  productSearchForm: any = {
    query: '',
    page: 1,
    limit: 10
  }
  loadProductLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.productSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.products = res.lists ?? []
      this.productTotal = res.total ?? 0

      this.productOptions = this.products.map((a: any) => ({
        value: a.productCode,
        label: `${a.productCode} - ${a.productName} (${a.brandName})`,
      }))

      return this.productOptions
    } catch (error) {
      console.error('API request failed', error)
      throw error
    }
  }

  // ASSET Type data
  assetTypeTotal: number = 0
  assetTypes: any = []
  assetTypeOptions: ListOption[] = [];
  assetTypeSearchForm: any = {
    typeCode: '',
    typeName: '',
    page: 1,
    limit: 10
  }
  loadAssetTypeLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.assetTypeSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.assetTypes = res.lists ?? [];
      this.assetTypeTotal = res.total ?? 0;

      this.assetTypeOptions = this.assetTypes.map((a: any) => ({
        value: a.typeCode,
        label: `${a.typeCode} - ${a.typeName}`,
      }));

      return this.assetTypeOptions;
    } catch (error) {
      console.error('API request failed', error);
      throw error;
    }
  }

  // Department data
  deptTotal: number = 0
  depts: any = []
  deptOptions: ListOption[] = [];
  deptSearchForm: any = {
    deptCode: '',
    deptName: '',
    page: 1,
    limit: 10
  }
  loadDeptLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.deptSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.depts = res.data.lists ?? [];
      this.deptTotal = res.data.total ?? 0;

      this.deptOptions = this.depts.map((a: any) => ({
        value: a.deptCode,
        label: `${a.deptCode} - ${a.deptName}`,
      }));

      return this.deptOptions;
    } catch (error) {
      console.error('API request failed', error);
      throw error;
    }
  }

  // Vendor data
  vendorTotal: number = 0
  vendors: any = []
  vendorOptions: ListOption[] = [];
  vendorSearchForm: any = {
    vendorCode: '',
    vendorName: '',
    page: 1,
    limit: 10
  }
  loadVendorLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.vendorSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.vendors = res.lists ?? [];
      this.vendorTotal = res.total ?? 0;

      this.vendorOptions = this.vendors.map((a: any) => ({
        value: a.vendorCode,
        label: `${a.vendorCode} - ${a.vendorName}`,
      }));

      return this.vendorOptions;
    } catch (error) {
      console.error('API request failed', error);
      throw error;
    }
  }

  // Stock take data
  stockTakeTotal: number = 0
  stockTakes: any = []
  stockTakeOptions: ListOption[] = [];
  stockTakeSearchForm: any = {
    stockTakeCode: '',
    stockTakeName: '',
    page: 1,
    limit: 10
  }
  loadStockTakeLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.stockTakeSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.stockTakes = res.lists ?? [];
      this.stockTakeTotal = res.total ?? 0;

      this.stockTakeOptions = this.stockTakes.map((a: any) => ({
        value: a.actionName,
        label: `${a.actionName}`,
      }));

      return this.stockTakeOptions;
    } catch (error) {
      console.error('API request failed', error);
      throw error;
    }
  }

  // Location data
  placeTotal: number = 0
  places: any = []
  placeOptions: ListOption[] = [];
  placeSearchForm: any = {
    query: '',
    page: 1,
    limit: 10
  }
  loadPlaceLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.placeSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.places = res.lists ?? [];
      this.placeTotal = res.total ?? 0;

      this.placeOptions = this.places.map((a: any) => ({
        value: a.placeCode,
        label: `${a.placeName} - ${a.address} - ${a.country} ${a.zipCode}`,
      }));

      return this.placeOptions;
    } catch (error) {
      console.error('API request failed', error);
      throw error;
    }
  }

  // Invoice data
  invoiceTotal: number = 0
  invoices: any = []
  invoiceOptions: ListOption[] = [];
  invoiceSearchForm: any = {
    number: '',
    page: 1,
    limit: 10
  }
  loadInvoiceLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.invoiceSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.invoices = res.lists ?? []
      this.invoiceTotal = res.total ?? 0

      this.invoiceOptions = this.invoices.map((a: any) => ({
        value: a.number,
        label: `${a.number} (${a.location.placeName})`,
      }));

      return this.invoiceOptions
    } catch (error) {
      console.error('API request failed', error)
      throw error
    }
  }

  // Invoice data
  returnInvoiceTotal: number = 0
  returnInvoices: any = []
  returnInvoiceOptions: ListOption[] = [];
  returnInvoiceSearchForm: any = {
    caseNumber: '',
    invoiceNumber: '',
    page: 1,
    limit: 10
  }
  loadReturnInvoiceLists = async (
    url: string,
    body: Record<string, unknown>,
  ): Promise<ListOption[]> => {
    const payload = {
      ...this.returnInvoiceSearchForm,
      ...body,
    }

    try {
      const res = await postApiWithAuth(url, payload);

      this.returnInvoices = res.lists ?? []
      this.returnInvoiceTotal = res.total ?? 0

      this.returnInvoiceOptions = this.returnInvoices.map((a: any) => ({
        value: a.returnCaseNumber,
        label: `${a.returnCaseNumber} (${a.invoiceNumber}$ - ${a.location.placeName})`,
      }));

      return this.returnInvoiceOptions
    } catch (error) {
      console.error('API request failed', error)
      throw error
    }
  }

  visible(): Job[] {
    return this.jobs.filter(
      (j: any) =>
        (this.filter() === 'all') &&
        (
          j.title +
          ' ' +
          j.number +
          ' ' +
          j.category +
          ' ' +
          j.tags.join(' ') +
          ' ' +
          j.labels.join(' ')
        )
          .toLowerCase()
          .includes(this.search().toLowerCase()),
    );
  }
  done(): number {
    return this.jobs.filter((j) => j.status === 5).length
  }
  active(): number {
    return this.jobs.filter((j) => j.status === 2).length;
  }

  cards(status: number) {
    return this.visible().filter((j) => j.status === status);
  }

  countDone(job: Job) {
    return job.checkLists &&job.checkLists.filter((c) => c.status === 1 && c.checked).length;
  }
  countAllDone(job: Job) {
    return job.checkLists && job.checkLists.filter((c) => c.status === 1).length;
  }

  initials(id: string) {
    return this.members.find((m: any) => m._id === id)?.initial || id.slice(0, 1);
  }

  async open(job: Job) {
    if (job && job._id) {
      this.selected = await getApiWithAuth(`/tasks/get-one/${job._id}`);
     // this.tagText =  job.tags.join(', ');
    //  this.labelText = job.labels.join(', ');
    } else {
       this.selected = structuredClone(job);
    }
    this.tab = 'details';
    this.error = '';
    this.comment = '';
    this.checklist = '';
  }

  create(status = 1) {
    const job = newJob(status)
    this.open(job)
  }

  changeStatus(job: Job | null | undefined, status: number) {
    if (!job?._id) return

    this.jobs = this.jobs.map(item =>
      item._id === job._id ? { ...item, status } : item,
    )
  }

  async save() {
    if (!this.selected) return;
    const j = this.selected;
    if (!j.title.trim()) {
      this.error = 'Please enter a title.';
      return;
    }
    if (j.jobStartDate && j.jobEndDate && j.jobEndDate < j.jobStartDate) {
      this.error = 'End date must be after the start date.';
      return;
    }

    j.tags = this.tagText
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
    j.labels = this.labelText
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);

    const finalData = j

    if (finalData._id !== '') {
      const res = await postApiWithAuth('/tasks/update', finalData)
      if (res.error) {
        this.notify('Update task failed')
      } else {
        this.selected = null;
        this.notify('Task updated successfully')
        this.loadJobLists()
      }
    } else {
      const res = await postApiWithAuth('/tasks/create', finalData)
      if (res.error) {
        this.notify('Create task failed')
      } else {
        this.selected = null;
        this.notify('Task saved successfully')
        this.loadJobLists()
      }
    }
  }

  async remove() {
    if (this.selected && confirm('Confirm delete this task？')) {
      const url = `/tasks/void/${this.selected._id}`
      const res = await deleteApiWithAuth(url)
      if (res.error) {
        this.notify(res.error)
      } else {
        this.selected = null;
        this.notify(res.message)
        this.loadJobLists()
      }
    }
  }

 savingMove = false;

  async drop(event: CdkDragDrop<Job[]>, status: number) {

      if (this.savingMove) return

      const moved = event.item.data as Job

      const target = event.container.data.filter(job => job._id !== moved._id);
      target.splice(event.currentIndex, 0, moved);

      const previous = target[event.currentIndex - 1];
      const next = target[event.currentIndex + 1];

      if (!(this.yourUserData.superAdmin || this.yourUserData.manager) && (status > 3 || target[0].status === 4 || target[0].status === 5)) {
        this.failNotify('You are not allowed to move this task')
        return
      } else {
        let sortOrder: number;

        if (previous && next) {
          sortOrder = previous.sortOrder + (next.sortOrder - previous.sortOrder) / 2
        } else if (previous) {
          sortOrder = previous.sortOrder + 10
        } else if (next) {
          sortOrder = next.sortOrder - 10
        } else {
          sortOrder = 10
        }

        const previousJobs = structuredClone(this.jobs);
        const list = [...event.container.data];

        if (event.previousContainer === event.container) {
          moveItemInArray(list, event.previousIndex, event.currentIndex);
        } else {
          list.splice(event.currentIndex, 0, { ...moved, status });
        }

        if (moved.status === status) return

        this.savingMove = true

      try {
          await postApiWithAuth('/tasks/change-status', {
            _id: moved._id,
            status,
            sortOrder
          })
          this.loadJobLists()
          this.notify('Updated successfully')
      } catch (error) {
          this.jobs = previousJobs
          this.failNotify('Update failed, the tasks have been fallback')
          console.error(error)
      } finally {
        this.savingMove = false
      }
    }
  }
  toggleMember(id: string) {
    if (!this.selected) return;
    const a = this.selected.userIds;
    this.selected.userIds = a.includes(id) ? a.filter((x) => x !== id) : [...a, id];
  }
  addCheck() {
    if (this.checklist.trim() && this.selected) {
      this.selected.checkLists.push({ name: this.checklist.trim(), checked: false, status: 1 });
      this.checklist = '';
    }
  }
  addComment() {
    if (this.comment.trim() && this.selected) {
      this.selected.comments.push({ 
        userId: this.yourUserData._id, 
        comment: this.comment.trim(), 
        status: 1, 
        username: this.yourUserData.username, 
        avatar: this.yourUserData.avatar 
      })
      this.comment = ''
    }
  }
  async attach(event: Event) {
    const input = event.target as HTMLInputElement;
    const target = this.selected;
    if (!target) return;
    for (const file of Array.from(input.files || [])) {
      if (file.size > 2 * 1024 * 1024) {
        this.error = 'Per file size limit is 2MB';
        continue;
      }
      try {
        const base64 = await new Promise<string>((resolve, reject) => {
          const r = new FileReader();
          r.onload = () => resolve(String(r.result));
          r.onerror = reject;
          r.readAsDataURL(file);
        });
        target.attachments.push({
          fileName: file.name,
          fileType: file.type,
          base64,
          updatedAt: new Date().toISOString(),
          status: 1
        });
      } catch {
        this.failNotify('Attachment upload failed, please try again.')
      }
    }
    input.value = '';
  }
  download(a: Job['attachments'][number]) {
    const bytes = atob(a.base64.split(',')[1]);
    const url = URL.createObjectURL(
      new Blob([Uint8Array.from(bytes, (c) => c.charCodeAt(0))], {
        type: 'application/octet-stream',
      }),
    );
    const link = document.createElement('a');
    link.href = url;
    link.download = a.fileName;
    link.click();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
  }
  export() {
    const url = URL.createObjectURL(
      new Blob([JSON.stringify(this.jobs, null, 2)], { type: 'application/json' }),
    );
    const a = document.createElement('a');
    a.href = url;
    a.download = 'flowboard-jobs.json';
    a.click();
    URL.revokeObjectURL(url);
    this.notify('Output JSON');
  }
  notify(s: string) {
    this.toast.set(s);
    setTimeout(() => this.toast.set(''), 3500);
  }

  failNotify(s: string) {
    this.failed.set(s);
    setTimeout(() => this.failed.set(''), 4000);
  }
}