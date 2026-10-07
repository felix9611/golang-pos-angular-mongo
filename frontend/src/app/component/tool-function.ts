export function findMenuItem(roleData: any, name: string, path: string) {
    for (const role of roleData) {
      const found = role.menuLists.find((menu: any) => 
        menu.name.toLowerCase() === name.toLowerCase() && 
        menu.path.toLowerCase().includes(path.toLowerCase())
      )
      if (found) {
        return found
      }
    }
    return null // If no match is found
}


export function findMenuItemByPath(roleData: any, path: string) {
  for (const role of roleData) {
    const found = role.menuLists.find((menu: any) => 
      menu.path.toLowerCase().includes(path.toLowerCase())
    )
    if (found) {
      return found
    }
  }
  return null
}