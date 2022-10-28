import { useState } from "react"
import { api } from "../api"
import { ImageButton, InputTextBox } from "../component/input"
import { SwipeIndicatorWidget } from "../component/widget"
import { ModifiableDropDown } from "./component/StoreDropDown"
import { GroceryItemModel } from './schema'

const DragToDeleteTolerance = 60

type AddGroceryItemProps = {
  id: string
  store: string
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
    })
  }

  return (<InputTextBox id={props.id} value="" className="grocery-item-add-root theme-primary-light"
    inputAttrs={{ autoCapitalize: "none" }}
    placeholder="Add a grocery" clearOnChange={true} onChange={s => save(s)} />)
}

type GroceryItemProps = {
  model: GroceryItemModel
  stores: string[]
}

type UpdateRequest = {
  uid: string,
  item?: string,
  store?: string,
}

export function GroceryItem(props: GroceryItemProps) {

  const [initX, setInitX] = useState(0)
  const [deltaX, setDeltaX] = useState(0)
  const [initY, setInitY] = useState(0)
  const cssStyle = (deltaX !== 0) ? { transform: `translateX(${deltaX}px)` } : undefined

  const onUpdate = (req: UpdateRequest) => {
    api.post("grocery/update", req)
  }

  return (<div style={cssStyle} className="grocery-item-root"
    onTouchStart={e => { setInitX(e.touches[0].clientX); setInitY(e.touches[0].clientY) }}
    onTouchMove={e => {
      if (Math.abs(e.touches[0].clientY - initY) < 30) {
        setDeltaX(e.touches[0].clientX - initX);
      } else {
        setDeltaX(0);
      }
    }}
    onTouchEnd={e => {
      if (Math.abs(deltaX) >= DragToDeleteTolerance) {
        api.post("grocery/delete", { uid: props.model.uid })
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
    <ImageButton className="grocery-item-collect" alt="collect" src="/icons/done.svg"
      onClick={e => api.post("grocery/collect", { uid: props.model.uid, is_collected: true })} />
    <InputTextBox value={props.model.item} inputAttrs={{ autoCapitalize: "none" }} className="grocery-item-input"
      onChange={i => onUpdate({ uid: props.model.uid, item: i })} />
    <ModifiableDropDown className="grocery-item-store" value={props.model.store} options={props.stores}
      dropDownButtonOnLeft={true} placeholder="Store"
      onChange={s => onUpdate({ uid: props.model.uid, store: s })} />
  </div >
  )
}

export default GroceryItem
