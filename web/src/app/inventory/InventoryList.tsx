import { useContext, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { api } from "../api";
import { AppCacheContext } from "../AppCache";
import { InputTextBox } from "../component/input";
import InventoryItem from "./InventoryItem";
import InventoryModal from "./InventoryModal";
import { InventoryCategories } from "./schema";

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
      <Link to="/inventory/non-perishables" >
        <img alt="Non-perishables inventory" src="/icons/non-perishables.svg" id="inventory-list-non-perishables-btn"
          className={"inventory-list-sub-menu-icon" + (filter === "non-perishables" ? " theme-primary-light" : "")} />
      </Link>
      <Link to="/inventory/non-cooking" >
        <img alt="Non-cooking inventory" src="/icons/non-cooking.svg" id="inventory-list-non-cooking-btn"
          className={"inventory-list-sub-menu-icon" + (filter === "non-cooking" ? " theme-primary-light" : "")} />
      </Link>
    </div>
  )
}

function InventoryList() {

  const nav = useNavigate()
  const { filter } = useParams()
  const { inventory } = useContext(AppCacheContext)
  const [currentItem, setCurrentItem] = useState<string | null>(null)

  const displayFilter = (f: string): string => {
    return f.split("-").map(s => s[0].toUpperCase() + s.substring(1).toLowerCase()).join("-")
  }

  useEffect(() => { if (filter && !InventoryCategories.has(filter)) nav("/inventory") }, [filter, nav])
  const displayItems = filter ? inventory.filter(i => { return i.category === filter }) : inventory

  return (<div className="inventory-list-root">

    <div className="inventory-list-add-item theme-primary-light">
      <InputTextBox value="" className="inventory-list-add-item-input" clearOnChange={true}
        inputAttrs={{ autoCapitalize: "none" }}
        placeholder={`Add to ${filter ? displayFilter(filter) : "Inventory"}`} onChange={s => {
          const newItem = s.trim().toLowerCase()
          if (newItem.length === 0 || inventory.map(i => i.item).indexOf(newItem) >= 0) {
            return
          }

          api.post("inventory/create", { item: s, category: filter })
        }} />
    </div>

    <div className="inventory-list-items">{
      displayItems.map(i => {
        return (<InventoryItem key={i.uid} model={i} onShowMore={i => { setCurrentItem(i.uid) }} />)
      })}
      {displayItems.length === 0 ?
        <div className="inventory-list-empty">{
          filter ? `Nothing in ${displayFilter(filter)}` : "Nothing in your inventory"
        }</div>
        : null
      }
    </div>

    {currentItem != null ?
      <InventoryModal model={inventory.filter(i => i.uid === currentItem)[0]}
        closeCallback={() => { setCurrentItem(null) }} />
      : null}
  </div>)

}

export default InventoryList