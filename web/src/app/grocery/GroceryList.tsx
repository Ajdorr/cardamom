import "./grocery.css"
import { useContext, useMemo, useState } from 'react'
import { api } from '../api';
import { ImageButton } from '../component/input';
import { GroceryItem, AddGroceryItem } from './GroceryItem'
import { StoreDropDown } from "./component/StoreDropDown";
import { AppCacheContext } from "../AppCache";

export default function GroceryList() {

  const { grocery } = useContext(AppCacheContext)
  const [selectedStore, setSelectedStore] = useState("")

  const collectedItems = grocery.filter(i => i.is_collected)
  const uniqueStores = useMemo(() => {
    let stores = grocery.map(i => i.store).filter((s, i, a) => a.indexOf(s) === i && s.length > 0)
    stores.splice(0, 0, "")
    return stores
  }, [grocery])

  const displayedItems = useMemo(() => {
    let uncollectedItems = grocery.filter(i => !i.is_collected)
    return (selectedStore !== "") ? uncollectedItems.filter(i => i.store === selectedStore) : uncollectedItems
  }, [grocery, selectedStore])

  return (
    <div className="grocery-list-root">
      <StoreDropDown options={uniqueStores} value={selectedStore} id={"grocery-list-store"}
        className="grocery-list-store theme-primary-light" onChange={setSelectedStore}
      />

      <AddGroceryItem id="grocery-list-add-item" store={selectedStore} />

      <div className="grocery-list-items">
        {displayedItems.length > 0 ? displayedItems.map(i => {
          return (<GroceryItem key={i.uid} model={i} stores={uniqueStores} />)
        }) : <div className="grocery-list-items-empty"><span>No grocery items in your list</span></div>}
      </div>

      <div className="grocery-list-collected-divider theme-primary-light">
        <div className="grocery-list-collected-space"><span>Collected Groceries</span></div>
        <ImageButton className="grocery-list-collected-clear-all"
          src="/icons/delete-all.svg" alt="clear" onClick={e => api.post("grocery/clear")} />
      </div>

      <div className="grocery-list-collected-items">
        {collectedItems.length > 0 ? collectedItems.map(i => {
          return (<CollectedGroceryItem key={i.uid} uid={i.uid} item={i.item} store={i.store} />)
        }) : <div className="grocery-list-collected-empty">No collected items</div>
        }
      </div>
    </div>
  )
}

type CollectedGroceryItemProps = {
  uid: string,
  item: string,
  store: string,
}

function CollectedGroceryItem(props: CollectedGroceryItemProps) {
  return (<div className="grocery-list-collected-root">
    <span className="grocery-collected-item">{props.item}</span>
    <span className="grocery-collected-store">{props.store}</span>
    <ImageButton src="/icons/undo.svg" alt="undo" onClick={e => {
      api.post("grocery/collect", { uid: props.uid, is_collected: false })
    }} />
  </div>)
}