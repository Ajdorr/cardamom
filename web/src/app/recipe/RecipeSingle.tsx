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
}

function RecipeSingle(props: RecipeSingleProps) {

  const { recipeUid } = useParams()
  const nav = useNavigate()
  const [recipe, setRecipe] = useState<RecipeModel>({
    uid: "",
    created_at: "",
    update_at: "",
    user_uid: "",
    name: "",
    description: "",
    meal: "",
    instructions: [],
    ingredients: [],
  })

  const createRecipe = () => {
    api.post("recipe/create", CreateRecipeRequest(recipe))
      .then(rsp => {
        nav(`/recipe/${rsp.data.uid}`)
      })
  }

  const updateRecipe = (r: RecipeModel) => {
    if (props.isCreate) {
      setRecipe(r)
    } else {
      api.post("recipe/update", UpdateRecipeRequest(r))
        .then(rsp => {
          setRecipe(rsp.data)
        }).catch(console.log)
    }
  }

  const setRecipeState = ({ name, meal, description }: UpdateRecipe) => {

    var newRecipe = { ...recipe }

    if (name) newRecipe.name = name
    if (meal) newRecipe.meal = meal
    if (description) newRecipe.description = description

    updateRecipe(newRecipe)
  }

  const addRecipeInstruction = (value: string) => {
    var newRecipe = { ...recipe }
    newRecipe.instructions = [...newRecipe.instructions, {
      uid: "",
      created_at: "",
      update_at: "",
      user_uid: "",
      recipe_uid: "",
      order: newRecipe.instructions.length + 1,
      meal: "",
      text: value,

    }]
    updateRecipe(newRecipe)
  }

  const setRecipeInstruction = (ndx: number, value: string) => {
    var newRecipe = { ...recipe }
    newRecipe.instructions[ndx].text = value
    updateRecipe(newRecipe)
  }

  const removeRecipeInstruction = (ndx: number) => {
    var newRecipe = { ...recipe }
    newRecipe.instructions = newRecipe.instructions.filter((_, i) => { return i !== ndx })
    newRecipe.instructions.forEach((v, i) => v.order = i + 1)
    updateRecipe(newRecipe)
  }

  const addRecipeIngredient = () => {
    var newRecipe = { ...recipe }
    newRecipe.ingredients = [...newRecipe.ingredients, {
      uid: "",
      created_at: "",
      update_at: "",
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
    if (unit) { newRecipe.ingredients[ndx].unit = unit }
    if (item) { newRecipe.ingredients[ndx].item = item }

    updateRecipe(newRecipe)
  }

  const removeRecipeIngredients = (ndx: number) => {
    var newRecipe = { ...recipe }
    newRecipe.ingredients = newRecipe.ingredients.filter((_, i) => { return i !== ndx })
    updateRecipe(newRecipe)
  }

  useEffect(() => {
    if (!props.isCreate) {
      api.post("recipe/read", { uid: recipeUid }).then(rsp => {
        setRecipe(rsp.data)
      }).catch(e => {
        console.log(e)
        nav("/recipe/list")
      })
    }
  }, [recipeUid, nav, props.isCreate])

  return (<div className="reciple-single-root theme-background">
    <div className="reciple-single-name-meal">
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
          return (<RecipeIngredient key={i} value={ingre.item} quantity={ingre.quantity} unit={ingre.unit}
            onQuantityChange={s => setRecipeIngredient(i, { quantity: s })}
            onUnitChange={s => setRecipeIngredient(i, { unit: s })}
            onValueChange={s => setRecipeIngredient(i, { item: s })}
            onDelete={() => removeRecipeIngredients(i)} />)
        })
      }

        <ImageButton alt="Add ingredient" src="/icons/plus.svg" className="recipe-single-ingredient-add"
          onClick={e => addRecipeIngredient()} />
    </div>

    <div className="recipe-single-instruction-list theme-focus">
      <div className="format-font-small">Instructions</div>
      {
        recipe.instructions.map((instr, i) => {
          return (<RecipeInstruction key={i} clearOnChange={false} order={instr.order.toString()}
            value={instr.text} onChange={s => setRecipeInstruction(i, s)} onDelete={() => removeRecipeInstruction(i)} />)
        })
      }
      <RecipeInstruction value="" order="*" clearOnChange={true} placeholder="Add new instruction"
        onChange={s => addRecipeInstruction(s)} />
    </div>

    {props.isCreate ? <ImageButton alt="Save recipe" src="/icons/save.svg" onClick={e => createRecipe()} /> : null}
  </div>)

}

export default RecipeSingle