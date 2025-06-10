import "./css/recipe-search.css"
import { useState } from "react"
import { Link } from "react-router-dom"
import { api } from "../../core/api"
import { FormDropDown } from "../../component/form"
import { ImageButton, InputTextBox } from "../../component/input"
import { MealTypes, RecipeModel } from "./schema"

function nullIfEmpty(s: string): string | null {
  if (s.length === 0) {
    return null
  }

  return s
}

export default function RecipeSearch() {
  const [name, setName] = useState("")
  const [meal, setMeal] = useState("")
  const [description, setDescription] = useState("")
  const [ingredient, setIngredient] = useState("")
  const [showAdvanced, setShowAdvanced] = useState(false)

  const [recipes, setRecipes] = useState<RecipeModel[]>([])

  let SearchMealTypes = new Map(MealTypes)
  SearchMealTypes.set("", "")

  const search = (criteria: any) => {
    let data = {
      name: nullIfEmpty(name),
      meal: nullIfEmpty(meal),
      description: nullIfEmpty(description),
      ingredient: nullIfEmpty(ingredient),
    }
    api.post("recipe/search", { ...data, ...criteria }
    ).then(rsp => { setRecipes(rsp.data) })
  }

  let recipeListEle = null
  if (recipes.length > 0) {
    recipeListEle = recipes.map(r => {
      return (
        <div key={r.uid} className="recipe-search-list-element-root">
          <Link to={`/recipe/edit/${r.uid}`} >
            <div className="recipe-search-list-element-name"><span>{r.name}</span></div>
          </Link>
        </div>
      )
    })
  } else if (name.length > 0 || meal.length > 0 || description.length > 0 || ingredient.length > 0) {
    recipeListEle = (<div className="recipe-search-list-empty">
      <span>Looks like we couldn't find anything that matches your criteria.</span>
    </div>)
  }


  return (<div className="recipe-search-list-root">
    <div className="recipe-search-menu theme-primary-light">
      <InputTextBox value={name} placeholder="Recipe Name" className="recipe-search-menu-name"
        inputAttrs={{ autoCapitalize: "none" }} onChange={e => {
          setName(e);
          if (!showAdvanced) {
            if (e.length > 0) {
              api.post("recipe/search", { name: e }).then(rsp => { setRecipes(rsp.data) })
            }
          } else {
            search({ name: nullIfEmpty(e) })
          }
        }} />
      <ImageButton alt="Show more search options" src="/icons/drop-down.svg"
        className="recipe-search-menu-show-advanced" onClick={e => { setShowAdvanced(!showAdvanced) }} />

      <div style={{ display: showAdvanced ? "flex" : "none" }} className="recipe-search-advanced-menu" >
        <FormDropDown label="Meal" options={SearchMealTypes} value={meal} className="recipe-search-advanced-meal"
          onChange={e => { setMeal(e); search({ meal: nullIfEmpty(e) }) }} />
        <InputTextBox value={description} placeholder="Description" className="recipe-search-advanced-description"
          inputAttrs={{ autoCapitalize: "none" }} onChange={e => { setDescription(e); search({ description: nullIfEmpty(e) }) }} />
        <InputTextBox value={ingredient} placeholder="Ingredient" className="recipe-search-advanced-ingredient"
          inputAttrs={{ autoCapitalize: "none" }} onChange={e => { setIngredient(e); search({ ingredient: nullIfEmpty(e) }) }} />
      </div>
    </div>

    <div className="recipe-search-list-recipes">{recipeListEle}</div>
  </div >)
}

