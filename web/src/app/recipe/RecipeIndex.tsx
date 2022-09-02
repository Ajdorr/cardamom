import "./recipe.css"
import { Component } from "react"
import { Link, Route, Routes } from "react-router-dom"
import RecipeList from "./RecipeList"
import RecipeSingle from "./RecipeSingle"

class RecipeIndex extends Component {

  render() {
    return (<div className="recipe-index-root">
      <div className="recipe-index-create theme-primary">
        <Link to="/recipe/create" >
          <img alt="New Recipe" src="/icons/add-note.svg" className="recipe-index-create-img" />
        </Link>
      </div>

      <Routes>
        <Route path="/" element={<RecipeList />} />
        <Route path="list" element={<RecipeList />} />
        <Route path="create" element={<RecipeSingle isCreate={true} />} />
        <Route path=":recipeUid" element={<RecipeSingle isCreate={false} />} />
      </Routes>
    </div>)
  }

}

export default RecipeIndex