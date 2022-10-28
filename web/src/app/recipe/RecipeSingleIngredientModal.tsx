import { api } from "../api"
import { useContext } from "react"
import { AppCacheContext } from "../AppCache"
import { ImageButton } from "../component/input"
import ModalPanel from "../component/Modal"

type RecipeIngredientMoreModalProps = {
  recipeIngredients: string[]
  onClose: () => void
}

export default function RecipeIngredientMoreModal(props: RecipeIngredientMoreModalProps) {

  const { grocery, inventory } = useContext(AppCacheContext)
  const groceryItems = grocery.filter(i => !i.is_collected).map(i => i.item)
  const inventoryItems = inventory.map(i => i.item)

  return <ModalPanel closeCallback={props.onClose}>
    <div className="recipe-ingredient-modal-root">
      <div className="recipe-ingredient-modal-missing-ingredients">
        <div className="recipe-ingredient-modal-missing-ingredients-grid">

          <div className="format-font-small">Ingredients</div>
          <div className="recipe-ingredient-modal-button-header format-font-subscript">Grocery</div> <div className="recipe-ingredient-modal-button-header format-font-subscript">Inventory</div>

          {props.recipeIngredients.map((ingre, i) => {
            return (<div key={i} className="recipe-ingredient-missing-element">
              <div className="recipe-ingredient-missing-item"><span className="format-font-small">{ingre}</span></div>

              <ImageButton alt="add to grocery" src="/icons/cart-add.svg" className="recipe-ingredient-modal-add-ingredient"
                disabled={groceryItems.indexOf(ingre) >= 0} onClick={e => {
                  api.post("grocery/create", { "item": ingre })
                }} />

              <ImageButton alt="add to inventory" className="recipe-ingredient-missing-add-inventory"
                src={inventoryItems.indexOf(ingre) >= 0 ? "/icons/no-inventory.svg" : "/icons/inventory.svg"}
                onClick={e => {
                  let model = inventory.filter(i => i.item === ingre)[0]
                  if (model) {
                    api.post("inventory/update", { "uid": model.uid, "in_stock": !model.in_stock })
                  } else {
                    api.post("inventory/create", { "item": ingre })
                  }
                }} />
            </div>)
          })}
        </div>

        <div className="recipe-ingredient-modal-add-all">
          <div className="recipe-ingredient-modal-add-all-grocery">

            <ImageButton alt="add all to grocery" src="/icons/cart-add.svg" className="recipe-ingredient-modal-add-all-btn"
              onClick={e => {
                let missingIngredients = props.recipeIngredients.filter(i => groceryItems.indexOf(i) < 0)
                api.post("grocery/create-batch", { "items": missingIngredients })
              }} />

            <div className="format-font-subscript">Add all to Grocery</div>
          </div>

          <div className="recipe-ingredient-modal-add-all-inventory">

            <ImageButton alt="add all to grocery" src="/icons/inventory.svg" className="recipe-ingredient-modal-add-all-btn"
              onClick={e => {
                let missingIngredients = props.recipeIngredients.filter(i => inventoryItems.indexOf(i) < 0)
                api.post("inventory/create-batch", { "items": missingIngredients })
              }} />

            <div className="format-font-subscript">Add all to Inventory</div>
          </div>
        </div>

      </div>
    </div>
  </ModalPanel>
}