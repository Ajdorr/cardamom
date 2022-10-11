import "./inventory.css"
import { ImageButton, InputTextBox } from "../component/input"
import { api } from "../api";
import { InventoryItemModel } from "./schema";

type InventoryItemProps = {
  model: InventoryItemModel
  onUpdate: (i: InventoryItemModel) => void
  onRemove: (i: InventoryItemModel) => void
}

function InventoryItem(props: InventoryItemProps) {

  return (<div className="inventory-item-root">

    <InputTextBox value={props.model.item} className="inventory-item-input"
      onChange={s => {
        api.post("inventory/update", { uid: props.model.uid, item: s }).then(rsp => {
          props.onUpdate(rsp.data)
        })
      }} />

    <ImageButton src="icons/delete.svg" className="inventory-item-unstock" alt="unstock"
      onClick={() => {
        api.post("inventory/update", { uid: props.model.uid, in_stock: false }).then(rsp => {
          props.onRemove(rsp.data)
        })
      }} />

  </div>)
}

export default InventoryItem