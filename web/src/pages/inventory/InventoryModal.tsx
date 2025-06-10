import "./inventory.css"
import ModalPanel from "../../component/Modal"
import { InventoryCategories, InventoryItemModel } from "./schema"
import { FormDropDown } from "../../component/form"
import { api } from "../../core/api"
import { ImageButton, InputTextBox } from "../../component/input"
import { useContext } from "react"
import { AppCacheContext } from "@core/AppCache"

type InventoryModalProps = {
  model: InventoryItemModel
  closeCallback: () => void
}

export default function InventoryModal(props: InventoryModalProps) {

  const { grocery } = useContext(AppCacheContext)

  const updateItem = (input: { item?: string, category?: string, in_stock?: boolean }) => {
    api.post("inventory/update", { uid: props.model.uid, ...input })
  }

  return (<ModalPanel closeCallback={props.closeCallback}>
    <div className="inventory-modal-root">
      <InputTextBox value={props.model.item} className="inventory-modal-item"
        inputAttrs={{ autoCapitalize: "none" }} onChange={e => { updateItem({ item: e }) }} />
      <FormDropDown label="Category" options={InventoryCategories} value={props.model.category}
        className="inventory-modal-category" onChange={e => { updateItem({ category: e }) }} />
      <div className="inventory-modal-actions">
        <ImageButton alt="Add to grocery" src="/icons/cart-add.svg" className="inventory-modal-add-to-grocery-btn"
          disabled={grocery.map(i => i.item).indexOf(props.model.item) >= 0}
          onClick={e => { api.post("grocery/create", { item: props.model.item }) }} />
        <ImageButton alt="Delete item" src="/icons/delete.svg" className="inventory-modal-delete-btn"
          onClick={e => { updateItem({ in_stock: false }); props.closeCallback() }} />
      </div>
    </div>
  </ModalPanel>)
}