import { useEffect, useState } from "react";
import { api } from "../api";
import { InputTextBox } from "../component/input";
import InventoryItem from "./InventoryItem";
import { InventoryItemModel } from "./schema";

function InventoryList() {

  const [items, setItems] = useState<InventoryItemModel[]>([])

  const updateInventoryList = (newItems: InventoryItemModel[]) => {
    setItems(newItems)
  }

  const updateInventoryItem = (item: InventoryItemModel) => {
    updateInventoryList(items.map(i => {
      if (i.uid === item.uid) {
        return item
      } else {
        return i
      }
    }))
  }

  const refresh = () => {
    api.post("inventory/list").then(rsp => {
      setItems(rsp.data)
    }).catch(e => console.log(e))
  }

  useEffect(() => { refresh() }, [])

  return (<div className="inventory-list-root">

    <div className="inventory-list-add-item theme-primary-light">
      <InputTextBox value="" className="inventory-list-add-item-input" clearOnChange={true}
        placeholder="Add grocery item" onChange={s => {
          const newItem = s.trim()
          if (newItem.length === 0) {
            return
          }

          api.post("inventory/create", { item: s }).then(rsp => {
            updateInventoryList([...items, rsp.data])
          }).catch(e => console.log(e))
        }} />
    </div>

    <div className="inventory-list-items">{items.map(i => {
      return (<InventoryItem key={i.uid} model={i}
        onUpdate={i => updateInventoryItem(i)}
        onRemove={vic => updateInventoryList(items.filter(i => i.uid !== vic.uid))} />)
    })}</div>
  </div>)

}

export default InventoryList