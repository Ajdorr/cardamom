import { AxiosResponse } from "axios";
import { createContext, useCallback, useEffect, useState } from "react";
import { api } from "./api";
import { GroceryItemModel } from "./grocery/schema";
import { InventoryItemModel } from "./inventory/schema";

const cacheValidDurationMs = 15 * 60 * 1000
var cacheTTL = 0

interface CacheContextStruct {
  inventory: InventoryItemModel[]
  grocery: GroceryItemModel[]
}

export const AppCacheContext = createContext<CacheContextStruct>({ grocery: [], inventory: [] })

export interface AppCacheProps {
  children?: React.ReactNode;
}

export default function AppCache(props: AppCacheProps) {

  const [groceryItems, setGroceryItems] = useState<GroceryItemModel[]>([])
  const [inventoryItems, setInventoryItems] = useState<InventoryItemModel[]>([])

  const updateGrocery = useCallback((...items: GroceryItemModel[]) => {
    let newItems = [...groceryItems]
    for (let item of items) {
      let ndx = newItems.map(i => i.uid).indexOf(item.uid)
      if (ndx >= 0) {
        newItems[ndx] = item
      } else {
        newItems.push(item)
      }
    }
    setGroceryItems(newItems)
  }, [groceryItems])

  const updateInventory = useCallback((...items: InventoryItemModel[]) => {
    let newItems = [...inventoryItems]
    for (let item of items) {
      let ndx = newItems.map(i => i.uid).indexOf(item.uid)
      if (ndx >= 0) {
        newItems[ndx] = item
      } else {
        newItems.push(item)
      }
    }

    setInventoryItems(newItems.filter(i => i.in_stock))
  }, [inventoryItems])

  const onResponse = useCallback((rsp: AxiosResponse<any, any>) => {

    if (cacheTTL < Date.now()) {
      cacheTTL = Date.now() + cacheValidDurationMs
      api.post("grocery/list")
      api.post("inventory/list")
      return
    }

    switch (rsp.config.url) {
      // Grocery
      case "grocery/create":
      case "grocery/update":
        updateGrocery(rsp.data)
        break;
      case "grocery/collect":
        updateGrocery(rsp.data.grocery_item)
        updateInventory(rsp.data.inventory_item)
        break;
      case "grocery/create-batch":
        updateGrocery(...rsp.data)
        break;
      case "grocery/list":
        setGroceryItems(rsp.data)
        break;
      case "grocery/delete":
        let reqData = JSON.parse(rsp.config.data)
        if (reqData) {
          setGroceryItems(groceryItems.filter(i => i.uid !== reqData.uid))
        }
        break;
      case "grocery/clear":
        setGroceryItems(groceryItems.filter(i => !i.is_collected))
        break;

      // Inventory
      case "inventory/create":
      case "inventory/update":
        updateInventory(rsp.data)
        break;
      case "inventory/create-batch":
        updateInventory(...rsp.data)
        break;
      case "inventory/list":
        setInventoryItems(rsp.data)
        break;
      case "inventory/delete":
        setInventoryItems(inventoryItems.filter(i => i.uid !== rsp.data.uid))
        break;
      default:
        break;
    }
  }, [groceryItems, inventoryItems, updateGrocery, updateInventory])

  // eslint-disable-next-line
  useEffect(() => {
    const intId = api.interceptors.response.use(
      rsp => { onResponse(rsp); return Promise.resolve(rsp) },
      err => { return Promise.reject(err) }
    )

    // Initially cache the grocery and inventory lists
    return () => api.interceptors.response.eject(intId)
  }, [onResponse])

  // eslint-disable-next-line
  useEffect(() => {
    api.post("grocery/list")
    api.post("inventory/list")
  }, [])

  return (
    <AppCacheContext.Provider value={{ grocery: groceryItems, inventory: inventoryItems }}>
      {props.children}
    </AppCacheContext.Provider>
  )
}