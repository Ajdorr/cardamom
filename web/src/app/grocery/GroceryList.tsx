import "./grocery.css"
import { Component } from 'react'
import { api } from '../api';
import { ImageButton } from '../component/input';
import { GroceryItem, AddGroceryItem } from './GroceryItem'
import { GroceryItemModel } from './schema'
import { ModifiableDropDown } from "./component/DropDown";

type GroceryState = {
  items: GroceryItemModel[],
  uniqueStores: string[],
  selectedStore: string
}

class GroceryList extends Component<{}, GroceryState> {

  constructor(props: any) {
    super(props)
    this.state = {
      items: [],
      uniqueStores: [],
      selectedStore: ""
    }
  }

  updateGroceryItem(item: GroceryItemModel) {
    this.updateGroceryList(this.state.items.map(i => {
      if (i.uid === item.uid) {
        return item
      } else {
        return i
      }
    }))
  }

  updateGroceryList(itemList: GroceryItemModel[]) {
    const stores = itemList.map(i => i.store)
      .filter((s, i, a) => a.indexOf(s) === i && s.length > 0)
    stores.sort()
    this.setState({ items: itemList, uniqueStores: stores })
  }

  clearAll() {
    api.post("grocery/clear").then(rsp => {
      this.refresh()
    }).catch(e => {
      console.log(e) // FIXME
    })
  }

  refresh() {
    api.post("grocery/list").then(rsp => {
      this.updateGroceryList(rsp.data)
    }).catch(e => {
      console.log(e) // FIXME
    })
  }

  componentDidMount() {
    this.refresh()
  }

  render() {
    const collectedItems = this.state.items.filter(i => i.is_collected)
    var displayedItems = this.state.items.filter(i => !i.is_collected)
    displayedItems = (this.state.selectedStore !== "") ?
      displayedItems.filter(i => i.store === this.state.selectedStore) : displayedItems

    return (
      <div className="grocery-list-root">

        <ModifiableDropDown options={this.state.uniqueStores} value={this.state.selectedStore} id={"grocery-list-store"}
          className="grocery-list-store theme-primary-light" displayClear={true}
          onChange={s => this.setState({ selectedStore: s })}
        />

        <AddGroceryItem id="grocery-list-add-item" store={this.state.selectedStore}
          onAdd={newItem => {
            // Only update if its a new item
            if (this.state.items.map(i => i.item).indexOf(newItem.item) < 0) {
              this.updateGroceryList([newItem, ...this.state.items])
            } else {
              this.updateGroceryItem(newItem)
            }
          }
          } />

        <div className="grocery-list-items">
          {displayedItems.length > 0 ? displayedItems.map(i => {
            return (<GroceryItem key={i.uid} model={i} stores={this.state.uniqueStores}
              onUpdate={i => this.updateGroceryItem(i)}
            />)
          }) : <div className="grocery-list-items-empty"><span>No grocery items in your list</span></div>}
        </div>

        <div className="grocery-list-collected-divider theme-primary-light">
          <div className="grocery-list-collected-space"> </div>
          <ImageButton className="grocery-list-collected-clear-all"
            src="icons/delete-all.svg" alt="clear" onClick={e => this.clearAll()} />
        </div>

        <div className="grocery-list-collected-items">
          {collectedItems.length > 0 ? collectedItems.map(i => {
            return (<CollectedGroceryItem key={i.uid} uid={i.uid} item={i.item} store={i.store}
              onUndo={i => this.updateGroceryItem(i)} />)
          }) : <div className="grocery-list-collected-empty">No collected items</div>
          }
        </div>
      </div>
    )
  }
}

type CollectedGroceryItemProps = {
  uid: string,
  item: string,
  store: string,
  onUndo: (item: GroceryItemModel) => void
}

function CollectedGroceryItem(props: CollectedGroceryItemProps) {
  return (<div className="grocery-list-collected-root">
    <span className="grocery-collected-item">{props.item}</span>
    <span className="grocery-collected-store">{props.store}</span>
    <ImageButton src="icons/undo.svg" alt="undo" onClick={e => {

      api.post("grocery/collect", { uid: props.uid, is_collected: false }).then(rsp => {
        props.onUndo(rsp.data)
      }).catch(e => {
        console.log(e)
      })
    }} />
  </div>)
}

export default GroceryList