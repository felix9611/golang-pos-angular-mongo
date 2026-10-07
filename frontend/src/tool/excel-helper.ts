// import * as XLSXS from 'xlsx-style'
import * as XLSX from 'xlsx'
import { saveAs } from 'file-saver'
import { NzUploadFile } from 'ng-zorro-antd/upload'

export function readExcelFile(file: any) {
  return new Promise((resolve: any, reject: any) => {
    // check if real file
    const realFile: File | undefined = file.originFileObj || file.raw || file
    if (!(realFile instanceof Blob)) {
      return reject(new Error('Invalid file format. File must be a Blob or File.'))
    }

    // get file type / sub name
    const extension = realFile.name.split('.').pop()?.toLowerCase()
    const allowedTypes = ['xlsx', 'xls', 'csv']

    if (!allowedTypes.includes(extension!)) {
      return reject({ message: 'Format not supported! Only Excel files are allowed.' });
    }

    const reader: FileReader = new FileReader()
    const result: any = []

    reader.onload = (e: any) => {
      const data = e.target?.result
      const wb = XLSX.read(data, { type: 'array' })

      wb.SheetNames.forEach((sheetName: string) => {
        result.push({
          sheetName,
          sheet: XLSX.utils.sheet_to_json(wb.Sheets[sheetName])
        })
      })

      resolve(result.length > 1 ? result[0] : result[0].sheet)
    };

    reader.onerror = (error: any) => reject(error)

    // reading doc
    reader.readAsArrayBuffer(realFile)
  })
}


export function formatJson (header: any, filterVal: any, jsonData: any) {
    const updatedArray = jsonData.map((obj: any) => {
      const newObj: any = {}
  
      header.forEach((oldKey: string, index: number) => {
        
        const newKey = filterVal[index]
        console.log(obj[oldKey])
        if (obj[oldKey] !== undefined) {
          newObj[newKey] = obj[oldKey]
        }
      })
      return newObj
    })
    
    return updatedArray
  
    /* return jsonData.map((v: any)=> {
      const obj: any = {}
      header.forEach((h, i) => { 
        const anyD: any = [filterVal[i]]
        
        const newData: any = v[h]
        console.log(newData, 'data')
        obj[anyD] = newData
        
      })
      return obj
    }) */
}
  
export function formatJsonToSheet(filterVal: string[], jsonData: any[]) {
  if (!Array.isArray(jsonData) || !jsonData.length) {
    return []
  }

  return jsonData.map((row: any) =>
    filterVal.map((key: string) => {
      const value = key.split('.').reduce((obj, k) => (obj || {})[k], row)
    
      return value !== undefined && value !== null ? value : ''
    })
  )
}


export function downloadTempExcelFile(
    excelHeader: any, 
    fileName: string, 
    excelStyle?: any,
    headerColSeetting? : any
  ) {
    const ws = XLSX.utils.aoa_to_sheet([excelHeader])
    if (excelStyle) {
      if (excelStyle) {
        for (const [key] of Object.entries(ws)) {
          if (key !== '!cols' && key !== '!ref') {
            ws[key].s = excelStyle
          }
        }
      }
      if(headerColSeetting) {
        ws['!cols'] = headerColSeetting['!cols']
      }
    }
    XLSX.utils.sheet_add_aoa(ws, [], { origin: 'A2' })
    let wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1')
    const wopts: any = { bookType: 'xlsx', bookSST: false , type: 'binary' }
    const fileEx = XLSX.write(wb, { ...wopts })
    saveAs(new Blob([s2ab(fileEx)],{type:""}), fileName)
  }
export function saveJsonToExcel(
  headers: any, 
  data: any, 
  excelHeader: any, 
  fileName: string, 
  excelStyle?: any,
  headerColSeetting?: any
) {
  
  const dataSet = formatJsonToSheet(headers, data)
  const fullData = [excelHeader, ...dataSet]
  const ws = XLSX.utils.aoa_to_sheet(fullData)

  if (headerColSeetting && headerColSeetting['!cols']) {
    ws['!cols'] = headerColSeetting['!cols']
  }

  if (excelStyle) {
    for (const key of Object.keys(ws)) {
      if (!key.startsWith('!')) {
        ws[key].s = excelStyle
      }
    }
  }

  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Sheet1')
  const fileEx = XLSX.write(wb, { bookType: 'xlsx', type: 'array' })
  
  const blob = new Blob([fileEx], { 
    type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8' 
  })
  const exportFileName = fileName.endsWith('.xlsx') ? fileName : `${fileName}.xlsx`
  saveAs(blob, exportFileName)
}
  
function s2ab(s: any) {
    var buf = new ArrayBuffer(s.length)
    var view = new Uint8Array(buf)
    for (var i=0; i!=s.length; ++i) view[i] = s.charCodeAt(i) & 0xFF
    return buf
}