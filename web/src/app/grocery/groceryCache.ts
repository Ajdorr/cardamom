import { api } from "../api"
import { GroceryItemModel } from "./schema"

const groceryListKey = "grocery.list"

export const setGroceryCache = (newItems: GroceryItemModel[]) => {
  localStorage.setItem(groceryListKey, JSON.stringify(newItems))
}

export function getGroceryCache(): Promise<GroceryItemModel[]> {
  return new Promise<GroceryItemModel[]>((resolve, reject) => {
    const cache = localStorage.getItem(groceryListKey)
    if (cache === null) {
      api.post("grocery/list").then(rsp => {
        localStorage.setItem(groceryListKey, JSON.stringify(rsp.data))
        return resolve(rsp.data)
      }).catch(reject)
    } else {
      return resolve(JSON.parse(cache))
    }
  })
}

export function invalidateGroceryCache(newItems: GroceryItemModel[]) {
  localStorage.removeItem(groceryListKey)
}

export const updateGroceryCache = (...newItems: GroceryItemModel[]) => {
  const cache = localStorage.getItem(groceryListKey)
  if (cache === null) {
    return
  }

  let cachedItems: GroceryItemModel[] = JSON.parse(cache)
  for (let newItem of newItems) {
    let ndx = cachedItems.map(i => i.uid).indexOf(newItem.uid)
    if (ndx >= 0) {
      cachedItems[ndx] = newItem
    } else {
      cachedItems.push(newItem)
    }
  }

  localStorage.setItem(groceryListKey, JSON.stringify(cachedItems))
}

export const deleteGroceryCache = (victim: GroceryItemModel) => {
  const cache = localStorage.getItem(groceryListKey)
  if (cache === null) {
    return
  }

  let items: GroceryItemModel[] = JSON.parse(cache)
  localStorage.setItem(
    groceryListKey, JSON.stringify(items.filter(i => { return i.uid !== victim.uid })))
}
