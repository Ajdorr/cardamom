import { api } from "../api"
import { InventoryItemModel } from "./schema"

const inventoryListKey = "inventory.list"
const alwaysAvailableIngredients = ["water"]

export const setInventory = (inventory: string[]) => {
  inventory.push(...alwaysAvailableIngredients)
  localStorage.setItem(inventoryListKey, JSON.stringify(inventory))
}

export function getInventory(): Promise<string[]> {
  return new Promise<string[]>((resolve, reject) => {
    const inventoryCache = localStorage.getItem(inventoryListKey)
    if (inventoryCache === null) {
      api.post("inventory/list").then(rsp => {
        let inventory = rsp.data.map((i: InventoryItemModel) => i.item)
        inventory.push(...alwaysAvailableIngredients)
        localStorage.setItem(inventoryListKey, JSON.stringify(inventory))
        return resolve(inventory)
      }).catch(reject)
    } else {
      return resolve(JSON.parse(inventoryCache))
    }
  })
}