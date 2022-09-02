import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import { api } from "../api"
import { RecipeModel } from "./schema"

function RecipeList() {
  const [recipes, setRecipes] = useState<RecipeModel[]>([])

  useEffect(() => {
    api.post("recipe/list").then(rsp => {
      setRecipes(rsp.data)
    }).catch(e => console.log(e))
  }, [])

  return (<div className="recipe-list-root">
    {recipes.length > 0 ? recipes.map(i => {
      return <RecipeListElement key={i.uid} recipe={i} />
    }) : <div className="recipe-list-empty">No recipes</div>}
  </div>)
}

type RecipeListElementProps = {
  recipe: RecipeModel
}

function RecipeListElement(props: RecipeListElementProps) {

  return (
    <div className="recipe-list-element-root">
      <Link to={`/recipe/${props.recipe.uid}`} >
        <div className="recipe-list-element-name"><span>{props.recipe.name}</span></div>
        <div className="recipe-list-element-meal"><span>{props.recipe.meal}</span></div>
      </Link>
    </div>
  )
}

export default RecipeList