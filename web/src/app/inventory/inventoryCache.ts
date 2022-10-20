import { api } from "../api"
import { InventoryItemModel } from "./schema"

const inventoryListKey = "inventory.list"

export const setInventory = (inventory: string[]) => {
  localStorage.setItem(inventoryListKey, JSON.stringify(inventory))
}

export function getInventory(): Promise<string[]> {
  return new Promise<string[]>((resolve, reject) => {
    const inventoryCache = localStorage.getItem(inventoryListKey)
    if (inventoryCache === null) {
      api.post("inventory/list").then(rsp => {
        const inventoryList = rsp.data.map((i: InventoryItemModel) => i.item)
        localStorage.setItem(inventoryListKey, JSON.stringify(inventoryList))
        return resolve(inventoryList)
      }).catch(reject)
    } else {
      return resolve(JSON.parse(inventoryCache))
    }
  })
}