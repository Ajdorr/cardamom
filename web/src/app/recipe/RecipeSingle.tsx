import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { api } from "../api"
import { FormDropDown, FormText, FormTextArea } from "../component/form"
import { ImageButton } from "../component/input"
import { RecipeIngredient, RecipeInstruction } from "./component/RecipeComp"
import { CreateRecipeRequest, MealTypes, RecipeModel, Units, UpdateIngredient, UpdateRecipe } from "./schema"

type RecipeSingleProps = {
  isCreate: boolean
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
      .then(rsp => { nav(`/recipe/edit/${rsp.data.uid}`) })
  }

  const updateRecipe = (r: UpdateRecipe) => {
    var newRecipe = { ...recipe }

    if (r.name) newRecipe.name = r.name
    if (r.meal) newRecipe.meal = r.meal
    if (r.description !== undefined && r.description !== null) { newRecipe.description = r.description }
    if (r.instructions) newRecipe.instructions = r.instructions

    setRecipe(newRecipe)
    if (!props.isCreate) {
      api.post("recipe/update", { uid: recipe.uid, ...r }).then(rsp => {
        rsp.data.ingredients = recipe.ingredients
        setRecipe(rsp.data)
      })
    }
  }

  const createIngredient = () => {
    if (props.isCreate) {
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
        optional: false,
        modifier: null,
      }]
      setRecipe(newRecipe)
    } else {
      api.post("recipe/ingredient/create", {
        recipe_uid: recipe.uid,
        quantity: 1,
        unit: "cup",
        item: "",
        order: recipe.ingredients.length,
      }).then(rsp => {
        var newRecipe = { ...recipe }
        newRecipe.ingredients = [...newRecipe.ingredients, rsp.data]
        setRecipe(newRecipe)
      })
    }
  }

  const updateIngredient = (ndx: number, updateRequest: UpdateIngredient) => {
    var newRecipe = { ...recipe }
    newRecipe.ingredients[ndx] = { ...newRecipe.ingredients[ndx], ...updateRequest }
    setRecipe(newRecipe)
    if (!props.isCreate) {
      api.post("recipe/ingredient/update", { uid: recipe.ingredients[ndx].uid, ...updateRequest }).then(rsp => {
        newRecipe.ingredients[ndx] = rsp.data
        setRecipe(newRecipe)
      })
    }
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

    setRecipe(newRecipe)
    if (!props.isCreate) {
      api.post("recipe/ingredient/reorder",
        {
          uid: newRecipe.uid,
          ingredient_uids: newRecipe.ingredients.map(r => r.uid)
        }).then(rsp => { setRecipe(rsp.data) })
    }
  }

  const deleteIngredient = (ndx: number) => {
    if (!props.isCreate) {
      api.post("recipe/ingredient/delete", { uid: recipe.ingredients[ndx].uid })
    }
    var newRecipe = { ...recipe }
    newRecipe.ingredients = newRecipe.ingredients.filter((_, i) => { return i !== ndx })
    setRecipe(newRecipe)
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
        onChange={s => updateRecipe({ name: s })} />
      <FormDropDown label="Meal" value={recipe.meal} className="recipe-single-meal theme-focus"
        options={MealTypes} onChange={s => updateRecipe({ meal: s })} />
    </div>

    <FormTextArea label="Description" value={recipe.description} className="recipe-single-desc theme-focus"
      onChange={s => updateRecipe({ description: s })} />

    <div className="recipe-single-ingredient-list theme-focus">
      <div className="format-font-small">Ingredients</div>
      {
        recipe.ingredients.map((ingre, i) => {
          return (<RecipeIngredient key={i} model={ingre}
            className={i === ingreNewNdx ? indicatorClass : undefined}
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

      <ImageButton alt="Add ingredient" src="/icons/plus.svg" className="recipe-single-ingredient-add"
        onClick={e => createIngredient()} />
    </div>

    <div className="recipe-single-instruction-list theme-focus">
      <div className="format-font-small">Instructions</div>
      <RecipeInstruction value={recipe.instructions} onChange={s => updateRecipe({ instructions: s })} />
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