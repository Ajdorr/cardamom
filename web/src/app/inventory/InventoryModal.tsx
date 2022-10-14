import "./inventory.css"
import ModalPanel from "../component/Modal"
import { InventoryCategories, InventoryItemModel } from "./schema"
import { FormDropDown } from "../component/form"
import { api } from "../api"
import { ImageButton, InputTextBox } from "../component/input"

type InventoryModalProps = {
  model: InventoryItemModel
  onUpdate: (model: InventoryItemModel) => void
  onUnstock: (model: InventoryItemModel) => void
  onClose: () => void
}

export default function InventoryModal(props: InventoryModalProps) {

  const updateItem = (input: { item?: string, category?: string, in_stock?: boolean }) => {
    api.post("inventory/update", { uid: props.model.uid, ...input }).then(rsp => {
      if (rsp.data.in_stock) {
        props.onUpdate(rsp.data)
      } else {
        props.onUnstock(rsp.data)
        props.onClose()
      }
    })
  }

  return (<ModalPanel onClose={props.onClose}>
    <div className="inventory-modal-root">
      <InputTextBox value={props.model.item} className="inventory-modal-item"
        onChange={e => { updateItem({ item: e }) }} />
      <FormDropDown label="Category" options={InventoryCategories} value={props.model.category}
        className="inventory-modal-category" onChange={e => { updateItem({ category: e }) }} />
      <div className="inventory-modal-delete">
        <ImageButton alt="Delete item" src="/icons/delete.svg" className="inventory-modal-delete-btn"
          onClick={e => { updateItem({ in_stock: false }) }} />
      </div>
    </div>
  </ModalPanel>)
}