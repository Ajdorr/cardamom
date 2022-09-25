import "./recipe.css"
import { Component } from "react"
import { Route, Routes } from "react-router-dom"
import RecipeList from "./RecipeList"
import RecipeSingle from "./RecipeSingle"
import RecipeShuffle from "./RecipeShuffle"
import RecipeTrashList from "./RecipeTrashList"

class RecipeIndex extends Component {

  render() {
    return (<div className="recipe-index-root">

      <Routes>
        <Route path="/" element={<RecipeList />} />
        <Route path="list" element={<RecipeList />} />
        <Route path="create" element={<RecipeSingle isCreate={true} />} />
        <Route path="available" element={<RecipeShuffle />} />
        <Route path="trash" element={<RecipeTrashList />} />
        <Route path=":recipeUid" element={<RecipeSingle isCreate={false} />} />
      </Routes>
    </div>)
  }

}

export default RecipeIndex