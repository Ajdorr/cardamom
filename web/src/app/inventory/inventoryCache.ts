import { api } from "../api"
import { InventoryItemModel } from "./schema"

const inventoryListKey = "inventory.list"

export function setInventoryCache(newItems: InventoryItemModel[]) {
  localStorage.setItem(inventoryListKey, JSON.stringify(newItems))
}

export function getInventoryCache(): Promise<InventoryItemModel[]> {
  return new Promise<InventoryItemModel[]>((resolve, reject) => {
    const inventoryCache = localStorage.getItem(inventoryListKey)
    if (inventoryCache === null) {
      api.post("inventory/list").then(rsp => {
        localStorage.setItem(inventoryListKey, JSON.stringify(rsp.data))
        return resolve(rsp.data)
      }).catch(reject)
    } else {
      return resolve(JSON.parse(inventoryCache))
    }
  })
}

export function invalidateInventoryCache() {
  localStorage.removeItem(inventoryListKey)
}

export function updateInventoryCache(...newItems: InventoryItemModel[]) {
  const inventoryCache = localStorage.getItem(inventoryListKey)
  if (inventoryCache === null) {
    return
  }

  let cachedItems: InventoryItemModel[] = JSON.parse(inventoryCache)
  for (let newItem of newItems) {
    let ndx = cachedItems.map(i => i.uid).indexOf(newItem.uid)
    if (ndx >= 0) {
      cachedItems[ndx] = newItem
    } else {
      cachedItems.push(newItem)
    }
  }

  localStorage.setItem(inventoryListKey, JSON.stringify(cachedItems))
}

export function deleteFromInventoryCache(victim: InventoryItemModel) {
  const inventoryCache = localStorage.getItem(inventoryListKey)
  if (inventoryCache === null) {
    return
  }

  let items: InventoryItemModel[] = JSON.parse(inventoryCache)
  localStorage.setItem(
    inventoryListKey, JSON.stringify(items.filter(i => { return i.uid !== victim.uid })))
}
