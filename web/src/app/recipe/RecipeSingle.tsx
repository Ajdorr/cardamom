import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { api } from "../api"
import { FormDropDown, FormText, FormTextArea } from "../component/form"
import { ImageButton } from "../component/input"
import { RecipeIngredient, RecipeInstruction } from "./component/RecipeComp"
import { CreateRecipeRequest, MealTypes, RecipeModel, Units, UpdateRecipeRequest } from "./schema"

type RecipeSingleProps = {
  isCreate: boolean
}

interface UpdateRecipe {
  name?: string
  meal?: string
  description?: string
  instructions?: string
}

function RecipeSingle(props: RecipeSingleProps) {

  const { recipeUid } = useParams()
  const nav = useNavigate()

  const [ingreNewNdx, setIngreNewNdx] = useState(-1)
  const [indicatorClass, setIndicatorClass] = useState("theme-indicator-top")

  const [recipe, setRecipe] = useState<RecipeModel>({
    uid: "",
    created_at: "",
    updated_at: "",
    is_trashed: false,
    user_uid: "",
    name: "",
    description: "",
    meal: "breakfast",
    instructions: "",
    ingredients: [],
  })

  const createRecipe = () => {
    api.post("recipe/create", CreateRecipeRequest(recipe))
      .then(rsp => {
        nav(`/recipe/edit/${rsp.data.uid}`)
      })
  }

  const updateRecipe = (r: RecipeModel) => {
    if (props.isCreate) {
      setRecipe(r)
    } else {
      setRecipe(r)
      api.post("recipe/update", UpdateRecipeRequest(r)).then(rsp => {
        setRecipe(rsp.data)
      })
    }
  }

  const setRecipeState = ({ name, meal, description, instructions }: UpdateRecipe) => {

    var newRecipe = { ...recipe }

    if (name) newRecipe.name = name
    if (meal) newRecipe.meal = meal
    if (description !== undefined && description !== null) { newRecipe.description = description }
    if (instructions) newRecipe.instructions = instructions

    updateRecipe(newRecipe)
  }

  const addRecipeIngredient = () => {
    var newRecipe = { ...recipe }
    newRecipe.ingredients = [...newRecipe.ingredients, {
      uid: "",
      created_at: "",
      updated_at: "",
      user_uid: "",
      recipe_uid: "",
      order: 0,
      meal: "",
      item: "",
      unit: Units[0],
      quantity: 1,
    }]
    updateRecipe(newRecipe)
  }

  const setRecipeIngredient = (
    ndx: number,
    { quantity, unit, item }: { quantity?: string, unit?: string, item?: string }) => {
    var newRecipe = { ...recipe }

    if (quantity) { newRecipe.ingredients[ndx].quantity = quantity }
    if (item) { newRecipe.ingredients[ndx].item = item }
    if (unit !== null && unit !== undefined) {
      newRecipe.ingredients[ndx].unit = unit.length > 0 ? unit : null
    }

    updateRecipe(newRecipe)
  }

  const removeRecipeIngredients = (ndx: number) => {
    var newRecipe = { ...recipe }
    newRecipe.ingredients = newRecipe.ingredients.filter((_, i) => { return i !== ndx })
    updateRecipe(newRecipe)
  }

  const reorderIngredient = (ndx: number, delta: number) => {
    if ((ndx === 0 && delta < 0) || (ndx === recipe.ingredients.length - 1 && delta > 0)) {
      return
    }

    var newRecipe = { ...recipe }
    const ingre = recipe.ingredients[ndx]
    if (delta > 0) {
      newRecipe.ingredients.splice(ndx + delta + 1, 0, ingre)
      newRecipe.ingredients.splice(ndx, 1)
    } else if (delta < 0) {
      newRecipe.ingredients.splice(ndx, 1)
      newRecipe.ingredients.splice(ndx + delta, 0, ingre)
    }

    updateRecipe(newRecipe)
  }

  useEffect(() => {
    if (props.isCreate) {
      setRecipe({
        uid: "",
        created_at: "",
        updated_at: "",
        is_trashed: false,
        user_uid: "",
        name: "",
        description: "",
        meal: "breakfast",
        instructions: "",
        ingredients: [],
      })
    } else {
      api.post("recipe/read", { uid: recipeUid }).then(rsp => {
        setRecipe(rsp.data)
      }).catch(e => {
        nav("/recipe/list")
      })
    }
  }, [recipeUid, nav, props.isCreate])

  return (<div className="recipe-single-root theme-background">

    <div className="recipe-single-name-meal">
      <FormText label="Name" value={recipe.name} className="recipe-single-name theme-focus"
        onChange={s => setRecipeState({ name: s })} />
      <FormDropDown label="Meal" value={recipe.meal} className="recipe-single-meal theme-focus"
        options={MealTypes} onChange={s => setRecipeState({ meal: s })} />
    </div>

    <FormTextArea label="Description" value={recipe.description} className="recipe-single-desc theme-focus"
      onChange={s => setRecipeState({ description: s })} />

    <div className="recipe-single-ingredient-list theme-focus">
      <div className="format-font-small">Ingredients</div>
      {
        recipe.ingredients.map((ingre, i) => {
          return (<RecipeIngredient key={i} value={ingre.item} quantity={ingre.quantity}
            className={i === ingreNewNdx ? indicatorClass : undefined}
            unit={ingre.unit ? ingre.unit : ""}
            onQuantityChange={s => setRecipeIngredient(i, { quantity: s })}
            onUnitChange={s => setRecipeIngredient(i, { unit: s })}
            onValueChange={s => setRecipeIngredient(i, { item: s })}
            onReorder={d => { setIngreNewNdx(-1); reorderIngredient(i, d) }}
            onMove={d => {
              setIngreNewNdx(d !== 0 ? i + d : -1);
              setIndicatorClass(d < 0 ? "theme-indicator-top" : "theme-indicator-bottom")
            }}
            onDelete={() => removeRecipeIngredients(i)}
          />)
        })
      }

      <ImageButton alt="Add ingredient" src="/icons/plus.svg" className="recipe-single-ingredient-add"
        onClick={e => addRecipeIngredient()} />
    </div>

    <div className="recipe-single-instruction-list theme-focus">
      <div className="format-font-small">Instructions</div>
      <RecipeInstruction value={recipe.instructions} onChange={s => setRecipeState({ instructions: s })} />
    </div>

    {props.isCreate ?
      <div className="recipe-single-save theme-focus">
        <ImageButton alt="Save recipe" src="/icons/save.svg"
          onClick={e => createRecipe()} />
      </div>
      : null}
  </div>)

}

export default RecipeSingle