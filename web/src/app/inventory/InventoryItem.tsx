import "./inventory.css"
import { ImageButton, InputTextBox } from "../component/input"
import { api } from "../api";
import { InventoryItemModel } from "./schema";

type InventoryItemProps = {
  model: InventoryItemModel
  onUpdate: (i: InventoryItemModel) => void
  onShowMore: (i: InventoryItemModel) => void
}

function InventoryItem(props: InventoryItemProps) {

  return (<div className="inventory-item-root">

    <InputTextBox value={props.model.item} className="inventory-item-input"
      onChange={s => {
        api.post("inventory/update", { uid: props.model.uid, item: s }).then(rsp => {
          props.onUpdate(rsp.data)
        })
      }} />

    <ImageButton src="/icons/more-vertical.svg" className="inventory-item-show-more" alt="show-more"
      onClick={() => { props.onShowMore(props.model) }} />

  </div>)
}

export default InventoryItem