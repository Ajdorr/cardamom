import { api } from "../api"
import { ImageButton, InputTextBox } from "../component/input"
import { ModifiableDropDown } from "./component/DropDown"
import { GroceryItemModel } from './schema'

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
    }).catch(e => { console.log(e) })
  }

  return (<InputTextBox id={props.id} value="" className="grocery-item-add-root theme-primary-light"
    placeholder="Add a grocery" clearOnChange={true} onChange={s => save(s)} />)
}

type GroceryItemProps = {
  model: GroceryItemModel
  stores: string[]
  onUpdate: (i: GroceryItemModel) => void
}

type UpdateRequest = {
  uid: string,
  item?: string,
  store?: string,
}

export function GroceryItem(props: GroceryItemProps) {

  const onUpdate = (req: UpdateRequest) => {

    // Prevent screen jitter
    props.model.item = req.item ? req.item : props.model.item
    props.model.store = req.store ? req.store : props.model.store
    props.onUpdate(props.model)

    api.post("grocery/update", req).then(rsp => {
      props.onUpdate(rsp.data)
    }).catch(e => {
      console.log(e) // FIXME
    })
  }

  const collectItem = () => {
    api.post("grocery/collect", { uid: props.model.uid, is_collected: true }).then(rsp => {
      props.onUpdate(rsp.data)
    }).catch(e => {
      console.log(e) // FIXME
    })
  }

  return (<div className="grocery-item-root" >
    <ImageButton className="grocery-item-collect" alt="collect" src="icons/done.svg" onClick={e => collectItem()} />
    <InputTextBox value={props.model.item} className="grocery-item-input" onChange={i => onUpdate({ uid: props.model.uid, item: i })} />
    <ModifiableDropDown className="grocery-item-store" value={props.model.store} options={props.stores}
      onChange={s => onUpdate({ uid: props.model.uid, store: s })} />
  </div>
  )
}

export default GroceryItem
