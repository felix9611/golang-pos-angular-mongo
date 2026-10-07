export function buildTreeForUI(data: any[]) {
    const map = new Map<string, any>()
    const tree: any[] = []
  
    // Initialize Map with transformed objects
    data.forEach(item => {
      map.set(String(item._id), { 
        title: item.name, 
        key: String(item._id), 
        children: [] 
      })
    })
  
    // Build hierarchical structure
    data.forEach(item => {
      const parentId = String(item.mainId)
      if (parentId && map.has(parentId)) {
        map.get(parentId)!.children.push(map.get(String(item._id)));
      } else {
        tree.push(map.get(String(item._id)))
      }
    })
  

    function markLeafNodes(nodes: any[]) {
      nodes.forEach(node => {
        if (node.children.length === 0) {
          node.isLeaf = true
        } else {
          markLeafNodes(node.children)
        }
      });
    }
  
    markLeafNodes(tree)
  
    return tree
}