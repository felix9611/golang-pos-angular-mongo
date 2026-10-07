const monthMap: Record<string, number> = {
  "Jan": 0, "Feb": 1, "Mar": 2, "Apr": 3, "May": 4, "Jun": 5,
  "Jul": 6, "Aug": 7, "Sep": 8, "Oct": 9, "Nov": 10, "Dec": 11
}

export function transformData(rawData: any[], chartType: string, showInLegend: boolean, keyName: string, valueName: string, dateKeyName: string[]): any[] {
  const monthMap: Record<string, number> = {
      "Jan": 0, "Feb": 1, "Mar": 2, "Apr": 3, "May": 4, "Jun": 5,
      "Jul": 6, "Aug": 7, "Sep": 8, "Oct": 9, "Nov": 10, "Dec": 11
  };

  const dataMap: Record<string, { name: string; type: string; showInLegend: boolean ; xValueType?: string, dataPoints: { x: number, y: number }[] }> = {}

  rawData.forEach((data: any) => {
    if (!data[dateKeyName[0]] || !data[dateKeyName[1]] || data[valueName] === 0) return;
    
    if (!dataMap[data[keyName]]) {
      dataMap[data[keyName]] = { name: data[keyName], type: chartType, showInLegend: showInLegend, xValueType: "dateTime", dataPoints: [] };
    }

    const year = parseInt(data[dateKeyName[0]], 10);
    const month = monthMap[data[dateKeyName[1]]];
    
    dataMap[data[keyName]].dataPoints.push({ x: new Date(year, month).getTime(), y: data[valueName] });
  })

  return Object.values(dataMap);
}

export function transformDataNoDate(rawData: any[], chartType: string, showInLegend: boolean, keyName: string, valueName: string): any[] {


  let dataMap: any = { type: chartType, showInLegend: showInLegend, dataPoints: [] }

  rawData.forEach((data: any) => {
    if (data[valueName] === 0) return;
    
    console.log(data[valueName])
    dataMap.dataPoints.push({ label: data[keyName], y: data[valueName] });
  })

  return [dataMap]
}

export function transformDataPointsOnly(rawData: any[], keyName: string, valueName: string): any[] {


  let dataPoints: any = []

  rawData.forEach((data: any) => {
    if (data[valueName] === 0) return;
    

    dataPoints.push({ label: data[keyName], y: data[valueName] });
  })

  return dataPoints
}

export function transformDate(rawData: any[], valueName: string, dateKeyName: string[]): any[] {
  const monthMap: Record<string, number> = {
      "January": 0, "February": 1, "March": 2, "April": 3, "May": 4, "June": 5,
      "July": 6, "August": 7, "September": 8, "October": 9, "November": 10, "December": 11
  };

  let dataPoints: any = []

  // const dataMap: Record<string, { name: string; type: string; showInLegend: boolean ; xValueType?: string, dataPoints: { x: number, y: number }[] }> = {}

  rawData.forEach((data: any) => {
   /* if (!data[dateKeyName[0]] || !data[dateKeyName[1]] || data[valueName] === 0) return;
    
    if (!dataMap[data[keyName]]) {
      dataMap[data[keyName]] = { name: data[keyName], type: chartType, showInLegend: showInLegend, xValueType: "dateTime", dataPoints: [] };
    } */

    const year = parseInt(data[dateKeyName[0]], 10);
    const month = monthMap[data[dateKeyName[1]]];
    
    dataPoints.push({ x: new Date(year, month).getTime(), y: data[valueName] });
  })

  return Object.values(dataPoints);
}