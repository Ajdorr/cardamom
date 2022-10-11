import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { api } from "../api"
import { TextButton } from "../component/input"
import { Theme } from "../theme"
import { RecipeModel } from "./schema"

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
  const [recipes, setRecipes] = useState<RecipeModel[]>([])
  const [status, setStatus] = useState(
    "It looks like you don't have any available recipes. Try adding more recipes or ingredients into your inventory")
  const nav = useNavigate()

  useEffect(() => {
    api.post("recipe/available").then(rsp => {
      setRecipes(shuffle(rsp.data))
    })
  }, [])

  if (recipes.length === 0) {
    return (<div className="recipe-shuffle-root theme-background">
      <div className="recipe-shuffle-empty theme-focus">
        <span>{status}</span>
      </div>
    </div>)
  }

  const recipe = recipes[recipes.length - 1]

  return (<div className="recipe-shuffle-root recipe-shuffle-root-grid theme-background">

    <div className="recipe-shuffle-name theme-focus"><span>{recipe.name}</span></div>
    <div className="recipe-shuffle-meal theme-focus"><span>{recipe.meal}</span></div>
    <div className="recipe-shuffle-desc theme-focus"><span>{recipe.description}</span></div>
    <div className="recipe-shuffle-ingredient-list theme-focus">
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

    <div className="recipe-shuffle-action-panel theme-focus">

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