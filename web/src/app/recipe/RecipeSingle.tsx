import "./css/recipe-single.css"
import { useContext, useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { api } from "../api"
import { FormDropDown, FormText, FormTextArea } from "../component/form"
import { ImageButton } from "../component/input"
import { AlwaysAvailableIngredients } from "../inventory/schema"
import { RecipeIngredient, RecipeInstruction } from "./component/RecipeComp"
import { MealTypes, RecipeModel, UpdateIngredient, UpdateRecipe } from "./schema"
import { AppCacheContext } from "../AppCache"
import RecipeSingleIngredientModal from "./RecipeSingleIngredientModal"

export default function RecipeSingle() {

  const nav = useNavigate()
  const { recipeUid } = useParams()
  const { inventory } = useContext(AppCacheContext)

  const [ingreNewNdx, setIngreNewNdx] = useState(-1)
  const [indicatorClass, setIndicatorClass] = useState("theme-indicator-top")
  const [areIngredientsEditable, setEditableIngredients] = useState(false)
  const [displayIngredientMore, setDisplayIngredientMore] = useState(false)

  const [recipe, setRecipe] = useState<RecipeModel>({
    uid: "",
    created_at: "",
    updated_at: "",
    is_trashed: false,
    user_uid: "",
    name: "",
    description: "",
    meal: "dinner",
    instructions: "",
    ingredients: [],
  })

  const updateRecipe = (r: UpdateRecipe) => {
    let newRecipe = { ...recipe }

    if (r.name) newRecipe.name = r.name
    if (r.meal) newRecipe.meal = r.meal
    if (r.description !== undefined && r.description !== null) { newRecipe.description = r.description }
    if (r.instructions) newRecipe.instructions = r.instructions

    setRecipe(newRecipe)
    api.post("recipe/update", { uid: recipe.uid, ...r }).then(rsp => {
      rsp.data.ingredients = recipe.ingredients
      setRecipe(rsp.data)
    })
  }

  const createIngredient = () => {
    api.post("recipe/ingredient/create", {
      recipe_uid: recipe.uid,
      quantity: 1,
      unit: "cup",
      item: "",
      order: recipe.ingredients.length,
    }).then(rsp => {
      let newRecipe = { ...recipe }
      newRecipe.ingredients = [...newRecipe.ingredients, rsp.data]
      setRecipe(newRecipe)
    })
  }

  const updateIngredient = (ndx: number, updateRequest: UpdateIngredient) => {
    let newRecipe = { ...recipe }
    newRecipe.ingredients[ndx] = { ...newRecipe.ingredients[ndx], ...updateRequest }
    setRecipe(newRecipe)
    api.post("recipe/ingredient/update", { uid: recipe.ingredients[ndx].uid, ...updateRequest }).then(rsp => {
      newRecipe.ingredients[ndx] = rsp.data
      setRecipe(newRecipe)
    })
  }

  const reorderIngredient = (ndx: number, delta: number) => {
    if ((ndx === 0 && delta < 0) || (ndx === recipe.ingredients.length - 1 && delta > 0)) {
      return
    }

    let newRecipe = { ...recipe }
    const ingre = recipe.ingredients[ndx]
    if (delta > 0) {
      newRecipe.ingredients.splice(ndx + delta + 1, 0, ingre)
      newRecipe.ingredients.splice(ndx, 1)
    } else if (delta < 0) {
      newRecipe.ingredients.splice(ndx, 1)
      newRecipe.ingredients.splice(ndx + delta, 0, ingre)
    }

    setRecipe(newRecipe)
    api.post("recipe/ingredient/reorder",
      {
        uid: newRecipe.uid,
        ingredient_uids: newRecipe.ingredients.map(r => r.uid)
      }).then(rsp => { setRecipe(rsp.data) })
  }

  const deleteIngredient = (ndx: number) => {
    api.post("recipe/ingredient/delete", { uid: recipe.ingredients[ndx].uid })
    let newRecipe = { ...recipe }
    newRecipe.ingredients = newRecipe.ingredients.filter((_, i) => { return i !== ndx })
    setRecipe(newRecipe)
  }

  useEffect(() => {
    api.post("recipe/read", { uid: recipeUid }).then(rsp => {
      setRecipe(rsp.data)

    }).catch(e => {
      nav("/recipe/list")
    })
  }, [recipeUid, nav])

  const totalInventory = inventory.map(i => i.item)
  totalInventory.push(...AlwaysAvailableIngredients)

  return (<div className="recipe-single-root theme-background">

    <div className="recipe-single-name-meal">
      <FormText label="Name" value={recipe.name} className="recipe-single-name theme-focus"
        inputAttrs={{ autoCapitalize: "words" }} onChange={s => updateRecipe({ name: s })} />
      <FormDropDown label="Meal" value={recipe.meal} className="recipe-single-meal theme-focus"
        options={MealTypes} onChange={s => updateRecipe({ meal: s })} />
    </div>

    <FormTextArea label="Description" value={recipe.description} className="recipe-single-desc theme-focus"
      onChange={s => updateRecipe({ description: s })} />

    <div className="recipe-single-ingredients theme-focus">
      <div className="recipe-single-ingredient-header">
        <span className="recipe-single-ingredient-header-title format-font-small">Ingredients</span>
        <ImageButton alt="more ingredient options" src={areIngredientsEditable ? "/icons/done.svg" : "/icons/edit.svg"}
          className="recipe-single-ingredient-edit" onClick={e => { setEditableIngredients(!areIngredientsEditable) }} />
        <ImageButton alt="more ingredient options" src="/icons/more-horizontal.svg"
          className="recipe-single-ingredient-more" onClick={e => { setDisplayIngredientMore(true) }} />
      </div>
      <div className="recipe-single-ingredient-list">
        {
          recipe.ingredients.map((ingre, i) => {
            return (<RecipeIngredient key={i} model={ingre}
              className={i === ingreNewNdx ? indicatorClass : undefined}
              isDraggable={areIngredientsEditable}
              isInInventory={totalInventory.indexOf(ingre.item) >= 0}
              onChange={s => updateIngredient(i, s)}
              onReorderComplete={d => { setIngreNewNdx(-1); reorderIngredient(i, d) }}
              onReorderMove={d => {
                setIngreNewNdx(d !== 0 ? i + d : -1);
                setIndicatorClass(d < 0 ? "theme-indicator-top" : "theme-indicator-bottom")
              }}
              onDelete={() => deleteIngredient(i)}
            />)
          })
        }
      </div>

      <ImageButton alt="Add ingredient" src="/icons/plus.svg" className="recipe-single-ingredient-add"
        onClick={e => createIngredient()} />
    </div>

    <div className="recipe-single-instruction-list theme-focus">
      <div className="format-font-small">Instructions</div>
      <RecipeInstruction value={recipe.instructions} onChange={s => updateRecipe({ instructions: s })} />
    </div>

    {displayIngredientMore ?
      <RecipeSingleIngredientModal onClose={() => setDisplayIngredientMore(false)}
        recipeIngredients={recipe.ingredients.map(i => i.item)} />
      : null}
  </div>)

}