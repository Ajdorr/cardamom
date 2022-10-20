import { useState } from "react"
import { api } from "../api"
import { ImageButton, InputTextBox } from "../component/input"
import { SwipeIndicatorWidget } from "../component/widget"
import { ModifiableDropDown } from "./component/DropDown"
import { GroceryItemModel } from './schema'

const DragToDeleteTolerance = 60

type AddGroceryItemProps = {
  id: string
  store: string
  onAdd: (item: GroceryItemModel) => void
}

export function AddGroceryItem(props: AddGroceryItemProps) {
  const save = function (s: string) {
    const item = s.trim()
    if (item.length === 0) {
      return
    }

    api.post("grocery/create", {
      item: item,
      store: props.store.length > 0 ? props.store : null
    }).then(rsp => {
      props.onAdd(rsp.data)
    })
  }

  return (<InputTextBox id={props.id} value="" className="grocery-item-add-root theme-primary-light"
    inputAttrs={{ autoCapitalize: "none" }}
    placeholder="Add a grocery" clearOnChange={true} onChange={s => save(s)} />)
}

type GroceryItemProps = {
  model: GroceryItemModel
  stores: string[]
  onUpdate: (i: GroceryItemModel) => void
  onDelete: (i: GroceryItemModel) => void
}

type UpdateRequest = {
  uid: string,
  item?: string,
  store?: string,
}

export function GroceryItem(props: GroceryItemProps) {

  const [initX, setInitX] = useState(0)
  const [deltaX, setDeltaX] = useState(0)
  const cssStyle = (deltaX !== 0) ? { transform: `translateX(${deltaX}px)` } : undefined

  const onUpdate = (req: UpdateRequest) => {

    // Prevent screen jitter
    props.model.item = req.item ? req.item : props.model.item
    props.model.store = req.store ? req.store : props.model.store
    props.onUpdate(props.model)

    api.post("grocery/update", req).then(rsp => {
      props.onUpdate(rsp.data)
    })
  }

  const collectItem = () => {
    api.post("grocery/collect", { uid: props.model.uid, is_collected: true }).then(rsp => {
      props.onUpdate(rsp.data)
    })
  }

  return (<div style={cssStyle} className="grocery-item-root"
    onTouchStart={e => { setInitX(e.touches[0].clientX) }}
    onTouchMove={e => { setDeltaX(e.touches[0].clientX - initX); }}
    onTouchEnd={e => {
      if (Math.abs(deltaX) >= DragToDeleteTolerance) {
        api.post("grocery/delete", { uid: props.model.uid }).then(rsp => {
          props.onDelete(props.model)
        })
      }
      setInitX(0); setDeltaX(0)
    }}
  >
    {deltaX < 0 ?
      <div style={{
        width: `${Math.abs(-deltaX)}px`,
        transform: `translateX(${Math.abs(deltaX)}px)`,
        right: "0"
      }}
        className="grocery-item-delete-indicator">
        {Math.abs(deltaX) > 40 ? <img alt="delete indicator" src="/icons/delete.svg" /> : null}
      </div> : null
    }
    {deltaX > 0 ? <SwipeIndicatorWidget className="grocery-item-collect-indicator"
      deltaX={deltaX} height={40} iconSrc="/icons/done.svg" />
      : null}
    <ImageButton className="grocery-item-collect" alt="collect" src="/icons/done.svg" onClick={e => collectItem()} />
    <InputTextBox value={props.model.item} inputAttrs={{ autoCapitalize: "none" }} className="grocery-item-input"
      onChange={i => onUpdate({ uid: props.model.uid, item: i })} />
    <ModifiableDropDown className="grocery-item-store" value={props.model.store} options={props.stores}
      dropDownButtonOnLeft={true} placeholder="Store"
      onChange={s => onUpdate({ uid: props.model.uid, store: s })} />
  </div >
  )
}

export default GroceryItem
