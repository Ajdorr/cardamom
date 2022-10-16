import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { api } from "../api"
import { FormDropDown } from "../component/form"
import { TextButton } from "../component/input"
import { Theme } from "../theme"
import { MealTypes, RecipeModel } from "./schema"

function shuffle(a: Array<any>): Array<any> {
  let i = a.length

  while (i !== 0) {
    const rndx = Math.floor(Math.random() * i)
    i--
    [a[i], a[rndx]] = [a[rndx], a[i]]
  }

  return a
}

function RecipeShuffle() {

  const [meal, setMeal] = useState("dinner")
  const [recipes, setRecipes] = useState<RecipeModel[]>([])
  const [status, setStatus] = useState(
    "It looks like you don't have any available recipes. Try adding more recipes or ingredients into your inventory")
  const nav = useNavigate()

  useEffect(() => {
    api.post("recipe/available", {
      meal: meal ? meal : null
    }).then(rsp => {
      setRecipes(shuffle(rsp.data))
    })
  }, [meal])

  if (recipes.length === 0) {
    return (<div className="recipe-shuffle-root theme-background">
      <div className="recipe-shuffle-empty theme-focus">
        <span>{status}</span>
      </div>
    </div>)
  }

  let MealChoices = new Map(MealTypes)
  MealChoices.set("", "None")
  const recipe = recipes[recipes.length - 1]

  return (<div className="recipe-shuffle-root recipe-shuffle-list theme-background">

    <div className="recipe-shuffle-form theme-focus">
      <FormDropDown label="Meal" options={MealChoices} value={meal} className="recipe-shuffle-form-meal"
        onChange={setMeal} />
    </div>

    <div className="recipe-shuffle-recipe theme-focus">
      <div className="recipe-shuffle-name format-font-header-small theme-primary-light"><span>{recipe.name}</span></div>
      <div className="recipe-shuffle-desc"><span>{recipe.description}</span></div>
      <div className="recipe-shuffle-ingredient-list">
        {
          recipe.ingredients.map(i => {
            return (<div key={i.uid} className="recipe-shuffle-ingredient">
              <span className="recipe-shuffle-ingredient-quantity">{i.quantity}</span>
              <span className="recipe-shuffle-ingredient-unit">{i.unit}</span>
              <span className="recipe-shuffle-ingredient-item">{i.item}</span>
            </div>)
          })
        }
      </div>
      <div className="recipe-shuffle-instructions"><span>{recipe.instructions}</span></div>
    </div>

    <div className="recipe-shuffle-action theme-focus">

      <TextButton label="Yeah!" theme={Theme.Primary} className="recipe-shuffle-action-btn"
        onClick={() => nav(`/recipe/edit/${recipe.uid}`)}
      />

      <TextButton label="Nah." theme={Theme.Surface}
        className="recipe-shuffle-action-btn"
        onClick={() => {
          var rs = [...recipes]
          rs.pop()
          if (rs.length === 0) {
            setStatus("That's all folks!")
          }
          setRecipes(rs)
        }} />
    </div>
  </div>)
}

export default RecipeShuffle