import "./inventory.css"
import { ImageButton, InputTextBox } from "../../component/input"
import { api } from "../../core/api";
import { InventoryItemModel } from "./schema";

type InventoryItemProps = {
  model: InventoryItemModel
  onShowMore: (i: InventoryItemModel) => void
}

function InventoryItem(props: InventoryItemProps) {

  return (<div className="inventory-item-root">

    <InputTextBox value={props.model.item} className="inventory-item-input"
      inputAttrs={{ autoCapitalize: "none" }} onChange={s => {
        api.post("inventory/update", { uid: props.model.uid, item: s })
      }} />

    <ImageButton src="/icons/more-vertical.svg" className="inventory-item-show-more" alt="inventory-show-more"
      onClick={() => { props.onShowMore(props.model) }} />

  </div>)
}

export default InventoryItem