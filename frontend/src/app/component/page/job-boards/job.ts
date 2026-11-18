export interface Job {
  _id: string
  number: string
  category: string
  title: string
  description: string
  jobStartDate: string
  jobEndDate: string
  jobTimeRemain: string
  tags: string[]
  labels: string[]
  locationIds: string[]
  attachments: {
    fileName: string
    fileType: string
    base64: string
    updatedAt: string
    status: number
  }[]
  comments: {
    userId: string
    comment: string
    status: number
    username: string
    avatar: string
  }[]
  checkLists: {
    name: string;
    checked: boolean
    status: number
  }[]
  userIds: string[]
  jobRelationData: {
    productCode: string
    typeCode: string
    deptCode: string
    placeCode: string
    vendorCode: string
    stockTakeName: string
    invoiceNumber: string
    returnInvoiceCaseNumber: string
  }
  jobValid: boolean
  status: number
  users: {
    username: string
    avatarBase64: string
  }[]
  sortOrder: number
  commentsCount?: number
  attachmentsCount?: number
  checklistsFinishedCount?: number
  checklistsCount?: number
}
export const newJob = (status = 1): Job => ({
  _id: '',
  number: '',
  category: '一般工作',
  title: '',
  description: '',
  jobStartDate: '',
  jobEndDate: '',
  jobTimeRemain: '',
  tags: [],
  labels: [],
  attachments: [],
  comments: [],
  checkLists: [],
  userIds: [],
  jobRelationData: {
    productCode: '',
    typeCode: '',
    deptCode: '',
    placeCode: '',
    vendorCode: '',
    stockTakeName: '',
    invoiceNumber: '',
    returnInvoiceCaseNumber: '',
  },
  jobValid: false,
  status,
  users: [],
  sortOrder: 0,
  locationIds: [],
});
export function demoJobs(): Job[] {
  const rows: [string, string, number, string, string, string][] = [
    [
      '建立資產管理系統介面',
      '設計',
      1,
      '設計新的資產總覽，讓團隊更快掌握設備狀態。',
      '高優先',
      '11-12',
    ],
    ['整理第四季設備盤點清單', '營運', 1, '核對各部門設備資料及資產編號。', '盤點', '11-15'],
    ['更新供應商聯絡資料', '行政', 1, '確認主要供應商的聯絡資料及服務範圍。', '一般', '11-18'],
    ['辦公室網絡設備升級', '工程', 2, '升級交換器並確認各樓層的網絡連線。', '高優先', '11-06'],
    ['開發工作單 API', '開發', 2, '完成工作單的新增、更新及狀態轉換功能。', '開發', '11-10'],
    ['會議室設備保養', '維修', 2, '檢查投影機、音響及視訊會議設備。', '保養', '11-08'],
    ['驗收新一批手提電腦', '資產', 3, '確認設備規格、數量及保養文件。', '驗收', '11-05'],
    ['測試庫存報表匯出功能', '開發', 3, '檢查 CSV 匯出的欄位及資料格式。', '開發', '11-07'],
    ['完成十月份資產盤點', '營運', 4, '已完成所有部門的資產核對。', '盤點', '10-30'],
    ['更新團隊操作指引', '行政', 4, '整理最新的工作流程及操作步驟。', '文件', '10-29'],
  ];
  return rows.map((r, i) => ({
    ...newJob(r[2]),
    number: `JOB-${String(i + 1).padStart(3, '0')}`,
    title: r[0],
    category: r[1],
    description: r[3],
    labels: [r[4]],
    tags: i % 2 ? ['內部'] : ['Q4'],
    jobStartDate: '2026-11-01',
    jobEndDate: `2026-${r[5]}`,
    userIds: [['01', '02', '03'][i % 3]],
    checkLists:
      i % 3 === 0
        ? [
            { name: '確認工作範圍', checked: true, status: 0 },
            { name: '完成執行及驗證', checked: r[2] === 4, status: 0 },
          ]
        : [],
    comments:
      i % 2 === 0
        ? [
            {
              userId: '02',
              comment: '已確認工作安排，會按計劃跟進。',
              status: 0,
              username: '使用者 02',
              avatar: '',
            },
          ]
        : [],
  }));
}
