import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { api } from "../api";
import { InputTextBox } from "../component/input";
import InventoryItem from "./InventoryItem";
import InventoryModal from "./InventoryModal";
import { InventoryCategories, InventoryItemModel } from "./schema";

export function InventoryMenu() {

  const { filter } = useParams()

  return (
    <div className="inventory-list-sub-menu">
      <Link to="/inventory">
        <img alt="All inventory" src="/icons/inventory.svg" id="inventory-list-inventory-btn"
          className={"inventory-list-sub-menu-icon" + (filter === undefined ? " theme-primary-light" : "")} />
      </Link>
      <Link to="/inventory/cooking" >
        <img alt="Cooking inventory" src="/icons/cooking.svg" id="inventory-list-cooking-btn"
          className={"inventory-list-sub-menu-icon" + (filter === "cooking" ? " theme-primary-light" : "")} />
      </Link>
      <Link to="/inventory/spices" >
        <img alt="Spices inventory" src="/icons/spices.svg" id="inventory-list-spices-btn"
          className={"inventory-list-sub-menu-icon" + (filter === "spices" ? " theme-primary-light" : "")} />
      </Link>
      <Link to="/inventory/sauces" >
        <img alt="Sauces and Condiments inventory" src="/icons/sauces.svg" id="inventory-list-sauces-btn"
          className={"inventory-list-sub-menu-icon" + (filter === "sauces" ? " theme-primary-light" : "")} />
      </Link>
      <Link to="/inventory/non-cooking" >
        <img alt="Non-cooking inventory" src="/icons/non-cooking.svg" id="inventory-list-non-cooking-btn"
          className={"inventory-list-sub-menu-icon" + (filter === "non-cooking" ? " theme-primary-light" : "")} />
      </Link>
    </div>
  )
}

function InventoryList() {

  const { filter } = useParams()
  const nav = useNavigate()
  const [items, setItems] = useState<InventoryItemModel[]>([])
  const [currentItem, setCurrentItem] = useState<InventoryItemModel | null>(null)

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

  const getDisplayItems = () => {
    return filter ? items.filter(i => { return i.category === filter }) : items
  }

  const refresh = () => {
    api.post("inventory/list").then(rsp => {
      setItems(rsp.data)
    })
  }

  useEffect(() => { refresh() }, [])
  useEffect(() => {
    if (filter && !InventoryCategories.has(filter)) {
      nav("/inventory")
    }
  }, [filter, nav])

  return (<div className="inventory-list-root">

    <div className="inventory-list-add-item theme-primary-light">
      <InputTextBox value="" className="inventory-list-add-item-input" clearOnChange={true}
        placeholder="Add inventory item" onChange={s => {
          const newItem = s.trim().toLowerCase()
          if (newItem.length === 0 || items.map(i => i.item).indexOf(newItem) >= 0) {
            return
          }

          api.post("inventory/create", { item: s, category: filter }).then(rsp => {
            updateInventoryList([...items, rsp.data])
          })
        }} />
    </div>

    <div className="inventory-list-items">{getDisplayItems().map(i => {
      return (<InventoryItem key={i.uid} model={i}
        onUpdate={i => updateInventoryItem(i)}
        onShowMore={i => { setCurrentItem(i) }} />)
    })}</div>

    {currentItem != null ?
      <InventoryModal model={currentItem} onClose={() => { setCurrentItem(null) }}
        onUpdate={i => { updateInventoryItem(i); setCurrentItem(i); }}
        onUnstock={vic => updateInventoryList(items.filter(i => i.uid !== vic.uid))} />
      : null}
  </div>)

}

export default InventoryList